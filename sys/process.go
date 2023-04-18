package sys

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sync"

	fileimpl "github.com/iansmith/parigot/apiimpl/file/go_"
	logimpl "github.com/iansmith/parigot/apiimpl/log/go_"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	"github.com/iansmith/parigot/sys/dep"
	"github.com/iansmith/parigot/sys/jspatch"

	wasmtime "github.com/bytecodealliance/wasmtime-go/v7"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service interface {
	IsServer() bool
	IsLocal() bool
	IsRemote() bool
	GetName() string
	GetArg() []string
	GetEnv() []string
	GetPath() string
	GetModule() *wasmtime.Module
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

	module       *wasmtime.Module
	linkage      []wasmtime.AsExtern
	memPtr       uintptr
	instance     *wasmtime.Instance
	parent       *wasmtime.Store
	syscall      *syscallReadWrite
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
func NewProcessFromMicroservice(parentStore *wasmtime.Store, m Service, ctx *DeployContext) (*Process, error) {

	rt := newRuntime(ctx)
	// split mode
	logViewer := &logimpl.LogViewerImpl{}
	fileSvc := &fileimpl.FileSvcImpl{}

	lastProcessId++
	id := lastProcessId
	proc := &Process{
		id:              id,
		parent:          parentStore,
		module:          m.GetModule(),
		linkage:         nil,
		memPtr:          0,
		instance:        nil,
		syscall:         rt.syscall,
		running:         false,
		reachedRunBlock: false,
		exited:          false,
		microservice:    m,
		path:            m.GetPath(),

		callCh: make(chan *callContext),
		//runCh:  make(chan bool),
	}
	proc.key = NewDepKeyFromProcess(proc)

	l, err := proc.checkLinkage(rt, logViewer, fileSvc)
	if err != nil {
		return nil, err
	}
	proc.linkage = l

	instance, err :=
		wasmtime.NewInstance(parentStore, proc.module, l)
	if err != nil {
		return nil, err
	}
	proc.instance = instance

	ext := instance.GetExport(parentStore, "mem")
	if ext.Memory() == nil {
		return nil, fmt.Errorf("'mem' export is not a memory object")
	}

	memptr := uintptr(ext.Memory().Data(parentStore))
	proc.memPtr = memptr
	rt.SetMemPtr(memptr)
	// split mode
	logViewer.SetWasmMem(memptr)
	fileSvc.SetWasmMem(memptr)
	rt.SetProcess(proc)

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
	print(fmt.Sprintf("process %s exiting\n", p))
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

func (p *Process) checkLinkage(rt *Runtime, lv *logimpl.LogViewerImpl, fs *fileimpl.FileSvcImpl) ([]wasmtime.AsExtern, error) {

	// all available funcs end up in here
	available := make(map[string]*wasmtime.Func)
	addSupportedFunctions(p.parent, available, rt)
	addSplitModeFunctions(p.parent, available, lv, fs)

	// result of checking the linkage
	linkage := []wasmtime.AsExtern{}

	// walk all the module's imports
	for _, imp := range p.module.Imports() {
		n := "$$ANON$$"
		if imp.Name() != nil {
			n = *imp.Name()
		}
		importName := fmt.Sprintf("%s.%s", imp.Module(), n)
		ext, ok := available[importName]
		if !ok {
			return nil, fmt.Errorf("unable to find linkage for %s in module %s", importName, p.path)
		} else {
			linkage = append(linkage, ext)
		}
	}
	return linkage, nil
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
	procPrint("START ", "we have been loaded/started by the runner: %s", p)
	var err error
	procPrint("START", "start of args  %s", p)
	startOfArgs := wasmStartAddr + int32(0)
	if p == nil {
		panic("process is nil!")
	}
	// p.lock.Lock()
	// defer p.lock.Unlock()

	procPrint("START", "get buffer from args and env  %s", p)
	p.argvBuffer, p.argv, err = GetBufferFromArgsAndEnv(p.microservice, startOfArgs)
	if err != nil {
		code = int(ExitCodeArgsTooLarge)
		return
	}
	p.argc = int32(len(p.microservice.GetArg()))

	procPrint("START", "create wasm mem %s", p)

	wasmMem := jspatch.NewWasmMem(p.memPtr)
	wasmMem.SetInt32(wasmStartAddr-int32(4), p.argv)
	wasmMem.CopyToMemAddr(startOfArgs, p.argvBuffer.Bytes())

	procPrint("START", "get export %s", p)
	start := p.instance.GetExport(p.parent, "run")
	if start == nil {
		p.SetExitCode(int(ExitCodeNoStartSymbol))
		return p.ExitCode()
	}
	procPrint("START", "defer %s", p)

	defer func(proc *Process) {
		r := recover()
		if r != nil {
			procPrint("START ******** ", "INSIDE defer %s, %+v", proc, r)
			e, ok := r.(*syscallmsg.ExitRequest)
			procPrint("START ", "INSIDE defer exit req %+v, ok %v", r.(*syscallmsg.ExitRequest), ok)
			if ok {
				code = int(e.GetCode())
				proc.SetExitCode(code)
				procPrint("Start/Exit", "exiting with code %d", e.GetCode())
			} else {
				p.SetExitCode(int(ExitCodePanic))
				code = int(ExitCodePanic)
				print(fmt.Sprintf("golang (not WASM) panic '%v'\n", r))
			}
		}
	}(p)
	procPrint("START", "getting start func %s", p)
	f := start.Func()
	procPrint("START", "calling start func %s", p)
	result, err := f.Call(p.parent, p.argc, p.argv)
	procPrint("END ", "process %s has completed: %v, %v", p, result, err)

	if err != nil {
		p.SetExitCode(int(ExitCodeTrapped))
		procPrint("END ", "process %s trapped: %v, exit code %d", p, err, p.ExitCode())
		return int(ExitCodeTrapped)
	}
	if result == nil {
		procPrint("END ", "process %s finished (exit code %d)", p, p.ExitCode())
		p.SetExited(true)
		return p.ExitCode()

	}
	procPrint("END ", "process %s finished normally: %+v", p, result)
	procPrint("END ", "going to sleep now")
	ch := make(chan struct{})
	<-ch
	return p.ExitCode()
}

func procPrint(method string, spec string, arg ...interface{}) {
	if processVerbose {
		part1 := fmt.Sprintf("PROCESS:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		logimpl.ProcessLogRequest(
			&logmsg.LogRequest{
				Level:   logmsg.LogLevel_LOG_LEVEL_INFO,
				Message: part1 + part2 + "\n",
				Stamp:   timestamppb.Now(), // xxx should use the kernel calls
			}, true, false, false, nil)

	}
}
