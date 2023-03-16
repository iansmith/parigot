package tree

import (
	"bytes"
)

type FuncInvoc struct {
	Name    *DocIdOrVar
	Actual  []*FuncActual
	Builtin bool
}

func NewFuncInvoc(n *DocIdOrVar, actual []*FuncActual) *FuncInvoc {
	return &FuncInvoc{Name: n, Actual: actual}
}

func (f *FuncInvoc) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString(f.Name.Name + "(")
	for i := 0; i < len(f.Actual); i++ {
		buf.WriteString(f.Actual[i].String())
		buf.WriteString(",")
	}
	buf.WriteString(")")
	return buf.String()
}

type FuncActual struct {
	Var     string
	Literal string
}

func (f *FuncActual) String() string {
	if f.Var != "" {
		return f.Var
	}
	return f.Literal
}

func NewFuncActual(v, lit string) *FuncActual {
	if v == "" && lit == "" {
		panic("unable to understand function actual (neither variable reference nor literal)")
	}
	if v != "" && lit != "" {
		panic("unable to understand function actual (both variable reference and literal)")
	}
	if v != "" {
		return &FuncActual{Var: v}
	}
	return &FuncActual{Literal: lit}
}
