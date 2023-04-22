package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/iansmith/parigot/helper/antlr"

	v4 "github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type TypeDetail struct {
	Param  []TypeName
	Result TypeName
}
type FuncDetail struct {
	Package string
	Name    string
	TypeNum int
	FuncNum int
}

var TypeInfo = make(map[int]*TypeDetail)
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
func (t TypeName) GoType() string {
	switch t {
	case i32:
		return "int32"
	case i64:
		return "int64"
	case f32:
		return "float32"
	case f64:
		return "float64"
	}
	panic("unknown type")
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
	list := item.GetItem_()[0].List[1:]

	for _, stmt := range list {
		if isListWithFirstAtom(stmt, "type") {
			processType(stmt)
		}
		if isListWithFirstAtom(stmt, "import") {
			imp := processImport(stmt)
			if imp != nil {
				FuncInfo[imp.FuncNum] = imp
			}
		}
	}

	generateFile()
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
	param, result := processFuncTypeParamsReturn(l[2].List)
	TypeInfo[num] = &TypeDetail{Param: param, Result: result}
}

func processImport(line *Item) *FuncDetail {
	if line.List[1].Atom == nil || line.List[2].Atom == nil {
		panic("unexpect formatting of import")
	}
	if line.List[1].Atom.String == "" || line.List[2].Atom.String == "" {
		panic("import is missing the name of the import (package or func name)")
	}
	pkg := strings.TrimPrefix(strings.TrimSuffix(line.List[1].Atom.String, "\""), "\"")
	fn := strings.TrimPrefix(strings.TrimSuffix(line.List[2].Atom.String, "\""), "\"")

	if line.List[3].List[0].Atom.Symbol != "func" {
		//log.Printf("ignorning import: %s", line.List[3].List[0])
		return nil
	}
	return processImportDef(line.List[3:], pkg, fn)
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
func processImportDef(f []*Item, pkg, name string) *FuncDetail {
	if name == "__memory_base" {
		return nil
	}

	fnNum := f[0].List[1].List[0].Atom.Number
	typeNum := f[0].List[2].List[0].Atom.Number
	return &FuncDetail{
		Package: pkg,
		Name:    name,
		TypeNum: typeNum,
		FuncNum: fnNum,
	}
}

func processFuncTypeParamsReturn(f []*Item) ([]TypeName, TypeName) {
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

var file = `
package sys

import (
	"log"

	wasmtime "github.com/bytecodealliance/wasmtime-go/v7"
)

$1

func addEmscriptenFuncs(store wasmtime.Storelike, result map[string]*wasmtime.Func, rt *Runtime) {
$2
}
`

func generateFile() {
	decl := &bytes.Buffer{}
	for _, f := range FuncInfo {
		decl.WriteString(fmt.Sprintf("// %s.%s => wasm function #%d\n", f.Package, f.Name, f.FuncNum))
		decl.WriteString(fmt.Sprintf("func %s_%s(", f.Package, f.Name))
		tInfo := TypeInfo[f.TypeNum]
		for i, p := range tInfo.Param {
			decl.WriteString(fmt.Sprintf("p%d %s", i, p.GoType()))
			if i != len(tInfo.Param)-1 {
				decl.WriteString(",")
			}
		}
		decl.WriteString(")")
		if tInfo.Result.String() != "" {
			decl.WriteString(" " + tInfo.Result.GoType())
		}
		decl.WriteString("{\n")
		decl.WriteString(fmt.Sprintf("\tlog.Printf(\"call to --> %s_%s:", f.Package, f.Name))
		for i, p := range tInfo.Param {
			if p.String() == "i32" || p.String() == "i64" {
				decl.WriteString("0x%x")
			} else {
				decl.WriteString("%4.4f")
			}
			if i != len(tInfo.Param)-1 {
				decl.WriteString(",")
			}
		}
		decl.WriteString("\",")

		for i := range tInfo.Param {
			decl.WriteString(fmt.Sprintf("p%d", i))
			if i != len(tInfo.Param)-1 {
				decl.WriteString(",")
			}
		}
		decl.WriteString(")\n")
		decl.WriteString("\treturn 0\n")
		decl.WriteString("}\n\n")
	}

	wrap := &bytes.Buffer{}

	for _, f := range FuncInfo {
		wrap.WriteString(fmt.Sprintf("\tresult[\"%s.%s\"] = wasmtime.WrapFunc(store,%s_%s)\n", f.Package, f.Name, f.Package, f.Name))
	}

	x := strings.Replace(file, "$1", decl.String(), 1)
	x = strings.Replace(x, "$2", wrap.String(), 1)
	fmt.Print(x)
}
