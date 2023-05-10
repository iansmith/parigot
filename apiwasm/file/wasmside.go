package main

import (
	"fmt"

	"github.com/iansmith/parigot/apiwasm/syscall"
	"github.com/iansmith/parigot/g/file/v1"
	"github.com/iansmith/parigot/g/log/v1"
	filemsg "github.com/iansmith/parigot/g/msg/file/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"

	"google.golang.org/protobuf/proto"
)

//
// These two functions are the parts of this service that
// have to be implemented on the host.
//
// see apigo/file/go_/filehost.go

// go:wasm-module file
// go:export open
func Open(*filemsg.OpenRequest) *filemsg.OpenResponse

// go:wasm-module file
// go:export load_test_data
func LoadTestData(*filemsg.LoadTestDataRequest) *filemsg.LoadTestDataResponse

//go:export parigot_main
//go:linkname parigot_main
func parigot_main() {
	file.ExportFileServiceOrPanic()
	log.RequireLogServiceOrPanic()

	file.RunFileService(&myFileServer{})
}

type myFileServer struct {
	logger log.LogService
}

// Do you setup of myFileServer in here, not in main
func (m *myFileServer) Ready() bool {
	_ = syscall.Run(&syscallmsg.RunRequest{Wait: true})
	var err error
	m.logger, err = log.LocateLogService()
	if err != nil {
		panic("unable to locate the log:" + err.Error())
	}

	return true
}

// This file contains the "setup" code that builds a payload that will be sent to the other part of
// this service.  That other part is the one that runs natively on the host machine.  This code runs
// in WASM.  Note that this code returns a "normal" go error; the code generator is
// go specific and the code generator determines the signature here.
func (m *myFileServer) Open(pctx *protosupportmsg.Pctx, in *filemsg.OpenRequest) (*filemsg.OpenResponse, error) {
	out := Open(in)
	id := lib.Unmarshal(out.Id)
	if id.IsErrorType() && id.IsError() {
		return nil, lib.NewPerrorFromId("open goside call", id)
	}
	return out, nil
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
func (m *myFileServer) LoadTestData(pctx *protosupportmsg.Pctx, inProto proto.Message) (proto.Message, error) {
	resp := filemsg.LoadTestResponse{}
	// your IDE may become confuse and show an error because of the tricks we are doing to call LogRequestHandler
	spayload := splitutil.SendReceiveSingleProto(callImpl, inProto, &resp, go_.FileSvcLoad)
	if splitutil.IsErrorInSinglePayload(spayload) {
		err := splitutil.NewPerrorFromSinglePayload(spayload)
		return nil, err
	}
	return &resp, nil
}
