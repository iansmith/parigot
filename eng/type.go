package eng

import (
	"context"
	"errors"

	"github.com/iansmith/parigot/api/shared/id"
	"github.com/tetratelabs/wazero/api"
)

// ErrNotFound is returned when a resource was requested by name
// and it could not be found.
var ErrNotFound = errors.New("unable to find requested object")

// ErrOutOfRange is returned when you attempt to read from guest memory
// addresses that are out of range.
var ErrOutOfRange = errors.New("attempted to read memory that is out of range")

// Environment is an interface that allows a module to be initialized properly.
// This should contain inofrmation that would normally come from the environment
// has to be explicitly created because we are in a wasm container.
type Environment interface {
	Host() id.HostId
	Environment() map[string]string
	Arg() []string
}

type Engine interface {
	NewModuleFromFile(ctx context.Context, path string, env Environment) (Module, error)
	// AddSupportedFunc defines a function that is implemented on the host.
	// The only version AddSupportedFunc() that does not have the suffix
	// is the one that is the "standard" one for exchanging protobufs.
	AddSupportedFunc(ctx context.Context, pkg, name string, fn func(context.Context, api.Module, []uint64))
	// AddSupportedFunc_i32_v is AddSupportedFunc with one input parameter (int32)
	// and no return value.
	AddSupportedFunc_i32_v(ctx context.Context, pkg, name string, fn func(context.Context, api.Module, []uint64))
	// AddSupportedFunc_i32_v is AddSupportedFunc with one input parameter (int32)
	// and no return value.
	AddSupportedFunc_7i32_v(ctx context.Context, pkg, name string, fn func(context.Context, api.Module, []uint64))

	// InstantiateHost module maps the module called into the memory space
	// of the engine.  Any exported functions by AddSupportFunc*() will
	// available with a host implementation but callable from the guest.
	InstantiateHostModule(ctx context.Context, pkg string) (Instance, error)
	// InstanceByName returns an already instantiated instance or an error.
	// The instance must have come from the NewInstance() call at
	// the time it was created.  If the name cannot be found the return
	// value is the error NotFound.
	InstanceByName(ctx context.Context, name string) (Instance, error)
	// HasHostSideFunction returns true if the package named pkg has declared
	// host-side functions in its Init() method. If this is false, there is
	// no reason to instantiate a host module.
	HasHostSideFunction(ctx context.Context, pkg string) bool
}

type Module interface {
	// NewInstance creates an Intance in the host machine and
	// if this returns without error that instance has been
	// initialized, linked, and start_ called without error.
	// XXX The timezone parameter must be passed in from the
	// XXX outside.  The timezoneDir is a path on the _host_ that should be
	// XXX mounted into the guest fs as /tz.  Usually want this to
	// XXX to be GOROOT/lib/time.
	NewInstance(ctx context.Context, timezone string, timezoneDir string) (Instance, error)
	// Name returns the path of the binary that was loaded.
	Name() string
}

type Instance interface {
	Name() string
	// Memory currently only can handle one exported memory
	// section.  With the returned MemoryExtern it is possible to
	// write and read from the guest's memory.
	Memory(ctx context.Context) ([]MemoryExtern, error)
	// Exported function returns a FunctionExtern associated with this module.
	// Empty string is a valid name for a function in Wasm, even if that is
	// a really bad idea.
	//
	// This functions NotFound if the function cannot be found inside
	// the instance.
	//Function(ctx context.Context, name string) (FunctionExtern, error)
	// GetEntryPointExport returns the entry point function
	// even if you don't know its exact name.
	EntryPoint(ctx context.Context) (EntryPointExtern, error)

	// ValueReturns a reference to a value from the instance. This is
	// used to grab values from a previously running instance.  Using this
	// while the instance is running is undefined and a bad idea.
	Value(ctx context.Context, name string) (ExternValue, error)
}

// Extern is the base interface all the Extern types.  It only
// has a name for use in debugging.
type Extern interface {
	Name() string
}

type ExternValue interface {
	GetU64() uint64
	GetU16() uint16
}

// FunctionExtern represents an exposed function be the module.  Note that Go modules compiled by
// the go compiler (not tinygo) only expose "_start" and nothing else.  Normally, you use
// EntryPointExtern to access this function.
type FunctionExtern interface {
	// Call invokes this function with the value given. Note that this
	// function is on the guest side, but this function Call is invoked
	// on the host side. Any error returned is related to problems
	// on the host side.
	Call(ctx context.Context, arg ...uint64) ([]uint64, error)
}

// MemoryExtern is a wrapper around a Wasm Memory object from an instance.
// Currently there is support for only one Memory object per module.
// Note that the host cannot be sure if the guest that owns this
// Memory object is 32bit or 64bit and further does not if the
// memory layout of the guest is little endian or big endian.
type MemoryExtern interface {
	Extern
	// WriteUint64LittleEndian is a utility for copying a uint64 into the
	// (only) memory of the instance.
	WriteUint64LittleEndian(memoryOffset uint32, value uint64)
	// ReadBytes reads the number of bytes requested from the memory
	// associated with this Extern.  There is no way to know if the data
	// is valid.  This moving data from the guest world to the host world.
	// If there is an error, it is likely to be an OutOfRange error.
	ReadBytes(memoryOffset uint32, len uint32) ([]byte, error)
	// ReturnData allocates space for a return value in the host address
	// space.  The values provided should be a serialized protobuf structure
	// and an error id.  Either one, but not both of these can be nil.
	//
	// It is expected that the memory returned from this function
	// will be marked as "dont collect" for the guest GC and will
	// be released by the host code that actually uses the result.
	//ReturnData(ctx context.Context, msg proto.Message, idOrError id.Id) (int32, error)
}

type EntryPointExtern interface {
	FunctionExtern
	// Run has extra parameters that are specific to the paritcular wasm engine.
	Run(ctx context.Context, argv []string, extra interface{}) (uint8, error)
}

type Utility interface {
	// DecodeI32 converts a 64 bit quantity received from a guest
	// call into a int32 that is ok to use on the client side.
	DecodeI32(value uint64) int32
	// DecodeU32 converts a 64 bit quantity received from a guest
	// call into a unsigned int32 that is ok to use on the client side.
	DecodeU32(value uint64) uint32
}

// Util is a set of utility functions that is only available
// (non-nil) after you have called New*Engine. If New*Engine()
// is called multiple times, this value can be set to a
// different value, but it should not make a difference
// if Utility is implemented properly.
var Util Utility
