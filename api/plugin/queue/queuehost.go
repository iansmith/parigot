package queue

import (
	"context"
	"database/sql"
	"regexp"
	"time"

	apiplugin "github.com/iansmith/parigot/api/plugin"
	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/eng"

	"github.com/iansmith/parigot/g/protosupport/v1"
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

var queueSvc *queueSvcImpl

// This impl is part of the GO side, it is not visible in WASM.
type queueSvcImpl struct {
	db      *sql.DB
	queries *Queries
	ctx     context.Context
}

type QueuePlugin struct{}

// This init functions points the host functions at the functions that
// are the ones to a short setup before calling the real implementation.
func (*QueuePlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "queue", "create_queue_", createQueueHost)
	e.AddSupportedFunc(ctx, "queue", "delete_queue_", deleteQueueHost)
	e.AddSupportedFunc(ctx, "queue", "length_", lengthHost)
	e.AddSupportedFunc(ctx, "queue", "locate_", locateHost)
	e.AddSupportedFunc(ctx, "queue", "mark_done_", markDoneHost)
	e.AddSupportedFunc(ctx, "queue", "receive_", receiveHost)
	e.AddSupportedFunc(ctx, "queue", "send_", sendHost)

	var errId queue.QueueErr
	queueSvc, errId = newQueueSvc(ctx)
	if int32(errId) != 0 {
		pcontext.Fatalf(ctx, "QueueSvc: unable to start:%s", queue.QueueErr_name[int32(errId)])
		return false
	}

	return true
}

func hostBase[T proto.Message, U proto.Message](ctx context.Context, fnName string, fn func(context.Context, T, U) int32,
	m api.Module, stack []uint64, req T, resp U) {
	defer func() {
		if r := recover(); r != nil {
			print(">>>>>>>> Trapped recover in set up for   ", fnName, "<<<<<<<<<<\n")
		}
	}()
	apiplugin.InvokeImplFromStack(ctx, fnName, m, stack, fn, req, resp)
}

func createQueueHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &queue.CreateQueueRequest{}
	resp := &queue.CreateQueueResponse{}

	hostBase(ctx, "Create",
		queueSvc.create, m, stack, req, resp)
}

func deleteQueueHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &queue.DeleteQueueRequest{}
	resp := &queue.DeleteQueueResponse{}

	hostBase(ctx, "Delete", queueSvc.delete, m, stack, req, resp)
}
func lengthHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &queue.LengthRequest{}
	resp := &queue.LengthResponse{}

	hostBase(ctx, "Length", queueSvc.length,
		m, stack, req, resp)
}
func locateHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &queue.LocateRequest{}
	resp := &queue.LocateResponse{}

	hostBase(ctx, "Locate", queueSvc.locate,
		m, stack, req, resp)
}
func markDoneHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &queue.MarkDoneRequest{}
	resp := &queue.MarkDoneResponse{}

	hostBase(ctx, "MarkDone", queueSvc.markDone, m, stack, req, resp)
}
func sendHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &queue.SendRequest{}
	resp := &queue.SendResponse{}

	hostBase(ctx, "Send", queueSvc.send, m, stack,
		req, resp)
}
func receiveHost(ctx context.Context, m api.Module, stack []uint64) {
	req := &queue.ReceiveRequest{}
	resp := &queue.ReceiveResponse{}

	hostBase(ctx, "Receive", queueSvc.receive, m, stack,
		req, resp)
}

//
// real implementation starts here
//

// newQueueSvc returns an initialized Queue service.
func newQueueSvc(ctx context.Context) (*queueSvcImpl, queue.QueueErr) {
	newCtx := pcontext.ServerGoContext(ctx)

	var err error
	q := &queueSvcImpl{}
	//open db
	q.db, err = sql.Open("sqlite3", inMemDSN)
	if err != nil {
		return nil, queue.QueueErr_InternalError
	}
	// create tables
	if _, err := q.db.ExecContext(context.Background(), ddl); err != nil {
		return nil, queue.QueueErr_InternalError
	}

	q.queries = New(q.db)
	q.ctx = newCtx
	return q, queue.QueueErr_NoError
}

