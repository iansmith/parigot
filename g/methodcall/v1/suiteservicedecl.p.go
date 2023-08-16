//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: methodcall/v1/suite.proto

package methodcall


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
// MethodCallSuite from methodcall/v1/suite.proto
//
//service interface
type MethodCallSuite interface {
    Exec(ctx context.Context,in *ExecRequest) *FutureExec  
    SuiteReport(ctx context.Context,in *SuiteReportRequest) *FutureSuiteReport   
    Ready(context.Context,id.ServiceId) *future.Base[bool]
}

type Client interface {
    Exec(ctx context.Context,in *ExecRequest) *FutureExec  
    SuiteReport(ctx context.Context,in *SuiteReportRequest) *FutureSuiteReport   
}

// Client difference from MethodCallSuite: Ready() 
type Client_ struct {
    *client.BaseService
}
// Check that Client_ is a Client.
var _ = Client(&Client_{})

//
// method: MethodCallSuite.Exec 
//
type FutureExec struct {
    Method *future.Method[*ExecResponse,MethodCallSuiteErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureExec) CompleteMethod(ctx context.Context,a proto.Message, e int32) syscall.KernelErr{
    out:=&ExecResponse{}
    if a!=nil {
        if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
            return syscall.KernelErr_UnmarshalFailed
        }
    }
    f.Method.CompleteMethod(ctx,out,MethodCallSuiteErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureExec)Success(sfn func (proto.Message)) {
    x:=func(m *ExecResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureExec)Failure(ffn func (int32)) {
    x:=func(err MethodCallSuiteErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureExec)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureExec)Cancel()   {
    f.Method.Cancel()
}
func NewFutureExec() *FutureExec {
    f:=&FutureExec{
        Method: future.NewMethod[*ExecResponse,MethodCallSuiteErr](nil,nil),
    } 
    return f
}
func (i *Client_) Exec(ctx context.Context, in *ExecRequest) *FutureExec { 
    mid, ok := i.BaseService.MethodIdByName("Exec")
    if !ok {
        f:=NewFutureExec()
        f.CompleteMethod(ctx,nil,1)/*dispatch error*/
    }
    _,cid,kerr:= i.BaseService.Dispatch(mid,in) 
    f:=NewFutureExec()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1)/*dispatch error*/
        return f
     }
    syscallguest.MatchCompleter(ctx,syscallguest.CurrentHostId(),cid,f)
    return f
}

//
// method: MethodCallSuite.SuiteReport 
//
type FutureSuiteReport struct {
    Base *future.Base[MethodCallSuiteErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureSuiteReport) CompleteMethod(ctx context.Context,a proto.Message, e int32) syscall.KernelErr{
    f.Base.Set(MethodCallSuiteErr(e)) 
    return syscall.KernelErr_NoError

} 
func (f *FutureSuiteReport)Success(sfn func (proto.Message)) {
    // no way for this to be called
} 

func (f *FutureSuiteReport)Failure(ffn func (int32)) {
    x:=func(err MethodCallSuiteErr) {
        ffn(int32(err))
    }
    f.Base.Handle(x) 
}

func (f *FutureSuiteReport)Completed() bool  {
    return f.Base.Completed()

}
func (f *FutureSuiteReport)Cancel()   {
    f.Base.Cancel()
}
func NewFutureSuiteReport() *FutureSuiteReport {
    f:=&FutureSuiteReport{
        Base: future.NewBase[MethodCallSuiteErr](),
    } 
    return f
}
func (i *Client_) SuiteReport(ctx context.Context, in *SuiteReportRequest) *FutureSuiteReport { 
    mid, ok := i.BaseService.MethodIdByName("SuiteReport")
    if !ok {
        f:=NewFutureSuiteReport()
        f.CompleteMethod(ctx,nil,1)/*dispatch error*/
    }
    _,cid,kerr:= i.BaseService.Dispatch(mid,in) 
    f:=NewFutureSuiteReport()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1)/*dispatch error*/
        return f
     }
    syscallguest.MatchCompleter(ctx,syscallguest.CurrentHostId(),cid,f)
    return f
}  
