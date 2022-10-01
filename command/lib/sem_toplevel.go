package lib

import (
	"fmt"
)

type TopLevelDefT int

const (
	TypeDefT   TopLevelDefT = 1
	ImportDefT TopLevelDefT = 2
	FuncDefT   TopLevelDefT = 3
)

// TopLevelDef instances write a carriage return after they are printed.
type TopLevelDef interface {
	IndentedStringer
	TopLevelType() TopLevelDefT
}

// TypeDef represents WAT like this:   (type (;9;) (func (param i64 i32 i32 i32) (result i64)))
type TypeDef struct {
	Annotation *TypeAnnotation
	Func       *FuncSpec
}

func (t *TypeDef) TopLevelType() TopLevelDefT {
	return TypeDefT
}

func (t *TypeDef) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString(fmt.Sprintf("(type %s %s)\n", t.Annotation.String(),
		t.Func.String()))
	return buf.String()
}

// ImportDef represents WAT like this: (import "wasi_snapshot_preview1" "fd_write" (func $runtime.fd_write (type 4)))
type ImportDef struct {
	ModuleName  string
	ImportedAs  string
	FuncNameRef *FuncNameRef
}

func (i *ImportDef) TopLevelType() TopLevelDefT {
	return ImportDefT
}

func (i *ImportDef) IndentedString(indented int) string {
	buf := NewIndentedBuffer(indented)
	buf.WriteString(fmt.Sprintf("(import \"%s\" \"%s\" %s)\n",
		i.ModuleName, i.ImportedAs, i.FuncNameRef.String()))
	return buf.String()

}

// FuncDef represents WAT like this:   (func $_*internal/task.gcData_.swap (type 1) (param i32) (local i32) ...statements )
type FuncDef struct {
	Name  string
	Type  *TypeRef
	Param *ParamDef
	Local *LocalDef
	Code  []Stmt
}

func (f *FuncDef) TopLevelType() TopLevelDefT {
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
	if f.Local != nil {
		buf.WriteString("\n")
		for i := 0; i < indented+2; i++ {
			buf.WriteString(" ")
		}
		buf.WriteString(f.Local.String())
	}
	buf.WriteString(stmtsToString(f.Code, indented))
	// strangely this is not terminated with a "\n"
	return buf.String() + ")"
}

func (f *FuncDef) AddStmt(s Stmt) {
	f.Code = append(f.Code, s)
}
