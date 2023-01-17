package sys

import (
	"fmt"
	"log"
	"os"
	"sync"

	wasmtime "github.com/bytecodealliance/wasmtime-go/v3"
	"github.com/iansmith/parigot/command/runner/runner"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
)

// A DeployContext represents a deployment during the process of starting it up.
// A context holds the processes that are used by other parts of the system.
type DeployContext struct {
	config     *runner.DeployConfig
	engine     *wasmtime.Engine
	notify     *sync.Map
	process    *sync.Map
	nameserver NameServer
}

// Flip this flag for more detailed output from the runner.
var runnerVerbose = false || os.Getenv("PARIGOT_VERBOSE") != ""

// NewDeployContext returns a new, initialized DeployContext object or an error.
// This function can be thought of as the bridge between the configuration
// of the deploy (runner.DeployConfig) and the running state of the deployment
// which is represented by DeployContext.  This context can be used to create
// processes and start them running.
func NewDeployContext(conf *runner.DeployConfig) (*DeployContext, error) {
	// this config is for setting options that are global to the whole WASM world, like SetWasmThreads (ugh!)
	wasmConfig := wasmtime.NewConfig()

	engine := wasmtime.NewEngineWithConfig(wasmConfig)
	// load the images from disk and make sure they are valid modules
	if err := conf.LoadAllModules(engine); err != nil {
		return nil, err
	}
	// our notify map is shared by the nameserver
	notifyMap := &sync.Map{}
	processMap := &sync.Map{}
	var ns NameServer
	if !conf.Flag.Remote {
		raw := NewLocalNameServer(notifyMap)
		raw.NSCore = NewNSCore(true)
		ns = raw
		// xxx FIX ME! evil hack
		LocalNS = raw
	} else {
		panic("not implemented yet, remote nameserever start up")
	}

	return &DeployContext{
		config:     conf,
		engine:     engine,
		process:    processMap,
		notify:     notifyMap,
		nameserver: ns,
		//	remoteSpec: rs,
	}, nil
}

func (c *DeployContext) Process() *sync.Map {
	return c.process
}

// CreateAllProcess returns an error if it could not create a process (and an underlying store) for each
// module that was configured.  CreateAllProcess does not start the processes running, see Start()
// for that.
func (c *DeployContext) CreateAllProcess() error {
	// create processes and check linkage for each user program
	for _, name := range c.config.AllName() {
		m := c.config.Microservice[name]
		procPrint("CreateAlllProcess", "create process %s", name)
		store := wasmtime.NewStore(c.engine)
		mod := c.config.Module(name)
		if mod == nil {
			panic("unable to find (internal) module for " + name)
		}
		p, err := NewProcessFromMicroservice(store, m, c)
		if err != nil {
			return fmt.Errorf("unable to create process from module (%s): %v", name, err)
		}
		c.process.Store(name, p)
		ch := make(chan bool)
		c.notify.Store(name, ch)
	}
	return nil
}

// StartServer takes all the processes that were created with CreateAllProcess
// and are marked as servers and starts them.  It also returns a list of
// names that are the names of the "main" program(s).  Main process can be
// proper main programs (run to completion programs) or tests, depending on
// the value of the flags in the configuration.  If there was an error,
// this function returns a nil for the list of main programs and it
// returns the exit code to be used when exiting.
func (c *DeployContext) StartServer() ([]string, int) {
	mainList := []string{}
	for _, f := range c.config.Microservice {
		procAny, ok := c.process.Load(f.Name())
		if !ok {
			panic("unable to find (internal) process for name " + f.Name())
		}
		if (c.config.Flag.TestMode && f.Test) || (!c.config.Flag.TestMode && f.Main) {
			mainList = append(mainList, f.Name())
		}
		name := f.Name()
		if f.Server {
			contextPrint(logmsg.LogLevel_LOG_LEVEL_DEBUG, "StartingServer", "StartProcess creating goroutine for server process %s at Start()", name)
			go func(p *Process, serverProcessName string) {
				code := p.Start()
				p.SetExitCode(code)
				contextPrint(logmsg.LogLevel_LOG_LEVEL_ERROR, "StartingServer", "server process %s exited with code %d", serverProcessName, code)
			}(procAny.(*Process), name)
		}
	}
	if len(mainList) == 0 {
		return nil, int(ExitCodeNoMain)
	}
	return mainList, 0
}

// StartMain runs a main program (one that is not a server and usually expected
// to terminate) and returns the error code provided by the main program.  Note
// that this function is run synchronously, not on a goroutine.
func (c *DeployContext) StartMain(mainProg string) (int, error) {
	procAny, ok := c.process.Load(mainProg)
	if !ok {
		return 0, fmt.Errorf("main program '%s' not found", mainProg)
	}
	proc := procAny.(*Process)
	code := proc.Start()
	proc.SetExitCode(code)
	return code, nil
}

func (d *DeployContext) NotifyMap() *sync.Map {
	return d.notify
}

func contextPrint(level logmsg.LogLevel, method, spec string, arg ...interface{}) {
	if runnerVerbose {
		part1 := fmt.Sprintf("RUNNER:%s ", method)
		part2 := fmt.Sprintf(spec, arg...)
		log.Printf("%s%s", part1, part2)
	}
}
