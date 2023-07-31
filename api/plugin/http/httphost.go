package http

import (
	"context"

	apiplugin "github.com/iansmith/parigot/api/plugin"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/eng"
	"github.com/iansmith/parigot/g/http/v1"
	"github.com/tetratelabs/wazero/api"
	"google.golang.org/protobuf/proto"
)

type HttpPlugin struct{}

type httpSvcImpl struct {
	ctx context.Context
}

var httpSvc *httpSvcImpl

func (*HttpPlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "http", "get_", getHost)

	newHttpSvc(ctx)
	return true
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

func getHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &http.GetRequest{}
	resp := &http.GetResponse{}

	hostBase(ctx, "[http]get", httpSvc.get, m, stack, req, resp)
}

func newHttpSvc(ctx context.Context) *httpSvcImpl {
	newCtx := pcontext.ServerGoContext(ctx)

	h := &httpSvcImpl{
		ctx: newCtx,
	}

	return h
}

func (h *httpSvcImpl) get(ctx context.Context, req *http.GetRequest, resp *http.GetResponse) int32 {
	return int32(http.HttpErr_NoError)
}
