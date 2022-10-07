package main

import (
	"fmt"
	"github.com/iansmith/parigot/sys/abi_impl"
	"log"
	"os"
	"unsafe"

	"github.com/iansmith/parigot/abi/go/abi"
	"github.com/iansmith/parigot/abi/go/jspatch"
	"github.com/iansmith/parigot/abi/go/tinygopatch"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
)

var libs = []string{}

var memPtr uintptr

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("pass one wasm file on the command line")
	}
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)
	module, err := wasmtime.NewModuleFromFile(engine, os.Args[1])
	check(err)
	wrappers := generateWrappersForABI(store)
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
			log.Printf("linked %s", importName)
			linkage = append(linkage, ext)
		}
	}

	if linkFailed {
		os.Exit(1)
	}
	instance, err := wasmtime.NewInstance(store, module, linkage)
	check(err)
	ext := instance.GetExport(store, "memory")
	memPtr = uintptr(ext.Memory().Data(store))

	log.Printf("about to start")
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

	addABIToStoreExternref(store, result)

	result["wasi_snapshot_preview1.fd_write"] = wasmtime.WrapFunc(store, abi.FdWrite)
	result["env.syscall/js.valueGet"] = wasmtime.WrapFunc(store, jspatch.ValueGet)
	result["env.syscall/js.valuePrepareString"] = wasmtime.WrapFunc(store, jspatch.ValuePrepareString)
	result["env.syscall/js.valueLoadString"] = wasmtime.WrapFunc(store, jspatch.ValueLoadString)
	result["env.syscall/js.finalizeRef"] = wasmtime.WrapFunc(store, jspatch.FinalizeRef)
	result["env.syscall/js.stringVal"] = wasmtime.WrapFunc(store, jspatch.StringVal)
	result["env.syscall/js.valueSet"] = wasmtime.WrapFunc(store, jspatch.ValueSet)
	result["env.syscall/js.valueLength"] = wasmtime.WrapFunc(store, jspatch.ValueLength)
	result["env.syscall/js.valueIndex"] = wasmtime.WrapFunc(store, jspatch.ValueIndex)
	result["env.syscall/js.valueCall"] = wasmtime.WrapFunc(store, jspatch.ValueCall)
	result["env.syscall/js.valueSetIndex"] = wasmtime.WrapFunc(store, jspatch.ValueSetIndex)
	result["env.syscall/js.valueNew"] = wasmtime.WrapFunc(store, jspatch.ValueNew)

	result["env.runtime.ticks"] = wasmtime.WrapFunc(store, tinygopatch.Ticks)
	result["env.runtime.sleepTicks"] = wasmtime.WrapFunc(store, tinygopatch.SleepTicks)
	return result
}

func strConvert(mem uintptr, ptr int32, length int32) string {
	buf := make([]byte, length)
	for i := int32(0); i < length; i++ {
		b := (*byte)(unsafe.Pointer(memPtr + uintptr(ptr+i)))
		buf[i] = *b
	}
	s := string(buf)
	return s
}

func addABIToStoreExternref(store wasmtime.Storelike, linkage map[string]*wasmtime.Func) {

	linkage["parigot_abi.OutputString"] = wasmtime.WrapFunc(store, func(p0, p1 int32) {
		fmt.Print(strConvert(memPtr, p0, p1))
		return
	})
	linkage["parigot_abi.Now"] = wasmtime.WrapFunc(store, abi_impl.Now)
	linkage["parigot_abi.Exit"] = wasmtime.WrapFunc(store, abi_impl.Exit)
}