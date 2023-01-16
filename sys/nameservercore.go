package sys

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/sys/backdoor"
	"github.com/iansmith/parigot/sys/dep"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var nscoreVerbose = false || os.Getenv("PARIGOT_VERBOSE") != ""

// NScore is used by both the local and remote (net) name server implementations
// to manage all the dependencies and handle require, export, and runWait.
// NSCore has a lock.  This lock is only used for fields that are not using
// sync.Map.
//
// All of the maps here are using sync.Map because it is optimized for the case
// where data is written once and read many times.  This is true for all three
// of these maps.
type NSCore struct {
	// key: package name, like "foo.bar", value: another sync.Map
	// inner map has key: service name like "baz" and a serviceData
	// the package registry is only used when looking up a service by fully qualified name
	// like "foo.bar.baz"
	packageRegistry *sync.Map
	// key: string rep of the service id, value: *serviceData
	serviceRegistry *sync.Map
	// key: string rep of a call id, value: *callContext
	inFlight           *sync.Map
	serviceCounter_    int
	dependencyGraph_   *dep.DepGraph
	alreadyExported_   []string
	useLocalServiceId_ bool
	lock               *sync.Mutex
}

func NewNSCore(useLocalServiceId bool) *NSCore {
	return &NSCore{
		packageRegistry:    &sync.Map{},
		serviceRegistry:    &sync.Map{},
		dependencyGraph_:   dep.NewDepGraph(),
		serviceCounter_:    7, //first 8 are reserved
		alreadyExported_:   []string{},
		useLocalServiceId_: useLocalServiceId,
		inFlight:           &sync.Map{},
		lock:               &sync.Mutex{},
	}
}

// DepKeyImpl is the implementation of DepKey. It is used so that we can refer
// to a "key" that represents a service in the same way, no matter if that service
// is local or remote. In the local case proc will be non-nil otherwise addr is
// the host and port of the remote service.
type DepKeyImpl struct {
	proc *Process
	addr string
}

func (n *NSCore) nextServiceCounter() int {
	n.lock.Lock()
	defer n.lock.Unlock()
	next := n.serviceCounter_ + 1
	n.serviceCounter_++
	return next
}

func (n *NSCore) addAlreadyExported(export ...string) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.alreadyExported_ = append(n.alreadyExported_, export...)
}

func (n *NSCore) useLocalServiceId() bool {
	n.lock.Lock()
	defer n.lock.Unlock()
	return n.useLocalServiceId_
}

func (n *NSCore) copyAlreadyExported() []string {
	n.lock.Lock()
	defer n.lock.Unlock()
	result := make([]string, len(n.alreadyExported_))
	copy(result, n.alreadyExported_)
	return result
}

func (d *DepKeyImpl) Name() string {
	if d.proc != nil {
		return d.proc.microservice.GetName()
	}
	return d.addr // close enough to a name
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

// walkDependencyGraph should be used to obtain the correct lock and then
// change the edge holders.  This function does not allow you to add or
// remove entire nodes.
func (n *NSCore) walkDependencyGraph(fn func(key string, eh *dep.EdgeHolder) bool) {
	n.dependencyGraph_.Walk(fn)
}

type KeyNSPair struct {
	Key        dep.DepKey
	NameServer NameServer
}

func NewKeyNSPair(k dep.DepKey, NameServer NameServer) *KeyNSPair {
	return &KeyNSPair{Key: k, NameServer: NameServer}
}

type ServiceData struct {
	serviceId lib.Id
	closed    bool
	exported  bool
	method    *sync.Map
	//map[string]lib.Id
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
		method:    &sync.Map{},
		key:       nil,
	}
}

// newServiceId is called to create a new service id.  If n.useLocalServiceId
// then the service id's content is not random but small (increasing) int.
func (n *NSCore) newServiceId() lib.Id {
	if n.useLocalServiceId() {
		latest := n.nextServiceCounter()
		return lib.LocalId[*protosupportmsg.ServiceId](uint64(latest))
	}
	_ = n.nextServiceCounter()
	return lib.NewId[*protosupportmsg.ServiceId]()
}

// ServiceData returns the serviceData for the service that has a given ID or nil
// if that id is not known.
func (n *NSCore) ServiceData(serviceId lib.Id) *ServiceData {
	sData, ok := n.serviceRegistry.Load(serviceId.String())
	if !ok {
		return nil
	}
	return sData.(*ServiceData)
}

