package main

import (
	"context"
	"log/slog"
	"unsafe"

	"github.com/iansmith/parigot/api/guest"
	"github.com/iansmith/parigot/api/shared/id"
	file "github.com/iansmith/parigot/g/file/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/lib/go/future"
)

var _ = unsafe.Sizeof([]byte{})
var logger *slog.Logger

func main() {
	f := &myFileSvc{}
	binding, fut, ctx, sid := file.Init([]lib.MustRequireFunc{}, f)
	logger = slog.New(guest.NewParigotHandler(sid))

	fut.Success(func(_ *syscall.LaunchResponse) {
		logger.Info("file service guest side started correctly")
	})
	kerr := file.Run(ctx, binding, file.TimeoutInMillis, nil)
	logger.Error("error while waiting for file service calls", slog.String("syscall.KernelErr", syscall.KernelErr_name[int32(kerr)]))
}

type myFileSvc struct{}

func (f *myFileSvc) Ready(ctx context.Context, _ id.ServiceId) *future.Base[bool] {
	logger.Debug("Ready reached in file service")
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

func (f *myFileSvc) Stat(ctx context.Context, in *file.StatRequest) *file.FutureStat {
	return file.StatHost(ctx, in)
}
