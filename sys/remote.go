package sys

import (
	"context"

	pcontext "github.com/iansmith/parigot/context"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/sys/dep"
)

const sysVerbose = true

// remoteSyscall is wrapped around a netNameServer which contains a local nameserver inside it.
type remoteSyscall struct {
	nameServer *NSProxy
}

func newRemoteSysCall(ns *NSProxy) *remoteSyscall {
	return &remoteSyscall{
		nameServer: ns,
	}
}

func (r *remoteSyscall) Bind(ctx context.Context, p *Process, packagePath, service, method string) (lib.Id, lib.Id) {
	sysPrint(ctx, "BIND", "[remote] bind for method %s on (%s.%s)", method, packagePath, service)
	return sharedBind(r.nameServer.NSCore, p, packagePath, service, method)
}

func (r *remoteSyscall) Export(ctx context.Context, key dep.DepKey, pkg, service string) lib.Id {
	return sharedExport(ctx, r.nameServer, key, pkg, service)
}

func (r *remoteSyscall) CallService(ctx context.Context, key dep.DepKey, info *callContext) (*syscallmsg.ReturnValueRequest, lib.Id, string) {
	return sharedCallService(ctx, r.nameServer, key, info)
}

// func (r *remoteSyscall) HandleMethod(ctx context.Context,key dep.DepKey, pkgPath, service, method string) (lib.Id, lib.Id) {
// 	return sharedHandleMethod(r.nameServer, key.(*DepKeyImpl).proc, pkgPath, service, method)
// }

func (r *remoteSyscall) FindMethodByName(ctx context.Context, caller dep.DepKey, sid lib.Id, method string) (*callContext, lib.Id, string) {
	return r.nameServer.FindMethodByName(ctx, caller, sid, method)
}

func (r *remoteSyscall) GetInfoForCallId(ctx context.Context, cid lib.Id) *callContext {
	return r.nameServer.GetInfoForCallId(ctx, cid)
}
func (r *remoteSyscall) GetService(ctx context.Context, key dep.DepKey, pkgPath, service string) (lib.Id, lib.Id, string) {
	return sharedGetService(ctx, r.nameServer, key, pkgPath, service)
}
func (r *remoteSyscall) Require(ctx context.Context, key dep.DepKey, pkgPath, service string) lib.Id {
	return sharedRequire(ctx, r.nameServer, key, pkgPath, service)
}
func (r *remoteSyscall) RunBlock(ctx context.Context, key dep.DepKey) (bool, lib.Id) {
	return r.nameServer.RunBlock(ctx, key)
}
func (l *remoteSyscall) BlockUntilCall(ctx context.Context, key dep.DepKey, canTimeout bool) *callContext {
	info := l.nameServer.BlockUntilCall(ctx, key, canTimeout)
	// this loop is because we get the "error" case as a nil
	for info == nil {
		info = l.nameServer.BlockUntilCall(ctx, key, canTimeout)
	}
	return info
}
func sysPrint(ctx context.Context, funcName, spec string, rest ...interface{}) {
	if sysVerbose {
		pcontext.LogFullf(ctx, pcontext.Debug, pcontext.Parigot, funcName, spec, rest...)
	}
}
