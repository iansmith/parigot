package eng

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
	"unsafe"

	apishared "github.com/iansmith/parigot/api/shared"
	"github.com/iansmith/parigot/api/shared/id"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/sys"
)

var wazerologger *slog.Logger

var _ = unsafe.Pointer(nil)

const maxWasmFile = 0x1024 * 0x1024 * 0x20

// ErrorOrIdOffset is how far PAST the pointer that points to
// a ReturnValue
const ErrorOrIdOffset = 8

type wazeroEng struct {
	rt           wazero.Runtime
	builder      map[string]wazero.HostModuleBuilder
	instantiated map[string]api.Module
}
type wazeroModule struct {
	parent *wazeroEng
	cm     wazero.CompiledModule
	name   string
	host   bool
	env    Environment
}

type wazeroInstance struct {
	parent     *wazeroModule
	mod        api.Module
	returnData *wazeroFunctionExtern
	mem        []*wazeroMemoryExtern
}

type wazeroExtern struct {
	name string
}

type wazeroMemoryExtern struct {
	*wazeroExtern
	parent *wazeroInstance
	m      api.Memory
}

type wazeroFunctionExtern struct {
	*wazeroExtern
	parent *wazeroInstance
	fn     api.Function
}

type wazeroValueExtern struct {
	*wazeroExtern
	parent    *wazeroInstance
	valueName string
	g         api.Global
}

type wazeroEntryPointExtern struct {
	*wazeroFunctionExtern
}

// NewWazeroEngine creates a new eng.Instance that uses wazer as the
// underlying wasm compiler/interpreter.
func NewWaZeroEngine(ctx context.Context, conf wazero.RuntimeConfig) Engine {
	e := &wazeroEng{
		instantiated: make(map[string]api.Module),
	}
	opt := &slog.HandlerOptions{}
	opt.Level = slog.LevelWarn
	wazerologger = slog.New(slog.NewTextHandler(os.Stdout, opt)).With("wazero", true)

	if conf != nil {
		e.rt = wazero.NewRuntimeWithConfig(ctx, conf)
		//wasi_snapshot_preview1.MustInstantiate(ctx, e.rt)
	} else {
		e.rt = wazero.NewRuntime(ctx)
		//wasi_snapshot_preview1.MustInstantiate(ctx, e.rt)
	}

	// XXX need to make wasi optional
	wasiBuilder := fakeWasiAddFunc(e)

	_, err := wasiBuilder.Instantiate(ctx)
	if err != nil {
		log.Fatalf("failed to instantiate wasi override: %v", err)
	}

	// put our utils in place
	Util = &wazeroUtil{}

	// builder is used to create host functions and a host module for
	// each element of this map
	e.builder = make(map[string]wazero.HostModuleBuilder)
	return e
}

func (w *wazeroModule) Name() string {
	return w.name
}

// InstanceByName creates a new eng.Instance based on information
// found inside the engine.
func (e *wazeroEng) InstanceByName(ctx context.Context, name string) (Instance, error) {
	mod, ok := e.instantiated[name]
	if !ok {
		wazerologger.Error("unable to find instance", "name", name, "num choices", len(e.instantiated))
		return nil, ErrNotFound
	}
	inst := &wazeroInstance{
		parent: nil,
		mod:    mod,
	}
	if err := inst.addInstanceInternalFunctions(ctx); err != nil {
		return nil, err
	}
	if err := inst.addMemory(ctx); err != nil {
		return nil, err
	}
	return inst, nil
}

func newWazeroExtern(name string) *wazeroExtern {
	return &wazeroExtern{name: name}
}

func (e *wazeroExtern) Name() string {
	return e.name
}

// // Note that in wasm the name of the function can be "".
func (e *wazeroInstance) Function(name string) (FunctionExtern, error) {
	f := e.mod.ExportedFunction(name)
	if f == nil {
		wazerologger.Error("unable to find exported Function", "name", name, "instance", e.Name(), "module", e.parent.Name())
		return nil, ErrNotFound
	}
	return &wazeroFunctionExtern{
		wazeroExtern: newWazeroExtern(f.Definition().DebugName()),
		parent:       e,
		fn:           f,
	}, nil
}

func (e *wazeroFunctionExtern) Call(ctx context.Context, param ...uint64) ([]uint64, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Call generated panic: %v", r)
		}
	}()
	result, err := e.fn.Call(ctx, param...)
	if err != nil {
		if exitErr, ok := err.(*sys.ExitError); ok {
			// This means your module exited with non-zero code!
			// You can get the code with exitErr.ExitCode()
			wazerologger.Info("exit code returned", "code", exitErr.ExitCode(), "returned values", len(result))
			return []uint64{uint64(exitErr.ExitCode())}, nil
		} else {
			wazerologger.Error("return value from process wasn't exit error", "type", fmt.Sprintf("%T", err))
		}
	}
	return result, nil
}
func (e *wazeroMemoryExtern) WriteUint64LittleEndian(memoryOffset uint32, value uint64) {
	e.m.WriteUint64Le(memoryOffset, value)
}

