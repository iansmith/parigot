package sys

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/command/runner/runner"
	"github.com/iansmith/parigot/eng"

	"github.com/tetratelabs/wazero"
)

// A DeployContext represents a deployment during the process of starting it up.
// A context holds the processes that are used by other parts of the system.
type DeployContext struct {
	config  *runner.DeployConfig
	engine  eng.Engine
	notify  *sync.Map
	process *sync.Map
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

	engine := eng.NewWaZeroEngine(context.Background(), wazero.NewRuntimeConfig())

	// our notify map is shared by the nameserver
	notifyMap := &sync.Map{}
	processMap := &sync.Map{}

	depCtx := &DeployContext{
		config:  conf,
		engine:  engine,
		process: processMap,
		notify:  notifyMap,
	}

	return depCtx, nil
}

func (c *DeployContext) Timezone() string {
	return c.config.Timezone
}

func (c *DeployContext) LoadAllModules(e eng.Engine) error {
	return c.config.LoadAllModules(e)
}

func (c *DeployContext) Process() *sync.Map {
	return c.process
}

// CreateAllProcess returns an error if it could not create a process (and an underlying store) for each
// module that was configured.  CreateAllProcess does not start the processes running, see Start()
// for that.
func (c *DeployContext) CreateAllProcess() error {
	// load the parigot syscalls, this is done based on the config in the .toml file
	err := LoadPluginAndAddHostFunc(context.Background(),
		c.config.ParigotLibPath, c.config.ParigotLibSymbol, c.engine, "parigot")
	if err != nil {
		return err
	}

	// load wasm files, implicitly checks them and converts them to binary
	if err := c.LoadAllModules(c.engine); err != nil {
		panic(fmt.Sprintf("unable to load modules in preparation for launch: %v", err))
	}

	// create processes
	for _, name := range c.config.AllName() {
		m := c.config.Microservice[name]
		hid := id.NewHostId()
		p, err := NewProcessFromMicroservice(c.engine, m, c, hid)
		if err != nil {
			return fmt.Errorf("unable to create process from module '%s': %v", name, err)
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
func (c *DeployContext) StartServer(ctx context.Context) ([]string, int) {
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
			go func(p *Process, serverProcessName string) {
				code := p.Start()
				p.SetExitCode(code)
				if code != 0 {
					log.Printf("service exited ", "name", serverProcessName, "code", code)
				}

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

func (d *DeployContext) instantiateBuiltinHostFunc(ctx context.Context) error {
	for _, name := range []string{"parigot"} {
		if _, err := d.engine.InstantiateHostModule(ctx, name); err != nil {
			return err
		}
	}
	return nil
}

func contextPrint(spec string, arg ...interface{}) {
	if runnerVerbose {
		log.Printf(spec, arg...)
	}
}
