package main

import (
	"fmt"
	"os"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
)

func main() {
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)
	wd, _ := os.Getwd()
	print("cwd:", wd, "\n")
	module, err := wasmtime.NewModuleFromFile(engine, "build/hello-go.wasm")
	check(err)
	print("exports---\n")
	for _, exp := range module.Exports() {
		print("exp:", exp.Name(), "\n")
	}
	print("imports---\n")
	for _, imp := range module.Imports() {
		print("imp:", imp.Module(), ",", *imp.Name(), "\n")
	}
	instance, err := wasmtime.NewInstance(store, module, []wasmtime.AsExtern{})
	check(err)

	gcd := instance.GetExport(store, "gcd").Func()
	val, err := gcd.Call(store, 6, 27)
	check(err)
	fmt.Printf("gcd(6, 27) = %d\n", val.(int32))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
