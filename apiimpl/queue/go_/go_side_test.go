package go_

import (
	"context"
	"database/sql"
	"testing"
	"time"

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

const sender = "sender"
const payload = "payload"

func TestSendAndReceive(t *testing.T) {
	impl := createQueueService(t)
	qid := createQueueSuccess(t, impl, "sendAndRecv")
	qidM := lib.Marshal[protosupportmsg.QueueId](qid)

	// these kernel ids are just for use as content in the messages (the sender and payload)
	var senderK *protosupportmsg.KernelErrorId
	var payloadK [2]*protosupportmsg.KernelErrorId
	var payloadId [2]lib.Id
	senderK, payloadK[0], payloadK[1] = sendTwoMessagesForContent(t, impl, qid)
	senderId := lib.Unmarshal(senderK)
	payloadId[0] = lib.Unmarshal(payloadK[0])
	payloadId[1] = lib.Unmarshal(payloadK[1])

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

	iterationName := []string{sender, payload}

	for i := range []int{0, 1} {
		// the sender is real?
		if receiveResp.Message == nil || receiveResp.Message[i].GetSender() == nil {
			t.Errorf("no sender found on received message: %d", i)
			t.FailNow()
		}
		// really have payload?
		if receiveResp.Message[i].GetPayload() == nil {
			t.Errorf("no payload found on received message: %d", i)
			t.FailNow()
		}
		// check sender and payload
		for j, a := range []*anypb.Any{receiveResp.Message[i].GetSender(), receiveResp.Message[i].GetPayload()} {
			iterName := iterationName[j]
			candidateContent := protosupportmsg.KernelErrorId{}
			err := a.UnmarshalTo(&candidateContent)
			if err != nil {
				t.Errorf("unable to unmarshal ith %s (%d): %v", iterName, i, err)
			}
			candidateId := lib.Unmarshal(&candidateContent)
			if iterName == sender {
				if !candidateId.Equal(senderId) {
					t.Errorf("mismached sender on ith message (%d): %s vs %s", i, senderId.Short(), candidateId.Short())
				}
			} else {
				// deal with the lack of ordering by checking both
				if !candidateId.Equal(payloadId[0]) && !candidateId.Equal(payloadId[1]) {
					t.Errorf("unable to match payload on ith message (%d): %s", i, candidateId.Short())
				}
			}
		}
	}
}

func TestLenAndMarkdone(t *testing.T) {
	impl := createQueueService(t)
	qid := createQueueSuccess(t, impl, "lenAndMarkdone")
	qidM := lib.Marshal[protosupportmsg.QueueId](qid)

	// queue starts empty
	testLen(t, impl, qid, 0)

	rcvResp := queuemsg.ReceiveResponse{}
	//send two messages and test len
	sendAndReceiveMessages(t, impl, qid, &rcvResp, 2)
	// sendTwoMessagesForContent(t, impl, qid)
	// testLen(t, impl, qid, 2)

	// rcvReq := queuemsg.ReceiveRequest{
	// 	Id:           qidM,
	// 	MessageLimit: 2,
	// }
	// errId, errDetail := impl.QueueSvcReceiveImpl(&rcvReq, &rcvResp)
	// if errId != nil {
	// 	t.Errorf("unable to receive messages: %s: %s", errId.Short(), errDetail)
	// 	t.FailNow()
	// }

	if len(rcvResp.Message) != 2 {
		t.Errorf("expected to find 2 messages in queue, but found %d", len(rcvResp.Message))
	}
	// receive doesn't change the queue length
	testLen(t, impl, qid, 2)

	doneReq := queuemsg.MarkDoneRequest{}
	doneReq.Id = qidM
	doneResp := queuemsg.MarkDoneResponse{}
	doneReq.Msg = make([]*protosupportmsg.QueueMsgId, len(rcvResp.Message))
	for i := 0; i < len(doneReq.Msg); i++ {
		doneReq.Msg[i] = rcvResp.Message[i].GetMsgId()
	}
	errId, errDetail := impl.QueueSvcMarkDoneImpl(&doneReq, &doneResp)
	if errId != nil {
		t.Errorf("unable to mark messages as done: %s: %s", errId.Short(), errDetail)
		t.FailNow()
	}
	// markDone does change the queue length
	testLen(t, impl, qid, 0)

}

func TestReceivedTimeAndCount(t *testing.T) {
	name := "receive_and_count"
	impl := createQueueService(t)
	qid := createQueueSuccess(t, impl, name)
	qidM := lib.Marshal[protosupportmsg.QueueId](qid)

	rcvResp := queuemsg.ReceiveResponse{}
	//send two messages and test len
	sendAndReceiveMessages(t, impl, qid, &rcvResp, 2)
	//neither of these messages has been received before
	t1, t2, c1, c2 := twoReceivedMessagesToReceiveTimeAndCount(&rcvResp)
	if !t1.IsZero() || !t2.IsZero() {
		t.Errorf("unexpected Received value for a message: %#v", rcvResp.Message[0].GetReceived().AsTime())
	}
	if c1 != 0 || c2 != 0 {
		t.Errorf("unexpected ReceiveCount value for a message")
	}
	// read them again and check values since this is 2nd time of reading

	rcvReq := queuemsg.ReceiveRequest{
		Id:           qidM,
		MessageLimit: 2,
	}
	//receive again and testResult
	errId, errDetail := impl.QueueSvcReceiveImpl(&rcvReq, &rcvResp)
	if errId != nil {
		t.Errorf("unable to receive messages: %s: %s", errId.Short(), errDetail)
		t.FailNow()
	}
	t1, t2, c1, c2 = twoReceivedMessagesToReceiveTimeAndCount(&rcvResp)
	if t1.IsZero() || t2.IsZero() {
		t.Errorf("unexpected zero time value for message delivered before")
	}
	if c1 != 1 || c2 != 1 {
		t.Errorf("unexpected receive count for message delivered before")
	}
}

func TestLocate(t *testing.T) {
	name := "locate_test_queue"
	impl := createQueueService(t)
	qid := createQueueSuccess(t, impl, name)
	notFound := lib.NewQueueError(lib.QueueNotFound)
	req := queuemsg.LocateRequest{}
	resp := queuemsg.LocateResponse{}
	req.QueueName = "bleah"
	id, err := impl.QueueSvcLocateImpl(&req, &resp)
	if !id.Equal(notFound) {
		t.Errorf("locate produced something other than notFound for bleah: %s: %s", id.Short(), err)
	}
	req.QueueName = name
	id, err = impl.QueueSvcLocateImpl(&req, &resp)
	if id != nil {
		t.Errorf("locate failed: %s: %s", id.Short(), err)
	}
	candidate := lib.Unmarshal(resp.GetId())
	if !qid.Equal(candidate) {
		t.Errorf("locate failed, wrong id returned: expected %s but got %s", candidate.Short(), id.Short())
	}
}

//
// HELPERS
//

func testLen(t *testing.T, impl *QueueSvcImpl, qid lib.Id, expected int) {
	//t.Helper()
	qidM := lib.Marshal[protosupportmsg.QueueId](qid)
	lReq := queuemsg.LengthRequest{Id: qidM}
	lResp := queuemsg.LengthResponse{}
	id, deets := impl.QueueSvcLengthImpl(&lReq, &lResp)
	if id != nil {
		t.Errorf("unable check length of queue: %s: %s", id.Short(), deets)
		t.FailNow()
	}
	respQid := lib.Unmarshal(lResp.Id)
	if !respQid.Equal(qid) {
		t.Errorf("mismatched queue ids, expected %s but got %s", qid.Short(), respQid.Short())
	}
}

func sendTwoMessagesForContent(t *testing.T, impl *QueueSvcImpl, qid lib.Id) (*protosupportmsg.KernelErrorId, *protosupportmsg.KernelErrorId, *protosupportmsg.KernelErrorId) {
	t.Helper()

	qidM := lib.Marshal[protosupportmsg.QueueId](qid)

	var p0, p1, s anypb.Any
	kidSender := lib.Marshal[protosupportmsg.KernelErrorId](lib.NewKernelError(lib.KernelNamespaceExhausted))
	kidContent1 := lib.Marshal[protosupportmsg.KernelErrorId](lib.NewKernelError(lib.KernelDependencyCycle))
	kidContent2 := lib.Marshal[protosupportmsg.KernelErrorId](lib.NewKernelError(lib.KernelDataTooLarge))
	err1 := p0.MarshalFrom(kidContent1)
	err2 := p1.MarshalFrom(kidContent2)
	errSender := s.MarshalFrom(kidSender)
	if err1 != nil || err2 != nil || errSender != nil {
		t.Errorf("unable to marshal content")
		t.FailNow()
	}

	message0 := &queuemsg.QueueMsg{Id: qidM, Sender: &s, Payload: &p0}
	message1 := &queuemsg.QueueMsg{Id: qidM, Sender: &s, Payload: &p1}
	resp := &queuemsg.SendResponse{}
	sendMessageTestResult(t, impl, qid, []*queuemsg.QueueMsg{message0, message1}, resp, 0, 2)

	return kidSender, kidContent1, kidContent2

}

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
func sendAndReceiveMessages(t *testing.T, impl *QueueSvcImpl, qid lib.Id, rcvResp *queuemsg.ReceiveResponse, expectedCount int) {
	t.Helper()

	//send two messages and test len
	sendTwoMessagesForContent(t, impl, qid)
	testLen(t, impl, qid, expectedCount)
	qidM := lib.Marshal[protosupportmsg.QueueId](qid)

	rcvReq := queuemsg.ReceiveRequest{
		Id:           qidM,
		MessageLimit: int32(expectedCount),
	}
	//receive and testResult
	errId, errDetail := impl.QueueSvcReceiveImpl(&rcvReq, rcvResp)
	if errId != nil {
		t.Errorf("unable to receive messages: %s: %s", errId.Short(), errDetail)
		t.FailNow()
	}

}

func twoReceivedMessagesToReceiveTimeAndCount(rcvResp *queuemsg.ReceiveResponse) (time.Time, time.Time, int32, int32) {
	t1 := rcvResp.Message[0].GetReceived().AsTime()
	t2 := rcvResp.Message[1].GetReceived().AsTime()
	c1 := rcvResp.Message[0].GetReceiveCount()
	c2 := rcvResp.Message[1].GetReceiveCount()
	return t1, t2, c1, c2
}
