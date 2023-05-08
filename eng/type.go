package eng

import "github.com/iansmith/parigot/sys/jspatch"

type Engine interface {
	NewModuleFromFile(path string) (Module, error)
	AddSupportedFunc(name string, fn any)
}

type Config struct {
	OptLevel int /* 0=normal, 1=high, 2 = low */
	NoDebug  bool
}

type Module interface {
	NewInstance() (Instance, error)
}

type Instance interface {
	Name() string
	GetMemoryExport() (MemoryExtern, error)
	GetEntryPointExport() (EntryPointExtern, error)
}

type Extern interface {
	Name() string
}

type MemoryExtern interface {
	Extern
	Memptr() uintptr
}

type EntryPointExtern interface {
	Extern
	Run(argv []string, highMem int32, mem *jspatch.WasmMem) (any, error)
}
