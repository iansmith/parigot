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
	sysPrint("BIND", "about to jump to sharedBind")
	return sharedBind(r.nameServer.local, p, packagePath, service, method)
}

func (r *remoteSyscall) Export(key dep.DepKey, pkg, service string) lib.Id {
	return sharedExport(r.nameServer, key, pkg, service)
}

func (r *remoteSyscall) CallService(key dep.DepKey, info *callInfo) (*resultInfo, lib.Id) {
	return sharedCallService(r.nameServer, key, info)
}

func (r *remoteSyscall) RunNotify(key dep.DepKey) {
	//nothing to do we don't use the runreader
}

func (r *remoteSyscall) HandleMethod(key dep.DepKey, pkgPath, service, method string) (lib.Id, lib.Id) {
	return sharedHandleMethod(r.nameServer, key.(*DepKeyImpl).proc, pkgPath, service, method)
}

func (r *remoteSyscall) FindMethodByName(caller dep.DepKey, sid lib.Id, method string) *callContext {
	return r.nameServer.FindMethodByName(caller, sid, method)
}

func (r *remoteSyscall) GetProcessForCallId(cid lib.Id) dep.DepKey {
	panic("GetProcessForCallId on remote syscall")
}
func (r *remoteSyscall) GetService(key dep.DepKey, pkgPath, service string) (lib.Id, lib.Id) {
	return sharedGetService(r.nameServer, key, pkgPath, service)
}
func (r *remoteSyscall) Require(key dep.DepKey, pkgPath, service string) lib.Id {
	return sharedRequire(r.nameServer, key, pkgPath, service)
}
func (r *remoteSyscall) RunBlock(key dep.DepKey) (bool, lib.Id) {
	return r.nameServer.RunBlock(key)
}
func (l *remoteSyscall) BlockUntilCall(key dep.DepKey) *callInfo {
	print("xxx block until call hit on remote ", key.String(), "\n")
	return l.nameServer.BlockUntilCall(key)
}
