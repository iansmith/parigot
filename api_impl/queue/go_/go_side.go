package go_

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"github.com/iansmith/parigot/api_impl/splitutil"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/sys/jspatch"
	"google.golang.org/protobuf/proto"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed sqlc/schema.sql
var ddl string

var inMemDSN = "file:test.db?cache=shared&mode=memory"

// This impl is part of the GO side, it is not visible in WASM.
type QueueSvcImpl struct {
	mem     *jspatch.WasmMem //not used when this is being tested
	db      *sql.DB
	queries *Queries
}

var noErr = lib.NewQueueError(lib.QueueNoError)

// NewQueueSvc returns an initialized Queue service and two nils or
// a nil service and two bits of error information.  When running unit
// tests, it is ok to pass nil to this function.
func NewQueueSvc(mem *jspatch.WasmMem) (*QueueSvcImpl, lib.Id, string) {
	q := &QueueSvcImpl{mem: mem}

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
		return lib.NewQueueError(lib.QueueInternalError), fmt.Sprintf("error: %v", err)
	}

	qid := lib.NewQueueId()
	param := CreateIdToKeyMappingParams{
		IDLow:    sql.NullInt64{Int64: int64(qid.Low()), Valid: true},
		IDHigh:   sql.NullInt64{Int64: int64(qid.High()), Valid: true},
		QueueKey: sql.NullInt64{Int64: mod.ID, Valid: true},
	}
	_, err = q.queries.CreateIdToKeyMapping(context.Background(), param)
	if err != nil {
		return lib.NewQueueError(lib.QueueInternalError), fmt.Sprintf("insert error: %v", err)
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
		return lib.NewQueueError(lib.QueueInternalError),
			fmt.Sprintf("unable to find db key for %s: %v",
				lib.Unmarshal(req.GetId()), err)
	}
	err = q.queries.DeleteQueue(context.Background(), r.QueueKey.Int64)
	if err != nil {
		return lib.NewQueueError(lib.QueueInternalError),
			fmt.Sprintf("unable to delete row for key %d: %v",
				r.QueueKey.Int64, err)
	}
	// just return the internal id, not the row id
	resp.Id = req.GetId()
	return nil, ""
}
func (q *QueueSvcImpl) QueueSvcSendImpl(req *queuemsg.SendRequest, resp *queuemsg.SendResponse) (lib.Id, string) {

	u := lib.Unmarshal[*protosupportmsg.QueueId](req.GetId())

	p := getKeyFromQueueIdParams{
		IDLow:  sql.NullInt64{Int64: int64(u.Low()), Valid: true},
		IDHigh: sql.NullInt64{Int64: int64(u.High()), Valid: true},
	}
	result, err := q.queries.getKeyFromQueueId(context.Background(), p)
	if err != nil {
		return lib.NewQueueError(lib.QueueInternalError), fmt.Sprintf("unable to find internal key from id (%s): %v", u.Short(), err)
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
				failedId = lib.NewQueueError(lib.QueueInternalError)
				failedDetail = fmt.Sprintf("unable to flatten sender: %v", err)
				alreadyFailed = true
				fail = append(fail, current)
				continue
			}
		}
		// flatten payload
		payloadBytes, err := proto.Marshal(current.GetPayload())
		if err != nil {
			return lib.NewQueueError(lib.QueueInternalError), fmt.Sprintf("unable to flatten payload: %v", err)
		}
		params := CreateMessageParams{
			IDLow:    sql.NullInt64{Int64: int64(mid.Low()), Valid: true},
			IDHigh:   sql.NullInt64{Int64: int64(mid.High()), Valid: true},
			QueueKey: result.QueueKey,
			Sender:   senderBytes,
			Payload:  payloadBytes,
		}
		_, err = q.queries.CreateMessage(context.Background(), params)
		if err != nil {
			alreadyFailed = true
			fail = append(fail, current)
			failedId = lib.NewQueueError(lib.QueueInternalError)
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

// QueueSvcCreateQueue is the wrapper around the true implementation that
// knows how to pull requests and return results/errors to the WASM side.
func (q *QueueSvcImpl) QueueSvcCreateQueue(sp int32) {
	req := queuemsg.CreateQueueRequest{}
	errId, errDetail := splitutil.StackPointerToRequest(q.mem, sp, &req)
	if errId != nil {
		splitutil.ErrorResponse(q.mem, sp, errId, errDetail)
		return
	}
	resp := queuemsg.CreateQueueResponse{}
	id, errDetail := q.QueueSvcCreateQueueImpl(&req, &resp)
	if id != nil {
		splitutil.ErrorResponse(q.mem, sp, id, errDetail)
		return
	}
	splitutil.RespondSingleProto(q.mem, sp, &resp)
}

func (q *QueueSvcImpl) QueueSvcDeleteQueue(sp int32) {
	req := queuemsg.DeleteQueueRequest{}
	errId, errDetail := splitutil.StackPointerToRequest(q.mem, sp, &req)
	if errId != nil {
		splitutil.ErrorResponse(q.mem, sp, errId, errDetail)
		return
	}
	resp := queuemsg.DeleteQueueResponse{}

	splitutil.RespondSingleProto(q.mem, sp, &resp)

}
func (q *QueueSvcImpl) QueueSvcLength(sp int32) {
}
func (q *QueueSvcImpl) QueueSvcLocate(sp int32) {
}
func (q *QueueSvcImpl) QueueSvcMarkDone(sp int32) {
}
func (q *QueueSvcImpl) QueueSvcSend(sp int32) {
	req := queuemsg.SendRequest{}
	errId, errDetail := splitutil.StackPointerToRequest(q.mem, sp, &req)
	if errId != nil {
		splitutil.ErrorResponse(q.mem, sp, errId, errDetail)
		return
	}
	resp := queuemsg.SendResponse{}
	splitutil.RespondSingleProto(q.mem, sp, &resp)
}
func (q *QueueSvcImpl) QueueSvcReceive(sp int32) {
}
