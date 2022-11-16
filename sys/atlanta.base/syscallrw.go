package sys

import (
	"fmt"
	"log"
	"os"
	"unsafe"

	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/jspatch"
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
	proc *Process // this is US

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

// func (a *syscallReadWrite) OutputString(s string) {
// 	print(s)
// }

// func (a *syscallReadWrite) JSNotImplemented(s string) {
// 	print(s)
// 	os.Exit(1)
// }

// func (a *syscallReadWrite) JSHandleEvent(e int32) {
// 	print("JSHandleEvent\n")
// 	os.Exit(1)
// }

func (a *syscallReadWrite) Exit(sp int32) {
	log.Printf("exit called")
	a.proc.exited = true
	// xxx fixme how can we cause a return here instead of using exit()?
	os.Exit(int(0))
}

// func (a *syscallReadWrite) Now(retVal int32) {
// 	print("Now")
// 	os.Exit(1)
// }

// func (a *syscallReadWrite) SetNow(_ int64, _ bool) {
// 	print("SetNow\n")
// 	os.Exit(1)
// }

func (s *syscallReadWrite) Locate(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)
	sysPrint("Locate", "wasmptr %x,true=%x", wasmPtr, s.mem.TrueAddr(int32(wasmPtr)))
	pkg := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.LocatePayload{}.PkgPtr),
		unsafe.Offsetof(lib.LocatePayload{}.PkgLen))

	service := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.LocatePayload{}.ServicePtr),
		unsafe.Offsetof(lib.LocatePayload{}.ServiceLen))

	sysPrint("LOCATE ", "locate requested for %s.%s", pkg, service)

	sid, err := s.procToSysCall().GetService(NewDepKeyFromProcess(s.proc), pkg, service)

	if err != nil && err.IsError() {
		s.Write64BitPair(wasmPtr,
			unsafe.Offsetof(lib.LocatePayload{}.ErrorPtr), err)
		return
	}
	s.Write64BitPair(wasmPtr,
		unsafe.Offsetof(lib.LocatePayload{}.ErrorPtr), lib.NoKernelErr())

	s.Write64BitPair(wasmPtr,
		unsafe.Offsetof(lib.LocatePayload{}.ServiceIdPtr), sid)
}

func DebugPrint(ct int32) {
	log.Printf("----DebugPrint %d", ct)
}

