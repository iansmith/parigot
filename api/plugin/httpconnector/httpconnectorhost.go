package httpconnector

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/eng"
	phttp "github.com/iansmith/parigot/g/http/v1"
	"github.com/iansmith/parigot/g/httpconnector/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/sys/kernel"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const bufferSizeOnChan = 8

const port = ":9000"

const maxBodySize = 0xfffff // about a meg

var logger *slog.Logger

type HttpConnectorPlugin struct {
}

type httpConnectorImpl struct {
	httpChan chan *http.Request
	kernelCh chan *syscall.ReadOneResponse
}

type ParigotRequestWrapper interface {
	GetRequest() *phttp.HttpRequest
}

var handleMethod id.MethodId

func (h *HttpConnectorPlugin) Init(ctx context.Context, e eng.Engine) bool {
	//????
	//e.AddSupportedFunc(ctx, "httpconnector", "handle_", handleHost)

	// setup logger
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})).With("plugin", "httpconnector")
	logger.Info("about to create connector")

	// create the connector
	connector = newHtttpConnectorImpl()

	// run the loop that is listening on the http port (9000)
	logger.Info("about to run connector")
	go runHttpListener(connector)

	kernel.K.AddReceiver(connector)

	logger.Info("httpconnecter Init() complete")

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
	logger.Info("setting up the handler")
	http.HandleFunc("/", h1)

	slog.Info("start the HTTP listerner")
	http.ListenAndServe(port, nil)
}

func (h *httpConnectorImpl) makeReadOneResult(req *http.Request) *syscall.ReadOneResponse {

	var parigotReq *anypb.Any
	var err httpconnector.HttpConnectorErr
	log.Printf("dispatching a call to %s", req.Method)
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
	hid, sid, mid, err := findHttpConnector()
	logger.Info("finished finding HTTP connector")
	if err != httpconnector.HttpConnectorErr_NoError {
		return nil
	}
	handleReq := &httpconnector.HandleRequest{
		HttpMethod: req.Method,
		ServiceId:  sid.Marshal(),
		MethodId:   mid.Marshal(),
		ReqAny:     parigotReq,
	}
	slog.Info("created handle req", "method", mid.Short())

	if hid.IsZeroOrEmptyValue() || sid.IsZeroOrEmptyValue() || mid.IsZeroOrEmptyValue() {
		logger.Error("failed to get id set back corectly from MethodDetail", "host", hid,
			"service", sid, "method", mid)
		return nil
	}

	handleAny := &anypb.Any{}
	if err := handleAny.MarshalFrom(handleReq); err != nil {
		logger.Error("unable to marshal handle req into any", "error", err)
	}
	logger.Info("about to hit the response to the client side")
	bundle := &syscall.MethodBundle{
		HostId:    hid.Marshal(), // why is this needed?
		ServiceId: sid.Marshal(),
		MethodId:  mid.Marshal(),
		CallId:    id.NewCallId().Marshal(),
	}

	// do we really need this dispatch?
	dispReq := &syscall.DispatchRequest{
		Bundle: bundle,
		Param:  handleAny,
	}
	dispResp := &syscall.DispatchResponse{}
	logger.Info("sending dispatch!")
	kerr := kernel.K.Dispatch(dispReq, dispResp)
	if kerr != syscall.KernelErr_NoError {
		logger.Error("unable to dispatch Handle() message in httpconnector", "kernel error", kerr)
	}
	log.Printf("finished the send of the message")

	rd := &syscall.ReadOneResponse{
		Timeout:       false,
		Bundle:        bundle,
		ParamOrResult: handleAny,
		ResultErr:     0,
		Resolved:      nil,
		Exit:          nil,
	}
	return rd
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

func (h *httpConnectorImpl) Ch() chan *syscall.ReadOneResponse {
	return h.kernelCh
}

func convertFieldsHttp(req *http.Request) *phttp.HttpRequest {
	result := &phttp.HttpRequest{}
	result.Url = req.URL.String()
	result.Header = make(map[string]string)
	for key, values := range req.Header {
		value := values[0] // ignore dups
		result.Header[key] = value
	}
	result.Trailer = make(map[string]string)
	for key, values := range req.Trailer {
		value := values[0] // ignore dups
		result.Trailer[key] = value
	}
	body := req.Body
	if body != nil {
		rd := io.LimitReader(body, maxBodySize)
		buf, err := io.ReadAll(rd)
		if err != nil {
			logger.Error("unable to read body from HTTP request", "error", err)
			return nil
		}
		result.Body = buf
		body.Close()
	}
	return result
}

// findHttpConnector attempts find some service that implements httpconnector
func findHttpConnector() (id.HostId, id.ServiceId, id.MethodId, httpconnector.HttpConnectorErr) {
	locReq := &syscall.LocateRequest{
		PackageName: "httpconnector.v1",
		ServiceName: "httpconnector",
		CalledBy:    nil,
	}
	locResp := &syscall.LocateResponse{}
	log.Printf("xxx -- about do real locate")
	kerr := kernel.K.Locate(locReq, locResp)
	if kerr != syscall.KernelErr_NoError {
		logger.Error("unable to find service with (internal) Locate",
			"package", locReq.GetPackageName(), "service", locReq.GetServiceName())
		return id.HostIdZeroValue(), id.ServiceIdZeroValue(), id.MethodIdZeroValue(), httpconnector.HttpConnectorErr_NoReceiver
	}
	hid := id.UnmarshalHostId(locResp.GetHostId())
	sid := id.UnmarshalServiceId(locResp.GetServiceId())
	methMap := locResp.GetBinding()
	mid := id.MethodIdZeroValue()
	for _, binding := range methMap {
		if binding.GetMethodName() == "Handle" {
			mid = id.UnmarshalMethodId(binding.MethodId)
			break
		}
	}
	if mid.IsZeroValue() {
		logger.Error("unable to find method 'Handle' for service",
			"package", locReq.GetPackageName(), "service", locReq.GetServiceName())
		return id.HostIdZeroValue(), id.ServiceIdZeroValue(), id.MethodIdZeroValue(), httpconnector.HttpConnectorErr_NoReceiver
	}
	logger.Info("xxx -- received locate", "host", hid.Short(), "service", sid.Short(),
		"method", mid.Short())
	return hid, sid, mid, httpconnector.HttpConnectorErr_NoError
}
