package syscall

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"time"
	_ "unsafe"

	"github.com/iansmith/parigot/apiplugin"
	"github.com/iansmith/parigot/apishared"
	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/eng"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/tetratelabs/wazero/api"
)

var serviceNameToId = make(map[string]id.ServiceId)
var serviceIdToName = make(map[string]string)

var serviceIdToMethodNameMap = make(map[string]map[string]id.MethodId)
var serviceIdToMethodIdMap = make(map[string]map[string]string)

var pairIdToChannel = make(map[string]chan anypb.Any)

type SyscallPlugin struct {
}

var ParigotInitialize apiplugin.ParigotInit = &SyscallPlugin{}

func (*SyscallPlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "parigot", "locate_", locate)
	e.AddSupportedFunc(ctx, "parigot", "dispatch_", dispatch)
	e.AddSupportedFunc(ctx, "parigot", "block_until_call_", blockUntilCall)
	e.AddSupportedFunc(ctx, "parigot", "bind_method_", bindMethod)
	e.AddSupportedFunc(ctx, "parigot", "launch_", launch)
	e.AddSupportedFunc(ctx, "parigot", "export_", export)
	e.AddSupportedFunc(ctx, "parigot", "return_value_", returnValue)
	e.AddSupportedFunc(ctx, "parigot", "require_", require)
	e.AddSupportedFunc(ctx, "parigot", "register_", register)
	e.AddSupportedFunc(ctx, "parigot", "exit_", exit)
	e.AddSupportedFunc(ctx, "parigot", "read_one_", readOne)

	return true
}

//
// Syscall host implementations
//

func exportImpl(ctx context.Context, req *syscall.ExportRequest, resp *syscall.ExportResponse) int32 {
	for _, fullyQualified := range req.GetService() {
		sid, _ := coordinator().SetService(ctx, fullyQualified.GetPackagePath(), fullyQualified.GetService(), false)

		if coordinator().Export(ctx, sid.Id()) == nil {
			return int32(syscall.KernelErr_NotFound)
		}
	}
	return int32(syscall.KernelErr_NoError)

}

func launchImpl(ctx context.Context, req *syscall.LaunchRequest, resp *syscall.LaunchResponse) int32 {
	sid := id.UnmarshalServiceId(req.GetServiceId())
	return int32(coordinator().Launch(ctx, sid))
}

func bindMethodImpl(ctx context.Context, req *syscall.BindMethodRequest, resp *syscall.BindMethodResponse) int32 {
	sid := id.UnmarshalServiceId(req.GetServiceId())
	mid, err := addMethodByName(ctx, sid, req.GetMethodName())
	if err != syscall.KernelErr_NoError {
		return int32(err)
	}
	resp.MethodId = mid.Marshal()
	pairIdToChannel[sid.String()+mid.String()] = make(chan anypb.Any)
	return int32(syscall.KernelErr_NoError)
}

var exampleTime time.Time

