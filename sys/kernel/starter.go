package kernel

import (
	"fmt"
	"math/rand"
	"strings"

	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/syscall/v1"
)

// starter manages the data structures are getting the services to come up in
// the correct order.
type starter struct {
	serviceIsLaunched    map[string]struct{}
	serviceToWaiting     map[string][]FQName
	pkgToServiceImpl     map[string]map[string][]hostServiceBinding
	serviceToExports     map[string][]FQName
	serviceToFulfillment map[string]map[string]map[string]hostServiceBinding
	stringToService      map[string]id.ServiceId
	serviceToMethMap     map[string]map[string]id.MethodId

	waitList    []launchData
	pendingList []launchCompleteBundle
	runList     []id.ServiceId

	//matchImpl callMatcher
}

// launchData is the information needed about a service that is in
// the waiting state.  We need the information here so we can
// complete the launch request when th
type launchData struct {
	sid id.ServiceId
	cid id.CallId
	hid id.HostId
}

func (l *launchData) String() string {
	return fmt.Sprintf("LauchData{%s,%s,%s}", l.sid.Short(), l.cid.Short(), l.hid.Short())
}

var _ Starter = &starter{}

// NewStarter returns an implementation of the starter that assumes
// we have global knowlege.
func NewStarter() *starter {
	return &starter{
		serviceIsLaunched:    map[string]struct{}{},
		serviceToWaiting:     map[string][]FQName{},
		pkgToServiceImpl:     map[string]map[string][]hostServiceBinding{},
		serviceToExports:     map[string][]FQName{},
		serviceToFulfillment: map[string]map[string]map[string]hostServiceBinding{},
		stringToService:      map[string]id.ServiceId{},
		serviceToMethMap:     map[string]map[string]id.MethodId{},
		waitList:             []launchData{},
		pendingList:          []launchCompleteBundle{},
		runList:              []id.ServiceId{},
	}
}

func (s *starter) Export(hid id.HostId, sid id.ServiceId, fqn FQName) syscall.KernelErr {
	// not registered
	if _, ok := s.stringToService[sid.String()]; !ok {
		klog.Errorf("ailed to find string %s", sid.String())
		return syscall.KernelErr_BadId
	}

	pkg := fqn.Pkg
	name := fqn.Name
	pkg2map, ok := s.pkgToServiceImpl[pkg]
	if !ok {
		pkg2map = make(map[string][]hostServiceBinding)
		s.pkgToServiceImpl[pkg] = pkg2map
	}
	allBind, ok := pkg2map[name]
	if !ok {
		allBind = []hostServiceBinding{}
		pkg2map[name] = allBind
	}
	allBind = append(allBind, hostServiceBinding{
		service: sid,
		host:    hid,
	})
	pkg2map[name] = allBind
	s.pkgToServiceImpl[pkg] = pkg2map

	s.findRunnable()
	return syscall.KernelErr_NoError
}

func (s *starter) Launch(sid id.ServiceId, cid id.CallId, hid id.HostId, mid id.MethodId) syscall.KernelErr {
	if hid.IsZeroOrEmptyValue() || cid.IsZeroOrEmptyValue() || sid.IsEmptyValue() {
		klog.Errorf("launch failed because of bad id (%s,%s,%s)",
			sid.Short(), cid.Short(), hid.Short())
		return syscall.KernelErr_BadId
	}
	if !mid.Equal(apishared.LaunchMethod) {
		klog.Errorf("launch failed because method id doesn't match %s,%s",
			mid.Short(), apishared.LaunchMethod.Short())
		return syscall.KernelErr_BadId
	}
	ld := launchData{
		sid: sid,
		cid: cid,
		hid: hid,
	}

	s.waitList = append(s.waitList, ld)
	s.serviceIsLaunched[sid.String()] = struct{}{}
	s.findRunnable()

	return syscall.KernelErr_NoError
}

// Dead code?
func (s *starter) Register(hid id.HostId, sid id.ServiceId, debugName string) syscall.KernelErr {
	s.stringToService[sid.String()] = sid
	return syscall.KernelErr_NoError
}

