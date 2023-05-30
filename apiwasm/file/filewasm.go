package main

import (
	"context"
	"unsafe"

	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/file/v1"
)

var _ = unsafe.Sizeof([]byte{})

func main() {
	ctx := pcontext.GuestContext(pcontext.NewContextWithContainer(context.Background(), "[filewasm]main"))
	file.MustRegisterFileService(ctx)
	s := file.NewSimpleFileService(ctx, ready)
	file.RunFileService(ctx, s)
}

func ready(ctx context.Context, _ *file.SimpleFileService) bool {
	file.WaitFileServiceOrPanic()
	return true
}
