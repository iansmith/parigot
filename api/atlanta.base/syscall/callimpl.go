// This package is a thin wrapper around kernel functionality intended to be run by clients in WASM.
package syscall

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/api/proto/g/pb/call"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	pbsys "github.com/iansmith/parigot/api/proto/g/pb/syscall"
	"github.com/iansmith/parigot/api/splitutil"
	"github.com/iansmith/parigot/lib"

	"google.golang.org/protobuf/proto"
)

// Your IDE may complain about calls to functions in call_js.s and/or calljs.go.  It may claim that these
// are not defined when in fact they are defined, if tricky.    If it really bothers you, most likely you
// can change this be setting the tag "js".  This needs to be defined to get the code in calljs.go.

type callImpl struct {
}

func (l *callImpl) Exit(in *call.ExitRequest) {
	exit(in.Code)
}

// Flip this switch for debug output that starts with libparigot (like libc) and looks like this:
// libparigot:DISPATCH preparing for result return value, 69
// Output from the "BLOCKUNTILCALL" method can get interleaved with the output of the SYSCALL on another
//
// process that is sending a message like this:                  vvvvvvvvvvv
// SYSCALL[DISPATCH,mem-7f0524000000,[proc-9:storeclient.p.wasm]]:libparigot:BLOCKUNTILCALL got result from other process [c-9deb99],[k-000005] with sizes pctx=0,result=0
// params ready (468800,1024) and (468400,1024)SYSCALL[DISPATCH,mem-7f0524000000,[proc-9:storeclient.p.wasm]]:telling the  caller the size of the result and pctx [0,0]
// So that should have been:
// SYSCALL[DISPATCH,mem-7f0524000000,[proc-9:storeclient.p.wasm]]:params ready (468800,1024) and (468400,1024)SYSCALL[DISPATCH,mem-7f0524000000,[proc-9:storeclient.p.wasm]]:telling the  caller the size of the result and pctx [0,0]
// libparigot:BLOCKUNTILCALL got result from other process [c-9deb99],[k-000005] with sizes pctx=0,result=0
// This is because the terminal is not synchronized and these are in different processes (gouroutines).
var libparigotVerbose = true

// Locate is a kernel request that returns either a reference to the service
// or an error.  In the former case, the token returned can be used with Dispatch()
// to make a call on a remote service.  It is implicit in the use of this call that
// the caller wants to be a client of the service in question.  This call can
// be made by clients or servers, but in either case the code in question becomes
// a client of the named service.
//
//go:noinline
func (l *callImpl) Locate(req *pbsys.LocateRequest) (*pbsys.LocateResponse, error) {
	resp := pbsys.LocateResponse{}
	id, err := splitutil.SendReceiveSingleProto(req, &resp, locate)
	if err != nil {
		return nil, err
	}
	if id != nil && id.IsError() {
		// xxx this is bad, swallowing the real error and converting to text
		return nil, lib.NewPerrorFromId("failed to locate properly", id)
	}
	return &resp, nil
}

func stringToTwoInt64s(s string) (int64, int64) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return int64(sh.Data), int64(sh.Len)
}

func sliceToTwoInt64s(b []byte) (int64, int64) {
	slh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return int64(slh.Data), int64(slh.Len)
}

// Dispatch is the primary means that a caller can send an RPC message.
// If you are in local development mode, this call is handled by the kernel
// itself, otherwise it implies a remote procedure call.  This method
// checks the returned response for errors. If there are errors inside the
// result they are pulled out and returned in the error parameter.  Thus
// if the error parameter is nil, the Dispatch() occurred successfully.
// This is code that runs on the WASM side.
func (l *callImpl) Dispatch(req *pbsys.DispatchRequest) (*pbsys.DispatchResponse, error) {
	resp := pbsys.DispatchResponse{}

	libprint("CallImpl.Dispatch", "info from dispatch request: %#v", req)
	id, err := splitutil.SendReceiveSingleProto(req, &resp, dispatch)
	if err != nil {
		return nil, err
	}
	if id != nil {
		if id.IsErrorType() {
			if id.IsError() {
				// xxx this is bad, swallowing the real error and converting to text
				return nil, lib.NewPerrorFromId("failed to dispatch properly", id)
			}
		} else {
			panic(fmt.Sprintf("response to dispatch is unexpected id type: %v, %s", id.IsErrorType(), id.Short()))
		}
	}
	return &resp, nil
}

// BindMethodIn binds a method that only has an in parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func (l *callImpl) BindMethodIn(in *call.BindMethodRequest, _ func(*protosupport.Pctx, proto.Message) error) (*call.BindMethodResponse, error) {
	return l.bindMethodByName(in, call.MethodDirection_METHOD_DIRECTION_IN)
}

