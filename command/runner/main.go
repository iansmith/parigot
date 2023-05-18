package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
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
	config, err := runner.Parse(flag.Arg(0), flg)
	if err != nil {
		log.Fatalf("failed to parse configuration file %s: %v", flag.Arg(0), err)

	}
	ctx := pcontext.ServerGoContext(context.TODO(), "main")
	defer pcontext.Dump(ctx)

	// the deploy context creation also creates any needed nameservers
	deployCtx, err := sys.NewDeployContext(ctx, config)
	if err != nil {
		log.Fatalf("unable to create deploy context: %v", err)
	}

	if err := deployCtx.CreateAllProcess(ctx); err != nil {
		log.Fatalf("unable to create process: %v", err)
	}
	main, code := deployCtx.StartServer(ctx)
	if main == nil {
		if code != 0 {
			log.Printf("server startup returned error code %d", code)
			panic("os.Exit() with code " + fmt.Sprint(code))
		}
	}

	go func() {
		var buf bytes.Buffer
		for {
			buf.Reset()
			time.Sleep(15 * time.Second)
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
		if err != nil {
			pcontext.Logf(ctx, pcontext.Error, "could not start main program:%v", err)
			return
		}
		log.Printf("logging return code of %d from %s [%v]", code, mainProg, err)
		if code != 0 {
			pcontext.Logf(ctx, pcontext.Error, "main program '%s' exited with code %d", mainProg, code)
			return
		}
	}
	if len(main) > 1 {
		pcontext.Logf(ctx, pcontext.Info, "all main programs completed successfully")
	} else {
		pcontext.Logf(ctx, pcontext.Info, "main program completed successfully")
	}
	os.Exit(8)
}
