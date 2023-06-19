package main

import (
	"context"
	"math"
	"unsafe"

	"github.com/iansmith/parigot/apishared/id"
	"github.com/iansmith/parigot/apiwasm"
	pcontext "github.com/iansmith/parigot/context"
	methodcall "github.com/iansmith/parigot/g/methodcall/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/test/func/methodcall/impl/foo/const_"
)

var _ = unsafe.Sizeof([]byte{})

const pathPrefix = "/parigotvirt/"

func main() {
	req := []apiwasm.MustRequireFunc{}
	ctx := pcontext.CallTo(pcontext.SourceContext(pcontext.NewContextWithContainer(context.Background(), "fooservice.Main"), pcontext.Guest), "fooservice.Main")
	foo := &fooServer{}
	binding := methodcall.InitFoo(ctx, req, foo)
	var kerr syscall.KernelErr
	for {
		if len(binding.Pair()) == 0 {
			panic("BINDING EMPTY")
		}
		kerr = methodcall.ReadOneAndCallFoo(ctx, binding, methodcall.TimeoutInMillisFoo)
		if kerr == syscall.KernelErr_ReadOneTimeout {
			pcontext.Infof(ctx, "waiting for calls to foo service")
			continue
		}
		if kerr == syscall.KernelErr_NoError {
			continue
		}
		break
	}
	pcontext.Errorf(ctx, "error while waiting for foo service calls: %s", syscall.KernelErr_name[int32(kerr)])

}

// this type better implement methodcall.v1.FooService
type fooServer struct {
}

//
// This file contains the true implementations--the server side--for the methods
// defined in foo.proto.
//

func (f *fooServer) AddMultiply(ctx context.Context, req *methodcall.AddMultiplyRequest) (*methodcall.AddMultiplyResponse,
	methodcall.MethodCallSuiteErr) {
	//f.log(pctx, pblog.LogLevel_LOG_LEVEL_DEBUG, "received call for fooServer.AddMultiply")
	resp := &methodcall.AddMultiplyResponse{}
	if req.IsAdd {
		resp.Result = req.Value0 + req.Value1
	} else {
		resp.Result = req.Value0 * req.Value1
	}
	return resp, methodcall.MethodCallSuiteErr_NoError
}

func (f *fooServer) LucasSequence(ctx context.Context) (*methodcall.LucasSequenceResponse, methodcall.MethodCallSuiteErr) {
	pcontext.Debugf(ctx, "LucasSequence", "received call for fooServer.LucasSequence")
	resp := &methodcall.LucasSequenceResponse{}
	seq := make([]int32, const_.LucasSize) // -2 because first two are given
	seq[0] = 2
	seq[1] = 1
	for i := 2; i < const_.LucasSize; i++ {
		seq[i] = seq[i-1] + seq[i-2]
	}
	resp.Sequence = seq
	return resp, methodcall.MethodCallSuiteErr_NoError
}

// Newton-Raphson method, terms values beyond about 4 are silly
func (f *fooServer) WritePi(ctx context.Context, req *methodcall.WritePiRequest) methodcall.MethodCallSuiteErr {
	pcontext.Debugf(ctx, "WritePi", "received call for fooServer.AddMultiply")

	if req.GetTerms() < 1 {
		return methodcall.MethodCallSuiteErr_BadTerms
	}
	runningTotal := 3.0 // k==0 term

	for k := 1; k <= int(req.GetTerms()); k++ {
		runningTotal = runningTotal - math.Tan(runningTotal)
		pcontext.Debugf(ctx, "WritePi", "%f", runningTotal)
	}
	pcontext.Infof(ctx, "WritePi result:", "%f", runningTotal)
	return methodcall.MethodCallSuiteErr_NoError
}

// Ready is a check, if this returns false the library will abort and not attempt to run this service.
// Normally this is used to do LocateXXX() calls that are needed for
// the operation of the service.
func (f *fooServer) Ready(_ context.Context, _ id.ServiceId) bool {
	return true
}
