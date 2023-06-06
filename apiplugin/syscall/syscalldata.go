package main

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/yourbasic/graph"
)

var _depData *syscallDataImpl

func depData() *syscallDataImpl {
	if _depData == nil {
		_depData = newSyscallDataImpl()
	}
	return _depData
}

//
// serviceImpl
//

// The default lock discipline for this type is that you should call a method
// on this type when the lock is unlocked.  The methods you can call when you
// DO have the lock are demarcated by the NoLock suffix. Any function that
// does assert the lock should make sure to release it before returning.

type serviceImpl struct {
	pkg, name string
	id        id.ServiceId
	runReady  bool
	exported  bool
	started   bool
	parent    *syscallDataImpl
	runCh     chan struct{}
	lock      *sync.Mutex
}

func newServiceImpl(pkg, name string, sid id.ServiceId, parent *syscallDataImpl, isClient bool) *serviceImpl {
	result := &serviceImpl{
		pkg:      pkg,
		name:     name,
		id:       sid,
		runReady: false,
		parent:   parent,
		exported: isClient,
		lock:     new(sync.Mutex),
		runCh:    make(chan struct{}),
	}
	_ = Service(result)
	return result
}

// wakeUp causes a send on the servicImpl's runCh and thus check to see
// if it can run now.  This method IS in use, even if VSCode seems confused
// about that.
func (s *serviceImpl) wakeUp() {
	print("attempting wakeup for ", s.name, "\n")
	s.runCh <- struct{}{}
}

// Id returns the id of this service.
func (s *serviceImpl) Id() id.ServiceId {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.id
}

// Name returns the name, not the fully qualified name, of this service.
func (s *serviceImpl) Name() string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.name
}

// Package returns the package name, not the fully qualified name, of this service.
func (s *serviceImpl) Package() string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.pkg
}

func (s *serviceImpl) RunRequested() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.runReady
}

func (s *serviceImpl) Exported() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.exported
}

func (s *serviceImpl) Run(ctx context.Context) bool {
	// slightly tricky with the lock
	s.lock.Lock()
	s.runReady = true
	s.lock.Unlock()
	return s.waitToRun(ctx)
}

func (s *serviceImpl) canRun(ctx context.Context) bool {
	if !s.Exported() {
		return false
	}
	if !s.RunRequested() {
		return false
	}
	pcontext.Debugf(ctx, "trying to see if %s [%s] can run now ", s.Short(), s.Name())
	depCheck := s.parent.checkNodesBehindForReady(ctx, s.String())
	return depCheck
}

// waitToRun waits until the timeout expires or until it receives a wake
// up call and a check for the ability to run successfully is made. It returns
// false if it is returning because of a timeout. Note that this function
// does not lock so that other things can proceed concurrently.
func (s *serviceImpl) waitToRun(ctx context.Context) bool {
	if s.canRun(ctx) {
		pcontext.Debugf(ctx, "%s is immediately ready to run", s.name)
		return true
	}
	pcontext.Debugf(ctx, "Timeout loop started for %s", s.Name())
	for {
		select {
		case <-s.runCh:
			if s.canRun(ctx) {
				s.lock.Lock()
				s.started = true
				s.lock.Unlock()
				return true
			}
		case <-time.After(5 * time.Second):
			return false
		}
	}
}

func (s *serviceImpl) export() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.exported = true
}
func (s *serviceImpl) String() string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.id.String()
}

func (s *serviceImpl) Short() string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.id.Short()
}

func (s *serviceImpl) Started() bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.started
}

//
// syscallDataImpl
//

// The default lock discipline for this type is that you should call a method
// on this type when the lock is unlocked.  The methods you can call when you
// DO have the lock are demarcated by the NoLock suffix. Any function that
// does assert the lock should make sure to release it before returning.

type syscallDataImpl struct {
	sidStringToService          map[string]Service
	packageNameToServiceNameMap map[string]map[string]Service
	depGraph                    *graph.Mutable
	vertexName                  map[string]int
	lock                        *sync.Mutex
}

