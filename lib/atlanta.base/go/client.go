package lib

import (
	"github.com/iansmith/parigot/g/parigot/kernel"
	"github.com/iansmith/parigot/lib/interface_"
	"google.golang.org/protobuf/proto"
)

type ClientSideService struct {
	svc    interface_.Id
	caller string
}

func NewClientSideService(svc ServiceId) *ClientSideService {
	return &ClientSideService{
		svc: svc,
	}
}

func (c *ClientSideService) SetCaller(caller string) {
	c.caller = caller
}

func (c *ClientSideService) Dispatch(method string, in proto.Message, out *kernel.DispatchResponse) error {
	return c.DispatchFull(nil, method, c.caller, in, out)
}

func (c *ClientSideService) DispatchFull(ctx interface_.Pctx, method string, caller string, in proto.Message, out *kernel.DispatchResponse) error {

	if ctx == nil {
		ctx = NewPctx()
	}

	req := &kernel.DispatchRequest{
		ServiceSid: 0,
		Method:     method,
		Caller:     caller,
	}
	b, err := ctx.ToBytes()
	if err != nil {
		panic("unable to marshal protobuf pctx:" + err.Error())
	}
	req.InPctx = b
	b, err = proto.Marshal(in)
	if err != nil {
		panic("unable to marshal protobuf dispatch request:" + err.Error())
	}
	req.InBlob = b
	resp := kernel.DispatchResponse{}
	err = Dispatch(req, &resp)

	if err != nil {
		return NewPerrorFromError("internal error in dispatch", err)
	}
	return nil
}
