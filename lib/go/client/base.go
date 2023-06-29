package client

import (
	"bytes"
	"context"
	"fmt"
	"log"

	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// BaseService is a type that is used for all client (call origin)
// side implementations.  This includes any client in the
// guest codespace. This object primarily just receives
// message call requests and this type sends it to the
// kernel.
type BaseService struct {
	svc   id.ServiceId
	smMap *lib.ServiceMethodMap
}

func NewBaseService(ctx context.Context, id id.ServiceId, sm *lib.ServiceMethodMap) *BaseService {
	if len(sm.Call()) == 0 {
		pcontext.Infof(ctx, "NewBaseService: binding size is zero")
	}
	return &BaseService{
		svc:   id,
		smMap: sm,
	}
}

func (c *BaseService) ServiceId() id.ServiceId {
	return c.svc
}

func (c *BaseService) ServiceMethodMap() *lib.ServiceMethodMap {
	return c.smMap
}

func (c *BaseService) MethodIdByName(str string) (id.MethodId, bool) {
	mid := c.smMap.MethodNameToId(c.ServiceId(), str)
	if mid.IsEmptyValue() {
		return mid, false
	}
	return mid, true
}

// String() returns a useful string for debugging a client side service.
// This includes all the known methods for the service.
func (c *BaseService) String() string {
	buf := &bytes.Buffer{}
	for _, pair := range c.smMap.Call() {
		mid := id.UnmarshalMethodId(pair.GetMethodId())
		buf.WriteString(fmt.Sprintf("%s:%s ", c.smMap.MethodIdToName(mid), mid.Short()))
	}
	msg := fmt.Sprintf("[BaseService: sid:%s, methods(%d):%s]", c.svc.Short(),
		c.smMap.Len(), buf.String())
	log.Println(msg)
	return msg
}

// Dispatch is called by every client side "method" on the client side
// service. This funciton is the one that make a system call to the
// kernel and prepares for handling the result.
func (c *BaseService) Dispatch(method id.MethodId, param proto.Message) (id.CallId, syscall.KernelErr) {
	var a *anypb.Any
	var err error
	if param != nil {
		a, err = anypb.New(param)
		if err != nil {
			ctx := pcontext.NewContextWithContainer(pcontext.GuestContext(context.Background()), "dispatch")
			pcontext.Errorf(ctx, "failed in call to dispatch: %v", err)
			pcontext.Dump(ctx)
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
		HostId:    lib.CurrentHostId().Marshal(),
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

func (c *BaseService) Launch() (*syscall.LaunchResponse, syscall.KernelErr) {
	req := syscall.LaunchRequest{
		ServiceId: c.svc.Marshal(),
	}
	out, err := syscallguest.Launch(&req)
	if err != 0 {
		return nil, err
	}
	return out, syscall.KernelErr_NoError
}
