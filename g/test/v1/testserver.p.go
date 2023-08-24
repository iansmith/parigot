//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: test/v1/test.proto

package test




import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"runtime/debug"
    "unsafe" 
    // this set of imports is _unrelated_ to the particulars of what the .proto imported... those are above
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"  
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/api/shared/id"
	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/lib/go/future"
	"github.com/iansmith/parigot/lib/go/client"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/proto"

)
var _ =  unsafe.Sizeof([]byte{})

 
func LaunchTest(ctx context.Context, sid id.ServiceId, impl Test) *future.Base[bool] {

	readyResult:=future.NewBase[bool]()

	ready:=impl.Ready(ctx,sid)
	ready.Handle(func (b bool) {
		if b {
			readyResult.Set(true)			
			return
		}
		slog.Error("Unable to start test.v1.Test, Ready returned false")
		readyResult.Set(false)
	})

	return readyResult
}

// Note that  InitTest returns a future, but the case of failure is covered
// by this definition so the caller need only deal with Success case.
// The context passed here does not need to contain a logger, one will be created.
func InitTest(require []lib.MustRequireFunc, impl Test) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture, context.Context, id.ServiceId){
	// tricky, this context really should not be used but is
	// passed so as to allow printing if things go wrong
	ctx, myId := MustRegisterTest()
	MustExportTest(context.Background(),myId)
	if len(require)>0 {
		for _, f := range require {
			f(ctx, myId)
		}
	}
	smmap, launchF:=MustLaunchServiceTest(ctx, myId, impl)
	launchF.Failure(func (err syscall.KernelErr) {
		t:=syscall.KernelErr_name[int32(err)]
		slog.Error("launch failure on call Test","error",t)
		lib.ExitSelf(ctx, 1, myId)
	})
	return smmap,launchF, ctx,myId
}
func RunTest(ctx context.Context,
	binding *lib.ServiceMethodMap, timeoutInMillis int32, bg lib.Backgrounder) syscall.KernelErr{
	defer func() {
		if r := recover(); r != nil {
			s, ok:=r.(string)
			if !ok && s!=apishared.ControlledExit {
				slog.Error("RunTest: trapped a panic in the guest side", "recovered", r)
			}
		}
	}()
	var kerr syscall.KernelErr
	for {
		kerr:=ReadOneAndCallTest(ctx, binding, timeoutInMillis)
		if kerr == syscall.KernelErr_ReadOneTimeout {
			if bg==nil {
				continue
			}
			slog.Info("calling backgrounder of Test")
			bg.Background(ctx)
			continue
		}
		if kerr == syscall.KernelErr_NoError {
			continue
		}
		break
	}
	slog.Error("error while waiting for Test service calls", "error",syscall.KernelErr_name[int32(kerr)])
	return kerr
}
// Increase this value at your peril!
// Decreasing this value may make your overall program more responsive if you have many services.
var TimeoutInMillisTest = int32(50)

func ReadOneAndCallTest(ctx context.Context, binding *lib.ServiceMethodMap, 
	timeoutInMillis int32) syscall.KernelErr{
	req:=syscall.ReadOneRequest{}
	hid:= syscallguest.CurrentHostId()

	req.TimeoutInMillis = timeoutInMillis
	req.HostId = hid.Marshal()
	resp, err:=syscallguest.ReadOne(ctx, &req)
	if err!=syscall.KernelErr_NoError {
		return err
	}
	// is timeout?
	if resp.Timeout {
		return syscall.KernelErr_ReadOneTimeout
	}

	// check for finished futures from within our address space
	ctx, t:=lib.CurrentTime(ctx)
	syscallguest.ExpireMethod(ctx,t)

	// is a promise being completed that was fulfilled somewhere else
	if r:=resp.GetResolved(); r!=nil {
		cid:=id.UnmarshalCallId(r.GetCallId())
		defer func() {
			if r:=recover(); r!=nil {
				sid:=id.UnmarshalServiceId(resp.GetBundle().GetServiceId())
				mid:=id.UnmarshalMethodId(resp.GetBundle().GetMethodId())
				log.Printf("completing method %s on service %s failed due to panic: '%s', exiting",
					mid.Short(), sid.Short(), r)
				debug.PrintStack()
				syscallguest.Exit(ctx, &syscall.ExitRequest{
					Pair: &syscall.ExitPair {
						ServiceId: sid.Marshal(),
						Code: 2,
					},
				})
			}
		}()
		syscallguest.CompleteCall(ctx, syscallguest.CurrentHostId(),cid,r.GetResult(), r.GetResultError())
		return syscall.KernelErr_NoError
	}

	// its a method call from another address space
	sid:=id.UnmarshalServiceId(resp.GetBundle().GetServiceId())
	mid:=id.UnmarshalMethodId(resp.GetBundle().GetMethodId())
	cid:=id.UnmarshalCallId(resp.GetBundle().GetCallId())

	//if mid.Equal(apishared.ExitMethod) {
	// log.Printf("xxx -- got an exit marked read one %s", hid.Short())
	//	os.Exit(51)
	//}

	// we let the invoker handle the unmarshal from anypb.Any because it
	// knows the precise type to be consumed
	fn:=binding.Func(sid,mid)
	if fn==nil {
		slog.Error("unable to find binding for method %s on service, ignoring","mid",mid.Short(),"sid", sid.Short())
		return syscall.KernelErr_NoError
	}
	fut:=fn.Invoke(ctx,resp.GetParamOrResult())
	fut.Success(func (result proto.Message){
		rvReq:=&syscall.ReturnValueRequest{}
		rvReq.Bundle=&syscall.MethodBundle{}

		rvReq.Bundle.CallId= cid.Marshal()
		rvReq.Bundle.HostId= syscallguest.CurrentHostId().Marshal()
		var a anypb.Any
		if err:=a.MarshalFrom(result); err!=nil {
			slog.Error("unable to marshal result for return value request")
			return
		}
		rvReq.Result = &a
		rvReq.ResultError = 0
		syscallguest.ReturnValue(ctx, rvReq) // nowhere for return value to go
	})
	fut.Failure(func (err int32) {
		rvReq:=&syscall.ReturnValueRequest{}
		rvReq.Bundle=&syscall.MethodBundle{}

		rvReq.Bundle.CallId= cid.Marshal()
		rvReq.Bundle.HostId= syscallguest.CurrentHostId().Marshal()
		rvReq.ResultError = err
		syscallguest.ReturnValue(ctx,rvReq) // nowhere for return value to go
	})
	return syscall.KernelErr_NoError

}

