// This package is a thin wrapper around kernel functionality. Its primary
// function is to check responses for errors and return them as a go
// error (type Perror).

package lib

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/g/pb/call"
	"github.com/iansmith/parigot/g/pb/log"
	"github.com/iansmith/parigot/g/pb/protosupport"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
var libparigotVerbose = false

// Locate is a kernel request that returns either a reference to the service
// or an error.  In the former case, the token returned can be used with Dispatch()
// to make a call on a remote service.  It is implicit in the use of this call that
// the caller wants to be a client of the service in question.  This call can
// be made by clients or servers, but in either case the code in question becomes
// a client of the named service.
//
//go:noinline
func (l *callImpl) Locate(in *call.LocateRequest) (*call.LocateResponse, error) {
	out := new(call.LocateResponse)
	out.ErrorId = &protosupport.KernelErrorId{High: 6, Low: 7}
	out.ErrorId.High = 1
	out.ErrorId.Low = 2
	out.ServiceId = &protosupport.ServiceId{High: 3, Low: 4}

	detail := new(LocatePayload)
	pkgSh := (*reflect.StringHeader)(unsafe.Pointer(&in.PackageName))
	detail.PkgPtr = int64(pkgSh.Data)
	detail.PkgLen = int64(pkgSh.Len)

	serviceSh := (*reflect.StringHeader)(unsafe.Pointer(&in.ServiceName))
	detail.ServicePtr = int64(serviceSh.Data)
	detail.ServiceLen = int64(serviceSh.Len)
	// choosing low's addr means you ADD 8 to get to the high
	detail.ErrorPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Low))
	detail.ServiceIdPtr = (*[2]int64)(unsafe.Pointer(&out.ServiceId.Low))

	u := uintptr(unsafe.Pointer(detail))
	locate(int32(u))
	// marshal them back together
	svcDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.ServiceIdPtr))))
	sid := ServiceIdFromUint64(uint64(svcDataPtr[1]), uint64(uint64(svcDataPtr[0])))
	out.ServiceId = MarshalServiceId(sid)
	locErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.ErrorPtr))))

	err := NewKernelError(KernelErrorCode(locErrDataPtr[0]))
	// in case the caller walks the structure repacks a new protobuf
	out.ServiceId = MarshalServiceId(sid)
	out.ErrorId = MarshalKernelErrId(err)

	if err.IsError() {
		return out, NewPerrorFromId("failed to locate properly", err)
	}
	return out, nil
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
func (l *callImpl) Dispatch(in *call.DispatchRequest) (*call.DispatchResponse, error) {
	out := new(call.DispatchResponse)

	detail := &DispatchPayload{}
	detail.ServiceId[0] = int64(in.ServiceId.GetLow())
	detail.ServiceId[1] = int64(in.ServiceId.GetHigh())

	detail.MethodPtr, detail.MethodLen = stringToTwoInt64s(in.Method)
	detail.CallerPtr, detail.CallerLen = stringToTwoInt64s(in.Caller)

	if in.GetInPctx() == nil {
		in.InPctx = &protosupport.PCtx{}
	}
	in.GetInPctx().Event = append(in.GetInPctx().GetEvent(),
		&protosupport.PCtxEvent{
			Line: []*protosupport.PCtxMessage{
				{
					Stamp:   timestamppb.Now(),
					Level:   log.LogLevel_LOGLEVEL_INFO,
					Message: fmt.Sprintf("Call of %s by %s", in.Method, in.Caller),
				},
			},
		})
	b, err := proto.Marshal(in.InPctx)
	if err != nil {
		return nil, NewPerrorFromId("marshal of PCtx for Dispatch()", NewKernelError(KernelMarshalFailed))
	}
	detail.PctxPtr, detail.PctxLen = sliceToTwoInt64s(b)

	b, err = proto.Marshal(in.Param)
	if err != nil {
		return nil, NewPerrorFromId("marshal of any for Dispatch()", NewKernelError(KernelMarshalFailed))
	}

	detail.ParamPtr, detail.ParamLen = sliceToTwoInt64s(b)

	resultPctx := make([]byte, GetMaxMessageSize())
	detail.OutPctxPtr, detail.OutPctxLen = sliceToTwoInt64s(resultPctx)

	resultPtr := make([]byte, GetMaxMessageSize())
	detail.ResultPtr, detail.ResultLen = sliceToTwoInt64s(resultPtr)
	if detail.ResultLen != int64(GetMaxMessageSize()) {
		panic("GetMaxMessageSize() should be the result length!")
	}
	out.ErrorId = &protosupport.KernelErrorId{High: 1, Low: 2}

	detail.ErrorPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Low))

	// THE CALL
	u := uintptr(unsafe.Pointer(detail))
	dispatch(int32(u))

	// we need to process the dispatch error first because if there was
	// an error, it could be that the pointers were not used
	dispatchErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.ErrorPtr))))
	derr := NewKernelError(KernelErrorCode(dispatchErrDataPtr[0]))
	if derr.IsError() {
		return nil, NewPerrorFromId("dispatch error", derr)
	}

	// no error sent back to use, now we will attempt to unmarshal
	// try the outpctx
	libprint("DISPATCH ", "preparing for pctx return value")
	out.OutPctx = &protosupport.PCtx{}
	if detail.OutPctxLen > 0 {
		err = proto.Unmarshal(resultPctx[:detail.OutPctxLen], out.OutPctx)
	} else {
		out.OutPctx = nil
	}
	if err != nil {
		return nil, NewPerrorFromId("unmarshal of PCtx in Dispatch()", NewKernelError(KernelUnmarshalFailed))
	}

	libprint("DISPATCH ", "preparing for result return value, %d",
		detail.ResultLen)
	out.Result = &anypb.Any{}
	err = proto.Unmarshal(resultPtr[:detail.ResultLen], out.Result)
	if err != nil {
		libprint("DISPATCH ", "unmarshal err %v with result len %d", err, detail.ResultLen)
		return nil, NewPerrorFromId("unmarshal of result in Dispatch()", NewKernelError(KernelUnmarshalFailed))
	}
	libprint("DISPATCH ", "returning our result %s --  done with %s", out.Result, in.Method)
	return out, nil
}

