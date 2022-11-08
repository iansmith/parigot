package sys

import (
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/dep"
)

// remoteSyscall is wrapped around a netNameServer which contains a local nameserver inside it.
type remoteSyscall struct {
	nameServer *NetNameServer
}

func newRemoteSysCall(ns *NetNameServer) *remoteSyscall {
	return &remoteSyscall{
		nameServer: ns,
	}
}

func (r *remoteSyscall) Bind(p *Process, packagePath, service, method string) (lib.Id, lib.Id) {
	return sharedBind(r.nameServer.local, p, packagePath, service, method)
}

func (r *remoteSyscall) Export(key dep.DepKey, pkg, service string) lib.Id {
	return sharedExport(r.nameServer, key, pkg, service)
}

func (r *remoteSyscall) RunNotify(key dep.DepKey) {
	//nothing to do we don't use the runreader
}

func (r *remoteSyscall) HandleMethod(key dep.DepKey, pkgPath, service, method string) (lib.Id, lib.Id) {
	return sharedHandleMethod(r.nameServer, key.(*depKeyImpl).proc, pkgPath, service, method)
}

func (r *remoteSyscall) FindMethodByName(caller *Process, sid lib.Id, method string) *callContext {
	return sharedFindMethodByName(r.nameServer.local, caller, sid, method)
}

func (r *remoteSyscall) GetProcessForCallId(p *Process, cid lib.Id) *Process {
	panic("GetProcessForCallId on remote syscall")
}
func (r *remoteSyscall) GetService(p *Process, pkgPath, service string) (lib.Id, lib.Id) {
	panic("GetService on remote syscall")
}
func (r *remoteSyscall) Require(key dep.DepKey, pkgPath, service string) lib.Id {
	panic("Require on remote syscall")
}
func (r *remoteSyscall) RunBlock(_ dep.DepKey) bool {
	panic("run block on remote should be calling the network?")
}
