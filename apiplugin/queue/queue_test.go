package main

import (
	"context"
	"testing"

	pcontext "github.com/iansmith/parigot/context"
	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
	"github.com/iansmith/parigot/g/queue/v1"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestCreateAndDelete(t *testing.T) {
	svc, errId := newQueueSvc((context.Background()))
	if errId.IsError() {
		t.FailNow()
	}
	// create bad name
	testQueueCreate(t, svc, "foo$bar", "bad name", true, 4)

	// create good name
	qid := testQueueCreate(t, svc, "foobar", "good name", false, 0)
	testQueueCreate(t, svc, "foobar", "duplicate name", true, 8)

	// delete it
	testQueueDelete(t, svc, qid, "simple delete", false, 0)
	// delete 2nd time
	testQueueDelete(t, svc, qid, "2nd delete", true, 7)

	// create it again
	qid2nd := testQueueCreate(t, svc, "foobar", "good name", false, 0)
	if qid.Equal(qid2nd) {
		t.Errorf("unexpected that second creation of a deleted queue gives same id")
	}
}

// Uses two dummy values:
// a BoolValue as the "sender" of this message
// a StringValue as the "payload" of this message
func TestQueueHappyPath(t *testing.T) {
	ctx := pcontext.DevNullContext(context.Background())

	svc, err := newQueueSvc((context.Background()))
	if err.IsError() {
		t.FailNow()
	}
	qid := setupQueue(t, svc)

	// send a test message
	msg := testCreateMessage(t, true, "this is a test", qid)
	sendReq := &queuemsg.SendRequest{
		Id: qid.Marshal(),
		Msg: []*queuemsg.QueueMsg{
			msg,
		},
	}
	sendResp := &queuemsg.SendResponse{}
	errId := svc.send(ctx, sendReq, sendResp)
	if errId.IsError() {
		t.Errorf("unable to send message")
	}
	savedMid := queue.MustUnmarshalQueueMsgId(sendResp.Succeed[0])
	//check queue len
	qlen := testQueueLen(t, svc, qid)
	if qlen != 1 {
		t.Errorf("bad queue length, expected 1 but got %d", qlen)
		t.FailNow()
	}
	// receive a message
	receiveReq := &queuemsg.ReceiveRequest{
		Id:           qid.Marshal(),
		MessageLimit: 1,
	}
	receiveResp := &queuemsg.ReceiveResponse{}
	errId = svc.receive(ctx, receiveReq, receiveResp)
	if errId.IsError() {
		t.Errorf("unable to receive message")
	}
	// check length of RECEIVED messages, not queue len
	lm := len(receiveResp.Message)
	if lm != 1 {
		t.Errorf("wrong number of messages received, got %d but expected %d", lm, 1)
	}
	//sendId := queue.MustUnmarshalQueueId(msg.GetMsgId())
	rmsg := receiveResp.Message[0]
	recvId := queue.MustUnmarshalQueueMsgId(rmsg.GetMsgId())
	if !savedMid.Equal(recvId) {
		t.Errorf("mismatched send (%s) and receive (%s) ids", savedMid.Short(), recvId.Short())
	}
}

//
// helpers
//

func testQueueDelete(t *testing.T, svc *queueSvcImpl, qid queue.QueueId, msg string, errorExpected bool, errorCode uint16) {
	ctx := pcontext.DevNullContext(context.Background())

	delReq := &queuemsg.DeleteQueueRequest{}
	delResp := &queuemsg.DeleteQueueResponse{}
	delReq.Id = qid.Marshal()
	err := svc.delete(ctx, delReq, delResp)
	qerr := queue.NewQueueErrIdFromRaw(err)
	if errorExpected {
		if !err.IsError() {
			t.Errorf("expected error from delete queue (%s): %s: %s", qid.Short(), msg, qerr.Short())
		}
		if qerr.ErrorCode() != 7 {
			t.Errorf("wrong error code from delete queue: %s: %s", msg, qerr.Short())
		}
		return
	}
	if qerr.IsError() {
		t.Errorf("uxexpected error from delete queue: %s: %s", msg, qerr.Short())
	}
	candidate := queue.MustUnmarshalQueueId(delResp.GetId())
	if !qid.Equal(candidate) {
		t.Errorf("created and deleted ids don't match")
	}

}

func testQueueCreate(t *testing.T, svc *queueSvcImpl, name, msg string, errorExpected bool, expectedCode uint16) queue.QueueId {
	ctx := pcontext.DevNullContext(context.Background())

	t.Helper()
	creat := &queuemsg.CreateQueueRequest{}
	creat.QueueName = name
	resp := &queuemsg.CreateQueueResponse{}
	err := queue.NewQueueErrIdFromRaw(svc.create(ctx, creat, resp))
	if errorExpected {
		if !err.IsError() {
			t.Errorf("expected error: %s :%s", msg, err.Short())
		}
		if queue.QueueErrIdCode(expectedCode) != err.ErrorCode() {
			t.Errorf("wrong code: %s, expected %d but got %d",
				msg, expectedCode, err.ErrorCode())
		}
		return queue.ZeroValueQueueId()
	}
	// no error expected case
	if err.IsError() {
		t.Errorf("unexpected error: %s :%s", msg, err.Short())
		return queue.ZeroValueQueueId()
	}

	return queue.MustUnmarshalQueueId(resp.GetId())
}

func setupQueue(t *testing.T, svc *queueSvcImpl) queue.QueueId {
	t.Helper()
	ctx := pcontext.DevNullContext(context.Background())

	creat := &queuemsg.CreateQueueRequest{}
	creat.QueueName = "unitTest"
	resp := &queuemsg.CreateQueueResponse{}
	err := queue.NewQueueErrIdFromRaw(svc.create(ctx, creat, resp))
	if err.IsError() {
		t.Errorf("expected queue to be created")
	}
	qid := queue.MustUnmarshalQueueId(resp.GetId())
	return qid
}

func testQueueLen(t *testing.T, svc *queueSvcImpl, qid queue.QueueId) int {
	ctx := pcontext.DevNullContext(context.Background())

	req := &queuemsg.LengthRequest{
		Id: qid.Marshal(),
	}
	resp := &queuemsg.LengthResponse{}
	rawErr := svc.length(ctx, req, resp)
	err := queue.NewQueueErrIdFromRaw(rawErr)
	if err.IsError() {
		t.Errorf("unable to get a clean response from length")
		t.FailNow()
	}
	return int(resp.GetLength())
}

func testCreateMessage(t *testing.T, senderValue bool, payloadValue string, qid queue.QueueId) *queuemsg.QueueMsg {
	sender, err := anypb.New(wrapperspb.Bool(senderValue))
	if err != nil {
		t.Errorf("unable to create sender protobuf because of error in marshal: %s", err.Error())
		t.FailNow()
	}
	payload, err := anypb.New(wrapperspb.String(payloadValue))
	if err != nil {
		t.Errorf("unable to create payload protobuf because of error in marshal: %s", err.Error())
		t.FailNow()
	}
	msg := &queuemsg.QueueMsg{
		Id:           qid.Marshal(),
		MsgId:        nil, // filled in by send
		ReceiveCount: 0,
		Received:     nil,
		Sender:       sender,
		Sent:         nil,
		Payload:      payload,
	}
	return msg
}
