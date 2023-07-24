package syscall

import (
	"context"
	"fmt"
	"log"
	"reflect"
	_ "unsafe"

	apiplugin "github.com/iansmith/parigot/api/plugin"
	"github.com/iansmith/parigot/api/plugin/syscall/wheeler"
	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/eng"
	syscall "github.com/iansmith/parigot/g/syscall/v1"

	"github.com/tetratelabs/wazero/api"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var pairIdToChannel = make(map[string]chan CallInfo)

// CallInfo is sent to the channels that represent service/method calls.
type CallInfo struct {
	cid   id.CallId
	param *anypb.Any
}

type SyscallPlugin struct {
}

var ParigotInitialize apiplugin.ParigotInit = &SyscallPlugin{}

func (*SyscallPlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "parigot", "locate_", locate)
	e.AddSupportedFunc(ctx, "parigot", "dispatch_", dispatch)
	e.AddSupportedFunc(ctx, "parigot", "block_until_call_", blockUntilCall)
	e.AddSupportedFunc(ctx, "parigot", "bind_method_", bindMethod)
	e.AddSupportedFunc(ctx, "parigot", "launch_", launch)
	e.AddSupportedFunc(ctx, "parigot", "export_", export)
	e.AddSupportedFunc(ctx, "parigot", "return_value_", returnValue)
	e.AddSupportedFunc(ctx, "parigot", "require_", require)
	e.AddSupportedFunc(ctx, "parigot", "register_", register)
	e.AddSupportedFunc(ctx, "parigot", "exit_", exit)
	e.AddSupportedFunc(ctx, "parigot", "read_one_", readOne)

	return true
}

func fqServiceName(p, s string) string {
	return fmt.Sprintf("%s.%s", p, s)
}

//
// Syscall host implementations
//

func exportImpl(ctx context.Context, req *syscall.ExportRequest, resp *syscall.ExportResponse) int32 {
	hid := id.UnmarshalHostId(req.GetHostId())

	///////
	// retCh := make(chan wheeler.OutProtoPair)
	// inPair := wheeler.InProtoPair{
	// 	Msg: req,
	// 	Ch:  retCh,
	// }
	// wheeler.In() <- inPair
	// out := <-retCh
	// if out.Err != 0 {
	// 	log.Printf("error in export impl: %s", syscall.KernelErr_name[int32(out.Err)])
	// }
	// log.Printf("xxxx --- %+v, %d", out.Msg, out.Err)
	////
	for _, fullyQualified := range req.GetService() {
		sid, _ := startCoordinator().SetService(ctx, fullyQualified.GetPackagePath(), fullyQualified.GetService(), false)

		sid.Export()

		fqs := fqServiceName(fullyQualified.GetPackagePath(), fullyQualified.GetService())
		if kerr := finder().AddHost(fqs, hid); kerr != syscall.KernelErr_NoError {
			return int32(kerr)
		}
	}

	return int32(syscall.KernelErr_NoError)
}

func exitImpl(ctx context.Context, req *syscall.ExitRequest, resp *syscall.ExitResponse) int32 {
	//return handleByWheeler(req, resp)

	return int32(0x7fffff00 | (req.Code & 0xff))
}

func launchImpl(ctx context.Context, req *syscall.LaunchRequest, resp *syscall.LaunchResponse) int32 {
	sid := id.UnmarshalServiceId(req.GetServiceId())
	return int32(startCoordinator().Launch(ctx, sid))
}

func bindMethodImpl(ctx context.Context, req *syscall.BindMethodRequest, resp *syscall.BindMethodResponse) int32 { //syscall.KernelErr {
	return int32(handleByWheeler(req, resp))
}

func readOneImpl(ctx context.Context, req *syscall.ReadOneRequest, resp *syscall.ReadOneResponse) int32 {
	hid := id.UnmarshalHostId(req.HostId)
	rc, err := matcher().Ready(hid)
	if err != syscall.KernelErr_NoError {
		return int32(err)
	}
	// we favor resolving calls, which may be a terrible idea
	if rc != nil {
		resp.Timeout = false
		resp.Call = nil
		resp.Param = nil
		resp.Resolved = rc
		return int32(syscall.KernelErr_NoError)
	}

	cases := []reflect.SelectCase{}
	// now we are going to listen for a message on one of the channels
	// we can also timeout.  the order of these, sadly, matters and
	// the service/method listeners must go first because the index
	// of the channel is how we figure out how to dispatch the method.
	timeoutChoice, exitChoice := -1, -1

	mcl := newMethodCallListener(req)
	cases = append(cases, mcl.Case()...)

	tl := newTimeoutListener(req.TimeoutInMillis)
	c := tl.Case()
	if len(c) != 0 {
		timeoutChoice = len(cases) // we are about to fill the spot
	}
	cases = append(cases, c...)

	exitChannel := (chan int32)(nil)
	el := NewExitListener(exitChannel)
	c = el.Case()
	if len(c) != 0 {
		exitChoice = len(cases)
	}
	cases = append(cases, c...)

	if len(cases) == 0 { // very unlikely since there is the possibility of exit
		resp.Call = nil
		resp.Param = nil
		resp.Timeout = false
		resp.Exit = false
		return int32(syscall.KernelErr_NoError)
	}

	// run the select
	chosen, value, ok := reflect.Select(cases)
	// ok will be true if the channel has not been closed.
	if !ok {
		return int32(syscall.KernelErr_KernelConnectionFailed)
	}
	switch chosen {
	case timeoutChoice:
		tl.Handle(value, chosen, resp)
	case exitChoice:
		el.Handle(value, chosen, resp)
	default:
		mcl.Handle(value, chosen, resp)
	}

	return int32(syscall.KernelErr_NoError)
}

