package tinygopatch

//go:wasm-module env
//export runtime.ticks
func Ticks() float64 {
	abi.TinygoNotImplemented("Tinygo not implement: runtime.ticks()")
	return 0
}

//go:wasm-module env
//export runtime.sleepTicks
func SleepTicks(float64) {
	abi.TinygoNotImplemented("Tinygo not implemented: runtime.sleepTicks")
}
