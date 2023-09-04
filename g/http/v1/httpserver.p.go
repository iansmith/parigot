//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: http/v1/http.proto

package http




import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"runtime/debug"
    "unsafe"

 
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
 

func Launch(ctx context.Context, sid id.ServiceId, impl Http) *future.Base[bool] {

	readyResult:=future.NewBase[bool]()

	ready:=impl.Ready(ctx,sid)
	ready.Handle(func (b bool) {
		if b {
			readyResult.Set(true)			
			return
		}
		slog.Error("Unable to start http.v1.Http, Ready returned false")
		readyResult.Set(false)
	})

	return readyResult
}

// Note that  Init returns a future, but the case of failure is covered
// by this definition so the caller need only deal with Success case.
// The context passed here does not need to contain a logger, one will be created.
func Init(require []lib.MustRequireFunc, impl Http) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture, context.Context, id.ServiceId){ 
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
		slog.Error("launch failure on call Http","error",t)
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
				slog.Error("Run Http: trapped a panic in the guest side", "recovered", r)
				debug.PrintStack()
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
			slog.Info("calling backgrounder of Http")
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
		slog.Error("Http, readOneAndCall:unable to find binding for method on service, ignoring","mid",mid.Short(),"sid", sid.Short(),
			"current host",syscallguest.CurrentHostId().Short())
		return syscall.KernelErr_NoError
	}
	fut:=fn.Invoke(ctx,resp.GetParamOrResult())
	// if we get a nil, the intention is that the invocation be ignored
	if fut==nil {
		slog.Warn("ignoring call result for call","call",cid.Short())
		return syscall.KernelErr_NoError
	}
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


func bind(ctx context.Context,sid id.ServiceId, impl Http) (*lib.ServiceMethodMap, syscall.KernelErr) {
	smmap:=lib.NewServiceMethodMap()
	var mid id.MethodId
	var bindReq *syscall.BindMethodRequest
	var resp *syscall.BindMethodResponse
	var err syscall.KernelErr
//
// http.v1.Http.Get
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
	smmap.AddServiceMethod(sid,mid,"Http","Get",
		GenerateGetInvoker(impl))
//
// http.v1.Http.Post
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Post"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Http","Post",
		GeneratePostInvoker(impl))
//
// http.v1.Http.Put
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Put"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Http","Put",
		GeneratePutInvoker(impl))
//
// http.v1.Http.Delete
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Delete"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Http","Delete",
		GenerateDeleteInvoker(impl))
//
// http.v1.Http.Head
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Head"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Http","Head",
		GenerateHeadInvoker(impl))
//
// http.v1.Http.Options
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Options"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Http","Options",
		GenerateOptionsInvoker(impl))
//
// http.v1.Http.Patch
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Patch"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Http","Patch",
		GeneratePatchInvoker(impl))
//
// http.v1.Http.Connect
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Connect"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Http","Connect",
		GenerateConnectInvoker(impl))
//
// http.v1.Http.Trace
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Trace"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Http","Trace",
		GenerateTraceInvoker(impl)) 
	return smmap,syscall.KernelErr_NoError
}
 

// Locate finds a reference to the client interface of Http.  
func Locate(ctx context.Context,sid id.ServiceId) (Client,syscall.KernelErr) {
    cs, kerr:=client.LocateDynamic(ctx, "http.v1","http", sid)
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
    normal:="unable to locate http.v1.http:"+name
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
	debugName:=fmt.Sprintf("%s.%s","http.v1","http")
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
        slog.Error("unable to register","package","http.v1","service name","http")
        panic("unable to register "+"http")
    }
    return context.Background(), sid
}

func MustRequire(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Require1(ctx, "http.v1","http",sid)
    if err!=syscall.KernelErr_NoError {
        if err==syscall.KernelErr_DependencyCycle{
            slog.Error("unable to require because it creates a dependcy loop","package","http.v1","service name","http","error",syscall.KernelErr_name[int32(err)])
            panic("require http.v1.http creates a dependency loop")
        }
        slog.Error("unable to require","package","http.v1","service name","http","error",syscall.KernelErr_name[int32(err)])
        panic("not able to require http.v1.http:"+syscall.KernelErr_name[int32(err)])
    }
}

func MustExport(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Export1(ctx,"http.v1","http",sid)
    if err!=syscall.KernelErr_NoError{
        slog.Error("unable to export","package","http.v1","service name","http")
        panic("not able to export http.v1.http:"+syscall.KernelErr_name[int32(err)])
    }
}


