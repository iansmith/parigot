package apiwasm

import (
	"fmt"
	"unsafe"

	"github.com/iansmith/parigot/id"
	"github.com/iansmith/parigot/sharedconst"

	"google.golang.org/protobuf/proto"
)

// ReturnData is the data returned by call to the "host side"
// of the system.  Unless you are changing host functions in
// serious ways, you probably don't need to be concerned with this.
type ReturnData struct {
	length, ptr int32
	idErr       [2]int64
}

// ToError takes a ReturnData object and returns either nil
// for no error, or an error object (actually a PError) that is
// derived from the error code.
func (r *ReturnData) ToError(msg string) error {
	high := uint64(r.idErr[0])
	low := uint64(r.idErr[0])
	idOrError := id.NewIdCopy(high, low)
	if idOrError.ErrorCode() == 0 {
		return nil
	}
	return id.NewPerrorFromId(msg, idOrError)
}

func Decode[U proto.Message]() {
	// do unmarshal here
}

func NewReturnData(ptr int32) *ReturnData {
	self := (*ReturnData)(unsafe.Pointer(uintptr(ptr)))
	return self
}

func toByteBuffer(msg proto.Message) (int32, *byte, error) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return 0, nil, fmt.Errorf("unable to marshal bytes for WASM: %v", err)
	}
	return int32(len(buf)), unsafe.SliceData(buf), nil
}

func fromByteBuffer(length, ptr int32, msg proto.Message) error {
	b := (*byte)(unsafe.Pointer(uintptr(ptr)))
	buffer := unsafe.Slice(b, length)
	err := proto.Unmarshal(buffer, msg)
	if err != nil {
		return fmt.Errorf("unable to unmarshal bytes from WASM: %v", err)

	}
	return nil
}

func WasmCallNativeInOut[U proto.Message, V proto.Message](in proto.Message, out proto.Message, fn func(int32, int32) int32) error {
	return wasmCallNative[U, V](in, out, fn)
}
func WasmCallNativeIn[U proto.Message, V proto.Message](in proto.Message, fn func(int32, int32) int32) error {
	return wasmCallNative[U, V](in, nil, fn)
}
func WasmCallNativeOut[U proto.Message, V proto.Message](out proto.Message, fn func(int32, int32) int32) error {
	return wasmCallNative[U, V](nil, out, fn)
}

func wasmCallNative[U proto.Message, V proto.Message](in proto.Message, out proto.Message, fn func(int32, int32) int32) error {
	var outRaw int32
	if in != nil {
		length, ptr, err := toByteBuffer(in)
		if err != nil {
			return err
		}
		if length == 0 {
			return fmt.Errorf("unable convert proto.Message to byte buffer, resulting buffer is zero sized")
		}
		ptr32 := int32(uintptr(unsafe.Pointer(ptr)))
		outRaw = fn(length, ptr32)
	} else {
		outRaw = fn(0, 0)
	}
	rd := NewReturnData(outRaw)
	if out == nil {
		return nil
	}
	err := fromByteBuffer(rd.length, rd.ptr, out)
	if err != nil {
		return err
	}
	return rd.ToError("wasmCallNative failed:")
}

// This global variable is here *JUST* to thwart the GC.
var globalPtr = make(map[uintptr][]byte)

//go:export apiwasm new_return_data_with_buffer
func NewReturnDataWithBuffer(size int32) int32 {
	s := unsafe.Sizeof(ReturnData{})
	off := unsafe.Offsetof(ReturnData{}.idErr)
	if s != sharedconst.ReturnDataSize || off != sharedconst.ReturnDataSize {
		panic(fmt.Sprintf("running on unknown wasm host, size of return data should be %d but is %d, offset of idErr should be %d but is %d ",
			sharedconst.ReturnDataSize, s, sharedconst.ReturnDataSize, off))
	}
	buffer := make([]byte, size)
	data := unsafe.SliceData(buffer)
	rawRt := make([]byte, s)
	rt := (*ReturnData)(unsafe.Pointer(unsafe.SliceData(rawRt)))
	rt.length = size
	rt.ptr = int32(uintptr(unsafe.Pointer(data)))
	rt.idErr = [2]int64{}
	rawRtUintptr := uintptr((unsafe.Pointer(&rawRt)))
	bufferUintptr := uintptr((unsafe.Pointer(&buffer)))
	globalPtr[rawRtUintptr] = rawRt
	globalPtr[bufferUintptr] = buffer

	return int32(uintptr(unsafe.Pointer(rt)))
}

//go:export apiwasm new_string
func NewString(size int32) int32 {
	x := make([]byte, size)
	ptr := uintptr(unsafe.Pointer(unsafe.SliceData(x)))
	key := uintptr(unsafe.Pointer(&x))
	globalPtr[key] = x
	return int32(uintptr(ptr))
}
