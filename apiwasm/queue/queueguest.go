package main

import (
	"context"
	"unsafe"

	pcontext "github.com/iansmith/parigot/context"
	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
	gqueue "github.com/iansmith/parigot/g/queue/v1"
)

var _ = unsafe.Sizeof([]byte{})

func main() {
	ctx := pcontext.GuestContext(pcontext.NewContextWithContainer(context.Background(), "[queuewasm]main"))
	gqueue.MustRegisterQueueService(ctx)
	gqueue.MustExportQueueService(ctx)
	gqueue.RunQueueService(ctx, &myQueueSvc{})
}

type myQueueSvc struct {
}

func (m *myQueueSvc) Ready(ctx context.Context) bool {
	return true
}
func (m *myQueueSvc) CreateQueue(ctx context.Context, in *queuemsg.CreateQueueRequest) (*queuemsg.CreateQueueResponse, gqueue.QueueErrId) {
	return gqueue.CreateQueueHost(in)
}
func (m *myQueueSvc) Locate(ctx context.Context, in *queuemsg.LocateRequest) (*queuemsg.LocateResponse, gqueue.QueueErrId) {
	return gqueue.LocateHost(in)
}
func (m *myQueueSvc) DeleteQueue(ctx context.Context, in *queuemsg.DeleteQueueRequest) (*queuemsg.DeleteQueueResponse, gqueue.QueueErrId) {
	return gqueue.DeleteQueueHost(in)
}
func (m *myQueueSvc) Receive(ctx context.Context, in *queuemsg.ReceiveRequest) (*queuemsg.ReceiveResponse, gqueue.QueueErrId) {
	return gqueue.ReceiveHost(in)
}
func (m *myQueueSvc) MarkDone(ctx context.Context, in *queuemsg.MarkDoneRequest) (*queuemsg.MarkDoneResponse, gqueue.QueueErrId) {
	return gqueue.MarkDoneHost(in)
}
func (m *myQueueSvc) Length(ctx context.Context, in *queuemsg.LengthRequest) (*queuemsg.LengthResponse, gqueue.QueueErrId) {
	return gqueue.LengthHost(in)
}
func (m *myQueueSvc) Send(ctx context.Context, in *queuemsg.SendRequest) (*queuemsg.SendResponse, gqueue.QueueErrId) {
	return gqueue.SendHost(in)
}