// BindMethodIn binds a method that only has an in parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func (l *callImpl) BindMethodIn(in *call.BindMethodRequest, _ func(Pctx, proto.Message) error) (*call.BindMethodResponse, error) {
	return l.bindMethodByName(in, call.MethodDirection_MethodDirectionIn)
}

// BindMethodInNoPctx binds a method that only has an in parameter and does not
// use the Pctx mechanism for logging.  This may, in fact, be a terrible idea but one
// cannot write a separate logger server with having this.
// xxxfixme: temporary? Should this be a different kernel call?
func (l *callImpl) BindMethodInNoPctx(in *call.BindMethodRequest, _ func(proto.Message) error) (*call.BindMethodResponse, error) {
	return l.bindMethodByName(in, call.MethodDirection_MethodDirectionIn)
}

// BindMethodOut binds a method that only has an out parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func (l *callImpl) BindMethodOut(in *call.BindMethodRequest, _ func(Pctx) (proto.Message, error)) (*call.BindMethodResponse, error) {
	return l.bindMethodByName(in, call.MethodDirection_MethodDirectionOut)
}

// BindMethodBoth binds a method that has both an in and out parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func (l *callImpl) BindMethodBoth(in *call.BindMethodRequest, _ func(Pctx, proto.Message) (proto.Message, error)) (*call.BindMethodResponse, error) {
	return l.bindMethodByName(in, call.MethodDirection_MethodDirectionBoth)
}

func (l *callImpl) bindMethodByName(in *call.BindMethodRequest, dir call.MethodDirection) (*call.BindMethodResponse, error) {
	out := new(call.BindMethodResponse)

	out.ErrorId = &protosupport.KernelErrorId{High: 6, Low: 7}
	out.ErrorId.High = 1
	out.ErrorId.Low = 2

	out.MethodId = &protosupport.MethodId{High: 10, Low: 11}

	detail := new(BindPayload)
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
	detail.MethodId = (*[2]int64)(unsafe.Pointer(&out.MethodId.Low))
	detail.ErrorPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Low))

	// THE CALL
	u := uintptr(unsafe.Pointer(detail))
	bindMethod(int32(u))

	// check for in band errors
	kernelErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.ErrorPtr))))
	kerr := NewKernelError(KernelErrorCode(kernelErrDataPtr[0]))
	if kerr.IsError() {
		return nil, NewPerrorFromId("bind error", kerr)
	}

	methodDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.MethodId))))
	mid := MethodIdFromUint64(uint64(methodDataPtr[1]), uint64(uint64(methodDataPtr[0])))
	out.MethodId = MarshalMethodId(mid)

	out.ErrorId = MarshalKernelErrId(NoKernelErr())

	return out, nil
}

