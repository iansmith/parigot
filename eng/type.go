package eng

type Engine interface {
	NewModuleFromFile(path string) (Module, error)
	// AddSupportedFunc defines a function but the type of the function def is engine specific so
	// we have to use interface{} for the function ptr.
	AddSupportedFunc(pkg, name string, fn interface{})
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
	GetFunction(pkg, name string) (Function, error)
	Allocate(size uint32) (uintptr, error)
	Free(ptr uintptr) error
}

type Function interface {
	Call(int32) int32
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
	// Run has extra parameters that are specific to the paritcular wasm engine.
	Run(argv []string, extra interface{}) (any, error)
}
