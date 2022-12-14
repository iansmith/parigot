package main

import (
	"fmt"
	"time"

	"github.com/iansmith/parigot/api/proto/g/log"
	pb "github.com/iansmith/parigot/api/proto/g/pb/log"
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

func (m *myLogServer) Log(pctx lib.Pctx, inProto proto.Message) error {
	req := inProto.(*pb.LogRequest)
	print("xxxxlogloglog ", fmt.Sprintf("%s:%s", req.Stamp.AsTime().Format(time.RFC3339), req.GetMessage()))
	return nil
}
