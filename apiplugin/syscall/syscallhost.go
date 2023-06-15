package syscall

import (
	"context"
	"log"
	_ "unsafe"

	"github.com/iansmith/parigot/apiplugin"
	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/eng"
	syscall "github.com/iansmith/parigot/g/syscall/v1"

	"github.com/tetratelabs/wazero/api"
)

////-----  //go:linkname ParigotInitialize ParigotInitialize

type SyscallPlugin struct{}

var ParigotInitialize apiplugin.ParigotInit = &SyscallPlugin{}

func (*SyscallPlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "parigot", "locate_", locate)
	e.AddSupportedFunc(ctx, "parigot", "dispatch_", dispatch)
	e.AddSupportedFunc(ctx, "parigot", "block_until_call_", blockUntilCall)
	e.AddSupportedFunc(ctx, "parigot", "bind_method_", bindMethod)
	e.AddSupportedFunc(ctx, "parigot", "run_", run)
	e.AddSupportedFunc(ctx, "parigot", "export_", export)
	e.AddSupportedFunc(ctx, "parigot", "return_value_", returnValue)
	e.AddSupportedFunc(ctx, "parigot", "require_", require)
	e.AddSupportedFunc(ctx, "parigot", "register_", register)
	e.AddSupportedFunc(ctx, "parigot", "exit_", exit)

	return true
}

//
// Syscall host implementations
//

func exportImpl(ctx context.Context, req *syscall.ExportRequest, resp *syscall.ExportResponse) int32 {
	for _, fullyQualified := range req.GetService() {
		sid, _ := coordinator().SetService(ctx, fullyQualified.GetPackagePath(), fullyQualified.GetService(), false)

		if coordinator().Export(ctx, sid.Id()) == nil {
			return int32(syscall.KernelErr_NotFound)
		}
	}
	return int32(syscall.KernelErr_NoError)

}

func runImpl(ctx context.Context, req *syscall.RunRequest, resp *syscall.RunResponse) int32 {
	sid := id.UnmarshalServiceId(req.GetServiceId())
	return int32(coordinator().Run(ctx, sid))
}

func locateImpl(ctx context.Context, req *syscall.LocateRequest, resp *syscall.LocateResponse) int32 {
	pcontext.Debugf(ctx, "start of locate impl: req is sender=%v,%v", id.UnmarshalServiceId(req.CalledBy),
		req.GetPackageName()+"."+req.GetServiceName())
	svc, ok := coordinator().SetService(ctx, req.GetPackageName(), req.GetServiceName(), false)
	if ok {
		return int32(syscall.KernelErr_NotFound)
	}
	calledBy := id.UnmarshalServiceId(req.CalledBy)
	if !coordinator().PathExists(ctx, calledBy.String(), svc.String()) {
		return int32(syscall.KernelErr_NotRequired)
	}
	svcId := svc.Id()
	pcontext.Debugf(ctx, "at end of locate, we are returning %s", svcId.Short())
	resp.ServiceId = svcId.Marshal()
	return int32(syscall.KernelErr_NoError)
}

func registerImpl(ctx context.Context, req *syscall.RegisterRequest, resp *syscall.RegisterResponse) int32 {
	svc, _ := coordinator().SetService(ctx, req.Fqs.GetPackagePath(), req.Fqs.GetService(), req.GetIsClient())
	resp.Id = svc.Id().Marshal()
	return int32(syscall.KernelErr_NoError)
}

func requireImpl(ctx context.Context, req *syscall.RequireRequest, resp *syscall.RequireResponse) int32 {
	src := id.UnmarshalServiceId(req.GetSource())
	fqn := req.GetDest()
	//log.Printf("xxx--- require impl, source = %s", src.Short())

	for _, fullyQualified := range fqn {
		dest, _ := coordinator().SetService(ctx, fullyQualified.GetPackagePath(), fullyQualified.GetService(), false)
		// if ok {
		// 	pcontext.Infof(ctx, "requireImpl: created new service id %s.%s => %s", fullyQualified.GetPackagePath(),
		// 		fullyQualified.GetService(), dest.Short())
		// }
		kerr := coordinator().Import(ctx, src, dest.Id())
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
	log.Printf("dispatch 0x%x", stack)
}

func blockUntilCall(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("blockUntilCall 0x%x", stack)
}
func bindMethod(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("bindMethod 0x%x", stack)
}
func run(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.RunRequest{}
	resp := (*syscall.RunResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]export", m, stack, runImpl, req, resp)
}
func export(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.ExportRequest{}
	resp := (*syscall.ExportResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]export", m, stack, exportImpl, req, resp)
}
func returnValue(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("returnValue 0x%x", stack)
}

func require(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.RequireRequest{}
	resp := (*syscall.RequireResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]require", m, stack, requireImpl, req, resp)

}

func register(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.RegisterRequest{}
	resp := &syscall.RegisterResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[syscall]register", m, stack, registerImpl, req, resp)
}

func exit(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("exit 0x%x", stack)
	panic("exit called ")
}
