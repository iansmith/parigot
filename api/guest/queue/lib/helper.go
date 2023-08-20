package lib

import (
	"context"

	"github.com/iansmith/parigot/api/guest"
	"github.com/iansmith/parigot/lib/go/future"
	"google.golang.org/protobuf/proto"

	"github.com/iansmith/parigot/g/queue/v1"
)

// FindOrCreateQueue gets a queue by name, creating it if necessary.  The return
// value will be the queue.QueueIdZeroValue if there was any error,
// which is a bit dodgy.
func FindOrCreateQueue(ctx context.Context, queueHandle queue.Client, name string) *future.Base[queue.QueueId] {
	req := queue.LocateRequest{}
	req.QueueName = name

	guest.Log(ctx).Info("FindOrCreateQueue: looking for queue '%s'...", name)

	qidFuture := future.NewBase[queue.QueueId]()

	locateFuture := queueHandle.Locate(ctx, &req)
	locateFuture.Success(func(raw proto.Message) {
		resp := raw.(*queue.LocateResponse)
		qidFuture.Set(queue.UnmarshalQueueId(resp.Id))
	})
	locateFuture.Failure(func(raw int32) {
		qerr := queue.QueueErr(raw)
		if qerr != queue.QueueErr_NotFound {
			qidFuture.Set(queue.QueueIdZeroValue())
		}
		createReq := &queue.CreateQueueRequest{
			QueueName: name,
		}
		fcreate := queueHandle.CreateQueue(ctx, createReq)
		fcreate.Success(func(raw proto.Message) {
			resp := raw.(*queue.CreateQueueResponse)
			qidFuture.Set(queue.UnmarshalQueueId(resp.Id))
		})
		fcreate.Failure(func(raw int32) {
			guest.Log(ctx).Error("unable to create queue for testing: %s", queue.QueueErr_name[raw])
			qidFuture.Set(queue.QueueIdZeroValue()) //xxx hack with zero value
		})
	})
	return qidFuture
}
