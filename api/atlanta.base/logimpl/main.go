package main

import (
	"github.com/iansmith/parigot/api/logimpl/go_"
	"github.com/iansmith/parigot/api/proto/g/log"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	"github.com/iansmith/parigot/api/splitutil"
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
//

func (m *myLogServer) Log(pctx *protosupport.Pctx, inProto proto.Message) error {
	var err error
	u, err := splitutil.SendSingleProto(inProto)
	if err != nil {
		return err
	}
	go_.LogRequestViaSocket(int32(u))
	return nil
}
