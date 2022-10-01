package lib

import (
	"bytes"
	"fmt"
)

// FuncSpec represents something like this: (func (param i64 i32 i32 i32) (result i64))
// and is used in inside the TypeDef node.
type FuncSpec struct {
	Param  *ParamDef
	Result *ResultDef
}

func (f *FuncSpec) String() string {
	var buf bytes.Buffer
	buf.WriteString("(func")
	if f.Param != nil {
		buf.WriteString(" " + f.Param.String())
	}
	if f.Result != nil {
		buf.WriteString(" " + f.Result.String())
	}
	buf.WriteString(")")
	return buf.String()
}

// ParamDef represents something that looks like this: (param i32 i32)
type ParamDef struct {
	Type *TypeNameSeq
}

func (f *ParamDef) String() string {
	return fmt.Sprintf("(param %s)", f.Type.String())
}

// ResultDef represents something that looks like this: (result i32)
// There can be multiple results.
type ResultDef struct {
	Type *TypeNameSeq
}

func (r *ResultDef) String() string {
	return fmt.Sprintf("(result %s)", r.Type.String())
}

// FuncNameRef represents something that looks like this:  (func $syscall/js.stringVal (type 12)))
// and is used by the ImportDef struct.
type FuncNameRef struct {
	Name string
	Type *TypeRef
}

func (f *FuncNameRef) String() string {
	n := f.Name
	if n != "" {
		n = " " + n
	}
	if f.Type != nil {
		return fmt.Sprintf("(func%s %s)", n, f.Type.String())
	}

	return fmt.Sprintf("(func%s)", n)
}

// TypeRef represents something like (type 4)
type TypeRef struct {
	Num int
}

func (t *TypeRef) String() string {
	return fmt.Sprintf("(type %d)", t.Num)
}

// LocalDef represents something like (local i32 i32 i64)
type LocalDef struct {
	Type *TypeNameSeq
}

func (l *LocalDef) String() string {
	return fmt.Sprintf("(local %s)", l.Type.String())
}
