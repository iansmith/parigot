package plugin

import (
	"context"
	"encoding/binary"
	"log"
	"runtime/debug"

	apishared "github.com/iansmith/parigot/api/shared"
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
		if !m.WriteByte(offset, fourByte[i]) {
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
// them in the guest side memory.  The first return value is true only if
// the address provided on the stack is out of bounds or the
// address plus the data that is written is out of bounds.  The second
// error return indicates that a system call performing an
// exit has been called.  All
// other errors are returned in-band as kernel errors.
func pushResponseToStack(ctx context.Context, m api.Module, resp proto.Message, respErr int32, stack []uint64) bool {
	log.Printf("pushResponseToStack: %T, %v", resp, resp != nil)
	errPtr := eng.Util.DecodeU32(stack[3])
	if respErr&0x7fffff00 != 0 {
		errCopy := respErr
		respErr = 0
		if !writeErr32Guest(m.Memory(), errPtr, errCopy) {
			return true
		}
	} else {
		if !writeErr32Guest(m.Memory(), errPtr, respErr) {
			return true
		}
	}
	if respErr != 0 {
		stack[0] = NoReturnedStruct
		return false
	}

	buf, err := proto.Marshal(resp)
	if err != nil {
		if writeFailed := writeKernelErrToGuest(m.Memory(), errPtr, syscall.KernelErr_MarshalFailed); writeFailed {
			return true
		}
		stack[0] = NoReturnedStruct
		return false
	}
	if len(buf) > apishared.GuestReceiveBufferSize {
		writeFailed := writeKernelErrToGuest(m.Memory(), errPtr, syscall.KernelErr_DataTooLarge)
		return writeFailed
	}
	outPtr := eng.Util.DecodeU32(stack[2])
	if !m.Memory().Write(outPtr, buf) {
		return true
	}
	stack[0] = windUpLenAndPtr(uint32(len(buf)), outPtr)
	return false
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

	defer func() {
		if r := recover(); r != nil {
			log.Printf("host side panic from inside function %s: %v (%T)", name, r, t)
			debug.PrintStack()
		}
	}()
	if !pullRequestFromStack(ctx, m, t, stack) { //consumes 0 and 1, 3 of stack
		return
	}
	kerr := fn(ctx, t, u)
	badWrite := pushResponseToStack(ctx, m, u, int32(kerr), stack)
	if badWrite {
		panic("unable to push response back to guest memory")
	}
}

func HostBase[T proto.Message, U proto.Message](ctx context.Context, fnName string,
	fn func(context.Context, T, U) int32, m api.Module, stack []uint64, req T, resp U) {
	defer func() {
		if r := recover(); r != nil {
			print(">>>>>>>> Trapped recover in set up for   ", fnName, "<<<<<<<<<<\n")
		}
	}()
	InvokeImplFromStack(ctx, fnName, m, stack, fn, req, resp)
}
