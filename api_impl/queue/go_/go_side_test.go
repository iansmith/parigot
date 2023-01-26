package go_

import (
	"testing"

	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
	lib "github.com/iansmith/parigot/lib/go"
)

// WARNING: This type of unit test can ONLY be written against the
// WARNING: go side of a split impl service.  The WASM side cannot
// WARNING: run `go test` because it requires a filesystem and many
// WARNING: other things.  The test api is designed for the wasm side.

func TestCreateDelete(t *testing.T) {
	// we can only call the methods that don't take sp, but use req/resp
	impl, id, err := NewQueueSvc(nil)
	if id != nil {
		t.Logf("error creating queue service:%v, %s", id.Short(), err)
		t.FailNow()
	}

	var detail string
	id, detail = impl.QueueSvcCreateQueueImpl(&queuemsg.CreateQueueRequest{
		QueueName: "",
	}, &queuemsg.CreateQueueResponse{})
	if !id.IsError() || id.ErrorCode() != lib.QueueInvalidName {
		t.Errorf("expected error with Invalid name, but got: %s,%s", id.Short(), detail)
	}

	resp := &queuemsg.CreateQueueResponse{}
	id, detail = impl.QueueSvcCreateQueueImpl(&queuemsg.CreateQueueRequest{
		QueueName: "abc123",
	}, resp)
	if id != nil {
		t.Errorf("unexpected error from createQueue: %s,%s", id.Short(), detail)
	}
	if resp.GetId() == nil {
		t.Errorf("id not found in createQueue response")
	}
	t.Logf("id of queue: %s, %s", id.Short(), id.String())
}
