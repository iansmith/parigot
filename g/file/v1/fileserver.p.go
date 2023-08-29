//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: file/v1/file.proto

package file




import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"runtime/debug"
    "unsafe"

 
    // this set of imports is _unrelated_ to the particulars of what the .proto imported... those are above
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"  
	"github.com/iansmith/parigot/api/shared/id"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/lib/go/future"
	apishared "github.com/iansmith/parigot/api/shared"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"github.com/iansmith/parigot/lib/go/client"
)
var _ =  unsafe.Sizeof([]byte{})
 

func Launch(ctx context.Context, sid id.ServiceId, impl File) *future.Base[bool] {

	readyResult:=future.NewBase[bool]()

	ready:=impl.Ready(ctx,sid)
	ready.Handle(func (b bool) {
		if b {
			readyResult.Set(true)			
			return
		}
		slog.Error("Unable to start file.v1.File, Ready returned false")
		readyResult.Set(false)
	})

	return readyResult
}

// Note that  Init returns a future, but the case of failure is covered
// by this definition so the caller need only deal with Success case.
// The context passed here does not need to contain a logger, one will be created.
func Init(require []lib.MustRequireFunc, impl File) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture, context.Context, id.ServiceId){ 
	// tricky, this context really should not be used but is
	// passed so as to allow printing if things go wrong
	ctx, myId := MustRegister()
	MustExport(context.Background(),myId)
	if len(require)>0 {
		for _, f := range require {
			f(ctx, myId)
		}
	}
	smmap, launchF:=MustLaunchService(ctx, myId, impl)
	launchF.Failure(func (err syscall.KernelErr) {
		t:=syscall.KernelErr_name[int32(err)]
		slog.Error("launch failure on call File","error",t)
		lib.ExitSelf(ctx, 1, myId)
	})
	return smmap,launchF, ctx,myId
}
func Run(ctx context.Context,
	binding *lib.ServiceMethodMap, timeoutInMillis int32, bg lib.Backgrounder) syscall.KernelErr{
	defer func() {
		if r := recover(); r != nil {
			s, ok:=r.(string)
			if !ok && s!=apishared.ControlledExit {
				slog.Error("Run File: trapped a panic in the guest side", "recovered", r)
				debug.PrintStack()
			}
		}
	}()
	var kerr syscall.KernelErr
	for {
		kerr:=ReadOneAndCall(ctx, binding, timeoutInMillis)
		if kerr == syscall.KernelErr_ReadOneTimeout {
			if bg==nil {
				continue
			}
			slog.Info("calling backgrounder of File")
			bg.Background(ctx)
			continue
		}
		if kerr == syscall.KernelErr_NoError {
			continue
		}
		break
	}
	slog.Error("error while waiting for  service calls", "error",syscall.KernelErr_name[int32(kerr)])
	return kerr
}
// Increase this value at your peril!
// Decreasing this value may make your overall program more responsive if you have many services.
var TimeoutInMillis = int32(50)

