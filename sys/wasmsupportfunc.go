package sys

import (
	filego "github.com/iansmith/parigot/apigo/file/go_"
	queuego "github.com/iansmith/parigot/apigo/queue/go_"
)

type WasmtimeSupportFunc struct {
	rt    *Runtime
	file  *filego.FileSvcImpl
	queue *queuego.QueueSvcImpl
}

func NewWasmtimeSupportFunc(ctx *DeployContext) *WasmtimeSupportFunc {
	q, errId, errInfo := queuego.NewQueueSvc()
	if errId != nil || errInfo != "" {
		panic("unable to create queue service " + errInfo)
	}
	return &WasmtimeSupportFunc{
		rt:    newRuntime(ctx),
		file:  filego.NewFileSvcImpl(),
		queue: q,
	}
}
func (w *WasmtimeSupportFunc) SetProcess(p *Process) {
	w.rt.SetProcess(p)
}

func (w *WasmtimeSupportFunc) SetMemPtr(u uintptr) {
	w.rt.SetMemPtr(u)
	w.file.SetMemPtr(u)
	w.queue.SetMemPtr(u)
}
