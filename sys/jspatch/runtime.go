package jspatch

import (
	"fmt"
	"time"
)

type RuntimePatch struct {
	mem *WasmMem
}

func NewRuntimePatchWithMemPtr(memptr uintptr) *RuntimePatch {
	print("new runtime patch", memptr, "\n")
	return &RuntimePatch{mem: NewWasmMem(memptr)}
}
func NewRuntimePatch() *RuntimePatch {
	return &RuntimePatch{}
}

func (r *RuntimePatch) Nanotime1(sp int32) {
	r.mem.SetInt64(sp+8, time.Now().UnixNano())
}

func (j *RuntimePatch) SetMemPtr(m uintptr) {
	//print(fmt.Sprintf("set mem ptr: %x\n", m))
	j.mem = NewWasmMem(m)
}

func (r *RuntimePatch) GetRandomData(sp int32) {
	print("xxx ignoring call to GetRandomData", "\n")
	//return
	// b := r.mem.LoadSlice(sp + 8)
	// print("xxxx GetRandomData", uint32(len(b)), "\n")
	// _, _ = rand.Read(b) //docs say no returned error
}
func (r *RuntimePatch) WallTime(sp int32) {
	secs := time.Now().Unix()
	r.mem.SetInt64(sp+8, secs)
	nanos := time.Now().UnixNano()
	nanos = nanos % 1000000000
	r.mem.SetInt32(sp+16, int32(nanos)) // this is big enough, max size is 1B and 31 bits is 2B
	//log.Printf("walltime executed %d,%d", secs, nanos)
}
func (r *RuntimePatch) ScheduleTimeoutEvent(sp int32) {
	t := r.mem.GetInt64(sp + 8)
	print(fmt.Sprintf("got a call to ScheduleTimeoutEvent but time is %x, %d\n", t, t/1000))
	panic("got a scheduleTimeoutEvent call, aborting")
}
func (r *RuntimePatch) ClearTimeoutEvent(sp int32) {
	panic("got a clearTimeoutEvent call, aborting")
}

func (r *RuntimePatch) GoDebug(sp int32) {
	panic("got a go.debug call, aborting")
}

func (r *RuntimePatch) ResetMemoryDataView(sp int32) {
	i0 := r.mem.GetInt64(sp + 8)
	i1 := r.mem.GetInt64(sp + 16)
	i2 := r.mem.GetInt64(sp + 24)
	print(fmt.Sprintf("reset memory data view: %x, %x, %x\n", i0, i1, i2))
}
