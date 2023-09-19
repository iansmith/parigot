package main

import (
	"context"
	"encoding/binary"
	"log/slog"

	"github.com/iansmith/parigot/api/guest"
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/nutsdb/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	libnuts "github.com/iansmith/parigot/lib/go/nutsdb"
	nutsdbguest "github.com/iansmith/parigot/lib/go/nutsdb"

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

const firstValue = 1789
const secondValue = 1776

func afterLaunch(ctx context.Context, _ *syscall.LaunchResponse, myServiceId id.ServiceId, fut *syscallguest.LaunchFuture) {

	//client := nutsdb.MustLocate(ctx, myId)

	openfut := libnuts.Open(ctx, "testprogram")
	// if open went ok, we do the first key/value write
	openfut.Method.Success(func(req *nutsdb.OpenResponse) {
		nid := nutsdb.UnmarshalNutsDBId(req.GetNutsdbId())
		logger.Info("nutsdb opened successfully", "nuts db id", nid.Short())

		client := nutsdb.MustLocate(ctx, myServiceId)
		writeFut := nutsdbguest.WritePair(ctx, client, nid, "/foo", "bar", createValue(firstValue))
		writeFut.Method.Failure(func(err nutsdb.NutsDBErr) {
			writeFailed(ctx, int32(err), myServiceId)
		})
		writeFut.Method.Success(func(resp *nutsdb.WritePairResponse) {
			writeSecond(ctx, resp, client, myServiceId, nid)
		})
	})
	//Handle negative outcome of open()
	openfut.Method.Failure(func(err nutsdb.NutsDBErr) {
		logger.Error("failed to open nutsdb", "error", nutsdb.NutsDBErr_name[int32(err)])
		lib.ExitSelf(ctx, 1, myServiceId)
	})
}

func writeFailed(ctx context.Context, err int32, myServiceId id.ServiceId) {
	logger.Error("unable to write pair to db", "error", nutsdb.NutsDBErr_name[err])
	lib.ExitSelf(ctx, 1, myServiceId)
}
func readFailed(ctx context.Context, err int32, myServiceId id.ServiceId) {
	logger.Error("unable to read pair from db", "error", nutsdb.NutsDBErr_name[err])
	lib.ExitSelf(ctx, 1, myServiceId)
}

func writeSecond(ctx context.Context, resp *nutsdb.WritePairResponse, client nutsdb.Client, myServiceId id.ServiceId, nutsId nutsdb.NutsDBId) {
	confirmValueOrFail(ctx, resp.GetPair(), myServiceId, firstValue, "wrote first value")

	//note we are using a different BUCKET here
	secondFut := nutsdbguest.WritePair(ctx, client, nutsId, "/quux", "foo", createValue(secondValue))
	secondFut.Failure(func(err int32) {
		writeFailed(ctx, err, myServiceId)
	})
	secondFut.Method.Success(func(resp *nutsdb.WritePairResponse) {
		readFirst(ctx, resp, client, myServiceId, nutsId)
	})
}
func getValue(pair *nutsdb.Pair) uint64 {
	return binary.LittleEndian.Uint64(pair.GetValue())
}
func createValue(val uint64) []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, val)
	return data
}

func readFirst(ctx context.Context, resp *nutsdb.WritePairResponse, client nutsdb.Client, myServiceId id.ServiceId, nutsId nutsdb.NutsDBId) {
	confirmValueOrFail(ctx, resp.GetPair(), myServiceId, secondValue, "wrote second pair")

	readFut := nutsdbguest.ReadPair(ctx, client, nutsId, "/foo", "bar", nil)
	readFut.Failure(func(err int32) {
		readFailed(ctx, err, myServiceId)
	})
	readFut.Method.Success(func(resp *nutsdb.ReadPairResponse) {
		readSecond(ctx, resp, client, myServiceId, nutsId)
	})
}

func confirmValueOrFail(ctx context.Context, pair *nutsdb.Pair, myServiceId id.ServiceId, expected uint64, msg string) {
	val := getValue(pair)
	if val != expected {
		logger.Error("whoa! expected same value I put in, aborting", "expected", expected, "received", val)
		lib.ExitSelf(ctx, 1, myServiceId)
	}
	logger.Info(msg+" successfully", "bucket path", pair.GetBucketPath(),
		"key", pair.GetKey(), "value", val)

}

func badReadSucceeded(ctx context.Context, myServiceId id.ServiceId) {
	logger.Error("expected to fail reading nutsdb, but did not")
	lib.ExitSelf(ctx, 1, myServiceId)

}
func readSecond(ctx context.Context, resp *nutsdb.ReadPairResponse, client nutsdb.Client, myServiceId id.ServiceId, nutsId nutsdb.NutsDBId) {
	confirmValueOrFail(ctx, resp.GetPair(), myServiceId, firstValue, "read first pair")

	// same key different package
	readFut := nutsdbguest.ReadPair(ctx, client, nutsId, "/quux", "foo", nil)

	readFut.Method.Failure(func(err nutsdb.NutsDBErr) {
		readFailed(ctx, int32(err), myServiceId)
	})
	readFut.Method.Success(func(resp *nutsdb.ReadPairResponse) {
		badReadFirst(ctx, resp, client, myServiceId, nutsId)
	})
}

func badReadFirst(ctx context.Context, resp *nutsdb.ReadPairResponse, client nutsdb.Client, myServiceId id.ServiceId, nutsId nutsdb.NutsDBId) {
	confirmValueOrFail(ctx, resp.GetPair(), myServiceId, secondValue, "read second pair")

	// same key different package
	readFut := nutsdbguest.ReadPair(ctx, client, nutsId, "/quux", "fleazil", nil)
	readFut.Method.Failure(func(err nutsdb.NutsDBErr) {
		badReadSecond(ctx, nutsdb.NutsDBErr(err), client, myServiceId, nutsId)
	})
	readFut.Method.Success(func(resp *nutsdb.ReadPairResponse) {
		badReadSucceeded(ctx, myServiceId)
	})

}
func confirmErrorOrFail(ctx context.Context, err nutsdb.NutsDBErr, expected nutsdb.NutsDBErr, myServiceId id.ServiceId) {
	if err != expected {
		logger.Error("whoa! wrong error value found", "expected", expected, "received", err)
		lib.ExitSelf(ctx, 1, myServiceId)
	}
	logger.Info("successfully checked the error code on a bad read")
}

func badReadSecond(ctx context.Context, err nutsdb.NutsDBErr, client nutsdb.Client, myServiceId id.ServiceId, nutsId nutsdb.NutsDBId) {
	confirmErrorOrFail(ctx, err, nutsdb.NutsDBErr_PairNotFound, myServiceId)

	// same different key correct package
	readFut := nutsdbguest.ReadPair(ctx, client, nutsId, "/quux", "baz", nil)
	readFut.Method.Failure(func(err nutsdb.NutsDBErr) {
		prepareExit(ctx, err, client, myServiceId, nutsId)
	})
	readFut.Method.Success(func(resp *nutsdb.ReadPairResponse) {
		badReadSucceeded(ctx, myServiceId)
	})

}

func prepareExit(ctx context.Context, err nutsdb.NutsDBErr, client nutsdb.Client, myServiceId id.ServiceId, nutsId nutsdb.NutsDBId) {
	confirmErrorOrFail(ctx, err, nutsdb.NutsDBErr_PairNotFound, myServiceId)
	lib.ExitSelf(ctx, 1, myServiceId)
}
