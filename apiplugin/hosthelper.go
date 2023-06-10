package apiplugin

import (
	"context"
	"encoding/binary"
	"runtime/debug"

	"github.com/iansmith/parigot/apishared"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/eng"
	syscall "github.com/iansmith/parigot/g/syscall/v1"

	"github.com/tetratelabs/wazero/api"
	"google.golang.org/protobuf/proto"
)

var NoReturnedStruct = uint64(0)

// windUpLenAndPtr takes two uint32s and makes a single
// uint64 from them.  It makes the length the upper 32 bits.
// This combined value is used when passing a response between the
// host and guest.
func windUpLenAndPtr(length, ptr uint32) uint64 {
	ptr64 := uint64(ptr)
	length64 := uint64(length)
	length64 <<= 32
	if ptr < 0xff {
		debug.PrintStack()
	}
	return length64 | ptr64
}

// writeErr32Guest takes a 32 bit error value and writes it
// into the offset provided. This returns true only if has
// an out of bounds write.
func writeErr32Guest(m api.Memory, offset uint32, err int32) bool {
	fourByte := make([]byte, 4)
	binary.LittleEndian.PutUint32(fourByte, uint32(err))
	for i := 0; i < 4; i++ {
		if !m.WriteByte(offset, fourByte[4-i]) {
			return false
		}
		offset++
	}
	return true
}

// writeKernelErrorToGuest is a wrapper around writeErr32Guest
// that take the specific type of a KernelErr as the error
// value to be written.
func writeKernelErrToGuest(m api.Memory, offset uint32, kerr syscall.KernelErr) bool {
	return writeErr32Guest(m, offset, int32(kerr))
}

// pushResponseToStack does the work of taking a response and a possible error and placing
// them in the guest side memory.  This returns true only if
// the address provided on the stack is out of bounds are the
// address plus the data that is written is out of bounds.  All
// other errors are returned in-band as kernel errors.
func pushResponseToStack(ctx context.Context, m api.Module, resp proto.Message, respErr int32, stack []uint64) bool {
	errPtr := eng.Util.DecodeU32(stack[3])
	if !writeErr32Guest(m.Memory(), errPtr, respErr) {
		return false
	}
	if respErr != 0 {
		stack[0] = NoReturnedStruct
		return true
	}

	buf, err := proto.Marshal(resp)
	if err != nil {
		if !writeKernelErrToGuest(m.Memory(), errPtr, syscall.KernelErr_MarshalFailed) {
			return false
		}
		stack[0] = NoReturnedStruct
		return true
	}
	if len(buf) > apishared.GuestReceiveBufferSize {
		if !writeKernelErrToGuest(m.Memory(), errPtr, syscall.KernelErr_DataTooLarge) {
			return false
		}
		return true
	}
	outPtr := eng.Util.DecodeU32(stack[2])
	if !m.Memory().Write(outPtr, buf) {
		return false
	}
	stack[0] = windUpLenAndPtr(uint32(len(buf)), outPtr)
	return true
}

// pullRequest from stack uses the stack provided to pull the
// parameters from a guest to the host side.  The values in the
// stack are (in order from 0) the length of the input,
// a pointer to the input, ignored value, and
// a ptr to where the return error should be placed.
// All the pointers given are in the guest address space.
// This function only returns true if there was an attempt
// to read a memory address from the client that is out of
// bounds.  Other errors are returned in band.
func pullRequestFromStack(ctx context.Context, m api.Module, req proto.Message, stack []uint64) bool {
	inLen := eng.Util.DecodeU32(stack[0])
	inPtr := eng.Util.DecodeU32(stack[1])
	kerrPtr := eng.Util.DecodeU32(stack[3])
	hostByte, ok := m.Memory().Read(inPtr, inLen)
	if !ok {
		if !writeKernelErrToGuest(m.Memory(), kerrPtr,
			syscall.KernelErr_UnmarshalFailed) {
			return false
		}
		stack[0] = NoReturnedStruct
		return true
	}
	err := proto.Unmarshal(hostByte, req)
	if err != nil {
		if !writeKernelErrToGuest(m.Memory(), kerrPtr, syscall.KernelErr_UnmarshalFailed) {
			return false
		}
		stack[0] = NoReturnedStruct
	}
	return true
}

// InvokeImplFromStack is the primary interface between
// host code and the machinery to communicate with the guest.
// This function takes the parameters like an api.Module and
// a section of the stack provided by Wazero and then reads and
// writes an input and output protocol buffer to the guest
// memory.   The name is provided here just for human error
// messages.  This functions uses pullRequestFromStack and
// and pushResponseToStack to do the actual work of encoding
// and decoding the values to/from the guest memory.
func InvokeImplFromStack[T proto.Message, U proto.Message](ctx context.Context, name string, m api.Module, stack []uint64,
	fn func(context.Context, T, U) int32, t T, u U) {
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

// ManufactureHostContext is a helper to return a context
// configured for the given function name and set the source to
// ServerGo.
func ManufactureHostContext(ctx context.Context, funcName string) context.Context {
	return pcontext.CallTo(pcontext.ServerGoContext(pcontext.NewContextWithContainer(ctx, "ManufactureHostContext")), funcName)
}
