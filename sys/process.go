package sys

import (
	"fmt"
	"log"
	"path/filepath"

	fileimpl "github.com/iansmith/parigot/api_impl/file/go_"
	logimpl "github.com/iansmith/parigot/api_impl/log/go_"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"

	wasmtime "github.com/bytecodealliance/wasmtime-go/v3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Flip this switch to see debug messages from the process.
var processVerbose = true || envVerbose != ""

var lastProcessId = 7

type Process struct {
	id   int
	path string

	module   *wasmtime.Module
	linkage  []wasmtime.AsExtern
	memPtr   uintptr
	instance *wasmtime.Instance
	parent   *wasmtime.Store
	syscall  *syscallReadWrite

	waiter     bool
	reachedRun bool
	exited     bool
	server     bool
	remote     bool
	local      *bool

	callCh chan *callContext
	runCh  chan bool

	exitCode int
}

// NewProcessFromMod does not handle concurrent use. It assumes that each call to this
// method is called from the same thread/goroutine, in sequence.  This is, effectively,
// a loader for the os.  xxxfixme this really should be safe to use in multiple go routines ... then we
// could have a repl
func NewProcessFromMod(parentStore *wasmtime.Store, mod *wasmtime.Module, path string,
	isServer, isMain, isRemote bool) (*Process, error) {

	if isServer && isMain {
		panic("unable to understand process that is both a server and a main:" + path)
	}

	rt := newRuntime()
	// split mode
	logViewer := &logimpl.LogViewerImpl{}
	fileSvc := &fileimpl.FileSvcImpl{}

	lastProcessId++
	id := lastProcessId
	proc := &Process{
		id:         id,
		path:       path,
		parent:     parentStore,
		module:     mod,
		linkage:    nil,
		memPtr:     0,
		instance:   nil,
		syscall:    rt.syscall,
		waiter:     false,
		reachedRun: false,
		exited:     false,
		server:     isServer,
		remote:     isRemote,

		callCh: make(chan *callContext),
		runCh:  make(chan bool),
	}

	l, err := proc.checkLinkage(rt, logViewer, fileSvc)
	if err != nil {
		return nil, err
	}
	proc.linkage = l

	instance, err := wasmtime.NewInstance(parentStore, mod, l)
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

func (p *Process) IsServer() bool {
	// if we have a remote spec, then we are remote
	return p.server
}

func (p *Process) Exit() {
	print(fmt.Sprintf("process %s exiting\n", p))
	p.exited = true
}

func (p *Process) String() string {
	dir, file := filepath.Split(p.path)
	if dir == "" {
		dir = "."
	}

	return fmt.Sprintf("[proc-%d:%s]", p.id, file)
}

func (p *Process) ReachedStart() bool {
	return p.reachedRun
}
func (p *Process) Exited() bool {
	return p.exited
}
func (p *Process) IsWaiter() bool {
	return p.waiter
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
			// if strings.HasPrefix(importName, "go.parigot") {
			// 	//procPrint("CHECKLINKAGE ", "linked %s into module %s", importName, p.path)
			// }
			linkage = append(linkage, ext)
		}
	}
	return linkage, nil
}

func (p *Process) SetExitCode(code int) {
	p.exitCode = code
	p.exited = true
}

// Run() is used to let a process proceed with running.  This is
// called when we discover all his requirements have been met.
func (p *Process) Run() {
	procPrint("RUN", "trying to tell %s to run, everything is ok", p)
	p.runCh <- true
	procPrint("RUN", "process %s running", p)
}

// Start invokes the wasm interp and returns an error code if this is a "main" process.
func (p *Process) Start() (code int) {
	procPrint("START ", "we have been loaded/started by the runner: %s", p)
	start := p.instance.GetExport(p.parent, "run")
	if start == nil {
		log.Printf("unable to start process based on %s, can't fid start symbol", p.path)
		p.SetExitCode(55)
		return p.exitCode
	}
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(*syscallmsg.ExitRequest)
			if ok {
				p.exitCode = int(e.GetCode())
				code = p.exitCode
				procPrint("Start/Exit", "exiting with code %d", e.GetCode())
			} else {
				p.SetExitCode(254)
				code = p.exitCode
				procPrint("Start/Exit", "trapped a panic: %v, exiting with code 254", r)
			}
		}
	}()
	f := start.Func()
	procPrint("START ", "calling the entry point for proc %s", p)
	result, err := f.Call(p.parent, 0, 0)
	procPrint("END ", "process %s has completed: %v, %v", p, result, err)

	if err != nil {
		p.SetExitCode(253)
		procPrint("END ", "process %s trapped: %v, exit code %d", p, err, p.exitCode)
		return p.exitCode
	}
	if result == nil {
		procPrint("END ", "process %s finished (exit code %d)", p, p.exitCode)
		p.exited = true
		return p.exitCode

	}
	procPrint("END ", "process %s finished normally: %+v", p, result)
	p.exitCode = 0
	return p.exitCode
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
			}, true, false, nil)
		//print(part1 + part2 + "\n")

	}
}
