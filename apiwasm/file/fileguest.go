package main

import (
	"context"
	"unsafe"

	"github.com/iansmith/parigot/apishared/id"
	"github.com/iansmith/parigot/apiwasm"
	pcontext "github.com/iansmith/parigot/context"
	file "github.com/iansmith/parigot/g/file/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
)

var _ = unsafe.Sizeof([]byte{})

func main() {
	ctx := pcontext.CallTo(pcontext.SourceContext(context.Background(), pcontext.Guest), "fileguest.Main")
	f := &myFileSvc{}
	binding := file.Init(ctx, []apiwasm.MustRequireFunc{}, f)
	var kerr syscall.KernelErr
	for {
		kerr = file.ReadOneAndCall(ctx, binding, file.TimeoutInMillisFile)
		if kerr == syscall.KernelErr_ReadOneTimeout {
			pcontext.Infof(ctx, "waiting for calls to file service")
			continue
		}
		if kerr == syscall.KernelErr_NoError {
			continue
		}
		break
	}
	pcontext.Errorf(ctx, "error while waiting for file service calls: %s", syscall.KernelErr_name[int32(kerr)])
}

type myFileSvc struct{}

func (f *myFileSvc) Ready(ctx context.Context, _ id.ServiceId) *lib.Future {
	pcontext.Debugf(ctx, "Ready reached in file service")
	ff := lib.NewFuture[any, file.FileErr](nil, nil)
	ff.CompleteCall(nil, 0)
	return ff
}

func (f *myFileSvc) Open(ctx context.Context, in *file.OpenRequest) (*file.OpenResponse, file.FileErr) {
	return file.OpenHost(in)
}

func (f *myFileSvc) Create(ctx context.Context, in *file.CreateRequest) (*file.CreateResponse, file.FileErr) {
	return file.CreateHost(in)

}

func (f *myFileSvc) Close(ctx context.Context, in *file.CloseRequest) (*file.CloseResponse, file.FileErr) {
	return file.CloseHost(in)
}

func (f *myFileSvc) LoadTestData(ctx context.Context, in *file.LoadTestDataRequest) (*file.LoadTestDataResponse, file.FileErr) {
	return file.LoadTestDataHost(in)

}
