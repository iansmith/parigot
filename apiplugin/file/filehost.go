package main

import (
	"context"
	"unsafe"

	"github.com/iansmith/parigot/apiplugin"
	"github.com/iansmith/parigot/apishared/id"
	"github.com/iansmith/parigot/eng"
	"github.com/iansmith/parigot/g/file/v1"
	filemsg "github.com/iansmith/parigot/g/msg/file/v1"
	"github.com/iansmith/parigot/sys"

	"github.com/tetratelabs/wazero/api"
)

type filePlugin struct{}

var _ = unsafe.Sizeof([]byte{})

var ParigiotInitialize sys.ParigotInit = &filePlugin{}

func (*filePlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "file", "open_", open)
	e.AddSupportedFunc(ctx, "file", "load_test_data_", loadTestData)
	return true
}

// true native implementation of open
func openImpl(ctx context.Context, in *filemsg.OpenRequest, out *filemsg.OpenResponse) id.IdRaw {
	return file.FileErrIdNoErr.Raw()
}

// true native implementation of load test data
func loadTestDataImpl(ctx context.Context, in *filemsg.LoadTestDataRequest, out *filemsg.LoadTestDataResponse) id.IdRaw {
	return file.FileErrIdNoErr.Raw()
}

func open(ctx context.Context, m api.Module, stack []uint64) {
	req := &filemsg.OpenRequest{}
	resp := &filemsg.OpenResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[file]open", m, stack, openImpl, req, resp)
	return

}

func loadTestData(ctx context.Context, m api.Module, stack []uint64) {
	// xxxx should be pointing at the plugin code for load test data ,not open
	req := &filemsg.LoadTestDataRequest{}
	resp := &filemsg.LoadTestDataResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[file]loadTestData", m, stack, loadTestDataImpl, req, resp)
	return
}
