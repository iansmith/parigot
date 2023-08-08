package httpconnector

import (
	"context"
	"io"
	"net/http"
	"time"

	apiplugin "github.com/iansmith/parigot/api/plugin"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/eng"
	"github.com/iansmith/parigot/g/httpconnector/v1"
	"github.com/tetratelabs/wazero/api"
	"google.golang.org/protobuf/proto"
)

const port = ":9000"

type HttpConnectorPlugin struct{}

var httpCnt *httpConnectorImpl

type httpConnectorImpl struct {
	ctx      context.Context
	httpChan chan *http.Request
}

func (*HttpConnectorPlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "httpconnector", "check_", checkHost)

	httpCnt = newHttpCntSvc(ctx)
	go runHttpListener(httpCnt)

	return true
}

func newHttpCntSvc(ctx context.Context) *httpConnectorImpl {
	newCtx := pcontext.ServerGoContext(ctx)

	c := &httpConnectorImpl{
		ctx:      newCtx,
		httpChan: make(chan *http.Request),
	}

	return c
}

// runHttpListener starts an HTTP listener on port 9000 and listens for incoming requests from outside.
// When a request is received, it is sent to the httpChan channel of the provided httpConnectorImpl instance.
func runHttpListener(hCnt *httpConnectorImpl) {
	h1 := func(w http.ResponseWriter, req *http.Request) {
		if req != nil {
			hCnt.httpChan <- req
		}
		io.WriteString(w, "start running HTTP listener!")
	}
	http.HandleFunc("/", h1)

	http.ListenAndServe(port, nil)
}

func hostBase[T proto.Message, U proto.Message](ctx context.Context, fnName string,
	fn func(context.Context, T, U) int32, m api.Module, stack []uint64, req T, resp U) {
	defer func() {
		if r := recover(); r != nil {
			print(">>>>>>>> Trapped recover in set up for   ", fnName, "<<<<<<<<<<\n")
		}
	}()
	apiplugin.InvokeImplFromStack(ctx, fnName, m, stack, fn, req, resp)
}

func checkHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &httpconnector.CheckRequest{}
	resp := &httpconnector.CheckResponse{}

	hostBase(ctx, "[httpconnector]check", httpCnt.check, m, stack, req, resp)
}

// Check waits for a response from the http channel.
// It sets the Success field of the CheckResponse to true if the http channel responds.
func (hCnt *httpConnectorImpl) check(ctx context.Context, req *httpconnector.CheckRequest, resp *httpconnector.CheckResponse) int32 {
	resp.Success = false
	// Check that the http connector and httpChannel are properly initialized.
	// They are always good, just be careful
	if hCnt == nil || hCnt.httpChan == nil {
		pcontext.Errorf(ctx, "HTTP connector is not properly initialized")
		return int32(httpconnector.HttpConnectorErr_InternalError)
	}

	select {
	case <-hCnt.httpChan:
		pcontext.Infof(ctx, "req is ready")
		resp.Success = true
	case <-time.After(time.Millisecond * 40):
		pcontext.Infof(ctx, "http connector time out")
	}
	return int32(httpconnector.HttpConnectorErr_NoError)
}
