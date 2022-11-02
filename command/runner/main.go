package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/iansmith/parigot/sys"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
)

var libFile *string = flag.String("l", "", "the filename that has the list of wasm modules to load")

var libs = []string{}

func main() {
	flag.Parse()

	// xxx should be configuring this engine with my settings
	engine := wasmtime.NewEngine()

	// better error messages from the path than some ptr
	modToPath := make(map[*wasmtime.Module]string)

	evilHackIndex := 0
	// libraries
	var libs []*wasmtime.Module
	if *libFile != "" {
		var err error
		libs, err = readLibList(engine, modToPath)
		if err != nil {
			log.Fatalf("failed reading library list file (%s): %v", *libFile, err)
		}
	}
	libs = append(libs, walkArgs(engine, modToPath)...)
	evilHackIndex = len(libs)

	// check that we actually got something
	if len(libs) == 0 {
		log.Fatalf("unable to find any .wasm module files to load, pass filenames on the command line or use the -l option")
	}

	// store is shared by all the instances
	store := wasmtime.NewStore(engine)

	proc := []*sys.Process{}
	// create processes and check linkage for each one
	for i, lib := range libs {
		p, err := sys.NewProcessFromMod(store, lib, modToPath[lib])
		if err != nil {
			log.Fatalf("unable to create process from module (%s): %v", modToPath[lib], err)
		}
		proc = append(proc, p)
		if i == evilHackIndex {
			// this hack is a poor substitue for a topo sort which is what we should be doing
			// this tries to let all the processes that are in the lib file get started before
			// doin the rest... if the lib file is the standard lib, that almost might be sensible
			time.Sleep(500 * time.Millisecond) //ugh, cough, gack, choke, puke
		}
	}

	log.Printf("DONE")
	os.Exit(0)
}

func startup(store wasmtime.Storelike, module *wasmtime.Module, linkage []wasmtime.AsExtern) {
	// instance, err := wasmtime.NewInstance(store, module, linkage)
	// check(err)
	// ext := instance.GetExport(store, "mem")
	// mptr := uintptr(ext.Memory().Data(store))

	// // tell everybody about the memory
	// jsEnv.SetMemPtr(mptr)
	// runtimeEnv.SetMemPtr(mptr)
	// wasiEnv.SetMemPtr(mptr)
	// syscall.SetMemPtr(mptr)

	// start := instance.GetExport(store, "run")
	// if start == nil {
	// 	log.Fatalf("couldn't find start symbol")
	// }
	// f := start.Func()
	// _, err = f.Call(store, 0, 0)
	// check(err)
	// log.Printf("done with success!\n")

}

// func checkLinkage(wrappers map[string]*wasmtime.Func, module *wasmtime.Module) []wasmtime.AsExtern {
// 	linkage := []wasmtime.AsExtern{}
// 	for _, imp := range module.Imports() {
// 		n := "$$ANON$$"
// 		if imp.Name() != nil {
// 			n = *imp.Name()
// 		}
// 		importName := fmt.Sprintf("%s.%s", imp.Module(), n)
// 		ext, ok := wrappers[importName]
// 		if !ok {
// 			log.Printf("unable to find linkage for %s in module %s", importName, moduleMap[module])
// 			return nil
// 		} else {
// 			log.Printf("linked %s", importName)
// 			linkage = append(linkage, ext)
// 		}
// 	}
// 	return linkage
// }
