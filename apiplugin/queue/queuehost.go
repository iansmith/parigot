package main

import (
	"context"
	"log"

	"github.com/iansmith/parigot/eng"

	"github.com/tetratelabs/wazero/api"
)

type queuePlugin struct{}

var ParigotInitialize = queuePlugin{}

func (*queuePlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "queue", "create_queue_", createQueueHost)
	e.AddSupportedFunc(ctx, "queue", "delete_queue_", deleteQueueHost)
	e.AddSupportedFunc(ctx, "queue", "length_", lengthHost)
	e.AddSupportedFunc(ctx, "queue", "locate_", locateHost)
	e.AddSupportedFunc(ctx, "queue", "mark_done_", markDoneHost)
	e.AddSupportedFunc(ctx, "queue", "receive_", receiveHost)
	e.AddSupportedFunc(ctx, "queue", "send_", sendHost)
	return true
}

func createQueueHost(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("queue.createQueue 0x%x", stack)
}

func deleteQueueHost(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("queue.deleteQueue 0x%x", stack)
}
func lengthHost(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("queue.length 0x%x", stack)
}
func locateHost(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("queue.locate 0x%x", stack)
}
func markDoneHost(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("queue.markDone 0x%x", stack)
}
func sendHost(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("queue.send 0x%x", stack)
}
func receiveHost(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("queue.receive 0x%x", stack)
}
