package sys

import (
	"fmt"
	"sync"

	"github.com/iansmith/parigot/lib"
)

// Flip this switch to get extra debug information from the nameserver when it is doing
// various lookups.
var nameserverVerbose = true

const MaxService = 127

type callContext struct {
	mid    lib.Id   // the method id this call is going to be made TO
	target *Process // the process this call is going to be made TO
	cid    lib.Id   // call id that should be be used by the caller to match results
	sender *Process // the process this call is going to be made FROM
}

type packageData struct {
	service map[string]*serviceData
}

func newPackageData() *packageData {
	return &packageData{
		service: make(map[string]*serviceData),
	}
}

type serviceData struct {
	serviceId         lib.Id
	closed            bool
	exported          bool
	method            map[string]lib.Id
	methodIdToProcess map[string] /*really method id*/ *Process
}

func newServiceData() *serviceData {
	return &serviceData{
		serviceId:         nil,
		method:            make(map[string]lib.Id),
		methodIdToProcess: make(map[string]*Process),
	}
}

type NameServer struct {
	packageRegistry        map[string]*packageData
	serviceCounter         int
	serviceIdToServiceData map[string] /*really service id*/ *serviceData // accelerator only, we could walk to find this
	inFlight               []*callContext
	lock                   *sync.RWMutex
	dependencyGraph        map[string]*edgeHolder
	alreadyExported        []string
	runCh                  chan *Process
}

func NewNameServer() *NameServer {
	return &NameServer{
		lock:                   new(sync.RWMutex),
		packageRegistry:        make(map[string]*packageData),
		serviceIdToServiceData: make(map[string]*serviceData),
		inFlight:               []*callContext{},
		serviceCounter:         7, //first 8 are reserved
		dependencyGraph:        make(map[string]*edgeHolder),
		alreadyExported:        []string{},
		runCh:                  make(chan *Process),
	}
}

// RegisterClientService connects a packagePath.service with a particular service id.
// This is used by client side code and is ususally called via the init() method in
// generated code.
func (n *NameServer) RegisterClientService(packagePath string, service string) (lib.Id, lib.Id) {
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.serviceCounter >= MaxService {
		return nil, lib.NewKernelError(lib.KernelNamespaceExhausted)
	}

	pData, ok := n.packageRegistry[packagePath]
	if !ok {
		pData = newPackageData()
		n.packageRegistry[packagePath] = pData
	}
	sData, ok := pData.service[service]
	if !ok {
		sData = newServiceData()
		n.serviceCounter++
		sData.serviceId = lib.ServiceIdFromUint64(0, uint64(n.serviceCounter))
		n.serviceIdToServiceData[sData.serviceId.String()] = sData
		pData.service[service] = sData
	}
	return sData.serviceId, nil
}

