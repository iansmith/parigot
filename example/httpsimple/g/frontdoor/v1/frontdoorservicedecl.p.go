//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: frontdoor/v1/frontdoor.proto

package frontdoor


import(
    "context"

    "github.com/iansmith/parigot/g/httpconnector/v1"
  

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
// Frontdoor from frontdoor/v1/frontdoor.proto
//
//service interface
type Frontdoor interface {
    Handle(ctx context.Context,in *httpconnector.HandleRequest) *FutureHandle   
    Ready(context.Context,id.ServiceId) *future.Base[bool]
}

type Client interface {
    Handle(ctx context.Context,in *httpconnector.HandleRequest) *FutureHandle   
}

// Client difference from Frontdoor: Ready() 
type Client_ struct {
    *client.BaseService
}
// Check that Client_ is a Client.
var _ = Client(&Client_{})

//
// method: Frontdoor.Handle 
//
type FutureHandle struct {
    Method *future.Method[*httpconnector.HandleResponse,FrontdoorErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureHandle) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&httpconnector.HandleResponse{}
    if a!=nil {
        tmp, ok:=a.(*httpconnector.HandleResponse)
        if !ok {
            cvt:=a.(*anypb.Any)
            if cvt!=nil {
                if err:=cvt.UnmarshalTo(out); err!=nil {
                    return syscall.KernelErr_UnmarshalFailed
                }
            }
        } else {
            proto.Merge(out,tmp)
        }
    }
    f.Method.CompleteMethod(ctx,out,FrontdoorErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureHandle)Success(sfn func (proto.Message)) {
    x:=func(m *httpconnector.HandleResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureHandle)Failure(ffn func (int32)) {
    x:=func(err FrontdoorErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureHandle)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureHandle)Cancel()   {
    f.Method.Cancel()
}
func NewFutureHandle() *FutureHandle {
    f:=&FutureHandle{
        Method: future.NewMethod[*httpconnector.HandleResponse,FrontdoorErr](nil,nil),
    } 
    return f
}
func (i *Client_) Handle(ctx context.Context, in *httpconnector.HandleRequest) *FutureHandle { 
    mid, ok := i.BaseService.MethodIdByName("Handle")
    if !ok {
        f:=NewFutureHandle()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureHandle()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1,syscallguest.CurrentHostId())/*dispatch error*/
        return f
     }

    ctx, t:=lib.CurrentTime(ctx)
    source:=syscallguest.CurrentHostId()
    syscallguest.MatchCompleter(ctx,t,source,targetHid,cid,f)
    return f
}  
