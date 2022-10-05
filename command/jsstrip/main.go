package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iansmith/parigot/command/transform"
)

const (
	WasiModule    = "wasi_snapshot_preview1"
	parigotModule = "parigot_abi"
	JSModule      = "env"
	FdWrite       = "fd_write"

	parigotFilename           = "parigot-transformed.wat"
	jsNotImplementedImpl      = "$github.com/iansmith/parigot/abi.JSNotImplemented"
	tinygoNotImplementedImpl  = "$github.com/iansmith/parigot/abi.TinyGoNotImplemented"
	jsNotImplementedImpl1     = "$github.com/iansmith/parigot/abi.JSNotImplemented1"
	tinygoNotImplementedImpl1 = "$github.com/iansmith/parigot/abi.TinyGoNotImplemented1"
	jsEventImpl               = "$github.com/iansmith/parigot/abi.JSHandleEvent"

	jsNotImplImportName      = "JSNotImplemented"
	tinygoNotImplImportName  = "TinyGoNotImplemented"
	jsNotImplImportName1     = "JSNotImplemented1"
	tinygoNotImplImportName1 = "TinyGoNotImplemented1"
	jsEventImportName        = "JSHandleEvent"

	jsHandleEvent = "$syscall/js.handleEvent"
)

var importNameToImplName = map[string]string{
	tinygoNotImplImportName:  tinygoNotImplementedImpl,
	jsNotImplImportName:      jsNotImplementedImpl,
	tinygoNotImplImportName1: tinygoNotImplementedImpl1,
	jsNotImplImportName1:     jsNotImplementedImpl1,
	jsEventImportName:        jsEventImpl,
}

var importNameToTypeNumber = map[string]int{
	tinygoNotImplImportName: 0,
	jsNotImplImportName:     0,
	jsEventImportName:       1,
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

func main() {

	if len(os.Args) != 3 {
		log.Fatalf("unable to understand arguments, should provide exactly two arguments (input wasm file, output wasm file)")
	}

	_, baseArg0 := filepath.Split(os.Args[0])
	tmp, err := os.MkdirTemp("", baseArg0)
	if err != nil {
		log.Fatalf("cannot create temp dir: %v", err)
	}

	watVersion, err := convertWasmToWat(tmp, os.Args[1])
	if err != nil {
		os.Exit(1)
	}
	modifiedWat, err := parigotProcessing(watVersion, tmp)
	if err != nil {
		os.Exit(1)
	}
	err = convertWatToWasm(tmp, modifiedWat, os.Args[2])
	if err != nil {
		os.Exit(1)
	}
	//os.RemoveAll(tmp)
	log.Printf("%s -> %s\n", os.Args[1], watVersion)
	log.Printf("%s -> %s\n", watVersion, filepath.Join(tmp, parigotFilename))
	log.Printf("%s -> %s\n", filepath.Join(tmp, parigotFilename), os.Args[2])
	os.Exit(0)
}

// this is dead code with respect to the compiler tinygo, but we are doing binary
// code patching, so we are creating new "live" functions and they need to be imported
func deadCode(name string) transform.TopLevel {
	imp := &transform.ImportDef{
		ModuleName: parigotModule,
		ImportedAs: name,
		FuncNameRef: &transform.FuncNameRef{
			Name: importNameToImplName[name],
			Type: &transform.TypeRef{
				Num: importNameToTypeNumber[name],
			},
		},
	}
	return imp
}

func changeJSEventHandle(stmt transform.Stmt) transform.Stmt {
	if stmt.StmtType() != transform.OpStmtT ||
		stmt.(transform.Op).OpType() != transform.CallT {
		return stmt
	}
	call := stmt.(*transform.CallOp)
	if call.Arg != jsHandleEvent {
		return stmt
	}
	// we have the target
	call.Arg = jsEventImpl
	return stmt
}

func patchingPass(mod *transform.Module) {
	for _, fnName := range []string{jsNotImplImportName, tinygoNotImplImportName, jsEventImportName, jsNotImplImportName1, tinygoNotImplementedImpl1} {
		addToplevelToModule(mod, deadCode(fnName))
	}
	changeStatementInModule(mod, changeJSEventHandle)
}

func strippingPass(mod *transform.Module) {
	changeToplevelInModule(mod, transform.ImportDefT, changeWasiImportToEmulation)
	changeToplevelInModule(mod, transform.ImportDefT, deleteJSImports)
	changeToplevelInModule(mod, transform.FuncDefT, changeFunctionsToNotImplemented)
	changeToplevelInModule(mod, transform.FuncDefT, deleteFunctionDefinitions)
}

func parigotMangleMethod(category, name string) string {
	return fmt.Sprintf("%s_%s", category, name)
}
func changeWasiImportToEmulation(tl transform.TopLevel) transform.TopLevel {
	imp := tl.(*transform.ImportDef)
	if imp.ModuleName == WasiModule && imp.ImportedAs == FdWrite {
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

var jsErrorCall = &transform.CallOp{
	Arg: jsNotImplementedImpl,
}
var tgErrorCall = &transform.CallOp{
	Arg: tinygoNotImplementedImpl,
}

func changeFunctionsToNotImplemented(tl transform.TopLevel) transform.TopLevel {
	if tl.TopLevelType() != transform.FuncDefT {
		return tl
	}
	fn := tl.(*transform.FuncDef)
	switch fn.Name {
	case "$runtime.printitf":
		fn.Code = []transform.Stmt{tgErrorCall}
	case "$_syscall/js.Value_.String",
		"$_syscall/js.Type_String",
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