// GetService is used when you know the package (full name) and the service within
// that package you want to find.  If it was found, you'll get back the serviceId.
// It's better to use the serviceId if possible as that lookup (like serviceFromServiceId)
// is faster.
func (n *NSCore) GetService(pkgPath, service string) (lib.Id, lib.KernelErrorCode) {
	pDataAny, ok := n.packageRegistry.Load(pkgPath)
	if !ok {
		return nil, lib.KernelNotFound
	}
	pData := pDataAny.(*sync.Map)
	sDataAny, ok := pData.Load(service)
	if !ok {
		return nil, lib.KernelNotFound
	}
	sData := sDataAny.(*ServiceData)
	return sData.serviceId, lib.KernelNoError
}

// CloseService is used to indicate that 1) the given service will not have
// more methods being registered to it and thus NotFound can be given for any
// methods not know after this point and 2) that the given service exists.
func (n *NSCore) CloseService(key dep.DepKey, pkgPath, service string) lib.Id {
	sData := n.create(key, pkgPath, service)
	sData.closed = true
	return nil
}

// validatePkgAndService is a utiliy to verify pkg and service name.
func (n *NSCore) validatePkgAndService(pkgPath, service string) (*ServiceData, lib.Id) {
	pDataAny, ok := n.packageRegistry.Load(pkgPath)
	if !ok {
		return nil, lib.NewKernelError(lib.KernelNotFound)
	}
	pData := pDataAny.(*sync.Map)
	sDataAny, ok := pData.Load(service)
	if !ok {
		return nil, lib.NewKernelError(lib.KernelNotFound)
	}
	sData := sDataAny.(*ServiceData)
	return sData, nil
}

// GetSData is a convenience wrapper around validatePkgAndService that returns
// nil instead of an error id.  This is useful in go code that is checking on
// the SData internally.
func (n *NSCore) GetSData(pkgPath, service string) *ServiceData {
	sData, err := n.validatePkgAndService(pkgPath, service)
	if err != nil && err.IsError() {
		return nil
	}
	return sData
}

// create is called by client code that wants to be sure that a given
// package and service is known.  It is, in some sense, the
// opposite of validatePackageAndService.
func (n *NSCore) create(key dep.DepKey, pkgPath, service string) *ServiceData {
	sid := n.newServiceId()
	return n.CreateWithSid(key, pkgPath, service, sid)
}

// DEBUG_DumpSIDTables prints the packages and their services to the terminal.
// Note this only will show output if nscoreVerbose is true or the environment
// variable PARIGOT_VERBOSE!=""
func (n *NSCore) DEBUG_DumpSIDTables() {
	n.packageRegistry.Range(func(pkg, innerMapAny any) bool {
		nscorePrint("DUMP", "nscore package %s -> %p", pkg, innerMapAny)
		innerMap := innerMapAny.(*sync.Map)
		innerMap.Range(func(service, sDataAny any) bool {
			sdata := sDataAny.(*ServiceData)
			nscorePrint("DUMP", "\t nscore service %s -> %p", service, sdata)
			nscorePrint("DUMP", "\t service id on sdata %s", sdata.serviceId.Short())
			return true
		})
		return true
	})
}

// createWithSid means that the caller wants to pick the service id for the
// new service being created.
func (n *NSCore) CreateWithSid(key dep.DepKey, pkgPath, service string, sid lib.Id) *ServiceData {
	var pData *sync.Map
	nscorePrint("create with sid", "1 gid=%x", GetGID())
	pDataAny, ok := n.packageRegistry.Load(pkgPath)
	if !ok {
		pData = &sync.Map{}
		n.packageRegistry.Store(pkgPath, pData)
	} else {
		pData = pDataAny.(*sync.Map)
	}
	nscorePrint("create with sid", "2 gid=%x", GetGID())
	var sData *ServiceData
	sDataAny, ok := pData.Load(service)
	nscorePrint("create with sid", "3 gid=%x", GetGID())
	if !ok {
		sData = NewServiceData(sid)
		pData.Store(service, sData)
		n.serviceRegistry.Store(sData.serviceId.String(), sData)
		nscorePrint("CREATE ", "created new service record: %s for %s.%s",
			sid.Short(), pkgPath, service)
		sData.key = key
	} else {
		sData = sDataAny.(*ServiceData)
	}
	nscorePrint("create with sid", "4 gid=%x", GetGID())

	return sData
}

