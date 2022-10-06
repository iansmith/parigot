package abi_impl

import (
	"github.com/iansmith/parigot/abi/go/abi"
	"os"
	"time"
)

func OutputString(s string) {
	print(s)
}

func JSNotImplemented() {
	abi.Exit(1)
}

func JSHandleEvent(e int32) {
	abi.OutputString("JSHandleEvent\n")
	JSNotImplemented()
}

func TinyGoNotImplemented() {
	abi.Exit(1)
}

func Exit(i int32) {
	os.Exit(int(i))
}

func Exit_(mem uintptr, i int32) {
	Exit(i)
}

func Now() int64 {
	return time.Now().UnixNano()
}

func SetNow(_ int64) {
	print("SetNow\n")
}
