package sys

import (
	"fmt"
	"sync"
	"time"

	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/sys/backdoor"
	"github.com/iansmith/parigot/sys/dep"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Flip this switch to get extra debug information from the nameserver when it is doing
// various lookups.
var nameserverVerbose = false || envVerbose != ""

const MaxService = 127

const parigotNameserverRemoteAddress = "parigot_ns:13330"

var LocalNS *LocalNameServer
var NetNS *NSProxy

// NameServer should probably be renamed. There are two implementations of this interface, one for
// the local (all in one process) case and another for the remote case (across a network).  Some
// of the functions in SysCall are actually delegated down to here.  This is typically done to
// have different behaviors in the local and remote cases.
type NameServer interface {
	//HandleMethod(p *Process, pkgPath, service, method string) (lib.Id, lib.Id)
	Export(key dep.DepKey, pkgPath, service string) lib.Id
	Require(key dep.DepKey, pkgPath, service string) lib.Id
	CloseService(key dep.DepKey, pkgPath, service string) lib.Id
	RunBlock(key dep.DepKey) (bool, lib.Id)
	RunIfReady(key dep.DepKey) []dep.DepKey
	GetService(key dep.DepKey, pkgPath, service string) (lib.Id, lib.Id, string)
	StartFailedInfo() string
	FindMethodByName(key dep.DepKey, serviceId lib.Id, method string) (*callContext, lib.Id, string)
	CallService(dep.DepKey, *callContext) (*syscallmsg.ReturnValueRequest, lib.Id, string)
	GetInfoForCallId(target lib.Id) *callContext
	ExitWhenInFlightEmpty() bool
}

type callContext struct {
	mid          lib.Id                              // the method id this call is going to be made TO
	method       string                              // if the call is remote our LOCAL mid wont mean squat, the remote needs the name
	target       dep.DepKey                          // the process/addr this call is going to be made TO
	cid          lib.Id                              // call id that should be be used by the caller to match results
	sender       dep.DepKey                          // the process/addr this call is going to be made FROM
	sid          lib.Id                              // service that is being called
	respCh       chan *syscallmsg.ReturnValueRequest // this is where to send the return results
	param        *anypb.Any                          // where to put the param data
	pctx         *protosupportmsg.Pctx               // where to put the previous pctx
	timedOut     bool                                // set to true when we waited on a call for a while and didn't get anything
	exitAfterUse bool                                // this is set to true ONLY when the nscore has requested it AND the inflight queue is empty
}

type LocalNameServer struct {
	*NSCore
	notify *sync.Map
}

// NewLocalNameServer creates a new name service implementation for the local or
// "all in one process" case.  We have to use a sync.Map here as the map of
// name of service to notify channel because the map is shared with the
// DelpoyContext.
func NewLocalNameServer(notify *sync.Map) *LocalNameServer {
	return &LocalNameServer{
		notify: notify,
		NSCore: NewNSCore(true),
	}
}

func mapToContent(sm *sync.Map) (int, []string) {
	count := 0
	result := []string{}
	sm.Range(func(key, val any) bool {
		count++
		s := fmt.Sprintf("%s:%T:%v", key, val, val)
		result = append(result, s)
		return true
	})
	return count, result
}

// FindMethodByName is called by the client side when doing a dispatch.  This is where the client
// exchanges a (service id,name) pair for the appropriate call context.  The call context is used
// by the calling client to 1) know where to send the message and 2) how to block waiting on
// the result.  Note that the callContext is created by this function and registered
// with the in flight list.  If there was an error the last two return values will
// indicate the error.  If there was no error, the last two return values will
// nil,"".
func (n *LocalNameServer) FindMethodByName(caller dep.DepKey, serviceId lib.Id, name string) (*callContext, lib.Id, string) {
	// we are NOT holding the lock here
	sData := n.ServiceData(serviceId)
	if sData == nil {
		return nil, lib.NewKernelError(lib.KernelNotFound),
			fmt.Sprintf("could not find service for %s", serviceId.String())
	}
	mid, ok := sData.method.Load(name)
	if !ok {
		return nil, lib.NewKernelError(lib.KernelNotFound),
			fmt.Sprintf("could not find method %s on service %s", name, serviceId.String())
	}
	cc := &callContext{
		method: name,
		sid:    serviceId,
		mid:    mid.(lib.Id),
		respCh: make(chan *syscallmsg.ReturnValueRequest),
		cid:    lib.NewId[*protosupportmsg.CallId](),
		sender: caller,
		target: sData.key,
	}
	n.NSCore.addCallContextMapping(cc.cid, cc)
	return cc, nil, ""
}

// GetService can be called by either a client or a server. If this returns without error, the resulting
// serviceId can be used to be a client of the requested service.
func (n *LocalNameServer) GetService(_ dep.DepKey, pkgPath, service string) (lib.Id, lib.Id, string) {
	return n.NSCore.GetService(pkgPath, service)
}

// GetProcessForCallId is used to match up responses with requests.  It
// walks the in-flight calls and if it finds the target cid it returns
// it and removes it from the in-flight list.
func (n *LocalNameServer) GetInfoForCallId(target lib.Id) *callContext {
	return n.NSCore.getContextForCallId(target)
}

// CloseService is called by a server to inform us (via lib
// and syscall) that there are no more methods to be registered
// for this service. This can fail if the service was already
// closed or the service cannot be found and if so, we return
// the appropriate kernel error to the caller wrapped in a
// lib.Error.
func (n *LocalNameServer) CloseService(key dep.DepKey, pkgPath, service string) lib.Id {
	return n.NSCore.CloseService(key, pkgPath, service)
}

// Exports is used to inform the nameserver that a particular process
// exports the given service.  It returns a kernel error id
// if the service cannot be found or has already been exported
// by another server.
func (n *LocalNameServer) Export(key dep.DepKey, pkgPath, service string) lib.Id {
	return n.NSCore.Export(key, pkgPath, service, nil)
}

// RunBlock is used to wait until the requirements of this process have been
// met. In any case, RunBlock uses RunIfReady to check to see if any
// other processes (possibly including this one) are ready to run.
//
// Tricky: this is called by the goroutine for a server process and it results
// in a block on the notify channel. Later, some *other* goroutine, representing
// a client or a server will end up calling this function and it will unblock
// the previous caller by writing to the notify channel.
func (n *LocalNameServer) RunBlock(key dep.DepKey) (bool, lib.Id) {
	readyList := n.NSCore.RunIfReady(key)
	go n.possiblyUnblock(readyList)

	proc := key.(*DepKeyImpl).proc
	if proc == nil {
		panic("unable to find the process associated with " + key.String())
	}

	myProc := key.(*DepKeyImpl).proc
	myProc.SetReachedRunBlock(true)
	myName := myProc.microservice.GetName()

	chAny, ok := n.notify.Load(myName)
	if !ok {
		panic("unable to find the notification channel for " + myName)
	}
	ch := chAny.(chan bool)
	fromChan := <-ch

	return fromChan, nil
}

// possiblyUnblock is run in a separate goroutine so it can send a message then
// return.  The message is sent to any member of readyList that is not already
// running.
func (n *LocalNameServer) possiblyUnblock(readyList []dep.DepKey) {
	for _, ready := range readyList {
		readyProc := ready.(*DepKeyImpl).proc
		readyName := readyProc.microservice.GetName()
		if readyProc.IsRunning() {
			backdoor.Log(&logmsg.LogRequest{
				Message: fmt.Sprintf("about to skip ready list entry '%s', already marked as running\n", readyName),
				Level:   logmsg.LogLevel_LOG_LEVEL_INFO,
				Stamp:   timestamppb.Now(), // xxx fixme(iansmith) use kernel now
			}, true, false, false, nil)
			continue // nothing to do
		}
		chAny, ok := n.notify.Load(readyName)
		if !ok {
			panic("unable to find the notify channel associated with " + readyName)
		}
		ch := chAny.(chan bool)
		readyProc.SetRunning(true)
		ch <- true
	}
}

// Require is used to inform the nameserver that a particular process
// requires the given service.
func (n *LocalNameServer) Require(key dep.DepKey, pkgPath, service string) lib.Id {
	id := n.NSCore.Require(key, pkgPath, service)
	return id
}

// RunIfReady checks to see if any process is ready to run because all its
// dependencies are satisfied.  If there are ready processes, they are returned.
func (n *LocalNameServer) RunIfReady(key dep.DepKey) []dep.DepKey {
	return n.NSCore.RunIfReady(key)
}

// ExitWhenInFlightEmpty is a switch that says only the calls currently in progress
// should be allowed to complete and the the caller of the last in flight request should
// be told to exit.  This returns true if the in flight queue is empty now.
//
// xxx fixme(iansmith)
// Racey? Could the in flight queue have an entry now but by the time we get
// to the "next" call being processed for a return value we could have an empty
// queue so we never signal a caller to exit.
//
// Tricky: If there are _new_ requests that start after this is called but before the
// "final" request has been processed, then this flag remains on and calls keep getting
// processed until the in flight list is actually empty.  This is usually what you
// want but it could be a problem in a client that has some repetitive task that
// actually "starve" the attempt to ExitWhenInFlightEmpty().
func (n *LocalNameServer) ExitWhenInFlightEmpty() bool {
	return n.NSCore.ExitWhenInFlightEmpty()
}

// SendAbortMessage is used to tell processes that are waiting to run that their
// dependencies could not be fulfilled.  This can only be done when using this
// nameserver.
//
// Because the processes are blocked on their notify channel, we can send a false
// through the run channel to tell them to give up.  We have to use Walk()
// here to walk through all the dependencies and leave the graph unchanged.
func (n *LocalNameServer) sendAbortMessage() {
	panic("sendAbortMessage")
	n.walkDependencyGraph(func(key string, value *dep.EdgeHolder) bool {
		p := value.Key().(*DepKeyImpl).proc
		if p.ReachedRunBlock() && !p.Running() && !p.exited {
			p.runCh <- false
		}
		return true
	})
}

// CallService is used by a process A to signal another process B.  Process A
// has created the callContext and pushes it through the callCh that B is blocked on.
// Included in the callContext is another channel that B can use to send the result,
// in the form a returnValueRequest to A.
//
// We use the returnValueRequest so this path is the same as it would be in the case
// of a remote call.  In this the local case, we *could* just pass the result back
// from B to A.
func (l *LocalNameServer) CallService(key dep.DepKey, ctx *callContext) (*syscallmsg.ReturnValueRequest, lib.Id, string) {
	proc := key.(*DepKeyImpl).proc
	proc.callCh <- ctx
	result := <-ctx.respCh
	return result, nil, ""
}

// BlockUntilCall implements the stopping of a program until a method is
// called.  Because this all implemented locally in this case, we just
// block on our callCh and wait for another process to use CallService to
// signal us that they need one of our methods.
func (l *LocalNameServer) BlockUntilCall(key dep.DepKey, canTimeout bool) *callContext {
	if canTimeout { //simple case
		v := <-key.(*DepKeyImpl).proc.callCh
		return v
	}
	// complex case
	select {
	case v := <-key.(*DepKeyImpl).proc.callCh:
		v.timedOut = false
		return v
	case <-time.After(1 * time.Second):
		v := &callContext{}
		v.timedOut = true
		return v
	}
}

func nameserverPrint(methodName string, format string, arg ...interface{}) {
	if nameserverVerbose {
		part1 := fmt.Sprintf("NAMESERVER:%s", methodName)
		part2 := fmt.Sprintf(format, arg...)
		req := logmsg.LogRequest{
			Level:   logmsg.LogLevel_LOG_LEVEL_DEBUG,
			Stamp:   timestamppb.Now(), //xxx fix me should be using the kernel for this
			Message: part1 + part2,
		}
		backdoor.Log(&req, true, false, false, nil)
	}
}
