package main

import (
	"context"
	"fmt"
	"math"

	"github.com/iansmith/parigot/g/methodcall/v1"
	lib "github.com/iansmith/parigot/lib/go"

	pcontext "github.com/iansmith/parigot/context"
	methodcallmsg "github.com/iansmith/parigot/g/msg/methodcall/v1"
	"github.com/iansmith/parigot/test/func/methodcall/impl/foo/const_"
)

//go:export parigot_main
//go:linkname parigot_main
func parigot_main() {
	lib.FlagParseCreateEnv()

	bg := context.Background()
	methodcall.ExportFooServiceOrPanic()
	methodcall.RequireFooServiceOrPanic(bg)
	s := &fooServer{}
	methodcall.RunFooService(s)
}

// this type better implement methodcall.v1.FooService
type fooServer struct {
}

//
// This file contains the true implementations--the server side--for the methods
// defined in foo.proto.
//

func (f *fooServer) AddMultiply(ctx context.Context, req *methodcallmsg.AddMultiplyRequest) (*methodcallmsg.AddMultiplyResponse, error) {
	//f.log(pctx, pblog.LogLevel_LOG_LEVEL_DEBUG, "received call for fooServer.AddMultiply")
	resp := &methodcallmsg.AddMultiplyResponse{}
	if req.IsAdd {
		resp.Result = req.Value0 + req.Value1
	} else {
		resp.Result = req.Value0 * req.Value1
	}
	return resp, nil
}

func (f *fooServer) LucasSequence(ctx context.Context) (*methodcallmsg.LucasSequenceResponse, error) {
	pcontext.Debugf(ctx, "LucasSequence", "received call for fooServer.LucasSequence")
	resp := &methodcallmsg.LucasSequenceResponse{}
	seq := make([]int32, const_.LucasSize) // -2 because first two are given
	seq[0] = 2
	seq[1] = 1
	for i := 2; i < const_.LucasSize; i++ {
		seq[i] = seq[i-1] + seq[i-2]
	}
	resp.Sequence = seq
	return resp, nil
}

// Newton-Raphson method, terms values beyond about 4 are silly
func (f *fooServer) WritePi(ctx context.Context, req *methodcallmsg.WritePiRequest) error {
	pcontext.Debugf(ctx, "WritePi", "received call for fooServer.AddMultiply")

	if req.GetTerms() < 1 {
		return fmt.Errorf("number of terms in WritePi must be a positive integer")
	}
	runningTotal := 3.0 // k==0 term

	for k := 1; k <= int(req.GetTerms()); k++ {
		runningTotal = runningTotal - math.Tan(runningTotal)
	}
	pcontext.Debugf(ctx, "WritePi", "%f", runningTotal)
	return nil
}

// Ready is a check, if this returns false the library will abort and not attempt to run this service.
// Normally, this is used to block using the lib.Run() call.  This call will wait until all the required
// services are ready.
func (f *fooServer) Ready(ctx context.Context) bool {
	methodcall.WaitFooServiceOrPanic()

	return true
}
