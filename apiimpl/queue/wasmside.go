package main

import (
	"fmt"
	golog "log"

	"github.com/iansmith/parigot/apiimpl/splitutil"
	"github.com/iansmith/parigot/apiimpl/syscall"
	"github.com/iansmith/parigot/g/log/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
	"github.com/iansmith/parigot/g/queue/v1"
	lib "github.com/iansmith/parigot/lib/go"

	"google.golang.org/protobuf/proto"
)

var callImpl = syscall.NewCallImpl()

type QueueReqPtr interface {
	proto.Message
	*queuemsg.CreateQueueRequest |
		*queuemsg.DeleteQueueRequest |
		*queuemsg.LengthRequest |
		*queuemsg.LocateRequest |
		*queuemsg.MarkDoneRequest |
		*queuemsg.SendRequest |
		*queuemsg.ReceiveRequest
}
type QueueRespPtr interface {
	proto.Message
	*queuemsg.CreateQueueResponse |
		*queuemsg.DeleteQueueResponse |
		*queuemsg.LengthResponse |
		*queuemsg.LocateResponse |
		*queuemsg.MarkDoneResponse |
		*queuemsg.SendResponse |
		*queuemsg.ReceiveResponse
}

//go:linkname QueueServiceCreateHandler
func QueueSvcCreateHandlerL1([]byte, []byte) lib.Id

//go:linkname QueueServiceDeleteHandler
func QueueSvcDeleteHandlerL1([]byte, []byte) lib.Id

// QueueSvcLengthHandler
func QueueSvcLengthHandlerL1([]byte, []byte) lib.Id

// QueueSvcLocateHandler
func QueueSvcLocateHandlerL1([]byte, []byte) lib.Id

// QueueSvcMarkDoneHandler
func QueueSvcMarkDoneHandlerL1([]byte, []byte) lib.Id

// QueueSvcReceiveHandler
func QueueSvcReceiveHandlerL1([]byte, []byte) lib.Id

// QueueSvcSendHandler
func QueueSvcSendHandlerL1([]byte, []byte) lib.Id

//go:export parigot_main
//go:linkname parigot_main
func parigot_main() {
	lib.FlagParseCreateEnv()

	queue.ExportQueueServiceOrPanic()
	log.RequireLogServiceOrPanic()

	queue.RunQueueService(&myQueueServer{})
}

// Note: myQueueServer is the WASM representation for this service, but it has
// Note: little to do because the service is almost entirely implemented on the
// Note: go side.  Almost all the methods of myQueueServer are simply calls to
// Note: simpleRequestResponse which just deals with moving the request to the
// Note: go side, and then pulling the response back here (WASM).

type myQueueServer struct {
	logger log.LogService
}

func (m *myQueueServer) Ready() bool {
	golog.Printf("queue server, ready called")
	queue.WaitQueueServiceOrPanic()
	m.logger = log.LocateLogServiceOrPanic()
	return true
}

func simpleRequestResponse(in, out proto.Message, fn func(int32)) (proto.Message, error) {
	// your IDE may become confuse and show an error because of the tricks we are doing to call QueueSvcCreateQueue
	spayload := splitutil.SendReceiveSingleProto(callImpl, in, out, fn)
	if splitutil.IsErrorInSinglePayload(spayload) {
		return nil, splitutil.NewPerrorFromSinglePayload(spayload)
	}
	return out, nil
}

func QueueCallL2[T QueueReqPtr, U QueueRespPtr](_ *protosupportmsg.Pctx, in T, out U, fn func([]byte, []byte) lib.Id) error {
	inBuffer, err := proto.Marshal(in)
	if err != nil {
		return lib.NewPerrorFromId("QueueCall", lib.NewKernelError(lib.KernelMarshalFailed))
	}
	if len(inBuffer) > int(lib.GetMaxMessageSize()) {
		inBuffer = nil
		return lib.NewPerrorFromId("QueueCall", lib.NewKernelError(lib.KernelDataTooLarge))
	}
	outBuffer := make([]byte, lib.GetMaxMessageSize())
	id := fn(inBuffer, outBuffer)
	if id != nil && lib.IdRepresentsError(id.High(), id.Low()) {
		inBuffer = nil
		return lib.NewPerrorFromId("QueueCall", id)
	}

	err = proto.Unmarshal(outBuffer, out)
	if err != nil {
		inBuffer = nil
		return lib.NewPerrorFromId("QueueCall", lib.NewKernelError(lib.KernelUnmarshalFailed))
	}

	return nil
}

