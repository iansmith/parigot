package main

import (
	"fmt"
	wasmtime "github.com/bytecodealliance/wasmtime-go"
	"log"
	"os"
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
