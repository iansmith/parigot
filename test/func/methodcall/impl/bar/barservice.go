package main

import (
	"context"
	"unsafe"

	"github.com/iansmith/parigot/g/methodcall/v1"
	lib "github.com/iansmith/parigot/lib/go"

	methodcallmsg "github.com/iansmith/parigot/g/msg/methodcall/v1"
)

var _ = unsafe.Sizeof([]byte{})

func main() {

}

//go:export parigot_main
//go:linkname parigot_main
func parigot_main() {
	lib.FlagParseCreateEnv()

	bg := context.Background()
	methodcall.ExportBarServiceOrPanic()
	methodcall.RequireFooServiceOrPanic(bg)
	s := &barServer{}
	methodcall.RunBarService(s)
}

// this type better implement methodcall.v1.BarService
type barServer struct {
	foo methodcall.FooServiceClient
}

//
// This file contains the true implementations--the server side--for the methods
// defined in bar.proto.
//

func (b *barServer) Accumulate(ctx context.Context, req *methodcallmsg.AccumulateRequest) (*methodcallmsg.AccumulateResponse, error) {
	resp := &methodcallmsg.AccumulateResponse{}
	if len(req.Value) == 0 {
		resp.Product = 0
		resp.Sum = 0
		return resp, nil
	}

	reqAdd := &methodcallmsg.AddMultiplyRequest{
		IsAdd: true,
	}
	reqMul := &methodcallmsg.AddMultiplyRequest{
		IsAdd: false,
	}

	reqAdd.Value1 = 0 //identity to start
	reqMul.Value1 = 1 // identity to start

	var respAdd, respMul *methodcallmsg.AddMultiplyResponse
	var err error
	for i := 0; i < len(req.GetValue()); i++ {
		reqAdd.Value0 = req.GetValue()[i]
		respAdd, err = b.foo.AddMultiply(reqAdd)
		if err != nil {
			return nil, err
		}
		// b.log(nil, pblog.LogLevel_LOG_LEVEL_DEBUG, "ADD (%d,%d) iteration #%d, result add %d",
		// 	reqAdd.GetValue0(), reqAdd.GetValue1(), i, respAdd.Result)
		reqAdd.Value1 = respAdd.GetResult()

		/// multiply
		reqMul.Value0 = req.GetValue()[i]
		respMul, err = b.foo.AddMultiply(reqMul)
		if err != nil {
			return nil, err
		}
		// b.log(nil, pblog.LogLevel_LOG_LEVEL_DEBUG, "MUL (%d,%d) iteration #%d, result mul %d",
		// 	reqMul.GetValue0(), reqMul.GetValue1(), i, respMul.Result)
		reqMul.Value1 = respMul.GetResult()
	}
	resp.Product = respMul.GetResult()
	resp.Sum = respAdd.GetResult()
	return resp, nil
}

// Ready is a check, if this returns false the library will abort and not attempt to run this service.
// Normally, this is used to block using the lib.Run() call.  This call will wait until all the required
// services are ready.
func (b *barServer) Ready(ctx context.Context) bool {
	methodcall.WaitBarServiceOrPanic()
	b.foo = methodcall.LocateFooServiceOrPanic(ctx)
	return true
}
