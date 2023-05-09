package eng

import (
	"fmt"
	"log"

	"github.com/iansmith/parigot/sys/jspatch"

	"github.com/bytecodealliance/wasmtime-go/v8"
)

// This is the implementation of the eng types for the wasmtime
// embedding in go.
type wasmtimeEngine struct {
	e       *wasmtime.Engine
	s       *wasmtime.Store
	funcMap map[string]*wasmtime.Func
}

type wasmtimeModule struct {
	owner *wasmtimeEngine
	m     *wasmtime.Module
	name  string
}

type wasmtimeInstance struct {
	owner *wasmtimeEngine
	i     *wasmtime.Instance
	name  string
}

type wasmtimeAsExtern struct {
	e     wasmtime.AsExtern
	name  string
	owner *wasmtimeInstance
}

type wasmtimeMemExtern struct {
	*wasmtimeAsExtern
}

type wasmtimeEntryPointExtern struct {
	*wasmtimeAsExtern
}

// Implementation functions
func (i *wasmtimeInstance) Name() string {
	return i.name
}

func (i *wasmtimeInstance) GetMemoryExport() (MemoryExtern, error) {
	inst := i.i
	external := inst.GetExport(i.owner.s, "mem")

	if external == nil {
		return nil, fmt.Errorf("module %s does not export mem, so not the golang compiler", i.Name())
	}
	if external.Memory() == nil {
		return nil, fmt.Errorf("module %s has an exported symbol 'mem' but it's not a memory object", i.Name())
	}
	return &wasmtimeMemExtern{
		wasmtimeAsExtern: &wasmtimeAsExtern{
			e:     external,
			owner: i},
	}, nil
}
func (i *wasmtimeInstance) GetEntryPointExport() (EntryPointExtern, error) {
	inst := i.i
	external := inst.GetExport(i.owner.s, "run")
	if external == nil {
		return nil, fmt.Errorf("module %s does not export 'run', so compiler not understood", i.Name())
	}
	if external.Func() == nil {
		return nil, fmt.Errorf("module %s has an exported symbol 'run' but it's not a function object", i.Name())
	}
	return &wasmtimeEntryPointExtern{
		wasmtimeAsExtern: &wasmtimeAsExtern{
			owner: i,
			e:     external,
		},
	}, nil
}
func NewWasmtimeEngine(c *Config) Engine {
	conf := wasmtime.NewConfig()
	conf.SetDebugInfo(!c.NoDebug)
	switch c.OptLevel {
	case 0:
		conf.SetCraneliftOptLevel(wasmtime.OptLevelSpeed)
	case 1:
		conf.SetCraneliftOptLevel(wasmtime.OptLevelSpeedAndSize)
	case 2:
		conf.SetCraneliftOptLevel(wasmtime.OptLevelNone)
	}
	eng := wasmtime.NewEngineWithConfig(conf)
	result := &wasmtimeEngine{
		e:       eng,
		s:       wasmtime.NewStore(eng),
		funcMap: make(map[string]*wasmtime.Func),
	}
	return result
}

func (e *wasmtimeEngine) NewModuleFromFile(path string) (Module, error) {
	m, err := wasmtime.NewModuleFromFile(e.e, path)
	if err != nil {
		return nil, err
	}
	// for _, export := range m.Exports() {
	// 	_, ok := e.funcMap[export.Name()]
	// 	log.Printf("module %s has EXPORT: %s (%v)", path, export.Name(), ok)
	// }

	return &wasmtimeModule{
		owner: e,
		m:     m,
		name:  path,
	}, nil
}

func (e *wasmtimeEngine) AddSupportedFunc(name string, fn any) {
	e.funcMap[name] = wasmtime.WrapFunc(e.s, fn)
}

func (m *wasmtimeModule) NewInstance() (Instance, error) {
	engine := m.owner

	importCandidate := make(map[string]struct{})
	imports := m.m.Imports()
	for _, imp := range imports {
		name := "$ANON$"
		n := imp.Name()
		if n == nil {
			log.Printf("import has nil name, using $ANON$")
		} else if *n == "" {
			log.Printf("import has no name, using $ANON$")
		} else {
			name = *n
		}
		if imp.Type().FuncType() == nil {
			log.Printf("name is %s but not a function (global? %v)", name, imp.Type().GlobalType() != nil)
			continue
		}
		importCandidate[name] = struct{}{}
	}
	opened := []wasmtime.AsExtern{}
	seen := make(map[string]struct{})

	i := 0
	for n, f := range m.owner.funcMap {
		if _, ok := importCandidate[n]; !ok {
			continue
		}
		seen[n] = struct{}{}
		opened = append(opened, f)
		i++
	}
	for k := range importCandidate {
		delete(seen, k)
	}
	for k := range seen {
		log.Printf("WARNING: no import found for %s", k)
	}
	inst, err := wasmtime.NewInstance(engine.s, m.m, opened)
	if err != nil {
		return nil, err
	}
	return &wasmtimeInstance{i: inst, owner: m.owner}, nil
}

func (e *wasmtimeMemExtern) Memptr() uintptr {
	mem := e.e.(*wasmtime.Extern).Memory()
	return uintptr(mem.Data(e.owner.owner.s))

}
func (e *wasmtimeAsExtern) Name() string {
	return e.name

}

func (e *wasmtimeEntryPointExtern) Run(argv []string, memHigh int32, mem *jspatch.WasmMem) (any, error) {

	f := e.e.(*wasmtime.Extern).Func()
	// argc := int32(len(argv))
	// sizeOfPtrs := 4 * (argc + 1)
	// addrOfArgv := memHigh - sizeOfPtrs // including a zero

	// addrData := make([]int32, len(argv))
	// running := addrOfArgv
	// for i := 0; i < len(argv); i++ {
	// 	addrData[i] = running - int32(len(argv[i])+1)
	// 	log.Printf("xxxx --- Run data for %d @ %x", i, addrData[i])
	// 	b := []byte(argv[i])
	// 	for j := int32(0); j < int32(len(argv[i])); j++ {
	// 		mem.SetUint8(addrData[i]+j, b[j])
	// 	}
	// 	mem.SetUint8(addrData[i]+int32(len(argv[i])), 0)
	// 	running -= int32(len(argv) + 1)
	// 	running %= 4 //aligment, moving upward
	// 	if running < 0 {
	// 		return nil, fmt.Errorf("total size of arguments/environment is too large, max is %d bytes", memHigh)
	// 	}
	// }
	// // copied all the argv
	// for i := 0; i < len(argv); i++ {
	// 	log.Printf("xxxx --- Run ptr for %d @ %x", i, addrData[i])
	// 	mem.SetInt32(addrOfArgv+(4*int32(i)), addrData[i])
	// }
	return f.Call(e.owner.owner.s, int32(0), int32(0))
}