func (s *syscallReadWrite) Dispatch(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)
	sysPrint("Dispatch", "wasmptr %x,true=%x", wasmPtr, s.mem.TrueAddr(int32(wasmPtr)))

	low, high := s.Read64BitPair(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.ServiceId))

	sid := lib.ServiceIdFromUint64(uint64(high), uint64(low))

	method := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.DispatchPayload{}.MethodPtr),
		unsafe.Offsetof(lib.DispatchPayload{}.MethodLen))

	caller := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.DispatchPayload{}.CallerPtr),
		unsafe.Offsetof(lib.DispatchPayload{}.CallerLen))

	pctx := s.ReadSlice(wasmPtr,
		unsafe.Offsetof(lib.DispatchPayload{}.PctxPtr),
		unsafe.Offsetof(lib.DispatchPayload{}.PctxLen))

	param := s.ReadSlice(wasmPtr,
		unsafe.Offsetof(lib.DispatchPayload{}.ParamPtr),
		unsafe.Offsetof(lib.DispatchPayload{}.ParamLen))

	sysPrint("DISPATCH", "about to FindByName: %s,%s", sid.Short(), method)
	// this call sets up the call context and it's tricky because the key abstraction is used
	// for both the source and des, but one is local (Process) and the other is remote (addr)
	source := NewDepKeyFromProcess(s.proc)
	callCtx := s.procToSysCall().FindMethodByName(source, sid, method)

	if callCtx == nil {
		sysPrint("DISPATCH", "FindMethodByName failed for %s,%s", sid.Short(), method)
		s.sendKernelErrorFromDispatch(wasmPtr, lib.KernelNotFound)
		return
	}
	sysPrint("DISPATCH", "FindMethodByName done and OK: %s from '%s'",
		callCtx.cid.Short(), caller)

	destParam := make([]byte, len(param))
	destPctx := make([]byte, len(pctx))
	copy(destParam, param)
	if len(pctx) > 0 {
		copy(destPctx, pctx)
	} else {
		destPctx = []byte{}
	}
	// send the message...
	callInfo := &callInfo{
		mid:    callCtx.mid,
		cid:    callCtx.cid,
		sid:    callCtx.sid,
		method: callCtx.method,
		sender: callCtx.sender,
		param:  destParam,
		pctx:   destPctx,
	}

	// the magic: send the call value to the other process
	resultInfo, kerr := s.procToSysCall().CallService(callCtx.target, callInfo)
	if kerr != nil && kerr.IsError() {
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.ErrorPtr),
			kerr)
		return
	}
	if resultInfo == nil {
		sysPrint("DISPATCH", "waiting for result from other process: %s", s.proc.String())
		// wait for the other process to message us back with a result... note that we should be
		// timing this out after some period of time.  xxx fixme  This situation occurs specifically
		// if the callee (the server implementation) cannot receive the data because is it too large.
		resultInfo = <-s.proc.resultCh
	}

	sysPrint("DISPATCH", "got result from other process %s,%s with sizes pctx=%d,result=%d", resultInfo.cid.Short(), resultInfo.errorId.Short(),
		len(resultInfo.pctx), len(resultInfo.result))

	// we have to check BOTH of the length values we were given to make
	// sure our results will fit
	resultLen := s.mem.GetInt64(int32(wasmPtr) +
		int32(unsafe.Offsetof(lib.DispatchPayload{}.ResultLen)))

	// we can't fit the result, so we signal error and abort
	if len(resultInfo.result) > int(resultLen) {
		s.sendKernelErrorFromDispatch(wasmPtr, lib.KernelDataTooLarge)
		return
	}

	pctxLen := s.ReadInt64(wasmPtr,
		unsafe.Offsetof(lib.DispatchPayload{}.OutPctxLen))

	// we can't fit the pctx, so we signal error and abort
	if len(resultInfo.result) > int(pctxLen) {
		s.sendKernelErrorFromDispatch(wasmPtr, lib.KernelDataTooLarge)
		return
	}

	sysPrint("DISPATCH", "telling the  caller the size of the result and pctx [%d,%d]",
		len(resultInfo.result), len(resultInfo.pctx))

	// tell the caller how big result is
	resultLen = int64(len(resultInfo.result))
	s.WriteInt64(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.ResultLen),
		int64(resultLen))

	// tell the caller how big the pctx is
	pctxLen = int64(len(resultInfo.pctx))
	s.WriteInt64(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.OutPctxLen),
		int64(pctxLen))

	if pctxLen > 0 {
		s.CopyToPtr(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.OutPctxPtr), resultInfo.pctx)
	}

	if resultLen > 0 {
		s.CopyToPtr(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.ResultPtr), resultInfo.result)
		sysPrint("DISPATCH ", "copied %d bytes to original caller", len(resultInfo.result))
	}

	noErr := lib.NoKernelErr() // the lack of an error
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.ErrorPtr),
		noErr)

	sysPrint("DISPATCH ", "completed call %s", method)

}

func (s *syscallReadWrite) sendKernelErrorFromDispatch(wasmPtr int64, code lib.KernelErrorCode) {
	dispErr := lib.NewKernelError(code)
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.ErrorPtr),
		dispErr)
	return
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
	sysPrint("BINDMETHOD ", "wasmptr %x, true=%x", wasmPtr, s.mem.TrueAddr(int32(wasmPtr)))

	pkg := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.BindPayload{}.PkgPtr),
		unsafe.Offsetof(lib.BindPayload{}.PkgLen))

	service := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.BindPayload{}.ServicePtr),
		unsafe.Offsetof(lib.BindPayload{}.ServiceLen))

	method := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.BindPayload{}.MethodPtr),
		unsafe.Offsetof(lib.BindPayload{}.MethodLen))

	sysPrint("BINDMETHOD", "about to bind %s in service %s", method, service)
	mid, kerr := s.procToSysCall().Bind(s.proc, pkg, service, method)
	if kerr != nil {
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.ErrorPtr),
			kerr)
		return
	}

	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.MethodId), mid)
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.ErrorPtr),
		lib.NoKernelErr())

	sysPrint("BINDMETHOD", "bind completed, %s bound to %s", method, mid.Short())
}

