package splitutil

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"log"
	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/apiimpl/netconst"
	lib "github.com/iansmith/parigot/lib/go"
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
	InPtr        int64
	InLen        int64
	OutPtr       int64
	OutLen       int64
	ErrPtr       [2]int64
	ErrDetailLen int64
	ErrDetail    int64
}

var kerrNone = lib.Unmarshal(lib.NoKernelError())

var callImpl lib.Call

// NewSinglePayload is used to allocate the space for the SinglePayload as well as the needed return values
// during a call from WASM to the go language side.  This code is run in WASM.
func NewSinglePayload() *SinglePayload {
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

func IsErrorInSinglePayload(ptr *SinglePayload) bool {
	errRtn := lib.IdRepresentsError(uint64(ptr.ErrPtr[0]), uint64(ptr.ErrPtr[1]))
	return errRtn
}

// This is the top level entry point for the WASM side.  It sends the proto given and fills in the resp
// provided.  It uses the c parameter to get access to the system calls
// and the fn to implemented the desired functionality after decoding
// req and before encoding resp.  It returns two logical entities, a boolean indicating
// if the error has already been sent and an error pair describing the error.  If
// the boolean is true, then the error pair is really just informational because
// the error is ALREADY recorded in the return object.  It returns false,nil,""
// if everything is ok.  If the bool is true, it means that the implementation of
// the function fn actually went ahead and set the error code during its execution.

// this function's signature was updated because we no longer use errors
// on the client side (in the go sense) until we are ready to go back to user
// code.  thus, we care that the internal response (ptr, SinglePayload) has the
// error code.
func SendReceiveSingleProto(c lib.Call, req, resp proto.Message, fn func(int32)) *SinglePayload {
	if callImpl == nil {
		callImpl = c
	}
	spayload := SendSingleProto(req)
	if IsErrorInSinglePayload(spayload) {
		return spayload
	}
	u := uintptr(unsafe.Pointer(spayload))
	fn(int32(u))
	// check to see if this is an returned error
	//print(fmt.Sprintf("WASM SIDE found an error?? in error ptr 0x%x,0x%x", ptr.ErrPtr[0], ptr.ErrPtr[1]))
	if IsErrorInSinglePayload(spayload) {
		return spayload
	}
	// if they returned nothing, we are done
	if spayload.OutLen == 0 {
		return spayload
	}
	var byteBuffer []byte
	wasmSideSlice := (*reflect.SliceHeader)(unsafe.Pointer(&byteBuffer))
	wasmSideSlice.Data = uintptr(spayload.OutPtr)
	wasmSideSlice.Len = int(spayload.OutLen)
	wasmSideSlice.Cap = int(spayload.OutLen)
	id, detail := DecodeSingleProto(byteBuffer, resp)
	if id != nil {
		formatErrorResult(spayload, id, detail)
	}
	return spayload
}

// SendSingleProto is a utility function creating an initializng the InPtr and InLen fields based on a
// single protobuf object.  This code is run in WASM.  This function uses encodeSingleProto for converting
// the proto->bytes. It returns either a pointer to the 32bit addr of the resulting payload structure or
// an error pair.
func SendSingleProto(req proto.Message) *SinglePayload {
	size := proto.Size(req)
	if size > 0 {
		if size+netconst.TrailerSize+netconst.FrontMatterSize >= netconst.ReadBufferSize {
			retVal := NewSinglePayload()
			kid := lib.NewKernelError(lib.KernelDataTooLarge)
			formatErrorResult(retVal, kid, "not enough space for call argument")
			return retVal
		}
	}
	payload := NewSinglePayload()

	// when calling encodeSingleProto from the WASM side here, we are initializing the SinglePayload so
	// we just copy the returned values into that structure.
	buffer, id, detail := encodeSingleProto(req, size)
	if id != nil {
		formatErrorResult(payload, id, detail)
		return payload
	}
	payload.InPtr, payload.InLen = SliceToTwoInt64s(buffer)
	return payload
}

func formatErrorResult(s *SinglePayload, id lib.Id, msg string) {
	s.ErrPtr[0] = int64(id.High())
	s.ErrPtr[1] = int64(id.Low())
	buffer := []byte(msg)
	s.ErrDetailLen = int64(len(buffer))
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buffer))
	s.ErrDetail = int64(sh.Data)

}

