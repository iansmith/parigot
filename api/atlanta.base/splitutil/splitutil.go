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
	pbfile "github.com/iansmith/parigot/api/proto/g/pb/file"
	pblog "github.com/iansmith/parigot/api/proto/g/pb/log"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	"github.com/iansmith/parigot/sys/jspatch"

	"google.golang.org/protobuf/proto"
)

type SplitUtilSinglePayload struct {
	Ptr int64
	Len int64
}

var emptyReturnProto = protosupport.CallId{}.Id
var DecodeError = errors.New("decoding error")

type SplitReqType interface {
	pblog.LogRequest
	pbfile.FileOpenRequest
}

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
	payload := &SplitUtilSinglePayload{}
	payload.Ptr, payload.Len = sliceToTwoInt64s(buffer)
	// we have to convert this because pointers are only 32 bits in most wasm compilers
	u := uintptr(unsafe.Pointer(payload))
	return u, nil
}

// DecodeProto decodes a buffer obtained when the client side drops the payload (above) to us
// inside the go world.  The return value of T should NOT be touched/used if error!=nil.
func DecodeProto[T SplitReqType](buffer []byte) (*T, error) {

	m := binary.LittleEndian.Uint64(buffer[0:8])
	if m != netconst.MagicStringOfBytes {
		log.Printf("unable to print log message, bad magic number %x", m)
		return nil, DecodeError
	}
	l := binary.LittleEndian.Uint32(buffer[8:12])
	if l >= uint32(netconst.ReadBufferSize) {
		log.Printf("unable to print log message, very large log message [%d bytes]", l)
		return nil, DecodeError
	}
	size := int(l)
	var req T
	objBuffer := buffer[netconst.FrontMatterSize : netconst.FrontMatterSize+size]
	if err := proto.Unmarshal(objBuffer, req); err != nil {
		log.Printf("unable to print log message, request could not be unmarshaled: %v", err)
		return nil, DecodeError
	}
	result := crc32.Checksum(objBuffer, netconst.KoopmanTable)
	expected := binary.LittleEndian.Uint32(buffer[netconst.FrontMatterSize+size : netconst.FrontMatterSize+size+4])
	if expected != result {
		log.Printf("unable to print log message, bad checksum found on log request")
		return nil, DecodeError
	}
	return &req, nil

}

func ReadSlice(mem *jspatch.WasmMem, structPtr int64, dataOffset uintptr, lenOffset uintptr) []byte {
	return mem.LoadSliceWithLenAddr(int32(structPtr)+int32(dataOffset),
		int32(structPtr)+int32(lenOffset))
}

func sliceToTwoInt64s(b []byte) (int64, int64) {
	slh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return int64(slh.Data), int64(slh.Len)
}
