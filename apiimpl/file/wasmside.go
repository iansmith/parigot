package main

import (
	"fmt"

	"github.com/iansmith/parigot/apiimpl/file/go_"
	"github.com/iansmith/parigot/apiimpl/splitutil"
	"github.com/iansmith/parigot/apiimpl/syscall"
	"github.com/iansmith/parigot/g/file/v1"
	"github.com/iansmith/parigot/g/log/v1"
	filemsg "github.com/iansmith/parigot/g/msg/file/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"

	"google.golang.org/protobuf/proto"
)

var callImpl = syscall.NewCallImpl()

func main() {
	print("FILE FILE FILE main1\n")
	// we export and require services before the call to file.Run()... our call to the Run() system call is in Ready()
	if _, err := callImpl.Export1("file", "FileService"); err != nil {
		print("FILE FILE FILE export file" + err.Error() + "\n")
		for i := 0; i < 15; i++ {
			print("\n")
		}
	}
	print("FILE FILE FILE main2\n")
	if _, err := callImpl.Require1("log", "LogService"); err != nil {
		print("FILE FILE FILE require1 file" + err.Error() + "\n")
		for i := 0; i < 15; i++ {
			print("\n")
		}
	}
	print("FILE FILE FILE main3\n")

	file.RunFileService(&myFileServer{})
}

type myFileServer struct {
	logger log.LogService
}

func (m *myFileServer) Ready() bool {
	print("FILE FILE FILE ready\n")
	if _, err := callImpl.Run(&syscallmsg.RunRequest{Wait: true}); err != nil {
		print("ready: error in attempt to signal Run: ", err.Error(), "\n")
		return false
	}
	var err error
	m.logger, err = log.LocateLogService()
	if err != nil {
		panic("unable to locate the log:" + err.Error())
	}

	print("FILE FILE FILE erady done\n")

	return true
}

// This file contains the "setup" code that builds a payload that will be sent to the other part of
// this service.  That other part is the one that runs natively on the host machine.  This code runs
// in WASM.  Note that this code returns a "normal" go error; the code generator is
// go specific and the code generator determines the signature here.
func (m *myFileServer) Open(pctx *protosupportmsg.Pctx, inProto proto.Message) (proto.Message, error) {
	resp := filemsg.OpenResponse{}
	// your IDE may become confuse and show an error because of the tricks we are doing to call LogRequestHandler
	_, errId, errDetail := splitutil.SendReceiveSingleProto(callImpl, inProto, &resp, go_.FileSvcOpen)
	if errId != nil {
		return nil, lib.NewPerrorFromId(errDetail, errId)
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

// log is a utility for sending messages to the logging service.
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

// LoadTest is a method that will read a directory _on the host machine_ and load
// it into an in-memory filesystem.  This is used only for testing.  Once the
// LoadTest has completed successfully the in memory filesystem can be used with other
// file-related calls.
func (m *myFileServer) LoadTest(pctx *protosupportmsg.Pctx, inProto proto.Message) (proto.Message, error) {
	resp := filemsg.LoadTestResponse{}
	// your IDE may become confuse and show an error because of the tricks we are doing to call LogRequestHandler
	_, errId, errDetail := splitutil.SendReceiveSingleProto(callImpl, inProto, &resp, go_.FileSvcLoad)
	in := inProto.(*filemsg.LoadTestRequest)
	if errId != nil {
		m.log(nil, "in WASM fileserver.Load('%s'), error trying to return: %s (%s)", in.Path, errId.Short(), errDetail)
		return nil, lib.NewPerrorFromId(errDetail, errId)
	}
	return &resp, nil
}
