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

	parigotFilename = "parigot-transformed.wat"

	tinygoNotImplementedImpl = "$github.com/iansmith/parigot/abi.TinyGoNotImplemented"
	jsNotImplementedImpl     = "$github.com/iansmith/parigot/abi.JSNotImplemented"

	jsNotImplImportName     = "JSNotImplemented"
	tinygoNotImplImportName = "TinyGoNotImplemented"
)

var importNameToImplName = map[string]string{
	tinygoNotImplImportName: tinygoNotImplementedImpl,
	jsNotImplImportName:     jsNotImplementedImpl,
}

var importNameToTypeNumber = map[string]int{
	tinygoNotImplImportName: 0,
	jsNotImplImportName:     0,
}

var typeToReturnValue = make(map[int]int)
var replacementFuncTypeTable = make(map[string]int)
var importFuncTypeTable = make(map[string]int)

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

// addNotImplImports adds our two functions to the list of imports from our
// module that will be linked in later.
func addNotImplImports(mod *transform.Module) {
	for _, name := range []string{jsNotImplImportName, tinygoNotImplImportName} {
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
		addToplevelToModule(mod, imp)
	}
}
func modifyJSImports(tl transform.TopLevel) {
	fn := tl.(*transform.ImportDef)
	if isJSFunction(fn.FuncNameRef.Name) {
		log.Printf("js function imported: (%s:%s) %s:%d\n", fn.ModuleName, fn.ImportedAs, fn.FuncNameRef.Name, fn.FuncNameRef.Type.Num)
		importFuncTypeTable[fn.ImportedAs] = fn.FuncNameRef.Type.Num
	}
}

func mapTypeReturnValues(tl transform.TopLevel) {
	td := tl.(*transform.TypeDef)
	typeToReturnValue[td.Annotation] = 0
	if td.Func.Result != nil {
		if len(td.Func.Result.Type.Name) > 1 {
			panic("unable to handle mapping return types with multiple values")
		}
		rType := td.Func.Result.Type.Name[0]
		log.Printf("return type of type %d is %s", td.Annotation, rType)
		switch rType {
		case "i32":
			typeToReturnValue[td.Annotation] = 32
		case "i64":
			typeToReturnValue[td.Annotation] = 64
		default:
			panic("unable to handle mapping return types with type " + rType)
		}
	}
}
func transformation(mod *transform.Module) {
	findToplevelInModule(mod, transform.FuncDefT, mapReplacementFunctions)
	findToplevelInModule(mod, transform.TypeDefT, mapTypeReturnValues)
	addNotImplImports(mod)
	findToplevelInModule(mod, transform.ImportDefT, modifyJSImports)
	changeToplevelInModule(mod, transform.FuncDefT, changeFuncsToNotImplemented)
}

var returnStmt = &transform.ZeroOp{
	Op: "return",
}

func changeFuncsToNotImplemented(tl transform.TopLevel) transform.TopLevel {
	fn := tl.(*transform.FuncDef)
	if isJSFunction(fn.Name) {
		_, ok := replacementFuncTypeTable[fn.Name]
		log.Printf("replace: %s->%d->%d [%v]", fn.Name, replacementFuncTypeTable[fn.Name],
			typeToReturnValue[replacementFuncTypeTable[fn.Name]], ok)
		if typeToReturnValue[replacementFuncTypeTable[fn.Name]] != 0 {
			op := "i32.const"
			if typeToReturnValue[replacementFuncTypeTable[fn.Name]] == 64 {
				op = "i64.const"
			}
			// for now, we always return
			argOp := &transform.ArgOp{
				Op:     op,
				IntArg: new(int),
			}
			*argOp.IntArg = 0
			fn.Code = []transform.Stmt{
				jsErrorCall,
				argOp,
				returnStmt,
			}
		} else {
			fn.Code = []transform.Stmt{
				jsErrorCall,
			}
		}
	}
	return fn
}

func isJSFunction(name string) bool {
	return strings.Index(name, "/js") != -1
}

func parigotMangleMethod(category, name string) string {
	return fmt.Sprintf("%s_%s", category, name)
}

var jsErrorCall = &transform.CallOp{
	Arg: jsNotImplementedImpl,
}

func mapReplacementFunctions(level transform.TopLevel) {
	fn := level.(*transform.FuncDef)
	if isJSFunction(fn.Name) {
		log.Printf("%s:type %d\n", fn.Name, fn.Type.Num)
		replacementFuncTypeTable[fn.Name] = fn.Type.Num
	}
}
