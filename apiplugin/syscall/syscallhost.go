package main

import (
	"context"
	"log"

	"github.com/iansmith/parigot/apiplugin"
	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/eng"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"

	"github.com/tetratelabs/wazero/api"
)

type syscallPlugin struct{}

var ParigiotInitialize = syscallPlugin{}

// xxx global vairable kinda sucks
var currentEng eng.Engine

func (*syscallPlugin) Init(ctx context.Context, e eng.Engine) bool {
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

func exportImpl(ctx context.Context, req *syscallmsg.ExportRequest, resp *syscallmsg.ExportResponse) id.IdRaw {
	for _, fullyQualified := range req.GetService() {
		sid, ok := depData().SetService(ctx, fullyQualified.GetPackagePath(), fullyQualified.GetService())
		if !ok {
			pcontext.Debugf(ctx, "created new service because of export %s.%s", fullyQualified.GetPackagePath(),
				fullyQualified.GetService(), sid.Short())
		}
		pcontext.Debugf(ctx, "exported service id %s.%s => %s", fullyQualified.GetPackagePath(),
			fullyQualified.GetService(), sid.Short())
		if depData().Export(ctx, sid.Id()) == nil {
			return id.NewKernelErrId(id.KernelNotFound).Raw()
		}
	}
	return id.KernelErrIdNoErr.Raw()

}
func returnKernelErrorForIdErr(ctx context.Context, ierr id.IdErr) id.KernelErrId {
	pcontext.Errorf(ctx, "unable to unmrarshal the service id in a require request: %s", ierr.Short())
	// xxx this isn't a great error because the problem is really an Id error
	return id.NewKernelErrId(id.KernelErrIdUnmarshalError)

}

func runImpl(ctx context.Context, req *syscallmsg.RunRequest, resp *syscallmsg.RunResponse) id.IdRaw {
	sid, idErr := id.UnmarshalServiceId(req.GetServiceId())
	if idErr.IsError() {
		returnKernelErrorForIdErr(ctx, idErr)
	}
	depData().Run(ctx, sid)
	return id.KernelErrIdNoErr.Raw()
}

func registerImpl(ctx context.Context, req *syscallmsg.RegisterRequest, resp *syscallmsg.RegisterResponse) id.IdRaw {
	svc, _ := depData().SetService(ctx, req.Fqs.GetPackagePath(), req.Fqs.GetService())
	resp.Id = svc.Id().Marshal()
	return id.KernelErrIdNoErr.Raw()
}

func requireImpl(ctx context.Context, req *syscallmsg.RequireRequest, resp *syscallmsg.RequireResponse) id.IdRaw {
	src, idErr := id.UnmarshalServiceId(req.GetSource())
	if idErr.IsError() {
		returnKernelErrorForIdErr(ctx, idErr)
	}
	fqn := req.GetDest()
	for _, fullyQualified := range fqn {
		dest, ok := depData().SetService(ctx, fullyQualified.GetPackagePath(), fullyQualified.GetService())
		if ok {
			pcontext.Infof(ctx, "requireImpl: created new service id %s.%s => %s", fullyQualified.GetPackagePath(),
				fullyQualified.GetService(), dest.Short())
		}
		if !depData().Import(ctx, src, dest.Id()) {
			return id.NewKernelErrId(id.KernelNotFound).Raw()
		}
	}
	return id.KernelErrIdNoErr.Raw()
}

//
// Syscall marshal/unmarshal for each system call
//

func locate(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("locate %s 0=0x%x, 1=0x%x", m.Name(), stack[0], stack[1])
	stack[0] = 82
	return
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
	req := &syscallmsg.RunRequest{}
	resp := (*syscallmsg.RunResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]export", m, stack, runImpl, req, resp)
}
func export(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscallmsg.ExportRequest{}
	resp := (*syscallmsg.ExportResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]export", m, stack, exportImpl, req, resp)
}
func returnValue(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("returnValue 0x%x", stack)
}

func require(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscallmsg.RequireRequest{}
	resp := (*syscallmsg.RequireResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]require", m, stack, requireImpl, req, resp)
	return
}

func register(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscallmsg.RegisterRequest{}
	resp := &syscallmsg.RegisterResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[syscall]register", m, stack, registerImpl, req, resp)
	return
}

func exit(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("exit 0x%x", stack)
	panic("exit called ")
}
