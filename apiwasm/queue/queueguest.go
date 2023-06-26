package main

import (
	"context"
	"unsafe"

	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"
	queue "github.com/iansmith/parigot/g/queue/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
)

var _ = unsafe.Sizeof([]byte{})

func main() {
	ctx := pcontext.CallTo(pcontext.SourceContext(context.Background(), pcontext.Guest), "fileguest.Main")
	f := &myQueueSvc{}
	binding := queue.Init(ctx, []lib.MustRequireFunc{}, f)
	var kerr syscall.KernelErr
	for {
		kerr = queue.ReadOneAndCall(ctx, binding, queue.TimeoutInMillis)
		if kerr == syscall.KernelErr_ReadOneTimeout {
			pcontext.Infof(ctx, "waiting for calls to queue service")
			continue
		}
		if kerr == syscall.KernelErr_NoError {
			continue
		}
		break
	}
	pcontext.Errorf(ctx, "error while waiting for queue service calls: %s", syscall.KernelErr_name[int32(kerr)])
}

type myQueueSvc struct {
}

func (m *myQueueSvc) Ready(ctx context.Context, _ id.ServiceId) *lib.BaseFuture[bool] {
	n := lib.NewBaseFuture[bool]()
	n.Set(true)
	return n
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
