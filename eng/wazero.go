package eng

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

const maxWasmFile = 0x1024 * 0x1024 * 0x20

const EntryPoint = "parigot_main"

type wazeroEng struct {
	r       wazero.Runtime
	builder map[string]wazero.HostModuleBuilder
}
type wazeroModule struct {
	parent *wazeroEng
	cm     wazero.CompiledModule
}

type wazeroInstance struct {
	parent       *wazeroModule
	m            api.Module
	malloc, free api.Function
}

type wazeroExtern struct {
	parent *wazeroInstance
}

type wazeroMemoryExtern struct {
	*wazeroExtern
	m api.MemoryDefinition
}

type wazeroEntryPointExtern struct {
	*wazeroExtern
	f api.Function
}

var bg = context.Background()

func NewWaZeroEngine(conf wazero.RuntimeConfig) Engine {
	e := &wazeroEng{}
	if conf != nil {
		e.r = wazero.NewRuntimeWithConfig(bg, conf)
	} else {
		e.r = wazero.NewRuntime(bg)
	}
	wasi_snapshot_preview1.MustInstantiate(bg, e.r)
	e.builder = make(map[string]wazero.HostModuleBuilder)
	return e
}
func (e *wazeroEng) NewModuleFromFile(path string) (Module, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	limited := io.LimitReader(fp, maxWasmFile)
	all, err := io.ReadAll(limited)
	if err != nil {
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("unable to stat %s, trying to check for file larger than %d", path, len(all))
		}
		if info.Size() >= int64(len(all)) {
			return nil, fmt.Errorf("wasm file too large %s, file is %d bytes, limit is %d bytes", path, info.Size(), len(all))
		}
	}
	mod, err := e.r.CompileModule(bg, all)
	if err != nil {
		return nil, err
	}
	return &wazeroModule{cm: mod, parent: e}, nil
}
func (m *wazeroModule) Name() string {
	return m.cm.Name()
}

func (e *wazeroExtern) Name() string {
	return e.parent.Name()
}

func (i *wazeroInstance) GetMemoryExport() (MemoryExtern, error) {
	def := i.m.ExportedMemoryDefinitions()
	if def == nil {
		return nil, fmt.Errorf("module '%s' has no exported memory definitions", i.Name())
	}
	if len(def) > 1 {
		panic("parigot currently only supports one memory export, but you should file a ticket to remind us to fix that")
	}
	for k, v := range def {
		log.Printf("xxx -- memory defs '%s',%+v", k, v)
	}
	return &wazeroMemoryExtern{m: def["0"]}, nil
}

func (m *wazeroMemoryExtern) Memptr() uintptr {
	return 0
}
func (m *wazeroMemoryExtern) Name() string {
	return m.wazeroExtern.parent.Name()
}

func (i *wazeroInstance) Name() string {
	return i.m.Name()
}

func (i *wazeroInstance) GetEntryPointExport() (EntryPointExtern, error) {
	fn := i.m.ExportedFunction(EntryPoint)
	if fn == nil {
		return nil, fmt.Errorf("unable to find exported symbol '%s' in %s", EntryPoint, i.Name())
	}
	epoint := &wazeroEntryPointExtern{}
	epoint.wazeroExtern = &wazeroExtern{parent: i}
	epoint.f = fn

	return epoint, nil
}

func (m *wazeroModule) NewInstance() (Instance, error) {
	mod, err := m.parent.r.InstantiateModule(bg, m.cm, nil)
	if err != nil {
		return nil, err
	}
	malloc := mod.ExportedFunction("malloc")
	if malloc == nil {
		return nil, fmt.Errorf("unable to find exported symbol '%s' in %s", "malloc", m.Name())
	}
	free := mod.ExportedFunction(EntryPoint)
	if free == nil {
		return nil, fmt.Errorf("unable to find exported symbol '%s' in %s", "free", m.Name())
	}

	i := &wazeroInstance{
		parent: m,
		m:      mod,
		malloc: malloc,
		free:   free,
	}
	return i, nil
}

func (e *wazeroEng) AddSupportedFunc(pkg, name string, raw interface{}) {
	mod, ok := e.builder[pkg]
	if !ok {
		mod = e.r.NewHostModuleBuilder(pkg)
		e.builder[pkg] = mod
	}
	fn := raw.(func(uint32))
	err := mod.NewFunctionBuilder().WithFunc(fn).Export(name)
	if err != nil {
		panic(fmt.Sprintf("unable to create supported (host) function %s.%s: %v", pkg, name, err))
	}
}

func (e *wazeroEntryPointExtern) Run(argv []string, extra interface{}) (any, error) {
	argptr := make([]uint64, len(argv))
	for i, arg := range argv {
		enc := api.EncodeI32(int32(len(arg)))
		result, err := e.parent.malloc.Call(bg, enc)
		if err != nil {
			return nil, err
		}
		ptr := result[0]
		argptr[i] = ptr
	}
	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	defer func() {
		for _, ptr := range argptr {
			e.parent.free.Call(bg, ptr)
		}
	}()
	result, err := e.f.Call(bg, argptr...)
	if err != nil {
		return nil, err
	}
	log.Printf("\n\ngot a return value... ")
	for i, r := range result {
		log.Printf("\tresult %02d:%x", i, r)
	}
	return nil, nil
}
