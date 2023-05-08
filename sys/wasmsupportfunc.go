package sys

import (
	filego "github.com/iansmith/parigot/apiimpl/file/go_"
	loggo "github.com/iansmith/parigot/apiimpl/log/go_"
	queuego "github.com/iansmith/parigot/apiimpl/queue/go_"
)

type WasmtimeSupportFunc struct {
	rt    *Runtime
	file  *filego.FileSvcImpl
	log   *loggo.LogViewerImpl
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
		log:   loggo.NewLogViewerImpl(),
		queue: q,
	}
}
func (w *WasmtimeSupportFunc) SetProcess(p *Process) {
	w.rt.SetProcess(p)
}

func (w *WasmtimeSupportFunc) SetMemPtr(u uintptr) {
	w.rt.SetMemPtr(u)
	w.file.SetMemPtr(u)
	w.log.SetMemPtr(u)
	w.queue.SetMemPtr(u)
}
