package main

import (
	"fmt"

	"github.com/iansmith/parigot/command/transform"
)

var abi = []transform.TypeDescriptor{
	transform.FuncToDescriptor(abiimpl.OutputString),
	transform.FuncToDescriptor(abiimpl.Now),
	transform.FuncToDescriptor(abiimpl.Exit),
}

func main() {
	fmt.Print("// Code generated DO NOT EDIT\n")
	fmt.Print("package main\n")
	fmt.Print("import (\n")
	fmt.Print("\twasmtime \"github.com/bytecodealliance/wasmtime-go\"\n")
	fmt.Print("\t\"github.com/iansmith/parigot/sys/abi_impl\"\n")
	fmt.Print(")\n")
	fmt.Printf("func addABIToStore(store wasmtime.Storelike, memPtr uintptr, linkage map[string]*wasmtime.Func) {\n\n")
	for _, desc := range abi {
		fmt.Print(desc.APIWrapper())
	}
	fmt.Printf("}\n")
}
