package eng

import "context"

type Engine interface {
	NewModuleFromFile(ctx context.Context, path string) (Module, error)
	// AddSupportedFunc defines a function but the type of the function def is engine specific so
	// we have to use interface{} for the function ptr.
	AddSupportedFunc(ctx context.Context, pkg, name string, fn interface{})
}

type Module interface {
	NewInstance(ctx context.Context) (Instance, error)
}

type Instance interface {
	Name() string
	GetMemoryExport(ctx context.Context) (MemoryExtern, error)
	GetEntryPointExport(ctx context.Context) (EntryPointExtern, error)
	GetFunction(ctx context.Context, pkg, name string) (Function, error)
	Allocate(ctx context.Context, size uint32) (uintptr, error)
	Free(ctx context.Context, ptr uintptr) error
}

type Function interface {
	Call(ctx context.Context, i int32) int32
}

type Extern interface {
	Name() string
}

type MemoryExtern interface {
	Extern
}

type EntryPointExtern interface {
	Extern
	// Run has extra parameters that are specific to the paritcular wasm engine.
	Run(ctx context.Context, argv []string, extra interface{}) (any, error)
}
