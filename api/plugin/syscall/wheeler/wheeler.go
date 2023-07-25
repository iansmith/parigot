package wheeler

import (
	"context"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type serviceState int8

// Wheeler is an opaque interface that talk to over the provided
// channel.  It implements the system calls for parigot.
type Wheeler interface {
	In() chan InProtoPair
}

const (
	maxRunListSize  = 256
	maxWaitListSize = 256
)

// the only wheeler
var _wheeler *wheeler

// InstallWheeler should be called one time to
// install the wheeler with the given exitCh.  It is assumed
// that the caller will be monitoring the exit channel.
// The exit channel receives information about exits requested
// or required from the running program.  If the exit code
// is 0 to 192, this a requested exit with that code.  If
// the value is > 192 then this is an abort.  Values
// from 193 to 255 are reserved and the value 255 indicates
// a trapped error (trapped via recover), likely a panic.
//
// Values with the 256 bit set (bit 8) are exits caused
// by exit codes taken from the syscall.KernelErr set and the
// value of the error is ored into the low order 4 bits.
//
// The context provided is saved and used for all calls on the
// wheeler.
func InstallWheeler(ctx context.Context, exitCh chan int32) {
	_wheeler = newWheeler(ctx, exitCh)
}

// OutProtoPair is the return type of a message to the wheeler.
// it is sent through the channel given as part of the request.
// If the err != 0, the msg should be ignored.   If the err
// is 0, then the msg will be non-nil.
type OutProtoPair struct {
	A   *anypb.Any
	Err syscall.KernelErr
}

// InProtoPair is a request and the channel to send the error
// or response to.
type InProtoPair struct {
	Msg proto.Message
	Ch  chan OutProtoPair
}

// hostServiceBinding creates a connection between a given
// service id and the host that it lives on.
type hostServiceBinding struct {
	service id.ServiceId
	host    id.HostId
}

// fqName is a fully qualified name of a service, analagous to
// syscall.FullyQualifiedName.
type fqName struct {
	pkg, name string
}

// launchData is the information needed about a service that is in
// the waiting state.  We need the information here so we can
// complete the launch request when th
type launchData struct {
	sid  id.ServiceId
	resp *syscall.LaunchResponse // currently empty
	ch   chan OutProtoPair
}

// wheeler is the type that implements system calls.  It actually
// reads a channel and responds to the requests one by one.
// One can think of it as a wheel in that multiple different
// callers are all trying to get info to the wheeler and it
// is spinning around to take each one in turn.
type wheeler struct {
	ch                   chan InProtoPair
	ctx                  context.Context
	exitCh               chan int32
	pkgToServiceImpl     map[string]map[string][]hostServiceBinding
	hostToService        map[string][]id.ServiceId
	serviceToHost        map[string]id.HostId
	serviceToMethMap     map[string]map[string]id.MethodId
	serviceToFQName      map[string]fqName
	pkgToRegistration    map[string]map[string][]hostServiceBinding
	hostToRegistration   map[string][]id.ServiceId
	stringToService      map[string]id.ServiceId
	serviceToWaiting     map[string][]fqName
	serviceToFulfillment map[string]map[string]map[string]hostServiceBinding
	serviceIsLaunched    map[string]struct{}
	waitList             []launchData
	runList              []id.ServiceId
}

// newWheeler returns a Wheeler and should only be called--
// exactly one time--from  InstallWheeler.
func newWheeler(ctx context.Context, exitCh chan int32) *wheeler {
	w := &wheeler{
		exitCh:               exitCh,
		ctx:                  ctx,
		ch:                   make(chan InProtoPair),
		pkgToServiceImpl:     make(map[string]map[string][]hostServiceBinding),
		hostToService:        make(map[string][]id.ServiceId),
		serviceToHost:        make(map[string]id.HostId),
		serviceToMethMap:     make(map[string]map[string]id.MethodId),
		serviceToFQName:      make(map[string]fqName),
		pkgToRegistration:    make(map[string]map[string][]hostServiceBinding),
		hostToRegistration:   make(map[string][]id.ServiceId),
		stringToService:      make(map[string]id.ServiceId),
		serviceToWaiting:     make(map[string][]fqName),
		serviceIsLaunched:    make(map[string]struct{}),
		serviceToFulfillment: make(map[string]map[string]map[string]hostServiceBinding),
		waitList:             []launchData{},
		runList:              []id.ServiceId{},
	}
	go w.Run()
	return w
}

// In is the implementation of the only method in the Wheeler interface.
// It just returns the channel that people can write on.
func In() chan InProtoPair {
	return _wheeler.ch
}

// errorf is a convenience for writing an error the context given to
// us at creation. if the context is nil, this does nothing.
func (w *wheeler) errorf(spec string, rest ...interface{}) {
	if w.ctx == nil {
		return
	}
	pcontext.Errorf(w.ctx, spec, rest...)
	pcontext.Dump(w.ctx)
}

// requirementsMet checks to see if there are services that
// meet the requirements requested *and* these services are
// in the run list.
func (w *wheeler) requirementsMet(sid id.ServiceId) bool {
	//if you are not launched, then your requirements cannot be met
	_, ok := w.serviceIsLaunched[sid.String()]
	if !ok {
		log.Printf("xxx -- can't be ready to run (%s) because not launched yet", sid.String())
		return false
	}

	neededList := w.serviceToWaiting[sid.String()]
	fulfilled := make(map[string]map[string]hostServiceBinding)
	log.Printf("xxxx -- checking requirements (0) %s has %+v", sid.Short(), neededList)
	for _, need := range neededList {
		log.Printf("xxx -- checking requirements (1) %s who needs %s.%s", sid.Short(), need.pkg, need.name)
		smap, ok := w.pkgToServiceImpl[need.pkg]
		if !ok {
			return false
		}
		implList, ok := smap[need.name]
		if !ok {
			return false
		}
		running := false
		var winner hostServiceBinding
		for _, impl := range implList {
			if w.isRunning(impl.service) {
				log.Printf("xxx -- %s has dep [%s,%s] that is running (searching for %s.%s)", sid.Short(), impl.host.Short(), impl.service.Short(), need.pkg, need.name)
				winner = impl
				running = true
				break
			} else {
				log.Printf("xxx -- %s has dep [%s,%s] that is NOT running (searchin for %s.%s)", sid.Short(), impl.host.Short(), impl.service.Short(), need.pkg, need.name)
				running = false
			}
		}
		if !running {
			return false
		}
		s, ok := fulfilled[need.pkg]
		if !ok {
			s = make(map[string]hostServiceBinding)
			fulfilled[need.pkg] = s
		}
		_, ok = s[need.name]
		if !ok {
			s[need.name] = winner
		} else {
			w.errorf("unexpected already existing fulfilled value for %s.%s", need.pkg, need.name)
		}
	}
	// if len(fulfilled) == 0 {
	// 	log.Printf("xxx requirements met, would have nulled fulfilled")
	// 	fulfilled = nil
	// }
	// return fulfilled
	return true
}

// isRunning checks to see if a particular service is on the running list.
func (w *wheeler) isRunning(sid id.ServiceId) bool {
	running := false
	for _, r := range w.runList {
		if r.Equal(sid) {
			running = true
			break
		}
	}
	return running
}

// findRunnable walks the waiting list and if the service on the
// waiting list is ready to run, it moves the service to the running
// list and notifies the channel that is waiting on the response.
// This function runs until there are no more changes.
func (w *wheeler) findRunnable() {
	log.Printf("xxxx --- find Runnable")
	change := true
	var result []launchData
	if len(w.waitList) == 0 {
		return
	}
	for change {
		change = false
		result = []launchData{}
		log.Printf("xxxx --- find Runnable: waitlist size %d", len(w.waitList))
		for _, wait := range w.waitList {
			if w.requirementsMet(wait.sid) {
				change = true
				w.runList = append(w.runList, wait.sid)
				if !w.notifyRun(wait) {
					w.errorf("unable to send response back to client via the launch response channel")
					return
				}
				log.Printf("xxx -- we moved %s to run list", wait.sid.Short())
			} else {
				result = append(result, wait)
			}
		}
		w.waitList = result
		if change {
			log.Printf("xxx -- something is ready now, go through all waiting again")
		}
	}
	log.Printf("xxxx --- find Runnable: finished with waitlist of %d", len(w.waitList))

}

// launch causes the calling service to be put on the waiting list. It
// will check to see if there are any eligible to run services at that point.
// This function return NoError to indicate to the caller that no
// value should be returned to the caller, as the finish of this
// is handled by findRunnable.
func (w *wheeler) launch(req *syscall.LaunchRequest, respChan chan OutProtoPair) (*anypb.Any, syscall.KernelErr) {
	sid := id.UnmarshalServiceId(req.GetServiceId())
	if sid.IsZeroOrEmptyValue() {
		return nil, syscall.KernelErr_BadId
	}
	ld := launchData{
		sid:  sid,
		ch:   respChan,
		resp: &syscall.LaunchResponse{},
	}
	w.waitList = append(w.waitList, ld)
	w.serviceIsLaunched[sid.String()] = struct{}{}
	log.Printf("xxxx ---  we are checking to see if the newly launched service is is runnable (%s)", sid.Short())
	w.findRunnable()
	if w.isRunning(sid) {
		resp := &syscall.LaunchResponse{}
		resp.IsRunning = true
		resp.CallId = id.CallIdZeroValue().Marshal()
		resp.HostId = id.HostIdZeroValue().Marshal()
		return returnResponseOrMarshalError(w, resp)
	}
	resp := &syscall.LaunchResponse{}
	resp.IsRunning = false
	resp.CallId = id.NewCallId().Marshal()
	hid, ok := w.serviceToHost[sid.String()]
	if !ok {
		return nil, syscall.KernelErr_NotFound
	}
	resp.HostId = hid.Marshal()
	return returnResponseOrMarshalError(w, resp)
}

// notifyRun creates the response bundle, an OutProtoPair, and
// sends it to the code that is waiting on the result.  This returns
// false if there was a problem, meaning that the service that is
// being notified should be removed from the run list (since we could
// not tell it that it is now running).
func (w *wheeler) notifyRun(launch launchData) bool {
	pair := OutProtoPair{}

	pair.Err = syscall.KernelErr_NoError
	pair.A = &anypb.Any{}
	result := true

	err := pair.A.MarshalFrom(launch.resp)
	if err != nil {
		w.errorf("unable to start service %s, failed to marshal for launch response: %s",
			launch.sid.Short(), err.Error())
		pair.A = nil
		pair.Err = syscall.KernelErr_MarshalFailed
		result = false
	}
	log.Printf("telling %s that the launch data is %+v, and go for launch", launch.sid.Short(), pair.A)
	// send them their package
	launch.ch <- pair
	return result
}

// export implements the export for all the given types.
// It binds each type to the hostname provided.  It also does
// a check of the runnable services, since calling export on a service
// can change the state of _others_ dependencies.
func (w *wheeler) export(req *syscall.ExportRequest) (*anypb.Any, syscall.KernelErr) {
	sid := id.UnmarshalServiceId(req.ServiceId)
	hid := id.UnmarshalHostId(req.HostId)
	resp := &syscall.ExportResponse{}

	if hid.IsZeroOrEmptyValue() {
		return nil, syscall.KernelErr_BadId
	}
	if sid.IsZeroOrEmptyValue() {
		return nil, syscall.KernelErr_BadId
	}
	// not registered
	if _, ok := w.stringToService[sid.String()]; !ok {
		return nil, syscall.KernelErr_BadId
	}

	w.addHost(hid, sid)

	for _, fqn := range req.GetService() {
		pkg := fqn.GetPackagePath()
		name := fqn.GetService()
		log.Printf("xxx --- export %s: %s.%s", sid.Short(), pkg, name)
		pkg2map, ok := w.pkgToServiceImpl[pkg]
		if !ok {
			pkg2map = make(map[string][]hostServiceBinding)
			w.pkgToServiceImpl[pkg] = pkg2map
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
		w.pkgToServiceImpl[pkg] = pkg2map
		w.serviceToFQName[sid.String()] = fqName{fqn.PackagePath, fqn.Service}
	}
	w.findRunnable()
	return returnResponseOrMarshalError(w, resp)
}

// addHost adds the given hid to our to list of hosts if we have
// not seen it before.  It also adds the given service to the mappings
// between hosts and what services they have and what host a given service
// is on.
func (w *wheeler) addHost(hid id.HostId, sid id.ServiceId) {
	w.serviceToHost[sid.String()] = hid
	allSvc, ok := w.hostToService[hid.String()]
	if !ok {
		allSvc = []id.ServiceId{}
		w.hostToService[hid.String()] = []id.ServiceId{}
	}
	allSvc = append(allSvc, sid)
	w.hostToService[hid.String()] = allSvc
}

// checkHost should be called periodically to validate that the services
// on a host are still running.
func (w *wheeler) checkHost(hid id.HostId, allSvc []id.ServiceId) {
	//ignored for now, should be a status check to see if the
	//alleged services are still there
}

// exit is a call that has a response, but is actually called for
// effect.  wheeler notifies its creator via the channel that was
// provided at creation time.  The caller (receiver) should take
// actions to gracefully shutdown the appropriate services.
func (w *wheeler) exit(req *syscall.ExitRequest) (*anypb.Any, syscall.KernelErr) {
	resp := &syscall.ExitResponse{}

	code := int32(192)
	if req.GetCode() >= 0 && req.GetCode() <= 192 {
		code = req.GetCode()
	}
	w.exitCh <- code

	resp.Code = code
	return returnResponseOrMarshalError(w, resp)
}

// register creates an entry in wheeler data structures that
// are per-service.  It must be called by a service before that
// service can export interfaces or launch.
func (w *wheeler) register(req *syscall.RegisterRequest) (*anypb.Any, syscall.KernelErr) {
	pkg := req.Fqs.GetPackagePath()
	name := req.Fqs.GetService()
	hid := id.UnmarshalHostId(req.GetHostId())

	if hid.IsZeroOrEmptyValue() {
		w.errorf("host Id in register() is zero or empty: %s", hid.Short())
		return nil, syscall.KernelErr_BadId
	}

	sMap, ok := w.pkgToRegistration[pkg]
	if !ok {
		sMap = make(map[string][]hostServiceBinding)
		w.pkgToRegistration[pkg] = sMap
	}
	b, ok := sMap[name]
	if !ok {
		b = []hostServiceBinding{}
		sMap[name] = b
	}
	sid := id.NewServiceId()
	existedPreviously := false
	//check it's not already there
	for _, bind := range b {
		if bind.host.Equal(hid) {
			sid = bind.service
			existedPreviously = true
			break
		}
	}
	if !existedPreviously {
		b = append(b, hostServiceBinding{host: hid, service: sid})
		sMap[name] = b
		w.stringToService[sid.String()] = sid
	}
	resp := &syscall.RegisterResponse{}
	resp.ExistedPreviously = existedPreviously
	resp.Id = sid.Marshal()

	return returnResponseOrMarshalError(w, resp)
}

// bindMethod creates a mapping in the tables about what methods
// a given (concrete) service has.  This registration should be done
func (w *wheeler) bindMethod(req *syscall.BindMethodRequest) (*anypb.Any, syscall.KernelErr) {
	sid := id.UnmarshalServiceId(req.GetServiceId())

	if sid.IsZeroOrEmptyValue() {
		return nil, syscall.KernelErr_BadId
	}

	methName := req.GetMethodName()
	resp := &syscall.BindMethodResponse{}

	methMap, ok := w.serviceToMethMap[sid.String()]
	if !ok {
		methMap = make(map[string]id.MethodId)
		w.serviceToMethMap[sid.String()] = methMap
	}
	mid, ok := methMap[methName]
	log.Printf("xxxx -- method map %+v (%v)", methMap, ok)
	if !ok {
		mid = id.NewMethodId()
		methMap[methName] = mid
	}
	resp.MethodId = mid.Marshal()
	return returnResponseOrMarshalError(w, resp)
}

// serviceByName looks up a service by the name it used at registration time.
func (w *wheeler) serviceByName(req *syscall.ServiceByNameRequest) (*anypb.Any, syscall.KernelErr) {
	hb := []*syscall.HostBinding{}

	nameMap, ok := w.pkgToRegistration[req.Fqs.PackagePath]
	if ok {
		list, ok := nameMap[req.Fqs.Service]
		if ok {
			hb = make([]*syscall.HostBinding, len(list))
			for i, elem := range list {
				current := &syscall.HostBinding{
					ServiceId: elem.service.Marshal(),
					HostId:    elem.host.Marshal(),
				}
				hb[i] = current
			}
		}
	}
	resp := &syscall.ServiceByNameResponse{}
	resp.Binding = hb
	return returnResponseOrMarshalError(w, resp)
}

// require is the way a service expresses what _other_ services it needs.
// any service name that is going to be looked up with "Locate" should have
// had the service doing said lookup require it beforehand.
func (w *wheeler) require(req *syscall.RequireRequest) (*anypb.Any, syscall.KernelErr) {
	resp := &syscall.RequireResponse{}
	src := id.UnmarshalServiceId(req.GetSource())
	dest := make([]fqName, len(req.GetDest()))
	log.Printf("xxxx -- require: source %s has %d targets", src.Short(), len(dest))
	for i, d := range req.GetDest() {
		curr := fqName{
			pkg:  d.PackagePath,
			name: d.Service,
		}
		dest[i] = curr
	}
	wait, ok := w.serviceToWaiting[src.String()]
	if !ok {
		wait = []fqName{}
	}
	for i := 0; i < len(dest); i++ {
		wait = append(wait, dest[i])
	}
	w.serviceToWaiting[src.String()] = wait

	// if err := dfs(); err != syscall.KernelErr_NoError {
	// 	return nil, err
	// }
	return returnResponseOrMarshalError(w, resp)
}

// dfs is a depth first search starting at the first parameter.  This checks
// to see if the first parameter implies a loop because it includes itself.
func dfs(current, target fqName, path []fqName, isStart bool) syscall.KernelErr {
	return syscall.KernelErr_NoError
}

// returnResponseOrMarshalError is a convenience wrapper around marshalling
// a response and returning the result, which might be a marshal error if
// we can't create the anypb.Any.  This is used by most of the methods
// of wheeler.
func returnResponseOrMarshalError[T proto.Message](w *wheeler, resp T) (*anypb.Any, syscall.KernelErr) {
	a := &anypb.Any{}
	if err := a.MarshalFrom(resp); err != nil {
		w.errorf("marshal failed: %s", err.Error())
		return nil, syscall.KernelErr_MarshalFailed
	}
	return a, syscall.KernelErr_NoError
}

// locate is a foundational call in parigot. It converts a fully qualified
// string (like foo.v1.Foo) into a reference to service and host.  This call
// fails if nobody has exported the service requested.  Since you should always
// call Launch() before Locate() there is no issue with timing.  Launch blocks
// until all dependencies are met.
func (w *wheeler) locate(req *syscall.LocateRequest) (*anypb.Any, syscall.KernelErr) {
	resp := &syscall.LocateResponse{}
	sid := id.UnmarshalServiceId(req.GetCalledBy())
	pkg := req.GetPackageName()
	name := req.GetServiceName()
	log.Printf("xxx -- locate! %s trying to locate %s.%s", sid.Short(), pkg, name)
	sMap, ok := w.pkgToServiceImpl[pkg]
	if !ok {
		w.errorf("failed to find service that was requested in Locate (0): %s.%s", strings.ToUpper(pkg), name)
		return nil, syscall.KernelErr_NotFound
	}
	target, ok := sMap[name]
	if !ok {
		w.errorf("failed to find service that was requested in Locate (1): %s.%s", pkg, strings.ToUpper(name))
		return nil, syscall.KernelErr_NotFound
	}
	if len(target) == 0 {
		w.errorf("failed to find service that was requested in Locate (2): %s.%s", strings.ToUpper(pkg), strings.ToUpper(name))
		return nil, syscall.KernelErr_NotFound

	}
	chosen := target[0]
	if len(target) > 1 {
		n := rand.Intn(len(target))
		chosen = target[n]
	}

	resp.ServiceId = chosen.service.Marshal()
	resp.HostId = chosen.host.Marshal()
	meth, ok := w.serviceToMethMap[chosen.service.String()]
	if !ok {
		// no methods?
		w.errorf("unable to find any methods associated with service %s (%s.%s)", chosen.service.Short(), pkg, name)
	}
	resp.Binding = []*syscall.MethodBinding{}
	for name, meth := range meth {
		mb := &syscall.MethodBinding{
			MethodName: name,
			MethodId:   meth.Marshal(),
		}
		resp.Binding = append(resp.Binding, mb)
	}
	return returnResponseOrMarshalError(w, resp)
}

// dependencyExists is a means of  requesting to know what dependencies
// a given service has.  One version takes a fully qualified name and
// this function will see if that name is something that was  required
// by the service.  The other version traces through all the _running_
// dependencies of a given service to see if there is a connection
// between the two services.
func (w *wheeler) dependencyExists(req *syscall.DependencyExistsRequest) (*anypb.Any, syscall.KernelErr) {
	source := id.UnmarshalServiceId(req.SourceServiceId)
	if source.IsZeroOrEmptyValue() {
		return nil, syscall.KernelErr_BadId
	}
	dest := id.UnmarshalServiceId(req.DestServiceId)
	if dest.IsZeroOrEmptyValue() {
		// we are in the case where it's a name
		list, ok := w.serviceToWaiting[source.String()]
		if !ok {
			resp := &syscall.DependencyExistsResponse{Exists: false}
			return returnResponseOrMarshalError(w, resp)
		}
		for _, elem := range list {
			if elem.name == req.GetServiceName().Service && elem.pkg == req.GetServiceName().PackagePath {
				resp := &syscall.DependencyExistsResponse{Exists: true}
				return returnResponseOrMarshalError(w, resp)
			}
		}
		resp := &syscall.DependencyExistsResponse{Exists: false}
		return returnResponseOrMarshalError(w, resp)
	}
	// this case is the dep edge case
	panic("not yet implemented")
}

// serviceById looks up a service by the string representation of its id.
// this mostly useful if you want to send a service over the wire--you
// can send the string representation and the receiver can use this to
// instantiate the service on their side.
func (w *wheeler) serviceById(req *syscall.ServiceByIdRequest) (*anypb.Any, syscall.KernelErr) {
	h, ok := w.serviceToHost[req.ServiceId]

	resp := &syscall.ServiceByIdResponse{}

	// does it have a host?
	if ok {
		sid, ok := w.stringToService[req.ServiceId]
		//seen it before (registered)?
		if ok {
			resp.Binding = &syscall.HostBinding{ServiceId: sid.Marshal(), HostId: h.Marshal()}
		}
	}
	return returnResponseOrMarshalError(w, resp)
}

// Run repeatedly reads InProtoPairs from the channel associated with
// the wheeler.  It looks at the type of the request inside the InProtoPair
// and dispatches req to the correct function to handle it.  In most cases,
// it also takes the return values from the called function and packages
// them up into an OutProtoPair and sends through the channel to the originator
// of the call.
func (w *wheeler) Run() {
	for {
		in := <-w.ch
		desc := in.Msg.ProtoReflect().Descriptor()
		var result *anypb.Any
		var err syscall.KernelErr
		switch desc.FullName() {
		case "syscall.v1.ExportRequest":
			result, err = w.export((*syscall.ExportRequest)(in.Msg.(*syscall.ExportRequest)))
		case "syscall.v1.ExitRequest":
			result, err = w.exit((*syscall.ExitRequest)(in.Msg.(*syscall.ExitRequest)))
		case "syscall.v1.BindMethodRequest":
			result, err = w.bindMethod((*syscall.BindMethodRequest)(in.Msg.(*syscall.BindMethodRequest)))
		case "syscall.v1.RegisterRequest":
			result, err = w.register((*syscall.RegisterRequest)(in.Msg.(*syscall.RegisterRequest)))
		case "syscall.v1.ServiceByNameRequest":
			result, err = w.serviceByName((*syscall.ServiceByNameRequest)(in.Msg.(*syscall.ServiceByNameRequest)))
		case "syscall.v1.ServiceByIdRequest":
			result, err = w.serviceById((*syscall.ServiceByIdRequest)(in.Msg.(*syscall.ServiceByIdRequest)))
		case "syscall.v1.RequireRequest":
			result, err = w.require((*syscall.RequireRequest)(in.Msg.(*syscall.RequireRequest)))
		case "syscall.v1.DependencyExistsRequest":
			result, err = w.dependencyExists((*syscall.DependencyExistsRequest)(in.Msg.(*syscall.DependencyExistsRequest)))
		case "syscall.v1.LocateRequest":
			result, err = w.locate((*syscall.LocateRequest)(in.Msg.(*syscall.LocateRequest)))
		case "syscall.v1.LaunchRequest":
			result, err = w.launch((*syscall.LaunchRequest)(in.Msg.(*syscall.LaunchRequest)), in.Ch)
		default:
			log.Printf("ERROR! wheeler received unknown type %s", desc.FullName())
			continue
		}
		var a anypb.Any
		e := a.MarshalFrom(result)
		if e != nil {
			select {
			case in.Ch <- OutProtoPair{nil, syscall.KernelErr_MarshalFailed}:
			case <-time.After(1 * time.Second):
				w.errorf("unable to reach client (0) that requested respose (%s) with marshal failed error %s", desc.FullName(), e.Error())
				close(in.Ch)
			}
			return
		}
		outPair := OutProtoPair{
			A:   result,
			Err: err,
		}
		select {
		case in.Ch <- outPair:
		case <-time.After(1 * time.Second):
			w.errorf("unable to reach client (1) that requested response (%T) with response %+v", result, result)
			close(in.Ch)
		}
	}
}
