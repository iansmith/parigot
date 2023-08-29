package main

import (
	"context"
	"log/slog"
	"unsafe"

	"github.com/iansmith/parigot/api/guest"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/httpconnector/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/lib/go/future"
)

var _ = unsafe.Sizeof([]byte{})

var logger *slog.Logger

func main() {
	//hCnt := &myConnector{}
	binding, fut, ctx, sid := httpconnector.Init([]lib.MustRequireFunc{}, nil)
	logger = slog.New(guest.NewParigotHandler(sid))
	fut.Success(func(_ *syscall.LaunchResponse) {
		logger.Info("httpconnector service guest side started correctly")
	})
	// this is effectively a useless loop but without it, we would exit
	kerr := httpconnector.Run(ctx, binding, httpconnector.TimeoutInMillis, nil)
	logger.Error("error while waiting for httpconnector service calls", "kernel error", syscall.KernelErr_name[int32(kerr)])
}

type myConnector struct{}

func (_ *myConnector) Ready(ctx context.Context, _ id.ServiceId) *future.Base[bool] {
	logger.Info("Ready reached in httpconnector service")
	return future.NewBaseWithValue[bool](true)
}
