package syscall

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"unsafe"

	"github.com/iansmith/parigot/api_impl/splitutil"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Your IDE may complain about calls to functions in call_js.s and/or calljs.go.  It may claim that these
// are not defined when in fact they are defined, if tricky.    If it really bothers you, most likely you
// can change this be setting the tag "js".  This needs to be defined to get the code in calljs.go.

// callImpl is a thin wrapper around kernel functionality intended to be run by clients in WASM.
type callImpl struct {
}

var envVerbose = os.Getenv("PARIGOT_VERBOSE")

// Flip this switch for debug output.
var libparigotVerbose = false || envVerbose != ""

// Locate is a kernel request that returns either a reference to the service
// or an error.  In the former case, the token returned can be used with Dispatch()
// to make a call on a remote service.  It is implicit in the use of this call that
// the caller wants to be a client of the service in question.  This call can
// be made by clients or servers, but in either case the code in question becomes
// a client of the named service.
//
//go:noinline
func (l *callImpl) Locate(req *syscallmsg.LocateRequest) (*syscallmsg.LocateResponse, error) {
	resp := syscallmsg.LocateResponse{}
	return splitImplementation[*syscallmsg.LocateRequest, *syscallmsg.LocateResponse](l, req, &resp, locate)
}

// Exit is called from the WASM side to cause the WASM program to exit.  This is implemented by causing
// the WASM code to panic and then using recover to catch it and then the program is stopped and the kernel
// will marke it dead and so forth.
func (l *callImpl) Exit(in *syscallmsg.ExitRequest) {
	// resp := syscallmsg.ExitResponse{}
	// splitImplementation[*syscallmsg.ExitRequest, *syscallmsg.ExitResponse](l, in, &resp, exit)
	// panic(resp)
}

// Dispatch is the primary means that a caller can send an RPC message.
// If you are in local development mode, this call is handled by the kernel
// itself, otherwise it implies a remote procedure call.  This method
// checks the returned response for errors. If there are errors inside the
// result they are pulled out and returned in the error parameter.  Thus
// if the error parameter is nil, the Dispatch() occurred successfully.
// This is code that runs on the WASM side.
func (l *callImpl) Dispatch(req *syscallmsg.DispatchRequest) (*syscallmsg.DispatchResponse, error) {
	resp := syscallmsg.DispatchResponse{}
	return splitImplementation[*syscallmsg.DispatchRequest, *syscallmsg.DispatchResponse](l, req, &resp, dispatch)

}

// BlockUntilCall is used to block a process until a request is received from another process.  Even when
// all the "processes" are in a single process for debugging, the BlockUntilCall is for the same purpose.
func (l *callImpl) BlockUntilCall(in *syscallmsg.BlockUntilCallRequest, canTimeout bool) (*syscallmsg.BlockUntilCallResponse, error) {
	resp := &syscallmsg.BlockUntilCallResponse{}
	return splitImplementation[*syscallmsg.BlockUntilCallRequest, *syscallmsg.BlockUntilCallResponse](l, in, resp, blockUntilCall)
}

// BindMethodIn binds a method that only has an in parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func (l *callImpl) BindMethodIn(in *syscallmsg.BindMethodRequest, _ func(*protosupportmsg.Pctx, proto.Message) error) (*syscallmsg.BindMethodResponse, error) {
	return l.bindMethodByName(in, syscallmsg.MethodDirection_METHOD_DIRECTION_IN)
}

// BindMethodInNoPctx binds a method that only has an in parameter and does not
// use the Pctx mechanism for logging.  This may, in fact, be a terrible idea but one
// cannot write a separate logger server with having this.
// xxxfixme: temporary? Should this be a different kernel call?
func (l *callImpl) BindMethodInNoPctx(in *syscallmsg.BindMethodRequest, _ func(proto.Message) error) (*syscallmsg.BindMethodResponse, error) {
	return l.bindMethodByName(in, syscallmsg.MethodDirection_METHOD_DIRECTION_IN)
}

