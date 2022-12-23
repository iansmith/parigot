package main

import (
	"errors"
	"fmt"

	"github.com/iansmith/parigot/api/fileimpl/go_"
	"github.com/iansmith/parigot/api/proto/g/file"
	"github.com/iansmith/parigot/api/proto/g/log"
	pb "github.com/iansmith/parigot/api/proto/g/pb/file"
	pblog "github.com/iansmith/parigot/api/proto/g/pb/log"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	"github.com/iansmith/parigot/api/splitutil"
	"github.com/iansmith/parigot/lib"

	"google.golang.org/protobuf/proto"
)

func main() {
	// we export and require services before the call to file.Run()... our call to the Run() system call is in Ready()
	if _, err := lib.Export1("file", "File"); err != nil {
		panic("ready: error in attempt to export api.Log: " + err.Error())
	}
	if _, err := lib.Require1("log", "Log"); err != nil {
		panic("ready: error in attempt to export api.Log: " + err.Error())
	}

	file.Run(&myFileServer{})
}

type myFileServer struct {
	logger log.Log
}

func (m *myFileServer) Ready() bool {
	if _, err := lib.Run(false); err != nil {
		print("ready: error in attempt to signal Run: ", err.Error(), "\n")
		return false
	}
	var err error
	m.logger, err = log.LocateLog()
	if err != nil {
		panic("unable to locate the log:" + err.Error())
	}
	return true

}

// This file contains the "setup" code that builds a payload that will be sent to the other part of
// this service.  That other part is the one that runs natively on the host machine.  This code runs
// in WASM.
func (m *myFileServer) Open(pctx *protosupport.Pctx, inProto proto.Message) (proto.Message, error) {

	resp := pb.OpenResponse{}
	// your IDE may become confuse and show an error because of the tricks we are doing to call LogRequestHandler
	errId, err := splitutil.SendReceiveSingleProto(inProto, &resp, go_.FileSvcOpen)
	if err != nil {
		return nil, err
	}
	if errId != nil {
		return nil, errors.New("internal error:" + errId.Short())
	}
	return &resp, nil
}

func (m *myFileServer) Close(pctx *protosupport.Pctx, inProto proto.Message) (proto.Message, error) {
	_ = inProto.(*pb.CloseRequest)
	panic("Close")
}
func (m *myFileServer) Create(pctx *protosupport.Pctx, inProto proto.Message) (proto.Message, error) {
	_ = inProto.(*pb.CreateRequest)
	panic("Create")
}
func (m *myFileServer) log(pctx *protosupport.Pctx, spec string, rest ...interface{}) {
	s := fmt.Sprintf(spec, rest...)
	req := pblog.LogRequest{
		Stamp:   pctx.GetNow(),
		Level:   pblog.LogLevel_LOG_LEVEL_DEBUG,
		Message: "myFileServer:" + s,
	}
	err := m.logger.Log(&req)
	if err != nil {
		print("unable to log ", s, ":", err.Error(), "\n")
	}
}

func (m *myFileServer) Load(pctx *protosupport.Pctx, inProto proto.Message) (proto.Message, error) {
	resp := pb.LoadResponse{}
	// your IDE may become confuse and show an error because of the tricks we are doing to call LogRequestHandler
	errId, err := splitutil.SendReceiveSingleProto(inProto, &resp, go_.FileSvcLoad)
	if err != nil {
		return nil, err
	}
	if errId != nil {
		return nil, errors.New("internal error:" + errId.Short())
	}
	return &resp, nil
}
