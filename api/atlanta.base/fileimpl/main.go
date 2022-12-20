package main

import (
	"github.com/iansmith/parigot/api/fileimpl/go_"
	"github.com/iansmith/parigot/api/proto/g/file"
	pb "github.com/iansmith/parigot/api/proto/g/pb/file"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	"github.com/iansmith/parigot/api/splitutil"
	"github.com/iansmith/parigot/lib"

	"google.golang.org/protobuf/proto"
)

func main() {
	file.Run(&myFileServer{})
}

type myFileServer struct{}

func (m *myFileServer) Ready() bool {
	if _, err := lib.Export1("file", "File"); err != nil {
		print("ready: error in attempt to export api.Log: ", err.Error(), "\n")
		return false
	}
	if _, err := lib.Run(false); err != nil {
		print("ready: error in attempt to signal Run: ", err.Error(), "\n")
		return false
	}
	return true

}

// This file contains the "setup" code that builds a payload that will be sent to the other part of
// this service.  That other part is the one that runs natively on the host machine.
// This code receives an encoded single protobuf object from the wasm-side client code.
func (m *myFileServer) Open(pctx *protosupport.Pctx, inProto proto.Message) (proto.Message, error) {
	u, err := splitutil.SendSingleProto(inProto)
	if err != nil {
		return nil, err
	}
	go_.FileSvcOpen(int32(u))
	return nil, nil
}

func (m *myFileServer) Close(pctx *protosupport.Pctx, inProto proto.Message) (proto.Message, error) {
	_ = inProto.(*pb.CloseRequest)
	panic("Close")
}
func (m *myFileServer) Create(pctx *protosupport.Pctx, inProto proto.Message) (proto.Message, error) {
	_ = inProto.(*pb.CreateRequest)
	panic("Create")
}
