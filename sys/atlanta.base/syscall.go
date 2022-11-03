package sys

import (
	"fmt"
	"log"
	"os"
	"unsafe"

	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/jspatch"
)

type SysCall struct {
	mem        *jspatch.WasmMem
	nameServer *NameServer
	proc       *Process // this is US
}

// SetMemPtr has to be separated out because at the time this object is created, we don't yet
// know the memory address that is the memPtr.
func (s *SysCall) SetMemPtr(m uintptr) {
	s.mem = jspatch.NewWasmMem(m)
}

// SetProcess has to be separated out because at the time this object is created, we don't yet
// know the process that this syscall works on behalf of.  This is because the process needs a
// syscall (this one) to be created.
func (s *SysCall) SetProcess(p *Process) {
	s.proc = p
}

func NewSysCall(ns *NameServer) *SysCall {
	return &SysCall{
		nameServer: ns,
	}
}

func (a *SysCall) OutputString(s string) {
	print(s)
}

func (a *SysCall) JSNotImplemented(s string) {
	print(s)
	os.Exit(1)
}

func (a *SysCall) JSHandleEvent(e int32) {
	print("JSHandleEvent\n")
	os.Exit(1)
}

func (a *SysCall) Exit(sp int32) {
	log.Printf("exit called")
	os.Exit(int(0))
}

func (a *SysCall) Now(retVal int32) {
	print("Now")
	os.Exit(1)
}

func (a *SysCall) SetNow(_ int64, _ bool) {
	print("SetNow\n")
	os.Exit(1)
}
func (s *SysCall) Register(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)

	pkg := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.RegPayload{}.PkgPtr),
		unsafe.Offsetof(lib.RegPayload{}.PkgLen))

	service := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.RegPayload{}.ServicePtr),
		unsafe.Offsetof(lib.RegPayload{}.ServiceLen))

	log.Printf("registration for %s.%s", pkg, service)

	sid, err := s.nameServer.RegisterClientService(pkg, service)

	// did the ns complain?
	if err != nil {
		s.Write64BitPair(wasmPtr,
			unsafe.Offsetof(lib.RegPayload{}.ErrorPtr), err)
		return
	}
	// tell the client there is no error
	s.Write64BitPair(wasmPtr,
		unsafe.Offsetof(lib.RegPayload{}.ErrorPtr), lib.NoKernelErr())
	s.Write64BitPair(wasmPtr,
		unsafe.Offsetof(lib.RegPayload{}.ServiceIdPtr), sid)

}

func (s *SysCall) Locate(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)
	pkg := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.LocatePayload{}.PkgPtr),
		unsafe.Offsetof(lib.LocatePayload{}.PkgLen))

	service := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.LocatePayload{}.ServicePtr),
		unsafe.Offsetof(lib.LocatePayload{}.ServiceLen))

	log.Printf("locate requested for %s.%s", pkg, service)

	sid, err := s.nameServer.GetService(pkg, service)

	if err != nil {
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

func (s *SysCall) Dispatch(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)
	low, high := s.Read64BitPair(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.ServiceId))

	sid := lib.ServiceIdFromUint64(uint64(high), uint64(low))

	method := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.DispatchPayload{}.MethodPtr),
		unsafe.Offsetof(lib.DispatchPayload{}.MethodLen))

	caller := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.DispatchPayload{}.CallerPtr),
		unsafe.Offsetof(lib.DispatchPayload{}.CallerLen))
	log.Printf("read the two strings: %s,%s", method, caller)
	pctx := s.ReadSlice(wasmPtr,
		unsafe.Offsetof(lib.DispatchPayload{}.PctxPtr),
		unsafe.Offsetof(lib.DispatchPayload{}.PctxLen))

	param := s.ReadSlice(wasmPtr,
		unsafe.Offsetof(lib.DispatchPayload{}.ParamPtr),
		unsafe.Offsetof(lib.DispatchPayload{}.ParamLen))
	log.Printf("read the two slices: %d,%d", len(pctx), len(param))

	// we ask the nameserver to find us the appropriate process and method so we can call
	// the other side... the nameserver also assigns us a call id
	callCtx := s.nameServer.FindMethodByName(sid, method, s.proc)
	print(fmt.Sprintf("xxx FindMethodByName done: %+v\n", callCtx))

	if callCtx == nil {
		s.sendKernelErrorFromDispatch(wasmPtr, lib.KernelNotFound)
		return
	}

	// send the message...
	callInfo := &callInfo{
		mid:    callCtx.mid,
		cid:    callCtx.cid,
		caller: callCtx.sender, // also s.proc
		param:  param,
		pctx:   pctx,
	}

	log.Printf("sending the call info to other process: %+v", callInfo)
	// the magic: send the call value to the other process
	callCtx.target.callCh <- callInfo

	log.Printf("waiting for result from other process: %+v", s.proc)
	// wait for the other process to message us back with a result... note that we should be
	// timing this out after some period of time.  xxx fixme  This situation occurs specifically
	// if the callee (the server implementation) cannot receive the data because is it too large.
	resultInfo := <-s.proc.resultCh

	log.Printf("got result from other process (YAY!): %+v", resultInfo)

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
	}

	noErr := lib.NoKernelErr() // the lack of an error
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.ErrorPtr),
		noErr)

}

