package httpconnector

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

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

// var httpChan = make(chan *http.Request)

func (*HttpConnectorPlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "httpconnector", "check_", checkHost)

	newHttpCntSvc(ctx)
	go runHttpListener()

	return true
}

func newHttpCntSvc(ctx context.Context) *httpConnectorImpl {
	newCtx := pcontext.ServerGoContext(ctx)

	hCnt := &httpConnectorImpl{
		ctx:      newCtx,
		httpChan: make(chan *http.Request),
	}

	return hCnt
}

func runHttpListener() {
	h1 := func(w http.ResponseWriter, req *http.Request) {
		// hCnt.httpChan <- req
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
	select {
	case <-hCnt.httpChan:
		log.Println("req is")
	case <-time.After(10 * time.Second):
		log.Println("http connector time out")
	}
	return int32(httpconnector.HttpConnectorErr_NoError)
}
