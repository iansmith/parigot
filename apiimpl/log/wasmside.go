package main

import (
	"fmt"

	"github.com/iansmith/parigot/apiimpl/log/go_"
	"github.com/iansmith/parigot/apiimpl/splitutil"
	"github.com/iansmith/parigot/apiimpl/syscall"
	"github.com/iansmith/parigot/g/log/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"

	"google.golang.org/protobuf/proto"
)

var callImpl = syscall.NewCallImpl()

func main() {
	lib.FlagParseCreateEnv()
	// you need to put Require and Export calls in here, but put Run() call in Ready()
	if _, err := callImpl.Export1("log", "LogService"); err != nil {
		panic("myLogServer:ready: error in attempt to export api.Log: " + err.Error())
	}
	log.RunLogService(&myLogServer{})
}

type myLogServer struct{}

func (m *myLogServer) Ready() bool {
	if _, err := callImpl.Run(&syscallmsg.RunRequest{Wait: true}); err != nil {
		panic("myLogServer: ready: error in attempt to signal Run: " + err.Error())
	}
	return true
}

//
// This file contains the "setup" code that builds a payload that will be sent to the other part of
// this service.  That other part is the one that runs natively on the host machine.
// We discard the pb.LogResponse{} since there is no content inside it.

func (m *myLogServer) Log(pctx *protosupportmsg.Pctx, inProto proto.Message) error {

	resp := logmsg.LogResponse{}
	// your IDE may become confuse and show an error because of the tricks we are doing to call LogRequestHandler
	spayload := splitutil.SendReceiveSingleProto(callImpl, inProto, &resp, go_.LogRequestHandler)
	if splitutil.IsErrorInSinglePayload(spayload) {
		err := splitutil.NewPerrorFromSinglePayload(spayload)
		return err
	}
	req, ok := inProto.(*logmsg.LogRequest)
	if !ok {
		print("MAYDAY! trying ot log an error, but failed!\n")
		return nil // xxx ugh swallowing this? But what else can be done?
	}
	if req.Level == logmsg.LogLevel_LOG_LEVEL_FATAL {
		print(fmt.Sprintf("xxx -- FATAL FATAL FATAL --> %s\n", inProto.(*logmsg.LogRequest).Message))
		//callImpl.Exit(&syscallmsg.ExitRequest{Code: 7})
	}
	return nil
}