func returnValueImpl(ctx context.Context, req *syscall.ReturnValueRequest, resp *syscall.ReturnValueResponse) int32 {
	cid := id.UnmarshalCallId(req.GetCallId())
	kerr := matcher().Response(cid, req.Result, req.ResultError)
	return int32(kerr)
}

func locateImpl(ctx context.Context, req *syscall.LocateRequest, resp *syscall.LocateResponse) int32 {
	svc, ok := startCoordinator().SetService(ctx, req.GetPackageName(),
		req.GetServiceName(), false)
	if ok {
		return int32(syscall.KernelErr_NotFound)
	}
	calledBy := id.UnmarshalServiceId(req.CalledBy)
	if !startCoordinator().PathExists(ctx, calledBy.String(), svc.String()) {
		return int32(syscall.KernelErr_NotRequired)
	}
	host := finder().FindByName(fqServiceName(req.GetPackageName(), req.GetServiceName()))
	if host == nil {
		return int32(syscall.KernelErr_NotFound)
	}
	svcId := svc.Id()
	resp.ServiceId = svcId.Marshal()
	resp.Binding = svc.Method()
	resp.HostId = host.hid.Marshal()
	return int32(syscall.KernelErr_NoError)
}

func dispatchImpl(ctx context.Context, req *syscall.DispatchRequest, resp *syscall.DispatchResponse) int32 {
	sid := id.UnmarshalServiceId(req.GetServiceId())
	mid := id.UnmarshalMethodId(req.GetMethodId())
	cid := id.UnmarshalCallId(req.GetCallId())
	hid := id.UnmarshalHostId(req.GetHostId())

	matcher().Dispatch(hid, cid)

	target := pairIdToChannel[makeSidMidCombo(sid, mid)]
	if target == nil {
		// should this have a special error?
		return int32(syscall.KernelErr_NotFound)
	}

	resp.CallId = cid.Marshal()

	cm := CallInfo{
		cid:   cid,
		param: req.GetParam(),
	}
	target <- cm
	return int32(syscall.KernelErr_NoError)
}

func registerImpl(ctx context.Context, req *syscall.RegisterRequest, resp *syscall.RegisterResponse) int32 { //syscall.KernelErr {
	return int32(handleByWheeler(req, resp))
}

func requireImpl(ctx context.Context, req *syscall.RequireRequest, resp *syscall.RequireResponse) int32 {
	if req.GetDest() == nil {
		return 0
	}
	src := id.UnmarshalServiceId(req.GetSource())
	fqn := req.GetDest()

	for _, fullyQualified := range fqn {
		dest, _ := startCoordinator().SetService(ctx, fullyQualified.GetPackagePath(), fullyQualified.GetService(), false)
		kerr := startCoordinator().Import(ctx, src, dest.Id())
		if int32(kerr) != 0 {
			pcontext.Errorf(ctx, "kernel error returned from import: %d", kerr)
			return int32(kerr)
		}
	}
	return int32(syscall.KernelErr_NoError)
}

//
// Syscall marshal/unmarshal for each system call
//

func locate(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.LocateRequest{}
	resp := &syscall.LocateResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[syscall]locate", m, stack, locateImpl, req, resp)
}

func dispatch(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.DispatchRequest{}
	resp := &syscall.DispatchResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[syscall]dispatch", m, stack, dispatchImpl, req, resp)
}

func blockUntilCall(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("blockUntilCall 0x%x", stack)
}
func bindMethod(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.BindMethodRequest{}
	resp := &syscall.BindMethodResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[syscall]bindMethod", m, stack, bindMethodImpl, req, resp)
}

func launch(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.LaunchRequest{}
	resp := (*syscall.LaunchResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]export", m, stack, launchImpl, req, resp)
}
func export(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.ExportRequest{}
	resp := (*syscall.ExportResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]export", m, stack, exportImpl, req, resp)
}
func returnValue(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.ReturnValueRequest{}
	resp := (*syscall.ReturnValueResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]returnValue", m, stack, returnValueImpl, req, resp)
}

func require(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.RequireRequest{}
	resp := (*syscall.RequireResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]require", m, stack, requireImpl, req, resp)

}
func readOne(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.ReadOneRequest{}
	resp := &syscall.ReadOneResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[syscall]readOne", m, stack, readOneImpl, req, resp)

}

func register(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.RegisterRequest{}
	resp := &syscall.RegisterResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[syscall]register", m, stack, registerImpl, req, resp)
}

func exit(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.ExitRequest{}
	resp := &syscall.ExitResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[syscall]register", m, stack, exitImpl, req, resp)
}

func makeSidMidCombo(sid id.ServiceId, mid id.MethodId) string {
	return sid.String() + "," + mid.String()
}

func handleByWheeler[T proto.Message, U proto.Message](t T, u U) syscall.KernelErr {
	retCh := make(chan wheeler.OutProtoPair)
	inPair := wheeler.InProtoPair{
		Ch: retCh,
	}
	inPair.Msg = t
	wheeler.In() <- inPair
	out := <-retCh
	if out.Err != 0 {
		log.Printf("error in wheeler impl: %T, %s", t, syscall.KernelErr_name[int32(out.Err)])
	}
	err := out.A.UnmarshalTo(u)
	if err != nil {
		return syscall.KernelErr_MarshalFailed
	}
	return out.Err
}
