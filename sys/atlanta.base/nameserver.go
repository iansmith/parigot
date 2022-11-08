package sys

import (
	"fmt"
	"sync"

	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/dep"
)

// Flip this switch to get extra debug information from the nameserver when it is doing
// various lookups.
var nameserverVerbose = false

const MaxService = 127

// These are the two nameservers.  They share a runNotifyChannel and created
// by a call to InitNameServer()
var LocalNS *LocalNameServer
var NetNS *NetNameServer

type NameServer interface {
	HandleMethod(p *Process, pkgPath, service, method string) (lib.Id, lib.Id)
	Export(key dep.DepKey, pkgPath, service string) lib.Id
	CloseService(pkgPath, service string) lib.Id
	RunNotify(key dep.DepKey)
	RunBlock(key dep.DepKey) bool
	RunIfReady(key dep.DepKey)
	StartFailedInfo() string
}

type callContext struct {
	mid    lib.Id   // the method id this call is going to be made TO
	target *Process // the process this call is going to be made TO
	cid    lib.Id   // call id that should be be used by the caller to match results
	sender *Process // the process this call is going to be made FROM
}
type LocalNameServer struct {
	*NSCore
	inFlight    []*callContext
	lock        *sync.RWMutex
	runNotifyCh chan *KeyNSPair
}

func NewLocalNameServer(runNotifyChannel chan *KeyNSPair) *LocalNameServer {
	return &LocalNameServer{
		lock:        new(sync.RWMutex),
		inFlight:    []*callContext{},
		runNotifyCh: runNotifyChannel,
		NSCore:      NewNSCore(),
	}
}

// FindMethodByName is called by the client side when doing a dispatch.  This is where the client
// exchanges a service.id,name pair for the appropriate call context.  The call context is used
// by the calling client to 1) know where to send the message and 2) how to block waiting on
// the result.  The return result here is the property of the nameserver, don't mess with it,
// just read it.
func (n *LocalNameServer) FindMethodByName(caller *Process, serviceId lib.Id, name string) *callContext {
	n.lock.Lock()
	defer n.lock.Unlock()

	sData, ok := n.serviceIdToServiceData[serviceId.String()]
	if !ok {
		return nil
	}
	mid, ok := sData.method[name]
	if !ok {
		return nil
	}
	target, ok := sData.methodIdToImpl[mid.String()]
	if !ok {
		return nil
	}
	cc := &callContext{
		mid:    mid,
		target: target.(*depKeyImpl).proc,
		cid:    lib.NewCallId(),
		sender: caller,
	}
	nameserverPrint("FINDMETHODBYNAME", "adding in flight rpc call %s and %s",
		cc.cid.Short(), cc.sender.String())
	n.inFlight = append(n.inFlight, cc)
	return cc
}

// HandleMethod is called by the server side to indicate that it will handle a particular
// method call on a particular service.
func (n *LocalNameServer) HandleMethod(proc *Process, pkgPath, service, method string) (lib.Id, lib.Id) {
	n.lock.Lock()
	defer n.lock.Unlock()

	nameserverPrint("HANDLEMETHOD", "adding method in the nameserver in process %s",
		proc.String())
	return n.NSCore.HandleMethod(newDepKeyFromProcess(proc), pkgPath, service, method)

}

// GetService can be called by either a client or a server. If this returns without error, the resulting
// serviceId can be used to be a client of the requested service.
func (n *LocalNameServer) GetService(pkgPath, service string) (lib.Id, lib.Id) {
	n.lock.RLock()
	defer n.lock.RUnlock()

	return n.NSCore.GetService(pkgPath, service)
}

// GetProcessForCallId is used to match up responses with requests.  It
// walks the in-flight calls and if it finds the target cid it returns
// it and removes it from the in-flight list.
func (n *LocalNameServer) GetProcessForCallId(target lib.Id) *Process {
	n.lock.Lock()
	defer n.lock.Unlock()

	for i, cctx := range n.inFlight {
		nameserverPrint("GETPROCESSFORCALLID", "checking in-flight rpc calls, cctx #%d, with target %s versus %s", i, target.Short(),
			cctx.cid.Short())
		if cctx.cid.Equal(target) {
			n.inFlight[i] = n.inFlight[len(n.inFlight)-1]
			n.inFlight = n.inFlight[:len(n.inFlight)-1]
			// xxxfix me should we be checking the method id as well?
			return cctx.sender
		}
	}
	return nil
}

// CloseService is called by a server to inform us (via lib
// and syscall) that there are no more methods to be registered
// for this service. This can fail if the service was already
// closed or the service cannot be found and if so, we return
// the appropriate kernel error to the caller wrapped in a
// lib.Error.
func (n *LocalNameServer) CloseService(pkgPath string, service string) lib.Id {
	n.lock.Lock()
	defer n.lock.Unlock()

	return n.NSCore.CloseService(pkgPath, service)
}

// Exports is used to inform the nameserver that a particular process
// exports the given service.  It returns a kernel error id
// if the service cannot be found or has already been exported
// by another server.
func (n *LocalNameServer) Export(key dep.DepKey, pkgPath, service string) lib.Id {
	n.lock.Lock()
	defer n.lock.Unlock()

	return n.NSCore.Export(key, pkgPath, service)
}

// Require is used to inform the nameserver that a particular process
// requires the given service.
func (n *LocalNameServer) Require(key dep.DepKey, pkgPath, service string) lib.Id {
	n.lock.Lock()
	defer n.lock.Unlock()

	return n.NSCore.Require(key, pkgPath, service)
}

func (n *LocalNameServer) RunIfReady(key dep.DepKey) {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.NSCore.RunIfReady(key, func(key dep.DepKey) {
		key.(*depKeyImpl).proc.Run()
	})
}

func (n *LocalNameServer) WaitingToRun() int {
	n.lock.RLock()
	defer n.lock.RUnlock()

	return n.NSCore.WaitingToRun()
}

func (l *LocalNameServer) StartFailedInfo() string {
	if l.NSCore.WaitingToRun() > 0 {
		l.sendAbortMessage()
	}
	return l.NSCore.StartFailedInfo()
}

// SendAbortMessage is used to tell processes that are waiting to run that their
// dependencies could not be fulfilled.  This can only be done when using this
// nameserver.
func (n *LocalNameServer) sendAbortMessage() {
	for _, v := range n.dependencyGraph.AllEdge() {
		p := v.Key().(*depKeyImpl).proc
		if p.reachedRun {
			if !p.exited {
				p.runCh <- false
			}
		}
	}
}

// This is called by a proc that is local to shove itself and this nameserver
// to the run reader.
func (l *LocalNameServer) RunNotify(key dep.DepKey) {
	l.runNotifyCh <- NewKeyNSPair(key, l)
}

// This is called by a proc that is local and this blocks until the nameserver
// signals to us.
func (l *LocalNameServer) RunBlock(key dep.DepKey) bool {
	b := <-key.(*depKeyImpl).proc.runCh
	return b
}

// InitNameServers initializes the two nameservers with a shared channel that
// is used to implement RunNotify.
func InitNameServer(runNotifyChannel chan *KeyNSPair, local, remote bool) {
	if local {
		LocalNS = NewLocalNameServer(runNotifyChannel)
	}
	if remote {
		NetNS = NewNetNameserver(LocalNS, "parigot_ns:13330")
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
		print(part1, part2, "\n")
	}
}
