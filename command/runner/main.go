package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/iansmith/parigot/command/runner/runner"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	"github.com/iansmith/parigot/sys"
	"github.com/iansmith/parigot/sys/backdoor"

	"google.golang.org/protobuf/types/known/timestamppb"
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

	// the deploy context creation also creates any needed nameservers
	ctx, err := sys.NewDeployContext(config)
	if err != nil {
		log.Fatalf("unable to create deploy context: %v", err)
	}
	if err := ctx.CreateAllProcess(); err != nil {
		log.Fatalf("unable to create process: %v", err)
	}
	main, code := ctx.StartServer()
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
			ctx.Process().Range(func(keyAny, valueAny any) bool {
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
		code, err := ctx.StartMain(mainProg)
		if err != nil {
			log.Fatalf("could not start main program:%v", err)
		}
		if code != 0 {
			log.Fatalf("main program '%s' exited with code %d", mainProg, code)
		}
	}
	log.Printf("size of main is %+v", main)
	if len(main) > 1 {
		backdoor.Log(&logmsg.LogRequest{
			Level:   logmsg.LogLevel_LOG_LEVEL_INFO,
			Stamp:   timestamppb.Now(),
			Message: "all main programs completed successfully",
		}, false, true, false, nil)
	} else {
		backdoor.Log(&logmsg.LogRequest{
			Level:   logmsg.LogLevel_LOG_LEVEL_INFO,
			Stamp:   timestamppb.Now(),
			Message: "main program completed successfully",
		}, false, true, false, nil)
		log.Printf("xxx backdoor log, main completed\n")
	}
	os.Exit(8)
}
