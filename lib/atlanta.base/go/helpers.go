package lib

import (
	"github.com/iansmith/parigot/g/pb/call"
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

// ReturnValueEncode is a layer on top of ReturnValue.  This is here because
// there are number of cases and doing this in this library means
// the code generator can be much simpler.  It just passes all the
// information into here, and this function sorts it out.
func ReturnValueEncode(cid, mid Id, marshalError, execError error, out proto.Message, pctx Pctx) (*call.ReturnValueResponse, error) {
	var err error
	var a anypb.Any
	// xxxfixme we should be doing an examination of execError to see if it is a lib.Perror
	// xxxfixme and if it is, we should be pushing the user error back the other way
	rv := &call.ReturnValueRequest{}
	rv.Call = MarshalCallId(cid)
	rv.Method = MarshalMethodId(mid)
	rv.ErrorId = MarshalKernelErrId(NoKernelErr()) // just to allocate the space
	if marshalError != nil || execError != nil {
		if marshalError != nil {
			rv.ErrorMessage = marshalError.Error()
		} else {
			rv.ErrorMessage = execError.Error()
		}
		// we use this because the error didn't come from INSIDE
		// the kernel itself, see below for more
		rv.ErrorId = MarshalKernelErrId(NoKernelErr())
		goto encodeError
	}
	// these are the mostly normal cases, but they can go hawywire
	// due to marshalling
	rv.PctxBuffer, err = pctx.Marshal()
	if err != nil {
		goto internalMarshalProblem
	}
	err = a.MarshalFrom(out)
	if err != nil {
		goto internalMarshalProblem
	}
	rv.ResultBuffer, err = proto.Marshal(&a)
	if err != nil {
		goto internalMarshalProblem
	}

	libprint("RETURNVALUEENCODE ", "size of result buffer %d, out %s", len(rv.ResultBuffer), out)
	return CallConnection().ReturnValue(rv)
internalMarshalProblem:
	// this is an internal error, so we signal it the opposite way we did the others at the top
	rv.ErrorMessage = ""
	rv.ErrorId = MarshalKernelErrId(NewKernelError(KernelMarshalFailed))
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
