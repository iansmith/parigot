package sys

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	wasmtime "github.com/bytecodealliance/wasmtime-go/v3"
	fileimpl "github.com/iansmith/parigot/api/fileimpl/go_"
	logimpl "github.com/iansmith/parigot/api/logimpl/go_"
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/dep"
	"google.golang.org/protobuf/types/known/anypb"
)

// Flip this switch to see debug messages from the process.
var processVerbose = true

var lastProcessId = 7

// callInfo is the data that is actually passed through the channel to the waiting server
// side.  This contains the channel to send the response on.
type callInfo struct {
	mid    lib.Id     // method id
	cid    lib.Id     // call id
	sid    lib.Id     // service id
	param  []byte     // can be nil
	pctx   []byte     // can be nil for optimization reasons
	sender dep.DepKey // who sent this message, so we can return result
	method string     // name of the method being called
	respCh chan *anypb.Any
}

// resultInfo is the response that the recipient of a call sends back to the originator.
type resultInfo struct {
	cid     lib.Id
	mid     lib.Id
	errorId lib.Id
	result  []byte //can be nil
	pctx    []byte // can be nil for optimization reasons

}

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
	local      *bool

	callCh   chan *callContext
	resultCh chan *resultInfo
	runCh    chan bool
}

// NewProcessFromMod does not handle concurrent use. It assumes that each call to this
// method is called from the same thread/goroutine, in sequence.  This is, effectively,
// a loader for the os.  xxxfixme this really should be safe to use in multiple go routines ... then we
// could have a repl
func NewProcessFromMod(parentStore *wasmtime.Store, mod *wasmtime.Module, path string,
	rs *RemoteSpec) (*Process, error) {

	rt := newRuntime(rs)
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
		server:     false,

		callCh:   make(chan *callContext),
		resultCh: make(chan *resultInfo),
		runCh:    make(chan bool),
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

	if rt.spec.IsRemote(proc) {
		proc.server = true
	}

	return proc, nil
}

func (p *Process) IsServer() bool {
	// if we have a remote spec, then we are remote
	return p.server
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
			if strings.HasPrefix(importName, "go.parigot") {
				procPrint("CHECKLINKAGE ", "linked %s into module %s", importName, p.path)
			}
			linkage = append(linkage, ext)
		}
	}
	return linkage, nil
}

// Run() is used to let a process go past its call to "start()" in the API.  This is
// called when we discover all his requirements have been met.
func (p *Process) Run() {
	procPrint("RUN", "trying to tell %s to run, everything is ok", p)
	p.runCh <- true
	procPrint("RUN", "process %s running", p)
}

func (p *Process) Start() {
	procPrint("START ", "we have been loaded/started by the runner: %s", p)
	start := p.instance.GetExport(p.parent, "run")
	if start == nil {
		log.Printf("unable to start process based on %s, can't fid start symbol", p.path)
		return
	}
	f := start.Func()
	procPrint("START ", "calling the entry point (%+v,%T), for proc %s (parent %v)",
		f, f, p, p.parent)
	result, err := f.Call(p.parent, 0, 0)
	p.exited = true
	if err != nil {
		procPrint("START ", "process %s trapped: %v", p, err)
	} else {
		if result == nil {
			procPrint("START ", "process %s finished", p)
		} else {
			procPrint("START ", "process %s finished: %+v", p, result)
		}
	}
	procPrint("START ", "exiting...")
	// xxx fixme, we need to do process cleanup here
	return
}

func procPrint(method string, spec string, arg ...interface{}) {
	if processVerbose {
		part1 := fmt.Sprintf("PROCESS:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}
