package sys

import (
	"bytes"
	"fmt"
	"log"
	"sync"

	"github.com/iansmith/parigot/lib"
)

// Flip this switch to get extra debug information from the nameserver when it is doing
// various lookups.
var nameserverVerbose = false

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
}

func NewNameServer() *NameServer {
	return &NameServer{
		lock:                   new(sync.RWMutex),
		packageRegistry:        make(map[string]*packageData),
		serviceIdToServiceData: make(map[string]*serviceData),
		inFlight:               []*callContext{},
		serviceCounter:         7, //first 8 are reserved
		dependencyGraph:        make(map[string]*edgeHolder),
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

	nameserverPrint("REQUIRE", "process %s exports %s.%s",
		proc.String(), pkgPath, service)

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
	name := fmt.Sprintf("%s.%s", pkgPath, service)
	for _, r := range node.require {
		if r == name {
			return lib.NewPerrorFromId("already required",
				lib.NewKernelError(lib.KernelServiceAlreadyRequired))
		}
	}
	node.require = append(node.require, name)
	return nil
}

func nameserverPrint(methodName string, format string, arg ...interface{}) {
	if nameserverVerbose {
		part1 := fmt.Sprintf("NAMESERVER:%s", methodName)
		part2 := fmt.Sprintf(format, arg...)
		print(part1, part2, "\n")
	}
}

func (n *NameServer) Start(proc *Process) {
	n.lock.Lock()
	defer n.lock.Unlock()
	// the only time things can change is when a new process calls start
	// so we only to make sure here that we find all the eligible processes to start

	node := n.dependencyGraph[proc.String()]
	candidateList := []*edgeHolder{node}

	for len(candidateList) > 0 {
		candidate := candidateList[0]
		// are we ready to run?
		if candidate.isReady() {
			// remove candidate from list
			if len(candidateList) == 1 {
				candidateList = nil
			} else {
				candidateList = candidateList[1:]
			}
			delete(n.dependencyGraph, candidate.proc.String())
			nameserverPrint("START", "process %s is ready to run", proc.String())
			// we are ready, so lets process his exports through the list of waiting processes
			for _, other := range n.dependencyGraph {
				changed := other.removeRequire(candidate.export)
				if changed {
					candidateList = append(candidateList, other)
				}
			}
			candidate.proc.Run("")
		}
	}
}

func (n *NameServer) WaitingToRun() int {
	n.lock.RLock()
	defer n.lock.RUnlock()

	return len(n.dependencyGraph)
}

type depPair struct {
	proc       *Process
	service    string
	linkExport string
}

func (d *depPair) String() string {
	if d.linkExport == "" {
		return fmt.Sprintf("require %s in %s;", d.service, d.proc.String())
	}
	return fmt.Sprintf("require %s (via export of %s) in %s;", d.service, d.linkExport, d.proc.String())
}

// DependencyLoop does a DFS looking for a cycle.
func (n *NameServer) dependencyLoop(cand *depPair, seen []*depPair) string {
	log.Printf("---------------- dependency loop ---------\n")
	log.Printf("candidate: %s\n", cand.String())
	for _, s := range seen {
		log.Printf("\tseen %s\n", s.String())
	}
	log.Printf("---------------- END ---------\n")
	found := false
	//check to see if we found it
	for _, pair := range seen {
		if cand.service == pair.service {
			found = true
			break
		}
	}
	if found {
		log.Printf("got it! %s\n", cand.String())
		var buf bytes.Buffer
		for _, pair := range seen {
			buf.WriteString(pair.String())
		}
		// add the loop end
		buf.WriteString(cand.String())
		log.Printf("got it part 2: %s", buf.String())
		return buf.String()
	}
	//we didn't find it so we are now in seen
	seen = append(seen, cand)
	child := n.findCandidateProcessesByExport(cand)
	for _, c := range child {
		loop := n.dependencyLoop(c, seen)
		if loop != "" {
			return loop
		}
	}
	log.Printf("RET EMPTY %s", cand.String())
	return ""
}

// findCandidateProcessesByExport returns the set of processes that export a particular
// service.
func (n *NameServer) findCandidateProcessesByExport(pair *depPair) []*depPair {
	candidateList := []*depPair{}
	for _, node := range n.dependencyGraph {
		for _, exp := range node.export {
			if exp == pair.service {
				for _, r := range node.require {
					candidateList = append(candidateList,
						&depPair{proc: node.proc,
							linkExport: exp,
							service:    r})
				}
				break
			}
		}
	}
	return candidateList
}

func (n *NameServer) getLoopContent() string {
	if len(n.dependencyGraph) == 0 {
		panic("should not be sending scanning for loop when every process is running!")
	}
	candidateList := []*depPair{}
	// we want to try all combos
	for _, v := range n.dependencyGraph {
		for _, req := range v.require {
			candidateList = append(candidateList, &depPair{v.proc, req, ""})
		}
	}
	for _, candidate := range candidateList {
		loop := n.dependencyLoop(candidate, []*depPair{})
		log.Printf("deploop returned '%s'", loop)
		if loop != "" {
			return loop
		}
	}
	return ""
}

// getDeadNodeContent returns a list of the nodes that cannot possibly be fulfilled.
func (n *NameServer) getDeadNodeContent() string {
	candidateList := []*depPair{}
	// we want to try all combos
	for _, v := range n.dependencyGraph {
		for _, req := range v.require {
			candidateList = append(candidateList, &depPair{v.proc, req, ""})
		}
	}
	//now build a list of the possible exports of the whole graph
	possibleExport := []string{}
	for _, v := range n.dependencyGraph {
		for _, export := range v.export {
			possibleExport = append(possibleExport, export)
		}
	}
	log.Printf("possible exports: %+v", possibleExport)
	// strip out all the candidates who could be unlocked by the export
	result := []*depPair{}
outer:
	for _, candidate := range candidateList {
		for _, export := range possibleExport {
			if candidate.service == export {
				continue outer
			}
		}
		log.Printf("candidate passed all exports: %+v", candidate)
		result = append(result, candidate)
	}
	var buf bytes.Buffer
	for _, r := range result {
		buf.WriteString(r.String())
	}
	return buf.String()
}

func (n *NameServer) SendLoopMessage() {
	return
}