func (s *starter) Ready() (launchCompleteBundle, int) {
	if len(s.pendingList) == 0 {
		lcb := launchCompleteBundle{
			hid:      id.HostIdZeroValue(),
			cid:      id.CallIdZeroValue(),
			sid:      id.ServiceIdZeroValue(),
			hasCycle: false,
		}
		return lcb, len(s.waitList)
	}
	zero := s.pendingList[0]
	if len(s.pendingList) == 1 {
		s.pendingList = nil
	} else {
		s.pendingList = s.pendingList[1:]
	}
	s.runList = append(s.runList, zero.sid)
	return zero, len(s.waitList)
}

// hostServiceBinding creates a connection between a given
// service id and the host that it lives on.
type hostServiceBinding struct {
	service id.ServiceId
	host    id.HostId
}

// requirementsMet checks to see if there are services that
// meet the requirements requested *and* these services are
// in the run list.  The return value is nil in the case where
// the sid has not had it its requirements met.  Otherwise,
// the return value is a map of maps whose keys are the package
// name and service name of a requirement and the final value
// is the hostServiceBinding that fulfilled the requirement.
func (s *starter) requirementsMet(sid id.ServiceId) map[string]map[string]hostServiceBinding {
	//if you are not launched, then your requirements cannot be met
	_, ok := s.serviceIsLaunched[sid.String()]
	if !ok {
		klog.Errorf("can't be ready to run (%s) because not launched yet", sid.String())
		return nil
	}

	neededList := s.serviceToWaiting[sid.String()]
	fulfilled := make(map[string]map[string]hostServiceBinding)
	for _, need := range neededList {
		smap, ok := s.pkgToServiceImpl[need.Pkg]
		if !ok {
			return nil
		}
		implList, ok := smap[need.Name]
		if !ok {
			return nil
		}
		running := false
		var winner hostServiceBinding
		for _, impl := range implList {
			if s.isRunningOrPending(impl.service) {
				winner = impl
				running = true
				break
			} else {
				running = false
			}
		}
		if !running {
			return nil
		}
		sv, ok := fulfilled[need.Pkg]
		if !ok {
			sv = make(map[string]hostServiceBinding)
			fulfilled[need.Pkg] = sv
		}
		_, ok = sv[need.Name]
		if !ok {
			sv[need.Name] = winner
		} else {
			klog.Errorf("unexpected already existing fulfilled value for %s.%s", need.Pkg, need.Name)
		}
	}
	return fulfilled
}

