package abi_impl

import (
	"os"
	"reflect"
	"time"
	"unsafe"
)

func OutputString(s string) {
	print(s)
}

func strConvert(mem uintptr, ptr int32, length int32) string {
	addr := mem + uintptr(ptr)
	var data []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = addr
	sh.Len = int(length)
	sh.Cap = int(length)
	// this assumes there is no GC running!
	return string(data)
}

func OutputString_(mem uintptr, ptr int32, length int32) {
	OutputString(strConvert(mem, ptr, length))
}

func JSNotImplemented() {
	print("JSNotImplemented\n")
	panic("JSNotImplemented")
}

func JSHandleEvent(e int32) {
	print("JSHandleEvent\n")
	panic("JSHandleEvent")
}

func TinyGoNotImplemented() {
	print("TinyGoNotImplemented\n")
	panic("TinyGoNotImplemented")
}

func Exit(i int32) {
	os.Exit(int(i))
}

func Exit_(i int32) {
	Exit(i)
}

func Now() int64 {
	return time.Now().UnixNano()
}

func Now_() int64 {
	return Now()
}

func SetNow(_ int64) {
	print("SetNow\n")
}
