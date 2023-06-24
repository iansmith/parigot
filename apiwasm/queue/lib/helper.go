package lib

import (
	"context"

	pcontext "github.com/iansmith/parigot/context"
	lib "github.com/iansmith/parigot/lib/go"

	"github.com/iansmith/parigot/g/queue/v1"
)

func FindOrCreateQueue(ctx context.Context, queueHandle queue.Client, name string) *lib.Future {
	req := queue.LocateRequest{}
	req.QueueName = name
	pcontext.Infof(ctx, "FindOrCreateQueue: looking for queue '%s'...", name)

	qidFuture := lib.NewFutureId()

	locateFuture := queueHandle.Locate(ctx, &req)
	locateFuture.Success(func(resp *queue.LocateResponse) {
		qidFuture.CompleteCall(resp.Id, 0)
	})
	locateFuture.Failure(queue.Future2LocateFailure[queue.Client, string](func(qerr queue.QueueErr, handle queue.Client, n string) {
		if qerr != queue.QueueErr_NotFound {
			return // somebody else's problem
		}
		createReq := &queue.CreateQueueRequest{
			QueueName: name,
		}
		fcreate := queueHandle.CreateQueue(ctx, createReq)
		fcreate.Success(func(resp *queue.CreateQueueResponse) {
			qidFuture.CompleteCall(resp.Id, 0)
		})
		fcreate.Failure(func(qe queue.QueueErr) {

		})
	}, queueHandle, name))
	return qidFuture
}
