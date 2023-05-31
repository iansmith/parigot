package lib

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/iansmith/parigot/apishared/id"
	"github.com/iansmith/parigot/apiwasm/syscall"
	pcontext "github.com/iansmith/parigot/context"
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
	if err.IsError() {
		return nil, err
	}
	return out, id.KernelErrIdNoErr
}

// Require1 is a thin wrapper over syscall.Require so it's easy
// to require things by their name.  This is used by the code generator
// primarily.
func Require1(pkg, name string, source id.ServiceId) (*syscallmsg.RequireResponse, id.KernelErrId) {
	fqs := &syscallmsg.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}
	in := &syscallmsg.RequireRequest{
		Dest:   []*syscallmsg.FullyQualifiedService{fqs},
		Source: source.Marshal(),
	}
	resp, err := syscall.Require(in)
	return resp, err
}

// Export1 is a thin wrapper over syscall.Export so it's easy
// to export things by their name.  This is used by the code generator
// primarily.
func Export1(pkg, name string) (*syscallmsg.ExportResponse, id.KernelErrId) {
	fqs := &syscallmsg.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}
	in := &syscallmsg.ExportRequest{
		Service: []*syscallmsg.FullyQualifiedService{fqs},
	}
	return syscall.Export(in)
}

// MustRegisterClient should be used by the "main" function of a client
// program that is not service itself, in other words it is a client only.
// If you are a service, you should use the automagically generated code
// MustRegister<BLAH>().
func MustRegisterClient(ctx context.Context) id.ServiceId {
	pkg := "Client"
	name := fmt.Sprintf("program%03d", rand.Intn(999))
	sid := register(ctx, pkg, name, true)
	return sid
}

func register(ctx context.Context, pkg, name string, isClient bool) id.ServiceId {
	req := &syscallmsg.RegisterRequest{}
	fqs := &syscallmsg.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}
	req.Fqs = fqs
	req.IsClient = isClient
	resp, err := syscall.Register(req)
	if err.IsError() {
		pcontext.Fatalf(ctx, "unable to register %s.%s: %s", pkg, name, err.Short())
		panic("registration error")
	}
	sid, errId := id.UnmarshalServiceId(resp.GetId())
	if errId.IsError() {
		pcontext.Fatalf(ctx, "unable to unmarshal service id for %s.%s: %s", pkg, name, err.Short())
		panic("registration error")
	}
	return sid

}
