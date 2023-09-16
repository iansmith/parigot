//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: queue/v1/queue.proto

package queue


import(
    "context" 

    "github.com/iansmith/parigot/lib/go/future"  
    "github.com/iansmith/parigot/lib/go/client"  
    "github.com/iansmith/parigot/api/shared/id"
    "google.golang.org/protobuf/proto"
    syscallguest "github.com/iansmith/parigot/api/guest/syscall" 
    syscall "github.com/iansmith/parigot/g/syscall/v1" 
    "github.com/iansmith/parigot/lib/go"  
    "google.golang.org/protobuf/types/known/anypb"


)  
//
// Queue from queue/v1/queue.proto
//
//service interface
type Queue interface {
    CreateQueue(ctx context.Context,in *CreateQueueRequest) *FutureCreateQueue  
    Locate(ctx context.Context,in *LocateRequest) *FutureLocate  
    DeleteQueue(ctx context.Context,in *DeleteQueueRequest) *FutureDeleteQueue  
    Receive(ctx context.Context,in *ReceiveRequest) *FutureReceive  
    MarkDone(ctx context.Context,in *MarkDoneRequest) *FutureMarkDone  
    Length(ctx context.Context,in *LengthRequest) *FutureLength  
    Send(ctx context.Context,in *SendRequest) *FutureSend   
    Ready(context.Context,id.ServiceId) *future.Base[bool]
}

type Client interface {
    CreateQueue(ctx context.Context,in *CreateQueueRequest) *FutureCreateQueue  
    Locate(ctx context.Context,in *LocateRequest) *FutureLocate  
    DeleteQueue(ctx context.Context,in *DeleteQueueRequest) *FutureDeleteQueue  
    Receive(ctx context.Context,in *ReceiveRequest) *FutureReceive  
    MarkDone(ctx context.Context,in *MarkDoneRequest) *FutureMarkDone  
    Length(ctx context.Context,in *LengthRequest) *FutureLength  
    Send(ctx context.Context,in *SendRequest) *FutureSend   
}

// Client difference from Queue: Ready() 
type Client_ struct {
    *client.BaseService
}
// Check that Client_ is a Client.
var _ = Client(&Client_{})

//
// method: Queue.CreateQueue 
//
type FutureCreateQueue struct {
    Method *future.Method[*CreateQueueResponse,QueueErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureCreateQueue) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&CreateQueueResponse{}
    if a!=nil {
        tmp, ok:=a.(*CreateQueueResponse)
        if !ok {
            if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
                return syscall.KernelErr_UnmarshalFailed
            }
        } else {
            proto.Merge(out,tmp)
        }
    }
    f.Method.CompleteMethod(ctx,out,QueueErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureCreateQueue)Success(sfn func (proto.Message)) {
    x:=func(m *CreateQueueResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureCreateQueue)Failure(ffn func (int32)) {
    x:=func(err QueueErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureCreateQueue)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureCreateQueue)Cancel()   {
    f.Method.Cancel()
}
func NewFutureCreateQueue() *FutureCreateQueue {
    f:=&FutureCreateQueue{
        Method: future.NewMethod[*CreateQueueResponse,QueueErr](nil,nil),
    } 
    return f
}
func (i *Client_) CreateQueue(ctx context.Context, in *CreateQueueRequest) *FutureCreateQueue { 
    mid, ok := i.BaseService.MethodIdByName("CreateQueue")
    if !ok {
        f:=NewFutureCreateQueue()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureCreateQueue()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1,syscallguest.CurrentHostId())/*dispatch error*/
        return f
     }

    ctx, t:=lib.CurrentTime(ctx)
    source:=syscallguest.CurrentHostId()
    syscallguest.MatchCompleter(ctx,t,source,targetHid,cid,f)
    return f
}

