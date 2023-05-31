//go:build wasip1

package main

import (
	"context"
	"unsafe"

	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/file/v1"
	gfile "github.com/iansmith/parigot/g/file/v1"
	filemsg "github.com/iansmith/parigot/g/msg/file/v1"
)

var _ = unsafe.Sizeof([]byte{})

func main() {
	ctx := pcontext.GuestContext(pcontext.NewContextWithContainer(context.Background(), "[filewasm]main"))
	gfile.MustRegisterFileService(ctx)
	gfile.MustExportFileService(ctx)
	gfile.RunFileService(ctx, &myFileSvc{})
}

type myFileSvc struct{}

var myImpl *myFileSvc = &myFileSvc{}

func (f *myFileSvc) Ready(ctx context.Context) bool {
	pcontext.Debugf(ctx, "Ready reached in file service")
	return true
}

func (f *myFileSvc) Open(ctx context.Context, in *filemsg.OpenRequest) (*filemsg.OpenResponse, file.FileErrId) {
	return gfile.OpenHost(in)
}

func (f *myFileSvc) Create(ctx context.Context, in *filemsg.CreateRequest) (*filemsg.CreateResponse, file.FileErrId) {
	return gfile.CreateHost(in)

}

func (f *myFileSvc) Close(ctx context.Context, in *filemsg.CloseRequest) (*filemsg.CloseResponse, file.FileErrId) {
	return gfile.CloseHost(in)
}

func (f *myFileSvc) LoadTestData(ctx context.Context, in *filemsg.LoadTestDataRequest) (*filemsg.LoadTestDataResponse, file.FileErrId) {
	return gfile.LoadTestDataHost(in)

}
