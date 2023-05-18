package id

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"runtime/debug"

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

	elementIdLetter = 'E'

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
		*protosupportmsg.TestErrorId |
		*protosupportmsg.ElementId
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
		protosupportmsg.TestErrorId |
		protosupportmsg.ElementId
}

// NoError() creates an id of the given type with the "error type" but with no error as the value.
func NoError[T AllIdPtr]() Id {
	var t T
	letter, _ := typeToLetter(t)
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
	letter, ok := typeToLetter(t)
	if ok {
		return newFromErrorCode(uint64(err), letter)
	}
	log.Printf("created an error which is not associated with a type (%T), this is likely to be wrong", t)
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

// NewQueueMsgId returns a queue id, initialized for use.
func NewQueueMsgId() Id {
	return newIdRand(queueMsgLetter)
}

// NewElementId returns a element id, initialized for use.
func NewElementId() Id {
	return newIdRand(elementIdLetter)
}

// newFromErrorCode is a convenience wrapper around creating an id which represents an error for
// any of the wrapper types.
func newFromErrorCode(code uint64, letter byte) Id {
	id := idBaseFromConst(code, code != 0, letter)
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

	letter, ok := typeToLetter(t)
	if ok {
		if letter != h[7] {
			panic(fmt.Sprintf("type of id does not match letter, expected %x but got %x", letter, h[7]))
		}
		switch letter {
		case 'k', 'Q', 'F':
			if h[6] == 0 {
				panic(fmt.Sprintf("xxx NewFrom64BitPair failed because bad high uint:%x,%x\n", high, low))
			}
		}
		if h[6]&1 == 1 {
			// this is an error so we need to create it with the error code
			h[6] &= 0xfe
			return newFromErrorCode(low, letter)
		}
		return idFromUint64Pair(high, low, letter)
	}
	// this is the case where we did NOT get a letter we understood
	if h[7] == 0 {
		panic(fmt.Sprintf("unable to understand attempt to create Id from golang type %T", t))
	}
	//print("converted from unknown type to '" + string([]byte{h[7]}) + "'\n")
	return idFromUint64Pair(high, low, h[7])
}

// NewIdCopy is dangerous and most callers should use "NewFrom64BitPair".
// This call does no checking of the validity of the values provided.  NewFrom64BitPair
// does, and anytime you are not getting handed raw bits, you should probably use
// the NewFrom64BitPair.
func NewIdCopy(high, low uint64) Id {
	return &IdBase{
		h: high,
		l: low,
	}
}

// IdRepresentsError checks the error bit.  We do not want errors marked
// with the error bit (in byte 6 of high) when they have no code.
// func IdRepresentsError(high, low uint64) bool {
// 	var h [8]byte
// 	if high == 0 && low == 0 {
// 		return false
// 	}

//		binary.LittleEndian.PutUint64(h[:], high)
//		//print(fmt.Sprintf("IdRepresentsError: 0x%x, 0x%x, [%x,%x] %v\n", high, low, h[7], h[6], h[6]&1 == 1))
//		if h[6]&1 == 1 {
//			lowBits := low & 0xffff
//			if lowBits == 0 {
//				panic("badly formed error id:" + fmt.Sprintf("%x,%x", h[7], low))
//			}
//			return lowBits != 0
//		}
//		return false
//	}

// RepresentsError returns the error code that was encoded in the id
// or it returns 0 for no error.  It will print warnings to the terminal
// in cases where the id's information is inconsistent.
func (i IdBase) RepresentsError() int16 {
	var h [8]byte

	binary.LittleEndian.PutUint64(h[:], i.High())
	print(fmt.Sprintf("RepresentsError: 0x%x, 0x%x, [%x,%x] %v\n", i.High(), i.Low(), h[7], h[6], h[6]&1 == 1))
	errorBit := h[6]&1 == 1
	lowBits := int16(i.Low() & 0xffff)
	marker := h[7]
	if errorBit {
		if marker < 32 || marker > 127 {
			print(fmt.Sprintf("Badly formed id, error marked but type code not ascii character (%d)!", marker))
		}
		if lowBits == 0 {
			print("Badly formed id, error marked but no error code; assuming no error")
			return 0
		}
		return lowBits
	}
	if lowBits == 0 {
		if marker < 32 || marker > 127 {
			print(fmt.Sprintf("Badly formed id, type code not ascii character (%d)\n", marker))
		}
		return 0
	}
	if i.Low() != 0 {
		print(fmt.Sprintf("Badly formed id, no error bit or error code but has low region data (%d,%x)\n", marker, i.Low()))
	}
	// we know low bits are not zero at this point
	if marker < 32 || marker > 127 {
		print(fmt.Sprintf("Badly formed id, error code found (%d), but no type code; returning no error\n", lowBits))
		return 0
	}
	print(fmt.Sprintf("Badly formed id, error code found (%d), but no error flag; adding error flag\n", lowBits))
	h[6] = 1 << 7
	i.h = binary.LittleEndian.Uint64(h[:])
	return lowBits
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
	asAny := interface{}(result)
	letter, ok := typeToLetter(asAny)
	if highByte != letter || !ok {
		panic(fmt.Sprintf("mismatched letters in id, type letter does not match the letter found or could not understand type %T", result))
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
	letter, ok := typeToLetter(t)
	if !ok {
		log.Printf("unexpected type provided to NewId (T%), that is likely an mistake", t)
	}
	return newIdRand(letter)
}

// LocalId() should not be used by user programs.  LocalId() returns a new id for the given type.
// The value filled in the next integer in a sequence starting at one.  This is used by some parts
// of the kernel when creating a new Id() in an effort to ease debugging.  User level code should
// always use NewId().
func LocalId[T AllIdPtr](v uint64) Id {
	var t T
	letter, ok := typeToLetter(t)
	if !ok {
		panic(fmt.Sprintf("unable to understand call to LocalId because it contains a bad type %T", t))
	}
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
	case *protosupportmsg.ElementId:
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
	case *protosupportmsg.ElementId:
		return v.GetId()
	}
	panic("unknown id type")
}

// typeToLetter returns the single ascii character (as a byte) that represents a given type of id.
func typeToLetter(i interface{}) (byte, bool) {
	switch t := i.(type) {
	case *protosupportmsg.CallId:
		return callIdLetter, true
	case *protosupportmsg.DBConnId:
		return dbConnIdLetter, true
	case *protosupportmsg.DBIOId:
		return dbIOIdLetter, true
	case *protosupportmsg.FileId:
		return fileIdLetter, true
	case *protosupportmsg.FileErrorId:
		return fileErrorIdLetter, true
	case *protosupportmsg.FileIOId:
		return fileIOIdLetter, true
	case *protosupportmsg.MethodId:
		return methodIdLetter, true
	case *protosupportmsg.ServiceId:
		return serviceIdLetter, true
	case *protosupportmsg.KernelErrorId:
		return kernelErrorIdLetter, true
	case *protosupportmsg.QueueErrorId:
		return queueErrorIdLetter, true
	case *protosupportmsg.QueueId:
		return queueIdLetter, true
	case *protosupportmsg.QueueMsgId:
		return queueMsgLetter, true
	case *protosupportmsg.TestErrorId:
		return testErrorIdLetter, true
	case *protosupportmsg.ElementId:
		return elementIdLetter, true
	case *protosupportmsg.BaseId:
		if t == nil {
			print(fmt.Sprintf("t is nil!\n"))
		}
		var h [8]byte
		binary.LittleEndian.PutUint64(h[:], t.High)
		l := (h[7])
		print(fmt.Sprintf("converted unknown type of id to '%c'\n", l))
		return l, true
	}
	log.Printf("uknown id type %T", i)
	debug.PrintStack()
	return '?', false
}
