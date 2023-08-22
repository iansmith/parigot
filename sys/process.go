package sys

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/eng"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
)

type Service interface {
	IsServer() bool
	IsLocal() bool
	IsRemote() bool
	GetName() string
	GetArg() []string
	GetEnv() []string
	GetWasmPath() string
	GetModule() eng.Module
	GetPluginPath() string
	GetPlugin() *plugin.Plugin
	GetPluginSymbol() string
}

type ParigotExitCode int

const (
	ExitCodeArgsTooLarge  ParigotExitCode = 252
	ExitCodeNoStartSymbol ParigotExitCode = 251
	ExitCodePanic         ParigotExitCode = 254
	ExitCodeTrapped       ParigotExitCode = 253
	ExitCodeNoMain        ParigotExitCode = 255
)

// Flip this switch to see debug messages from the process.
var processVerbose = false

var lastProcessId = 7

type Process struct {
	id   int
	path string

	lock sync.Mutex

	module   eng.Module
	instance eng.Instance
	engine   eng.Engine

	microservice Service

	hid id.HostId // our name

	// requirementsMet bool
	// reachedRunBlock bool
	running   bool
	exited    bool
	exitCode_ int //really only 0-192

	// argv       int32 //ptr
	// argc       int32
	// argvBuffer *bytes.Buffer

	exitChan chan int32

	runCh chan bool
}

// NewProcessFromMicroservice does not handle concurrent use. It assumes that each call to this
// method is called from the same thread/goroutine, in sequence.  This is, effectively,
// a loader for the os.  xxxfixme this really should be safe to use in multiple go routines ... then we
// could have a repl??
func NewProcessFromMicroservice(engine eng.Engine, m Service, ctx *DeployContext, hid id.HostId) (*Process, error) {

	lastProcessId++
	id := lastProcessId
	proc := &Process{
		id:       id,
		engine:   engine,
		module:   m.GetModule(),
		instance: nil,
		running:  false,
		//reachedRunBlock: false,
		exited:       false,
		microservice: m,
		path:         m.GetWasmPath(),
		exitChan:     make(chan int32),
		hid:          hid,
	}

	if m.GetPluginPath() != "" {
		err := LoadPluginAndAddHostFunc(context.Background(),
			m.GetPluginPath(), m.GetPluginSymbol(), engine, m.GetName())
		if err != nil {
			return nil, err
		}
	}

	instance, err := proc.module.NewInstance(context.Background(), ctx.config.Timezone)
	if err != nil {
		return nil, err
	}

	proc.instance = instance
	return proc, nil
}

func LoadPluginAndAddHostFunc(ctx context.Context, pluginPath string, pluginSymbol string, engine eng.Engine, name string) error {
	i, err := LoadPlugin(ctx, pluginPath, pluginSymbol, name)
	if err != nil {
		return err
	}
	if !i.Init(ctx, engine) {
		return fmt.Errorf("unable to load plugin: %v", err.Error())
	}

	if _, err := engine.InstantiateHostModule(ctx, name); err != nil {
		return fmt.Errorf("instantiate host module failed: %s", err.Error())
	}
	return nil
}

// func (p *Process) RequirementsMet() bool {
// 	return p.requirementsMet
// }

func (p *Process) IsServer() bool {
	// if we have a remote spec, then we are remote
	return p.microservice.IsServer()
}

func (p *Process) Exit() {
	p.lock.Lock()
	defer p.lock.Unlock()
	log.Printf("process %s exiting\n", p)
	p.exited = true
}

func (p *Process) String() string {
	p.lock.Lock()
	defer p.lock.Unlock()

	dir, file := filepath.Split(p.path)
	if dir == "" {
		dir = "."
	}

	return fmt.Sprintf("[proc-%d:%s:%s]", p.id, p.microservice.GetName(), file)
}

// func (p *Process) SetReachedRunBlock(r bool) {
// 	p.lock.Lock()
// 	defer p.lock.Unlock()

// 	p.reachedRunBlock = r
// }

// func (p *Process) ReachedRunBlock() bool {
// 	p.lock.Lock()
// 	defer p.lock.Unlock()

// 	return p.reachedRunBlock
// }

func (p *Process) Running() bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.running
}

func (p *Process) SetExited(e bool) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.exited = e
}
func (p *Process) Exited() bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.exited
}

func (p *Process) SetExitCode(code int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.exitCode_ = code
	p.exited = true
}

func (p *Process) ExitCode() int {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.exitCode_
}
func (p *Process) IsRunning() bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.running
}

func (p *Process) SetRunning(r bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.running = r
}

// Run() is used to let a process proceed with running.  This is
// called when we discover all his requirements have been met.
func (p *Process) Run() {
	p.runCh <- true
}

// Start invokes the wasm interp and returns an error code if this is a "main" process.
func (p *Process) Start() (code int) {
	if p == nil {
		panic("unable to Start when there is no process (p==nil)")
	}
	procPrint("start process: %s", p)
	var err error
	procPrint("start of args  %+v", p.microservice.GetArg())

	procPrint("get entry point")
	start, err := p.instance.EntryPoint(context.Background())
	if err != nil {
		panic(err)
	}
	// exitCode, err := p.instance.Value(context.Background(), "exit_code_")
	// if err != nil {
	// 	panic(err)
	// }
	procPrint("defer %s (%v)", p, start != nil)

	defer func(proc *Process) {
		r := recover()
		if r != nil {
			log.Printf("defer caught it %T, %v", r, r.(string))
			log.Printf("flush")

			procPrint("INSIDE defer %s, %+v", proc, r)
			e, ok := r.(*syscall.ExitRequest)
			procPrint("INSIDE defer exit req %+v, ok %v", r.(*syscall.ExitRequest), ok)
			if ok {
				code = int(e.Pair.GetCode())
				proc.SetExitCode(code)
				procPrint("INSIDE DEFER exiting with code %d", e.Pair.GetCode())
			} else {
				p.SetExitCode(int(ExitCodePanic))
				code = int(ExitCodePanic)
				procPrint("golang (not WASM) panic '%v'\n", r)
			}
		}
	}(p)
	procPrint("calling start func %s", p)
	var info interface{}

	retVal, err := start.Run(context.Background(), p.microservice.GetArg(), info)
	procPrint("process %s has completed: result=%v, err=%v", p, retVal, err)

	//procPrint("xxxx what is exit code?  %0xd", exitCode.GetU16())

	if err != nil {
		p.SetExitCode(int(ExitCodeTrapped))
		procPrint("process %s trapped: %v, exit code %d", p, err, p.ExitCode())
		return int(ExitCodeTrapped)
	}
	if retVal == nil {
		procPrint("process %s finished w/no return value (exit code %d)", p, p.ExitCode())
		p.SetExited(true)
		return p.ExitCode()
	}
	procPrint("process %s finished normally: %+v", p, retVal)
	ch := make(chan struct{})
	<-ch
	return p.ExitCode()
}

func procPrint(spec string, arg ...interface{}) {
	if processVerbose {
		log.Printf(spec, arg...)
	}
}

// func SetupContextFor(cont *pcontext.LogContainer, funcName string) context.Context {

// 	ctx := context.WithValue(context.Background(), pcontext.ParigotTime, time.Now())
// 	ctx = context.WithValue(ctx, pcontext.ParigotFunc, funcName)
// 	ctx = context.WithValue(ctx, pcontext.ParigotSource, pcontext.ServerGo)
// 	ctx = context.WithValue(ctx, pcontext.ParigotLogContainer, cont)
// 	return ctx
// }
