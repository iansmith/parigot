package sys

import (
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/dep"
)

type SysCall interface {
	Bind(p *Process, packagePath, service, method string) (lib.Id, lib.Id)
	Export(key dep.DepKey, packagePath, service string) lib.Id
	Require(key dep.DepKey, packagePath, service string) lib.Id
	RunBlock(key dep.DepKey) bool
	RunNotify(key dep.DepKey)
	GetProcessForCallId(p *Process, cid lib.Id) *Process
	FindMethodByName(p *Process, sid lib.Id, method string) *callContext
	GetService(p *Process, packagePath, service string) (lib.Id, lib.Id)
}