func (s *starter) detectCycle(hid id.HostId, sid id.ServiceId) bool {
	export := s.serviceToExports[sid.String()]
	if len(export) == 0 { // no exports
		return false
	}
	for svc, fulfilled := range s.serviceToFulfillment {
		for pkg, nMap := range fulfilled {
			for _, ex := range export {
				if ex.Pkg == pkg {
					for name, candidate := range nMap {
						if name == ex.Name {
							if sid.Equal(candidate.service) &&
								hid.Equal(candidate.host) {
								klog.Errorf("cycle detected between %s and %s, because %s imports %s.%s and %s exports that service",
									sid.Short(), svc, svc, pkg, name, sid.Short())
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

// isRunningOrPending checks to see if a particular service is either on the
// running list or will be on it soon via the pending list.
func (s *starter) isRunningOrPending(sid id.ServiceId) bool {
	for _, r := range s.runList {
		if r.Equal(sid) {
			return true
		}
	}
	for _, l := range s.pendingList {
		if l.sid.Equal(sid) {
			return true
		}
	}
	return false
}

// findRunnable walks the waiting list and if the service on the
// waiting list is ready to run, it moves the service to the running
// list and notifies the channel that is waiting on the response.
// This function runs until there are no more changes.
func (s *starter) findRunnable() {
	change := true
	var result []launchData
	if len(s.waitList) == 0 {
		return
	}
	for change {
		change = false
		result = []launchData{}
		for _, wait := range s.waitList {
			fulfillment := s.requirementsMet(wait.sid)
			if fulfillment != nil {
				change = true
				s.serviceToFulfillment[wait.sid.String()] = fulfillment
				s.moveToPending(wait)
			} else {
				result = append(result, wait)
			}
		}
		s.waitList = result
	}
}

// moveToPending takes a service from the from the waiting list
// to the pending list. The pending list is where services go
// that are waiting to run based on the kernel asking for a
// service via Ready().  Once they are returned by Ready()
// they are put on the running list.
func (s *starter) moveToPending(launch launchData) syscall.KernelErr {
	pp := launchCompleteBundle{
		sid:      launch.sid,
		hasCycle: false,
		hid:      launch.hid,
		cid:      launch.cid,
	}
	if s.detectCycle(launch.hid, launch.sid) {
		klog.Errorf("started service %s, but cycle detected in depndencies",
			launch.sid.Short())
		pp.hasCycle = true
	}
	s.pendingList = append(s.pendingList, pp)
	return syscall.KernelErr_NoError
}

func (s *starter) Bind(_ id.HostId, sid id.ServiceId, mid id.MethodId, methodName string) syscall.KernelErr {
	sMap, ok := s.serviceToMethMap[sid.String()]
	if !ok {
		sMap = make(map[string]id.MethodId)
		s.serviceToMethMap[sid.String()] = sMap

	}
	sMap[methodName] = mid
	return syscall.KernelErr_NoError
}

// locate is a foundational call in parigot. It converts a fully qualified
// string (like foo.v1.Foo) into a reference to service and host.  This call
// fails if nobody has exported the service requested.  Since you should always
// call Launch() before Locate() there is no issue with timing.  Launch blocks
// until all dependencies are met.
func (s *starter) Locate(req *syscall.LocateRequest, resp *syscall.LocateResponse) syscall.KernelErr {
	sid := id.ServiceIdEmptyValue()
	if req.GetCalledBy() != nil {
		sid = id.UnmarshalServiceId(req.GetCalledBy())
	}
	pkg := req.GetPackageName()
	name := req.GetServiceName()
	sMap, ok := s.pkgToServiceImpl[pkg]

	if !ok {
		klog.Errorf("failed to find service that was requested in Locate (0): %s.%s", strings.ToUpper(pkg), name)
		return syscall.KernelErr_NotFound
	}

	if !sid.IsEmptyValue() && !s.checkAlreadyRequired(sid, pkg, name) {
		klog.Errorf("service %s did not require service %s.%s but imports it", sid.Short(),
			pkg, name)
		return syscall.KernelErr_NotRequired
	}

	target, ok := sMap[name]
	if !ok {
		klog.Errorf("failed to find service that was requested in Locate (1): %s.%s", pkg, strings.ToUpper(name))
		return syscall.KernelErr_NotFound
	}
	if len(target) == 0 {
		klog.Errorf("failed to find service that was requested in Locate (2): %s.%s", strings.ToUpper(pkg), strings.ToUpper(name))
		return syscall.KernelErr_NotFound

	}
	chosen := target[0]
	if len(target) > 1 {
		n := rand.Intn(len(target))
		chosen = target[n]
	}

	resp.ServiceId = chosen.service.Marshal()
	resp.HostId = chosen.host.Marshal()
	meth, ok := s.serviceToMethMap[chosen.service.String()]
	if !ok {
		// no methods?
		klog.Errorf("unable to find any methods associated with service %s (%s.%s)", chosen.service.Short(), pkg, name)
	}
	resp.Binding = []*syscall.MethodBinding{}
	for name, meth := range meth {
		mb := &syscall.MethodBinding{
			MethodName: name,
			MethodId:   meth.Marshal(),
		}
		resp.Binding = append(resp.Binding, mb)
	}
	return syscall.KernelErr_NoError
}

// checkAlreadyRequired returns true when the waiting list for the service
// provided includes pkg.name.
func (s *starter) checkAlreadyRequired(sid id.ServiceId, pkg, name string) bool {
	waiting := s.serviceToWaiting[sid.String()]
	found := false
	for _, wait := range waiting {
		if wait.Name == name && wait.Pkg == pkg {
			found = true
			break
		}
	}
	return found
}

// require
func (s *starter) Require(req *syscall.RequireRequest, resp *syscall.RequireResponse) syscall.KernelErr {
	src := id.UnmarshalServiceId(req.GetSource())
	dest := make([]FQName, len(req.GetDest()))
	for i, d := range req.GetDest() {
		curr := FQName{
			Pkg:  d.PackagePath,
			Name: d.Service,
		}
		dest[i] = curr
	}
	wait, ok := s.serviceToWaiting[src.String()]
	if !ok {
		wait = []FQName{}
	}
	for i := 0; i < len(dest); i++ {
		wait = append(wait, dest[i])
	}
	s.serviceToWaiting[src.String()] = wait

	return syscall.KernelErr_NoError

}
