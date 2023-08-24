package lib

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"

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
	return context.Background(), sid
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

// ExitSelf sends a request to exit and when the future is completed it will
// exit.  Valid code values must be in the range 0 to 192, inclusive.  If ExitSelf()
// is called by the last running guest program then the entire process will exit.
func ExitSelf(ctx context.Context, code int32, myId id.ServiceId) *syscallguest.ExitFuture {
	return exitImpl(ctx, code, myId, false)
}

// ExitAll sends a request to exit on all hosts and when the future is completed, it
// will exit.  Valid code values must be in the range 0 to 192, inclusive.
func ExitAll(ctx context.Context, code int32, myId id.ServiceId) *syscallguest.ExitFuture {
	return exitImpl(ctx, code, myId, true)
}

func exitImpl(ctx context.Context, code int32, myId id.ServiceId, all bool) *syscallguest.ExitFuture {
	hid := syscallguest.CurrentHostId()
	cid := id.NewCallId()
	req := &syscall.ExitRequest{
		HostId: hid.Marshal(),
		Pair: &syscall.ExitPair{
			ServiceId: myId.Marshal(),
			Code:      code,
		},
		CallId:      cid.Marshal(),
		MethodId:    apishared.ExitMethod.Marshal(),
		ShutdownAll: all,
	}
	return syscallguest.Exit(ctx, req)
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
