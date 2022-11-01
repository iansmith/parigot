package sys

import (
	"log"
	"os"
	"unsafe"

	"github.com/iansmith/parigot/g/pb/parigot"
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/jspatch"

	"demo/vvv/proto/g/vvv/pb"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const MaxService = 127

type SysCall struct {
	mem *jspatch.WasmMem
}

func (s *SysCall) SetMemPtr(m uintptr) {
	s.mem = jspatch.NewWasmMem(m)
}

func NewSysCall() *SysCall {
	return &SysCall{}
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

var packageRegistry = make(map[string]map[string]lib.Id)
var serviceCounter = 7 // last USED service num

func (s *SysCall) ReadString(structPtr int64, dataOffset uintptr, lenOffset uintptr) string {
	return s.mem.LoadStringWithLen(int32(structPtr)+int32(dataOffset), int32(structPtr)+int32(lenOffset))
}
func (s *SysCall) ReadInt64(structPtr int64, dataOffset uintptr) int64 {
	return s.mem.GetInt64(int32(structPtr) + int32(dataOffset))
}
func (s *SysCall) WriteInt64(structPtr int64, dataOffset uintptr, value int64) {
	s.mem.SetInt64(int32(structPtr)+int32(dataOffset), value)
}

func (s *SysCall) Write64BitPair(structPtr int64, dataOffset uintptr, id lib.Id) {
	derefed := s.mem.GetInt32(int32(structPtr + int64(dataOffset)))
	// write the error info back to client
	s.mem.SetInt64(derefed, int64(id.Low()))
	s.mem.SetInt64(derefed+8, int64(id.High()))
}
func (s *SysCall) Read64BitPair(structPtr int64, dataOffset uintptr) (int64, int64) {
	low := s.mem.GetInt64(int32(structPtr + int64(dataOffset)))
	high := s.mem.GetInt64(int32(structPtr + int64(dataOffset) + 8))
	return low, high
}
func (s *SysCall) ReadSlice(structPtr int64, dataOffset uintptr, lenOffset uintptr) []byte {
	return s.mem.LoadSliceWithLenAddr(int32(structPtr)+int32(dataOffset),
		int32(structPtr)+int32(lenOffset))
}
func (s *SysCall) CopyToPtr(structPtr int64, dataOffset uintptr, content []byte) {
	s.mem.CopyToPtr(int32(structPtr)+int32(dataOffset), content)
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

	//errOff := unsafe.Offsetof(lib.RegDetail{}.OutErrPtr)
	//errInd := s.mem.GetInt32(int32(wasmPtr + int64(errOff)))
	//
	//offSvcPtr := unsafe.Offsetof(lib.RegDetail{}.OutServiceIdPtr)
	//svcInd := s.mem.GetInt32(int32(wasmPtr + int64(offSvcPtr)))

	serviceRegistry, ok := packageRegistry[pkg]
	if !ok {
		serviceRegistry = make(map[string]lib.Id)
		packageRegistry[pkg] = serviceRegistry
	}
	_, duplicate := serviceRegistry[service]

	var regErr lib.Id
	switch {
	case serviceCounter > MaxService:
		regErr = lib.NewRegisterErr(lib.RegisterNamespaceExhausted)
	case duplicate:
		regErr = lib.NewRegisterErr(lib.RegisterAlreadyRegistered)
	default:
		regErr = lib.NoRegisterErr()
	}
	// assign sid if no error
	if !regErr.IsError() {
		sid := lib.ServiceIdFromUint64(0, uint64(serviceCounter+1))
		serviceRegistry[service] = sid
		serviceCounter++
		//send back the data to client
		s.Write64BitPair(wasmPtr,
			unsafe.Offsetof(lib.RegPayload{}.OutServiceIdPtr), sid)
		//s.mem.SetInt64(svcInd, int64(sid.Low()))
		//s.mem.SetInt64(svcInd+8, int64(sid.High()))
	}
	s.Write64BitPair(wasmPtr,
		unsafe.Offsetof(lib.RegPayload{}.OutErrPtr), regErr)

	// write the error info back to client
	//s.mem.SetInt64(errInd, int64(regErr.Low()))
	//s.mem.SetInt64(errInd+8, int64(regErr.High()))

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

	var locErr *lib.Id
	var serviceId *lib.Id
	serviceRegistry, ok := packageRegistry[pkg]
	if !ok {
		l := lib.NewLocateErr(lib.LocateNotFound)
		locErr = &l
	} else {
		id, ok := serviceRegistry[service]
		if !ok {
			l := lib.NewLocateErr(lib.LocateNotFound)
			locErr = &l
		} else {
			serviceId = &id
		}
	}
	// no we need to write the results back to the client
	switch {
	case locErr != nil:
		s.Write64BitPair(wasmPtr,
			unsafe.Offsetof(lib.LocatePayload{}.OutErrPtr), *locErr)
	case serviceId != nil:
		l := lib.NoLocateErr()
		locErr = &l
		s.Write64BitPair(wasmPtr,
			unsafe.Offsetof(lib.LocatePayload{}.OutErrPtr), *locErr)
		s.Write64BitPair(wasmPtr,
			unsafe.Offsetof(lib.LocatePayload{}.OutServiceIdPtr), *serviceId)
	default:
		panic("did not create an error or service id when performing locate")
	}
	return
}

func DebugPrint(ct int32) {
	log.Printf("----DebugPrint %d", ct)
}

func (s *SysCall) Dispatch(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)
	low, high := s.Read64BitPair(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.ServiceId))

	_ = lib.ServiceIdFromUint64(uint64(high), uint64(low))

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
	fakeResult, err := anypb.New(&pb.RevenueResponse{})
	if err != nil {
		s.sendDispatchError(wasmPtr, lib.DispatchTooLarge)
		return
	}
	fakePctx := &parigot.PCtx{
		Event: []*parigot.PCtxEvent{
			&parigot.PCtxEvent{
				Line: []*parigot.PCtxMessage{
					&parigot.PCtxMessage{
						Stamp:   timestamppb.Now(),
						Level:   parigot.LogLevel_LOGLEVEL_DEBUG,
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
		s.sendDispatchError(wasmPtr, lib.DispatchTooLarge)
		return
	}
	pctxLen := s.ReadInt64(wasmPtr,
		unsafe.Offsetof(lib.DispatchPayload{}.OutPctxLen))

	// we can't fit the pctx, so we signal error and abort
	if int64(proto.Size(fakePctx)) > pctxLen {
		s.sendDispatchError(wasmPtr, lib.DispatchTooLarge)
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
		s.sendDispatchError(wasmPtr, lib.DispatchMarshalFailed)
		return
	}
	s.CopyToPtr(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.OutPctxPtr), buf)

	// write the pctx
	buf, err = proto.Marshal(fakeResult)
	if err != nil {
		s.sendDispatchError(wasmPtr, lib.DispatchMarshalFailed)
		log.Printf("Dispatch: failed to marshal a result for which we had enough space: %v", err)
		return
	}

	s.CopyToPtr(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.ResultPtr), buf)

	dispErr := lib.NoDispatchErr() // the lack of an error

	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.ErrorPtr),
		dispErr)

}

func (s *SysCall) sendDispatchError(wasmPtr int64, code lib.DispatchErrCode) {
	dispErr := lib.NewDispatchErr(code)
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.DispatchPayload{}.ErrorPtr),
		dispErr)
	return

}
