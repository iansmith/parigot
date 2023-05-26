package syscall

import (
	"fmt"

	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"

	"google.golang.org/protobuf/proto"
)

// xxx possibly dead now?
type syscallPtrIn interface {
	proto.Message
	*syscallmsg.LocateRequest |
		*syscallmsg.DispatchRequest |
		*syscallmsg.BlockUntilCallRequest |
		*syscallmsg.BindMethodRequest |
		*syscallmsg.RunRequest |
		*syscallmsg.RequireRequest |
		*syscallmsg.ExportRequest |
		*syscallmsg.ReturnValueRequest |
		*syscallmsg.ExitRequest
}

// xxx possibly dead now?
type syscallPtrOut interface {
	proto.Message
	*syscallmsg.LocateResponse |
		*syscallmsg.DispatchResponse |
		*syscallmsg.BlockUntilCallResponse |
		*syscallmsg.BindMethodResponse |
		*syscallmsg.RunResponse |
		*syscallmsg.RequireResponse |
		*syscallmsg.ExportResponse |
		*syscallmsg.ReturnValueResponse |
		*syscallmsg.ExitResponse
}

// Locate is the means of aquiring a handle to a particular service.
// Most users will not want this interface, but rather will use the
// auto generated method LocateFooOrPanic() for getting an initial
// handle the Foo service.
//
// func Locate(*syscallmsg.LocateRequest) *syscallmsg.LocateResponse
//
//xxxgo:wasm-module parigot
//xxxgo:export locate
//go:wasmimport parigot locate_
func Locate_(int32, int32) int32
func Locate(in *syscallmsg.LocateRequest) (*syscallmsg.LocateResponse, error) {
	out := &syscallmsg.LocateResponse{}
	err := error(nil)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Dispatch is the primary means that a caller can send an RPC message.
// If you are in local development mode, this call is handled by the kernel
// itself, otherwise it implies a remote procedure call.  This method
// checks the returned response for errors. If there are errors inside the
// result they are pulled out and returned in the error parameter.  Thus
// if the error parameter is nil, the Dispatch() occurred successfully.
// This is code that runs on the WASM side.
//
//xxxgo:wasm-module parigot
//xxxgo:export dispatch
//go:wasmimport parigot dispatch_
func Dispatch_(int32, int32) int32
func Dispatch(in *syscallmsg.DispatchRequest) (*syscallmsg.DispatchResponse, error) {
	out := &syscallmsg.DispatchResponse{}
	err := error(nil)
	if err != nil {
		return nil, err
	}
	return out, nil
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
//xxxgo:wasm-module parigot
//xxxgo:export bindMethod
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
//xxxgo:wasm-module parigot
//xxxgo:export run
//go:wasmimport parigot run_
func Run_(int32, int32) int32
func Run(in *syscallmsg.RunRequest) (*syscallmsg.RunResponse, error) {
	out := &syscallmsg.RunResponse{}
	err := error(nil)
	if err != nil {
		return nil, fmt.Errorf("Run_ failed:%v", err)
	}
	return out, nil
}

// Export is a declaration that a service implements a particular interface.
// This is not needed by most user code that will use queue.ExportQueueServiceOrPanic()
// to export itself as the queue service.
//
// func Export(*syscallmsg.ExportRequest) (*syscallmsg.ExportResponse,error)
//
//xxxgo:wasm-module parigot
//xxxgo:export export
//go:wasmimport parigot export_
func Export_(int32, int32) int32
func Export(in *syscallmsg.ExportRequest) (*syscallmsg.ExportResponse, error) {
	out := &syscallmsg.ExportResponse{}
	err := error(nil)
	if err != nil {
		return nil, fmt.Errorf("Export_ failed:%v", err)
	}
	return out, nil
}

// ReturnValue is not a call that user code should be using. It is the
// mechanism for a return value to be communicated back to the caller
// from the caller.  User code will typically use the wrappers around
// this that make the method calls looking synchronous.
//
// func ReturnValue(*syscallmsg.ReturnValueRequest) *syscallmsg.ReturnValueResponse
//
//xxxgo:wasm-module parigot
//xxxgo:export return_value
//go:wasmimport parigot return_value_
func ReturnValue_(int32, int32) int32
func ReturnValue(in *syscallmsg.ReturnValueRequest) (*syscallmsg.ReturnValueResponse, error) {
	out := &syscallmsg.ReturnValueResponse{}
	err := error(nil)
	if err != nil {
		return nil, fmt.Errorf("ReturnValue_ failed:%v", err)
	}
	return out, nil
}

// Require is a declaration that a service needs a particular interface.
// This is not needed by most user code that will use queue.ImpleQueueServiceOrPanic()
// to import the queue service.
//
// func Require(*syscallmsg.RequireRequest) *syscallmsg.RequireResponse
//
//xxxgo:wasm-module parigot
//xxxgo:export require
//go:wasmimport parigot require_
func Require_(int32, int32) int32
func Require(in *syscallmsg.RequireRequest) (*syscallmsg.RequireResponse, error) {
	out := &syscallmsg.RequireResponse{}
	err := error(nil)
	if err != nil {
		return nil, fmt.Errorf("Require_ failed:%v", err)
	}
	return out, nil
}

// Exit is called from the WASM side to cause the WASM program to exit.  This is implemented by causing
// the WASM code to panic and then using recover to catch it and then the program is stopped and the kernel
// will marke it dead and so forth.
//
// func Exit(*syscallmsg.ExitRequest) *syscallmsg.ExitResponse
//
//xxxgo:wasm-module parigot
//xxxgo:export exit
//go:wasmimport parigot exit
func Exit_(int32, int32) int32

func Exit(in *syscallmsg.ExitRequest) (*syscallmsg.ExitResponse, error) {
	out := &syscallmsg.ExitResponse{}
	err := error(nil)
	if err != nil {
		return nil, fmt.Errorf("Exit_ failed:%v", err)
	}
	return out, nil
}

// RegisterExport is how a wasm-implemented, guest function is made avaialble to be called at any time from
// the host side.  Note that this function must be called by the same guest-side goroutine that was
// created for WasmExport.
//
//go:noescape
//go:wasmimport parigot register_export_
func RegisterExport(nameHeader uint32, //reflect.StringHeader
	poolHeader uint32, //*reflect.SliceHeader
	is32Bit uint32,
	buffer uint32, //*reflect.SliceHeader
	exclusiveBufferSizePtr uint32, // *int32
	flagPtr uint32, // *[2]int32
	turnPtr uint32) // *int32
