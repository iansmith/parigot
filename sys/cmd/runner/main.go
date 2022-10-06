package main

import (
	"fmt"
	wasmtime "github.com/bytecodealliance/wasmtime-go"
	"github.com/iansmith/parigot/abi/go/abi"
	"log"
	"os"
	"unsafe"
)

var libs = []string{}

var memPtr unsafe.Pointer

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("pass one wasm file on the command line")
	}
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)
	module, err := wasmtime.NewModuleFromFile(engine, os.Args[1])
	check(err)
	wrappers := generateWrappersForABI(store)
	//for _, exp := range module.Exports() {
	//	log.Printf("exp: %s", exp.Name())
	//}
	linkFailed := false
	linkage := []wasmtime.AsExtern{}
	for _, imp := range module.Imports() {
		n := "$$ANON$$"
		if imp.Name() != nil {
			n = *imp.Name()
		}
		importName := fmt.Sprintf("%s.%s", imp.Module(), n)
		ext, ok := wrappers[importName]
		if !ok {
			log.Printf("unable to find linkage for %s", importName)
			linkFailed = true
		} else {
			linkage = append(linkage, ext)
		}
	}

	if linkFailed {
		os.Exit(1)
	}
	instance, err := wasmtime.NewInstance(store, module, linkage)
	check(err)
	ext := instance.GetExport(store, "memory")
	memPtr = ext.Memory().Data(store)

	start := instance.GetExport(store, "_start").Func()
	_, err = start.Call(store)
	check(err)
	print("done with success!\n")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func generateWrappersForABI(store wasmtime.Storelike) map[string]*wasmtime.Func {
	var result = make(map[string]*wasmtime.Func)
	result["parigot_abi.TinyGoNotImplemented"] = wasmtime.WrapFunc(store, abi.TinyGoNotImplemented)
	result["parigot_abi.JSHandleEvent"] = wasmtime.WrapFunc(store, abi.JSHandleEvent)
	result["parigot_abi.JSNotImplemented"] = wasmtime.WrapFunc(store, abi.JSNotImplemented)
	result["parigot_abi.SetNow"] = wasmtime.WrapFunc(store, abi.SetNow)
	//result["parigot_abi.Now"] = wasmtime.WrapFunc(store, abi.Now)
	result["parigot_abi.NowConvert"] = wasmtime.WrapFunc(store, abi.NowConvert)
	//result["parigot_abi.OutputString"] = wasmtime.WrapFunc(store, abi.OutputString)
	//result["parigot_abi.OutputString"] = wasmtime.WrapFunc(store, abi.OutputStringConvert)
	result["parigot_abi.OutputString"] = wasmtime.WrapFunc(store, func(a int32, b int32) {
		print("a=", a, " b=", b, "\n")
		ptr := uintptr(memPtr)
		for i := int32(0); i < b; i++ {
			addr := ptr + uintptr(a+i)
			p := (*byte)(unsafe.Pointer(addr))
			fmt.Printf("%c", *p)
		}
	})

	result["parigot_abi.Exit"] = wasmtime.WrapFunc(store, abi.Exit)
	result["wasi_snapshot_preview1.fd_write"] = wasmtime.WrapFunc(store, abi.FdWrite)

	result["env.syscall/js.valueGet"] = wasmtime.WrapFunc(store, abi.ValueGet)
	result["env.syscall/js.valuePrepareString"] = wasmtime.WrapFunc(store, abi.ValuePrepareString)
	result["env.syscall/js.valueLoadString"] = wasmtime.WrapFunc(store, abi.ValueLoadString)
	result["env.syscall/js.finalizeRef"] = wasmtime.WrapFunc(store, abi.FinalizeRef)
	result["env.syscall/js.stringVal"] = wasmtime.WrapFunc(store, abi.StringVal)
	result["env.syscall/js.valueSet"] = wasmtime.WrapFunc(store, abi.ValueSet)
	result["env.syscall/js.valueLength"] = wasmtime.WrapFunc(store, abi.ValueLength)
	result["env.syscall/js.valueIndex"] = wasmtime.WrapFunc(store, abi.ValueIndex)
	result["env.syscall/js.valueCall"] = wasmtime.WrapFunc(store, abi.ValueCall)
	result["env.syscall/js.valueSetIndex"] = wasmtime.WrapFunc(store, abi.ValueSetIndex)
	result["env.syscall/js.valueNew"] = wasmtime.WrapFunc(store, abi.ValueNew)

	result["env.runtime.ticks"] = wasmtime.WrapFunc(store, abi.Ticks)
	result["env.runtime.sleepTicks"] = wasmtime.WrapFunc(store, abi.SleepTicks)
	return result
}
