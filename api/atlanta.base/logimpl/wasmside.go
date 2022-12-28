package main

import (
	"github.com/iansmith/parigot/api/logimpl/go_"
	"github.com/iansmith/parigot/api/proto/g/log"
	pb "github.com/iansmith/parigot/api/proto/g/pb/log"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	"github.com/iansmith/parigot/api/splitutil"
	"github.com/iansmith/parigot/api/syscall"

	"google.golang.org/protobuf/proto"
)

func main() {
	// you need to put Require and Export calls in here, but put Run() call in Ready()
	if _, err := syscall.Export1("log", "Log"); err != nil {
		panic("myLogServer:ready: error in attempt to export api.Log: " + err.Error())
	}
	log.Run(&myLogServer{})
}

type myLogServer struct{}

func (m *myLogServer) Ready() bool {
	if _, err := syscall.Run(false); err != nil {
		panic("myLogServer: ready: error in attempt to signal Run: " + err.Error())
	}

	return true
}

//
// This file contains the "setup" code that builds a payload that will be sent to the other part of
// this service.  That other part is the one that runs natively on the host machine.
// We discard the pb.LogResponse{} since there is no content inside it.

func (m *myLogServer) Log(pctx *protosupport.Pctx, inProto proto.Message) error {
	resp := pb.LogResponse{}
	// your IDE may become confuse and show an error because of the tricks we are doing to call LogRequestHandler
	errId, err := splitutil.SendReceiveSingleProto(inProto, &resp, go_.LogRequestHandler)
	if err != nil {
		return err
	}
	if errId != nil {
		print("unable to log due to kernel error:", errId.Short())
	}
	return nil
}
