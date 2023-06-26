package main

import (
	"context"
	"math"
	"unsafe"

	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"
	foo "github.com/iansmith/parigot/g/methodcall/foo/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/test/func/methodcall/impl/foo/const_"
)

var _ = unsafe.Sizeof([]byte{})

const pathPrefix = "/parigotvirt/"

func main() {
	require := []lib.MustRequireFunc{}
	ctx := pcontext.CallTo(pcontext.SourceContext(pcontext.NewContextWithContainer(context.Background(), "fooservice.Main"), pcontext.Guest), "fooservice.Main")
	fooServ := &fooServer{}
	binding := foo.Init(ctx, require, fooServ)
	kerr := foo.Run(ctx, binding, foo.TimeoutInMillis, nil)
	pcontext.Errorf(ctx, "error while waiting for foo service calls: %s", syscall.KernelErr_name[int32(kerr)])
}

// this type better implement methodcall.v1.FooService
type fooServer struct {
}

//
// This file contains the true implementations--the server side--for the methods
// defined in foo.proto.
//

func (f *fooServer) AddMultiply(ctx context.Context, req *foo.AddMultiplyRequest) (*foo.AddMultiplyResponse, foo.FooErr) {
	resp := &foo.AddMultiplyResponse{}
	if req.IsAdd {
		resp.Result = req.Value0 + req.Value1
	} else {
		resp.Result = req.Value0 * req.Value1
	}
	return resp, foo.FooErr_NoError
}

func (f *fooServer) LucasSequence(ctx context.Context) (*foo.LucasSequenceResponse, foo.FooErr) {
	pcontext.Debugf(ctx, "LucasSequence", "received call for fooServer.LucasSequence")
	resp := &foo.LucasSequenceResponse{}
	seq := make([]int32, const_.LucasSize) // -2 because first two are given
	seq[0] = 2
	seq[1] = 1
	for i := 2; i < const_.LucasSize; i++ {
		seq[i] = seq[i-1] + seq[i-2]
	}
	resp.Sequence = seq
	return resp, foo.FooErr_NoError
}

// Newton-Raphson method, terms values beyond about 4 are silly
func (f *fooServer) WritePi(ctx context.Context, req *foo.WritePiRequest) foo.FooErr {
	pcontext.Debugf(ctx, "WritePi", "received call for fooServer.AddMultiply")

	if req.GetTerms() < 1 || req.GetTerms() > 4 {
		return foo.FooErr_BadParamWritePi
	}
	runningTotal := 3.0 // k==0 term

	for k := 1; k <= int(req.GetTerms()); k++ {
		runningTotal = runningTotal - math.Tan(runningTotal)
		pcontext.Debugf(ctx, "WritePi", "%f", runningTotal)
	}
	pcontext.Infof(ctx, "WritePi result:", "%f", runningTotal)
	return foo.FooErr_NoError
}

// Ready is a check, if this returns false the library will abort and not attempt to run this service.
// Normally this is used to do LocateXXX() calls that are needed for
// the operation of the service.
func (f *fooServer) Ready(_ context.Context, _ id.ServiceId) *lib.BaseFuture[bool] {
	fut := lib.NewBaseFuture[bool]()
	fut.Set(true)
	return fut
}
