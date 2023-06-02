package main

import (
	"context"
	"database/sql"
	"log"
	"regexp"
	"time"

	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/eng"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
	"github.com/iansmith/parigot/g/queue/v1"

	"github.com/tetratelabs/wazero/api"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed sqlc/schema.sql
var ddl string

var inMemDSN = "file:test.db?cache=shared&mode=memory"

const longForm = "2006-01-02 15:04:05"

// This impl is part of the GO side, it is not visible in WASM.
type queueSvcImpl struct {
	db      *sql.DB
	queries *Queries
	ctx     context.Context
}

// XXXXX why can't I use this?
// QueueInternalError
var internalErr = queue.NewQueueErrId(5)
var queueSvc *queueSvcImpl

type queuePlugin struct{}

var ParigotInitialize = queuePlugin{}

func (*queuePlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "queue", "create_queue_", createQueueHost)
	e.AddSupportedFunc(ctx, "queue", "delete_queue_", deleteQueueHost)
	e.AddSupportedFunc(ctx, "queue", "length_", lengthHost)
	e.AddSupportedFunc(ctx, "queue", "locate_", locateHost)
	e.AddSupportedFunc(ctx, "queue", "mark_done_", markDoneHost)
	e.AddSupportedFunc(ctx, "queue", "receive_", receiveHost)
	e.AddSupportedFunc(ctx, "queue", "send_", sendHost)

	var errId queue.QueueErrId
	queueSvc, errId = newQueueSvc(ctx)
	if errId.IsError() {
		pcontext.Fatalf(ctx, "QueueSvc: unable to start:%s", errId.Short())
		return false
	}

	return true
}

func createQueueHost(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("queue.createQueue 0x%x", stack)
}

func deleteQueueHost(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("queue.deleteQueue 0x%x", stack)
}
func lengthHost(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("queue.length 0x%x", stack)
}
func locateHost(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("queue.locate 0x%x", stack)
}
func markDoneHost(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("queue.markDone 0x%x", stack)
}
func sendHost(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("queue.send 0x%x", stack)
}
func receiveHost(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("queue.receive 0x%x", stack)
}

// newQueueSvc returns an initialized Queue service.
func newQueueSvc(ctx context.Context) (*queueSvcImpl, queue.QueueErrId) {
	newCtx := pcontext.ServerGoContext(ctx)

	var err error
	q := &queueSvcImpl{}
	//open db
	q.db, err = sql.Open("sqlite3", inMemDSN)
	if err != nil {
		return nil, internalErr
	}
	// create tables
	if _, err := q.db.ExecContext(context.Background(), ddl); err != nil {
		return nil, internalErr
	}

	q.queries = New(q.db)
	q.ctx = newCtx
	return q, queue.QueueErrIdNoErr
}

var legalName = regexp.MustCompile(`[A-Za-z0-9_\-\.]+`)

func (q *queueSvcImpl) validateName(name string) bool {
	return legalName.MatchString(name)
}

// QueueSvcCreateQueueImpl is separate from the "real" call of
// QueueSvcCreateQueue so it is easy to test.  The return values
// will be nil, "" if there was no error.  If there was an error
// these are ready to be returned to the WASM side.
func (q *queueSvcImpl) QueueSvcCreateQueueImpl(req *queuemsg.CreateQueueRequest,
	resp *queuemsg.CreateQueueResponse) queue.QueueErrId {

	ctx := pcontext.CallTo(q.ctx, "QueueSvcCreateQueueImpl")

	if !q.validateName(req.GetQueueName()) {
		pcontext.Errorf(ctx, "queue name is not valid: '%s'", req.GetQueueName())
		return queue.NewQueueErrId(4) //xxx QueueErrIdCode
	}
	mod, err := q.queries.CreateQueue(context.Background(), req.GetQueueName())
	if err != nil {
		return internalErr
	}

	qid := queue.NewQueueId()
	param := CreateIdToKeyMappingParams{
		IDLow:    sql.NullInt64{Int64: int64(qid.Low()), Valid: true},
		IDHigh:   sql.NullInt64{Int64: int64(qid.High()), Valid: true},
		QueueKey: sql.NullInt64{Int64: mod.ID, Valid: true},
	}
	_, err = q.queries.CreateIdToKeyMapping(context.Background(), param)
	if err != nil {
		pcontext.Errorf(ctx, "unable to create id to key mapping: %s", err.Error())
		return internalErr
	}
	resp.Id = qid.Marshal()
	return queue.QueueErrIdNoErr
}

