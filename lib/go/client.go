package lib

import (
	"context"

	"github.com/iansmith/parigot/apiwasm/syscall"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ClientSideService struct {
	svc    Id
	caller string
}

func NewClientSideService(ctx context.Context, id Id, caller string) *ClientSideService {
	return &ClientSideService{
		svc:    id,
		caller: caller,
	}
}

func (c *ClientSideService) SetCaller(caller string) {
	c.caller = caller
}

// Shorthand to make it cleaner for the calls from a client side proxy.
func (c *ClientSideService) Dispatch(method string, param proto.Message) (*syscallmsg.DispatchResponse, error) {
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
	in := &syscallmsg.DispatchRequest{
		ServiceId: Marshal[protosupportmsg.ServiceId](c.svc),
		Caller:    c.caller,
		Method:    method,
		InPctx:    nil,
		Param:     a,
	}
	out := syscall.Dispatch(in)
	// xxx should be checking for error value
	return out, nil
}

func (c *ClientSideService) Run() (*syscallmsg.RunResponse, error) {
	req := syscallmsg.RunRequest{
		Wait: true,
	}
	out := syscall.Run(&req)
	// xxx should be checking for error value
	return out, nil
}

// Require1 is a thin wrapper over syscall.Require so it's easy
// to require things by their name.  This is used by the code generator
// primarily.
func Require1(pkg, name string) (*syscallmsg.RequireResponse, error) {
	fqs := &syscallmsg.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}
	in := &syscallmsg.RequireRequest{
		Service: []*syscallmsg.FullyQualifiedService{fqs},
	}
	out := syscall.Require(in)
	return out, nil
}

// Export1 is a thin wrapper over syscall.Export so it's easy
// to export things by their name.  This is used by the code generator
// primarily.
func Export1(pkg, name string) (*syscallmsg.ExportResponse, error) {
	fqs := &syscallmsg.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}
	in := &syscallmsg.ExportRequest{
		Service: []*syscallmsg.FullyQualifiedService{fqs},
	}
	out := syscall.Export(in)
	return out, nil
}