// FindMethodByName is called by the client side when doing a dispatch.  This is where the client
// exchanges a service.id,name pair for the appropriate call context.  The call context is used
// by the calling client to 1) know where to send the message and 2) how to block waiting on
// the result.  The return result here is the property of the nameserver, don't mess with it,
// just read it.
func (n *NameServer) FindMethodByName(serviceId lib.Id, name string, caller *Process) *callContext {
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
	target, ok := sData.methodIdToProcess[mid.String()]
	if !ok {
		return nil
	}
	cc := &callContext{
		mid:    mid,
		target: target,
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
func (n *NameServer) HandleMethod(pkgPath, service, method string, proc *Process) (lib.Id, lib.Id) {
	n.lock.Lock()
	defer n.lock.Unlock()

	nameserverPrint("HANDLEMETHOD", "adding method in the nameserver in process %s",
		proc.String())
	pData, ok := n.packageRegistry[pkgPath]
	if !ok {
		pData = newPackageData()
		n.packageRegistry[pkgPath] = pData
	}
	sData, ok := pData.service[service]
	if !ok {
		sData = newServiceData()
		n.serviceCounter++
		sData.serviceId = lib.ServiceIdFromUint64(0, uint64(n.serviceCounter))
		pData.service[service] = sData
		n.serviceIdToServiceData[sData.serviceId.String()] = sData
	}
	result := lib.NewMethodId()
	nameserverPrint("HANDLEMETHOD", "assigning %s to the method %s in service %s",
		result, method, service)
	sData.method[method] = result
	sData.methodIdToProcess[result.String()] = proc

	// xxx fix me, should be able to realize that a method does not exist and reject the attempt to
	// handle it

	//xxx fix me, there should be a limit on the number of methods per service
	return result, nil
}

// GetService can be called by either a client or a server. If this returns without error, the resulting
// serviceId can be used to be a client of the requested service.
func (n *NameServer) GetService(pkgPath, service string) (lib.Id, lib.Id) {
	n.lock.RLock()
	defer n.lock.RUnlock()

	pData, ok := n.packageRegistry[pkgPath]
	if !ok {
		return nil, lib.NewKernelError(lib.KernelNotFound)
	}
	sData, ok := pData.service[service]
	if !ok {
		return nil, lib.NewKernelError(lib.KernelNotFound)
	}
	return sData.serviceId, nil
}

// GetProcessForCallId is used to match up responses with requests.  It
// walks the in-flight calls and if it finds the target cid it returns
// it and removes it from the in-flight list.
func (n *NameServer) GetProcessForCallId(target lib.Id) *Process {
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
func (n *NameServer) CloseService(pkgPath string, service string) lib.Error {
	n.lock.Lock()
	defer n.lock.Unlock()

	nameserverPrint("CLOSESERVICE", "closing %s.%s",
		pkgPath, service)
	sData, err := n.validatePkgAndService(pkgPath, service)
	if err != nil {
		return err
	}
	if sData.closed {
		return lib.NewPerrorFromId("already closed",
			lib.NewKernelError(lib.KernelServiceAlreadyClosedOrExported))

	}
	sData.closed = true
	return nil
}

// validatePkgAndService makes sure the package and service are valid.
// It does not lock, so callers must hold the lock before they call this
// function.
func (n *NameServer) validatePkgAndService(pkgPath, service string) (*serviceData, lib.Error) {

	pData, ok := n.packageRegistry[pkgPath]
	if !ok {
		return nil, lib.NewPerrorFromId("no such package",
			lib.NewKernelError(lib.KernelNotFound))
	}
	sData, ok := pData.service[service]
	if !ok {
		return nil, lib.NewPerrorFromId("no such package",
			lib.NewKernelError(lib.KernelNotFound))
	}
	return sData, nil
}

// Exports is used to inform the nameserver that a particular process
// exports the given service.  It returns a kernel error id inside the
// lib.Error if the service cannot be found or has already been exported
// by another server.
func (n *NameServer) Export(proc *Process, pkgPath, service string) lib.Error {
	n.lock.Lock()
	defer n.lock.Unlock()

	nameserverPrint("EXPORT", "process %s exports %s.%s",
		proc.String(), pkgPath, service)
	sData, err := n.validatePkgAndService(pkgPath, service)
	if err != nil {
		return err
	}
	if !sData.closed {
		// xxxfix me, should this be an error?
	}
	if sData.exported {
		return lib.NewPerrorFromId("already exported",
			lib.NewKernelError(lib.KernelServiceAlreadyClosedOrExported))
	}
	sData.exported = true
	node, ok := n.dependencyGraph[proc.String()]
	if !ok {
		node = &edgeHolder{
			proc:    proc,
			export:  []string{},
			require: []string{},
		}
		n.dependencyGraph[proc.String()] = node
	}
	node.export = append(node.export, fmt.Sprintf("%s.%s", pkgPath, service))
	return nil
}

// Require is used to inform the nameserver that a particular process
// requires the given service.
func (n *NameServer) Require(proc *Process, pkgPath, service string) lib.Error {
	n.lock.Lock()
	defer n.lock.Unlock()

	alreadyExported := false

	nameserverPrint("REQUIRE", "process %s requires %s.%s",
		proc.String(), pkgPath, service)

	name := fmt.Sprintf("%s.%s", pkgPath, service)
	for _, s := range n.alreadyExported {
		if s == name {
			nameserverPrint("REQUIRE", "process %s required %s.%s but it is already exported",
				proc.String(), pkgPath, service)
			alreadyExported = true
		}
	}

	pData, ok := n.packageRegistry[pkgPath]
	if !ok {
		pData = newPackageData()
		n.packageRegistry[pkgPath] = pData
	}
	sData, ok := pData.service[service]
	if !ok {
		sData = newServiceData()
		n.serviceCounter++
		sData.serviceId = lib.ServiceIdFromUint64(0, uint64(n.serviceCounter))
		pData.service[service] = sData
		n.serviceIdToServiceData[sData.serviceId.String()] = sData
	}
	node, ok := n.dependencyGraph[proc.String()]
	if !ok {
		node = &edgeHolder{
			proc:    proc,
			export:  []string{},
			require: []string{},
		}
		n.dependencyGraph[proc.String()] = node
	}
	if !alreadyExported {
		for _, r := range node.require {
			if r == name {
				return lib.NewPerrorFromId("already required",
					lib.NewKernelError(lib.KernelServiceAlreadyRequired))
			}
		}
		node.require = append(node.require, name)
	}
	return nil
}

func nameserverPrint(methodName string, format string, arg ...interface{}) {
	if nameserverVerbose {
		part1 := fmt.Sprintf("NAMESERVER:%s", methodName)
		part2 := fmt.Sprintf(format, arg...)
		print(part1, part2, "\n")
	}
}

// RunReader's job is to listen on the channel for messages to the nameserver that some
// user program is ready to run.  When that happens we send his process information through
// the machinery to figure out dependencies and so on.
func (n *NameServer) RunReader() {
	for {
		nameserverPrint("RUNREADER ", "about to read from the channel")
		p := <-n.runCh
		nameserverPrint("RUNREADER ", "calling nameserver.Run on %s", p)
		n.Run(p)
	}
}

func (n *NameServer) Run(proc *Process) {
	n.lock.Lock()
	defer n.lock.Unlock()

	nameserverPrint("RUN ", "-----> called for process=%s,waiter=%v", proc, proc.waiter)
	// the only time things can change is when a new process calls run
	// so we only to make sure here that we find all the eligible processes to run

	node := n.dependencyGraph[proc.String()]
	candidateList := []*edgeHolder{node}

	nameserverPrint("RUN ", "node %s (%d req,%d exp) and dep-graph has %d total entries", proc, len(node.require), len(node.export), len(n.dependencyGraph))
	for len(candidateList) > 0 {
		candidate := candidateList[0]
		// remove candidate from list
		if len(candidateList) == 1 {
			candidateList = nil
		} else {
			candidateList = candidateList[1:]
		}
		// are we ready to run?
		if candidate.isReady() {
			nameserverPrint("RUN ", "candidate %s is ready to run", candidate.proc)
			delete(n.dependencyGraph, candidate.proc.String())
			// we are ready, so lets process his exports through the list of waiting processes
			for _, other := range n.dependencyGraph {
				changed := other.removeRequire(candidate.export)
				if changed {
					candidateList = append(candidateList, other)
					nameserverPrint("RUN ", "candidate list changed")
				}
			}
			n.alreadyExported = append(n.alreadyExported, candidate.export...)
			nameserverPrint("RUN", "already exported updated: %+v", n.alreadyExported)

			nameserverPrint("RUN ", "running %s", candidate.proc)
			candidate.proc.Run()
		} else {
			nameserverPrint("RUN ", "%s is not ready to run, number of candidates left is %d", candidate.proc, len(candidateList))
		}
	}
	nameserverPrint("RUN ", "we are returning to the RunReader loop")
}

func (n *NameServer) WaitingToRun() int {
	n.lock.RLock()
	defer n.lock.RUnlock()

	return len(n.dependencyGraph)
}

// SendAbortMessage is used to tell processes that are waiting to run that their
// dependencies could not be fulfilled.
func (n *NameServer) SendAbortMessage() {
	for _, v := range n.dependencyGraph {
		if v.proc.reachedRun {
			if !v.proc.exited {
				v.proc.runCh <- false
			}
		}
	}
}
