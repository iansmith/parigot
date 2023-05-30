package id

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"reflect"
	"unsafe"

	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	"github.com/tetratelabs/wazero/api"
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
	// StringRaw returns the same thing as string, but without the enclosing
	// square brackets. String() is intended for humans, StringRaw() is not.
	StringRaw() string
	// IsError returns true if the id given represents an error.
	IsError() bool
	// ErrorCode returns the error code that was encoded in the id
	// or it returns 0 for no error.  It will print warnings to the terminal
	// in cases where the id's information is inconsistent.
	ErrorCode() uint16
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
	h, l uint64
}

// func newIdBaseFromKnown(high uint64, low uint64, isError bool, letter byte) {
// 	base := &IdBase{
// 		h: high,
// 		l: low,
// 	}
// 	highBytes := int64ToByteSlice(base.h)
// 	highBytes[7] = letter
// 	highBytes[6] = 0
// 	if isError {
// 		highBytes[6] = 1
// 	}
// }

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
	if i.IsErrorType() && !i.IsError() {
		return "[NoErr]"
		//return fmt.Sprintf("[%c-NoErr-]", key)
	}
	lowBytes := int64ToByteSlice(i.l)
	return fmt.Sprintf("[%c-%02x%02x%02x]", key, lowBytes[2], lowBytes[1], lowBytes[0])
}

// ErrorCode returns the error code contained in the lowest 2 bytes of the
// lower half of this id.  It panics if the id is not an error type.
func (i *IdBase) ErrorCode() uint16 {
	if !i.IsErrorType() {
		panic("called ErrorCode() on a non error type:" + i.Short())
	}
	return uint16(i.l & 0xffff)
}

func (i *IdBase) IsErrorType() bool {
	highBytes := int64ToByteSlice(i.h)
	if highBytes[7] == 0x6b && highBytes[6] == 0 {
		panic(fmt.Sprintf("combination not allowed! kerr but not an error %x,%x:", i.High(), i.Low()))
	}
	return highBytes[6]&1 == 1
}

func (i *IdBase) StringRaw() string {
	copy := i.h
	highByte := int64ToByteSlice(copy)
	key := highByte[7]
	highByte[7] = 0
	highByte[6] = 0 //lowest bit is true for an error type
	valueHigh := binary.LittleEndian.Uint64(highByte)
	two := valueHigh >> 32
	four := valueHigh & 0xffffffff

	lowPart := make([]uint64, 4)
	lowPart[3] = i.l & 0xffff
	lowPart[2] = (i.l >> 16) & 0xffff
	lowPart[1] = (i.l >> 24) & 0xffff
	lowPart[0] = (i.l >> 28) & 0xffff

	return fmt.Sprintf("%c-%04x-%08x:%04x-%04x-%04x-%04x", key, two, four,
		lowPart[0], lowPart[1], lowPart[2], lowPart[3])
}

func (i *IdBase) String() string {
	return fmt.Sprintf("[%s]", i.StringRaw())
}

func (i *IdBase) IsError() bool {
	if !i.IsErrorType() {
		print("------------- IS ERROR CALLED ON NO ERROR TYPE --------------\n")
		return false
	}
	// and with bytes 0-5 turned on
	high := i.h & 0xffffffffffff
	return high != 0 || i.l != 0
}

func (i *IdBase) Equal(other Id) bool {
	return i.h == other.High() && i.l == other.Low()
}

func (i *IdBase) Low() uint64 {
	return i.l
}
func (i *IdBase) High() uint64 {
	return i.h
}

// idBaseFromConst is useful for generating a new id from a given low order uint64.
func idBaseFromConst(i uint64, isErrType bool, letter byte) *IdBase {
	buf := make([]byte, 8)
	buf[7] = letter
	buf[6] = 0
	if isErrType {
		buf[6] = 1 //low bit
	}
	return &IdBase{
		h: binary.LittleEndian.Uint64(buf),
		l: i,
	}
}

// ID RAW
type IdRaw struct {
	high, low uint64
}

type idNameInfo interface {
	ShortString() string
	Letter() byte
	IsError() bool
}

type ErrIdResponse interface {
	IsError() bool
	Short() string
	String() string
	Raw() IdRaw
}

//
// ID ROOT
//

type IdRoot[T idNameInfo] struct {
	high, low uint64
}

func NewIdRoot[T idNameInfo]() IdRoot[T] {
	var t T
	high := rand.Uint64()
	low := rand.Uint64()
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(high))
	buf[7] = t.Letter()
	buf[6] = 0 // no low order bitset
	high = binary.LittleEndian.Uint64(buf)
	id := IdRoot[T]{
		high: high,
		low:  low,
	}
	return id
}
func NewIdRootFromRaw[T idNameInfo](r IdRaw) IdRoot[T] {
	return IdRoot[T]{high: r.high, low: r.low}
}
func ZeroValue[T idNameInfo]() IdRoot[T] {
	var t T
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(0xffffffffffff))
	buf[7] = t.Letter()
	if t.IsError() {
		buf[6] = 0x80
	}
	h := binary.LittleEndian.Uint64(buf)
	l := uint64(0xffffffffffffffff)
	return IdRoot[T]{high: h, low: l}
}

func NewIdRootError[T idNameInfo](code IdRootErrorCode) IdRoot[T] {
	var t T
	high := uint64(0)
	low := uint64(code)
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, high)
	buf[7] = t.Letter()
	buf[6] = 0 & 0x80
	high = binary.LittleEndian.Uint64(buf)
	id := IdRoot[T]{
		high: high,
		low:  low,
	}
	return id
}

