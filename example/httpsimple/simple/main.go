package main

import (
	"context"
	"log/slog"

	"github.com/iansmith/parigot/api/guest"
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/example/httpsimple/g/simple/v1"
	"github.com/iansmith/parigot/g/http/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/lib/go/future"

	_ "time/tzdata"
)

var logger *slog.Logger

// All services have a main.  The services should not really "start", however
// until Ready() is called because their dependencies are only guaranteed
// to be up when Ready() is reached.
func main() {

	// the implementation of the service (no state right now)
	impl := &myService{}

	// Init initiaizes a service and normally receives a list of functions
	// that indicate dependencies, but we don't have any here.
	binding, fut, ctx, sid := simple.Init([]lib.MustRequireFunc{http.MustRequire}, impl)

	logger = slog.New(guest.NewParigotHandler(sid))

	// Init() covers the case of a launch failure, so we only deal with success
	fut.Success(func(_ *syscall.LaunchResponse) {
		logger.Info("simple service launched successfully", "simple service host", syscallguest.CurrentHostId().Short())
	})

	// Run waits for calls to our methods and should not return.
	// The context provided here is passed through to calls on your methods.
	kerr := simple.Run(ctx, binding, simple.TimeoutInMillis, nil)

	// Should not happen.
	if kerr != syscall.KernelErr_NoError {
		logger.Error("error caused run to exit in simple", "kernel error", syscall.KernelErr_name[int32(kerr)])
	}

}

// myService is the true implementation of the greeting service.
type myService struct{}

// test at compile time that myService has appropriate methods.
var _ = simple.Simple(&myService{})

func (m *myService) get(ctx context.Context, req *http.GetRequest) (*http.GetResponse, http.HttpErr) {
	msg := []byte("hello, http protocal")
	resp := &http.GetResponse{
		Response: &http.HttpResponse{
			StatusCode:    200,
			Header:        map[string]string{},
			Body:          msg,
			ContentLength: int32(len(msg)),
			Trailer:       map[string]string{},
		},
	}
	return resp, http.HttpErr_NoError
}

func (m *myService) Get(ctx context.Context, req *http.GetRequest) *http.FutureGet {
	resp, err := m.get(ctx, req)
	fut := http.NewFutureGet()
	if err != http.HttpErr_NoError {
		fut.Method.CompleteMethod(ctx, nil, err)
	} else {
		// err is NoError
		fut.Method.CompleteMethod(ctx, resp, err)
	}
	return fut
}

// Ready simply returns an already completed future with the value
// false because it does not have anything to do.  Many Ready()
// functions use this function to MustLocateXXX() calls to obtain
// references to other services.  The second parameter is
// passed here with the ServiceId of myService (the receiver
// of this method call) but it is not needed.
func (m *myService) Ready(_ context.Context, _ id.ServiceId) *future.Base[bool] {
	fut := future.NewBase[bool]()
	fut.Set(true)
	return fut
}