// QueueSvcDeleteQueueImpl is separate from the "real" call of
// QueueSvcDeleteQueue so it is easy to test.  The return values
// will be an QueueErrId.
func (q *queueSvcImpl) QueueSvcDeleteQueueImpl(req *queuemsg.DeleteQueueRequest,
	resp *queuemsg.DeleteQueueResponse) queue.QueueErrId {

	ctx := pcontext.CallTo(q.ctx, "QueueSvcDeleteQueueImpl")

	qid := queue.MustUnmarshalQueueId(req.Id)
	r, err := q.queries.getKeyFromQueueId(ctx,
		getKeyFromQueueIdParams{
			IDLow:  sql.NullInt64{Int64: int64(qid.Low()), Valid: true},
			IDHigh: sql.NullInt64{Int64: int64(qid.High()), Valid: true},
		})
	if err != nil {
		pcontext.Errorf(ctx, "unable to find db key for %s, %v:", qid.Short(), err.Error())
		return internalErr
	}
	err = q.queries.DeleteQueue(context.Background(), r.QueueKey.Int64)
	if err != nil {
		pcontext.Errorf(ctx, "unable to delete row for key %s: %v:", qid.Short(), err.Error())
		return internalErr
	}
	// just return the internal id, not the row id
	resp.Id = req.GetId()
	return queue.QueueErrIdNoErr
}

// getRowidForId is an internal function to convert a queue id
// into a row id that can be used as a key in other queries.  If anything went wrong
// it will return an internal error.
func (q *queueSvcImpl) getRowidForId(ctx context.Context, qid queue.QueueId) (queue.RowId, queue.QueueErrId) {

	ctx = pcontext.CallTo(q.ctx, "getRowForId")

	p := getKeyFromQueueIdParams{
		IDLow:  sql.NullInt64{Int64: int64(qid.Low()), Valid: true},
		IDHigh: sql.NullInt64{Int64: int64(qid.High()), Valid: true},
	}
	result, err := q.queries.getKeyFromQueueId(ctx, p)
	if err != nil {
		pcontext.Errorf(ctx, "unable to find internal key from id (%s): %v", qid.Short(), err)
		return queue.ZeroValueRowId(), internalErr
	}
	//?? result.QueueKey.Int64,
	rid := id.NewIdTyped[queue.DefRow](0, uint64(result.QueueKey.Int64))
	return queue.RowId(rid), queue.QueueErrIdNoErr

}

// QueueSvcSendImpl is separate from the "real" call of
// QueueSvcSend so it is easy to test.  If there was an error
// it returns the successfully sent and failed messages. If you are
// trying to send many messages, be sure to look at these two lists
// because partial failure is possible.
func (q *queueSvcImpl) QueueSvcSendImpl(req *queuemsg.SendRequest, resp *queuemsg.SendResponse) queue.QueueErrId {
	ctx := pcontext.CallTo(q.ctx, "QueueSvcSendImpl")

	qid := queue.MustUnmarshalQueueId(req.GetId())

	rid, err := q.getRowidForId(ctx, qid)
	if err.IsError() {
		return err
	}
	succeed := []*protosupportmsg.IdRaw{}
	fail := []*queuemsg.QueueMsg{}

	alreadyFailed := false
	var failedOn *protosupportmsg.IdRaw

	for _, current := range req.GetMsg() {
		if alreadyFailed {
			fail = append(fail, current)
			continue
		}
		mid := queue.NewQueueMsgId()
		// flatten sender
		var senderBytes []byte
		var err error
		if current.GetSender() != nil {
			senderBytes, err = proto.Marshal(current.GetSender())
			if err != nil {
				pcontext.Errorf(ctx, "unable to flatten sender: %v", err)
				alreadyFailed = true
				fail = append(fail, current)
				failedOn = current.GetMsgId()
				continue
			}
		}
		// flatten payload
		payloadBytes, err := proto.Marshal(current.GetPayload())
		if err != nil {
			pcontext.Errorf(ctx, "unable to flatten payload: %v", err)
			return internalErr
		}
		params := CreateMessageParams{
			IDLow:    sql.NullInt64{Int64: int64(mid.Low()), Valid: true},
			IDHigh:   sql.NullInt64{Int64: int64(mid.High()), Valid: true},
			QueueKey: sql.NullInt64{Int64: int64(rid.Low()), Valid: true},
			Sender:   senderBytes,
			Payload:  payloadBytes,
		}

		/// XXX what about the failed id?
		// is this right?
		var fileMsgId queue.QueueMsgId
		_, err = q.queries.CreateMessage(ctx, params)
		if err != nil {
			alreadyFailed = true
			fail = append(fail, current)
			failedOn = fileMsgId.Marshal()
			pcontext.Errorf(ctx, "could not create message: %v", err)
			continue
		}
		succeed = append(succeed, mid.Marshal())
	}
	if alreadyFailed {
		return internalErr
	}
	resp.Fail = fail
	resp.Succeed = succeed
	resp.FailedOn = failedOn
	return queue.QueueErrIdNoErr
}

