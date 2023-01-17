package main

import (
	"errors"
	"os"

	"github.com/iansmith/parigot/api_impl/log/go_"
	"github.com/iansmith/parigot/api_impl/splitutil"
	"github.com/iansmith/parigot/api_impl/syscall"
	"github.com/iansmith/parigot/g/log/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	pbsys "github.com/iansmith/parigot/g/msg/syscall/v1"
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
	if _, err := callImpl.Run(&pbsys.RunRequest{Wait: true}); err != nil {
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
	errId, err := splitutil.SendReceiveSingleProto(callImpl, inProto, &resp, go_.LogRequestHandler)
	if err != nil {
		return err
	}
	if errId != nil {
		return errors.New("Log() failed:" + errId.Short())
	}

	req, ok := inProto.(*logmsg.LogRequest)
	if !ok {
		return nil
	}
	if req.Level == logmsg.LogLevel_LOG_LEVEL_FATAL {
		os.Exit(7)
	}
	return nil
}
