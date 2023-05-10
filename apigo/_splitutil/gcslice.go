//go:build !tinygo
// +build !tinygo

package splitutil

import (
	"reflect"
	"unsafe"
)

// convertSlice has to return a SliceHeader and this is for the "normal" case
// of compiling host go code.

func changeUnderlyingSlice(slice []byte, l int, data uintptr) {
	// slice header is different in gc and tinygo
	ptr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	ptr.Len = l
	ptr.Cap = l
	ptr.Data = data
}
