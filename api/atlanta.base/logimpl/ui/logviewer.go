//go:build !js
// +build !js

package ui

import (
	"log"

	"github.com/iansmith/parigot/sys/jspatch"
)

type LogViewerImpl struct {
	mem *jspatch.WasmMem
}

//go:noinline
func (l *LogViewerImpl) LogRequestViaSocket(sp int32) {
	wasmPtr := l.mem.GetInt64(sp + 8)
	log.Printf("LogRequestViaSocket  wasmptr %x,true=%x", wasmPtr, l.mem.TrueAddr(int32(wasmPtr)))

	buffer := l.ReadSlice(wasmPtr, 0, 0)

	log.Printf("LogRequestViaSocket, size of buffer: %d", len(buffer))

}

func (s *LogViewerImpl) ReadSlice(structPtr int64, dataOffset uintptr, lenOffset uintptr) []byte {
	return s.mem.LoadSliceWithLenAddr(int32(structPtr)+int32(dataOffset),
		int32(structPtr)+int32(lenOffset))
}

func (l *LogViewerImpl) SetWasmMem(m *jspatch.WasmMem) {
	l.mem = m
}