func (l *callImpl) Run(in *call.RunRequest) (*call.RunResponse, error) {
	out := new(call.RunResponse)

	// allocate space for any error
	out.ErrorId = &protosupport.KernelErrorId{High: 6, Low: 7}
	out.ErrorId.High = 1
	out.ErrorId.Low = 2

	detail := new(RunPayload)
	detail.Wait = 0
	if in.Wait {
		detail.Wait = 1
	}
	detail.KernelErrorPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Low))

	// THE CALL, and the walls came down...
	u := uintptr(unsafe.Pointer(detail))
	run(int32(u))

	kernelErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.KernelErrorPtr))))
	kerr := NewKernelError(KernelErrorCode(kernelErrDataPtr[0]))
	if kerr.IsError() {
		out.ErrorId = MarshalKernelErrId(kerr)
		return out, NewPerrorFromId("kernel failed to start your process", kerr)
	}
	out.ErrorId = MarshalKernelErrId(NoKernelErr())
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
		detail := new(ExportPayload)
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
		kerr := NewKernelError(KernelErrorCode(kernelErrDataPtr[0]))
		if kerr.IsError() {
			if isExport {
				return NewPerrorFromId("export error", kerr)
			} else {
				return NewPerrorFromId("require error", kerr)
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
	out.ErrorId = &protosupport.KernelErrorId{High: 6, Low: 7}
	out.ErrorId.High = 1
	out.ErrorId.Low = 2
	ptr := (*[2]int64)(unsafe.Pointer(&out.ErrorId.Low))
	err := l.exportOrRequire(in.GetService(), ptr, true)
	if err != nil {
		return nil, err
	}
	out.ErrorId = MarshalKernelErrId(NoKernelErr())
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
	out.ErrorId = &protosupport.KernelErrorId{High: 6, Low: 7}
	out.ErrorId.High = 1
	out.ErrorId.Low = 2
	ptr := (*[2]int64)(unsafe.Pointer(&out.ErrorId.Low))
	err := l.exportOrRequire(in.GetService(), ptr, false)
	if err != nil {
		return nil, err
	}
	out.ErrorId = MarshalKernelErrId(NoKernelErr())
	return out, nil
}

// xxx this may be a bad idea.  this is probably only temporary til I can work out if we
// xxx should support this at all. currently, it is used by terminal logger.  I made a copy
// xxx of the other bindMethodByName so it would be easier to delete this one.
func (l *callImpl) bindMethodByNameNoPctx(in *call.BindMethodRequest, dir call.MethodDirection) (*call.BindMethodResponse, error) {
	out := new(call.BindMethodResponse)

	out.ErrorId = &protosupport.KernelErrorId{High: 6, Low: 7}
	out.ErrorId.High = 1
	out.ErrorId.Low = 2

	out.MethodId = &protosupport.MethodId{High: 10, Low: 11}

	detail := new(BindPayload)
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
	detail.MethodId = (*[2]int64)(unsafe.Pointer(&out.MethodId.Low))
	detail.ErrorPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Low))

	// THE CALL
	u := uintptr(unsafe.Pointer(detail))
	bindMethod(int32(u))

	// check for in band errors
	kernelErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.ErrorPtr))))
	kerr := NewKernelError(KernelErrorCode(kernelErrDataPtr[0]))
	if kerr.IsError() {
		return nil, NewPerrorFromId("bind error", kerr)
	}

	methodDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.MethodId))))
	mid := MethodIdFromUint64(uint64(methodDataPtr[1]), uint64(uint64(methodDataPtr[0])))
	out.MethodId = MarshalMethodId(mid)

	out.ErrorId = MarshalKernelErrId(NoKernelErr())

	return out, nil
}

