// This package is a thin wrapper around kernel functionality intended to be run by clients in WASM.
package syscall

import (
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
	"strings"
	"unsafe"

	"github.com/iansmith/parigot/api/proto/g/pb/call"
	pblog "github.com/iansmith/parigot/api/proto/g/pb/log"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	pbsys "github.com/iansmith/parigot/api/proto/g/pb/syscall"
	"github.com/iansmith/parigot/api/splitutil"
	"github.com/iansmith/parigot/lib"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Your IDE may complain about calls to functions in call_js.s and/or calljs.go.  It may claim that these
// are not defined when in fact they are defined, if tricky.    If it really bothers you, most likely you
// can change this be setting the tag "js".  This needs to be defined to get the code in calljs.go.

type callImpl struct {
	// logger log.Log
}

func (l *callImpl) Exit(in *call.ExitRequest) {
	exit(in.Code)
}

var envVerbose = os.Getenv("PARIGOT_VERBOSE")

// Flip this switch for debug output.
var libparigotVerbose = true || envVerbose != ""

// Locate is a kernel request that returns either a reference to the service
// or an error.  In the former case, the token returned can be used with Dispatch()
// to make a call on a remote service.  It is implicit in the use of this call that
// the caller wants to be a client of the service in question.  This call can
// be made by clients or servers, but in either case the code in question becomes
// a client of the named service.
//
//go:noinline
func (l *callImpl) Locate(req *pbsys.LocateRequest) (*pbsys.LocateResponse, error) {

	resp := pbsys.LocateResponse{}
	return splitImplementation[*pbsys.LocateRequest, *pbsys.LocateResponse](l, req, &resp, locate)
	// l.log("Locate", "client side Locate called %s.%s", req.GetPackageName(), req.GetServiceName())
	// id, err := splitutil.SendReceiveSingleProto(l, req, &resp, locate)

	// if err != nil {
	// 	return nil, err
	// }

	// if checkIdForError(id) {
	// 	return nil, idErrorToPerror(id, "failed to locate properly")
	// }

	// return &resp, nil
}

// sliceToTwoInt64s is a utility used to populate a payload.
func sliceToTwoInt64s(b []byte) (int64, int64) {
	slh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return int64(slh.Data), int64(slh.Len)
}

// Dispatch is the primary means that a caller can send an RPC message.
// If you are in local development mode, this call is handled by the kernel
// itself, otherwise it implies a remote procedure call.  This method
// checks the returned response for errors. If there are errors inside the
// result they are pulled out and returned in the error parameter.  Thus
// if the error parameter is nil, the Dispatch() occurred successfully.
// This is code that runs on the WASM side.
func (l *callImpl) Dispatch(req *pbsys.DispatchRequest) (*pbsys.DispatchResponse, error) {
	resp := pbsys.DispatchResponse{}
	return splitImplementation[*pbsys.DispatchRequest, *pbsys.DispatchResponse](l, req, &resp, dispatch)

	// id, err := splitutil.SendReceiveSingleProto(l, req, &resp, dispatch)
	// if err != nil {
	// 	return nil, err
	// }
	// if checkIdForError(id) {
	// 	return nil, idErrorToPerror(id, "failed to dispatch properly")
	// }
	// return &resp, nil
}

// checkIdForError returns true in the when it found an error value in the provided id, false otherwise.
func checkIdForError(id lib.Id) bool {
	if id != nil {
		if id.IsErrorType() {
			if id.IsError() {
				return true
			}
		} else {
			panic(fmt.Sprintf("response is unexpected id type: isErrorType=%v, id=%s", id.IsErrorType(), id.Short()))
		}
	}
	return false
}

// idErrorToPerror returns an error suitable for returning to user code.
func idErrorToPerror(id lib.Id, message string) lib.Error {
	return lib.NewPerrorFromId(message, id)
}

// BindMethodIn binds a method that only has an in parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func (l *callImpl) BindMethodIn(in *call.BindMethodRequest, _ func(*protosupport.Pctx, proto.Message) error) (*call.BindMethodResponse, error) {
	return l.bindMethodByName(in, call.MethodDirection_METHOD_DIRECTION_IN)
}

// BindMethodInNoPctx binds a method that only has an in parameter and does not
// use the Pctx mechanism for logging.  This may, in fact, be a terrible idea but one
// cannot write a separate logger server with having this.
// xxxfixme: temporary? Should this be a different kernel call?
func (l *callImpl) BindMethodInNoPctx(in *call.BindMethodRequest, _ func(proto.Message) error) (*call.BindMethodResponse, error) {
	return l.bindMethodByName(in, call.MethodDirection_METHOD_DIRECTION_IN)
}

