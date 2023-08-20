//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: methodcall/v1/suite.proto

package methodcall




import (
	"context"
	"fmt"
    "unsafe" 
    // this set of imports is _unrelated_ to the particulars of what the .proto imported... those are above
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"  
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/api/shared/id"
	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/lib/go/future"
	"github.com/iansmith/parigot/lib/go/client"
	"github.com/iansmith/parigot/api/guest"  

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/proto"

)
var _ =  unsafe.Sizeof([]byte{})
 
func Launch(ctx context.Context, sid id.ServiceId, impl MethodCallSuite) *future.Base[bool] {

	readyResult:=future.NewBase[bool]()

	ready:=impl.Ready(ctx,sid)
	ready.Handle(func (b bool) {
		if b {
			readyResult.Set(true)			
			return
		}
		guest.Log(ctx).Error("Unable to start methodcall.v1.MethodCallSuite, Ready returned false")
		readyResult.Set(false)
	})

	return readyResult
}

// Note that  Init returns a future, but the case of failure is covered
// by this definition so the caller need only deal with Success case.
// The context passed here does not need to contain a logger, one will be created.
func Init(require []lib.MustRequireFunc, impl MethodCallSuite) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture, context.Context, id.ServiceId){
	defer func() {
		if r := recover(); r != nil {
			guest.Log(context.Background()).Info("InitMethodCallSuite: trapped a panic in the guest side","recovered", r)
		}
	}()

	// tricky, this context really should not be used but is
	// passed so as to allow printing if things go wrong
	ctx, myId := MustRegister()
	MustExport(context.Background(),myId)
	if len(require)>0 {
		for _, f := range require {
			f(ctx, myId)
		}
	}
	smmap, launchF:=MustLaunchService(ctx, myId, impl)
	launchF.Failure(func (err syscall.KernelErr) {
		t:=syscall.KernelErr_name[int32(err)]
		guest.Log(ctx).Error("launch failure on call MethodCallSuite","error",t)
		lib.ExitClient(ctx, 1, myId, "unable to Launch in Init:"+t,
			"unable to call Exit in Init:"+t)
	})
	return smmap,launchF, ctx,myId
}
func Run(ctx context.Context,
	binding *lib.ServiceMethodMap, timeoutInMillis int32, bg lib.Backgrounder) syscall.KernelErr{
	defer func() {
		if r := recover(); r != nil {
			s, ok:=r.(string)
			if !ok && s!=apishared.ControlledExit {
				guest.Log(ctx).Error("Run: trapped a panic in the guest side", "recovered", r)
			}
		}
	}()
	var kerr syscall.KernelErr
	for {
		kerr:=ReadOneAndCall(ctx, binding, timeoutInMillis)
		if kerr == syscall.KernelErr_ReadOneTimeout {
			if bg==nil {
				continue
			}
			guest.Log(ctx).Info("calling backgrounder of MethodCallSuite")
			bg.Background(ctx)
			continue
		}
		if kerr == syscall.KernelErr_NoError {
			continue
		}
		break
	}
	guest.Log(ctx).Error("error while waiting for  service calls", "error",syscall.KernelErr_name[int32(kerr)])
	return kerr
}
// Increase this value at your peril!
// Decreasing this value may make your overall program more responsive if you have many services.
var TimeoutInMillis = int32(50)

