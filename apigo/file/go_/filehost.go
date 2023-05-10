package go_

import (
	"log"

	"github.com/iansmith/parigot/eng"
)

func ParigotInit(e eng.Engine, inst eng.Instance) {
	e.AddSupportedFunc("file", "open", wrapFunc(inst, open))
	e.AddSupportedFunc("file", "load_testt_data", wrapFunc(inst, loadTestData))
}

func wrapFunc(i eng.Instance, fn func(eng.Instance, int32) int32) func(int32) int32 {
	return func(x int32) int32 {
		return fn(i, x)
	}
}

func open(inst eng.Instance, ptr int32) int32 {
	log.Printf("file.open 0x%x", ptr)
	return 0
}

func loadTestData(inst eng.Instance, ptr int32) int32 {
	log.Printf("file.dispatch 0x%x", ptr)
	return 0
}
