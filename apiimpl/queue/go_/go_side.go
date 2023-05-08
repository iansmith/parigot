package go_

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/iansmith/parigot/apiimpl/splitutil"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/sys/jspatch"
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

var internalErr = lib.NewQueueError(lib.QueueInternalError)

// This impl is part of the GO side, it is not visible in WASM.
type QueueSvcImpl struct {
	mem     *jspatch.WasmMem //not used when this is being tested
	db      *sql.DB
	queries *Queries
}

// NewQueueSvc returns an initialized Queue service and two nils or
// a nil service and two bits of error information.  When running unit
// tests, it is ok to pass nil to this function.
func NewQueueSvc() (*QueueSvcImpl, lib.Id, string) {
	q := &QueueSvcImpl{mem: nil}

	var err error
	//open db
	q.db, err = sql.Open("sqlite3", inMemDSN)
	if err != nil {
		return nil, lib.NewQueueError(lib.QueueInternalError),
			fmt.Sprintf("unable to create in memory database: %v", err)
	}
	// create tables
	if _, err := q.db.ExecContext(context.Background(), ddl); err != nil {
		return nil, lib.NewQueueError(lib.QueueInternalError),
			fmt.Sprintf("unable to create in memory database tables: %v", err)
	}

	q.queries = New(q.db)

	return q, nil, ""
}

var legalName = regexp.MustCompile(`[A-Za-z0-9_\-\.]+`)

func (q *QueueSvcImpl) validateName(name string) bool {
	return legalName.MatchString(name)
}

func (q *QueueSvcImpl) SetMemPtr(ptr uintptr) {
	q.mem = jspatch.NewWasmMem(ptr)
}

// QueueSvcCreateQueueImpl is separate from the "real" call of
// QueueSvcCreateQueue so it is easy to test.  The return values
// will be nil, "" if there was no error.  If there was an error
// these are ready to be returned to the WASM side.
func (q *QueueSvcImpl) QueueSvcCreateQueueImpl(req *queuemsg.CreateQueueRequest,
	resp *queuemsg.CreateQueueResponse) (lib.Id, string) {

	if !q.validateName(req.GetQueueName()) {
		return lib.NewQueueError(lib.QueueInvalidName),
			fmt.Sprintf("%s is not a valid name, should be a sequenc of only alphanumeric and '_','-','.'", req.GetQueueName())
	}
	mod, err := q.queries.CreateQueue(context.Background(), req.GetQueueName())
	if err != nil {
		return internalErr, fmt.Sprintf("error: %v", err)
	}

	qid := lib.NewQueueId()
	param := CreateIdToKeyMappingParams{
		IDLow:    sql.NullInt64{Int64: int64(qid.Low()), Valid: true},
		IDHigh:   sql.NullInt64{Int64: int64(qid.High()), Valid: true},
		QueueKey: sql.NullInt64{Int64: mod.ID, Valid: true},
	}
	_, err = q.queries.CreateIdToKeyMapping(context.Background(), param)
	if err != nil {
		return internalErr, fmt.Sprintf("insert error: %v", err)
	}
	resp.Id = lib.Marshal[protosupportmsg.QueueId](qid)
	return nil, ""
}

// QueueSvcDeleteQueueImpl is separate from the "real" call of
// QueueSvcDeleteQueue so it is easy to test.  The return values
// will be nil, "" if there was no error.  If there was an error
// these are ready to be returned to the WASM side.
func (q *QueueSvcImpl) QueueSvcDeleteQueueImpl(req *queuemsg.DeleteQueueRequest,
	resp *queuemsg.DeleteQueueResponse) (lib.Id, string) {

	qid := lib.Unmarshal(req.GetId())
	r, err := q.queries.getKeyFromQueueId(context.Background(),
		getKeyFromQueueIdParams{
			IDLow:  sql.NullInt64{Int64: int64(qid.Low()), Valid: true},
			IDHigh: sql.NullInt64{Int64: int64(qid.High()), Valid: true},
		})
	if err != nil {
		return internalErr,
			fmt.Sprintf("unable to find db key for %s: %v",
				lib.Unmarshal(req.GetId()), err)
	}
	err = q.queries.DeleteQueue(context.Background(), r.QueueKey.Int64)
	if err != nil {
		return internalErr,
			fmt.Sprintf("unable to delete row for key %d: %v",
				r.QueueKey.Int64, err)
	}
	// just return the internal id, not the row id
	resp.Id = req.GetId()
	return nil, ""
}

