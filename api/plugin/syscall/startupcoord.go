package syscall

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	"github.com/yourbasic/graph"
)

//
// rawLocal
//

// The default lock discipline for this type is that you should call a method
// on this type when the lock is unlocked.  The methods you can call when you
// DO have the lock are demarcated by the NoLock suffix. Any function that
// does assert the lock should make sure to release it before returning.

type rawLocal struct {
	sidStringToService          map[string]Service
	packageNameToServiceNameMap map[string]map[string]Service
	depGraph                    *graph.Mutable
	vertexName                  map[string]int
	lock                        *sync.Mutex
}

func newSyscallDataImpl() *rawLocal {
	g := graph.New(0)
	impl := &rawLocal{
		sidStringToService:          make(map[string]Service),
		packageNameToServiceNameMap: make(map[string]map[string]Service),
		depGraph:                    g,
		lock:                        new(sync.Mutex),
		vertexName:                  make(map[string]int),
	}
	_ = SyscallData(impl)
	return impl
}

// SetService is both the "check if the service exists" and create the service
// in one function.  If the service is found as already existing, it returns the
// service and false.  If the service is not found it is created and the service
// is returned and the value false.
func (s *rawLocal) SetService(ctx context.Context, package_, name string, client bool) (Service, bool) {

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

	result := newStartupService(package_, name, svcId, s, client)

	nmap[name] = result
	if !s.addVertex(ctx, result.String()) {
		return nil, false
	}
	s.sidStringToService[result.String()] = result
	if result != nil {
		if result.id.IsEmptyValue() || result.id.IsZeroValue() {
			pcontext.Errorf(ctx, "Service Id error, bad id returned from syscall data")
		}
	} else {
		pcontext.Errorf(ctx, "result of set service is nil?")
	}
	return result, true

}

// ServiceByName takes the full name of a service as a package and a service name
// and returns the Service that represents that name.  If the given package and
// name cannot be found it returns nil.
func (s *rawLocal) ServiceByName(ctx context.Context, package_, name string) Service {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.serviceByNameNoLock(ctx, package_, name)
}

