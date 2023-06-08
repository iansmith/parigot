package lib

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/iansmith/parigot/apishared/id"
	syscallguest "github.com/iansmith/parigot/apiwasm/syscall"
	pcontext "github.com/iansmith/parigot/context"
	syscall "github.com/iansmith/parigot/g/syscall/v1"

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
func (c *ClientSideService) Dispatch(method string, param proto.Message) (*syscall.DispatchResponse, id.KernelErrId) {
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
	in := &syscall.DispatchRequest{
		ServiceId: c.svc.Marshal(),
		Caller:    c.caller,
		Method:    method,
		Param:     a,
	}
	return syscallguest.Dispatch(in)
}

func (c *ClientSideService) Run() (*syscall.RunResponse, id.KernelErrId) {
	req := syscall.RunRequest{
		Wait: true,
	}
	out, err := syscallguest.Run(&req)
	if err.IsError() {
		return nil, err
	}
	return out, id.KernelErrIdNoErr
}

// Require1 is a thin wrapper over syscall.Require so it's easy
// to require things by their name.  This is used by the code generator
// primarily.
func Require1(pkg, name string, source id.ServiceId) (*syscall.RequireResponse, id.KernelErrId) {
	fqs := &syscall.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}
	in := &syscall.RequireRequest{
		Dest:   []*syscall.FullyQualifiedService{fqs},
		Source: source.Marshal(),
	}
	resp, err := syscallguest.Require(in)
	return resp, err
}

// Export1 is a thin wrapper over syscall.Export so it's easy
// to export things by their name.  This is used by the code generator
// primarily.
func Export1(pkg, name string) (*syscall.ExportResponse, id.KernelErrId) {
	fqs := &syscall.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}
	in := &syscall.ExportRequest{
		Service: []*syscall.FullyQualifiedService{fqs},
	}
	resp, kerr := syscallguest.Export(in)
	return resp, kerr
}

// MustRegisterClient should be used by the "main" function of a client
// program that is not service itself, in other words it is a client only.
// If you are a service, you should use the automagically generated code
// MustRegister<BLAH>().
func MustRegisterClient(ctx context.Context) id.ServiceId {
	pkg := "Client"
	name := fmt.Sprintf("program%03d", rand.Intn(999))
	sid := register(ctx, pkg, name, true)
	pcontext.Debugf(ctx, "client faux service created: %s.%s", pkg, name)
	return sid
}

func register(ctx context.Context, pkg, name string, isClient bool) id.ServiceId {
	req := &syscall.RegisterRequest{}
	fqs := &syscall.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}
	req.Fqs = fqs
	req.IsClient = isClient
	resp, err := syscallguest.Register(req)
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
