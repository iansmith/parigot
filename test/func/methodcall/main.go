package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/iansmith/parigot/apiimpl/syscall"
	"github.com/iansmith/parigot/g/log/v1"
	"github.com/iansmith/parigot/g/methodcall/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	methodcallmsg "github.com/iansmith/parigot/g/msg/methodcall/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	testmsg "github.com/iansmith/parigot/g/msg/test/v1"
	"github.com/iansmith/parigot/g/test/v1"
	lib "github.com/iansmith/parigot/lib/go"
	const_ "github.com/iansmith/parigot/test/func/methodcall/impl/foo/const_"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var callImpl = syscall.NewCallImpl()

var exitCode = int32(0)

func main() {
	lib.FlagParseCreateEnv()

	if _, err := callImpl.Require1("methodcall", "BarService"); err != nil {
		panic("unable to require bar service: " + err.Error())
	}
	if _, err := callImpl.Require1("log", "LogService"); err != nil {
		panic("unable to require log service: " + err.Error())
	}
	if _, err := callImpl.Require1("test", "TestService"); err != nil {
		panic("unable to require log service: " + err.Error())
	}
	if _, err := callImpl.Require1("methodcall", "FooService"); err != nil {
		panic("unable to require foo service: " + err.Error())
	}
	if _, err := callImpl.Export1("test.v1", "UnderTestService"); err != nil {
		panic("unable to require foo service: " + err.Error())
	}
	test.RunUnderTestService(underTestServer)
}

