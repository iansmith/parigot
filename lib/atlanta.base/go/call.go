package lib

import (
	pbcall "github.com/iansmith/parigot/api/proto/g/pb/call"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"

	"google.golang.org/protobuf/proto"
)

type Call interface {
	Exit(in *pbcall.ExitRequest)
	Locate(in *pbcall.LocateRequest) (*pbcall.LocateResponse, error)
	Dispatch(in *pbcall.DispatchRequest) (*pbcall.DispatchResponse, error)
	BindMethodIn(in *pbcall.BindMethodRequest, _ func(*protosupport.Pctx, proto.Message) error) (*pbcall.BindMethodResponse, error)
	BindMethodOut(in *pbcall.BindMethodRequest, _ func(*protosupport.Pctx) (proto.Message, error)) (*pbcall.BindMethodResponse, error)
	BindMethodBoth(in *pbcall.BindMethodRequest, _ func(*protosupport.Pctx, proto.Message) (proto.Message, error)) (*pbcall.BindMethodResponse, error)
	Run(in *pbcall.RunRequest) (*pbcall.RunResponse, error)
	Export(in *pbcall.ExportRequest) (*pbcall.ExportResponse, error)
	Require(in *pbcall.RequireRequest) (*pbcall.RequireResponse, error)
	BlockUntilCall(in *pbcall.BlockUntilCallRequest) (*pbcall.BlockUntilCallResponse, error)
	ReturnValue(in *pbcall.ReturnValueRequest) (*pbcall.ReturnValueResponse, error)
}

var connector Call = newCallImpl()

func CallConnection() Call {
	return connector
}
