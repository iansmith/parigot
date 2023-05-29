package main

import (
	"context"
	"log"

	"github.com/iansmith/parigot/apishared"
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
	e.AddSupportedFunc(ctx, "parigot", "exit_", exit)
	e.AddSupportedFunc(ctx, "parigot", "exit_", exit)

	return true
}

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
	log.Printf("run 0x%x", stack)
}
func export(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("export 0x%x", stack)
}
func returnValue(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("returnValue 0x%x", stack)
}

func requireImpl(ctx context.Context, req *syscallmsg.RequireRequest, resp *syscallmsg.RequireResponse) id.KernelErrId {
	pcontext.Debugf(ctx, "got a call to require impl")
	fqn := req.GetService()
	for _, fullyQualified := range fqn {
		_, ok := SData.SetService(fullyQualified.GetPackagePath(), fullyQualified.GetService())
		if !ok {
			pcontext.Errorf(ctx, "expected to be able to establish service %s.%s but failed",
				fullyQualified.GetPackagePath(), fullyQualified.GetService())
		} else {
			pcontext.Infof(ctx, "created service %s.%s", fullyQualified.GetPackagePath(), fullyQualified.GetService())
		}
	}
	return id.KernelErrIdNoErr
}

func require(ctx context.Context, m api.Module, stack []uint64) {
	pcontext.Debugf(ctx, "xxx -- required received %04x, %04x, %04x, %04x", stack[0], stack[1], stack[2], stack[3])
	req := &syscallmsg.RequireRequest{}
	resp := (*syscallmsg.RequireResponse)(nil)
	invokeImplFromStack(ctx, "[syscall]require", m, stack, requireImpl, req, resp)
	return
}

func exit(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("exit 0x%x", stack)
	panic("exit called ")
}

func readStringFromGuest(mem api.Memory, nameOffset int32) string {
	l, ok := mem.ReadUint32Le(uint32(nameOffset))
	if !ok {
		panic("unable to read the length of a string from the guest")
	}
	data := uint32(nameOffset + apishared.WasmWidth)
	ptr, ok := mem.ReadUint32Le(data)
	if !ok {
		panic("unable to read the data pointer of a string from the guest")
	}
	result := make([]byte, int(l))
	for i := uint32(0); i < l; i++ {
		b, ok := mem.ReadByte(ptr + i)
		if !ok {
			panic("unable to read the data of a string from the guest")
		}
		result[int(i)] = b
	}
	return string(result)
}
