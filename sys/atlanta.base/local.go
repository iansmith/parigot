package sys

import (
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/dep"
)

type localSysCall struct {
	nameServer *LocalNameServer
}

// sharedBind uses localNameServer so the remote syscall cannot (in error) pass its netNameServer, it has to
// pass the nested localNameServer inside the netNameServer
func sharedBind(ns *LocalNameServer, p *Process, packagePath, service, method string) (lib.Id, lib.Id) {
	sysPrint("SHAREDBIND", "about to hit handle method")
	return ns.HandleMethod(p, packagePath, service, method)
}

// sharedFindMethodByName uses localNameServer so the remote syscall cannot (in error) pass its netNameServer, it has to
// pass the nested localNameServer inside the netNameServer
func sharedFindMethodByName(ns *LocalNameServer, caller *Process, sid lib.Id, method string) *callContext {
	return ns.FindMethodByName(caller, sid, method)
}

func sharedExport(ns NameServer, key dep.DepKey, pkg, service string) lib.Id {
	if kerr := ns.CloseService(pkg, service); kerr != nil {
		return kerr
	}
	if kerr := ns.Export(key, pkg, service); kerr != nil {
		return kerr
	}
	return nil

}

func sharedHandleMethod(ns NameServer, proc *Process, pkg, service, method string) (lib.Id, lib.Id) {
	return ns.HandleMethod(proc, pkg, service, method)
}

func newLocalSysCall(ns *LocalNameServer) *localSysCall {
	return &localSysCall{nameServer: ns}
}

func (l *localSysCall) Bind(proc *Process, packagePath, service, method string) (lib.Id, lib.Id) {
	return sharedBind(l.nameServer, proc, packagePath, service, method)
}

func (l *localSysCall) Export(key dep.DepKey, pkg, service string) lib.Id {
	return sharedExport(l.nameServer, key, pkg, service)
}

func (l *localSysCall) RunNotify(key dep.DepKey) {
	l.nameServer.RunNotify(key)
}
func (l *localSysCall) RunBlock(key dep.DepKey) bool {
	return l.nameServer.RunBlock(key)
}

func (l *localSysCall) FindMethodByName(caller *Process, serviceId lib.Id, name string) *callContext {
	return sharedFindMethodByName(l.nameServer, caller, serviceId, name)
}
func (l *localSysCall) GetProcessForCallId(p *Process, cid lib.Id) *Process {
	return l.nameServer.GetProcessForCallId(cid)
}
func (l *localSysCall) GetService(p *Process, pkgPath, service string) (lib.Id, lib.Id) {
	return l.nameServer.GetService(pkgPath, service)
}
func (l *localSysCall) Require(key dep.DepKey, pkgPath, service string) lib.Id {
	return l.nameServer.Require(key, pkgPath, service)
}
