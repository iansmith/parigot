package main

import (
	"sync"
	"time"

	"github.com/dominikbraun/graph"
	"github.com/iansmith/parigot/apishared/id"
)

var depData = newSyscallDataImpl()

var runWaitTimeout = time.Duration(10) * time.Second

type Service interface {
	Id() id.ServiceId
	Name() string
	Package() string
	Short() string
	String() string
	RunRequested() bool
	Exported() bool
	Run() bool
}

type SyscallData interface {
	//ServiceByName looks up a service and returns it based on the
	//values package_ and name.  If this returns nil, the service could
	//not be found.
	ServiceByName(package_, name string) Service
	//ServiceById looks up a service and returns it based on the
	//value sid.  If this returns nil, the service could
	//not be found.
	ServiceById(id.ServiceId) Service
	//ServiceByIdString looks up a service based on the printed representation
	//of the service id.  If the service cannot be found ServiceByIdString
	//returns nil.
	ServiceByIdString(string) Service
	// SetService puts a service into SyscallData.  This should only be
	// called once for each package_ and name pair. It returns the
	// ServiceId for the service named, creating a new one if necessary.
	// If the bool result is false, then the pair already existed and
	// we made no changes to it.
	SetService(package_, name string) (Service, bool)
	// Export finds a service by the given sid and then marks that
	// service as being exported. This function returns nil if
	// there is no such service.
	Export(svc id.ServiceId) Service
	// Import introduces a dendency between the sourge and dest
	// services. Thus,  dest must be running before source can run.
	// This function returns false if the services cannot be found
	// or the introduction of this edge would produce a cyle.
	Import(src, dest id.ServiceId) bool
}

//
// serviceImpl
//

type serviceImpl struct {
	pkg, name string
	id        id.ServiceId
	runReady  bool
	exported  bool
	runCh     chan struct{}
	lock      *sync.Mutex
}

func newServiceImpl(pkg, name string, sid id.ServiceId) *serviceImpl {
	result := &serviceImpl{
		pkg:      pkg,
		name:     name,
		id:       sid,
		runReady: false,
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

func (s *serviceImpl) Run() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.runReady = true
	return s.waitToRun()
}

func (s *serviceImpl) canRun() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Exported() && s.RunRequested()
}

// waitToRun waits until the timeout expires or until it receives a wake
// up call and a check for the ability to run successfully is made. It returns
// false if it is returning because of a timeout. Note that this function
// does not lock so that other things can proceed concurrently.
func (s *serviceImpl) waitToRun() bool {
	if s.canRun() {
		return true
	}
	for {
		select {
		case <-s.runCh:
			if s.canRun() {
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

func (s *syscallDataImpl) SetService(package_, name string) (Service, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	svc := s.serviceByNameNoLock(package_, name)
	if svc != nil {
		return svc, false
	}
	svcId := id.NewServiceId()
	nmap, ok := s.packageNameToServiceNameMap[package_]
	if !ok {
		nmap = make(map[string]Service)
		s.packageNameToServiceNameMap[package_] = nmap
	}
	result := newServiceImpl(package_, name, svcId)
	nmap[name] = result
	s.sidStringToService[result.String()] = result
	s.depGraph.AddVertex(result.String())
	return result, true

}

func (s *syscallDataImpl) ServiceByName(package_, name string) Service {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.serviceByNameNoLock(package_, name)
}
func (s *syscallDataImpl) serviceByNameNoLock(package_, name string) Service {

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

func (s *syscallDataImpl) ServiceById(sid id.ServiceId) Service {
	return s.ServiceByIdString(sid.String())
}
func (s *syscallDataImpl) serviceByIdStringNoLock(sid string) Service {
	svc, ok := s.sidStringToService[sid]
	if !ok {
		return nil
	}
	return svc
}

func (s *syscallDataImpl) ServiceByIdString(sid string) Service {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.serviceByIdStringNoLock(sid)
}

func (s *syscallDataImpl) Export(svcId id.ServiceId) Service {
	s.lock.Lock()
	defer s.lock.Unlock()

	svc := s.serviceByIdStringNoLock(svcId.String())
	if svc == nil {
		return nil
	}
	svc.(*serviceImpl).export()
	s.notifyNodedBehindNoLock(svc.String())
	return svc
}

func (s *syscallDataImpl) notifyNodedBehindNoLock(svcid string) {
	topo, _ := graph.TopologicalSort(s.depGraph)
	for _, str := range topo {
		svc := s.ServiceByIdString(str)
		if svcid == svc.String() {
			break
		}
		svc.(*serviceImpl).wakeUp()
	}
}

func (s *syscallDataImpl) Import(src, dest id.ServiceId) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	serviceSource := s.serviceByIdStringNoLock(src.String())
	if serviceSource == nil {
		return false
	}
	serviceDest := s.serviceByIdStringNoLock(dest.String())
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
