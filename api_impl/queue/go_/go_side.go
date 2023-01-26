package go_

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"

	"github.com/iansmith/parigot/api_impl/splitutil"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/sys/jspatch"

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
	log.Printf("params are %+v", param)
	_, err = q.queries.CreateIdToKeyMapping(context.Background(), param)
	if err != nil {
		return lib.NewQueueError(lib.QueueInternalError), fmt.Sprintf("insert error: %v", err)
	}
	resp.Id = lib.Marshal[*protosupportmsg.QueueId](qid)
	return nil, ""
}

// QueueSvcCreateQueue is the wrapper around the true implementation that
// knows how to pull requests and return results/errors to the WASM side.
func (q *QueueSvcImpl) QueueSvcCreateQueue(sp int32) {
	req := queuemsg.CreateQueueRequest{}
	err := splitutil.StackPointerToRequest(q.mem, sp, &req)
	if err != nil {
		splitutil.ErrorResponse(q.mem, sp, lib.NewQueueError(lib.QueueInternalError),
			fmt.Sprintf("unable to decode request from stack: %v", err))
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
}
