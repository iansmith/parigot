package main

import (
	"flag"
	"log"
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

	// the singular nameserver
	nameServer := sys.NewNameServer()

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

	// check that we actually got something
	if len(libs) == 0 {
		log.Fatalf("unable to find any .wasm module files to load, pass filenames on the command line or use the -l option")
	}

	// store is shared by all the instances
	store := wasmtime.NewStore(engine)

	proc := []*sys.Process{}
	maxModules := 0
	// create processes and check linkage for each one
	for _, lib := range libs {
		p, err := sys.NewProcessFromMod(store, lib, modToPath[lib], nameServer)
		if err != nil {
			log.Fatalf("unable to create process from module (%s): %v", modToPath[lib], err)
		}
		maxModules++
		go startProcess(p)
		proc = append(proc, p)
	}
	// we are the period check for being done
	for {
		if nameServer.WaitingToRun() == maxModules {
			// we have a loop and need to abort
			nameServer.SendLoopMessage()
			log.Fatalf("unable to run all the modules, dependency loop found")
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func startProcess(p *sys.Process) {
	p.Start()
}
