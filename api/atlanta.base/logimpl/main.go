package main

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/api/logimpl/ui"
	"github.com/iansmith/parigot/api/netconst"
	"github.com/iansmith/parigot/api/proto/g/log"
	pb "github.com/iansmith/parigot/api/proto/g/pb/log"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	"github.com/iansmith/parigot/lib"

	"google.golang.org/protobuf/proto"
)

func main() {
	log.Run(&myLogServer{})
}

type myLogServer struct{}

func (m *myLogServer) Ready() bool {
	if _, err := lib.Export1("log", "Log"); err != nil {
		print("ready: error in attempt to export api.Log: ", err.Error(), "\n")
		return false
	}
	if _, err := lib.Run(false); err != nil {
		print("ready: error in attempt to signal Run: ", err.Error(), "\n")
		return false
	}
	return true

}

//
// This file contains the "setup" code that builds a payload that will be sent to the other part of
// this service.  That other part is the one that runs natively on the host machine.
// This code receives
//

func (m *myLogServer) Log(pctx *protosupport.Pctx, inProto proto.Message) error {
	var err error
	req := inProto.(*pb.LogRequest)
	size := proto.Size(req)
	buffer := make([]byte, netconst.FrontMatterSize)
	if len(buffer) >= netconst.ReadBufferSize {
		panic("log message too large to fit in receive buffer:" + fmt.Sprint(len(buffer)))
	}
	binary.LittleEndian.PutUint64(buffer[:8], netconst.MagicStringOfBytes)
	binary.LittleEndian.PutUint32(buffer[8:netconst.FrontMatterSize], uint32(size))
	buffer, err = proto.MarshalOptions{}.MarshalAppend(buffer, req)
	if err != nil {
		return err
	}

	result := crc32.Checksum(buffer[:netconst.FrontMatterSize+size], netconst.KoopmanTable)
	buffer = append(buffer, []byte{0, 0, 0, 0}...) //space for the crc
	binary.LittleEndian.PutUint32(buffer[netconst.FrontMatterSize+size:], uint32(result))
	payload := &ui.LogViewerPayload{}
	payload.Ptr, payload.Len = sliceToTwoInt64s(buffer)
	// we have to convert this because pointers are only 32 bits in most wasm compilers
	u := uintptr(unsafe.Pointer(payload))

	ui.LogRequestViaSocket(int32(u))
	return nil
}

func sliceToTwoInt64s(b []byte) (int64, int64) {
	slh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return int64(slh.Data), int64(slh.Len)
}
