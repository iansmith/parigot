package abiimpl

import (
	"log"
	"os"

	"github.com/iansmith/parigot/lib"
)

type AbiImpl struct {
}

func NewAbiImpl() *AbiImpl {
	return &AbiImpl{}
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

var packageRegistry = make(map[string]map[string]lib.ServiceId)
var serviceCounter = 1

func (a *AbiImpl) Register(retVal int32, protoPackage string, service string) {
	serviceRegistry, ok := packageRegistry[protoPackage]
	if !ok {
		serviceRegistry = make(map[string]lib.ServiceId)
	}
	sid := lib.NewServiceFromInt(byte(serviceCounter + 1))
	serviceCounter++
	if serviceCounter > 255 {
		panic("service counter exceeded maximum of 255")
	}
	serviceRegistry[service] = sid
	log.Printf("Register: What's the retval? %x", retVal)
	os.Exit(1)
}

func (a *AbiImpl) Locate(retVal int32, pkg, service string) {
	log.Printf("Locate %s,%s->%x", pkg, service, retVal)
	os.Exit(1)
}

func DebugPrint(ct int32) {
	log.Printf("----DebugPrint %d", ct)
}

func (a *AbiImpl) Dispatch(retval int32, sid int64 /*xxx*/, method string, blob []byte) {
	print("Dispatch")
	os.Exit(1)
}