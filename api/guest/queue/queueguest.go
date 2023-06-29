package main

import (
	"context"
	"unsafe"

	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	queue "github.com/iansmith/parigot/g/queue/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/lib/go/future"
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

func (m *myQueueSvc) Ready(ctx context.Context, _ id.ServiceId) *future.Base[bool] {
	return future.NewBaseWithValue[bool](true)
}
func (m *myQueueSvc) CreateQueue(ctx context.Context, in *queue.CreateQueueRequest) *queue.FutureCreateQueue {
	return queue.CreateQueueHost(ctx, in)
}
func (m *myQueueSvc) Locate(ctx context.Context, in *queue.LocateRequest) *queue.FutureLocate {
	return queue.LocateHost(ctx, in)
}
func (m *myQueueSvc) DeleteQueue(ctx context.Context, in *queue.DeleteQueueRequest) *queue.FutureDeleteQueue {
	return queue.DeleteQueueHost(ctx, in)
}
func (m *myQueueSvc) Receive(ctx context.Context, in *queue.ReceiveRequest) *queue.FutureReceive {
	return queue.ReceiveHost(ctx, in)
}
func (m *myQueueSvc) MarkDone(ctx context.Context, in *queue.MarkDoneRequest) *queue.FutureMarkDone {
	return queue.MarkDoneHost(ctx, in)
}
func (m *myQueueSvc) Length(ctx context.Context, in *queue.LengthRequest) *queue.FutureLength {
	return queue.LengthHost(ctx, in)
}
func (m *myQueueSvc) Send(ctx context.Context, in *queue.SendRequest) *queue.FutureSend {
	return queue.SendHost(ctx, in)
}
