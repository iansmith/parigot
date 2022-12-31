package sys

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unsafe"

	ilog "github.com/iansmith/parigot/api/logimpl/go_"
	pblog "github.com/iansmith/parigot/api/proto/g/pb/log"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	pbsys "github.com/iansmith/parigot/api/proto/g/pb/syscall"
	"github.com/iansmith/parigot/api/splitutil"
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/jspatch"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Flip this switch for debug output.
var envVerbose = os.Getenv("PARIGOT_VERBOSE")

var syscallVerbose = true || envVerbose != ""

// syscallReadWrite is the code that reads the parameters from the client side and responds to
// the client side via the same parameters. In between it calls either remote or local to implement
// the system call.
type syscallReadWrite struct {
	mem  *jspatch.WasmMem
	proc *Process // this is OUR PROCESS

	localSysCall  *localSysCall
	remoteSysCall *remoteSyscall
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

func NewSysCallRW() *syscallReadWrite {
	return &syscallReadWrite{}
}

// Exit does an os.Exit(0) which is bad.  There really should be some more signaling because sometimes one service going down
// needs to alert other services or monitoring software.  Further, exiting the whloe program is probably not what
// you want in the "all in one address space version" since it brings _all_ the services down.
func (a *syscallReadWrite) Exit(sp int32) {
	log.Printf("exit called")
	a.proc.exited = true
	// xxx fixme how can we cause a return here instead of using exit()?
	os.Exit(int(0))
}

// Locate is the syste call thet finds the service requested (or returns an error if it cannot be found )
// and creates an implementation of the proper interface to allow the caller to talk to that service no
// matter if that service is across the network or in the same address space.
func (s *syscallReadWrite) Locate(sp int32) {
	resp := pbsys.LocateResponse{}
	req := pbsys.LocateRequest{}
	err := splitutil.StackPointerToRequest(s.mem, sp, &req)
	if err != nil {
		return // the error return code is already set
	}
	sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "LOCATE ", "locate requested for %s.%s",
		req.GetPackageName(), req.GetServiceName())

	sid, kcode := s.procToSysCall().GetService(NewDepKeyFromProcess(s.proc),
		req.GetPackageName(), req.GetServiceName())

	if kcode != lib.KernelNoError {
		print(fmt.Sprintf("xxx Locate got kcode back from locate:%x... headed to error response\n", kcode))
		splitutil.ErrorResponse(s.mem, sp, kcode)
		return
	}
	resp.ServiceId = lib.Marshal[protosupport.ServiceId](sid)
	sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "Locate", "respond single proto about to be called with %#v,size %d",
		&resp, proto.Size(&resp))
	splitutil.RespondSingleProto(s.mem, sp, &resp)
}

// Dispatch is the way that a client invokes and RPC to another service.  This code is on the kernel
// side (go implementation of kernel).
func (s *syscallReadWrite) Dispatch(sp int32) {
	resp := pbsys.DispatchResponse{}
	req := pbsys.DispatchRequest{}
	err := splitutil.StackPointerToRequest(s.mem, sp, &req)
	if err != nil {
		return // the error return code is already set
	}
	key := NewDepKeyFromProcess(s.proc)
	sid := lib.Unmarshal(req.GetServiceId())
	sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "Dispatch ", "find method by name requested for %s.%s",
		sid.Short(), req.Method)
	ctx := s.procToSysCall().FindMethodByName(key, sid, req.Method)
	ctx.param = req.Param
	sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "Dispatch", "after making ctx with FindMethodByName:method=%s,cid=%s,mid=%s,type of param=%s", ctx.method, ctx.cid.Short(), ctx.mid.Short(), ctx.param.TypeUrl)

	if ctx.pctx == nil {
		sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "Dispatch ", "xxx call context does not have pctx set but size of param is %d",
			proto.Size(ctx.param))
	}

	// this call is the machinery for making a call to another service
	sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "Dispatch", "about to hit call service")
	retReq := s.procToSysCall().CallService(ctx.target, ctx)
	if retReq.Result != nil {
		sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "Dispatch ", "--result info returned from call service: %#v, with result type %s and size %d", retReq,
			retReq.Result.TypeUrl, proto.Size(retReq.Result))
	} else {
		sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "Dispatch ", "--result info  retValue is nil")
	}
	if retReq.ExecErrorId != nil {
		errid := lib.Unmarshal(retReq.ExecErrorId)
		if errid != nil {
			if errid.IsErrorType() && errid.IsError() {
				kernErr := lib.KernelErrorCode(errid.Low())
				splitutil.ErrorResponse(s.mem, sp, kernErr)
				sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "Dispatch ", "error id from other side of the call:%s",
					errid.Short())
				return
			} else {
				panic("dispatch is unable to understand result error of type:" + fmt.Sprintf("%T", errid))
			}
		}
	}
	if retReq.MarshalError != "" {
		sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "Dispatch ", "marshal error from other side of the call:%s",
			retReq.MarshalError)
		splitutil.ErrorResponse(s.mem, sp, lib.KernelMarshalFailed)
		return
	}
	if retReq.ExecError != "" {
		sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "Dispatch ", "exec error from other side of the call:%s",
			retReq.ExecError)
		splitutil.ErrorResponse(s.mem, sp, lib.KernelExecError)
		return
	}
	resp.OutPctx = retReq.Pctx
	resp.Result = retReq.Result
	resp.MethodId = retReq.Method
	resp.CallId = retReq.Call
	if resp.Result == nil {
		sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "Dispatch", "returning single proto but resp.Result is nil")
		splitutil.RespondSingleProto(s.mem, sp, &resp)
	} else {
		sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "Dispatch", "returning single proto and resp.Result is %T", resp.Result)
		splitutil.RespondSingleProto(s.mem, sp, &resp)
	}
}

