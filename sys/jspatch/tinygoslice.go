//go:build tinygo
// +build tinygo

package jspatch

import (
	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/sys/jspatch/jsemul"
)

// convertSlice has to return a SliceHeader but tinygo (for no obvious reason) decided
// to change the definition to have uinptr values for Len and Cap.
func convertSlice(slice []byte, str jsemul.JsObject) *reflect.SliceHeader {
	ptr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	content := str.String()
	// XXX this has to be changed between tinygo and gc
	// XXX uintptr(len(content)) -> len(content)
	l := uintptr(len(content))
	ptr.Len = l
	ptr.Cap = l
	ptr.Data = uintptr(unsafe.Pointer((&content)))
	return ptr
}