func (e *wazeroEng) NewModuleFromFile(ctx context.Context, path string, env Environment) (Module, error) {
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
	mod, err := e.rt.CompileModule(ctx, all)

	if err != nil {
		return nil, err
	}
	return &wazeroModule{cm: mod, parent: e, name: path, env: env}, nil
}

func (i *wazeroInstance) addMemory(ctx context.Context) error {
	ext, err := i.memoryExportJustOne(ctx)
	if err != nil {
		return err
	}
	result := make([]*wazeroMemoryExtern, len(ext))
	for j, m := range ext {
		result[j] = m.(*wazeroMemoryExtern)
	}
	i.mem = result
	return nil
}

// Ready turns a length and a pointer in the guest address space into bytes
// in the host address space
func (i *wazeroMemoryExtern) ReadBytes(memoryOffset, length uint32) ([]byte, error) {
	b, outOfRange := i.m.Read(memoryOffset, length)
	if outOfRange {
		return nil, ErrOutOfRange
	}
	return b, nil
}

func (i *wazeroInstance) memoryExportJustOne(ctx context.Context) ([]MemoryExtern, error) {
	def := i.mod.ExportedMemoryDefinitions()
	if def == nil {
		return nil, fmt.Errorf("module '%s' has no exported memory definitions", i.Name())
	}
	if len(def) > 1 {
		panic("parigot currently only supports one memory export, but you should file a ticket to remind us to fix that")
	}
	candidate := i.mod.ExportedMemory("memory")
	if candidate == nil {
		return nil, fmt.Errorf("can't find memory object 'memory'")
	}
	single := newMemoryExtern(candidate, i, "'memory'")
	return []MemoryExtern{single}, nil
}

func newMemoryExtern(mem api.Memory, i *wazeroInstance, name string) MemoryExtern {
	return &wazeroMemoryExtern{
		wazeroExtern: newWazeroExtern(fmt.Sprintf("memory[%s]", name)),
		parent:       i,
		m:            mem,
	}
}

func (m *wazeroMemoryExtern) Name() string {
	return m.wazeroExtern.Name()
}

func (i *wazeroInstance) Name() string {
	return i.mod.Name()
}

func (i *wazeroInstance) Value(ctx context.Context, valueName string) (ExternValue, error) {
	g := i.mod.ExportedGlobal(valueName)
	if g == nil {
		return nil, fmt.Errorf("unable to find value '%s' in module '%s' ",
			valueName, i.parent.name)
	}
	return &wazeroValueExtern{g: g, valueName: valueName}, nil
}

func (e *wazeroValueExtern) GetU64() uint64 {
	return e.g.Get()
}

func (e *wazeroValueExtern) GetU16() uint16 {
	raw := e.g.Get()
	u32 := api.DecodeU32(raw)
	return uint16(u32 & 0xffff)
}

func (i *wazeroInstance) EntryPoint(ctx context.Context) (EntryPointExtern, error) {
	fn := i.mod.ExportedFunction(apishared.EntryPointSymbol)
	if fn == nil {
		wazerologger.Error("unable to find exported entry point", "name", apishared.EntryPointSymbol, "module", i.mod.Name())
		return nil, ErrNotFound
	}
	ext := newFunctionExtern(fn, i, fn.Definition().DebugName()).(*wazeroFunctionExtern)

	entry := &wazeroEntryPointExtern{
		wazeroFunctionExtern: ext,
	}
	return entry, nil
}

func newFunctionExtern(fn api.Function, i *wazeroInstance, name string) FunctionExtern {
	return &wazeroFunctionExtern{
		wazeroExtern: newWazeroExtern(name),
		parent:       i,
		fn:           fn,
	}
}

func (m *wazeroEng) InstantiateHostModule(ctx context.Context, pkg string) (Instance, error) {
	b, ok := m.builder[pkg]
	if !ok {
		wazerologger.Info("not instantiating host module", "package", pkg)
		return nil, nil // they have not declared any functions
	}
	inst, err := b.Instantiate(ctx)
	if err != nil {
		wazerologger.Error("failed to create instantiated host module", err, "package", pkg)
		return nil, err
	}
	mod := &wazeroModule{parent: m, host: true}
	return &wazeroInstance{
		parent: mod,
		mod:    inst,
	}, nil
}

func (m *wazeroInstance) Memory(ctx context.Context) ([]MemoryExtern, error) {
	result := make([]MemoryExtern, len(m.mem))
	for i, mem := range m.mem {
		result[i] = mem
	}
	return result, nil
}

