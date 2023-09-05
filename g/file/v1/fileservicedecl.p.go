//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: file/v1/file.proto

package file


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
// File from file/v1/file.proto
//
//service interface
type File interface {
    Open(ctx context.Context,in *OpenRequest) *FutureOpen  
    Create(ctx context.Context,in *CreateRequest) *FutureCreate  
    Close(ctx context.Context,in *CloseRequest) *FutureClose  
    LoadTestData(ctx context.Context,in *LoadTestDataRequest) *FutureLoadTestData  
    Read(ctx context.Context,in *ReadRequest) *FutureRead  
    Write(ctx context.Context,in *WriteRequest) *FutureWrite  
    Delete(ctx context.Context,in *DeleteRequest) *FutureDelete  
    Stat(ctx context.Context,in *StatRequest) *FutureStat   
    Ready(context.Context,id.ServiceId) *future.Base[bool]
}

type Client interface {
    Open(ctx context.Context,in *OpenRequest) *FutureOpen  
    Create(ctx context.Context,in *CreateRequest) *FutureCreate  
    Close(ctx context.Context,in *CloseRequest) *FutureClose  
    LoadTestData(ctx context.Context,in *LoadTestDataRequest) *FutureLoadTestData  
    Read(ctx context.Context,in *ReadRequest) *FutureRead  
    Write(ctx context.Context,in *WriteRequest) *FutureWrite  
    Delete(ctx context.Context,in *DeleteRequest) *FutureDelete  
    Stat(ctx context.Context,in *StatRequest) *FutureStat   
}

// Client difference from File: Ready() 
type Client_ struct {
    *client.BaseService
}
// Check that Client_ is a Client.
var _ = Client(&Client_{})

