package lib

import (
	"fmt"
	"runtime/debug"

	"github.com/iansmith/parigot/api/proto/g/pb/call"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var helperVerbose = true

// ReturnValueEncode is a layer on top of ReturnValue.  This functions exists because
// there are number of cases and doing this in this library means the code generator
// can be much simpler.  It just passes all the information into here, and this function
// sorts it out.
func ReturnValueEncode(callImpl Call, cid, mid Id, marshalError, execError error, out proto.Message, pctx *protosupport.Pctx) (*call.ReturnValueResponse, error) {
	var err error
	var a anypb.Any
	// xxxfixme we should be doing an examination of execError to see if it is a Perror
	// xxxfixme and if it is, we should be pushing the user error back the other way
	rv := &call.ReturnValueRequest{}
	rv.Call = Marshal[protosupport.CallId](cid)
	rv.Method = Marshal[protosupport.MethodId](mid)
	rv.ErrorId = NoKernelError() // just to allocate the space
	if marshalError != nil || execError != nil {
		helperprint("ReturnValueEncode ", "[%s],[%s]: marshalError? %v execError ? %v",
			cid.Short(), mid.Short(), marshalError, execError)
		debug.PrintStack()
		print("END OF STACK")
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

	helperprint("RETURNVALUEENCODE ", "size of result buffer %d, out %s",
		len(rv.ResultBuffer), out)
	return callImpl.ReturnValue(rv)
internalMarshalProblem:
	helperprint("RETURNVALUEENCODE ", "internal encoding error: %v", err)

	// this is an internal error, so we signal it the opposite way we did the others at the top
	rv.ErrorMessage = ""
	rv.ErrorId = NoKernelError()
encodeError:
	rv.PctxBuffer = []byte{}
	rv.ResultBuffer = []byte{}
	return callImpl.ReturnValue(rv)
}

func helperprint(fnName string, spec string, arg ...interface{}) {
	if helperVerbose {
		part1 := fmt.Sprintf("HELPER:%s ", fnName)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}