// BindMethodOut binds a method that only has an out parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func (l *callImpl) BindMethodOut(in *call.BindMethodRequest, _ func(*protosupport.Pctx) (proto.Message, error)) (*call.BindMethodResponse, error) {
	return l.bindMethodByName(in, call.MethodDirection_METHOD_DIRECTION_OUT)
}

// BindMethodBoth binds a method that has both an in and out parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func (l *callImpl) BindMethodBoth(in *call.BindMethodRequest, _ func(*protosupport.Pctx, proto.Message) (proto.Message, error)) (*call.BindMethodResponse, error) {
	return l.bindMethodByName(in, call.MethodDirection_METHOD_DIRECTION_BOTH)
}

func (l *callImpl) bindMethodByName(in *call.BindMethodRequest, dir call.MethodDirection) (*call.BindMethodResponse, error) {
	out := new(call.BindMethodResponse)

	out.ErrorId = lib.NoKernelError()

	detail := new(lib.BindPayload)
	sh := (*reflect.StringHeader)(unsafe.Pointer(&in.ProtoPackage))
	detail.PkgPtr = int64(sh.Data)
	detail.PkgLen = int64(sh.Len)
	sh = (*reflect.StringHeader)(unsafe.Pointer(&in.Service))
	detail.ServicePtr = int64(sh.Data)
	detail.ServiceLen = int64(sh.Len)
	sh = (*reflect.StringHeader)(unsafe.Pointer(&in.Method))
	detail.MethodPtr = int64(sh.Data)
	detail.MethodLen = int64(sh.Len)
	detail.Direction = int64(dir)

	// need to allocate the space for result
	out.MethodId = &protosupport.MethodId{}
	out.ErrorId = &protosupport.KernelErrorId{}

	out.MethodId.Id = &protosupport.BaseId{}
	out.ErrorId.Id = &protosupport.BaseId{}

	detail.MethodId = (*[2]int64)(unsafe.Pointer(&out.MethodId.Id.Low))
	detail.ErrorPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Id.Low))

	// THE CALL
	u := uintptr(unsafe.Pointer(detail))
	bindMethod(int32(u))

	// check for in band errors
	kernelErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.ErrorPtr))))
	kerr := lib.NewKernelError(lib.KernelErrorCode(kernelErrDataPtr[0]))
	if kerr.IsError() {
		return nil, lib.NewPerrorFromId("bind error", kerr)
	}

	methodDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.MethodId))))
	mid := lib.NewFrom64BitPair[*protosupport.MethodId](uint64(methodDataPtr[1]), uint64(uint64(methodDataPtr[0])))
	out.MethodId = lib.Marshal[protosupport.MethodId](mid)
	out.ErrorId = lib.NoKernelError()

	return out, nil
}

func (l *callImpl) Run(in *call.RunRequest) (*call.RunResponse, error) {
	out := new(call.RunResponse)
	detail := new(lib.RunPayload)
	detail.Wait = 0
	if in.Wait {
		detail.Wait = 1
	}
	out.ErrorId = &protosupport.KernelErrorId{
		Id: &protosupport.BaseId{},
	}
	detail.KernelErrorPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Id.Low))

	// THE CALL, and the walls came down...
	u := uintptr(unsafe.Pointer(detail))
	run(int32(u))

	kernelErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.KernelErrorPtr))))
	kerr := lib.NewKernelError(lib.KernelErrorCode(kernelErrDataPtr[0]))
	if kerr != nil && kerr.IsError() {
		out.ErrorId = lib.Marshal[protosupport.KernelErrorId](kerr)
		return out, lib.NewPerrorFromId("kernel failed to start your process", kerr)
	}
	out.ErrorId = lib.NoKernelError()
	return out, nil
}