func UnmarshalProtobuf[T idNameInfo](msg *protosupportmsg.IdRaw) (IdRoot[T], IdErr) {
	var t T
	h := msg.GetHigh()
	l := msg.GetLow()
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, h)
	if buf[7] != t.Letter() {
		return IdRoot[T]{0, 0}, NewIdErr(IdErrTypeMismatch)
	}
	if buf[7] == 0 {
		return IdRoot[T]{0, 0}, NewIdErr(IdErrNoType)
	}
	return UnmarshalRaw[T](IdRaw{high: h, low: l})
}

func UnmarshalRaw[T idNameInfo](r IdRaw) (IdRoot[T], IdErr) {
	return IdRoot[T](r), NoIdErr
}

func MarshalProtobuf[T idNameInfo](i IdRoot[T]) *protosupportmsg.IdRaw {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, i.high)
	h := binary.LittleEndian.Uint64(buf)
	binary.LittleEndian.PutUint64(buf, i.low)
	l := binary.LittleEndian.Uint64(buf)
	return &protosupportmsg.IdRaw{High: h, Low: l}
}

func (i IdRoot[T]) Short() string {
	var t T
	raw := i.Raw()
	low := raw.LowLE()
	high := raw.HighLE()

	if i.IsZeroValue() {
		return fmt.Sprintf("[%s-zero]", t.ShortString())
	}
	if high[6]&0x80 != 0 {
		if high[0] == 0 && high[1] == 0 {
			return fmt.Sprintf("[%s-NoErr]", t.ShortString())
		}
		return fmt.Sprintf("[%s-%02x%02x%02x%02x]", low[3], low[2], low[1], low[0])
	}
	return fmt.Sprintf("[%s-%02x%02x%02x]", t.ShortString(), low[2], low[1], low[0])
}

func (i IdRoot[T]) WriteGuestLe(mem api.Memory, offset uint32) bool {
	return IdRaw(i).WriteGuestLe(mem, offset)
}

func (i IdRoot[T]) String() string {
	var t T
	highByte := i.Raw().HighLE()
	highByte[6] &= 0x80 //high bit is true for an error type
	two := (i.high >> 32) & 0xffff
	four := i.high & 0xffffffff

	lowPart := make([]uint16, 4)
	lowPart[3] = uint16(i.low & 0xffff)
	lowPart[2] = uint16((i.low >> 16) & 0xffff)
	lowPart[1] = uint16((i.low >> 32) & 0xffff)
	lowPart[0] = uint16((i.low >> 48) & 0xffff)

	return fmt.Sprintf("%s-%04x-%08x:%04x-%04x-%04x-%04x", t.ShortString(), two, four,
		lowPart[0], lowPart[1], lowPart[2], lowPart[3])

}

func (i IdRoot[T]) Raw() IdRaw {
	return Raw(i)
}

func (i IdRoot[T]) Equal(other IdRoot[T]) bool {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, i.high)
	if i.IsZeroValue() || other.IsZeroValue() {
		return false // can't be compared to anything
	}
	return i.high == other.high && i.low == other.low
}
func (i IdRoot[T]) IsError() bool {
	return i.Raw().IsError()
}

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

func Raw[T idNameInfo](i IdRoot[T]) IdRaw {
	return IdRaw(i)
}

// NewRawId is dangerous in that it performs no checks about the validity of the
// data provided. Its use is discourage.  It usually needed when interacting
// with networking code since at the network, you simply receive 16 bytes of
// data and you need to turn it back into an id of some sort. The given bytes
// will be converted to LittleEndian.
func NewRawId(h, l uint64) IdRaw {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, h)
	hle := binary.LittleEndian.Uint64(buf)
	binary.LittleEndian.PutUint64(buf, l)
	lle := binary.LittleEndian.Uint64(buf)
	return IdRaw{high: hle, low: lle}
}

func (r IdRaw) LowLE() []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(r.low))
	return buf
}
func (r IdRaw) HighLE() []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(r.high))
	return buf
}

// String applied to an IdRaw not generally recommended.  It is far better to call it on an
// instance IdRoot[T] because that print out things nicely.  This is really
// only useful when you *dont* know the type, such as when you have read it off
// the wire.
func (r IdRaw) String() string {
	highByte := r.HighLE()
	key := highByte[7]
	highByte[6] &= 0x80 //high bit is true for an error type
	two := (r.high >> 32) & 0xffff
	four := r.high & 0xffffffff

	lowPart := make([]uint64, 4)
	lowPart[3] = r.low & 0xffff
	lowPart[2] = (r.low >> 16) & 0xffff
	lowPart[1] = (r.low >> 32) & 0xffff
	lowPart[0] = (r.low >> 48) & 0xffff

	return fmt.Sprintf("%c-%04x-%08x:%04x-%04x-%04x-%04x", key, two, four,
		lowPart[0], lowPart[1], lowPart[2], lowPart[3])
}

// IsError tests if this object 1) has the error mark set (byte 6 of high), so it is an error
// type and 2, that in the low order 2 bytes of low non zero.
func (r IdRaw) IsError() bool {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(r.high))
	if buf[6]&0x80 == 0 {
		return false
	}
	if buf[0] == 0 && buf[1] == 0 {
		return false
	}
	return true
}

// WriteGuestLe is a function that copies this into a memory space pointed to by
// by the offset into memory mem.
func (r IdRaw) WriteGuestLe(mem api.Memory, offset uint32) bool {
	if !mem.WriteUint64Le(offset, r.high) {
		return false
	}
	if !mem.WriteUint64Le(offset+8, r.low) {
		return false
	}
	return true
}
