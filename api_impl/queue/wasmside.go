package queue

import (
	"errors"
	"fmt"

	"github.com/iansmith/parigot/api_impl/splitutil"
	"github.com/iansmith/parigot/api_impl/syscall"
	"github.com/iansmith/parigot/g/log/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	"github.com/iansmith/parigot/g/queue/v1"

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

func (m *myQueueServer) CreateQueue(_ *protosupportmsg.Pctx, in proto.Message) (proto.Message, error) {
	resp := queuemsg.CreateQueueResponse{}

	// your IDE may become confuse and show an error because of the tricks we are doing to call LogRequestHandler
	errId, err := splitutil.SendReceiveSingleProto(callImpl, in, &resp, go_.QueueSvcCreateQueue)
	if err != nil {
		return nil, err
	}
	if errId != nil {
		return nil, errors.New("internal error:" + errId.Short())
	}
	return &resp, nil

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
