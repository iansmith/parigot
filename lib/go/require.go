package lib

import (
	"context"
	"fmt"
	"math/rand"

	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
	"github.com/iansmith/parigot/api/shared/id"
	pcontext "github.com/iansmith/parigot/context"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
)

// Export1 is a thin wrapper over syscall.Export so it's easy
// to export things by their name.  This is used by the code generator
// primarily.
func Export1(pkg, name string, serviceId id.ServiceId) (*syscall.ExportResponse, syscall.KernelErr) {
	fqs := &syscall.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}
	in := &syscall.ExportRequest{
		ServiceId: serviceId.Marshal(),
		Service:   []*syscall.FullyQualifiedService{fqs},
		HostId:    CurrentHostId().Marshal(),
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

// MustInitClient is for clients only.  In other words, you should
// only use this function if you do not implement services, just
// use them.  A common case of this is a demo program or a program
// that performs a one off task.  This function wraps MustRegisterClient
// and panics if things go wrong.
func MustInitClient(ctx context.Context, requirement []MustRequireFunc) id.ServiceId {

	myId := MustRegisterClient(ctx)
	for _, f := range requirement {
		f(ctx, myId)
	}
	syscallguest.MustSatisfyWait(ctx, myId)

	launchreq := &syscall.LaunchRequest{
		ServiceId: myId.Marshal(),
	}
	_, err := syscallguest.Launch(launchreq)
	if err != syscall.KernelErr_NoError {
		panic(fmt.Sprintf("unable to launch client service: %s",
			syscall.KernelErr_name[int32(err)]))
	}

	return myId
}

func MustRunClient(ctx context.Context, timeoutInMillis int32) syscall.KernelErr {
	var err syscall.KernelErr
	for {
		err = clientOnlyReadOneAndCall(ctx, nil, timeoutInMillis)
		if err != syscall.KernelErr_NoError && err != syscall.KernelErr_ReadOneTimeout {
			break
		}
	}
	return err
}

func clientOnlyReadOneAndCall(ctx context.Context, binding *ServiceMethodMap,
	timeoutInMillis int32) syscall.KernelErr {
	req := syscall.ReadOneRequest{}

	// setup a request to read an incoming message
	req.TimeoutInMillis = timeoutInMillis
	req.HostId = CurrentHostId().Marshal()
	resp, err := syscallguest.ReadOne(&req)
	if err != syscall.KernelErr_NoError {
		return err
	}
	// is timeout?
	if resp.Timeout {
		return syscall.KernelErr_ReadOneTimeout
	}

	// check for finished futures from within our address space
	ExpireMethod(ctx)

	// is a promise being completed that was fulfilled somewhere else
	if r := resp.GetResolved(); r != nil {
		cid := id.UnmarshalCallId(r.GetCallId())
		CompleteCall(ctx, cid, r.GetResult(), r.GetResultError())
		return syscall.KernelErr_NoError
	}

	return syscall.KernelErr_NoError
}
