//go:build !tinygo
// +build !tinygo

package jspatch

import (
	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/sys/jspatch/jsemul"
)

// convertSlice has to return a SliceHeader and this is for the "normal" case
// of compiling host go code.

func convertSlice(slice []byte, str jsemul.JsObject) *reflect.SliceHeader {
	ptr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	content := str.String()
	// XXX this has to be changed between tinygo and gc
	// XXX uintptr(len(content)) -> len(content)
	l := len(content)
	ptr.Len = l
	ptr.Cap = l
	ptr.Data = uintptr(unsafe.Pointer((&content)))
	return ptr
}
