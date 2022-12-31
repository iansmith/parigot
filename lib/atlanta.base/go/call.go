package lib

import (
	pbcall "github.com/iansmith/parigot/api/proto/g/pb/call"
	pblog "github.com/iansmith/parigot/api/proto/g/pb/log"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	pbsys "github.com/iansmith/parigot/api/proto/g/pb/syscall"

	"google.golang.org/protobuf/proto"
)

// Call is an interface that defines the set of calls that WASM code make into the kernel.  Note that
// look similar to the RPC but the mechanism here is the same as a split service.  These calls are
// always made from WASM code and primarily from code generated by our protoc plugin which wraps these
// nicely for the appropriate language.
//
// Call has an analogue called SysCall that is the means by which the kernel receives these calls.
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
	BlockUntilCall(in *pbsys.BlockUntilCallRequest) (*pbsys.BlockUntilCallResponse, error)
	ReturnValue(in *pbsys.ReturnValueRequest) (*pbsys.ReturnValueResponse, error)
	Export1(pkg, name string) (*pbcall.ExportResponse, error)
	Require1(pkg, name string) (*pbcall.RequireResponse, error)

	// Use of this function is discouraged. This function uses a backdoor to reach the logging service
	// and does not go through the normal LocateLog() process that can allow better/different implementation
	// of said service.  This is intended only for debugging when inside parigot's implementation.
	BackdoorLog(in *pblog.LogRequest) (*pblog.LogResponse, error)
}
