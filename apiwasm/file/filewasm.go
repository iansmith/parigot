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

const (
	// NoError means just what it sounds like.  All Ids that are errors represent
	// no error as 0.
	NoError FileErrCode = 0
	// BadPath means that the given path (filename) is not valid.
	BadPath FileErrCode = 1
	// NotFound means that the given path could point to a file (it is valid)
	// but the path given could not be found on the current filesystem.
	NotFound FileErrCode = 2
)

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