//
// method: File.Open 
//
type FutureOpen struct {
    Method *future.Method[*OpenResponse,FileErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureOpen) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&OpenResponse{}
    if a!=nil {
        if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
            return syscall.KernelErr_UnmarshalFailed
        }
    }
    f.Method.CompleteMethod(ctx,out,FileErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureOpen)Success(sfn func (proto.Message)) {
    x:=func(m *OpenResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureOpen)Failure(ffn func (int32)) {
    x:=func(err FileErr) {
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
        Method: future.NewMethod[*OpenResponse,FileErr](nil,nil),
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
// method: File.Create 
//
type FutureCreate struct {
    Method *future.Method[*CreateResponse,FileErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureCreate) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&CreateResponse{}
    if a!=nil {
        if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
            return syscall.KernelErr_UnmarshalFailed
        }
    }
    f.Method.CompleteMethod(ctx,out,FileErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureCreate)Success(sfn func (proto.Message)) {
    x:=func(m *CreateResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureCreate)Failure(ffn func (int32)) {
    x:=func(err FileErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureCreate)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureCreate)Cancel()   {
    f.Method.Cancel()
}
func NewFutureCreate() *FutureCreate {
    f:=&FutureCreate{
        Method: future.NewMethod[*CreateResponse,FileErr](nil,nil),
    } 
    return f
}
func (i *Client_) Create(ctx context.Context, in *CreateRequest) *FutureCreate { 
    mid, ok := i.BaseService.MethodIdByName("Create")
    if !ok {
        f:=NewFutureCreate()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureCreate()
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
// method: File.Close 
//
type FutureClose struct {
    Method *future.Method[*CloseResponse,FileErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureClose) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&CloseResponse{}
    if a!=nil {
        if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
            return syscall.KernelErr_UnmarshalFailed
        }
    }
    f.Method.CompleteMethod(ctx,out,FileErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureClose)Success(sfn func (proto.Message)) {
    x:=func(m *CloseResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureClose)Failure(ffn func (int32)) {
    x:=func(err FileErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureClose)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureClose)Cancel()   {
    f.Method.Cancel()
}
func NewFutureClose() *FutureClose {
    f:=&FutureClose{
        Method: future.NewMethod[*CloseResponse,FileErr](nil,nil),
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
// method: File.LoadTestData 
//
type FutureLoadTestData struct {
    Method *future.Method[*LoadTestDataResponse,FileErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureLoadTestData) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&LoadTestDataResponse{}
    if a!=nil {
        if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
            return syscall.KernelErr_UnmarshalFailed
        }
    }
    f.Method.CompleteMethod(ctx,out,FileErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureLoadTestData)Success(sfn func (proto.Message)) {
    x:=func(m *LoadTestDataResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureLoadTestData)Failure(ffn func (int32)) {
    x:=func(err FileErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureLoadTestData)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureLoadTestData)Cancel()   {
    f.Method.Cancel()
}
func NewFutureLoadTestData() *FutureLoadTestData {
    f:=&FutureLoadTestData{
        Method: future.NewMethod[*LoadTestDataResponse,FileErr](nil,nil),
    } 
    return f
}
func (i *Client_) LoadTestData(ctx context.Context, in *LoadTestDataRequest) *FutureLoadTestData { 
    mid, ok := i.BaseService.MethodIdByName("LoadTestData")
    if !ok {
        f:=NewFutureLoadTestData()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureLoadTestData()
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
// method: File.Read 
//
type FutureRead struct {
    Method *future.Method[*ReadResponse,FileErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureRead) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&ReadResponse{}
    if a!=nil {
        if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
            return syscall.KernelErr_UnmarshalFailed
        }
    }
    f.Method.CompleteMethod(ctx,out,FileErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureRead)Success(sfn func (proto.Message)) {
    x:=func(m *ReadResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureRead)Failure(ffn func (int32)) {
    x:=func(err FileErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureRead)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureRead)Cancel()   {
    f.Method.Cancel()
}
func NewFutureRead() *FutureRead {
    f:=&FutureRead{
        Method: future.NewMethod[*ReadResponse,FileErr](nil,nil),
    } 
    return f
}
func (i *Client_) Read(ctx context.Context, in *ReadRequest) *FutureRead { 
    mid, ok := i.BaseService.MethodIdByName("Read")
    if !ok {
        f:=NewFutureRead()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureRead()
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
// method: File.Write 
//
type FutureWrite struct {
    Method *future.Method[*WriteResponse,FileErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureWrite) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&WriteResponse{}
    if a!=nil {
        if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
            return syscall.KernelErr_UnmarshalFailed
        }
    }
    f.Method.CompleteMethod(ctx,out,FileErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureWrite)Success(sfn func (proto.Message)) {
    x:=func(m *WriteResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureWrite)Failure(ffn func (int32)) {
    x:=func(err FileErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureWrite)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureWrite)Cancel()   {
    f.Method.Cancel()
}
func NewFutureWrite() *FutureWrite {
    f:=&FutureWrite{
        Method: future.NewMethod[*WriteResponse,FileErr](nil,nil),
    } 
    return f
}
func (i *Client_) Write(ctx context.Context, in *WriteRequest) *FutureWrite { 
    mid, ok := i.BaseService.MethodIdByName("Write")
    if !ok {
        f:=NewFutureWrite()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureWrite()
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
// method: File.Delete 
//
type FutureDelete struct {
    Method *future.Method[*DeleteResponse,FileErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureDelete) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&DeleteResponse{}
    if a!=nil {
        if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
            return syscall.KernelErr_UnmarshalFailed
        }
    }
    f.Method.CompleteMethod(ctx,out,FileErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureDelete)Success(sfn func (proto.Message)) {
    x:=func(m *DeleteResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureDelete)Failure(ffn func (int32)) {
    x:=func(err FileErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureDelete)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureDelete)Cancel()   {
    f.Method.Cancel()
}
func NewFutureDelete() *FutureDelete {
    f:=&FutureDelete{
        Method: future.NewMethod[*DeleteResponse,FileErr](nil,nil),
    } 
    return f
}
func (i *Client_) Delete(ctx context.Context, in *DeleteRequest) *FutureDelete { 
    mid, ok := i.BaseService.MethodIdByName("Delete")
    if !ok {
        f:=NewFutureDelete()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureDelete()
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
// method: File.Stat 
//
type FutureStat struct {
    Method *future.Method[*StatResponse,FileErr]
} 

// This is the same API for output needed or not because of the Completer interface.
// Note that the return value refers to the process of the setup/teardown, not the
// execution of the user level code.
func (f * FutureStat) CompleteMethod(ctx context.Context,a proto.Message, e int32, orig id.HostId) syscall.KernelErr{
    out:=&StatResponse{}
    if a!=nil {
        if err:= a.(*anypb.Any).UnmarshalTo(out); err!=nil {
            return syscall.KernelErr_UnmarshalFailed
        }
    }
    f.Method.CompleteMethod(ctx,out,FileErr(e)) 
    return syscall.KernelErr_NoError

}
func (f *FutureStat)Success(sfn func (proto.Message)) {
    x:=func(m *StatResponse){
        sfn(m)
    }
    f.Method.Success(x)
} 

func (f *FutureStat)Failure(ffn func (int32)) {
    x:=func(err FileErr) {
        ffn(int32(err))
    }
    f.Method.Failure(x) 
}

func (f *FutureStat)Completed() bool  {
    return f.Method.Completed()

}
func (f *FutureStat)Cancel()   {
    f.Method.Cancel()
}
func NewFutureStat() *FutureStat {
    f:=&FutureStat{
        Method: future.NewMethod[*StatResponse,FileErr](nil,nil),
    } 
    return f
}
func (i *Client_) Stat(ctx context.Context, in *StatRequest) *FutureStat { 
    mid, ok := i.BaseService.MethodIdByName("Stat")
    if !ok {
        f:=NewFutureStat()
        f.CompleteMethod(ctx,nil,1,syscallguest.CurrentHostId())/*dispatch error*/
    }
    targetHid,cid,kerr:= i.BaseService.Dispatch(ctx,mid,in) 
    f:=NewFutureStat()
    if kerr!=syscall.KernelErr_NoError{
        f.CompleteMethod(ctx,nil, 1,syscallguest.CurrentHostId())/*dispatch error*/
        return f
     }

    ctx, t:=lib.CurrentTime(ctx)
    source:=syscallguest.CurrentHostId()
    syscallguest.MatchCompleter(ctx,t,source,targetHid,cid,f)
    return f
}  
