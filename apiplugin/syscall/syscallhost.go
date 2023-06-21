package syscall

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"
	_ "unsafe"

	"github.com/iansmith/parigot/apiplugin"
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

var pairIdToChannel = make(map[string]chan CallInfo)

// retValPair is the two return values that need to go back to the
// proper call site on the client side. Calls to returnValue will
// set these for later transmission to client.
type retValPair struct {
	Result *anypb.Any
	Error  int32
}

// key in pending calls is cid
var pendingCalls = make(map[string]DispatchInfo)

// DispatchInfo
type DispatchInfo struct {
	cid id.CallId
	ch  chan retValPair
}

// CallInfo is sent to the channels that represent service/method calls.
type CallInfo struct {
	cid   id.CallId
	param *anypb.Any
}

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
	mid := id.NewMethodId()
	resp.MethodId = mid.Marshal()
	svc := coordinator().ServiceById(ctx, sid)
	svc.AddMethod(req.GetMethodName(), mid)

	pairIdToChannel[makeSidMidCombo(sid, mid)] = make(chan CallInfo, 1)
	return int32(syscall.KernelErr_NoError)
}

func readOneImpl(ctx context.Context, req *syscall.ReadOneRequest, resp *syscall.ReadOneResponse) int32 {
	numCases := len(req.Pair) + len(pendingCalls)
	if req.TimeoutInMillis >= 0 {
		numCases++
	}
	if numCases == 0 {
		resp.Pair = nil
		return int32(syscall.KernelErr_NoError)
	}
	flatten := make([]DispatchInfo, len(pendingCalls))
	cases := make([]reflect.SelectCase, numCases)
	firstPending := len(req.Pair)
	for i, pair := range req.Pair {
		svc := id.UnmarshalServiceId(pair.ServiceId)
		meth := id.UnmarshalMethodId(pair.MethodId)
		combo := makeSidMidCombo(svc, meth)
		ch := pairIdToChannel[combo]
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	count := len(req.Pair)
	for _, info := range pendingCalls {
		cases[count] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(info.ch)}
		flatten[count-len(req.Pair)] = info
		count++
	}
	if req.TimeoutInMillis >= 0 {
		ch := time.After(time.Duration(req.TimeoutInMillis) * time.Millisecond)
		cases[len(req.Pair)+len(pendingCalls)] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	chosen, value, ok := reflect.Select(cases)
	// ok will be true if the channel has not been closed.
	if !ok {
		return int32(syscall.KernelErr_KernelConnectionFailed)
	}
	// is timeout?
	if chosen == len(req.Pair)+len(pendingCalls) {
		resp.Timeout = true
		return int32(syscall.KernelErr_NoError)
	}
	// is service/method call?
	resp.Timeout = false
	if chosen < firstPending {
		resp.DispatchResult = false
		pair := req.Pair[chosen]
		resp.Pair.ServiceId = pair.ServiceId
		resp.Pair.MethodId = pair.MethodId
		if !value.IsNil() {
			resp.Param = value.Interface().(*anypb.Any)
		}
		return int32(syscall.KernelErr_NoError)
	}
	// must be a dispatch
	resp.DispatchResult = true
	resp.Pair = nil
	msg := value.Interface().(*retValPair)
	resp.CallId = flatten[chosen].cid.Marshal()
	resp.Param = msg.Result
	resp.ResultError = msg.Error
	return int32(syscall.KernelErr_NoError)
}

func returnValueImpl(ctx context.Context, req *syscall.ReturnValueRequest, resp *syscall.ReturnValueResponse) int32 {
	cid := id.UnmarshalCallId(req.GetCallId())
	ret := retValPair{
		Result: req.GetResult(),
		Error:  req.GetResultError(),
	}
	call, ok := pendingCalls[cid.String()]
	if !ok {
		resp.Matched = false
		return int32(0)
	}
	call.ch <- ret
	resp.Matched = true
	delete(pendingCalls, cid.String())
	return int32(0)
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
	resp.Binding = svc.Method()
	return int32(syscall.KernelErr_NoError)
}

func dispatchImpl(ctx context.Context, req *syscall.DispatchRequest, resp *syscall.DispatchResponse) int32 {
	sid := id.UnmarshalServiceId(req.GetServiceId())
	mid := id.UnmarshalMethodId(req.GetMethodId())
	cid := id.UnmarshalCallId(req.GetCallId())

	target := pairIdToChannel[makeSidMidCombo(sid, mid)]
	if target == nil {
		return int32(syscall.KernelErr_NotFound)
	}

	resp.CallId = cid.Marshal()

	ch := make(chan retValPair, 1)
	pendingCalls[cid.String()] = DispatchInfo{
		cid: cid,
		ch:  ch,
	}
	cm := CallInfo{
		cid:   cid,
		param: req.GetParam(),
	}
	target <- cm
	return int32(syscall.KernelErr_NoError)
}

func registerImpl(ctx context.Context, req *syscall.RegisterRequest, resp *syscall.RegisterResponse) int32 {
	svc, ok := coordinator().SetService(ctx, req.Fqs.GetPackagePath(), req.Fqs.GetService(), req.GetIsClient())
	resp.ExistedPreviously = !ok
	resp.Id = svc.Id().Marshal()

	sname := fmt.Sprintf("%s.%s", req.Fqs.GetPackagePath(), req.Fqs.GetService())
	serviceNameToId[sname] = svc.Id()
	serviceIdToName[svc.Id().String()] = sname

	serviceIdToMethodNameMap[svc.Id().String()] = make(map[string]id.MethodId)
	serviceIdToMethodIdMap[svc.Id().String()] = make(map[string]string)

	return int32(syscall.KernelErr_NoError)
}

func requireImpl(ctx context.Context, req *syscall.RequireRequest, resp *syscall.RequireResponse) int32 {
	if req.GetDest() == nil {
		log.Printf("ignoring call to Require because the require list is empty")
		return 0
	}
	src := id.UnmarshalServiceId(req.GetSource())
	fqn := req.GetDest()

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
	req := &syscall.ReturnValueRequest{}
	resp := (*syscall.ReturnValueResponse)(nil)
	apiplugin.InvokeImplFromStack(ctx, "[syscall]returnValue", m, stack, returnValueImpl, req, resp)

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

func makeSidMidCombo(sid id.ServiceId, mid id.MethodId) string {
	return sid.String() + "," + mid.String()
}