// getRowidForId is an internal function to convert a queue id (in marshaled form)
// into a row id that can be used as a key in other queries.  If anything went wrong
// it will return non-nil and a detail string.
func (q *QueueSvcImpl) getRowidForId(id *protosupportmsg.QueueId) (int64, lib.Id, string) {
	u := lib.Unmarshal(id)

	p := getKeyFromQueueIdParams{
		IDLow:  sql.NullInt64{Int64: int64(u.Low()), Valid: true},
		IDHigh: sql.NullInt64{Int64: int64(u.High()), Valid: true},
	}
	result, err := q.queries.getKeyFromQueueId(context.Background(), p)
	if err != nil {
		return 44, internalErr, fmt.Sprintf("unable to find internal key from id (%s): %v", u.Short(), err)
	}
	return result.QueueKey.Int64, nil, ""

}

// QueueSvcSendImpl is separate from the "real" call of
// QueueSvcSend so it is easy to test.  The return values
// will be nil, "" if there was no error.  If there was an error
// it returns the successfully sent and failed messages. If you are
// trying to send many messages, be sure to look at these two lists
// because partial failure is possible.
func (q *QueueSvcImpl) QueueSvcSendImpl(req *queuemsg.SendRequest, resp *queuemsg.SendResponse) (lib.Id, string) {
	queueKey, id, detail := q.getRowidForId(req.GetId())
	if id != nil {
		return id, detail
	}
	succeed := []*protosupportmsg.QueueMsgId{}
	fail := []*queuemsg.QueueMsg{}

	alreadyFailed := false
	var failedId lib.Id
	var failedDetail string

	for _, current := range req.GetMsg() {
		if alreadyFailed {
			fail = append(fail, current)
			continue
		}
		mid := lib.NewId[*protosupportmsg.QueueMsgId]()
		// flatten sender
		var senderBytes []byte
		var err error
		if current.GetSender() != nil {
			senderBytes, err = proto.Marshal(current.GetSender())
			if err != nil {
				failedId = internalErr
				failedDetail = fmt.Sprintf("unable to flatten sender: %v", err)
				alreadyFailed = true
				fail = append(fail, current)
				continue
			}
		}
		// flatten payload
		payloadBytes, err := proto.Marshal(current.GetPayload())
		if err != nil {
			return internalErr, fmt.Sprintf("unable to flatten payload: %v", err)
		}
		params := CreateMessageParams{
			IDLow:    sql.NullInt64{Int64: int64(mid.Low()), Valid: true},
			IDHigh:   sql.NullInt64{Int64: int64(mid.High()), Valid: true},
			QueueKey: sql.NullInt64{Int64: queueKey, Valid: true},
			Sender:   senderBytes,
			Payload:  payloadBytes,
		}

		_, err = q.queries.CreateMessage(context.Background(), params)
		if err != nil {
			alreadyFailed = true
			fail = append(fail, current)
			failedId = internalErr
			failedDetail = fmt.Sprintf("could not create message: %v", err)
			continue
		}
		qMsgId := lib.Marshal[protosupportmsg.QueueMsgId](mid)
		succeed = append(succeed, qMsgId)
	}
	if alreadyFailed {
		return failedId, failedDetail
	}
	resp.Fail = fail
	resp.Succeed = succeed
	return nil, ""
}

// QueueSvcLengthImpl is separate from the "real" call of
// QueueSvcLength so it is easy to test.  The return values
// will be nil, "" if there was no error.  If there was an error you'll
// get the error code and detail.  If there was no error, the response
// object will have the apporimate length of the queue requested.
func (q *QueueSvcImpl) QueueSvcLengthImpl(req *queuemsg.LengthRequest, resp *queuemsg.LengthResponse) (lib.Id, string) {
	queueKey, id, detail := q.getRowidForId(req.GetId())
	if id != nil {
		return id, detail
	}

	count, err := q.queries.Length(context.Background(), sql.NullInt64{Int64: queueKey, Valid: true})
	if err != nil {
		return internalErr, err.Error()
	}
	resp.Id = req.Id
	resp.Length = count
	return nil, ""
}

