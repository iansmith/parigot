package main

import (
	"context"
	"fmt"
	"testing"

	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/queue/v1"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const queueNameInTest = "unit_test"

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
	// cleanup to leave queue in empty state
	testQueueDelete(t, svc, qid2nd, "simple delete", false, 0)
}

const payloadString = "this is a test"

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

	if testQueueLen(t, svc, qid) != 0 {
		t.Errorf("newly created queue should have no contents")
	}

	// send a test message
	senderValue := true
	msg := testCreateMessage(t, senderValue, payloadString, qid)
	sendReq := &queue.SendRequest{
		Id: qid.Marshal(),
		Msg: []*queue.QueueMsg{
			msg,
		},
	}
	sendResp := &queue.SendResponse{}
	errId := svc.send(ctx, sendReq, sendResp)
	if errId.IsError() {
		t.Errorf("unable to send message")
	}
	if testQueueLen(t, svc, qid) != 1 {
		t.Errorf("only one message sent, so should have 1 message")
	}
	savedMid := queue.MustUnmarshalQueueMsgId(sendResp.Succeed[0])
	//check queue len
	qlen := testQueueLen(t, svc, qid)
	if qlen != 1 {
		t.Errorf("bad queue length, expected 1 but got %d", qlen)
		t.FailNow()
	}
	// receive a message
	receiveReq := &queue.ReceiveRequest{
		Id:           qid.Marshal(),
		MessageLimit: 1,
	}
	receiveResp := &queue.ReceiveResponse{}
	receiveMessageNoError(t, svc, ctx, receiveReq, receiveResp)
	receivedOneCheckContent(t, ctx, qid, savedMid, svc, senderValue, payloadString)

	// read again, should receive same message again since we did not mark done
	receiveMessageNoError(t, svc, ctx, receiveReq, receiveResp)
	receivedOneCheckContent(t, ctx, qid, savedMid, svc, senderValue, payloadString)

	// check queue len
	if testQueueLen(t, svc, qid) != 1 {
		t.Errorf("read message twice and should have not changed queue content")
	}
	// mark done
	testMarkdone(t, ctx, svc, qid, savedMid)

	// mark done a second time has no effect
	testMarkdone(t, ctx, svc, qid, savedMid)

	ql := testQueueLen(t, svc, qid)
	if ql != 0 {
		t.Errorf("expected queue to be empty now but has %d messages", ql)
	}

	// cleanup to leave queue in empty state
	testQueueDelete(t, svc, qid, "simple delete", false, 0)

}

