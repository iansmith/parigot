package id

import (
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/tetratelabs/wazero/api"
)

type NameInfo interface {
	ShortString() string
	Letter() byte
}

// IdRoot[T] is used to indicate an id in parigot.
// This type is usually "covered up" by generated code
// that will given it a name lile FileId or QueueId so
// these types cannot be compared for equality, assigned
// to each other, and similar, despite being the same
// underlying type.
type IdRoot[T NameInfo] struct {
	high, low uint64
}

// NewIdRoot[T] returns a new Id of type[T] fille with
// 120 bits of randomness.
func NewIdRoot[T NameInfo]() IdRoot[T] {
	var t T
	high := rand.Uint64()
	low := rand.Uint64()
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(high))
	buf[7] = t.Letter()
	high = binary.LittleEndian.Uint64(buf)
	id := IdRoot[T]{
		high: high,
		low:  low,
	}
	return id
}

// NewIdType is dangerous in that it performs no checks about the validity of the
// data provided. Its use is discouraged.  It will obey the
// values provided by the T type parameter regarding the
// high order byte, even if this value is provided in h.
func NewIdTyped[T NameInfo](h, l uint64) IdRoot[T] {
	var t T
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, h)
	buf[7] = t.Letter()
	upper := binary.LittleEndian.Uint64(buf)
	upper |= (h & 0xffffffffffff)
	return IdRoot[T]{high: upper, low: l}
}

// ZeroValue is a special value of an id. Thes should
// be returned when an Id should not be used, such as when
// an error is also returned from a function.  Note that
// the zero value is not the same as the empty value.
func ZeroValue[T NameInfo]() IdRoot[T] {
	var t T
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(0xffffffffffff))
	buf[7] = t.Letter()
	h := binary.LittleEndian.Uint64(buf)
	l := uint64(0xffffffffffffffff)
	return IdRoot[T]{high: h, low: l}
}

// Short returns a short string representing this value.
// Strings returned represent the low order 16 bits of the
// this Id.  However, for debugging, this amount data is
// enough to uniquely identify a particular id.  If you want to
// see the entire 120 bits, then use String().
func (i IdRoot[T]) Short() string {
	var t T

	if i.IsZeroValue() {
		return fmt.Sprintf("[%s-zero]", t.ShortString())
	}
	if i.IsEmptyValue() {
		return fmt.Sprintf("[%s-empty]", t.ShortString())
	}

	low := make([]byte, 8)
	binary.LittleEndian.PutUint64(low, i.low)

	return fmt.Sprintf("[%s-%02x%02x]", t.ShortString(), low[1], low[0])
}

// WriteGuestLe will write an id into the guest memory when
// running on the host.  This always writes the data in Little
// Endian format.
func (i IdRoot[V]) WriteGuestLe(mem api.Memory, offset uint32) bool {
	if !mem.WriteUint64Le(offset, i.high) {
		return false
	}
	if !mem.WriteUint64Le(offset+8, i.low) {
		return false
	}
	return true
}

// String() returns a string that contains the short string
// name of this id (like "file" or "queue") and then 5 groups
// of numbers.  From left to right these are bytes[4-6] of
// the high part, bytes[0-3] of the high part, bytes [4-7]
// of the low part, bytes[2-3] of the low part, and then low
// order two bytes of the low part.  The low part is printed
// this way so the last section of the string match the portion
// of the id printed by Short()
func (i IdRoot[T]) String() string {
	var t T
	threehigh := (i.high >> 32) & 0xffffff
	fourlow := (i.high) & 0xffffffff

	lowPart := make([]uint32, 3)
	lowPart[2] = uint32(i.low & 0xffff)
	lowPart[1] = uint32((i.low >> 24) & 0xffff)
	lowPart[0] = uint32((i.low >> 32) & 0xffffffff)

	return fmt.Sprintf("%s-%06x-%08x:%08x-%04x-%04x",
		t.ShortString(), threehigh, fourlow,
		lowPart[0], lowPart[1], lowPart[2])
}

// Equal will compare two ids for equality.  At this level
// it can compare _any_ two ids, but most users will be using
// generated code that disallows comparisons between
// id types. Note that the empty value and the zero value
// are not equal to anything, including each other.
func (i IdRoot[T]) Equal(other IdRoot[T]) bool {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, i.high)
	if i.IsZeroValue() || other.IsZeroValue() {
		return false // can't be compared to anything
	}
	return i.high == other.high && i.low == other.low
}

// IsEmptyValue tells you if the given id is actually just full of zeros in
// the 15 data bytes.  This almost means that the caller gave you a bad id, since
// the chance of all 15 data bytes being zero is very low, and vastly lower
// than somebody forgetting to initialize a value.
func (i IdRoot[T]) IsEmptyValue() bool {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, i.high)
	for i := 0; i < 6; i++ {
		if buf[i] != 0 {
			return false
		}
	}
	return i.low == 0
}

// Use of High() is not recommended for user code.  It returns the
// high 8 bytes of the 128bit id.
func (i IdRoot[T]) High() uint64 {
	return i.high
}

// Use of Low() is not recommended for user code.  It returns the
// low 8 bytes of the 128bit id.
func (i IdRoot[T]) Low() uint64 {
	return i.low
}

// IsZeroValue checks an id value to see if it is the
// bit pattern of the zero value.
func (i IdRoot[T]) IsZeroValue() bool {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, i.high)
	for i := 0; i < 6; i++ {
		if buf[i] != 0xff {
			return false
		}
	}
	return i.low == 0xffffffffffffffff
}

// IsZeroOrEmptyValue returns true if the value this
// is called on is either the zero or empty value.

func (i IdRoot[T]) IsZeroOrEmptyValue() bool {
	return i.IsZeroValue() || i.IsEmptyValue()
}
