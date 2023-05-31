package main

import (
	"context"
	"unsafe"

	"github.com/iansmith/parigot/apiwasm"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/file/v1"
	filemsg "github.com/iansmith/parigot/g/msg/file/v1"
)

var _ = unsafe.Sizeof([]byte{})

func main() {
	ctx := pcontext.GuestContext(pcontext.NewContextWithContainer(context.Background(), "[filewasm]main"))
	file.MustRegisterFileService(ctx)
	s := file.NewWrapFileService(ctx, myImpl)
	file.RunFileService(ctx, s)
}

type myFileSvc struct{}

var myImpl *myFileSvc = &myFileSvc{}

func (f *myFileSvc) Ready(ctx context.Context) bool {
	pcontext.Debugf(ctx, "Ready reached in file service")
	return true
}

func (f *myFileSvc) Open(ctx context.Context, in *filemsg.OpenRequest) (*filemsg.OpenResponse, file.FileErrId) {
	outProtoPtr := &filemsg.OpenResponse{}
	defer pcontext.Dump(ctx)
	return apiwasm.ClientSide(ctx, inPtr, outProtoPtr, Run_)
}

func (f *myFileSvc) Create(ctx context.Context, in *filemsg.CreateRequest) (*filemsg.CreateResponse, file.FileErrId) {

}

func (f *myFileSvc) Close(ctx context.Context, in *filemsg.CloseRequest) (*filemsg.CloseResponse, file.FileErrId) {

}

func (f *myFileSvc) LoadTestData(ctx context.Context, in *filemsg.LoadTestDataRequest) (*filemsg.LoadTestDataResponse, file.FileErrId) {

}