func TestLocateManyMessages(t *testing.T) {
	ctx := pcontext.DevNullContext(context.Background())

	svc, err := newQueueSvc((context.Background()))
	if err.IsError() {
		t.FailNow()
	}
	qid := setupQueue(t, svc)

	senderValue := true
	payload := make([]string, 10)
	md := make([]queue.QueueMsgId, 10)

	message := make([]*queue.QueueMsg, 10)
	for i := 0; i < 10; i++ {
		payload[i] = fmt.Sprintf("[i=%d]", i)
		message[i] = testCreateMessage(t, senderValue, payload[i], qid)
	}
	req := queue.SendRequest{}
	resp := queue.SendResponse{}
	req.Msg = message
	req.Id = qid.Marshal()

	rawErr := svc.send(ctx, &req, &resp)
	qErr := queue.NewQueueErrIdFromRaw(rawErr)
	if qErr.IsError() {
		t.Errorf("unexpected error from send: %s", qErr.Short())
		t.FailNow()
	}
	lsucc := len(resp.GetSucceed())
	if lsucc != 10 {
		t.Errorf("expected to send 10 but only sent %d", lsucc)
		t.FailNow()
	}
	for i := 0; i < 10; i++ {
		md[i] = queue.MustUnmarshalQueueMsgId(resp.GetSucceed()[i])
	}
	//
	// switch to receive side
	//
	locRequest := queue.LocateRequest{
		QueueName: queueNameInTest,
	}
	locResponse := queue.LocateResponse{}
	rawErr = svc.locate(ctx, &locRequest, &locResponse)
	qErr = queue.NewQueueErrIdFromRaw((rawErr))
	if qErr.IsError() {
		t.Errorf("unexpected failure of locate: %s", qErr.Short())
	}
	locId := queue.MustUnmarshalQueueId(locResponse.GetId())
	if !qid.Equal(locId) {
		t.Errorf("mismatched ids from create (%s) and locate (%s)", qid.Short(), locId.Short())
	}
	for i := 0; i < 9; i += 3 {
		rcvReq := &queue.ReceiveRequest{}
		rcvResp := &queue.ReceiveResponse{}
		rcvReq.Id = locId.Marshal()
		rcvReq.MessageLimit = 3
		rawErr := svc.receive(ctx, rcvReq, rcvResp)
		qErr := queue.NewQueueErrIdFromRaw(rawErr)
		if qErr.IsError() {
			t.Errorf("unable successfully receive iteration %d, %s", i, qErr.Short())
		}
		if len(rcvResp.GetMessage()) != 3 {
			t.Errorf("could not retrieve 3 messages successfully (sent %d)", len(rcvResp.GetMessage()))
		}
		dead := []*proto.IdRaw{}
		num := 3
		if len(rcvResp.Message) < 3 {
			t.Logf("warning, expected to get three elements back from the queue, but got %d", len(rcvResp.Message))
			num = len(rcvResp.Message)
		}
		for j := 0; j < num; j++ {
			msg := rcvResp.Message[j]
			var payloadWrap wrapperspb.StringValue
			if err := msg.GetPayload().UnmarshalTo(&payloadWrap); err != nil {
				t.Errorf("unable to unwrap and unmarshal payload")
				t.FailNow()
			}
			var found bool
			var index int
			found, index, payload = isInList[string](payload, payloadWrap.GetValue(), true)
			if !found {
				t.Errorf("unable to find payload %s in list of payloads", payloadWrap.GetValue())
				t.FailNow()
			}
			dead = append(dead, md[index].Marshal())
			target := queue.NewQueueMsgIdFromProto(msg.GetMsgId())
			found, _, md = isInListId(md, target, true)
			if !found {
				t.Errorf("unable to find msg id %s in list of payloads", target)
				t.FailNow()
			}
		}
		//mark this iteration done
		mdReq := queue.MarkDoneRequest{}
		mdResp := queue.MarkDoneResponse{}
		mdReq.Msg = dead
		mdReq.Id = locId.Marshal()
		rawErr = svc.markDone(ctx, &mdReq, &mdResp)
		mdErr := queue.NewQueueErrIdFromRaw(rawErr)
		if mdErr.IsError() {
			t.Errorf("unable to mark done successfully: %s", mdErr.Short())
		}
		if len(mdResp.GetUnmodified()) != 0 {
			t.Errorf("expected to mark done all but %d were left", len(mdResp.GetUnmodified()))
		}
		ql := testQueueLen(t, svc, locId)
		if ql != 10-(i+3) {
			t.Errorf("expected %d but got %d elements left [%d]", 10-(i+3), ql, len(payload))
			t.FailNow()
		}
	}
	// there should be one queue item left
	ql := testQueueLen(t, svc, qid)
	if ql != 1 {
		t.Errorf("expected to have one left but had %d", ql)
	}

}

//
// helpers
//