// QueueSvcMarkeDoneImpl is separate from the "real" call of
// QueueSvcMarkDone so it is easy to test.  The return values
// will be nil, "" if there was no error.  If there was an error you'll
// get the error code and detail.  This method is used for marking
// messages as having been processed successfully. Note that if you do not
// mark a message as "done" you can receive the message additional times
// on future calls to Receive().  This implies that it is best to fully
// process a message in a single function call and mark it done when completed
// even if that implies an error was generated.  If you must process a
// message for more time than a single call, you'll need to hold state in
// a database or similar so you can resume processing at the right point.
func (q *QueueSvcImpl) QueueSvcMarkDoneImpl(req *queuemsg.MarkDoneRequest, resp *queuemsg.MarkDoneResponse) (lib.Id, string) {
	// xxx fixme(iansmith) sqlc doesn't seem to understand UPDATE FROM so I have
	// xxx to do this in two steps.  It's not really dangerous because the
	// xxx queue keys are written once and left until deleted
	queueKey, id, detail := q.getRowidForId(req.GetId())
	if id != nil {
		return id, detail
	}
	// should be inside a transaction: xxx fixme(iansmith)
	for _, m := range req.Msg {
		mid := lib.Unmarshal(m)
		low := mid.Low()
		high := mid.High()
		p := MarkDoneParams{
			QueueKey: sql.NullInt64{Int64: queueKey, Valid: true},
			IDLow:    sql.NullInt64{Int64: int64(low), Valid: true},
			IDHigh:   sql.NullInt64{Int64: int64(high), Valid: true},
		}
		err := q.queries.MarkDone(context.Background(), p)
		if err != nil {
			return internalErr, err.Error()
		}
	}
	return nil, ""
}

// QueueSvcLocateImpl is separate from the "real" call of
// QueueSvcLocate so it is easy to test.  The return values
// will be nil, "" if there was no error.  If there was an error you'll
// get the error code and detail, typically QueueErrNotFound.  If things
// are ok, this returns the queue id for a given name in the response.
func (q *QueueSvcImpl) QueueSvcLocateImpl(req *queuemsg.LocateRequest, resp *queuemsg.LocateResponse) (lib.Id, string) {
	row, err := q.queries.Locate(context.Background(), req.QueueName)
	if err != nil {
		return lib.NewQueueError(lib.QueueNotFound), err.Error()
	}
	h := uint64(row.IDHigh.Int64)
	l := uint64(row.IDLow.Int64)
	// type here is just for validation of the char
	qid := lib.NewFrom64BitPair[*protosupportmsg.QueueId](h, l)
	resp.Id = lib.Marshal[protosupportmsg.QueueId](qid)
	return nil, ""
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
func (q *QueueSvcImpl) QueueSvcReceiveImpl(req *queuemsg.ReceiveRequest, resp *queuemsg.ReceiveResponse) (lib.Id, string) {
	qid := lib.Unmarshal(req.GetId())
	// we have to do this because we can't do UPDATE FROM in the sql queries with sqlc
	queueKey, id, detail := q.getRowidForId(req.GetId())
	if id != nil {
		return id, detail
	}

	p := RetrieveMessageParams{
		IDHigh: sql.NullInt64{Int64: int64(qid.High()), Valid: true},
		IDLow:  sql.NullInt64{Int64: int64(qid.Low()), Valid: true},
	}
	resultMsg, err := q.queries.RetrieveMessage(context.Background(), p)
	if err != nil {
		return internalErr, fmt.Sprintf("failed retreiving messages from queue: %v", err)
	}
	if len(resultMsg) == 0 {
		resp.Id = req.GetId()
		resp.Message = nil
		return nil, ""
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
				return internalErr, fmt.Sprintf("failed to understand lastReceived time: %v", err)
			}
		}
		sent, err := time.Parse(longForm, result.OriginalSent.String)
		if err != nil {
			return internalErr, fmt.Sprintf("failed to understand lastReceived time: %v", err)
		}
		var senderAny, payloadAny anypb.Any
		err = proto.Unmarshal(resultMsg[i].Sender, &senderAny)
		if err != nil {
			return internalErr, fmt.Sprintf("unable to create sender proto: %v", err)
		}
		err = proto.Unmarshal(resultMsg[i].Payload, &payloadAny)
		if err != nil {
			return internalErr, fmt.Sprintf("unable to create payload proto: %v", err)
		}

		m := queuemsg.QueueMsg{
			Id: req.GetId(),
			MsgId: lib.Marshal[protosupportmsg.QueueMsgId](lib.NewFrom64BitPair[*protosupportmsg.QueueMsgId](
				uint64(result.IDHigh.Int64), uint64(result.IDLow.Int64))),
			ReceiveCount: int32(result.ReceivedCount.Int64),
			Received:     timestamppb.New(recv),
			Sent:         timestamppb.New(sent),
			Sender:       &senderAny,
			Payload:      &payloadAny,
		}
		resultList[i] = &m
		up := UpdateMessageRetrievedParams{
			QueueKey: sql.NullInt64{Int64: queueKey, Valid: true},
			IDLow:    result.IDLow,
			IDHigh:   result.IDHigh,
		}
		err = q.queries.UpdateMessageRetrieved(context.Background(), up)
		if err != nil {
			return internalErr, err.Error()
		}
	}
	resp.Id = req.GetId()
	resp.Message = resultList
	return nil, ""
}

