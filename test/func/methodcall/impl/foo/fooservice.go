package main

import (
	"context"
	"math"
	"unsafe"

	"github.com/iansmith/parigot/g/methodcall/v1"

	pcontext "github.com/iansmith/parigot/context"
	methodcallmsg "github.com/iansmith/parigot/g/msg/methodcall/v1"
	"github.com/iansmith/parigot/test/func/methodcall/impl/foo/const_"
)

var _ = unsafe.Sizeof([]byte{})

const pathPrefix = "/parigotvirt/"

func main() {
	ctx := pcontext.NewContextWithContainer(context.Background(), "[fooservice]main")
	ctx = pcontext.CallTo(pcontext.GuestContext(ctx), "[foo]main")
	defer pcontext.Dump(ctx)
	pcontext.Debugf(ctx, "started main open")
	myId := methodcall.MustRegisterFooService(ctx)
	methodcall.MustExportFooService(ctx)
	pcontext.Debugf(ctx, "finished exported from foo, who is %s", myId.Short())
	methodcall.RunFooService(ctx, &fooServer{})
}

// this type better implement methodcall.v1.FooService
type fooServer struct {
}

//
// This file contains the true implementations--the server side--for the methods
// defined in foo.proto.
//

func (f *fooServer) AddMultiply(ctx context.Context, req *methodcallmsg.AddMultiplyRequest) (*methodcallmsg.AddMultiplyResponse, methodcall.MethodcallErrId) {
	//f.log(pctx, pblog.LogLevel_LOG_LEVEL_DEBUG, "received call for fooServer.AddMultiply")
	resp := &methodcallmsg.AddMultiplyResponse{}
	if req.IsAdd {
		resp.Result = req.Value0 + req.Value1
	} else {
		resp.Result = req.Value0 * req.Value1
	}
	return resp, methodcall.MethodcallErrIdNoErr
}

func (f *fooServer) LucasSequence(ctx context.Context) (*methodcallmsg.LucasSequenceResponse, methodcall.MethodcallErrId) {
	pcontext.Debugf(ctx, "LucasSequence", "received call for fooServer.LucasSequence")
	resp := &methodcallmsg.LucasSequenceResponse{}
	seq := make([]int32, const_.LucasSize) // -2 because first two are given
	seq[0] = 2
	seq[1] = 1
	for i := 2; i < const_.LucasSize; i++ {
		seq[i] = seq[i-1] + seq[i-2]
	}
	resp.Sequence = seq
	return resp, methodcall.MethodcallErrIdNoErr
}

// Newton-Raphson method, terms values beyond about 4 are silly
func (f *fooServer) WritePi(ctx context.Context, req *methodcallmsg.WritePiRequest) methodcall.MethodcallErrId {
	pcontext.Debugf(ctx, "WritePi", "received call for fooServer.AddMultiply")

	if req.GetTerms() < 1 {
		return methodcall.NewMethodcallErrId(MethodcallErrIdBadTerms)
	}
	runningTotal := 3.0 // k==0 term

	for k := 1; k <= int(req.GetTerms()); k++ {
		runningTotal = runningTotal - math.Tan(runningTotal)
		pcontext.Debugf(ctx, "WritePi", "%f", runningTotal)
	}
	pcontext.Infof(ctx, "WritePi result:", "%f", runningTotal)
	return methodcall.ZeroValueMethodcallErrId()
}

// Ready is a check, if this returns false the library will abort and not attempt to run this service.
// Normally, this is used to block using the lib.Run() call.  This call will wait until all the required
// services are ready.
func (f *fooServer) Ready(ctx context.Context) bool {
	methodcall.WaitFooServiceOrPanic()

	return true
}
