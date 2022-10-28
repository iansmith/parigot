package sys

import (
	"log"
	"os"
	"unsafe"

	"github.com/iansmith/parigot/g/pb/kernel"
	"github.com/iansmith/parigot/g/pb/parigot"
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/jspatch"
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

func (a *SysCall) TinygoNotImplemented(s string) {
	print(s)
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

func (s *SysCall) getRegisterRequest(offset int32) *kernel.RegisterRequest {
	wasmPtr := s.mem.GetInt32(offset)
	offPkg := unsafe.Offsetof(kernel.RegisterRequest{}.ProtoPackage)
	pkg := s.mem.LoadString(wasmPtr + int32(offPkg))
	offSvc := unsafe.Offsetof(kernel.RegisterRequest{}.Service)
	svc := s.mem.LoadString(wasmPtr + int32(offSvc))
	return &kernel.RegisterRequest{
		ProtoPackage: pkg,
		Service:      svc,
	}
}

func (s *SysCall) fillRegisterResponse(offset int32,
	regErrRaw lib.Id,
	sidRaw lib.Id) {

	regErr := lib.MarshalRegisterErrId(regErrRaw)
	sid := lib.MarshalServiceId(sidRaw)

	wasmPtr := s.mem.GetInt32(offset)
	offErr := unsafe.Offsetof(kernel.RegisterResponse{}.ErrorId)
	errIdPtr := s.mem.GetInt32(wasmPtr + int32(offErr))
	lowErrOff := unsafe.Offsetof(parigot.RegisterErrorId{}.Low)
	errValuePtr := uintptr(errIdPtr + int32(lowErrOff))
	s.mem.SetInt64(int32(errValuePtr), int64(regErr.Low))
	s.mem.SetInt64(int32(errValuePtr-8), int64(regErr.High))

	offSid := unsafe.Offsetof(kernel.RegisterResponse{}.ServiceId)
	sidPtr := s.mem.GetInt32(wasmPtr + int32(offSid))
	lowSidOff := unsafe.Offsetof(kernel.RegisterResponse{}.ServiceId.Low)
	sidValuePtr := uintptr(sidPtr) + lowSidOff
	s.mem.SetInt64(int32(sidValuePtr), int64(sid.Low))
	s.mem.SetInt64(int32(sidValuePtr-8), int64(sid.High))
}

func (s *SysCall) Register(sp int32) {
	//rreq := s.getRegisterRequest(sp + 8)
	//serviceRegistry, ok := packageRegistry[rreq.ProtoPackage]
	//if !ok {
	//	serviceRegistry = make(map[string]lib.Id)
	//}
	//sid := lib.ServiceIdFromUint64(0, uint64(serviceCounter+1))
	//serviceCounter++
	//if serviceCounter > MaxService {
	//	panic("service count exceeded maximum of " + fmt.Sprint(MaxService))
	//}
	//serviceRegistry[rreq.Service] = sid
	//s.fillRegisterResponse(sp+16, lib.NewRegisterErr(0), sid)
	//pkg := s.mem.LoadStringWithLen(pkgPtr, pkgLen)
	//svc := s.mem.LoadStringWithLen(svcPtr, svcLen)
	//log.Printf("got register at ABI:%s,%s", pkg, svc)
	wasmPtr := s.mem.GetInt64(sp + 8)

	offPkg := unsafe.Offsetof(lib.RegDetail{}.PkgPtr)
	offLen := unsafe.Offsetof(lib.RegDetail{}.PkgLen)
	pkg := s.mem.LoadStringTwoPtrs(int32(wasmPtr)+int32(offPkg), int32(wasmPtr)+int32(offLen))

	offService := unsafe.Offsetof(lib.RegDetail{}.ServicePtr)
	serviceLen := unsafe.Offsetof(lib.RegDetail{}.ServiceLen)
	service := s.mem.LoadStringTwoPtrs(int32(wasmPtr)+int32(offService), int32(wasmPtr)+int32(serviceLen))
	log.Printf("registration for %s.%s", pkg, service)

	errOff := unsafe.Offsetof(lib.RegDetail{}.OutErrPtr)
	errInd := s.mem.GetInt32(int32(wasmPtr + int64(errOff)))

	offSvcPtr := unsafe.Offsetof(lib.RegDetail{}.OutServiceIdPtr)
	svcInd := s.mem.GetInt32(int32(wasmPtr + int64(offSvcPtr)))

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
		s.mem.SetInt64(svcInd, int64(sid.Low()))
		s.mem.SetInt64(svcInd+8, int64(sid.High()))
	}

	// write the error info back to client
	s.mem.SetInt64(errInd, int64(regErr.Low()))
	s.mem.SetInt64(errInd+8, int64(regErr.High()))

}

func (a *SysCall) Locate(sp int32) {
	log.Printf("locate called")

	//log.Printf("Locate %s,%s->%x", pkg, service, retVal)
	os.Exit(1)
}

func DebugPrint(ct int32) {
	log.Printf("----DebugPrint %d", ct)
}

func (a *SysCall) Dispatch(int32) {
	log.Printf("dispatch called")
	//print("Dispatch")
	os.Exit(1)
}
