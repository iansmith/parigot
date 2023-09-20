package main

import (
	"context"
	"log/slog"
	"strings"

	"github.com/iansmith/parigot/api/guest"
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/example/httpsimple/g/simple/v1"
	"github.com/iansmith/parigot/g/http/v1"
	"github.com/iansmith/parigot/g/httpconnector/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"google.golang.org/protobuf/types/known/anypb"

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
	myId id.ServiceId
}

var _ httpconnector.HttpConnector = &myService{}

func main() {

	// the implementation of the service (no state right now)
	impl := &myService{}

	// Init initiaizes a service and normally receives a list of functions
	// that indicate dependencies.  We depend only on simple.
	binding, fut, ctx, sid :=
		httpconnector.Init([]lib.MustRequireFunc{http.MustRequire}, impl)

	// create a logger
	logger = slog.New(guest.NewParigotHandler(sid))

	fut.Failure(func(err syscall.KernelErr) {
		logger.Error("failed to launch the frontdoor service", "kernel error", syscall.KernelErr_name[int32(err)])
		lib.ExitSelf(ctx, 1, sid)
	})

	fut.Success(func(resp *syscall.LaunchResponse) {
		logger.Info("frontdoor launched successfully",
			slog.String("frontdoor host", syscallguest.CurrentHostId().Short()))
	})

	// does not return except when the world is broken
	err := httpconnector.Run(ctx, binding, simple.TimeoutInMillis, nil)

	logger.Error("we got an error that is not normal from ReadOneAndCallClient", "kernel error", syscall.KernelErr_name[int32(err)])

	lib.ExitSelf(ctx, 1, sid)
}

var notGet = "not permitted, you jerk!\n"

func (m *myService) Handle(ctx context.Context, req *httpconnector.HandleRequest) *httpconnector.FutureHandle {

	if req.GetHttpMethod() != "GET" {
		fh := httpconnector.NewFutureHandle()
		result := &httpconnector.HandleResponse{
			HttpStatus:   401,
			HttpResponse: []byte(notGet),
			Header: map[string]string{
				"try-different-verb": "GET",
			},
		}
		a := &anypb.Any{}
		if e := a.MarshalFrom(result); e != nil {
			slog.Error("unable to marshal", "error", e)
			return nil
		}
		if e := fh.CompleteMethod(ctx, a, int32(httpconnector.HttpConnectorErr_NoError), syscallguest.CurrentHostId()); e != syscall.KernelErr_NoError {
			slog.Error("unable to complete Handle() method future", "kernel error", syscall.KernelErr_name[int32(e)])
			// not much we can do here, panic?
			return nil
		}
		return fh
	}

	httpClient, err := http.Locate(ctx, m.myId)
	if err == syscall.KernelErr_NotFound {
		badfut := httpconnector.NewFutureHandle()
		badfut.Method.CompleteMethod(ctx, nil, httpconnector.HttpConnectorErr_NoReceiver)
		return badfut
	}
	if err != syscall.KernelErr_NoError {
		slog.Error("unable to locate an http receiver", "kernel error", syscall.KernelErr_name[int32(err)])
		return nil
	}
	futGet := httpClient.Get(ctx, &http.GetRequest{
		Request: &http.HttpRequest{
			Url:     req.GetUrl(),
			Header:  map[string]string{},
			Body:    []byte{},
			Trailer: map[string]string{},
		},
	})

	futHandle := httpconnector.NewFutureHandle()

	futGet.Method.Success(func(resp *http.GetResponse) {
		futHandle.Method.CompleteMethod(ctx, &httpconnector.HandleResponse{
			HttpStatus:   resp.Response.StatusCode,
			Header:       resp.Response.Header,
			HttpResponse: resp.Response.Body,
		}, httpconnector.HttpConnectorErr_NoError)
	})
	futGet.Method.Failure(func(err http.HttpErr) {
		slog.Info("error received from GET receiver", "name", http.HttpErr_name[int32(err)])
		futHandle.Method.CompleteMethod(ctx, nil, httpconnector.HttpConnectorErr_ReceiverFailed)
	})
	if strings.HasSuffix(req.Url, "/exit") {
		logger.Info("exit requested")
		lib.ExitSelf(ctx, 1, m.myId)
	}
	return futHandle
}

func (m *myService) Ready(ctx context.Context, sid id.ServiceId) *future.Base[bool] {
	m.myId = sid
	return future.NewBaseWithValue[bool](true)
}