func newSyscallDataImpl() *syscallDataImpl {
	g := graph.New(0)
	impl := &syscallDataImpl{
		sidStringToService:          make(map[string]Service),
		packageNameToServiceNameMap: make(map[string]map[string]Service),
		depGraph:                    g,
		lock:                        new(sync.Mutex),
		vertexName:                  make(map[string]int),
	}
	_ = SyscallData(impl)
	return impl
}

func (s *syscallDataImpl) SetService(ctx context.Context, package_, name string, client bool) (Service, bool) {

	svc := s.ServiceByName(ctx, package_, name)
	if svc != nil {
		return svc, false
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	svcId := id.NewServiceId()
	nmap, ok := s.packageNameToServiceNameMap[package_]
	if !ok {
		nmap = make(map[string]Service)
		s.packageNameToServiceNameMap[package_] = nmap
	}

	if svcId.IsEmptyValue() || svcId.IsZeroValue() {
		print("Service Id error, bad id generated by random!\n")
	}

	result := newServiceImpl(package_, name, svcId, s, client)

	nmap[name] = result
	s.sidStringToService[result.String()] = result
	if !s.addVertex(ctx, result.String()) {
		return nil, false
	}
	if result != nil {
		if result.id.IsEmptyValue() || result.id.IsZeroValue() {
			pcontext.Errorf(ctx, "Service Id error, bad id returned from syscall data")
		}
	} else {
		pcontext.Errorf(ctx, "result of set service is nil?")
	}
	pcontext.Debugf(ctx, "created service via SetService: %s [%s]", result.id.Short(), result.Name())
	return result, true

}

func (s *syscallDataImpl) ServiceByName(ctx context.Context, package_, name string) Service {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.serviceByNameNoLock(ctx, package_, name)
}

func (s *syscallDataImpl) serviceByNameNoLock(ctx context.Context, package_, name string) Service {

	nameMap, ok := s.packageNameToServiceNameMap[package_]
	if !ok {
		return nil
	}
	svc, ok := nameMap[name]
	if !ok {
		return nil
	}
	return svc
}

func (s *syscallDataImpl) ServiceById(ctx context.Context, sid id.ServiceId) Service {
	return s.ServiceByIdString(ctx, sid.String())
}
func (s *syscallDataImpl) serviceByIdStringNoLock(ctx context.Context, sid string) Service {
	svc, ok := s.sidStringToService[sid]
	if !ok {
		return nil
	}
	return svc
}

func (s *syscallDataImpl) ServiceByIdString(ctx context.Context, sid string) Service {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.serviceByIdStringNoLock(ctx, sid)
}

func (s *syscallDataImpl) Export(ctx context.Context, svcId id.ServiceId) Service {
	svc := s.ServiceByIdString(ctx, svcId.String())
	if svc == nil {
		return nil
	}

	svc.(*serviceImpl).export()
	if svc.(*serviceImpl).canRun(ctx) {
		s.notifyNodesBehindForReady(ctx, svc.String())
		pcontext.Debugf(ctx, "service %s [%s] is ready to run due to export", svc.Short(), svc.Name())
	}
	return svc
}

// topoSort does a topological sort of all the nodes we currently know about.
// This function asserts the lock to ensure that while the topo algorihtm is
// running no part of the graph is disturbed.  It returns what I HOPE is a copy
// of the content of the graph vertices.
func (s *syscallDataImpl) topoSort(ctx context.Context) []string {
	s.lock.Lock()
	defer s.lock.Unlock()

	result, ok := graph.TopSort(s.depGraph)
	if !ok {
		panic("topolical sort could not be generated, likely it has a cycle")
	}
	name := make([]string, len(result))
	for i, v := range result {
		str, ok := reverseMap(s.vertexName, v)
		if !ok {
			panic("badly formed dependency graph, cant find vertex:" + fmt.Sprint(v))
		}
		name[len(result)-1-i] = str
	}
	return name
}

// notifyNodesBehind walks the topologically ordered vertices, looking for any nodes
// that are predecessors of the given node and notifying them to check their status for running.
// This returns false only when the service cannot be found.
func (s *syscallDataImpl) notifyNodesBehindForReady(ctx context.Context, svcid string) bool {
	topo := s.topoSort(ctx)                     // locks
	if s.ServiceByIdString(ctx, svcid) == nil { //locks
		return false
	}
	for _, str := range topo {
		if str == svcid {
			pcontext.Debugf(ctx, "covered all predecessors of %s", svcid)
			return true
		}
		svc := s.ServiceByIdString(ctx, str)
		if svc == nil {
			pcontext.Fatalf(ctx, "internal error trying to notify nodes to check their can run status, can't find  "+str)
			panic("internal error trying to notify nodes to check their can run status, can't find  " + str)
		}
		svc.(*serviceImpl).wakeUp()
	}
	panic("did not find id " + svcid + " in the list of vertices")
}

// checkNodesBehindForReady walks the topologically ordered vertices, looking for any nodes
// that are predecessors of the given node and testing to see if they are started. If not
// we return false.  If all the predecessorys are started, then we return true.  If the
// serviceId cannot be found, we return true.
func (s *syscallDataImpl) checkNodesBehindForReady(ctx context.Context, svcid string) bool {
	topo := s.topoSort(ctx)                     // locks
	if s.ServiceByIdString(ctx, svcid) == nil { //locks
		return true
	}
	pcontext.Debugf(ctx, "topo: %+v", topo)
	success := true
	for _, str := range topo {
		pcontext.Debugf(ctx, "topo entry: %s", str)
		if str == svcid {
			pcontext.Debugf(ctx, "topo success")
			return success
		}
		svc := s.ServiceByIdString(ctx, str) //locks
		if svc == nil {
			panic("unable to walk the dep graph looking for predeessors, a predecessor could not be found: " + str)
		}
		if svc.Exported() && svc.RunRequested() && svc.Started() {
			pcontext.Debugf(ctx, "  -- check completed %s [%s]", svc.Short(), svc.Name())
			continue
		}
		pcontext.Debugf(ctx, "topo failure")
		return false
	}
	panic("did not find id " + svcid + " in the list of vertices")
}

// checkNodesInFront walks the topologically ordered vertices, looking for any nodes
// that have satisfied dependencies and are net yet running. It returns false only if the svcid
// cannot be found.
func (s *syscallDataImpl) checkNodesInFront(ctx context.Context, svcid string) bool {
	topo := s.topoSort(ctx)
	foundSelf := false
	for _, str := range topo {
		svc := s.ServiceByIdString(ctx, str)
		if svc == nil {
			return false
		}
		if svcid == svc.String() {
			foundSelf = true
			continue
		}
		if !foundSelf {
			continue
		}
		if svc.Started() {
			continue
		}
		pcontext.Infof(ctx, "notifying service %s to check its ready state", svc.Short())
		// wake em up
		svc.(*serviceImpl).wakeUp()
	}
	return true
}

// Import adds a node in the dependency graph between src and dest.
// It returns a kerr if either the source or dest cannot be found; it
// returns a kerr if the new edge would create a cycle.
func (s *syscallDataImpl) Import(ctx context.Context, src, dest id.ServiceId) id.KernelErrId {

	serviceSource := s.ServiceByIdString(ctx, src.String())
	if serviceSource == nil {
		return id.NewKernelErrId(id.KernelNotFound)
	}
	serviceDest := s.ServiceByIdString(ctx, dest.String())
	if serviceDest == nil {
		return id.NewKernelErrId(id.KernelNotFound)
	}
	srcString := serviceSource.String()
	destString := serviceDest.String()
	// lock for the graph
	s.lock.Lock()
	defer s.lock.Unlock()
	if !graph.Acyclic(s.depGraph) {
		panic("graph is already cyclic, some previous edge was added without checking")
	}
	ok := s.addEdge(ctx, srcString, destString)
	if !ok {
		return id.NewKernelErrId(id.KernelNotFound)
	}
	if !graph.Acyclic(s.depGraph) {
		pcontext.Errorf(ctx, "acyclic check failed, removing %s->%s",
			src.Short(), dest.Short())
		// remove the edge so no cycles
		s.removeEdge(ctx, s.vertexName[srcString], s.vertexName[destString])
		// no need to check these again for existence, remove edge would not have worked
		srcV := s.vertexName[srcString]
		destV := s.vertexName[destString]
		path, _ := graph.ShortestPath(s.depGraph, destV, srcV)
		buf := &bytes.Buffer{}
		// discover the cycle
		for _, vertex := range path {
			n, bool := reverseMap(s.vertexName, vertex)
			if !bool {
				panic("badly formed graph in doing cycle calculation")
			}
			buf.WriteString(n + "\n")
		}
		pcontext.Debugf(ctx, "cycle:\n%s", buf.String())
		return id.NewKernelErrId(id.KernelDependencyCycle)

	}
	return id.KernelErrIdNoErr
}

func (s *syscallDataImpl) removeEdge(ctx context.Context, v, u int) {
	s.depGraph.Delete(v, u)
}
func (s *syscallDataImpl) addVertex(ctx context.Context, name string) bool {
	prevOrder := s.depGraph.Order()
	newG := graph.New(prevOrder + 1)
	for v := 0; v < prevOrder; v++ {
		s.depGraph.Visit(v, func(w int, _ int64) bool {
			newG.Add(v, w)
			return false
		})
	}
	s.depGraph = newG

	_, ok := s.vertexName[name]
	if ok {
		pcontext.Errorf(ctx, "attempt to add vertext %s ignored, vertex already in graph", name)
		return true
	}
	s.vertexName[name] = prevOrder
	return true
}

func (s *syscallDataImpl) addEdge(ctx context.Context, src, dest string) bool {
	srcV, srcOk := s.vertexName[src]
	destV, destOk := s.vertexName[dest]
	if !srcOk || !destOk {
		text := "neither are graph vertices"
		if srcOk && !destOk {
			text = "destination not a graph vertex"
		}
		if !srcOk && destOk {
			text = "source not a graph vertex"
		}
		pcontext.Errorf(ctx, "attempt to create edge (%s,%s) rejected, %s", src, dest, text)
		return false
	}
	s.depGraph.Add(srcV, destV)
	return true
}

// Run blocks the caller on a particular service being ready to run.  Note that
// function does not assert the lock.
func (s *syscallDataImpl) Run(ctx context.Context, sid id.ServiceId) bool {
	service := s.ServiceById(ctx, sid)
	if service == nil {
		return false
	}
	return service.Run(ctx)
}

// PathExists returns true if there is a path from src to dest
// following dependency edges.  Not that this implies that a
// service source requiring foo, and service foo requiring bar, will
// return true for a call for PathExists(source,bar).
// This means that carefully crafted require's that know the
// depgraph of other services will work, but seems unnecessary.
func (s *syscallDataImpl) PathExists(ctx context.Context, src, dest string) bool {

	// lock for the graph
	s.lock.Lock()
	defer s.lock.Unlock()
	data := DFS(s.depGraph)
	destV, ok := s.vertexName[dest]
	if !ok {
		pcontext.Errorf(ctx, "unable find vertex, %s, can't check for path existence", dest)
	}
	curr := destV
	for curr != -1 {
		cand, ok := reverseMap(s.vertexName, curr)
		if !ok {
			panic("badly formed dependency graph, can't find " + fmt.Sprint(curr))
		}
		// msg := ""
		// if cand == src {
		// 	msg = ":WINNER!"
		// }
		//pcontext.Debugf(ctx, "current is %s%s", cand, msg)
		if cand == src {
			return true
		}
		curr = data.Prev[curr]
	}
	pcontext.Errorf(ctx, "import but no require: %s -> %s", src, dest)
	return false
}

func reverseMap(dep map[string]int, i int) (string, bool) {
	for k, v := range dep {
		if v == i {
			return k, true
		}
	}
	return "", false
}
