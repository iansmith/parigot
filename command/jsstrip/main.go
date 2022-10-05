package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iansmith/parigot/command/transform"
)

const (
	WasiModule        = "wasi_snapshot_preview1"
	parigotModule     = "parigot_abi"
	JSModule          = "env"
	FdWrite           = "fd_write"
	oldWasiCallTarget = "$runtime.fd_write"

	parigotFilename = "parigot-transformed.wat"

	tinygoNotImplementedImpl = "$github.com/iansmith/parigot/abi.TinyGoNotImplemented"
	jsNotImplementedImpl     = "$github.com/iansmith/parigot/abi.JSNotImplemented"

	jsNotImplImportName     = "JSNotImplemented"
	tinygoNotImplImportName = "TinyGoNotImplemented"
	parigotBaseFnName       = "$github.com/iansmith/parigot/abi/go/abi"
)

var importNameToImplName = map[string]string{
	tinygoNotImplImportName: tinygoNotImplementedImpl,
	jsNotImplImportName:     jsNotImplementedImpl,
}

var importNameToTypeNumber = map[string]int{
	tinygoNotImplImportName: 0,
	jsNotImplImportName:     0,
}

var jsErrorCall = &transform.CallOp{
	Arg: jsNotImplementedImpl,
}

type retValueType = int32

const (
	// xxx this should be extended so each byte of the
	// xxx int32 represents a specific return type for the
	// xxx case of multiple return values
	retNone retValueType = 1
	retI32  retValueType = 2
	retI64  retValueType = 3
	retF32  retValueType = 4
	retF64  retValueType = 5
)

func retValueTToString(r retValueType) string {
	switch r {
	case retNone:
		return "None"
	case retI32:
		return "i32"
	case retI64:
		return "i64"
	case retF32:
		return "f32"
	case retF64:
		return "f64"
	default:
		panic("unknown return value type" + fmt.Sprint(r))
	}
}

var typeNumberToReturnValueType = make(map[int]retValueType)
var replacementFuncTypeTable = make(map[string]int)
var importFuncTypeTable = make(map[string]int)

// for the trap version
var asExternalRetTypeToNewFunc = make(map[retValueType]int)
var maxTypeNum = -1

// command line args
var outputFile *string = flag.String("o", "", "set to the output file you want to produce, otherwise output goes to stdout")
var useTrapForParigot *bool = flag.Bool("trap", false, "set if you want generate kernel traps for all parigot_abi calls")

