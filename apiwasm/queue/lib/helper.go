package lib

import (
	"context"

	pcontext "github.com/iansmith/parigot/context"

	"github.com/iansmith/parigot/g/queue/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
)

func FindOrCreateQueue(ctx context.Context, queueSvc queue.ClientQueue, name string) (queue.QueueId, queue.QueueErr) {
	req := queue.LocateRequest{}
	req.QueueName = name
	pcontext.Infof(ctx, "FindOrCreateQueue: looking for queue '%s'...", name)
	//real := queueSvc.(*queue.ClientQueue_)
	//smmap := real.ServiceMethodMap()
	// for _, v := range smmap.Pair() {
	// 	sid := id.UnmarshalServiceId(v.ServiceId)
	// 	mid := id.UnmarshalMethodId(v.MethodId)
	// }
	afterLocate := func(ctx context.Context, resp *queue.LocateResponse, err queue.QueueErr) syscall.KernelErr {

	}
	queueSvc.Locate(ctx, &req, afterLocate, locateErr)
	if err != queue.QueueErr_NoError && err == queue.QueueErr_NotFound {
		// it's a not found, so create it
		pcontext.Infof(ctx, "FindOrCreateQueue: looking for queue '%s'...", name)
		createReq := queue.CreateQueueRequest{}
		createReq.QueueName = name
		createResp, err := queueSvc.CreateQueue(ctx, &createReq)
		if err != queue.QueueErr_NoError {
			return queue.QueueIdZeroValue(), err
		}
		qid := queue.UnmarshalQueueId(createResp.GetId())
		return qid, queue.QueueErr_NoError
	}
	if err != queue.QueueErr_NoError {
		return queue.QueueIdZeroValue(), err
	}
	qid := queue.UnmarshalQueueId(resp.Id)
	return qid, queue.QueueErr_NoError
}