func ReadOneAndCall(ctx context.Context, binding *lib.ServiceMethodMap, 
	timeoutInMillis int32) syscall.KernelErr{
	req:=syscall.ReadOneRequest{}
	hid:= syscallguest.CurrentHostId()

	req.TimeoutInMillis = timeoutInMillis
	req.HostId = hid.Marshal()
	resp, err:=syscallguest.ReadOne(ctx, &req)
	if err!=syscall.KernelErr_NoError {
		return err
	}
	// is timeout?
	if resp.Timeout {
		return syscall.KernelErr_ReadOneTimeout
	}

	// check for finished futures from within our address space
	ctx, t:=lib.CurrentTime(ctx)
	syscallguest.ExpireMethod(ctx,t)

	// is a promise being completed that was fulfilled somewhere else
	if r:=resp.GetResolved(); r!=nil {
		cid:=id.UnmarshalCallId(r.GetCallId())
		defer func() {
			if r:=recover(); r!=nil {
				sid:=id.UnmarshalServiceId(resp.GetBundle().GetServiceId())
				mid:=id.UnmarshalMethodId(resp.GetBundle().GetMethodId())
				log.Printf("completing method %s on service %s failed due to panic: '%s', exiting",
					mid.Short(), sid.Short(), r)
				debug.PrintStack()
				syscallguest.Exit(ctx, &syscall.ExitRequest{
					Pair: &syscall.ExitPair {
						ServiceId: sid.Marshal(),
						Code: 2,
					},
				})
			}
		}()
		syscallguest.CompleteCall(ctx, syscallguest.CurrentHostId(),cid,r.GetResult(), r.GetResultError())
		return syscall.KernelErr_NoError
	}

	// its a method call from another address space
	sid:=id.UnmarshalServiceId(resp.GetBundle().GetServiceId())
	mid:=id.UnmarshalMethodId(resp.GetBundle().GetMethodId())
	cid:=id.UnmarshalCallId(resp.GetBundle().GetCallId())

	//if mid.Equal(apishared.ExitMethod) {
	// log.Printf("xxx -- got an exit marked read one %s", hid.Short())
	//	os.Exit(51)
	//}

	// we let the invoker handle the unmarshal from anypb.Any because it
	// knows the precise type to be consumed
	fn:=binding.Func(sid,mid)
	if fn==nil {
		slog.Error("File, readOneAndCall:unable to find binding for method on service, ignoring","mid",mid.Short(),"sid", sid.Short(),
			"current host",syscallguest.CurrentHostId())
		return syscall.KernelErr_NoError
	}
	fut:=fn.Invoke(ctx,resp.GetParamOrResult())
	fut.Success(func (result proto.Message){
		rvReq:=&syscall.ReturnValueRequest{}
		rvReq.Bundle=&syscall.MethodBundle{}

		rvReq.Bundle.CallId= cid.Marshal()
		rvReq.Bundle.HostId= syscallguest.CurrentHostId().Marshal()
		var a anypb.Any
		if err:=a.MarshalFrom(result); err!=nil {
			slog.Error("unable to marshal result for return value request")
			return
		}
		rvReq.Result = &a
		rvReq.ResultError = 0
		syscallguest.ReturnValue(ctx, rvReq) // nowhere for return value to go
	})
	fut.Failure(func (err int32) {
		rvReq:=&syscall.ReturnValueRequest{}
		rvReq.Bundle=&syscall.MethodBundle{}

		rvReq.Bundle.CallId= cid.Marshal()
		rvReq.Bundle.HostId= syscallguest.CurrentHostId().Marshal()
		rvReq.ResultError = err
		syscallguest.ReturnValue(ctx,rvReq) // nowhere for return value to go
	})
	return syscall.KernelErr_NoError

}


func bind(ctx context.Context,sid id.ServiceId, impl File) (*lib.ServiceMethodMap, syscall.KernelErr) {
	smmap:=lib.NewServiceMethodMap()
	var mid id.MethodId
	var bindReq *syscall.BindMethodRequest
	var resp *syscall.BindMethodResponse
	var err syscall.KernelErr
//
// file.v1.File.Open
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Open"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"File","Open",
		GenerateOpenInvoker(impl))
//
// file.v1.File.Create
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Create"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"File","Create",
		GenerateCreateInvoker(impl))
//
// file.v1.File.Close
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Close"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"File","Close",
		GenerateCloseInvoker(impl))
//
// file.v1.File.LoadTestData
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "LoadTestData"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"File","LoadTestData",
		GenerateLoadTestDataInvoker(impl))
//
// file.v1.File.Read
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Read"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"File","Read",
		GenerateReadInvoker(impl))
//
// file.v1.File.Write
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Write"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"File","Write",
		GenerateWriteInvoker(impl))
//
// file.v1.File.Delete
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Delete"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"File","Delete",
		GenerateDeleteInvoker(impl))
//
// file.v1.File.Stat
//

	bindReq = &syscall.BindMethodRequest{}
	bindReq.HostId = syscallguest.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Stat"
	resp, err=syscallguest.BindMethod(ctx, bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"File","Stat",
		GenerateStatInvoker(impl)) 
	return smmap,syscall.KernelErr_NoError
}
 