func LaunchService(ctx context.Context, sid id.ServiceId, impl  Http) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture,syscall.KernelErr) {
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
func MustLaunchService(ctx context.Context, sid id.ServiceId, impl Http) (*lib.ServiceMethodMap, *syscallguest.LaunchFuture) {
 
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

//go:wasmimport http get_
func Get_(int32,int32,int32,int32) int64
func GetHost(ctx context.Context,inPtr *GetRequest) *FutureGet {
	outProtoPtr := (*GetResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Get_)
	f:=NewFutureGet()
	f.CompleteMethod(ctx,ret,raw, syscallguest.CurrentHostId())
	return f
} 

//go:wasmimport http post_
func Post_(int32,int32,int32,int32) int64
func PostHost(ctx context.Context,inPtr *PostRequest) *FuturePost {
	outProtoPtr := (*PostResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Post_)
	f:=NewFuturePost()
	f.CompleteMethod(ctx,ret,raw, syscallguest.CurrentHostId())
	return f
} 

//go:wasmimport http put_
func Put_(int32,int32,int32,int32) int64
func PutHost(ctx context.Context,inPtr *PutRequest) *FuturePut {
	outProtoPtr := (*PutResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Put_)
	f:=NewFuturePut()
	f.CompleteMethod(ctx,ret,raw, syscallguest.CurrentHostId())
	return f
} 

//go:wasmimport http delete_
func Delete_(int32,int32,int32,int32) int64
func DeleteHost(ctx context.Context,inPtr *DeleteRequest) *FutureDelete {
	outProtoPtr := (*DeleteResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Delete_)
	f:=NewFutureDelete()
	f.CompleteMethod(ctx,ret,raw, syscallguest.CurrentHostId())
	return f
} 

//go:wasmimport http head_
func Head_(int32,int32,int32,int32) int64
func HeadHost(ctx context.Context,inPtr *HeadRequest) *FutureHead {
	outProtoPtr := (*HeadResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Head_)
	f:=NewFutureHead()
	f.CompleteMethod(ctx,ret,raw, syscallguest.CurrentHostId())
	return f
} 

//go:wasmimport http options_
func Options_(int32,int32,int32,int32) int64
func OptionsHost(ctx context.Context,inPtr *OptionsRequest) *FutureOptions {
	outProtoPtr := (*OptionsResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Options_)
	f:=NewFutureOptions()
	f.CompleteMethod(ctx,ret,raw, syscallguest.CurrentHostId())
	return f
} 

//go:wasmimport http patch_
func Patch_(int32,int32,int32,int32) int64
func PatchHost(ctx context.Context,inPtr *PatchRequest) *FuturePatch {
	outProtoPtr := (*PatchResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Patch_)
	f:=NewFuturePatch()
	f.CompleteMethod(ctx,ret,raw, syscallguest.CurrentHostId())
	return f
} 

//go:wasmimport http connect_
func Connect_(int32,int32,int32,int32) int64
func ConnectHost(ctx context.Context,inPtr *ConnectRequest) *FutureConnect {
	outProtoPtr := (*ConnectResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Connect_)
	f:=NewFutureConnect()
	f.CompleteMethod(ctx,ret,raw, syscallguest.CurrentHostId())
	return f
} 

//go:wasmimport http trace_
func Trace_(int32,int32,int32,int32) int64
func TraceHost(ctx context.Context,inPtr *TraceRequest) *FutureTrace {
	outProtoPtr := (*TraceResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Trace_)
	f:=NewFutureTrace()
	f.CompleteMethod(ctx,ret,raw, syscallguest.CurrentHostId())
	return f
}   

// This is interface for invocation.

type invokeGet struct {
    fn func(context.Context,*GetRequest) *FutureGet
}

func (t *invokeGet) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&GetRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateGetInvoker(impl Http) future.Invoker {
	return &invokeGet{fn:impl.Get} 
}

// This is interface for invocation.

type invokePost struct {
    fn func(context.Context,*PostRequest) *FuturePost
}

func (t *invokePost) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&PostRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GeneratePostInvoker(impl Http) future.Invoker {
	return &invokePost{fn:impl.Post} 
}

// This is interface for invocation.

type invokePut struct {
    fn func(context.Context,*PutRequest) *FuturePut
}

func (t *invokePut) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&PutRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GeneratePutInvoker(impl Http) future.Invoker {
	return &invokePut{fn:impl.Put} 
}

// This is interface for invocation.

type invokeDelete struct {
    fn func(context.Context,*DeleteRequest) *FutureDelete
}

func (t *invokeDelete) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&DeleteRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateDeleteInvoker(impl Http) future.Invoker {
	return &invokeDelete{fn:impl.Delete} 
}

// This is interface for invocation.

type invokeHead struct {
    fn func(context.Context,*HeadRequest) *FutureHead
}

func (t *invokeHead) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&HeadRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateHeadInvoker(impl Http) future.Invoker {
	return &invokeHead{fn:impl.Head} 
}

// This is interface for invocation.

type invokeOptions struct {
    fn func(context.Context,*OptionsRequest) *FutureOptions
}

func (t *invokeOptions) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&OptionsRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateOptionsInvoker(impl Http) future.Invoker {
	return &invokeOptions{fn:impl.Options} 
}

// This is interface for invocation.

type invokePatch struct {
    fn func(context.Context,*PatchRequest) *FuturePatch
}

func (t *invokePatch) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&PatchRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GeneratePatchInvoker(impl Http) future.Invoker {
	return &invokePatch{fn:impl.Patch} 
}

// This is interface for invocation.

type invokeConnect struct {
    fn func(context.Context,*ConnectRequest) *FutureConnect
}

func (t *invokeConnect) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&ConnectRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateConnectInvoker(impl Http) future.Invoker {
	return &invokeConnect{fn:impl.Connect} 
}

// This is interface for invocation.

type invokeTrace struct {
    fn func(context.Context,*TraceRequest) *FutureTrace
}

func (t *invokeTrace) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&TraceRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateTraceInvoker(impl Http) future.Invoker {
	return &invokeTrace{fn:impl.Trace} 
}  
