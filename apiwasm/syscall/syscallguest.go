package syscall

import (
	"context"
	"os"

	"github.com/iansmith/parigot/apiplugin"
	"github.com/iansmith/parigot/apishared/id"
	"github.com/iansmith/parigot/apiwasm"
	pcontext "github.com/iansmith/parigot/context"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
)

// Locate is the means of aquiring a handle to a particular service.
// Most users will not want this interface, but rather will use the
// auto generated method LocateFooOrPanic() for getting an initial
// handle the Foo service.
//
// func Locate(*syscallmsg.LocateRequest) *syscallmsg.LocateResponse
//
//go:wasmimport parigot locate_
func Locate_(int32, int32, int32, int32) int64
func Locate(inPtr *syscallmsg.LocateRequest) (*syscallmsg.LocateResponse, id.KernelErrId) {
	outProtoPtr := &syscallmsg.LocateResponse{}
	ctx := apiwasm.ManufactureGuestContext("[syscall]Locate")
	defer pcontext.Dump(ctx)

	lr, errIdRaw, signal := apiwasm.ClientSide(ctx, inPtr, outProtoPtr, Locate_)
	if signal {
		os.Exit(1)
	}
	kerr := id.KernelErrId(errIdRaw)
	return lr, kerr
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

// Run is starts a service (or a guest application) running. Note that
// this may not return immediately and may fail entirely.  For most user
// code this is not used because user code usually uses file.MustFileServiceRun() to
// block service File until it is cleared to run.
//
//go:wasmimport parigot run_
func Run_(int32, int32, int32, int32) int64
func Run(inPtr *syscallmsg.RunRequest) (*syscallmsg.RunResponse, id.KernelErrId) {
	outProtoPtr := (*syscallmsg.RunResponse)(nil)
	ctx := apiwasm.ManufactureGuestContext("[syscall]Run")
	defer pcontext.Dump(ctx)
	rr, err, signal := apiwasm.ClientSide[*syscallmsg.RunRequest, *syscallmsg.RunResponse](ctx, inPtr, outProtoPtr, Run_)
	if signal {
		os.Exit(1)
	}

	return rr, id.KernelErrId(err)
}

// Export is a declaration that a service implements a particular interface.
// This is not needed by most user code that will use queue.ExportQueueServiceOrPanic()
// to export itself as the queue service.
//
//go:wasmimport parigot export_
func Export_(int32, int32, int32, int32) int64
func Export(inPtr *syscallmsg.ExportRequest) (*syscallmsg.ExportResponse, id.KernelErrId) {
	outProtoPtr := (*syscallmsg.ExportResponse)(nil)
	ctx := apiwasm.ManufactureGuestContext("[syscall]Export")
	defer pcontext.Dump(ctx)
	er, err, signal := apiwasm.ClientSide(ctx, inPtr, outProtoPtr, Require_)
	if signal {
		os.Exit(1)
	}
	return er, id.KernelErrId(err)
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
func Require_(int32, int32, int32, int32) int64
func Require(inPtr *syscallmsg.RequireRequest) (*syscallmsg.RequireResponse, id.KernelErrId) {
	outProtoPtr := (*syscallmsg.RequireResponse)(nil)
	ctx := apiwasm.ManufactureGuestContext("[syscall]Require")
	defer pcontext.Dump(ctx)
	rr, err, signal := apiwasm.ClientSide(ctx, inPtr, outProtoPtr, Require_)
	if signal {
		os.Exit(1)
	}
	return rr, id.KernelErrId((err))
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

// Register should be called before any other services are
// Required, Exported, or Located.
//
//go:wasmimport parigot register_
func Register_(int32, int32, int32, int32) int64

func Register(inPtr *syscallmsg.RegisterRequest) (*syscallmsg.RegisterResponse, id.KernelErrId) {
	outProtoPtr := &syscallmsg.RegisterResponse{}
	ctx := apiplugin.ManufactureHostContext(context.Background(), "[syscall]Register")
	defer pcontext.Dump(ctx)
	rr, kid, signal := apiwasm.ClientSide(ctx, inPtr, outProtoPtr, Register_)
	if signal {
		os.Exit(1)
	}
	return rr, id.KernelErrId(kid)
}