// Export tells the core that the given key is associated with the implementation
// of the given pkgPath.Service.  The service must exist or a KernelNotFound error
// wil result.  This function locks because it needs to be sure the dependency
// graph is not changed out from under it.
func (n *NSCore) Export(key dep.DepKey, pkgPath, service string, newSid lib.Id) lib.Id {

	nscorePrint("EXPORT ", "process %s exports %s.%s",
		key.String(), pkgPath, service)

	// this region needs to lock because of the changes to sData which is returned
	// above, but without the lock... so this is
	lockRegion := func() lib.Id {
		n.lock.Lock()
		defer n.lock.Unlock()

		sData, err := n.validatePkgAndService(pkgPath, service)
		if err != nil {
			return err
		}
		// should we throw an error if this is not closed yet?

		if sData.exported {
			return lib.NewKernelError(lib.KernelServiceAlreadyClosedOrExported)
		}
		if newSid != nil {
			nscorePrint("EXPORT ", "rewriting the sid of the service, now that it is exported: %s", newSid.Short())
			oldSid := sData.serviceId
			sData.serviceId = newSid
			n.serviceRegistry.Delete(oldSid.String())
			n.serviceRegistry.Store(newSid.String(), sData)
		}
		//mark exported
		sData.exported = true
		sData.key = key
		node, ok := n.dependencyGraph_.GetEdge(key)
		if !ok {
			node = dep.NewEdgeHolder(key)
			n.dependencyGraph_.PutEdge(key, node)
		}
		node.AddExport(fmt.Sprintf("%s.%s", pkgPath, service))
		return nil
	}
	return lockRegion()
}

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// Require indicates a required package that must be up (key) before the given
// service (pkgPath.service) can start running.  This function locks because it
// needs exclusive access to the dep graph.
func (n *NSCore) Require(key dep.DepKey, pkgPath, service string) lib.Id {
	nscorePrint("Require", "entered part1  + gid=%x", GetGID())
	alreadyExported := false
	nscorePrint("Require ", "process %s requires %s.%s",
		key.String(), pkgPath, service)

	name := fmt.Sprintf("%s.%s", pkgPath, service)
	nscorePrint("Require", "entered part2 %x", GetGID())
	// this check always fails in the remote case, but n.alreadyExported==nil so it doesn't really hurt
	// n.copyAlreadyExported asserts the lock on NSCore so we can't lock before this
	// XXXYYY alreadyExported is not locked?
	for _, s := range n.alreadyExported_ {
		print(fmt.Sprintf("part2a %s\n", s))
		if s == name {
			nscorePrint("Require ", "process %s required %s.%s but it is already exported",
				key.String(), pkgPath, service)
			alreadyExported = true
		}
	}
	nscorePrint("Require ", "entered part3 gid= %x", GetGID())
	g := GetGID()
	nscorePrint("Require", "xxx about to lock region %x\n", g)
	lockRegion := func() lib.Id {
		n.lock.Lock()
		defer n.lock.Unlock()
		nscorePrint("require.lockRegion", "part 1 on gid %x", g)

		// we create the namespaces if they are not there yet because the exporter
		// may not have registered yet
		n.create(key, pkgPath, service)
		nscorePrint("require.lockRegion", "part 1a on gid %x", g)

		node, ok := n.dependencyGraph_.GetEdge(key)
		if !ok {
			nscorePrint("require.lockRegion", "part 1b on gid %x", g)
			node = dep.NewEdgeHolder(key)
			nscorePrint("require.lockRegion", "part 1c on gid %x", g)
			n.dependencyGraph_.PutEdge(key, node)
			nscorePrint("require.lockRegion", "part 1d on gid %x", g)
		}
		nscorePrint("Require", "process %s did a require for %s.%s and already exported? %v (gid %x)",
			key.String(), pkgPath, service, alreadyExported, g)
		if !alreadyExported {
			found := false
			nscorePrint("lockRegion", "gid %x about to hit the walk", g)
			node.WalkRequire(func(s string) bool {
				if s == name {
					found = true
					return false
				}
				return true
			})
			if found {
				nscorePrint("Require ", "process %s ERROR RET (%x)", key.String(), g)
				nscorePrint("xxx require 4: %s vs %s", key.String(), name)
				print(fmt.Sprintf("xxx about to leave ERR lockRegion %x\n", g))
				return lib.NewKernelError(lib.KernelServiceAlreadyRequired)
			}
			nscorePrint("Require ", "about to add require %x, %s", g, name)
			node.AddRequire(name)
		}
		nscorePrint("Require ", "process %s  NORMAL RET %x", key.String(), GetGID())
		print(fmt.Sprintf("xxx about to leave lockRegion %x\n", g))
		return nil
	}
	return lockRegion()
}