func main() {

	flag.Parse()

	if len(os.Args) == 1 {
		log.Fatalf("unable to understand arguments, providing an input WASM file is required")
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
	err = convertWatToWasm(tmp, modifiedWat, outputFile)
	if err != nil {
		os.Exit(1)
	}
	//os.RemoveAll(tmp)
	log.Printf("%s -> %s\n", os.Args[1], watVersion)
	log.Printf("%s -> %s\n", watVersion, filepath.Join(tmp, parigotFilename))
	log.Printf("%s -> %s\n", filepath.Join(tmp, parigotFilename), os.Args[2])
	os.Exit(0)
}
func transformation(mod *transform.Module) {
	findToplevelInModule(mod, transform.FuncDefT, mapReplacementFunctions)
	findToplevelInModule(mod, transform.TypeDefT, mapTypeReturnValues)
	addNotImplImports(mod)
	changeToplevelInModule(mod, transform.ImportDefT, modifyWasiImport)
	changeStatementInModule(mod, modifyWasiCall)
	findToplevelInModule(mod, transform.ImportDefT, modifyJSImports)
	changeToplevelInModule(mod, transform.FuncDefT, changeFuncsToNotImplemented)
	// this is only needed for causing the ABI to trap to the kernel
	if *useTrapForParigot {
		findToplevelInModule(mod, transform.ImportDefT, mapNeededNewTypesForExternalRef)
		addNewExternRefTypes(mod)
		changeToplevelInModule(mod, transform.ImportDefT, changeImportsToTrap)
	}

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
		addToplevelToModule(mod, imp, true)
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
	typeNumberToReturnValueType[td.Annotation] = 0
	if td.Func.Result != nil {
		if len(td.Func.Result.Type.Name) > 1 {
			panic("unable to handle mapping return types with multiple values")
		}
		rType := td.Func.Result.Type.Name[0]
		log.Printf("return type of type %d is %s", td.Annotation, rType)
		switch rType {
		case "i32":
			typeNumberToReturnValueType[td.Annotation] = 32
		case "i64":
			typeNumberToReturnValueType[td.Annotation] = 64
		default:
			panic("unable to handle mapping return types with type " + rType)
		}
	}
}
func modifyWasiCall(stmt transform.Stmt) transform.Stmt {
	if stmt.StmtType() == transform.OpStmtT && stmt.(transform.Op).OpType() == transform.CallT {
		call := stmt.(*transform.CallOp)
		if call.Arg == oldWasiCallTarget {
			call.Arg = parigotBaseFnName + ".FdWrite"
			log.Printf("adjusted call from runtime.fd_write to abi.FdWrite")
		}
	}
	return stmt
}

func modifyWasiImport(tl transform.TopLevel) transform.TopLevel {
	idef := tl.(*transform.ImportDef)
	if idef.ModuleName == WasiModule && idef.ImportedAs == FdWrite {
		idef.FuncNameRef.Name = parigotBaseFnName + ".FdWrite"
		log.Printf("flipped WASI interface to parigot emulation")
	}
	return idef
}

var returnStmt = &transform.ZeroOp{
	Op: "return",
}

func changeFuncsToNotImplemented(tl transform.TopLevel) transform.TopLevel {
	fn := tl.(*transform.FuncDef)
	if isJSFunction(fn.Name) {
		_, ok := replacementFuncTypeTable[fn.Name]
		log.Printf("replace: %s->%d->%d [%v]", fn.Name, replacementFuncTypeTable[fn.Name],
			typeNumberToReturnValueType[replacementFuncTypeTable[fn.Name]], ok)
		if typeNumberToReturnValueType[replacementFuncTypeTable[fn.Name]] != 0 {
			op := "i32.const"
			if typeNumberToReturnValueType[replacementFuncTypeTable[fn.Name]] == 64 {
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

func isParigotModule(name string) bool {
	return name == parigotModule
}

func mapReplacementFunctions(level transform.TopLevel) {
	fn := level.(*transform.FuncDef)
	if isJSFunction(fn.Name) {
		log.Printf("%s:type %d\n", fn.Name, fn.Type.Num)
		replacementFuncTypeTable[fn.Name] = fn.Type.Num
		if fn.Type.Num > maxTypeNum {
			maxTypeNum = fn.Type.Num
		}
	}
}

// this is used to populate typeNumberToReturnValueType
func mapNeededNewTypesForExternalRef(tl transform.TopLevel) {
	idef := tl.(*transform.ImportDef)
	if idef.ModuleName == parigotModule {
		// called for effect
		importedFunctionToExternalRefTypeNum(idef.FuncNameRef.Name)
	}
}

func importedFunctionToExternalRefTypeNum(name string) int {
	typeNum, ok := replacementFuncTypeTable[name]
	if !ok {
		panic("unable to find a type number for " + name)
	}
	retT := typeNumberToReturnValueType[typeNum]
	// check tosee if we already have this
	ourTypeNum, seenBefore := asExternalRetTypeToNewFunc[retT]
	if !seenBefore {
		// we need to create this function
		maxTypeNum++
		asExternalRetTypeToNewFunc[retT] = maxTypeNum
		ourTypeNum = maxTypeNum
	}
	return ourTypeNum
}

func changeImportsToTrap(tl transform.TopLevel) transform.TopLevel {
	idef := tl.(*transform.ImportDef)
	tModifiedTypeNum := importedFunctionToExternalRefTypeNum(idef.FuncNameRef.Name)
	if idef.ModuleName != parigotModule && idef.FuncNameRef.Type.Num != tModifiedTypeNum {
		log.Printf("updating for trap, changed %s from type %d to %d", idef.FuncNameRef.Name,
			idef.FuncNameRef.Type.Num, tModifiedTypeNum)
		idef.FuncNameRef.Type.Num = tModifiedTypeNum
	}
	return tl
}

func addNewExternRefTypes(mod *transform.Module) {
	for retT, typeNum := range asExternalRetTypeToNewFunc {
		tdef := &transform.TypeDef{
			Annotation: typeNum,
			Func: &transform.FuncSpec{
				Param: &transform.ParamDef{
					Type: &transform.TypeNameSeq{
						Name: []string{"externalRef"},
					}},
				Result: &transform.ResultDef{
					Type: &transform.TypeNameSeq{
						Name: []string{retValueTToString(retT)},
					},
				},
			},
		}
		addToplevelToModule(mod, tdef, true)
	}
}
