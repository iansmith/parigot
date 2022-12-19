package lib

import (
	"github.com/iansmith/parigot/api/proto/g/pb/call"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ClientSideService struct {
	svc    Id
	caller string
	pctx   *protosupport.Pctx
}

func NewClientSideService(id Id, caller string) *ClientSideService {
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
func (c *ClientSideService) Dispatch(method string, param proto.Message) (*call.DispatchResponse, error) {
	var a *anypb.Any
	var err error
	if param != nil {
		a, err = anypb.New(param)
		if err != nil {
			return nil, NewPerrorFromError("unable to convert param for dispatch into Any", err)
		}
	}

	in := &call.DispatchRequest{
		ServiceId: Marshal[protosupport.ServiceId](c.svc),
		Caller:    c.caller,
		Method:    method,
		InPctx:    c.pctx,
		Param:     a,
	}
	// xxx this should be going through dispatch anyway
	return CallConnection().Dispatch(in)
}
