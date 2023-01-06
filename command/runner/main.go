package main

import (
	"flag"
	"log"

	"github.com/iansmith/parigot/command/runner/runner"
	"github.com/iansmith/parigot/sys"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatalf("unable to deploy application, no deployment config (toml format) provided")
	}
	if flag.NArg() > 1 {
		log.Fatalf("unable to deploy application, too many deployment configuration files provided (%d)", flag.NArg())
	}
	config, err := runner.Parse(flag.Arg(0))
	if err != nil {
		log.Fatalf("failed to parse configuration file %s: %v", flag.Arg(0), err)
	}
	log.Printf("config is ok:%s", flag.Arg(0))
	ctx, err := runner.NewContext(config)
	if err != nil {
		log.Fatalf("unable to create context: %v", err)
	}
	if err := ctx.CreateProcess(); err != nil {
		log.Fatalf("unable to create process: %v", err)
	}

	sys.InitNameServer(ctx.NotifyCh, !config.Remote, config.Remote)
	// This go routine's only purpose is to accept run requests from user programs running in a local
	// configuration (all in one process); it's called the run reader.
	//
	// Every process calls Run(). The run impl inside the kernel then hits RunNotify() which is
	// a noop if we are in the network case.  In the local case, the RunNotify() of the localNameserver
	// shoves it's information (sys.KeyNSPair) through the channel (notifyCh) and returns.
	//
	// This code then which is waiting on the notification of ANY process notifying it that
	// it has reached Run() and is now blocked (via RunBlock()) calls RunIfReady on whoever
	// sent it the notice via the notifyCh.
	//
	// RunIfReady() checks to see if any processes are ready to run (given their requires and
	// and exports) and if there are one or more ready, the name server ultimately will end up
	// calling Process.Run().  Process.Run() uses the runCh to "let go" of a process that had
	// previously called RunBlock().  Note that the process receiving the "let go" message can't
	// be also the sender.  The implication is that a single server running with any Require()
	// will just stop, which is what you'd expect.
	go func() {
		for {
			pair := <-ctx.NotifyCh
			log.Printf("calling nameserver.Run on %s", pair.Key)
			pair.NameServer.RunIfReady(pair.Key)
		}
	}()

	if code := ctx.Start(); code != 0 {
		log.Printf("main function returned error code %d", code)
	}

	log.Printf("DONE")
}

// create processes and check linkage for each user program
// for _, lib := range libs {
// 	store := wasmtime.NewStore(engine)
// 	p, err := sys.NewProcessFromMod(store, lib, modToPath[lib], rs)
// 	if err != nil {
// 		log.Fatalf("unable to create process from module (%s): %v", modToPath[lib], err)
// 	}
// 	maxModules++
// 	runnerPrint("MAIN ", "starting goroutine for process %s", p)
// 	go func(p *sys.Process) {
// 		p.Start()
// 		return
// 	}(p)
// 	proc = append(proc, p)
// }

// // we are the periodic check for things getting up ok
// everbodyStarted := false
// iter := 0
// for iter < secondsBeforeStartupFailed {
// 	startedCount := 0
// 	for _, p := range proc {
// 		if p.ReachedStart() {
// 			startedCount++
// 		}
// 	}
// 	if startedCount == len(proc) {
// 		everbodyStarted = true
// 		break
// 	}
// 	runnerPrint("MAIN ", "startedCount %d, len %d", startedCount, len(proc))
// 	time.Sleep(100 * time.Millisecond)
// 	iter++
// }
// if !everbodyStarted {
// 	log.Printf("we waited %d seconds, but some processes failed to start...", iter)
// 	for _, p := range proc {
// 		if !p.ReachedStart() {
// 			log.Printf("\t%s", p)
// 		}
// 	}
// 	log.Fatalf("aborting because processes failed to start")
// } else {
// 	everybodyExited := true
// 	for _, p := range proc {
// 		if p.IsWaiter() && !p.Exited() {
// 			everybodyExited = false
// 			break
// 		}
// 	}
// 	if !everybodyExited {
// 		localInfo, netInfo := sys.StartFailedInfo()
// 		if localInfo != "" || netInfo != "" {
// 			if localInfo != "" {
// 				log.Printf("dependency problem:\n%s", localInfo)
// 			}
// 			if netInfo != "" {
// 				log.Printf("dependency problem:\n%s", netInfo)
// 			}
// 			log.Fatalf("aborting due to dependency issues")
// 		}
// 	}
// }
// for {
// 	allDead := true
// 	for _, p := range proc {
// 		//print(fmt.Sprintf("proc %s: waiter? %v exited? %v server? %v\n", p, p.IsWaiter(), p.Exited(), p.IsServer()))
// 		if (p.IsWaiter() && !p.Exited()) || p.IsServer() {
// 			allDead = false
// 			break
// 		}
// 	}
// 	if allDead {
// 		log.Printf("all processes have exited, finished")
// 		os.Exit(0)
// 	}
// 	time.Sleep(250 * time.Millisecond)
// }
//}
