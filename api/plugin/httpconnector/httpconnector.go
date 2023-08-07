package httpconnector

import (
	"context"
	"io"
	"log"
	"net/http"

	apiplugin "github.com/iansmith/parigot/api/plugin"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/eng"
	"github.com/iansmith/parigot/g/httpconnector/v1"
	"github.com/tetratelabs/wazero/api"
	"google.golang.org/protobuf/proto"
)

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

	pcontext.Infof(ctx, "httpconnector plugin initialized, Init() in host")

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

func runHttpListener(hCnt *httpConnectorImpl) {
	h1 := func(w http.ResponseWriter, req *http.Request) {
		if req != nil {
			hCnt.httpChan <- req
		}
		io.WriteString(w, "listening!")
	}
	http.HandleFunc("/", h1)

	http.ListenAndServe(":9000", nil)
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

func (hCnt *httpConnectorImpl) check(ctx context.Context, req *httpconnector.CheckRequest, resp *httpconnector.CheckResponse) int32 {
	resp.Method = "test"

	pcontext.Infof(ctx, "http connector check is called %+v", req)

	if hCnt == nil {
		log.Println("http connector is nil")
		return int32(httpconnector.HttpConnectorErr_InternalError)
	}

	if hCnt.httpChan == nil {
		log.Println("httpChan is nil")
		return int32(httpconnector.HttpConnectorErr_InternalError)
	}
	select {
	case <-hCnt.httpChan:
		log.Println("req is ready")
	// case <-time.After(10 * time.Second):
	// 	log.Println("http connector time out")
	default:
		log.Println("the channel is not ready to receive the data")
	}
	return int32(httpconnector.HttpConnectorErr_NoError)
}
