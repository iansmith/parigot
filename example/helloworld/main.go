package main

import (
	"context"
	"log/slog"

	"github.com/iansmith/parigot/api/guest"
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/example/helloworld/g/greeting/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
)

var timeoutInMillis = int32(50)

func main() {

	// Logging system needs help with panics, so we trap and Dump() the log data.
	defer func() {
		if r := recover(); r != nil {
			if r == apishared.ControlledExit {
				return
			}
			slog.Error("helloworld: trapped a panic in the guest side:", "recovered error:", r)
		}
	}()
	// get context and logger
	ctx, myId := lib.MustInitClient([]lib.MustRequireFunc{greeting.MustRequire})
	fut := lib.LaunchClient(ctx, myId)

	fut.Failure(func(err syscall.KernelErr) {
		guest.Log(ctx).Error("failed to launch the hello world service: %s", syscall.KernelErr_name[int32(err)])
		lib.ExitClient(ctx, 1, myId, "exiting due to launch failure", "unable to exit, forcing exit with os.Exit() after failure to Launch()")
	})

	fut.Success(func(resp *syscall.LaunchResponse) {
		guest.Log(ctx).Info("hello world launched successfully",
			slog.String("hello world host", syscallguest.CurrentHostId().Short()))
		afterLaunch(ctx, resp, myId, fut)
	})

	// MustRunClient should never return.  Timeout in millis is used
	// for the question of how long should we "wait" for a network call
	// before doing something else.
	var err syscall.KernelErr
	for {
		err = lib.ReadOneAndCallClient(ctx, nil, timeoutInMillis)
		if err != syscall.KernelErr_NoError && err != syscall.KernelErr_ReadOneTimeout {
			break
		}
	}
	guest.Log(ctx).Error("we got an error that is not normal: %s", syscall.KernelErr_name[int32(err)])

	// Should not happen.
	guest.Log(ctx).Error("failed inside run: %s", syscall.KernelErr_name[int32(err)])
	lib.ExitClient(ctx, 1, myId, "failed in MustRunClient", "failed trying to exit after failure in MustRunClient")
}

func afterLaunch(ctx context.Context, _ *syscall.LaunchResponse, myId id.ServiceId, fut *syscallguest.LaunchFuture) {

	greetService := greeting.MustLocate(ctx, myId)

	req := &greeting.FetchGreetingRequest{
		Tongue: greeting.Tongue_French,
	}

	// Make the call to the greeting service.
	greetFuture := greetService.FetchGreeting(ctx, req)

	// Handle positive outcome.
	greetFuture.Method.Success(func(response *greeting.FetchGreetingResponse) {
		guest.Log(ctx).Info(response.GetGreeting() + " world")
		lib.ExitClient(ctx, 1, myId, "exiting after successful call to greeting.FetchGreeting",
			"failed trying to exit after success, so forcing exit with os.Exit()")
	})

	//Handle negative outcome.
	greetFuture.Method.Failure(func(err greeting.GreetErr) {
		guest.Log(ctx).Error("failed to fetch greeting: %s", greeting.GreetErr_name[int32(err)])
		lib.ExitClient(ctx, 1, myId, "exiting because we failed to call greet.FetchGreeting",
			"tried to exit after failed call to greet.FetchGreeting, failed so forcing exit with os.Exit()")
	})
}
