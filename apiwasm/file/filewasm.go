package main

import (
	"context"
	"unsafe"

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
	s := file.NewSimpleFileService(ready)
	file.RunFileService(s)
}

func ready(ctx context.Context, _ *file.SimpleFileService) bool {
	file.WaitFileServiceOrPanic()
	return true
}
