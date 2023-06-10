package main

import (
	"context"
	"unsafe"

	pcontext "github.com/iansmith/parigot/context"
	file "github.com/iansmith/parigot/g/file/v1"
)

var _ = unsafe.Sizeof([]byte{})

func main() {
	ctx := pcontext.GuestContext(pcontext.NewContextWithContainer(context.Background(), "[filewasm]main"))
	file.MustRegisterFileService(ctx)
	file.MustExportFileService(ctx)
	file.RunFileService(ctx, &myFileSvc{})
}

type myFileSvc struct{}

var myImpl *myFileSvc = &myFileSvc{}

func (f *myFileSvc) Ready(ctx context.Context) bool {
	pcontext.Debugf(ctx, "Ready reached in file service")
	return true
}

func (f *myFileSvc) Open(ctx context.Context, in *file.OpenRequest) (*file.OpenResponse, file.FileErrId) {
	return file.OpenHost(in)
}

func (f *myFileSvc) Create(ctx context.Context, in *file.CreateRequest) (*file.CreateResponse, file.FileErrId) {
	return file.CreateHost(in)

}

func (f *myFileSvc) Close(ctx context.Context, in *file.CloseRequest) (*file.CloseResponse, file.FileErrId) {
	return file.CloseHost(in)
}

func (f *myFileSvc) LoadTestData(ctx context.Context, in *file.LoadTestDataRequest) (*file.LoadTestDataResponse, file.FileErrId) {
	return file.LoadTestDataHost(in)

}
