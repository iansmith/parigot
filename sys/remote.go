package sys

import (
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/sys/dep"
)

// remoteSyscall is wrapped around a netNameServer which contains a local nameserver inside it.
type remoteSyscall struct {
	nameServer *NSProxy
}

func newRemoteSysCall(ns *NSProxy) *remoteSyscall {
	return &remoteSyscall{
		nameServer: ns,
	}
}

func (r *remoteSyscall) Bind(p *Process, packagePath, service, method string) (lib.Id, lib.Id) {
	sysPrint(logmsg.LogLevel_LOG_LEVEL_INFO, "BIND", "[remote] bind for method %s on (%s.%s)", method, packagePath, service)
	return sharedBind(r.nameServer.NSCore, p, packagePath, service, method)
}

func (r *remoteSyscall) Export(key dep.DepKey, pkg, service string) lib.Id {
	return sharedExport(r.nameServer, key, pkg, service)
}

func (r *remoteSyscall) CallService(key dep.DepKey, info *callContext) *syscallmsg.ReturnValueRequest {
	return sharedCallService(r.nameServer, key, info)
}

func (r *remoteSyscall) RunNotify(key dep.DepKey) {
	//nothing to do we don't use the runreader
}

// func (r *remoteSyscall) HandleMethod(key dep.DepKey, pkgPath, service, method string) (lib.Id, lib.Id) {
// 	return sharedHandleMethod(r.nameServer, key.(*DepKeyImpl).proc, pkgPath, service, method)
// }

func (r *remoteSyscall) FindMethodByName(caller dep.DepKey, sid lib.Id, method string) *callContext {
	return r.nameServer.FindMethodByName(caller, sid, method)
}

func (r *remoteSyscall) GetInfoForCallId(cid lib.Id) *callContext {
	return r.nameServer.GetInfoForCallId(cid)
}
func (r *remoteSyscall) GetService(key dep.DepKey, pkgPath, service string) (lib.Id, lib.KernelErrorCode) {
	return sharedGetService(r.nameServer, key, pkgPath, service)
}
func (r *remoteSyscall) Require(key dep.DepKey, pkgPath, service string) lib.Id {
	return sharedRequire(r.nameServer, key, pkgPath, service)
}
func (r *remoteSyscall) RunBlock(key dep.DepKey) (bool, lib.Id) {
	return r.nameServer.RunBlock(key)
}
func (l *remoteSyscall) BlockUntilCall(key dep.DepKey) *callContext {
	info := l.nameServer.BlockUntilCall(key)
	// this loop is because we get the "error" case as a nil
	for info == nil {
		info = l.nameServer.BlockUntilCall(key)
	}
	return info
}
