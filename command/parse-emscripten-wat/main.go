package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"

	"github.com/iansmith/parigot/helper/antlr"

	v4 "github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type FuncDetail struct {
	Param  []TypeName
	Result TypeName
}

var FuncInfo = make(map[int]*FuncDetail)

type TypeName string

const (
	i32 TypeName = "i32"
	i64          = "i34"
	f32          = "f32"
	f64          = "f64"
)

func (t TypeName) String() string {
	return string(t)
}
func ToTypeName(s string) TypeName {
	switch s {
	case "i32", "i64", "f32", "f64":
		return TypeName(s)
	}
	panic("unknown type name " + s)
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("you should supply exactly one file, containing wat, to this program")
	}
	file := flag.Arg(0)
	log.Printf("file is %s", file)
	lexer := NewsexprLexer(nil)
	parser := NewsexprParser(nil)
	build := NewSexprBuilder(file)
	el := antlr.AntlrSetupLexParse(file, lexer.BaseLexer, parser.BaseParser)
	item := parser.Sexpr()
	v4.ParseTreeWalkerDefault.Walk(build, item)
	if len(item.GetItem_()) != 1 {
		log.Fatalf("unable to understand sexpr that has wrong number of items %d", len(item.GetItem_()))
	}
	if el.Failed() {
		antlr.AntlrFatalf("failed due to syntax errors trying to parse sexpression")
	}
	//log.Printf("xxx %s", item.GetItem_()[0].String())
	list := item.GetItem_()[0].List[1:]

	log.Printf("inner list len %d", len(list))
	for _, stmt := range list {
		if isListWithFirstAtom(stmt, "type") {
			processType(stmt)
		}
	}

	log.Printf("finished")

}

func isListWithFirstAtom(curr *Item, s string) bool {
	if curr.List == nil {
		return false
	}
	l := curr.List
	if len(l) == 0 {
		return false
	}
	first := l[0]
	if first.Atom == nil || first.Atom.Symbol == "" {
		return false
	}
	if first.Atom.Symbol != s {
		return false
	}
	return true
}

func processType(curr *Item) {

	l := curr.List
	num := l[1].List[0].Atom.Number
	param, result := processFunc(l[2].List)
	FuncInfo[num] = &FuncDetail{Param: param, Result: result}
	log.Printf("#%d(%+v)", num, FuncInfo[num])
}

func processImport(line *Item) {

}

func (i *Item) String() string {
	if i == nil {
		panic("nil item, trying to convert to string")
	}
	var buf bytes.Buffer
	i.StringToBuffer(&buf, 0, false)
	return buf.String()
}

func (i *Item) StringToBuffer(b *bytes.Buffer, indent int, useCR bool) {
	if i == nil {
		panic("nil item, trying to convert to string")
	}
	if i.Atom != nil {
		if i.Atom.String != "" {
			b.WriteString(i.Atom.String)
			return
		}
		if i.Atom.Symbol != "" {
			b.WriteString(i.Atom.Symbol)
			return
		}
		s := fmt.Sprintf("%d", i.Atom.Number)
		if !i.Atom.CommentNum {
			b.WriteString(s)
			return
		}
		b.WriteString(fmt.Sprintf(";%d;", i.Atom.Number))
		return
	}
	if i.List == nil {
		panic("neither list nor atom found in item!")
	}

	if useCR {
		b.WriteString("\n")
		for in := 0; in < indent; in++ {
			b.WriteString(" ")
		}
	}
	b.WriteString("(")
	for j, elem := range i.List {
		elem.StringToBuffer(b, indent+2, useCR)
		if j != len(i.List)-1 {
			b.WriteString(" ")
		}
	}
	b.WriteString(")")

	if useCR {
		b.WriteString("\n")
		for in := 0; in < indent; in++ {
			b.WriteString(" ")
		}

	}
}

func processFunc(f []*Item) ([]TypeName, TypeName) {
	param := []TypeName{}
	var result TypeName

	if !isListWithFirstAtom(&Item{List: f}, "func") {
		panic("sent a non func to be processed")
	}

	// this is the case of no paramaters and no result (call for effect)
	if len(f) == 1 {
		return param, result
	}
	second := f[1]
	funcSpec := second.List[1:]
	if second.List[0].Atom.Symbol == "param" {
		for i := 0; i < len(funcSpec); i++ {
			if funcSpec[i].Atom == nil {
				panic("unable to understand func definition")
			}
			if funcSpec[i].Atom.Symbol == "" {
				panic("unable to understand atom in param list")
			}
			param = append(param, ToTypeName(funcSpec[i].Atom.Symbol))
		}
		if len(f) == 3 {
			third := f[2]
			if third.List == nil {
				panic("unable to understand return value")
			}
			result = ToTypeName(third.List[1].Atom.Symbol)
		}
	} else {
		if second.List == nil {
			panic("unable to understand return value")
		}
		result = ToTypeName(second.List[1].Atom.Symbol)
	}

	return param, result
}
