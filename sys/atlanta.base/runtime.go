package sys

import (
	"log"

	"github.com/iansmith/parigot/sys/jspatch"
)

type Runtime struct {
	// setup the objects that the kernel needs to do things like handle the golang
	// startup code that expects js and wasi
	jsEnv      *jspatch.JSPatch
	wasiEnv    *jspatch.WasiPatch
	runtimeEnv *jspatch.RuntimePatch
	syscall    *SysCall
}

func newRuntime(nameServer *NameServer) *Runtime {
	return &Runtime{
		jsEnv:      jspatch.NewJSPatch(),
		wasiEnv:    jspatch.NewWasiPatch(),
		runtimeEnv: jspatch.NewRuntimePatch(),
		syscall:    NewSysCall(nameServer),
	}
}

func (r *Runtime) SetProcess(p *Process) {
	log.Printf("RT set process %x", p)
	r.syscall.SetProcess(p)
}

func (r *Runtime) SetMemPtr(memPtr uintptr) {
	log.Printf("RT set memPtr %x", memPtr)
	r.jsEnv.SetMemPtr(memPtr)
	r.wasiEnv.SetMemPtr(memPtr)
	r.runtimeEnv.SetMemPtr(memPtr)
	r.syscall.SetMemPtr(memPtr)
}
