package syscall

import (
	"context"
	"log"
	"os"

	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/lib/go/exit"

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

	cid := id.UnmarshalCallId(dr.GetCallId())
	targetHid := id.UnmarshalHostId(dr.GetTargetHostId())
	log.Printf("xxx -- client side completed of dispatch finished (target=%s, current=%s): %+v, err=>%v", targetHid.Short(), CurrentHostId().Short(), cid.Short(), err)
	comp := getCompleter(CurrentHostId(), cid)
	log.Printf("xxx -- client side completed: (%s,%s)-> comp %v", CurrentHostId().Short(), cid.Short(), comp)
	// in band error?
	kerr := syscall.KernelErr(err)
	if kerr != syscall.KernelErr_NoError {
		log.Printf(" dispatch error in client side %s", syscall.KernelErr_name[int32(kerr)])
		return nil, kerr
	}

	// somebody else died?
	if signal {
		log.Printf("xxx Dispatch exiting due to signal")
		os.Exit(1)
	}

	return dr, syscall.KernelErr_NoError
}

// Run is starts a service (or a guest application) running. Note that
// this may not return immediately and may fail entirely.  For most user
// code this is not used because user code usually uses file.MustFileServiceRun() to
// block service File until it is cleared to run.
//
//go:wasmimport parigot launch_
func Launch_(int32, int32, int32, int32) int64
func Launch(inPtr *syscall.LaunchRequest) *LaunchFuture {
	outProtoPtr := (*syscall.LaunchResponse)(nil)
	ctx := ManufactureGuestContext("[syscall]Launch")
	defer pcontext.Dump(ctx)
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
	_, err, signal :=
		ClientSide(ctx, inPtr, outProtoPtr, Launch_)
	if signal {
		log.Printf("xxx Launch exiting due to signal")
		os.Exit(1)
	}
	if err != 0 {
		m := NewLaunchFuture()
		m.CompleteMethod(ctx, nil, syscall.KernelErr(err))
		return m
	}
	fut := NewLaunchFuture()
	comp := NewLaunchCompleter(fut)
	MatchCompleter(hid, cid, comp)
	return fut
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

// Exit is called from the WASM side to cause the WASM program, or all the WASM
// programs, to exit.  The future the future is called when the exit is recognized,
// but it is not called when the actual shutdown occurs.  The future given here
// is called when the Exit() itself as has been completed.  For something run
// just before the program stops, use AtExit.
//
//go:wasmimport parigot exit_
func Exit_(int32, int32, int32, int32) int64
func Exit(exitReq *syscall.ExitRequest) *ExitFuture {
	//	ctx := ManufactureGuestContext("[syscall]Exit")

	hid := CurrentHostId()
	cid := id.NewCallId()
	mid := apishared.ExitMethod

	exitReq.CallId = cid.Marshal()
	exitReq.HostId = hid.Marshal()
	exitReq.MethodId = mid.Marshal()

	outProtoPtr := &syscall.ExitResponse{}
	errResp := standardGuestSide(exitReq, outProtoPtr, Exit_, "Exit")
	if errResp != 0 {
		ef := NewExitFuture()
		ef.fut.CompleteMethod(context.Background(), nil, errResp)
		return ef
	}
	ef := NewExitFuture()
	comp := NewExitCompleter(ef)
	MatchCompleter(hid, cid, comp)
	return ef
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

// SynchronousExit is a request that is sent to a service to tell the service it will
// exit shortly (order of milliseconds) and resources should be cleaned up.
// Note that this can happen when another service actually made the Exit() call.
//
//go:wasmimport parigot synchronous_exit_
func SynchronousExit_(int32, int32, int32, int32) int64
func SynchronousExit(in *syscall.SynchronousExitRequest) (*syscall.SynchronousExitResponse, syscall.KernelErr) {
	out := &syscall.SynchronousExitResponse{}

	ctx := pcontext.NewContextWithContainer(context.Background(), "SynchronousExit")
	if standardGuestSide(in, out, SynchronousExit_, "SyncExit") != syscall.KernelErr_NoError {
		pcontext.Errorf(ctx, "unable to exit cleanly, aborting")
		pcontext.Dump(ctx)
		os.Exit(1)
	}
	exit.ExecuteAtExit(ctx)
	panic(apishared.ControlledExit)
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
