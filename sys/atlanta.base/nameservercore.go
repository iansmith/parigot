package sys

import (
	"fmt"
	"strings"

	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/dep"
)

var nscoreVerbose = false

// NScore is used by both the local and remote (net) name server implementations
// to manage all the dependencies and handle require, export, and runWait.
// NSCore does not have a lock; it expects the higher level to handle that.
type NSCore struct {
	packageRegistry        map[string]*packageData
	serviceCounter         int
	serviceIdToServiceData map[string] /*really service id*/ *ServiceData // accelerator only, we could walk to find this
	dependencyGraph        *dep.DepGraph
	alreadyExported        []string
	useLocalServiceId      bool
	inFlight               map[string]*callContext // key is really lib.Id
}

func NewNSCore(useLocalServiceId bool) *NSCore {
	return &NSCore{
		packageRegistry:        make(map[string]*packageData),
		serviceIdToServiceData: make(map[string]*ServiceData),
		dependencyGraph:        dep.NewDepGraph(),
		serviceCounter:         7, //first 8 are reserved
		alreadyExported:        []string{},
		useLocalServiceId:      useLocalServiceId,
		inFlight:               make(map[string]*callContext), // key is really lib.Id
	}
}

type DepKeyImpl struct {
	proc *Process
	addr string
}

func (d *DepKeyImpl) String() string {
	if d.proc != nil {
		return d.proc.String()
	}
	return d.addr
}
func (d *DepKeyImpl) IsKey() bool {
	return true
}

func NewDepKeyFromProcess(p *Process) *DepKeyImpl {
	return &DepKeyImpl{proc: p}
}
func NewDepKeyFromAddr(a string) *DepKeyImpl {
	return &DepKeyImpl{addr: a}
}

func (n *NSCore) DependencyGraph() *dep.DepGraph {
	return n.dependencyGraph
}

type KeyNSPair struct {
	Key        dep.DepKey
	NameServer NameServer
}

func NewKeyNSPair(k dep.DepKey, NameServer NameServer) *KeyNSPair {
	return &KeyNSPair{Key: k, NameServer: NameServer}
}

type packageData struct {
	service map[string]*ServiceData
}

func newPackageData() *packageData {
	return &packageData{
		service: make(map[string]*ServiceData),
	}
}

type ServiceData struct {
	serviceId lib.Id
	closed    bool
	exported  bool
	method    map[string]lib.Id
	//methodIdToImpl map[string]dep.DepKey
	key dep.DepKey
}

func (s *ServiceData) GetServiceId() lib.Id {
	return s.serviceId
}

func (s *ServiceData) GetKey() dep.DepKey {
	return s.key
}

func NewServiceData(sid lib.Id) *ServiceData {
	return &ServiceData{
		serviceId: sid,
		method:    make(map[string]lib.Id),
		//methodIdToImpl: make(map[string]dep.DepKey),
		key: nil,
	}
}

// newServiceId is called to create a new service id with a strategy based
// on the useLocalServiceId.
func (n *NSCore) newServiceId() lib.Id {
	if n.useLocalServiceId {
		n.serviceCounter++
		return lib.LocalId[*protosupport.ServiceId](uint64(n.serviceCounter))
	}
	return lib.NewId[*protosupport.ServiceId]()
}

