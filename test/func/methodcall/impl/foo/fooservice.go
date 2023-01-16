package main

import (
	"fmt"
	"math"
	"time"

	"github.com/iansmith/parigot/g/methodcall/v1"
	lib "github.com/iansmith/parigot/lib/go"

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
	lib.FlagParseCreateEnv()
	// for i := 0; i < flag.NArg(); i++ {
	// 	print(fmt.Sprintf("xxx foo %d=>%s\n", i, flag.Arg(i)))
	// }

	//if things need to be required/exported you need to force them to the ready state BEFORE calling run()
	if _, err := callImpl.Require1("log", "LogService"); err != nil {
		panic("unable to require log service: " + err.Error())
	}
	if _, err := callImpl.Export1("methodcall", "FooService"); err != nil {
		panic("unable to export methodcall.Foo: " + err.Error())
	}
	// one cannot initialize the fields of fooServer{} here, must wait until Ready() is called
	methodcall.RunFooService(&fooServer{})
}

// this type better implement methodcall.v1.FooService
type fooServer struct {
	logger log.LogService
}

//
// This file contains the true implementations--the server side--for the methods
// defined in foo.proto.
//

func (f *fooServer) AddMultiply(pctx *protosupportmsg.Pctx, in protoreflect.ProtoMessage) (protoreflect.ProtoMessage, error) {
	req := in.(*methodcallmsg.AddMultiplyRequest)
	//f.log(pctx, pblog.LogLevel_LOG_LEVEL_DEBUG, "received call for fooServer.AddMultiply")
	resp := &methodcallmsg.AddMultiplyResponse{}
	if req.IsAdd {
		resp.Result = req.Value0 + req.Value1
	} else {
		resp.Result = req.Value0 * req.Value1
	}
	return resp, nil
}

// we don't provide a size to LucasSequence (via LucasSequenceRequest) becaues we are testing response
// with no input request
const lucasSize = 50

func (f *fooServer) LucasSequence(pctx *protosupportmsg.Pctx) (protoreflect.ProtoMessage, error) {
	f.log(pctx, pblog.LogLevel_LOG_LEVEL_DEBUG, "received call for fooServer.LucasSequence")
	resp := &methodcallmsg.LucasSequenceResponse{}
	seq := make([]int32, lucasSize)
	seq[0] = 2
	seq[1] = 1
	for i := 2; i < lucasSize-2; i++ {
		seq[i] = seq[i-1] + seq[i-2]
	}
	resp.Sequence = seq
	return resp, nil
}

// Newton-Raphson method, terms values beyond about 4 are silly
func (f *fooServer) WritePi(pctx *protosupportmsg.Pctx, in protoreflect.ProtoMessage) error {
	req := in.(*methodcallmsg.WritePiRequest)
	f.log(pctx, pblog.LogLevel_LOG_LEVEL_DEBUG, "received call for fooServer.AddMultiply")

	if req.GetTerms() < 1 {
		return fmt.Errorf("number of terms in WritePi must be a positive integer")
	}
	runningTotal := 3.0

	for k := 1; k <= int(req.GetTerms()); k++ {
		runningTotal = runningTotal - math.Tan(runningTotal)
	}
	f.log(pctx, pblog.LogLevel_LOG_LEVEL_INFO, "%f", runningTotal)
	return nil
}

// Ready is a check, if this returns false the library will abort and not attempt to run this service.
// Normally, this is used to block using the lib.Run() call.  This call will wait until all the required
// services are ready.
func (f *fooServer) Ready() bool {
	print("zzz in foo server about to  run()\n")
	if _, err := callImpl.Run(&syscallmsg.RunRequest{Wait: true}); err != nil {
		print("ready: error in attempt to signal Run: ", err.Error(), "\n")
		return false
	}
	logger, err := log.LocateLogService()
	if err != nil {
		print("ERROR trying to create log client: ", err.Error(), "\n")
		return false
	}
	f.logger = logger
	f.log(nil, pblog.LogLevel_LOG_LEVEL_DEBUG, "foo service: about to return from ready")
	return true
}

func (f *fooServer) log(pctx *protosupportmsg.Pctx, level pblog.LogLevel, spec string, rest ...interface{}) {
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
	f.logger.Log(&req)
}
