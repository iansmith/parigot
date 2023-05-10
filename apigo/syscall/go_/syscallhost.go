package syscall

import (
	"log"

	"github.com/iansmith/parigot/eng"
)

func ParigotInit(e eng.Engine, inst eng.Instance) {
	e.AddSupportedFunc("parigot", "locate", wrapFunc(inst, locate))
	e.AddSupportedFunc("parigot", "dispatch", wrapFunc(inst, dispatch))
	e.AddSupportedFunc("parigot", "blockUntilCall", wrapFunc(inst, blockUntilCall))
	e.AddSupportedFunc("parigot", "bindMethod", wrapFunc(inst, bindMethod))
	e.AddSupportedFunc("parigot", "run", wrapFunc(inst, run))
	e.AddSupportedFunc("parigot", "export", wrapFunc(inst, export))
	e.AddSupportedFunc("parigot", "return_value", wrapFunc(inst, returnValue))
	e.AddSupportedFunc("parigot", "require", wrapFunc(inst, require))
}

func wrapFunc(i eng.Instance, fn func(eng.Instance, int32) int32) func(int32) int32 {
	return func(x int32) int32 {
		return fn(i, x)
	}
}

func locate(inst eng.Instance, ptr int32) int32 {
	log.Printf("locate 0x%x", ptr)
	return 0
}

func dispatch(inst eng.Instance, ptr int32) int32 {
	log.Printf("dispatch 0x%x", ptr)
	return 0
}
func blockUntilCall(inst eng.Instance, ptr int32) int32 {
	log.Printf("blockUntilCall 0x%x", ptr)
	return 0
}
func bindMethod(inst eng.Instance, ptr int32) int32 {
	log.Printf("bindMethod 0x%x", ptr)
	return 0
}
func run(inst eng.Instance, ptr int32) int32 {
	log.Printf("run 0x%x", ptr)
	return 0
}
func export(inst eng.Instance, ptr int32) int32 {
	log.Printf("export 0x%x", ptr)
	return 0
}
func returnValue(inst eng.Instance, ptr int32) int32 {
	log.Printf("returnValue 0x%x", ptr)
	return 0
}
func require(inst eng.Instance, ptr int32) int32 {
	log.Printf("require 0x%x", ptr)
	return 0
}
func exit(inst eng.Instance, ptr int32) int32 {
	log.Printf("exit 0x%x", ptr)
	panic("exit called ")
	return 0
}
