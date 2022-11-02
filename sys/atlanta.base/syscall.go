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

const MaxService = 127

var packageRegistry = make(map[string]*packageData)
var serviceCounter = 7 // last USED service num

type packageData struct {
	service map[string]*serviceData
}

func newPackageData() *packageData {
	return &packageData{
		service: make(map[string]*serviceData),
	}
}

type serviceData struct {
	serviceId lib.Id
	method    map[string]lib.Id
}

func newServiceData() *serviceData {
	return &serviceData{
		serviceId: nil,
		method:    make(map[string]lib.Id),
	}
}

type SysCall struct {
	mem     *jspatch.WasmMem
	pkgData map[string]*packageData
}

func (s *SysCall) SetMemPtr(m uintptr) {
	s.mem = jspatch.NewWasmMem(m)
}

func NewSysCall() *SysCall {
	return &SysCall{
		pkgData: make(map[string]*packageData),
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

	pData, ok := s.pkgData[pkg]
	if !ok {
		pData = newPackageData()
		packageRegistry[pkg] = pData
	}
	_, duplicate := pData.service[service]

	var regErr lib.Id
	switch {
	case serviceCounter >= MaxService:
		regErr = lib.NewKernelError(lib.KernelNamespaceExhausted)
	case duplicate:
		regErr = lib.NewKernelError(lib.KernelAlreadyRegistered)
	default:
		regErr = lib.NoKernelErr()
	}
	// assign sid if no error
	if !regErr.IsError() {
		sid := lib.ServiceIdFromUint64(0, uint64(serviceCounter+1))
		serviceCounter++
		sData := newServiceData()
		sData.serviceId = sid
		pData.service[service] = sData
		//send back the data to client
		s.Write64BitPair(wasmPtr,
			unsafe.Offsetof(lib.RegPayload{}.ServiceIdPtr), sid)
	}
	// either case, we want to tell the client error status
	s.Write64BitPair(wasmPtr,
		unsafe.Offsetof(lib.RegPayload{}.ErrorPtr), regErr)
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
	pData, ok := s.pkgData[pkg]
	if !ok {
		l := lib.NewKernelError(lib.KernelNotFound)
		locErr = &l
	} else {
		sData, ok := pData.service[service]
		if !ok {
			l := lib.NewKernelError(lib.KernelNotFound)
			locErr = &l
		} else {
			serviceId = &sData.serviceId
		}
	}
	// now we need to write the results back to the client
	switch {
	case locErr != nil:
		s.Write64BitPair(wasmPtr,
			unsafe.Offsetof(lib.LocatePayload{}.ErrorPtr), *locErr)
	case serviceId != nil:
		l := lib.NoKernelErr()
		locErr = &l
		s.Write64BitPair(wasmPtr,
			unsafe.Offsetof(lib.LocatePayload{}.ErrorPtr), *locErr)
		s.Write64BitPair(wasmPtr,
			unsafe.Offsetof(lib.LocatePayload{}.ServiceIdPtr), *serviceId)
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

	pData, ok := s.pkgData[pkg]
	if !ok {
		//we allow this because we don't know if client or server will connect first
		pData = newPackageData()
		s.pkgData[pkg] = pData
	}
	sData, ok := pData.service[service]
	if !ok {
		sData = newServiceData()
		pData.service[service] = sData
	}
	mid := lib.NewMethodId()
	sData.method[method] = mid
	noErr := lib.NoKernelErr() // the lack of an error

	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.MethodId), mid)
	s.Write64BitPair(wasmPtr, unsafe.Offsetof(lib.BindPayload{}.ErrorPtr),
		noErr)

	log.Printf("bind completed")
}
