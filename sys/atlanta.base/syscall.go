package sys

import (
	"log"
	"os"
	"unsafe"

	pblog "github.com/iansmith/parigot/g/pb/log"
	"github.com/iansmith/parigot/g/pb/parigot"
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/jspatch"

	"demo/vvv/proto/g/vvv/pb"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SysCall struct {
	mem        *jspatch.WasmMem
	nameServer *nameServer
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

func NewSysCall(ns *nameServer) *SysCall {
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

	params := s.ReadSlice(wasmPtr,
		unsafe.Offsetof(lib.DispatchPayload{}.ParamPtr),
		unsafe.Offsetof(lib.DispatchPayload{}.ParamLen))
	log.Printf("read the two slices: %d,%d", len(pctx), len(params))

	// xxx fix me
	// this is where we should be doing the call
	// xxx fix me
	log.Printf("aaa reached the call in dispatch[%d,%d]", len(pctx), len(params))
	methodId, proc := s.nameServer.FindMethodByName(sid, method)
	if methodId == nil || proc == nil {
		s.sendKernelErrorFromDispatch(wasmPtr, lib.KernelNotFound)
	}

	fakeResult, err := anypb.New(&pb.RevenueResponse{})
	if err != nil {
		s.sendKernelErrorFromDispatch(wasmPtr, lib.KernelDispatchTooLarge)
		return
	}
	fakePctx := &parigot.PCtx{
		Event: []*parigot.PCtxEvent{
			&parigot.PCtxEvent{
				Line: []*parigot.PCtxMessage{
					&parigot.PCtxMessage{
						Stamp:   timestamppb.Now(),
						Level:   pblog.LogLevel_LOGLEVEL_DEBUG,
						Message: "faked from inside the kernel",
					},
				},
			},
		},
	}

	// we have to check BOTH of the length values we were given to make
	// sure our results will fit
	resultLen := s.mem.GetInt64(int32(wasmPtr) +
		int32(unsafe.Offsetof(lib.DispatchPayload{}.ResultLen)))

	// we can't fit the result, so we signal error and abort
	if int64(proto.Size(fakeResult)) > resultLen {
		s.sendKernelErrorFromDispatch(wasmPtr, lib.KernelDispatchTooLarge)
		return
	}
	pctxLen := s.ReadInt64(wasmPtr,
		unsafe.Offsetof(lib.DispatchPayload{}.OutPctxLen))

	// we can't fit the pctx, so we signal error and abort
	if int64(proto.Size(fakePctx)) > pctxLen {
		s.sendKernelErrorFromDispatch(wasmPtr, lib.KernelDispatchTooLarge)
	}

	// tell the caller how big result is
	resultLen = int64(proto.Size(fakeResult))
	s.WriteInt64(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.ResultLen),
		int64(resultLen))

	// tell the caller how big the pctx is
	pctxLen = int64(proto.Size(fakePctx))
	s.WriteInt64(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.OutPctxLen),
		int64(pctxLen))

	// get the pctx bytes
	buf, err := proto.Marshal(fakePctx)
	if err != nil {
		s.sendKernelErrorFromDispatch(wasmPtr, lib.KernelMarshalFailed)
		return
	}
	s.CopyToPtr(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.OutPctxPtr), buf)

	// write the pctx
	buf, err = proto.Marshal(fakeResult)
	if err != nil {
		s.sendKernelErrorFromDispatch(wasmPtr, lib.KernelMarshalFailed)
		log.Printf("Dispatch: failed to marshal a result for which we had enough space: %v", err)
		return
	}

	s.CopyToPtr(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.ResultPtr), buf)

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
	log.Printf("BindMethod reached-------")
	wasmPtr := s.mem.GetInt64(sp + 8)

	log.Printf("BindMethod reached-------%x", wasmPtr)
	pkg := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.BindPayload{}.PkgPtr),
		unsafe.Offsetof(lib.BindPayload{}.PkgLen))
	log.Printf("BindMethod reached------%s", pkg)

	service := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.BindPayload{}.ServicePtr),
		unsafe.Offsetof(lib.BindPayload{}.ServiceLen))

	method := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.BindPayload{}.MethodPtr),
		unsafe.Offsetof(lib.BindPayload{}.MethodLen))

	log.Printf("about to hit the name service: %s,%s", service, method)
	mid, err := s.nameServer.HandleMethod(pkg, service, method, s.proc)
	if err != nil {
		s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.ErrorPtr),
			err)
		return

	}

	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.MethodId), mid)
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.ErrorPtr),
		lib.NoKernelErr())

	log.Printf("bind completed")
}

// BlockUntilCall is used by servers to block themselves until the kernel sends them a
// message.
func (s *SysCall) BlockUntilCall(sp int32) {
}

// Return value is used by servers to register the return value for a particular function
// call on a method they implement.
func (s *SysCall) ReturnValue(sp int32) {
}
