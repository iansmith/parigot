package main

import (
	"bufio"
	"github.com/iansmith/parigot/command/transform"
	"log"
	"os"
	"strings"
)

type replaceFn struct {
	pair         []*replacement
	orig         map[string]*replacement
	repl         map[string]*replacement
	nextFnNumber int // this is the number of the NEXT function added
}

type replacement struct {
	name1, name2          string
	importpkg, importname string
	fn1, fn2              *transform.FuncDef
	importFn              *transform.ImportDef
	newFuncNumber         int
}

func newReplaceFn(funcName1, funcName2 string, fp *os.File) *replaceFn {
	pair := []*replacement{}
	if fp != nil {
		// read the file
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			t := scanner.Text()
			parts := strings.Split(t, ",")
			if len(parts) != 2 {
				if !strings.HasPrefix(t, "#") && t != "" {
					log.Printf("\tignoring line '%s' read from file", t)
				}
				continue
			}
			parts[0] = strings.TrimSpace(parts[0])
			parts[1] = strings.TrimSpace(parts[1])
			rp := &replacement{name1: parts[0], name2: parts[1]}
			pair = append(pair, rp)
		}
		if scanner.Err() != nil {
			log.Fatalf("scanner failed reading replacement function file: %v", scanner.Err())
		}
	} else {
		rp := &replacement{name1: funcName1, name2: funcName2}
		pair = []*replacement{rp}
	}
	if len(pair) > 1 {
		log.Printf("\tsurgery ready for %d replacements", len(pair))
	}
	orig := make(map[string]*replacement)
	repl := make(map[string]*replacement)
	for _, rp := range pair {
		orig[rp.name1] = rp
		repl[rp.name2] = rp
	}
	return &replaceFn{
		pair: pair,
		orig: orig,
		repl: repl,
	}
}

func (r *replaceFn) validateFunctions(mod *transform.Module) {
	count := 0
	findToplevelInModule(mod, transform.ImportDefT, func(tl transform.TopLevel) {
		count++
	})
	findToplevelInModule(mod, transform.FuncDefT, func(tl transform.TopLevel) {
		count++
	})
	log.Printf("xxx number of functions counted %d", count)
	r.nextFnNumber = count
	findToplevelInModule(mod, transform.FuncDefT, func(tl transform.TopLevel) {
		fdef := tl.(*transform.FuncDef)
		if fdef.Name == nil {
			return
		}
		rp, ok := r.orig[*fdef.Name]
		if ok {
			if rp.fn1 != nil {
				panic("two hits on the same function definition of " + *fdef.Name)
			}
			rp.fn1 = fdef
			return
		}
		rp, ok = r.repl[*fdef.Name]
		if ok {
			if rp.fn2 != nil {
				panic("two hits on the same function definition of " + *fdef.Name)
			}
			rp.fn2 = fdef
			return
		}
	})
	ok := true
	for _, rp := range r.pair {
		if rp.fn1 == nil {
			log.Printf("unable to locate original function '%s'", rp.name1)
			ok = false
		}
		if rp.fn2 == nil {
			r.createImportForName(rp)
			r.createImportDefForName(rp)
		}
		// check the types
		if rp.fn2 != nil && rp.fn1.Type.Num != rp.fn2.Type.Num {
			log.Printf("cannot replace function '%s' with '%s' because they have different type signatures",
				rp.name1, rp.name2)
			ok = false
		}
	}
	if !ok {
		log.Fatalf("aborting, did not find all functions with matching types")
	}
}
func (r *replaceFn) addImports(module *transform.Module) {
	for _, rp := range r.pair {
		if rp.importFn != nil {
			module.AppendTopLevelDef(rp.importFn)
		}
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
	rp, ok := r.orig[*call.ArgName]
	if !ok {
		return []transform.Stmt{stmt}
	}
	log.Printf("'%s' -> '%s'", rp.name1, rp.name2)
	if rp.importFn != nil {
		op := &transform.CallOp{
			ArgNum: new(int),
		}
		*op.ArgNum = rp.newFuncNumber
		return []transform.Stmt{
			op,
		}
	}
	op := &transform.CallOp{
		ArgName: new(string),
	}
	*op.ArgName = *rp.fn2.Name
	return []transform.Stmt{
		op,
	}
}

func (r *replaceFn) createImportForName(rp *replacement) {
	parts := strings.Split(rp.name2, ".")
	if len(parts) != 2 {
		rp.importpkg = "env"
		rp.importname = rp.name2
	} else {
		rp.importpkg = parts[0]
		rp.importname = parts[1]
	}
	rp.name2 = ""
}

func (r *replaceFn) createImportDefForName(rp *replacement) {
	fnNumber := new(int)
	*fnNumber = r.nextFnNumber
	rp.newFuncNumber = *fnNumber
	r.nextFnNumber++
	idef := &transform.ImportDef{
		ModuleName: rp.importpkg,
		ImportedAs: rp.importname,
		FuncNameRef: &transform.FuncNameRef{
			Name:   nil,
			Number: fnNumber,
			Type: &transform.TypeRef{
				Num: rp.fn1.Type.Num,
			},
		},
	}
	rp.importFn = idef
}
