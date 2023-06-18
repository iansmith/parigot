package main

import (
	"context"
	"testing"

	"github.com/iansmith/parigot/apishared/id"
	"github.com/iansmith/parigot/apiwasm/syscall"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/test/v1"
	lib "github.com/iansmith/parigot/lib/go"

	methodcall "github.com/iansmith/parigot/g/methodcall/v1"
	sysg "github.com/iansmith/parigot/g/syscall/v1"

	const_ "github.com/iansmith/parigot/test/func/methodcall/impl/foo/const_"
)

var exitCode = int32(0)

func manufactureContext(name string) context.Context {
	return pcontext.NewContextWithContainer(pcontext.CallTo(pcontext.GuestContext(context.Background()), name), name)
}

var myServiceId id.ServiceId
var foo methodcall.ClientFoo
var bar methodcall.ClientBar

func main() {
	prevMain()
}
func prevMain() {
	ctx := manufactureContext("[methodcall]main")
	defer func() {
		if r := recover(); r != nil {
			pcontext.Errorf(ctx, "methodcall: trapped a panic in the guest side: %v", r)
		}
		pcontext.Dump(ctx)
	}()
	pcontext.Debugf(ctx, "program started")
	myServiceId = lib.MustRegisterClient(ctx)
	methodcall.MustRequireFoo(pcontext.CallTo(ctx, "Require"), myServiceId)
	test.MustRequireTest(ctx, myServiceId)
	methodcall.MustRequireBar(ctx, myServiceId)

	syscall.MustSatisfyWait(ctx, myServiceId)

	underTestServer.testSvc = test.MustLocateTest(ctx, myServiceId)
	underTestServer.foo = methodcall.MustLocateFoo(ctx, myServiceId)
	underTestServer.bar = methodcall.MustLocateBar(ctx, myServiceId)

	kerr := test.LaunchUnderTest(ctx, myServiceId, underTestServer)
	if kerr != sysg.KernelErr_NoError {
		pcontext.Errorf(ctx, "unable to start the under test service: %s", sysg.KernelErr_name[int32(kerr)])
	}
}

// TestAddMulitply is a test of a function that has both input and output.
func (m *myUnderTestServer) TestAddMultiply(ctx context.Context, t *testing.T) {
	fn := func(t *testing.T, value0, value1, sum, product int32) {
		req := &methodcall.AddMultiplyRequest{
			Value0: value0,
			Value1: value1,
			IsAdd:  true,
		}
		resp, err := m.foo.AddMultiply(ctx, req)
		if err != methodcall.MethodCallSuiteErr_NoError {
			t.Fatalf("error in AddMultiply [add]:%v", methodcall.MethodCallSuiteErr_name[int32(err)])
		}
		if resp.Result != sum {
			t.Fatalf("bad result for add, expected %d but got %d", sum, resp.Result)
		}

		req.IsAdd = false
		resp, err = m.foo.AddMultiply(ctx, req)
		if err != methodcall.MethodCallSuiteErr_NoError {
			t.Fatalf("error in AddMultiply [mult]:%s", methodcall.MethodCallSuiteErr_name[int32(err)])
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

		req := methodcall.AccumulateRequest{
			Value: rest,
		}
		resp, err := m.bar.Accumulate(ctx, &req)
		if err != methodcall.MethodCallSuiteErr_NoError {
			t.Errorf("received error from call to Accumulate: %v", methodcall.MethodCallSuiteErr_name[int32(err)])
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
	if err != methodcall.MethodCallSuiteErr_NoError {
		t.Logf("outside func f1")
		t.Errorf("received error from call to LucasSequence: %v", methodcall.MethodCallSuiteErr_name[int32(err)])
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
	testSvc test.ClientTest
	foo     methodcall.ClientFoo
	bar     methodcall.ClientBar
}

// Ready is a check, if this returns false the library will abort and not attempt to run this service.
// Normally this is used to do LocateXXX() calls that are needed for
// the operation of the service.
func (m *myUnderTestServer) Ready(ctx context.Context, sid id.ServiceId) bool {
	m.foo = methodcall.MustLocateFoo(ctx, myServiceId)
	m.bar = methodcall.MustLocateBar(ctx, myServiceId)
	m.testSvc = test.MustLocateTest(ctx, myServiceId)
	if err := m.setupTests(ctx); err != test.TestErr_NoError {
		pcontext.Logf(ctx, pcontext.Error, "test setup failed:%s", test.TestErr_name[int32(err)])
		return false
	}
	return true
}
func (m *myUnderTestServer) Exec(ctx context.Context, req *test.ExecRequest) (*test.ExecResponse, test.TestErr) {
	pcontext.Debugf(ctx, "Exec", "got an exec call %s.%s.%s", req.GetPackage(), req.GetService(), req.GetName())
	resp := &test.ExecResponse{}
	return resp, test.TestErr_NoError
}

func (m *myUnderTestServer) setupTests(ctx context.Context) test.TestErr {
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
	resp, err := m.testSvc.AddTestSuite(ctx, addReq)
	if err != test.TestErr_NoError {
		pcontext.Logf(ctx, pcontext.Error, "AddTestSuite:%s", test.TestErr_name[int32(err)])
		return err
	}
	pcontext.Logf(ctx, pcontext.Info, "AddTestSuite success: %+v", resp.Succeeded)
	startResp, err := m.testSvc.Start(ctx, &test.StartRequest{})
	if err != test.TestErr_NoError {
		pcontext.Logf(ctx, pcontext.Error, "testSvc.Start():%s", test.TestErr_name[int32(err)])
		return err
	}
	if startResp.GetRegexFailed() {
		pcontext.Logf(ctx, pcontext.Error, "Regexp Failed in filter:%v", err)
		return test.TestErr_RegexpFailed
	}
	pcontext.Logf(ctx, pcontext.Info, "Start() success: started %v tests", startResp.GetNumTest())
	return test.TestErr_NoError
}
