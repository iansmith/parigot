package lib

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"
)

// Id is a type representing a global identifier in parigot.  They are composed of
// a character ('x') and a number.  In production, that number is 112 bits of
// randomness or a small integer.  The small integer case in production is for
// error ids, to indicate a call has failed in an "expected" way.
// In development, the numbers are always small integers, so printing them out
// is easier.  In all cases, the character that is the highest order byte indicates
// the type of thing the id represents. Ids of different types with the same number
// are not equal.
type Id interface {
	// Short returns a short string for debugging, like [s-6a29].  The number is
	// the last 2 bytes of the full id number.
	Short() string
	// String returns a long string that uniquely identifies this id.  This is
	// usualy something like [r-xx-xxxxxxxx-xxxxxxxx-xxxx-xxxx] where all the x's
	// are hex digits.  Note that Short() is the equivalent of the first and last
	// five characters of string.  If the number is a small integer, the leading
	// zeros are omitted.
	String() string
	// IsError returns true if this is an error type id and there is an error.  It returns
	// false if this is an error type id and there is no error (0 value).  If
	// this is not an error type id, it panics.
	IsError() bool
	// Type returns the name of the type of id, like "service" or "locate error"
	Type() string
	// Equal returns true if the two ids are of the same type and have the same number.
	Equal(Id) bool
	// High returns the high order uint64 of this id. Please don't use this unless you
	// are sending things over the wire.
	High() uint64
	// Low returns the low order uint64 of this id. Please don't use this unless you
	// are sending things over the wire.
	Low() uint64
}

type IdBase struct {
	h, l      uint64
	isErrType bool
	name      string
	letter    byte
}

func int64ToByteSlice(i uint64) []byte {
	x := (uintptr)(unsafe.Pointer(&i))
	sh := reflect.SliceHeader{
		Data: x,
		Len:  8,
		Cap:  8,
	}
	s := *(*[]byte)(unsafe.Pointer(&sh))
	return s
}

// Short returns a short string that is useful for debugging when you want to know what
// an Id is.  It prints the type of thing, plus the last 3 bytes in hex.  Given that
// the three bytes represents 2^24, you would need to be VERY unlucky to end up with
// a collision that could cause confusion when the short versions are printed.
func (i *IdBase) Short() string {
	highBytes := int64ToByteSlice(i.h)
	key := highBytes[7]
	if i.isErrType && !i.IsError() {
		return fmt.Sprintf("[%c-NoErr-]", key)
	}
	lowBytes := int64ToByteSlice(i.l)
	return fmt.Sprintf("[%c-%02x%02x%02x]", key, lowBytes[2], lowBytes[1], lowBytes[0])
}

func (i *IdBase) String() string {
	highByte := int64ToByteSlice(i.h)
	key := highByte[7]
	highByte[7] = 0
	highByte[6] = 0 //reserved for future use
	valueHigh := binary.LittleEndian.Uint64(highByte)
	two := valueHigh >> 32
	four := valueHigh & 0xffffffff

	lowPart := make([]uint64, 4)
	lowPart[3] = i.l & 0xffff
	lowPart[2] = (i.l >> 16) & 0xffff
	lowPart[1] = (i.l >> 24) & 0xffff
	lowPart[0] = (i.l >> 28) & 0xffff

	return fmt.Sprintf("[%c-%02x-%02x:%02x-%02x-%02x-%02x]", key, two, four,
		lowPart[0], lowPart[1], lowPart[2], lowPart[3])
}

func (i *IdBase) IsError() bool {
	if !i.isErrType {
		panic("IsError called on a non-error type")
	}
	high := i.h & 0xffffffffffff
	return high != 0 || i.l != 0
}

func (i *IdBase) Type() string {
	return i.name
}

func (i *IdBase) Equal(other Id) bool {
	if i.IsError() != other.IsError() {
		return false
	}
	return i.h == other.High() && i.l == other.Low()
}

func (i *IdBase) Low() uint64 {
	return i.l
}
func (i *IdBase) High() uint64 {
	return i.h
}
func idBaseFromConst(i uint64, isErrType bool, name string, letter byte) *IdBase {
	buf := make([]byte, 8)
	buf[7] = letter
	buf[6] = 0 //reserved
	return &IdBase{
		h:         binary.LittleEndian.Uint64(buf),
		l:         i,
		isErrType: isErrType,
		name:      name,
		letter:    letter,
	}
}
