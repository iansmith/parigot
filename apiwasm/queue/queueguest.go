package main

import (
	"context"
	"unsafe"

	pcontext "github.com/iansmith/parigot/context"
	queue "github.com/iansmith/parigot/g/queue/v1"
)

var _ = unsafe.Sizeof([]byte{})

func main() {
	ctx := pcontext.GuestContext(pcontext.NewContextWithContainer(context.Background(), "[queuewasm]main"))
	queue.MustRegisterQueueService(ctx)
	queue.MustExportQueueService(ctx)
	queue.RunQueueService(ctx, &myQueueSvc{})
}

type myQueueSvc struct {
}

func (m *myQueueSvc) Ready(ctx context.Context) bool {
	return true
}
func (m *myQueueSvc) CreateQueue(ctx context.Context, in *queue.CreateQueueRequest) (*queue.CreateQueueResponse, queue.QueueErrId) {
	return queue.CreateQueueHost(in)
}
func (m *myQueueSvc) Locate(ctx context.Context, in *queue.LocateRequest) (*queue.LocateResponse, queue.QueueErrId) {
	return queue.LocateHost(in)
}
func (m *myQueueSvc) DeleteQueue(ctx context.Context, in *queue.DeleteQueueRequest) (*queue.DeleteQueueResponse, queue.QueueErrId) {
	return queue.DeleteQueueHost(in)
}
func (m *myQueueSvc) Receive(ctx context.Context, in *queue.ReceiveRequest) (*queue.ReceiveResponse, queue.QueueErrId) {
	return queue.ReceiveHost(in)
}
func (m *myQueueSvc) MarkDone(ctx context.Context, in *queue.MarkDoneRequest) (*queue.MarkDoneResponse, queue.QueueErrId) {
	return queue.MarkDoneHost(in)
}
func (m *myQueueSvc) Length(ctx context.Context, in *queue.LengthRequest) (*queue.LengthResponse, queue.QueueErrId) {
	return queue.LengthHost(in)
}
func (m *myQueueSvc) Send(ctx context.Context, in *queue.SendRequest) (*queue.SendResponse, queue.QueueErrId) {
	return queue.SendHost(in)
}
