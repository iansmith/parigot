package main

import (
	"fmt"
	"log"
	"os"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
	"github.com/iansmith/parigot/sys/abi"
)

var libs = []string{}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("pass one wasm file")
	}
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)
	module, err := wasmtime.NewModuleFromFile(engine, os.Args[1])
	check(err)
	for _, exp := range module.Exports() {
		print("exp:", exp.Name(), "\n")
	}
	for _, imp := range module.Imports() {
		print("imp:", imp.Module(), ",", *imp.Name(), "\n")
	}
	wrappers := generateWrappersForABI(store)
	instance, err := wasmtime.NewInstance(store, module, wrappers)
	check(err)

	start := instance.GetExport(store, "_start").Func()
	result, err := start.Call(store)
	check(err)
	fmt.Printf("done! %v\n", result)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func generateWrappersForABI(store wasmtime.Storelike) []wasmtime.AsExtern {
	result := []wasmtime.AsExtern{}
	result = append(result, wasmtime.WrapFunc(store, abi.TinyGoNotImplemented))
	result = append(result, wasmtime.WrapFunc(store, abi.JSHandleEvent))
	result = append(result, wasmtime.WrapFunc(store, abi.JSNotImplemented))
	result = append(result, wasmtime.WrapFunc(store, abi.SetNow))
	result = append(result, wasmtime.WrapFunc(store, abi.Now))
	result = append(result, wasmtime.WrapFunc(store, abi.OutputString))

	return result
}
