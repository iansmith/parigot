package main

import (
	"fmt"
	"github.com/iansmith/parigot/command/transform"
	"log"
)

type replaceFn struct {
	funcName1 string
	funcName2 string
	fn1       *transform.FuncDef
	fn2       *transform.FuncDef
}

func newReplaceFn(funcName1, funcName2 string) *replaceFn {
	return &replaceFn{
		funcName1: funcName1,
		funcName2: funcName2,
	}
}

func (r *replaceFn) validateFunctions(mod *transform.Module) {

	findToplevelInModule(mod, transform.FuncDefT, func(tl transform.TopLevel) {
		fdef := tl.(*transform.FuncDef)
		if fdef.Name != nil && *fdef.Name == r.funcName1 {
			if r.fn1 != nil {
				panic(fmt.Sprintf("found two different functions with the same name: %s", fdef.Name))
			}
			r.fn1 = fdef
		}
		if fdef.Name != nil && *fdef.Name == r.funcName2 {
			if r.fn2 != nil {
				panic(fmt.Sprintf("found two different functions with the same name:%s", fdef.Name))
			}
			r.fn2 = fdef
		}
	})
	errStr := ""
	if r.fn1 == nil {
		errStr += fmt.Sprintf("unable to find function named '%s'", r.funcName1)
	}
	if r.fn2 == nil {
		errStr += fmt.Sprintf(" unable to find function named '%s'", r.funcName2)
	}
	if errStr != "" {
		log.Fatalf("relpacefn failed: %s", errStr)
	}
	// check the types
	if r.fn1.Type.Num != r.fn2.Type.Num {
		log.Fatalf("replacefn failed, the two functions have different signatures")
	}
}

func (r *replaceFn) replace(stmt transform.Stmt, _ int) []transform.Stmt {
	if stmt.StmtType() != transform.OpStmtT {
		return []transform.Stmt{stmt}
	}
	if stmt.(transform.Op).OpType() != transform.CallT {
		return []transform.Stmt{stmt}
	}
	call := stmt.(*transform.CallOp)
	if call.Arg != r.funcName1 {
		return []transform.Stmt{stmt}
	}
	log.Printf("found and replaced function in : %s", call.IndentedString(0))
	return []transform.Stmt{
		&transform.CallOp{
			Arg: r.funcName2,
		},
	}
}
