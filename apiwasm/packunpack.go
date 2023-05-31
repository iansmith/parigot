package apiwasm

import (
	"context"
	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/apishared"
	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"

	"google.golang.org/protobuf/proto"
)

// Manufacture context is used to setup the context for a given state that makes sense for this, the
// Guest side of the wire.  You pass the name of the function you are constructing this in.
func ManufactureGuestContext(fn string) context.Context {
	result := pcontext.NewContextWithContainer(context.Background(), fn)
	result = pcontext.GuestContext(result)
	return pcontext.CallTo(result, fn)
}

// ClientSide does the marshalling and unmarshalling needed to read the T given,
// write the U given, and return the KernelErrId properly. It does these
// manipulations so you can call a lower level function that is implemented by
// the host.
func ClientSide[T proto.Message, U proto.Message](ctx context.Context, t T, u U, fn func(int32, int32, int32, int32) int64) (U, id.IdRaw) {
	var outErr id.IdRaw
	outProtoPtr := u
	outErrPtr := &outErr
	var nilU U

	buf, err := proto.Marshal(t)
	if err != nil {
		return nilU, id.NewKernelErrId(id.KernelMarshalFailed).Raw()
	}
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	length := int32(len(buf))
	req := int32(sh.Data)
	val := reflect.ValueOf(u)
	if val.Kind() != reflect.Ptr {
		panic("client side of syscall passed a proto.Message that is not a pointer")
	}

	outBuf := make([]byte, apishared.GuestReceiveBufferSize)
	sh = (*reflect.SliceHeader)(unsafe.Pointer(&outBuf))
	out := int32(sh.Data)
	errPtr := int32(uintptr(unsafe.Pointer(outErrPtr)))
	wrapped := fn(length, req, out, errPtr)
	if outErr.IsError() {
		return nilU, outErr
	}
	l, ptr := unwindLenAndPtr(uint64(wrapped))

	if int32(ptr) != out { //sanity
		panic("mismatched pointers in host call/return")
	}
	if unsafe.Pointer(asPtr(u)) == nil {
		myId := id.KernelErrIdNoErr.Raw()

		return u, myId
	}
	if err := proto.Unmarshal(outBuf[:l], u); err != nil {
		myId := id.NewKernelErrId((2))
		outErr = myId.Raw()
	}
	return outProtoPtr, id.NewKernelErrId((0)).Raw()
}

func asPtr[T proto.Message](t T) uintptr {
	val := reflect.ValueOf(t)
	if val.Kind() != reflect.Pointer {
		panic("should never call the standard processing of a client side syscall with a value, always use a pointer")
	}
	return val.Pointer()
}

func unwindLenAndPtr(ret uint64) (uint32, uint32) {
	len64 := ret
	len64 >>= 32
	len32 := uint32(len64)
	ptr64 := ret
	ptr64 &= 0xffffffff
	ptr32 := uint32(ptr64)
	return len32, ptr32
}