// BindMethodInNoPctx binds a method that only has an in parameter and does not
// use the Pctx mechanism for logging.  This may, in fact, be a terrible idea but one
// cannot write a separate logger server with having this.
// xxxfixme: temporary? Should this be a different kernel call?
func (l *callImpl) BindMethodInNoPctx(in *call.BindMethodRequest, _ func(proto.Message) error) (*call.BindMethodResponse, error) {
	return l.bindMethodByName(in, call.MethodDirection_METHOD_DIRECTION_IN)
}

// BindMethodOut binds a method that only has an out parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func (l *callImpl) BindMethodOut(in *call.BindMethodRequest, _ func(*protosupport.Pctx) (proto.Message, error)) (*call.BindMethodResponse, error) {
	return l.bindMethodByName(in, call.MethodDirection_METHOD_DIRECTION_OUT)
}

// BindMethodBoth binds a method that has both an in and out parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func (l *callImpl) BindMethodBoth(in *call.BindMethodRequest, _ func(*protosupport.Pctx, proto.Message) (proto.Message, error)) (*call.BindMethodResponse, error) {
	return l.bindMethodByName(in, call.MethodDirection_METHOD_DIRECTION_BOTH)
}

func (l *callImpl) bindMethodByName(in *call.BindMethodRequest, dir call.MethodDirection) (*call.BindMethodResponse, error) {
	out := new(call.BindMethodResponse)

	out.ErrorId = lib.NoKernelError()

	detail := new(lib.BindPayload)
	sh := (*reflect.StringHeader)(unsafe.Pointer(&in.ProtoPackage))
	detail.PkgPtr = int64(sh.Data)
	detail.PkgLen = int64(sh.Len)
	sh = (*reflect.StringHeader)(unsafe.Pointer(&in.Service))
	detail.ServicePtr = int64(sh.Data)
	detail.ServiceLen = int64(sh.Len)
	sh = (*reflect.StringHeader)(unsafe.Pointer(&in.Method))
	detail.MethodPtr = int64(sh.Data)
	detail.MethodLen = int64(sh.Len)
	detail.Direction = int64(dir)

	// need to allocate the space for result
	out.MethodId = &protosupport.MethodId{}
	out.ErrorId = &protosupport.KernelErrorId{}

	out.MethodId.Id = &protosupport.BaseId{}
	out.ErrorId.Id = &protosupport.BaseId{}

	detail.MethodId = (*[2]int64)(unsafe.Pointer(&out.MethodId.Id.Low))
	detail.ErrorPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Id.Low))

	// THE CALL
	u := uintptr(unsafe.Pointer(detail))
	bindMethod(int32(u))

	// check for in band errors
	kernelErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.ErrorPtr))))
	kerr := lib.NewKernelError(lib.KernelErrorCode(kernelErrDataPtr[0]))
	if kerr.IsError() {
		return nil, lib.NewPerrorFromId("bind error", kerr)
	}

	methodDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.MethodId))))
	mid := lib.NewFrom64BitPair[*protosupport.MethodId](uint64(methodDataPtr[1]), uint64(uint64(methodDataPtr[0])))
	out.MethodId = lib.Marshal[protosupport.MethodId](mid)
	out.ErrorId = lib.NoKernelError()

	return out, nil
}

func (l *callImpl) Run(in *call.RunRequest) (*call.RunResponse, error) {
	out := new(call.RunResponse)

	detail := new(lib.RunPayload)
	detail.Wait = 0
	if in.Wait {
		detail.Wait = 1
	}
	out.ErrorId = &protosupport.KernelErrorId{
		Id: &protosupport.BaseId{},
	}
	detail.KernelErrorPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Id.Low))

	// THE CALL, and the walls came down...
	u := uintptr(unsafe.Pointer(detail))
	run(int32(u))

	kernelErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.KernelErrorPtr))))
	kerr := lib.NewKernelError(lib.KernelErrorCode(kernelErrDataPtr[0]))
	if kerr.IsError() {
		out.ErrorId = lib.Marshal[protosupport.KernelErrorId](kerr)
		return out, lib.NewPerrorFromId("kernel failed to start your process", kerr)
	}
	out.ErrorId = lib.NoKernelError()
	return out, nil
}

