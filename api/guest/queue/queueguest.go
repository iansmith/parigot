package main

import (
	"context"
	"log/slog"
	"unsafe"

	"github.com/iansmith/parigot/api/guest"
	"github.com/iansmith/parigot/api/shared/id"
	queue "github.com/iansmith/parigot/g/queue/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/lib/go/future"
)

var _ = unsafe.Sizeof([]byte{})

var logger *slog.Logger

func main() {
	f := &myQueueSvc{}
	binding, fut, ctx, sid := queue.Init([]lib.MustRequireFunc{}, f)
	logger = slog.New(guest.NewParigotHandler(sid))

	fut.Success(func(_ *syscall.LaunchResponse) {
		var kerr syscall.KernelErr
		for {
			kerr = queue.ReadOneAndCall(ctx, binding, queue.TimeoutInMillis)
			if kerr == syscall.KernelErr_ReadOneTimeout {
				logger.Info("waiting for calls to queue service")
				continue
			}
			if kerr == syscall.KernelErr_NoError {
				continue
			}
			break
		}
		logger.Error("error while waiting for queue service calls", slog.String("syscall.KernelErr", syscall.KernelErr_name[int32(kerr)]))
	})
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
