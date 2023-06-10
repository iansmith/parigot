package id

import (
	"encoding/binary"
	"fmt"
	"math/rand"

	protosupportmsg "github.com/iansmith/parigot/g/protosupport/v1"
	"github.com/tetratelabs/wazero/api"
)

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

func (f IdRoot[T]) ErrorCode() int32 {
	low := f.Low() & 0xffffffff
	return int32(low)
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
		if high[6]&0x80 != 0 {
			return fmt.Sprintf("[%s-NoErr]", t.ShortString())
		}
		return fmt.Sprintf("[%s-empty]", t.ShortString())
	}
	if high[6]&0x80 != 0 {
		if low[0] == 0 && low[1] == 0 {
			return fmt.Sprintf("[%s-NoErr]", t.ShortString())
		}
		return fmt.Sprintf("[%s-%02x%02x]", t.ShortString(), low[1], low[0])
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

type IdEqualer interface {
	EqualId(other IdEqualer) bool
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
	fourhigh := (i.high >> 32) & 0xffffffff
	fourlow := (i.high) & 0xffffffff

	lowPart := make([]uint32, 4)
	lowPart[3] = uint32(i.low & 0xffffff)
	lowPart[2] = uint32((i.low >> 24) & 0xff)
	lowPart[1] = uint32((i.low >> 32) & 0xffff)
	lowPart[0] = uint32((i.low >> 48) & 0xffff)

	return fmt.Sprintf("%s-%08x-%08x:%04x-%04x-%02x-%06x",
		"raw", fourhigh, fourlow,
		lowPart[0], lowPart[1], lowPart[2], lowPart[3])
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
	if r.low&0xffff == 0 {
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
