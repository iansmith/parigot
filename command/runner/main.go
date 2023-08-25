package main

// at no point can this package or anything it depends on import anything in api/guest
import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/iansmith/parigot/command/runner/runner"
	"github.com/iansmith/parigot/sys"
	"github.com/iansmith/parigot/sys/kernel"
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
	}

	config, err := runner.Parse(flag.Arg(0), flg)
	if err != nil {
		log.Fatalf("failed to parse configuration file %s: %v", flag.Arg(0), err)

	}
	ok := false
	kernel.K, ok = kernel.InitSingle()
	if !ok {
		log.Fatalf("unable to create kernel, aborting")
	}

	// the deploy context creation also creates any needed nameservers
	deployCtx, err := sys.NewDeployContext(config)
	if err != nil {
		log.Fatalf("unable to create deploy context: %v", err)
	}

	if err := deployCtx.CreateAllProcess(); err != nil {
		log.Fatalf("unable to create process in main: %v", err)
	}

	main, code := deployCtx.StartServer(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		sig := <-c
		log.Printf("received signal %+v", sig)
		os.Exit(1)
	}()

	if main == nil {
		if code != 0 {
			log.Printf("servers are running")
		}
		for {
			// a year
			time.Sleep(8760 * time.Hour)
		}
	} else {
		for _, mainProg := range main {
			code, err := deployCtx.StartMain(mainProg)
			if code == 253 && err == nil {
				//pcontext.Fatalf(ctx, "code failed (usually a panic) in execution of  program %s (code %d) -- can be host or guest", mainProg, code)
			} else if code != 0 {
				slog.Info("main exited", "name", mainProg, "code", code, "error?", err != nil)
			}
		}
		if len(main) > 1 {
			log.Printf(
				"all main programs completed successfully")
		} else {
			log.Printf("main program '%s' completed successfully", main[0])
		}
	}
}