// BindMethodOut binds a method that only has an out parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func (l *callImpl) BindMethodOut(in *syscallmsg.BindMethodRequest, _ func(*protosupportmsg.Pctx) (proto.Message, error)) (*syscallmsg.BindMethodResponse, error) {
	return l.bindMethodByName(in, syscallmsg.MethodDirection_METHOD_DIRECTION_OUT)
}

// BindMethodBoth binds a method that has both an in and out parameter.  This should
// only be called by servers because it provides the implementation of the
// method in question.  The returned response includes a MethodId and an error.
// If there was an error, it is pulled out and returned in the 2nd result here.
// MethodIds are opaque tokens that the kernel uses to communicate to an
// implementing server which method has been invoked.
func (l *callImpl) BindMethodBoth(in *syscallmsg.BindMethodRequest, _ func(*protosupportmsg.Pctx, proto.Message) (proto.Message, error)) (*syscallmsg.BindMethodResponse, error) {
	return l.bindMethodByName(in, syscallmsg.MethodDirection_METHOD_DIRECTION_BOTH)
}

// bindMethodByName is the implementation of all three of the Bind* calls.
func (l *callImpl) bindMethodByName(in *syscallmsg.BindMethodRequest, dir syscallmsg.MethodDirection) (*syscallmsg.BindMethodResponse, error) {
	in.Direction = dir
	out := &syscallmsg.BindMethodResponse{}
	return splitImplementation[*syscallmsg.BindMethodRequest, *syscallmsg.BindMethodResponse](l, in, out, bindMethod)
}

// Run is used only at startup.  A call is made to Run when a program has notified the kernel of all the
// requires and exports it has.  When Run() is called, the program calling
func (l *callImpl) Run(in *syscallmsg.RunRequest) (*syscallmsg.RunResponse, error) {
	out := new(syscallmsg.RunResponse)
	r, e := splitImplementation[*syscallmsg.RunRequest, *syscallmsg.RunResponse](l, in, out, run)
	return r, e
}

// Export is the way that a server can express that it is done binding methods
// the service and it is ready to export it.  This call does not block.  If the input
// structure has multiple services in it, this call will repeatedly call
// the kernel and it will abort and return the error at the first failure.
func (l *callImpl) Export(in *syscallmsg.ExportRequest) (*syscallmsg.ExportResponse, error) {
	resp := syscallmsg.ExportResponse{}
	return splitImplementation[*syscallmsg.ExportRequest, *syscallmsg.ExportResponse](l, in, &resp, export)
}

// Returnvalue is a call that a service (implementation of a function) uses to return values when a
// function it implements the function. This is means that the request to "send" a return value must
// go from the WASM side to the go side.  Thus the calls that start a function call (Dispatch and
// CallService) also have to send control *back* to WASM from their go implementation.
func (l *callImpl) ReturnValue(req *syscallmsg.ReturnValueRequest) (*syscallmsg.ReturnValueResponse, error) {
	resp := syscallmsg.ReturnValueResponse{}
	return splitImplementation(l, req, &resp, returnValue)
}

// Require is the way that a client or server can express that uses a particular
// interface.  This call does not block.  If the input structure has multiple
// services in it, this call will repeatedly call the kernel and it will abort
// and return the error at the first failure.
func (l *callImpl) Require(in *syscallmsg.RequireRequest) (*syscallmsg.RequireResponse, error) {
	resp := syscallmsg.RequireResponse{}
	return splitImplementation(l, in, &resp, require)
}

