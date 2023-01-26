package sys

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/iansmith/parigot/api_impl/splitutil"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/sys/backdoor"
	"github.com/iansmith/parigot/sys/jspatch"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Flip this switch for debug output.
var envVerbose = os.Getenv("PARIGOT_VERBOSE")

var syscallVerbose = false || envVerbose != ""

// syscallReadWrite is the code that reads the parameters from the client side and responds to
// the client side via the same parameters. In between it calls either remote or local to implement
// the system call.
type syscallReadWrite struct {
	mem  *jspatch.WasmMem
	proc *Process // this is our process

	localSysCall  *localSysCall
	remoteSysCall *remoteSyscall

	ns NameServer
}

func splitImplRetOne[T, U proto.Message](mem *jspatch.WasmMem, sp int32, req T, resp U,
	fn func(t T, u U) (lib.Id, string)) {
	err := splitutil.StackPointerToRequest(mem, sp, req)
	if err != nil {
		return // the error return code is already set
	}
	id, errDetail := fn(req, resp)
	if id != nil {
		splitutil.ErrorResponse(mem, sp, id, errDetail)
		return
	}
	splitutil.RespondSingleProto(mem, sp, resp)
}

// SetMemPtr has to be separated out because at the time this object is created, we don't yet
// know the memory address that is the memPtr.
func (s *syscallReadWrite) SetMemPtr(m uintptr) {
	s.mem = jspatch.NewWasmMem(m)
}

// SetProcess has to be separated out because at the time this object is created, we don't yet
// know the process that this syscall works on behalf of.  This is because the process needs a
// syscall (this one) to be created.
func (s *syscallReadWrite) SetProcess(p *Process) {
	s.proc = p
}

func NewSysCallRW(ns NameServer) *syscallReadWrite {
	return &syscallReadWrite{ns: ns}
}

// Exit causes the WASM program to exit.  This is done by marking the process as dead and then
// doing panic to force the stack to unroll.... but that has to be done on the client (WASM)
// side, not here.
func (s *syscallReadWrite) Exit(sp int32) {
	panic("HERE")
	// req := &syscallmsg.ExitRequest{}
	// resp := &syscallmsg.ExitResponse{}
	// s.proc.exited = true
	// splitImplRetOne(s.mem, sp, req, resp,
	// 	func(req *syscallmsg.ExitRequest, resp *syscallmsg.ExitResponse) lib.KernelErrorCode {
	// 		code := req.GetCode()
	// 		if code > 192 || code < 0 {
	// 			code = 192
	// 		}
	// 		resp.Code = code
	// 		if s.ns.ExitWhenInFlightEmpty() { // tell the infra to exit after the current
	// 			return lib.KernelExitRequest
	// 		}
	// 		return lib.KernelNoError
	// 	})
}

// Locate is the syste call thet finds the service requested (or returns an error if it cannot be found )
// and creates an implementation of the proper interface to allow the caller to talk to that service no
// matter if that service is across the network or in the same address space.
func (s *syscallReadWrite) Locate(sp int32) {
	resp := syscallmsg.LocateResponse{}
	req := syscallmsg.LocateRequest{}
	err := splitutil.StackPointerToRequest(s.mem, sp, &req)
	if err != nil {
		return // the error return code is already set
	}
	sid, id, errDetail := s.procToSysCall().GetService(NewDepKeyFromProcess(s.proc),
		req.GetPackageName(), req.GetServiceName())

	if id != nil {
		splitutil.ErrorResponse(s.mem, sp, id, errDetail)
		return
	}
	resp.ServiceId = lib.Marshal[protosupportmsg.ServiceId](sid)
	splitutil.RespondSingleProto(s.mem, sp, &resp)
}

// Dispatch is the way that a client invokes and RPC to another service.  This code is on the kernel
// side (go implementation of kernel).
func (s *syscallReadWrite) Dispatch(sp int32) {
	resp := syscallmsg.DispatchResponse{}
	req := syscallmsg.DispatchRequest{}
	err := splitutil.StackPointerToRequest(s.mem, sp, &req)
	if err != nil {
		return // the error return code is already set
	}
	key := NewDepKeyFromProcess(s.proc)
	sid := lib.Unmarshal(req.GetServiceId())
	ctx, id, errDetail := s.procToSysCall().FindMethodByName(key, sid, req.Method)
	if id != nil {
		print("xxxctx is nil in Dispatch:" + key.String() + "," + req.Method + "\n")
		splitutil.ErrorResponse(s.mem, sp, id, errDetail)
		return
	}
	ctx.param = req.Param

	// this call is the machinery for making a call to another service
	retReq, id, errDetail := s.procToSysCall().CallService(ctx.target, ctx)
	if id != nil {
		splitutil.ErrorResponse(s.mem, sp, id, errDetail)
		return
	}
	resp.OutPctx = retReq.Pctx
	resp.Result = retReq.Result
	resp.MethodId = retReq.Method
	resp.CallId = retReq.Call
	if resp.Result == nil {
		splitutil.RespondSingleProto(s.mem, sp, &resp)
	} else {
		splitutil.RespondSingleProto(s.mem, sp, &resp)
	}
}

