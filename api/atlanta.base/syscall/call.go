package syscall

import (
	pbcall "github.com/iansmith/parigot/api/proto/g/pb/call"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	pbsys "github.com/iansmith/parigot/api/proto/g/pb/syscall"

	"google.golang.org/protobuf/proto"
)

// Call is an interface that defines the set of calls that WASM code make into the kernel.  Note that
// look similar to the RPC but the mechanism here is the same as a split service.  These calls are
// always made from WASM code and primarily form code generated by our protoc plugin which wraps these
// nicely for the appropriate language.
type Call interface {
	Exit(in *pbcall.ExitRequest)
	Locate(in *pbsys.LocateRequest) (*pbsys.LocateResponse, error)
	Dispatch(in *pbsys.DispatchRequest) (*pbsys.DispatchResponse, error)
	BindMethodIn(in *pbcall.BindMethodRequest, _ func(*protosupport.Pctx, proto.Message) error) (*pbcall.BindMethodResponse, error)
	BindMethodOut(in *pbcall.BindMethodRequest, _ func(*protosupport.Pctx) (proto.Message, error)) (*pbcall.BindMethodResponse, error)
	BindMethodBoth(in *pbcall.BindMethodRequest, _ func(*protosupport.Pctx, proto.Message) (proto.Message, error)) (*pbcall.BindMethodResponse, error)
	Run(in *pbcall.RunRequest) (*pbcall.RunResponse, error)
	Export(in *pbcall.ExportRequest) (*pbcall.ExportResponse, error)
	Require(in *pbcall.RequireRequest) (*pbcall.RequireResponse, error)
	BlockUntilCall(in *pbcall.BlockUntilCallRequest) (*pbcall.BlockUntilCallResponse, error)
	ReturnValue(in *pbcall.ReturnValueRequest) (*pbcall.ReturnValueResponse, error)
}

var connector Call

func CallConnection() Call {
	if connector == nil {
		connector = newCallImpl()
	}
	return connector
}
