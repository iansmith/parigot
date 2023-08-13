package main

// at no point can this package or anything it depends on import anything in api/guest
import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/iansmith/parigot/api/plugin/syscall/kernel"
	"github.com/iansmith/parigot/api/plugin/syscall/wheeler"
	"github.com/iansmith/parigot/command/runner/runner"
	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/syscall/v1"
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
	}
	ctx := pcontext.NewContextWithContainer(context.Background(), "runner:main")
	ctx = pcontext.CallTo(pcontext.InternalParigot(ctx), "main")
	defer pcontext.Dump(ctx)

	config, err := runner.Parse(ctx, flag.Arg(0), flg)
	if err != nil {
		log.Fatalf("failed to parse configuration file %s: %v", flag.Arg(0), err)

	}
	// create the syscall implementation
	exitCh := make(chan *syscall.ExitPair)
	wheeler.InstallWheeler(ctx, exitCh)
	ok := false
	kernel.K, ok = kernel.InitSingle()
	if !ok {
		log.Fatalf("unable to create kernel, aborting")
	}

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
		}
		log.Printf("goroutine exited")
	}(ctx)

	for _, mainProg := range main {
		code, err := deployCtx.StartMain(ctx, mainProg)
		if code == 253 && err == nil {
			//pcontext.Fatalf(ctx, "code failed (usually a panic) in execution of  program %s (code %d) -- can be host or guest", mainProg, code)
		} else if code != 0 {
			pcontext.Infof(ctx, "main exited from %s with code %d and error? %v", mainProg, code, err != nil)
		} else {
			pcontext.Infof(ctx, "program %s finished (code %d, error is not nil %v)", mainProg, code, err == nil)
		}
		return
	}
	pcontext.Dump(ctx)
	if len(main) > 1 {
		pcontext.Logf(ctx, pcontext.Info,
			"all main programs completed successfully")
	} else {
		pcontext.Logf(ctx, pcontext.Info, "main program completed successfully")
	}
	log.Printf("xxx main program exiting")
	for {

	}
	//os.Exit(8)
}
