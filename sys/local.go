package sys

import (
	"context"

	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/sys/dep"
)

type localSysCall struct {
	nameServer *LocalNameServer
}

// sharedBind takes an NSCore, not a nameserver.
func sharedBind(core *NSCore, p *Process, packagePath, service, method string) (lib.Id, lib.Id) {
	key := NewDepKeyFromProcess(p)
	mid := core.FindOrCreateMethodId(key, packagePath, service, method)
	return mid, nil
}

func sharedExport(ctx context.Context, ns NameServer, key dep.DepKey, pkg, service string) lib.Id {
	if kerr := ns.CloseService(ctx, key, pkg, service); kerr != nil && kerr.IsError() {
		return kerr
	}
	if kerr := ns.Export(ctx, key, pkg, service); kerr != nil && kerr.IsError() {
		return kerr
	}
	return nil
}

//	func sharedHandleMethod(ns NameServer, proc *Process, pkg, service, method string) (lib.Id, lib.Id) {
//		return ns.HandleMethod(proc, pkg, service, method)
//	}
func sharedRequire(ctx context.Context, ns NameServer, key dep.DepKey, pkg, service string) lib.Id {
	return ns.Require(ctx, key, pkg, service)
}
func sharedGetService(ctx context.Context, ns NameServer, key dep.DepKey, pkg, service string) (lib.Id, lib.Id, string) {
	return ns.GetService(ctx, key, pkg, service)
}
func sharedCallService(ctx context.Context, ns NameServer, key dep.DepKey, info *callContext) (*syscallmsg.ReturnValueRequest, lib.Id, string) {
	return ns.CallService(ctx, key, info)
}

func newLocalSysCall(ns *LocalNameServer) *localSysCall {
	return &localSysCall{nameServer: ns}
}

func (l *localSysCall) Bind(ctx context.Context, proc *Process, packagePath, service, method string) (lib.Id, lib.Id) {
	return sharedBind(ctx, l.nameServer.NSCore, proc, packagePath, service, method)
}

func (l *localSysCall) Export(ctx context.Context, key dep.DepKey, pkg, service string) lib.Id {
	return sharedExport(ctx, l.nameServer, key, pkg, service)
}

func (l *localSysCall) CallService(ctx context.Context, key dep.DepKey, info *callContext) (*syscallmsg.ReturnValueRequest, lib.Id, string) {
	return sharedCallService(l.nameServer, key, info)
}
func (l *localSysCall) GetService(ctx context.Context, key dep.DepKey, pkgPath, service string) (lib.Id, lib.Id, string) {
	return sharedGetService(ctx, l.nameServer, key, pkgPath, service)
}
func (l *localSysCall) Require(ctx context.Context, key dep.DepKey, pkgPath, service string) lib.Id {
	return sharedRequire(ctx, l.nameServer, key, pkgPath, service)
}
func (l *localSysCall) RunBlock(ctx context.Context, key dep.DepKey) (bool, lib.Id) {
	return l.nameServer.RunBlock(ctx, key)
}

func (l *localSysCall) FindMethodByName(ctx context.Context, caller dep.DepKey, serviceId lib.Id, method string) (*callContext, lib.Id, string) {
	return l.nameServer.FindMethodByName(ctx, caller, serviceId, method)
}
func (l *localSysCall) GetInfoForCallId(ctx context.Context, cid lib.Id) *callContext {
	return l.nameServer.GetInfoForCallId(ctx, cid)
}
func (l *localSysCall) BlockUntilCall(ctx context.Context, key dep.DepKey, canTimeout bool) *callContext {
	info := l.nameServer.BlockUntilCall(ctx, key, canTimeout)
	// this loop is because we get the "error" case as a nil
	for info == nil {
		info = l.nameServer.BlockUntilCall(ctx, key, canTimeout)
	}
	return info
}