// RunIfReady checks the dependencies of the given service (key) and if the service
// is ready to run--all the dependencies of the service are met--then it returns the
// key of the newly ready service.  This call has no effect if the service given by
// key is not yet ready to run.
//
// Because this process of determining if a given service is ready to run may
// involve _other_ services being declared ready to run (and thus returned from
// this function) we have to be careful to insure deterministic ordering
// of the results.  This function checks on each iteration to see if any
// new services are ready to run.  The newly ready services are sorted alphabetically
// by the name of the service used in the deployment file.  This ensures that in the case
// where multiple services are "made ready" by the same iteration of this process
// the order of results is known.
//
// RunIfReady locks to make sure that the depnedency graph isn't changed out from
// under it.
func (n *NSCore) RunIfReady(key dep.DepKey) []dep.DepKey {
	n.lock.Lock()
	defer n.lock.Unlock()

	result := []dep.DepKey{}

	nscorePrint("RunIfReady ", "-----> called for key=%s  <------", key)
	// the only time things can change is when a new process calls run
	// so we only to make sure here that we find all the eligible processes to run
	print(fmt.Sprintf("zzzz starting runIfReady in core for %s\n", key))
	node, ok := n.dependencyGraph_.GetEdge(key)
	if !ok {
		nscorePrint("RunIfReady ", "Ignoring request to check on key %s -- no edges assuming ready to run", key.String())
		nscorePrint("RunIfReady ", "state of graph %#v", n.dependencyGraph_.AllEdge())
		return []dep.DepKey{key}
	}
	candidateList := []*dep.EdgeHolder{node}

	nscorePrint("RunIfReady ", "node %s (%d req,%d exp) and dep-graph has %d total entries",
		key, node.RequireLen(), node.ExportLen(), n.dependencyGraph_.Len())
	for len(candidateList) > 0 {
		print(fmt.Sprintf("zzz runIfReady %s has %d candidates in list\n", key, len(candidateList)))
		newCandidates := []*dep.EdgeHolder{}
		readyList := []string{}
		readyMap := make(map[string]*dep.EdgeHolder)
		candidate := candidateList[0]
		// remove candidate from list
		if len(candidateList) == 1 {
			candidateList = nil
		} else {
			candidateList = candidateList[1:]
		}
		// is candidate ready to run?
		if candidate.IsReady() {
			nscorePrint("RunIfReady ", "candidate %s is ready to run", candidate.Key())
			key := candidate.Key()
			n.dependencyGraph_.Del(key)
			// we are ready, so lets process his exports through the list of waiting processes
			exports := candidate.Export()
			print(fmt.Sprintf("zzz about to walk graph for exports %s: %+v\n", key, exports))
			n.dependencyGraph_.Walk(func(key string, other *dep.EdgeHolder) bool {
				print(fmt.Sprintf("zzz WALK %s.RemoveRequire(%+v)\n", other.Key(), exports))
				changed := other.RemoveRequire(exports)
				if changed {
					newCandidates = append(newCandidates, other)
					nscorePrint("RunIfReady ", "candidate list changed")
				}
				return true
			})
			print(fmt.Sprintf("zzz done with walk for exports %s: %+v\n", key, newCandidates))
			print(fmt.Sprintf("zzz about to call addAlreadyExported for %s: %+v\n", key, candidate.Export()))
			n.addAlreadyExported(candidate.Export()...)
			nscorePrint("RunIfReady ", "already exported updated: %+v", n.copyAlreadyExported())

			nscorePrint("RunIfReady ", "%s is on ready list", candidate.Key())
			print(fmt.Sprintf("zzz about to add %s to ready list to ready list %+v\n", candidate.Key().Name(), readyList))
			readyList = append(readyList, candidate.Key().Name())
			print(fmt.Sprintf("zzz about to add %s ready map %+v\n", candidate.Key().Name(), readyMap))
			readyMap[candidate.Key().Name()] = candidate
			print(fmt.Sprintf("zzz ready map for %s has %d entries, list is %+v\n", key, len(readyMap), readyList))
		} else {
			nscorePrint("RunIfReady ", "%s is not ready to run, number of candidates left is %d", candidate.Key(), len(candidateList))
		}
		//update datastructures and start in alpha order
		sort.Strings(readyList)
		for _, readyName := range readyList {
			nscorePrint("RunIfReady ", "adding %s to result list", readyName)
			result = append(result, readyMap[readyName].Key())
			print(fmt.Sprintf("zzz added %s and new result list has %d entries\n", readyMap[readyName].Key(), len(result)))
		}
		candidateList = append(candidateList, newCandidates...)
		print(fmt.Sprintf("zzz updated candidate list for %s is %+v\n", key, candidateList))
	}
	print(fmt.Sprintf("zzzz returning from runIfReady in core for %s\n", key))

	nscorePrint("RunIfReady ", "returning:"+key.String())
	return result
}

