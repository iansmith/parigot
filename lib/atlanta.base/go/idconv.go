package lib

import (
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/iansmith/parigot/g/pb/parigot"
)

const (
	kernelErrorIdLetter  = 'k'
	serviceIdLetter      = 's'
	methodIdLetter       = 'm'
	callIdLetter         = 'c'
	developerIdLetter    = 'y'
	developerErrorLetter = 'z'
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
	// KernelDataTooLarge means that the size of some part of remote call was bigger
	// than the buffer allocated to receive it.  This could be a problem either on the call or
	// the return. When this error occurs, none of the client allocated buffers/pointers
	// sent in the syscall payload are touched.
	KernelDataTooLarge KernelErrorCode = 4
	// KernelMarshalFailed is an internal error of the kernel. This means that
	// a marshal of a protobuf has failed.  This is only used in situations
	// that are internel to the kernel--if user code misbehaves in this fashion
	// an error is sent to the program _from_ the kernel about the failure.
	KernelMarshalFailed KernelErrorCode = 5
	// KernelUnmarshal failed is exactly as KernelMarshalFailed, but for unpacking
	// data.
	KernelUnmarshalFailed KernelErrorCode = 6
	// KernelCallerUnavailable means that the kernel could not find the original caller
	// requeted the computation for which results have been provided.  It is most likely
	// because the caller was killed, exited or timed out.
	KernelCallerUnavailable KernelErrorCode = 7
	// KernelServiceAlreadyClosedOrExported means that some process has already reported
	// the service in question as closed or has already expressed that it is
	// exporting (implementing this service).  This is very likely a case where there
	// are two servers that think they are or should be implementing the same service.
	KernelServiceAlreadyClosedOrExported KernelErrorCode = 8
	// KernelServiceAlreadyRequired means that this same process has already
	// required the given service.
	KernelServiceAlreadyRequired KernelErrorCode = 9
	// KernelDependencyCycle means that no deterministic startup ordering
	// exists for the set of exports and requires in use.  In other words,
	// you must refactor your program so that you do not have a cyle to make
	// it come up cleanly.
	KernelDependencyCycle KernelErrorCode = 10
)

// NoKernelErr returns a kernel error id, with the value set to zero, or no error.
func NoKernelErr() Id {
	return newFromErrorCode(0, kernelErrorIdLetter)
}

// NoDeveloperError returns a developer error id, with the value set to zero, or no error.
func NoDeveloperErr() Id {
	return newFromErrorCode(0, developerErrorLetter)
}

// NewKernelError returns a kernel error id with the value (low 16 bits) set to the
// error code.
func NewKernelError(kerr KernelErrorCode) Id {
	return newFromErrorCode(uint64(kerr), kernelErrorIdLetter)
}

// NewDeveloperError returns a developer error id with the value (low 16 bits) set to the
// error code.  Developers may use this as they please.
func NewUserError(code int16) Id {
	return newFromErrorCode(uint64(code), developerErrorLetter)
}

func newFromErrorCode(code uint64, letter byte) Id {
	id := idBaseFromConst(code, true, letter)
	return id
}

// ServiceIdFromUint64 is useful when dealing with wire values.  When you receive the two
// uint64s, you can use this to turn it into a service id.  This should not be used
// otherwise.
func ServiceIdFromUint64(high uint64, low uint64) Id {
	return idFromUint64(high, low, serviceIdLetter)
}

// DeveloperIdFromUint64 is useful when dealing with wire values.  When you receive the two
// uint64s, you can use this to turn it into a developer id.  This should not be used
// otherwise.
func UserIdFromUint64(high uint64, low uint64) Id {
	return idFromUint64(high, low, developerIdLetter)
}

// CallIdFromUint64 is useful when dealing with wire values.  When you receive the two
// uint64s, you can use this to turn it into a call id.  This should not be used
// otherwise.
func CallIdFromUint64(high uint64, low uint64) Id {
	return idFromUint64(high, low, callIdLetter)
}

// MethodIdFromUint64 is useful when dealing with wire values.  When you receive the two
// uint64s, you can use this to turn it into a method id.  This should not be used
// otherwise.
func MethodIdFromUint64(high uint64, low uint64) Id {
	return idFromUint64(high, low, methodIdLetter)
}

// NewCallId returns a new call id, initialized with the value (low 112 bits) derived
// from the source of randomness.
func NewCallId() Id {
	return newIdRand(callIdLetter)
}

// NewMethodId returns a new method id, initialized with the value (low 112 bits) derived
// from the source of randomness.
func NewMethodId() Id {
	return newIdRand(methodIdLetter)
}

// NewServiceId returns a new service id, initialized with the value (low 112 bits) derived
// from the source of randomness.
func NewServiceId() Id {
	return newIdRand(serviceIdLetter)
}

// NewDeveloperId returns a new user id, initialized with the value (low 112 bits) derived
// from the source of randomness.  Developers may use this as they please.
func NewDeveloperId() Id {
	return newIdRand(developerIdLetter)
}

// newIdRand computes a new id for any given letter with the value drawn from the source of
// randomness.
func newIdRand(letter byte) Id {
	high := rand.Uint64()
	low := rand.Uint64()
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(high))
	buf[7] = letter
	buf[6] = 0 // no low order bitset
	high = binary.LittleEndian.Uint64(buf)
	id := &IdBase{
		h: high,
		l: low,
	}
	return id
}

