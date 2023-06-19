package lib

import (
	"context"
	"log"

	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"

	"github.com/iansmith/parigot/g/queue/v1"
)

func FindOrCreateQueue(ctx context.Context, queueSvc queue.Queue, name string) (queue.QueueId, queue.QueueErr) {
	req := queue.LocateRequest{}
	req.QueueName = name
	pcontext.Infof(ctx, "FindOrCreateQueue: looking for queue '%s'...", name)
	real := queueSvc.(*queue.ClientQueue_)
	smmap := real.ServiceMethodMap()
	log.Printf("size of the Pair list is %d", len(smmap.Pair()))
	for _, v := range smmap.Pair() {
		sid := id.UnmarshalServiceId(v.ServiceId)
		mid := id.UnmarshalMethodId(v.MethodId)
		log.Printf("xxxx --- %s, %s", sid.Short(), mid.Short())
	}
	resp, err := queueSvc.Locate(ctx, &req)
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
