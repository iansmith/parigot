package main

import (
	"context"
	"fmt"
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
	// that indicate dependencies, but we don't have any here. We *ARE* the http
	// implementation, we don't depend on it.
	binding, fut, ctx, sid := http.Init([]lib.MustRequireFunc{}, impl)

	logger = slog.New(guest.NewParigotHandler(sid))

	// Init() covers the case of a launch failure, so we only deal with success
	fut.Success(func(_ *syscall.LaunchResponse) {
		logger.Info("simple service launched successfully", "simple service host", syscallguest.CurrentHostId().Short())
	})

	// Run waits for calls to our methods and should not return.
	// The context provided here is passed through to calls on your methods.
	kerr := http.Run(ctx, binding, simple.TimeoutInMillis, nil)

	// Should not happen.star
	if kerr != syscall.KernelErr_NoError {
		logger.Error("error caused run to exit in simple", "kernel error", syscall.KernelErr_name[int32(kerr)])
	}

}

// myService is the true implementation of the greeting service.
type myService struct {
	myId id.ServiceId
}

// test at compile time that myService has appropriate methods.
var _ = http.Http(&myService{})
var start = []byte("hello, http protocol")

func (m *myService) get(ctx context.Context, req *http.GetRequest) (*http.GetResponse, http.HttpErr) {
	msg := fmt.Sprintf("%s (%s %s)\n", start, "GET", req.Request.Url)
	resp := &http.GetResponse{
		Response: &http.HttpResponse{
			StatusCode:    200,
			Header:        map[string]string{},
			Body:          []byte(msg),
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
// true because it does not have anything to do.  Many Ready()
// functions use this function to MustLocateXXX() calls to obtain
// references to other services.  The second parameter is
// passed here with the ServiceId of myService (the receiver
// of this method call) but it is not needed.
func (m *myService) Ready(_ context.Context, sid id.ServiceId) *future.Base[bool] {
	m.myId = sid
	fut := future.NewBase[bool]()
	fut.Set(true)
	return fut
}

//
// These methods are required by the API, but are not used.
//

func (m *myService) Post(ctx context.Context, req *http.PostRequest) *http.FuturePost {
	panic("Post")
}

func (m *myService) Put(ctx context.Context, req *http.PutRequest) *http.FuturePut {
	panic("Put")
}

func (m *myService) Delete(ctx context.Context, req *http.DeleteRequest) *http.FutureDelete {
	panic("Delete")
}

func (m *myService) Head(ctx context.Context, req *http.HeadRequest) *http.FutureHead {
	panic("Head")
}

func (m *myService) Options(ctx context.Context, req *http.OptionsRequest) *http.FutureOptions {
	panic("Options")
}

func (m *myService) Patch(ctx context.Context, req *http.PatchRequest) *http.FuturePatch {
	panic("Patch")
}

func (m *myService) Connect(ctx context.Context, req *http.ConnectRequest) *http.FutureConnect {
	panic("Connect")
}

func (m *myService) Trace(ctx context.Context, req *http.TraceRequest) *http.FutureTrace {
	panic("Trace")
}