// BlockUntilCall is used by servers to block themselves until some other process sends them a
// message.
func (s *syscallReadWrite) BlockUntilCall(sp int32) {
	call := <-s.proc.callCh

	sysPrint("BLOCKUNTILCALL ", "received a call: %s,%s", call.cid.Short(), call.mid.Short())

	wasmPtr := s.mem.GetInt64(sp + 8)
	sysPrint("BLOCKUNTILCALL ", "wasmptr %x,true=%x", wasmPtr, s.mem.TrueAddr(int32(wasmPtr)))

	// check that we can fit the values
	availablePctxLen := s.ReadInt64(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.PctxLen))
	sysPrint("BLOCKUNTILCALL ", "size of available buffer for pctx: %d, need %d", availablePctxLen,
		len(call.pctx))
	if availablePctxLen > 0 && int64(len(call.pctx)) > availablePctxLen {
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.ErrorPtr),
			lib.NewKernelError(lib.KernelDataTooLarge))
		return
	}
	availableParamLen := s.ReadInt64(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.ParamLen))
	sysPrint("BLOCKUNTILCALL ", "size of available buffer for param: %d, need %d", availableParamLen,
		len(call.param))
	if int64(len(call.param)) > availableParamLen {
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.ErrorPtr),
			lib.NewKernelError(lib.KernelDataTooLarge))
		return
	}
	sysPrint("BLOCKUNTILCALL ", "Block until call checked the sizes and they are ok")

	// write the sizes of the incoming values and if size >0 copy the data to the pointer given
	if availablePctxLen == 0 {
		sysPrint("BLOCKUNTILCALL ", "ignoring pctx in this call because callee says can't accept it")
		s.mem.SetInt64(int32(wasmPtr+int64(unsafe.Offsetof(lib.BlockPayload{}.PctxLen))), int64(0))
	} else {
		// we want to send the PCTX value and they said ok
		s.mem.SetInt64(int32(wasmPtr+int64(unsafe.Offsetof(lib.BlockPayload{}.PctxLen))), int64(len(call.pctx)))
		if len(call.pctx) > 0 {
			s.CopyToPtr(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.PctxPtr), call.pctx)
		} else {
			// this is the reverse of the previous, this is because the caller sent no PCTX
			sysPrint("BLOCKUNTILCALL ", "skipping pctx, size is zero")

		}
	}
	s.mem.SetInt64(int32(wasmPtr+int64(unsafe.Offsetof(lib.BlockPayload{}.ParamLen))), int64(len(call.param)))
	if len(call.param) > 0 {
		s.CopyToPtr(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.ParamPtr), call.param)
	} else {
		sysPrint("BLOCKUNTILCALL", "skipping param, size is zero")
	}

	sysPrint("BLOCKUNTILCALL", "Block until copied the data, and it is ok")

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
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.ErrorPtr), lib.NoKernelErr())

	sysPrint("BLOCKUNTILCALL", "Block until finished OK")
	return // the server goes back to work
}

// Return value is used by servers to register the return value for a particular function
// call on a method they implement.
func (s *syscallReadWrite) ReturnValue(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)
	sysPrint("RETURNVALUE", "wasmptr %x,true=%x", wasmPtr, s.mem.TrueAddr(int32(wasmPtr)))

	// we just have to create the structure and send it through the correct channel
	info := &resultInfo{}
	low, high := s.Read64BitPair(wasmPtr, unsafe.Offsetof(lib.ReturnValuePayload{}.CallId))
	info.cid = lib.CallIdFromUint64(uint64(high), uint64(low))
	low, high = s.Read64BitPair(wasmPtr, unsafe.Offsetof(lib.ReturnValuePayload{}.MethodId))
	info.mid = lib.MethodIdFromUint64(uint64(high), uint64(low))
	low, _ = s.Read64BitPair(wasmPtr, unsafe.Offsetof(lib.ReturnValuePayload{}.KernelErrorPtr))
	info.errorId = lib.NewKernelError(lib.KernelErrorCode(low))
	sysPrint("RETURNVALUE", "Errorfound in the call? %s", info.errorId)
	// if pctx len is 0 this is a no op
	info.pctx = s.ReadSlice(wasmPtr, unsafe.Offsetof(lib.ReturnValuePayload{}.PctxPtr),
		unsafe.Offsetof(lib.ReturnValuePayload{}.PctxLen))
	info.result = s.ReadSlice(wasmPtr, unsafe.Offsetof(lib.ReturnValuePayload{}.ResultPtr),
		unsafe.Offsetof(lib.ReturnValuePayload{}.ResultLen))

	sysPrint("RETURNVALUE", "searching for process assocated with the call of %s", info.cid.Short())
	key := s.procToSysCall().GetProcessForCallId(info.cid)
	if key == nil {
		sysPrint("RETURNVALUE", "unable to find process that called %s", info.cid.Short())
		kerr := lib.NewKernelError(lib.KernelCallerUnavailable)
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.ReturnValuePayload{}.KernelErrorPtr),
			kerr)
		return
	}
	sysPrint("RESULTVALUE ", "computed info, found channel, sending results: %d,%d for result and pctx data",
		len(info.result), len(info.pctx))

	key.(*DepKeyImpl).proc.resultCh <- info
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.ReturnValuePayload{}.KernelErrorPtr),
		lib.NoKernelErr())
	sysPrint("RESULTVALUE ", "finished")
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
// of that service.
func (s *syscallReadWrite) Export(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)
	sysPrint("EXPORT", "wasmptr %x,true=%x", wasmPtr, s.mem.TrueAddr(int32(wasmPtr)))

	pkg := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.ExportPayload{}.PkgPtr),
		unsafe.Offsetof(lib.ExportPayload{}.PkgLen))

	service := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.ExportPayload{}.ServicePtr),
		unsafe.Offsetof(lib.ExportPayload{}.ServiceLen))

	sysPrint("EXPORT", "about to tell the nameserver closeservice and  export %s.%s [syscall->%T]", pkg, service,
		s.procToSysCall())
	kerr := s.procToSysCall().Export(NewDepKeyFromProcess(s.proc), pkg, service)
	if kerr != nil && kerr.IsError() {
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.ExportPayload{}.KernelErrorPtr),
			kerr)
		return
	}
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.ExportPayload{}.KernelErrorPtr),
		lib.NoKernelErr())

	sysPrint("EXPORT", "done")

}