// encodeSingleProto takes a given req and encodes it with all the trimmings include the preceding magic
// bytes and the 32 bit crc at the end.  This probably is overkill for things being transferred from a client
// program in the same address space, as inadvertent errors are highly unlikely. It returns the content of
// the encoded value or an error pair.
//
// This code is called BOTH from the go side and the wasm side.
func encodeSingleProto(req proto.Message, size int) ([]byte, lib.Id, string) {
	//frontmatter
	buffer := make([]byte, netconst.FrontMatterSize)
	binary.LittleEndian.PutUint64(buffer[:8], netconst.MagicStringOfBytes)
	binary.LittleEndian.PutUint32(buffer[8:netconst.FrontMatterSize], uint32(size))
	// append network form
	buffer, err := proto.MarshalOptions{}.MarshalAppend(buffer, req)
	if err != nil {
		return nil, lib.NewKernelError(lib.KernelEncodeError),
			fmt.Sprintf("unable to marshal request:%v", err)
	}
	// compute checksum
	result := crc32.Checksum(buffer[netconst.FrontMatterSize:netconst.FrontMatterSize+size], netconst.KoopmanTable)
	buffer = append(buffer, []byte{0, 0, 0, 0}...) //space for the crc
	// put checksum in buffer
	binary.LittleEndian.PutUint32(buffer[netconst.FrontMatterSize+size:], uint32(result))
	return buffer, nil, ""
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

	log.Printf("xxx --- Respond Single Proto1: %#v\n", resp)
	size := proto.Size(resp)
	fullSize := int64(netconst.TrailerSize + netconst.FrontMatterSize + size)
	// how much space do we have?
	offsetForLen := int32(unsafe.Offsetof(SinglePayload{}.OutLen))
	available := mem.GetInt64(wasmPtr + offsetForLen)

	// can't fit the response?
	if fullSize >= available {
		ErrorResponse(mem, wasmPtr, lib.NewKernelError(lib.KernelDataTooLarge),
			fmt.Sprintf("unable to fit data of size %d in buffer of size %d", fullSize, available))
		return
	}
	// encode the proto into a buffer... this resulting pointer to the buffer could be 32 or 64 bits
	buffer, id, detail := encodeSingleProto(resp, size)
	if id != nil {
		// this can only happen on some type of protobuf encoding issue
		log.Printf("xxx --- Respond Single Proto2A: %s\n", id)
		ErrorResponse(mem, wasmPtr, id, detail)
		return
	}
	u := unsafe.Pointer(&buffer)
	sh := (*reflect.SliceHeader)(u)
	spayload := (*SinglePayload)(unsafe.Pointer(sh.Data))

	ptrOffset := unsafe.Offsetof(SinglePayload{}.OutPtr)
	// this is tricky: we have to COPY the bytes from the go side to the wasm side bc the pointer
	// returned as buffer is in the GO address space
	CopyToPtr(mem, int64(wasmPtr), ptrOffset, buffer)

	log.Printf("xxx --- Respond Single Proto3: InPtr 0x%0x, InLen 0x%0x, ErrPtr[0] 0x%0x, ErrPtr[1] 0x%0x, ErrDetailLen 0x%0x, ErrDetail: 0x%0x\n",
		spayload.InPtr, spayload.InLen, spayload.ErrPtr[0], spayload.ErrPtr[1], spayload.ErrDetailLen, spayload.ErrDetail)
	// tell the caller the length
	mem.SetInt64(wasmPtr+offsetForLen, fullSize)
}

