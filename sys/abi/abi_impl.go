package abi

import (
	"os"
	"time"
)

func OutputString(s string) {
	print(s)
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

func Exit(i int64) {
	print("Exit\n")
	os.Exit(int(i))
}

func Now() int64 {
	print("Now\n")
	return time.Now().Unix()
}

func SetNow(_ int64) {
	print("SetNow\n")
}
