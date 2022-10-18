package tinygopatch

import "os"

//go:wasm-module env
//export runtime.ticks
func Ticks() float64 {
	print("runtime.ticks()")
	//os.Exit(1)
	return 0.0
}

//go:wasm-module env
//export runtime.sleepTicks
func SleepTicks(float64) {
	print("Tinygo not implemented: runtime.sleepTicks")
	os.Exit(1)
}
