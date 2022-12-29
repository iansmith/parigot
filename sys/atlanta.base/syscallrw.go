package sys

import (
	"fmt"
	"log"
	"os"
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
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Flip this switch for debug output that starts with SYSCALL and looks like this:
// SYSCALL[DISPATCH,mem-7f0524000000,[proc-9:storeclient.p.wasm]]:about to FindByName: [s-000008],BestOfAllTime,true
// The output is prefixed in the [] with the name of the syscall, the memory pointer for that
// process, and the process name/number in its own [].
//
// Look at the doc for libparigotVerbose to see about interleaving issue with syscallVerbose and
// libparigotVerbose.
var syscallVerbose = true

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
	ctx := s.procToSysCall().FindMethodByName(key, sid, req.Method)
	ctx.param, err = proto.Marshal(req.Param)
	if err != nil {
		splitutil.ErrorResponse(s.mem, sp, lib.KernelUnmarshalFailed)
		return
	}

	sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "Dispatch ", "find method by name requested for %s.%s",
		sid.Short(), req.Method)

	if ctx.pctx == nil {
		sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "Dispatch ", "xxx call context does not have pctx set but size of param is %d, any passed in is '%s'",
			len(ctx.param), req.Param.TypeUrl)
	}

	// this call is the machinery for making a call to another service
	resultInfo, errid := s.procToSysCall().CallService(ctx.target, ctx)
	sysPrint(pblog.LogLevel_LOG_LEVEL_DEBUG, "Dispatch ", "result info returned from call service: %#v", resultInfo)
	if errid != nil {
		if errid.IsErrorType() && errid.IsError() {
			kernErr := lib.KernelErrorCode(errid.Low())
			splitutil.ErrorResponse(s.mem, sp, kernErr)
			return
		} else {
			panic("dispatch is unable to understand result error of type:" + fmt.Sprintf("%T", errid))
		}
	}
	if resultInfo.errorId != nil {
		if resultInfo.errorId.IsErrorType() {
			if resultInfo.errorId.IsError() {
				kernErr := lib.KernelErrorCode(resultInfo.errorId.Low())
				splitutil.ErrorResponse(s.mem, sp, kernErr)
				return
			}
			// IsError is false when the error value is 0 (no error)
		} else {
			panic("dispatch is unable to understand result of inband error type in resultInfo:" + fmt.Sprintf("%T, %s", resultInfo.errorId, resultInfo.errorId.Short()))
		}
	}
	if len(resultInfo.pctx) > 0 {
		resp.OutPctx = &protosupport.Pctx{}
		if err := proto.Unmarshal(resultInfo.pctx, resp.OutPctx); err != nil {
			splitutil.ErrorResponse(s.mem, sp, lib.KernelUnmarshalFailed)
			return
		}
	}
	if len(resultInfo.result) > 0 {
		var a anypb.Any
		if err := proto.Unmarshal(resultInfo.result, &a); err != nil {
			splitutil.ErrorResponse(s.mem, sp, lib.KernelUnmarshalFailed)
			return

		}
		resp.Result = &a
	}
	resp.MethodId = lib.Marshal[protosupport.MethodId](resultInfo.mid)
	resp.CallId = lib.Marshal[protosupport.CallId](resultInfo.cid)
	splitutil.RespondSingleProto(s.mem, sp, &resp)
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

