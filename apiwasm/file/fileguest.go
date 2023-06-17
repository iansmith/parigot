package main

import (
	"context"
	"unsafe"

	"github.com/iansmith/parigot/apishared/id"
	"github.com/iansmith/parigot/apiwasm"
	pcontext "github.com/iansmith/parigot/context"
	file "github.com/iansmith/parigot/g/file/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
)

var _ = unsafe.Sizeof([]byte{})

func main() {
	GoFile()
}

func GoFile() {
	ctx := pcontext.GuestContext(pcontext.NewContextWithContainer(context.Background(), "[filewasm]main"))
	myId := file.MustRegisterFile(ctx)
	file.MustExportFile(ctx)

	svc := &myFileSvc{}
	allDead := apiwasm.NewParigotWaitGroup("[main]File")
	file.MustWaitSatisfiedFile(ctx, myId, svc, allDead)
	kerr := file.StartFile(ctx, myId, svc)
	if kerr != syscall.KernelErr_NoError {
		pcontext.Errorf(ctx, "unable to start File: %s", syscall.KernelErr_name[int32(kerr)])
	}
	allDead.Wait()
}

type myFileSvc struct{}

var myImpl *myFileSvc = &myFileSvc{}

func (f *myFileSvc) Ready(ctx context.Context, _ id.ServiceId) bool {
	pcontext.Debugf(ctx, "Ready reached in file service")
	return true
}

func (f *myFileSvc) Open(ctx context.Context, in *file.OpenRequest) (*file.OpenResponse, file.FileErr) {
	return file.OpenHost(in)
}

func (f *myFileSvc) Create(ctx context.Context, in *file.CreateRequest) (*file.CreateResponse, file.FileErr) {
	return file.CreateHost(in)

}

func (f *myFileSvc) Close(ctx context.Context, in *file.CloseRequest) (*file.CloseResponse, file.FileErr) {
	return file.CloseHost(in)
}

func (f *myFileSvc) LoadTestData(ctx context.Context, in *file.LoadTestDataRequest) (*file.LoadTestDataResponse, file.FileErr) {
	return file.LoadTestDataHost(in)

}
