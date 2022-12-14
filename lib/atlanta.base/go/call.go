package lib

import (
	"github.com/iansmith/parigot/api/proto/g/pb/call"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"

	"google.golang.org/protobuf/proto"
)

type Call interface {
	Exit(in *call.ExitRequest)
	Locate(in *call.LocateRequest) (*call.LocateResponse, error)
	Dispatch(in *call.DispatchRequest) (*call.DispatchResponse, error)
	BindMethodIn(in *call.BindMethodRequest, _ func(*protosupport.Pctx, proto.Message) error) (*call.BindMethodResponse, error)
	BindMethodOut(in *call.BindMethodRequest, _ func(*protosupport.Pctx) (proto.Message, error)) (*call.BindMethodResponse, error)
	BindMethodBoth(in *call.BindMethodRequest, _ func(*protosupport.Pctx, proto.Message) (proto.Message, error)) (*call.BindMethodResponse, error)
	Run(in *call.RunRequest) (*call.RunResponse, error)
	Export(in *call.ExportRequest) (*call.ExportResponse, error)
	Require(in *call.RequireRequest) (*call.RequireResponse, error)
	BlockUntilCall(in *call.BlockUntilCallRequest) (*call.BlockUntilCallResponse, error)
	ReturnValue(in *call.ReturnValueRequest) (*call.ReturnValueResponse, error)
}

var connector Call = newCallImpl()

func CallConnection() Call {
	return connector
}
