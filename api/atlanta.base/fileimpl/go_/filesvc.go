package go_

import (
	"errors"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"unicode"
	"unsafe"

	pb "github.com/iansmith/parigot/api/proto/g/pb/file"
	"github.com/iansmith/parigot/api/splitutil"
	"github.com/iansmith/parigot/lib"
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

	buffer := splitutil.ReadSlice(l.mem, wasmPtr,
		unsafe.Offsetof(splitutil.SinglePayload{}.InPtr),
		unsafe.Offsetof(splitutil.SinglePayload{}.InLen))
	req := pb.OpenRequest{}
	err := splitutil.ExtractProtoFromBytes(buffer, &req)
	if err != nil {
		kerr := lib.NewKernelError(lib.KernelUnmarshalFailed)
		// low goes first
		l.mem.SetInt64(int32(wasmPtr)+int32(unsafe.Offsetof(splitutil.SinglePayload{}.ErrPtr)),
			int64(kerr.Low()))
		// high is 8 bytes higher
		l.mem.SetInt64(int32(wasmPtr)+int32(unsafe.Offsetof(splitutil.SinglePayload{}.ErrPtr)+8),
			int64(kerr.High()))
	}
	log.Printf("xxx FileSvcOpen path to file %s", req.GetPath())
}

func (l *FileSvcImpl) SetWasmMem(ptr uintptr) {
	l.mem = jspatch.NewWasmMem(ptr)
}

// ValidatePathForParigot checks for and avoids many of the common pitfalls in pathnames but is not
// a substitute for chroot() or similar types jailing of user code implemented by the underlying operating
// systems.
func ValidatePathForParigot(path string, op string) (string, error) {

	cleaned := filepath.Clean(path)
	noFrontSlash := cleaned
	if strings.HasPrefix(cleaned, "/") {
		noFrontSlash = noFrontSlash[1:]
	} else {
		return "", &fs.PathError{
			Op:   op,
			Path: path,
			Err:  errors.New("parigot requires fully qualified path names"),
		}
	}
	if !fs.ValidPath(noFrontSlash) {
		return "", &fs.PathError{
			Op:   op,
			Path: path,
			Err:  errors.New("failed ValidPath() test"),
		}
	}
	part := filepath.SplitList(cleaned)
	if len(part) == 0 {
		return "", &fs.PathError{
			Op:   op,
			Path: path,
			Err:  errors.New("empty path"), //can this actually happen?
		}
	}
	for _, element := range part {
		for _, r := range element {
			if unicode.IsControl(r) {
				return "", &fs.PathError{
					Op:   op,
					Path: path,
					Err:  errors.New("control characters not allowed"),
				}
			}
		}
	}
	return cleaned, nil
}
