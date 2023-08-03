//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: httpconnector/v1/httpconnector.proto

package httpconnector


import(
    "context" 

    // this set of imports is _unrelated_ to the particulars of what the .proto imported... those are above
    "github.com/iansmith/parigot/lib/go/future"  
    "github.com/iansmith/parigot/lib/go/client"  
    "github.com/iansmith/parigot/api/shared/id"
    syscall "github.com/iansmith/parigot/g/syscall/v1" 
    syscallguest "github.com/iansmith/parigot/api/guest/syscall" 

    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/types/known/anypb"



)  
//
// HttpConnector from httpconnector/v1/httpconnector.proto
//
//service interface
type HttpConnector interface {
    Check(ctx context.Context,in *CheckRequest) *FutureCheck   
    Ready(context.Context,id.ServiceId) *future.Base[bool]
}

type Client interface {
    Check(ctx context.Context,in *CheckRequest) *FutureCheck   
}

// Client difference from HttpConnector: Ready() 
type Client_ struct {
    *client.BaseService
}
// Check that Client_ is a Client.
var _ = Client(&Client_{})

//
// method: HttpConnector.Check 
//
type FutureCheck struct {
    Method *future.Method[*CheckResponse,HttpConnectorErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureCheck) CompleteMethod(ctx context.Context,a proto.Message, e int32) syscall.KernelErr{
    out:=&CheckResponse{}
    if a!=nil {
        if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
            return syscall.KernelErr_UnmarshalFailed
        }
    }
    f.Method.CompleteMethod(ctx,out,HttpConnectorErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureCheck)Success(sfn func (proto.Message)) {
    x:=func(m *CheckResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureCheck)Failure(ffn func (int32)) {
    x:=func(err HttpConnectorErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}
func NewFutureCheck() *FutureCheck {
    f:=&FutureCheck{
        Method: future.NewMethod[*CheckResponse,HttpConnectorErr](nil,nil),
    } 
    return f
}
func (i *Client_) Check(ctx context.Context, in *CheckRequest) *FutureCheck { 
    mid, ok := i.BaseService.MethodIdByName("Check")
    if !ok {
        f:=NewFutureCheck()
        f.CompleteMethod(ctx,nil,1)/*dispatch error*/
    }
    cid,kerr:= i.BaseService.Dispatch(mid,in) 
    f:=NewFutureCheck()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1)/*dispatch error*/
        return f
     }
    syscallguest.MatchCompleter(cid,f)
    return f
}  
