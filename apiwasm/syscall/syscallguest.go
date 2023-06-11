package syscall

import (
	"context"
	"fmt"
	"os"

	"github.com/iansmith/parigot/apiplugin"
	"github.com/iansmith/parigot/apishared/id"
	"github.com/iansmith/parigot/apiwasm"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/syscall/v1"
)

// Locate is the means of aquiring a handle to a particular service.
// Most users will not want this interface, but rather will use the
// auto generated method LocateFooOrPanic() for getting an initial
// handle the Foo service.
//
// func Locate(*syscall.LocateRequest) *syscall.LocateResponse
//
//go:wasmimport parigot locate_
func Locate_(int32, int32, int32, int32) int64
func Locate(inPtr *syscall.LocateRequest) (*syscall.LocateResponse, syscall.KernelErr) {
	outProtoPtr := &syscall.LocateResponse{}
	ctx := apiwasm.ManufactureGuestContext("[syscall]Locate")
	defer pcontext.Dump(ctx)

	lr, errIdRaw, signal :=
		apiwasm.ClientSide(ctx, inPtr, outProtoPtr, Locate_)
	kerr := syscall.KernelErr(errIdRaw)
	if signal {
		os.Exit(1)
	}
	if kerr != syscall.KernelErr_NoError {
		return nil, kerr

	}
	return lr, syscall.KernelErr_NoError
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
func Dispatch(in *syscall.DispatchRequest) (*syscall.DispatchResponse, syscall.KernelErr) {
	out := &syscall.DispatchResponse{}
	// err := error(nil)
	// if err != nil {
	// 	return nil, err
	// }
	return out, syscall.KernelErr_NoError
}

// BlockUntilCall is used to block a process until a request is received from another process.  Even when
// all the "processes" are in a single process for debugging, the BlockUntilCall is for the same purpose.
//
// func BlockUntilCall(*syscall.BlockUntilCallRequest) *syscall.BlockUntilCallResponse
//
//go:wasmimport parigot block_until_call_
func BlockUntilCall_(int32, int32) int32

func BlockUntilCall(in *syscall.BlockUntilCallRequest) (*syscall.BlockUntilCallResponse, error) {
	out := &syscall.BlockUntilCallResponse{}
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
// func BindMethod(*syscall.BindMethodRequest) *syscall.BindMethodResponse
//
//go:wasmimport parigot bind_method_
func BindMethod_(int32, int32) int32

func BindMethod(in *syscall.BindMethodRequest) (*syscall.BindMethodResponse, error) {
	out := &syscall.BindMethodResponse{}
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
func Run(inPtr *syscall.RunRequest) (*syscall.RunResponse, syscall.KernelErr) {
	outProtoPtr := (*syscall.RunResponse)(nil)
	ctx := apiwasm.ManufactureGuestContext("[syscall]Run")
	defer pcontext.Dump(ctx)
	sid := id.UnmarshalServiceId(inPtr.GetServiceId())
	if sid.IsZeroOrEmptyValue() {
		return nil, syscall.KernelErr_BadId
	}
	rr, err, signal :=
		apiwasm.ClientSide(ctx, inPtr, outProtoPtr, Run_)
	if signal {
		os.Exit(1)
	}
	if err != 0 {
		return nil, syscall.KernelErr(err)
	}

	return rr, syscall.KernelErr_NoError
}

// Export is a declaration that a service implements a particular interface.
// This is not needed by most user code that will use queue.ExportQueueServiceOrPanic()
// to export itself as the queue service.
//
//go:wasmimport parigot export_
func Export_(int32, int32, int32, int32) int64
func Export(inPtr *syscall.ExportRequest) (*syscall.ExportResponse, syscall.KernelErr) {
	outProtoPtr := (*syscall.ExportResponse)(nil)
	ctx := apiwasm.ManufactureGuestContext("[syscall]Export")
	defer pcontext.Dump(ctx)
	er, err, signal := apiwasm.ClientSide(ctx, inPtr, outProtoPtr, Export_)
	if signal {
		os.Exit(1)
	}
	return er, syscall.KernelErr(err)
}

// ReturnValue is not a call that user code should be using. It is the
// mechanism for a return value to be communicated back to the caller
// from the caller.  User code will typically use the wrappers around
// this that make the method calls looking synchronous.
//
//go:wasmimport parigot return_value_
// func ReturnValue_(int32, int32) int32
// func ReturnValue(in *syscall.ReturnValueRequest) (*syscall.ReturnValueResponse, id.Id) {
// 	out := &syscall.ReturnValueResponse{}
// 	// err := error(nil)
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("ReturnValue_ failed:%v", err)
// 	// }
// 	return out, nil
// }

// Require is a declaration that a service needs a particular interface.
// This is not needed by most user code that will use queue.ImpleQueueServiceOrPanic()
// to import the queue service.
//
//go:wasmimport parigot require_
func Require_(int32, int32, int32, int32) int64
func Require(inPtr *syscall.RequireRequest) (*syscall.RequireResponse, syscall.KernelErr) {
	outProtoPtr := (*syscall.RequireResponse)(nil)
	ctx := apiwasm.ManufactureGuestContext("[syscall]Require")
	defer pcontext.Dump(ctx)
	rr, err, signal := apiwasm.ClientSide(ctx, inPtr, outProtoPtr, Require_)
	kerr := syscall.KernelErr((err))
	if signal {
		os.Exit(1)
	}

	return rr, kerr
}

// Exit is called from the WASM side to cause the WASM program to exit.  This is implemented by causing
// the WASM code to panic and then using recover to catch it and then the program is stopped and the kernel
// will marke it dead and so forth.
//
//go:wasmimport parigot exit
func Exit_(int32, int32) int32
func Exit(in *syscall.ExitRequest) (*syscall.ExitResponse, syscall.KernelErr) {
	out := &syscall.ExitResponse{}
	// err := error(nil)
	// if err != nil {
	// 	return nil, //fmt.Errorf("Exit_ failed:%v", err)
	// }
	return out, syscall.KernelErr_NoError
}

// Register should be called before any other services are
// Required, Exported, or Located.
//
//go:wasmimport parigot register_
func Register_(int32, int32, int32, int32) int64

func Register(inPtr *syscall.RegisterRequest) (*syscall.RegisterResponse, syscall.KernelErr) {
	outProtoPtr := &syscall.RegisterResponse{}
	ctx := apiplugin.ManufactureHostContext(context.Background(), "[syscall]Register")
	defer pcontext.Dump(ctx)
	rr, kid, signal := apiwasm.ClientSide(ctx, inPtr, outProtoPtr, Register_)
	if signal {
		os.Exit(1)
	}
	return rr, syscall.KernelErr(kid)
}

// MustSatisfyWait is a convenience wrapper around creating a RunRequest and
// using the Run syscall.  MustSatisfyWait is a better name for what goes on
// in the course of a Run() call.
func MustSatisfyWait(ctx context.Context, sid id.ServiceId) {
	req := &syscall.RunRequest{
		Wait:      true,
		ServiceId: sid.Marshal(),
	}
	pcontext.Debugf(ctx, "about to call satisy wait .............. %s", sid.Short())
	_, err := Run(req)
	if err != 0 {
		panic(fmt.Sprintf("failed to run successfully:%s",
			syscall.KernelErr_name[int32(err)]))
	}
}
