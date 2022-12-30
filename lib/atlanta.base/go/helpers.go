package lib

import (
	"fmt"
	"os"

	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	pbsys "github.com/iansmith/parigot/api/proto/g/pb/syscall"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var envVerbose = os.Getenv("PARIGOT_VERBOSE")
var helperVerbose = true || envVerbose != ""

// ReturnValueEncode is a layer on top of ReturnValue.  This functions exists because
// there are number of cases and doing this in this library means the code generator
// can be much simpler.  It just passes all the information into here, and this function
// sorts it out.
func ReturnValueEncode(callImpl Call, cid, mid Id, marshalError, execError error, out proto.Message, pctx *protosupport.Pctx) (*pbsys.ReturnValueResponse, error) {
	var err error
	var a anypb.Any
	// xxxfixme we should be doing an examination of execError to see if it is a Perror
	// xxxfixme and if it is, we should be pushing the user error back the other way
	rv := &pbsys.ReturnValueRequest{}
	rv.Call = Marshal[protosupport.CallId](cid)
	rv.Method = Marshal[protosupport.MethodId](mid)
	if marshalError != nil || execError != nil {
		if marshalError != nil {
			rv.MarshalError = marshalError.Error()
		} else {
			rv.ExecError = execError.Error()
			perr, ok := execError.(Error)
			if ok {
				rv.ExecErrorId = Marshal[protosupport.BaseId](perr.Id())
			}
		}
		goto encodeError
	}
	// these are the mostly normal cases, but they can go hawywire
	// due to marshalling
	rv.Pctx = pctx
	if out != nil {
		err = a.MarshalFrom(out)
		if err != nil {
			goto internalMarshalProblem
		}
		rv.Result = &a
	} else {
		rv.Result = nil
	}
	return callImpl.ReturnValue(rv)
internalMarshalProblem:
	// this is an internal error, so we signal it the opposite way we did the others at the top
	rv.MarshalError = "ReturnValueEncode:internal marshaling error"
encodeError:
	rv.Pctx = nil
	rv.Result = nil
	return callImpl.ReturnValue(rv)
}

func helperprint(fnName string, spec string, arg ...interface{}) {
	if helperVerbose {
		part1 := fmt.Sprintf("HELPER:%s ", fnName)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}
