//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: simple/v1/simple.proto

package simple




import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"runtime/debug"
    "unsafe"

"github.com/iansmith/parigot/g/http/v1" 
    // this set of imports is _unrelated_ to the particulars of what the .proto imported... those are above
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"  
	"github.com/iansmith/parigot/api/shared/id"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/lib/go/future"
	apishared "github.com/iansmith/parigot/api/shared"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"github.com/iansmith/parigot/lib/go/client"
)
var _ =  unsafe.Sizeof([]byte{})
 

func Launch(ctx context.Context, sid id.ServiceId, impl Simple) *future.Base[bool] {

	readyResult:=future.NewBase[bool]()

	ready:=impl.Ready(ctx,sid)
	ready.Handle(func (b bool) {
		if b {
			readyResult.Set(true)			
			return
		}
		slog.Error("Unable to start simple.v1.Simple, Ready returned false")
		readyResult.Set(false)
	})

	return readyResult
}

// Note that  Init returns a future, but the case of failure is covered
// by this definition so the caller need only deal with Success case.
// The context passed here does not need to contain a logger, one will be created.
func Init(require []lib.MustRequireFunc, impl Simple) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture, context.Context, id.ServiceId){
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
		slog.Error("launch failure on call Simple","error",t)
		lib.ExitSelf(ctx, 1, myId)
	})
	return smmap,launchF, ctx,myId
}
func Run(ctx context.Context,
	binding *lib.ServiceMethodMap, timeoutInMillis int32, bg lib.Backgrounder) syscall.KernelErr{
	defer func() {
		if r := recover(); r != nil {
			s, ok:=r.(string)
			if !ok && s!=apishared.ControlledExit {
				slog.Error("Run: trapped a panic in the guest side", "recovered", r)
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
			slog.Info("calling backgrounder of Simple")
			bg.Background(ctx)
			continue
		}
		if kerr == syscall.KernelErr_NoError {
			continue
		}
		break
	}
	slog.Error("error while waiting for  service calls", "error",syscall.KernelErr_name[int32(kerr)])
	return kerr
}
// Increase this value at your peril!
// Decreasing this value may make your overall program more responsive if you have many services.
var TimeoutInMillis = int32(50)

func ReadOneAndCall(ctx context.Context, binding *lib.ServiceMethodMap, 
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


func bind(ctx context.Context,sid id.ServiceId, impl Simple) (*lib.ServiceMethodMap, syscall.KernelErr) {
	smmap:=lib.NewServiceMethodMap()
	var mid id.MethodId
	var bindReq *syscall.BindMethodRequest
	var resp *syscall.BindMethodResponse
	var err syscall.KernelErr
//
// simple.v1.Simple.Get
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Get"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Simple","Get",
		GenerateGetInvoker(impl)) 
	return smmap,syscall.KernelErr_NoError
}
 

// Locate finds a reference to the client interface of simple.  
func Locate(ctx context.Context,sid id.ServiceId) (Client,syscall.KernelErr) {
    cs, kerr:=client.LocateDynamic(ctx, "simple.v1","simple", sid)
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
    normal:="unable to locate simple.v1.simple:"+name
    if err!=0 {
        if err == syscall.KernelErr_NotRequired {
            slog.Error("service was located, but it was not required")
            panic("locate attempted on a service that was not required")
        }
        panic(normal)
    }
    return result
}


func Register() (id.ServiceId, syscall.KernelErr){
    req := &syscall.RegisterRequest{}
	debugName:=fmt.Sprintf("%s.%s","simple.v1","simple")
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
func MustRegister() (context.Context,id.ServiceId) {
    sid, err:=Register()
    if err!=syscall.KernelErr_NoError {
        slog.Error("unable to register","package","simple.v1","service name","simple")
        panic("unable to register "+"simple")
    }
    return context.Background(), sid
}

func MustRequire(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Require1(ctx, "simple.v1","simple",sid)
    if err!=syscall.KernelErr_NoError {
        if err==syscall.KernelErr_DependencyCycle{
            slog.Error("unable to require because it creates a dependcy loop","package","simple.v1","service name","simple","error",syscall.KernelErr_name[int32(err)])
            panic("require simple.v1.simple creates a dependency loop")
        }
        slog.Error("unable to require","package","simple.v1","service name","simple","error",syscall.KernelErr_name[int32(err)])
        panic("not able to require simple.v1.simple:"+syscall.KernelErr_name[int32(err)])
    }
}

func MustExport(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Export1(ctx,"simple.v1","simple",sid)
    if err!=syscall.KernelErr_NoError{
        slog.Error("unable to export","package","simple.v1","service name","simple")
        panic("not able to export simple.v1.simple:"+syscall.KernelErr_name[int32(err)])
    }
}

func LaunchService(ctx context.Context, sid id.ServiceId, impl Simple) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture,syscall.KernelErr) {

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
	fut:=syscallguest.Launch(ctx,req)


    return smmap,fut,syscall.KernelErr_NoError

}

func MustLaunchService(ctx context.Context, sid id.ServiceId, impl Simple) (*lib.ServiceMethodMap, *syscallguest.LaunchFuture) {
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

//go:wasmimport simple get_
func Get_(int32,int32,int32,int32) int64
func GetHost(ctx context.Context,inPtr *http.GetRequest) *FutureGet {
	outProtoPtr := (*http.GetResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Get_)
	f:=NewFutureGet()
	f.CompleteMethod(ctx,ret,raw)
	return f
}  

// This is interface for invocation.

type invokeGet struct {
    fn func(context.Context,*http.GetRequest) *FutureGet
}

func (t *invokeGet) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&http.GetRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateGetInvoker(impl Simple) future.Invoker {
	return &invokeGet{fn:impl.Get} 
}  