func ReadOneAndCall(ctx context.Context, binding *lib.ServiceMethodMap, 
	timeoutInMillis int32) syscall.KernelErr{
	req:=syscall.ReadOneRequest{}

	req.TimeoutInMillis = timeoutInMillis
	req.HostId = syscallguest.CurrentHostId().Marshal()
	resp, err:=syscallguest.ReadOne(&req)
	if err!=syscall.KernelErr_NoError {
		return err
	}
	// is timeout?
	if resp.Timeout {
		return syscall.KernelErr_ReadOneTimeout
	}

	// check for finished futures from within our address space
	syscallguest.ExpireMethod(ctx)

	// is a promise being completed that was fulfilled somewhere else
	if r:=resp.GetResolved(); r!=nil {
		cid:=id.UnmarshalCallId(r.GetCallId())
		syscallguest.CompleteCall(ctx, syscallguest.CurrentHostId(),cid,r.GetResult(), r.GetResultError())
		return syscall.KernelErr_NoError
	}

	// its a method call from another address space
	sid:=id.UnmarshalServiceId(resp.GetBundle().GetServiceId())
	mid:=id.UnmarshalMethodId(resp.GetBundle().GetMethodId())
	cid:=id.UnmarshalCallId(resp.GetBundle().GetCallId())

	if mid.Equal(apishared.ExitMethod) {
		panic(apishared.ControlledExit)
	}
	// we let the invoker handle the unmarshal from anypb.Any because it
	// knows the precise type to be consumed
	fn:=binding.Func(sid,mid)
	if fn==nil {
		guest.Log(ctx).Error("unable to find binding for method %s on service, ignoring","mid",mid.Short(),"sid", sid.Short())
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
			guest.Log(ctx).Error("unable to marshal result for return value request")
			return
		}
		rvReq.Result = &a
		rvReq.ResultError = 0
		syscallguest.ReturnValue(rvReq) // nowhere for return value to go
	})
	fut.Failure(func (err int32) {
		rvReq:=&syscall.ReturnValueRequest{}
		rvReq.Bundle=&syscall.MethodBundle{}

		rvReq.Bundle.CallId= cid.Marshal()
		rvReq.Bundle.HostId= syscallguest.CurrentHostId().Marshal()
		rvReq.ResultError = err
		syscallguest.ReturnValue(rvReq) // nowhere for return value to go
	})
	return syscall.KernelErr_NoError

}

func bind(ctx context.Context,sid id.ServiceId, impl MethodCallSuite) (*lib.ServiceMethodMap, syscall.KernelErr) {
	smmap:=lib.NewServiceMethodMap()
	var mid id.MethodId
	var bindReq *syscall.BindMethodRequest
	var resp *syscall.BindMethodResponse
	var err syscall.KernelErr
//
// methodcall.v1.MethodCallSuite.Exec
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Exec"
	resp, err=syscallguest.BindMethod(bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"MethodCallSuite","Exec",
		GenerateExecInvoker(impl))
//
// methodcall.v1.MethodCallSuite.SuiteReport
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "SuiteReport"
	resp, err=syscallguest.BindMethod(bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"MethodCallSuite","SuiteReport",
		GenerateSuiteReportInvoker(impl)) 
	return smmap,syscall.KernelErr_NoError
}

// Locate finds a reference to the client interface of method_call_suite.  
func Locate(ctx context.Context,sid id.ServiceId) (Client,syscall.KernelErr) {
    cs, kerr:=client.LocateDynamic(ctx, "methodcall.v1","method_call_suite", sid)
    if kerr!=syscall.KernelErr_NoError{
        return nil, kerr
    }
    return &Client_{
        BaseService: cs,
    },syscall.KernelErr_NoError
}

func MustLocate(ctx context.Context, sid id.ServiceId) Client {
    result, err:=Locate(ctx, sid)
    name:=syscall.KernelErr_name[int32(err)]
    normal:="unable to locate methodcall.v1.method_call_suite:"+name
    if err!=0 {
        if err == syscall.KernelErr_NotRequired {
            guest.Log(ctx).Error("service was located, but it was not required")
            panic("locate attempted on a service that was not required")
        }
        panic(normal)
    }
    return result
}


func Register() (id.ServiceId, syscall.KernelErr){
    req := &syscall.RegisterRequest{}
	debugName:=fmt.Sprintf("%s.%s","methodcall.v1","method_call_suite")
	req.HostId = syscallguest.CurrentHostId().Marshal()
	req.DebugName = debugName

	resp, err := syscallguest.Register(req)
    if err!=syscall.KernelErr_NoError{
        return id.ServiceIdZeroValue(), err
    }
    sid:=id.UnmarshalServiceId(resp.ServiceId)
    if sid.IsZeroOrEmptyValue() {
        panic("received bad service Id from register")
    }

    return sid,syscall.KernelErr_NoError
}
func MustRegister() (context.Context,id.ServiceId) {
    sid, err:=Register()
    if err!=syscall.KernelErr_NoError {
        guest.Log(context.Background()).Error("unable to register","package","methodcall.v1","service name","method_call_suite")
        panic("unable to register "+"method_call_suite")
    }
    return guest.NewContextWithLogger(sid), sid
}

