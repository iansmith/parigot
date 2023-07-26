//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: file/v1/file.proto

package file




import (
	"context"
    "unsafe" 
    // this set of imports is _unrelated_ to the particulars of what the .proto imported... those are above
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"  
	pcontext "github.com/iansmith/parigot/context"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/api/shared/id"
	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/lib/go/future"
	"github.com/iansmith/parigot/lib/go/client"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/proto"

)
var _ =  unsafe.Sizeof([]byte{})
 
func Launch(ctx context.Context, sid id.ServiceId, impl File) *future.Base[bool] {

	defer func() {
		pcontext.Dump(ctx)
	}()

	readyResult:=future.NewBase[bool]()

	ready:=impl.Ready(ctx,sid)
	ready.Handle(func (b bool) {
		if b {
			readyResult.Set(true)			
			return
		}
		pcontext.Errorf(ctx,"Unable to start file.v1.File, Ready returned false")
		pcontext.Dump(ctx)
		readyResult.Set(false)
	})

	return readyResult
}

// Note that  Init returns a future, but the case of failure is covered
// by this definition so the caller need only deal with Success case.
func Init(ctx context.Context,require []lib.MustRequireFunc, impl File) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture, id.ServiceId){
	defer func() {
		if r := recover(); r != nil {
			pcontext.Infof(ctx, "InitFile: trapped a panic in the guest side: %v", r)
		}
		pcontext.Dump(ctx)
	}()

	myId := MustRegister(ctx)
	MustExport(ctx,myId)
	if len(require)>0 {
		for _, f := range require {
			f(ctx, myId)
		}
	}
	smmap, launchF:=MustLaunchService(ctx, myId, impl)
	launchF.Failure(func (err syscall.KernelErr) {
		t:=syscall.KernelErr_name[int32(err)]
		pcontext.Errorf(ctx, "launch failure on call File:%s",t)
		lib.ExitClient(ctx, 1, myId, "unable to Launch in Init:"+t,
			"unable to call Exit in Init:"+t)
	})
	return smmap,launchF, myId
}
func Run(ctx context.Context,
	binding *lib.ServiceMethodMap, timeoutInMillis int32, bg lib.Backgrounder) syscall.KernelErr{
	defer func() {
		if r := recover(); r != nil {
			pcontext.Infof(ctx, "Run: trapped a panic in the guest side: %v", r)
		}
		pcontext.Dump(ctx)
	}()
	var kerr syscall.KernelErr
	for {
		pctx:=pcontext.CallTo(ctx,"ReadOneAndCall")
		kerr:=ReadOneAndCall(pctx, binding, timeoutInMillis)
		pcontext.Dump(pctx)
		if kerr == syscall.KernelErr_ReadOneTimeout {
			if bg==nil {
				continue
			}
			pcontext.Infof(ctx,"calling backgrounder of File")
			bgctx:=pcontext.CallTo(ctx,"Background")
			bg.Background(bgctx)
			pcontext.Dump(bgctx)
			continue
		}
		if kerr == syscall.KernelErr_NoError {
			continue
		}
		break
	}
	pcontext.Errorf(ctx, "error while waiting for  service calls: %s", syscall.KernelErr_name[int32(kerr)])
	return kerr
}

var TimeoutInMillis = int32(10000)

