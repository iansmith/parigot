//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: methodcall/bar/v1/bar.proto

package bar


import(
    "context" 

    // this set of imports is _unrelated_ to the particulars of what the .proto imported... those are above
    lib "github.com/iansmith/parigot/lib/go"  
    "github.com/iansmith/parigot/lib/go/future"  
    "github.com/iansmith/parigot/lib/go/client"  
    "github.com/iansmith/parigot/api/shared/id"
    syscall "github.com/iansmith/parigot/g/syscall/v1" 

    "google.golang.org/protobuf/proto"


)  
//
// Bar from methodcall/bar/v1/bar.proto
//
//service interface
type Bar interface {
    Accumulate(ctx context.Context,in *AccumulateRequest) *FutureAccumulate   
    Ready(context.Context,id.ServiceId) *future.Base[bool]
}

type Client interface {
    Accumulate(ctx context.Context,in *AccumulateRequest) *FutureAccumulate   
}

// Client difference from Bar: Ready() 
type Client_ struct {
    *client.BaseService
}
// Check that Client_ is a Client.
var _ = Client(&Client_{})

//
// method: Bar.Accumulate 
//
type FutureAccumulate struct {
    Method *future.Method[*AccumulateResponse,BarErr]
} 

// This is the same API for output needed or not because of the Completer interface.
func (f * FutureAccumulate) CompleteMethod(ctx context.Context,a proto.Message, e int32) {
    result:= a.(*AccumulateResponse)
    f.Method.CompleteMethod(ctx,result,BarErr(e)) 
}
func (f *FutureAccumulate)Success(sfn func (proto.Message)) {
    x:=func(m *AccumulateResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureAccumulate)Failure(ffn func (int32)) {
    x:=func(err BarErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}
func NewFutureAccumulate() *FutureAccumulate {
    f:=&FutureAccumulate{
        Method: future.NewMethod[*AccumulateResponse,BarErr](nil,nil),
    } 
    return f
}
func (i *Client_) Accumulate(ctx context.Context, in *AccumulateRequest) *FutureAccumulate { 
    mid, ok := i.BaseService.MethodIdByName("Accumulate")
    if !ok {
        f:=NewFutureAccumulate()
        f.CompleteMethod(ctx,nil,1)/*dispatch error*/
    }
    cid,kerr:= i.BaseService.Dispatch(mid,in) 
    f:=NewFutureAccumulate()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1)/*dispatch error*/
        return f
     }
    lib.MatchCompleter(cid,f)
    return f
}  
