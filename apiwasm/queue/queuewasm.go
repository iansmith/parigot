package main

import (
	"context"
	"unsafe"

	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/queue/v1"
	lib "github.com/iansmith/parigot/lib/go"
)

var _ = unsafe.Sizeof([]byte{})

func main() {

}

//go:export parigot_main
//go:linkname parigot_main
func parigot_main() {
	lib.FlagParseCreateEnv()

	ctx := pcontext.ClientContext(pcontext.NewContextWithContainer(context.Background(), "[queuewasm]main"))

	queue.ExportQueueServiceOrPanic()
	s := queue.NewSimpleQueueService(ctx, ready)
	queue.RunQueueService(ctx, s)
}

func ready(ctx context.Context, _ *queue.SimpleQueueService) bool {
	queue.WaitQueueServiceOrPanic()
	return true
}
