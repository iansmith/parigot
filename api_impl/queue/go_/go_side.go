package go_

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/iansmith/parigot/api_impl/splitutil"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/sys/jspatch"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"

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
func (q *QueueSvcImpl) getRowidForId(id *protosupportmsg.QueueId) (int64, lib.Id, string) {
	u := lib.Unmarshal[*protosupportmsg.QueueId](id)

	p := getKeyFromQueueIdParams{
		IDLow:  sql.NullInt64{Int64: int64(u.Low()), Valid: true},
		IDHigh: sql.NullInt64{Int64: int64(u.High()), Valid: true},
	}
	result, err := q.queries.getKeyFromQueueId(context.Background(), p)
	if err != nil {
		return 44, lib.NewQueueError(lib.QueueInternalError), fmt.Sprintf("unable to find internal key from id (%s): %v", u.Short(), err)
	}
	return result.QueueKey.Int64, nil, ""

}
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
			QueueKey: sql.NullInt64{Int64: queueKey, Valid: true},
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

const longForm = "2006-01-02 15:04:05"

func (q *QueueSvcImpl) QueueSvcReceiveImpl(req *queuemsg.ReceiveRequest, resp *queuemsg.ReceiveResponse) (lib.Id, string) {
	queueKey, id, detail := q.getRowidForId(req.GetId())
	if id != nil {
		return id, detail
	}
	resultMsg, err := q.queries.RetrieveMessage(context.Background(), sql.NullInt64{Int64: queueKey, Valid: true})
	if err != nil {
		return lib.NewQueueError(lib.QueueInternalError), fmt.Sprintf("failed retreiving messages from queue: %v", err)
	}
	max := int(req.GetMessageLimit())
	if max < 1 {
		max = 1
	}
	if max > 3 {
		max = 3
	}
	resultList := make([]*queuemsg.QueueMsg, max)

	for i := 0; i < max; i++ {
		result := resultMsg[i]
		var recv time.Time
		if result.LastReceived.String != "" {
			recv, err = time.Parse(longForm, result.LastReceived.String)
			if err != nil {
				return lib.NewQueueError(lib.QueueInternalError), fmt.Sprintf("failed to understand lastReceived time: %v", err)
			}
		}
		sent, err := time.Parse(longForm, result.OriginalSent.String)
		if err != nil {
			return lib.NewQueueError(lib.QueueInternalError), fmt.Sprintf("failed to understand lastReceived time: %v", err)
		}
		m := queuemsg.QueueMsg{
			Id:           req.GetId(),
			MsgId:        lib.Marshal[protosupportmsg.QueueMsgId](lib.NewFrom64BitPair[*protosupportmsg.QueueMsgId](uint64(result.IDHigh.Int64), uint64(result.IDLow.Int64))),
			ReceiveCount: int32(result.ReceivedCount.Int64),
			Received:     timestamppb.New(recv),
			Sent:         timestamppb.New(sent),
		}
		resultList[i] = &m
	}
	resp.Id = req.GetId()
	resp.Message = resultList
	return nil, ""
}

type allReq interface {
	*queuemsg.CreateQueueRequest |
		*queuemsg.DeleteQueueRequest |
		*queuemsg.LengthRequest |
		*queuemsg.LocateRequest |
		*queuemsg.MarkDoneRequest |
		*queuemsg.SendRequest |
		*queuemsg.ReceiveRequest
}
type allResp interface {
	*queuemsg.CreateQueueResponse |
		*queuemsg.DeleteQueueResponse |
		*queuemsg.LengthResponse |
		*queuemsg.LocateResponse |
		*queuemsg.MarkDoneResponse |
		*queuemsg.SendResponse |
		*queuemsg.ReceiveResponse
}

// generic
func queueOp[T allReq, U allResp](q *QueueSvcImpl, sp int32, op func(T, U) (lib.Id, string)) {
	var req T
	var resp U
	errId, errDetail := splitutil.StackPointerToRequest(q.mem, sp, protoreflect.ProtoMessage(req))
	if errId != nil {
		splitutil.ErrorResponse(q.mem, sp, errId, errDetail)
		return
	}
	id, errDetail := op(req, resp)
	if id != nil {
		splitutil.ErrorResponse(q.mem, sp, id, errDetail)
		return
	}
	splitutil.RespondSingleProto(q.mem, sp, protoreflect.ProtoMessage(resp))

}

// WRAPPERS
func (q *QueueSvcImpl) QueueSvcCreateQueue(sp int32) {
	queueOp(q, sp, q.QueueSvcCreateQueueImpl)
}
func (q *QueueSvcImpl) QueueSvcDeleteQueue(sp int32) {
	queueOp(q, sp, q.QueueSvcDeleteQueueImpl)
}
func (q *QueueSvcImpl) QueueSvcLength(sp int32) {
}
func (q *QueueSvcImpl) QueueSvcLocate(sp int32) {
}
func (q *QueueSvcImpl) QueueSvcMarkDone(sp int32) {
}
func (q *QueueSvcImpl) QueueSvcSend(sp int32) {
	queueOp(q, sp, q.QueueSvcSendImpl)
}

func (q *QueueSvcImpl) QueueSvcReceive(sp int32) {
	queueOp(q, sp, q.QueueSvcReceiveImpl)
}
