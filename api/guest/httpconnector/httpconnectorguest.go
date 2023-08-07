package main

import (
	"context"
	"unsafe"

	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/httpconnector/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/lib/go/future"
)

var _ = unsafe.Sizeof([]byte{})

func main() {
	ctx := pcontext.NewContextWithContainer(context.Background(), "[httpconnectorguest]main")
	hCnt := &myHttpCnt{}
	binding, fut, _ := httpconnector.Init(ctx, []lib.MustRequireFunc{}, hCnt)
	fut.Success(func(_ *syscall.LaunchResponse) {
		pcontext.Infof(ctx, "httpconnector service guest side started correctly")
	})

	pcontext.Infof(ctx, "start Run in httpconnector guest side")
	kerr := httpconnector.Run(ctx, binding, httpconnector.TimeoutInMillis, nil)

	pcontext.Infof(ctx, "reach the point after Run in httpconnector guest side %s", syscall.KernelErr_name[int32(kerr)])

	pcontext.Errorf(ctx, "error while waiting for httpconnector service calls: %s", syscall.KernelErr_name[int32(kerr)])
}

type myHttpCnt struct{}

func (hCnt *myHttpCnt) Ready(ctx context.Context, _ id.ServiceId) *future.Base[bool] {
	pcontext.Debugf(ctx, "Ready reached in httpconnector service")
	return future.NewBaseWithValue[bool](true)
}

func (hCnt *myHttpCnt) Check(ctx context.Context, in *httpconnector.CheckRequest) *httpconnector.FutureCheck {
	return httpconnector.CheckHost(ctx, in)
}
