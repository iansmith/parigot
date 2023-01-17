package sys

import (
	"github.com/iansmith/parigot/sys/jspatch"
)

type Runtime struct {
	// setup the objects that the kernel needs to do things like handle the golang
	// startup code that expects js and wasi
	jsEnv      *jspatch.JSPatch
	wasiEnv    *jspatch.WasiPatch
	runtimeEnv *jspatch.RuntimePatch
	syscall    *syscallReadWrite
	//spec       *RemoteSpec
}

// type RemoteSpec struct {
// 	remote []string
// 	local  []string
// }

// func (r *RemoteSpec) IsRemote(proc *Process) bool {
// 	for _, path := range r.remote {
// 		if path == proc.path {
// 			return true
// 		}
// 	}
// 	return false
// }
// func (r *RemoteSpec) IsLocal(proc *Process) bool {
// 	for _, path := range r.local {
// 		if path == proc.path {
// 			return true
// 		}
// 	}
// 	return false
// }

// func NewRemoteSpec(remote []string, local []string) *RemoteSpec {
// 	if len(remote) > 0 && len(local) > 0 {
// 		panic("mixed mode of runtime/syscall is not yet supported")
// 	}
// 	return &RemoteSpec{local: local, remote: remote}
// }

func newRuntime(ctx *DeployContext) *Runtime {
	return &Runtime{
		jsEnv:      jspatch.NewJSPatch(),
		wasiEnv:    jspatch.NewWasiPatch(),
		runtimeEnv: jspatch.NewRuntimePatch(),
		syscall:    NewSysCallRW(ctx.nameserver),
		//spec:       spec,
	}
}

func (r *Runtime) SetProcess(p *Process) {
	r.syscall.SetProcess(p)
}

func (r *Runtime) SetMemPtr(memPtr uintptr) {
	r.jsEnv.SetMemPtr(memPtr)
	r.wasiEnv.SetMemPtr(memPtr)
	r.runtimeEnv.SetMemPtr(memPtr)
	r.syscall.SetMemPtr(memPtr)
}
