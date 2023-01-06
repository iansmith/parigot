package main

import (
	"fmt"
	"time"

	"github.com/iansmith/parigot/g/methodcall/v1"

	"github.com/iansmith/parigot/api_impl/syscall"
	"github.com/iansmith/parigot/g/log/v1"
	pblog "github.com/iansmith/parigot/g/msg/log/v1"
	methodcallmsg "github.com/iansmith/parigot/g/msg/methodcall/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var callImpl = syscall.NewCallImpl()

func main() {
	//if things need to be required/exported you need to force them to the ready state BEFORE calling run()
	if _, err := callImpl.Require1("log", "LogService"); err != nil {
		panic("unable to require log service: " + err.Error())
	}
	if _, err := callImpl.Require1("methodcall", "FooService"); err != nil {
		panic("unable to require methodcall.FooService service: " + err.Error())
	}
	if _, err := callImpl.Export1("methodcall", "BarService"); err != nil {
		panic("unable to export methodcall.BarService: " + err.Error())
	}
	// one cannot initialize the fields of fooServer{} here, must wait until Ready() is called
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

func (f *barServer) Accumulate(pctx *protosupportmsg.Pctx, in protoreflect.ProtoMessage) (protoreflect.ProtoMessage, error) {
	req := in.(*methodcallmsg.AccumulateRequest)
	f.log(pctx, pblog.LogLevel_LOG_LEVEL_DEBUG, "received call for barServer.AccumulateMultiply")
	resp := &methodcallmsg.AccumulateResponse{}

	reqAdd := &methodcallmsg.AddMultiplyRequest{
		IsAdd: true,
	}
	reqMul := &methodcallmsg.AddMultiplyRequest{
		IsAdd: false,
	}
	addTerm := make([]int32, 2)
	mulTerm := make([]int32, 2)
	addTerm[0] = req.GetValue()[0]
	addTerm[1] = req.GetValue()[1]
	mulTerm[0] = req.GetValue()[0]
	mulTerm[1] = req.GetValue()[1]

	accSum := addTerm[0] + addTerm[1]
	accProduct := mulTerm[0] * mulTerm[1]

	for i := 1; i < len(req.GetValue()); i++ {
		reqAdd.Value0 = addTerm[0]
		reqAdd.Value1 = addTerm[1]
		resp, err := f.foo.AddMultiply(reqAdd)
		if err != nil {
			return nil, err
		}
		accSum += resp.Result
		addTerm[0] = addTerm[1]
		addTerm[1] = resp.Result

		reqMul.Value0 = mulTerm[0]
		reqMul.Value1 = mulTerm[1]
		resp, err = f.foo.AddMultiply(reqMul)
		if err != nil {
			return nil, err
		}
		accProduct += resp.Result
		mulTerm[0] = mulTerm[1]
		mulTerm[1] = resp.Result
	}
	resp.Product = accProduct
	resp.Sum = accSum
	return resp, nil

}

// Ready is a check, if this returns false the library will abort and not attempt to run this service.
// Normally, this is used to block using the lib.Run() call.  This call will wait until all the required
// services are ready.
func (b *barServer) Ready() bool {
	if _, err := callImpl.Run(&syscallmsg.RunRequest{Wait: true}); err != nil {
		print("ready: error in attempt to signal Run: ", err.Error(), "\n")
		return false
	}
	logger, err := log.LocateLogService()
	if err != nil {
		print("ERROR trying to create log client: ", err.Error(), "\n")
		return false
	}
	foo, err := methodcall.LocateFooService(logger)
	if err != nil {
		print("ERROR trying to create foo client: ", err.Error(), "\n")
		return false
	}
	b.logger = logger
	b.foo = foo
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
