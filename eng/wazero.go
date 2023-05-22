package eng

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/id"
	"github.com/iansmith/parigot/sharedconst"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"google.golang.org/protobuf/proto"
)

var _ = unsafe.Pointer(nil)

const maxWasmFile = 0x1024 * 0x1024 * 0x20

// ErrorOrIdOffset is how far PAST the pointer that points to
// a ReturnValue
const ErrorOrIdOffset = 8

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

var bg = context.Background()

// NewWazeroEngine creates a new eng.Instance that uses wazer as the
// underlying wasm compiler/interpreter.
func NewWaZeroEngine(ctx context.Context, conf wazero.RuntimeConfig) Engine {
	e := &wazeroEng{}
	if conf != nil {
		e.r = wazero.NewRuntimeWithConfig(bg, conf)
	} else {
		e.r = wazero.NewRuntime(bg)
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
		log.Printf("unable to find instance %s (%d entries)", name, len(e.instantiated))
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

// Note that in wasm the name of the function can be "".
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
		return nil, err
	}
	return result, nil
}

// ReturnData creates a space for a protobuf return value inside the memory space of this
// the parent instance.  This allocates space for the return value and returns
// a pointer to the space as an int32. The error return here is
// for the HOST side only, meaning that some part of the process of creating
// the space in the GUEST side failed.  You must supply at least one of
// msg and err.
func (e *wazeroMemoryExtern) ReturnData(ctx context.Context, msg proto.Message, err id.Id) (int32, error) {
	if msg == nil && !err.IsError() {
		return 0, fmt.Errorf("ReturnSpace called on %s but neither return value nor error value provided", e.Name())
	}
	buf, marshErr := proto.Marshal(msg)
	if err != nil {
		return 0, marshErr
	}
	l := uint64(len(buf))
	if l == 0 && !err.IsError() {
		return 0, fmt.Errorf("ReturnSpace called on %s but empty data value and no error value provided", e.Name())
	}

	res, callErr := e.parent.returnData.Call(ctx, l)
	if callErr != nil {
		return 0, fmt.Errorf("%s called on %s but call to %s failed: %v", sharedconst.ReturnDataName, e.parent.Name(), e.Name(), callErr.Error())
	}
	if len(res) != 1 {
		return 0, fmt.Errorf("%s called on %s but wrong number of return values (%d)", sharedconst.ReturnDataName, e.Name(), len(res))
	}
	result := Util.DecodeI32(res[0])
	if err == nil {
		err = id.NewKernelError(id.KernelNoError)
	}
	// we are doing this for a 32 bit machine, hope this works
	e.WriteUint64LittleEndian(uint32(result+sharedconst.ReturnDataIdErrOffset), err.High())
	e.WriteUint64LittleEndian(uint32(result+sharedconst.ReturnDataIdErrOffset), err.High())
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
	mod, err := e.r.CompileModule(bg, all)

	if err != nil {
		return nil, err
	}
	return &wazeroModule{cm: mod, parent: e, name: path}, nil
}

func (i *wazeroInstance) addMemory(ctx context.Context) error {
	ext, err := i.memoryExportJustOne()
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

func (i *wazeroInstance) memoryExportJustOne() ([]MemoryExtern, error) {
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
	single := newMemoryExtern(i.m.ExportedMemory("0"), i, "\"0\"")
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
	entry := &wazeroEntryPointExtern{
		wazeroFunctionExtern: newFunctionExtern(i.main.fn, i, sharedconst.EntryPoint).(*wazeroFunctionExtern),
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
		log.Printf("xxx -- created instantiated host module %s", pkg)
		return nil, err
	}
	log.Printf("xxx -- created instantiated host module %s", pkg)
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
	log.Printf("about to inst module")
	mod, err := m.parent.r.InstantiateModule(bg, m.cm, wazero.NewModuleConfig())
	log.Printf("done inst module: %v", err)
	if err != nil {
		return nil, err
	}
	i := &wazeroInstance{
		parent:     m,
		m:          mod,
		returnData: nil,
	}
	if err := i.addInstanceInternalFunctions(ctx); err != nil {
		return nil, err
	}
	if err := i.addMemory(ctx); err != nil {
		return nil, err
	}
	log.Printf("adding xxx --- new module %s", mod.Name())
	m.parent.instantiated[mod.Name()] = i.m
	return i, nil
}

func (e *wazeroEng) addSupportFuncAnyType(ctx context.Context, pkg, name string, fn api.GoModuleFunction, iType []api.ValueType, oType []api.ValueType) {
	mod, ok := e.builder[pkg]
	if !ok {
		log.Printf("xxx -- Add supported func, adding new module %s", pkg)
		mod = e.r.NewHostModuleBuilder(pkg)
		e.builder[pkg] = mod
	}
	mod.NewFunctionBuilder().WithGoModuleFunction(fn, iType, oType).Export(name)
}

func (e *wazeroEng) AddSupportedFunc(ctx context.Context, pkg, name string, raw func(context.Context, api.Module, []uint64)) {
	e.addSupportFuncAnyType(ctx, pkg, name, api.GoModuleFunc(raw), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32})
}
func (e *wazeroEng) AddSupportedFunc_i32_v(ctx context.Context, pkg, name string, raw func(context.Context, api.Module, []uint64)) {
	e.addSupportFuncAnyType(ctx, pkg, name, api.GoModuleFunc(raw), []api.ValueType{api.ValueTypeI32}, nil)
}
func (e *wazeroEng) AddSupportedFunc_7i32_v(ctx context.Context, pkg, name string, raw func(context.Context, api.Module, []uint64)) {
	e.addSupportFuncAnyType(ctx, pkg, name, api.GoModuleFunc(raw), []api.ValueType{
		api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, nil)
}

func (e *wazeroEntryPointExtern) Run(ctx context.Context, argv []string, extra interface{}) (any, error) {
	result, err := e.wazeroFunctionExtern.Call(bg, []uint64{0}...)
	if err != nil {
		return nil, err
	}
	log.Printf("\n\ngot a return value... ")
	for i, r := range result {
		log.Printf("\tresult %02d:%x", i, r)
	}
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
	funcExt, expErr := i.Function(ctx, sharedconst.ReturnDataName)
	if expErr != nil {
		return expErr
	}
	i.returnData = funcExt.(*wazeroFunctionExtern)

	// funcExt, expErr = i.Function(ctx, sharedconst.ReturnDataName)
	// if expErr != nil {
	// 	return expErr
	// }
	// i.main = funcExt.(*wazeroFunctionExtern)

	// funcExt, expErr = i.Function(ctx, sharedconst.NewStringName)
	// if expErr != nil {
	// 	return expErr
	// }
	// i.newString = funcExt.(*wazeroFunctionExtern)

	return nil
}

// Right now, part of the code will not build with go 1.20 or 1.21 and
// so we are limited to go1.19.  That version did not have the sliceData
// utility function in unsafe.
func sliceData(p []byte) uintptr {
	x := unsafe.Pointer(&p)
	sh := (*reflect.SliceHeader)(x)
	return sh.Data
}
