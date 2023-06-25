package lib

import (
	"context"

	pcontext "github.com/iansmith/parigot/context"
	lib "github.com/iansmith/parigot/lib/go"

	"github.com/iansmith/parigot/g/queue/v1"
)

func FindOrCreateQueue(ctx context.Context, queueHandle queue.Client, name string) *lib.BaseFuture[queue.QueueId] {
	req := queue.LocateRequest{}
	req.QueueName = name
	pcontext.Infof(ctx, "FindOrCreateQueue: looking for queue '%s'...", name)

	qidFuture := lib.NewBaseFuture[queue.QueueId]()

	locateFuture := queueHandle.Locate(ctx, &req)
	locateFuture.Success(func(resp *queue.LocateResponse) {
		qidFuture.Set(queue.UnmarshalQueueId(resp.Id))
	})
	locateFuture.Failure(queue.Future2LocateFailure[queue.Client, string](func(qerr queue.QueueErr, handle queue.Client, n string) {
		if qerr != queue.QueueErr_NotFound {
			qidFuture.Set(queue.QueueIdZeroValue())
			return // somebody will have to handle the error
		}
		createReq := &queue.CreateQueueRequest{
			QueueName: name,
		}
		fcreate := queueHandle.CreateQueue(ctx, createReq)
		fcreate.Success(func(resp *queue.CreateQueueResponse) {
			qidFuture.Set(queue.UnmarshalQueueId(resp.Id))
		})
		fcreate.Failure(func(qe queue.QueueErr) {
			pcontext.Errorf(ctx, "unable to create queue for testing!")
			qidFuture.Set(queue.QueueIdZeroValue())
		})
	}, queueHandle, name))
	return qidFuture
}