// Require is used when client or server wishes to indicate that it consumes
// a service.  This becomes part of the dependency graph.
func (s *syscallReadWrite) Require(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)
	sysPrint("REQUIRE", "wasmptr %x,true=%x", wasmPtr, s.mem.TrueAddr(int32(wasmPtr)))

	// we just have to create the structure and send it through the correct channel
	pkg := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.RequirePayload{}.PkgPtr),
		unsafe.Offsetof(lib.RequirePayload{}.PkgLen))

	service := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.RequirePayload{}.ServicePtr),
		unsafe.Offsetof(lib.RequirePayload{}.ServiceLen))

	sysPrint("REQUIRE", "telling nameserver %s requires %s.%s", s.proc, pkg, service)

	kerr := s.procToSysCall().Require(NewDepKeyFromProcess(s.proc), pkg, service)

	if kerr != nil && kerr.IsError() {
		sysPrint("REQUIRE", " nameserver failed require of %s.%s", pkg, service)
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.RequirePayload{}.KernelErrorPtr),
			kerr)
		return
	}
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.RequirePayload{}.KernelErrorPtr),
		lib.NoKernelErr())

}

// Run is used to start up the processes in a deterministic order. It will
// fail and return an error if there are problems getting all the require and export
// requests to match up.
func (s *syscallReadWrite) Run(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)
	sysPrint("RUN", "wasmptr %x,true=%x", wasmPtr, s.mem.TrueAddr(int32(wasmPtr)))
	w := s.mem.GetInt64(int32(wasmPtr) + int32(unsafe.Offsetof(lib.RunPayload{}.Wait)))
	wait := false
	if w != 0 {
		wait = true
	}
	s.proc.waiter = wait
	s.proc.reachedRun = true

	key := NewDepKeyFromProcess(s.proc)
	s.procToSysCall().RunNotify(key)

	sysPrint("RUN", "%s is blocked on channel for run confirmation", s.proc)
	// block until we are told to proceed
	ok, kerr := s.procToSysCall().RunBlock(key)
	if kerr != nil && kerr.IsError() {
		sysPrint("RUN", "%s cannot run, error %s and ok %v, aborting...", s.proc, kerr, ok)
		return
	}
	sysPrint("RUN", "process %s read from the run channel %v", s.proc, ok)
	if !ok {
		sysPrint("RUN", "we are now ready to run, but have been told to abort by nameserver, %s", s.proc)
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.RunPayload{}.KernelErrorPtr),
			lib.NewKernelError(lib.KernelDependencyCycle))
	}
	sysPrint("RUN", "we are now ready to run, %s", s.proc)
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.RunPayload{}.KernelErrorPtr),
		lib.NoKernelErr())
}

func sysPrint(call, spec string, arg ...interface{}) {
	if syscallVerbose {
		p1 := fmt.Sprintf("SYSCALL[%s]:", call)
		print(p1, fmt.Sprintf(spec, arg...), "\n")
	}
}
