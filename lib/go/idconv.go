package lib

import (
	"encoding/binary"
	"fmt"
	"math/rand"

	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
)

const (
	kernelErrorIdLetter = 'k'

	queueErrorIdLetter = 'Q'
	queueIdLetter      = 'q'
	queueMsgLetter     = 'w'

	serviceIdLetter = 's'
	methodIdLetter  = 'm'
	callIdLetter    = 'c'
	dbConnIdLetter  = 'd'
	dbIOIdLetter    = 'e'

	fileIdLetter      = 'f'
	fileErrorIdLetter = 'F'
	fileIOIdLetter    = 'i'

	testErrorIdLetter = 'u'
)

type KernelErrorCode uint16
type QueueErrorCode uint16
type FileErrorCode uint16
type TestErrorCode uint16

type AllIdPtr interface {
	*protosupportmsg.CallId |
		*protosupportmsg.DBConnId |
		*protosupportmsg.DBIOId |
		*protosupportmsg.FileErrorId |
		*protosupportmsg.FileId |
		*protosupportmsg.FileIOId |
		*protosupportmsg.MethodId |
		*protosupportmsg.QueueErrorId |
		*protosupportmsg.QueueId |
		*protosupportmsg.QueueMsgId |
		*protosupportmsg.ServiceId |
		*protosupportmsg.KernelErrorId |
		*protosupportmsg.BaseId |
		*protosupportmsg.TestErrorId
}

type AllId interface {
	protosupportmsg.CallId |
		protosupportmsg.DBConnId |
		protosupportmsg.DBIOId |
		protosupportmsg.FileErrorId |
		protosupportmsg.FileId |
		protosupportmsg.FileIOId |
		protosupportmsg.MethodId |
		protosupportmsg.QueueErrorId |
		protosupportmsg.QueueId |
		protosupportmsg.QueueMsgId |
		protosupportmsg.ServiceId |
		protosupportmsg.KernelErrorId |
		protosupportmsg.BaseId |
		protosupportmsg.TestErrorId
}

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
func NoKernelError() *protosupportmsg.KernelErrorId {
	return NoErrorMarshaled[protosupportmsg.KernelErrorId, *protosupportmsg.KernelErrorId]()
}

// NewKernelError is just a convenience wrapper around NewError() that has the correct constant type
// its parameter.  Note that this will accept "NoKernelError" and behave correctly.
func NewKernelError(code KernelErrorCode) Id {
	if code == KernelNoError {
		return NoError[*protosupportmsg.KernelErrorId]()
	}
	return NewError[*protosupportmsg.KernelErrorId](uint16(code))
}

// NewFileError is just a convenience wrapper around NewError() that has the correct constant type
// its parameter.  Note that this will accept "NoFileError" and behave correctly.
func NewFileError(code FileErrorCode) Id {
	if code == FileNoError {
		return NoError[*protosupportmsg.FileErrorId]()
	}
	return NewError[*protosupportmsg.FileErrorId](uint16(code))
}

// NewQueueError is just a convenience wrapper around NewError() that has the
// correct constant type on its parameter.  Note that this will
// accept "NoQueueError" as its parameter.
func NewQueueError(code QueueErrorCode) Id {
	if code == QueueNoError {
		return NoError[*protosupportmsg.QueueErrorId]()
	}
	return NewError[*protosupportmsg.QueueErrorId](uint16(code))
}

// NewTestError is just a convenience wrapper around NewError() that has the
// correct constant type on its parameter.  Note that this will
// accept "NoTestError" as its parameter.
func NewTestError(code TestErrorCode) Id {
	if code == TestNoError {
		return NoError[*protosupportmsg.TestErrorId]()
	}
	return NewError[*protosupportmsg.QueueErrorId](uint16(code))
}

// NewQueueId returns a queue id, initialized for use.
func NewQueueId() Id {
	return newIdRand(queueIdLetter)
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

// Unmarshal converts protosupportmsg.*Id -> lib.id
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

// Marshal converts lib.id -> protosupportmsg.*Id
func Marshal[T AllId](id Id) *T {
	if id == nil {
		return nil
	}
	highBytes := int64ToByteSlice(id.High())
	highByte := highBytes[7]

	inner := &protosupportmsg.BaseId{
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
func addInnerIdToType(i interface{}, inner *protosupportmsg.BaseId) {
	switch v := i.(type) {
	case *protosupportmsg.CallId:
		v.Id = inner
	case *protosupportmsg.DBConnId:
		v.Id = inner
	case *protosupportmsg.DBIOId:
		v.Id = inner
	case *protosupportmsg.FileId:
		v.Id = inner
	case *protosupportmsg.FileIOId:
		v.Id = inner
	case *protosupportmsg.MethodId:
		v.Id = inner
	case *protosupportmsg.ServiceId:
		v.Id = inner
	case *protosupportmsg.KernelErrorId:
		v.Id = inner
	case *protosupportmsg.QueueErrorId:
		v.Id = inner
	case *protosupportmsg.QueueId:
		v.Id = inner
	case *protosupportmsg.QueueMsgId:
		v.Id = inner
	default:
		panic("unknown id type")
	}
}

// addInnerToType adds the given "inner" id (the real data) to an existing wrapper type that is
// one of the ones visible to user code.
func typeToInnerId(i interface{}) *protosupportmsg.BaseId {
	switch v := i.(type) {
	case *protosupportmsg.CallId:
		return v.GetId()
	case *protosupportmsg.DBConnId:
		return v.GetId()
	case *protosupportmsg.DBIOId:
		return v.GetId()
	case *protosupportmsg.FileId:
		return v.GetId()
	case *protosupportmsg.FileIOId:
		return v.GetId()
	case *protosupportmsg.MethodId:
		return v.GetId()
	case *protosupportmsg.ServiceId:
		return v.GetId()
	case *protosupportmsg.KernelErrorId:
		return v.GetId()
	case *protosupportmsg.QueueErrorId:
		return v.GetId()
	case *protosupportmsg.QueueId:
		return v.GetId()
	case *protosupportmsg.QueueMsgId:
		return v.GetId()
	case *protosupportmsg.FileErrorId:
		return v.GetId()
	case *protosupportmsg.TestErrorId:
		return v.GetId()
	}
	panic("unknown id type")
}

// typeToLetter returns the single ascii character (as a byte) that represents a given type of id.
func typeToLetter(i interface{}) byte {
	switch i.(type) {
	case *protosupportmsg.CallId:
		return callIdLetter
	case *protosupportmsg.DBConnId:
		return dbConnIdLetter
	case *protosupportmsg.DBIOId:
		return dbIOIdLetter
	case *protosupportmsg.FileId:
		return fileIdLetter
	case *protosupportmsg.FileErrorId:
		return fileErrorIdLetter
	case *protosupportmsg.FileIOId:
		return fileIOIdLetter
	case *protosupportmsg.MethodId:
		return methodIdLetter
	case *protosupportmsg.ServiceId:
		return serviceIdLetter
	case *protosupportmsg.KernelErrorId:
		return kernelErrorIdLetter
	case *protosupportmsg.QueueErrorId:
		return queueErrorIdLetter
	case *protosupportmsg.QueueId:
		return queueIdLetter
	case *protosupportmsg.QueueMsgId:
		return queueMsgLetter
	case *protosupportmsg.TestErrorId:
		return testErrorIdLetter
	}
	panic("unknown id type:" + fmt.Sprintf("%T", i))
}
