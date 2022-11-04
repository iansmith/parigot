// This package is a thin wrapper around kernel functionality. Its primary
// function is to check responses for errors and return them as a go
// error (type Perror).

package lib

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/g/pb/kernel"
	"github.com/iansmith/parigot/g/pb/log"
	"github.com/iansmith/parigot/g/pb/parigot"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Exit(in *kernel.ExitRequest) {
	exit(in)
}

// Register calls the kernel to register the given type. The Id returned is only
// useful when the error is nil.  This should only be called by clients of the
// interface being registered.  Usually this is done automatically by the init()
// method of the generated client side code.
//
//go:noinline
func Register(in *kernel.RegisterRequest, out *kernel.RegisterResponse) (Id, error) {
	out.ErrorId = &parigot.KernelErrorId{High: 6, Low: 7}
	out.ErrorId.High = 1
	out.ErrorId.Low = 2
	out.ServiceId = &parigot.ServiceId{High: 3, Low: 4}

	detail := new(RegPayload)
	pkgSh := (*reflect.StringHeader)(unsafe.Pointer(&in.ProtoPackage))
	detail.PkgPtr = int64(pkgSh.Data)
	detail.PkgLen = int64(pkgSh.Len)

	serviceSh := (*reflect.StringHeader)(unsafe.Pointer(&in.Service))
	detail.ServicePtr = int64(serviceSh.Data)
	detail.ServiceLen = int64(serviceSh.Len)
	// choosing low's addr means you ADD 8 to get to the high
	detail.ErrorPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Low))
	detail.ServiceIdPtr = (*[2]int64)(unsafe.Pointer(&out.ServiceId.Low))

	u := uintptr(unsafe.Pointer(detail))
	register(int32(u))

	// marshal them back together
	svcDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.ServiceIdPtr))))
	sid := ServiceIdFromUint64(uint64(svcDataPtr[1]), uint64(uint64(svcDataPtr[0])))
	out.ServiceId = MarshalServiceId(sid)
	regErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.ErrorPtr))))
	err := NewKernelError(KernelErrorCode(regErrDataPtr[0]))
	// in case the caller walks the structure repacks a new protobuf
	out.ServiceId = MarshalServiceId(sid)
	out.ErrorId = MarshalKernelErrId(err)

	if err.IsError() {
		return sid, NewPerrorFromId("failed to register properly", err)
	}
	return sid, nil
}

// Locate is a kernel request that returns either a reference to the service
// or an error.  In the former case, the token returned can be used with Dispatch()
// to make a call on a remote service.  It is implicit in the use of this call that
// the caller wants to be a client of the service in question.  This call can
// be made by clients or servers, but in either case the code in question becomes
// a client of the named service.
//
//go:noinline
func Locate(in *kernel.LocateRequest, out *kernel.LocateResponse) (Id, error) {
	out.ErrorId = &parigot.KernelErrorId{High: 6, Low: 7}
	out.ErrorId.High = 1
	out.ErrorId.Low = 2
	out.ServiceId = &parigot.ServiceId{High: 3, Low: 4}

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
		return sid, NewPerrorFromId("failed to locate properly", err)
	}
	return sid, nil
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
func Dispatch(in *kernel.DispatchRequest) (*kernel.DispatchResponse, error) {
	out := new(kernel.DispatchResponse)

	detail := &DispatchPayload{}
	detail.ServiceId[0] = int64(in.ServiceId.GetLow())
	detail.ServiceId[1] = int64(in.ServiceId.GetHigh())

	detail.MethodPtr, detail.MethodLen = stringToTwoInt64s(in.Method)
	detail.CallerPtr, detail.CallerLen = stringToTwoInt64s(in.Caller)

	if in.GetInPctx() == nil {
		in.InPctx = &parigot.PCtx{}
	}
	in.GetInPctx().Event = append(in.GetInPctx().GetEvent(),
		&parigot.PCtxEvent{
			Line: []*parigot.PCtxMessage{
				&parigot.PCtxMessage{
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
	out.ErrorId = &parigot.KernelErrorId{High: 1, Low: 2}

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
	out.OutPctx = &parigot.PCtx{}
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
		libprint("DISPATCH", "unmarshal err %v with result len %d", err, detail.ResultLen)
		return nil, NewPerrorFromId("unmarshal of result in Dispatch()", NewKernelError(KernelUnmarshalFailed))
	}
	libprint("DISPATCH", "returning our result %s, done with %s", out.Result, in.Method)
	return out, nil
}

// BindMethodIn binds a method that only has an in parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func BindMethodIn(in *kernel.BindMethodRequest, _ func(Pctx, proto.Message) error) (*kernel.BindMethodResponse, error) {
	return bindMethodByName(in, kernel.MethodDirection_MethodDirectionIn)
}

// BindMethodInNoPctx binds a method that only has an in parameter and does not
// use the Pctx mechanism for logging.  This may, in fact, be a terrible idea but one
// cannot write a separate logger server with having this.
// xxxfixme: temporary? Should this be a different kernel call?
func BindMethodInNoPctx(in *kernel.BindMethodRequest, _ func(proto.Message) error) (*kernel.BindMethodResponse, error) {
	return bindMethodByName(in, kernel.MethodDirection_MethodDirectionIn)
}

// BindMethodOut binds a method that only has an out parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func BindMethodOut(in *kernel.BindMethodRequest, _ func(Pctx) (proto.Message, error)) (*kernel.BindMethodResponse, error) {
	return bindMethodByName(in, kernel.MethodDirection_MethodDirectionOut)
}

// BindMethodBoth binds a method that has both an in and out parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func BindMethodBoth(in *kernel.BindMethodRequest, _ func(Pctx, proto.Message) (proto.Message, error)) (*kernel.BindMethodResponse, error) {
	return bindMethodByName(in, kernel.MethodDirection_MethodDirectionBoth)
}

