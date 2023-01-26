package sys

import (
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
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
	GetInfoForCallId(cid lib.Id) *callContext
	// FindMethodByName searches for a method name with a known service.  If the method
	// name is found, a call context for that method will be returned plus a nil, "" pair
	// indicating no error.  If there was an error, the first return value will be
	// nil and the remaining two will give the error info.
	FindMethodByName(key dep.DepKey, sid lib.Id, method string) (*callContext, lib.Id, string)
	// GetService looks up the given service (packagePath.service) and returns either
	// the service id in the first parameter or an error pair of an error id and an
	// error detail.
	GetService(key dep.DepKey, packagePath, service string) (lib.Id, lib.Id, string)
	// CallService, in the success case, returns a return value request and nil, "" for the
	// error values.  The ReturnValueRequest is ready to be use in a call
	// to ReturnValue().  If there was an error, the return value is nil and the
	// last two error values will be filled in.
	CallService(key dep.DepKey, info *callContext) (*syscallmsg.ReturnValueRequest, lib.Id, string)
	BlockUntilCall(key dep.DepKey, canTimeout bool) *callContext
}
