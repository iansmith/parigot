package lib

import (
	"encoding/binary"

	"github.com/iansmith/parigot/g/pb/parigot"
)

const (
	kernelErrorIdLetter = 'k'
	serviceIdLetter     = 's'
)

type KernelErrorCode uint16

const (
	// KernelAlreadyRegistered means that the a package, service or method
	// has already been registered and the attempted 2nd registration has been
	// rejected.
	KernelAlreadyRegistered KernelErrorCode = 1
	// KernelServiceNamespaceExhausted is returned when the kernel can no
	// along accept additional packages, services, or methods.  This is used
	// primarily to thwart attempts at DOS attacks.
	KernelNamespaceExhausted KernelErrorCode = 2
	// KernelNotFound means that a package, service, or method that was requested
	// could not be found.
	KernelNotFound KernelErrorCode = 3
	// KernelDispatchTooLarge means that the result of some part of a
	// Dispatch() call was bigger than the buffer we allocated to hold it.
	// When this error occurs, none of the client allocated buffers are
	// touched.
	KernelDispatchTooLarge KernelErrorCode = 4
	// KernelMarshalFailed is an internal error of the kernel. This means that
	// a marshal of a protobuf has failed.  This is only used in situations
	// that are internel to the kernel--if user code misbehaves in this fashion
	// an error is sent to the program _from_ the kernel about the failure.
	KernelMarshalFailed KernelErrorCode = 5
	// KernelUnmarshal failed is exactly as KernelMarshalFailed, but for unpacking
	// data.
	KernelUnmarshalFailed KernelErrorCode = 6
)

func NoKernelErr() Id {
	return newFromErrorCode(0, "kernelErrorId", kernelErrorIdLetter)
}

func NewKernelError(kerr KernelErrorCode) Id {
	return newFromErrorCode(uint64(kerr), "kernelErrorId", kernelErrorIdLetter)
}

func newFromErrorCode(code uint64, name string, letter byte) Id {
	id := idBaseFromConst(code, true, name, letter)
	return id
}

func ServiceIdFromUint64(high uint64, low uint64) Id {
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

func UnmarshalKernelErrorId(sid *parigot.KernelErrorId) Id {
	return &IdBase{h: sid.GetHigh(), l: sid.GetLow(), isErrType: true,
		name: "kernelErrId", letter: kernelErrorIdLetter}
}
func MarshalKernelErrId(id Id) *parigot.KernelErrorId {
	return &parigot.KernelErrorId{
		High: id.High(),
		Low:  id.Low(),
	}
}
