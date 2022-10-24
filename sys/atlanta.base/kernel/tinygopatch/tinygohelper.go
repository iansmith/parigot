package tinygopatch

import (
	"encoding/binary"
	"github.com/iansmith/parigot/sys/kernel/jspatch"
	"log"
	"math"
	"time"
	"unsafe"
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

func WasiWriteFd(fd int32, iovec int32, len int32, written int32) int32 {
	//log.Printf("fd_write called:  %d, %x, %d %x", fd, iovec, len, written)
	memPtr := jspatch.MemPtr()
	buf := make([]byte, len*8)
	for i := int32(0); i < len*8; i++ {
		b := (*byte)(unsafe.Pointer(memPtr + uintptr(iovec+i)))
		buf[i] = *b
	}
	u := binary.LittleEndian.Uint32(buf[0:4])
	v := binary.LittleEndian.Uint32(buf[4:8])
	s := jspatch.StrConvert(memPtr, int32(u), int32(v))
	log.Printf("iovec output: %s", s)
	b := (*byte)(unsafe.Pointer(memPtr + uintptr(written)))
	slice := unsafe.Slice(b, 4)
	binary.LittleEndian.PutUint32(slice, v)
	log.Printf("wrote result")
	return 0
}

func WasiProcExit(i0 int32) {
	log.Printf("fd_proc_exit called:  %d", i0)
}
