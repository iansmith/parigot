// This package is a thin wrapper around kernel functionality. Its primary
// function is to check responses for errors and return them as a go
// error (type Perror).

package lib

import (
	"github.com/iansmith/parigot/g/pb/parigot"
	"reflect"
	"unsafe"
	_ "unsafe"

	"github.com/iansmith/parigot/g/pb/kernel"
)

func Exit(in *kernel.ExitRequest) {
	exit(in)
}

type RegDetail struct {
	PkgPtr          int32     // in p0a
	PkgLen          int32     // in p0b
	ServicePtr      int32     // in p1a
	ServiceLen      int32     // in p1b
	OutErrPtr       *[2]int64 // out p0
	OutServiceIdPtr *[2]int64 // out p1
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

	detail := new(RegDetail)
	pkgSh := (*reflect.StringHeader)(unsafe.Pointer(&in.ProtoPackage))
	detail.PkgPtr = int32(pkgSh.Data)
	detail.PkgLen = int32(pkgSh.Len)

	serviceSh := (*reflect.StringHeader)(unsafe.Pointer(&in.Service))
	detail.ServicePtr = int32(serviceSh.Data)
	detail.ServiceLen = int32(serviceSh.Len)
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

func Locate(in *kernel.LocateRequest, out *kernel.LocateResponse) (Id, error) {
	locate(in, out)
	return nil, nil
}

func Dispatch(in *kernel.DispatchRequest, out *kernel.DispatchResponse) {
	dispatch(in, out)
}

//go:noinline
//go:linkname locate parigot.locate_
func locate(in *kernel.LocateRequest, out *kernel.LocateResponse) int32

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