func testQueueDelete(t *testing.T, svc *queueSvcImpl, qid queue.QueueId, msg string, errorExpected bool, errorCode int32) {
	ctx := pcontext.DevNullContext(context.Background())

	delReq := &queue.DeleteQueueRequest{}
	delResp := &queue.DeleteQueueResponse{}
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

func testQueueCreate(t *testing.T, svc *queueSvcImpl, name, msg string, errorExpected bool, expectedCode int32) queue.QueueId {
	ctx := pcontext.DevNullContext(context.Background())

	t.Helper()
	creat := &queue.CreateQueueRequest{}
	creat.QueueName = name
	resp := &queue.CreateQueueResponse{}
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

	creat := &queue.CreateQueueRequest{}
	creat.QueueName = queueNameInTest
	resp := &queue.CreateQueueResponse{}
	err := queue.NewQueueErrIdFromRaw(svc.create(ctx, creat, resp))
	if err.IsError() {
		t.Errorf("expected queue to be created")
	}
	qid := queue.MustUnmarshalQueueId(resp.GetId())
	return qid
}

func testQueueLen(t *testing.T, svc *queueSvcImpl, qid queue.QueueId) int {
	ctx := pcontext.DevNullContext(context.Background())

	req := &queue.LengthRequest{
		Id: qid.Marshal(),
	}
	resp := &queue.LengthResponse{}
	rawErr := svc.length(ctx, req, resp)
	err := queue.NewQueueErrIdFromRaw(rawErr)
	if err.IsError() {
		t.Errorf("unable to get a clean response from length")
		t.FailNow()
	}
	return int(resp.GetLength())
}

func testCreateMessage(t *testing.T, senderValue bool, payloadValue string, qid queue.QueueId) *queue.QueueMsg {
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
	msg := &queue.QueueMsg{
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

func receivedOneCheckContent(t *testing.T, ctx context.Context, qid queue.QueueId,
	mid queue.QueueMsgId, svc *queueSvcImpl, sendValue bool, payloadValue string) {
	// receive a message
	receiveReq := &queue.ReceiveRequest{
		Id:           qid.Marshal(),
		MessageLimit: 1,
	}
	receiveResp := &queue.ReceiveResponse{}
	rawId := svc.receive(ctx, receiveReq, receiveResp)
	qerr := queue.NewQueueErrIdFromRaw(rawId)
	if qerr.IsError() {
		t.Errorf("unable to receive queue messages: %s", qerr.Short())
	}
	// testing number received, not number in queue
	if len(receiveResp.Message) != 1 {
		t.Errorf("wrong number of messages received, expected 1 but got %d", len(receiveResp.Message))
		t.FailNow()
	}
	errId := svc.receive(ctx, receiveReq, receiveResp)
	if errId.IsError() {
		qerrId := queue.NewQueueErrIdFromRaw(errId)
		t.Errorf("unable to receive message: %s", qerrId.Short())
		t.FailNow()
	}
	msg := receiveResp.Message[0]
	candMid := queue.MustUnmarshalQueueMsgId(msg.GetMsgId())
	if !candMid.Equal(mid) {
		t.Errorf("mismatched send (%s) and received id (%s)", mid.Short(), candMid.Short())
	}
	var wrapSender wrapperspb.BoolValue
	if err := msg.Sender.UnmarshalTo(&wrapSender); err != nil {
		t.Errorf("unmarshal of sender failed: %s", err.Error())
		t.FailNow()
	}
	if wrapSender.GetValue() != sendValue {
		t.Errorf("mismatched sender values")
	}
	var wrapPayload wrapperspb.StringValue
	if err := msg.Payload.UnmarshalTo(&wrapPayload); err != nil {
		t.Errorf("unmarshal of payload failed: %s", err.Error())
		t.FailNow()
	}
	if wrapPayload.GetValue() != payloadValue {
		t.Errorf("mismatched sender values (%s and %s)",
			wrapPayload.GetValue(), payloadValue)
	}
}

func receiveMessageNoError(t *testing.T, svc *queueSvcImpl, ctx context.Context,
	req *queue.ReceiveRequest, resp *queue.ReceiveResponse) {
	errId := svc.receive(ctx, req, resp)
	if errId.IsError() {
		t.Errorf("unable to receive message")
	}
}

func testMarkdone(t *testing.T, ctx context.Context, svc *queueSvcImpl, qid queue.QueueId, savedMid queue.QueueMsgId) {
	mdReq := &queue.MarkDoneRequest{}
	mdResp := &queue.MarkDoneResponse{}
	mdReq.Id = qid.Marshal()
	mdReq.Msg = []*proto.IdRaw{savedMid.Marshal()}

	rawErr := svc.markDone(ctx, mdReq, mdResp)
	qErr := queue.NewQueueErrIdFromRaw(rawErr)
	if qErr.IsError() {
		t.Errorf("unable to mark message %s done", savedMid.Short())
		t.FailNow()
	}
	candQid := queue.MustUnmarshalQueueId(mdResp.GetId())
	if !candQid.Equal(qid) {
		t.Errorf("mismatched original (%s) and received (%s) queue id",
			qid.Short(), candQid.Short())
	}
	if len(mdResp.Unmodified) != 0 {
		t.Errorf("expected to have no unmodified messages, but had %d", len(mdResp.Unmodified))
	}
}

func compareEqual[T comparable](self, other T) bool {
	return self == other
}
func isInList[T comparable](list []T, cand T, remove bool) (bool, int, []T) {
	return isInListAny[T](list, cand, remove, compareEqual[T])
}
func isInListId(list []queue.QueueMsgId, cand queue.QueueMsgId, remove bool) (bool, int, []queue.QueueMsgId) {
	comp := func(t, u queue.QueueMsgId) bool { return t.Equal(u) }
	return isInListAny(list, cand, remove, comp)
}

func isInListAny[T any](list []T, cand T, remove bool, equalFn func(T, T) bool) (bool, int, []T) {

	if len(list) == 0 {
		return false, -2, list
	}
	found := -2
	for i := 0; i < len(list); i++ {
		s := list[i]
		if equalFn(s, list[i]) {
			found = i
			break
		}
	}
	if found < 0 {
		return false, -2, list
	}
	if !remove {
		return true, found, list
	}
	if len(list) == 1 {
		return true, found, nil
	}
	if found == 0 {
		return true, found, list[1:]
	}
	// can't be the only element, checked that
	if found == len(list)-1 {
		return true, found, list[:found]
	}

	return true, found, append(list[:found], list[found+1:]...)
}
