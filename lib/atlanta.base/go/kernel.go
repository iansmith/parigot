// This package is a thin wrapper around kernel functionality. Its primary
// function is to check responses for errors and return them as a go
// error (type Perror).

package lib

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/g/pb/kernel"
	"github.com/iansmith/parigot/g/pb/parigot"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Exit(in *kernel.ExitRequest) {
	exit(in)
}

// Register calls the kernel to register the given type. The Id returned is only
// useful when the error is nil.
//
//go:noinline
func Register(in *kernel.RegisterRequest, out *kernel.RegisterResponse) (Id, error) {
	out.ErrorId = &parigot.RegisterErrorId{High: 6, Low: 7}
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
	detail.OutErrPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Low))
	detail.OutServiceIdPtr = (*[2]int64)(unsafe.Pointer(&out.ServiceId.Low))

	u := uintptr(unsafe.Pointer(detail))
	register(int32(u))

	// marshal them back together
	svcDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.OutServiceIdPtr))))
	sid := ServiceIdFromUint64(uint64(svcDataPtr[1]), uint64(uint64(svcDataPtr[0])))
	out.ServiceId = MarshalServiceId(sid)
	regErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.OutErrPtr))))
	err := NewRegisterErr(RegisterErrCode(regErrDataPtr[0]))
	// in case the caller walks the structure repacks a new protobuf
	out.ServiceId = MarshalServiceId(sid)
	out.ErrorId = MarshalRegisterErrId(err)

	if err.IsError() {
		return sid, NewPerrorFromId("failed to register properly", err)
	}
	return sid, nil
}

//go:noinline
func Locate(in *kernel.LocateRequest, out *kernel.LocateResponse) (Id, error) {
	out.ErrorId = &parigot.LocateErrorId{High: 6, Low: 7}
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
	detail.OutErrPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Low))
	detail.OutServiceIdPtr = (*[2]int64)(unsafe.Pointer(&out.ServiceId.Low))

	u := uintptr(unsafe.Pointer(detail))
	locate(int32(u))
	// marshal them back together
	svcDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.OutServiceIdPtr))))
	sid := ServiceIdFromUint64(uint64(svcDataPtr[1]), uint64(uint64(svcDataPtr[0])))
	out.ServiceId = MarshalServiceId(sid)
	locErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.OutErrPtr))))

	err := NewLocateErr(LocateErrCode(locErrDataPtr[0]))
	// in case the caller walks the structure repacks a new protobuf
	out.ServiceId = MarshalServiceId(sid)
	out.ErrorId = MarshalLocateErrId(err)

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
// itself, otherwise it implies a remote procedure call.  Note that the
// kernel calls, including this one, have an in-band and out-of-band error
// return.  If there was no error out-of-band the second return value will
// be nil.  However, inside the result there may more error values to check.
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
					Level:   parigot.LogLevel_LOGLEVEL_INFO,
					Message: fmt.Sprintf("Call of %s by %s", in.Method, in.Caller),
				},
			},
		})
	b, err := proto.Marshal(in.InPctx)
	if err != nil {
		return nil, NewPerrorFromId("marshal of PCtx for Dispatch()", NewProtoErr(ProtoMarshalFailed))
	}
	detail.PctxPtr, detail.PctxLen = sliceToTwoInt64s(b)

	b, err = proto.Marshal(in.Param)
	if err != nil {
		return nil, NewPerrorFromId("marshal of any for Dispatch()", NewProtoErr(ProtoMarshalFailed))
	}

	detail.ParamPtr, detail.ParamLen = sliceToTwoInt64s(b)

	resultPctx := make([]byte, GetMaxMessageSize())
	detail.OutPctxPtr, detail.OutPctxLen = sliceToTwoInt64s(resultPctx)

	resultPtr := make([]byte, GetMaxMessageSize())
	detail.ResultPtr, detail.ResultLen = sliceToTwoInt64s(resultPtr)
	if detail.ResultLen != int64(GetMaxMessageSize()) {
		panic("GetMaxMessageSize() should be the result length!")
	}
	out.ErrorId = &parigot.DispatchErrorId{High: 1, Low: 2}

	detail.ErrorPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Low))

	// THE CALL
	u := uintptr(unsafe.Pointer(detail))
	dispatch(int32(u))

	// we need to process the dispatch error first because if there was
	// an error, it could be that the pointers were not used
	print("P9", detail.OutPctxPtr, "\n")
	dispatchErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.ErrorPtr))))
	derr := NewDispatchErr(DispatchErrCode(dispatchErrDataPtr[0]))
	if derr.IsError() {
		return nil, NewPerrorFromId("dispatch error", derr)
	}

	// no error sent back to use, now we will attempt to unmarshal
	// try the outpctx
	out.OutPctx = &parigot.PCtx{}
	err = proto.Unmarshal(resultPctx[:detail.OutPctxLen], out.OutPctx)
	if err != nil {
		return nil, NewPerrorFromId("unmarshal of PCtx in Dispatch()", NewProtoErr(ProtoUnmarshalFailed))
	}

	out.Result = &anypb.Any{}
	err = proto.Unmarshal(resultPtr[:detail.ResultLen], out.Result)
	if err != nil {
		return nil, NewPerrorFromId("unmarshal of result in Dispatch()", NewProtoErr(ProtoUnmarshalFailed))
	}
	return out, nil
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
//go:linkname exit parigot.exit_
func exit(in *kernel.ExitRequest) int32