func (s *syscallReadWrite) sendKernelErrorFromBind(wasmPtr int64, code lib.KernelErrorCode) {
	dispErr := lib.NewKernelError(code)
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.ErrorPtr),
		dispErr)
	return
}

// BindMethod is used to indicate the function that will handle a given method.  We don't
// actually have a handle to the function pointer, we give out MethodIds instead.
func (s *syscallReadWrite) BindMethod(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)
	sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "BINDMETHOD ", "wasmptr %x, true=%x", wasmPtr, s.mem.TrueAddr(int32(wasmPtr)))

	pkg := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.BindPayload{}.PkgPtr),
		unsafe.Offsetof(lib.BindPayload{}.PkgLen))

	service := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.BindPayload{}.ServicePtr),
		unsafe.Offsetof(lib.BindPayload{}.ServiceLen))

	method := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.BindPayload{}.MethodPtr),
		unsafe.Offsetof(lib.BindPayload{}.MethodLen))

	sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "BINDMETHOD", "about to bind %s in service %s", method, service)
	mid, kerr := s.procToSysCall().Bind(s.proc, pkg, service, method)
	if kerr != nil {
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.ErrorPtr),
			kerr)
		return
	}

	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.MethodId), mid)
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.ErrorPtr),
		lib.NoError[*protosupport.KernelErrorId]())

	sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "BINDMETHOD", "bind completed, %s bound to %s", method, mid.Short())
}

// BlockUntilCall is used by servers to block themselves until some other process sends them a message.
func (s *syscallReadWrite) BlockUntilCall(sp int32) {
	sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "BlockUntilCall--", "starting on %s", s.proc.String())
	resp := pbsys.BlockUntilCallResponse{}
	req := pbsys.BlockUntilCallRequest{}
	err := splitutil.StackPointerToRequest(s.mem, sp, &req)
	if err != nil {
		return // the error return code is already set
	}

	sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "BlockUntilCall--", "about to hit the nameserver %s", s.proc.String())
	// tricky, we have to pass a process here
	call := s.procToSysCall().BlockUntilCall(NewDepKeyFromProcess(s.proc))
	resp.Param = call.param
	resp.Pctx = call.pctx
	resp.Call = lib.Marshal[protosupport.CallId](call.cid)
	resp.Method = lib.Marshal[protosupport.MethodId](call.mid)
	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "BLOCKUNTILCALL", "Block until call client side finished OK with call to %s and param?%d", call.method, proto.Size(resp.Param))
	return // the server goes back to work
}

// Return value is used by server implementations to set the return value for a particular function
// call on a method they implement.
func (s *syscallReadWrite) ReturnValue(sp int32) {
	req := pbsys.ReturnValueRequest{}
	err := splitutil.StackPointerToRequest(s.mem, sp, &req)
	cid := lib.Unmarshal(req.Call)
	mid := lib.Unmarshal(req.Method)
	if err != nil {
		return // the error return code is already set
	}
	sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "ReturnValue ", "request accepted for method %s, callId %s",
		mid.Short(), cid.Short())

	ctx := s.procToSysCall().GetInfoForCallId(cid)
	if ctx == nil {
		sysPrint(pblog.LogLevel_LOG_LEVEL_ERROR, "RETURNVALUE", "unable to find process/addr that called %s", cid.Short())
		splitutil.ErrorResponse(s.mem, sp, lib.KernelCallerUnavailable)
		return
	}
	callerProc := ctx.sender.(*DepKeyImpl).proc

	if callerProc == nil {
		sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "RETURNVALUE ", "no caller proc, caller addr is %s", ctx.sender.(*DepKeyImpl).addr)
	}
	ctx.respCh <- &req
	splitutil.RespondEmpty(s.mem, sp)

	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "RESULTVALUE ", "sent return request and finished with no error")
	return
}

