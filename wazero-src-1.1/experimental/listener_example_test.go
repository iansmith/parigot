package experimental_test

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"sort"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/internal/wasm"
)

// listenerWasm was generated by the following:
//
//	cd testdata; wat2wasm --debug-names listener.wat
//
//go:embed logging/testdata/listener.wasm
var listenerWasm []byte

// uniqGoFuncs implements both FunctionListenerFactory and FunctionListener
type uniqGoFuncs map[string]struct{}

// callees returns the go functions called.
func (u uniqGoFuncs) callees() []string {
	ret := make([]string, 0, len(u))
	for k := range u {
		ret = append(ret, k)
	}
	// Sort names for consistent iteration
	sort.Strings(ret)
	return ret
}

// NewListener implements FunctionListenerFactory.NewListener
func (u uniqGoFuncs) NewListener(def api.FunctionDefinition) experimental.FunctionListener {
	if def.GoFunction() == nil {
		return nil // only track go funcs
	}
	return u
}

// Before implements FunctionListener.Before
func (u uniqGoFuncs) Before(ctx context.Context, _ api.Module, def api.FunctionDefinition, _ []uint64, _ experimental.StackIterator) context.Context {
	u[def.DebugName()] = struct{}{}
	return ctx
}

// After implements FunctionListener.After
func (u uniqGoFuncs) After(context.Context, api.Module, api.FunctionDefinition, error, []uint64) {}

// This shows how to make a listener that counts go function calls.
func Example_customListenerFactory() {
	u := uniqGoFuncs{}

	// Set context to one that has an experimental listener
	ctx := context.WithValue(context.Background(), experimental.FunctionListenerFactoryKey{}, u)

	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx) // This closes everything this Runtime created.

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	mod, err := r.Instantiate(ctx, listenerWasm)
	if err != nil {
		log.Panicln(err)
	}

	for i := 0; i < 5; i++ {
		if _, err = mod.ExportedFunction("rand").Call(ctx, 4); err != nil {
			log.Panicln(err)
		}
	}

	// A Go function was called multiple times, but we should only see it once.
	for _, f := range u.callees() {
		fmt.Println(f)
	}

	// Output:
	// wasi_snapshot_preview1.fd_write
	// wasi_snapshot_preview1.random_get
}

func Example_stackIterator() {
	it := &fakeStackIterator{}

	for it.Next() {
		fmt.Println("function:", it.FunctionDefinition().DebugName(), "args", it.Args())
	}

	// Output:
	// function: fn0 args [1 2 3]
	// function: fn1 args []
	// function: fn2 args [4]
}

type fakeStackIterator struct {
	iteration int
	def       api.FunctionDefinition
	args      []uint64
}

func (s *fakeStackIterator) Next() bool {
	switch s.iteration {
	case 0:
		s.def = &mockFunctionDefinition{debugName: "fn0"}
		s.args = []uint64{1, 2, 3}
	case 1:
		s.def = &mockFunctionDefinition{debugName: "fn1"}
		s.args = []uint64{}
	case 2:
		s.def = &mockFunctionDefinition{debugName: "fn2"}
		s.args = []uint64{4}
	case 3:
		return false
	}
	s.iteration++
	return true
}

func (s *fakeStackIterator) FunctionDefinition() api.FunctionDefinition {
	return s.def
}

func (s *fakeStackIterator) Args() []uint64 {
	return s.args
}

type mockFunctionDefinition struct {
	debugName string
	*wasm.FunctionDefinition
}

func (f *mockFunctionDefinition) DebugName() string {
	return f.debugName
}

func (f *mockFunctionDefinition) ParamTypes() []wasm.ValueType {
	return []wasm.ValueType{}
}

func (f *mockFunctionDefinition) ResultTypes() []wasm.ValueType {
	return []wasm.ValueType{}
}
