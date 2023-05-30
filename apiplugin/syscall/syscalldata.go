package main

import (
	"context"
	"sync"
	"time"

	"github.com/dominikbraun/graph"
	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"
)

var depData = newSyscallDataImpl()

var runWaitTimeout = time.Duration(10) * time.Second

//
// serviceImpl
//

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

func newServiceImpl(pkg, name string, sid id.ServiceId, parent *syscallDataImpl) *serviceImpl {
	result := &serviceImpl{
		pkg:      pkg,
		name:     name,
		id:       sid,
		runReady: false,
		parent:   parent,
		lock:     new(sync.Mutex),
		runCh:    make(chan struct{}),
	}
	_ = Service(result)
	return result
}

// wakeUp causes a send on the servicImpl's runCh and thus check to see
// if it can run now.
func (s *serviceImpl) wakeUp() {
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
	s.lock.Lock()
	defer s.lock.Unlock()
	s.runReady = true
	return s.waitToRun(ctx)
}

func (s *serviceImpl) canRun(ctx context.Context) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Exported() && s.RunRequested() && s.parent.checkNodesBehindForReadyNoLock(ctx, s.String())
}

// waitToRun waits until the timeout expires or until it receives a wake
// up call and a check for the ability to run successfully is made. It returns
// false if it is returning because of a timeout. Note that this function
// does not lock so that other things can proceed concurrently.
func (s *serviceImpl) waitToRun(ctx context.Context) bool {
	if s.canRun(ctx) {
		return true
	}
	for {
		select {
		case <-s.runCh:
			if s.canRun(ctx) {
				s.started = true
				return true
			}
		case <-time.After(1 * time.Minute):
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

type syscallDataImpl struct {
	sidStringToService          map[string]Service
	packageNameToServiceNameMap map[string]map[string]Service
	depGraph                    graph.Graph[string, string]
	lock                        *sync.Mutex
}

func newSyscallDataImpl() *syscallDataImpl {
	g := graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())
	impl := &syscallDataImpl{
		sidStringToService:          make(map[string]Service),
		packageNameToServiceNameMap: make(map[string]map[string]Service),
		depGraph:                    g,
		lock:                        new(sync.Mutex),
	}
	_ = SyscallData(impl)
	return impl
}

func (s *syscallDataImpl) SetService(ctx context.Context, package_, name string) (Service, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	svc := s.serviceByNameNoLock(ctx, package_, name)
	if svc != nil {
		return svc, false
	}
	svcId := id.NewServiceId()
	nmap, ok := s.packageNameToServiceNameMap[package_]
	if !ok {
		nmap = make(map[string]Service)
		s.packageNameToServiceNameMap[package_] = nmap
	}
	result := newServiceImpl(package_, name, svcId, s)

	nmap[name] = result
	s.sidStringToService[result.String()] = result
	s.depGraph.AddVertex(result.String())
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
	s.lock.Lock()
	defer s.lock.Unlock()

	svc := s.serviceByIdStringNoLock(ctx, svcId.String())
	if svc == nil {
		return nil
	}
	svc.(*serviceImpl).export()
	pcontext.Infof(ctx, "export completed, service %s", svc.Short())
	if svc.(*serviceImpl).canRun(ctx) {

	}
	return svc
}

// topoSort does a topological sort of all the nodes we currently know about.
// This function asserts the lock to ensure that while the topo algorihtm is
// running no part of the graph is disturbed.  It returns what I HOPE is a copy
// of the content of the graph vertices.
func (s *syscallDataImpl) topoSort() []string {
	s.lock.Lock()
	defer s.lock.Unlock()

	topo, _ := graph.TopologicalSort(s.depGraph)
	return topo
}

// notifyNodesBehindForReadyNoLock walks the topologically ordered vertices, looking for any nodes
// that are predecessors of the given node and notifying them to check their status for running.
// This returns false only when the service cannot be found.
func (s *syscallDataImpl) notifyNodesBehindForReadyNoLock(ctx context.Context, svcid string) bool {
	topo := s.topoSort()                        // locks
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
	return true
}

// checkNodesBehindForReadyNoLock walks the topologically ordered vertices, looking for any nodes
// that are predecessors of the given node and testing to see if they are started. If not
// we return false.  If all the predecessorys are started, then we return true.  If the
// serviceId cannot be found, we return true.
func (s *syscallDataImpl) checkNodesBehindForReadyNoLock(ctx context.Context, svcid string) bool {

	topo := s.topoSort()                        // locks
	if s.ServiceByIdString(ctx, svcid) == nil { //locks
		return true
	}
	success := true
	for _, str := range topo {
		if str == svcid {
			pcontext.Debugf(ctx, "%s is now ready to run", svcid)
			return success
		}
		svc := s.ServiceByIdString(ctx, str) //locks
		if svc == nil {
			panic("unable to walk the dep graph looking for predeessors, a predecessor could not be found: " + str)
		}
		if svc.Exported() && svc.RunRequested() && svc.Started() {
			continue
		}
		return false
	}
	panic("did not find id " + svcid + " in the list of vertices")
}

// notifyNodesInFrontNoLock walks the topologically ordered vertices, looking for any nodes
// that have satisfied dependencies and are net yet running. It returns false only if the svcid
// cannot be found.
func (s *syscallDataImpl) checkNodesInFrontNoLock(ctx context.Context, svcid string) bool {
	topo := s.topoSort()
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

func (s *syscallDataImpl) Import(ctx context.Context, src, dest id.ServiceId) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	serviceSource := s.serviceByIdStringNoLock(ctx, src.String())
	if serviceSource == nil {
		return false
	}
	serviceDest := s.serviceByIdStringNoLock(ctx, dest.String())
	if serviceDest == nil {
		return false
	}
	err := s.depGraph.AddEdge(serviceSource.String(), serviceDest.String())
	if err == nil {
		return true
	}
	switch err {
	case graph.ErrEdgeCreatesCycle:
		return false
	case graph.ErrEdgeAlreadyExists:
		return true
	case graph.ErrEdgeNotFound:
		panic("internal error in dependency graph construction")
	default:
		panic("unexpected graph error in dependency graph:" + err.Error())
	}
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