func (l *callImpl) exportOrRequire(fqs []*call.FullyQualifiedService, errorPtr *[2]int64, isExport bool) error {
	for _, s := range fqs {
		pkg := s.GetPackagePath()
		svc := s.GetService()

		if isExport {
			libprint("EXPORT", "exporting service %s.%s", pkg, svc)
		} else {
			libprint("REQUIRE", "requiring service %s.%s", pkg, svc)
		}
		detail := new(lib.ExportPayload)
		sh := (*reflect.StringHeader)(unsafe.Pointer(&pkg))
		detail.PkgPtr = int64(sh.Data)
		detail.PkgLen = int64(sh.Len)
		sh = (*reflect.StringHeader)(unsafe.Pointer(&svc))
		detail.ServicePtr = int64(sh.Data)
		detail.ServiceLen = int64(sh.Len)
		detail.KernelErrorPtr = errorPtr //(*[2]int64)(unsafe.Pointer(&out.ErrorId.Low))

		// THE CALL
		u := uintptr(unsafe.Pointer(detail))
		if isExport {
			export(int32(u))
		} else {
			require(int32(u))
		}

		// check for in band errors
		kernelErrDataPtr := detail.KernelErrorPtr //(*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.KernelErrorPtr))))
		kerr := lib.NewKernelError(lib.KernelErrorCode(kernelErrDataPtr[0]))
		if kerr.IsError() {
			if isExport {
				return lib.NewPerrorFromId("export error", kerr)
			} else {
				return lib.NewPerrorFromId("require error", kerr)
			}

		}
	}
	return nil

}

// Export is the way that a server can express that it is done binding methods
// the service and it is ready to export it.  This call does not block.  If the input
// structure has multiple services in it, this call will repeatedly call
// the kernel and it will abort and return the error at the first failure.
func (l *callImpl) Export(in *call.ExportRequest) (*call.ExportResponse, error) {
	out := new(call.ExportResponse)
	// allocate space for any error
	out.ErrorId = lib.NoKernelError()
	ptr := (*[2]int64)(unsafe.Pointer(&out.ErrorId.Id.Low))
	err := l.exportOrRequire(in.GetService(), ptr, true)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Require is the way that a client or server can express that uses a particular
// interface.  This call does not block.  If the input structure has multiple
// services in it, this call will repeatedly call the kernel and it will abort
// and return the error at the first failure.
func (l *callImpl) Require(in *call.RequireRequest) (*call.RequireResponse, error) {
	libprint("REQUIRE ", "request to require %d services", len(in.Service))
	out := new(call.RequireResponse)
	// allocate space for any error
	out.ErrorId = lib.NoKernelError()
	ptr := (*[2]int64)(unsafe.Pointer(&out.ErrorId.Id.Low))
	err := l.exportOrRequire(in.GetService(), ptr, false)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// xxx this may be a bad idea.  this is probably only temporary til I can work out if we
// xxx should support this at all. currently, it is used by terminal logger.  I made a copy
// xxx of the other bindMethodByName so it would be easier to delete this one.
func (l *callImpl) bindMethodByNameNoPctx(in *call.BindMethodRequest, dir call.MethodDirection) (*call.BindMethodResponse, error) {
	out := new(call.BindMethodResponse)

	out.ErrorId = lib.NoKernelError()

	detail := new(lib.BindPayload)
	sh := (*reflect.StringHeader)(unsafe.Pointer(&in.ProtoPackage))
	detail.PkgPtr = int64(sh.Data)
	detail.PkgLen = int64(sh.Len)
	sh = (*reflect.StringHeader)(unsafe.Pointer(&in.Service))
	detail.ServicePtr = int64(sh.Data)
	detail.ServiceLen = int64(sh.Len)
	sh = (*reflect.StringHeader)(unsafe.Pointer(&in.Method))
	detail.MethodPtr = int64(sh.Data)
	detail.MethodLen = int64(sh.Len)
	detail.Direction = int64(dir)
	detail.MethodId = (*[2]int64)(unsafe.Pointer(&out.MethodId.Id.Low))
	detail.ErrorPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Id.Low))

	// THE CALL
	u := uintptr(unsafe.Pointer(detail))
	bindMethod(int32(u))

	// check for in band errors
	kernelErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.ErrorPtr))))
	kerr := lib.NewKernelError(lib.KernelErrorCode(kernelErrDataPtr[0]))
	if kerr.IsError() {
		return nil, lib.NewPerrorFromId("bind error", kerr)
	}

	methodDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.MethodId))))
	mid := lib.NewFrom64BitPair[*protosupport.MethodId](uint64(methodDataPtr[1]), uint64(uint64(methodDataPtr[0])))
	out.MethodId = lib.Marshal[protosupport.MethodId](mid)

	out.ErrorId = lib.NoKernelError()

	return out, nil
}

func libprint(call, format string, arg ...interface{}) {
	if libparigotVerbose {
		part1 := fmt.Sprintf("libparigot:%s", call)
		part2 := fmt.Sprintf(format, arg...)
		print(part1, " ", part2, "\n")
	}
}

