package main

import (
	"bufio"
	"github.com/iansmith/parigot/command/transform"
	"log"
	"os"
	"strings"
)

type replaceFn struct {
	pair []*replacementPair
	orig map[string]*replacementPair
	repl map[string]*replacementPair
}

type replacementPair struct {
	name1, name2 string
	fn1, fn2     *transform.FuncDef
}

func newReplaceFn(funcName1, funcName2 string, fp *os.File) *replaceFn {
	pair := []*replacementPair{}
	if fp != nil {
		// read the file
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			t := scanner.Text()
			parts := strings.Split(t, ",")
			if len(parts) != 2 {
				if !strings.HasPrefix(t, "#") {
					log.Printf("ignoring line '%s' read from file", t)
				}
				continue
			}
			parts[0] = strings.TrimSpace(parts[0])
			parts[1] = strings.TrimSpace(parts[1])
			rp := &replacementPair{name1: parts[0], name2: parts[1]}
			pair = append(pair, rp)
		}
		if scanner.Err() != nil {
			log.Fatalf("scanner failed reading replacement function file: %v", scanner.Err())
		}
	} else {
		rp := &replacementPair{name1: funcName1, name2: funcName2}
		pair = []*replacementPair{rp}
	}
	orig := make(map[string]*replacementPair)
	repl := make(map[string]*replacementPair)
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
			log.Printf("unable to locate replacement function '%s'", rp.name2)
			ok = false
		}
		// check the types
		if rp.fn1.Type.Num != rp.fn2.Type.Num {
			log.Fatalf("cannot replace function '%s' with '%s' because they have different type signatures",
				rp.name1, rp.name2)
		}
	}
	if !ok {
		log.Fatalf("aborting, did not find all functions with matching types")
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
	rp, ok := r.orig[call.Arg]
	if !ok {
		return []transform.Stmt{stmt}
	}
	log.Printf("'%s' -> '%s'", rp.name1, rp.name2)
	return []transform.Stmt{
		&transform.CallOp{
			Arg: *rp.fn2.Name,
		},
	}
}
