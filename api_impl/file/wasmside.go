package main

import (
	"errors"
	"fmt"

	"github.com/iansmith/parigot/api_impl/file/go_"
	"github.com/iansmith/parigot/api_impl/splitutil"
	"github.com/iansmith/parigot/api_impl/syscall"
	"github.com/iansmith/parigot/g/file/v1"
	"github.com/iansmith/parigot/g/log/v1"
	filemsg "github.com/iansmith/parigot/g/msg/file/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"

	"google.golang.org/protobuf/proto"
)

var callImpl = syscall.NewCallImpl()

func main() {
	// we export and require services before the call to file.Run()... our call to the Run() system call is in Ready()
	if _, err := callImpl.Export1("file", "FileService"); err != nil {
		panic("ready: error in attempt to export api.Log: " + err.Error())
	}
	if _, err := callImpl.Require1("log", "LogService"); err != nil {
		panic("ready: error in attempt to export api.Log: " + err.Error())
	}

	file.RunFileService(&myFileServer{})
}

type myFileServer struct {
	logger log.LogService
}

func (m *myFileServer) Ready() bool {
	if _, err := callImpl.Run(&syscallmsg.RunRequest{Wait: true}); err != nil {
		print("ready: error in attempt to signal Run: ", err.Error(), "\n")
		return false
	}
	var err error
	m.logger, err = log.LocateLogService()
	if err != nil {
		panic("unable to locate the log:" + err.Error())
	}
	return true
}

// This file contains the "setup" code that builds a payload that will be sent to the other part of
// this service.  That other part is the one that runs natively on the host machine.  This code runs
// in WASM.
func (m *myFileServer) Open(pctx *protosupportmsg.Pctx, inProto proto.Message) (proto.Message, error) {
	resp := filemsg.OpenResponse{}
	// your IDE may become confuse and show an error because of the tricks we are doing to call LogRequestHandler
	errId, err := splitutil.SendReceiveSingleProto(callImpl, inProto, &resp, go_.FileSvcOpen)
	if err != nil {
		return nil, err
	}
	if errId != nil {
		return nil, errors.New("internal error:" + errId.Short())
	}
	return &resp, nil
}

func (m *myFileServer) Close(pctx *protosupportmsg.Pctx, inProto proto.Message) (proto.Message, error) {
	_ = inProto.(*filemsg.CloseRequest)
	panic("Close")
}
func (m *myFileServer) Create(pctx *protosupportmsg.Pctx, inProto proto.Message) (proto.Message, error) {
	_ = inProto.(*filemsg.CreateRequest)
	panic("Create")
}
func (m *myFileServer) log(pctx *protosupportmsg.Pctx, spec string, rest ...interface{}) {
	s := fmt.Sprintf(spec, rest...)
	req := logmsg.LogRequest{
		Stamp:   pctx.GetNow(),
		Level:   logmsg.LogLevel_LOG_LEVEL_DEBUG,
		Message: "myFileServer:" + s,
	}
	err := m.logger.Log(&req)
	if err != nil {
		print("unable to log ", s, ":", err.Error(), "\n")
	}
}

func (m *myFileServer) LoadTest(pctx *protosupportmsg.Pctx, inProto proto.Message) (proto.Message, error) {
	resp := filemsg.LoadTestResponse{}
	// your IDE may become confuse and show an error because of the tricks we are doing to call LogRequestHandler
	errId, err := splitutil.SendReceiveSingleProto(callImpl, inProto, &resp, go_.FileSvcLoad)
	in := inProto.(*filemsg.LoadTestRequest)
	if err != nil {
		m.log(nil, "in WASM fileserver.Load('%s'), error trying to return: %v", in.Path, err)
		return nil, err
	}
	if errId != nil {
		m.log(nil, "in WASM fileserver.Load('%s') error found: %s", in.Path, errId.Short())
		return nil, errors.New("internal error:" + errId.Short())
	}
	return &resp, nil
}
