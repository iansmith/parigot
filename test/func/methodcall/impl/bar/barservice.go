package main

import (
	"fmt"
	golog "log"
	"time"

	"github.com/iansmith/parigot/g/methodcall/v1"

	"github.com/iansmith/parigot/apiimpl/syscall"
	"github.com/iansmith/parigot/g/log/v1"
	pblog "github.com/iansmith/parigot/g/msg/log/v1"
	methodcallmsg "github.com/iansmith/parigot/g/msg/methodcall/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var callImpl = syscall.NewCallImpl()

func main() {
	//lib.FlagParseCreateEnv()
	print(" xxxx --- about to start doing require in Bar service\n")
	log.RequireLogServiceOrPanic()
	methodcall.RequireFooServiceOrPanic()

	// one cannot initialize the fields of barServer{} here, must wait until Ready() is called
	methodcall.ExportBarServiceOrPanic()
	print(" xxxx --- about to run bar service from main\n")
	methodcall.RunBarService(&barServer{})
}

// this type better implement methodcall.v1.BarService
type barServer struct {
	logger log.LogService
	foo    *methodcall.FooServiceClient
}

//
// This file contains the true implementations--the server side--for the methods
// defined in bar.proto.
//

func (b *barServer) Accumulate(pctx *protosupportmsg.Pctx, in protoreflect.ProtoMessage) (protoreflect.ProtoMessage, error) {
	req := in.(*methodcallmsg.AccumulateRequest)
	//f.log(pctx, pblog.LogLevel_LOG_LEVEL_DEBUG, "received call for barServer.AccumulateMultiply")
	resp := &methodcallmsg.AccumulateResponse{}
	if len(req.Value) == 0 {
		resp.Product = 0
		resp.Sum = 0
		return resp, nil
	}

	reqAdd := &methodcallmsg.AddMultiplyRequest{
		IsAdd: true,
	}
	reqMul := &methodcallmsg.AddMultiplyRequest{
		IsAdd: false,
	}

	reqAdd.Value1 = 0 //identity to start
	reqMul.Value1 = 1 // identity to start

	var respAdd, respMul *methodcallmsg.AddMultiplyResponse
	var err error
	for i := 0; i < len(req.GetValue()); i++ {
		reqAdd.Value0 = req.GetValue()[i]
		respAdd, err = b.foo.AddMultiply(reqAdd)
		if err != nil {
			return nil, err
		}
		// b.log(nil, pblog.LogLevel_LOG_LEVEL_DEBUG, "ADD (%d,%d) iteration #%d, result add %d",
		// 	reqAdd.GetValue0(), reqAdd.GetValue1(), i, respAdd.Result)
		reqAdd.Value1 = respAdd.GetResult()

		/// multiply
		reqMul.Value0 = req.GetValue()[i]
		respMul, err = b.foo.AddMultiply(reqMul)
		if err != nil {
			return nil, err
		}
		// b.log(nil, pblog.LogLevel_LOG_LEVEL_DEBUG, "MUL (%d,%d) iteration #%d, result mul %d",
		// 	reqMul.GetValue0(), reqMul.GetValue1(), i, respMul.Result)
		reqMul.Value1 = respMul.GetResult()
	}
	resp.Product = respMul.GetResult()
	resp.Sum = respAdd.GetResult()
	//b.log(nil, pblog.LogLevel_LOG_LEVEL_DEBUG, "final tally--- sum=%d prod=%d",
	//	resp.GetProduct(), resp.GetSum())
	return resp, nil
}

// Ready is a check, if this returns false the library will abort and not attempt to run this service.
// Normally, this is used to block using the lib.Run() call.  This call will wait until all the required
// services are ready.
func (b *barServer) Ready() bool {
	golog.Printf("xxx--- ready called on bar")
	methodcall.WaitBarServiceOrPanic()

	b.logger = log.LocateLogServiceOrPanic()
	b.foo = methodcall.LocateFooServiceOrPanic(b.logger)
	return true
}

func (b *barServer) log(pctx *protosupportmsg.Pctx, level pblog.LogLevel, spec string, rest ...interface{}) {
	n := time.Now()
	if pctx != nil && !pctx.GetNow().AsTime().IsZero() {
		n = pctx.GetNow().AsTime() // xxx need to use kernel time or better use the pctx itself
	}
	msg := fmt.Sprintf(spec, rest...)
	req := pblog.LogRequest{
		Stamp:   timestamppb.New(n),
		Level:   level,
		Message: msg,
	}
	b.logger.Log(&req)
}