// // Use of this function is discouraged.  This is intended only for debugging the parigot implementation
// // itself. User code should use the normal LocateLog() and the Log service.
// func (l *callImpl) BackdoorLog(in *logmsg.LogRequest) (*logmsg.LogResponse, error) {
// 	// this implementation is unique to this call because it will cause an infinite recursion
// 	// if we end using any function that call this logging facility.
// 	out := &logmsg.LogResponse{}
// 	sp := splitutil.NewSinglePayload()
// 	buff, err := proto.Marshal(in)
// 	if err != nil {
// 		return nil, err
// 	}
// 	sp.InPtr, sp.InLen = sliceToTwoInt64s(buff)
// 	backdoorLog(int32(uintptr(unsafe.Pointer(sp))))

// 	var high [8]byte
// 	binary.LittleEndian.PutUint64(high[:], uint64(sp.ErrPtr[0]))
// 	var errId lib.Id
// 	if high[6]&1 == 1 && sp.ErrPtr[1] != 0 {
// 		if high[7] != 107 {
// 			panic("returned error was not a kernel error")
// 		}
// 		errId = lib.NewKernelError(lib.KernelErrorCode(sp.ErrPtr[1]))
// 	}
// 	if checkIdForError(errId) {
// 		return nil, idErrorToPerror(errId, "backdoor log failed")
// 	}
// 	return out, nil
// }

// log is used to by the callImpl code to get debug messages on the terminal.
func (l *callImpl) log(funcName string, spec string, rest ...interface{}) {
	p1 := fmt.Sprintf("callImpl.%s", funcName)
	p2 := fmt.Sprintf(spec, rest...)
	if !strings.HasSuffix(p2, "\n") {
		p2 += "\n"
	}
	req := logmsg.LogRequest{
		Stamp:   timestamppb.Now(), //xxx should be using the kernel version
		Level:   logmsg.LogLevel_LOG_LEVEL_DEBUG,
		Message: p1 + p2,
	}
	print("xxx call impl on client side: %s\n", req.Message)
	// _, err := l.BackdoorLog(&req)
	// if err != nil {
	// 	panic("backdoorLog failed:" + err.Error())
	// }
}

// splitImplementation is the implementation for any client side program that calls callImpl.  This
// is based on the splitutil.SendReceiveSingleProto in terms of how the data moves around.   This
// code simply sets up the request and response structures and then sends these to the SendReceiveSingleProto
// which does the work of flattening the content of the parameters and sends that blob of bytes to the
// kernel.  Similarly, it does the unpacking of the result sent from the Kernel, including handling
// errors.
func splitImplementation[T proto.Message, U proto.Message](l *callImpl, req T, resp U, fn func(int32)) (U, error) {
	var zeroValForU U
	_, errId, errDetail := splitutil.SendReceiveSingleProto(l, req, resp, fn)
	if errId != nil {
		return zeroValForU, lib.NewPerrorFromId(errDetail, errId)
	}
	return resp, nil
}

// Export1 is a wrapper around Export which makes it easy to say you export a single
// service. It does not change any of the Export behavior.
func (l *callImpl) Export1(packagePath, service string) (*syscallmsg.ExportResponse, error) {
	fqSvc := &syscallmsg.FullyQualifiedService{
		PackagePath: packagePath, Service: service}
	req := &syscallmsg.ExportRequest{}
	req.Service = []*syscallmsg.FullyQualifiedService{fqSvc}
	return l.Export(req)
}

// Require1 is a wrapper around Require which makes it easy to say you require a single
// service. It does not change any of the Require behavior.
func (l *callImpl) Require1(packagePath, service string) (*syscallmsg.RequireResponse, error) {
	fqSvc := &syscallmsg.FullyQualifiedService{
		PackagePath: packagePath, Service: service}
	req := &syscallmsg.RequireRequest{}
	req.Service = []*syscallmsg.FullyQualifiedService{fqSvc}
	return l.Require(req)
}

func NewCallImpl() lib.Call {
	return &callImpl{}
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

// sliceToTwoInt64s is a utility used to populate a payload.
func sliceToTwoInt64s(b []byte) (int64, int64) {
	slh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return int64(slh.Data), int64(slh.Len)
}