// Locate finds a reference to the client interface of file.  
func Locate(ctx context.Context,sid id.ServiceId) (Client,syscall.KernelErr) {
    cs, kerr:=client.LocateDynamic(ctx, "file.v1","file", sid)
    if kerr!=syscall.KernelErr_NoError{
        return nil, kerr
    }
    return &Client_{
        BaseService: cs,
    },syscall.KernelErr_NoError
}

func MustLocate(ctx context.Context, sid id.ServiceId) Client {
    result, err:=Locate(ctx, sid)
    name:=syscall.KernelErr_name[int32(err)]
    normal:="unable to locate file.v1.file:"+name
    if err!=0 {
        if err == syscall.KernelErr_NotRequired {
            slog.Error("service was located, but it was not required")
            panic("locate attempted on a service that was not required")
        }
        panic(normal)
    }
    return result
}


func Register() (id.ServiceId, syscall.KernelErr){
    req := &syscall.RegisterRequest{}
	debugName:=fmt.Sprintf("%s.%s","file.v1","file")
	req.HostId = syscallguest.CurrentHostId().Marshal()
	req.DebugName = debugName

	resp, err := syscallguest.Register(context.Background(), req)
    if err!=syscall.KernelErr_NoError{
        return id.ServiceIdZeroValue(), err
    }
    sid:=id.UnmarshalServiceId(resp.ServiceId)
    if sid.IsZeroOrEmptyValue() {
        panic("received bad service Id from register")
    }

    return sid,syscall.KernelErr_NoError
}
func MustRegister() (context.Context,id.ServiceId) {
    sid, err:=Register()
    if err!=syscall.KernelErr_NoError {
        slog.Error("unable to register","package","file.v1","service name","file")
        panic("unable to register "+"file")
    }
    return context.Background(), sid
}

func MustRequire(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Require1(ctx, "file.v1","file",sid)
    if err!=syscall.KernelErr_NoError {
        if err==syscall.KernelErr_DependencyCycle{
            slog.Error("unable to require because it creates a dependcy loop","package","file.v1","service name","file","error",syscall.KernelErr_name[int32(err)])
            panic("require file.v1.file creates a dependency loop")
        }
        slog.Error("unable to require","package","file.v1","service name","file","error",syscall.KernelErr_name[int32(err)])
        panic("not able to require file.v1.file:"+syscall.KernelErr_name[int32(err)])
    }
}

func MustExport(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Export1(ctx,"file.v1","file",sid)
    if err!=syscall.KernelErr_NoError{
        slog.Error("unable to export","package","file.v1","service name","file")
        panic("not able to export file.v1.file:"+syscall.KernelErr_name[int32(err)])
    }
}


func LaunchService(ctx context.Context, sid id.ServiceId, impl  File) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture,syscall.KernelErr) {
	smmap, err:=bind(ctx,sid, impl)
	if err!=0{
		return  nil,nil,syscall.KernelErr(err)
	}

	cid:=id.NewCallId()
	req:=&syscall.LaunchRequest{
		ServiceId: sid.Marshal(),
		CallId: cid.Marshal(),
		HostId: syscallguest.CurrentHostId().Marshal(),
		MethodId: apishared.LaunchMethod.Marshal(),
	}
	fut:=syscallguest.Launch(ctx,req)


    return smmap,fut,syscall.KernelErr_NoError

}
func MustLaunchService(ctx context.Context, sid id.ServiceId, impl File) (*lib.ServiceMethodMap, *syscallguest.LaunchFuture) {
 
    smmap,fut,err:=LaunchService(ctx,sid,impl)
    if err!=syscall.KernelErr_NoError {
        panic("Unable to call LaunchService successfully: "+syscall.KernelErr_name[int32(err)])
    }
    return smmap,fut
}


// If you want to implement part of your server in host cost you should call 
// <methodName>Host from your server implementation. These will be optimized 
// away by the compiler if you don't use them--in other words, if you want to 
// implement everything on the guest side).
//  

