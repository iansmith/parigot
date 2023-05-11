package go_

import (
	"log"

	"github.com/iansmith/parigot/eng"
)

func ParigotInit(e eng.Engine, inst eng.Instance) {
	e.AddSupportedFunc("queue", "create_queue", wrapFunc(inst, createQueue))
	e.AddSupportedFunc("queue", "delete_queue", wrapFunc(inst, deleteQueue))
	e.AddSupportedFunc("queue", "length", wrapFunc(inst, length))
	e.AddSupportedFunc("queue", "locate", wrapFunc(inst, locate))
	e.AddSupportedFunc("queue", "mark_done", wrapFunc(inst, markDone))
	e.AddSupportedFunc("queue", "receive", wrapFunc(inst, receive))
	e.AddSupportedFunc("queue", "send", wrapFunc(inst, send))
}

func wrapFunc(i eng.Instance, fn func(eng.Instance, int32) int32) func(int32) int32 {
	return func(x int32) int32 {
		return fn(i, x)
	}
}

func createQueue(inst eng.Instance, ptr int32) int32 {
	log.Printf("queue.createQueue 0x%x", ptr)
	return 0
}

func deleteQueue(inst eng.Instance, ptr int32) int32 {
	log.Printf("queue.deleteQueue 0x%x", ptr)
	return 0
}
func length(inst eng.Instance, ptr int32) int32 {
	log.Printf("queue.length 0x%x", ptr)
	return 0
}
func locate(inst eng.Instance, ptr int32) int32 {
	log.Printf("queue.locate 0x%x", ptr)
	return 0
}
func markDone(inst eng.Instance, ptr int32) int32 {
	log.Printf("queue.markDone 0x%x", ptr)
	return 0
}
func send(inst eng.Instance, ptr int32) int32 {
	log.Printf("queue.send 0x%x", ptr)
	return 0
}
func receive(inst eng.Instance, ptr int32) int32 {
	log.Printf("queue.receive 0x%x", ptr)
	return 0
}
