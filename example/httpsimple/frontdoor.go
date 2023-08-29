package main

import (
	"context"
	"log"
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

var _ httpconnector.HttpConnector = &myService{}

func main() {

	// the implementation of the service (no state right now)
	impl := &myService{}

	// Init initiaizes a service and normally receives a list of functions
	// that indicate dependencies, but we don't have any here.
	_, fut, ctx, sid :=
		frontdoor.Init([]lib.MustRequireFunc{simple.MustRequire, httpconnector.MustRequire}, impl)

	// create a logger
	logger = slog.New(guest.NewParigotHandler(sid))

	smmap := kernelSetup(ctx, sid, impl)

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
		err = lib.ReadOneAndCallClient(ctx, smmap, timeoutInMillis)
		if err != syscall.KernelErr_NoError && err != syscall.KernelErr_ReadOneTimeout {
			break
		}
	}
	logger.Error("we got an error that is not normal from ReadOneAndCallClient", "kernel error", syscall.KernelErr_name[int32(err)])

	lib.ExitSelf(ctx, 1, sid)
}

func (m *myService) Handle(ctx context.Context, req *httpconnector.HandleRequest) *httpconnector.FutureHandle {
	log.Printf("IN HANDLE")
	return httpconnector.NewFutureHandle()
}

func (m *myService) Ready(ctx context.Context, sid id.ServiceId) *future.Base[bool] {
	m.simple = simple.MustLocate(ctx, sid)
	return future.NewBaseWithValue[bool](true)
}

// kernelSetup is a utility routine for doing some tasks that are needed,
// but normally generated by parigot at build time.  However, because we are
// using a reverse API we have to induce the structures we want inside the
// kernel.
func kernelSetup(ctx context.Context, sid id.ServiceId, impl httpconnector.HttpConnector) *lib.ServiceMethodMap {
	// synchronous: Export the name of the reverse interface
	_, kerr := lib.Export1(ctx, "httpconnector.v1", "httpconnector", sid)
	if kerr != syscall.KernelErr_NoError {
		logger.Error("unable to export the httpconnector interface")
		lib.ExitSelf(ctx, 1, sid)
	}

	// sync: bind the method we implement
	bindReq := &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Handle"
	bindResp, err := syscallguest.BindMethod(ctx, bindReq)
	if err != syscall.KernelErr_NoError {
		logger.Error("unable to bind the httpconnector to Handle implementation")
		lib.ExitSelf(ctx, 1, sid)
	}
	mid := id.UnmarshalMethodId(bindResp.GetMethodId())
	smmap := lib.NewServiceMethodMap()

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid, mid, "HttpConnector", "Handle", httpconnector.GenerateHandleInvoker(impl))

	logger.Info("success binding method 'Handle'", "method", mid.Short())
	return smmap

}