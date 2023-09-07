package httpconnector

import (
	"bytes"
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
	"github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/sys/kernel"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const port = ":9000"

const maxBodySize = 0xfffff // about a meg

var logger *slog.Logger

type HttpConnectorPlugin struct {
}

//var _ = kernel.GeneralReceiver(&httpConnectorImpl{})

type httpConnectorImpl struct {
	//kernelCh   chan *syscall.ReadOneResponse
	resultChan chan *httpconnector.HandleResponse

	host id.HostId
}

// type ParigotRequestWrapper interface {
// 	GetRequest() *phttp.HttpRequest
// }

func (h *HttpConnectorPlugin) Init(ctx context.Context, e eng.Engine, host id.HostId) bool {
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})).With("plugin", "httpconnector")

	// create the connector
	connector = newHtttpConnectorImpl(host)

	// run the loop that is listening on the http port (9000)
	go connector.runHttpListener()

	//kernel.K.AddReceiver(connector)

	return true
}

var connector *httpConnectorImpl

func newHtttpConnectorImpl(host id.HostId) *httpConnectorImpl {

	c := &httpConnectorImpl{
		//kernelCh:   make(chan *syscall.ReadOneResponse, bufferSizeOnChan),
		resultChan: make(chan *httpconnector.HandleResponse),
		host:       host,
	}

	return c
}

// runHttpListener starts an HTTP listener on port 9000 and listens for incoming requests from outside.
// When a request is received, it is sent to the httpChan channel of the provided httpConnectorImpl instance.
func (c *httpConnectorImpl) runHttpListener() {
	h1 := func(w http.ResponseWriter, req *http.Request) {
		if req != nil {
			ok := connector.callHostDispatch(req, w)
			if ok {
				ourResp := <-c.resultChan
				if ourResp != nil {
					c.writeResponse(ourResp, w, req)
				}
			}
		}
	}
	http.HandleFunc("/", h1)
	http.ListenAndServe(port, nil)
}

func (h *httpConnectorImpl) writeResponse(result *httpconnector.HandleResponse, writer http.ResponseWriter, req *http.Request) {
	for k, v := range result.Header {
		writer.Header().Set(k, v)
	}
	raw := result.GetHttpResponse()
	writer.WriteHeader(int(result.GetHttpStatus()))
	rd := bytes.NewBuffer(raw)
	l, err := io.Copy(writer, rd)
	if err != nil {
		logger.Error("io.Copy failed", "error", err)
		return
	}
	if err != nil {
		logger.Error("unable to write bytes of response", "error", err)
		return
	}
	logger.Info("response written", "status", result.GetHttpStatus(), "result bytes", l, "method", req.Method, "url", req.URL.Path)
}

func (h *httpConnectorImpl) callHostDispatch(req *http.Request, httpresp http.ResponseWriter) bool {

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
		return false
	}
	if err != httpconnector.HttpConnectorErr_NoError {
		logger.Error("swallowing httpconnector error", "httpconnector error", httpconnector.HttpConnectorErr_name[int32(err)])
		return false
	}
	hid, sid, mid, err := findHttpConnector()
	if err != httpconnector.HttpConnectorErr_NoError {
		return false
	}

	handleReq := &httpconnector.HandleRequest{
		Url:        req.RequestURI,
		HttpMethod: req.Method,
		ServiceId:  sid.Marshal(),
		MethodId:   mid.Marshal(),
		ReqAny:     parigotReq,
	}

	if hid.IsZeroOrEmptyValue() || sid.IsZeroOrEmptyValue() || mid.IsZeroOrEmptyValue() {
		logger.Error("failed to get id set back corectly from MethodDetail", "host", hid,
			"service", sid, "method", mid)
		return false
	}

	handleAny := &anypb.Any{}
	if err := handleAny.MarshalFrom(handleReq); err != nil {
		logger.Error("unable to marshal handle req into any", "error", err)
		return false
	}

	bundle := &syscall.MethodBundle{
		HostId:    hid.Marshal(), // why is this needed?
		ServiceId: sid.Marshal(),
		MethodId:  mid.Marshal(),
		CallId:    id.NewCallId().Marshal(),
	}

	dispReq := &syscall.DispatchRequest{
		Bundle: bundle,
		Param:  handleAny,
	}
	dispResp := &syscall.DispatchResponse{}

	wrapper := &wrapper{httpresp}
	kerr := kernel.K.HostDispatch(dispReq, dispResp, func(rc *syscall.ResolvedCall) {
		h.callbackFromHandle(httpresp, rc)
	}, wrapper)
	if kerr != syscall.KernelErr_NoError {
		logger.Error("unable to dispatch Handle() message in httpconnector", "kernel error", kerr)
		return false
	}
	return true
}

type wrapper struct {
	inner io.Writer
}

func (w *wrapper) Write(data []byte) (int, error) {
	slog.Warn("somebody is messing with the writer", "string", string(data), "size", len(data))
	return w.inner.Write(data)
}

func (h *httpConnectorImpl) callbackFromHandle(writer http.ResponseWriter, rc *syscall.ResolvedCall) {

	//xxx this should probably be doing bookkeeping so that we can respond to
	//xxx http messages in the order received. this would require keeping a list
	//xxx of outstanding requests and their associated call ids so if
	//xxx we get an out of order response, we should just hold it until
	//xxx it reaches the front of the list.

	if rc.GetResult() == nil {
		// better have an error cause all we can do is show it
		logger.Error("failed to get response from Handle",
			"error", httpconnector.HttpConnectorErr_name[rc.GetResultError()])
		h.resultChan <- nil
		return
	}

	result := &httpconnector.HandleResponse{}
	if e := rc.GetResult().UnmarshalTo(result); e != nil {
		logger.Error("unmarshal failed getting response from Handle",
			"error", e.Error())
		h.resultChan <- nil
		return
	}

	h.resultChan <- result
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
	return 50
}

// func (h *httpConnectorImpl) Ch() chan *syscall.ReadOneResponse {
// 	return h.kernelCh
// }

func (h *httpConnectorImpl) HostId() id.HostId {
	return h.host
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
		ServiceName: "http_connector", // this is the name in WASM terms
		CalledBy:    nil,
	}
	locResp := &syscall.LocateResponse{}
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
	return hid, sid, mid, httpconnector.HttpConnectorErr_NoError
}
