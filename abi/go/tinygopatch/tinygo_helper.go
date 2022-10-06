package tinygopatch

import "github.com/iansmith/parigot/abi/go/abi"

//go:wasm-module env
//export runtime.ticks
func Ticks() float64 {
	abi.TinyGoNotImplemented()
	return 0
}

//go:wasm-module env
//export runtime.sleepTicks
func SleepTicks(float64) {
	abi.TinyGoNotImplemented()
}