// TestAddMulitply is a test of a function that has both input and output.
func (m *myUnderTestServer) TestAddMultiply(t *testing.T) {
	callImpl.Exit(&syscallmsg.ExitRequest{Code: 0})
	fn := func(t *testing.T, value0, value1, sum, product int32) {
		req := &methodcallmsg.AddMultiplyRequest{
			Value0: value0,
			Value1: value1,
			IsAdd:  true,
		}
		resp, err := m.foo.AddMultiply(req)
		if err != nil {
			t.Fatalf("error in AddMultiply [add]:%v", err)
		}
		if resp.Result != sum {
			t.Fatalf("bad result for add, expected %d but got %d", sum, resp.Result)
		}

		req.IsAdd = false
		resp, err = m.foo.AddMultiply(req)
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
func (m *myUnderTestServer) TestAccumulate(t *testing.T) {

	fn := func(t *testing.T, sum, prod int32, rest ...int32) {
		req := methodcallmsg.AccumulateRequest{
			Value: rest,
		}
		resp, err := m.bar.Accumulate(&req)
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
	// accumulate starts with identity values
	t.Run("accumulate4", func(t *testing.T) {
		fn(t, 44, 44, 44)
	})
}

// TestLucas returns the members for some unspecified number of iterations.
func (m *myUnderTestServer) TestLucas(t *testing.T) {
	t.Logf("outside func early\n")
	result, err := m.foo.LucasSequence()
	if err != nil {
		t.Logf("outside func f1\n")
		t.Errorf("received error from call to LucasSequence: %v", err)
		t.Fail()
	}
	member := result.GetSequence()[const_.LucasSize]
	t.Logf("member inside lucas test %d\n", member)
	m.logger.Log(&logmsg.LogRequest{Stamp: timestamppb.Now(),
		Level:   logmsg.LogLevel_LOG_LEVEL_DEBUG,
		Message: fmt.Sprintf("lucas sequence: %+v (%d)", result.GetSequence(), member),
	})
	if member != 141422324 {
		t.Logf("outside func f2\n")
		t.Logf("unexpected value in lucas sequence (index %d): got %d but expected %d\n",
			const_.LucasSize-1, member, 141422324)
	}
	m.logger.Log(&logmsg.LogRequest{Stamp: timestamppb.Now(),
		Level:   logmsg.LogLevel_LOG_LEVEL_DEBUG,
		Message: fmt.Sprintf("lucas sequence: %+v (%d)", result.GetSequence(), member),
	})
	t.Logf("outside func\n")
}

/////////////////
// UNDER TEST
////////////////

// UnderTest
var underTestServer = &myUnderTestServer{}

type myUnderTestServer struct {
	logger  *log.LogServiceClient
	testSvc *test.TestServiceClient
	foo     *methodcall.FooServiceClient
	bar     *methodcall.BarServiceClient
}

func (m *myUnderTestServer) Ready() bool {
	if _, err := callImpl.Run(&syscallmsg.RunRequest{Wait: true}); err != nil {
		panic("myLogServer: ready: error in attempt to signal Run: " + err.Error())
	}
	// now we need to setup the stuff we need to run tests
	var err error
	log.Printf("READY1")
	m.logger, err = log.LocateLogService()
	if err != nil {
		panic("unable to get logger with locate")
	}
	m.foo, err = methodcall.LocateFooService(m.logger)
	if err != nil {
		m.logError("LocateFooServer:", err)
		return false
	}
	m.bar, err = methodcall.LocateBarService(m.logger)
	if err != nil {
		m.logError("LocateBarServer:", err)
		return false
	}
	m.testSvc, err = test.LocateTestService(m.logger)
	if err != nil {
		m.logError("LocateTestServer:", err)
		return false
	}
	log.Printf("READY2")
	if err := m.setupTests(); err != nil {
		m.logError("test setup failed:", err)
		return false
	}
	return true
}
func (m *myUnderTestServer) Exec(pctx *protosupportmsg.Pctx, inProto proto.Message) (proto.Message, error) {
	req := inProto.(*testmsg.ExecRequest)
	log.Printf("EXEC1")
	m.logInfo(fmt.Sprintf("got an exec call %s.%s.%s", req.GetPackage(), req.GetService(), req.GetName()))
	resp := &testmsg.ExecResponse{}
	return resp, nil
}

func (m *myUnderTestServer) setupTests() error {
	if err := m.logger.Log(&logmsg.LogRequest{
		Stamp:   timestamppb.Now(), // xxx use kernel now()
		Level:   logmsg.LogLevel_LOG_LEVEL_DEBUG,
		Message: "Testing logger is functioning ok.",
	}); err != nil {
		panic("error trying to log in methodcalltest")
	}
	log.Printf("setupTests")

	addReq := &testmsg.AddTestSuiteRequest{
		Suite: []*testmsg.SuiteInfo{
			{
				PackagePath:  "methodcall",
				Service:      "Main",
				FunctionName: []string{"TestAddMultiply"},
			},
		},
		ExecPackage: "test.v1",
		ExecService: "UnderTestService",
	}
	resp, err := m.testSvc.AddTestSuite(addReq)
	if err != nil {
		m.logError("testSvc.AddTestSuite", err)
		return err
	}
	m.logInfo(fmt.Sprintf("AddTestSuite success: %+v", resp.Succeeded))
	startResp, err := m.testSvc.Start(&testmsg.StartRequest{})
	if err != nil {
		m.logError("testSvc.Start()", err)
		return err
	}
	if startResp.GetRegexFailed() {
		m.logError("RegexpFailed in filter", err)
		return fmt.Errorf("RegexpFailed in filter")
	}
	m.logInfo(fmt.Sprintf("Start() success: started %v tests", startResp.GetNumTest()))
	return nil
}

func (m *myUnderTestServer) logError(msg string, err error) {
	req := &logmsg.LogRequest{
		Level:   logmsg.LogLevel_LOG_LEVEL_ERROR,
		Stamp:   timestamppb.Now(),
		Message: fmt.Sprintf("%s:%s", msg, err.Error()),
	}
	m.logger.Log(req)
}

func (m *myUnderTestServer) logInfo(msg string) {
	req := &logmsg.LogRequest{
		Level:   logmsg.LogLevel_LOG_LEVEL_INFO,
		Stamp:   timestamppb.Now(),
		Message: msg,
	}
	m.logger.Log(req)
}