func (s *syscallReadWrite) procToSysCall() SysCall {
	if *s.proc.local {
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
	wasmPtr := s.mem.GetInt64(sp + 8)
	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "EXPORT", "wasmptr %x,true=%x", wasmPtr, s.mem.TrueAddr(int32(wasmPtr)))

	pkg := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.ExportPayload{}.PkgPtr),
		unsafe.Offsetof(lib.ExportPayload{}.PkgLen))

	service := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.ExportPayload{}.ServicePtr),
		unsafe.Offsetof(lib.ExportPayload{}.ServiceLen))

	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "EXPORT", "about to tell the nameserver closeservice and  export %s.%s [syscall->%T]", pkg, service,
		s.procToSysCall())
	kerr := s.procToSysCall().Export(NewDepKeyFromProcess(s.proc), pkg, service)
	if kerr != nil && kerr.IsError() {
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.ExportPayload{}.KernelErrorPtr),
			kerr)
		return
	}
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.ExportPayload{}.KernelErrorPtr),
		lib.NoError[*protosupport.KernelErrorId]())

	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "EXPORT", "done")

}

// Require is used when client or server wishes to indicate that it consumes
// a service.  This becomes part of the dependency graph.
func (s *syscallReadWrite) Require(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)
	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "REQUIRE", "wasmptr %x,true=%x", wasmPtr, s.mem.TrueAddr(int32(wasmPtr)))

	// we just have to create the structure and send it through the correct channel
	pkg := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.RequirePayload{}.PkgPtr),
		unsafe.Offsetof(lib.RequirePayload{}.PkgLen))

	service := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.RequirePayload{}.ServicePtr),
		unsafe.Offsetof(lib.RequirePayload{}.ServiceLen))

	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "REQUIRE", "telling nameserver %s requires %s.%s", s.proc, pkg, service)

	kerr := s.procToSysCall().Require(NewDepKeyFromProcess(s.proc), pkg, service)

	if kerr != nil && kerr.IsError() {
		sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "REQUIRE", " nameserver failed require of %s.%s", pkg, service)
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.RequirePayload{}.KernelErrorPtr),
			kerr)
		return
	}
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.RequirePayload{}.KernelErrorPtr),
		lib.NoError[*protosupport.KernelErrorId]())

}

// Run is used to start up the processes in a deterministic order. It will
// fail and return an error if there are problems getting all the require and export
// requests to match up.
func (s *syscallReadWrite) Run(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)
	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "RUN", "wasmptr %x,true=%x", wasmPtr, s.mem.TrueAddr(int32(wasmPtr)))
	w := s.mem.GetInt64(int32(wasmPtr) + int32(unsafe.Offsetof(lib.RunPayload{}.Wait)))
	wait := false
	if w != 0 {
		wait = true
	}
	s.proc.waiter = wait
	s.proc.reachedRun = true

	key := NewDepKeyFromProcess(s.proc)
	s.procToSysCall().RunNotify(key)

	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "RUN", "%s is blocked on channel for run confirmation", s.proc)
	// block until we are told to proceed
	ok, kerr := s.procToSysCall().RunBlock(key)
	if kerr != nil && kerr.IsError() {
		sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "RUN", "%s cannot run, error %s and ok %v, aborting...", s.proc, kerr, ok)
		return
	}
	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "RUN", "process %s read from the run channel %v", s.proc, ok)
	if !ok {
		sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "RUN", "we are now ready to run, but have been told to abort by nameserver, %s", s.proc)
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.RunPayload{}.KernelErrorPtr),
			lib.NewKernelError(lib.KernelDependencyCycle))
	}
	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "RUN", "we are now ready to run, %s", s.proc)
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.RunPayload{}.KernelErrorPtr),
		lib.NoError[*protosupport.KernelErrorId]())
}

func sysPrint(level pblog.LogLevel, call, spec string, arg ...interface{}) {
	if syscallVerbose {
		spec = "%s:" + spec
		arg = append([]interface{}{call}, arg...)
		msg := fmt.Sprintf(spec, arg...)
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		req := &pblog.LogRequest{
			Stamp:   timestamppb.New(time.Now()),
			Level:   level,
			Message: msg,
		}
		ilog.ProcessLogRequest(req, true, false, nil)

	}
}

// BackdoorLog is only for use by internal parigot debugging messages.
func (s *syscallReadWrite) BackdoorLog(sp int32) {
	req := &pblog.LogRequest{}
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
	ilog.ProcessLogRequest(req, true, false, nil)
	s.mem.SetInt32(int32(wasmPtr)+int32(unsafe.Offsetof(splitutil.SinglePayload{}.OutPtr)), 0)
	s.mem.SetInt32(int32(wasmPtr)+int32(unsafe.Offsetof(splitutil.SinglePayload{}.OutLen)), 0)

}
