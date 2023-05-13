package main

import (
	"context"
	"log"

	"github.com/iansmith/parigot/eng"
	"github.com/iansmith/parigot/sys"
)

type queuePlugin struct{}

var ParigiotInitialize sys.ParigotInit = &queuePlugin{}

func (*queuePlugin) Init(ctx context.Context, e eng.Engine, inst eng.Instance) bool {
	e.AddSupportedFunc(ctx, "queue", "create_queue", wrapFunc(inst, createQueueHost))
	e.AddSupportedFunc(ctx, "queue", "delete_queue", wrapFunc(inst, deleteQueueHost))
	e.AddSupportedFunc(ctx, "queue", "length", wrapFunc(inst, lengthHost))
	e.AddSupportedFunc(ctx, "queue", "locate", wrapFunc(inst, locateHost))
	e.AddSupportedFunc(ctx, "queue", "mark_done", wrapFunc(inst, markDoneHost))
	e.AddSupportedFunc(ctx, "queue", "receive", wrapFunc(inst, receiveHost))
	e.AddSupportedFunc(ctx, "queue", "send", wrapFunc(inst, sendHost))
	return true
}

func wrapFunc(i eng.Instance, fn func(eng.Instance, int32) int32) func(int32) int32 {
	return func(x int32) int32 {
		return fn(i, x)
	}
}

func createQueueHost(inst eng.Instance, ptr int32) int32 {
	log.Printf("queue.createQueue 0x%x", ptr)
	return 0
}

func deleteQueueHost(inst eng.Instance, ptr int32) int32 {
	log.Printf("queue.deleteQueue 0x%x", ptr)
	return 0
}
func lengthHost(inst eng.Instance, ptr int32) int32 {
	log.Printf("queue.length 0x%x", ptr)
	return 0
}
func locateHost(inst eng.Instance, ptr int32) int32 {
	log.Printf("queue.locate 0x%x", ptr)
	return 0
}
func markDoneHost(inst eng.Instance, ptr int32) int32 {
	log.Printf("queue.markDone 0x%x", ptr)
	return 0
}
func sendHost(inst eng.Instance, ptr int32) int32 {
	log.Printf("queue.send 0x%x", ptr)
	return 0
}
func receiveHost(inst eng.Instance, ptr int32) int32 {
	log.Printf("queue.receive 0x%x", ptr)
	return 0
}