// BlockUntilCall is used by servers to block themselves until some other process sends them a
// message.
func (s *syscallReadWrite) BlockUntilCall(sp int32) {
	// tricky, we have to pass a process here?
	call := s.procToSysCall().BlockUntilCall(NewDepKeyFromProcess(s.proc))

	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "BLOCKUNTILCALL ", "received a call: %s,%s", call.cid.Short(), call.method)

	wasmPtr := s.mem.GetInt64(sp + 8)
	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "BLOCKUNTILCALL ", "wasmptr %x,true=%x", wasmPtr, s.mem.TrueAddr(int32(wasmPtr)))

	// check that we can fit the values
	availablePctxLen := s.ReadInt64(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.PctxLen))
	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "BLOCKUNTILCALL ", "size of available buffer for pctx: %d, need %d", availablePctxLen,
		len(call.pctx))
	if availablePctxLen > 0 && int64(len(call.pctx)) > availablePctxLen {
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.ErrorPtr),
			lib.NewKernelError(lib.KernelDataTooLarge))
		return
	}
	availableParamLen := s.ReadInt64(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.ParamLen))
	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "BLOCKUNTILCALL ", "size of available buffer for param: %d, need %d", availableParamLen,
		len(call.param))
	if int64(len(call.param)) > availableParamLen {
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.ErrorPtr),
			lib.NewKernelError(lib.KernelDataTooLarge))
		return
	}
	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "BLOCKUNTILCALL ", "Block until call checked the sizes and they are ok")

	// write the sizes of the incoming values and if size >0 copy the data to the pointer given
	if availablePctxLen == 0 {
		sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "BLOCKUNTILCALL ", "ignoring pctx in this call because callee says can't accept it")
		s.mem.SetInt64(int32(wasmPtr+int64(unsafe.Offsetof(lib.BlockPayload{}.PctxLen))), int64(0))
	} else {
		// we want to send the PCTX value and they said ok
		s.mem.SetInt64(int32(wasmPtr+int64(unsafe.Offsetof(lib.BlockPayload{}.PctxLen))), int64(len(call.pctx)))
		if len(call.pctx) > 0 {
			s.CopyToPtr(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.PctxPtr), call.pctx)
		} else {
			// this is the reverse of the previous, this is because the caller sent no PCTX
			sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "BLOCKUNTILCALL ", "skipping pctx, size is zero")

		}
	}
	s.mem.SetInt64(int32(wasmPtr+int64(unsafe.Offsetof(lib.BlockPayload{}.ParamLen))), int64(len(call.param)))
	if len(call.param) > 0 {
		s.CopyToPtr(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.ParamPtr), call.param)
	} else {
		sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "BLOCKUNTILCALL", "skipping param, size is zero")
	}

	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "BLOCKUNTILCALL", "Block until copied the data, and it is ok")

	// xxx fixme
	//direction := s.ReadInt64(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.Direction))

	// xxx fixme, should we be aborting here? Returning an error?
	// xxx also, it's not clear that this should not be just looked up rather than included in the payload
	// switch kernel.MethodDirection(direction) {
	// case kernel.MethodDirection_MethodDirectionIn, kernel.MethodDirection_MethodDirectionBoth:
	// 	if len(call.param) == 0 {
	// 		log.Printf("Warning: sent a zero size parametr to a function expecting input parameters")
	// 	}
	// }

	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.MethodId), call.mid)
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.CallId), call.cid)
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.ErrorPtr),
		lib.NoError[*protosupport.KernelErrorId]())

	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "BLOCKUNTILCALL", "Block until finished OK")
	return // the server goes back to work
}

// Return value is used by servers to register the return value for a particular function
// call on a method they implement.
func (s *syscallReadWrite) ReturnValue(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)
	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "RETURNVALUE", "wasmptr %x,true=%x", wasmPtr, s.mem.TrueAddr(int32(wasmPtr)))

	// we just have to create the structure and send it through the correct channel
	info := &resultInfo{}
	low, high := s.Read64BitPair(wasmPtr, unsafe.Offsetof(lib.ReturnValuePayload{}.CallId))
	info.cid = lib.NewFrom64BitPair[*protosupport.CallId](uint64(high), uint64(low))
	low, high = s.Read64BitPair(wasmPtr, unsafe.Offsetof(lib.ReturnValuePayload{}.MethodId))
	info.mid = lib.NewFrom64BitPair[*protosupport.MethodId](uint64(high), uint64(low))
	low, _ = s.Read64BitPair(wasmPtr, unsafe.Offsetof(lib.ReturnValuePayload{}.KernelErrorPtr))
	info.errorId = lib.NewKernelError(lib.KernelErrorCode(low))
	if !info.errorId.Equal(lib.NewError[*protosupport.KernelErrorId](uint16(lib.KernelNoError))) {
		short := info.errorId.Short()
		sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "RETURNVALUE", "Error found in the call? %s", short)
	}
	// if pctx len is 0 this is a no op
	info.pctx = s.ReadSlice(wasmPtr, unsafe.Offsetof(lib.ReturnValuePayload{}.PctxPtr),
		unsafe.Offsetof(lib.ReturnValuePayload{}.PctxLen))
	info.result = s.ReadSlice(wasmPtr, unsafe.Offsetof(lib.ReturnValuePayload{}.ResultPtr),
		unsafe.Offsetof(lib.ReturnValuePayload{}.ResultLen))

	callInfo := s.procToSysCall().GetInfoForCallId(info.cid)
	if callInfo == nil {
		sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "RETURNVALUE", "unable to find process/addr that called %s", info.cid.Short())
		kerr := lib.NewKernelError(lib.KernelCallerUnavailable)
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.ReturnValuePayload{}.KernelErrorPtr),
			kerr)
		return
	}
	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "RESULTVALUE ", "computed info, found channel, sending results of sizes: %d,%d for result and pctx data",
		len(info.result), len(info.pctx))

	callerProc := callInfo.sender.(*DepKeyImpl).proc

	if callerProc == nil {
		sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "RETURNVALUE ", "no caller proc, caller addr is %s", callInfo.sender.(*DepKeyImpl).addr)
	}
	callInfo.respCh <- info

	noerr := lib.NoError[*protosupport.KernelErrorId]()
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.ReturnValuePayload{}.KernelErrorPtr), noerr)

	sysPrint(pblog.LogLevel_LOG_LEVEL_INFO, "RESULTVALUE ", "finished with no error: %s", noerr.Short())
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
		req := &pblog.LogRequest{
			Stamp:   timestamppb.New(time.Now()),
			Level:   level,
			Message: fmt.Sprintf(spec, arg...),
		}
		ilog.ProcessLogRequest(req, true, false, nil)

	}
}
