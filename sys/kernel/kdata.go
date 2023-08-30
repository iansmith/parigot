package kernel

import (
	"sync"

	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/syscall/v1"
)

// kdata is one of the core kernel data structures. The data structures
// this object holds are primarily for the ability to block and correctly
// return when new input event has happened.
type kdata struct {
	lock sync.Mutex

	rawRecv []GeneralReceiver

	// computed based on what we actually get passed
	reg []Registrar

	exporter []Exporter

	// computed based on what we actually get passed
	bind []Binder

	ns    Nameserver
	start Starter

	match callMatcher

	// useful channel
	cancel chan bool
}

var _ Kernel = &kdata{}

// newKData returns an initialized kernel
func newKData() *kdata {
	return &kdata{
		cancel: make(chan bool),
		match:  newCallMatcher(),
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
	for _, candidate := range []interface{}{r, f, n, st} {
		if e, ok := candidate.(Exporter); ok {
			k.exporter = append(k.exporter, e)
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
	// we don't want to lock here because we could block somebody
	// else who is reading from the same channel

	sid := id.UnmarshalServiceId(req.GetBundle().GetServiceId())
	//hid := id.UnmarshalHostId(req.GetBundle().GetHostId())
	mid := id.UnmarshalMethodId(req.GetBundle().GetMethodId())

	targetHid := k.Nameserver().FindHost(sid)
	if targetHid.IsZeroOrEmptyValue() {
		return syscall.KernelErr_BadId
	}
	cid := id.UnmarshalCallId(req.GetBundle().GetCallId())
	k.matcher().Dispatch(targetHid, cid, mid)
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
	k.matcher().Dispatch(hid, cid, mid)
	return k.start.Launch(sid, cid, hid, mid)
}

// Export binds a particular serviceid to a given name.  The name is the name
// of an interface that allows the service to be found by other services.
func (k *kdata) Export(req *syscall.ExportRequest, resp *syscall.ExportResponse) syscall.KernelErr {
	sid := id.UnmarshalServiceId(req.GetServiceId())
	hid := id.UnmarshalHostId(req.GetHostId())
	if hid.IsZeroOrEmptyValue() {
		return syscall.KernelErr_BadId
	}
	if sid.IsZeroOrEmptyValue() {
		return syscall.KernelErr_BadId
	}

	fqsn := req.GetService()
	for _, fqs := range fqsn {
		p := fqs.GetPackagePath()
		n := fqs.GetService()
		fqn := FQName{Pkg: p, Name: n}
		for _, exp := range k.exporter {
			kerr := exp.Export(hid, sid, fqn)
			if kerr != syscall.KernelErr_NoError {
				return kerr
			}
		}
	}
	return syscall.KernelErr_NoError
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

// func dumpRC(rc *syscall.ResolvedCall) string {
// 	hid := id.UnmarshalHostId(rc.GetHostId())
// 	cid := id.UnmarshalCallId(rc.GetCallId())
// 	mid := id.UnmarshalMethodId(rc.GetMethodId())
// 	return fmt.Sprintf("RC[%s,%s,%s]", hid.Short(), cid.Short(), mid.Short())
// }

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

func (k *kdata) Exit(req *syscall.ExitRequest, resp *syscall.ExitResponse) syscall.KernelErr {
	k.lock.Lock()
	defer k.lock.Unlock()

	cid := id.UnmarshalCallId(req.GetCallId())
	hid := id.UnmarshalHostId(req.GetHostId())
	mid := id.UnmarshalMethodId(req.GetMethodId())

	// save for later
	k.matcher().Dispatch(hid, cid, mid)
	if req.ShutdownAll {
		for _, host := range k.ns.AllHosts() {
			kerr := k.matcher().Dispatch(host, cid, mid)
			if kerr != syscall.KernelErr_NoError {
				return syscall.KernelErr_ExitFailed
			}
		}
	}
	resp.Pair = req.Pair

	return syscall.KernelErr_NoError
}
