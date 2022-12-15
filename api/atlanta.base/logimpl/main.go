package main

import (
	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/api/logimpl/ui"
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
// This file contains the true implementations--the server side--for the method(s)
// defined in log.proto.
//

func (m *myLogServer) Log(pctx *protosupport.Pctx, inProto proto.Message) error {
	req := inProto.(*pb.LogRequest)
	buffer, err := proto.Marshal(req)
	if err != nil {
		return err
	}
	sh := &reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&buffer[0])),
		Len:  len(buffer),
		Cap:  len(buffer),
	}
	// we have to convert this because pointers are only 32 bits in most wasm compilers
	u := uintptr(unsafe.Pointer(sh))

	ui.LogRequestViaSocket(int32(u))
	return nil
}
