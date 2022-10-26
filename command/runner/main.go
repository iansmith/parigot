package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"time"
	"unsafe"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
	"github.com/iansmith/parigot/sys/kernel"
	"github.com/iansmith/parigot/sys/kernel/jspatch"
	"github.com/iansmith/parigot/sys/kernel/tinygopatch"
)

var libs = []string{}

var memPtr *uintptr

func main() {
	type x struct{}
	log.Printf("??? %s", reflect.TypeOf(x{}).PkgPath())
	mainNormal()
}
func mainTest() {
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)
	module, err := wasmtime.NewModuleFromFile(engine, os.Args[1])
	check(err)
	wrappers := make(map[string]*wasmtime.Func)
	testWrappers(store, wrappers)
	linkage := checkLinkage(wrappers, module)
	if linkage == nil {
		os.Exit(1)
	}
	startup(store, module, linkage)
}
func checkAddr(p int32, l int32) {
	log.Printf("check addr %x,%d", p, l)
}

func testWrappers(store wasmtime.Storelike, result map[string]*wasmtime.Func) {
	result["env.checkAddr"] = wasmtime.WrapFunc(store, checkAddr)
	result["wasi_snapshot_preview1.fd_write"] = wasmtime.WrapFunc(store,
		func(x int32, y int32, z int32, a int32) int32 {
			log.Printf("fd_write %x,%x,%x,%x", x, y, z, a)
			return 0
		})
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

	// add functions
	memPtr = new(uintptr)
	jspatch.SetMemPtr(memPtr)
	_ /*impl*/ = abiimpl.NewAbiImpl(memPtr)
	///xxx fixme
	//g.SetCaller(impl)
	//g.WasmTimeWrapABI(impl, store, wrappers)
	tinygoPatch(store, wrappers)
	jsPatch(store, wrappers)

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
	*memPtr = uintptr(ext.Memory().Data(store))

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

// temporary while we are getting rid of JS linkage
func jsPatch(store wasmtime.Storelike, result map[string]*wasmtime.Func) {
	result["env.syscall/js.valueSetIndex"] = wasmtime.WrapFunc(store, jspatch.ValueSetIndex)
	result["wasi_snapshot_preview1.fd_write"] = wasmtime.WrapFunc(store, tinygopatch.WasiWriteFd)
	result["wasi_snapshot_preview1.proc_exit"] = wasmtime.WrapFunc(store, tinygopatch.WasiProcExit)
	result["env.syscall/js.valueGet"] = wasmtime.WrapFunc(store, jspatch.ValueGet)
	result["env.syscall/js.valuePrepareString"] = wasmtime.WrapFunc(store, jspatch.ValuePrepareString)
	result["env.syscall/js.valueLoadString"] = wasmtime.WrapFunc(store, jspatch.ValueLoadString)
	result["env.syscall/js.finalizeRef"] = wasmtime.WrapFunc(store, jspatch.FinalizeRef)
	result["env.syscall/js.stringVal"] = wasmtime.WrapFunc(store, jspatch.StringVal)
	result["env.syscall/js.valueSet"] = wasmtime.WrapFunc(store, jspatch.ValueSet)
	//result["env.syscall/js.valueLength"] = wasmtime.WrapFunc(store, jspatch.ValueLength)
	//result["env.syscall/js.valueIndex"] = wasmtime.WrapFunc(store, jspatch.ValueIndex)
	result["env.syscall/js.valueCall"] = wasmtime.WrapFunc(store, jspatch.ValueCall)
	result["env.syscall/js.valueNew"] = wasmtime.WrapFunc(store, jspatch.ValueNew)

	result["parigot.debugprint"] = wasmtime.WrapFunc(store, abiimpl.DebugPrint)
	result["go.debug"] = wasmtime.WrapFunc(store, func(x int32) {
		print("got a go debug call", x, "\n")
	})
	result["go.runtime.resetMemoryDataView"] = wasmtime.WrapFunc(store, func(x int32) {
		print("got a resetMemoryDataView call", x, "\n")
	})
	result["go.runtime.wasmExit"] = wasmtime.WrapFunc(store, func(sp int32) {
		log.Printf("wasmExit: %d", jspatch.GetInt32(*memPtr, sp+8))
	})
	result["go.runtime.wasmWrite"] = wasmtime.WrapFunc(store, func(sp int32) {
		_ = jspatch.GetInt64(*memPtr, sp+8)
		p := jspatch.GetInt64(*memPtr, sp+16)
		n := jspatch.GetInt32(*memPtr, sp+24)
		content := make([]byte, n)
		ptr := (*memPtr + uintptr(p))
		for i := int32(0); i < n; i++ {
			content[i] = *((*byte)(unsafe.Pointer(ptr + uintptr(i))))
		}
		//log.Printf("wasm write on file descriptor %d with %x and len %d: %s", fd, p, n, string(content))
		fmt.Printf("%s", string(content))
	})
	result["go.runtime.nanotime1"] = wasmtime.WrapFunc(store, func(sp int32) {
		jspatch.SetInt64(*memPtr, sp+8, time.Now().UnixNano())
	})
	result["go.runtime.walltime"] = wasmtime.WrapFunc(store, func(x int32) {
		print("got a walltime call", x, "\n")
	})
	result["go.runtime.scheduleTimeoutEvent"] = wasmtime.WrapFunc(store, func(x int32) {
		print("got a schedule timeout event call", x, "\n")
	})
	result["go.runtime.clearTimeoutEvent"] = wasmtime.WrapFunc(store, func(x int32) {
		print("got a CLEAR timeout event call", x, "\n")
	})
	result["go.runtime.getRandomData"] = wasmtime.WrapFunc(store, func(sp int32) {
		b := jspatch.LoadSlice(*memPtr, sp+8)
		_, _ = rand.Read(b) //docs say no returned error
	})
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
	return *memPtr
}