func (m *wazeroModule) NewInstance(ctx context.Context, timezone string, timezoneDir string, hid id.HostId) (Instance, error) {
	args := []string{}
	envp := make(map[string]string)

	if m.env != nil {
		args = m.env.Arg()
		envp = m.env.Environment()
	}

	wazerologger.Info("created new wasm module instance", "name", m.name, "host", hid.Short())
	conf := wazero.NewModuleConfig().
		WithStartFunctions().
		WithName(m.Name()).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		// WithStdout(newRawLineReader(rawLineContext, pcontext.GuestOut)).
		// WithStderr(newRawLineReader(rawLineContext, pcontext.GuestErr)).
		// WithStdin(os.Stdin). // xxx this should probably be fixed to be in config file
		WithRandSource(rand.Reader).
		WithFSConfig(wazero.NewFSConfig()).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithArgs(strings.Join(append([]string{m.name}, args...), " ")).
		WithEnv("HOSTID_HIGH", fmt.Sprintf("%x", hid.High())).
		WithEnv("HOSTID_LOW", fmt.Sprintf("%x", hid.Low())).
		WithEnv("TZ", timezone)

	for k, v := range envp {
		conf.WithEnv(k, v)
	}

	mod, err := m.parent.rt.InstantiateModule(ctx, m.cm, conf)
	if err != nil {
		wazerologger.Error("instantiate module failed in wazero runtime", "error", err, "name", m.name)
		return nil, err
	}
	wazerologger.Info("instantiate module success", "name", m.name)

	i := &wazeroInstance{
		parent:     m,
		mod:        mod,
		returnData: nil,
	}
	if err := i.addInstanceInternalFunctions(ctx); err != nil {
		return nil, err
	}
	if err := i.addMemory(ctx); err != nil {
		return nil, err
	}

	m.parent.instantiated[mod.Name()] = i.mod
	return i, nil
}

func (e *wazeroEng) AddBuilder(ctx context.Context, builderName string) wazero.HostModuleBuilder {
	b, ok := e.builder[builderName]
	if ok {
		return b
	}
	mod := e.rt.NewHostModuleBuilder(builderName)
	e.builder[builderName] = mod
	return mod
}

func (e *wazeroEng) addSupportFuncAnyType(ctx context.Context, pkg, name string, fn api.GoModuleFunction, iType []api.ValueType, oType []api.ValueType) {
	mod, ok := e.builder[pkg]
	if !ok {
		mod = e.AddBuilder(ctx, pkg)
	}
	mod.NewFunctionBuilder().WithGoModuleFunction(fn, iType, oType).Export(name)
}

func (e *wazeroEng) HasHostSideFunction(ctx context.Context, pkg string) bool {

	_, ok := e.builder[pkg]
	return ok
}

func (e *wazeroEng) AddSupportedFunc(ctx context.Context, pkg, name string, raw func(context.Context, api.Module, []uint64)) {
	e.addSupportedFunc_i32i32i32i32_i64(ctx, pkg, name, api.GoModuleFunc(raw))
}
func (e *wazeroEng) addSupportedFunc_i32i32i32i32_i64(ctx context.Context, pkg, name string, raw func(context.Context, api.Module, []uint64)) {
	e.addSupportFuncAnyType(ctx, pkg, name, api.GoModuleFunc(raw), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI64})
}

func (e *wazeroEng) AddSupportedFunc_i32_v(ctx context.Context, pkg, name string, raw func(context.Context, api.Module, []uint64)) {
	e.addSupportFuncAnyType(ctx, pkg, name, api.GoModuleFunc(raw), []api.ValueType{api.ValueTypeI32}, nil)
}
func (e *wazeroEng) AddSupportedFunc_7i32_v(ctx context.Context, pkg, name string, raw func(context.Context, api.Module, []uint64)) {
	e.addSupportFuncAnyType(ctx, pkg, name, api.GoModuleFunc(raw), []api.ValueType{
		api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, nil)
}

func (e *wazeroEntryPointExtern) Run(ctx context.Context, argv []string, extra interface{}) (uint8, error) {
	retVal, err := e.wazeroFunctionExtern.Call(ctx)
	if retVal[0] > 192 {
		return uint8(retVal[0] & 0xff), err
	}
	// if err != nil {
	// 	return psys.ExitCodeTrapped, err
	// }
	// if len(result) > 0 {
	// 	for i, r := range result {
	// 		wazerologger.Info("return value from entry point", "index", i, "result", r)
	// 	}
	// }
	code := uint8(retVal[0] & 0xff)
	return code, nil
}

// Utility
type wazeroUtil struct{}

func (u *wazeroUtil) DecodeI32(value uint64) int32 {
	return api.DecodeI32(value)
}
func (u *wazeroUtil) DecodeU32(value uint64) uint32 {
	return api.DecodeU32(value)
}

func (i *wazeroInstance) addInstanceInternalFunctions(ctx context.Context) error {
	return nil
}
