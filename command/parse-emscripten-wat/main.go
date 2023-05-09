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
type ImportDetail struct {
	Package string
	Name    string
	TypeNum int
	FuncNum int

	IsGlobal      bool
	GlobalNum     int
	GlobalMutable bool
	GlobalType    TypeName
}

var TypeInfo = make(map[int]*TypeDetail)
var ImportFuncInfo = make(map[int]*ImportDetail)
var ImportGlobalInfo = make(map[int]*ImportDetail)

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
	case "i64":
		return "int64"
	case f32:
		return "float32"
	case f64:
		return "float64"
	}
	panic("unknown type" + fmt.Sprintf("%T,'%s'", t, t))
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
				if imp.IsGlobal {
					ImportGlobalInfo[imp.GlobalNum] = imp
				} else {
					ImportFuncInfo[imp.FuncNum] = imp
				}
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

func processImport(line *Item) *ImportDetail {
	if line.List[1].Atom == nil || line.List[2].Atom == nil {
		panic("unexpect formatting of import")
	}
	if line.List[1].Atom.String == "" || line.List[2].Atom.String == "" {
		panic("import is missing the name of the import (package or func name)")
	}
	pkg := strings.TrimPrefix(strings.TrimSuffix(line.List[1].Atom.String, "\""), "\"")
	fn := strings.TrimPrefix(strings.TrimSuffix(line.List[2].Atom.String, "\""), "\"")

	if line.List[3].List[0].Atom.Symbol == "func" {
		return processImportDefFunc(line.List[3:], pkg, fn)
	}
	if line.List[3].List[0].Atom.Symbol == "global" {
		//log.Printf("n is %d", line.List[3].List[1].List[0].Atom.Number)
		n := line.List[3].List[1].List[0].Atom.Number
		if line.List[3].List[2].Atom != nil {
			log.Printf("about to process not mutable %s", line.List[3].List[2])
			return processImportDefGlobal(line.List[3].List[2].Atom.Symbol, false, n, pkg, fn)
		}
		m := line.List[3].List[2].List[0].Atom.Symbol == "mut"
		t := line.List[3].List[2].List[1].Atom.Symbol
		log.Printf("about to process %v,%s", m, t)
		return processImportDefGlobal(t, m, n, pkg, fn)
	}
	return nil
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
func processImportDefGlobal(globalT string, mut bool, globalNum int, pkg, name string) *ImportDetail {

	id := &ImportDetail{
		Package:       pkg,
		Name:          name,
		TypeNum:       0,
		FuncNum:       0,
		GlobalNum:     globalNum,
		GlobalMutable: mut,
		GlobalType:    ToTypeName(globalT),
		IsGlobal:      true,
	}
	return id

}
func processImportDefFunc(f []*Item, pkg, name string) *ImportDetail {

	fnNum := f[0].List[1].List[0].Atom.Number
	typeNum := f[0].List[2].List[1].Atom.Number
	// log.Printf("ideffunc: %s, -- %d", f[0].List[2].List[0], f[0].List[2].List[1].Atom.Number)
	// log.Printf("processImportDefFunc %s.%s=>%d", pkg, name, typeNum)
	return &ImportDetail{
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

	wasmtime "github.com/bytecodealliance/wasmtime-go/v8"
)

$1

func addEmscriptenFuncs(store wasmtime.Storelike, result map[string]*wasmtime.Func, rt *Runtime) {
$2
}

func addEmscriptenGlobals(store wasmtime.Storelike, result map[string]*wasmtime.Global) {
	var valType *wasmtime.ValType
	var gType *wasmtime.GlobalType
	var g *wasmtime.Global
	var err error

$3
}
	
`

func generateFile() {
	decl := &bytes.Buffer{}
	for _, f := range ImportFuncInfo {
		if f.IsGlobal {
			continue
		}
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
		if tInfo.Result != "" {
			//log.Printf("f is %s.%s and tInfo.Result is %s", f.Package, f.Name, tInfo.Result)
			decl.WriteString("\treturn 0\n")
		}
		decl.WriteString("}\n\n")
	}

	wrap := &bytes.Buffer{}

	for _, f := range ImportFuncInfo {
		if f.IsGlobal {
			continue
		}
		wrap.WriteString(fmt.Sprintf("\tresult[\"%s.%s\"] = wasmtime.WrapFunc(store,%s_%s)\n", f.Package, f.Name, f.Package, f.Name))
	}

	global := &bytes.Buffer{}

	for _, f := range ImportGlobalInfo {
		if !f.IsGlobal {
			continue
		}
		typeSuffix := "I32"
		switch f.GlobalType {
		case i32:
			break
		case i64:
			typeSuffix = "I64"
		case f64:
			typeSuffix = "F64"
		case f32:
			typeSuffix = "F32"
		}
		mut := "false"
		if f.GlobalMutable {
			mut = "true"
		}
		typeWithSuffix := "wasmtime.Kind" + typeSuffix
		valWithSuffix := "wasmtime.Val" + typeSuffix
		global.WriteString(fmt.Sprintf("\tvalType=wasmtime.NewValType(%s)\n", typeWithSuffix))
		global.WriteString(fmt.Sprintf("\tgType=wasmtime.NewGlobalType(valType,%s)\n", mut))
		global.WriteString(fmt.Sprintf("\tg,err=wasmtime.NewGlobal(store,gType,%s(0))\n", valWithSuffix))
		global.WriteString("\tif err!=nil {\n")
		global.WriteString("\t\tpanic(err.Error())\n")
		global.WriteString("\t}\n")
		global.WriteString(fmt.Sprintf("\tresult[\"%s.%s\"]=g\n", f.Package, f.Name))
	}

	x := strings.Replace(file, "$1", decl.String(), 1)
	x = strings.Replace(x, "$2", wrap.String(), 1)
	x = strings.Replace(x, "$3", global.String(), 1)
	fmt.Print(x)
}
