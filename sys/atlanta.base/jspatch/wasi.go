package jspatch

import (
	"fmt"
	"log"
)

type WasiPatch struct {
	mem *WasmMem
}

func NewWasiPatchWithMemPtr(memptr uintptr) *WasiPatch {
	return &WasiPatch{mem: NewWasmMem(memptr)}
}
func NewWasiPatch() *WasiPatch {
	return &WasiPatch{}
}
func (w *WasiPatch) SetMemPtr(m uintptr) {
	w.mem = NewWasmMem(m)
}

func (w *WasiPatch) WasiWrite(sp int32) {
	_ = w.mem.GetInt64(sp + 8)
	content := w.mem.LoadString(sp + 16)
	fmt.Printf("%s", string(content))
}
func (w *WasiPatch) WasiExit(sp int32) {
	log.Printf("ignoring! wasmExit: %d", w.mem.GetInt32(sp+8))
}
