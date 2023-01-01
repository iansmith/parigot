package lib

import (
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
)

const (
	kernelErrorIdLetter = 'k'
	serviceIdLetter     = 's'
	methodIdLetter      = 'm'
	callIdLetter        = 'c'
	dbConnIdLetter      = 'D'
	dbIOIdLetter        = 'I'
	fileIdLetter        = 'f'
	fileIOIdLetter      = 'i'
)

type KernelErrorCode uint16

type AllIdPtr interface {
	*protosupport.CallId |
		*protosupport.DBConnId |
		*protosupport.DBIOId |
		*protosupport.FileId |
		*protosupport.FileIOId |
		*protosupport.MethodId |
		*protosupport.ServiceId |
		*protosupport.KernelErrorId |
		*protosupport.BaseId
}

type AllId interface {
	protosupport.CallId |
		protosupport.DBConnId |
		protosupport.DBIOId |
		protosupport.FileId |
		protosupport.FileIOId |
		protosupport.MethodId |
		protosupport.ServiceId |
		protosupport.KernelErrorId |
		protosupport.BaseId
}

const (
	// NoKernelError means just what it sounds like.  All Ids that are errors represent
	// no error as 0.
	KernelNoError KernelErrorCode = 0
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
	// the return.
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
	// KernelNetworkFailed means that we successfully connected to the nameserver, but failed
	// during the communication process itself.
	KernelNetworkFailed KernelErrorCode = 11
	// KernelNetworkConnectionLost means that our internal connection to the remote nameserver
	// was either still working but has lost "sync" in the protocol or the connection has
	// become entirely broken.  The kernel will close the connection to remote nameserver
	// and reestablish it after this error.
	KernelNetworkConnectionLost KernelErrorCode = 12
	// KernelDataTooSmall means that the kernel was speaking some protocol with a remote server,
	// such as a remote nameserver, and data read from the remote said was smaller than the protocol
	// dictated, e.g. it did not contain a checksum after a data block.
	KernelDataTooSmall KernelErrorCode = 13
	// KernelConnectionFailed means that the attempt to open a connection to a remote
	// service has failed to connect.
	KernelConnectionFailed KernelErrorCode = 14
	// KernelNSRetryFailed means that we tried twice to reach the nameserver with
	// the given request, but both times could not do so.
	KernelNSRetryFailed KernelErrorCode = 15
	// KernelBadPath means that the path cannot correspond to any Parigot file. This does not mean
	// we tried and succeeded or failed to read the file, it means only that we could determine by
	// looking at the path that this cannot succeed.
	KernelBadPath KernelErrorCode = 16
	// KernelExecError means that we received a response from the implenter of a particular
	// service's function and the execution of that function failed.
	KernelExecError KernelErrorCode = 17
	// KernelBadId means received something from your code that was supposed to be an error and
	// it did not have the proper mark on it (IsErrorType()).
	KernelBadId KernelErrorCode = 18
	// KernelDependencyFailure means that the dependency infrastructure has failed.  This is different
	// than when a user creates bad set of depedencies (KernelDependencyCycle).
	KernelDependencyFailure KernelErrorCode = 19
)

// NoError() creates an id of the given type with the "error type" but with no error as the value.
func NoError[T AllIdPtr]() Id {
	var t T
	letter := typeToLetter(t)
	return newFromErrorCode(0, letter)
}

// NoErrorMarhsalled is a wrapper to call NoError() and then Marshal() on the given type.
func NoErrorMarshaled[T AllId, P AllIdPtr]() *T {
	return Marshal[T](NoError[P]())
}

// NewError creates an error of the given type that contains the given error code. which only fills
// the low 16 bits of the low uint64 of the id.
func NewError[T AllIdPtr](err uint16) Id {
	var t T
	letter := typeToLetter(t)
	return newFromErrorCode(uint64(err), letter)
}

// NoKernelError is a wrapper to create a no error from the kernel and marshal it.
func NoKernelError() *protosupport.KernelErrorId {
	return NoErrorMarshaled[protosupport.KernelErrorId, *protosupport.KernelErrorId]()
}

// NewKernelError is just a convenience wrapper around NewError() that has the correct constant type
// its parameter.  Note that this will accept "NoKer"
func NewKernelError(code KernelErrorCode) Id {
	if code == KernelNoError {
		return NoError[*protosupport.KernelErrorId]()
	}
	return NewError[*protosupport.KernelErrorId](uint16(code))
}

// newFromErrorCode is a convenience wrapper around creating an id which represents an error for
// any of the wrapper types.
func newFromErrorCode(code uint64, letter byte) Id {
	id := idBaseFromConst(code, true, letter)
	return id
}

