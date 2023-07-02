//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: queue/v1/queue.proto

package queue




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
 
func Launch(ctx context.Context, sid id.ServiceId, impl Queue) *future.Base[bool] {

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
		pcontext.Errorf(ctx,"Unable to start queue.v1.Queue, Ready returned false")
		pcontext.Dump(ctx)
		readyResult.Set(false)
	})

	return readyResult
}

func Init(ctx context.Context,require []lib.MustRequireFunc, impl Queue) *lib.ServiceMethodMap{
	defer func() {
		if r := recover(); r != nil {
			pcontext.Infof(ctx, "InitQueue: trapped a panic in the guest side: %v", r)
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
			pcontext.Errorf(ctx, "ready call on Queue failed")
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
			pcontext.Infof(ctx,"calling backgrounder of Queue")
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
		var a *anypb.Any
		if err:=a.MarshalFrom(result); err!=nil {
			pcontext.Errorf(ctx, "unable to marshal result for return value request")
			return
		}
		rvReq.Result = a
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

func bind(ctx context.Context,sid id.ServiceId, impl Queue) (*lib.ServiceMethodMap, syscall.KernelErr) {
	smmap:=lib.NewServiceMethodMap()
	var mid id.MethodId
	var bindReq *syscall.BindMethodRequest
	var resp *syscall.BindMethodResponse
	var err syscall.KernelErr
//
// queue.v1.Queue.CreateQueue
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = lib.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Open"
	resp, err=syscallguest.BindMethod(bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Queue","CreateQueue",
		GenerateCreateQueueInvoker(impl))
//
// queue.v1.Queue.Locate
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = lib.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Open"
	resp, err=syscallguest.BindMethod(bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Queue","Locate",
		GenerateLocateInvoker(impl))
//
// queue.v1.Queue.DeleteQueue
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = lib.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Open"
	resp, err=syscallguest.BindMethod(bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Queue","DeleteQueue",
		GenerateDeleteQueueInvoker(impl))
//
// queue.v1.Queue.Receive
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = lib.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Open"
	resp, err=syscallguest.BindMethod(bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Queue","Receive",
		GenerateReceiveInvoker(impl))
//
// queue.v1.Queue.MarkDone
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = lib.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Open"
	resp, err=syscallguest.BindMethod(bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Queue","MarkDone",
		GenerateMarkDoneInvoker(impl))
//
// queue.v1.Queue.Length
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = lib.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Open"
	resp, err=syscallguest.BindMethod(bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Queue","Length",
		GenerateLengthInvoker(impl))
//
// queue.v1.Queue.Send
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = lib.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Open"
	resp, err=syscallguest.BindMethod(bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"Queue","Send",
		GenerateSendInvoker(impl)) 
	pcontext.Dump(ctx)
	return smmap,syscall.KernelErr_NoError
}

// Locate finds a reference to the client interface of queue.  
func Locate(ctx context.Context,sid id.ServiceId) (Client,syscall.KernelErr) {
    cs, kerr:=client.LocateDynamic(ctx, "queue.v1","queue", sid)
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
    normal:="unable to locate queue.v1.queue:"+name
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
		PackagePath: "queue.v1",
		Service:     "queue",
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
        pcontext.Fatalf(ctx,"unable to register %s.%s","queue.v1","queue")
        panic("unable to register "+"queue")
    }
    return sid
}

func MustRequire(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Require1("queue.v1","queue",sid)
    if err!=syscall.KernelErr_NoError {
        if err==syscall.KernelErr_DependencyCycle{
            pcontext.Errorf(ctx,"unable to require %s.%s because it creates a dependcy loop: %s","queue.v1","queue",syscall.KernelErr_name[int32(err)])
            panic("require queue.v1.queue creates a dependency loop")
        }
        pcontext.Errorf(ctx,"unable to require %s.%s:%s","queue.v1","queue",syscall.KernelErr_name[int32(err)])
        panic("not able to require queue.v1.queue:"+syscall.KernelErr_name[int32(err)])
    }
}

func MustExport(ctx context.Context) {
    _, err:=lib.Export1("queue.v1","queue")
    if err!=syscall.KernelErr_NoError{
        pcontext.Fatalf(ctx, "unable to export %s.%s","queue.v1","queue")
        panic("not able to export queue.v1.queue:"+syscall.KernelErr_name[int32(err)])
    }
}

func WaitSatisfied(ctx context.Context, sid id.ServiceId, impl Queue) (*lib.ServiceMethodMap,syscall.KernelErr) {
	smmap, err:=bind(ctx,sid, impl)
	if err!=0{
		return  nil,syscall.KernelErr(err)
	}

    s:=sid.Marshal()
	syscallguest.Launch(&syscall.LaunchRequest{ServiceId:s })

    return smmap,syscall.KernelErr_NoError
}

func MustWaitSatisfied(ctx context.Context, sid id.ServiceId, impl Queue) *lib.ServiceMethodMap {
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

//go:wasmimport queue create_queue_
func CreateQueue_(int32,int32,int32,int32) int64
func CreateQueueHost(ctx context.Context,inPtr *CreateQueueRequest) *FutureCreateQueue {
	outProtoPtr := (*CreateQueueResponse)(nil)
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, CreateQueue_)
	if signal {
		pcontext.Infof(ctx, "CreateQueue exiting because of parigot signal")
		pcontext.Dump(ctx)
		syscallguest.Exit(1)
	}
	f:=NewFutureCreateQueue()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport queue locate_
func Locate_(int32,int32,int32,int32) int64
func LocateHost(ctx context.Context,inPtr *LocateRequest) *FutureLocate {
	outProtoPtr := (*LocateResponse)(nil)
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Locate_)
	if signal {
		pcontext.Infof(ctx, "Locate exiting because of parigot signal")
		pcontext.Dump(ctx)
		syscallguest.Exit(1)
	}
	f:=NewFutureLocate()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport queue delete_queue_
func DeleteQueue_(int32,int32,int32,int32) int64
func DeleteQueueHost(ctx context.Context,inPtr *DeleteQueueRequest) *FutureDeleteQueue {
	outProtoPtr := (*DeleteQueueResponse)(nil)
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, DeleteQueue_)
	if signal {
		pcontext.Infof(ctx, "DeleteQueue exiting because of parigot signal")
		pcontext.Dump(ctx)
		syscallguest.Exit(1)
	}
	f:=NewFutureDeleteQueue()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport queue receive_
func Receive_(int32,int32,int32,int32) int64
func ReceiveHost(ctx context.Context,inPtr *ReceiveRequest) *FutureReceive {
	outProtoPtr := (*ReceiveResponse)(nil)
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Receive_)
	if signal {
		pcontext.Infof(ctx, "Receive exiting because of parigot signal")
		pcontext.Dump(ctx)
		syscallguest.Exit(1)
	}
	f:=NewFutureReceive()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport queue mark_done_
func MarkDone_(int32,int32,int32,int32) int64
func MarkDoneHost(ctx context.Context,inPtr *MarkDoneRequest) *FutureMarkDone {
	outProtoPtr := (*MarkDoneResponse)(nil)
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, MarkDone_)
	if signal {
		pcontext.Infof(ctx, "MarkDone exiting because of parigot signal")
		pcontext.Dump(ctx)
		syscallguest.Exit(1)
	}
	f:=NewFutureMarkDone()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport queue length_
func Length_(int32,int32,int32,int32) int64
func LengthHost(ctx context.Context,inPtr *LengthRequest) *FutureLength {
	outProtoPtr := (*LengthResponse)(nil)
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Length_)
	if signal {
		pcontext.Infof(ctx, "Length exiting because of parigot signal")
		pcontext.Dump(ctx)
		syscallguest.Exit(1)
	}
	f:=NewFutureLength()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport queue send_
func Send_(int32,int32,int32,int32) int64
func SendHost(ctx context.Context,inPtr *SendRequest) *FutureSend {
	outProtoPtr := (*SendResponse)(nil)
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Send_)
	if signal {
		pcontext.Infof(ctx, "Send exiting because of parigot signal")
		pcontext.Dump(ctx)
		syscallguest.Exit(1)
	}
	f:=NewFutureSend()
	f.CompleteMethod(ctx,ret,raw)
	return f
}  

// This is interface for invocation.
type invokeCreateQueue struct {
    fn func(context.Context,*CreateQueueRequest) *FutureCreateQueue
}

func (t *invokeCreateQueue) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx CreateQueueRequest and 'CreateQueueRequest{}' why empty?
    in:=&CreateQueueRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateCreateQueueInvoker(impl Queue) future.Invoker {
	return &invokeCreateQueue{fn:impl.CreateQueue} 
}

// This is interface for invocation.
type invokeLocate struct {
    fn func(context.Context,*LocateRequest) *FutureLocate
}

func (t *invokeLocate) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx LocateRequest and 'LocateRequest{}' why empty?
    in:=&LocateRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateLocateInvoker(impl Queue) future.Invoker {
	return &invokeLocate{fn:impl.Locate} 
}

// This is interface for invocation.
type invokeDeleteQueue struct {
    fn func(context.Context,*DeleteQueueRequest) *FutureDeleteQueue
}

func (t *invokeDeleteQueue) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx DeleteQueueRequest and 'DeleteQueueRequest{}' why empty?
    in:=&DeleteQueueRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateDeleteQueueInvoker(impl Queue) future.Invoker {
	return &invokeDeleteQueue{fn:impl.DeleteQueue} 
}

// This is interface for invocation.
type invokeReceive struct {
    fn func(context.Context,*ReceiveRequest) *FutureReceive
}

func (t *invokeReceive) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx ReceiveRequest and 'ReceiveRequest{}' why empty?
    in:=&ReceiveRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateReceiveInvoker(impl Queue) future.Invoker {
	return &invokeReceive{fn:impl.Receive} 
}

// This is interface for invocation.
type invokeMarkDone struct {
    fn func(context.Context,*MarkDoneRequest) *FutureMarkDone
}

func (t *invokeMarkDone) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx MarkDoneRequest and 'MarkDoneRequest{}' why empty?
    in:=&MarkDoneRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateMarkDoneInvoker(impl Queue) future.Invoker {
	return &invokeMarkDone{fn:impl.MarkDone} 
}

// This is interface for invocation.
type invokeLength struct {
    fn func(context.Context,*LengthRequest) *FutureLength
}

func (t *invokeLength) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx LengthRequest and 'LengthRequest{}' why empty?
    in:=&LengthRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateLengthInvoker(impl Queue) future.Invoker {
	return &invokeLength{fn:impl.Length} 
}

// This is interface for invocation.
type invokeSend struct {
    fn func(context.Context,*SendRequest) *FutureSend
}

func (t *invokeSend) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx SendRequest and 'SendRequest{}' why empty?
    in:=&SendRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateSendInvoker(impl Queue) future.Invoker {
	return &invokeSend{fn:impl.Send} 
}  
