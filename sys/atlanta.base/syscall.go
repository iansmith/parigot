package sys

import (
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/dep"
)

type SysCall interface {
	Bind(p *Process, packagePath, service, method string) (lib.Id, lib.Id)
	Export(key dep.DepKey, packagePath, service string) lib.Id
	Require(key dep.DepKey, packagePath, service string) lib.Id
	RunBlock(key dep.DepKey) (bool, lib.Id)
	RunNotify(key dep.DepKey)
	GetInfoForCallId(cid lib.Id) *callContext
	FindMethodByName(key dep.DepKey, sid lib.Id, method string) *callContext
	GetService(key dep.DepKey, packagePath, service string) (lib.Id, lib.Id)
	CallService(key dep.DepKey, info *callContext) (*resultInfo, lib.Id)
	BlockUntilCall(key dep.DepKey) *callContext
}
