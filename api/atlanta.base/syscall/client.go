package syscall

import (
	pbcall "github.com/iansmith/parigot/api/proto/g/pb/call"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	"github.com/iansmith/parigot/lib"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ClientSideService struct {
	svc    lib.Id
	caller string
	pctx   *protosupport.Pctx
}

func NewClientSideService(id lib.Id, caller string) *ClientSideService {
	return &ClientSideService{
		svc:    id,
		caller: caller,
		pctx:   &protosupport.Pctx{},
	}
}

func (c *ClientSideService) SetCaller(caller string) {
	c.caller = caller
}

func (c *ClientSideService) SetPctx(pctx *protosupport.Pctx) {
	c.pctx = pctx
}

// Shorthand to make it cleaner for the calls from a client side proxy.
func (c *ClientSideService) Dispatch(method string, param proto.Message) (*pbcall.DispatchResponse, error) {
	var a *anypb.Any
	var err error
	if param != nil {
		a, err = anypb.New(param)
		if err != nil {
			return nil, lib.NewPerrorFromError("unable to convert param for dispatch into Any", err)
		}
	}
	if c.svc == nil {
		panic("cannot dispatch to a nil service! client side service field 'svc' is nil")
	}
	in := &pbcall.DispatchRequest{
		ServiceId: lib.Marshal[protosupport.ServiceId](c.svc),
		Caller:    c.caller,
		Method:    method,
		InPctx:    c.pctx,
		Param:     a,
	}
	// xxx this should be going through dispatch anyway
	return CallConnection().Dispatch(in)
}