func MustRequire(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Require1("methodcall.v1","method_call_suite",sid)
    if err!=syscall.KernelErr_NoError {
        if err==syscall.KernelErr_DependencyCycle{
            guest.Log(ctx).Error("unable to require because it creates a dependcy loop","package","methodcall.v1","service name","method_call_suite","error",syscall.KernelErr_name[int32(err)])
            panic("require methodcall.v1.method_call_suite creates a dependency loop")
        }
        guest.Log(ctx).Error("unable to require","package","methodcall.v1","service name","method_call_suite","error",syscall.KernelErr_name[int32(err)])
        panic("not able to require methodcall.v1.method_call_suite:"+syscall.KernelErr_name[int32(err)])
    }
}

func MustExport(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Export1("methodcall.v1","method_call_suite",sid)
    if err!=syscall.KernelErr_NoError{
        guest.Log(ctx).Error("unable to export","package","methodcall.v1","service name","method_call_suite")
        panic("not able to export methodcall.v1.method_call_suite:"+syscall.KernelErr_name[int32(err)])
    }
}

func LaunchService(ctx context.Context, sid id.ServiceId, impl MethodCallSuite) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture,syscall.KernelErr) {
	smmap, err:=bind(ctx,sid, impl)
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
	fut:=syscallguest.Launch(req)

    return smmap,fut,syscall.KernelErr_NoError
}

func MustLaunchService(ctx context.Context, sid id.ServiceId, impl MethodCallSuite) (*lib.ServiceMethodMap, *syscallguest.LaunchFuture) {
    smmap,fut,err:=LaunchService(ctx,sid,impl)
    if err!=syscall.KernelErr_NoError {
        panic("Unable to call LaunchService successfully: "+syscall.KernelErr_name[int32(err)])
    }
    return smmap,fut
}


// If you want to implement part of your server in host cost you should call 
// <methodName>Host from your server implementation. These will be optimized 
// away by the compiler if you don't use them--in other words, if you want to 
// implement everything on the guest side).
// 

//go:wasmimport methodcall exec_
func Exec_(int32,int32,int32,int32) int64
func ExecHost(ctx context.Context,inPtr *ExecRequest) *FutureExec {
	outProtoPtr := (*ExecResponse)(nil)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Exec_)
	if signal {
		guest.Log(ctx).Info("Exec exiting because of parigot signal")
		lib.ExitClient(ctx, 1, id.NewServiceId(), "xxx warning, no implementation of unsolicited exit",
			"xxx warning, no implementation of unsolicited exit and failed trying to exit")
	}
	f:=NewFutureExec()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport methodcall suite_report_
func SuiteReport_(int32,int32,int32,int32) int64
func SuiteReportHost(ctx context.Context,inPtr *SuiteReportRequest) *FutureSuiteReport {
	outProtoPtr := (*SuiteReportResponse)(nil)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, SuiteReport_)
	if signal {
		guest.Log(ctx).Info("SuiteReport exiting because of parigot signal")
		lib.ExitClient(ctx, 1, id.NewServiceId(), "xxx warning, no implementation of unsolicited exit",
			"xxx warning, no implementation of unsolicited exit and failed trying to exit")
	}
	f:=NewFutureSuiteReport()
	f.CompleteMethod(ctx,ret,raw)
	return f
}  

// This is interface for invocation.
type invokeExec struct {
    fn func(context.Context,*ExecRequest) *FutureExec
}

func (t *invokeExec) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx ExecRequest and 'ExecRequest{}' why empty?
    in:=&ExecRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        guest.Log(ctx).Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateExecInvoker(impl MethodCallSuite) future.Invoker {
	return &invokeExec{fn:impl.Exec} 
}

// This is interface for invocation.
type invokeSuiteReport struct {
    fn func(context.Context,*SuiteReportRequest) *FutureSuiteReport
}

func (t *invokeSuiteReport) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx SuiteReportRequest and '' why empty?
    in:=&SuiteReportRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        guest.Log(ctx).Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateSuiteReportInvoker(impl MethodCallSuite) future.Invoker {
	return &invokeSuiteReport{fn:impl.SuiteReport} 
}  