// queueSvcLengthImpl is separate from the "real" call of
// QueueSvcLength so it is easy to test.  The return values
// will be nil, "" if there was no error.  If there was an error you'll
// get the error code and detail.  If there was no error, the response
// object will have the apporimate length of the queue requested.
func (q *queueSvcImpl) QueueSvcLengthImpl(req *queuemsg.LengthRequest, resp *queuemsg.LengthResponse) queue.QueueErrId {
	ctx := pcontext.CallTo(q.ctx, "QueueSvcLengthImpl")

	qid := queue.MustUnmarshalQueueId(req.GetId())
	rowId, err := q.getRowidForId(ctx, qid)
	if err.IsError() {
		pcontext.Errorf(ctx, "could not find the queue %s: %s", qid.Short(), err.Short())
		return err
	}

	count, oldErr := q.queries.Length(ctx, sql.NullInt64{Int64: int64(rowId.Low()), Valid: true})
	if oldErr != nil {
		pcontext.Errorf(ctx, "failed length query: on queue:%s: %s", qid.Short(), oldErr.Error())
		return internalErr
	}
	resp.Id = req.Id
	resp.Length = count
	return queue.QueueErrIdNoErr
}

// QueueSvcMarkeDoneImpl is separate from the "real" call of
// QueueSvcMarkDone so it is easy to test.   This method is used for marking
// messages as having been processed successfully. Note that if you do not
// mark a message as "done" you can receive the message additional times
// on future calls to Receive().  This implies that it is best to fully
// process a message in a single function call and mark it done when completed
// even if that implies an error was generated.  If you must process a
// message for more time than a single call, you'll need to hold state in
// a database or similar so you can resume processing at the right point.
func (q *queueSvcImpl) QueueSvcMarkDoneImpl(req *queuemsg.MarkDoneRequest, resp *queuemsg.MarkDoneResponse) queue.QueueErrId {
	ctx := pcontext.CallTo(q.ctx, "QueueSvcMarkDoneImpl")

	// xxx fixme(iansmith) sqlc doesn't seem to understand UPDATE FROM so I have
	// xxx to do this in two steps.  It's not really dangerous because the
	// xxx queue keys are written once and left until deleted
	qid := queue.MustUnmarshalQueueId(req.GetId())
	rowId, errId := q.getRowidForId(ctx, qid)
	if errId.IsError() {
		return errId
	}
	// should be inside a transaction: xxx fixme(iansmith)
	for _, m := range req.Msg {
		mid := queue.MustUnmarshalQueueMsgId(m)
		p := MarkDoneParams{
			QueueKey: sql.NullInt64{Int64: int64(rowId.Low()), Valid: true},
			IDLow:    sql.NullInt64{Int64: int64(mid.Low()), Valid: true},
			IDHigh:   sql.NullInt64{Int64: int64(mid.High()), Valid: true},
		}
		err := q.queries.MarkDone(context.Background(), p)
		if err != nil {
			pcontext.Errorf(ctx, "failed mark done: on queue:%s: %s", qid.Short(), err.Error())
			return internalErr
		}
	}
	return queue.QueueErrIdNoErr
}

// QueueSvcLocateImpl is separate from the "real" call of
// QueueSvcLocate so it is easy to test.  The return values
// will be nil, "" if there was no error.  If there was an error you'll
// get the error code and detail, typically QueueErrNotFound.  If things
// are ok, this returns the queue id for a given name in the response.
func (q queueSvcImpl) QueueSvcLocateImpl(req *queuemsg.LocateRequest, resp *queuemsg.LocateResponse) queue.QueueErrId {
	ctx := pcontext.CallTo(q.ctx, "QueueSvcLocateImpl")
	row, err := q.queries.Locate(ctx, req.QueueName)
	if err != nil {
		pcontext.Errorf(ctx, "error trying to locate queue %s: %s", req.QueueName, err.Error())
		// QueueNotFound XXXX fix me! why not constant
		return queue.NewQueueErrId(7)
	}
	h := uint64(row.IDHigh.Int64)
	l := uint64(row.IDLow.Int64)
	// type here is just for validation of the char
	qid := queue.QueueId(id.NewIdTyped[queue.DefQueue](h, l))
	resp.Id = qid.Marshal()
	return queue.QueueErrIdNoErr
}

