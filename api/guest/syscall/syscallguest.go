package syscall

import (
	"context"
	"time"

	"github.com/iansmith/parigot/api/guest"
	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/syscall/v1"

	"google.golang.org/protobuf/proto"
)

//go:wasmimport parigot exit_code_
var _exitCode uint16

// Locate is the means of aquiring a handle to a particular service.
// Most users will not want this interface, but rather will use the
// auto generated method LocateFooOrPanic() for getting an initial
// handle the Foo service.
//
// func Locate(*syscall.LocateRequest) *syscall.LocateResponse
//
//go:wasmimport parigot locate_
func Locate_(int32, int32, int32, int32) int64
func Locate(ctx context.Context, inPtr *syscall.LocateRequest) (*syscall.LocateResponse, syscall.KernelErr) {
	outProtoPtr := &syscall.LocateResponse{}

	lr, errIdRaw, _ :=
		ClientSide(ctx, inPtr, outProtoPtr, Locate_)

	kerr := syscall.KernelErr(errIdRaw)
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
func Dispatch(ctx context.Context, inPtr *syscall.DispatchRequest) (*syscall.DispatchResponse, syscall.KernelErr) {
	outProtoPtr := &syscall.DispatchResponse{}

	dr, err, _ :=
		ClientSide(ctx, inPtr, outProtoPtr, Dispatch_)

	// in band error?
	kerr := syscall.KernelErr(err)
	if kerr != syscall.KernelErr_NoError {
		guest.Log(ctx).Error("dispatch error in client side", "kernel error", syscall.KernelErr_name[int32(kerr)])
		return nil, kerr
	}

	return dr, syscall.KernelErr_NoError
}

// xxx fixme(iansmith) this is horrible.
// we have a loop in our deps, and this breaks it.
func copyOfCurrentTime(ctx context.Context) time.Time {
	t := ctx.Value("parigot_time")
	if t != nil && !t.(time.Time).IsZero() {
		if localTimeZone != nil {
			return t.(time.Time).In(localTimeZone)
		}
		return t.(time.Time)
	}
	if localTimeZone != nil {
		return time.Now().In(localTimeZone)
	}
	return time.Now()
}

// Run is starts a service (or a guest application) running. Note that
// this may not return immediately and may fail entirely.  For most user
// code this is not used because user code usually uses file.MustFileServiceRun() to
// block service File until it is cleared to run.
//
//go:wasmimport parigot launch_
func Launch_(int32, int32, int32, int32) int64
func Launch(ctx context.Context, inPtr *syscall.LaunchRequest) *LaunchFuture {
	outProtoPtr := (*syscall.LaunchResponse)(nil)
	sid := id.UnmarshalServiceId(inPtr.GetServiceId())
	hid := id.UnmarshalHostId(inPtr.GetHostId())
	cid := id.UnmarshalCallId(inPtr.GetCallId())
	mid := id.UnmarshalMethodId(inPtr.GetMethodId())

	if sid.IsZeroOrEmptyValue() || hid.IsZeroOrEmptyValue() || cid.IsZeroOrEmptyValue() || mid.IsZeroOrEmptyValue() {
		lf := NewLaunchFuture()
		lf.fut.CompleteMethod(ctx, nil, syscall.KernelErr_BadId)
	}

	inPtr.CallId = cid.Marshal()
	inPtr.HostId = hid.Marshal()
	_, err, _ :=
		ClientSide(ctx, inPtr, outProtoPtr, Launch_)
	if err != 0 {
		m := NewLaunchFuture()
		m.CompleteMethod(ctx, nil, syscall.KernelErr(err))
		return m
	}
	fut := NewLaunchFuture()
	comp := NewLaunchCompleter(fut)
	MatchCompleter(ctx, copyOfCurrentTime(ctx), hid, cid, comp)
	return fut
}

// Export is a declaration that a service implements a particular interface.
// This is not needed by most user code that will use queue.ExportQueueServiceOrPanic()
// to export itself as the queue service.
//
//go:wasmimport parigot export_
func Export_(int32, int32, int32, int32) int64
func Export(ctx context.Context, inPtr *syscall.ExportRequest) (*syscall.ExportResponse, syscall.KernelErr) {
	outProtoPtr := (*syscall.ExportResponse)(nil)
	return outProtoPtr, standardGuestSide(ctx, inPtr, outProtoPtr, Export_, "Export")
}

// ReturnValue is for providing return values for calls that have
// been made on the local service.
//
//go:wasmimport parigot return_value_
func ReturnValue_(int32, int32, int32, int32) int64
func ReturnValue(ctx context.Context, in *syscall.ReturnValueRequest) syscall.KernelErr {
	out := &syscall.ReturnValueResponse{}
	return standardGuestSide(ctx, in, out, ReturnValue_, "ReturnValue")
}

// Require is a declaration that a service needs a particular interface.
// This is not needed by most user code that will use queue.ImpleQueueServiceOrPanic()
// to import the queue service.
//
//go:wasmimport parigot require_
func Require_(int32, int32, int32, int32) int64
func Require(ctx context.Context, inPtr *syscall.RequireRequest) (*syscall.RequireResponse, syscall.KernelErr) {
	outProtoPtr := (*syscall.RequireResponse)(nil)
	if len(inPtr.GetDest()) == 0 {
		return nil, syscall.KernelErr_NoError
	}
	return outProtoPtr, standardGuestSide(ctx, inPtr, outProtoPtr, Require_, "Require")
}

// Exit is called from the WASM side to cause the WASM program, or all the WASM
// programs, to exit. It does not return.
//
//go:wasmimport parigot exit_
func Exit_(int32, int32, int32, int32) int64
func Exit(ctx context.Context, exitReq *syscall.ExitRequest) *ExitFuture {

	hid := CurrentHostId()
	cid := id.NewCallId()
	mid := apishared.ExitMethod

	exitReq.CallId = cid.Marshal()
	exitReq.HostId = hid.Marshal()
	exitReq.MethodId = mid.Marshal()
	outProtoPtr := &syscall.ExitResponse{}
	standardGuestSide(ctx, exitReq, outProtoPtr, Exit_, "Exit")

	_exitCode = uint16(exitReq.Pair.Code)
	processExit(exitReq.Pair)
	//won't happen
	return nil
}

func processExit(pair *syscall.ExitPair) {
	panic(apishared.ControlledExit)
}

// Register should be called before any other services are
// Required, Exported, or Located.
//
//go:wasmimport parigot register_
func Register_(int32, int32, int32, int32) int64

func Register(ctx context.Context, inPtr *syscall.RegisterRequest) (*syscall.RegisterResponse, syscall.KernelErr) {
	outProtoPtr := &syscall.RegisterResponse{}
	return outProtoPtr, standardGuestSide(ctx, inPtr, outProtoPtr, Register_, "Register")
}

// BindMethod is the way that a particular service gets associated with
// a given method id. This is normally not needed by user code because the
// generated code for any service will call this automatically.
//
// func BindMethod(*syscall.BindMethodRequest) *syscall.BindMethodResponse
//
//go:wasmimport parigot bind_method_
func BindMethod_(int32, int32, int32, int32) int64

func BindMethod(ctx context.Context, in *syscall.BindMethodRequest) (*syscall.BindMethodResponse, syscall.KernelErr) {
	resp := &syscall.BindMethodResponse{}
	return resp, standardGuestSide(ctx, in, resp, BindMethod_, "BindMethod")
}

func MustBindMethodName(ctx context.Context, in *syscall.BindMethodRequest) id.MethodId {
	tmp, kerr := BindMethod(ctx, in)
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
func ReadOne(ctx context.Context, in *syscall.ReadOneRequest) (*syscall.ReadOneResponse, syscall.KernelErr) {
	out := &syscall.ReadOneResponse{}
	return out, standardGuestSide(ctx, in, out, ReadOne_, "ReadOne")
}

// standardGuestSide is a wrapper around ClientSide that knows how
// to handle  error return among return values.
func standardGuestSide[T proto.Message, U proto.Message](ctx context.Context, in T, out U, fn func(int32, int32, int32, int32) int64, name string) syscall.KernelErr {
	_, err, _ := ClientSide(ctx, in, out, fn)
	return syscall.KernelErr(err)
}