// BindMethod is used to indicate the function that will handle a given method.  We don't
// actually have a handle to the function pointer, we give out MethodIds instead.
func (s *syscallReadWrite) BindMethod(sp int32) {
	req := syscallmsg.BindMethodRequest{}
	resp := syscallmsg.BindMethodResponse{}
	splitImplRetOne(s.mem, sp, &req, &resp,
		func(req *syscallmsg.BindMethodRequest, resp *syscallmsg.BindMethodResponse) (lib.Id, string) {
			mid, kerr := s.procToSysCall().Bind(s.proc, req.GetProtoPackage(), req.GetService(), req.GetMethod())
			if kerr != nil {
				return kerr, fmt.Sprintf("unable to bind %s with method %s", req.GetProtoPackage(), req.GetMethod())
			}
			resp.MethodId = lib.Marshal[protosupportmsg.MethodId](mid)
			return nil, ""
		})
}

// BlockUntilCall is used by servers to block themselves until some other process sends them a message.
func (s *syscallReadWrite) BlockUntilCall(sp int32) {
	resp := syscallmsg.BlockUntilCallResponse{}
	req := syscallmsg.BlockUntilCallRequest{}
	splitImplRetOne(s.mem, sp, &req, &resp,
		func(req *syscallmsg.BlockUntilCallRequest, resp *syscallmsg.BlockUntilCallResponse) (lib.Id, string) {
			call := s.procToSysCall().BlockUntilCall(NewDepKeyFromProcess(s.proc), req.CanTimeout)
			if call.timedOut {
				resp.TimedOut = true
				return nil, ""
			}
			resp.TimedOut = false
			resp.Param = call.param
			resp.Pctx = call.pctx
			resp.Call = lib.Marshal[protosupportmsg.CallId](call.cid)
			resp.Method = lib.Marshal[protosupportmsg.MethodId](call.mid)
			return nil, ""
		})
}

func splitImplRetEmpty[T proto.Message](mem *jspatch.WasmMem, sp int32, req T, fn func(t T) (lib.Id, string)) {
	err := splitutil.StackPointerToRequest(mem, sp, req)
	if err != nil {
		return // the error return code is already set
	}
	id, errDetail := fn(req)
	if id != nil {
		splitutil.ErrorResponse(mem, sp, id, errDetail)
		return
	}
	splitutil.RespondEmpty(mem, sp)
}

// Return value is used by server implementations to set the return value for a particular function
// call on a method they implement.
func (s *syscallReadWrite) ReturnValue(sp int32) {
	req := syscallmsg.ReturnValueRequest{}
	resp := syscallmsg.ReturnValueResponse{}
	splitImplRetOne(s.mem, sp, &req, &resp, func(t *syscallmsg.ReturnValueRequest, u *syscallmsg.ReturnValueResponse) (lib.Id, string) {
		cid := lib.Unmarshal(req.GetCall())
		ctx := s.procToSysCall().GetInfoForCallId(cid)
		if ctx == nil {
			sysPrint(logmsg.LogLevel_LOG_LEVEL_WARNING, "RETURNVALUE ", "no record of that call (caller addr %v)", ctx.sender.(*DepKeyImpl).addr)
			return lib.NewKernelError(lib.KernelCallerUnavailable),
				fmt.Sprintf("no record of that call (caller addr %v)", ctx.sender.(*DepKeyImpl).addr)
		}
		callerProc := ctx.sender.(*DepKeyImpl).proc
		if callerProc == nil {
			sysPrint(logmsg.LogLevel_LOG_LEVEL_WARNING, "RETURNVALUE ", "no caller proc, caller addr is %s", ctx.sender.(*DepKeyImpl).addr)
			return lib.NewKernelError(lib.KernelCallerUnavailable),
				fmt.Sprintf("no caller proc, caller addr is %s", ctx.sender.(*DepKeyImpl).addr)
		}
		ctx.respCh <- &req

		return nil, ""
	})
}

