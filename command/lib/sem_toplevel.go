package lib

import (
	"fmt"
)

type TopLevelT int

const (
	TypeDefT   TopLevelT = 1
	ImportDefT TopLevelT = 2
	FuncDefT   TopLevelT = 3
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
	buf.WriteString(fmt.Sprintf("(import \"%s\" \"%s\" %s)",
		i.ModuleName, i.ImportedAs, i.FuncNameRef.String()))
	return buf.String()

}

// FuncDef represents WAT like this:   (func $_*internal/task.gcData_.swap (type 1) (param i32) (local i32) ...statements )
type FuncDef struct {
	Name   string
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
	buf.WriteString(fmt.Sprintf("(func %s",
		f.Name))
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
