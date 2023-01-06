package main

import (
	"fmt"
	"time"

	"github.com/iansmith/parigot/api_impl/syscall"
	"github.com/iansmith/parigot/g/log/v1"
	"github.com/iansmith/parigot/g/methodcall/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var logger log.LogService
var callImpl = syscall.NewCallImpl()
var foo methodcall.FooServiceClient

//go:noinline
func main() {
	//flag.Parse() <--- can't do this until we get startup args figured out

	if _, err := callImpl.Require1("methodcall", "Foo"); err != nil {
		panic("unable to require foo service: " + err.Error())
	}
	if _, err := callImpl.Require1("methodcall", "Bar"); err != nil {
		panic("unable to require bar service: " + err.Error())
	}
	if _, err := callImpl.Require1("log", "Log"); err != nil {
		panic("unable to require log service: " + err.Error())
	}
	if _, err := callImpl.Run(&syscallmsg.RunRequest{Wait: true}); err != nil {
		panic("error starting client process:" + err.Error())
	}

	var err error
	logger, err = log.LocateLogService()
	if err != nil {
		panic(fmt.Sprintf("failed to locate log:%v", err))
	}

}

func Log(spec string, arg ...interface{}) {
	req := &logmsg.LogRequest{
		Stamp:   timestamppb.New(time.Now()), // xxx fix me, should be using the kernel
		Level:   logmsg.LogLevel_LOG_LEVEL_INFO,
		Message: fmt.Sprintf(spec, arg...),
	}
	if logger == nil {
		print("Internal error in methodcall test: logger is nil! "+fmt.Sprintf(spec, arg...), "\n")
		return
	}
	if err := logger.Log(req); err != nil {
		print("methodcall: error in log call:", err.Error(), "\n")
	}
}