// DecodeSingleProto decodes a buffer obtained when the client side drops the payload (above) to us
// inside the go world. This code is run on BOTH the go and WASM sides to extract a protobuf object from
// the bytes.  Because this can run on the WASM side, it is not safe to use fmt.Printf() or log.Printf().
// This expects the buffer to have been encoded by encodeSingleProto and thus includes the extra bells
// and whistles like a CRC and a magic number.
//
// Note: you must pass the pointer to an allocated and empty protobuf structure here as the obj.
func DecodeSingleProto(buffer []byte, obj proto.Message) (lib.Id, string) {
	m := binary.LittleEndian.Uint64(buffer[0:8])
	if m != netconst.MagicStringOfBytes {
		return lib.NewKernelError(lib.KernelDecodeError), "Unable to find magic byte sequence"
	}
	l := binary.LittleEndian.Uint32(buffer[8:12])
	if l >= uint32(netconst.ReadBufferSize) {
		return lib.NewKernelError(lib.KernelDecodeError),
			fmt.Sprintf("Size of buffer to decode (%d) is too large (max is %d)", l, netconst.ReadBufferSize)
	}
	size := int(l)
	if size == 0 {
		buffer = nil
		return nil, ""
	}
	objBuffer := buffer[netconst.FrontMatterSize : netconst.FrontMatterSize+size]
	if err := proto.Unmarshal(objBuffer, obj); err != nil {
		return lib.NewKernelError(lib.KernelUnmarshalFailed), fmt.Sprintf("unable to unmarshal encoded object: %v", err)
	}
	result := crc32.Checksum(objBuffer, netconst.KoopmanTable)
	expected := binary.LittleEndian.Uint32(buffer[netconst.FrontMatterSize+size : netconst.FrontMatterSize+size+4])
	if expected != result {
		return lib.NewKernelError(lib.KernelDecodeError), "CRC check failed for bundle"
	}
	return nil, ""

}

// ErrorResponse takes the pointer to the stack (in the WASM address space) and sets
// the error contained in it to the code provided.  This is called by the GO side.  This should be
// use to signal errors back to the WASM code.
func ErrorResponse(mem *jspatch.WasmMem, sp int32, id lib.Id, errorDetail string) {
	wasmPtr := mem.GetInt64(sp + 8)
	errId := id
	// the [0] value is the high 8 bytes, the [1] the low 8 bytes
	mem.SetInt64(int32(wasmPtr)+int32(unsafe.Offsetof(SinglePayload{}.ErrPtr)),
		int64(errId.High()))
	// high is 8 bytes higher
	mem.SetInt64(int32(wasmPtr)+int32(unsafe.Offsetof(SinglePayload{}.ErrPtr)+8),
		int64(errId.Low()))
	highBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(highBytes, errId.High())
}

// StackPointerToRequest assumes that this function was passed a pointer to the WASM stack and that
// this code is running on the GO side. This function assumes that 8 bytes above the wasm SP is a pointer
// to a SinglePayload.  Again, the GO side runs this to pull data from the WASM side.  This
// fills in the proto with the data sent from the WASM side or it returns an error.  An error here
// means we could not decode the package sent from the WASM side.  If an error is returned, the
// caller can simply return, as the payload has been modified to tell the other side about the error.
func StackPointerToRequest(mem *jspatch.WasmMem, sp int32, req proto.Message) (lib.Id, string) {
	wasmPtr := mem.GetInt64(sp + 8)

	buffer := ReadSlice(mem, wasmPtr,
		unsafe.Offsetof(SinglePayload{}.InPtr),
		unsafe.Offsetof(SinglePayload{}.InLen))

	id, detail := DecodeSingleProto(buffer, req)
	if id != nil {
		ErrorResponse(mem, int32(wasmPtr), id, detail)
		return id, detail
	}
	return nil, ""
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

// newPerrorFromPayload decodes the single payload object provided and returns
// a perror based on it's content. This cannot live in the go/lib pcakage like
// the other Id/Error related helpers becase it creates an import cycle.  This function
// asumes you have already checked the input param and it is known to be an error.
func NewPerrorFromSinglePayload(sp *SinglePayload) lib.Error {
	buf := make([]byte, sp.ErrDetailLen)
	for i := 0; i < len(buf); i++ {
		str := (*byte)(unsafe.Pointer((uintptr(int(sp.ErrDetail) + i))))
		buf[i] = *str
	}
	h := uint64(sp.ErrPtr[0])
	l := uint64(sp.ErrPtr[1])
	id := lib.NewIdCopy(h, l) // DANGER!
	return lib.NewPerrorFromId(string(buf), id)
}
