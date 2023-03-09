package queue

import (
	"fmt"

	"github.com/iansmith/parigot/apiimpl/queue/go_"
	"github.com/iansmith/parigot/apiimpl/splitutil"
	"github.com/iansmith/parigot/apiimpl/syscall"
	"github.com/iansmith/parigot/g/log/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	"github.com/iansmith/parigot/g/queue/v1"
	lib "github.com/iansmith/parigot/lib/go"

	"google.golang.org/protobuf/proto"
)

var callImpl = syscall.NewCallImpl()

func main() {
	// we export and require services before the call to file.Run()... our call to the Run() system call is in Ready()
	if _, err := callImpl.Export1("file", "FileService"); err != nil {
		panic("ready: error in attempt to export api.Log: " + err.Error())
	}
	if _, err := callImpl.Require1("log", "LogService"); err != nil {
		panic("ready: error in attempt to export api.Log: " + err.Error())
	}

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
	if _, err := callImpl.Run(&syscallmsg.RunRequest{Wait: true}); err != nil {
		print("ready: error in attempt to signal Run from Queue service: ", err.Error(), "\n")
		return false
	}
	var err error
	m.logger, err = log.LocateLogService()
	if err != nil {
		panic("unable to locate the log:" + err.Error())
	}
	return true
}

func simpleRequestResponse(in, out proto.Message, fn func(int32)) (proto.Message, error) {
	// your IDE may become confuse and show an error because of the tricks we are doing to call QueueSvcCreateQueue
	_, errId, errDetail := splitutil.SendReceiveSingleProto(callImpl, in, out, fn)
	if errId != nil {
		return nil, lib.NewPerrorFromId(errDetail, errId)
	}
	return out, nil
}

func (m *myQueueServer) CreateQueue(_ *protosupportmsg.Pctx, in proto.Message) (proto.Message, error) {
	resp := &queuemsg.CreateQueueResponse{}
	// your IDE may become confuse and show an error because of the tricks we are doing to call QueueSvcCreateQueue
	return simpleRequestResponse(in, resp, go_.QueueSvcCreateHandler)
	// errId, err := splitutil.SendReceiveSingleProto(callImpl, in, &resp, go_.QueueSvcCreateQueue)
}

func (m *myQueueServer) DeleteQueue(_ *protosupportmsg.Pctx, in proto.Message) (proto.Message, error) {
	resp := &queuemsg.DeleteQueueResponse{}
	return simpleRequestResponse(in, resp, go_.QueueSvcDeleteHandler)
}

func (m *myQueueServer) Length(_ *protosupportmsg.Pctx, in proto.Message) (proto.Message, error) {
	resp := &queuemsg.LengthResponse{}
	return simpleRequestResponse(in, resp, go_.QueueSvcLengthHandler)
}

func (m *myQueueServer) Locate(_ *protosupportmsg.Pctx, in proto.Message) (proto.Message, error) {
	resp := &queuemsg.LocateResponse{}
	return simpleRequestResponse(in, resp, go_.QueueSvcLocateHandler)
}

func (m *myQueueServer) MarkDone(_ *protosupportmsg.Pctx, in proto.Message) (proto.Message, error) {
	resp := &queuemsg.MarkDoneResponse{}
	return simpleRequestResponse(in, resp, go_.QueueSvcMarkDoneHandler)
}
func (m *myQueueServer) Receive(_ *protosupportmsg.Pctx, in proto.Message) (proto.Message, error) {
	resp := &queuemsg.ReceiveResponse{}
	return simpleRequestResponse(in, resp, go_.QueueSvcReceiveHandler)
}

func (m *myQueueServer) Send(_ *protosupportmsg.Pctx, in proto.Message) (proto.Message, error) {
	resp := &queuemsg.SendResponse{}
	return simpleRequestResponse(in, resp, go_.QueueSvcSendHandler)
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
