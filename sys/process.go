package sys

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"plugin"
	"sync"
	"time"

	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/eng"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	"github.com/iansmith/parigot/sys/dep"
)

// ParigotInitSymbolName is the name of the function that parigot
// will call in a user-supplied plugin to do initialization. The
// type of this must be ParigotInit.
const ParigotInitSymbolName = "ParigotInitialize"

// ParigotInit is the interface that plugins must meet to be
// initialized. It is expected that they will use the supplied
// Engine in the call to Init to register Host functions.
type ParigotInit interface {
	Init(ctx context.Context, e eng.Engine, i eng.Instance) bool
}

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
var processVerbose = true || envVerbose != ""

var lastProcessId = 7

type Process struct {
	id   int
	path string

	lock sync.Mutex

	module   eng.Module
	instance eng.Instance
	engine   eng.Engine

	microservice Service
	key          dep.DepKey

	requirementsMet bool
	reachedRunBlock bool
	running         bool
	exited          bool
	exitCode_       int //really only 0-192

	argv       int32 //ptr
	argc       int32
	argvBuffer *bytes.Buffer

	callCh chan *callContext
	runCh  chan bool
}

// NewProcessFromMicroservice does not handle concurrent use. It assumes that each call to this
// method is called from the same thread/goroutine, in sequence.  This is, effectively,
// a loader for the os.  xxxfixme this really should be safe to use in multiple go routines ... then we
// could have a repl??
func NewProcessFromMicroservice(c context.Context, engine eng.Engine, m Service, ctx *DeployContext) (*Process, error) {

	lastProcessId++
	id := lastProcessId
	proc := &Process{
		id:              id,
		engine:          engine,
		module:          m.GetModule(),
		instance:        nil,
		running:         false,
		reachedRunBlock: false,
		exited:          false,
		microservice:    m,
		path:            m.GetWasmPath(),

		callCh: make(chan *callContext),
	}
	proc.key = NewDepKeyFromProcess(proc)
	instance, err := proc.module.NewInstance(c)
	if err != nil {
		return nil, err
	}
	proc.instance = instance

	sym, err := m.GetPlugin().Lookup(ParigotInitSymbolName)

	initFn := sym.(ParigotInit)
	cont := &pcontext.LogContainer{}
	initCtx := SetupContextFor(cont, "ParigiotInit")
	ok := initFn.Init(initCtx, engine, instance)
	if !ok {
		pcontext.Dump(cont)
		return nil, fmt.Errorf("unable to initialize plugin '%s'", m.GetPluginPath())
	}
	pcontext.Dump(cont)

	return proc, nil
}

func (p *Process) RequirementsMet() bool {
	return p.requirementsMet
}
func (p *Process) IsServer() bool {
	// if we have a remote spec, then we are remote
	return p.microservice.IsServer()
}

func (p *Process) Exit() {
	p.lock.Lock()
	defer p.lock.Unlock()
	panic("xxx exit() xxx")
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
func (p *Process) SetReachedRunBlock(r bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.reachedRunBlock = r
}

func (p *Process) ReachedRunBlock() bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.reachedRunBlock
}
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
func (p *Process) Start(ctx context.Context) (code int) {
	procPrint(ctx, "START ", "start process: %s", p)
	var err error
	procPrint(ctx, "START ", "start of args  %+v", p.microservice.GetArg())
	// startOfArgs := wasmStartAddr + int32(0)
	// if p == nil {
	// 	panic("process is nil!")
	// }
	// p.lock.Lock()
	// defer p.lock.Unlock()

	// procPrint("START ", "get buffer from args and env  %s", p)
	// p.argvBuffer, p.argv, err = GetBufferFromArgsAndEnv(p.microservice, startOfArgs)
	// if err != nil {
	// 	code = int(ExitCodeArgsTooLarge)
	// 	return
	// }
	// p.argc = int32(len(p.microservice.GetArg()))

	procPrint(ctx, "START", "get entry point")
	start, err := p.instance.GetEntryPointExport(ctx)
	if err != nil {
		panic(err)
	}
	procPrint(ctx, "START", "defer %s (%v)", p, start != nil)

	defer func(proc *Process) {
		r := recover()
		if r != nil {
			log.Printf("defer caught it %T, %v", r, r.(string))
			log.Printf("flush")

			procPrint(ctx, "START ******** ", "INSIDE defer %s, %+v", proc, r)
			e, ok := r.(*syscallmsg.ExitRequest)
			procPrint(ctx, "Start/Exit ", "INSIDE defer exit req %+v, ok %v", r.(*syscallmsg.ExitRequest), ok)
			if ok {
				code = int(e.GetCode())
				proc.SetExitCode(code)
				procPrint(ctx, "Start/Exit", "INSIDE DEFER exiting with code %d", e.GetCode())
			} else {
				p.SetExitCode(int(ExitCodePanic))
				code = int(ExitCodePanic)
				procPrint(ctx, "Start/Exit", "golang (not WASM) panic '%v'\n", r)
			}
		}
	}(p)
	procPrint(ctx, "START ", "calling start func %s", p)
	var info interface{}

	retVal, err := start.Run(ctx, p.microservice.GetArg(), info)
	procPrint(ctx, "END ", "process %s has completed: result=%v, err=%v", p, retVal, err)

	if err != nil {
		p.SetExitCode(int(ExitCodeTrapped))
		procPrint(ctx, "END ", "process %s trapped: %v, exit code %d", p, err, p.ExitCode())
		return int(ExitCodeTrapped)
	}
	if retVal == nil {
		procPrint(ctx, "END ", "process %s finished w/no return value (exit code %d)", p, p.ExitCode())
		p.SetExited(true)
		return p.ExitCode()
	}
	procPrint(ctx, "END ", "process %s finished normally: %+v", p, retVal)
	procPrint(ctx, "END ", "going to sleep now")
	ch := make(chan struct{})
	<-ch
	return p.ExitCode()
}

func procPrint(ctx context.Context, method string, spec string, arg ...interface{}) {
	if processVerbose {
		pcontext.LogFullf(ctx, pcontext.Debug, pcontext.Parigot, method, spec, arg...)
	}
}

func SetupContextFor(cont *pcontext.LogContainer, funcName string) context.Context {

	ctx := context.WithValue(context.Background(), pcontext.ParigotTime, time.Now())
	ctx = context.WithValue(ctx, pcontext.ParigotFunc, funcName)
	ctx = context.WithValue(ctx, pcontext.ParigotSource, pcontext.ServerGo)
	ctx = context.WithValue(ctx, pcontext.ParigotLogContainer, cont)
	return ctx
}
