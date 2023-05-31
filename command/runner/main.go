package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	"github.com/iansmith/parigot/command/runner/runner"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/sys"
)

var testMode *bool = flag.Bool("t", false, "turns testmode on, implies running services marked 'Test' in deploy config")
var remote *bool = flag.Bool("r", false, "run all services in separate address spaces")

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatalf("unable to deploy application, no deployment config (toml format) provided")
	}
	if flag.NArg() > 1 {
		log.Fatalf("unable to deploy application, too many deployment configuration files provided (%d)", flag.NArg())
	}
	flg := &runner.DeployFlag{
		TestMode: *testMode,
		Remote:   *remote,
	}
	defer func() {
		if r := recover(); r != nil {
			print("runner crashed:", fmt.Sprintf("%T %v", r, r), "\n")
			debug.PrintStack()
		}
	}()
	config, err := runner.Parse(flag.Arg(0), flg)
	if err != nil {
		log.Fatalf("failed to parse configuration file %s: %v", flag.Arg(0), err)

	}
	ctx := pcontext.NewContextWithContainer(context.Background(), "runner:main")
	ctx = pcontext.CallTo(pcontext.InternalParigot(ctx), "main")
	defer pcontext.Dump(ctx)

	// the deploy context creation also creates any needed nameservers
	deployCtx, err := sys.NewDeployContext(ctx, config)
	if err != nil {
		log.Fatalf("unable to create deploy context: %v", err)
	}

	if err := deployCtx.CreateAllProcess(ctx); err != nil {
		log.Fatalf("unable to create process in main: %v", err)
	}

	main, code := deployCtx.StartServer(pcontext.CallTo(ctx, "StartServer"))
	if main == nil {
		if code != 0 {
			log.Printf("server startup returned error code %d", code)
			panic("os.Exit() with code " + fmt.Sprint(code))
		}
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(context context.Context) {
		for range c {
			pcontext.Dump(context)
			os.Exit(2)
		}
	}(ctx)

	go func() {
		var buf bytes.Buffer
		for {
			buf.Reset()
			time.Sleep(60 * time.Second)
			deployCtx.Process().Range(func(keyAny, valueAny any) bool {
				key := keyAny.(string)
				proc := valueAny.(*sys.Process)
				buf.WriteString(fmt.Sprintf("process %20s:block=%v,run=%v,req met=%v, exited=%v\n",
					key, proc.ReachedRunBlock(), proc.Running(), proc.RequirementsMet(), proc.Exited()))
				return true
			})
			print("periodic check:-----------\n", buf.String(), "\n")
		}
	}()

	for _, mainProg := range main {
		code, err := deployCtx.StartMain(ctx, mainProg)
		if code == 253 {
			pcontext.Fatalf(ctx, "host side code failed (usually a panic) in execution of  program %s", mainProg)
		} else if code != 0 {
			pcontext.Infof(ctx, "main exited from %s with code %d and error %s", mainProg, code, fmt.Sprint(err))
		}
		return
	}
	if len(main) > 1 {
		pcontext.Logf(ctx, pcontext.Info,
			"all main programs completed successfully")
	} else {
		pcontext.Logf(ctx, pcontext.Info, "main program completed successfully")
	}
	os.Exit(8)
}
