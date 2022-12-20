package go_

import (
	"log"
	"unsafe"

	"github.com/iansmith/parigot/api/splitutil"
	"github.com/iansmith/parigot/sys/jspatch"
)

type FileSvcImpl struct {
	mem *jspatch.WasmMem
}

// This is the native code side of the logviewer.  It reads the payload from the WASM world and either
// dumps it to the terminal or sends it through the UD socket to the GUI.

//go:noinline
func (l *FileSvcImpl) FileSvcOpen(sp int32) {
	wasmPtr := l.mem.GetInt64(sp + 8)

	buffer := splitutil.ReadSlice(l.mem, wasmPtr, unsafe.Offsetof(splitutil.SplitUtilSinglePayload{}.Ptr),
		unsafe.Offsetof(splitutil.SplitUtilSinglePayload{}.Len))
	req, err := splitutil.DecodeProto[FileOpenRequest](buffer)
	log.Printf("xxx FileSvcOpen path to file %s", req.GetPath())
}
