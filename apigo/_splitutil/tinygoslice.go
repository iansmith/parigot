//go:build tinygo
// +build tinygo

package splitutil

import (
	"reflect"
	"unsafe"
)

// convertSlice has to return a SliceHeader but tinygo (for no obvious reason) decided
// to change the definition to have uinptr values for Len and Cap.

func changeUnderlyingSlice(slice []byte, l int, data uintptr) {
	// slice header is different in gc and tinygo
	ptr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	ptr.Len = uintptr(l)
	ptr.Cap = uintptr(l)
	ptr.Data = data
}
