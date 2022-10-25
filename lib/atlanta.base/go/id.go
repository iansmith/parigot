package lib

import (
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/iansmith/parigot/g/parigot"
	"github.com/iansmith/parigot/lib/interface_"
)

const service = 's'

const locateError = 'l'
const dispatchError = 'd'
const registerError = 'r'

const kernelServiceConst = 1

func init() {
	rand.New(rand.NewSource(0)) // reproducibility for dev
}

// 128 raw bits, 112 in use as number 16 as type (byte + reserved byte)
type id struct {
	high uint64
	low  uint64
}

// Short returns a short version of an id, like [s-1ae2].
func (i id) Short() string {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(buf, uint64(i.high))
	b := buf[7]
	buf[7] = 0
	buf[6] = 0
	buf[5] = 0
	buf[4] = 0
	buf[3] = 0
	buf[2] = 0
	return bufToPrintedRepSmall(b, buf)
}

func bufToPrintedRepSmall(b byte, buf []byte) string {
	found := -1
	for i := binary.MaxVarintLen64 - 1; i >= 0; i-- {
		if buf[i] != 0 {
			found = i
			break
		}
	}
	if found == -1 {
		switch b {
		case dispatchError, registerError, locateError:
			return fmt.Sprintf("[%c-NoError]", b)
		default:
			panic("should not have a zero value with id type " + fmt.Sprint(b))
		}
	}
	// get the small into out of the buf
	val := binary.LittleEndian.Uint64(buf[0 : found+1])
	return fmt.Sprintf("[%c-%04x]", val)
}

// String returns a long string that uniquely identifies this id.  This is
// usualy something like r-xx-xxxxxxxx-xxxxxxxx-xxxx-xxxx where all the x's
// are hex digits.  Note that Short() is the equivalent of the first and last
// five characters of string.  If the number is a small integer, the leading
// zeros are omitted.
func (i id) String() string {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(buf, uint64(i.high))
	b := buf[7]
	buf[7] = 0
	buf[6] = 0 //reserved for future expansion
	num := binary.LittleEndian.Uint64(buf)
	result := ""
	if num != 0 {
		leftover := binary.LittleEndian.Uint64(buf[4:6])
		rest := binary.LittleEndian.Uint64(buf[0:3])
		result = fmt.Sprintf("[%c-%02x-%04x-", b, leftover, rest)
	}
	binary.LittleEndian.PutUint64(buf, uint64(i.low))
	if num == 0 {
		return bufToPrintedRepSmall(b, buf)
	}
	// normal big case
	upper := binary.LittleEndian.Uint64(buf[4:8])
	lower := binary.LittleEndian.Uint64(buf[0:4])
	result += fmt.Sprintf("%04x-%04x]", upper, lower)
	return result
}

// Error() returns true if the type of this id is an error type and there
// is an error value.  0 is the only non-error value.
func (i id) Error() bool {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(buf, uint64(i.high))
	b := buf[7]
	switch b {
	case dispatchError, registerError, locateError:
	default:
		panic("cannot test this id type for error " + fmt.Sprint(b))
	}
	return i.low != 0
}

// Returns the name of the type of id, like "service" or "locate error"
func (i id) Type() string {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(buf, uint64(i.high))
	b := buf[7]
	switch b {
	case dispatchError:
		return "dispatch error"
	case registerError:
		return "register error"
	case locateError:
		return "locate error"
	case service:
		return "service"
	default:
		panic("unknown id type " + fmt.Sprint(b))
	}
}

// Equal returns true if the two ids are of the same type and have the same number.
func (i id) Equal(otherId interface_.Id) bool {
	bufI := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(bufI, uint64(i.high))
	bI := bufI[7]
	other := otherId.(id)
	bufOther := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(bufOther, uint64(other.high))
	bOther := bufOther[7]
	if bI != bOther {
		return false
	}
	return i.high == other.high && i.low == other.low
}

// ServiceId is returned when the locator succeeds at finding your service.
type ServiceId interface_.Id

// LocateError is return when the locator failed to find your service or
// had other problems.  It returns well known constants for the different
// types of errors it has.
type LocateError interface_.Id

// DispatchError is returned by a failed call to Dispatch in the ABI.
// It returns a well known constant for its errors.
type DispatchError interface_.Id

// RegisterError is returned by a failed call to Register in the ABI.
// It returns a well known constant for its errors.
type RegisterError interface_.Id

// NewService creates a new, random service id.
func NewService() ServiceId {
	return ServiceId(newId(service))
}

// NewService from int creates a service with the given number.  Note that the
// kernel's id is always 1.
func NewServiceFromInt(i byte) ServiceId {
	return idFromInt(service, i)
}

// NewDispatchError creates a dispatch error type from the given error code.
func NewDispatchError(errorCode byte) DispatchError {
	return DispatchError(idFromInt(dispatchError, errorCode))
}

// NewLocateError takes an error code and converts it to a Locator error.
func NewLocateError(errorCode byte) LocateError {
	return LocateError(idFromInt(locateError, errorCode))
}

// newId creates a new random Id with the given type
func newId(typ byte) interface_.Id {
	return newIdFromRaw(typ, rand.Uint64(), rand.Uint64())
}
func newIdFromRaw(typ byte, high uint64, low uint64) interface_.Id {
	highBuf := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(highBuf, high)
	highBuf[7] = typ
	highBuf[6] = 0 //reserved
	high = binary.LittleEndian.Uint64(highBuf)
	return id{high: high, low: low}
}

// idFromInt is used by code that wants to create an id in the right format, but with
// a given value, like an error code.
func idFromInt(s byte, value byte) id {
	i := uint64(value)
	var high uint64
	var low uint64
	highBuf := make([]byte, binary.MaxVarintLen64)
	lowBuf := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(lowBuf, uint64(i))
	low = binary.LittleEndian.Uint64(lowBuf)
	highBuf[7] = s
	highBuf[6] = 0 // reserved
	highBuf[5] = 0
	highBuf[4] = 0
	highBuf[3] = 0
	highBuf[2] = 0
	highBuf[1] = 0
	highBuf[0] = 0
	high = binary.LittleEndian.Uint64(highBuf)
	return id{high: high, low: low}
}

// FromServiceId converts from the wrapped protobuf version of service id
// to the normal version.
func FromServiceId(serviceId parigot.ServiceId) ServiceId {
	return newIdFromRaw(service, serviceId.GetHigh(), serviceId.GetLow())
}