func ReadOneAndCall(ctx context.Context, binding *lib.ServiceMethodMap, 
	timeoutInMillis int32) syscall.KernelErr{
	req:=syscall.ReadOneRequest{}

	// makes a copy
	for _, c := range binding.Call() {
		req.Call=append(req.Call, c)
	}

	req.TimeoutInMillis = timeoutInMillis
	req.HostId = lib.CurrentHostId().Marshal()
	resp, err:=syscallguest.ReadOne(&req)
	if err!=syscall.KernelErr_NoError {
		return err
	}
	// is timeout?
	if resp.Timeout {
		return syscall.KernelErr_ReadOneTimeout
	}

	// check for finished futures from within our address space
	syscallguest.ExpireMethod(ctx)

	// is a promise being completed that was fulfilled somewhere else
	if r:=resp.GetResolved(); r!=nil {
		cid:=id.UnmarshalCallId(r.GetCallId())
		syscallguest.CompleteCall(ctx, cid,r.GetResult(), r.GetResultError())
		return syscall.KernelErr_NoError
	}

	// its a method call from another address space
	sid:=id.UnmarshalServiceId(resp.GetCall().GetServiceId())
	mid:=id.UnmarshalMethodId(resp.GetCall().GetMethodId())
	cid:=id.UnmarshalCallId(resp.GetCallId())
	fn:=binding.Func(sid,mid)
	if fn==nil {
		pcontext.Errorf(ctx,"unable to find service/method pair binding: service=%s, method=%s",
			sid.Short(), mid.Short())
		return syscall.KernelErr_NotFound
	}
	// we let the invoker handle the unmarshal from anypb.Any because it
	// knows the precise type to be consumed
	fut:=fn.Invoke(ctx,resp.GetParam())
	fut.Success(func (result proto.Message){
		rvReq:=&syscall.ReturnValueRequest{}
		rvReq.CallId= cid.Marshal()
		rvReq.HostId= lib.CurrentHostId().Marshal()
		var a anypb.Any
		if err:=a.MarshalFrom(result); err!=nil {
			pcontext.Errorf(ctx, "unable to marshal result for return value request")
			return
		}
		rvReq.Result = &a
		rvReq.ResultError = 0
		syscallguest.ReturnValue(rvReq) // nowhere for return value to go
	})
	fut.Failure(func (err int32) {
		rvReq:=&syscall.ReturnValueRequest{}
		rvReq.CallId= cid.Marshal()
		rvReq.HostId= lib.CurrentHostId().Marshal()
		rvReq.ResultError = err
		syscallguest.ReturnValue(rvReq) // nowhere for return value to go
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
	bindReq.HostId = lib.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Open"
	resp, err=syscallguest.BindMethod(bindReq)
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
	bindReq.HostId = lib.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Create"
	resp, err=syscallguest.BindMethod(bindReq)
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
	bindReq.HostId = lib.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Close"
	resp, err=syscallguest.BindMethod(bindReq)
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
	bindReq.HostId = lib.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "LoadTestData"
	resp, err=syscallguest.BindMethod(bindReq)
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
	bindReq.HostId = lib.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Read"
	resp, err=syscallguest.BindMethod(bindReq)
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
	bindReq.HostId = lib.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Write"
	resp, err=syscallguest.BindMethod(bindReq)
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
	bindReq.HostId = lib.CurrentHostId().Marshal()
	bindReq.ServiceId = sid.Marshal()
	bindReq.MethodName = "Delete"
	resp, err=syscallguest.BindMethod(bindReq)
	if err!=syscall.KernelErr_NoError {
		return nil, err
	}
	mid=id.UnmarshalMethodId(resp.GetMethodId())

	// completer already prepared elsewhere
	smmap.AddServiceMethod(sid,mid,"File","Delete",
		GenerateDeleteInvoker(impl)) 
	pcontext.Dump(ctx)
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
        pcontext.Debugf(ctx,"kernel error was  %s",name)
        if err == syscall.KernelErr_NotRequired {
            pcontext.Errorf(ctx,"service was located, but it was not required")
            panic("locate attempted on a service that was not required")
        }
        panic(normal)
    }
    return result
}


func Register(ctx context.Context) (id.ServiceId, syscall.KernelErr){
    req := &syscall.RegisterRequest{}
	fqs := &syscall.FullyQualifiedService{
		PackagePath: "file.v1",
		Service:     "file",
	}
	req.Fqs = fqs
	req.HostId = lib.CurrentHostId().Marshal()

	resp, err := syscallguest.Register(req)
    if err!=syscall.KernelErr_NoError{
        return id.ServiceIdZeroValue(), err
    }
    sid:=id.UnmarshalServiceId(resp.Id)
    if sid.IsZeroOrEmptyValue() {
        panic("received bad service Id from register")
    }

    return sid,syscall.KernelErr_NoError
}
func MustRegister(ctx context.Context) id.ServiceId {
    sid, err:=Register(ctx)
    if err!=syscall.KernelErr_NoError {
        pcontext.Fatalf(ctx,"unable to register %s.%s","file.v1","file")
        panic("unable to register "+"file")
    }
    return sid
}

func MustRequire(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Require1("file.v1","file",sid)
    if err!=syscall.KernelErr_NoError {
        if err==syscall.KernelErr_DependencyCycle{
            pcontext.Errorf(ctx,"unable to require %s.%s because it creates a dependcy loop: %s","file.v1","file",syscall.KernelErr_name[int32(err)])
            panic("require file.v1.file creates a dependency loop")
        }
        pcontext.Errorf(ctx,"unable to require %s.%s:%s","file.v1","file",syscall.KernelErr_name[int32(err)])
        panic("not able to require file.v1.file:"+syscall.KernelErr_name[int32(err)])
    }
}

func MustExport(ctx context.Context, sid id.ServiceId) {
    _, err:=lib.Export1("file.v1","file",sid)
    if err!=syscall.KernelErr_NoError{
        pcontext.Fatalf(ctx, "unable to export %s.%s","file.v1","file")
        panic("not able to export file.v1.file:"+syscall.KernelErr_name[int32(err)])
    }
}

func LaunchService(ctx context.Context, sid id.ServiceId, impl File) (*lib.ServiceMethodMap,*syscallguest.LaunchFuture,syscall.KernelErr) {
	smmap, err:=bind(ctx,sid, impl)
	if err!=0{
		return  nil,nil,syscall.KernelErr(err)
	}
	cid:=id.NewCallId()
	req:=&syscall.LaunchRequest{
		ServiceId: sid.Marshal(),
		CallId: cid.Marshal(),
		HostId: lib.CurrentHostId().Marshal(),
		MethodId: apishared.LaunchMethod.Marshal(),
	}
	fut:=syscallguest.Launch(req)

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
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Open_)
	if signal {
		pcontext.Infof(ctx, "Open exiting because of parigot signal")
		pcontext.Dump(ctx)
		lib.ExitClient(ctx, 1, id.NewServiceId(), "xxx warning, no implementation of unsolicited exit",
			"xxx warning, no implementation of unsolicited exit and failed trying to exit")
	}
	f:=NewFutureOpen()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport file create_
func Create_(int32,int32,int32,int32) int64
func CreateHost(ctx context.Context,inPtr *CreateRequest) *FutureCreate {
	outProtoPtr := (*CreateResponse)(nil)
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Create_)
	if signal {
		pcontext.Infof(ctx, "Create exiting because of parigot signal")
		pcontext.Dump(ctx)
		lib.ExitClient(ctx, 1, id.NewServiceId(), "xxx warning, no implementation of unsolicited exit",
			"xxx warning, no implementation of unsolicited exit and failed trying to exit")
	}
	f:=NewFutureCreate()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport file close_
func Close_(int32,int32,int32,int32) int64
func CloseHost(ctx context.Context,inPtr *CloseRequest) *FutureClose {
	outProtoPtr := (*CloseResponse)(nil)
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Close_)
	if signal {
		pcontext.Infof(ctx, "Close exiting because of parigot signal")
		pcontext.Dump(ctx)
		lib.ExitClient(ctx, 1, id.NewServiceId(), "xxx warning, no implementation of unsolicited exit",
			"xxx warning, no implementation of unsolicited exit and failed trying to exit")
	}
	f:=NewFutureClose()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport file load_test_data_
func LoadTestData_(int32,int32,int32,int32) int64
func LoadTestDataHost(ctx context.Context,inPtr *LoadTestDataRequest) *FutureLoadTestData {
	outProtoPtr := (*LoadTestDataResponse)(nil)
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, LoadTestData_)
	if signal {
		pcontext.Infof(ctx, "LoadTestData exiting because of parigot signal")
		pcontext.Dump(ctx)
		lib.ExitClient(ctx, 1, id.NewServiceId(), "xxx warning, no implementation of unsolicited exit",
			"xxx warning, no implementation of unsolicited exit and failed trying to exit")
	}
	f:=NewFutureLoadTestData()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport file read_
func Read_(int32,int32,int32,int32) int64
func ReadHost(ctx context.Context,inPtr *ReadRequest) *FutureRead {
	outProtoPtr := (*ReadResponse)(nil)
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Read_)
	if signal {
		pcontext.Infof(ctx, "Read exiting because of parigot signal")
		pcontext.Dump(ctx)
		lib.ExitClient(ctx, 1, id.NewServiceId(), "xxx warning, no implementation of unsolicited exit",
			"xxx warning, no implementation of unsolicited exit and failed trying to exit")
	}
	f:=NewFutureRead()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport file write_
func Write_(int32,int32,int32,int32) int64
func WriteHost(ctx context.Context,inPtr *WriteRequest) *FutureWrite {
	outProtoPtr := (*WriteResponse)(nil)
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Write_)
	if signal {
		pcontext.Infof(ctx, "Write exiting because of parigot signal")
		pcontext.Dump(ctx)
		lib.ExitClient(ctx, 1, id.NewServiceId(), "xxx warning, no implementation of unsolicited exit",
			"xxx warning, no implementation of unsolicited exit and failed trying to exit")
	}
	f:=NewFutureWrite()
	f.CompleteMethod(ctx,ret,raw)
	return f
} 

