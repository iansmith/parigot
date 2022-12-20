package splitutil

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/api/netconst"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/jspatch"

	"google.golang.org/protobuf/proto"
)

// SinglePayload is the data structure passed from the WASM portion of a split mode service to
// the go portion.  SinglePayload has to be structured in a way that allows different programming
// languages to use it.  The WASM portion of the service builds an instance of this structure somewhere
// in memory, usually on the heap,  The caller should allocate space for the result buffer to be filled
// in by parigot. parigot will fill in the OutLen slot with the actual size of the result.  The size
// of the result buffer (probably netconst.ReadBufferSize) should be set by the caller so parigot
// kernel is aware of the amount of space available for the result.  The ErrId will be filled in by
// parigot so the caller is notified of errors.
type SinglePayload struct {
	InPtr  int64
	InLen  int64
	OutPtr int64
	OutLen int64
	ErrPtr *[2]int64
}

var emptyReturnProto = protosupport.CallId{}.Id

// DecodeError is returned when the result from the go language portion of the service cannot be
// understood.
var DecodeError = errors.New("decoding error")

// newSinglePayload is used to allocate the space for the SinglePayload as well as the needed return values
// during a call from WASM to the go language side.  This code is run in WASM.
func newSinglePayload() *SinglePayload {
	buffer := make([]byte, netconst.ReadBufferSize)
	ptr, l := SliceToTwoInt64s(buffer)
	errorId := lib.NoKernelError()

	return &SinglePayload{
		InPtr:  0,
		InLen:  0,
		OutPtr: ptr,
		OutLen: l,
		ErrPtr: (*[2]int64)(unsafe.Pointer(&errorId.Id.Low)),
	}

}

// SendSingleProto is a utility function for initializing the InPtr and InLen fields defined in SinglePayload.
// This code is run in WASM.
func SendSingleProto(req proto.Message) (uintptr, error) {
	size := proto.Size(req)
	buffer := make([]byte, netconst.FrontMatterSize)
	if len(buffer) >= netconst.ReadBufferSize {
		panic("log message too large to fit in receive buffer:" + fmt.Sprint(len(buffer)))
	}
	binary.LittleEndian.PutUint64(buffer[:8], netconst.MagicStringOfBytes)
	binary.LittleEndian.PutUint32(buffer[8:netconst.FrontMatterSize], uint32(size))
	buffer, err := proto.MarshalOptions{}.MarshalAppend(buffer, req)
	if err != nil {
		return 0, err
	}
	result := crc32.Checksum(buffer[netconst.FrontMatterSize:netconst.FrontMatterSize+size], netconst.KoopmanTable)
	buffer = append(buffer, []byte{0, 0, 0, 0}...) //space for the crc
	binary.LittleEndian.PutUint32(buffer[netconst.FrontMatterSize+size:], uint32(result))
	payload := newSinglePayload()
	payload.InPtr, payload.InLen = SliceToTwoInt64s(buffer)
	// we have to convert this because pointers are only 32 bits in most wasm compilers
	u := uintptr(unsafe.Pointer(payload))
	return u, nil
}

// ExtractProtoFromBytes decodes a buffer obtained when the client side drops the payload (above) to us
// inside the go world.  The return value of T should NOT be touched/used if error!=nil.  This code
// is run in WASM.
func ExtractProtoFromBytes(buffer []byte, req proto.Message) error {
	m := binary.LittleEndian.Uint64(buffer[0:8])
	if m != netconst.MagicStringOfBytes {
		log.Printf("unable to print log message, bad magic number %x", m)
		return DecodeError
	}
	l := binary.LittleEndian.Uint32(buffer[8:12])
	if l >= uint32(netconst.ReadBufferSize) {
		log.Printf("unable to print log message, very large log message [%d bytes]", l)
		return DecodeError
	}
	size := int(l)

	objBuffer := buffer[netconst.FrontMatterSize : netconst.FrontMatterSize+size]
	if err := proto.Unmarshal(objBuffer, req); err != nil {
		log.Printf("unable to print log message, request could not be unmarshaled: %v", err)
		return DecodeError
	}
	result := crc32.Checksum(objBuffer, netconst.KoopmanTable)
	expected := binary.LittleEndian.Uint32(buffer[netconst.FrontMatterSize+size : netconst.FrontMatterSize+size+4])
	if expected != result {
		log.Printf("unable process data received on the go side, bad checksum found on request")
		return DecodeError
	}
	return nil

}

// ReadSlice is a utility for reading a slice, given its data offset into a payload pointed to be structPtr.
func ReadSlice(mem *jspatch.WasmMem, structPtr int64, dataOffset uintptr, lenOffset uintptr) []byte {
	return mem.LoadSliceWithLenAddr(int32(structPtr)+int32(dataOffset),
		int32(structPtr)+int32(lenOffset))
}

// SliceToTwoInt64s is a utility for taking an array of bytes (a buffer) and converting it to two int64s
// which can be put into a payload.  The capacity of the slice is ignored.
func SliceToTwoInt64s(b []byte) (int64, int64) {
	slh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return int64(slh.Data), int64(slh.Len)
}

// Write64BitPair is a utility for a pair of 64 bit ints, the contents of a lib.Id, into a payload structure
// pointed to be structPtr.
func Write64BitPair(mem *jspatch.WasmMem, structPtr int64, dataOffset uintptr, id lib.Id) {
	derefed := mem.GetInt32(int32(structPtr + int64(dataOffset)))
	// write the error info back to client
	mem.SetInt64(derefed, int64(id.Low()))
	mem.SetInt64(derefed+8, int64(id.High()))
}
