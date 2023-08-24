package httpconnector

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/eng"
	phttp "github.com/iansmith/parigot/g/http/v1"
	"github.com/iansmith/parigot/g/httpconnector/v1"
	"github.com/iansmith/parigot/g/protosupport/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/sys/kernel"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const bufferSizeOnChan = 8

const port = ":9000"

const timeoutInMillis = 40

const maxBodySize = 0xffff

var logger *slog.Logger

type HttpConnectorPlugin struct {
	conn *httpConnectorImpl
}

type httpConnectorImpl struct {
	httpChan chan *http.Request
	kernelCh chan *syscall.ReadOneResponse
}

type ParigotRequestWrapper interface {
	GetRequest() *phttp.HttpRequest
}

func (h *HttpConnectorPlugin) Init(ctx context.Context, e eng.Engine) bool {
	//????
	//e.AddSupportedFunc(ctx, "httpconnector", "check_", checkHost)

	// setup logger
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})).With("plugin", "httpconnector")

	// create the connector
	connector = newHtttpConnectorImpl()

	// run the loop that is listening on the http port (9000)
	go runHttpListener(connector)

	kernel.K.AddReceiver(connector)

	return true
}

var connector *httpConnectorImpl

func newHtttpConnectorImpl() *httpConnectorImpl {

	c := &httpConnectorImpl{
		httpChan: make(chan *http.Request, bufferSizeOnChan),
		kernelCh: make(chan *syscall.ReadOneResponse, bufferSizeOnChan),
	}

	return c
}

// runHttpListener starts an HTTP listener on port 9000 and listens for incoming requests from outside.
// When a request is received, it is sent to the httpChan channel of the provided httpConnectorImpl instance.
func runHttpListener(connector *httpConnectorImpl) {
	h1 := func(w http.ResponseWriter, req *http.Request) {
		if req != nil {
			logger.Info("got HTTP request", "method", req.Method)
			result := connector.makeReadOneResult(req)
			if result != nil {
				connector.kernelCh <- result
			}
		}
	}
	http.HandleFunc("/", h1)

	slog.Info("start the HTTP listerner")
	http.ListenAndServe(port, nil)
}

func (h *httpConnectorImpl) makeReadOneResult(req *http.Request) *syscall.ReadOneResponse {

	var parigotReq *anypb.Any
	var err httpconnector.HttpConnectorErr
	switch strings.ToUpper(req.Method) {
	case "GET":
		get := &phttp.GetRequest{Request: convertFieldsHttp(req)}
		parigotReq, err = convertHttpReqToParigot(get)
	case "POST":
		post := &phttp.PostRequest{Request: convertFieldsHttp(req)}
		parigotReq, err = convertHttpReqToParigot(post)
	case "PUT":
		put := &phttp.PutRequest{Request: convertFieldsHttp(req)}
		parigotReq, err = convertHttpReqToParigot(put)
	case "DELETE":
		delete := &phttp.DeleteRequest{Request: convertFieldsHttp(req)}
		parigotReq, err = convertHttpReqToParigot(delete)
	case "OPTIONS":
		options := &phttp.OptionsRequest{Request: convertFieldsHttp(req)}
		parigotReq, err = convertHttpReqToParigot(options)
	case "PATCH":
		patch := &phttp.PatchRequest{Request: convertFieldsHttp(req)}
		parigotReq, err = convertHttpReqToParigot(patch)
	case "CONNECT":
		conn := &phttp.ConnectRequest{Request: convertFieldsHttp(req)}
		parigotReq, err = convertHttpReqToParigot(conn)
	case "TRACE":
		trace := &phttp.TraceRequest{Request: convertFieldsHttp(req)}
		parigotReq, err = convertHttpReqToParigot(trace)
	default:
		logger.Warn("knknown HTTP method type, ignoring", "method", req.Method)
		return nil
	}
	if err != httpconnector.HttpConnectorErr_NoError {
		logger.Error("swallowing httpconnector error", "httpconnector error", httpconnector.HttpConnectorErr_name[int32(err)])
		return nil
	}

	handleReq := &httpconnector.HandleRequest{
		HttpMethod: req.Method,
		ServiceId:  &protosupport.IdRaw{},
		MethodId:   &protosupport.IdRaw{},
		ReqAny:     parigotReq,
	}

	handleAny := &anypb.Any{}
	if err := handleAny.MarshalFrom(handleReq); err != nil {
		logger.Error("unable to marshal handle req into any", "error", err)
	}

	hid, sid, mid, kerr := kernel.K.Nameserver().MethodDetail(
		kernel.FQName{Pkg: "httpconnector.v1", Name: "httpconnector"}, "Handle")

	if kerr != syscall.KernelErr_NoError {
		logger.Error("failed to get details about connector handle method", "kernel error", kerr)
		return nil
	}
	if hid.IsZeroOrEmptyValue() || sid.IsZeroOrEmptyValue() || mid.IsZeroOrEmptyValue() {
		logger.Error("failed to get id set back corectly from MethodDetail", "host", hid,
			"service", sid, "method", mid)
		return nil
	}

	resp := &syscall.ReadOneResponse{
		Timeout: false,
		Bundle: &syscall.MethodBundle{
			HostId:    hid.Marshal(),
			ServiceId: sid.Marshal(),
			MethodId:  mid.Marshal(),
			CallId:    id.NewCallId().Marshal(),
		},
		ParamOrResult: handleAny,
		ResultErr:     0,
		Resolved:      nil,
		Exit:          nil,
	}
	return resp
}

