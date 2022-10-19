package tinygopatch

import (
	"log"
	"math"
	"time"
)

//go:wasm-module env
//export runtime.ticks
func Ticks() float64 {
	log.Printf("runtime.ticks(), just giving them 0")
	//os.Exit(1)
	return 0.0
}

//go:wasm-module env
//export runtime.sleepTicks
func SleepTicks(f float64) {
	log.Printf("Tinygo impl: runtime.sleepTicks %f (%f)", f, math.Trunc(f))
	time.Sleep(time.Duration(int(math.Trunc(f))) * time.Millisecond)
}

func WasiWriteFd(i0 int32, i1 int32, i2 int32, i3 int32) int32 {
	log.Printf("fd_write called:  %x, %d, %x %d", i0, i1, i2, i3)
	log.Printf("giving them 20")
	return 20
}
func WasiProcExit(i0 int32) {
	log.Printf("fd_proc_exit called:  %d", i0)
}
