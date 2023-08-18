package kernel

import (
	"sync"

	sharedapi "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/types/known/anypb"
)

// kdata is one of the core kernel data structures. The data structures
// this object holds are primarily for the ability to block and correctly
// return when new input event has happened.
type kdata struct {
	lock sync.Mutex

	rawRecv []GeneralReceiver

	// computed based on what we actually get passed
	reg []Registrar

	// computed based on what we actually get passed
	bind []Binder

	ns    Nameserver
	start Starter

	match callMatcher

	// this is so we can dispatch multiple messages at a time but
	// only one at a time PER GUEST
	serviceIdToWG map[string]sync.WaitGroup

	// useful channel
	cancel chan bool
}

var _ Kernel = &kdata{}

// newKData returns an initialized kernel
func newKData() *kdata {
	return &kdata{
		cancel:        make(chan bool),
		serviceIdToWG: make(map[string]sync.WaitGroup),
		match:         newCallMatcher(),
	}
}

func (k *kdata) AddReceiver(r GeneralReceiver) {
	k.lock.Lock()
	defer k.lock.Unlock()

	k.rawRecv = append(k.rawRecv, r)
}

// SetApproach sets a number of key subsystems in the kernel and should only be
// called when the kernel is a fresh state, as this call resets many internal data
// structures.
func (k *kdata) SetApproach(r GeneralReceiver, f GeneralReceiver, n Nameserver, st Starter) syscall.KernelErr {

	//k.rawSend = append(k.rawSend, s)
	k.rawRecv = append(k.rawRecv, r)
	k.rawRecv = append(k.rawRecv, f)

	k.ns = n
	k.start = st

	for _, candidate := range []interface{}{r, f, n, st} {
		if reg, ok := candidate.(Registrar); ok {
			k.reg = append(k.reg, reg)
		}
	}
	for _, candidate := range []interface{}{r, f, n, st} {
		if b, ok := candidate.(Binder); ok {
			k.bind = append(k.bind, b)
		}
	}
	return syscall.KernelErr_NoError
}

// Register is used to notify the kernel that a given service
// should be assigned a service id.  Note that this may reach multiple
// parts of the kernel based on the Registrar interface.
func (k *kdata) Register(req *syscall.RegisterRequest, resp *syscall.RegisterResponse) syscall.KernelErr {
	hid := id.UnmarshalHostId(req.GetHostId())
	sid := id.NewServiceId()
	debugName := req.GetDebugName()
	result := syscall.KernelErr_NoError
	for _, r := range k.reg {
		if kerr := r.Register(hid, sid, debugName); kerr != syscall.KernelErr_NoError {
			klog.Errorf("unexpected failure in registrar: %s", syscall.KernelErr_name[int32(kerr)])
			result = kerr
		}
	}
	resp.ServiceId = sid.Marshal()
	return result
}

// func newService(hid id.HostId, sid id.ServiceId, debugName string) syscall.KernelErr {
// 	return syscall.KernelErr_NoError
// }

func (k *kdata) matcher() callMatcher {
	return k.match
}

// Dispatch is used to send a call to a remote machine.  If this
// returns a kernel error it is because the dispatch call itself could
// not be made, not that the dispatch worked ok and an error was returned
// by the remote code.
func (k *kdata) Dispatch(req *syscall.DispatchRequest, resp *syscall.DispatchResponse) syscall.KernelErr {
	k.lock.Lock()
	defer k.lock.Unlock()

	sid := id.UnmarshalServiceId(req.GetBundle().GetServiceId())
	hid := id.UnmarshalHostId(req.GetBundle().GetHostId())

	targetHid := k.Nameserver().FindHost(sid)
	if targetHid.IsZeroOrEmptyValue() {
		return syscall.KernelErr_BadId
	}
	cid := id.UnmarshalCallId(req.GetBundle().GetCallId())
	k.matcher().Dispatch(hid, cid)
	ch := k.Nameserver().FindHostChan(targetHid)
	ch <- req
	resp.CallId = cid.Marshal()
	resp.TargetHostId = targetHid.Marshal()
	return syscall.KernelErr_NoError
}

