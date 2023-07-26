package main

import (
	"context"
	"unsafe"

	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	file "github.com/iansmith/parigot/g/file/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/lib/go/future"
)

var _ = unsafe.Sizeof([]byte{})

func main() {
	ctx := pcontext.CallTo(pcontext.SourceContext(context.Background(), pcontext.Guest), "fileguest.Main")
	f := &myFileSvc{}
	binding, fut, _ := file.Init(ctx, []lib.MustRequireFunc{}, f)
	fut.Success(func(_ *syscall.LaunchResponse) {
		var kerr syscall.KernelErr
		for {
			kerr = file.ReadOneAndCall(ctx, binding, file.TimeoutInMillis)
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
	})
}

type myFileSvc struct{}

func (f *myFileSvc) Ready(ctx context.Context, _ id.ServiceId) *future.Base[bool] {
	pcontext.Debugf(ctx, "Ready reached in file service")
	return future.NewBaseWithValue[bool](true)
}

func (f *myFileSvc) Open(ctx context.Context, in *file.OpenRequest) *file.FutureOpen {
	return file.OpenHost(ctx, in)
}

func (f *myFileSvc) Create(ctx context.Context, in *file.CreateRequest) *file.FutureCreate {
	return file.CreateHost(ctx, in)

}

func (f *myFileSvc) Close(ctx context.Context, in *file.CloseRequest) *file.FutureClose {
	return file.CloseHost(ctx, in)
}

func (f *myFileSvc) LoadTestData(ctx context.Context, in *file.LoadTestDataRequest) *file.FutureLoadTestData {
	return file.LoadTestDataHost(ctx, in)

}

func (f *myFileSvc) Read(ctx context.Context, in *file.ReadRequest) *file.FutureRead {
	return file.ReadHost(ctx, in)
}

func (f *myFileSvc) Write(ctx context.Context, in *file.WriteRequest) *file.FutureWrite {
	return file.WriteHost(ctx, in)

}

func (f *myFileSvc) Delete(ctx context.Context, in *file.DeleteRequest) *file.FutureDelete {
	return file.DeleteHost(ctx, in)
}
