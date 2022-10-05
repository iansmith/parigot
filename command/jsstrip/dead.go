package main

import (
	"log"
	"sort"
	"strings"

	"github.com/iansmith/parigot/command/transform"
)

var tgErrorCall = &transform.CallOp{
	Arg: tinygoNotImplementedImpl,
}

func modifyParigotAbiFuncNames(tl transform.TopLevel) transform.TopLevel {
	fn := tl.(*transform.ImportDef)
	if isParigotModule(fn.ModuleName) {
		log.Printf("parigot function imported: (%s:%s) %s:%d\n", fn.ModuleName, fn.ImportedAs, fn.FuncNameRef.Name, fn.FuncNameRef.Type.Num)
		if strings.HasPrefix(fn.FuncNameRef.Name, "$") {
			fn.FuncNameRef.Name = fn.FuncNameRef.Name[1:]
			log.Printf("    --> removed $ to make it %s", fn.FuncNameRef.Name)
		}
	}
	return fn
}

func mapAllImports(t transform.TopLevel) {
	idef := t.(*transform.ImportDef)
	name := idef.FuncNameRef.Name
	num := idef.FuncNameRef.Type.Num
	importedFuncNameToTypeNum[name] = num
}

// this is used to populate typeNumberToReturnValueType
func mapNeededNewTypesForExternalRef(tl transform.TopLevel) {
	idef := tl.(*transform.ImportDef)
	// called for effect
	importedFunctionToExternalRefTypeNum(idef.FuncNameRef.Name)
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
	allNewTypes := []int{}
	reverse := make(map[int]retValueType)
	for retT, typeNum := range asExternalRetTypeToNewFunc {
		allNewTypes = append(allNewTypes, typeNum)
		reverse[typeNum] = retT
	}
	sort.Ints(allNewTypes)
	for _, typeNum := range allNewTypes {
		retT := reverse[typeNum]
		tdef := &transform.TypeDef{
			Annotation: typeNum,
			Func: &transform.FuncSpec{
				Param: &transform.ParamDef{
					Type: &transform.TypeNameSeq{
						Name: []string{"externref"},
					}},
				Result: &transform.ResultDef{
					Type: &transform.TypeNameSeq{
						Name: []string{retValueTToString(retT)},
					},
				},
			},
		}
		// this is easier than trying to construct a diff tree
		if retT == retNone {
			//chop it off
			tdef.Func.Result = nil
		}
		addToplevelToModule(mod, tdef, true)
	}
}

func changeWasiImportToEmulation(tl transform.TopLevel) transform.TopLevel {
	imp := tl.(*transform.ImportDef)
	if imp.ModuleName == WasiModule && imp.ImportedAs == FdWrite {
		imp.ModuleName = "parigot_abi"
		//		imp.ImportedAs = parigotMangleMethod("WasiEmulation", "FdWrite")
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
