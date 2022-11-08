package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/iansmith/parigot/sys"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
)

// Flip this flag for more detailed output from the runner.
var runnerVerbose = true

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

	// this go routine's only purpose is to accept run requests from user programs
	go nameServer.RunReader()

	// create processes and check linkage for each user program
	for _, lib := range libs {
		p, err := sys.NewProcessFromMod(store, lib, modToPath[lib], nameServer)
		if err != nil {
			log.Fatalf("unable to create process from module (%s): %v", modToPath[lib], err)
		}
		maxModules++
		runnerPrint("MAIN ", "starting goroutine for process %s", p)
		go startProcess(p)
		proc = append(proc, p)
	}
	// we are the periodic check for things getting up ok
	everbodyStarted := false
	iter := 0
	for iter < 3 {
		startedCount := 0
		for _, p := range proc {
			if p.ReachedStart() {
				startedCount++
			}
		}
		if startedCount == len(proc) {
			break
		}
		runnerPrint("MAIN ", "startedCount %d, len %d", startedCount, len(proc))
		time.Sleep(1000 * time.Millisecond)
		iter++
	}
	if !everbodyStarted {
		log.Printf("we waited %d seconds, but some processes failed to start...", iter)
		for _, p := range proc {
			if !p.ReachedStart() {
				log.Printf("\t%s", p)
			}
		}
		log.Fatalf("aborting")
	} else {
		if nameServer.WaitingToRun() > 0 {
			nameServer.SendAbortMessage()
		}
		log.Printf("was not able to get all processes started due to export/require problems")
		loop := nameServer.GetLoopContent()
		dead := nameServer.GetDeadNodeContent()
		if loop != "" {
			loop = strings.Replace(loop, ";", "\n", -1)
			log.Printf("Loop discovered in the dependencies\n%s\n", loop)
		}
		if dead != "" {
			dead = strings.Replace(dead, ";", "\n", -1)
			log.Printf("Dead processes are processes that cannot start because no other process exports what they require\n%s\n", dead)
		}
		log.Fatalf("aborting due to export/require problems")
	}
	for {
		allDead := true
		for _, p := range proc {
			if p.IsWaiter() && !p.Exited() {
				allDead = false
				break
			}
		}
		if allDead {
			log.Printf("all processes have exited, finished")
		}
		time.Sleep(250 * time.Millisecond)
	}
}

func startProcess(p *sys.Process) {
	p.Start()
}

func runnerPrint(method, spec string, arg ...interface{}) {
	if runnerVerbose {
		part1 := fmt.Sprintf("RUNNER:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}
