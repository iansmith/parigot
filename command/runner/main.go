package main

import (
	"flag"
	"log"
	"os"
	"sync"
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
	var wg sync.WaitGroup
	// create processes and check linkage for each one
	for i, lib := range libs {
		p, err := sys.NewProcessFromMod(store, lib, modToPath[lib], nameServer)
		if err != nil {
			log.Fatalf("unable to create process from module (%s): %v", modToPath[lib], err)
		}
		wg.Add(1)
		go startProcess(p)
		proc = append(proc, p)
		if i == evilHackIndex {
			// this hack is a poor substitue for a topo sort which is what we should be doing
			// this tries to let all the processes that are in the lib file get started before
			// doin the rest... if the lib file is the standard lib, that almost might be sensible
			time.Sleep(500 * time.Millisecond) //ugh, cough, gack, choke, puke
		}
	}
	wg.Wait()
	log.Printf("all go routines have exited")
	os.Exit(0)
}

func startProcess(p *sys.Process) {
	p.Start()
}
