package sys

import (
	"fmt"
	"log"

	"github.com/bytecodealliance/wasmtime-go"
	"google.golang.org/protobuf/proto"
)

type callInfo struct {
	param proto.Message // can be nil
	pctx  proto.Message // can be nil for optimization reasons
}
type resultInfo struct {
	result proto.Message //can be nil
	pctx   proto.Message // can be nil for optimization reasons
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

func NewProcessFromMod(parentStore *wasmtime.Store, mod *wasmtime.Module, path string) (*Process, error) {

	rt := newRuntime()
	proc := &Process{
		path:     path,
		parent:   parentStore,
		module:   mod,
		linkage:  nil,
		memPtr:   0,
		instance: nil,
		syscall:  rt.syscall,
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
	log.Printf("xxx module %s has memptr %x", path, memptr)

	return proc, nil

}

func (p *Process) checkLinkage(rt *runtime) ([]wasmtime.AsExtern, error) {

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
			log.Printf("info: linked %s into module %s", importName, p.path)
			linkage = append(linkage, ext)
		}
	}
	return linkage, nil
}
