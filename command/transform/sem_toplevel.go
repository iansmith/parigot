package transform

import (
	"fmt"
)

type TopLevelT int

const (
	TypeDefT   TopLevelT = 1
	ImportDefT TopLevelT = 2
	FuncDefT   TopLevelT = 3
	TableDefT  TopLevelT = 4
	MemoryDefT TopLevelT = 5
	GlobalDefT TopLevelT = 6
	ExportDefT TopLevelT = 7
	ElemDefT   TopLevelT = 8
	DataDefT   TopLevelT = 9
)

// TopLevel instances are the decls at the top level of a module
type TopLevel interface {
	IndentedStringer
	TopLevelType() TopLevelT
}

// TypeDef represents WAT like this:   (type (;9;) (func (param i64 i32 i32 i32) (result i64)))
type TypeDef struct {
	Annotation int
	Func       *FuncSpec
}

func (t *TypeDef) TopLevelType() TopLevelT {
	return TypeDefT
}

func (t *TypeDef) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString(fmt.Sprintf("(type (;%d;) %s)", t.Annotation,
		t.Func.String()))
	return buf.String()
}

// ImportDef represents WAT like this: (import "wasi_snapshot_preview1" "fd_write" (func $runtime.fd_write (type 4)))
type ImportDef struct {
	ModuleName  string
	ImportedAs  string
	FuncNameRef *FuncNameRef
}

func (i *ImportDef) TopLevelType() TopLevelT {
	return ImportDefT
}

func (i *ImportDef) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString(fmt.Sprintf("(import \"%s\" \"%s\"",
		i.ModuleName, i.ImportedAs))
	if i.FuncNameRef != nil {
		buf.WriteString(" " + i.FuncNameRef.String())
	}
	buf.WriteString(")")
	return buf.String()

}

// FuncDef represents WAT like this:   (func $_*internal/task.gcData_.swap (type 1) (param i32) (local i32) ...statements )
type FuncDef struct {
	Name   *string
	Number *int
	Type   *TypeRef
	Param  *ParamDef
	Local  *LocalDef
	Result *ResultDef
	Code   []Stmt
}

func (f *FuncDef) TopLevelType() TopLevelT {
	return FuncDefT
}

func (f *FuncDef) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	s := ""
	if f.Name != nil {
		s = *f.Name
	} else {
		if f.Number == nil {
			panic("function has neither name nor number")
		}
		s = fmt.Sprint(*f.Number)
	}
	buf.WriteString(fmt.Sprintf("(func %s",
		s))
	if f.Type != nil {
		buf.WriteString(" " + f.Type.String())
	}
	if f.Param != nil {
		buf.WriteString(" " + f.Param.String())
	}
	if f.Result != nil {
		buf.WriteString(" " + f.Result.String())
	}
	if f.Local != nil {
		buf.WriteString("\n")
		for i := 0; i < indented+2; i++ {
			buf.WriteString(" ")
		}
		buf.WriteString(f.Local.String())
	}
	buf.WriteString(stmtsToString(f.Code, indented+2))
	// strangely this is not terminated with a "\n"
	return buf.String() + ")"
}

func (f *FuncDef) AddStmt(s Stmt) {
	f.Code = append(f.Code, s)
}

type TableDef struct {
	Type     *int
	Min, Max int
}

func (t *TableDef) TopLevelType() TopLevelT {
	return TableDefT
}
func (t *TableDef) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString("(table")
	if t.Type != nil {
		buf.WriteString(fmt.Sprintf(" (;%d;)", *t.Type))
	}
	buf.WriteString(fmt.Sprintf(" %d %d funcref)", t.Min, t.Max))
	return buf.String()
}

type MemoryDef struct {
	Type *int
	Size int
}

func (m *MemoryDef) TopLevelType() TopLevelT {
	return MemoryDefT
}
func (m *MemoryDef) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString("(memory")
	if m.Type != nil {
		buf.WriteString(fmt.Sprintf(" (;%d;)", *m.Type))
	}
	buf.WriteString(fmt.Sprintf(" %d)", m.Size))
	return buf.String()
}

type GlobalDef struct {
	Name    *string
	Value   Stmt
	Type    Stmt
	Special *SpecialIdT
	Anno    *int
}

func (g *GlobalDef) TopLevelType() TopLevelT {
	return GlobalDefT
}

func (g *GlobalDef) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString("(global ")
	if g.Name != nil {
		buf.WriteString(*g.Name + " ")
	}
	if g.Special != nil {
		buf.WriteString(g.Special.String() + " ")
	}
	if g.Anno != nil {
		buf.WriteString(fmt.Sprintf("(;%d;) ", *g.Anno))
	}
	buf.WriteString(fmt.Sprintf("(%s) (%s))", g.Type.IndentedString(0), g.Value.IndentedString(0)))
	return buf.String()
}

type ExportDef struct {
	Name   string
	Func   *FuncNameRef
	Memory *MemoryDef
}

func (e *ExportDef) TopLevelType() TopLevelT {
	return ExportDefT
}

func (e *ExportDef) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString(fmt.Sprintf("(export \"%s\"", e.Name))
	if e.Func != nil {
		buf.WriteString(fmt.Sprintf(" %s)", e.Func.String()))
	}
	if e.Memory != nil {
		buf.WriteString(fmt.Sprintf(" %s)", e.Memory.IndentedString(0)))
	}
	return buf.String()
}

type ElemDef struct {
	Const Stmt
	Ident []string
	Anno  *int
}

func (e *ElemDef) TopLevelType() TopLevelT {
	return ElemDefT
}

func (e *ElemDef) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString("(elem")
	if e.Anno != nil {
		buf.WriteString(fmt.Sprintf(" (;%d;)", *e.Anno))
	}
	buf.WriteString(fmt.Sprintf(" (%s) ", e.Const.IndentedString(0)))
	for i := 0; i < len(e.Ident); i++ {
		if i != 0 {
			buf.WriteString(" ")
		} else {
			buf.WriteString("func ")
		}
		buf.WriteString(e.Ident[i])
	}
	return buf.String() + ")"
}

type DataDef struct {
	Segment    string
	Const      Stmt
	QuotedData string
}

func (d *DataDef) TopLevelType() TopLevelT {
	return DataDefT
}

func (d *DataDef) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString(fmt.Sprintf("(data %s (%s) %s)", d.Segment, d.Const.IndentedString(0),
		d.QuotedData))
	return buf.String()
}