// NewFrom64BitPair creates a new instance of Id that has the given high and low uint64s but it
// also checks that the given high uint64 has the proper letter in the high byte.  If the given
// high uint64 does not match the proper letter for the type, it panics.  This method is probably
// not of interest to user code, even though it might look like it.  NewError() and NewRand()
// are probably the ones to consider instead.
func NewFrom64BitPair[T AllIdPtr](high, low uint64) Id {
	var h [8]byte
	binary.LittleEndian.PutUint64(h[:], high)
	var t T
	letter := typeToLetter(t)
	if letter != h[7] {
		panic(fmt.Sprintf("type of id does not match letter, expected %x but got %x", letter, h[7]))
	}
	if letter == 'k' && h[6] == 0 {
		panic(fmt.Sprintf("xxx NewFrom64BitPair failed because bad high uint:%x,%x\n", high, low))
	}
	if h[6]&1 == 1 {
		// this is an error so we need to create it with the error code
		return newFromErrorCode(low, letter)
	}
	return idFromUint64Pair(high, low, letter)
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

// idFromUint64Pair creates a non-error type of id given by letter from the two values provided.  Note
// that if you pass it a high with the error bit (in byte[6]) set, it will be cleared by this function.
func idFromUint64Pair(high uint64, low uint64, letter byte) Id {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, high)
	buf[7] = letter
	buf[6] = 0 // note, low order bit not set
	return &IdBase{h: binary.LittleEndian.Uint64(buf), l: low}
}

// Unmarshal onverts from the internal representation in memory (lib.Id) to the external, wire-format
// compatible version in protosupport.  This works for all types of Ids.
func Unmarshal[T AllIdPtr](wrapper T) Id {
	if wrapper == nil {
		return nil
	}
	inner := typeToInnerId(wrapper)
	var highByte [8]byte
	binary.LittleEndian.PutUint64(highByte[:], inner.GetHigh())
	if highByte[6]&0x01 == 1 {
		return newFromErrorCode(inner.GetLow(), byte(inner.GetAsciiValue()))
	}
	return idFromUint64Pair(inner.GetHigh(), inner.GetLow(), byte(inner.GetAsciiValue()))
}

// Marshal converts from the wire-format of protobufs to the in-memory format that is two 64 bit
// integers, in lib.Id.
func Marshal[T AllId](id Id) *T {
	if id == nil {
		return nil
	}
	highBytes := int64ToByteSlice(id.High())
	highByte := highBytes[7]

	inner := &protosupport.BaseId{
		High:       id.High(),
		Low:        id.Low(),
		AsciiValue: uint32(highByte),
	}
	result := new(T)
	var asAny interface{}
	asAny = result
	letter := typeToLetter(asAny)
	if highByte != letter {
		panic("mismatched letters in id, type letter does not match the letter found")
	}
	if asAny == nil {
		panic("attemp to add put id on an nil wrapper")
	}
	addInnerIdToType(asAny, inner)
	return result
}

// NewId() returns a new id for the given type. The value filled in is derived from a randomness
// source that is currently math.Rand which not as secure as crypto.Rand but has the advantage of
// being replayable.  Note that if you intended to create an id representing an error, use
// NewError() or NoError().
func NewId[T AllIdPtr]() Id {
	var t T
	letter := typeToLetter(t)
	return newIdRand(letter)
}

// LocalId() should not be used by user programs.  LocalId() returns a new id for the given type.
// The value filled in the next integer in a sequence starting at one.  This is used by some parts
// of the kernel when creating a new Id() in an effort to ease debugging.  User level code should
// always use NewId().
func LocalId[T AllIdPtr](v uint64) Id {
	var t T
	letter := typeToLetter(t)
	return idFromUint64Pair(0, v, letter)
}

// addInnerToType adds the given "inner" id (the real data) to an existing wrapper type that is
// one of the ones visible to user code.
func addInnerIdToType(i interface{}, inner *protosupport.BaseId) {
	switch v := i.(type) {
	case *protosupport.CallId:
		v.Id = inner
	case *protosupport.DBConnId:
		v.Id = inner
	case *protosupport.DBIOId:
		v.Id = inner
	case *protosupport.FileId:
		v.Id = inner
	case *protosupport.FileIOId:
		v.Id = inner
	case *protosupport.MethodId:
		v.Id = inner
	case *protosupport.ServiceId:
		v.Id = inner
	case *protosupport.KernelErrorId:
		v.Id = inner
	default:
		panic("unknown id type")
	}
}

// addInnerToType adds the given "inner" id (the real data) to an existing wrapper type that is
// one of the ones visible to user code.
func typeToInnerId(i interface{}) *protosupport.BaseId {
	switch v := i.(type) {
	case *protosupport.CallId:
		return v.GetId()
	case *protosupport.DBConnId:
		return v.GetId()
	case *protosupport.DBIOId:
		return v.GetId()
	case *protosupport.FileId:
		return v.GetId()
	case *protosupport.FileIOId:
		return v.GetId()
	case *protosupport.MethodId:
		return v.GetId()
	case *protosupport.ServiceId:
		return v.GetId()
	case *protosupport.KernelErrorId:
		return v.GetId()
	}
	panic("unknown id type")
}

// typeToLetter returns the single ascii character (as a byte) that represents a given type of id.
func typeToLetter(i interface{}) byte {
	switch i.(type) {
	case *protosupport.CallId:
		return callIdLetter
	case *protosupport.DBConnId:
		return dbConnIdLetter
	case *protosupport.DBIOId:
		return dbIOIdLetter
	case *protosupport.FileId:
		return fileIdLetter
	case *protosupport.FileIOId:
		return fileIOIdLetter
	case *protosupport.MethodId:
		return methodIdLetter
	case *protosupport.ServiceId:
		return serviceIdLetter
	case *protosupport.KernelErrorId:
		return kernelErrorIdLetter
	}
	panic("unknown id type")
}
