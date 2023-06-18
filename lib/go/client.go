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

// ClientSideService is a type that is used for all client (call origin)
// side implementations.  This includes any client in the
// guest codespace. This object primarily just receives
// message call requests and this type sends it to the
// kernel.
type ClientSideService struct {
	svc id.ServiceId
}

func NewClientSideService(ctx context.Context, id id.ServiceId) *ClientSideService {
	return &ClientSideService{
		svc: id,
	}
}

func (c *ClientSideService) ServiceId() id.ServiceId {
	return c.svc
}

// Shorthand to make it cleaner for the calls from a client side proxy.
func (c *ClientSideService) Dispatch(method id.MethodId, param proto.Message) (*syscall.DispatchResponse, syscall.KernelErr) {
	var a *anypb.Any
	var err error
	if param != nil {
		a, err = anypb.New(param)
		if err != nil {
			// do we want to have a special error type for this?
			return nil, syscall.KernelErr_MarshalFailed
		}
	}
	if c.svc.IsZeroOrEmptyValue() {
		panic("cannot dispatch to an unknown service! client side service field 'svc' is zero or empty")
	}
	if method.IsZeroOrEmptyValue() {
		print("!!!!!! METHOD IS ZERO !!!!!!\n")
		panic("cannot dispatch to an unknown method! client side service field 'method id' is zero or empty")
	}
	in := &syscall.DispatchRequest{
		ServiceId: c.svc.Marshal(),
		MethodId:  method.Marshal(),
		Param:     a,
	}
	return syscallguest.Dispatch(in)
}

func (c *ClientSideService) Launch() (*syscall.LaunchResponse, syscall.KernelErr) {
	req := syscall.LaunchRequest{
		ServiceId: c.svc.Marshal(),
	}
	out, err := syscallguest.Launch(&req)
	if err != 0 {
		return nil, err
	}
	return out, syscall.KernelErr_NoError
}

// Require1 is a thin wrapper over syscall.Require so it's easy
// to require things by their name.  This is used by the code generator
// primarily.
func Require1(pkg, name string, source id.ServiceId) (*syscall.RequireResponse, syscall.KernelErr) {
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
func Export1(pkg, name string) (*syscall.ExportResponse, syscall.KernelErr) {
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
	if err != 0 {
		pcontext.Fatalf(ctx, "unable to register %s.%s: %s", pkg, name,
			syscall.KernelErr_name[int32(err)])
		panic("registration error")
	}
	sid := id.UnmarshalServiceId(resp.GetId())
	return sid

}
