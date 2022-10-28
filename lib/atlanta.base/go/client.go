package lib

import (
	"github.com/iansmith/parigot/g/pb/kernel"
	"google.golang.org/protobuf/proto"
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

func (c *ClientSideService) Dispatch(method string, in proto.Message, out *kernel.DispatchResponse) error {
	return c.DispatchFull(nil, method, c.caller, in, out)
}

func (c *ClientSideService) DispatchFull(ctx Pctx, method string, caller string, in proto.Message, out *kernel.DispatchResponse) error {

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
	Dispatch(req, &resp)

	//if err != nil {
	//	return NewPerrorFromError("internal error in dispatch", err)
	//}
	return nil
}