func testbind(ctx context.Context,sid id.ServiceId, impl Test) (*lib.ServiceMethodMap, syscall.KernelErr) {
	smmap:=lib.NewServiceMethodMap()
	var mid id.MethodId
	var bindReq *syscall.BindMethodRequest
	var resp *syscall.BindMethodResponse
	var err syscall.KernelErr
//
// test.v1.Test.AddTestSuite
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "AddTestSuite"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Test","TestAddTestSuite",
		GenerateTestAddTestSuiteInvoker(impl))
//
// test.v1.Test.Start
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Start"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Test","TestStart",
		GenerateTestStartInvoker(impl)) 
	return smmap,syscall.KernelErr_NoError
}

// Locate finds a reference to the client interface of test.  
func LocateTest(ctx context.Context,sid id.ServiceId) (ClientTest,syscall.KernelErr) {
    cs, kerr:=client.LocateDynamic(ctx, "test.v1","test", sid)
    if kerr!=syscall.KernelErr_NoError{
        return nil, kerr
    }
    return &ClientTest_{
        BaseService: cs,
    },syscall.KernelErr_NoError
}

func MustLocateTest(ctx context.Context, sid id.ServiceId) ClientTest {
    result, err:=LocateTest(ctx, sid)
    name:=syscall.KernelErr_name[int32(err)]
    normal:="unable to locate test.v1.test:"+name
    if err!=0 {
        if err == syscall.KernelErr_NotRequired {
            slog.Error("service was located, but it was not required")
            panic("locate attempted on a service that was not required")
        }
        panic(normal)
    }
    return result
}


func RegisterTest() (id.ServiceId, syscall.KernelErr){
    req := &syscall.RegisterRequest{}
	debugName:=fmt.Sprintf("%s.%s","test.v1","test")
	req.HostId = syscallguest.CurrentHostId().Marshal()
	req.DebugName = debugName

	resp, err := syscallguest.Register(context.Background(), req)
    if err!=syscall.KernelErr_NoError{
        return id.ServiceIdZeroValue(), err
    }
    sid:=id.UnmarshalServiceId(resp.ServiceId)
    if sid.IsZeroOrEmptyValue() {
        panic("received bad service Id from register")
    }

    return sid,syscall.KernelErr_NoError
}
func MustRegisterTest() (context.Context,id.ServiceId) {
    sid, err:=RegisterTest()
    if err!=syscall.KernelErr_NoError {
        slog.Error("unable to register","package","test.v1","service name","test")
        panic("unable to register "+"test")
    }
    return context.Background(), sid
}

func MustRequireTest(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Require1(ctx, "test.v1","test",sid)
    if err!=syscall.KernelErr_NoError {
        if err==syscall.KernelErr_DependencyCycle{
            slog.Error("unable to require because it creates a dependcy loop","package","test.v1","service name","test","error",syscall.KernelErr_name[int32(err)])
            panic("require test.v1.test creates a dependency loop")
        }
        slog.Error("unable to require","package","test.v1","service name","test","error",syscall.KernelErr_name[int32(err)])
        panic("not able to require test.v1.test:"+syscall.KernelErr_name[int32(err)])
    }
}

