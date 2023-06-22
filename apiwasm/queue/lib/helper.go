package lib

import (
	"context"

	pcontext "github.com/iansmith/parigot/context"
	lib "github.com/iansmith/parigot/lib/go"

	"github.com/iansmith/parigot/g/queue/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
)

func FindOrCreateQueue(ctx context.Context, queueHandle queue.ClientQueue, name string) lib.NewPromise[queue.QueueId, queue.QueueErr] {
	req := queue.LocateRequest{}
	req.QueueName = name
	pcontext.Infof(ctx, "FindOrCreateQueue: looking for queue '%s'...", name)
	p := lib.NewPromise[queue.QueueId, queue.QueueErr]()
	queueHandle.Locate(ctx, &req).
		OnSuccess(ctx, func(resp *queue.LocateResponse) {
			p.Resolve(resp.GetId(), queue.QueueErr_None)
		}).
		OnFailure(ctx, func(err queue.QueueErr) {
			if err == queue.QueueErr_NotFound {

			}
		})

	//real := queueSvc.(*queue.ClientQueue_)
	//smmap := real.ServiceMethodMap()
	// for _, v := range smmap.Pair() {
	// 	sid := id.UnmarshalServiceId(v.ServiceId)
	// 	mid := id.UnmarshalMethodId(v.MethodId)
	// }
	afterLocate := func(ctx context.Context, resp *queue.LocateResponse, err queue.QueueErr) syscall.KernelErr {
		if err != queue.QueueErr_NoError {
			if err == queue.QueueErr_NotFound {

			}
		}
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
