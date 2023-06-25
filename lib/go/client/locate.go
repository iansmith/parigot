package client

import (
	"context"

	"github.com/iansmith/parigot/apishared/id"
	syscallguest "github.com/iansmith/parigot/apiwasm/syscall"
	pcontext "github.com/iansmith/parigot/context"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
)

func LocateDynamic(ctx context.Context, protoPkg, serviceName string, calledBy id.ServiceId) (*BaseService, syscall.KernelErr) {
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
	smmap := lib.NewServiceMethodMap()
	for _, pair := range resp.GetBinding() {
		mid := id.UnmarshalMethodId(pair.MethodId)
		smmap.AddServiceMethod(serviceId, mid,
			serviceName, pair.MethodName, nil)
	}
	cs := NewBaseService(ctx, serviceId, smmap)
	return cs, syscall.KernelErr_NoError

}
