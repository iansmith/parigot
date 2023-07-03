//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: methodcall/bar/v1/bar.proto

package bar




import (
	"context"
    "unsafe" 
    // this set of imports is _unrelated_ to the particulars of what the .proto imported... those are above
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"  
	pcontext "github.com/iansmith/parigot/context"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/lib/go/future"
	"github.com/iansmith/parigot/lib/go/client"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/proto"

)
var _ =  unsafe.Sizeof([]byte{})
 
func Launch(ctx context.Context, sid id.ServiceId, impl Bar) *future.Base[bool] {

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
		pcontext.Errorf(ctx,"Unable to start methodcall.bar.v1.Bar, Ready returned false")
		pcontext.Dump(ctx)
		readyResult.Set(false)
	})

	return readyResult
}

func Init(ctx context.Context,require []lib.MustRequireFunc, impl Bar) *lib.ServiceMethodMap{
	defer func() {
		if r := recover(); r != nil {
			pcontext.Infof(ctx, "InitBar: trapped a panic in the guest side: %v", r)
		}
		pcontext.Dump(ctx)
	}()

	myId := MustRegister(ctx)
	MustExport(ctx)
	if len(require)>0 {
		for _, f := range require {
			f(ctx, myId)
		}
	}
	smmap:=MustWaitSatisfied(ctx, myId, impl)
	launchF:=Launch(ctx, myId, impl)

	// kinda tricky: if this get resolved this exit occurs on the
	// call stack of the code that called Set().
	launchF.Handle(func (ready bool) {
		if !ready {
			pcontext.Errorf(ctx, "ready call on Bar failed")
			syscallguest.Exit(1)
		}
	})
	return smmap
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
			pcontext.Infof(ctx,"calling backgrounder of Bar")
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

var TimeoutInMillis = int32(500)

func ReadOneAndCall(ctx context.Context, binding *lib.ServiceMethodMap, 
	timeoutInMillis int32) syscall.KernelErr{
	req:=syscall.ReadOneRequest{}

	// makes a copy
	for _, c := range binding.Call() {
		req.Call=append(req.Call, c)
	}

	req.TimeoutInMillis = timeoutInMillis
	req.HostId = lib.CurrentHostId().Marshal()
	resp, err:=syscallguest.ReadOne(&req)
	if err!=syscall.KernelErr_NoError {
		return err
	}
	// is timeout?
	if resp.Timeout {
		return syscall.KernelErr_ReadOneTimeout
	}

	// check for finished futures from within our address space
	lib.ExpireMethod(ctx)

	// is a promise being completed that was fulfilled somewhere else
	if r:=resp.GetResolved(); r!=nil {
		cid:=id.UnmarshalCallId(r.GetCallId())
		lib.CompleteCall(ctx, cid,r.GetResult(), r.GetResultError())
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
		rvReq.HostId= lib.CurrentHostId().Marshal()
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
		rvReq.HostId= lib.CurrentHostId().Marshal()
		rvReq.ResultError = err
		syscallguest.ReturnValue(rvReq) // nowhere for return value to go
	})
	return syscall.KernelErr_NoError

}

func bind(ctx context.Context,sid id.ServiceId, impl Bar) (*lib.ServiceMethodMap, syscall.KernelErr) {
	smmap:=lib.NewServiceMethodMap()
	var mid id.MethodId
	var bindReq *syscall.BindMethodRequest
	var resp *syscall.BindMethodResponse
	var err syscall.KernelErr
//
// methodcall.bar.v1.Bar.Accumulate
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = lib.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Accumulate"
	resp, err=syscallguest.BindMethod(bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Bar","Accumulate",
		GenerateAccumulateInvoker(impl)) 
	pcontext.Dump(ctx)
	return smmap,syscall.KernelErr_NoError
}

// Locate finds a reference to the client interface of bar.  
func Locate(ctx context.Context,sid id.ServiceId) (Client,syscall.KernelErr) {
    cs, kerr:=client.LocateDynamic(ctx, "methodcall.bar.v1","bar", sid)
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
    normal:="unable to locate methodcall.bar.v1.bar:"+name
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
		PackagePath: "methodcall.bar.v1",
		Service:     "bar",
	}
	req.Fqs = fqs
    req.IsClient = false

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
        pcontext.Fatalf(ctx,"unable to register %s.%s","methodcall.bar.v1","bar")
        panic("unable to register "+"bar")
    }
    return sid
}

func MustRequire(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Require1("methodcall.bar.v1","bar",sid)
    if err!=syscall.KernelErr_NoError {
        if err==syscall.KernelErr_DependencyCycle{
            pcontext.Errorf(ctx,"unable to require %s.%s because it creates a dependcy loop: %s","methodcall.bar.v1","bar",syscall.KernelErr_name[int32(err)])
            panic("require methodcall.bar.v1.bar creates a dependency loop")
        }
        pcontext.Errorf(ctx,"unable to require %s.%s:%s","methodcall.bar.v1","bar",syscall.KernelErr_name[int32(err)])
        panic("not able to require methodcall.bar.v1.bar:"+syscall.KernelErr_name[int32(err)])
    }
}

func MustExport(ctx context.Context) {
    _, err:=lib.Export1("methodcall.bar.v1","bar")
    if err!=syscall.KernelErr_NoError{
        pcontext.Fatalf(ctx, "unable to export %s.%s","methodcall.bar.v1","bar")
        panic("not able to export methodcall.bar.v1.bar:"+syscall.KernelErr_name[int32(err)])
    }
}

func WaitSatisfied(ctx context.Context, sid id.ServiceId, impl Bar) (*lib.ServiceMethodMap,syscall.KernelErr) {
	smmap, err:=bind(ctx,sid, impl)
	if err!=0{
		return  nil,syscall.KernelErr(err)
	}

    s:=sid.Marshal()
	syscallguest.Launch(&syscall.LaunchRequest{ServiceId:s })

    return smmap,syscall.KernelErr_NoError
}

func MustWaitSatisfied(ctx context.Context, sid id.ServiceId, impl Bar) *lib.ServiceMethodMap {
    smmap,err:=WaitSatisfied(ctx,sid,impl)
    if err!=syscall.KernelErr_NoError {
        panic("Unable to call WaitSatisfied successfully: "+syscall.KernelErr_name[int32(err)])
    }
    return smmap
}


// If you want to implement part of your server in host cost you should call 
// <methodName>Host from your server implementation. These will be optimized 
// away by the compiler if you don't use them--in other words, if you want to 
// implement everything on the guest side).
// 

//go:wasmimport bar accumulate_
func Accumulate_(int32,int32,int32,int32) int64
func AccumulateHost(ctx context.Context,inPtr *AccumulateRequest) *FutureAccumulate {
	outProtoPtr := (*AccumulateResponse)(nil)
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Accumulate_)
	if signal {
		pcontext.Infof(ctx, "Accumulate exiting because of parigot signal")
		pcontext.Dump(ctx)
		syscallguest.Exit(1)
	}
	f:=NewFutureAccumulate()
	f.CompleteMethod(ctx,ret,raw)
	return f
}  

// This is interface for invocation.
type invokeAccumulate struct {
    fn func(context.Context,*AccumulateRequest) *FutureAccumulate
}

func (t *invokeAccumulate) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx AccumulateRequest and 'AccumulateRequest{}' why empty?
    in:=&AccumulateRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateAccumulateInvoker(impl Bar) future.Invoker {
	return &invokeAccumulate{fn:impl.Accumulate} 
}  
