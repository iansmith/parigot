package sys

import (
	"fmt"
	"sync"

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

// These are the two nameservers.  They share a runNotifyChannel and created
// by a call to InitNameServer()
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
	RunNotify(key dep.DepKey)
	RunBlock(key dep.DepKey) (bool, lib.Id)
	RunIfReady(key dep.DepKey)
	GetService(key dep.DepKey, pkgPath, service string) (lib.Id, lib.KernelErrorCode)
	StartFailedInfo() string
	FindMethodByName(key dep.DepKey, serviceId lib.Id, method string) *callContext
	CallService(dep.DepKey, *callContext) *syscallmsg.ReturnValueRequest
	GetInfoForCallId(target lib.Id) *callContext
}

type callContext struct {
	mid    lib.Id                              // the method id this call is going to be made TO
	method string                              // if the call is remote our LOCAL mid wont mean squat, the remote needs the name
	target dep.DepKey                          // the process/addr this call is going to be made TO
	cid    lib.Id                              // call id that should be be used by the caller to match results
	sender dep.DepKey                          // the process/addr this call is going to be made FROM
	sid    lib.Id                              // service that is being called
	respCh chan *syscallmsg.ReturnValueRequest // this is where to send the return results
	param  *anypb.Any                          // where to put the param data
	pctx   *protosupportmsg.Pctx               // where to put the previous pctx
}

type LocalNameServer struct {
	*NSCore
	lock        *sync.RWMutex
	runNotifyCh chan *KeyNSPair
}

// NewLocalNameServer creates a new name service implementation for the local or
// "all in one process" case.
func NewLocalNameServer(runNotifyChannel chan *KeyNSPair) *LocalNameServer {
	return &LocalNameServer{
		lock:        new(sync.RWMutex),
		runNotifyCh: runNotifyChannel,
		NSCore:      NewNSCore(true),
	}
}

// FindMethodByName is called by the client side when doing a dispatch.  This is where the client
// exchanges a (service id,name) pair for the appropriate call context.  The call context is used
// by the calling client to 1) know where to send the message and 2) how to block waiting on
// the result.  Note that the callContext is created by this function and registered
// with the in flight list.
func (n *LocalNameServer) FindMethodByName(caller dep.DepKey, serviceId lib.Id, name string) *callContext {
	// we are NOT holding the lock here
	sData := n.ServiceData(serviceId)
	if sData == nil {
		return nil
	}
	mid, ok := sData.method[name]
	if !ok {
		return nil
	}
	cc := &callContext{
		method: name,
		sid:    serviceId,
		mid:    mid,
		respCh: make(chan *syscallmsg.ReturnValueRequest),
		cid:    lib.NewId[*protosupportmsg.CallId](),
		sender: caller,
		target: sData.key,
	}
	n.NSCore.addCallContextMapping(cc.cid, cc)
	return cc
}

// GetService can be called by either a client or a server. If this returns without error, the resulting
// serviceId can be used to be a client of the requested service.
func (n *LocalNameServer) GetService(_ dep.DepKey, pkgPath, service string) (lib.Id, lib.KernelErrorCode) {
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

// Require is used to inform the nameserver that a particular process
// requires the given service.
func (n *LocalNameServer) Require(key dep.DepKey, pkgPath, service string) lib.Id {
	return n.NSCore.Require(key, pkgPath, service)
}

// RunIfReady blocks until it receives a callback that the given key, representing
// a process, has all its requirements satisfied.
func (n *LocalNameServer) RunIfReady(key dep.DepKey) {
	n.NSCore.RunIfReady(key, func(key dep.DepKey) {
		nscorePrint("RunIfReadyCallback ", "notifying run reader %s is ready", key.String())
		key.(*DepKeyImpl).proc.Run()
	})
}

// SendAbortMessage is used to tell processes that are waiting to run that their
// dependencies could not be fulfilled.  This can only be done when using this
// nameserver.
//
// Because the processes are blocked on their run channel, we can send a false
// through the run channel to tell them to give up.  We have to use Walk()
// here to walk through all the dependencies and leave the graph unchanged.
func (n *LocalNameServer) sendAbortMessage() {
	n.dependencyGraph.Walk(func(key string, value *dep.EdgeHolder) bool {
		p := value.Key().(*DepKeyImpl).proc
		if p.reachedRun {
			if !p.exited {
				p.runCh <- false
			}
		}
		return true
	})
}

// RunNotify is called by a proc that is local to shove itself and this nameserver
// to the run reader.  The run reader is a separate goroutine that listens for
// these notifications and then looks for any ready to run processes.
func (l *LocalNameServer) RunNotify(key dep.DepKey) {
	l.runNotifyCh <- NewKeyNSPair(key, l)
}

// RunBlock is called by a proc that is local and this blocks until the nameserver
// signals to us that our dependencies have been met.  Note that call is running
// on a different goroutine (the goroutine of the service who is blocked) than
// the RunIfReady() call than unblocks it.
func (l *LocalNameServer) RunBlock(key dep.DepKey) (bool, lib.Id) {
	nameserverPrint("RunBlock ", "localnameserver about to block on runCh for %s\n", key.String())
	b := <-key.(*DepKeyImpl).proc.runCh
	nameserverPrint("RunBlock", "localnameserver, block done for %s\n", key.String())
	return b, nil
}

// CallService is used by a process A to signal another process B.  Process A
// has created the callContext and pushes it through the callCh that B is blocked on.
// Included in the callContext is another channel that B can use to send the result,
// in the form a returnValueRequest to A.
//
// We use the returnValueRequest so this path is the same as it would be in the case
// of a remote call.  In this the local case, we *could* just pass the result back
// from B to A.
func (l *LocalNameServer) CallService(key dep.DepKey, ctx *callContext) *syscallmsg.ReturnValueRequest {
	proc := key.(*DepKeyImpl).proc
	proc.callCh <- ctx
	result := <-ctx.respCh
	return result
}

// BlockUntilCall implements the stopping of a program until a method is
// called.  Because this all implemented locally in this case, we just
// block on our callCh and wait for another process to use CallService to
// signal us that they need one of our methods.
func (l *LocalNameServer) BlockUntilCall(key dep.DepKey) *callContext {
	v := <-key.(*DepKeyImpl).proc.callCh
	return v
}

// InitNameServers initializes the two nameservers with a shared channel that
// is used to implement RunNotify for the local case.
func InitNameServer(runNotifyChannel chan *KeyNSPair, local, remote bool) {
	if local {
		LocalNS = NewLocalNameServer(runNotifyChannel)
	}
	if remote {
		NetNS = NewNSProxy(parigotNameserverRemoteAddress)
	}
}

// StartFailedInfo returns the pair of strings that results from calling the
// StartFailedInfo on each nameserver.  We have to do this on both nameservers
// because it is possible only one of them has a problem.
func StartFailedInfo() (string, string) {
	local := ""
	remote := ""
	if LocalNS != nil {
		local = LocalNS.StartFailedInfo()
	}
	if NetNS != nil {
		remote = NetNS.StartFailedInfo()
	}
	return local, remote
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
