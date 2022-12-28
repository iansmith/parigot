package splitutil

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
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
	ErrPtr [2]int64
}

var kerrNone = lib.Unmarshal(lib.NoKernelError())

// DecodeError is returned when the result from the go language portion of the service cannot be
// understood.
var DecodeError = errors.New("decoding error")

// newSinglePayload is used to allocate the space for the SinglePayload as well as the needed return values
// during a call from WASM to the go language side.  This code is run in WASM.
func newSinglePayload() *SinglePayload {
	buffer := make([]byte, netconst.ReadBufferSize)
	ptr, l := SliceToTwoInt64s(buffer)

	sp := &SinglePayload{
		InPtr:  0,
		InLen:  0,
		OutPtr: ptr,
		OutLen: l,
		ErrPtr: [2]int64{int64(kerrNone.High()), int64(kerrNone.Low())},
	}
	return sp
}

// This is the top level entry point for the WASM side.  It sends the proto given and fills in the resp
// provided.  It will return a lib.Id that is a KernelErrorId, an error or two nils.  If it is a
// lib.Id (1st) that error has been returned from the "other side" and the resp has not been filled in
// due to the error. If the error (2nd) is returned, then there was a problem sending the req so
// nothing has been sent at this point, or there was a problem unpacking the returned value.   In the
// latter case, the returned error will be DecodeError and the Go side responded without error but
// sent a bogus buffer of values.  Put another way, the error (2nd return value) cases are about
// problems that occurred while running wasm code and the first return value is about problems that
// occurred running go code. If this function returns two nils, everything went through the happy path
// and resp has been filled in with the object sent from the go side.  If the go side has nothing
// to return, the resp object is left unchanged.
func SendReceiveSingleProto(req, resp proto.Message, fn func(int32)) (lib.Id, error) {
	u, err := SendSingleProto(req)
	if err != nil {
		return nil, err
	}
	// u is a pointer to the SinglePayload, send payload through fn
	fn(int32(u))
	// check to see if this is an returned error
	ptr := (*SinglePayload)(unsafe.Pointer(u))
	errRtn := lib.NewFrom64BitPair[*protosupport.KernelErrorId](uint64(ptr.ErrPtr[0]), uint64(ptr.ErrPtr[1]))
	if !errRtn.Equal(kerrNone) {
		return errRtn, nil
	}
	// if they returned nothing, we are done
	if ptr.OutLen == 0 {
		return nil, nil
	}
	var byteBuffer []byte
	wasmSideSlice := (*reflect.SliceHeader)(unsafe.Pointer(&byteBuffer))
	wasmSideSlice.Data = uintptr(ptr.OutPtr)
	wasmSideSlice.Len = int(ptr.OutLen)
	wasmSideSlice.Cap = int(ptr.OutLen)

	err = DecodeSingleProto(byteBuffer, resp)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// SendSingleProto is a utility function creating an initializng the InPtr and InLen fields based on a
// single protobuf object.  This code is run in WASM.  This function uses encodeSingleProto for converting
// the proto->bytes. It returns either a pointer to the 32bit addr of the resulting payload structure or
// an error.
func SendSingleProto(req proto.Message) (uintptr, error) {
	size := proto.Size(req)

	if size+netconst.TrailerSize+netconst.FrontMatterSize >= netconst.ReadBufferSize {
		return 0, errors.New("log message too large to fit in receive buffer:" + fmt.Sprint(size))
	}
	payload := newSinglePayload()
	u := uintptr(unsafe.Pointer(payload))

	var err error
	// when calling encodeSingleProto from the WASM side here, we are initializing the SinglePayload so
	// we just copy the returned values into that structure.
	buffer, err := encodeSingleProto(req, size)
	if err != nil {
		return 0, err
	}
	payload.InPtr, payload.InLen = SliceToTwoInt64s(buffer)
	return u, nil
}

// encodeSingleProto takes a given req and encodes it with all the trimmings include the preceding magic
// bytes and the 32 bit crc at the end.  This probably is overkill for things being transferred from a client
// program in the same address space, as inadvertent errors are highly unlikely. It returns the content of
// the encoded value or an error.
//
// This code is called BOTH from the go side and the wasm side.
func encodeSingleProto(req proto.Message, size int) ([]byte, error) {
	//frontmatter
	buffer := make([]byte, netconst.FrontMatterSize)
	binary.LittleEndian.PutUint64(buffer[:8], netconst.MagicStringOfBytes)
	binary.LittleEndian.PutUint32(buffer[8:netconst.FrontMatterSize], uint32(size))
	// append network form
	buffer, err := proto.MarshalOptions{}.MarshalAppend(buffer, req)
	if err != nil {
		return nil, err
	}
	// compute checksum
	result := crc32.Checksum(buffer[netconst.FrontMatterSize:netconst.FrontMatterSize+size], netconst.KoopmanTable)
	buffer = append(buffer, []byte{0, 0, 0, 0}...) //space for the crc
	// put checksum in buffer
	binary.LittleEndian.PutUint32(buffer[netconst.FrontMatterSize+size:], uint32(result))

	return buffer, nil
}

// RespondEmpty is a way of signaling that you have no data to return, given a SinglePayload.
// This is run on the GO side.  Since the payload is initalized with a return error that is "no kernel error"
// this function leaves that field of the SinglePayload unmodified.
//
// Note that the wasmptr provided is a pointer into the WASM address space.
func RespondEmpty(mem *jspatch.WasmMem, sp int32) {
	wasmPtr := int32(mem.GetInt64(sp + 8))

	offset := int32(unsafe.Offsetof(SinglePayload{}.OutLen))
	mem.SetInt64(wasmPtr+offset, 0)
}

// RespondSingleProto is a way of returning a single proto (flattened to bytes) to the WASM world via
// payload.  This is called by the GO side.  This function will set the error return if it cannot fit
// the response message into space allocated in the payload or if the encode fails.  This function
// uses encodeSingleProto for doing the conversion proto->bytes.
//
// Note that the wasmPtr provided is the address in the WASM address space.
func RespondSingleProto(mem *jspatch.WasmMem, sp int32, resp proto.Message) {
	wasmPtr := int32(mem.GetInt64(sp + 8))

	size := proto.Size(resp)
	fullSize := int64(netconst.TrailerSize + netconst.FrontMatterSize + size)

	// how much space do we have?
	offset := int32(unsafe.Offsetof(SinglePayload{}.OutLen))
	available := mem.GetInt64(wasmPtr + offset)

	// can't fit the response?
	if fullSize >= available {
		ErrorResponse(mem, wasmPtr, lib.KernelDataTooLarge)
		return
	}

	// encode the proto into a buffer... this resulting pointer to the buffer could be 32 or 64 bits
	buffer, err := encodeSingleProto(resp, size)
	if err != nil {
		// this can only happen on some type of protobuf encoding issue
		ErrorResponse(mem, wasmPtr, lib.KernelMarshalFailed)
		return
	}

	ptrOffset := unsafe.Offsetof(SinglePayload{}.OutPtr)
	// this is tricky: we have to COPY the bytes from the go side to the wasm side bc the pointer
	// returned as buffer is in the GO address space
	CopyToPtr(mem, int64(wasmPtr), ptrOffset, buffer)

	// tell the caller the length
	mem.SetInt64(wasmPtr+offset, fullSize)
}

// DecodeSingleProto decodes a buffer obtained when the client side drops the payload (above) to us
// inside the go world. This code is run on BOTH the go and WASM sides to extract a protobuf object from
// the bytes.  Because this can run on the WASM side, it is not safe to use fmt.Printf() or log.Printf().
// This expects the buffer to have been encoded by encodeSingleProto and thus includes the extra bells
// and whistles like a CRC and a magic number.
//
// Note: you must pass the pointer to an allocated and empty protobuf structure here as the obj.
func DecodeSingleProto(buffer []byte, obj proto.Message) error {
	m := binary.LittleEndian.Uint64(buffer[0:8])
	if m != netconst.MagicStringOfBytes {
		return DecodeError
	}
	l := binary.LittleEndian.Uint32(buffer[8:12])
	if l >= uint32(netconst.ReadBufferSize) {
		return DecodeError
	}
	size := int(l)

	objBuffer := buffer[netconst.FrontMatterSize : netconst.FrontMatterSize+size]
	if err := proto.Unmarshal(objBuffer, obj); err != nil {
		return DecodeError
	}
	result := crc32.Checksum(objBuffer, netconst.KoopmanTable)
	expected := binary.LittleEndian.Uint32(buffer[netconst.FrontMatterSize+size : netconst.FrontMatterSize+size+4])
	if expected != result {
		return DecodeError
	}
	return nil

}

// ErrorResponse takes the pointer to the stack (in the WASM address space) and sets
// the error contained in it to the code provided.  This is called by the GO side.  This should be
// use to signal errors back to the WASM code.
func ErrorResponse(mem *jspatch.WasmMem, sp int32, code lib.KernelErrorCode) {
	wasmPtr := mem.GetInt64(sp + 8)
	kerr := lib.NewKernelError(code)
	// the [0] value is the high 8 bytes, the [1] the low 8 bytes
	mem.SetInt64(int32(wasmPtr)+int32(unsafe.Offsetof(SinglePayload{}.ErrPtr)),
		int64(kerr.High()))
	// high is 8 bytes higher
	mem.SetInt64(int32(wasmPtr)+int32(unsafe.Offsetof(SinglePayload{}.ErrPtr)+8),
		int64(kerr.Low()))
	print(fmt.Sprintf("xxx Error Response: we just set the values to %x,%x\n", kerr.High(), kerr.Low()), "\n")
}

// StackPointerToRequest assumes that this function was passed a pointer to the WASM stack and that
// this code is running on the GO side. This function assumes that 8 bytes above the wasm SP is a pointer
// to a SinglePayload.  Again, the GO side runs this to pull data from the WASM side.  This
// fills in the proto with the data sent from the WASM side or it returns an error.  An error here
// means we could not decode the package sent from the WASM side.  If an error is returned, the
// caller can simply return, as the payload has been modified to tell the other side about the error.
func StackPointerToRequest(mem *jspatch.WasmMem, sp int32, req proto.Message) error {
	wasmPtr := mem.GetInt64(sp + 8)

	buffer := ReadSlice(mem, wasmPtr,
		unsafe.Offsetof(SinglePayload{}.InPtr),
		unsafe.Offsetof(SinglePayload{}.InLen))

	err := DecodeSingleProto(buffer, req)
	if err != nil {
		ErrorResponse(mem, int32(wasmPtr), lib.KernelUnmarshalFailed)
		return err
	}
	return nil
}

//
// Utilities
//

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

func CopyToPtr(mem *jspatch.WasmMem, structPtr int64, dataOffset uintptr, content []byte) {
	mem.CopyToPtr(int32(structPtr)+int32(dataOffset), content)
}
