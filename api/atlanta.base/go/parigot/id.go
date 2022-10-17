package parigot

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano())) // reproducibility
}

const service = 's'
const locateError = 'l'
const dispatchError = 'd'

// AnyId is something using the Id format that is 64 bits, with the first 8
// being an ascii character. So given the hex value of an id you at least
// know what type of thing it is.  The other 56 bits are typically random.
type AnyId int64

// ServiceId is returned when the locator succeeds at finding your service.
type ServiceId AnyId

// LocateErrorId is return when the locator failed to find your service or
// had other problems.  It returns well known constants for the different
// types of errors it has.  Note that a call to Locate() in the ABI should
// be checked to see if the return value is a ServiceId or a LocateErrorId.
type LocateErrorId AnyId

// DispatchErrorId is returned by a failed call to Dispatch in the ABI.  Note
// that Dispatch always returns a buffer of bytes, but they will only be the
// correct size and format of a DispatchErrorId in the case of error.
type DispatchErrorId AnyId

// NewServiceId is used by the kernel to create a new service id in production.
// The value is generated randomly (56 bits).
func NewServiceId() ServiceId {
	return ServiceId(newId(service))
}

// NewServiceIdFromInt should only be used by the kernel to create a service
// id when debugging.
func NewServiceIdFromInt(i int64) ServiceId {
	return ServiceId(idFromInt(service, i))
}

// NewDispatchError requires an error code on the server side.  This is used
// by the Dispatch implentation to create an error id.
func NewDispatchErrorId(errorCode int64) DispatchErrorId {
	return DispatchErrorId(idFromInt(dispatchError, errorCode))
}

// NewLocateErrorId takes a previously thought to be ServiceId and changes
// it to a LocateErrorId.
func NewLocateErrorId(oldValue int64) LocateErrorId {
	return LocateErrorId(idFromInt(locateError, oldValue))
}

// IsServiceId is used by the client side service to test if an id is a
// service id.
func IsServiceId(v int64) bool {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(buf, uint64(v))
	return buf[7] == locateError
}

// NewDispatchErrorFromBytes is used by the ClientSideService to create an
// error from a bundle of bytes it got back from Dispatch.
func NewDispatchErrorFromBytes(b []byte) DispatchErrorId {
	if len(b) != 8 {
		panic(fmt.Sprintf("unable to understand dispatch error (%d bytes)"))
	}
	b[7] = dispatchError
	u := binary.LittleEndian.Uint64(b)
	return DispatchErrorId(u)
}

// BytesFromDispatchError is called by the kernel to convert a DispatchErrorId to a
// buffer of bytes.
func BytesFromDispatchError(id DispatchErrorId) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(buf, uint64(id))
	return buf
}

// AsIdShort formats the last two bytes and the letter indicator for printing.
// This is used for debugging sessions, not for production.
// Two bytes means 64K combinations so you have to very unlucky to get a collision.
func AsIdShort(n AnyId) string {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(buf, uint64(n))
	buf[6] = 0
	buf[5] = 0
	buf[4] = 0
	buf[3] = 0
	buf[2] = 0
	var b byte
	switch buf[7] {
	case service:
	case locateError:
	case dispatchError:
		b = buf[7]
	default:
		panic("unable to understand id (expected s,l,d)")
	}
	buf[7] = 0
	u := binary.LittleEndian.Uint64(buf)
	return fmt.Sprintf("%c-%04x", b, u)
}

// AsId formats an id with its letter prefix.  It uses the full 56 bits of the
// value and is suitable for logs in production.
func AsId(n AnyId) string {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(buf, uint64(n))
	var b byte
	switch buf[7] {
	case service:
	case locateError:
	case dispatchError:
		b = buf[7]
	default:
		panic("unable to understand id (expected s,l,d)")
	}
	u := binary.LittleEndian.Uint64(buf)
	return fmt.Sprintf("%c-%x", b, u)
}

// newId returns the new Id that has been generated randomly.
func newId(s byte) AnyId {
	u := rand.Uint64()
	return idFromInt(s, int64(u))
}

func idFromInt(s byte, i int64) AnyId {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(buf, uint64(i))
	buf[7] = s
	x := binary.LittleEndian.Uint64(buf)
	return AnyId(x)
}
