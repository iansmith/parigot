package client

import (
	"context"

	"github.com/iansmith/parigot/api/guest"
	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
	"github.com/iansmith/parigot/api/shared/id"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
)

// LocateDynamic is an important interface to the infrastructure's knowlege
// about the types and methods of services.  This method takes a package and
// a service name and returns a client side proxy that can call all the methods
// that the system knows about for the given service. This is usually used
// by the code generate to then wrap method declarations around to give the
// resulting client object an API that is more what is expected.  However, it
// can be used "uncooked" and its methods called via Dispatch().  This is
// the only mechanism by which you can call methods on services that are
// not known at compile time, e.g. looking up a service by package and name,
// then and calling methods on the service entirely by using strings.
func LocateDynamic(ctx context.Context, protoPkg, serviceName string, calledBy id.ServiceId) (*BaseService, syscall.KernelErr) {
	req := &syscall.LocateRequest{
		PackageName: protoPkg,
		ServiceName: serviceName,
		CalledBy:    calledBy.Marshal(),
	}

	resp, kerr := syscallguest.Locate(ctx, req)
	if kerr != syscall.KernelErr_NoError {
		guest.Log(ctx).Error("UnmarshalServiceId failed", "kernel error", syscall.KernelErr_name[int32(kerr)])
		return nil, kerr
	}
	serviceId := id.UnmarshalServiceId(resp.GetServiceId())
	smmap := lib.NewServiceMethodMap()

	for _, pair := range resp.GetBinding() {
		mid := id.UnmarshalMethodId(pair.MethodId)
		smmap.AddServiceMethod(serviceId, mid,
			serviceName, pair.MethodName, nil)
	}
	cs := NewBaseService(ctx, serviceId, smmap)
	return cs, syscall.KernelErr_NoError

}