func MustExportTest(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Export1(ctx,"test.v1","test",sid)
    if err!=syscall.KernelErr_NoError{
        slog.Error("unable to export","package","test.v1","service name","test")
        panic("not able to export test.v1.test:"+syscall.KernelErr_name[int32(err)])
    }
}

func LaunchServiceTest(ctx context.Context, sid id.ServiceId, impl Test) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture,syscall.KernelErr) {
	smmap, err:=testbind(ctx,sid, impl)
	if err!=0{
		return  nil,nil,syscall.KernelErr(err)
	}
	cid:=id.NewCallId()
	req:=&syscall.LaunchRequest{
		ServiceId: sid.Marshal(),
		CallId: cid.Marshal(),
		HostId: syscallguest.CurrentHostId().Marshal(),
		MethodId: apishared.LaunchMethod.Marshal(),
	}
	fut:=syscallguest.Launch(ctx,req)

    return smmap,fut,syscall.KernelErr_NoError
}

func MustLaunchServiceTest(ctx context.Context, sid id.ServiceId, impl Test) (*lib.ServiceMethodMap, *syscallguest.LaunchFuture) {
    smmap,fut,err:=LaunchServiceTest(ctx,sid,impl)
    if err!=syscall.KernelErr_NoError {
        panic("Unable to call LaunchService successfully: "+syscall.KernelErr_name[int32(err)])
    }
    return smmap,fut
}


// If you want to implement part of your server in host cost you should call 
// Test<methodName>Host from your server implementation. These will be optimized 
// away by the compiler if you don't use them--in other words, if you want to 
// implement everything on the guest side).
// 

//go:wasmimport test add_test_suiteTest_
func AddTestSuite_(int32,int32,int32,int32) int64
func AddTestSuiteTestHost(ctx context.Context,inPtr *AddTestSuiteRequest) *FutureTestAddTestSuite {
	outProtoPtr := (*AddTestSuiteResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, AddTestSuite_)
	f:=NewFutureTestAddTestSuite()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport test startTest_
func Start_(int32,int32,int32,int32) int64
func StartTestHost(ctx context.Context,inPtr *StartRequest) *FutureTestStart {
	outProtoPtr := (*StartResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Start_)
	f:=NewFutureTestStart()
	f.CompleteMethod(ctx,ret,raw)
	return f
}  

// This is interface for invocation.
type invokeTestAddTestSuite struct {
    fn func(context.Context,*AddTestSuiteRequest) *FutureTestAddTestSuite
}

func (t *invokeTestAddTestSuite) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx AddTestSuiteRequest and 'AddTestSuiteRequest{}' why empty?
    in:=&AddTestSuiteRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateTestAddTestSuiteInvoker(impl Test) future.Invoker {
	return &invokeTestAddTestSuite{fn:impl.TestAddTestSuite} 
}

// This is interface for invocation.
type invokeTestStart struct {
    fn func(context.Context,*StartRequest) *FutureTestStart
}

func (t *invokeTestStart) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx StartRequest and 'StartRequest{}' why empty?
    in:=&StartRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateTestStartInvoker(impl Test) future.Invoker {
	return &invokeTestStart{fn:impl.TestStart} 
} 

 
func LaunchMethodCallSuite(ctx context.Context, sid id.ServiceId, impl MethodCallSuite) *future.Base[bool] {

	readyResult:=future.NewBase[bool]()

	ready:=impl.Ready(ctx,sid)
	ready.Handle(func (b bool) {
		if b {
			readyResult.Set(true)			
			return
		}
		slog.Error("Unable to start test.v1.MethodCallSuite, Ready returned false")
		readyResult.Set(false)
	})

	return readyResult
}

// Note that  InitMethodCallSuite returns a future, but the case of failure is covered
// by this definition so the caller need only deal with Success case.
// The context passed here does not need to contain a logger, one will be created.
func InitMethodCallSuite(require []lib.MustRequireFunc, impl MethodCallSuite) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture, context.Context, id.ServiceId){
	// tricky, this context really should not be used but is
	// passed so as to allow printing if things go wrong
	ctx, myId := MustRegisterMethodCallSuite()
	MustExportMethodCallSuite(context.Background(),myId)
	if len(require)>0 {
		for _, f := range require {
			f(ctx, myId)
		}
	}
	smmap, launchF:=MustLaunchServiceMethodCallSuite(ctx, myId, impl)
	launchF.Failure(func (err syscall.KernelErr) {
		t:=syscall.KernelErr_name[int32(err)]
		slog.Error("launch failure on call MethodCallSuite","error",t)
		lib.ExitSelf(ctx, 1, myId)
	})
	return smmap,launchF, ctx,myId
}
func RunMethodCallSuite(ctx context.Context,
	binding *lib.ServiceMethodMap, timeoutInMillis int32, bg lib.Backgrounder) syscall.KernelErr{
	defer func() {
		if r := recover(); r != nil {
			s, ok:=r.(string)
			if !ok && s!=apishared.ControlledExit {
				slog.Error("RunMethodCallSuite: trapped a panic in the guest side", "recovered", r)
			}
		}
	}()
	var kerr syscall.KernelErr
	for {
		kerr:=ReadOneAndCallMethodCallSuite(ctx, binding, timeoutInMillis)
		if kerr == syscall.KernelErr_ReadOneTimeout {
			if bg==nil {
				continue
			}
			slog.Info("calling backgrounder of MethodCallSuite")
			bg.Background(ctx)
			continue
		}
		if kerr == syscall.KernelErr_NoError {
			continue
		}
		break
	}
	slog.Error("error while waiting for MethodCallSuite service calls", "error",syscall.KernelErr_name[int32(kerr)])
	return kerr
}
// Increase this value at your peril!
// Decreasing this value may make your overall program more responsive if you have many services.
var TimeoutInMillisMethodCallSuite = int32(50)

