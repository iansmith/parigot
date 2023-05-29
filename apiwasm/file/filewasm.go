package main

import (
	"context"
	"unsafe"

	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/file/v1"
	"github.com/iansmith/parigot/g/queue/v1"
	lib "github.com/iansmith/parigot/lib/go"
)

var _ = unsafe.Sizeof([]byte{})

func main() {

}

//go:wasmexport parigot_main parigot_main
func parigot_main() {
	lib.FlagParseCreateEnv()

	queue.ExportQueueServiceOrPanic()
	ctx := pcontext.ClientContext(pcontext.NewContextWithContainer(context.Background(), "[filewasm]main"))
	s := file.NewSimpleFileService(ctx, ready)
	file.RunFileService(ctx, s)
}

func ready(ctx context.Context, _ *file.SimpleFileService) bool {
	file.WaitFileServiceOrPanic()
	return true
}
