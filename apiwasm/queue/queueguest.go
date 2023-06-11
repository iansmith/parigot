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
	sid := queue.MustRegisterQueue(ctx)
	queue.MustExportQueue(ctx)
	queue.MustWaitSatisfiedQueue(sid)
	queue.RunQueue(ctx, &myQueueSvc{})
}

type myQueueSvc struct {
}

func (m *myQueueSvc) Ready(ctx context.Context) bool {
	return true
}
func (m *myQueueSvc) CreateQueue(ctx context.Context, in *queue.CreateQueueRequest) (*queue.CreateQueueResponse, queue.QueueErr) {
	return queue.CreateQueueHost(in)
}
func (m *myQueueSvc) Locate(ctx context.Context, in *queue.LocateRequest) (*queue.LocateResponse, queue.QueueErr) {
	return queue.LocateHost(in)
}
func (m *myQueueSvc) DeleteQueue(ctx context.Context, in *queue.DeleteQueueRequest) (*queue.DeleteQueueResponse, queue.QueueErr) {
	return queue.DeleteQueueHost(in)
}
func (m *myQueueSvc) Receive(ctx context.Context, in *queue.ReceiveRequest) (*queue.ReceiveResponse, queue.QueueErr) {
	return queue.ReceiveHost(in)
}
func (m *myQueueSvc) MarkDone(ctx context.Context, in *queue.MarkDoneRequest) (*queue.MarkDoneResponse, queue.QueueErr) {
	return queue.MarkDoneHost(in)
}
func (m *myQueueSvc) Length(ctx context.Context, in *queue.LengthRequest) (*queue.LengthResponse, queue.QueueErr) {
	return queue.LengthHost(in)
}
func (m *myQueueSvc) Send(ctx context.Context, in *queue.SendRequest) (*queue.SendResponse, queue.QueueErr) {
	return queue.SendHost(in)
}