func (n *NSCore) GetService(pkgPath, service string) (lib.Id, lib.Id) {
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

// CloseService is used to indicate that 1) the given service will not have
// more methods being registered to it and thus NotFound can be given for any
// methods not know after this point and 2) that the given service exists.
func (n *NSCore) CloseService(key dep.DepKey, pkgPath, service string) lib.Id {
	sData := n.create(key, pkgPath, service)
	sData.closed = true
	return nil
}

// validatePkgAndService is a utiliy to verify pkg and service name.  It should
// not be called directly, it should only be used by the functions of NSCore.
// It does not lock, like everything in NSCore.
func (n *NSCore) validatePkgAndService(pkgPath, service string) (*ServiceData, lib.Id) {
	pData, ok := n.packageRegistry[pkgPath]
	if !ok {
		return nil, lib.NewKernelError(lib.KernelNotFound)
	}
	sData, ok := pData.service[service]
	if !ok {
		return nil, lib.NewKernelError(lib.KernelNotFound)
	}
	return sData, nil
}

// GetSData is a convenience wrapper around validatePkgAndService
func (n *NSCore) GetSData(pkgPath, service string) *ServiceData {
	sData, err := n.validatePkgAndService(pkgPath, service)
	if err != nil && err.IsError() {
		return nil
	}
	return sData
}

// GetSDataById returns the service data for a given id.
func (n *NSCore) GetSDataById(sid lib.Id) *ServiceData {
	return n.serviceIdToServiceData[sid.String()]
}

// create is called by client code that wants to be sure that a given
// package is known.  It is, in some sense, the opposite of validatePackageAndService.
func (n *NSCore) create(key dep.DepKey, pkgPath, service string) *ServiceData {
	sid := n.newServiceId()
	return n.CreateWithSid(key, pkgPath, service, sid)
}

func (n *NSCore) DumpSIDTables() {
	for pkg, pData := range n.packageRegistry {
		nscorePrint("DUMP", "nscore package %s -> %p", pkg, pData)
		for service, sdata := range pData.service {
			nscorePrint("DUMP", "\t nscore service %s -> %p", service, sdata)
			nscorePrint("DUMP", "\t service id on sdata %s", sdata.serviceId.Short())
		}
	}
	for sid, sdata := range n.serviceIdToServiceData {
		nscorePrint("DUMP", "mappnig: service id %s -> %p, with sid %s", sid, sdata, sdata.serviceId.String())
	}
}

// createWithSid means that the caller wants to pick the service id for the
// new service being created.
func (n *NSCore) CreateWithSid(key dep.DepKey, pkgPath, service string, sid lib.Id) *ServiceData {
	pData, ok := n.packageRegistry[pkgPath]
	if !ok {
		pData = newPackageData()
		n.packageRegistry[pkgPath] = pData
	}
	sData, ok := pData.service[service]
	if !ok {
		sData = NewServiceData(sid)
		pData.service[service] = sData
		n.serviceIdToServiceData[sData.serviceId.String()] = sData
		nscorePrint("CREATE", "created new service record: %s==%s",
			sid.Short(), sData.serviceId.Short())
		sData.key = key
	}
	return sData
}

// Export tells the core that the given key is associated with the implementation
// of the given pkgPath.Service.  The service must exist or a KernelNotFound error
// wil result.
func (n *NSCore) Export(key dep.DepKey, pkgPath, service string, newSid lib.Id) lib.Id {
	nscorePrint("EXPORT", "process %s exports %s.%s",
		key.String(), pkgPath, service)
	sData, err := n.validatePkgAndService(pkgPath, service)
	if err != nil {
		return err
	}
	if !sData.closed {
		// xxxfix me, should this be an error?
	}
	if sData.exported {
		return lib.NewKernelError(lib.KernelServiceAlreadyClosedOrExported)
	}
	if newSid != nil {
		nscorePrint("EXPORT", "rewriting the sid of the service, now that it is exported: %s", newSid.Short())
		oldSid := sData.serviceId
		sData.serviceId = newSid
		delete(n.serviceIdToServiceData, oldSid.String())
		n.serviceIdToServiceData[newSid.String()] = sData
	}

	sData.exported = true
	sData.key = key
	node, ok := n.dependencyGraph.GetEdge(key)
	if !ok {
		node = dep.NewEdgeHolder(key)
		n.dependencyGraph.PutEdge(key, node)
	}
	node.AddExport(fmt.Sprintf("%s.%s", pkgPath, service))
	return nil
}

func (n *NSCore) Require(key dep.DepKey, pkgPath, service string) lib.Id {
	alreadyExported := false
	nscorePrint("REQUIRE", "process %s requires %s.%s",
		key.String(), pkgPath, service)

	name := fmt.Sprintf("%s.%s", pkgPath, service)

	// this check always fails in the remote case, but n.alreadyExported==nil so it doesn't really hurt
	for _, s := range n.alreadyExported {
		if s == name {
			nameserverPrint("REQUIRE", "process %s required %s.%s but it is already exported",
				key.String(), pkgPath, service)
			alreadyExported = true
		}
	}

	// we create the namespaces if they are not there yet because the exporter
	// may not have registered yet
	n.create(key, pkgPath, service)

	node, ok := n.dependencyGraph.GetEdge(key)
	if !ok {
		node = dep.NewEdgeHolder(key)
		n.dependencyGraph.PutEdge(key, node)
	}
	if !alreadyExported {
		for _, r := range node.Require() {
			if r == name {
				return lib.NewKernelError(lib.KernelServiceAlreadyRequired)
			}
		}
		node.AddRequire(name)
	}
	return nil
}

func (n *NSCore) RunIfReady(key dep.DepKey, fn func(dep.DepKey)) {
	nscorePrint("RUNIFREADY ", "-----> called for key=%s", key)
	// the only time things can change is when a new process calls run
	// so we only to make sure here that we find all the eligible processes to run

	node, ok := n.dependencyGraph.GetEdge(key)
	if !ok {
		nscorePrint("RUNIFREADY ", "Ignoring request to check on key %s", key)
		nscorePrint("RUNIFREADY ", "state of graph %#v", n.dependencyGraph.AllEdge())
		return
	}
	candidateList := []*dep.EdgeHolder{node}

	nscorePrint("RUNIFREADY ", "node %s (%d req,%d exp) and dep-graph has %d total entries",
		key, node.RequireLen(), node.ExportLen(), n.dependencyGraph.Len())
	for len(candidateList) > 0 {
		candidate := candidateList[0]
		// remove candidate from list
		if len(candidateList) == 1 {
			candidateList = nil
		} else {
			candidateList = candidateList[1:]
		}
		// is candidate ready to run?
		if candidate.IsReady() {
			nscorePrint("RUNIFREADY ", "candidate %s is ready to run", candidate.Key())
			key := candidate.Key()
			n.dependencyGraph.Del(key)
			// we are ready, so lets process his exports through the list of waiting processes
			for _, other := range n.dependencyGraph.AllEdge() {
				changed := other.RemoveRequire(candidate.Export())
				if changed {
					candidateList = append(candidateList, other)
					nscorePrint("RUNIFREADY ", "candidate list changed")
				}
			}
			n.alreadyExported = append(n.alreadyExported, candidate.Export()...)
			nscorePrint("RUNIFREADY ", "already exported updated: %+v", n.alreadyExported)

			nscorePrint("RUNIFREADY ", "notifying that %s is ready", candidate.Key())
			fn(candidate.Key())
		} else {
			nscorePrint("RUNIFREADY ", "%s is not ready to run, number of candidates left is %d", candidate.Key(), len(candidateList))
		}
	}
	nscorePrint("RUN ", "blocking completed")
}

func (n *NSCore) StartFailedInfo() string {
	result := ""
	if n.WaitingToRun() > 0 {
		nscorePrint("StartFailed ", "waiting to run %d", n.WaitingToRun())

		nscorePrint("StartFailed", "was not able to get all processes started due to export/require problems")
		loop := n.dependencyGraph.GetLoopContent()
		dead := n.dependencyGraph.GetDeadNodeContent()
		if loop != "" {
			loop = strings.Replace(loop, ";", "\n", -1)
			result += fmt.Sprintf("Loop discovered in the dependencies\n%s\n", loop)
		}
		if dead != "" {
			dead = strings.Replace(dead, ";", "\n", -1)
			result += fmt.Sprintf("Dead processes are processes that cannot start because no other process exports what they require:\n%s\n", dead)
		}
		result += fmt.Sprintf("aborting due to export/require problems")
	}
	return result
}

func (n *NSCore) WaitingToRun() int {
	return n.dependencyGraph.Len()
}

func (n *NSCore) FindOrCreateMethodId(key dep.DepKey, packagePath, service, method string) lib.Id {
	sData, err := n.validatePkgAndService(packagePath, service)
	if err != nil && err.IsError() {
		nscorePrint("FINDORCREATEMID", "we need to create a service data for %s.%s", packagePath, service)
		if !err.Equal(lib.NewKernelError(lib.KernelNotFound)) {
			nscorePrint("FINDORCREATEMID ", "WARN unable to understand error from validatePackage() %s", err.Short())
			return nil
		}
		sData = n.create(key, packagePath, service)
	}
	mid, ok := sData.method[method]
	if !ok {
		nscorePrint("FINDORCREATEMID", "we need to create a method id for %s.%s.%s", packagePath, service, method)
		mid = lib.NewId[*protosupport.MethodId]()
		sData.method[method] = mid
	}
	return mid

}

// func (n *NSCore) HandleMethod(key dep.DepKey, pkgPath, service, method string) (lib.Id, lib.Id) {
// 	// create the data for this package and service
// 	sData := n.create(key, pkgPath, service)
// 	result := lib.NewMethodId()
// 	nameserverPrint("HANDLEMETHOD", "assigning %s to the method %s in service %s",
// 		result, method, service)
// 	sData.method[method] = result
// 	sData.methodIdToImpl[result.String()] = key

// 	// xxx fixme, should be able to realize that a method does not exist and reject the attempt to
// 	// handle it

// 	//xxx fixme, there should be a limit on the number of methods per service
// 	return result, nil
// }

func nscorePrint(method, spec string, arg ...interface{}) {
	if nscoreVerbose {
		part1 := fmt.Sprintf("NSCore:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}

// getContextForCallId returns the callContext for a given call id.  The call context
// has all the information for all stages of processing.   If this function returns
// nil that's bad because it means that we couldn't find a reference to the cid,
// and that should not happen.  If we return the call context we _also_ have removed
// that cid from the mapping, since you should only do this lookup once.
func (n *NSCore) getContextForCallId(cid lib.Id) *callContext {
	result := n.inFlight[cid.String()]
	delete(n.inFlight, cid.String())
	return result
}

// addCallContextMapping is used to add a new mapping from cid (actually cid.String()) to
// the given call context. If the mapping already exists, this function panics.
func (n *NSCore) addCallContextMapping(cid lib.Id, cctx *callContext) {
	_, ok := n.inFlight[cid.String()]
	if ok {
		panic("found already existing element in inFlight mapping for " + cid.String())
	}
	n.inFlight[cid.String()] = cctx
}