func bindMethodByName(in *kernel.BindMethodRequest, dir kernel.MethodDirection) (*kernel.BindMethodResponse, error) {
	out := new(kernel.BindMethodResponse)

	out.ErrorId = &parigot.KernelErrorId{High: 6, Low: 7}
	out.ErrorId.High = 1
	out.ErrorId.Low = 2

	out.MethodId = &parigot.MethodId{High: 10, Low: 11}

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

// xxx this may be a bad idea.  this is probably only temporary til I can work out if we
// xxx should support this at all. currently, it is used by terminal logger.  I made a copy
// xxx of the other bindMethodByName so it would be easier to delete this one.
func bindMethodByNameNoPctx(in *kernel.BindMethodRequest, dir kernel.MethodDirection) (*kernel.BindMethodResponse, error) {
	out := new(kernel.BindMethodResponse)

	out.ErrorId = &parigot.KernelErrorId{High: 6, Low: 7}
	out.ErrorId.High = 1
	out.ErrorId.Low = 2

	out.MethodId = &parigot.MethodId{High: 10, Low: 11}

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
	part1 := fmt.Sprintf("libparigot:%s", call)
	part2 := fmt.Sprintf(format, arg...)
	print(part1, part2, "\n")
}

func BlockUntilCall(in *kernel.BlockUntilCallRequest) (*kernel.BlockUntilCallResponse, error) {
	out := &kernel.BlockUntilCallResponse{
		Method:  &parigot.MethodId{High: 12, Low: 13},
		Call:    &parigot.CallId{High: 22, Low: 33},
		ErrorId: &parigot.KernelErrorId{High: 7, Low: 8},
	}
	libprint("BLOCKUNTILCALL ", "out ptr %p", out)

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

	libprint("BLOCKUNTILCALL ", "unpacked the data %s,%s --- paramlen %d, pctxlen %d\n", mid.Short(), cid.Short(),
		out.ParamLen, out.PctxLen)
	return out, nil
}

// ReturnValueEncode is a layer on top of ReturnValue.  This is here because
// there are number of cases and doing this in this library means
// the code generator can be much simpler.  It just passes all the
// information into here, and this function sorts it out.
func ReturnValueEncode(cid, mid Id, marshalError, execError error, out proto.Message, pctx Pctx) (*kernel.ReturnValueResponse, error) {
	var err error
	var a anypb.Any
	// xxxfixme we should be doing an examination of execError to see if it is a lib.Perror
	// xxxfixme and if it is, we should be pushing the user error back the other way
	rv := &kernel.ReturnValueRequest{}
	rv.Call = MarshalCallId(cid)
	rv.Method = MarshalMethodId(mid)
	rv.ErrorId = MarshalKernelErrId(NoKernelErr()) // just to allocate the space
	if marshalError != nil || execError != nil {
		if marshalError != nil {
			rv.ErrorMessage = marshalError.Error()
		} else {
			rv.ErrorMessage = execError.Error()
		}
		// we use this because the error didn't come from INSIDE
		// the kernel itself, see below for more
		rv.ErrorId = MarshalKernelErrId(NoKernelErr())
		goto encodeError
	}
	// these are the mostly normal cases, but they can go hawywire
	// due to marshalling
	rv.PctxBuffer, err = pctx.Marshal()
	if err != nil {
		goto internalMarshalProblem
	}
	err = a.MarshalFrom(out)
	if err != nil {
		goto internalMarshalProblem
	}
	rv.ResultBuffer, err = proto.Marshal(&a)
	if err != nil {
		goto internalMarshalProblem
	}

	libprint("RETURNVALUEENCODE ", "size of result buffer %d, out %s", len(rv.ResultBuffer), out)
	return ReturnValue(rv)
internalMarshalProblem:
	// this is an internal error, so we signal it the opposite way we did the others at the top
	rv.ErrorMessage = ""
	rv.ErrorId = MarshalKernelErrId(NewKernelError(KernelMarshalFailed))
encodeError:
	rv.PctxBuffer = []byte{}
	rv.ResultBuffer = []byte{}
	return ReturnValue(rv)
}

func ReturnValue(in *kernel.ReturnValueRequest) (*kernel.ReturnValueResponse, error) {
	detail := &ReturnValuePayload{}

	detail.PctxPtr, detail.PctxLen = sliceToTwoInt64s(in.PctxBuffer)
	detail.ResultPtr, detail.ResultLen = sliceToTwoInt64s(in.ResultBuffer)

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
	return &kernel.ReturnValueResponse{
		ErrorId: big,
	}, nil
}

//go:noinline
//go:linkname locate parigot.locate_
func locate(int32)

//go:noinline
//go:linkname register parigot.register_
func register(int32)

//   go:noinline
//   go:linkname register2 go.parigot.register
//   func register2(int32)

//go:noinline
//go:linkname dispatch parigot.dispatch_
func dispatch(int32)

//go:noinline
//go:linkname bindMethod parigot.bind_method_
func bindMethod(int32)

//go:noinline
//go:linkname exit parigot.exit_
func exit(in *kernel.ExitRequest) int32

//go:noinline
//go:linkname blockUntilCall parigot.block_until_call_
func blockUntilCall(int32)

//go:noinline
//go:linkname returnValue parigot.return_value_
func returnValue(int32)
