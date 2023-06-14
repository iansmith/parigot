package eng

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"unsafe"

	apishared "github.com/iansmith/parigot/apishared"
	pcontext "github.com/iansmith/parigot/context"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"
	"github.com/tetratelabs/wazero/experimental/logging"
)

var _ = unsafe.Pointer(nil)

const maxWasmFile = 0x1024 * 0x1024 * 0x20

// ErrorOrIdOffset is how far PAST the pointer that points to
// a ReturnValue
const ErrorOrIdOffset = 8

var AsyncInteraction = NewAsyncClientInteraction(pcontext.ServerGoContext(pcontext.NewContextWithContainer(context.Background(), "asynchInteraction")))

type wazeroEng struct {
	r            wazero.Runtime
	builder      map[string]wazero.HostModuleBuilder
	instantiated map[string]api.Module
}
type wazeroModule struct {
	parent *wazeroEng
	cm     wazero.CompiledModule
	name   string
	host   bool
}

type wazeroInstance struct {
	parent     *wazeroModule
	m          api.Module
	returnData *wazeroFunctionExtern
	main       *wazeroFunctionExtern
	newString  *wazeroFunctionExtern
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

type wazeroEntryPointExtern struct {
	*wazeroFunctionExtern
}

var rawLineContext = pcontext.NewContextWithContainer(context.Background(), "var wazero:rawLineContext")
var bridgeWriter = newRawLineReader(rawLineContext, pcontext.Wazero)
var WithLogCtx = context.WithValue(context.Background(), experimental.FunctionListenerFactoryKey{},
	logging.NewHostLoggingListenerFactory(bridgeWriter, logging.LogScopeAll))

// NewWazeroEngine creates a new eng.Instance that uses wazer as the
// underlying wasm compiler/interpreter.
func NewWaZeroEngine(withLogContext context.Context, conf wazero.RuntimeConfig) Engine {
	e := &wazeroEng{
		instantiated: make(map[string]api.Module),
	}

	if conf != nil {
		e.r = wazero.NewRuntimeWithConfig(withLogContext, conf)
	} else {
		e.r = wazero.NewRuntime(withLogContext)
	}

	// XXX need to make wasi optional
	wasiBuilder := fakeWasiAddFunc(e)

	_, err := wasiBuilder.Instantiate(withLogContext)
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
		pcontext.Errorf(ctx, "unable to find instance %s (%d entries)", name, len(e.instantiated))
		return nil, ErrNotFound
	}
	inst := &wazeroInstance{
		parent: nil,
		m:      mod,
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
func (e *wazeroInstance) Function(ctx context.Context, name string) (FunctionExtern, error) {
	f := e.m.ExportedFunction(name)
	if f == nil {
		log.Printf("unable to find exported Function %s in inst '%s' mod '%s'", name, e.Name(), e.parent.Name())
		return nil, ErrNotFound
	}
	return &wazeroFunctionExtern{
		wazeroExtern: newWazeroExtern(f.Definition().DebugName()),
		parent:       e,
		fn:           f,
	}, nil
}

func (e *wazeroFunctionExtern) Call(ctx context.Context, param ...uint64) ([]uint64, error) {
	result, err := e.fn.Call(ctx, param...)
	if err != nil {
		pcontext.Errorf(ctx, "call returned an err >>> %s", err.Error())
		return nil, err
	}
	return result, nil
}
func (e *wazeroMemoryExtern) WriteUint64LittleEndian(memoryOffset uint32, value uint64) {
	e.m.WriteUint64Le(memoryOffset, value)
}

func (e *wazeroEng) NewModuleFromFile(ctx context.Context, path string) (Module, error) {
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
	mod, err := e.r.CompileModule(ctx, all)

	if err != nil {
		return nil, err
	}
	return &wazeroModule{cm: mod, parent: e, name: path}, nil
}

func (i *wazeroInstance) addMemory(ctx context.Context) error {
	ext, err := i.memoryExportJustOne(pcontext.CallTo(ctx, "memoryExportJustOne"))
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
	def := i.m.ExportedMemoryDefinitions()
	if def == nil {
		return nil, fmt.Errorf("module '%s' has no exported memory definitions", i.Name())
	}
	if len(def) > 1 {
		panic("parigot currently only supports one memory export, but you should file a ticket to remind us to fix that")
	}
	candidate := i.m.ExportedMemory("memory")
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
	return i.m.Name()
}

func (i *wazeroInstance) EntryPoint(ctx context.Context) (EntryPointExtern, error) {
	fn := i.m.ExportedFunction(apishared.EntryPointSymbol)
	if fn == nil {
		pcontext.Errorf(ctx, "unable to find exported function '%s' in module '%s'", apishared.EntryPointSymbol, i.m.Name())
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
		panic(fmt.Sprintf("unknown builder '%s'", pkg))
	}
	inst, err := b.Instantiate(ctx)
	if err != nil {
		pcontext.Errorf(ctx, "failed to create instantiated host module '%s': %v", pkg, err)
		return nil, err
	}
	pcontext.Logf(ctx, pcontext.Info, "created instantiated host module '%s'", pkg)
	mod := &wazeroModule{parent: m, host: true}
	return &wazeroInstance{
		parent: mod,
		m:      inst,
	}, nil
}

func (m *wazeroInstance) Memory(ctx context.Context) ([]MemoryExtern, error) {
	result := make([]MemoryExtern, len(m.mem))
	for i, mem := range m.mem {
		result[i] = mem
	}
	return result, nil
}

func (m *wazeroModule) NewInstance(ctx context.Context) (Instance, error) {
	fsConfig := wazero.NewFSConfig().WithFauxFs(AsyncInteraction, "/parigotvirt/")
	conf := wazero.NewModuleConfig().
		WithStartFunctions().
		WithName(m.Name()).
		WithStdout(newRawLineReader(rawLineContext, pcontext.GuestOut)).
		WithStderr(newRawLineReader(rawLineContext, pcontext.GuestErr)).
		WithStdin(os.Stdin). // xxx this should probably be fixed to be in config file
		WithRandSource(rand.Reader).
		WithFSConfig(fsConfig).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime()

	mod, err := m.parent.r.InstantiateModule(ctx, m.cm, conf)
	if err != nil {
		pcontext.Errorf(ctx, "ERR IS %s", err.Error())
		return nil, err
	}
	i := &wazeroInstance{
		parent:     m,
		m:          mod,
		returnData: nil,
	}
	if err := i.addInstanceInternalFunctions(pcontext.CallTo(ctx, "addInstanceInternalFunctions")); err != nil {
		return nil, err
	}
	if err := i.addMemory(ctx); err != nil {
		return nil, err
	}

	m.parent.instantiated[mod.Name()] = i.m
	return i, nil
}

func (e *wazeroEng) addSupportFuncAnyType(ctx context.Context, pkg, name string, fn api.GoModuleFunction, iType []api.ValueType, oType []api.ValueType) {
	mod, ok := e.builder[pkg]
	if !ok {
		mod = e.r.NewHostModuleBuilder(pkg)
		e.builder[pkg] = mod
	}
	mod.NewFunctionBuilder().WithGoModuleFunction(fn, iType, oType).Export(name)
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

func wazeroContext(ctx context.Context) context.Context {
	tmp := context.WithValue(ctx, pcontext.ParigotSource, pcontext.Wazero)
	return context.WithValue(tmp, pcontext.ParigotFunc, "wazerolog")
}

func (e *wazeroEntryPointExtern) Run(ctx context.Context, argv []string, extra interface{}) (any, error) {
	// go func(c context.Context, a *AsyncClientInteraction) {
	// 	time.Sleep(time.Duration(1) * time.Second)
	// 	c = pcontext.ServerGoContext(pcontext.CallTo(c, "Send(fake)"))
	// 	err := a.Send("methodcall.v1.AddMultiply", &methodcallmsg.AddMultiplyRequest{
	// 		Value0: 27,
	// 		Value1: 918,
	// 		IsAdd:  true,
	// 	})
	// 	if err != nil {
	// 		pcontext.Errorf(c, "unable to push bundle: %v", err)
	// 	}
	// 	pcontext.Dump(rawLineContext)
	// 	pcontext.Dump(AsyncInteraction.origCtx)
	// }(ctx, AsyncInteraction)
	result, err := e.wazeroFunctionExtern.Call(pcontext.NewContextWithContainer(pcontext.GuestContext(ctx), "call of run()"))
	if err != nil {
		return nil, err
	}
	pcontext.Debugf(ctx, "Run", "got a return value from entry point... ")
	for i, r := range result {
		pcontext.Debugf(ctx, "Run", "result %02d:%x", i, r)
	}
	pcontext.Dump(rawLineContext)
	pcontext.Dump(AsyncInteraction.origCtx)
	return nil, nil
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
