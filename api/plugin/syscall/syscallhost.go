package syscall

import (
	"context"
	"fmt"
	"log"
	"reflect"
	_ "unsafe"

	apiplugin "github.com/iansmith/parigot/api/plugin"
	"github.com/iansmith/parigot/api/plugin/syscall/kernel"
	"github.com/iansmith/parigot/api/plugin/syscall/wheeler"
	"github.com/iansmith/parigot/eng"
	syscall "github.com/iansmith/parigot/g/syscall/v1"

	"github.com/tetratelabs/wazero/api"
	"google.golang.org/protobuf/proto"
)

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
	e.AddSupportedFunc(ctx, "parigot", "synchronous_exit_", syncExit)

	return true
}

// channel structure for output, including debugging
var debugOutChan = true

type OutChannel struct {
	ch     chan wheeler.OutProtoPair
	origin string
	num    int
}

func fqServiceName(p, s string) string {
	return fmt.Sprintf("%s.%s", p, s)
}

//
// Syscall host implementations
//

func exportImpl(ctx context.Context, req *syscall.ExportRequest, resp *syscall.ExportResponse) int32 {
	return int32(handleByWheeler(req, resp))
}

func exitImpl(ctx context.Context, req *syscall.ExitRequest, resp *syscall.ExitResponse) int32 {
	err := handleByWheeler(req, resp)
	if err != syscall.KernelErr_NoError {
		return int32(err)
	}
	return int32(syscall.KernelErr_NoError)
}

func launchImpl(ctx context.Context, req *syscall.LaunchRequest, resp *syscall.LaunchResponse) int32 {
	err := handleByWheeler(req, resp)
	if err != syscall.KernelErr_NoError {
		return int32(err)
	}

	return int32(syscall.KernelErr_NoError)
}

func bindMethodImpl(ctx context.Context, req *syscall.BindMethodRequest, resp *syscall.BindMethodResponse) int32 { //syscall.KernelErr {
	return int32(handleByWheeler(req, resp))
}

func syncExitImpl(ctx context.Context, req *syscall.SynchronousExitRequest, resp *syscall.SynchronousExitResponse) int32 { //syscall.KernelErr {
	return int32(handleByWheeler(req, resp))
}

func readOneImpl(ctx context.Context, req *syscall.ReadOneRequest, resp *syscall.ReadOneResponse) int32 {
	return int32(handleByWheeler(req, resp))
}

func returnValueImpl(ctx context.Context, req *syscall.ReturnValueRequest, resp *syscall.ReturnValueResponse) int32 {
	return int32(handleByWheeler(req, resp))
}

func locateImpl(ctx context.Context, req *syscall.LocateRequest, resp *syscall.LocateResponse) int32 {
	return int32(handleByWheeler(req, resp))
}

func dispatchImpl(ctx context.Context, req *syscall.DispatchRequest, resp *syscall.DispatchResponse) int32 {
	return int32(handleByWheeler(req, resp))
}

func registerImpl(ctx context.Context, req *syscall.RegisterRequest, resp *syscall.RegisterResponse) int32 { //syscall.KernelErr {
	log.Printf("yyy -- register")
	kerr := kernel.K.Register(req, resp)
	log.Printf("yyy -- register done %d", kerr)
	return int32(handleByWheeler(req, resp))
}

func requireImpl(ctx context.Context, req *syscall.RequireRequest, resp *syscall.RequireResponse) int32 {
	if req.GetDest() == nil {
		return 0
	}
	return int32(handleByWheeler(req, resp))
}

// Syscall marshal/unmarshal for each system call
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
	apiplugin.InvokeImplFromStack(ctx, "[syscall]launch", m, stack, launchImpl, req, resp)
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

func syncExit(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.SynchronousExitRequest{}
	resp := &syscall.SynchronousExitResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[syscall]syncExit", m, stack, syncExitImpl, req, resp)
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

func handleByWheeler[T proto.Message, U proto.Message](t T, u U) syscall.KernelErr {
	retCh := make(chan wheeler.OutProtoPair, 1)
	inPair := wheeler.InProtoPair{
		Ch: retCh,
	}
	inPair.Msg = t
	wheeler.In() <- inPair
	out := <-retCh
	if out.Err != 0 {
		log.Printf("error in wheeler impl: %T, %s", t, syscall.KernelErr_name[int32(out.Err)])
		return out.Err
	}

	if out.A != nil {
		r := reflect.ValueOf(u)
		if !r.IsNil() {
			err := out.A.UnmarshalTo(u)
			if err != nil {
				return syscall.KernelErr_MarshalFailed
			}
		}
	}
	return out.Err
}
