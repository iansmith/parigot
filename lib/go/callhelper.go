package lib

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/iansmith/parigot/api/guest"
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/api/shared/id"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
)

// Export1 is a thin wrapper over syscall.Export so it's easy
// to export things by their name.  This is used by the code generator
// primarily.
func Export1(ctx context.Context, pkg, name string, serviceId id.ServiceId) (*syscall.ExportResponse, syscall.KernelErr) {
	fqs := &syscall.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}
	in := &syscall.ExportRequest{
		ServiceId: serviceId.Marshal(),
		Service:   []*syscall.FullyQualifiedService{fqs},
		HostId:    syscallguest.CurrentHostId().Marshal(),
	}
	resp, kerr := syscallguest.Export(ctx, in)
	return resp, kerr
}

// MustRegisterClient should be used by the "main" function of a client
// program that is not service itself, in other words it is a client only.
// If you are a service, you should use the automagically generated code
// MustRegister<BLAH>().
func MustRegisterClient() (context.Context, id.ServiceId) {
	pkg := "Client"
	name := fmt.Sprintf("program%03d", rand.Intn(999))
	sid := register(context.Background(), pkg, name, true)
	ctx := guest.NewContextWithLogger(sid)
	return ctx, sid
}

// register is the setup of the guest-side registration. It uses
// syscall lib to actually make the call.
func register(ctx context.Context, pkg, name string, isClient bool) id.ServiceId {
	req := &syscall.RegisterRequest{}
	req.HostId = syscallguest.CurrentHostId().Marshal()
	req.DebugName = fmt.Sprintf("%s.%s", pkg, name)
	resp, err := syscallguest.Register(ctx, req)
	if err != 0 {
		panic("registration error")
	}
	sid := id.UnmarshalServiceId(resp.GetServiceId())
	return sid

}

// Require1 is a thin wrapper over syscall.Require so it's easy
// to require things by their name.  This is used by the code generator
// primarily.
func Require1(ctx context.Context, pkg, name string, source id.ServiceId) (*syscall.RequireResponse, syscall.KernelErr) {
	fqs := &syscall.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}
	in := &syscall.RequireRequest{
		Dest:   []*syscall.FullyQualifiedService{fqs},
		Source: source.Marshal(),
	}
	resp, err := syscallguest.Require(ctx, in)
	return resp, err
}

// MustInitClient is for clients only.  In other words, you should
// only use this function if you do not implement services, just
// use them.  A common case of this is a demo program or a program
// that performs a one off task.  This function wraps MustRegisterClient
// and panics if things go wrong.
func MustInitClient(requirement []MustRequireFunc) (context.Context, id.ServiceId) {

	ctx, myId := MustRegisterClient()
	for _, f := range requirement {
		f(ctx, myId)
	}
	return ctx, myId
}

// LaunchClient is a convienence wrapper around Launch() for clients that don't
// want to create their own request structure.
func LaunchClient(ctx context.Context, myId id.ServiceId) *syscallguest.LaunchFuture {
	cid := id.NewCallId()
	req := &syscall.LaunchRequest{
		HostId:    syscallguest.CurrentHostId().Marshal(),
		ServiceId: myId.Marshal(),
		CallId:    cid.Marshal(),
		MethodId:  apishared.LaunchMethod.Marshal(),
	}
	return syscallguest.Launch(ctx, req)
}

// ExitClient sends a request to exit and attaches hanndlers that print
// the given strings. It only forces the exit if the Exit call itself
// fails. Only the values from 0 to 192 are permissable as the code; other
// values will be changed to 192.
func ExitClient(ctx context.Context, code int32, myId id.ServiceId) syscall.KernelErr {
	hid := syscallguest.CurrentHostId()
	cid := id.NewCallId()
	log.Printf("xxx -- ExitClient reached on %s, exit call is %s", hid.Short(), cid.Short())
	req := &syscall.ExitRequest{
		HostId: hid.Marshal(),
		Pair: &syscall.ExitPair{
			ServiceId: myId.Marshal(),
			Code:      code,
		},
		CallId:   cid.Marshal(),
		MethodId: apishared.ExitMethod.Marshal(),
	}
	syscallguest.Exit(ctx, req)
	log.Printf("xxx -- finished call to Exit(), host %s, cid %s", hid.Short(), cid.Short())
	// ctx, t := CurrentTime(ctx)
	// syscallguest.MatchCompleter(ctx, t, hid, cid, &syscallguest.ExitCompleter{})
	// resp := &syscall.ExitResponse{
	// 	Pair: &syscall.ExitPair{
	// 		ServiceId: myId.Marshal(),
	// 		Code:      code,
	// 	},
	// }
	// a := &anypb.Any{}
	// if err := a.MarshalFrom(resp); err != nil {
	// 	return syscall.KernelErr_MarshalFailed
	// }
	// log.Printf("xxx completing call on host %s, with host,cid of %s,%s",
	// 	syscallguest.CurrentHostId().Short(), hid.Short(), cid.Short())
	// syscallguest.CompleteCall(ctx, hid, cid, a, int32(syscall.KernelErr_NoError))
	return syscall.KernelErr_NoError
}

func MustRunClient(ctx context.Context, timeoutInMillis int32) syscall.KernelErr {
	var err syscall.KernelErr
	for {
		err = ReadOneAndCallClient(ctx, nil, timeoutInMillis)
		if err != syscall.KernelErr_NoError && err != syscall.KernelErr_ReadOneTimeout {
			break
		}
	}
	return err
}

// ReadOneAndCallClient does the waiting for an incoming call and if one
// arrives, it dispatches the call to the appropriate method.  Similarly, it
// will detect and respond to finished futures.  It returns KernelErr_ReadOneTimeout
// if the waiting timed out, otherwise the value should be KernelErr_NoError or
// an appropriate error code.
func ReadOneAndCallClient(ctx context.Context, binding *ServiceMethodMap,
	timeoutInMillis int32) syscall.KernelErr {
	req := syscall.ReadOneRequest{}

	// setup a request to read an incoming message
	req.TimeoutInMillis = timeoutInMillis
	req.HostId = syscallguest.CurrentHostId().Marshal()
	resp, err := syscallguest.ReadOne(ctx, &req)
	if err != syscall.KernelErr_NoError {
		return err
	}
	// is timeout?
	if resp.Timeout {
		return syscall.KernelErr_ReadOneTimeout
	}

	if resp.Exit != nil && resp.Exit.ServiceId != nil {
		//sid := id.UnmarshalServiceId(resp.Exit.GetServiceId())
		os.Exit(int(resp.Exit.GetCode()))
	}

	// check for finished futures from within our address space
	ctx, t := CurrentTime(ctx)
	syscallguest.ExpireMethod(ctx, t)

	// is a promise being completed that was fulfilled somewhere else
	if r := resp.GetResolved(); r != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("trapped xxx panic ---'%v' on %s", r, syscallguest.CurrentHostId())
			}
		}()
		cid := id.UnmarshalCallId(r.GetCallId())
		syscallguest.CompleteCall(ctx, syscallguest.CurrentHostId(), cid, r.GetResult(), r.GetResultError())
		return syscall.KernelErr_NoError
	}

	return syscall.KernelErr_NoError
}
