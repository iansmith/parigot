package abiimpl

import (
	"os"
)

type AbiImpl struct {
	memoryPtr *uintptr
}

func NewAbiImpl(memptr *uintptr) *AbiImpl {
	return &AbiImpl{
		memoryPtr: memptr,
	}
}
func (a *AbiImpl) GetMemPtr() uintptr {
	return *a.memoryPtr
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

func (a *AbiImpl) Now(retVal int32) {
	print("Now")
	os.Exit(1)
}

func (a *AbiImpl) SetNow(_ int64, _ bool) {
	print("SetNow\n")
	os.Exit(1)
}

func (a *AbiImpl) Locate(retVal int32, pkg, service string) {
	print("Locate")
	os.Exit(1)
}

func (a *AbiImpl) Dispatch(retval int32, sid int64, method string, blob []byte) {
	print("Dispatch")
	os.Exit(1)
}
