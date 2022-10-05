package abi

//go:wasm-module env
//export runtime.ticks
func Ticks() float64 {
	TinyGoNotImplemented()
	return 0
}

//go:wasm-module env
//export runtime.sleepTicks
func SleepTicks(float64) {
	TinyGoNotImplemented()
}