func ReadOneAndCallMethodCallSuite(ctx context.Context, binding *lib.ServiceMethodMap, 
	timeoutInMillis int32) syscall.KernelErr{
	req:=syscall.ReadOneRequest{}
	hid:= syscallguest.CurrentHostId()

	req.TimeoutInMillis = timeoutInMillis
	req.HostId = hid.Marshal()
	resp, err:=syscallguest.ReadOne(ctx, &req)
	if err!=syscall.KernelErr_NoError {
		return err
	}
	// is timeout?
	if resp.Timeout {
		return syscall.KernelErr_ReadOneTimeout
	}

	// check for finished futures from within our address space
	ctx, t:=lib.CurrentTime(ctx)
	syscallguest.ExpireMethod(ctx,t)

	// is a promise being completed that was fulfilled somewhere else
	if r:=resp.GetResolved(); r!=nil {
		cid:=id.UnmarshalCallId(r.GetCallId())
		defer func() {
			if r:=recover(); r!=nil {
				sid:=id.UnmarshalServiceId(resp.GetBundle().GetServiceId())
				mid:=id.UnmarshalMethodId(resp.GetBundle().GetMethodId())
				log.Printf("completing method %s on service %s failed due to panic: '%s', exiting",
					mid.Short(), sid.Short(), r)
				debug.PrintStack()
				syscallguest.Exit(ctx, &syscall.ExitRequest{
					Pair: &syscall.ExitPair {
						ServiceId: sid.Marshal(),
						Code: 2,
					},
				})
			}
		}()
		syscallguest.CompleteCall(ctx, syscallguest.CurrentHostId(),cid,r.GetResult(), r.GetResultError())
		return syscall.KernelErr_NoError
	}

	// its a method call from another address space
	sid:=id.UnmarshalServiceId(resp.GetBundle().GetServiceId())
	mid:=id.UnmarshalMethodId(resp.GetBundle().GetMethodId())
	cid:=id.UnmarshalCallId(resp.GetBundle().GetCallId())

	//if mid.Equal(apishared.ExitMethod) {
	// log.Printf("xxx -- got an exit marked read one %s", hid.Short())
	//	os.Exit(51)
	//}

	// we let the invoker handle the unmarshal from anypb.Any because it
	// knows the precise type to be consumed
	fn:=binding.Func(sid,mid)
	if fn==nil {
		slog.Error("unable to find binding for method %s on service, ignoring","mid",mid.Short(),"sid", sid.Short())
		return syscall.KernelErr_NoError
	}
	fut:=fn.Invoke(ctx,resp.GetParamOrResult())
	fut.Success(func (result proto.Message){
		rvReq:=&syscall.ReturnValueRequest{}
		rvReq.Bundle=&syscall.MethodBundle{}

		rvReq.Bundle.CallId= cid.Marshal()
		rvReq.Bundle.HostId= syscallguest.CurrentHostId().Marshal()
		var a anypb.Any
		if err:=a.MarshalFrom(result); err!=nil {
			slog.Error("unable to marshal result for return value request")
			return
		}
		rvReq.Result = &a
		rvReq.ResultError = 0
		syscallguest.ReturnValue(ctx, rvReq) // nowhere for return value to go
	})
	fut.Failure(func (err int32) {
		rvReq:=&syscall.ReturnValueRequest{}
		rvReq.Bundle=&syscall.MethodBundle{}

		rvReq.Bundle.CallId= cid.Marshal()
		rvReq.Bundle.HostId= syscallguest.CurrentHostId().Marshal()
		rvReq.ResultError = err
		syscallguest.ReturnValue(ctx,rvReq) // nowhere for return value to go
	})
	return syscall.KernelErr_NoError

}

