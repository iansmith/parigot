package jspatch

import (
	"math/rand"
	"time"
)

type RuntimePatch struct {
	mem *wasmMem
}

func NewRuntimePatch(memptr uintptr) *RuntimePatch {
	return &RuntimePatch{mem: newWasmMem(memptr)}
}

func (r *RuntimePatch) Nanotime1(sp int32) {
	r.mem.setInt64(sp+8, time.Now().UnixNano())
}

func (r *RuntimePatch) GetRandomData(sp int32) {
	b := r.mem.loadSlice(sp + 8)
	_, _ = rand.Read(b) //docs say no returned error
}
func (r *RuntimePatch) WallTime(sp int32) {
	panic("got a walltime call, aborting")
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
