package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/iansmith/parigot/sys"

	wasmtime "github.com/bytecodealliance/wasmtime-go/v3"
)

// Flip this flag for more detailed output from the runner.
var runnerVerbose = false

var libFile *string = flag.String("f", "", "the filename that has the list of wasm modules to load")
var remote *bool = flag.Bool("r", false, "all services will use remote; use this flag for a docker swarm of microservices")
var libs = []string{}

var secondsBeforeStartupFailed = 15

func main() {
	flag.Parse()

	// xxx should be configuring this engine with my settings
	engine := wasmtime.NewEngine()

	// better error messages from the path than some ptr
	modToPath := make(map[*wasmtime.Module]string)

	// init the nameservers
	notifyCh := make(chan *sys.KeyNSPair)
	sys.InitNameServer(notifyCh, !(*remote), *remote)

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

	proc := []*sys.Process{}
	maxModules := 0

	// This go routine's only purpose is to accept run requests from user programs;
	// it's called the run reader.
	// This is called by the process saying "Run" which means "Block until my prereqs are ready".
	// Each time a process calls it and its nameserver get passed to here, and we
	// to use this indirect structure so the process can block waiting on a channel.
	go func() {
		for {
			pair := <-notifyCh
			runnerPrint("RUNREADER ", "calling nameserver.Run on %s", pair.Key)
			pair.NameServer.RunIfReady(pair.Key)
		}
	}()

	all := []string{}
	for _, lib := range libs {
		all = append(all, modToPath[lib])
	}

	rs := sys.NewRemoteSpec(nil, all)
	if *remote {
		rs = sys.NewRemoteSpec(all, nil)
	}

	// create processes and check linkage for each user program
	for _, lib := range libs {
		store := wasmtime.NewStore(engine)
		p, err := sys.NewProcessFromMod(store, lib, modToPath[lib], rs)
		if err != nil {
			log.Fatalf("unable to create process from module (%s): %v", modToPath[lib], err)
		}
		maxModules++
		runnerPrint("MAIN ", "starting goroutine for process %s", p)
		go func(p *sys.Process) {
			p.Start()
			p.Exit()
			return
		}(p)
		proc = append(proc, p)
	}

	// we are the periodic check for things getting up ok
	everbodyStarted := false
	iter := 0
	for iter < secondsBeforeStartupFailed {
		startedCount := 0
		for _, p := range proc {
			if p.ReachedStart() {
				startedCount++
			}
		}
		if startedCount == len(proc) {
			everbodyStarted = true
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
		log.Fatalf("aborting because processes failed to start")
	} else {
		everybodyExited := true
		for _, p := range proc {
			if p.IsWaiter() && !p.Exited() {
				everybodyExited = false
				break
			}
		}
		if !everybodyExited {
			localInfo, netInfo := sys.StartFailedInfo()
			if localInfo != "" || netInfo != "" {
				if localInfo != "" {
					log.Printf("dependency problem:\n%s", localInfo)
				}
				if netInfo != "" {
					log.Printf("dependency problem:\n%s", netInfo)
				}
				log.Fatalf("aborting due to dependency issues")
			}
		}
	}
	for {
		allDead := true
		for _, p := range proc {
			print(fmt.Sprintf("proc %s: waiter? %v exited? %v server? %v\n", p, p.IsWaiter(), p.Exited(), p.IsServer()))
			if (p.IsWaiter() && !p.Exited()) || p.IsServer() {
				allDead = false
				break
			}
		}
		if allDead {
			log.Printf("all processes have exited, finished")
			os.Exit(0)
		}
		time.Sleep(250 * time.Millisecond)
	}
}

func runnerPrint(method, spec string, arg ...interface{}) {
	if runnerVerbose {
		part1 := fmt.Sprintf("RUNNER:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}
