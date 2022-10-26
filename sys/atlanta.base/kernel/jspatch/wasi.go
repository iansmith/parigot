package jspatch

import (
	"fmt"
	"log"
)

type WasiPatch struct {
	mem *wasmMem
}

func NewWasiPatch(memptr uintptr) *WasiPatch {
	return &WasiPatch{mem: newWasmMem(memptr)}
}

func (w *WasiPatch) WasiWrite(sp int32) {
	_ = w.mem.getInt64(sp + 8)
	content := w.mem.loadString(sp + 16)
	fmt.Printf("%s", string(content))
}
func (w *WasiPatch) WasiExit(sp int32) {
	log.Printf("ignoring! wasmExit: %d", w.mem.getInt32(sp+8))
}
