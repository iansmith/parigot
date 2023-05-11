package main

import (
	"context"

	"github.com/iansmith/parigot/g/queue/v1"
	lib "github.com/iansmith/parigot/lib/go"
)

//go:export parigot_main
//go:linkname parigot_main
func parigot_main() {
	lib.FlagParseCreateEnv()

	queue.ExportQueueServiceOrPanic()
	s := queue.NewSimpleQueueService(ready)
	queue.RunQueueService(s)
}

func ready(ctx context.Context, _ *queue.SimpleQueueService) bool {
	queue.WaitQueueServiceOrPanic()
	return true
}
