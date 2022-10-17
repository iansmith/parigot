package main

import (
	"reflect"
	"unsafe"
)

//export checkAddr
func checkAddr(u unsafe.Pointer, l int32)

func main() {
	x := "foo"
	v := reflect.ValueOf(x)
	ptr := v.UnsafePointer()
	checkAddr(ptr, int32(len(x)))
}
