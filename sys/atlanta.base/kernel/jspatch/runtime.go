package jspatch

import (
	"log"
	"math/rand"
	"time"
)

type RuntimePatch struct {
	mem *wasmMem
}

func NewRuntimePatchWithMemPtr(memptr uintptr) *RuntimePatch {
	return &RuntimePatch{mem: newWasmMem(memptr)}
}
func NewRuntimePatch() *RuntimePatch {
	return &RuntimePatch{}
}

func (r *RuntimePatch) Nanotime1(sp int32) {
	r.mem.setInt64(sp+8, time.Now().UnixNano())
}

func (j *RuntimePatch) SetMemPtr(m uintptr) {
	j.mem = newWasmMem(m)
}

func (r *RuntimePatch) GetRandomData(sp int32) {
	b := r.mem.loadSlice(sp + 8)
	_, _ = rand.Read(b) //docs say no returned error
}
func (r *RuntimePatch) WallTime(sp int32) {
	secs := time.Now().Unix()
	r.mem.setInt64(sp+8, secs)
	nanos := time.Now().UnixNano()
	nanos = nanos % 1000000000
	r.mem.setInt32(sp+16, int32(nanos)) // this is big enough, max size is 1B and 31 bits is 2B
	log.Printf("walltime executed %d,%d", secs, nanos)
}
func (r *RuntimePatch) ScheduleTimeoutEvent(sp int32) {
	panic("got a scheduleTimeoutEvent call, aborting")
}
func (r *RuntimePatch) ClearTimeoutEvent(sp int32) {
	panic("got a clearTimeoutEvent call, aborting")
}

func (r *RuntimePatch) GoDebug(sp int32) {
	panic("got a go.debug call, aborting")
}

func (r *RuntimePatch) ResetMemoryDataView(sp int32) {
	panic("got a resetMemoryDataView call, aborting")
}