// CreateQueue is an L3 function, it should be machine generated.
func (m *myQueueServer) CreateQueue(pctx *protosupportmsg.Pctx, in_ proto.Message) (proto.Message, error) {
	in := in_.(*queuemsg.CreateQueueRequest)
	out := &queuemsg.CreateQueueResponse{}
	err := QueueCallL2[*queuemsg.CreateQueueRequest, *queuemsg.CreateQueueResponse](pctx, in, out, QueueSvcCreateHandlerL1)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteQueue is an L3 function, it should be machine generated.
func (m *myQueueServer) DeleteQueue(pctx *protosupportmsg.Pctx, in_ proto.Message) (proto.Message, error) {
	in := in_.(*queuemsg.DeleteQueueRequest)
	out := &queuemsg.DeleteQueueResponse{}
	err := QueueCallL2(pctx, in, out, QueueSvcDeleteHandlerL1)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Length is an L3 function, it should be machine generated
func (m *myQueueServer) Length(pctx *protosupportmsg.Pctx, in_ proto.Message) (proto.Message, error) {
	in := in_.(*queuemsg.LengthRequest)
	out := &queuemsg.LengthResponse{}
	err := QueueCallL2(pctx, in, out, QueueSvcDeleteHandlerL1)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Locate is an L3 function, it should be machine generated
func (m *myQueueServer) Locate(pctx *protosupportmsg.Pctx, in_ proto.Message) (proto.Message, error) {
	in := in_.(*queuemsg.LocateRequest)
	out := &queuemsg.LocateResponse{}
	err := QueueCallL2(pctx, in, out, QueueSvcLocateHandlerL1)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (m *myQueueServer) MarkDone(pctx *protosupportmsg.Pctx, in_ proto.Message) (proto.Message, error) {
	in := in_.(*queuemsg.CreateQueueRequest)
	out := &queuemsg.MarkDoneResponse{}
	err := QueueCallL2(pctx, in, out, QueueSvcMarkDoneHandlerL1)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (m *myQueueServer) Receive(pctx *protosupportmsg.Pctx, in_ proto.Message) (proto.Message, error) {
	in := in_.(*queuemsg.ReceiveRequest)
	out := &queuemsg.ReceiveResponse{}
	err := QueueCallL2(pctx, in, out, QueueSvcReceiveHandlerL1)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (m *myQueueServer) Send(pctx *protosupportmsg.Pctx, in_ proto.Message) (proto.Message, error) {
	in := in_.(*queuemsg.SendRequest)
	out := &queuemsg.SendResponse{}
	err := QueueCallL2(pctx, in, out, QueueSvcSendHandlerL1)
	if err != nil {
		return nil, err
	}
	return out, nil
}

//
// LOGGING UTILS
//

func (m *myQueueServer) logError(msg string, err error) {
	m.log(nil, logmsg.LogLevel_LOG_LEVEL_ERROR, msg+"%v", err)
}

func (m *myQueueServer) logInfo(spec string, rest ...interface{}) {
	m.log(nil, logmsg.LogLevel_LOG_LEVEL_INFO, spec, rest...)
}

func (m *myQueueServer) logDebug(spec string, rest ...interface{}) {
	m.log(nil, logmsg.LogLevel_LOG_LEVEL_DEBUG, spec, rest...)
}

func (m *myQueueServer) log(pctx *protosupportmsg.Pctx, level logmsg.LogLevel, spec string, rest ...interface{}) {
	s := fmt.Sprintf(spec, rest...)
	req := logmsg.LogRequest{
		Stamp:   pctx.GetNow(),
		Level:   level,
		Message: "myFileServer:" + s,
	}
	err := m.logger.Log(&req)
	if err != nil {
		print("unable to log ", s, ":", err.Error(), "\n")
	}
}