// sourceTwo64BitNumbers returns two numbers from the default (math) source.  Note that the
// docs of rand say that it safe for concurrent use, so there is no locking here.  If we switch
// to cryptographic source, this will probably need a lock.
func sourceTwo64BitNumbers() (uint64, uint64) {
	high := rand.Uint64()
	low := rand.Uint64()
	return high, low
}

// idFromUint64 creates a non-error type of id given by letter from the two values provided.
func idFromUint64(high uint64, low uint64, letter byte) Id {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, high)
	buf[7] = letter
	buf[6] = 0 // note, low order bit not set
	return &IdBase{h: binary.LittleEndian.Uint64(buf), l: low}
}

// UnmarshaServiceId goes from the protobuf concept of a service id (which is a large
// protobuf-based structure) to a simple Id.
func UnmarshalServiceId(sid *parigot.ServiceId) Id {
	result := &IdBase{h: sid.GetHigh(), l: sid.GetLow()}
	verifyIdType(result, serviceIdLetter)
	return result
}

// MarshalServiceId is used to convert a simple id into one suitable for use in a protobuf.
func MarshalServiceId(id Id) *parigot.ServiceId {
	return &parigot.ServiceId{
		High: id.High(),
		Low:  id.Low(),
	}
}

// UnmarshalDeveloperId goes from the protobuf concept of a developer id (which is a large
// protobuf-based structure) to a simple Id.
func UnmarshalDeveloperId(sid *parigot.DeveloperId) Id {
	result := &IdBase{h: sid.GetHigh(), l: sid.GetLow()}
	verifyIdType(result, serviceIdLetter)
	return result
}

// MarshalDeveloperId is used to convert a simple id into one suitable for use in a protobuf.
func MarshalDeveloperId(id Id) *parigot.DeveloperId {
	return &parigot.DeveloperId{
		High: id.High(),
		Low:  id.Low(),
	}
}

// UnmarshalMethodId goes from the protobuf concept of a method id (which is a large
// protobuf-based structure) to a simple Id.
func UnmarshalMethodId(mid *parigot.MethodId) Id {
	result := &IdBase{h: mid.GetHigh(), l: mid.GetLow()}
	verifyIdType(result, methodIdLetter)
	return result
}

// MarshalMethodId is used to convert a simple id into one suitable for use in a protobuf.
func MarshalMethodId(id Id) *parigot.MethodId {
	return &parigot.MethodId{
		High: id.High(),
		Low:  id.Low(),
	}
}

// UnmarshalCallId goes from the protobuf concept of a call id (which is a large
// protobuf-based structure) to a simple Id.
func UnmarshalCallId(mid *parigot.CallId) Id {
	result := &IdBase{h: mid.GetHigh(), l: mid.GetLow()}
	verifyIdType(result, callIdLetter)
	return result
}

// MarshalCallId is used to convert a simple id into one suitable for use in a protobuf.
func MarshalCallId(id Id) *parigot.CallId {
	return &parigot.CallId{
		High: id.High(),
		Low:  id.Low(),
	}
}

// UnmarshalKernelErrorId goes from the protobuf concept of a kernel error id (which is a large
// protobuf-based structure) to a simple Id.  This call verifies that the value provided from
// the protobuf is marked as an error type.
func UnmarshalKernelErrorId(sid *parigot.KernelErrorId) Id {
	result := &IdBase{h: sid.GetHigh(), l: sid.GetLow()}
	if !result.IsErrorType() {
		panic("unmarshaled a kernel id that was not marked as an error type")
	}
	return result
}

// MarshalKernelErrId is used to convert a simple id into one suitable for use in a protobuf.
func MarshalKernelErrId(id Id) *parigot.KernelErrorId {
	return &parigot.KernelErrorId{
		High: id.High(),
		Low:  id.Low(),
	}
}

// UnmarshalDeveloperErrorId goes from the protobuf concept of a developer error id (which is a large
// protobuf-based structure) to a simple Id.  This call verifies that the value provided from
// the protobuf is marked as an error type.
func UnmarshalDeveloperErrId(sid *parigot.DeveloperErrorId) Id {
	result := &IdBase{h: sid.GetHigh(), l: sid.GetLow()}
	if !result.IsErrorType() {
		panic("unmarshaled a developer id that was not marked as an error type")
	}
	return result
}

// MarshalDeveloperErrId is used to convert a simple id into one suitable for use in a protobuf.
func MarshalDeveloperErrId(id Id) *parigot.DeveloperErrorId {
	return &parigot.DeveloperErrorId{
		High: id.High(),
		Low:  id.Low(),
	}
}

// verifyIdType is used when unmarshalling or marshaling to make sure the type of the
// object in question is the one you expect.  It panics if the expectation is violated.
// It should not be used to verify ids that are of an error type.
func verifyIdType(id Id, expected byte) {
	slice := int64ToByteSlice(id.High())
	if slice[7] != expected {
		panic(fmt.Sprintf("expected id to be of type %c but was of type %c", expected, slice[7]))
	}
	if id.IsErrorType() {
		panic(fmt.Sprintf("id %s was not expected to be an error type, was expecting %c", id.Short(), expected))
	}
}
