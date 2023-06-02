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

type NameInfo interface {
	ShortString() string
	Letter() byte
	IsError() bool
}

//
// ID ROOT
//

type IdRoot[T NameInfo] struct {
	high, low uint64
}

func NewIdRoot[T NameInfo]() IdRoot[T] {
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
func NewIdRootFromRaw[T NameInfo](r IdRaw) IdRoot[T] {
	return IdRoot[T](r)
}

// NewRawId is dangerous in that it performs no checks about the validity of the
// data provided. Its use is discouraged.  It usually needed when interacting
// with networking code or when you must carefully choose id values for interacting
// with other systems (like a database). E.g. The Queue service uses a database to
// hold it's queue records, and so it forces the low order 64 bits to be the row id
// of the queue in question.
func NewIdTyped[T NameInfo](h, l uint64) IdRoot[T] {
	var t T
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, h)
	buf[7] = t.Letter()
	buf[6] = 0 // no low order bitset
	if t.IsError() {
		buf[6] = 0x80
	}
	upper := binary.LittleEndian.Uint64(buf)
	upper |= (h & 0xffffffffffff)
	return IdRoot[T](NewRawId(upper, l))
}

func ZeroValue[T NameInfo]() IdRoot[T] {
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

func NewIdRootError[T NameInfo](code IdRootErrorCode) IdRoot[T] {
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

func (f IdRoot[T]) ErrorCode() uint16 {
	low := f.Low() & 0xffff
	return uint16(low)
}

func UnmarshalProtobuf[T NameInfo](msg *protosupportmsg.IdRaw) (IdRoot[T], IdErr) {
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

func UnmarshalRaw[T NameInfo](r IdRaw) (IdRoot[T], IdErr) {
	return IdRoot[T](r), NoIdErr
}

func MarshalProtobuf[T NameInfo](i IdRoot[T]) *protosupportmsg.IdRaw {
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
	if i.IsEmptyValue() {
		return fmt.Sprintf("[%s-empty]", t.ShortString())
	}
	if high[6]&0x80 != 0 {
		if high[0] == 0 && high[1] == 0 {
			return fmt.Sprintf("[%s-NoErr]", t.ShortString())
		}
		return fmt.Sprintf("[%s-%02x%02x%02x%02x]", t.ShortString(), low[3], low[2], low[1], low[0])
	}
	return fmt.Sprintf("[%s-%02x%02x%02x]", t.ShortString(), low[2], low[1], low[0])
}

func (i IdRoot[V]) WriteGuestLe(mem api.Memory, offset uint32) bool {
	raw := i.Raw()
	return raw.WriteGuestLe(mem, offset)
}

func (i IdRoot[T]) String() string {
	var t T
	twohigh := (i.high >> 32) & 0xffff
	fourlow := (i.high) & 0xffffffff

	lowPart := make([]uint32, 4)
	lowPart[3] = uint32(i.low & 0xffffff)
	lowPart[2] = uint32((i.low >> 24) & 0xff)
	lowPart[1] = uint32((i.low >> 32) & 0xffff)
	lowPart[0] = uint32((i.low >> 48) & 0xffff)

	return fmt.Sprintf("%s-%04x-%08x:%04x-%04x-%02x-%06x",
		t.ShortString(), twohigh, fourlow,
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

// IsEmptyValue tells you if the given id is actually just full of zeros in
// the 14 data bytes.  This almost means that the caller gave you a bad it, since
// the chance of all 14 data bytes is very low, and vastly lower than somebody
// treating a KernelNoErr as an id value.
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
	return i.Raw().High()
}

// Use of Low() is not recommended for user code.  It returns the
// low 8 bytes of the 128bit id.
func (i IdRoot[T]) Low() uint64 {
	return i.Raw().Low()
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

func Raw[T NameInfo](i IdRoot[T]) IdRaw {
	return IdRaw(i)
}

// NewRawId is dangerous in that it performs no checks about the validity of the
// data provided. Its use is discouraged.  It usually needed when interacting
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
// the wire and are debugging.
func (i IdRaw) String() string {
	return "raw-128bit"
}

// Use of Low() is not recommended.  It returns the low bytes (8) of the 128 bit value.
// The only real need for this if you are dealing with a storage mechanism
// (like a database) that cannot handle 128 bit values, but you want to store
// the the id in two parrts, low and high.
func (i IdRaw) Low() uint64 {
	return i.low
}

// Use of High() is not recommended.  It returns the high bytes (8) of the 128 bit value.
// The only real need for this if you are dealing with a storage mechanism
// (like a database) that cannot handle 128 bit values, but you want to store
// the the id in two parrts, low and high.
func (i IdRaw) High() uint64 {
	return i.high
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