// serviceByNameNoLock returns the startupService that represents the service
// id given.  The service given is the String() value of the service id.  This
// function must be called with the lock held, since it does not lock itself.
func (s *rawLocal) serviceByNameNoLock(ctx context.Context, package_, name string) Service {
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

// ServiceById is just a convenience wrapper for ServiceByIdString for folks that
// have the service id they want to convert to a startupService.
func (s *rawLocal) ServiceById(ctx context.Context, sid id.ServiceId) Service {
	return s.ServiceByIdString(ctx, sid.String())
}
func (s *rawLocal) serviceByIdStringNoLock(ctx context.Context, sid string) Service {
	svc, ok := s.sidStringToService[sid]
	if !ok {
		return nil
	}
	return svc
}

func (s *rawLocal) ServiceByIdString(ctx context.Context, sid string) Service {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.serviceByIdStringNoLock(ctx, sid)
}

// Export causes the service given (in the id form) to marked as exported.
// It returns nil if the service is not found, otherwise it returns the startupService
// that represents the given id.
func (s *rawLocal) Export(ctx context.Context, svcId id.ServiceId) Service {
	svc := s.ServiceByIdString(ctx, svcId.String())
	if svc == nil {
		return nil
	}

	svc.Export()

	return svc
}

// Import adds an edge in the dependency graph between src and dest.
// It returns a kerr if either the source or dest cannot be found; it
// returns a kerr if the new edge would create a cycle.
func (s *rawLocal) Import(ctx context.Context, src, dest id.ServiceId) syscall.KernelErr {
	serviceSource := s.ServiceByIdString(ctx, src.String())
	if serviceSource == nil {
		sid, ok := s.SetService(ctx, serviceSource.Package(), serviceSource.Name(), false)
		if ok {
			pcontext.Infof(ctx, "startup coordinator: created service %s%s because of import", sid.Name(), sid.Short())
		}
	}
	serviceDest := s.ServiceByIdString(ctx, dest.String())
	if serviceDest == nil {
		return syscall.KernelErr_NotFound
	}
	srcString := serviceSource.String()
	destString := serviceDest.String()

	// graph does not lock itself
	s.lock.Lock()
	defer s.lock.Unlock()
	if !graph.Acyclic(s.depGraph) {
		panic("graph is already cyclic, some previous edge was added without checking")
	}
	ok := s.addEdge(ctx, srcString, destString)
	if !ok {
		return syscall.KernelErr_NotFound
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
		return syscall.KernelErr_DependencyCycle

	}
	return syscall.KernelErr_NoError
}

// removeEdge is really only useful when you have introduced an edge that
// creates a cycle.
func (s *rawLocal) removeEdge(ctx context.Context, v, u int) {
	s.depGraph.Delete(v, u)
}
func (s *rawLocal) addVertex(ctx context.Context, name string) bool {
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

// addEdge adds an edge to the dependency graph from src to dest.  This edge
// represents the idea that src cannot start until dest is started.
func (s *rawLocal) addEdge(ctx context.Context, src, dest string) bool {
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

// Launch blocks the caller on a particular service being ready to run.  Note that
// this function does not assert the lock.
func (s *rawLocal) Launch(ctx context.Context, sid id.ServiceId) syscall.KernelErr {
	service := s.ServiceById(ctx, sid)
	if service == nil {
		return syscall.KernelErr_NotFound
	}

	kerr := service.Run(ctx)
	if kerr != syscall.KernelErr_NoError {
		return kerr
	}

	s.notifyIncomingNeighbors(ctx, service)

	return syscall.KernelErr_NoError
}

// PathExists returns true if there is a path from src to dest
// following dependency edges.  Not that this implies that a
// service source requiring foo, and service foo requiring bar, will
// return true for a call for PathExists(source,bar).
// This means that carefully crafted require's that know the
// depgraph of other services will work, but seems wildly unnecessary.
func (s *rawLocal) PathExists(ctx context.Context, src, dest string) bool {

	// lock for the graph
	s.lock.Lock()
	defer s.lock.Unlock()
	srcV := s.vertexName[src]
	destV := s.vertexName[dest]

	if !s.searchEdges(srcV, destV) {
		pcontext.Errorf(ctx, "XXXlocate called, but no require given: %s -> %s", src, dest)
		return false
	}
	return true
}

//
// Graph routines
//

// reverseMap is kinda sucky.  This walks the elements of the dependency
// map looking for the name of the vertex that has the vertex number i.
// The number of dependency edges should be small so keeping another map seemed
// like overkill.
func reverseMap(dep map[string]int, i int) (string, bool) {
	for k, v := range dep {
		if v == i {
			return k, true
		}
	}
	return "", false
}

// inbound edges does not touch the lock.
func inboundEdges(g graph.Iterator, target int) []int {
	result := []int{}
	max := g.Order()
	for i := 0; i < max; i++ {
		if i == target {
			continue
		}
		g.Visit(i, func(w int, _ int64) (skip bool) {
			if w == target {
				result = append(result, i)
			}
			return
		})
	}
	return result
}

// This function assumes that callers are HOLDING the lock.
func (s *rawLocal) mustVertexNumToService(ctx context.Context, v int) Service {
	n, ok := reverseMap(s.vertexName, v)
	if !ok {
		fmt.Printf("xxxx looking for=>%d search failed, %+v\n", v, s.vertexName)
		panic("unable to find vertex name for index")
	}
	svc := s.serviceByIdStringNoLock(ctx, n)
	if svc == nil {
		panic("unable to find service by Id string")
	}
	return svc
}

// tricky: we use the lock in this function because we
// call a function that uses the lock.
func (s *rawLocal) mustServiceToVertexNum(ctx context.Context, svc Service) int {
	str := svc.String()
	i, ok := s.vertexName[str]
	if !ok {
		panic("unable to find vertex number for service id string")
	}
	return i
}

// This function is recursive and assumes the lock is NOT
// held when it is entered.
func (s *rawLocal) dependenciesStarted(ctx context.Context, v int) []int {
	fail := []int{}
	s.depGraph.Visit(v, func(w int, c int64) (skip bool) {
		cand := s.mustVertexNumToService(ctx, w)
		if !cand.Started() {
			fail = append(fail, w)
		}
		fail = append(fail, s.dependenciesStarted(ctx, w)...)
		return
	})
	return fail
}

// dfsDeps traces through the dependency graph depth first
// and looks for all reachable nodes and makes sure they are
// all started. If any are not started, it returns false.
// This is called from the startupService and it assumes
// the downstream functions are not touching the lock.
func (s *rawLocal) dfsDeps(ctx context.Context, sid Service) bool {

	i := s.mustServiceToVertexNum(ctx, sid)

	s.lock.Lock()
	defer s.lock.Unlock()

	fail := s.dependenciesStarted(ctx, i)
	return len(fail) == 0
}

// notifyIncomingNeighbors sends the wake up signal to
// to any node that has a dependecy on the node given. This is
// called when have successfully started up service sid.
func (s *rawLocal) notifyIncomingNeighbors(ctx context.Context, sid Service) {
	i := s.mustServiceToVertexNum(ctx, sid)

	s.lock.Lock()
	in := inboundEdges(s.depGraph, i)
	s.lock.Unlock()

	for _, c := range in {
		// mustVertex below asserts the lock
		svc := s.mustVertexNumToService(ctx, c)
		svc.WakeUp()
	}
}

// searchEdges looks for a path from source to dest.  It returns true
// if there such a path, otherwise false.
func (s *rawLocal) searchEdges(source, dest int) bool {
	found := false
	s.depGraph.Visit(source, func(w int, c int64) (skip bool) {
		if w == dest {
			found = true
			skip = true
			return
		}
		found = found || s.searchEdges(w, dest)
		return
	})
	return found

}
