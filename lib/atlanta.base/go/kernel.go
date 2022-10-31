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
	print(fmt.Sprintf("CLIENT result of Register %s,%s\n", sid.Short(), err.Short()))

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
	print(fmt.Sprintf("CLIENT result of Locate() %s,%s\n", sid.Short(), err.Short()))

	if err.IsError() {
		return sid, NewPerrorFromId("failed to locate properly", err)
	}
	return sid, nil
}

func Dispatch(in *kernel.DispatchRequest, out *kernel.DispatchResponse) {
	dispatch(in, out)
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
func dispatch(in *kernel.DispatchRequest, out *kernel.DispatchResponse) int32

//go:noinline
//go:linkname exit parigot.exit_
func exit(in *kernel.ExitRequest) int32