func libprint(call, format string, arg ...interface{}) {
	if libparigotVerbose {
		part1 := fmt.Sprintf("libparigot:%s", call)
		part2 := fmt.Sprintf(format, arg...)
		print(part1, part2, "\n")
	}
}

func (l *callImpl) BlockUntilCall(in *call.BlockUntilCallRequest) (*call.BlockUntilCallResponse, error) {
	out := &call.BlockUntilCallResponse{
		Method:  &protosupport.MethodId{High: 12, Low: 13},
		Call:    &protosupport.CallId{High: 22, Low: 33},
		ErrorId: &protosupport.KernelErrorId{High: 7, Low: 8},
	}

	payload := &BlockPayload{}

	if len(in.PctxBuffer) > 0 {
		payload.PctxPtr, payload.PctxLen = sliceToTwoInt64s(in.PctxBuffer)
	} else {
		payload.PctxPtr = 0
		payload.PctxLen = 0
	}
	payload.ParamPtr, payload.ParamLen = sliceToTwoInt64s(in.ParamBuffer)
	payload.ErrorPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Low))
	payload.MethodId = (*[2]int64)(unsafe.Pointer(&out.Method.Low))
	payload.CallId = (*[2]int64)(unsafe.Pointer(&out.Call.Low))
	libprint("BLOCKUNTILCALL ", "params ready (%x,%d) and (%x,%d)",
		payload.PctxPtr, payload.PctxLen, payload.ParamPtr, payload.ParamLen)

	// THE CALL
	u := uintptr(unsafe.Pointer(payload))
	blockUntilCall(int32(u))

	// unpack the result
	kernelErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(payload.ErrorPtr))))
	kerr := NewKernelError(KernelErrorCode(kernelErrDataPtr[0]))
	if kerr.IsError() {
		return nil, NewPerrorFromId("BlockUntilCall error", kerr)
	}

	callDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(payload.CallId))))
	cid := CallIdFromUint64(uint64(callDataPtr[1]), uint64(callDataPtr[0]))
	methDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(payload.MethodId))))
	mid := MethodIdFromUint64(uint64(methDataPtr[1]), uint64(methDataPtr[0]))

	out.Call = MarshalCallId(cid)
	out.Method = MarshalMethodId(mid)
	out.ErrorId = MarshalKernelErrId(NoKernelErr())

	// get the data
	out.ParamLen = int32(payload.ParamLen)
	out.PctxLen = int32(payload.PctxLen)

	libprint("BLOCKUNTILCALL ", "unpacked the data %s,%s --- paramlen %d, pctxlen %d", mid.Short(), cid.Short(),
		out.ParamLen, out.PctxLen)
	return out, nil
}

func (l *callImpl) ReturnValue(in *call.ReturnValueRequest) (*call.ReturnValueResponse, error) {
	detail := &ReturnValuePayload{}

	detail.PctxPtr, detail.PctxLen = sliceToTwoInt64s(in.PctxBuffer)
	detail.ResultPtr, detail.ResultLen = sliceToTwoInt64s(in.ResultBuffer)

	libprint("RETURNVALUE ", "buffers for pctx and result sending to kernel are size %d,%d",
		len(in.PctxBuffer), len(in.ResultBuffer))

	detail.MethodId[0] = int64(in.Method.GetLow())
	detail.MethodId[1] = int64(in.Method.GetHigh())
	detail.CallId[0] = int64(in.Call.GetLow())
	detail.CallId[1] = int64(in.Call.GetHigh())
	detail.KernelErrorPtr[0] = int64(in.ErrorId.GetLow())
	detail.KernelErrorPtr[1] = int64(in.ErrorId.GetHigh())

	u := uintptr(unsafe.Pointer(detail))
	returnValue(int32(u))
	// check to see the return value
	kerr := NewKernelError(KernelErrorCode(detail.KernelErrorPtr[0]))
	big := MarshalKernelErrId(kerr)
	if kerr.IsError() {
		return nil, NewPerrorFromId("failed to process return value", kerr)
	}
	return &call.ReturnValueResponse{
		ErrorId: big,
	}, nil
}

func newCallImpl() Call {
	return &callImpl{}
}
