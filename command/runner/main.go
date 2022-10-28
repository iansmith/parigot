package main

import (
	"fmt"
	"github.com/iansmith/parigot/sys"
	"log"
	"os"
	"reflect"

	"github.com/iansmith/parigot/sys/jspatch"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
)

var libs = []string{}

var jsEnv *jspatch.JSPatch
var wasiEnv *jspatch.WasiPatch
var runtimeEnv *jspatch.RuntimePatch
var syscall *sys.SysCall

func main() {
	type x struct{}
	log.Printf("??? %s", reflect.TypeOf(x{}).PkgPath())
	mainNormal()
}

func mainNormal() {
	if len(os.Args) != 2 {
		log.Fatalf("pass one wasm file on the command line")
	}
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)
	module, err := wasmtime.NewModuleFromFile(engine, os.Args[1])
	check(err)
	wrappers := make(map[string]*wasmtime.Func)

	jsEnv = jspatch.NewJSPatch()
	wasiEnv = jspatch.NewWasiPatch()
	runtimeEnv = jspatch.NewRuntimePatch()
	syscall = sys.NewSysCall()

	supportedFunctions(store, wrappers, syscall)
	// check that everything linked
	linkage := checkLinkage(wrappers, module)
	if linkage == nil {
		os.Exit(1)
	}
	startup(store, module, linkage)
}

func startup(store wasmtime.Storelike, module *wasmtime.Module, linkage []wasmtime.AsExtern) {
	instance, err := wasmtime.NewInstance(store, module, linkage)
	check(err)
	ext := instance.GetExport(store, "mem")
	mptr := uintptr(ext.Memory().Data(store))

	// tell everybody about the memory
	jsEnv.SetMemPtr(mptr)
	runtimeEnv.SetMemPtr(mptr)
	wasiEnv.SetMemPtr(mptr)
	syscall.SetMemPtr(mptr)

	start := instance.GetExport(store, "run")
	if start == nil {
		log.Fatalf("couldn't find start symbol")
	}
	f := start.Func()
	_, err = f.Call(store, 0, 0)
	check(err)
	log.Printf("done with success!\n")

}
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func checkLinkage(wrappers map[string]*wasmtime.Func, module *wasmtime.Module) []wasmtime.AsExtern {
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
			return nil
		} else {
			log.Printf("linked %s", importName)
			linkage = append(linkage, ext)
		}
	}
	return linkage
}

func supportedFunctions(store wasmtime.Storelike,
	result map[string]*wasmtime.Func,
	syscallImpl *sys.SysCall) {
	result["env.syscall/js.valueSetIndex"] = wasmtime.WrapFunc(store, jsEnv.ValueSetIndex)
	result["go.syscall/js.valueGet"] = wasmtime.WrapFunc(store, jsEnv.ValueGet)
	result["env.syscall/js.valuePrepareString"] = wasmtime.WrapFunc(store, jsEnv.ValuePrepareString)
	result["env.syscall/js.valueLoadString"] = wasmtime.WrapFunc(store, jsEnv.ValueLoadString)
	result["go.syscall/js.finalizeRef"] = wasmtime.WrapFunc(store, jsEnv.FinalizeRef)
	result["go.syscall/js.stringVal"] = wasmtime.WrapFunc(store, jsEnv.StringVal)
	result["go.syscall/js.valueSet"] = wasmtime.WrapFunc(store, jsEnv.ValueSet)
	result["go.syscall/js.valueDelete"] = wasmtime.WrapFunc(store, jsEnv.ValueDelete)
	result["go.syscall/js.valueIndex"] = wasmtime.WrapFunc(store, jsEnv.ValueIndex)
	result["go.syscall/js.valueSetIndex"] = wasmtime.WrapFunc(store, jsEnv.ValueSetIndex)
	result["go.syscall/js.valueLength"] = wasmtime.WrapFunc(store, jsEnv.ValueLength)
	result["go.syscall/js.valuePrepareString"] = wasmtime.WrapFunc(store, jsEnv.ValuePrepareString)
	result["go.syscall/js.valueCall"] = wasmtime.WrapFunc(store, jsEnv.ValueCall)
	result["go.syscall/js.valueInvoke"] = wasmtime.WrapFunc(store, jsEnv.ValueInvoke)
	result["go.syscall/js.valueNew"] = wasmtime.WrapFunc(store, jsEnv.ValueNew)
	result["go.syscall/js.valueLoadString"] = wasmtime.WrapFunc(store, jsEnv.ValueLoadString)
	result["go.syscall/js.valueInstanceOf"] = wasmtime.WrapFunc(store, jsEnv.ValueInstanceOf)
	result["go.syscall/js.copyBytesToGo"] = wasmtime.WrapFunc(store, jsEnv.CopyBytesToGo)
	result["go.syscall/js.copyBytesToJS"] = wasmtime.WrapFunc(store, jsEnv.CopyBytesToJS)
	result["go.runtime.resetMemoryDataView"] = wasmtime.WrapFunc(store, runtimeEnv.ResetMemoryDataView)
	result["go.runtime.wasmExit"] = wasmtime.WrapFunc(store, wasiEnv.WasiExit)
	result["go.runtime.wasmWrite"] = wasmtime.WrapFunc(store, wasiEnv.WasiWrite)
	result["go.runtime.nanotime1"] = wasmtime.WrapFunc(store, runtimeEnv.Nanotime1)
	result["go.runtime.walltime"] = wasmtime.WrapFunc(store, runtimeEnv.WallTime)
	result["go.runtime.scheduleTimeoutEvent"] = wasmtime.WrapFunc(store, runtimeEnv.ScheduleTimeoutEvent)
	result["go.runtime.clearTimeoutEvent"] = wasmtime.WrapFunc(store, runtimeEnv.ClearTimeoutEvent)
	result["go.runtime.getRandomData"] = wasmtime.WrapFunc(store, runtimeEnv.GetRandomData)
	result["parigot.debugprint"] = wasmtime.WrapFunc(store, sys.DebugPrint)
	result["go.debug"] = wasmtime.WrapFunc(store, runtimeEnv.GoDebug)

	//system calls
	result["go.parigot.locate_"] = wasmtime.WrapFunc(store, syscallImpl.Locate)
	result["go.parigot.register_"] = wasmtime.WrapFunc(store, syscallImpl.Register)
	result["go.parigot.exit_"] = wasmtime.WrapFunc(store, syscallImpl.Exit)
	result["go.parigot.dispatch_"] = wasmtime.WrapFunc(store, syscallImpl.Dispatch)
}