func methodCallSuitebind(ctx context.Context,sid id.ServiceId, impl MethodCallSuite) (*lib.ServiceMethodMap, syscall.KernelErr) {
	smmap:=lib.NewServiceMethodMap()
	var mid id.MethodId
	var bindReq *syscall.BindMethodRequest
	var resp *syscall.BindMethodResponse
	var err syscall.KernelErr
//
// test.v1.MethodCallSuite.Exec
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Exec"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"MethodCallSuite","MethodCallSuiteExec",
		GenerateMethodCallSuiteExecInvoker(impl))
//
// test.v1.MethodCallSuite.SuiteReport
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "SuiteReport"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"MethodCallSuite","MethodCallSuiteSuiteReport",
		GenerateMethodCallSuiteSuiteReportInvoker(impl)) 
	return smmap,syscall.KernelErr_NoError
}

// Locate finds a reference to the client interface of method_call_suite.  
func LocateMethodCallSuite(ctx context.Context,sid id.ServiceId) (ClientMethodCallSuite,syscall.KernelErr) {
    cs, kerr:=client.LocateDynamic(ctx, "test.v1","method_call_suite", sid)
    if kerr!=syscall.KernelErr_NoError{
        return nil, kerr
    }
    return &ClientMethodCallSuite_{
        BaseService: cs,
    },syscall.KernelErr_NoError
}

func MustLocateMethodCallSuite(ctx context.Context, sid id.ServiceId) ClientMethodCallSuite {
    result, err:=LocateMethodCallSuite(ctx, sid)
    name:=syscall.KernelErr_name[int32(err)]
    normal:="unable to locate test.v1.method_call_suite:"+name
    if err!=0 {
        if err == syscall.KernelErr_NotRequired {
            slog.Error("service was located, but it was not required")
            panic("locate attempted on a service that was not required")
        }
        panic(normal)
    }
    return result
}


func RegisterMethodCallSuite() (id.ServiceId, syscall.KernelErr){
    req := &syscall.RegisterRequest{}
	debugName:=fmt.Sprintf("%s.%s","test.v1","method_call_suite")
	req.HostId = syscallguest.CurrentHostId().Marshal()
	req.DebugName = debugName

	resp, err := syscallguest.Register(context.Background(), req)
    if err!=syscall.KernelErr_NoError{
        return id.ServiceIdZeroValue(), err
    }
    sid:=id.UnmarshalServiceId(resp.ServiceId)
    if sid.IsZeroOrEmptyValue() {
        panic("received bad service Id from register")
    }

    return sid,syscall.KernelErr_NoError
}
func MustRegisterMethodCallSuite() (context.Context,id.ServiceId) {
    sid, err:=RegisterMethodCallSuite()
    if err!=syscall.KernelErr_NoError {
        slog.Error("unable to register","package","test.v1","service name","method_call_suite")
        panic("unable to register "+"method_call_suite")
    }
    return context.Background(), sid
}

func MustRequireMethodCallSuite(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Require1(ctx, "test.v1","method_call_suite",sid)
    if err!=syscall.KernelErr_NoError {
        if err==syscall.KernelErr_DependencyCycle{
            slog.Error("unable to require because it creates a dependcy loop","package","test.v1","service name","method_call_suite","error",syscall.KernelErr_name[int32(err)])
            panic("require test.v1.method_call_suite creates a dependency loop")
        }
        slog.Error("unable to require","package","test.v1","service name","method_call_suite","error",syscall.KernelErr_name[int32(err)])
        panic("not able to require test.v1.method_call_suite:"+syscall.KernelErr_name[int32(err)])
    }
}

func MustExportMethodCallSuite(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Export1(ctx,"test.v1","method_call_suite",sid)
    if err!=syscall.KernelErr_NoError{
        slog.Error("unable to export","package","test.v1","service name","method_call_suite")
        panic("not able to export test.v1.method_call_suite:"+syscall.KernelErr_name[int32(err)])
    }
}

func LaunchServiceMethodCallSuite(ctx context.Context, sid id.ServiceId, impl MethodCallSuite) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture,syscall.KernelErr) {
	smmap, err:=methodCallSuitebind(ctx,sid, impl)
	if err!=0{
		return  nil,nil,syscall.KernelErr(err)
	}
	cid:=id.NewCallId()
	req:=&syscall.LaunchRequest{
		ServiceId: sid.Marshal(),
		CallId: cid.Marshal(),
		HostId: syscallguest.CurrentHostId().Marshal(),
		MethodId: apishared.LaunchMethod.Marshal(),
	}
	fut:=syscallguest.Launch(ctx,req)

    return smmap,fut,syscall.KernelErr_NoError
}

