package main

import (
	"fmt"
	"testing"

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
		Name: "TestAddMultiply",
		F:    TestAddMultiply,
	},
	{
		Name: "TestAccumulate",
		F:    TestAccumulate,
	},
}

func main() {
	lib.FlagParseCreateEnv()
	// print("xxx -- os.Getenv ", os.Getenv("PARIGOT_ENV"), "\n")
	// print("xxx -- lib.Getenv ", lib.Getenv("PARIGOT_ENV"), "\n")
	testing.Init()

	if _, err := callImpl.Require1("methodcall", "BarService"); err != nil {
		panic("unable to require bar service: " + err.Error())
	}
	if _, err := callImpl.Require1("log", "LogService"); err != nil {
		panic("unable to require log service: " + err.Error())
	}
	if _, err := callImpl.Require1("methodcall", "FooService"); err != nil {
		panic("unable to require foo service: " + err.Error())
	}
	print("zzz  in methodcall test, about to run()\n")
	if _, err := callImpl.Run(&syscallmsg.RunRequest{Wait: true}); err != nil {
		panic("error starting client process:" + err.Error())
	}
	var err error
	logger, err = log.LocateLogService()
	if err != nil {
		panic(fmt.Sprintf("failed to locate LogService:%v", err))
	}
	foo, err = methodcall.LocateFooService(logger)
	if err != nil {
		panic(fmt.Sprintf("failed to locate FooService:%v", err))
	}
	bar, err = methodcall.LocateBarService(logger)
	if err != nil {
		panic(fmt.Sprintf("failed to locate BarService:%v", err))
	}

	if err := logger.Log(&logmsg.LogRequest{
		Stamp:   timestamppb.Now(), // xxx use kernel now()
		Level:   logmsg.LogLevel_LOG_LEVEL_DEBUG,
		Message: "Testing logger isfunctioning ok.",
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

// TestAddMulitply is a test of a function that has both input and output.
func TestAddMultiply(t *testing.T) {
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
		t.Run("accumulate", func(t *testing.T) {
		})
	}
	t.Run("addMultiply1", func(t *testing.T) {
		fn(t, 3, 7, 10, 21)
	})
	t.Run("addMultiply2", func(t *testing.T) {
		fn(t, 0, 1, 1, 0)
	})
	t.Run("addMultiply3", func(t *testing.T) {
		fn(t, -10, 20, 10, -200)
	})
	t.Logf("xxx add multiply test-- done")
}

// TestAccumulate returns the sum and product of a list of integers.  The implementation
// inside BarService though actually calls AddMultiply on FooService to do the work. Thus
// this is a test that a service call can itself be implemented built on other service calls.
func TestAccumulate(t *testing.T) {

	fn := func(t *testing.T, sum, prod int32, rest ...int32) {
		req := methodcallmsg.AccumulateRequest{
			Value: rest,
		}
		resp, err := bar.Accumulate(&req)
		if err != nil {
			t.Errorf("received error from call to Accumulate: %v", err)
			t.FailNow()
		}
		if resp.GetProduct() != prod {
			t.Errorf("received wrong result from Accumulate: expected prod %d but got %d", prod, resp.GetProduct())
		}
		if resp.GetSum() != sum {
			t.Errorf("received wrong result from  Accumulate: expected sum %d but got %d", sum, resp.GetSum())
		}
	}
	t.Run("accumulate1", func(t *testing.T) {
		fn(t, 21, 720, 1, 2, 3, 4, 5, 6)
	})
	t.Run("accumulate2", func(t *testing.T) {
		fn(t, -1, -1, -1, 1, -1, 1, -1)
	})
	t.Run("accumulate3", func(t *testing.T) {
		fn(t, 0, 0)
	})
}