var legalName = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_\-\.]*$`)

func (q *queueSvcImpl) validateName(name string) bool {
	return legalName.MatchString(name)
}

// create is separate from the "real" call of
// createHost so it is easy to test.
func (q *queueSvcImpl) create(ctx context.Context, req *queue.CreateQueueRequest,
	resp *queue.CreateQueueResponse) int32 {

	if !q.validateName(req.GetQueueName()) {
		pcontext.Errorf(ctx, "queue name is not valid: '%s'", req.GetQueueName())
		return int32(queue.QueueErr_NotFound)
	}

	rawErr := q.testNameExists(ctx, req.GetQueueName())
	if rawErr != queue.QueueErr_NoError {
		pcontext.Errorf(ctx, "unable to test if queue exists: %s: %v", req.GetQueueName(), rawErr.String())
		return int32(rawErr)
	}

	mod, err := q.queries.CreateQueue(context.Background(), req.GetQueueName())
	if err != nil {
		return int32(queue.QueueErr_InternalError)
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
		return int32(queue.QueueErr_InternalError)
	}
	resp.Id = qid.Marshal()
	return int32(queue.QueueErr_NoError)
}

// delete is separate from the "real" call of
// deleteHost so it is easy to test.  The return value
// will be an QueueErrId.
func (q *queueSvcImpl) delete(ctx context.Context, req *queue.DeleteQueueRequest,
	resp *queue.DeleteQueueResponse) int32 {

	qid := queue.UnmarshalQueueId(req.Id)

	r, err := q.queries.getKeyFromQueueId(ctx,
		getKeyFromQueueIdParams{
			IDLow:  sql.NullInt64{Int64: int64(qid.Low()), Valid: true},
			IDHigh: sql.NullInt64{Int64: int64(qid.High()), Valid: true},
		})
	if err != nil {
		pcontext.Errorf(ctx, "unable to find db key for %s, %v (assuming not found):", qid.Short(), err.Error())
		return int32(queue.QueueErr_NotFound)
	}
	err = q.queries.DeleteQueue(context.Background(), r.QueueKey.Int64)
	if err != nil {
		pcontext.Errorf(ctx, "unable to delete row for key %s: %v:", qid.Short(), err.Error())
		return int32(queue.QueueErr_InternalError)
	}
	// just return the internal id, not the row id
	resp.Id = req.GetId()
	return int32(queue.QueueErr_NoError)
}

// getRowidForId is an internal function to convert a queue id
// into a row id that can be used as a key in other queries.  If anything went wrong
// it will return an internal error.
func (q *queueSvcImpl) getRowidForId(ctx context.Context, qid queue.QueueId) (queue.RowId, queue.QueueErr) {

	p := getKeyFromQueueIdParams{
		IDLow:  sql.NullInt64{Int64: int64(qid.Low()), Valid: true},
		IDHigh: sql.NullInt64{Int64: int64(qid.High()), Valid: true},
	}
	result, err := q.queries.getKeyFromQueueId(ctx, p)
	if err != nil {
		pcontext.Errorf(ctx, "unable to find internal key from id (%s): %v", qid.Short(), err)
		return queue.RowIdZeroValue(), queue.QueueErr_InternalError
	}
	//?? result.QueueKey.Int64,
	rid := id.NewIdTyped[queue.DefRow](0, uint64(result.QueueKey.Int64))
	return queue.RowId(rid), queue.QueueErr_NoError

}

func (q *queueSvcImpl) testNameExists(ctx context.Context, name string) queue.QueueErr {
	count, err := q.queries.TestNameExists(ctx, name)
	if err != nil {
		pcontext.Errorf(ctx, "unable to query table of queue names: %v", err)
		return queue.QueueErr_InternalError
	}
	if count != 0 {
		pcontext.Errorf(ctx, "attempt to create queue 2nd time: %s", name)
		return queue.QueueErr_AlreadyExists
	}
	return queue.QueueErr_NoError
}

// send is separate from the "real" call of
// sendHost so it is easy to test.  If there was an error
// it returns the successfully sent and failed messages. If you are
// trying to send many messages, be sure to look at these two lists
// because partial failure is possible.
func (q *queueSvcImpl) send(ctx context.Context, req *queue.SendRequest, resp *queue.SendResponse) int32 {

	qid := queue.UnmarshalQueueId(req.GetId())

	rid, err := q.getRowidForId(ctx, qid)
	if err != queue.QueueErr_NoError {
		return int32(err)
	}
	succeed := []*protosupport.IdRaw{}
	fail := []*queue.QueueMsg{}

	alreadyFailed := false
	var failedOn *protosupport.IdRaw

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
			return int32(queue.QueueErr_InternalError)
		}
		params := CreateMessageParams{
			IDLow:    sql.NullInt64{Int64: int64(mid.Low()), Valid: true},
			IDHigh:   sql.NullInt64{Int64: int64(mid.High()), Valid: true},
			QueueKey: sql.NullInt64{Int64: int64(rid.Low()), Valid: true},
			Sender:   senderBytes,
			Payload:  payloadBytes,
		}

		/// XXX what about the failed id?
		_, err = q.queries.CreateMessage(ctx, params)
		if err != nil {
			alreadyFailed = true
			fail = append(fail, current)
			failedOn = mid.Marshal()
			pcontext.Errorf(ctx, "could not create message: %v", err)
			continue
		}
		succeed = append(succeed, mid.Marshal())
	}
	if alreadyFailed {
		return int32(queue.QueueErr_InternalError)
	}
	resp.Fail = fail
	resp.Succeed = succeed
	resp.FailedOn = failedOn
	return int32(queue.QueueErr_NoError)
}

// length is separate from the "real" call of
// lengthHost so it is easy to test.  If there was an error you'll
// get a QueueErrId.  If there was no error, the response
// object will have the apporimate length of the queue requested.
func (q *queueSvcImpl) length(ctx context.Context, req *queue.LengthRequest, resp *queue.LengthResponse) int32 {

	qid := queue.UnmarshalQueueId(req.GetId())
	rowId, err := q.getRowidForId(ctx, qid)
	if err != queue.QueueErr_NoError {
		pcontext.Errorf(ctx, "could not find the queue %s: %s", qid.Short(), queue.QueueErr_name[int32(err)])
		return int32(err)
	}

	count, oldErr := q.queries.Length(ctx, sql.NullInt64{Int64: int64(rowId.Low()), Valid: true})
	if oldErr != nil {
		pcontext.Errorf(ctx, "failed length query: on queue:%s: %s", qid.Short(), oldErr.Error())
		return int32(queue.QueueErr_InternalError)
	}
	resp.Id = req.Id
	resp.Length = count
	return int32(queue.QueueErr_NoError)
}

// markDone is separate from the "real" call of
// markDoneHost so it is easy to test.   This method is used for marking
// messages as having been processed successfully. Note that if you do not
// mark a message as "done" you can receive the message additional times
// on future calls to Receive().  This implies that it is best to fully
// process a message in a single function call and mark it done when completed
// even if that implies an error was generated.  If you must process a
// message for more time than a single call, you'll need to hold state in
// a database or similar so you can resume processing at the right point.
func (q *queueSvcImpl) markDone(ctx context.Context, req *queue.MarkDoneRequest, resp *queue.MarkDoneResponse) int32 {

	// xxx fixme(iansmith) sqlc doesn't seem to understand UPDATE FROM so I have
	// xxx to do this in two steps.  It's not really dangerous because the
	// xxx queue keys are written once and left until deleted
	qid := queue.UnmarshalQueueId(req.GetId())
	rowId, err := q.getRowidForId(ctx, qid)
	if err != queue.QueueErr_NoError {
		return int32(err)
	}
	last := -1
	// should be inside a transaction: xxx fixme(iansmith)
	for i, m := range req.Msg {
		mid := queue.UnmarshalQueueMsgId(m)
		p := MarkDoneParams{
			QueueKey: sql.NullInt64{Int64: int64(rowId.Low()), Valid: true},
			IDLow:    sql.NullInt64{Int64: int64(mid.Low()), Valid: true},
			IDHigh:   sql.NullInt64{Int64: int64(mid.High()), Valid: true},
		}
		err := q.queries.MarkDone(ctx, p)
		if err != nil {
			pcontext.Errorf(ctx, "failed mark done: on queue:%s: %s", qid.Short(), err.Error())
			last = i
			break
		}
	}
	if last == -1 || len(req.Msg) == 0 { //everything went through went through
		resp.Unmodified = nil
	} else {
		resp.Unmodified = req.Msg[last:]
	}
	resp.Id = qid.Marshal()
	return int32(queue.QueueErr_NoError)
}

// locate is separate from the "real" call of
// locateHost so it is easy to test.   If there was an error you'll
// get the QueueErrId, typically NotFound.  If things
// are ok, this returns the queue id for a given name in the response.
func (q queueSvcImpl) locate(ctx context.Context, req *queue.LocateRequest, resp *queue.LocateResponse) int32 {
	row, err := q.queries.Locate(ctx, req.QueueName)
	if err != nil {
		pcontext.Errorf(ctx, "error trying to locate queue %s: %s", req.QueueName, err.Error())
		return int32(queue.QueueErr_NotFound)
	}
	h := uint64(row.IDHigh.Int64)
	l := uint64(row.IDLow.Int64)
	// type here is just for validation of the char
	qid := queue.QueueId(id.NewIdTyped[queue.DefQueue](h, l))
	resp.Id = qid.Marshal()
	return int32(queue.QueueErr_NoError)
}

// receive is separate from the "real" call of
// receiveHost so it is easy to test. If there was an error you'll
// get the queue.QueueErr.  This code will return some number of messages
// from zero to the requested maximum.  If the requested maximum is out of bounds
// it will be clipped to the range [1,4).  1 is the recommended value, and since
// the max is an integer type with a default of 0, it will be clipped to 1.
// It returns the content of the messages that are pending, in approximately the
// order sent, although this is not guaranteed.  Just retreiving messages is not
// enough to fully process them, you need to use markDone() to indicate that the
// item can be removed from the queue.
func (q *queueSvcImpl) receive(ctx context.Context, req *queue.ReceiveRequest, resp *queue.ReceiveResponse) int32 {
	qid := queue.UnmarshalQueueId(req.GetId())
	// we have to do this because we can't do UPDATE FROM in the sql queries with sqlc
	rowId, errId := q.getRowidForId(ctx, qid)
	if errId != queue.QueueErr_NoError {
		pcontext.Errorf(ctx, "error trying to find row for queue %s: %s", qid.Short(), queue.QueueErr_name[int32(errId)])
		return int32(errId)
	}

	p := RetrieveMessageParams{
		IDHigh: sql.NullInt64{Int64: int64(qid.High()), Valid: true},
		IDLow:  sql.NullInt64{Int64: int64(qid.Low()), Valid: true},
	}
	resultMsg, err := q.queries.RetrieveMessage(context.Background(), p)
	if err != nil {
		pcontext.Errorf(ctx, "error trying retreive messages in send on queue %s: %s (db error:%s)", qid.Short(),
			queue.QueueErr_name[int32(errId)], err.Error())
		return int32(queue.QueueErr_InternalError)
	}
	if len(resultMsg) == 0 {
		resp.Id = qid.Marshal()
		resp.Message = nil
		return int32(queue.QueueErr_InternalError)
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
	resultList := make([]*queue.QueueMsg, max)
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
				return int32(queue.QueueErr_InternalError)
			}
		}
		sent, err := time.Parse(longForm, result.OriginalSent.String)
		if err != nil {
			pcontext.Errorf(ctx, "failed to understand lastReceived time: %v", err)
			return int32(queue.QueueErr_InternalError)
		}
		var senderAny, payloadAny anypb.Any
		err = proto.Unmarshal(resultMsg[i].Sender, &senderAny)
		if err != nil {
			pcontext.Errorf(ctx, "unable to create sender proto: %v", err)
			return int32(queue.QueueErr_InternalError)
		}
		err = proto.Unmarshal(resultMsg[i].Payload, &payloadAny)
		if err != nil {
			pcontext.Errorf(ctx, "unable to create payload proto: %v", err)
			return int32(queue.QueueErr_UnmarshalFailed)
		}

		high := uint64(resultMsg[0].IDHigh.Int64)
		low := uint64(resultMsg[0].IDLow.Int64)

		messageId := queue.QueueIdFromPair(high, low)
		m := queue.QueueMsg{
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
			return int32(queue.QueueErr_InternalError)
		}
	}
	resp.Id = qid.Marshal()
	resp.Message = resultList
	return int32(queue.QueueErr_NoError)
}
