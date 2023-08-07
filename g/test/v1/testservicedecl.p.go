//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: test/v1/test.proto

package test


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
// Test from test/v1/test.proto
//
//service interface
type Test interface {
    TestAddTestSuite(ctx context.Context,in *AddTestSuiteRequest) *FutureTestAddTestSuite  
    TestStart(ctx context.Context,in *StartRequest) *FutureTestStart   
    Ready(context.Context,id.ServiceId) *future.Base[bool]
}

type ClientTest interface {
    TestAddTestSuite(ctx context.Context,in *AddTestSuiteRequest) *FutureTestAddTestSuite  
    TestStart(ctx context.Context,in *StartRequest) *FutureTestStart   
}

// ClientTest difference from Test: Ready() 
type ClientTest_ struct {
    *client.BaseService
}
// Check that Client_ is a Client.
var _ = ClientTest(&ClientTest_{})

//
// method: Test.AddTestSuite 
//
type FutureTestAddTestSuite struct {
    Method *future.Method[*AddTestSuiteResponse,TestErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureTestAddTestSuite) CompleteMethod(ctx context.Context,a proto.Message, e int32) syscall.KernelErr{
    out:=&AddTestSuiteResponse{}
    if a!=nil {
        if any, ok := a.(*anypb.Any); ok {
            if err:= any.UnmarshalTo(out); err!=nil {
                return syscall.KernelErr_UnmarshalFailed
            }
        } else {
            // `a` and `out` are the same type, so we can assign the values of a to out
            proto.Merge(out, a.(proto.Message))
        }
    }
    f.Method.CompleteMethod(ctx,out,TestErr(e)) 
    return syscall.KernelErr_NoError
}
func (f *FutureTestAddTestSuite)Success(sfn func (proto.Message)) {
    x:=func(m *AddTestSuiteResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureTestAddTestSuite)Failure(ffn func (int32)) {
    x:=func(err TestErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}
func NewFutureTestAddTestSuite() *FutureTestAddTestSuite {
    f:=&FutureTestAddTestSuite{
        Method: future.NewMethod[*AddTestSuiteResponse,TestErr](nil,nil),
    } 
    return f
}
func (i *ClientTest_) TestAddTestSuite(ctx context.Context, in *AddTestSuiteRequest) *FutureTestAddTestSuite { 
    mid, ok := i.BaseService.MethodIdByName("AddTestSuite")
    if !ok {
        f:=NewFutureTestAddTestSuite()
        f.CompleteMethod(ctx,nil,1)/*dispatch error*/
    }
    cid,kerr:= i.BaseService.Dispatch(mid,in) 
    f:=NewFutureTestAddTestSuite()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1)/*dispatch error*/
        return f
     }
    syscallguest.MatchCompleter(cid,f)
    return f
}

//
// method: Test.Start 
//
type FutureTestStart struct {
    Method *future.Method[*StartResponse,TestErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureTestStart) CompleteMethod(ctx context.Context,a proto.Message, e int32) syscall.KernelErr{
    out:=&StartResponse{}
    if a!=nil {
        if any, ok := a.(*anypb.Any); ok {
            if err:= any.UnmarshalTo(out); err!=nil {
                return syscall.KernelErr_UnmarshalFailed
            }
        } else {
            // `a` and `out` are the same type, so we can assign the values of a to out
            proto.Merge(out, a.(proto.Message))
        }
    }
    f.Method.CompleteMethod(ctx,out,TestErr(e)) 
    return syscall.KernelErr_NoError
}
func (f *FutureTestStart)Success(sfn func (proto.Message)) {
    x:=func(m *StartResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureTestStart)Failure(ffn func (int32)) {
    x:=func(err TestErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}
func NewFutureTestStart() *FutureTestStart {
    f:=&FutureTestStart{
        Method: future.NewMethod[*StartResponse,TestErr](nil,nil),
    } 
    return f
}
func (i *ClientTest_) TestStart(ctx context.Context, in *StartRequest) *FutureTestStart { 
    mid, ok := i.BaseService.MethodIdByName("Start")
    if !ok {
        f:=NewFutureTestStart()
        f.CompleteMethod(ctx,nil,1)/*dispatch error*/
    }
    cid,kerr:= i.BaseService.Dispatch(mid,in) 
    f:=NewFutureTestStart()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1)/*dispatch error*/
        return f
     }
    syscallguest.MatchCompleter(cid,f)
    return f
}  
//
// MethodCallSuite from test/v1/test.proto
//
//service interface
type MethodCallSuite interface {
    MethodCallSuiteExec(ctx context.Context,in *ExecRequest) *FutureMethodCallSuiteExec  
    MethodCallSuiteSuiteReport(ctx context.Context,in *SuiteReportRequest) *FutureMethodCallSuiteSuiteReport   
    Ready(context.Context,id.ServiceId) *future.Base[bool]
}

type ClientMethodCallSuite interface {
    MethodCallSuiteExec(ctx context.Context,in *ExecRequest) *FutureMethodCallSuiteExec  
    MethodCallSuiteSuiteReport(ctx context.Context,in *SuiteReportRequest) *FutureMethodCallSuiteSuiteReport   
}

// ClientMethodCallSuite difference from MethodCallSuite: Ready() 
type ClientMethodCallSuite_ struct {
    *client.BaseService
}
// Check that Client_ is a Client.
var _ = ClientMethodCallSuite(&ClientMethodCallSuite_{})

//
// method: MethodCallSuite.Exec 
//
type FutureMethodCallSuiteExec struct {
    Method *future.Method[*ExecResponse,TestErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureMethodCallSuiteExec) CompleteMethod(ctx context.Context,a proto.Message, e int32) syscall.KernelErr{
    out:=&ExecResponse{}
    if a!=nil {
        if any, ok := a.(*anypb.Any); ok {
            if err:= any.UnmarshalTo(out); err!=nil {
                return syscall.KernelErr_UnmarshalFailed
            }
        } else {
            // `a` and `out` are the same type, so we can assign the values of a to out
            proto.Merge(out, a.(proto.Message))
        }
    }
    f.Method.CompleteMethod(ctx,out,TestErr(e)) 
    return syscall.KernelErr_NoError
}
func (f *FutureMethodCallSuiteExec)Success(sfn func (proto.Message)) {
    x:=func(m *ExecResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureMethodCallSuiteExec)Failure(ffn func (int32)) {
    x:=func(err TestErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}
func NewFutureMethodCallSuiteExec() *FutureMethodCallSuiteExec {
    f:=&FutureMethodCallSuiteExec{
        Method: future.NewMethod[*ExecResponse,TestErr](nil,nil),
    } 
    return f
}
func (i *ClientMethodCallSuite_) MethodCallSuiteExec(ctx context.Context, in *ExecRequest) *FutureMethodCallSuiteExec { 
    mid, ok := i.BaseService.MethodIdByName("Exec")
    if !ok {
        f:=NewFutureMethodCallSuiteExec()
        f.CompleteMethod(ctx,nil,1)/*dispatch error*/
    }
    cid,kerr:= i.BaseService.Dispatch(mid,in) 
    f:=NewFutureMethodCallSuiteExec()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1)/*dispatch error*/
        return f
     }
    syscallguest.MatchCompleter(cid,f)
    return f
}

//
// method: MethodCallSuite.SuiteReport 
//
type FutureMethodCallSuiteSuiteReport struct {
    Base *future.Base[TestErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureMethodCallSuiteSuiteReport) CompleteMethod(ctx context.Context,a proto.Message, e int32) syscall.KernelErr{
    f.Base.Set(TestErr(e)) 
    return syscall.KernelErr_NoError
} 
func (f *FutureMethodCallSuiteSuiteReport)Success(sfn func (proto.Message)) {
    // no way for this to be called
} 

func (f *FutureMethodCallSuiteSuiteReport)Failure(ffn func (int32)) {
    x:=func(err TestErr) {
        ffn(int32(err))
    }
    f.Base.Handle(x) 
}
func NewFutureMethodCallSuiteSuiteReport() *FutureMethodCallSuiteSuiteReport {
    f:=&FutureMethodCallSuiteSuiteReport{
        Base: future.NewBase[TestErr](),
    } 
    return f
}
func (i *ClientMethodCallSuite_) MethodCallSuiteSuiteReport(ctx context.Context, in *SuiteReportRequest) *FutureMethodCallSuiteSuiteReport { 
    mid, ok := i.BaseService.MethodIdByName("SuiteReport")
    if !ok {
        f:=NewFutureMethodCallSuiteSuiteReport()
        f.CompleteMethod(ctx,nil,1)/*dispatch error*/
    }
    cid,kerr:= i.BaseService.Dispatch(mid,in) 
    f:=NewFutureMethodCallSuiteSuiteReport()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1)/*dispatch error*/
        return f
     }
    syscallguest.MatchCompleter(cid,f)
    return f
}  
//
// UnderTest from test/v1/test.proto
//
//service interface
type UnderTest interface {
    UnderTestExec(ctx context.Context,in *ExecRequest) *FutureUnderTestExec   
    Ready(context.Context,id.ServiceId) *future.Base[bool]
}

type ClientUnderTest interface {
    UnderTestExec(ctx context.Context,in *ExecRequest) *FutureUnderTestExec   
}

// ClientUnderTest difference from UnderTest: Ready() 
type ClientUnderTest_ struct {
    *client.BaseService
}
// Check that Client_ is a Client.
var _ = ClientUnderTest(&ClientUnderTest_{})

//
// method: UnderTest.Exec 
//
type FutureUnderTestExec struct {
    Method *future.Method[*ExecResponse,TestErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureUnderTestExec) CompleteMethod(ctx context.Context,a proto.Message, e int32) syscall.KernelErr{
    out:=&ExecResponse{}
    if a!=nil {
        if any, ok := a.(*anypb.Any); ok {
            if err:= any.UnmarshalTo(out); err!=nil {
                return syscall.KernelErr_UnmarshalFailed
            }
        } else {
            // `a` and `out` are the same type, so we can assign the values of a to out
            proto.Merge(out, a.(proto.Message))
        }
    }
    f.Method.CompleteMethod(ctx,out,TestErr(e)) 
    return syscall.KernelErr_NoError
}
func (f *FutureUnderTestExec)Success(sfn func (proto.Message)) {
    x:=func(m *ExecResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureUnderTestExec)Failure(ffn func (int32)) {
    x:=func(err TestErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}
func NewFutureUnderTestExec() *FutureUnderTestExec {
    f:=&FutureUnderTestExec{
        Method: future.NewMethod[*ExecResponse,TestErr](nil,nil),
    } 
    return f
}
func (i *ClientUnderTest_) UnderTestExec(ctx context.Context, in *ExecRequest) *FutureUnderTestExec { 
    mid, ok := i.BaseService.MethodIdByName("Exec")
    if !ok {
        f:=NewFutureUnderTestExec()
        f.CompleteMethod(ctx,nil,1)/*dispatch error*/
    }
    cid,kerr:= i.BaseService.Dispatch(mid,in) 
    f:=NewFutureUnderTestExec()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1)/*dispatch error*/
        return f
     }
    syscallguest.MatchCompleter(cid,f)
    return f
}  
