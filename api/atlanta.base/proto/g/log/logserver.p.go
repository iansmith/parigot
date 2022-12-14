package log

import (
	"fmt"

	pb "github.com/iansmith/parigot/api/proto/g/pb/log"

	"github.com/iansmith/parigot/api/proto/g/pb/call"
	"github.com/iansmith/parigot/lib"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

//
// LogServer
//

type LogServer interface {
	Log(pctx lib.Pctx, in proto.Message) error
	Ready() bool
}

// Log method ids
var logMethod lib.Id

var logServerVerbose = true

func Run(impl LogServer) {
	// register all methods
	method, err := logBind(impl)
	if err != nil {
		panic("failed to register method successfully: " + method + ":" + err.Error())
	}
	// allocate buffers for values coming back to us
	paramBuf := make([]byte, lib.GetMaxMessageSize())
	pctxBuf := make([]byte, lib.GetMaxMessageSize())

	// loop on handling calls
	for {
		//
		// wait for notification
		//
		resp, err := logBlockUntilCall(pctxBuf, paramBuf)
		if err != nil {
			// error is likely local to this process
			logPrint("RUN:primary for loop ", "Unable to dispatch method call: %v", err)
			continue
		}
		logPrint("RUN: primary for loop ", "block completed, got two values:%d,%d",
			resp.PctxLen, resp.ParamLen)
		//
		// incoming values, pctx and params
		//
		var pctxSlice []byte
		if resp.PctxLen == 0 {
			pctxSlice = []byte{}
		} else {
			pctxSlice = pctxBuf[:resp.PctxLen]
		}
		paramSlice := paramBuf[:resp.ParamLen]
		mid := lib.UnmarshalMethodId(resp.GetMethod())
		cid := lib.UnmarshalCallId(resp.GetCall())

		//
		// create the generic params, pctx and param
		//
		var pctx lib.Pctx
		err = nil
		if resp.PctxLen != 0 {
			pctx, err = lib.NewPctxFromBytes(pctxSlice)
		}
		if err != nil {
			logPrint("RUN: primary for loop ", "Unable to create Pctx for call: %v", err)
			continue
		}
		// a is an any that represents the params
		a := &anypb.Any{}
		err = proto.Unmarshal(paramSlice, a)
		if err != nil {
			logPrint("RUN: primary for loop ", "Unable to create parameters for call: %v", err)
			continue
		}

		//
		// pick the method
		//
		var marshalError, execError error
		var out proto.Message
		switch {
		case mid.Equal(logMethod):
			pctx.EventStart("---> call of Log.Log <---")
			req := &pb.LogRequest{}
			marshalError = a.UnmarshalTo(req)
			if marshalError != nil {
				break
			}

			execError = impl.Log(pctx, req)
			if execError != nil {
				break
			}
		}
		//
		// could be error, could be everything is cool, send to lib to figure it out
		//
		lib.ReturnValueEncode(cid, mid, marshalError, execError, out, pctx)
		// about to loop again
	}
	// wont reach here
}

func logBind(impl LogServer) (string, error) {

	resp, logerr := lib.CallConnection().BindMethodIn(&call.BindMethodRequest{
		ProtoPackage: "log",
		Service:      "Log",
		Method:       "Log",
	}, impl.Log)
	if logerr != nil {
		return "Log", logerr
	}
	logMethod = lib.UnmarshalMethodId(resp.GetMethodId())
	if !impl.Ready() {
		panic("unable to start Log because it failed Ready() check")
	}
	return "", nil
}

func logBlockUntilCall(pctx, param []byte) (*call.BlockUntilCallResponse, error) {

	req := &call.BlockUntilCallRequest{
		PctxBuffer:  pctx,
		ParamBuffer: param,
	}
	logPrint("logBlockUntilCall", "about to call BlockUntilCall on CallConnection1")
	logPrint("logBlockUntilCall", "about to call BlockUntilCall on CallConnection2")
	resp, err := lib.CallConnection().BlockUntilCall(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func logPrint(method string, spec string, arg ...interface{}) {
	if logServerVerbose {
		part1 := fmt.Sprintf("logServer:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}
