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
	// you need to put Require and Export calls in here, but put Run() call in Ready()
	if _, err := lib.Export1("log", "Log"); err != nil {
		panic("myLogServer:ready: error in attempt to export api.Log: " + err.Error())
	}
	log.Run(&myLogServer{})
}

type myLogServer struct{}

func (m *myLogServer) Ready() bool {
	if _, err := lib.Run(false); err != nil {
		panic("myLogServer: ready: error in attempt to signal Run: " + err.Error())
	}

	return true
}

//
// This file contains the "setup" code that builds a payload that will be sent to the other part of
// this service.  That other part is the one that runs natively on the host machine.
// Note that since there is no return value, we don't bother initializing or reading results from the
// results defined in splitutil.SinglePayload

func (m *myLogServer) Log(pctx *protosupport.Pctx, inProto proto.Message) error {
	var err error
	u, err := splitutil.SendSingleProto(inProto)
	if err != nil {
		return err
	}
	go_.LogRequestViaSocket(int32(u))
	return nil
}