//
// method: Queue.Locate 
//
type FutureLocate struct {
    Method *future.Method[*LocateResponse,QueueErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureLocate) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&LocateResponse{}
    if a!=nil {
        tmp, ok:=a.(*LocateResponse)
        if !ok {
            if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
                return syscall.KernelErr_UnmarshalFailed
            }
        } else {
            proto.Merge(out,tmp)
        }
    }
    f.Method.CompleteMethod(ctx,out,QueueErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureLocate)Success(sfn func (proto.Message)) {
    x:=func(m *LocateResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureLocate)Failure(ffn func (int32)) {
    x:=func(err QueueErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureLocate)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureLocate)Cancel()   {
    f.Method.Cancel()
}
func NewFutureLocate() *FutureLocate {
    f:=&FutureLocate{
        Method: future.NewMethod[*LocateResponse,QueueErr](nil,nil),
    } 
    return f
}
func (i *Client_) Locate(ctx context.Context, in *LocateRequest) *FutureLocate { 
    mid, ok := i.BaseService.MethodIdByName("Locate")
    if !ok {
        f:=NewFutureLocate()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureLocate()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1,syscallguest.CurrentHostId())/*dispatch error*/
        return f
     }

    ctx, t:=lib.CurrentTime(ctx)
    source:=syscallguest.CurrentHostId()
    syscallguest.MatchCompleter(ctx,t,source,targetHid,cid,f)
    return f
}

//
// method: Queue.DeleteQueue 
//
type FutureDeleteQueue struct {
    Method *future.Method[*DeleteQueueResponse,QueueErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureDeleteQueue) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&DeleteQueueResponse{}
    if a!=nil {
        tmp, ok:=a.(*DeleteQueueResponse)
        if !ok {
            if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
                return syscall.KernelErr_UnmarshalFailed
            }
        } else {
            proto.Merge(out,tmp)
        }
    }
    f.Method.CompleteMethod(ctx,out,QueueErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureDeleteQueue)Success(sfn func (proto.Message)) {
    x:=func(m *DeleteQueueResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureDeleteQueue)Failure(ffn func (int32)) {
    x:=func(err QueueErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureDeleteQueue)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureDeleteQueue)Cancel()   {
    f.Method.Cancel()
}
func NewFutureDeleteQueue() *FutureDeleteQueue {
    f:=&FutureDeleteQueue{
        Method: future.NewMethod[*DeleteQueueResponse,QueueErr](nil,nil),
    } 
    return f
}
func (i *Client_) DeleteQueue(ctx context.Context, in *DeleteQueueRequest) *FutureDeleteQueue { 
    mid, ok := i.BaseService.MethodIdByName("DeleteQueue")
    if !ok {
        f:=NewFutureDeleteQueue()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureDeleteQueue()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1,syscallguest.CurrentHostId())/*dispatch error*/
        return f
     }

    ctx, t:=lib.CurrentTime(ctx)
    source:=syscallguest.CurrentHostId()
    syscallguest.MatchCompleter(ctx,t,source,targetHid,cid,f)
    return f
}

//
// method: Queue.Receive 
//
type FutureReceive struct {
    Method *future.Method[*ReceiveResponse,QueueErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureReceive) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&ReceiveResponse{}
    if a!=nil {
        tmp, ok:=a.(*ReceiveResponse)
        if !ok {
            if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
                return syscall.KernelErr_UnmarshalFailed
            }
        } else {
            proto.Merge(out,tmp)
        }
    }
    f.Method.CompleteMethod(ctx,out,QueueErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureReceive)Success(sfn func (proto.Message)) {
    x:=func(m *ReceiveResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureReceive)Failure(ffn func (int32)) {
    x:=func(err QueueErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureReceive)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureReceive)Cancel()   {
    f.Method.Cancel()
}
func NewFutureReceive() *FutureReceive {
    f:=&FutureReceive{
        Method: future.NewMethod[*ReceiveResponse,QueueErr](nil,nil),
    } 
    return f
}
func (i *Client_) Receive(ctx context.Context, in *ReceiveRequest) *FutureReceive { 
    mid, ok := i.BaseService.MethodIdByName("Receive")
    if !ok {
        f:=NewFutureReceive()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureReceive()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1,syscallguest.CurrentHostId())/*dispatch error*/
        return f
     }

    ctx, t:=lib.CurrentTime(ctx)
    source:=syscallguest.CurrentHostId()
    syscallguest.MatchCompleter(ctx,t,source,targetHid,cid,f)
    return f
}

