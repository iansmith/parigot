package syscall

import (
	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/apishared"
	"github.com/iansmith/parigot/apishared/id"
	"google.golang.org/protobuf/proto"

	//"github.com/iansmith/parigot/apiwasm"

	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	//"google.golang.org/protobuf/proto"
)

func unwindLenAndPtr(ret uint64) (uint32, uint32) {
	len64 := ret
	len64 >>= 32
	len32 := uint32(len64)
	ptr64 := ret
	ptr64 &= 0xffffffff
	ptr32 := uint32(ptr64)
	return len32, ptr32
}

// clientSide does the marshalling and unmarshalling needed to read the T given,
// write the U given, and return the KernelErrId properly. It does these
// manipulations so you can call a lower level function that is implemented by
// the host.
func clientSide[T proto.Message, U proto.Message](t T, u U, fn func(int32, int32, int32, int32)) (U, id.KernelErrId) {
	var outErr id.KernelErrId
	outProtoPtr := u
	outErrPtr := &outErr
	var nilU U
	buf, err := proto.Marshal(t)
	if err != nil {
		return nilU, id.NewKernelErrId(id.KernelMarshalFailed)
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
	fn(length, req, out, errPtr)
	if outErr.IsError() {
		return nilU, outErr
	}
	if err := proto.Unmarshal(outBuf, u); err != nil {
		outErr = id.NewKernelErrId(id.KernelErrIdUnmarshalError)
	}
	return outProtoPtr, id.KernelErrIdNoErr
}

// Locate is the means of aquiring a handle to a particular service.
// Most users will not want this interface, but rather will use the
// auto generated method LocateFooOrPanic() for getting an initial
// handle the Foo service.
//
// func Locate(*syscallmsg.LocateRequest) *syscallmsg.LocateResponse
//
//go:wasmimport parigot locate_
func Locate_(int32, int32, int32, int32)
func Locate(inPtr *syscallmsg.LocateRequest) (*syscallmsg.LocateResponse, id.KernelErrId) {
	outProtoPtr := &syscallmsg.LocateResponse{}
	return clientSide(inPtr, outProtoPtr, Locate_)
}

// Dispatch is the primary means that a caller can send an RPC message.
// If you are in local development mode, this call is handled by the kernel
// itself, otherwise it implies a remote procedure call.  This method
// checks the returned response for errors. If there are errors inside the
// result they are pulled out and returned in the error parameter.  Thus
// if the error parameter is nil, the Dispatch() occurred successfully.
// This is code that runs on the WASM side.
//
//go:wasmimport parigot dispatch_
func Dispatch_(int32, int32) int32
func Dispatch(in *syscallmsg.DispatchRequest) (*syscallmsg.DispatchResponse, id.KernelErrId) {
	out := &syscallmsg.DispatchResponse{}
	// err := error(nil)
	// if err != nil {
	// 	return nil, err
	// }
	return out, id.KernelErrIdNoErr
}

// BlockUntilCall is used to block a process until a request is received from another process.  Even when
// all the "processes" are in a single process for debugging, the BlockUntilCall is for the same purpose.
//
// func BlockUntilCall(*syscallmsg.BlockUntilCallRequest) *syscallmsg.BlockUntilCallResponse
//
//xxxgo:wasm-module parigot
//xxxgo:export blockUntilCall
//go:wasmimport parigot block_until_call_
func BlockUntilCall_(int32, int32) int32

func BlockUntilCall(in *syscallmsg.BlockUntilCallRequest) (*syscallmsg.BlockUntilCallResponse, error) {
	out := &syscallmsg.BlockUntilCallResponse{}
	err := error(nil)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BindMethod is the way that a particular service gets associated with
// a given method id. This is normally not needed by user code because the
// generated code for any service will call this automatically.
//
// func BindMethod(*syscallmsg.BindMethodRequest) *syscallmsg.BindMethodResponse
//
//go:wasmimport parigot bind_method_
func BindMethod_(int32, int32) int32

func BindMethod(in *syscallmsg.BindMethodRequest) (*syscallmsg.BindMethodResponse, error) {
	out := &syscallmsg.BindMethodResponse{}
	err := error(nil)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Run is a request to start running. Note that this may not return
// immediately and may fail entirely.  For most user code this is not
// used because user code usually uses file.WaitFileServiceOrPanic() to
// block service File until it is cleared to run.
//
// func Run(*syscallmsg.RunRequest) *syscallmsg.RunResponse
//
//go:wasmimport parigot run_
func Run_(int32, int32) int32
func Run(in *syscallmsg.RunRequest) (*syscallmsg.RunResponse, id.IdRaw) {
	out := &syscallmsg.RunResponse{}
	// err := error(nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("Run_ failed:%v", err)
	// }
	return out, id.KernelErrIdNoErr.Raw()
}

// Export is a declaration that a service implements a particular interface.
// This is not needed by most user code that will use queue.ExportQueueServiceOrPanic()
// to export itself as the queue service.
//
//go:wasmimport parigot export_
func Export_(int32, int32) int32
func Export(in *syscallmsg.ExportRequest) (*syscallmsg.ExportResponse, id.Id) {
	out := &syscallmsg.ExportResponse{}
	// err := error(nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("Export_ failed:%v", err)
	// }
	return out, nil
}

// ReturnValue is not a call that user code should be using. It is the
// mechanism for a return value to be communicated back to the caller
// from the caller.  User code will typically use the wrappers around
// this that make the method calls looking synchronous.
//
//go:wasmimport parigot return_value_
func ReturnValue_(int32, int32) int32
func ReturnValue(in *syscallmsg.ReturnValueRequest) (*syscallmsg.ReturnValueResponse, id.Id) {
	out := &syscallmsg.ReturnValueResponse{}
	// err := error(nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("ReturnValue_ failed:%v", err)
	// }
	return out, nil
}

// Require is a declaration that a service needs a particular interface.
// This is not needed by most user code that will use queue.ImpleQueueServiceOrPanic()
// to import the queue service.
//
//go:wasmimport parigot require_
func Require_(int32, int32, int32, int32)
func Require(inPtr *syscallmsg.RequireRequest) (*syscallmsg.RequireResponse, id.KernelErrId) {
	print("xxxx GOT TO CLIENT SIDE REQUIRE\n")
	outProtoPtr := &syscallmsg.RequireResponse{}
	return clientSide(inPtr, outProtoPtr, Require_)
}

// Exit is called from the WASM side to cause the WASM program to exit.  This is implemented by causing
// the WASM code to panic and then using recover to catch it and then the program is stopped and the kernel
// will marke it dead and so forth.
//
//go:wasmimport parigot exit
func Exit_(int32, int32) int32

func Exit(in *syscallmsg.ExitRequest) (*syscallmsg.ExitResponse, id.IdRaw) {
	out := &syscallmsg.ExitResponse{}
	// err := error(nil)
	// if err != nil {
	// 	return nil, //fmt.Errorf("Exit_ failed:%v", err)
	// }
	return out, id.KernelErrIdNoErr.Raw()
}
