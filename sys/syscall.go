package sys

import (
	pbsys "github.com/iansmith/parigot/api/proto/g/pb/syscall"
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/dep"
)

// Syscall is the set of calls that the kernel responds to. This is the kernel-side analogue
// of the interface Call that defines the kernel's interface on the "other side".  The Call
// interface is what WASM programs sees, the SysCall interface is what the go implementation
// of the kernel sees.  Although these two interfaces are closely related, they are not
// identical.  The go implementation of the "handler" for these calls is in `syscallrw.go`
// and the wasm implementation (also implemented in go) is in `callimpl.go`.
//
// There are two implementations of this interface, localSysCall and remoteSysCall that share
// some parts of their implementations.  The localSysCall handles the case of implementing these
// functions for services in the same process and remoteSysCall handles the case of implementing
// these for services across a network from each other.
type SysCall interface {
	Bind(p *Process, packagePath, service, method string) (lib.Id, lib.Id)
	Export(key dep.DepKey, packagePath, service string) lib.Id
	Require(key dep.DepKey, packagePath, service string) lib.Id
	RunBlock(key dep.DepKey) (bool, lib.Id)
	RunNotify(key dep.DepKey)
	GetInfoForCallId(cid lib.Id) *callContext
	FindMethodByName(key dep.DepKey, sid lib.Id, method string) *callContext
	GetService(key dep.DepKey, packagePath, service string) (lib.Id, lib.KernelErrorCode)
	CallService(key dep.DepKey, info *callContext) *pbsys.ReturnValueRequest
	BlockUntilCall(key dep.DepKey) *callContext
}
