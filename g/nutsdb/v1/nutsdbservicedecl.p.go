//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: nutsdb/v1/nutsdb.proto

package nutsdb


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
// NutsDB from nutsdb/v1/nutsdb.proto
//
//service interface
type NutsDB interface {
    Open(ctx context.Context,in *OpenRequest) *FutureOpen  
    Close(ctx context.Context,in *CloseRequest) *FutureClose  
    ReadPair(ctx context.Context,in *ReadPairRequest) *FutureReadPair  
    WritePair(ctx context.Context,in *WritePairRequest) *FutureWritePair   
    Ready(context.Context,id.ServiceId) *future.Base[bool]
}

type Client interface {
    Open(ctx context.Context,in *OpenRequest) *FutureOpen  
    Close(ctx context.Context,in *CloseRequest) *FutureClose  
    ReadPair(ctx context.Context,in *ReadPairRequest) *FutureReadPair  
    WritePair(ctx context.Context,in *WritePairRequest) *FutureWritePair   
}

// Client difference from NutsDB: Ready() 
type Client_ struct {
    *client.BaseService
}
// Check that Client_ is a Client.
var _ = Client(&Client_{})

//
// method: NutsDB.Open 
//
type FutureOpen struct {
    Method *future.Method[*OpenResponse,NutsDBErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureOpen) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&OpenResponse{}
    if a!=nil {
        tmp, ok:=a.(*OpenResponse)
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
    f.Method.CompleteMethod(ctx,out,NutsDBErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureOpen)Success(sfn func (proto.Message)) {
    x:=func(m *OpenResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureOpen)VerifyRejectPresent() {
    f.Method.VerifyRejectPresent()
 
}

func (f *FutureOpen)Failure(ffn func (int32)) {
    x:=func(err NutsDBErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureOpen)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureOpen)Cancel()   {
    f.Method.Cancel()
}
func NewFutureOpen() *FutureOpen {
    f:=&FutureOpen{
        Method: future.NewMethod[*OpenResponse,NutsDBErr](nil,nil),
    } 
    return f
}
func (i *Client_) Open(ctx context.Context, in *OpenRequest) *FutureOpen { 
    mid, ok := i.BaseService.MethodIdByName("Open")
    if !ok {
        f:=NewFutureOpen()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureOpen()
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
// method: NutsDB.Close 
//
type FutureClose struct {
    Base *future.Base[NutsDBErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureClose) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    f.Base.Set(NutsDBErr(e)) 
    return syscall.KernelErr_NoError

} 
func (f *FutureClose)Success(sfn func (proto.Message)) {
    // no way for this to be called
} 

func (f *FutureClose)VerifyRejectPresent() { 
}

func (f *FutureClose)Failure(ffn func (int32)) {
    x:=func(err NutsDBErr) {
        ffn(int32(err))
    }
    f.Base.Handle(x) 
}

func (f *FutureClose)Completed() bool  {
    return f.Base.Completed()

}
func (f *FutureClose)Cancel()   {
    f.Base.Cancel()
}
func NewFutureClose() *FutureClose {
    f:=&FutureClose{
        Base: future.NewBase[NutsDBErr](),
    } 
    return f
}
func (i *Client_) Close(ctx context.Context, in *CloseRequest) *FutureClose { 
    mid, ok := i.BaseService.MethodIdByName("Close")
    if !ok {
        f:=NewFutureClose()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureClose()
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
// method: NutsDB.ReadPair 
//
type FutureReadPair struct {
    Method *future.Method[*ReadPairResponse,NutsDBErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureReadPair) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&ReadPairResponse{}
    if a!=nil {
        tmp, ok:=a.(*ReadPairResponse)
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
    f.Method.CompleteMethod(ctx,out,NutsDBErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureReadPair)Success(sfn func (proto.Message)) {
    x:=func(m *ReadPairResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureReadPair)VerifyRejectPresent() {
    f.Method.VerifyRejectPresent()
 
}

func (f *FutureReadPair)Failure(ffn func (int32)) {
    x:=func(err NutsDBErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureReadPair)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureReadPair)Cancel()   {
    f.Method.Cancel()
}
func NewFutureReadPair() *FutureReadPair {
    f:=&FutureReadPair{
        Method: future.NewMethod[*ReadPairResponse,NutsDBErr](nil,nil),
    } 
    return f
}
func (i *Client_) ReadPair(ctx context.Context, in *ReadPairRequest) *FutureReadPair { 
    mid, ok := i.BaseService.MethodIdByName("ReadPair")
    if !ok {
        f:=NewFutureReadPair()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureReadPair()
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
// method: NutsDB.WritePair 
//
type FutureWritePair struct {
    Method *future.Method[*WritePairResponse,NutsDBErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureWritePair) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&WritePairResponse{}
    if a!=nil {
        tmp, ok:=a.(*WritePairResponse)
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
    f.Method.CompleteMethod(ctx,out,NutsDBErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureWritePair)Success(sfn func (proto.Message)) {
    x:=func(m *WritePairResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureWritePair)VerifyRejectPresent() {
    f.Method.VerifyRejectPresent()
 
}

func (f *FutureWritePair)Failure(ffn func (int32)) {
    x:=func(err NutsDBErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureWritePair)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureWritePair)Cancel()   {
    f.Method.Cancel()
}
func NewFutureWritePair() *FutureWritePair {
    f:=&FutureWritePair{
        Method: future.NewMethod[*WritePairResponse,NutsDBErr](nil,nil),
    } 
    return f
}
func (i *Client_) WritePair(ctx context.Context, in *WritePairRequest) *FutureWritePair { 
    mid, ok := i.BaseService.MethodIdByName("WritePair")
    if !ok {
        f:=NewFutureWritePair()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureWritePair()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1,syscallguest.CurrentHostId())/*dispatch error*/
        return f
     }

    ctx, t:=lib.CurrentTime(ctx)
    source:=syscallguest.CurrentHostId()
    syscallguest.MatchCompleter(ctx,t,source,targetHid,cid,f)
    return f
}  
