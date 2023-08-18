package lib

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
	apishared "github.com/iansmith/parigot/api/shared"
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
		HostId:    syscallguest.CurrentHostId().Marshal(),
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
	req.HostId = syscallguest.CurrentHostId().Marshal()
	req.DebugName = fmt.Sprintf("%s.%s", pkg, name)
	resp, err := syscallguest.Register(req)
	if err != 0 {
		pcontext.Fatalf(ctx, "unable to register %s.%s: %s", pkg, name,
			syscall.KernelErr_name[int32(err)])
		panic("registration error")
	}
	sid := id.UnmarshalServiceId(resp.GetServiceId())
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
	return myId
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
	return syscallguest.Launch(req)
}

// ExitClient sends a request to exit and attaches hanndlers that print
// the given strings. It only forces the exit if the Exit call itself
// fails. Only the values from 0 to 192 are permissable as the code; other
// values will be changed to 192.
func ExitClient(ctx context.Context, code int32, myId id.ServiceId, msgSuccess, msgFailure string) {
	req := &syscall.ExitRequest{
		HostId: syscallguest.CurrentHostId().Marshal(),
		Pair: &syscall.ExitPair{
			ServiceId: myId.Marshal(),
			Code:      code,
		},
		CallId:   id.NewCallId().Marshal(),
		MethodId: apishared.ExitMethod.Marshal(),
	}
	syscallguest.Exit(req)
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
	resp, err := syscallguest.ReadOne(&req)
	if err != syscall.KernelErr_NoError {
		return err
	}
	// is timeout?
	if resp.Timeout {
		return syscall.KernelErr_ReadOneTimeout
	}

	if resp.Exit != nil && resp.Exit.ServiceId != nil {
		sid := id.UnmarshalServiceId(resp.Exit.GetServiceId())
		log.Printf("xxx %#v\n xxx -- %#v,,, %s", resp.Bundle, resp.Exit, sid.Short())
		panic(apishared.ControlledExit)
	}

	// check for finished futures from within our address space
	syscallguest.ExpireMethod(ctx)

	// is a promise being completed that was fulfilled somewhere else
	if r := resp.GetResolved(); r != nil {
		cid := id.UnmarshalCallId(r.GetCallId())
		syscallguest.CompleteCall(ctx, syscallguest.CurrentHostId(), cid, r.GetResult(), r.GetResultError())
		return syscall.KernelErr_NoError
	}

	return syscall.KernelErr_NoError
}
