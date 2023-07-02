package syscall

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/proto"
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
	ctx := ManufactureGuestContext("[syscall]Locate")
	defer pcontext.Dump(ctx)

	lr, errIdRaw, signal :=
		ClientSide(ctx, inPtr, outProtoPtr, Locate_)
	kerr := syscall.KernelErr(errIdRaw)
	if signal {
		pcontext.Fatalf(ctx, "xxxx Locate exiting due to signal")
		pcontext.Dump(ctx)
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
func Dispatch_(int32, int32, int32, int32) int64
func Dispatch(inPtr *syscall.DispatchRequest) (*syscall.DispatchResponse, syscall.KernelErr) {
	outProtoPtr := &syscall.DispatchResponse{}
	ctx := ManufactureGuestContext("[syscall]Dispatch")
	defer pcontext.Dump(ctx)

	dr, err, signal :=
		ClientSide(ctx, inPtr, outProtoPtr, Dispatch_)

	// in band error?
	kerr := syscall.KernelErr(err)
	if kerr != syscall.KernelErr_NoError {
		return nil, kerr
	}

	// somebody else died?
	if signal {
		log.Printf("xxx Dispatch exiting due to signal")
		os.Exit(1)
	}

	// normal case
	return dr, syscall.KernelErr_NoError
}

// Run is starts a service (or a guest application) running. Note that
// this may not return immediately and may fail entirely.  For most user
// code this is not used because user code usually uses file.MustFileServiceRun() to
// block service File until it is cleared to run.
//
//go:wasmimport parigot launch_
func Launch_(int32, int32, int32, int32) int64
func Launch(inPtr *syscall.LaunchRequest) (*syscall.LaunchResponse, syscall.KernelErr) {
	outProtoPtr := (*syscall.LaunchResponse)(nil)
	ctx := ManufactureGuestContext("[syscall]Run")
	defer pcontext.Dump(ctx)
	sid := id.UnmarshalServiceId(inPtr.GetServiceId())
	if sid.IsZeroOrEmptyValue() {
		return nil, syscall.KernelErr_BadId
	}
	rr, err, signal :=
		ClientSide(ctx, inPtr, outProtoPtr, Launch_)
	if signal {
		log.Printf("xxx Run exiting due to signal")
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
	return outProtoPtr, standardGuestSide(inPtr, outProtoPtr, Export_, "Export")
}

// ReturnValue is for providing return values for calls that have
// been made on the local service.
//
//go:wasmimport parigot return_value_
func ReturnValue_(int32, int32, int32, int32) int64
func ReturnValue(in *syscall.ReturnValueRequest) syscall.KernelErr {
	out := &syscall.ReturnValueResponse{}
	return standardGuestSide(in, out, ReturnValue_, "ReturnValue")
}

// Require is a declaration that a service needs a particular interface.
// This is not needed by most user code that will use queue.ImpleQueueServiceOrPanic()
// to import the queue service.
//
//go:wasmimport parigot require_
func Require_(int32, int32, int32, int32) int64
func Require(inPtr *syscall.RequireRequest) (*syscall.RequireResponse, syscall.KernelErr) {
	outProtoPtr := (*syscall.RequireResponse)(nil)
	if len(inPtr.GetDest()) == 0 {
		return nil, syscall.KernelErr_NoError
	}
	return outProtoPtr, standardGuestSide(inPtr, outProtoPtr, Require_, "Require")
}

// Exit is called from the WASM side to cause the WASM program to exit.  This is implemented by causing
// the WASM code to panic and then using recover to catch it and then the program is stopped and the kernel
// will marke it dead and so forth.
//
//go:wasmimport parigot exit_
func Exit_(int32, int32, int32, int32) int64
func Exit(err int32) (*syscall.ExitResponse, syscall.KernelErr) {
	in := &syscall.ExitRequest{
		Code: err,
	}
	out := &syscall.ExitResponse{}
	return out, standardGuestSide(in, out, Exit_, "Exit")
}

// Register should be called before any other services are
// Required, Exported, or Located.
//
//go:wasmimport parigot register_
func Register_(int32, int32, int32, int32) int64

func Register(inPtr *syscall.RegisterRequest) (*syscall.RegisterResponse, syscall.KernelErr) {
	outProtoPtr := &syscall.RegisterResponse{}
	return outProtoPtr, standardGuestSide(inPtr, outProtoPtr, Register_, "Register")
}

// MustSatisfyWait is a convenience wrapper around creating a RunRequest and
// using the Run syscall.  MustSatisfyWait is a better name for what goes on
// in the course of a Run() call.
func MustSatisfyWait(ctx context.Context, sid id.ServiceId) {
	req := &syscall.LaunchRequest{
		ServiceId: sid.Marshal(),
	}
	pcontext.Debugf(ctx, "about to call satisfy wait for .............. %s", sid.Short())
	_, err := Launch(req)
	if err != 0 {
		pcontext.Errorf(ctx, "Run failed of syscall.MustSatisfy wait: %s", syscall.KernelErr_name[int32(err)])
		panic(fmt.Sprintf("failed to run successfully:%s",
			syscall.KernelErr_name[int32(err)]))
	}
}

//
// unused?
//

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
func BindMethod_(int32, int32, int32, int32) int64

func BindMethod(in *syscall.BindMethodRequest) (*syscall.BindMethodResponse, syscall.KernelErr) {
	resp := &syscall.BindMethodResponse{}
	return resp, standardGuestSide(in, resp, BindMethod_, "BindMethod")
}

func MustBindMethodName(in *syscall.BindMethodRequest) id.MethodId {
	tmp, kerr := BindMethod(in)
	if kerr != syscall.KernelErr_NoError {
		panic("failed to bind method:" + in.GetMethodName() + ", error " + syscall.KernelErr_name[int32(kerr)])
	}
	return id.UnmarshalMethodId(tmp.MethodId)
}

// ReadOne checks to see if any of the service/method
// pairs have been called. Timeouts of negative values
// (forever) and 0 (instant check) are legal.
//
//go:wasmimport parigot read_one_
func ReadOne_(int32, int32, int32, int32) int64
func ReadOne(in *syscall.ReadOneRequest) (*syscall.ReadOneResponse, syscall.KernelErr) {
	out := &syscall.ReadOneResponse{}
	return out, standardGuestSide(in, out, ReadOne_, "ReadOne")
}

// standardGuestSide is a wrapper around ClientSide that knows how
// to handle the error return to do an immediate exit.
func standardGuestSide[T proto.Message, U proto.Message](in T, out U, fn func(int32, int32, int32, int32) int64, name string) syscall.KernelErr {
	ctx := ManufactureGuestContext("[guest syscall]" + name)
	defer pcontext.Dump(ctx)
	_, err, signal := ClientSide(ctx, in, out, fn)
	if signal {
		pcontext.Fatalf(ctx, "(syscall guest) %s method exiting due to signal", name)
		os.Exit(1)
	}

	return syscall.KernelErr(err)
}
