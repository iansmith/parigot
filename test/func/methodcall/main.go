package main

import (
	"context"
	"testing"
	"time"

	pcontext "github.com/iansmith/parigot/context"
	methg "github.com/iansmith/parigot/g/methodcall/v1"
	methodcallmsg "github.com/iansmith/parigot/g/msg/methodcall/v1"
	testmsg "github.com/iansmith/parigot/g/msg/test/v1"
	"github.com/iansmith/parigot/g/test/v1"
	const_ "github.com/iansmith/parigot/test/func/methodcall/impl/foo/const_"
)

var exitCode = int32(0)

const TestRegexpFail test.TestErrIdCode = test.TestErrIdGuestStart

func manufactureContext(name string) context.Context {
	return pcontext.NewContextWithContainer(pcontext.CallTo(pcontext.ClientContext(context.Background()), name), name)
}

func main() {
	ctx := manufactureContext("[methodcall]main")
	defer pcontext.Dump(ctx)
	// pcontext.Debugf(ctx, "main", "xxx main 2222 of methodcall test")
	// methg.RequireFooServiceOrPanic(ctx)
	// methg.RequireBarServiceOrPanic(ctx)
	// //test.RequireTestServiceOrPanic(bg)
	// queue.RequireQueueServiceOrPanic(ctx)

	// // now get handles to the services
	// methg.LocateBarServiceOrPanic(ctx)
	// methg.LocateFooServiceOrPanic(ctx)
	// queue.LocateQueueServiceOrPanic(ctx)

	for {
		pcontext.Debugf(ctx, "started up: main of methodcall test, waiting....")
		time.Sleep(5 * time.Second)
	}

	test.RunUnderTestService(ctx, underTestServer)
}

// TestAddMulitply is a test of a function that has both input and output.
func (m *myUnderTestServer) TestAddMultiply(ctx context.Context, t *testing.T) {
	fn := func(t *testing.T, value0, value1, sum, product int32) {
		req := &methodcallmsg.AddMultiplyRequest{
			Value0: value0,
			Value1: value1,
			IsAdd:  true,
		}
		resp, err := m.foo.AddMultiply(ctx, req)
		if err.IsError() {
			t.Fatalf("error in AddMultiply [add]:%v", err)
		}
		if resp.Result != sum {
			t.Fatalf("bad result for add, expected %d but got %d", sum, resp.Result)
		}

		req.IsAdd = false
		resp, err = m.foo.AddMultiply(ctx, req)
		if err.IsError() {
			t.Fatalf("error in AddMultiply [mult]:%s", err.Short())
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
		ctx := manufactureContext("TestAccumulate")

		req := methodcallmsg.AccumulateRequest{
			Value: rest,
		}
		resp, err := m.bar.Accumulate(ctx, &req)
		if err.IsError() {
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
	ctx := manufactureContext("[methodcall]TestLucas")

	result, err := m.foo.LucasSequence(ctx)
	if err.IsError() {
		t.Logf("outside func f1")
		t.Errorf("received error from call to LucasSequence: %v", err)
		t.Fail()
	}
	member := result.GetSequence()[const_.LucasSize]
	if member != 141422324 {
		t.Logf("outside func f2")
		t.Logf("unexpected value in lucas sequence (index %d): got %d but expected %d",
			const_.LucasSize-1, member, 141422324)
	}
	t.Logf("outside func")
}

/////////////////
// UNDER TEST
////////////////

// UnderTest
var underTestServer = &myUnderTestServer{}

type myUnderTestServer struct {
	testSvc test.TestServiceClient
	foo     methg.FooServiceClient
	bar     methg.BarServiceClient
}

func (m *myUnderTestServer) Ready(ctx context.Context) bool {
	m.foo = methg.LocateFooServiceOrPanic(ctx)
	m.bar = methg.LocateBarServiceOrPanic(ctx)
	m.testSvc = test.LocateTestServiceOrPanic(ctx)
	if err := m.setupTests(ctx); err.IsError() {
		pcontext.Logf(ctx, pcontext.Error, "test setup failed:", err)
		return false
	}
	return true
}
func (m *myUnderTestServer) Exec(ctx context.Context, req *testmsg.ExecRequest) (*testmsg.ExecResponse, test.TestErrId) {
	pcontext.Debugf(ctx, "Exec", "got an exec call %s.%s.%s", req.GetPackage(), req.GetService(), req.GetName())
	resp := &testmsg.ExecResponse{}
	return resp, test.TestErrIdNoErr
}

func (m *myUnderTestServer) setupTests(ctx context.Context) test.TestErrId {
	pcontext.Debugf(ctx, "setupTests", "setup tests reached")

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
	resp, err := m.testSvc.AddTestSuite(ctx, addReq)
	if err.IsError() {
		pcontext.Logf(ctx, pcontext.Error, "AddTestSuite:%v", err)
		return err
	}
	pcontext.Logf(ctx, pcontext.Info, "AddTestSuite success: %+v", resp.Succeeded)
	startResp, err := m.testSvc.Start(ctx, &testmsg.StartRequest{})
	if err.IsError() {
		pcontext.Logf(ctx, pcontext.Error, "testSvc.Start():%v", err)
		return err
	}
	if startResp.GetRegexFailed() {
		pcontext.Logf(ctx, pcontext.Error, "Regexp Failed in filter:%v", err)
		return test.NewTestErrId(TestRegexpFail)
	}
	pcontext.Logf(ctx, pcontext.Info, "Start() success: started %v tests", startResp.GetNumTest())
	return test.TestErrIdNoErr
}
