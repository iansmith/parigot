package sys

import (
	"fmt"
	"strings"

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
	serviceIdToServiceData map[string] /*really service id*/ *serviceData // accelerator only, we could walk to find this
	dependencyGraph        *dep.DepGraph
	alreadyExported        []string
}

func NewNSCore() *NSCore {
	return &NSCore{
		packageRegistry:        make(map[string]*packageData),
		serviceIdToServiceData: make(map[string]*serviceData),
		dependencyGraph:        dep.NewDepGraph(),
		serviceCounter:         7, //first 8 are reserved
		alreadyExported:        []string{},
	}
}

type depKeyImpl struct {
	proc *Process
	addr string
}

func (d *depKeyImpl) String() string {
	if d.proc != nil {
		return d.proc.String()
	}
	return d.addr
}
func (d *depKeyImpl) IsKey() bool {
	return true
}

func newDepKeyFromProcess(p *Process) *depKeyImpl {
	return &depKeyImpl{proc: p}
}
func newDepKeyFromAddr(a string) *depKeyImpl {
	return &depKeyImpl{addr: a}
}

type KeyNSPair struct {
	Key        dep.DepKey
	NameServer NameServer
}

func NewKeyNSPair(k dep.DepKey, NameServer NameServer) *KeyNSPair {
	return &KeyNSPair{Key: k, NameServer: NameServer}
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
	serviceId      lib.Id
	closed         bool
	exported       bool
	method         map[string]lib.Id
	methodIdToImpl map[string]dep.DepKey
}

func newServiceData() *serviceData {
	return &serviceData{
		serviceId:      nil,
		method:         make(map[string]lib.Id),
		methodIdToImpl: make(map[string]dep.DepKey),
	}
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

func (n *NSCore) CloseService(pkgPath, service string) lib.Id {
	sData, err := n.validatePkgAndService(pkgPath, service)
	if err != nil {
		return err
	}
	if sData.closed {
		return lib.NewKernelError(lib.KernelServiceAlreadyClosedOrExported)

	}
	sData.closed = true
	return nil
}

// validatePkgAndService is a utiliy to verify pkg and service name.  It should
// not be called directly, it should only be used by the functions of NSCore.
// It does not lock, like everything in NSCore.
func (n *NSCore) validatePkgAndService(pkgPath, service string) (*serviceData, lib.Id) {
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

func (n *NSCore) Export(key dep.DepKey, pkgPath, service string) lib.Id {
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
	sData.exported = true
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
	for _, s := range n.alreadyExported {
		if s == name {
			nameserverPrint("REQUIRE", "process %s required %s.%s but it is already exported",
				key.String(), pkgPath, service)
			alreadyExported = true
		}
	}

	// we create the namespaces if they are not there yet because the exporter
	// may not have registered yet
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

	node, _ := n.dependencyGraph.GetEdge(key)
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

			nscorePrint("RUNIFREADY ", "running %s", candidate.Key())
			fn(candidate.Key())
		} else {
			nscorePrint("RUNIFREADY ", "%s is not ready to run, number of candidates left is %d", candidate.Key(), len(candidateList))
		}
	}
	nscorePrint("RUN ", "we are returning to the RunReader loop")
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

func (n *NSCore) HandleMethod(key dep.DepKey, pkgPath, service, method string) (lib.Id, lib.Id) {
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
	sData.methodIdToImpl[result.String()] = key

	// xxx fixme, should be able to realize that a method does not exist and reject the attempt to
	// handle it

	//xxx fixme, there should be a limit on the number of methods per service
	return result, nil
}

func nscorePrint(method, spec string, arg ...interface{}) {
	if nscoreVerbose {
		part1 := fmt.Sprintf("NSCore:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}