func MustLaunchServiceMethodCallSuite(ctx context.Context, sid id.ServiceId, impl MethodCallSuite) (*lib.ServiceMethodMap, *syscallguest.LaunchFuture) {
    smmap,fut,err:=LaunchServiceMethodCallSuite(ctx,sid,impl)
    if err!=syscall.KernelErr_NoError {
        panic("Unable to call LaunchService successfully: "+syscall.KernelErr_name[int32(err)])
    }
    return smmap,fut
}


// If you want to implement part of your server in host cost you should call 
// MethodCallSuite<methodName>Host from your server implementation. These will be optimized 
// away by the compiler if you don't use them--in other words, if you want to 
// implement everything on the guest side).
// 

//go:wasmimport test execMethodCallSuite_
func Exec_(int32,int32,int32,int32) int64
func ExecMethodCallSuiteHost(ctx context.Context,inPtr *ExecRequest) *FutureMethodCallSuiteExec {
	outProtoPtr := (*ExecResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Exec_)
	f:=NewFutureMethodCallSuiteExec()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport test suite_reportMethodCallSuite_
func SuiteReport_(int32,int32,int32,int32) int64
func SuiteReportMethodCallSuiteHost(ctx context.Context,inPtr *SuiteReportRequest) *FutureMethodCallSuiteSuiteReport {
	outProtoPtr := (*SuiteReportResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, SuiteReport_)
	f:=NewFutureMethodCallSuiteSuiteReport()
	f.CompleteMethod(ctx,ret,raw)
	return f
}  

// This is interface for invocation.
type invokeMethodCallSuiteExec struct {
    fn func(context.Context,*ExecRequest) *FutureMethodCallSuiteExec
}

func (t *invokeMethodCallSuiteExec) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx ExecRequest and 'ExecRequest{}' why empty?
    in:=&ExecRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateMethodCallSuiteExecInvoker(impl MethodCallSuite) future.Invoker {
	return &invokeMethodCallSuiteExec{fn:impl.MethodCallSuiteExec} 
}

// This is interface for invocation.
type invokeMethodCallSuiteSuiteReport struct {
    fn func(context.Context,*SuiteReportRequest) *FutureMethodCallSuiteSuiteReport
}

func (t *invokeMethodCallSuiteSuiteReport) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx SuiteReportRequest and '' why empty?
    in:=&SuiteReportRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateMethodCallSuiteSuiteReportInvoker(impl MethodCallSuite) future.Invoker {
	return &invokeMethodCallSuiteSuiteReport{fn:impl.MethodCallSuiteSuiteReport} 
} 

 
func LaunchUnderTest(ctx context.Context, sid id.ServiceId, impl UnderTest) *future.Base[bool] {

	readyResult:=future.NewBase[bool]()

	ready:=impl.Ready(ctx,sid)
	ready.Handle(func (b bool) {
		if b {
			readyResult.Set(true)			
			return
		}
		slog.Error("Unable to start test.v1.UnderTest, Ready returned false")
		readyResult.Set(false)
	})

	return readyResult
}

// Note that  InitUnderTest returns a future, but the case of failure is covered
// by this definition so the caller need only deal with Success case.
// The context passed here does not need to contain a logger, one will be created.
func InitUnderTest(require []lib.MustRequireFunc, impl UnderTest) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture, context.Context, id.ServiceId){
	// tricky, this context really should not be used but is
	// passed so as to allow printing if things go wrong
	ctx, myId := MustRegisterUnderTest()
	MustExportUnderTest(context.Background(),myId)
	if len(require)>0 {
		for _, f := range require {
			f(ctx, myId)
		}
	}
	smmap, launchF:=MustLaunchServiceUnderTest(ctx, myId, impl)
	launchF.Failure(func (err syscall.KernelErr) {
		t:=syscall.KernelErr_name[int32(err)]
		slog.Error("launch failure on call UnderTest","error",t)
		lib.ExitSelf(ctx, 1, myId)
	})
	return smmap,launchF, ctx,myId
}
func RunUnderTest(ctx context.Context,
	binding *lib.ServiceMethodMap, timeoutInMillis int32, bg lib.Backgrounder) syscall.KernelErr{
	defer func() {
		if r := recover(); r != nil {
			s, ok:=r.(string)
			if !ok && s!=apishared.ControlledExit {
				slog.Error("RunUnderTest: trapped a panic in the guest side", "recovered", r)
			}
		}
	}()
	var kerr syscall.KernelErr
	for {
		kerr:=ReadOneAndCallUnderTest(ctx, binding, timeoutInMillis)
		if kerr == syscall.KernelErr_ReadOneTimeout {
			if bg==nil {
				continue
			}
			slog.Info("calling backgrounder of UnderTest")
			bg.Background(ctx)
			continue
		}
		if kerr == syscall.KernelErr_NoError {
			continue
		}
		break
	}
	slog.Error("error while waiting for UnderTest service calls", "error",syscall.KernelErr_name[int32(kerr)])
	return kerr
}
// Increase this value at your peril!
// Decreasing this value may make your overall program more responsive if you have many services.
var TimeoutInMillisUnderTest = int32(50)

