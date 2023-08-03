//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: httpconnector/v1/httpconnector.proto

package httpconnector




import (
	"context"
    "unsafe" 
    // this set of imports is _unrelated_ to the particulars of what the .proto imported... those are above
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"  
	pcontext "github.com/iansmith/parigot/context"
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
 
func Launch(ctx context.Context, sid id.ServiceId, impl HttpConnector) *future.Base[bool] {

	defer func() {
		pcontext.Dump(ctx)
	}()

	readyResult:=future.NewBase[bool]()

	ready:=impl.Ready(ctx,sid)
	ready.Handle(func (b bool) {
		if b {
			readyResult.Set(true)			
			return
		}
		pcontext.Errorf(ctx,"Unable to start httpconnector.v1.HttpConnector, Ready returned false")
		pcontext.Dump(ctx)
		readyResult.Set(false)
	})

	return readyResult
}

// Note that  Init returns a future, but the case of failure is covered
// by this definition so the caller need only deal with Success case.
func Init(ctx context.Context,require []lib.MustRequireFunc, impl HttpConnector) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture, id.ServiceId){
	defer func() {
		if r := recover(); r != nil {
			pcontext.Infof(ctx, "InitHttpConnector: trapped a panic in the guest side: %v", r)
		}
		pcontext.Dump(ctx)
	}()

	myId := MustRegister(ctx)
	MustExport(ctx,myId)
	if len(require)>0 {
		for _, f := range require {
			f(ctx, myId)
		}
	}
	smmap, launchF:=MustLaunchService(ctx, myId, impl)
	launchF.Failure(func (err syscall.KernelErr) {
		t:=syscall.KernelErr_name[int32(err)]
		pcontext.Errorf(ctx, "launch failure on call HttpConnector:%s",t)
		lib.ExitClient(ctx, 1, myId, "unable to Launch in Init:"+t,
			"unable to call Exit in Init:"+t)
	})
	return smmap,launchF, myId
}
func Run(ctx context.Context,
	binding *lib.ServiceMethodMap, timeoutInMillis int32, bg lib.Backgrounder) syscall.KernelErr{
	defer func() {
		if r := recover(); r != nil {
			pcontext.Infof(ctx, "Run: trapped a panic in the guest side: %v", r)
		}
		pcontext.Dump(ctx)
	}()
	var kerr syscall.KernelErr
	for {
		pctx:=pcontext.CallTo(ctx,"ReadOneAndCall")
		kerr:=ReadOneAndCall(pctx, binding, timeoutInMillis)
		pcontext.Dump(pctx)
		if kerr == syscall.KernelErr_ReadOneTimeout {
			if bg==nil {
				continue
			}
			pcontext.Infof(ctx,"calling backgrounder of HttpConnector")
			bgctx:=pcontext.CallTo(ctx,"Background")
			bg.Background(bgctx)
			pcontext.Dump(bgctx)
			continue
		}
		if kerr == syscall.KernelErr_NoError {
			continue
		}
		break
	}
	pcontext.Errorf(ctx, "error while waiting for  service calls: %s", syscall.KernelErr_name[int32(kerr)])
	return kerr
}
// Increase this value at your peril!
// Decreasing this value may make your overall program more responsive if you have many services.
var TimeoutInMillis = int32(50)

func ReadOneAndCall(ctx context.Context, binding *lib.ServiceMethodMap, 
	timeoutInMillis int32) syscall.KernelErr{
	req:=syscall.ReadOneRequest{}

	// makes a copy
	for _, c := range binding.Call() {
		req.Call=append(req.Call, c)
	}

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
		syscallguest.CompleteCall(ctx, cid,r.GetResult(), r.GetResultError())
		return syscall.KernelErr_NoError
	}

	// its a method call from another address space
	sid:=id.UnmarshalServiceId(resp.GetCall().GetServiceId())
	mid:=id.UnmarshalMethodId(resp.GetCall().GetMethodId())
	cid:=id.UnmarshalCallId(resp.GetCallId())
	fn:=binding.Func(sid,mid)
	if fn==nil {
		pcontext.Errorf(ctx,"unable to find service/method pair binding: service=%s, method=%s",
			sid.Short(), mid.Short())
		return syscall.KernelErr_NotFound
	}
	// we let the invoker handle the unmarshal from anypb.Any because it
	// knows the precise type to be consumed
	fut:=fn.Invoke(ctx,resp.GetParam())
	fut.Success(func (result proto.Message){
		rvReq:=&syscall.ReturnValueRequest{}
		rvReq.CallId= cid.Marshal()
		rvReq.HostId= syscallguest.CurrentHostId().Marshal()
		var a anypb.Any
		if err:=a.MarshalFrom(result); err!=nil {
			pcontext.Errorf(ctx, "unable to marshal result for return value request")
			return
		}
		rvReq.Result = &a
		rvReq.ResultError = 0
		syscallguest.ReturnValue(rvReq) // nowhere for return value to go
	})
	fut.Failure(func (err int32) {
		rvReq:=&syscall.ReturnValueRequest{}
		rvReq.CallId= cid.Marshal()
		rvReq.HostId= syscallguest.CurrentHostId().Marshal()
		rvReq.ResultError = err
		syscallguest.ReturnValue(rvReq) // nowhere for return value to go
	})
	return syscall.KernelErr_NoError

}