//go:wasmimport file open_
func Open_(int32,int32,int32,int32) int64
func OpenHost(ctx context.Context,inPtr *OpenRequest) *FutureOpen {
	outProtoPtr := (*OpenResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Open_)
	f:=NewFutureOpen()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport file create_
func Create_(int32,int32,int32,int32) int64
func CreateHost(ctx context.Context,inPtr *CreateRequest) *FutureCreate {
	outProtoPtr := (*CreateResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Create_)
	f:=NewFutureCreate()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport file close_
func Close_(int32,int32,int32,int32) int64
func CloseHost(ctx context.Context,inPtr *CloseRequest) *FutureClose {
	outProtoPtr := (*CloseResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Close_)
	f:=NewFutureClose()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport file load_test_data_
func LoadTestData_(int32,int32,int32,int32) int64
func LoadTestDataHost(ctx context.Context,inPtr *LoadTestDataRequest) *FutureLoadTestData {
	outProtoPtr := (*LoadTestDataResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, LoadTestData_)
	f:=NewFutureLoadTestData()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport file read_
func Read_(int32,int32,int32,int32) int64
func ReadHost(ctx context.Context,inPtr *ReadRequest) *FutureRead {
	outProtoPtr := (*ReadResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Read_)
	f:=NewFutureRead()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport file write_
func Write_(int32,int32,int32,int32) int64
func WriteHost(ctx context.Context,inPtr *WriteRequest) *FutureWrite {
	outProtoPtr := (*WriteResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Write_)
	f:=NewFutureWrite()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport file delete_
func Delete_(int32,int32,int32,int32) int64
func DeleteHost(ctx context.Context,inPtr *DeleteRequest) *FutureDelete {
	outProtoPtr := (*DeleteResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Delete_)
	f:=NewFutureDelete()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport file stat_
func Stat_(int32,int32,int32,int32) int64
func StatHost(ctx context.Context,inPtr *StatRequest) *FutureStat {
	outProtoPtr := (*StatResponse)(nil)
	ret, raw, _:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Stat_)
	f:=NewFutureStat()
	f.CompleteMethod(ctx,ret,raw)
	return f
}   

// This is interface for invocation.

type invokeOpen struct {
    fn func(context.Context,*OpenRequest) *FutureOpen
}

func (t *invokeOpen) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&OpenRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateOpenInvoker(impl File) future.Invoker {
	return &invokeOpen{fn:impl.Open} 
}

// This is interface for invocation.

type invokeCreate struct {
    fn func(context.Context,*CreateRequest) *FutureCreate
}

func (t *invokeCreate) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&CreateRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateCreateInvoker(impl File) future.Invoker {
	return &invokeCreate{fn:impl.Create} 
}

// This is interface for invocation.

type invokeClose struct {
    fn func(context.Context,*CloseRequest) *FutureClose
}

func (t *invokeClose) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&CloseRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateCloseInvoker(impl File) future.Invoker {
	return &invokeClose{fn:impl.Close} 
}

// This is interface for invocation.

type invokeLoadTestData struct {
    fn func(context.Context,*LoadTestDataRequest) *FutureLoadTestData
}

func (t *invokeLoadTestData) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&LoadTestDataRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateLoadTestDataInvoker(impl File) future.Invoker {
	return &invokeLoadTestData{fn:impl.LoadTestData} 
}

// This is interface for invocation.

type invokeRead struct {
    fn func(context.Context,*ReadRequest) *FutureRead
}

func (t *invokeRead) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&ReadRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateReadInvoker(impl File) future.Invoker {
	return &invokeRead{fn:impl.Read} 
}

// This is interface for invocation.

type invokeWrite struct {
    fn func(context.Context,*WriteRequest) *FutureWrite
}

func (t *invokeWrite) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&WriteRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateWriteInvoker(impl File) future.Invoker {
	return &invokeWrite{fn:impl.Write} 
}

// This is interface for invocation.

type invokeDelete struct {
    fn func(context.Context,*DeleteRequest) *FutureDelete
}

func (t *invokeDelete) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&DeleteRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateDeleteInvoker(impl File) future.Invoker {
	return &invokeDelete{fn:impl.Delete} 
}

// This is interface for invocation.

type invokeStat struct {
    fn func(context.Context,*StatRequest) *FutureStat
}

func (t *invokeStat) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
    in:=&StatRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        slog.Error("unmarshal inside Invoke() failed","error",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateStatInvoker(impl File) future.Invoker {
	return &invokeStat{fn:impl.Stat} 
}  