func ReadOneAndCallUnderTest(ctx context.Context, binding *lib.ServiceMethodMap, 
	timeoutInMillis int32) syscall.KernelErr{
	req:=syscall.ReadOneRequest{}
	hid:= syscallguest.CurrentHostId()

	req.TimeoutInMillis = timeoutInMillis
	req.HostId = hid.Marshal()
	resp, err:=syscallguest.ReadOne(ctx, &req)
	if err!=syscall.KernelErr_NoError {
		return err
	}
	// is timeout?
	if resp.Timeout {
		return syscall.KernelErr_ReadOneTimeout
	}

	// check for finished futures from within our address space
	ctx, t:=lib.CurrentTime(ctx)
	syscallguest.ExpireMethod(ctx,t)

	// is a promise being completed that was fulfilled somewhere else
	if r:=resp.GetResolved(); r!=nil {
		cid:=id.UnmarshalCallId(r.GetCallId())
		defer func() {
			if r:=recover(); r!=nil {
				sid:=id.UnmarshalServiceId(resp.GetBundle().GetServiceId())
				mid:=id.UnmarshalMethodId(resp.GetBundle().GetMethodId())
				log.Printf("completing method %s on service %s failed due to panic: '%s', exiting",
					mid.Short(), sid.Short(), r)
				debug.PrintStack()
				syscallguest.Exit(ctx, &syscall.ExitRequest{
					Pair: &syscall.ExitPair {
						ServiceId: sid.Marshal(),
						Code: 2,
					},
				})
			}
		}()
		syscallguest.CompleteCall(ctx, syscallguest.CurrentHostId(),cid,r.GetResult(), r.GetResultError())
		return syscall.KernelErr_NoError
	}

	// its a method call from another address space
	sid:=id.UnmarshalServiceId(resp.GetBundle().GetServiceId())
	mid:=id.UnmarshalMethodId(resp.GetBundle().GetMethodId())
	cid:=id.UnmarshalCallId(resp.GetBundle().GetCallId())

	//if mid.Equal(apishared.ExitMethod) {
	// log.Printf("xxx -- got an exit marked read one %s", hid.Short())
	//	os.Exit(51)
	//}

	// we let the invoker handle the unmarshal from anypb.Any because it
	// knows the precise type to be consumed
	fn:=binding.Func(sid,mid)
	if fn==nil {
		slog.Error("unable to find binding for method %s on service, ignoring","mid",mid.Short(),"sid", sid.Short())
		return syscall.KernelErr_NoError
	}
	fut:=fn.Invoke(ctx,resp.GetParamOrResult())
	fut.Success(func (result proto.Message){
		rvReq:=&syscall.ReturnValueRequest{}
		rvReq.Bundle=&syscall.MethodBundle{}

		rvReq.Bundle.CallId= cid.Marshal()
		rvReq.Bundle.HostId= syscallguest.CurrentHostId().Marshal()
		var a anypb.Any
		if err:=a.MarshalFrom(result); err!=nil {
			slog.Error("unable to marshal result for return value request")
			return
		}
		rvReq.Result = &a
		rvReq.ResultError = 0
		syscallguest.ReturnValue(ctx, rvReq) // nowhere for return value to go
	})
	fut.Failure(func (err int32) {
		rvReq:=&syscall.ReturnValueRequest{}
		rvReq.Bundle=&syscall.MethodBundle{}

		rvReq.Bundle.CallId= cid.Marshal()
		rvReq.Bundle.HostId= syscallguest.CurrentHostId().Marshal()
		rvReq.ResultError = err
		syscallguest.ReturnValue(ctx,rvReq) // nowhere for return value to go
	})
	return syscall.KernelErr_NoError

}

func underTestbind(ctx context.Context,sid id.ServiceId, impl UnderTest) (*lib.ServiceMethodMap, syscall.KernelErr) {
	smmap:=lib.NewServiceMethodMap()
	var mid id.MethodId
	var bindReq *syscall.BindMethodRequest
	var resp *syscall.BindMethodResponse
	var err syscall.KernelErr
//
// test.v1.UnderTest.Exec
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Exec"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"UnderTest","UnderTestExec",
		GenerateUnderTestExecInvoker(impl)) 
	return smmap,syscall.KernelErr_NoError
}

// Locate finds a reference to the client interface of under_test.  
func LocateUnderTest(ctx context.Context,sid id.ServiceId) (ClientUnderTest,syscall.KernelErr) {
    cs, kerr:=client.LocateDynamic(ctx, "test.v1","under_test", sid)
    if kerr!=syscall.KernelErr_NoError{
        return nil, kerr
    }
    return &ClientUnderTest_{
        BaseService: cs,
    },syscall.KernelErr_NoError
}

