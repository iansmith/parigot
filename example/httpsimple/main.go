package main

import (
	"context"
	"log/slog"

	"github.com/iansmith/parigot/api/guest"
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/example/httpsimple/g/frontdoor/v1"
	"github.com/iansmith/parigot/example/httpsimple/g/simple/v1"
	"github.com/iansmith/parigot/g/httpconnector/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/lib/go/future"

	_ "time/tzdata"
)

//
// The front door service that receives the initial request. It then forwards the
// received request to the simple service via parigot.
//

var timeoutInMillis = int32(50)
var logger *slog.Logger

type myService struct {
	simple simple.Client
}

var _ frontdoor.Frontdoor = &myService{}

func main() {

	// the implementation of the service (no state right now)
	impl := &myService{}

	// Init initiaizes a service and normally receives a list of functions
	// that indicate dependencies, but we don't have any here.
	_, fut, ctx, sid :=
		frontdoor.Init([]lib.MustRequireFunc{simple.MustRequire, httpconnector.MustRequire}, impl)

	logger = slog.New(guest.NewParigotHandler(sid))

	fut.Failure(func(err syscall.KernelErr) {
		logger.Error("failed to launch the frontdoor service: %s", syscall.KernelErr_name[int32(err)])
		lib.ExitSelf(ctx, 1, sid)
	})

	fut.Success(func(resp *syscall.LaunchResponse) {
		logger.Info("frontdoor launched successfully",
			slog.String("frontdoor host", syscallguest.CurrentHostId().Short()))
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

	lib.ExitSelf(ctx, 1, sid)
}

func (m *myService) handle(ctx context.Context, req *httpconnector.HandleRequest) (*httpconnector.HandleResponse, frontdoor.FrontdoorErr) {
	return &httpconnector.HandleResponse{}, frontdoor.FrontdoorErr_NoError
}

func (m *myService) Handle(ctx context.Context, req *httpconnector.HandleRequest) *frontdoor.FutureHandle {
	return frontdoor.NewFutureHandle()
}

func (m *myService) Ready(ctx context.Context, sid id.ServiceId) *future.Base[bool] {
	m.simple = simple.MustLocate(ctx, sid)
	return future.NewBaseWithValue[bool](true)
}
