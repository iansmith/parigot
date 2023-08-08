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

// Call the real implementation of the Get method for the http service
func (h *myHttpSvc) Get(ctx context.Context, in *http.GetRequest) *http.FutureGet {
	return http.GetHost(ctx, in)
}

// Call the real implementation of the Post method for the http service
func (h *myHttpSvc) Post(ctx context.Context, in *http.PostRequest) *http.FuturePost {
	return http.PostHost(ctx, in)
}

// Call the real implementation of the Put method for the http service
func (h *myHttpSvc) Put(ctx context.Context, in *http.PutRequest) *http.FuturePut {
	return http.PutHost(ctx, in)
}

// Call the real implementation of the Delete method for the http service
func (h *myHttpSvc) Delete(ctx context.Context, in *http.DeleteRequest) *http.FutureDelete {
	return http.DeleteHost(ctx, in)
}

// Call the real implementation of the Head method for the http service
func (h *myHttpSvc) Head(ctx context.Context, in *http.HeadRequest) *http.FutureHead {
	return http.HeadHost(ctx, in)
}

// Call the real implementation of the Connect method for the http service
func (h *myHttpSvc) Connect(ctx context.Context, in *http.ConnectRequest) *http.FutureConnect {
	return http.ConnectHost(ctx, in)
}

// Call the real implementation of the Options method for the http service
func (h *myHttpSvc) Options(ctx context.Context, in *http.OptionsRequest) *http.FutureOptions {
	return http.OptionsHost(ctx, in)
}

// Call the real implementation of the Patch method for the http service
func (h *myHttpSvc) Patch(ctx context.Context, in *http.PatchRequest) *http.FuturePatch {
	return http.PatchHost(ctx, in)
}

// Call the real implementation of the Trace method for the http service
func (h *myHttpSvc) Trace(ctx context.Context, in *http.TraceRequest) *http.FutureTrace {
	return http.TraceHost(ctx, in)
}
