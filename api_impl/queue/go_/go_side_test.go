package go_

import (
	"context"
	"database/sql"
	"testing"

	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"google.golang.org/protobuf/types/known/anypb"
)

// WARNING: This type of unit test can ONLY be written against the
// WARNING: go side of a split impl service.  The WASM side cannot
// WARNING: run `go test` because it requires a filesystem and many
// WARNING: other things.  The test api is designed for the wasm side.

func TestCreateDelete(t *testing.T) {
	// we can only call the methods that don't take sp, but use req/resp

	impl := createQueueService(t)

	// expect failure with no name
	id, detail := impl.QueueSvcCreateQueueImpl(&queuemsg.CreateQueueRequest{
		QueueName: "",
	}, &queuemsg.CreateQueueResponse{})
	if !id.IsError() || id.ErrorCode() != lib.QueueInvalidName {
		t.Errorf("expected error with Invalid name, but got: %s,%s", id.Short(), detail)
	}

	qid := createQueueSuccess(t, impl, "abc123")

	checkIdExistence(t, impl, qid, true)

	respDel := &queuemsg.DeleteQueueResponse{}
	errId, errDetail := impl.QueueSvcDeleteQueueImpl(&queuemsg.DeleteQueueRequest{
		Id: lib.Marshal[protosupportmsg.QueueId](qid),
	}, respDel)
	candidate := lib.Unmarshal(respDel.GetId())
	if !qid.Equal(candidate) {
		t.Errorf("error id returned : %v, %s", errId, errDetail)
	}
	checkIdExistence(t, impl, qid, false)
}

func TestSendAndReceive(t *testing.T) {
	impl := createQueueService(t)
	qid := createQueueSuccess(t, impl, "sendAndRecv")
	qidM := lib.Marshal[protosupportmsg.QueueId](qid)

	// these kernel ids are just for use as content in the messages
	kidSender := lib.Marshal[protosupportmsg.KernelErrorId](lib.NewKernelError(lib.KernelNamespaceExhausted))
	kidContent1 := lib.Marshal[protosupportmsg.KernelErrorId](lib.NewKernelError(lib.KernelDependencyCycle))
	kidContent2 := lib.Marshal[protosupportmsg.KernelErrorId](lib.NewKernelError(lib.KernelDataTooLarge))
	content1, content2, contentSender := anypb.Any{}, anypb.Any{}, anypb.Any{}
	err1 := content1.MarshalFrom(kidContent1)
	err2 := content2.MarshalFrom(kidContent2)
	errSender := contentSender.MarshalFrom(kidSender)
	if err1 != nil || err2 != nil || errSender != nil {
		t.Errorf("unable to marshal content")
		t.FailNow()
	}

	message := &queuemsg.QueueMsg{Id: qidM, Sender: &contentSender, Payload: &content1}
	resp := &queuemsg.SendResponse{}
	sendMessageTestResult(t, impl, qid, []*queuemsg.QueueMsg{message}, resp, 0, 1)

	message = &queuemsg.QueueMsg{Id: qidM, Sender: &contentSender, Payload: &content2}
	resp = &queuemsg.SendResponse{}
	sendMessageTestResult(t, impl, qid, []*queuemsg.QueueMsg{message}, resp, 0, 1)

	receiveReq := &queuemsg.ReceiveRequest{
		Id:           qidM,
		MessageLimit: 2,
	}
	receiveResp := &queuemsg.ReceiveResponse{}
	id, detail := impl.QueueSvcReceiveImpl(receiveReq, receiveResp)
	if id != nil {
		t.Errorf("failed to receive: %s, %s", id.Short(), detail)
		t.FailNow()
	}
	if len(receiveResp.GetMessage()) != 2 {
		t.Errorf("failed to receive, expected 2 messages but got %d", len(receiveResp.GetMessage()))
	}
}

//
// HELPERS
//

func sendMessageTestResult(t *testing.T, q *QueueSvcImpl, qid lib.Id, msg []*queuemsg.QueueMsg, resp *queuemsg.SendResponse, numFail, numSucc int) {
	t.Helper()

	req := &queuemsg.SendRequest{
		Id:  lib.Marshal[protosupportmsg.QueueId](qid),
		Msg: msg,
	}
	errId, detail := q.QueueSvcSendImpl(req, resp)
	if errId != nil {
		t.Errorf("unable to send messages: %s: %s", errId.Short(), detail)
		t.FailNow()
	}
	if len(resp.Fail) != numFail {
		t.Errorf("send failed, %d message listed as failed", len(resp.Fail))
	}
	if len(resp.Succeed) != numSucc {
		t.Errorf("send failed, %d message listed as succeeded", len(resp.Succeed))
	}
}

// createQueueService creates a new service and returns it.  If anything goes
// wrong, it uses FailNow().
func createQueueService(t *testing.T) *QueueSvcImpl {
	impl, id, err := NewQueueSvc(nil)
	if id != nil {
		t.Logf("error creating queue service:%v, %s", id.Short(), err)
		t.FailNow()
	}
	impl.queries.testDestroyAll(context.Background())
	return impl
}

// createQueueSuccess creates a queue and returns its id. If anything goes
// wrong, it uses FailNow().
func createQueueSuccess(t *testing.T, impl *QueueSvcImpl, name string) lib.Id {
	t.Helper()
	resp := &queuemsg.CreateQueueResponse{}
	id, detail := impl.QueueSvcCreateQueueImpl(&queuemsg.CreateQueueRequest{
		QueueName: name,
	}, resp)
	if id != nil {
		t.Errorf("unexpected error from createQueue: %s,%s", id.Short(), detail)
		t.FailNow()
	}
	if resp.GetId() == nil {
		t.Errorf("id not found in createQueue response")
		t.FailNow()
	}

	qid := lib.Unmarshal(resp.Id)
	return qid
}

func checkIdExistence(t *testing.T, impl *QueueSvcImpl, id lib.Id, exists bool) {
	t.Helper()
	idToKey, err := impl.queries.getKeyFromQueueId(context.Background(),
		getKeyFromQueueIdParams{
			IDLow:  sql.NullInt64{Int64: int64(id.Low()), Valid: true},
			IDHigh: sql.NullInt64{Int64: int64(id.High()), Valid: true},
		})
	if exists {
		if err != nil {
			t.Errorf("unable to find expected rowid for id %v", err)
		}
	} else {
		if err == nil {
			t.Errorf("did not expect to find row for id %s: %+v", id.Short(), idToKey)
		}
	}
}