// ReturnValue is used to finish a previous Dispatch call.  This is where the
// original caller will get his call completed.
func (k *kdata) ReturnValue(req *syscall.ReturnValueRequest, resp *syscall.ReturnValueResponse) syscall.KernelErr {
	cid := id.UnmarshalCallId(req.GetBundle().GetCallId())
	kerr := k.matcher().Response(cid, req.Result, req.ResultError)
	if kerr != syscall.KernelErr_NoError {
		return kerr
	}

	return syscall.KernelErr_NoError
}

// Launch logically causes a process to wait for all its dependencies to
// be ready.  In practice, it returns immediately and then finishes the
// process later.
func (k *kdata) Launch(req *syscall.LaunchRequest, resp *syscall.LaunchResponse) syscall.KernelErr {
	k.lock.Lock()
	defer k.lock.Unlock()

	sid := id.UnmarshalServiceId(req.GetServiceId())
	cid := id.UnmarshalCallId(req.GetCallId())
	hid := id.UnmarshalHostId(req.GetHostId())
	mid := id.UnmarshalMethodId(req.GetMethodId())

	// save for later
	k.matcher().Dispatch(hid, cid)
	return k.start.Launch(sid, cid, hid, mid)
}

// Export binds a particular serviceid to a given name.  The name is the name
// of an interface that allows the service to be found by other services.
func (k *kdata) Export(req *syscall.ExportRequest, resp *syscall.ExportResponse) syscall.KernelErr {
	return k.start.Export(req, resp)
}

// Locate is the constructor for parigot types.
func (k *kdata) Locate(req *syscall.LocateRequest, resp *syscall.LocateResponse) syscall.KernelErr {
	return k.start.Locate(req, resp)
}

// Require establishes a dependency(ies) from source to dest.
func (k *kdata) Require(req *syscall.RequireRequest, resp *syscall.RequireResponse) syscall.KernelErr {
	return k.start.Require(req, resp)
}

// BindMethod connects a method name to a method id.
func (k *kdata) BindMethod(req *syscall.BindMethodRequest, resp *syscall.BindMethodResponse) syscall.KernelErr {
	sid := id.UnmarshalServiceId(req.GetServiceId())
	hid := id.UnmarshalHostId(req.GetHostId())
	mid := id.NewMethodId()

	name := req.GetMethodName()

	for _, b := range k.bind {
		b.Bind(hid, sid, mid, name)
	}
	resp.MethodId = mid.Marshal()

	return syscall.KernelErr_NoError
}

// CancelRead should be call to gracefully exit its read loop.  This is
// used when you lock, make changes, unlock and then want the kernel
// to pick up your change.
func (k *kdata) CancelRead() {
	k.cancel <- true
}

func (k *kdata) Nameserver() Nameserver {
	return k.ns
}

func (k *kdata) responseReady(hid id.HostId, resp *syscall.ReadOneResponse) syscall.KernelErr {
	rc, err := k.matcher().Ready(hid)
	if err != syscall.KernelErr_NoError {
		return err
	}
	if rc == nil {
		resp.Resolved = nil
		return syscall.KernelErr_NoError
	}
	resp.Timeout = false
	resp.Bundle = &syscall.MethodBundle{}
	resp.ParamOrResult = nil
	resp.ResultErr = 0
	resp.Resolved = rc
	resp.Exit = &syscall.ExitPair{}

	return syscall.KernelErr_NoError
}

// exit does a dispatch, so it cannot lock

func (k *kdata) Exit(req *syscall.ExitRequest, resp *syscall.ExitResponse) syscall.KernelErr {
	all := k.ns.AllHosts()
	a := &anypb.Any{}
	err := a.MarshalFrom(req.Pair)
	if err != nil {
		klog.Errorf("unable to marshal pair into any: %v", err)
		return syscall.KernelErr_MarshalFailed
	}
	bundle := &syscall.MethodBundle{
		ServiceId: req.Pair.ServiceId,
		MethodId:  sharedapi.ExitMethod.Marshal(),
		CallId:    id.CallIdEmptyValue().Marshal(),
	}
	for _, host := range all {
		bundle.HostId = host.Marshal()
		dispatchReq := &syscall.DispatchRequest{
			Bundle: bundle,
			Param:  a,
		}
		resp := &syscall.DispatchResponse{}
		kerr := k.Dispatch(dispatchReq, resp)
		if err != nil {
			klog.Errorf("unable to dispatch exit message: %s", syscall.KernelErr_name[int32(kerr)])
			return syscall.KernelErr_ExitFailed
		}
	}
	return syscall.KernelErr_NoError
}