func convertHttpReqToParigot[T proto.Message](raw T) (*anypb.Any, httpconnector.HttpConnectorErr) {
	a := &anypb.Any{}
	if err := a.MarshalFrom(raw); err != nil {
		logger.Error("unable to marshal http req to parigot http req", "error", err)
		return nil, httpconnector.HttpConnectorErr_MarshalError
	}
	return a, httpconnector.HttpConnectorErr_NoError
}

func (h *httpConnectorImpl) TimeoutInMillis() int {
	return 20
}

// func checkHost(ctx context.Context, m api.Module, stack []uint64) {
// 	req := &httpconnector.CheckRequest{}
// 	resp := &httpconnector.CheckResponse{}

// 	hostBase(ctx, "[httpconnector]check", connector.check, m, stack, req, resp)
// }

// func handleHost(ctx context.Context, m api.Module, stack []uint64) {
// 	req := &httpconnector.HandleRequest{}
// 	resp := &httpconnector.HandleResponse{}

// 	hostBase(ctx, "[httpconnector]handle", connector.handleImpl, m, stack, req, resp)
// }

func (h *httpConnectorImpl) Ch() chan *syscall.ReadOneResponse {
	return h.kernelCh
}

// Check waits for a response from the http channel.
// It sets the Success field of the CheckResponse to true if the http channel responds.
// func (hCnt *httpConnectorImpl) check(ctx context.Context, req *httpconnector.CheckRequest, resp *httpconnector.CheckResponse) int32 {
// 	resp.Success = false
// 	// Check that the http connector and httpChannel are properly initialized.
// 	// They are always good, just be careful
// 	if hCnt == nil || hCnt.httpChan == nil {
// 		logger.Error("HTTP connector is not properly initialized")
// 		return int32(httpconnector.HttpConnectorErr_InternalError)
// 	}

// 	select {
// 	case <-hCnt.httpChan:
// 		logger.Info("req is ready")
// 		resp.Success = true
// 	case <-time.After(time.Millisecond * timeoutInMillis):
// 		logger.Info("http connector time out")
// 	}
// 	return int32(httpconnector.HttpConnectorErr_NoError)
// }

// func (c *httpConnectorImpl) handleImpl(ctx context.Context, req *httpconnector.HandleRequest, resp *httpconnector.HandleResponse) int32 {

// }

func convertFieldsHttp(req *http.Request) *phttp.HttpRequest {
	result := &phttp.HttpRequest{}
	result.Url = req.URL.String()
	hdr := req.Header
	for key, values := range hdr {
		value := values[0] // ignore dups
		result.Header[key] = value
	}
	trailer := req.Trailer
	for key, values := range trailer {
		value := values[0] // ignore dups
		result.Trailer[key] = value
	}
	body, err := req.GetBody()
	if err != nil {
		logger.Error("unable to get body from HTTP request", "error", err)
		return nil
	}
	rd := io.LimitReader(body, maxBodySize)
	buf, err := io.ReadAll(rd)
	if err != nil {
		logger.Error("unable to read body from HTTP request", "error", err)
		return nil
	}
	result.Body = buf
	return result
}