// Started failed info computes a string that is intended to be shown to the
// developer.  This string indicates the dependency cycle or the dead processes
// that have been found in the deployment file.  A dead process is one who
// has one or more dependencies that are not present at all in the set of
// services to be deployed.
// We lock in this function because we need to traverse the depnedency graph
// and it must be stable.
func (n *NSCore) StartFailedInfo() string {
	n.lock.Lock()
	defer n.lock.Unlock()

	result := ""
	if n.WaitingToRun() > 0 {
		nscorePrint("StartFailed ", "waiting to run %d", n.WaitingToRun())

		nscorePrint("StartFailed", "was not able to get all processes started due to export/require problems")
		loop := n.dependencyGraph_.GetLoopContent()
		dead := n.dependencyGraph_.GetDeadNodeContent()
		if loop != "" {
			loop = strings.Replace(loop, ";", "\n", -1)
			result += fmt.Sprintf("Loop discovered in the dependencies\n%s\n", loop)
		}
		if dead != "" {
			dead = strings.Replace(dead, ";", "\n", -1)
			result += fmt.Sprintf("Dead processes are processes that cannot start because no other process exports what they require:\n%s\n", dead)
		}
		result += "aborting due to export/require problems"
	}
	return result
}

// WaitingToRun returns the number of processes that are not yet started.
// We lock in this function for exclusive access to the dep graph.
func (n *NSCore) WaitingToRun() int {
	n.lock.Lock()
	defer n.lock.Unlock()

	return n.dependencyGraph_.Len()
}

// FindOrCreateMethodId takes a given service by name (packagePath.service) and
// either generates or finds the methodId for the given name.   This function can
// be avoided by caching the returned method id.
func (n *NSCore) FindOrCreateMethodId(key dep.DepKey, packagePath, service, method string) lib.Id {
	sData, err := n.validatePkgAndService(packagePath, service)
	if err != nil && err.IsError() {
		nscorePrint("FindOrCreateMethodId ", "we need to create a service data for %s.%s", packagePath, service)
		if !err.Equal(lib.NewKernelError(lib.KernelNotFound)) {
			nscorePrint("FindOrCreateMethodId ", "WARN unable to understand error from validatePackage() %s", err.Short())
			return nil
		}
		sData = n.create(key, packagePath, service)
	}
	var mid lib.Id
	methodAny, ok := sData.method.Load(method)
	if !ok {
		nscorePrint("FindOrCreateMethodId ", "we need to create a method id for %s.%s.%s", packagePath, service, method)
		mid = lib.NewId[*protosupportmsg.MethodId]()
		sData.method.Store(method, mid)
		nscorePrint("FindOrCreateMethodId ", "added method %s to sdata", mid.Short())
	} else {
		mid = methodAny.(lib.Id)
	}
	return mid

}

func nscorePrint(method, spec string, arg ...interface{}) {
	if nscoreVerbose {
		part1 := fmt.Sprintf("NSCore:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		req := &logmsg.LogRequest{
			Level:   logmsg.LogLevel_LOG_LEVEL_DEBUG,
			Stamp:   timestamppb.Now(), // xxx should use the kernel version
			Message: fmt.Sprintf("%s%s\n", part1, part2),
		}
		backdoor.Log(req, true, false, false, nil)
	}
}

// getContextForCallId returns the callContext for a given call id.  The call context
// has all the information for all stages of processing.   If this function returns
// nil that's bad because it means that we couldn't find a reference to the call id,
// and that should not happen.  If we return the call context we _also_ have removed
// that cid from the mapping, since you should only do this lookup once.
func (n *NSCore) getContextForCallId(cid lib.Id) *callContext {
	resultAny, ok := n.inFlight.LoadAndDelete(cid.String())
	if !ok {
		return nil // ugh, serious problem
	}
	result := resultAny.(*callContext)
	return result
}

// addCallContextMapping is used to add a new mapping from cid (actually cid.String()) to
// the given call context. If the mapping already exists, this function panics
// because that indicates that we are doing this for a second time for the same
// call.
func (n *NSCore) addCallContextMapping(cid lib.Id, cctx *callContext) {
	_, ok := n.inFlight.Load(cid.String())
	if ok {
		panic("found already existing element in inFlight mapping for " + cid.String())
	}
	n.inFlight.Store(cid.String(), cctx)
}
