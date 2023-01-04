package lib

import (
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	pbsys "github.com/iansmith/parigot/api/proto/g/pb/syscall"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ClientSideService struct {
	call   Call
	svc    Id
	caller string
	pctx   *protosupport.Pctx
	logger Log
}

func NewClientSideService(id Id, caller string, logger Log, call Call) *ClientSideService {
	return &ClientSideService{
		svc:    id,
		caller: caller,
		pctx:   &protosupport.Pctx{},
		logger: logger,
		call:   call,
	}
}

func (c *ClientSideService) SetCaller(caller string) {
	c.caller = caller
}

func (c *ClientSideService) SetPctx(pctx *protosupport.Pctx) {
	c.pctx = pctx
}

// Shorthand to make it cleaner for the calls from a client side proxy.
func (c *ClientSideService) Dispatch(method string, param proto.Message) (*pbsys.DispatchResponse, error) {
	var a *anypb.Any
	var err error
	if param != nil {
		a, err = anypb.New(param)
		if err != nil {
			return nil, NewPerrorFromError("unable to convert param for dispatch into Any", err)
		}
	}
	if c.svc == nil {
		panic("cannot dispatch to a nil service! client side service field 'svc' is nil")
	}
	in := &pbsys.DispatchRequest{
		ServiceId: Marshal[protosupport.ServiceId](c.svc),
		Caller:    c.caller,
		Method:    method,
		InPctx:    c.pctx,
		Param:     a,
	}
	return c.call.Dispatch(in)
}

func (c *ClientSideService) Run() (*pbsys.RunResponse, error) {
	req := pbsys.RunRequest{
		Wait: true,
	}
	return c.call.Run(&req)
}
