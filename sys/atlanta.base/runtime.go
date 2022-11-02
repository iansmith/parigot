package sys

import (
	"github.com/iansmith/parigot/sys/jspatch"
)

type runtime struct {
	// setup the objects that the kernel needs to do things like handle the golang
	// startup code that expects js and wasi
	jsEnv      *jspatch.JSPatch
	wasiEnv    *jspatch.WasiPatch
	runtimeEnv *jspatch.RuntimePatch
	syscall    *SysCall
}

func newRuntime() *runtime {
	ns := newNameServer()
	return &runtime{
		jsEnv:      jspatch.NewJSPatch(),
		wasiEnv:    jspatch.NewWasiPatch(),
		runtimeEnv: jspatch.NewRuntimePatch(),
		syscall:    NewSysCall(ns),
	}
}
