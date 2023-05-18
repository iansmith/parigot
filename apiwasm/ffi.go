package apiwasm

import (
	"fmt"
	"unsafe"

	"github.com/iansmith/parigot/id"

	"google.golang.org/protobuf/proto"
)

// ReturnData is the data returned by call to the "host side"
// of the system.  Unless you are changing host functions in
// serious ways, you probably don't need to be concerned with this.
type ReturnData struct {
	length, ptr int32
	idErr       [2]int64
}

// ReturnDataSizeWasm is the size of the structure ReturnData
// in bytes on BOTH the host and the wasm side.
const ReturnDataSize = 24

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