func (s *SysCall) sendKernelErrorFromDispatch(wasmPtr int64, code lib.KernelErrorCode) {
	dispErr := lib.NewKernelError(code)
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.ErrorPtr),
		dispErr)
	return
}

func (s *SysCall) sendKernelErrorFromBind(wasmPtr int64, code lib.KernelErrorCode) {
	dispErr := lib.NewKernelError(code)
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.ErrorPtr),
		dispErr)
	return
}

// BindMethod is used to indicate the function that will handle a given method.  We don't
// actually have a handle to the function pointer, we give out MethodIds instead.
func (s *SysCall) BindMethod(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)

	pkg := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.BindPayload{}.PkgPtr),
		unsafe.Offsetof(lib.BindPayload{}.PkgLen))

	service := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.BindPayload{}.ServicePtr),
		unsafe.Offsetof(lib.BindPayload{}.ServiceLen))

	method := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.BindPayload{}.MethodPtr),
		unsafe.Offsetof(lib.BindPayload{}.MethodLen))

	log.Printf("len of service: %d, len of method %d", len(service), len(method))
	log.Printf("about to hit the name service in Bind: %s,%s", service, method)
	mid, err := s.nameServer.HandleMethod(pkg, service, method, s.proc)
	if err != nil {
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.ErrorPtr),
			err)
		return
	}

	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.MethodId), mid)
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.ErrorPtr),
		lib.NoKernelErr())

	log.Printf("bind completed, %s bound to %s", method, mid.Short())
}

// BlockUntilCall is used by servers to block themselves until some other process sends them a
// message.
func (s *SysCall) BlockUntilCall(sp int32) {
	call := <-s.proc.callCh
	log.Printf("Block until call received a call: %+v", call)

	wasmPtr := s.mem.GetInt64(sp + 8)

	// check that we can fit the values
	availablePctxLen := s.ReadInt64(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.PctxLen))
	if int64(len(call.pctx)) > availablePctxLen {
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.ErrorPtr),
			lib.NewKernelError(lib.KernelDataTooLarge))
		return
	}
	availableParamLen := s.ReadInt64(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.ParamLen))
	if int64(len(call.param)) > availableParamLen {
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.ErrorPtr),
			lib.NewKernelError(lib.KernelDataTooLarge))
		return
	}
	log.Printf("Block until call checked the sizes and they are ok")

	// write the sizes of the incoming values and if size >0 copy the data to the pointer given
	s.mem.SetInt64(int32(wasmPtr+int64(unsafe.Offsetof(lib.BlockPayload{}.PctxLen))), int64(len(call.pctx)))
	if len(call.pctx) > 0 {
		s.CopyToPtr(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.PctxPtr), call.pctx)
	}
	s.mem.SetInt64(int32(wasmPtr+int64(unsafe.Offsetof(lib.BlockPayload{}.ParamLen))), int64(len(call.param)))
	if len(call.pctx) > 0 {
		s.CopyToPtr(wasmPtr, unsafe.Offsetof(lib.BlockPayload{}.ParamPtr), call.param)
	}

	log.Printf("Block until copied the data, and it is ok")

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

	log.Printf("Block until finished")
	return // the server goes back to work
}

// Return value is used by servers to register the return value for a particular function
// call on a method they implement.
func (s *SysCall) ReturnValue(sp int32) {
}
