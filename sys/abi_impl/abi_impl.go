package abi_impl

import (
	"os"
	"time"
)

func OutputString(s string) {
	print(s)
}

func JSNotImplemented() {
	Exit(1)
}

func JSHandleEvent(e int32) {
	OutputString("JSHandleEvent\n")
	JSNotImplemented()
}

func TinyGoNotImplemented() {
	Exit(1)
}

func Exit(i int32) {
	os.Exit(int(i))
}

func Now() int64 {
	return time.Now().UnixNano()
}

func SetNow(_ int64) {
	print("SetNow\n")
}
