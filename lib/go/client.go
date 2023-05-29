package lib

import (
	"context"

	"github.com/iansmith/parigot/apishared/id"
	"github.com/iansmith/parigot/apiwasm/syscall"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ClientSideService struct {
	svc    id.ServiceId
	caller string
}

func NewClientSideService(ctx context.Context, id id.ServiceId, caller string) *ClientSideService {
	return &ClientSideService{
		svc:    id,
		caller: caller,
	}
}

func (c *ClientSideService) SetCaller(caller string) {
	c.caller = caller
}

// Shorthand to make it cleaner for the calls from a client side proxy.
func (c *ClientSideService) Dispatch(method string, param proto.Message) (*syscallmsg.DispatchResponse, id.KernelErrId) {
	var a *anypb.Any
	var err error
	if param != nil {
		a, err = anypb.New(param)
		if err != nil {
			return nil, id.KernelErrIdNoErr
		}
	}
	if c.svc.IsZeroValue() {
		panic("cannot dispatch to a nil service! client side service field 'svc' is nil")
	}
	in := &syscallmsg.DispatchRequest{
		ServiceId: c.svc.Marshal(),
		Caller:    c.caller,
		Method:    method,
		InPctx:    nil,
		Param:     a,
	}
	return syscall.Dispatch(in)
}

func (c *ClientSideService) Run() (*syscallmsg.RunResponse, id.KernelErrId) {
	req := syscallmsg.RunRequest{
		Wait: true,
	}
	out, err := syscall.Run(&req)
	return out, id.NewKernelErrIdFromRaw(err)
}

// Require1 is a thin wrapper over syscall.Require so it's easy
// to require things by their name.  This is used by the code generator
// primarily.
func Require1(pkg, name string) (*syscallmsg.RequireResponse, id.Id) {
	fqs := &syscallmsg.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}
	in := &syscallmsg.RequireRequest{
		Service: []*syscallmsg.FullyQualifiedService{fqs},
	}
	return syscall.Require(in)
}

// Export1 is a thin wrapper over syscall.Export so it's easy
// to export things by their name.  This is used by the code generator
// primarily.
func Export1(pkg, name string) (*syscallmsg.ExportResponse, id.Id) {
	fqs := &syscallmsg.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}
	in := &syscallmsg.ExportRequest{
		Service: []*syscallmsg.FullyQualifiedService{fqs},
	}
	return syscall.Export(in)
}
