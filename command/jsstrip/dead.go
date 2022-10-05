package main

import (
	"github.com/iansmith/parigot/command/transform"
	"log"
	"strings"
)

var tgErrorCall = &transform.CallOp{
	Arg: tinygoNotImplementedImpl,
}

func changeWasiImportToEmulation(tl transform.TopLevel) transform.TopLevel {
	imp := tl.(*transform.ImportDef)
	if imp.ModuleName == WasiModule && imp.ImportedAs == FdWrite {
		imp.ModuleName = "parigot_abi"
		imp.ImportedAs = parigotMangleMethod("WasiEmulation", "FdWrite")
		return imp
	}
	return imp
}
func deleteJSImports(tl transform.TopLevel) transform.TopLevel {
	importDef := tl.(*transform.ImportDef)
	if importDef.ModuleName == JSModule {
		for _, s := range methodsToKillImport {
			if importDef.ImportedAs == s {
				return nil // drops the import
			}
		}
		log.Fatalf("to understand what to do with import of %s from module env (js runtime?)", importDef.ImportedAs)
	}
	return importDef
}

// poisoning? gunshots?
var methodsToKillImport = []string{
	"syscall/js.valueGet",
	"syscall/js.valuePrepareString",
	"syscall/js.valueLoadString",
	"syscall/js.finalizeRef",
	"syscall/js.stringVal",
	"syscall/js.valueSet",
	"syscall/js.valueLength",
	"syscall/js.valueIndex",
	"syscall/js.valueCall",
}

func changeFunctionsToNotImplemented(tl transform.TopLevel) transform.TopLevel {
	if tl.TopLevelType() != transform.FuncDefT {
		return tl
	}
	fn := tl.(*transform.FuncDef)
	switch fn.Name {
	case "$runtime.printitf":
		fn.Code = []transform.Stmt{tgErrorCall}
	case "$_syscall/js.Value_.String", // type 5
		"$_syscall/js.Type_String", // type 2
		"$_syscall/js.Value_.Get":
		fn.Code = []transform.Stmt{jsErrorCall}
	default:
		break
	}
	return fn
}

func deleteFunctionDefinitions(tl transform.TopLevel) transform.TopLevel {
	if tl.TopLevelType() != transform.FuncDefT {
		return tl
	}
	fn := tl.(*transform.FuncDef)
	switch fn.Name {
	case "$runtime.printitf",
		"$_syscall/js.Value_.String",
		"$_syscall/js.Type_String",
		"$_syscall/js.Value_.Get":
		return fn
	default:
		break
	}
	if strings.HasPrefix(fn.Name, "$_syscall/js.") ||
		strings.HasPrefix(fn.Name, "$syscall/js.") ||
		strings.HasPrefix(fn.Name, "$_struct_syscall/js") ||
		strings.HasPrefix(fn.Name, "$_*struct_syscall/js") ||
		strings.HasPrefix(fn.Name, "$_*syscall/js") {
		return nil
	}
	return fn
}
