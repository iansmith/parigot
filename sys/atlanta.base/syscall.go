package sys

import (
	"log"
	"os"
	"unsafe"

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

func (s *SysCall) ReadString(structPtr int64, dataOffset uintptr, lenOffset uintptr) string {
	return s.mem.LoadStringWithLen(int32(structPtr)+int32(dataOffset), int32(structPtr)+int32(lenOffset))
}
func (s *SysCall) Write64BitPair(structPtr int64, dataOffset uintptr, id lib.Id) {
	derefed := s.mem.GetInt32(int32(structPtr + int64(dataOffset)))
	// write the error info back to client
	s.mem.SetInt64(derefed, int64(id.Low()))
	s.mem.SetInt64(derefed+8, int64(id.High()))
}

func (s *SysCall) Register(sp int32) {
	wasmPtr := s.mem.GetInt64(sp + 8)

	pkg := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.RegDetail{}.PkgPtr),
		unsafe.Offsetof(lib.RegDetail{}.PkgLen))

	service := s.ReadString(wasmPtr,
		unsafe.Offsetof(lib.RegDetail{}.ServicePtr),
		unsafe.Offsetof(lib.RegDetail{}.ServiceLen))

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
			unsafe.Offsetof(lib.RegDetail{}.OutServiceIdPtr), sid)
		//s.mem.SetInt64(svcInd, int64(sid.Low()))
		//s.mem.SetInt64(svcInd+8, int64(sid.High()))
	}
	s.Write64BitPair(wasmPtr,
		unsafe.Offsetof(lib.RegDetail{}.OutErrPtr), regErr)

	// write the error info back to client
	//s.mem.SetInt64(errInd, int64(regErr.Low()))
	//s.mem.SetInt64(errInd+8, int64(regErr.High()))

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
