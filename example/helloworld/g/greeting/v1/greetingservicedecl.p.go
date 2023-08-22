//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: greeting/v1/greeting.proto

package greeting


import(
    "context" 

    // this set of imports is _unrelated_ to the particulars of what the .proto imported... those are above
    "github.com/iansmith/parigot/lib/go"  
    "github.com/iansmith/parigot/lib/go/future"  
    "github.com/iansmith/parigot/lib/go/client"  
    "github.com/iansmith/parigot/api/shared/id"
    syscall "github.com/iansmith/parigot/g/syscall/v1" 
    syscallguest "github.com/iansmith/parigot/api/guest/syscall" 

    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/types/known/anypb"



)  
//
// Greeting from greeting/v1/greeting.proto
//
//service interface
type Greeting interface {
    FetchGreeting(ctx context.Context,in *FetchGreetingRequest) *FutureFetchGreeting   
    Ready(context.Context,id.ServiceId) *future.Base[bool]
}

type Client interface {
    FetchGreeting(ctx context.Context,in *FetchGreetingRequest) *FutureFetchGreeting   
}

// Client difference from Greeting: Ready() 
type Client_ struct {
    *client.BaseService
}
// Check that Client_ is a Client.
var _ = Client(&Client_{})

//
// method: Greeting.FetchGreeting 
//
type FutureFetchGreeting struct {
    Method *future.Method[*FetchGreetingResponse,GreetErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureFetchGreeting) CompleteMethod(ctx context.Context,a proto.Message, e int32) syscall.KernelErr{
    out:=&FetchGreetingResponse{}
    if a!=nil {
        if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
            return syscall.KernelErr_UnmarshalFailed
        }
    }
    f.Method.CompleteMethod(ctx,out,GreetErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureFetchGreeting)Success(sfn func (proto.Message)) {
    x:=func(m *FetchGreetingResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureFetchGreeting)Failure(ffn func (int32)) {
    x:=func(err GreetErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureFetchGreeting)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureFetchGreeting)Cancel()   {
    f.Method.Cancel()
}
func NewFutureFetchGreeting() *FutureFetchGreeting {
    f:=&FutureFetchGreeting{
        Method: future.NewMethod[*FetchGreetingResponse,GreetErr](nil,nil),
    } 
    return f
}
func (i *Client_) FetchGreeting(ctx context.Context, in *FetchGreetingRequest) *FutureFetchGreeting { 
    mid, ok := i.BaseService.MethodIdByName("FetchGreeting")
    if !ok {
        f:=NewFutureFetchGreeting()
        f.CompleteMethod(ctx,nil,1)/*dispatch error*/
    }
    _,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureFetchGreeting()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1)/*dispatch error*/
        return f
     }

    ctx, t:=lib.CurrentTime(ctx)
    syscallguest.MatchCompleter(ctx,t,syscallguest.CurrentHostId(),cid,f)
    return f
}  
