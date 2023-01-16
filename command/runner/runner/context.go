package runner

import (
	"fmt"
	"log"
	"os"

	"github.com/iansmith/parigot/sys"

	wasmtime "github.com/bytecodealliance/wasmtime-go/v3"
)

// A context represents a running state for an application--which is itself a collection of WASM
// modules.  A context holds the processes that are used by other parts of the system.
type Context struct {
	config *DeployConfig
	engine *wasmtime.Engine
	// notify chan is used at startup in the single process case
	NotifyCh chan *sys.KeyNSPair
	process  map[string]*sys.Process
}

// Flip this flag for more detailed output from the runner.
var runnerVerbose = false || os.Getenv("PARIGOT_VERBOSE") != ""

// NewContext returns a new, initialized runner.Context object or an error.  The initialized Context
// can be used to create processes and start the processes running.
func NewContext(conf *DeployConfig) (*Context, error) {
	// this config is for setting options that are global to the whole WASM world, like SetWasmThreads (ugh!)
	wasmConfig := wasmtime.NewConfig()

	log.Printf("NewContext creating wasm engine")
	engine := wasmtime.NewEngineWithConfig(wasmConfig)
	// load the images from disk and make sure they are valid modules
	if err := conf.LoadAllModules(engine); err != nil {
		return nil, err
	}

	// right now all the proc must be either remote or local, not a mix
	notify := make(chan *sys.KeyNSPair)
	sys.InitNameServer(notify, !conf.Flag.Remote, conf.Flag.Remote)

	//var rs *sys.RemoteSpec

	return &Context{
		config:   conf,
		engine:   engine,
		NotifyCh: notify,
		process:  make(map[string]*sys.Process),
		//	remoteSpec: rs,
	}, nil
}

// CreateProcess returns an error if it could not create a process (and an underlying store) for each
// module that was configured.  CreateProcess does not _start_ the processes running, see Start()
// for that.
func (c *Context) CreateProcess() error {
	// create processes and check linkage for each user program
	for _, name := range c.config.AllName() {
		log.Printf("create process %s", name)
		store := wasmtime.NewStore(c.engine)
		mod := c.config.module[name]
		if mod == nil {
			panic("unable to find (internal) module " + name)
		}
		p, err := sys.NewProcessFromMod(store, mod, c.config.Microservice[name])
		if err != nil {
			return fmt.Errorf("unable to create process from module (%s): %v", name, err)
		}
		c.process[name] = p
	}
	return nil
}

// Start takes all the processes that were created with CreateProcess() and starts them on goroutines.  If any server process exits, a warning
// is printed.  If a main program exits, we return and we can terminate the whole deal.
func (c *Context) Start() int {
	mainList := []string{}
	for _, f := range c.config.Microservice {
		if (c.config.Flag.TestMode && f.Test) || (!c.config.Flag.TestMode && f.Main) {
			mod := c.process[f.name]
			if mod == nil {
				panic("unable to find (internal) process for name " + f.name)
			}
			mainList = append(mainList, f.name)
		}
		if f.Server {
			log.Printf("StartProcess creating goroutine for server process %s at Start()", f.name)
			go func(p *sys.Process, name string) {
				code := p.Start()
				log.Printf("warning: server process %s exited with code %d", name, code)
			}(c.process[f.name], f.name)
		}
	}
	if len(mainList) == 0 {
		return 255
	}

	for _, main := range mainList {
		proc := c.process[main]
		runnerPrint("CreateProcess ", "starting goroutine for main process %v at Start() %s", main, proc)
		code := proc.Start()
		if code != 0 {
			return code
		}
	}
	return 0
}

func runnerPrint(method, spec string, arg ...interface{}) {
	if runnerVerbose {
		part1 := fmt.Sprintf("RUNNER:%s ", method)
		part2 := fmt.Sprintf(spec, arg...)
		log.Printf(part1, part2)
	}
}