//go:wasmimport file delete_
func Delete_(int32,int32,int32,int32) int64
func DeleteHost(ctx context.Context,inPtr *DeleteRequest) *FutureDelete {
	outProtoPtr := (*DeleteResponse)(nil)
	defer pcontext.Dump(ctx)
	ret, raw, signal:= syscallguest.ClientSide(ctx, inPtr, outProtoPtr, Delete_)
	if signal {
		pcontext.Infof(ctx, "Delete exiting because of parigot signal")
		pcontext.Dump(ctx)
		lib.ExitClient(ctx, 1, id.NewServiceId(), "xxx warning, no implementation of unsolicited exit",
			"xxx warning, no implementation of unsolicited exit and failed trying to exit")
	}
	f:=NewFutureDelete()
	f.CompleteMethod(ctx,ret,raw)
	return f
}  

// This is interface for invocation.
type invokeOpen struct {
    fn func(context.Context,*OpenRequest) *FutureOpen
}

func (t *invokeOpen) Invoke(ctx context.Context,a *anypb.Any) future.Completer {
	// xxx OpenRequest and 'OpenRequest{}' why empty?
    in:=&OpenRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
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
	// xxx CreateRequest and 'CreateRequest{}' why empty?
    in:=&CreateRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
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
	// xxx CloseRequest and 'CloseRequest{}' why empty?
    in:=&CloseRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
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
	// xxx LoadTestDataRequest and 'LoadTestDataRequest{}' why empty?
    in:=&LoadTestDataRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
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
	// xxx ReadRequest and 'ReadRequest{}' why empty?
    in:=&ReadRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
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
	// xxx WriteRequest and 'WriteRequest{}' why empty?
    in:=&WriteRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
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
	// xxx DeleteRequest and 'DeleteRequest{}' why empty?
    in:=&DeleteRequest{}
    err:=a.UnmarshalTo(in)
    if err!=nil {
        pcontext.Errorf(ctx,"unmarshal inside Invoke() failed: %s",err.Error())
        return nil
    }
    return t.fn(ctx,in) 

}

func GenerateDeleteInvoker(impl File) future.Invoker {
	return &invokeDelete{fn:impl.Delete} 
}  