// QueueSvcReceiveImpl is separate from the "real" call of
// QueueSvcReceive so it is easy to test.  The return values
// will be nil, "" if there was no error.  If there was an error you'll
// get the error code and detail.  This code will return some number of messages
// from zero to the requested maximum.  If the requested maximum is out of bounds
// it will be clipped to the range [1,4).  1 is the recommended value, and since
// the max is an integer type with a default of 0, it will be clipped to 1.
// It returns the content of the messages that are pending, in approximately the
// order sent, although this is not guaranteed.  Just retreiving messages is not
// enough to fully process them, you need to use QueueSvcMarkDone() to indicate that the
// item can be removed from the queue.
func (q *queueSvcImpl) QueueSvcReceiveImpl(req *queuemsg.ReceiveRequest, resp *queuemsg.ReceiveResponse) queue.QueueErrId {
	ctx := pcontext.CallTo(q.ctx, "QueueSvcReceiveImpl")
	qid := queue.MustUnmarshalQueueId(req.GetId())
	// we have to do this because we can't do UPDATE FROM in the sql queries with sqlc
	rowId, errId := q.getRowidForId(ctx, qid)
	if errId.IsError() {
		pcontext.Errorf(ctx, "error trying to find row for queue %s: %s", qid.Short(), errId.Short())
		return errId
	}

	p := RetrieveMessageParams{
		IDHigh: sql.NullInt64{Int64: int64(qid.High()), Valid: true},
		IDLow:  sql.NullInt64{Int64: int64(qid.Low()), Valid: true},
	}
	resultMsg, err := q.queries.RetrieveMessage(context.Background(), p)
	if err != nil {
		pcontext.Errorf(ctx, "error trying retreive messages in send on queue %s: %s", qid.Short(), errId.Short())
		return internalErr
	}
	if len(resultMsg) == 0 {
		resp.Id = req.GetId()
		resp.Message = nil
		return queue.QueueErrIdNoErr
	}
	max := int(req.GetMessageLimit())
	if max < 1 {
		max = 1
	}
	if max > 3 {
		max = 3
	}
	if len(resultMsg) < max {
		max = len(resultMsg)
	}
	resultList := make([]*queuemsg.QueueMsg, max)
	n := max
	if len(resultMsg) < max {
		n = len(resultMsg)
	}
	// xxx fixme(iansmith): this should be in a transaction
	for i := 0; i < n; i++ {
		result := resultMsg[i]
		var recv time.Time
		if result.LastReceived.String != "" {
			recv, err = time.Parse(longForm, result.LastReceived.String)
			if err != nil {
				pcontext.Errorf(ctx, "failed to understand lastReceived time: %v", err)
				return internalErr
			}
		}
		sent, err := time.Parse(longForm, result.OriginalSent.String)
		if err != nil {
			pcontext.Errorf(ctx, "failed to understand lastReceived time: %v", err)
			return internalErr
		}
		var senderAny, payloadAny anypb.Any
		err = proto.Unmarshal(resultMsg[i].Sender, &senderAny)
		if err != nil {
			pcontext.Errorf(ctx, "unable to create sender proto: %v", err)
			return internalErr
		}
		err = proto.Unmarshal(resultMsg[i].Payload, &payloadAny)
		if err != nil {
			pcontext.Errorf(ctx, "unable to create payload proto: %v", err)
			return queue.NewQueueErrId(queue.QueueErrIdUnmarshalError)
		}

		messageId := queue.QueueMsgId(id.NewIdTyped[queue.DefQueueMsg](uint64(result.IDHigh.Int64), uint64(result.IDLow.Int64)))
		m := queuemsg.QueueMsg{
			Id:           req.GetId(),
			MsgId:        messageId.Marshal(),
			ReceiveCount: int32(result.ReceivedCount.Int64),
			Received:     timestamppb.New(recv),
			Sent:         timestamppb.New(sent),
			Sender:       &senderAny,
			Payload:      &payloadAny,
		}
		resultList[i] = &m
		up := UpdateMessageRetrievedParams{
			QueueKey: sql.NullInt64{Int64: int64(rowId.Low()), Valid: true},
			IDLow:    result.IDLow,
			IDHigh:   result.IDHigh,
		}
		err = q.queries.UpdateMessageRetrieved(context.Background(), up)
		if err != nil {
			pcontext.Errorf(ctx, "unable to update message data %s: %s", messageId.Short(),
				err.Error())
			return internalErr
		}
	}
	resp.Id = req.GetId()
	resp.Message = resultList
	return queue.QueueErrIdNoErr
}
