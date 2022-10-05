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
	FdWriteEmulName   = "FdWrite"
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

// command line args
var outputFile *string = flag.String("o", "", "set to the output file you want to produce, otherwise output goes to stdout")
var useTrapForParigot *bool = flag.Bool("trap", false, "set if you want generate kernel traps for all parigot_abi calls")

func main() {

	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalf("unable to understand arguments, providing an input WASM file is required")
	}

	_, baseArg0 := filepath.Split("jsstrip")
	tmp, err := os.MkdirTemp("", baseArg0)
	if err != nil {
		log.Fatalf("cannot create temp dir: %v", err)
	}

	watVersion, err := convertWasmToWat(tmp, flag.Arg(0))
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

	//// compute the number of type defs, because we may need to compute _new_ type defs
	//findToplevelInModule(mod, transform.TypeDefT, countTypeDefs)
	//
	//// this function looks for every js function (in the name) and notes what it's type num is
	//// it is building up replacementFuncTypeTable to hap the map we will need later
	//findToplevelInModule(mod, transform.FuncDefT, mapReplacementFunctions)
	//
	//// this function walks ever type number and stores its return value as a retTypeT
	//// in the table typeNumberToReturnValueType
	//findToplevelInModule(mod, transform.TypeDefT, mapTypeReturnValues)
	//
	//// we need to do the same thing for js imports that we did for js function definitons
	//// above... we are building up the data in importFuncTypeTable
	//findToplevelInModule(mod, transform.ImportDefT, mapJSImports)
	//
	//// we have to add two imports, one each for JS and TinyGo to say "not implemented"
	//// these are added via import and we don't need to change the set of types known,
	//// because these are clearly #0
	//addNotImplImports(mod)
	//
	//// walk the imports looking for a call into wasi and substituting our emulation code
	//// this modifies the imports it changes in place
	//changeToplevelInModule(mod, transform.ImportDefT, modifyWasiImport)
	//
	//// now we need to find every call to the wasi code and substitute our new function
	//// which is the same as we did with the import
	//changeStatementInModule(mod, modifyWasiCall)
	//
	//// now the big one, we walk all the function defs and rewrite the code to
	//// actually call our jsNotImplemented function
	//changeToplevelInModule(mod, transform.FuncDefT, changeFuncsToNotImplemented)
	//
	//xxxtemp hack to test
	changeToplevelInModule(mod, transform.ImportDefT, changeImportsToPointToConversions)
	changeStatementInModule(mod, changeCallsToUseConversion)

	// this is only needed for causing the ABI to trap to the kernel
	//if *useTrapForParigot {
	//	findToplevelInModule(mod, transform.ImportDefT, mapAllImports)
	//	findToplevelInModule(mod, transform.ImportDefT, mapNeededNewTypesForExternalRef)
	//	addNewExternRefTypes(mod)
	//	changeToplevelInModule(mod, transform.ImportDefT, changeImportsToTrap)
	//}

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

func mapJSImports(tl transform.TopLevel) {
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
			typeNumberToReturnValueType[td.Annotation] = retI32
		case "i64":
			typeNumberToReturnValueType[td.Annotation] = retI64
		case "f32":
			typeNumberToReturnValueType[td.Annotation] = retF32
		case "f64":
			typeNumberToReturnValueType[td.Annotation] = retF64
		default:
			panic("unable to handle mapping return types with type " + rType)
		}
	} else {
		typeNumberToReturnValueType[td.Annotation] = retNone
	}
}

var returnStmt = &transform.ZeroOp{
	Op: "return",
}

func isJSFunction(name string) bool {
	return strings.Index(name, "/js") != -1
}

func isParigotModule(name string) bool {
	return name == parigotModule
}

func importedFunctionToExternalRefTypeNum(name string) int {
	typeNum, ok := importedFuncNameToTypeNum[name]
	if !ok {
		panic("unable to find a type number for " + name)
	}
	retT := typeNumberToReturnValueType[typeNum]
	// check to see if we already have this
	ourTypeNum, seenBefore := asExternalRetTypeToNewFunc[retT]
	if !seenBefore {
		// we need to create this function
		maxTypeNum++
		asExternalRetTypeToNewFunc[retT] = maxTypeNum
		ourTypeNum = maxTypeNum
	}
	return ourTypeNum
}

func changeImportsToPointToConversions(tl transform.TopLevel) transform.TopLevel {
	idef := tl.(*transform.ImportDef)
	if idef.FuncNameRef.Name == "$github.com/iansmith/parigot/abi/go/abi.OutputString" {
		idef.FuncNameRef.Name = "$github.com/iansmith/parigot/abi/go/abi.OutputStringConvert"
		idef.FuncNameRef.Type.Num = 2
	}
	return idef
}
func changeCallsToUseConversion(stmt transform.Stmt) transform.Stmt {
	if stmt.StmtType() == transform.OpStmtT &&
		stmt.(transform.Op).OpType() == transform.CallT &&
		stmt.(*transform.CallOp).Arg == "$github.com/iansmith/parigot/abi/go/abi.OutputString" {
		log.Printf("adjusting call... %#v", stmt)
		stmt.(*transform.CallOp).Arg = "$github.com/iansmith/parigot/abi/go/abi.OutputStringConvert"
	}
	return stmt
}
