package lib

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/iansmith/parigot/apishared/id"
	"github.com/iansmith/parigot/apiwasm"
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
	svc   id.ServiceId
	smMap *apiwasm.ServiceMethodMap
	cont  map[string]Promise[proto.Message, int32]
}

func NewClientSideService(ctx context.Context, id id.ServiceId, sm *apiwasm.ServiceMethodMap) *ClientSideService {
	if len(sm.Call()) == 0 {
		log.Printf("NewClientSideService: binding is zero")
		//debug.PrintStack()
	}
	return &ClientSideService{
		svc:   id,
		smMap: sm,
		cont:  make(map[string]Promise[proto.Message, int32]),
	}
}

func (c *ClientSideService) ServiceId() id.ServiceId {
	return c.svc
}

func (c *ClientSideService) ServiceMethodMap() *apiwasm.ServiceMethodMap {
	return c.smMap
}
func (c *ClientSideService) Continuation(cid id.CallId, fn func(*anypb.Any, int32) syscall.KernelErr) {
	//c.contOut[cid.String()] = fn
}

// Complete call is used to connect the return results of a dispatch call
// the proper continues (success or failure).  This function returns
// a value that is ONLY regarding the setup/teardown behavior, not that of
// the called continuations.
func (c *ClientSideService) CompleteCall(cid id.CallId, a *anypb.Any, err int32) syscall.KernelErr {
	// fn, ok := c.cont[cid.String()]
	// if !ok {
	// 	return syscall.KernelErr_NotFound
	// }
	// e := fn(a, err)
	// delete(c.cont, cid.String())
	return syscall.KernelErr_NoError
}

// String() returns a useful stringn for debugging a client side service.
// This includes all the known methods for the service.
func (c *ClientSideService) String() string {
	buf := &bytes.Buffer{}
	for _, pair := range c.smMap.Call() {
		mid := id.UnmarshalMethodId(pair.GetMethodId())
		buf.WriteString(fmt.Sprintf("%s:%s ", c.smMap.MethodIdToName(mid), mid.Short()))
	}
	msg := fmt.Sprintf("[clientSideService: sid:%s, methods(%d):%s]", c.svc.Short(),
		c.smMap.Len(), buf.String())
	log.Println(msg)
	return msg
}

// Dispatch is called by every client side "method" on the client side
// service. This funciton is the one that make a system call to the
// kernel and prepares for handling the result.
func (c *ClientSideService) Dispatch(method id.MethodId, param proto.Message) (id.CallId, syscall.KernelErr) {
	var a *anypb.Any
	var err error
	if param != nil {
		a, err = anypb.New(param)
		if err != nil {
			// do we want to have a special error type for this?
			return id.CallIdZeroValue(), syscall.KernelErr_MarshalFailed
		}
	}
	if c.svc.IsZeroOrEmptyValue() {
		panic("cannot dispatch to an unknown service! client side service field 'svc' is zero or empty")
	}
	if method.IsZeroOrEmptyValue() {
		panic("cannot dispatch to an unknown method! client side service field 'method id' is zero or empty")
	}
	// this is where it all begins
	cid := id.NewCallId()
	in := &syscall.DispatchRequest{
		ServiceId: c.svc.Marshal(),
		MethodId:  method.Marshal(),
		CallId:    cid.Marshal(),
		Param:     a,
		HostId:    CurrentHostId().Marshal(),
	}
	resp, kerr := syscallguest.Dispatch(in)
	if kerr != syscall.KernelErr_NoError {
		return id.CallIdZeroValue(), kerr
	}
	cid2 := id.UnmarshalCallId(resp.GetCallId())
	if !cid.Equal(cid2) {
		panic("mismatched call ids in dispatch")
	}
	return cid, syscall.KernelErr_NoError
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
		HostId:  CurrentHostId().Marshal(),
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

func LocateDynamic(ctx context.Context, protoPkg, serviceName string, calledBy id.ServiceId) (*ClientSideService, syscall.KernelErr) {
	req := &syscall.LocateRequest{
		PackageName: protoPkg,
		ServiceName: serviceName,
		CalledBy:    calledBy.Marshal(),
	}

	resp, kerr := syscallguest.Locate(req)
	if kerr != syscall.KernelErr_NoError {
		pcontext.Errorf(ctx, "UnmarshalServiceId failed: %s", syscall.KernelErr_name[int32(kerr)])
		return nil, kerr
	}
	serviceId := id.UnmarshalServiceId(resp.GetServiceId())
	smmap := apiwasm.NewServiceMethodMap()
	for _, pair := range resp.GetBinding() {
		mid := id.UnmarshalMethodId(pair.MethodId)
		smmap.AddServiceMethod(serviceId, mid,
			serviceName, pair.MethodName, nil)
	}
	cs := NewClientSideService(ctx, serviceId, smmap)
	return cs, syscall.KernelErr_NoError

}
