package main

import (
	"context"
	"testing"

	"github.com/iansmith/parigot/apishared/id"
	"github.com/iansmith/parigot/apiwasm/syscall"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/test/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/lib/go/future"

	bar "github.com/iansmith/parigot/g/methodcall/bar/v1"
	foo "github.com/iansmith/parigot/g/methodcall/foo/v1"
	methodcall "github.com/iansmith/parigot/g/methodcall/v1"

	const_ "github.com/iansmith/parigot/test/func/methodcall/impl/foo/const_"
)

var exitCode = int32(0)

func manufactureContext(name string) context.Context {
	return pcontext.NewContextWithContainer(pcontext.CallTo(pcontext.GuestContext(pcontext.NewContextWithContainer(context.Background(), "methodcall.Main")), name), name)
}

var myServiceId id.ServiceId
var fooClient foo.Client
var barClient bar.Client

func main() {

	ctx := manufactureContext("[methodcall]main")
	defer func() {
		if r := recover(); r != nil {
			pcontext.Errorf(ctx, "methodcall: trapped a panic in the guest side: %v", r)
		}
		pcontext.Dump(ctx)
	}()
	pcontext.Debugf(ctx, "program started")
	myServiceId = lib.MustRegisterClient(ctx)
	foo.MustRequire(pcontext.CallTo(ctx, "Require"), myServiceId)
	test.MustRequireTest(ctx, myServiceId)
	bar.MustRequire(ctx, myServiceId)

	syscall.MustSatisfyWait(ctx, myServiceId)

	underTestServer.testSvc = test.MustLocateTest(ctx, myServiceId)
	underTestServer.fooClient = foo.MustLocate(ctx, myServiceId)
	underTestServer.barClient = bar.MustLocate(ctx, myServiceId)

	baseBool := test.LaunchUnderTest(ctx, myServiceId, underTestServer)
	baseBool.Handle(func(ok bool) {
		if !ok {
			pcontext.Errorf(ctx, "unable to start the under test service")
			return
		}
		pcontext.Infof(ctx, "LaunchUnderTest succeeded.")
	})
}

// TestAddMulitply is a test of a function that has both input and output.
func (m *myUnderTestServer) TestAddMultiply(ctx context.Context, t *testing.T) {
	fn := func(t *testing.T, value0, value1, sum, product int32) {
		req := &foo.AddMultiplyRequest{
			Value0: value0,
			Value1: value1,
			IsAdd:  true,
		}
		addMultFuture := m.fooClient.AddMultiply(ctx, req)
		addMultFuture.Failure(func(err foo.FooErr) {
			t.Errorf("error in AddMultiply [add]:%v", foo.FooErr_name[int32(err)])
		})
		addMultFuture.Success(func(resp *foo.AddMultiplyResponse) {
			if resp.Result != sum {
				t.Errorf("bad result for add, expected %d but got %d", sum, resp.Result)
			}
		})
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

		req := bar.AccumulateRequest{
			Value: rest,
		}
		futureAcc := m.barClient.Accumulate(ctx, &req)
		futureAcc.Failure(func(err bar.BarErr) {
			t.Errorf("received error from call to Accumulate: %v", methodcall.MethodCallSuiteErr_name[int32(err)])
			t.FailNow()
		})
		futureAcc.Success(func(resp *bar.AccumulateResponse) {
			if resp.GetProduct() != prod {
				t.Errorf("received wrong result from Accumulate: expected prod %d but got %d", prod, resp.GetProduct())
			}
			if resp.GetSum() != sum {
				t.Errorf("received wrong result from  Accumulate: expected sum %d but got %d", sum, resp.GetSum())
			}
		})
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

	lucasFuture := m.fooClient.LucasSequence(ctx)
	lucasFuture.Failure(func(err foo.FooErr) {
		t.Errorf("received error from call to LucasSequence: %v", methodcall.MethodCallSuiteErr_name[int32(err)])
	})
	lucasFuture.Success(func(result *foo.LucasSequenceResponse) {
		member := result.GetSequence()[const_.LucasSize]
		if member != 141422324 {
			t.Logf("outside func f2")
			t.Logf("unexpected value in lucas sequence (index %d): got %d but expected %d",
				const_.LucasSize-1, member, 141422324)
		}
	})
}

/////////////////
// UNDER TEST
////////////////

// UnderTest
var underTestServer = &myUnderTestServer{}

type myUnderTestServer struct {
	testSvc   test.ClientTest
	fooClient foo.Client
	barClient bar.Client
}

// Ready is a check, if this returns false the library will abort and not attempt to run this service.
// Normally this is used to do LocateXXX() calls that are needed for
// the operation of the service.
func (m *myUnderTestServer) Ready(ctx context.Context, sid id.ServiceId) *future.Base[bool] {
	m.fooClient = foo.MustLocate(ctx, myServiceId)
	m.barClient = bar.MustLocate(ctx, myServiceId)
	m.testSvc = test.MustLocateTest(ctx, myServiceId)
	setupFuture := m.setupTests(ctx)
	ready := future.NewBase[bool]()
	setupFuture.Success(func(_ *test.AddTestSuiteResponse) {
		pcontext.Infof(ctx, "test setup ok")
		ready.Set(true)
	})
	setupFuture.Failure(func(err test.TestErr) {
		pcontext.Infof(ctx, "test setup failed:%s", test.TestErr_name[int32(err)])
		ready.Set(false)
	})
	return ready
}
func (m *myUnderTestServer) Exec(ctx context.Context, req *test.ExecRequest) (*test.ExecResponse, test.TestErr) {
	pcontext.Debugf(ctx, "Exec", "got an exec call %s.%s.%s", req.GetPackage(), req.GetService(), req.GetName())
	resp := &test.ExecResponse{}
	return resp, test.TestErr_NoError
}

func (m *myUnderTestServer) setupTests(ctx context.Context) *test.FutureTestAddTestSuite {
	pcontext.Debugf(ctx, "setupTests reached in under test server")

	addReq := &test.AddTestSuiteRequest{
		Suite: []*test.SuiteInfo{
			{
				PackagePath:  "methodcall",
				Service:      "Main",
				FunctionName: []string{"TestAddMultiply"},
			},
		},
		ExecPackage: "test.v1",
		ExecService: "UnderTestService",
	}
	tsFuture := m.testSvc.TestAddTestSuite(ctx, addReq)
	tsFuture.Failure(func(err test.TestErr) {
		pcontext.Infof(ctx, "AddTestSuite:failed to add test suite %s", test.TestErr_name[int32(err)])
	})
	tsFuture.Success(func(resp *test.AddTestSuiteResponse) {
		pcontext.Infof(ctx, "AddTestSuite success: %+v",
			resp.GetSucceeded())

		startFuture := m.testSvc.TestStart(ctx, &test.StartRequest{})
		startFuture.Failure(func(err test.TestErr) {
			pcontext.Logf(ctx, pcontext.Error, "testSvc.Start():%s", test.TestErr_name[int32(err)])
		})
		startFuture.Success(func(resp *test.StartResponse) {
			if resp.GetRegexFailed() {
				pcontext.Logf(ctx, pcontext.Error, "Regexp Failed in filter")
			}
			pcontext.Logf(ctx, pcontext.Info, "Start() success: started %v tests", resp.GetNumTest())
		})
	})
	return tsFuture
}
