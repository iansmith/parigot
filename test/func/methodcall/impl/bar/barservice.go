package main

import (
	"context"
	"unsafe"

	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/methodcall/bar/v1"
	"github.com/iansmith/parigot/g/methodcall/foo/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
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
	foo foo.Client
}

//
// This file contains the true implementation--the server side--for the method
// defined in bar.proto.
//

func (b *barServer) Accumulate(ctx context.Context, req *bar.AccumulateRequest) (*bar.AccumulateResponse, bar.BarErr) {
	resp := &bar.AccumulateResponse{}
	if len(req.Value) == 0 {
		resp.Product = 0
		resp.Sum = 0
		return resp, bar.BarErr_NoError
	}

	reqAdd := &foo.AddMultiplyRequest{
		IsAdd: true,
	}
	reqMul := &foo.AddMultiplyRequest{
		IsAdd: false,
	}

	reqAdd.Value1 = 0 //identity to start
	reqMul.Value1 = 1 // identity to start

	var respAdd, respMul *foo.AddMultiplyResponse
	list := make([]foo.FutureAddMultiply, len(req.GetValue()))
	for i := 0; i < len(req.GetValue()); i++ {
		reqAdd.Value0 = req.GetValue()[i]
		futureAdd := b.foo.AddMultiply(ctx, reqAdd)
		list[i] = *futureAdd
		futureAdd.Failure(func(err foo.FooErr) {
			pcontext.Errorf(ctx, "add multiply of foo failed (add)")
		})
		futureAdd.Success(func(resp *foo.AddMultiplyResponse) {
			reqAdd.Value1 = respAdd.GetResult()
			/// multiply
			reqMul.Value0 = req.GetValue()[i]
		})
		futureMultiply := b.foo.AddMultiply(ctx, reqMul)
		futureMultiply.Failure(func(err foo.FooErr) {
			pcontext.Errorf(ctx, "add multiply of foo failed (mult)")
		})
		// b.log(nil, pblog.LogLevel_LOG_LEVEL_DEBUG, "MUL (%d,%d) iteration #%d, result mul %d",
		// 	reqMul.GetValue0(), reqMul.GetValue1(), i, respMul.Result)
		reqMul.Value1 = respMul.GetResult()
	}

	resp.Product = respMul.GetResult()
	resp.Sum = respAdd.GetResult()
	return resp, bar.BarErr_NoError
}

// Ready is a check, if this returns false the library will abort and not attempt to run this service.
// Normally this is used to do LocateXXX() calls that are needed for
// the operation of the service.
func (b *barServer) Ready(ctx context.Context, sid id.ServiceId) *lib.BaseFuture[bool] {
	b.foo = foo.MustLocate(ctx, sid)
	f := lib.NewBaseFuture[bool]()
	f.Set(true)
	return f
}