//
// method: Queue.MarkDone 
//
type FutureMarkDone struct {
    Method *future.Method[*MarkDoneResponse,QueueErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureMarkDone) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&MarkDoneResponse{}
    if a!=nil {
        tmp, ok:=a.(*MarkDoneResponse)
        if !ok {
            if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
                return syscall.KernelErr_UnmarshalFailed
            }
        } else {
            proto.Merge(out,tmp)
        }
    }
    f.Method.CompleteMethod(ctx,out,QueueErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureMarkDone)Success(sfn func (proto.Message)) {
    x:=func(m *MarkDoneResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureMarkDone)Failure(ffn func (int32)) {
    x:=func(err QueueErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureMarkDone)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureMarkDone)Cancel()   {
    f.Method.Cancel()
}
func NewFutureMarkDone() *FutureMarkDone {
    f:=&FutureMarkDone{
        Method: future.NewMethod[*MarkDoneResponse,QueueErr](nil,nil),
    } 
    return f
}
func (i *Client_) MarkDone(ctx context.Context, in *MarkDoneRequest) *FutureMarkDone { 
    mid, ok := i.BaseService.MethodIdByName("MarkDone")
    if !ok {
        f:=NewFutureMarkDone()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureMarkDone()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1,syscallguest.CurrentHostId())/*dispatch error*/
        return f
     }

    ctx, t:=lib.CurrentTime(ctx)
    source:=syscallguest.CurrentHostId()
    syscallguest.MatchCompleter(ctx,t,source,targetHid,cid,f)
    return f
}

//
// method: Queue.Length 
//
type FutureLength struct {
    Method *future.Method[*LengthResponse,QueueErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureLength) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&LengthResponse{}
    if a!=nil {
        tmp, ok:=a.(*LengthResponse)
        if !ok {
            if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
                return syscall.KernelErr_UnmarshalFailed
            }
        } else {
            proto.Merge(out,tmp)
        }
    }
    f.Method.CompleteMethod(ctx,out,QueueErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureLength)Success(sfn func (proto.Message)) {
    x:=func(m *LengthResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureLength)Failure(ffn func (int32)) {
    x:=func(err QueueErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureLength)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureLength)Cancel()   {
    f.Method.Cancel()
}
func NewFutureLength() *FutureLength {
    f:=&FutureLength{
        Method: future.NewMethod[*LengthResponse,QueueErr](nil,nil),
    } 
    return f
}
func (i *Client_) Length(ctx context.Context, in *LengthRequest) *FutureLength { 
    mid, ok := i.BaseService.MethodIdByName("Length")
    if !ok {
        f:=NewFutureLength()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureLength()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1,syscallguest.CurrentHostId())/*dispatch error*/
        return f
     }

    ctx, t:=lib.CurrentTime(ctx)
    source:=syscallguest.CurrentHostId()
    syscallguest.MatchCompleter(ctx,t,source,targetHid,cid,f)
    return f
}

//
// method: Queue.Send 
//
type FutureSend struct {
    Method *future.Method[*SendResponse,QueueErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureSend) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&SendResponse{}
    if a!=nil {
        tmp, ok:=a.(*SendResponse)
        if !ok {
            if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
                return syscall.KernelErr_UnmarshalFailed
            }
        } else {
            proto.Merge(out,tmp)
        }
    }
    f.Method.CompleteMethod(ctx,out,QueueErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureSend)Success(sfn func (proto.Message)) {
    x:=func(m *SendResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureSend)Failure(ffn func (int32)) {
    x:=func(err QueueErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureSend)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureSend)Cancel()   {
    f.Method.Cancel()
}
func NewFutureSend() *FutureSend {
    f:=&FutureSend{
        Method: future.NewMethod[*SendResponse,QueueErr](nil,nil),
    } 
    return f
}
func (i *Client_) Send(ctx context.Context, in *SendRequest) *FutureSend { 
    mid, ok := i.BaseService.MethodIdByName("Send")
    if !ok {
        f:=NewFutureSend()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureSend()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1,syscallguest.CurrentHostId())/*dispatch error*/
        return f
     }

    ctx, t:=lib.CurrentTime(ctx)
    source:=syscallguest.CurrentHostId()
    syscallguest.MatchCompleter(ctx,t,source,targetHid,cid,f)
    return f
}  
