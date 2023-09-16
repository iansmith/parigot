package main

import (
	"context"
	"log/slog"

	"github.com/iansmith/parigot/api/guest"
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/nutsdb/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	libnuts "github.com/iansmith/parigot/lib/go/nutsdb"

	_ "time/tzdata"
)

var timeoutInMillis = int32(50)
var logger *slog.Logger

func main() {

	// get context and logger
	ctx, myId := lib.MustInitClient([]lib.MustRequireFunc{nutsdb.MustRequire})
	fut := lib.LaunchClient(ctx, myId)
	logger = slog.New(guest.NewParigotHandler(myId))

	fut.Failure(func(err syscall.KernelErr) {
		logger.Error("failed to launch the kv service", "kernel error", syscall.KernelErr_name[int32(err)])
		lib.ExitSelf(ctx, 1, myId)
	})

	fut.Success(func(resp *syscall.LaunchResponse) {
		logger.Info("kv launched successfully",
			slog.String("kv host", syscallguest.CurrentHostId().Short()))
		afterLaunch(ctx, resp, myId, fut)
	})

	// This loop should never return.  Timeout in millis is used
	// for the question of how long should we "wait" for a network call
	// before doing something else.
	var err syscall.KernelErr
	for {
		err = lib.ReadOneAndCallClient(ctx, nil, timeoutInMillis)
		if err != syscall.KernelErr_NoError && err != syscall.KernelErr_ReadOneTimeout {
			break
		}
	}
	logger.Error("we got an error that is not normal from ReadOneAndCallClient", "kernel error", syscall.KernelErr_name[int32(err)])

	lib.ExitSelf(ctx, 1, myId)
}

func afterLaunch(ctx context.Context, _ *syscall.LaunchResponse, myId id.ServiceId, fut *syscallguest.LaunchFuture) {

	//client := nutsdb.MustLocate(ctx, myId)

	openfut := libnuts.Open(ctx, "testprogram")

	openfut.Method.Success(func(req *nutsdb.OpenResponse) {
		nid := nutsdb.UnmarshalNutsDBId(req.NutsdbId)
		logger.Info("nutsdb opened successfully", "nid", nid.Short())
	})

	//Handle negative outcome.
	openfut.Method.Failure(func(err nutsdb.NutsDBErr) {
		logger.Error("failed to open nutsdb", "error", nutsdb.NutsDBErr_name[int32(err)])
		lib.ExitSelf(ctx, 1, myId)
	})
}
