package lib

import (
	"context"

	pcontext "github.com/iansmith/parigot/context"

	"github.com/iansmith/parigot/g/queue/v1"
)

func FindOrCreateQueue(ctx context.Context, queueSvc queue.QueueService, name string) (queue.QueueId, queue.QueueErrId) {
	req := queue.LocateRequest{}
	req.QueueName = name
	pcontext.Infof(ctx, "FindOrCreateQueue: looking for queue '%s'...", name)
	resp, err := queueSvc.Locate(ctx, &req)
	if err.IsError() && err.ErrorCode() == 7 {
		// it's a not found, so create it
		pcontext.Infof(ctx, "FindOrCreateQueue: looking for queue '%s'...", name)
		createReq := queue.CreateQueueRequest{}
		createReq.QueueName = name
		createResp, err := queueSvc.CreateQueue(ctx, &createReq)
		if err.IsError() {
			return queue.ZeroValueQueueId(), err
		}
		qid := queue.MustUnmarshalQueueId(createResp.GetId())
		return qid, queue.QueueErrIdNoErr
	}
	if err.IsError() {
		return queue.ZeroValueQueueId(), err
	}
	qid := queue.MustUnmarshalQueueId(resp.Id)
	return qid, queue.QueueErrIdNoErr
}
