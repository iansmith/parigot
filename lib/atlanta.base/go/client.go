package lib

import (
	"github.com/iansmith/parigot/g/pb/kernel"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ClientSideService struct {
	svc    Id
	caller string
}

func NewClientSideService(id Id) *ClientSideService {
	return &ClientSideService{
		svc: id,
	}
}

func (c *ClientSideService) SetCaller(caller string) {
	c.caller = caller
}

// Shorthand to make it cleaner for the calls from a client side proxy.
func (c *ClientSideService) Dispatch(method string, param, result proto.Message) (*kernel.DispatchResponse, error) {
	a, err := anypb.New(param)
	if err != nil {
		return nil, NewPerrorFromError("unable to convert param for dispatch into Any", err)
	}
	in := &kernel.DispatchRequest{
		ServiceId: MarshalServiceId(c.svc),
		Caller:    c.caller,
		Method:    method,
		InPctx:    nil,
		Param:     a,
	}
	// xxx this should be going through dispatch anyway
	return Dispatch(in)
}
