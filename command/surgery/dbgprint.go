package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/iansmith/parigot/command/transform"
)

type dbgPrint struct {
	funcName string
	funcDef  *transform.FuncDef
	typeNum  int
}

func newDbgPrint(fn string) *dbgPrint {
	return &dbgPrint{
		funcName: fn,
		typeNum:  -2712,
	}
}

func (d *dbgPrint) hasFunction() bool {
	return d.funcDef != nil
}

func (d *dbgPrint) findFunction(tl transform.TopLevel) {
	fd := tl.(*transform.FuncDef)
	if fd.Name != nil && *fd.Name == d.funcName {
		d.funcDef = fd
	}
}

func (d *dbgPrint) listFunction() {
	for i, stmt := range d.funcDef.Code {
		log.Printf("%04d: %s", i, stmt.IndentedString(0))
	}
}

func (d *dbgPrint) hasFuncType() bool {
	return d.typeNum >= 0
}
func (d *dbgPrint) findFuncTypeForDbgPrint(tl transform.TopLevel) {
	td := tl.(*transform.TypeDef)
	if td.Func.Result != nil {
		return
	}
	if td.Func.Param == nil {
		return
	}
	p := td.Func.Param
	if len(p.Type.Name) != 1 {
		return
	}
	if p.Type.Name[0] != "i32" {
		return
	}
	d.typeNum = td.Annotation
}

func (d *dbgPrint) updateStatementsInFunc() {
	var outBuffer bytes.Buffer
	code := d.funcDef.Code
	count := 0
	d.funcDef.Code = changeStmtCodeOnly(code, 0, func(stmt transform.Stmt, _ int) []transform.Stmt {
		result := append([]transform.Stmt{stmt}, myStmts(count)...)
		outBuffer.WriteString(fmt.Sprintf("%05d", count, stmt.IndentedString(0)))
		count++
		return result
	})
	s := ""
	if d.funcDef.Name != nil {
		s = *d.funcDef.Name
	} else {
		if d.funcDef.Number == nil {
			panic("function has neither name nor number")
		}
		s = fmt.Sprint(*d.funcDef.Number)
	}
	log.Printf("instrumented %d statements in %s", count, s)
	fp, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("unable to create %s: %v", *outputFile, err)
	}
	_, err = io.Copy(fp, &outBuffer)
	if err != nil {
		log.Fatalf("unable to copy the output data to %s: %v", *outputFile, err)
	}
	fp.Close()
}

func myStmts(count int) []transform.Stmt {
	constStmt := &transform.ArgOp{
		Op:         "i32.const",
		IntArg:     new(int),
		FloatArg:   nil,
		BranchAnno: nil,
		ConstAnno:  nil,
		Special:    nil,
	}
	*constStmt.IntArg = count
	return []transform.Stmt{
		constStmt,
		&transform.CallOp{
			Arg: "$parigot.dbgprint",
		},
	}
}

func (dbg *dbgPrint) updateImports() {

}

func (dbg *dbgPrint) driver(mod *transform.Module) {
	findToplevelInModule(mod, transform.FuncDefT, dbg.findFunction)
	if !dbg.hasFunction() {
		log.Fatalf("unable to find function '%s' for dbgPrint operation", *fnName)
	}
	findToplevelInModule(mod, transform.TypeDefT, dbg.findFuncTypeForDbgPrint)
	if !dbg.hasFuncType() {
		log.Fatalf("unable to find a function type for our dbgprint call")
	}
	if *outputFile == "" {
		log.Fatalf("an outputfile needs to given with -o option (this produces a lot of output)")
	}
	dbg.updateStatementsInFunc()
	ourImport := &transform.ImportDef{
		ModuleName: "parigot",
		ImportedAs: "debugprint",
		FuncNameRef: &transform.FuncNameRef{
			Name: new(string),
			Type: &transform.TypeRef{
				Num: dbg.typeNum,
			},
		},
	}
	*ourImport.FuncNameRef.Name = "$parigot.dbgprint"
	addToplevelToModule(mod, ourImport, true)
}
