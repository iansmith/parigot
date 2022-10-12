package main

import (
	"fmt"
	"github.com/iansmith/parigot/sys/abiimpl"
	"log"
	"os"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
	"github.com/iansmith/parigot/abi/jspatch"
	"github.com/iansmith/parigot/abi/tinygopatch"
	"github.com/iansmith/parigot/command/runner/g"
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

// temporary while we are getting rid of JS linkage
func jsPatch(store wasmtime.Storelike, result map[string]*wasmtime.Func) {
	result["env.syscall/js.valueSetIndex"] = wasmtime.WrapFunc(store, jspatch.ValueSetIndex)
	//result["wasi_snapshot_preview1.fd_write"] = wasmtime.WrapFunc(store, abi.FdWrite)
	//result["env.syscall/js.valueGet"] = wasmtime.WrapFunc(store, jspatch.ValueGet)
	//result["env.syscall/js.valuePrepareString"] = wasmtime.WrapFunc(store, jspatch.ValuePrepareString)
	//result["env.syscall/js.valueLoadString"] = wasmtime.WrapFunc(store, jspatch.ValueLoadString)
	//result["env.syscall/js.finalizeRef"] = wasmtime.WrapFunc(store, jspatch.FinalizeRef)
	//result["env.syscall/js.stringVal"] = wasmtime.WrapFunc(store, jspatch.StringVal)
	//result["env.syscall/js.valueSet"] = wasmtime.WrapFunc(store, jspatch.ValueSet)
	//result["env.syscall/js.valueLength"] = wasmtime.WrapFunc(store, jspatch.ValueLength)
	//result["env.syscall/js.valueIndex"] = wasmtime.WrapFunc(store, jspatch.ValueIndex)
	//result["env.syscall/js.valueCall"] = wasmtime.WrapFunc(store, jspatch.ValueCall)
	//result["env.syscall/js.valueNew"] = wasmtime.WrapFunc(store, jspatch.ValueNew)
}

// temporary while we are getting rid of runtime of tinygo
func tinygoPatch(store wasmtime.Storelike, result map[string]*wasmtime.Func) {
	result["env.runtime.ticks"] = wasmtime.WrapFunc(store, tinygopatch.Ticks)
	result["env.runtime.sleepTicks"] = wasmtime.WrapFunc(store, tinygopatch.SleepTicks)
}

type abiWrapperForMemPtr struct {
	*abiimpl.AbiImpl
}

func (a *abiWrapperForMemPtr) GetMemPtr() uintptr {
	return memPtr
}

func generateWrappersForABI(store wasmtime.Storelike) map[string]*wasmtime.Func {
	var result = make(map[string]*wasmtime.Func)
	impl := &abiimpl.AbiImpl{}
	wrapper := &abiWrapperForMemPtr{AbiImpl: impl}
	g.WasmTimeWrapABI(wrapper, store, result)
	jsPatch(store, result)
	tinygoPatch(store, result)
	return result
}
