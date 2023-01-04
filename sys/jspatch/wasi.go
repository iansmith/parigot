package jspatch

import (
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
	content := w.mem.LoadString(sp + 16)
	print(string(content))
}
func (w *WasiPatch) WasiExit(sp int32) {
	// xxxxfix me, this should be implemented such that our WASM-level process returns
	// to end it's goroutine
	log.Printf("ignoring! wasmExit: %d", w.mem.GetInt32(sp+8))
}
