package abiimpl

import (
	"os"
	"time"
)

type AbiImpl struct {
}

func (a *AbiImpl) OutputString(s string) {
	print(s)
}

func (a *AbiImpl) JSNotImplemented(s string) {
	print(s)
	os.Exit(1)
}

func (a *AbiImpl) JSHandleEvent(e int32) {
	print("JSHandleEvent\n")
	os.Exit(1)
}

func (a *AbiImpl) TinygoNotImplemented(s string) {
	print(s)
	os.Exit(1)
}

func (a *AbiImpl) Exit(i int32) {
	os.Exit(int(i))
}

func (a *AbiImpl) Now() int64 {
	return time.Now().UnixNano()
}

func (a *AbiImpl) SetNow(_ int64, _ bool) {
	print("SetNow\n")
	os.Exit(1)
}

func (a *AbiImpl) Locate(team, service string) int64 {
	print("Locate")
	os.Exit(1)
	return int64(0)
}

func (a *AbiImpl) Dispatch(sid int64, method string, blob []byte) {
	print("Dispatch")
	os.Exit(1)
}
