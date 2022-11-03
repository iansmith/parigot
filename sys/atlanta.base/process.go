package sys

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/iansmith/parigot/lib"
)

var lastProcessId = 7

// callInfo is the data that is actually passed through the channel to the waiting server
// side.
type callInfo struct {
	mid    lib.Id
	cid    lib.Id
	param  []byte   // can be nil
	pctx   []byte   // can be nil for optimization reasons
	caller *Process // who sent this message, so we can send result
}

// resultInfo is the response that the recipient of a call sends back to the originator.
type resultInfo struct {
	cid    lib.Id
	result []byte //can be nil
	pctx   []byte // can be nil for optimization reasons
}

type Process struct {
	id       int
	path     string
	module   *wasmtime.Module
	linkage  []wasmtime.AsExtern
	memPtr   uintptr
	instance *wasmtime.Instance
	parent   *wasmtime.Store
	syscall  *SysCall

	callCh   chan *callInfo
	resultCh chan *resultInfo
}

// NewProcessFromMod does not handle concurrent use. It assumes that each call to this
// method is called from the same thread/goroutine, in sequence.  This is, effectively,
// a loader for the os.  xxxfixme this really should be safe to use in multiple go routines ... then we
// could have a repl
func NewProcessFromMod(parentStore *wasmtime.Store, mod *wasmtime.Module, path string, nameServer *NameServer) (*Process, error) {

	rt := newRuntime(nameServer)
	lastProcessId++
	id := lastProcessId
	proc := &Process{
		id:       id,
		path:     path,
		parent:   parentStore,
		module:   mod,
		linkage:  nil,
		memPtr:   0,
		instance: nil,
		syscall:  rt.syscall,

		callCh:   make(chan *callInfo),
		resultCh: make(chan *resultInfo),
	}

	l, err := proc.checkLinkage(rt)
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
	rt.SetProcess(proc)

	return proc, nil
}

func (p *Process) String() string {
	dir, file := filepath.Split(p.path)
	if dir == "" {
		dir = "."
	}
	return fmt.Sprintf("[proc-%d:%s]", p.id, file)
}

func (p *Process) checkLinkage(rt *Runtime) ([]wasmtime.AsExtern, error) {

	// all available funcs end up in here
	available := make(map[string]*wasmtime.Func)
	addSupportedFunctions(p.parent, available, rt)

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
				log.Printf("info: linked %s into module %s", importName, p.path)
			}
			linkage = append(linkage, ext)
		}
	}
	return linkage, nil
}

func (p *Process) Start() {
	start := p.instance.GetExport(p.parent, "run")
	if start == nil {
		log.Printf("unable to start process based on %s, can't fid start symbol", p.path)
		return
	}
	f := start.Func()
	result, err := f.Call(p.parent, 0, 0)
	if err != nil {
		log.Printf("process %d [%s] trapped: %v", p.id, p.path, err)
		// xxx fixme, we need to do process cleanup here
		return
	}
	if result == nil {
		log.Printf("process %d [%s] finished", p.id, p.path)
	} else {
		log.Printf("process %d [%s] fineshed: %+v", p.id, p.path, result)
	}
	// xxx fixme, we need to do process cleanup here
	return
}