func MustLocateUnderTest(ctx context.Context, sid id.ServiceId) ClientUnderTest {
    result, err:=LocateUnderTest(ctx, sid)
    name:=syscall.KernelErr_name[int32(err)]
    normal:="unable to locate test.v1.under_test:"+name
    if err!=0 {
        if err == syscall.KernelErr_NotRequired {
            slog.Error("service was located, but it was not required")
            panic("locate attempted on a service that was not required")
        }
        panic(normal)
    }
    return result
}


func RegisterUnderTest() (id.ServiceId, syscall.KernelErr){
    req := &syscall.RegisterRequest{}
	debugName:=fmt.Sprintf("%s.%s","test.v1","under_test")
	req.HostId = syscallguest.CurrentHostId().Marshal()
	req.DebugName = debugName

	resp, err := syscallguest.Register(context.Background(), req)
    if err!=syscall.KernelErr_NoError{
        return id.ServiceIdZeroValue(), err
    }
    sid:=id.UnmarshalServiceId(resp.ServiceId)
    if sid.IsZeroOrEmptyValue() {
        panic("received bad service Id from register")
    }

    return sid,syscall.KernelErr_NoError
}
func MustRegisterUnderTest() (context.Context,id.ServiceId) {
    sid, err:=RegisterUnderTest()
    if err!=syscall.KernelErr_NoError {
        slog.Error("unable to register","package","test.v1","service name","under_test")
        panic("unable to register "+"under_test")
    }
    return context.Background(), sid
}

func MustRequireUnderTest(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Require1(ctx, "test.v1","under_test",sid)
    if err!=syscall.KernelErr_NoError {
        if err==syscall.KernelErr_DependencyCycle{
            slog.Error("unable to require because it creates a dependcy loop","package","test.v1","service name","under_test","error",syscall.KernelErr_name[int32(err)])
            panic("require test.v1.under_test creates a dependency loop")
        }
        slog.Error("unable to require","package","test.v1","service name","under_test","error",syscall.KernelErr_name[int32(err)])
        panic("not able to require test.v1.under_test:"+syscall.KernelErr_name[int32(err)])
    }
}

func MustExportUnderTest(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Export1(ctx,"test.v1","under_test",sid)
    if err!=syscall.KernelErr_NoError{
        slog.Error("unable to export","package","test.v1","service name","under_test")
        panic("not able to export test.v1.under_test:"+syscall.KernelErr_name[int32(err)])
    }
}

func LaunchServiceUnderTest(ctx context.Context, sid id.ServiceId, impl UnderTest) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture,syscall.KernelErr) {
	smmap, err:=underTestbind(ctx,sid, impl)
	if err!=0{
		return  nil,nil,syscall.KernelErr(err)
	}
	cid:=id.NewCallId()
	req:=&syscall.LaunchRequest{
		ServiceId: sid.Marshal(),
		CallId: cid.Marshal(),
		HostId: syscallguest.CurrentHostId().Marshal(),
		MethodId: apishared.LaunchMethod.Marshal(),
	}
	fut:=syscallguest.Launch(ctx,req)

    return smmap,fut,syscall.KernelErr_NoError
}

func MustLaunchServiceUnderTest(ctx context.Context, sid id.ServiceId, impl UnderTest) (*lib.ServiceMethodMap, *syscallguest.LaunchFuture) {
    smmap,fut,err:=LaunchServiceUnderTest(ctx,sid,impl)
    if err!=syscall.KernelErr_NoError {
        panic("Unable to call LaunchService successfully: "+syscall.KernelErr_name[int32(err)])
    }
    return smmap,fut
}


// If you want to implement part of your server in host cost you should call 
// UnderTest<methodName>Host from your server implementation. These will be optimized 
// away by the compiler if you don't use them--in other words, if you want to 
// implement everything on the guest side).
// 

//go:wasmimport test exec_under_testUnderTest_
func ExecUnderTest_(int32,int32,int32,int32) int64
func ExecUnderTestUnderTestHost(ctx context.Context,inPtr *ExecRequest) *FutureUnderTestExec {
	outProtoPtr := (*ExecResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Exec_)
	f:=NewFutureUnderTestExec()
	f.CompleteMethod(ctx,ret,raw)
	return f
}  

// This is interface for invocation.
type invokeUnderTestExec struct {
    fn func(context.Context,*ExecRequest) *FutureUnderTestExec
}

func (t *invokeUnderTestExec) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx ExecRequest and 'ExecRequest{}' why empty?
    in:=&ExecRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateUnderTestExecInvoker(impl UnderTest) future.Invoker {
	return &invokeUnderTestExec{fn:impl.UnderTestExec} 
}  
