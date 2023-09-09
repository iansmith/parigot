package main

import (
	"context"
	"log/slog"
	"unsafe"

	"github.com/iansmith/parigot/api/guest"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/nutsdb/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/lib/go/future"
)

var _ = unsafe.Sizeof([]byte{})
var logger *slog.Logger

func main() {
	f := &myNutsDB{}
	binding, fut, ctx, sid := nutsdb.Init(nil, f)
	logger = slog.New(guest.NewParigotHandler(sid))

	fut.Success(func(_ *syscall.LaunchResponse) {
		logger.Info("nutsdb service guest side started correctly")
	})
	kerr := nutsdb.Run(ctx, binding, nutsdb.TimeoutInMillis, nil)
	logger.Error("error while waiting for nutsdb service calls", slog.String("syscall.KernelErr", syscall.KernelErr_name[int32(kerr)]))
}

type myNutsDB struct{}

func (f *myNutsDB) Ready(ctx context.Context, _ id.ServiceId) *future.Base[bool] {
	logger.Debug("Ready reached in nutsdb service")
	return future.NewBaseWithValue[bool](true)
}

func (f *myNutsDB) Open(ctx context.Context, in *nutsdb.OpenRequest) *nutsdb.FutureOpen {
	return nutsdb.OpenHost(ctx, in)
}

func (f *myNutsDB) Close(ctx context.Context, in *nutsdb.CloseRequest) *nutsdb.FutureClose {
	return nutsdb.CloseHost(ctx, in)

}

func (f *myNutsDB) ReadPair(ctx context.Context, in *nutsdb.ReadPairRequest) *nutsdb.FutureReadPair {
	return nutsdb.ReadPairHost(ctx, in)
}

func (f *myNutsDB) WritePair(ctx context.Context, in *nutsdb.WritePairRequest) *nutsdb.FutureWritePair {
	return nutsdb.WritePairHost(ctx, in)

}