func (s *syscallReadWrite) procToSysCall() SysCall {
	if s.proc.microservice.IsLocal() {
		if s.localSysCall != nil {
			return s.localSysCall
		}
		s.localSysCall = newLocalSysCall(LocalNS)
		return s.localSysCall
	}
	if s.remoteSysCall != nil {
		return s.remoteSysCall
	}
	s.remoteSysCall = newRemoteSysCall(NetNS)
	return s.remoteSysCall
}

// Export is used when a server wishes to express what services
// he implements and that he has finished binding all the methods
// of that service.q
func (s *syscallReadWrite) Export(sp int32) {
	req := &syscallmsg.ExportRequest{}
	splitImplRetEmpty(s.mem, sp, req, func(req *syscallmsg.ExportRequest) (lib.Id, string) {
		service := req.GetService()
		for _, svc := range service {
			// xxx  fixme what should we do in the face of some succeeding some not?
			kerr := s.procToSysCall().Export(NewDepKeyFromProcess(s.proc), svc.GetPackagePath(), svc.GetService())
			if kerr != nil {
				return kerr, fmt.Sprintf("unable to export %s.%s", svc.PackagePath, svc.Service)
			}
		}
		return nil, ""
	})
}

// Require is used when client or server wishes to indicate that it consumes
// a service.  This becomes part of the dependency graph.
func (s *syscallReadWrite) Require(sp int32) {
	req := &syscallmsg.RequireRequest{}
	splitImplRetEmpty(s.mem, sp, req, func(req *syscallmsg.RequireRequest) (lib.Id, string) {
		service := req.GetService()
		for _, svc := range service {
			// xxx  fixme what should we do in the face of some succeeding some not?
			kerr := s.procToSysCall().Require(NewDepKeyFromProcess(s.proc), svc.GetPackagePath(), svc.GetService())
			if kerr != nil {
				return kerr, fmt.Sprintf("failed to require %+v", req.GetService())
			}
		}
		return nil, ""
	})
}

// Run is used to start up the processes in a deterministic order. It will
// fail and return an error if there are problems getting all the require and export
// requests to match up.
func (s *syscallReadWrite) Run(sp int32) {
	req := &syscallmsg.RunRequest{}
	splitImplRetEmpty(s.mem, sp, req, func(req *syscallmsg.RunRequest) (lib.Id, string) {
		sysPrint(logmsg.LogLevel_LOG_LEVEL_DEBUG, "Run", "about to call new implementation of run inside nameserver")
		ok, err := s.ns.RunBlock(s.proc.key)
		if err != nil && err.IsErrorType() && err.IsError() {
			return err, "failed when doing RunBlock " + s.proc.key.String()
		}
		if !ok {
			return lib.NewKernelError(lib.KernelAbortRequest), "received abort request"
		}
		return nil, ""
	})
}

func sysPrint(level logmsg.LogLevel, call, spec string, arg ...interface{}) {
	if syscallVerbose {
		spec = "%s:" + spec
		arg = append([]interface{}{call}, arg...)
		msg := fmt.Sprintf(spec, arg...)
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		req := &logmsg.LogRequest{
			Stamp:   timestamppb.New(time.Now()),
			Level:   level,
			Message: msg,
		}
		backdoor.Log(req, true, false, false, nil)

	}
}

// BackdoorLog is only for use by internal parigot debugging messages.
func (s *syscallReadWrite) BackdoorLog(sp int32) {
	req := &logmsg.LogRequest{}
	wasmPtr := s.mem.GetInt64(sp + 8)

	buffer := splitutil.ReadSlice(s.mem, wasmPtr,
		unsafe.Offsetof(splitutil.SinglePayload{}.InPtr),
		unsafe.Offsetof(splitutil.SinglePayload{}.InLen))

	kerr := lib.NewKernelError(lib.KernelUnmarshalFailed)
	if err := proto.Unmarshal(buffer, req); err != nil {
		s.mem.SetInt32(int32(wasmPtr)+int32(unsafe.Offsetof(splitutil.SinglePayload{}.ErrPtr)+8),
			int32(kerr.High()))
		s.mem.SetInt32(int32(wasmPtr)+int32(unsafe.Offsetof(splitutil.SinglePayload{}.ErrPtr)+0),
			int32(kerr.Low()))
		return
	}
	backdoor.Log(req, true, false, false, nil)
	s.mem.SetInt32(int32(wasmPtr)+int32(unsafe.Offsetof(splitutil.SinglePayload{}.OutPtr)), 0)
	s.mem.SetInt32(int32(wasmPtr)+int32(unsafe.Offsetof(splitutil.SinglePayload{}.OutLen)), 0)

}
