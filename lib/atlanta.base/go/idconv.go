package lib

import (
	"encoding/binary"
	"log"

	"github.com/iansmith/parigot/g/pb/parigot"
)

const (
	registerIdErrLetter = 'r'
	serviceIdLetter     = 's'
	locateIdErrLetter   = 'l'
	dispatchIdErrLetter = 'd'
	protoIdErrLetter    = 'p'
)

type RegisterErrCode uint16
type DispatchErrCode uint16
type LocateErrCode uint16
type ProtoErrCode uint16

const (
	RegisterAlreadyRegistered  RegisterErrCode = 1
	RegisterNamespaceExhausted RegisterErrCode = 2

	DispatchNotFound DispatchErrCode = 1
	// DispatchTooLarge means that the result of some part of a
	// Dispatch() call was bigger than the buffer we allocated to hold it.
	// When this error occurs, none of the client allocated buffers are
	// touched.
	DispatchTooLarge DispatchErrCode = 2
	// DispatchMarshalFailed is an internal error of Dispatch. This means
	// Dispatch() was trying to marshal some data and the marshal failed.
	// Typically, it will output to the kernel logs the problem that occurred.
	DispatchMarshalFailed DispatchErrCode = 3

	LocateNotFound LocateErrCode = 1

	ProtoUnmarshalFailed = 1
	ProtoMarshalFailed   = 2
	ProtoDataTooLarge    = 3
)

func NoRegisterErr() Id {
	return newFromErrorCode(0, "registerErrorId", registerIdErrLetter)
}
func NoDispatchErr() Id {
	return newFromErrorCode(0, "dispatchErrorId", dispatchIdErrLetter)
}
func NoLocateErr() Id {
	return newFromErrorCode(0, "locationErrorId", locateIdErrLetter)
}

func NewRegisterErr(code RegisterErrCode) Id {
	return newFromErrorCode(uint64(code), "registerErrorId", registerIdErrLetter)
}
func NewLocateErr(code LocateErrCode) Id {
	return newFromErrorCode(uint64(code), "locateErrorId", locateIdErrLetter)
}
func NewDispatchErr(code DispatchErrCode) Id {
	log.Printf("createing new dispatch error with code %d", code)
	return newFromErrorCode(uint64(code), "dispatchErrorId", dispatchIdErrLetter)
}
func NewProtoErr(code ProtoErrCode) Id {
	return newFromErrorCode(uint64(code), "protoErrorID", protoIdErrLetter)
}

func newFromErrorCode(code uint64, name string, letter byte) Id {
	id := idBaseFromConst(code, true, name, letter)
	return id
}

func ServiceIdFromUint64(high, low uint64) Id {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, high)
	buf[7] = serviceIdLetter
	buf[6] = 0
	return &IdBase{h: binary.LittleEndian.Uint64(buf), l: low, isErrType: false,
		name: "serviceId", letter: serviceIdLetter}

}
func UnmarshalServiceId(sid *parigot.ServiceId) Id {
	return &IdBase{h: sid.GetHigh(), l: sid.GetLow(), isErrType: false,
		name: "serviceId", letter: serviceIdLetter}
}

func MarshalServiceId(id Id) *parigot.ServiceId {
	return &parigot.ServiceId{
		High: id.High(),
		Low:  id.Low(),
	}
}

func UnmarshalRegisterErrorId(sid *parigot.RegisterErrorId) Id {
	return &IdBase{h: sid.GetHigh(), l: sid.GetLow(), isErrType: true,
		name: "registerErrId", letter: registerIdErrLetter}
}
func MarshalRegisterErrId(id Id) *parigot.RegisterErrorId {
	return &parigot.RegisterErrorId{
		High: id.High(),
		Low:  id.Low(),
	}
}

func UnmarshalLocateErrorId(sid *parigot.LocateErrorId) Id {
	return &IdBase{h: sid.GetHigh(), l: sid.GetLow(), isErrType: true,
		name: "locateErrId", letter: locateIdErrLetter}
}

func MarshalLocateErrId(id Id) *parigot.LocateErrorId {
	return &parigot.LocateErrorId{
		High: id.High(),
		Low:  id.Low(),
	}
}

func UnmarshalDispatchErrId(sid *parigot.DispatchErrorId) Id {
	return &IdBase{h: sid.GetHigh(), l: sid.GetLow(), isErrType: true,
		name: "dispatchErrId", letter: dispatchIdErrLetter}
}
func MarshalDispatchErrId(id Id) *parigot.DispatchErrorId {
	return &parigot.DispatchErrorId{
		High: id.High(),
		Low:  id.Low(),
	}
}
