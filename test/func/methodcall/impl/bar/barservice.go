package main

import (
	"context"
	"os"
	"unsafe"

	apiwasm "github.com/iansmith/parigot/apiwasm"
	pcontext "github.com/iansmith/parigot/context"
	methodcall "github.com/iansmith/parigot/g/methodcall/v1"
)

var _ = unsafe.Sizeof([]byte{})

func main() {
	ctx := apiwasm.ManufactureGuestContext("[barservice]main")
	defer func() {
		pcontext.Dump(ctx)
		pcontext.Debugf(ctx, "trapped a panic in the guest side")
		os.Exit(1)
	}()

	myId := methodcall.MustRegisterBar(ctx)
	methodcall.MustExportBar(ctx)
	methodcall.MustRequireFoo(ctx, myId)
	methodcall.MustWaitSatisfiedBar(myId)
	b := &barServer{}
	b.foo = methodcall.MustLocateFoo(ctx, myId)
	methodcall.RunBar(ctx, myId, b)
}

// this type better implement methodcall.v1.BarService
type barServer struct {
	foo methodcall.Foo
}

//
// This file contains the true implementations--the server side--for the methods
// defined in bar.proto.
//

func (b *barServer) Accumulate(ctx context.Context, req *methodcall.AccumulateRequest) (*methodcall.AccumulateResponse,
	methodcall.MethodCallSuiteErr) {
	resp := &methodcall.AccumulateResponse{}
	if len(req.Value) == 0 {
		resp.Product = 0
		resp.Sum = 0
		return resp, methodcall.MethodCallSuiteErr_NoError
	}

	reqAdd := &methodcall.AddMultiplyRequest{
		IsAdd: true,
	}
	reqMul := &methodcall.AddMultiplyRequest{
		IsAdd: false,
	}

	reqAdd.Value1 = 0 //identity to start
	reqMul.Value1 = 1 // identity to start

	var respAdd, respMul *methodcall.AddMultiplyResponse
	for i := 0; i < len(req.GetValue()); i++ {
		reqAdd.Value0 = req.GetValue()[i]
		respAdd, multErr := b.foo.AddMultiply(ctx, reqAdd)
		if int32(multErr) != 0 {
			return nil, multErr
		}
		// b.log(nil, pblog.LogLevel_LOG_LEVEL_DEBUG, "ADD (%d,%d) iteration #%d, result add %d",
		// 	reqAdd.GetValue0(), reqAdd.GetValue1(), i, respAdd.Result)
		reqAdd.Value1 = respAdd.GetResult()

		/// multiply
		reqMul.Value0 = req.GetValue()[i]
		respMul, errId := b.foo.AddMultiply(ctx, reqMul)
		if int32(errId) != 0 {
			return nil, errId
		}
		// b.log(nil, pblog.LogLevel_LOG_LEVEL_DEBUG, "MUL (%d,%d) iteration #%d, result mul %d",
		// 	reqMul.GetValue0(), reqMul.GetValue1(), i, respMul.Result)
		reqMul.Value1 = respMul.GetResult()
	}
	resp.Product = respMul.GetResult()
	resp.Sum = respAdd.GetResult()
	return resp, methodcall.MethodCallSuiteErr_NoError
}

// Ready is a check, if this returns false the library will abort and not attempt to run this service.
// Normally, this is used to block using the lib.Run() call.  This call will wait until all the required
// services are ready.
func (b *barServer) Ready(ctx context.Context) bool {
	return true
}
