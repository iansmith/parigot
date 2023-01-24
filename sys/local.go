package sys

import (
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

func sharedExport(ns NameServer, key dep.DepKey, pkg, service string) lib.Id {
	if kerr := ns.CloseService(key, pkg, service); kerr != nil && kerr.IsError() {
		return kerr
	}
	if kerr := ns.Export(key, pkg, service); kerr != nil && kerr.IsError() {
		return kerr
	}
	return nil
}

//	func sharedHandleMethod(ns NameServer, proc *Process, pkg, service, method string) (lib.Id, lib.Id) {
//		return ns.HandleMethod(proc, pkg, service, method)
//	}
func sharedRequire(ns NameServer, key dep.DepKey, pkg, service string) lib.Id {
	return ns.Require(key, pkg, service)
}
func sharedGetService(ns NameServer, key dep.DepKey, pkg, service string) (lib.Id, lib.KernelErrorCode) {
	return ns.GetService(key, pkg, service)
}
func sharedCallService(ns NameServer, key dep.DepKey, info *callContext) *syscallmsg.ReturnValueRequest {
	return ns.CallService(key, info)
}

func newLocalSysCall(ns *LocalNameServer) *localSysCall {
	return &localSysCall{nameServer: ns}
}

func (l *localSysCall) Bind(proc *Process, packagePath, service, method string) (lib.Id, lib.Id) {
	return sharedBind(l.nameServer.NSCore, proc, packagePath, service, method)
}

func (l *localSysCall) Export(key dep.DepKey, pkg, service string) lib.Id {
	return sharedExport(l.nameServer, key, pkg, service)
}

func (l *localSysCall) CallService(key dep.DepKey, info *callContext) *syscallmsg.ReturnValueRequest {
	return sharedCallService(l.nameServer, key, info)
}
func (l *localSysCall) GetService(key dep.DepKey, pkgPath, service string) (lib.Id, lib.KernelErrorCode) {
	return sharedGetService(l.nameServer, key, pkgPath, service)
}
func (l *localSysCall) Require(key dep.DepKey, pkgPath, service string) lib.Id {
	return sharedRequire(l.nameServer, key, pkgPath, service)
}
func (l *localSysCall) RunBlock(key dep.DepKey) (bool, lib.Id) {
	return l.nameServer.RunBlock(key)
}

func (l *localSysCall) FindMethodByName(caller dep.DepKey, serviceId lib.Id, method string) *callContext {
	return l.nameServer.FindMethodByName(caller, serviceId, method)
}
func (l *localSysCall) GetInfoForCallId(cid lib.Id) *callContext {
	return l.nameServer.GetInfoForCallId(cid)
}
func (l *localSysCall) BlockUntilCall(key dep.DepKey, canTimeout bool) *callContext {
	info := l.nameServer.BlockUntilCall(key, canTimeout)
	// this loop is because we get the "error" case as a nil
	for info == nil {
		info = l.nameServer.BlockUntilCall(key, canTimeout)
	}
	return info
}
