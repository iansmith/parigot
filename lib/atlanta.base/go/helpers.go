package lib

import (
	"github.com/iansmith/parigot/api/proto/g/pb/call"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// Export1 is a wrapper around Export which makes it easy to say you export a single
// service. It does not change any of the Export behavior.
func Export1(packagePath, service string) (*call.ExportResponse, error) {
	fqSvc := &call.FullyQualifiedService{
		PackagePath: packagePath, Service: service}
	req := &call.ExportRequest{}
	req.Service = []*call.FullyQualifiedService{fqSvc}
	return CallConnection().Export(req)
}

// Require1 is a wrapper around Require which makes it easy to say you require a single
// service. It does not change any of the Require behavior.
func Require1(packagePath, service string) (*call.RequireResponse, error) {
	fqSvc := &call.FullyQualifiedService{
		PackagePath: packagePath, Service: service}
	req := &call.RequireRequest{}
	req.Service = []*call.FullyQualifiedService{fqSvc}
	return CallConnection().Require(req)
}

// ReturnValueEncode is a layer on top of ReturnValue.  This functions exists because
// there are number of cases and doing this in this library means the code generator
// can be much simpler.  It just passes all the information into here, and this function
// sorts it out.
func ReturnValueEncode(cid, mid Id, marshalError, execError error, out proto.Message, pctx *protosupport.Pctx) (*ReturnValueResponse, error) {
	libprint("RETURNVALUEENCODE ", "in return value %s, %s", cid.Short(), mid.Short())
	var err error
	var a anypb.Any
	// xxxfixme we should be doing an examination of execError to see if it is a lib.Perror
	// xxxfixme and if it is, we should be pushing the user error back the other way
	rv := &call.ReturnValueRequest{}
	rv.Call = Marshal[protosupport.CallId](cid)
	rv.Method = Marshal[protosupport.MethodId](mid)
	rv.ErrorId = NoKernelError() // just to allocate the space
	libprint("RETURNVALUEENCODE ", "marshalError? %v execError ? %v", marshalError, execError)
	if marshalError != nil || execError != nil {
		eText := ""
		if marshalError != nil {
			rv.ErrorMessage = marshalError.Error()
			eText += "marshalError:" + marshalError.Error()
		} else {
			rv.ErrorMessage = execError.Error()
			eText += "execError:" + execError.Error()
		}
		//pctx.Log(log.LogLevel_LOGLEVEL_ERROR, eText)
		//pctx.EventFinish()
		// we use this because the error didn't come from INSIDE
		// the kernel itself, see below for more
		rv.ErrorId = NoKernelError()
		goto encodeError
	}
	// these are the mostly normal cases, but they can go hawywire
	// due to marshalling
	//pctx.EventFinish()
	//libprint("RETURNVALUEENCODE -- log -- \n", pctx.Dump())
	rv.PctxBuffer, err = proto.Marshal(pctx)
	if err != nil {
		goto internalMarshalProblem
	}
	if out != nil {
		err = a.MarshalFrom(out)
		if err != nil {
			goto internalMarshalProblem
		}
		rv.ResultBuffer, err = proto.Marshal(&a)
		if err != nil {
			goto internalMarshalProblem
		}
	} else {
		rv.ResultBuffer = nil
	}

	libprint("RETURNVALUEENCODE ", "size of result buffer %d, out %s", len(rv.ResultBuffer), out)
	return CallConnection().ReturnValue(rv)
internalMarshalProblem:
	libprint("RETURNVALUEENCODE ", "internal encoding error: %v", err)

	// this is an internal error, so we signal it the opposite way we did the others at the top
	rv.ErrorMessage = ""
	rv.ErrorId = Marshal[protosupport.KernelErrorId](NewKernelError(KernelMarshalFailed))
encodeError:
	rv.PctxBuffer = []byte{}
	rv.ResultBuffer = []byte{}
	return CallConnection().ReturnValue(rv)
}

func Run(wait bool) (*call.RunResponse, error) {
	req := &call.RunRequest{
		Wait: wait,
	}
	return CallConnection().Run(req)
}

func Exit(code int32) {
	req := &call.ExitRequest{
		Code: code,
	}
	CallConnection().Exit(req)
	return // this will not happen, previous will cause the proc to die
}