func (l *callImpl) exportOrRequire(fqs []*call.FullyQualifiedService, errorPtr *[2]int64, isExport bool) error {
	for _, s := range fqs {
		pkg := s.GetPackagePath()
		svc := s.GetService()

		detail := new(lib.ExportPayload)
		sh := (*reflect.StringHeader)(unsafe.Pointer(&pkg))
		detail.PkgPtr = int64(sh.Data)
		detail.PkgLen = int64(sh.Len)
		sh = (*reflect.StringHeader)(unsafe.Pointer(&svc))
		detail.ServicePtr = int64(sh.Data)
		detail.ServiceLen = int64(sh.Len)
		detail.KernelErrorPtr = errorPtr //(*[2]int64)(unsafe.Pointer(&out.ErrorId.Low))

		// THE CALL
		u := uintptr(unsafe.Pointer(detail))
		if isExport {
			export(int32(u))
		} else {
			require(int32(u))
		}

		// check for in band errors
		kernelErrDataPtr := detail.KernelErrorPtr //(*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.KernelErrorPtr))))
		kerr := lib.NewKernelError(lib.KernelErrorCode(kernelErrDataPtr[0]))
		if kerr.IsError() {
			if isExport {
				return lib.NewPerrorFromId("export error", kerr)
			} else {
				return lib.NewPerrorFromId("require error", kerr)
			}

		}
	}
	return nil

}

// Export is the way that a server can express that it is done binding methods
// the service and it is ready to export it.  This call does not block.  If the input
// structure has multiple services in it, this call will repeatedly call
// the kernel and it will abort and return the error at the first failure.
func (l *callImpl) Export(in *call.ExportRequest) (*call.ExportResponse, error) {
	out := new(call.ExportResponse)
	// allocate space for any error
	out.ErrorId = lib.NoKernelError()
	ptr := (*[2]int64)(unsafe.Pointer(&out.ErrorId.Id.Low))
	err := l.exportOrRequire(in.GetService(), ptr, true)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Require is the way that a client or server can express that uses a particular
// interface.  This call does not block.  If the input structure has multiple
// services in it, this call will repeatedly call the kernel and it will abort
// and return the error at the first failure.
func (l *callImpl) Require(in *call.RequireRequest) (*call.RequireResponse, error) {
	out := new(call.RequireResponse)
	// allocate space for any error
	out.ErrorId = lib.NoKernelError()
	ptr := (*[2]int64)(unsafe.Pointer(&out.ErrorId.Id.Low))
	err := l.exportOrRequire(in.GetService(), ptr, false)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// xxx this may be a bad idea.  this is probably only temporary til I can work out if we
// xxx should support this at all. currently, it is used by terminal logger.  I made a copy
// xxx of the other bindMethodByName so it would be easier to delete this one.
func (l *callImpl) bindMethodByNameNoPctx(in *call.BindMethodRequest, dir call.MethodDirection) (*call.BindMethodResponse, error) {
	out := new(call.BindMethodResponse)

	out.ErrorId = lib.NoKernelError()

	detail := new(lib.BindPayload)
	sh := (*reflect.StringHeader)(unsafe.Pointer(&in.ProtoPackage))
	detail.PkgPtr = int64(sh.Data)
	detail.PkgLen = int64(sh.Len)
	sh = (*reflect.StringHeader)(unsafe.Pointer(&in.Service))
	detail.ServicePtr = int64(sh.Data)
	detail.ServiceLen = int64(sh.Len)
	sh = (*reflect.StringHeader)(unsafe.Pointer(&in.Method))
	detail.MethodPtr = int64(sh.Data)
	detail.MethodLen = int64(sh.Len)
	detail.Direction = int64(dir)
	detail.MethodId = (*[2]int64)(unsafe.Pointer(&out.MethodId.Id.Low))
	detail.ErrorPtr = (*[2]int64)(unsafe.Pointer(&out.ErrorId.Id.Low))

	// THE CALL
	u := uintptr(unsafe.Pointer(detail))
	bindMethod(int32(u))

	// check for in band errors
	kernelErrDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.ErrorPtr))))
	kerr := lib.NewKernelError(lib.KernelErrorCode(kernelErrDataPtr[0]))
	if kerr.IsError() {
		return nil, lib.NewPerrorFromId("bind error", kerr)
	}

	methodDataPtr := (*[2]int64)(unsafe.Pointer(uintptr(unsafe.Pointer(detail.MethodId))))
	mid := lib.NewFrom64BitPair[*protosupport.MethodId](uint64(methodDataPtr[1]), uint64(uint64(methodDataPtr[0])))
	out.MethodId = lib.Marshal[protosupport.MethodId](mid)

	out.ErrorId = lib.NoKernelError()

	return out, nil
}

