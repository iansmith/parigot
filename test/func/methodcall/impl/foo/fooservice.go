package main

import (
	"context"
	"log/slog"
	"math"
	"unsafe"

	"github.com/iansmith/parigot/api/guest"
	"github.com/iansmith/parigot/api/shared/id"
	foo "github.com/iansmith/parigot/g/methodcall/foo/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/lib/go/future"
	"github.com/iansmith/parigot/test/func/methodcall/impl/foo/const_"
)

var _ = unsafe.Sizeof([]byte{})

const pathPrefix = "/parigotvirt/"

var logger *slog.Logger

func main() {
	require := []lib.MustRequireFunc{}
	fooServ := &fooServer{}
	binding, fut, ctx, sid := foo.Init(require, fooServ)
	logger = slog.New(guest.NewParigotHandler(sid))
	fut.Success(func(_ *syscall.LaunchResponse) {
		kerr := foo.Run(ctx, binding, foo.TimeoutInMillis, nil)
		slog.Error("error while waiting for foo service calls", "kernel error", syscall.KernelErr_name[int32(kerr)])
	})
}

// this type better implement methodcall.v1.FooService
type fooServer struct {
}

//
// This file contains the true implementations--the server side--for the methods
// defined in foo.proto.
//

func (f *fooServer) AddMultiply(ctx context.Context, req *foo.AddMultiplyRequest) *foo.FutureAddMultiply {
	resp := &foo.AddMultiplyResponse{}
	if req.IsAdd {
		resp.Result = req.Value0 + req.Value1
	} else {
		resp.Result = req.Value0 * req.Value1
	}
	fut := foo.NewFutureAddMultiply()
	fut.Method.CompleteMethod(ctx, resp, foo.FooErr_NoError)
	return fut
}

func (f *fooServer) LucasSequence(ctx context.Context) *foo.FutureLucasSequence {
	logger.Debug("received call for fooServer.LucasSequence")
	resp := &foo.LucasSequenceResponse{}
	seq := make([]int32, const_.LucasSize) // -2 because first two are given
	seq[0] = 2
	seq[1] = 1
	for i := 2; i < const_.LucasSize; i++ {
		seq[i] = seq[i-1] + seq[i-2]
	}
	resp.Sequence = seq
	fut := foo.NewFutureLucasSequence()
	fut.Method.CompleteMethod(ctx, resp, foo.FooErr_NoError)
	return fut
}

// Newton-Raphson method, terms values beyond about 4 are silly
func (f *fooServer) WritePi(ctx context.Context, req *foo.WritePiRequest) *foo.FutureWritePi {
	logger.Debug("received call for fooServer.AddMultiply")

	fut := foo.NewFutureWritePi()
	if req.GetTerms() < 1 || req.GetTerms() > 4 {
		fut.Base.Set(foo.FooErr_BadParamWritePi)
		return fut
	}
	runningTotal := 3.0 // k==0 term

	for k := 1; k <= int(req.GetTerms()); k++ {
		runningTotal = runningTotal - math.Tan(runningTotal)
		logger.Debug("WritePi computing pi", "running total", runningTotal)
	}
	logger.Debug("WritePi finished computing pi", "running total", runningTotal)
	fut.Base.Set(foo.FooErr_NoError)
	return fut
}

// Ready is a check, if this returns false the library will abort and not attempt to run this service.
// Normally this is used to do LocateXXX() calls that are needed for
// the operation of the service.
func (f *fooServer) Ready(_ context.Context, _ id.ServiceId) *future.Base[bool] {
	return future.NewBaseWithValue[bool](true)
}
