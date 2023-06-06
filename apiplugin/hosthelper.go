package apiplugin

import (
	"context"
	"runtime/debug"

	"github.com/iansmith/parigot/apishared"
	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/eng"

	"github.com/tetratelabs/wazero/api"
	"google.golang.org/protobuf/proto"
)

var NoReturnedStruct = uint64(0)

func windUpLenAndPtr(length, ptr uint32) uint64 {
	ptr64 := uint64(ptr)
	length64 := uint64(length)
	length64 <<= 32
	if ptr < 0xff {
		debug.PrintStack()
	}
	return length64 | ptr64
}

func pushResponseToStack(ctx context.Context, m api.Module, resp proto.Message, kerr id.IdRaw, stack []uint64) bool {
	kerrPtr := eng.Util.DecodeU32(stack[3])
	if kerr.IsError() {
		kid := id.NewKernelErrIdFromRaw(kerr)
		kid.MustWriteGuestLe(m.Memory(), kerrPtr)
		stack[0] = NoReturnedStruct
		return true
	}
	buf, err := proto.Marshal(resp)
	if err != nil {
		e := id.NewKernelErrId(id.KernelMarshalFailed)
		e.MustWriteGuestLe(m.Memory(), kerrPtr)
		stack[0] = NoReturnedStruct
		return true
	}
	if len(buf) > apishared.GuestReceiveBufferSize {
		e := id.NewKernelErrId(id.KernelDataTooLarge)
		e.MustWriteGuestLe(m.Memory(), kerrPtr)
		stack[0] = NoReturnedStruct
		return true

	}
	outPtr := eng.Util.DecodeU32(stack[2])
	if !m.Memory().Write(outPtr, buf) {
		return false
	}
	stack[0] = windUpLenAndPtr(uint32(len(buf)), outPtr)
	return true
}

func pullRequestFromStack(ctx context.Context, m api.Module, req proto.Message, stack []uint64) bool {
	inLen := eng.Util.DecodeU32(stack[0])
	inPtr := eng.Util.DecodeU32(stack[1])
	kerrPtr := eng.Util.DecodeU32(stack[3])
	hostByte, ok := m.Memory().Read(inPtr, inLen)
	if !ok {
		kerr := id.NewKernelErrId(id.KernelGuestReadFailed)
		kerr.MustWriteGuestLe(m.Memory(), kerrPtr)
		stack[0] = NoReturnedStruct
		return false
	}
	err := proto.Unmarshal(hostByte, req)
	if err != nil {
		kerr := id.NewKernelErrId(id.KernelErrIdUnmarshalError)
		kerr.MustWriteGuestLe(m.Memory(), kerrPtr)
		stack[0] = NoReturnedStruct
		return false
	}
	return true
}

func InvokeImplFromStack[T proto.Message, U proto.Message](ctx context.Context, name string, m api.Module, stack []uint64,
	fn func(context.Context, T, U) id.IdRaw, t T, u U) {
	currCtx := ManufactureHostContext(ctx, name)
	defer func() {
		if r := recover(); r != nil {
			pcontext.Fatalf(currCtx, "host side panic from inside function %s: %v", name, r)
		}
		pcontext.Dump(currCtx)
	}()
	if !pullRequestFromStack(currCtx, m, t, stack) { //consumes 0 and 1, 3 of stack
		return
	}
	kerr := fn(currCtx, t, u)
	if !pushResponseToStack(currCtx, m, u, kerr, stack) {
		panic("unable to push response back to guest memory")
	}
}

func ManufactureHostContext(ctx context.Context, funcName string) context.Context {
	return pcontext.CallTo(pcontext.ServerGoContext(pcontext.NewContextWithContainer(ctx, "ManufactureHostContext")), funcName)
}