// Use of this function is discouraged.  This is intended only for debugging the parigot implementation
// itself. User code should use the normal LocateLog() and the Log service.
func (l *callImpl) BackdoorLog(in *pblog.LogRequest) (*pblog.LogResponse, error) {
	// this implementation is unique to this call because it will cause an infinite recursion
	// if we end using any function that call this logging facility.
	out := &pblog.LogResponse{}
	sp := splitutil.NewSinglePayload()
	buff, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}
	sp.InPtr, sp.InLen = sliceToTwoInt64s(buff)
	backdoorLog(int32(uintptr(unsafe.Pointer(sp))))

	var high [8]byte
	binary.LittleEndian.PutUint64(high[:], uint64(sp.ErrPtr[0]))
	var errId lib.Id
	if high[6]&1 == 1 && sp.ErrPtr[1] != 0 {
		if high[7] != 107 {
			panic("returned error was not a kernel error")
		}
		errId = lib.NewKernelError(lib.KernelErrorCode(sp.ErrPtr[1]))
	}
	if checkIdForError(errId) {
		return nil, idErrorToPerror(errId, "backdoor log failed")
	}
	return out, nil
}

func (l *callImpl) BlockUntilCall(in *pbsys.BlockUntilCallRequest) (*pbsys.BlockUntilCallResponse, error) {
	// this is JUST for reserving the space for the result to be placed into
	out := &pbsys.BlockUntilCallResponse{}
	id, err := splitutil.SendReceiveSingleProto(l, in, out, blockUntilCall)
	if err != nil {
		return nil, err
	}
	if checkIdForError(id) {
		return nil, idErrorToPerror(id, "BlockUntilCall failed")
	}
	return out, nil
}

func (l *callImpl) log(funcName string, spec string, rest ...interface{}) {
	p1 := fmt.Sprintf("callImpl.%s", funcName)
	p2 := fmt.Sprintf(spec, rest...)
	if !strings.HasSuffix(p2, "\n") {
		p2 += "\n"
	}
	req := pblog.LogRequest{
		Stamp:   timestamppb.Now(), //xxx should be using the kernel version
		Level:   pblog.LogLevel_LOG_LEVEL_DEBUG,
		Message: p1 + p2,
	}
	_, err := l.BackdoorLog(&req)
	if err != nil {
		panic("backdoorLog failed:" + err.Error())
	}
}

func splitImplementation[T proto.Message, U proto.Message](l *callImpl, req T, resp U, fn func(int32)) (U, error) {
	var zeroValForU U
	id, err := splitutil.SendReceiveSingleProto(l, req, resp, fn)
	if err != nil {
		return zeroValForU, err
	}
	if checkIdForError(id) {
		return zeroValForU, idErrorToPerror(id, "returnValue failed")
	}
	return resp, nil

}

func (l *callImpl) ReturnValue(req *pbsys.ReturnValueRequest) (*pbsys.ReturnValueResponse, error) {
	resp := pbsys.ReturnValueResponse{}
	return splitImplementation[*pbsys.ReturnValueRequest, *pbsys.ReturnValueResponse](l, req, &resp, returnValue)
	// l.log("ReturnValue", "client side ReturnValue called %T:%d ",
	// 	req.GetResult(), proto.Size(req.GetResult()))
	// id, err := splitutil.SendReceiveSingleProto(l, req, &resp, returnValue)
	// if err != nil {
	// 	return nil, err
	// }
	// if checkIdForError(id) {
	// 	return nil, idErrorToPerror(id, "returnValue failed")
	// }
	// return &resp, nil
}

// Export1 is a wrapper around Export which makes it easy to say you export a single
// service. It does not change any of the Export behavior.
func (l *callImpl) Export1(packagePath, service string) (*call.ExportResponse, error) {
	fqSvc := &call.FullyQualifiedService{
		PackagePath: packagePath, Service: service}
	req := &call.ExportRequest{}
	req.Service = []*call.FullyQualifiedService{fqSvc}
	return l.Export(req)
}

// Require1 is a wrapper around Require which makes it easy to say you require a single
// service. It does not change any of the Require behavior.
func (l *callImpl) Require1(packagePath, service string) (*call.RequireResponse, error) {
	fqSvc := &call.FullyQualifiedService{
		PackagePath: packagePath, Service: service}
	req := &call.RequireRequest{}
	req.Service = []*call.FullyQualifiedService{fqSvc}
	return l.Require(req)
}

func NewCallImpl() lib.Call {
	return &callImpl{}
}
