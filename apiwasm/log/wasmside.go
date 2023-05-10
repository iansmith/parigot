package main

import (
	"fmt"
	golog "log"

	"github.com/iansmith/parigot/apiimpl/log/go_"
	"github.com/iansmith/parigot/apiimpl/splitutil"
	"github.com/iansmith/parigot/apiimpl/syscall"
	"github.com/iansmith/parigot/g/log/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	lib "github.com/iansmith/parigot/lib/go"

	"google.golang.org/protobuf/proto"
)

var callImpl = syscall.NewCallImpl()

func bar() {
	return
}

//go:export parigot_main
//go:linkname parigot_main
func parigot_main() {
	bar()
	for i := 0; i < 10; i++ {
		golog.Printf("xxx --> LOG mainxxx()")
		golog.Printf("xxx 222 --> LOG mainxxx()")
	}

	lib.FlagParseCreateEnv()
	golog.Printf("xxx --> LOG main() 2")

	log.ExportLogServiceOrPanic()
	golog.Printf("xxx --> LOG main() 3")
	log.RunLogService(&myLogServer{})
	golog.Printf("xxx --> LOG main() 4")

}

type myLogServer struct{}

func (m *myLogServer) Ready() bool {
	log.WaitLogServiceOrPanic()
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
