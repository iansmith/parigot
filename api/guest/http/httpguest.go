package main

import (
	"context"
	"unsafe"

	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/http/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/lib/go/future"
)

var _ = unsafe.Sizeof([]byte{})

func main() {
	ctx := pcontext.NewContextWithContainer(context.Background(), "[httpguest]main")
	h := &myHttpSvc{}
	binding, fut, _ := http.Init(ctx, []lib.MustRequireFunc{}, h)
	fut.Success(func(_ *syscall.LaunchResponse) {
		pcontext.Infof(ctx, "http service guest side started correctly")
	})
	kerr := http.Run(ctx, binding, http.TimeoutInMillis, nil)
	pcontext.Errorf(ctx, "error while waiting for http service calls: %s", syscall.KernelErr_name[int32(kerr)])
}

type myHttpSvc struct{}

func (h *myHttpSvc) Ready(ctx context.Context, _ id.ServiceId) *future.Base[bool] {
	pcontext.Debugf(ctx, "Ready reached in http service")
	return future.NewBaseWithValue[bool](true)
}

func (h *myHttpSvc) Get(ctx context.Context, in *http.GetRequest) *http.FutureGet {
	return http.GetHost(ctx, in)
}
