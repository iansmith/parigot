package main

import (
	"context"
	"unsafe"

	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/methodcall/bar/v1"
	"github.com/iansmith/parigot/g/methodcall/foo/v1"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
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
	if kerr != syscall.KernelErr_NoError {
		pcontext.Fatalf(ctx, "RunBar exited with an error: %s", syscall.KernelErr_name[int32(kerr)])
	}
}

// this type better implement methodcall.v1.BarService
type barServer struct {
	foo foo.Client
}

//
// This file contains the true implementations--the server side--for the methods
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
	for i := 0; i < len(req.GetValue()); i++ {
		reqAdd.Value0 = req.GetValue()[i]
		respAdd, multErr := b.foo.AddMultiply(ctx, reqAdd)
		if int32(multErr) != 0 {
			return nil, bar.BarErr_AddMultFailed
		}
		// b.log(nil, pblog.LogLevel_LOG_LEVEL_DEBUG, "ADD (%d,%d) iteration #%d, result add %d",
		// 	reqAdd.GetValue0(), reqAdd.GetValue1(), i, respAdd.Result)
		reqAdd.Value1 = respAdd.GetResult()

		/// multiply
		reqMul.Value0 = req.GetValue()[i]
		respMul, errId := b.foo.AddMultiply(ctx, reqMul)
		if int32(errId) != 0 {
			return nil, bar.BarErr_AddMultFailed
		}
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
