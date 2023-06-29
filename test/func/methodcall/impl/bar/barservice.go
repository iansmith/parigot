package main

import (
	"context"
	"unsafe"

	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/methodcall/bar/v1"
	"github.com/iansmith/parigot/g/methodcall/foo/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/lib/go/future"
)

var _ = unsafe.Sizeof([]byte{})

type FullyQualifiedServiceName struct {
	PackageName, ServiceName string
}

func main() {
	req := []lib.MustRequireFunc{
		foo.MustRequire,
	}
	ctx := pcontext.CallTo(pcontext.SourceContext(context.Background(), pcontext.Guest), "barservice.Main")
	bServer := &barServer{}
	binding := bar.Init(ctx, req, bServer)
	kerr := bar.Run(ctx, binding, bar.TimeoutInMillis, nil)
	pcontext.Fatalf(ctx, "Run in Bar exited with an error: %s", syscall.KernelErr_name[int32(kerr)])

}

// this type better implement methodcall.v1.BarService
type barServer struct {
	fooClient foo.Client
}

var _ = bar.Bar(&barServer{})

//
// This file contains the true implementation--the server side--for the method
// defined in bar.proto.
//

// helper for generating success functions
func (b *barServer) generator(fn accumulator, ctx context.Context, req *bar.AccumulateRequest, i int, result *bar.FutureAccumulate) addMultSucc {
	r := func(r *foo.AddMultiplyResponse) {
		sf := b.generator(b.succFn, ctx, req, i+1, result)
		fut := fn(ctx, req, r, i, result)
		fut.Method.Success(sf)
	}
	return r
}

// addMultSucc is the typename of the "base" success function that the foo.FutureAddMultiply
// is expecting... this is the RESULT of the b.generator() function
type addMultSucc func(*foo.AddMultiplyResponse)

// convenience type to make names more clear, this what our internal model of
// a success function looks like, and it needs a bunch of params... this type gets
// wrapped and passed params by b.generator
type accumulator func(ctx context.Context, req *bar.AccumulateRequest, resp *foo.AddMultiplyResponse, i int, result *bar.FutureAccumulate) *foo.FutureAddMultiply

// our success function, each call uses successive versions of this func with i increasing by one but otherwise
// the params staying the same
func (b *barServer) succFn(ctx context.Context, req *bar.AccumulateRequest, resp *foo.AddMultiplyResponse, i int, result *bar.FutureAccumulate) *foo.FutureAddMultiply {
	if i < len(req.GetValue()) {
		// we just finished ith iter
		reqAdd := &foo.AddMultiplyRequest{}
		reqAdd.Value1 = resp.GetResult() //
		reqAdd.Value0 = req.GetValue()[i+1]
		reqAdd.IsAdd = true
		fut := b.fooClient.AddMultiply(ctx, reqAdd)
		fut.Method.Failure(func(err foo.FooErr) {
			result.Method.CompleteMethod(ctx, nil, bar.BarErr_AddMultFailed)
		})
		return fut
	}
	// finished all the given elements in req.GetValue()
	result.Method.Success(func(resp *bar.AccumulateResponse) {
		pcontext.Infof(ctx, "success for sum! %d", resp.GetSum())
	})
	return nil // wont be consumed
}

func (b *barServer) Accumulate(ctx context.Context, req *bar.AccumulateRequest) *bar.FutureAccumulate {
	//trivial case
	if len(req.Value) == 0 {
		resp := &bar.AccumulateResponse{}
		resp.Product = 1
		resp.Sum = 0
		f := bar.NewFutureAccumulate()
		f.CompleteMethod(ctx, resp, 0)
		return f
	}

	// ultimately this is the future for the WHOLE accumulate
	finalFut := bar.NewFutureAccumulate()

	// if anything goes wrong, we alse fail the finalFut
	finalFut.Method.Failure(func(err bar.BarErr) {
		pcontext.Errorf(ctx, "unable to compute accumulation sum: %s", bar.BarErr_name[int32(err)])
	})

	// initial startup
	reqAdd := &foo.AddMultiplyRequest{}
	reqAdd.Value1 = 0 //identity add
	reqAdd.Value0 = req.GetValue()[0]
	startSuccess := b.generator(b.succFn, ctx, req, 0, finalFut)
	futAdd := b.fooClient.AddMultiply(ctx, reqAdd)
	futAdd.Method.Success(startSuccess)

	return finalFut
}

// Ready is a check, if this returns false the library will abort and not attempt to run this service.
// Normally this is used to do LocateXXX() calls that are needed for
// the operation of the service.
func (b *barServer) Ready(ctx context.Context, sid id.ServiceId) *future.Base[bool] {
	b.fooClient = foo.MustLocate(ctx, sid)
	return future.NewBaseWithValue[bool](true)
}