func bind(ctx context.Context,sid id.ServiceId, impl HttpConnector) (*lib.ServiceMethodMap, syscall.KernelErr) {
	smmap:=lib.NewServiceMethodMap()
	var mid id.MethodId
	var bindReq *syscall.BindMethodRequest
	var resp *syscall.BindMethodResponse
	var err syscall.KernelErr
//
// httpconnector.v1.HttpConnector.Check
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Check"
	resp, err=syscallguest.BindMethod(bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"HttpConnector","Check",
		GenerateCheckInvoker(impl)) 
	pcontext.Dump(ctx)
	return smmap,syscall.KernelErr_NoError
}

// Locate finds a reference to the client interface of http_connector.  
func Locate(ctx context.Context,sid id.ServiceId) (Client,syscall.KernelErr) {
    cs, kerr:=client.LocateDynamic(ctx, "httpconnector.v1","http_connector", sid)
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
    normal:="unable to locate httpconnector.v1.http_connector:"+name
    if err!=0 {
        pcontext.Debugf(ctx,"kernel error was  %s",name)
        if err == syscall.KernelErr_NotRequired {
            pcontext.Errorf(ctx,"service was located, but it was not required")
            panic("locate attempted on a service that was not required")
        }
        panic(normal)
    }
    return result
}


func Register(ctx context.Context) (id.ServiceId, syscall.KernelErr){
    req := &syscall.RegisterRequest{}
	fqs := &syscall.FullyQualifiedService{
		PackagePath: "httpconnector.v1",
		Service:     "http_connector",
	}
	req.Fqs = fqs
	req.HostId = syscallguest.CurrentHostId().Marshal()

	resp, err := syscallguest.Register(req)
    if err!=syscall.KernelErr_NoError{
        return id.ServiceIdZeroValue(), err
    }
    sid:=id.UnmarshalServiceId(resp.Id)
    if sid.IsZeroOrEmptyValue() {
        panic("received bad service Id from register")
    }

    return sid,syscall.KernelErr_NoError
}
func MustRegister(ctx context.Context) id.ServiceId {
    sid, err:=Register(ctx)
    if err!=syscall.KernelErr_NoError {
        pcontext.Fatalf(ctx,"unable to register %s.%s","httpconnector.v1","http_connector")
        panic("unable to register "+"http_connector")
    }
    return sid
}

func MustRequire(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Require1("httpconnector.v1","http_connector",sid)
    if err!=syscall.KernelErr_NoError {
        if err==syscall.KernelErr_DependencyCycle{
            pcontext.Errorf(ctx,"unable to require %s.%s because it creates a dependcy loop: %s","httpconnector.v1","http_connector",syscall.KernelErr_name[int32(err)])
            panic("require httpconnector.v1.http_connector creates a dependency loop")
        }
        pcontext.Errorf(ctx,"unable to require %s.%s:%s","httpconnector.v1","http_connector",syscall.KernelErr_name[int32(err)])
        panic("not able to require httpconnector.v1.http_connector:"+syscall.KernelErr_name[int32(err)])
    }
}

func MustExport(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Export1("httpconnector.v1","http_connector",sid)
    if err!=syscall.KernelErr_NoError{
        pcontext.Fatalf(ctx, "unable to export %s.%s","httpconnector.v1","http_connector")
        panic("not able to export httpconnector.v1.http_connector:"+syscall.KernelErr_name[int32(err)])
    }
}

func LaunchService(ctx context.Context, sid id.ServiceId, impl HttpConnector) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture,syscall.KernelErr) {
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

func MustLaunchService(ctx context.Context, sid id.ServiceId, impl HttpConnector) (*lib.ServiceMethodMap, *syscallguest.LaunchFuture) {
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

//go:wasmimport httpconnector check_
func Check_(int32,int32,int32,int32) int64
func CheckHost(ctx context.Context,inPtr *CheckRequest) *FutureCheck {
	outProtoPtr := (*CheckResponse)(nil)
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Check_)
	if signal {
		pcontext.Infof(ctx, "Check exiting because of parigot signal")
		pcontext.Dump(ctx)
		lib.ExitClient(ctx, 1, id.NewServiceId(), "xxx warning, no implementation of unsolicited exit",
			"xxx warning, no implementation of unsolicited exit and failed trying to exit")
	}
	f:=NewFutureCheck()
	f.CompleteMethod(ctx,ret,raw)
	return f
}  

// This is interface for invocation.
type invokeCheck struct {
    fn func(context.Context,*CheckRequest) *FutureCheck
}

func (t *invokeCheck) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx CheckRequest and 'CheckRequest{}' why empty?
    in:=&CheckRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateCheckInvoker(impl HttpConnector) future.Invoker {
	return &invokeCheck{fn:impl.Check} 
}  
