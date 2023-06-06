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
		sid, _ := depData().SetService(ctx, fullyQualified.GetPackagePath(), fullyQualified.GetService(), false)

		if depData().Export(ctx, sid.Id()) == nil {
			return id.NewKernelErrId(id.KernelNotFound).Raw()
		}
	}
	return id.KernelErrIdNoErr.Raw()

}
func returnKernelErrorForIdErr(ctx context.Context, ierr id.IdErr) id.KernelErrId {
	pcontext.Errorf(ctx, "unable to unmarshal the service id in a require request: %s", ierr.Short())
	// xxx this isn't a great error because the problem is really an Id error
	return id.NewKernelErrId(id.KernelErrIdUnmarshalError)

}

func runImpl(ctx context.Context, req *syscallmsg.RunRequest, resp *syscallmsg.RunResponse) id.IdRaw {
	sid, idErr := id.UnmarshalServiceId(req.GetServiceId())
	if idErr.IsError() {
		returnKernelErrorForIdErr(ctx, idErr)
	}
	if !depData().Run(ctx, sid) {
		return id.IdRaw(id.NewKernelErrId(id.KernelDependencyCycle))
	}
	return id.KernelErrIdNoErr.Raw()
}

func locateImpl(ctx context.Context, req *syscallmsg.LocateRequest, resp *syscallmsg.LocateResponse) id.IdRaw {
	pcontext.Debugf(ctx, "start of locate impl: req is sender=%v,%v", id.MustUnmarshalServiceId(req.CalledBy),
		req.GetPackageName()+"."+req.GetServiceName())
	svc, ok := depData().SetService(ctx, req.GetPackageName(), req.GetServiceName(), false)
	if ok {
		return id.NewKernelErrId(id.KernelNotFound).Raw()
	}
	calledBy := id.MustUnmarshalServiceId(req.CalledBy)
	if !depData().PathExists(ctx, calledBy.String(), svc.String()) {
		return id.NewKernelErrId(id.KernelNotRequired).Raw()
	}
	svcId := svc.Id()
	pcontext.Debugf(ctx, "at end of locate, we are returning %s", svcId.Short())
	resp.ServiceId = svcId.Marshal()
	return id.IdRaw(id.KernelErrIdNoErr)
}

func registerImpl(ctx context.Context, req *syscallmsg.RegisterRequest, resp *syscallmsg.RegisterResponse) id.IdRaw {
	svc, _ := depData().SetService(ctx, req.Fqs.GetPackagePath(), req.Fqs.GetService(), req.GetIsClient())
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
		dest, ok := depData().SetService(ctx, fullyQualified.GetPackagePath(), fullyQualified.GetService(), false)
		if ok {
			pcontext.Infof(ctx, "requireImpl: created new service id %s.%s => %s", fullyQualified.GetPackagePath(),
				fullyQualified.GetService(), dest.Short())
		}
		kerr := depData().Import(ctx, src, dest.Id())
		if kerr.IsError() {
			pcontext.Errorf(ctx, "kernel error returned from import: %d", kerr.ErrorCode())
			return kerr.Raw()
		}
	}
	return id.KernelErrIdNoErr.Raw()
}

//
// Syscall marshal/unmarshal for each system call
//

func locate(ctx context.Context, m api.Module, stack []uint64) {
	defer func() {
		if r := recover(); r != nil {
			print("Trapped a panic in locate ", r, "\n")
		}
	}()
	req := &syscallmsg.LocateRequest{}
	resp := &syscallmsg.LocateResponse{}
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
	defer func() {
		if r := recover(); r != nil {
			print("Trapped a panic in run ", r, "\n")
		}
	}()
	req := &syscallmsg.RunRequest{}
	resp := (*syscallmsg.RunResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]export", m, stack, runImpl, req, resp)
}
func export(ctx context.Context, m api.Module, stack []uint64) {
	defer func() {
		if r := recover(); r != nil {
			print("trapped recover in export", r, "\n")
		}
	}()
	req := &syscallmsg.ExportRequest{}
	resp := (*syscallmsg.ExportResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]export", m, stack, exportImpl, req, resp)
}
func returnValue(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("returnValue 0x%x", stack)
}

func require(ctx context.Context, m api.Module, stack []uint64) {
	defer func() {
		if r := recover(); r != nil {
			print("FOUND RECOVERE ", r, "\n")
		}
	}()
	req := &syscallmsg.RequireRequest{}
	resp := (*syscallmsg.RequireResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]require", m, stack, requireImpl, req, resp)

}

func register(ctx context.Context, m api.Module, stack []uint64) {
	defer func() {
		if r := recover(); r != nil {
			print("FOUND RECOVERE ", r, "\n")
		}
	}()
	req := &syscallmsg.RegisterRequest{}
	resp := &syscallmsg.RegisterResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[syscall]register", m, stack, registerImpl, req, resp)

}

func exit(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("exit 0x%x", stack)
	panic("exit called ")
}