func (l *callImpl) BlockUntilCall(in *call.BlockUntilCallRequest) (*call.BlockUntilCallResponse, error) {
	// this is JUST for reserving the space for the result to be placed into
	out := &call.BlockUntilCallResponse{
		Method:  lib.NoErrorMarshaled[protosupport.MethodId, *protosupport.MethodId](),
		Call:    lib.NoErrorMarshaled[protosupport.CallId, *protosupport.CallId](),
		ErrorId: lib.NoKernelError(),
	}

	payload := &lib.BlockPayload{}

	if len(in.PctxBuffer) > 0 {
		payload.PctxPtr, payload.PctxLen = sliceToTwoInt64s(in.PctxBuffer)
	} else {
		payload.PctxPtr = 0
		payload.PctxLen = 0
	}
	payload.ParamPtr, payload.ParamLen = sliceToTwoInt64s(in.ParamBuffer)
	payload.ErrorPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Id.Low))
	payload.MethodId = (*[2]int64)(unsafe.Pointer(&out.Method.Id.Low))
	payload.CallId = (*[2]int64)(unsafe.Pointer(&out.Call.Id.Low))

	// THE CALL
	u := uintptr(unsafe.Pointer(payload))
	blockUntilCall(int32(u))
	libprint("BLOCKUNTILCALL ", "finished call")
	// unpack the result
	kernelErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(payload.ErrorPtr))))
	kerr := lib.NewKernelError(lib.KernelErrorCode(kernelErrDataPtr[0]))
	if kerr.IsError() {
		return nil, lib.NewPerrorFromId("BlockUntilCall error", kerr)
	}

	callDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(payload.CallId))))
	cid := lib.NewFrom64BitPair[*protosupport.CallId](uint64(callDataPtr[1]), uint64(callDataPtr[0]))
	methDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(payload.MethodId))))
	mid := lib.NewFrom64BitPair[*protosupport.MethodId](uint64(methDataPtr[1]), uint64(methDataPtr[0]))
	libprint("BLOCKUNTILCALL ", "mid computed %s", mid.Short())

	out.Call = lib.Marshal[protosupport.CallId](cid)
	out.Method = lib.Marshal[protosupport.MethodId](mid)
	out.ErrorId = lib.NoKernelError()

	// get the data
	out.ParamLen = int32(payload.ParamLen)
	out.PctxLen = int32(payload.PctxLen)

	libprint("BLOCKUNTILCALL ", "unpacked the data %s,%s --- paramlen %d, pctxlen %d", mid.Short(), cid.Short(),
		out.ParamLen, out.PctxLen)
	return out, nil
}

func (l *callImpl) ReturnValue(in *call.ReturnValueRequest) (*call.ReturnValueResponse, error) {
	detail := &lib.ReturnValuePayload{}

	detail.PctxPtr, detail.PctxLen = sliceToTwoInt64s(in.PctxBuffer)
	detail.ResultPtr, detail.ResultLen = sliceToTwoInt64s(in.ResultBuffer)

	libprint("RETURNVALUE ", "buffers for pctx and result sending to kernel are size %d,%d",
		len(in.PctxBuffer), len(in.ResultBuffer))

	detail.MethodId[0] = int64(in.Method.Id.GetLow())
	detail.MethodId[1] = int64(in.Method.Id.GetHigh())
	detail.CallId[0] = int64(in.Call.Id.GetLow())
	detail.CallId[1] = int64(in.Call.Id.GetHigh())
	detail.KernelErrorPtr[0] = int64(in.ErrorId.Id.GetLow())
	detail.KernelErrorPtr[1] = int64(in.ErrorId.Id.GetHigh())

	u := uintptr(unsafe.Pointer(detail))
	returnValue(int32(u))
	// check to see the return value
	kerr := lib.NewKernelError(lib.KernelErrorCode(detail.KernelErrorPtr[0]))
	big := lib.Marshal[protosupport.KernelErrorId](kerr)
	if kerr.IsError() {
		return nil, lib.NewPerrorFromId("failed to process return value", kerr)
	}
	return &call.ReturnValueResponse{
		ErrorId: big,
	}, nil
}

// Export1 is a wrapper around Export which makes it easy to say you export a single
// service. It does not change any of the Export behavior.
func (l *callImpl) Export1(packagePath, service string) (*call.ExportResponse, error) {
	fqSvc := &call.FullyQualifiedService{
		PackagePath: packagePath, Service: service}
	req := &call.ExportRequest{}
	req.Service = []*call.FullyQualifiedService{fqSvc}
	return l.Export(req)
}

// Require1 is a wrapper around Require which makes it easy to say you require a single
// service. It does not change any of the Require behavior.
func (l *callImpl) Require1(packagePath, service string) (*call.RequireResponse, error) {
	fqSvc := &call.FullyQualifiedService{
		PackagePath: packagePath, Service: service}
	req := &call.RequireRequest{}
	req.Service = []*call.FullyQualifiedService{fqSvc}
	return l.Require(req)
}

func NewCallImpl() lib.Call {
	return &callImpl{}
}