func readOneImpl(ctx context.Context, req *syscall.ReadOneRequest, resp *syscall.ReadOneResponse) int32 {
	numCases := len(req.Pair)
	if req.TimeoutInMillis >= 0 {
		numCases++
	}
	cases := make([]reflect.SelectCase, numCases)
	if numCases == 0 {
		resp.Pair = nil
		return int32(syscall.KernelErr_NoError)
	}
	for i, pair := range req.Pair {
		svc := id.UnmarshalServiceId(pair.ServiceId)
		meth := id.UnmarshalMethodId(pair.MethodId)
		combo := svc.String() + meth.String()
		ch := pairIdToChannel[combo]
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	if req.TimeoutInMillis >= 0 {
		ch := time.After(time.Duration(req.TimeoutInMillis) * time.Millisecond)
		cases[len(req.Pair)] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	chosen, value, ok := reflect.Select(cases)
	// ok will be true if the channel has not been closed.
	if !ok {
		return int32(syscall.KernelErr_KernelConnectionFailed)
	}
	if chosen == len(req.Pair) {
		resp.Timeout = true
		return int32(syscall.KernelErr_NoError)
	}
	resp.Timeout = false
	pair := req.Pair[chosen]
	resp.Pair.ServiceId = pair.ServiceId
	resp.Pair.MethodId = pair.MethodId
	if !value.IsNil() {
		resp.Param = value.Interface().(*anypb.Any)
	}
	return int32(syscall.KernelErr_NoError)
}

func locateImpl(ctx context.Context, req *syscall.LocateRequest, resp *syscall.LocateResponse) int32 {
	svc, ok := coordinator().SetService(ctx, req.GetPackageName(), req.GetServiceName(), false)
	if ok {
		return int32(syscall.KernelErr_NotFound)
	}
	calledBy := id.UnmarshalServiceId(req.CalledBy)
	if !coordinator().PathExists(ctx, calledBy.String(), svc.String()) {
		return int32(syscall.KernelErr_NotRequired)
	}
	svcId := svc.Id()
	resp.ServiceId = svcId.Marshal()
	return int32(syscall.KernelErr_NoError)
}

func dispatchImpl(ctx context.Context, req *syscall.DispatchRequest, resp *syscall.DispatchResponse) int32 {
	midToName, ok := serviceIdToMethodIdMap[req.GetServiceId().String()]
	if !ok {
		return int32(syscall.KernelErr_NotFound)
	}
	srvId := serviceIdToName[req.GetServiceId().String()]
	methId := midToName[req.GetMethodId().String()]

	svrName := serviceIdToName[srvId]
	methName := serviceIdToMethodIdMap[srvId][methId]
	fqName := fmt.Sprintf("%s.%s", svrName, methName)
	path := filepath.Join(apishared.FsName, fqName)

	result, cid, err := eng.AsyncInteraction.Dispatch(path, req.GetParam())
	if err != syscall.KernelErr_NoError {
		return int32(err)
	}
	resp.Result = result
	resp.CallId = cid.Marshal()
	return int32(syscall.KernelErr_NoError)
}

func registerImpl(ctx context.Context, req *syscall.RegisterRequest, resp *syscall.RegisterResponse) int32 {
	svc, _ := coordinator().SetService(ctx, req.Fqs.GetPackagePath(), req.Fqs.GetService(), req.GetIsClient())

	sname := fmt.Sprintf("%s.%s", req.Fqs.GetPackagePath(), req.Fqs.GetService())
	serviceNameToId[sname] = svc.Id()
	serviceIdToName[svc.Id().String()] = sname

	serviceIdToMethodNameMap[svc.Id().String()] = make(map[string]id.MethodId)
	serviceIdToMethodIdMap[svc.Id().String()] = make(map[string]string)

	resp.Id = svc.Id().Marshal()
	return int32(syscall.KernelErr_NoError)
}

func requireImpl(ctx context.Context, req *syscall.RequireRequest, resp *syscall.RequireResponse) int32 {
	src := id.UnmarshalServiceId(req.GetSource())
	fqn := req.GetDest()
	//log.Printf("xxx--- require impl, source = %s", src.Short())

	for _, fullyQualified := range fqn {
		dest, _ := coordinator().SetService(ctx, fullyQualified.GetPackagePath(), fullyQualified.GetService(), false)
		// if ok {
		// 	pcontext.Infof(ctx, "requireImpl: created new service id %s.%s => %s", fullyQualified.GetPackagePath(),
		// 		fullyQualified.GetService(), dest.Short())
		// }
		kerr := coordinator().Import(ctx, src, dest.Id())
		if int32(kerr) != 0 {
			pcontext.Errorf(ctx, "kernel error returned from import: %d", kerr)
			return int32(kerr)
		}
	}
	return int32(syscall.KernelErr_NoError)
}

//
// Syscall marshal/unmarshal for each system call
//

func locate(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.LocateRequest{}
	resp := &syscall.LocateResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[syscall]locate", m, stack, locateImpl, req, resp)
}

func dispatch(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.DispatchRequest{}
	resp := &syscall.DispatchResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[syscall]dispatch", m, stack, dispatchImpl, req, resp)
}

func blockUntilCall(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("blockUntilCall 0x%x", stack)
}
func bindMethod(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.BindMethodRequest{}
	resp := &syscall.BindMethodResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[syscall]bindMethod", m, stack, bindMethodImpl, req, resp)

}
func launch(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.LaunchRequest{}
	resp := (*syscall.LaunchResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]export", m, stack, launchImpl, req, resp)
}
func export(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.ExportRequest{}
	resp := (*syscall.ExportResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]export", m, stack, exportImpl, req, resp)
}
func returnValue(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("returnValue 0x%x", stack)
}

func require(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.RequireRequest{}
	resp := (*syscall.RequireResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]require", m, stack, requireImpl, req, resp)

}
func readOne(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.ReadOneRequest{}
	resp := (*syscall.ReadOneResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]readOne", m, stack, readOneImpl, req, resp)

}

func register(ctx context.Context, m api.Module, stack []uint64) {
	req := &syscall.RegisterRequest{}
	resp := &syscall.RegisterResponse{}
	apiplugin.InvokeImplFromStack(ctx, "[syscall]register", m, stack, registerImpl, req, resp)
}

func exit(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("exit 0x%x", stack)
	panic("exit called ")
}

// addMethodByName adds the new method name that is name inside the given service id,
// to the internal datastructures.  This creates an id for the method and returns
// that new id.  If the name already existed in our internal structures then
// we return the already know method id.  This method will return KernelErr_NotFound
// only in the case where the service given by the service id cannot be found.
func addMethodByName(ctx context.Context, serviceId id.ServiceId, methodName string) (id.MethodId, syscall.KernelErr) {

	methMapId, ok := serviceIdToMethodIdMap[serviceId.String()]
	methMapName, ok := serviceIdToMethodNameMap[serviceId.String()]
	if !ok {
		log.Printf("unable to find service %s, cannot add method %s", serviceId.Short(), methodName)
		return id.MethodIdZeroValue(), syscall.KernelErr_NotFound
	}
	var newId id.MethodId
	_, ok = methMapId[methodName]
	if !ok {
		newId = id.NewMethodId()
		methMapId[newId.String()] = methodName
		methMapName[methodName] = newId
	}
	return newId, syscall.KernelErr_NoError
}
