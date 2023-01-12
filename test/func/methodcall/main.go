package main

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/iansmith/parigot/api_impl/syscall"
	"github.com/iansmith/parigot/g/log/v1"
	"github.com/iansmith/parigot/g/methodcall/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	methodcallmsg "github.com/iansmith/parigot/g/msg/methodcall/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var logger *log.LogServiceClient
var foo *methodcall.FooServiceClient
var bar *methodcall.BarServiceClient

var callImpl = syscall.NewCallImpl()

var test = []testing.InternalTest{
	{
		Name: "TestFoo",
		F:    TestFoo,
	},
}

func main() {
	lib.FlagParseCreateEnv()
	// print("xxx -- os.Getenv ", os.Getenv("PARIGOT_ENV"), "\n")
	// print("xxx -- lib.Getenv ", lib.Getenv("PARIGOT_ENV"), "\n")
	print("xxx -- arg(0) ", flag.Arg(0), "\n")
	print("xxx -- arg(1) ", flag.Arg(1), "\n")
	testing.Init()
	print("after init\n")

	// if _, err := callImpl.Require1("methodcall", "FooService"); err != nil {
	// 	panic("unable to require foo service: " + err.Error())
	// }
	// if _, err := callImpl.Require1("methodcall", "BarService"); err != nil {
	// 	panic("unable to require bar service: " + err.Error())
	// }
	if _, err := callImpl.Require1("log", "LogService"); err != nil {
		panic("unable to require log service: " + err.Error())
	}
	// print(fmt.Sprintf("xx main of test about to head to callimpl.Run\n"))
	if _, err := callImpl.Run(&syscallmsg.RunRequest{Wait: true}); err != nil {
		panic("error starting client process:" + err.Error())
	}
	print("xxx -- got through run\n")
	var err error
	logger, err = log.LocateLogService()
	if err != nil {
		panic(fmt.Sprintf("failed to locate LogService:%v", err))
	}
	// foo, err = methodcall.LocateFooService(logger)
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to locate FooService:%v", err))
	// }
	// bar, err = methodcall.LocateBarService(logger)
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to locate BarService:%v", err))
	// }

	if err := logger.Log(&logmsg.LogRequest{
		Stamp:   timestamppb.Now(), // xxx use kernel now()
		Level:   logmsg.LogLevel_LOG_LEVEL_DEBUG,
		Message: "Success logging",
	}); err != nil {
		panic("error trying to log in methodcalltest")
	}
	matchFunc := func(pat, str string) (bool, error) {
		print(fmt.Sprintf("match func called with %s and %s\n", pat, str))
		return true, nil
	}
	// run tests
	testing.Main(matchFunc, test, nil, nil)

	// cleanup?
	callImpl.Exit(&syscallmsg.ExitRequest{
		Code: 0,
	})
}

func TestMain(m *testing.M) {
	print("in test main!!!!!!!!!!!!!\n")
}

func TestFoo(t *testing.T) {
	fn := func(t *testing.T, value0, value1, sum, product int32) {
		req := &methodcallmsg.AddMultiplyRequest{
			Value0: value0,
			Value1: value1,
			IsAdd:  true,
		}
		resp, err := foo.AddMultiply(req)
		if err != nil {
			t.Fatalf("error in AddMultiply [add]:%v", err)
		}
		if resp.Result != sum {
			t.Fatalf("bad result for add, expected %d but got %d", sum, resp.Result)
		}

		req.IsAdd = false
		resp, err = foo.AddMultiply(req)
		if err != nil {
			t.Fatalf("error in AddMultiply [mult]:%v", err)
		}
		if resp.Result != product {
			t.Fatalf("bad result for multiply, expected %d but got %d", product, resp.Result)
		}
		testMessage("OK addMultiply(%d,%d)", value0, value1)
	}
	fn(t, 3, 7, 10, 21)
	fn(t, 0, 1, 1, 0)
	fn(t, -10, 20, 10, -200)
}

func testMessage(spec string, arg ...interface{}) {
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
