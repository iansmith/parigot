package main

import (
	"context"
	"log"
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
	ctx := pcontext.CallTo(pcontext.SourceContext(context.Background(), pcontext.Guest), "httpguest.Main")

	h := &myHttpSvc{}

	_, fut, _ := http.Init(ctx, []lib.MustRequireFunc{}, h)

	panic("test1")
	fut.Success(func(_ *syscall.LaunchResponse) {
		log.Println("http service launched successfully")
	})
}

type myHttpSvc struct{}

func (h *myHttpSvc) Ready(ctx context.Context, _ id.ServiceId) *future.Base[bool] {
	pcontext.Debugf(ctx, "Ready reached in http service")
	return future.NewBaseWithValue[bool](true)
}

func (h *myHttpSvc) Get(ctx context.Context, in *http.GetRequest) *http.FutureGet {
	return http.GetHost(ctx, in)
}