//
// Generics
//

type allReq interface {
	proto.Message
	*queuemsg.CreateQueueRequest |
		*queuemsg.DeleteQueueRequest |
		*queuemsg.LengthRequest |
		*queuemsg.LocateRequest |
		*queuemsg.MarkDoneRequest |
		*queuemsg.SendRequest |
		*queuemsg.ReceiveRequest
}
type allResp interface {
	proto.Message
	*queuemsg.CreateQueueResponse |
		*queuemsg.DeleteQueueResponse |
		*queuemsg.LengthResponse |
		*queuemsg.LocateResponse |
		*queuemsg.MarkDoneResponse |
		*queuemsg.SendResponse |
		*queuemsg.ReceiveResponse
}

// generic
func queueOp[T allReq, U allResp](q *QueueSvcImpl, sp int32, req T, resp U, op func(T, U) (lib.Id, string)) {

	errId, errDetail := splitutil.StackPointerToRequest(q.mem, sp, req)
	if errId != nil {
		splitutil.ErrorResponse(q.mem, sp, errId, errDetail)
		return
	}
	id, errDetail := op(req, resp)
	if id != nil {
		print("QueueOp generate filling in id ", id.String(), " with ", errDetail, "\n")
		splitutil.ErrorResponse(q.mem, sp, id, errDetail)
		return
	}
	splitutil.RespondSingleProto(q.mem, sp, resp)

}

// WRAPPERS
func (q *QueueSvcImpl) QueueSvcCreateQueue(sp int32) {
	queueOp(q, sp, &queuemsg.CreateQueueRequest{}, &queuemsg.CreateQueueResponse{}, q.QueueSvcCreateQueueImpl)
}
func (q *QueueSvcImpl) QueueSvcDeleteQueue(sp int32) {
	queueOp(q, sp, &queuemsg.DeleteQueueRequest{}, &queuemsg.DeleteQueueResponse{}, q.QueueSvcDeleteQueueImpl)
}
func (q *QueueSvcImpl) QueueSvcLength(sp int32) {
	queueOp(q, sp, &queuemsg.LengthRequest{}, &queuemsg.LengthResponse{}, q.QueueSvcLengthImpl)
}
func (q *QueueSvcImpl) QueueSvcLocate(sp int32) {
	queueOp(q, sp, &queuemsg.LocateRequest{}, &queuemsg.LocateResponse{}, q.QueueSvcLocateImpl)
}
func (q *QueueSvcImpl) QueueSvcMarkDone(sp int32) {
	queueOp(q, sp, &queuemsg.MarkDoneRequest{}, &queuemsg.MarkDoneResponse{}, q.QueueSvcMarkDoneImpl)
}
func (q *QueueSvcImpl) QueueSvcSend(sp int32) {
	queueOp(q, sp, &queuemsg.SendRequest{}, &queuemsg.SendResponse{}, q.QueueSvcSendImpl)
}
func (q *QueueSvcImpl) QueueSvcReceive(sp int32) {
	queueOp(q, sp, &queuemsg.ReceiveRequest{}, &queuemsg.ReceiveResponse{}, q.QueueSvcReceiveImpl)
}

// dumpAll is only for debugging the sqlite3 queries, should not
// be used in any real code.
func dumpAll(impl *QueueSvcImpl, queue_key int64) {
	log.Printf("xxx----------------start")
	result, err := impl.queries.allMessages(context.Background(),
		sql.NullInt64{Int64: queue_key, Valid: true})
	if err != nil {
		panic("error trying to dumpAll")
	}
	for i := 0; i < len(result); i++ {
		r := result[i]
		done := "N/A"
		if r.MarkedDone.Valid {
			done = r.MarkedDone.String
		}
		rcvd := "N/A"
		if r.LastReceived.Valid {
			done = r.LastReceived.String
		}

		log.Printf("xxx %d: %x %x %x %s %s %s", i, r.QueueKey.Int64,
			r.IDHigh.Int64, r.IDLow.Int64, r.OriginalSent.String, done, rcvd)
	}
	log.Printf("xxx----------------end")
}
