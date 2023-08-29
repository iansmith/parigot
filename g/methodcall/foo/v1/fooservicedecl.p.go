//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: methodcall/foo/v1/foo.proto

package foo


import(
    "context" 

// no method? false

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
// Foo from methodcall/foo/v1/foo.proto
//
//service interface
type Foo interface {
    AddMultiply(ctx context.Context,in *AddMultiplyRequest) *FutureAddMultiply   
    LucasSequence(ctx context.Context) *FutureLucasSequence  
    WritePi(ctx context.Context,in *WritePiRequest) *FutureWritePi   
    Ready(context.Context,id.ServiceId) *future.Base[bool]
}

type Client interface {
    AddMultiply(ctx context.Context,in *AddMultiplyRequest) *FutureAddMultiply   
    LucasSequence(ctx context.Context) *FutureLucasSequence  
    WritePi(ctx context.Context,in *WritePiRequest) *FutureWritePi   
}

// Client difference from Foo: Ready() 
type Client_ struct {
    *client.BaseService
}
// Check that Client_ is a Client.
var _ = Client(&Client_{})

//
// method: Foo.AddMultiply 
//
type FutureAddMultiply struct {
    Method *future.Method[*AddMultiplyResponse,FooErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureAddMultiply) CompleteMethod(ctx context.Context,a proto.Message, e int32) syscall.KernelErr{
    out:=&AddMultiplyResponse{}
    if a!=nil {
        if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
            return syscall.KernelErr_UnmarshalFailed
        }
    }
    f.Method.CompleteMethod(ctx,out,FooErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureAddMultiply)Success(sfn func (proto.Message)) {
    x:=func(m *AddMultiplyResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureAddMultiply)Failure(ffn func (int32)) {
    x:=func(err FooErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureAddMultiply)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureAddMultiply)Cancel()   {
    f.Method.Cancel()
}
func NewFutureAddMultiply() *FutureAddMultiply {
    f:=&FutureAddMultiply{
        Method: future.NewMethod[*AddMultiplyResponse,FooErr](nil,nil),
    } 
    return f
}
func (i *Client_) AddMultiply(ctx context.Context, in *AddMultiplyRequest) *FutureAddMultiply { 
    mid, ok := i.BaseService.MethodIdByName("AddMultiply")
    if !ok {
        f:=NewFutureAddMultiply()
        f.CompleteMethod(ctx,nil,1)/*dispatch error*/
    }
    _,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureAddMultiply()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1)/*dispatch error*/
        return f
     }

    ctx, t:=lib.CurrentTime(ctx)
    syscallguest.MatchCompleter(ctx,t,syscallguest.CurrentHostId(),cid,f)
    return f
}

//
// method: Foo.LucasSequence 
//
type FutureLucasSequence struct {
    Method *future.Method[*LucasSequenceResponse,FooErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureLucasSequence) CompleteMethod(ctx context.Context,a proto.Message, e int32) syscall.KernelErr{
    out:=&LucasSequenceResponse{}
    if a!=nil {
        if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
            return syscall.KernelErr_UnmarshalFailed
        }
    }
    f.Method.CompleteMethod(ctx,out,FooErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureLucasSequence)Success(sfn func (proto.Message)) {
    x:=func(m *LucasSequenceResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureLucasSequence)Failure(ffn func (int32)) {
    x:=func(err FooErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureLucasSequence)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureLucasSequence)Cancel()   {
    f.Method.Cancel()
}
func NewFutureLucasSequence() *FutureLucasSequence {
    f:=&FutureLucasSequence{
        Method: future.NewMethod[*LucasSequenceResponse,FooErr](nil,nil),
    } 
    return f
}
func (i *Client_) LucasSequence(ctx context.Context) *FutureLucasSequence { 
    mid, ok := i.BaseService.MethodIdByName("LucasSequence")
    if !ok {
        f:=NewFutureLucasSequence()
        f.CompleteMethod(ctx,nil,1)/*dispatch error*/
    }
    _,cid,kerr:= i.BaseService.Dispatch(ctx,mid,nil) 
    f:=NewFutureLucasSequence()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1)/*dispatch error*/
        return f
     }

    ctx, t:=lib.CurrentTime(ctx)
    syscallguest.MatchCompleter(ctx,t,syscallguest.CurrentHostId(),cid,f)
    return f
}

//
// method: Foo.WritePi 
//
type FutureWritePi struct {
    Base *future.Base[FooErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureWritePi) CompleteMethod(ctx context.Context,a proto.Message, e int32) syscall.KernelErr{
    f.Base.Set(FooErr(e)) 
    return syscall.KernelErr_NoError

} 
func (f *FutureWritePi)Success(sfn func (proto.Message)) {
    // no way for this to be called
} 

func (f *FutureWritePi)Failure(ffn func (int32)) {
    x:=func(err FooErr) {
        ffn(int32(err))
    }
    f.Base.Handle(x) 
}

func (f *FutureWritePi)Completed() bool  {
    return f.Base.Completed()

}
func (f *FutureWritePi)Cancel()   {
    f.Base.Cancel()
}
func NewFutureWritePi() *FutureWritePi {
    f:=&FutureWritePi{
        Base: future.NewBase[FooErr](),
    } 
    return f
}
func (i *Client_) WritePi(ctx context.Context, in *WritePiRequest) *FutureWritePi { 
    mid, ok := i.BaseService.MethodIdByName("WritePi")
    if !ok {
        f:=NewFutureWritePi()
        f.CompleteMethod(ctx,nil,1)/*dispatch error*/
    }
    _,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureWritePi()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1)/*dispatch error*/
        return f
     }

    ctx, t:=lib.CurrentTime(ctx)
    syscallguest.MatchCompleter(ctx,t,syscallguest.CurrentHostId(),cid,f)
    return f
}  
