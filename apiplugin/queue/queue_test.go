package main

import (
	"context"
	"testing"

	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
)

func TestCreateAndDelete(t *testing.T) {

	svc, err := newQueueSvc((context.Background()))
	if err.IsError() {
		t.FailNow()
	}
	creat := &queuemsg.CreateQueueRequest{}
	creat.QueueName = "foo$bar"
	resp := &queuemsg.CreateQueueResponse{}
	err = svc.QueueSvcCreateQueueImpl(creat)
	if !err.IsError() {
		if err.Code()!=
		t.FailNowf() // bad name
	}
	creat.QueueName = "foobar"

}
