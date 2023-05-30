package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// counterWasm was generated by the following:
//
//	cd testdata; wat2wasm --debug-names counter.wat
//
//go:embed testdata/counter.wasm
var counterWasm []byte

// main shows how to share the same compilation cache across the multiple runtimes.
func main() {
	// Choose the context to use for function calls.
	ctx := context.Background()

	// Prepare a cache directory.
	cacheDir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.Panicln(err)
	}
	defer os.RemoveAll(cacheDir)

	// Initializes the new compilation cache with the cache directory.
	// This allows the compilation caches to be shared even across multiple OS processes.
	cache, err := wazero.NewCompilationCacheWithDir(cacheDir)
	if err != nil {
		log.Panicln(err)
	}
	defer cache.Close(ctx)

	// Creates a shared runtime config to share the cache across multiple wazero.Runtime.
	runtimeConfig := wazero.NewRuntimeConfig().WithCompilationCache(cache)

	// Creates two wazero.Runtimes with the same compilation cache.
	runtimeFoo := wazero.NewRuntimeWithConfig(ctx, runtimeConfig)
	runtimeBar := wazero.NewRuntimeWithConfig(ctx, runtimeConfig)

	// Instantiate two modules on separate Runtimes with identical configuration, which allows each instance
	// has the isolated states of "env" module.
	m1 := instantiateWithEnv(ctx, runtimeFoo)
	m2 := instantiateWithEnv(ctx, runtimeBar)

	for i := 0; i < 2; i++ {
		fmt.Printf("m1 counter=%d\n", counterGet(ctx, m1))
		fmt.Printf("m2 counter=%d\n", counterGet(ctx, m2))
	}
}

// count calls "counter.get" in the given namespace
func counterGet(ctx context.Context, mod api.Module) uint64 {
	results, err := mod.ExportedFunction("get").Call(ctx)
	if err != nil {
		log.Panicln(err)
	}
	return results[0]
}

// counter is an example showing state that needs to be independent per importing module.
type counter struct {
	counter uint32
}

func (e *counter) getAndIncrement() (ret uint32) {
	ret = e.counter
	e.counter++
	return
}

// instantiateWithEnv returns a module instance.
func instantiateWithEnv(ctx context.Context, r wazero.Runtime) api.Module {
	// Instantiate a new "env" module which exports a stateful function.
	c := &counter{}
	_, err := r.NewHostModuleBuilder("env").
		NewFunctionBuilder().WithFunc(c.getAndIncrement).Export("next_i32").
		Instantiate(ctx)
	if err != nil {
		log.Panicln(err)
	}

	// Instantiate the module that imports "env".
	mod, err := r.Instantiate(ctx, counterWasm)
	if err != nil {
		log.Panicln(err)
	}
	return mod
}
