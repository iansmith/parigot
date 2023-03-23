package tree

import (
	"bytes"
)

type FuncInvoc struct {
	Name                     *Ident
	Actual                   []*FuncActual
	Builtin                  bool
	AnonBody                 []TextItem
	LineNumber, ColumnNumber int
}

func NewFuncInvoc(n *Ident, actual []*FuncActual, line, col int) *FuncInvoc {
	return &FuncInvoc{Name: n, Actual: actual, LineNumber: line, ColumnNumber: col}
}

func (f *FuncInvoc) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString(f.Name.String() + "(")
	for i := 0; i < len(f.Actual); i++ {
		buf.WriteString(f.Actual[i].String())
		buf.WriteString(",")
	}
	buf.WriteString(")")
	return buf.String()
}

type FuncActual struct {
	Ref *ValueRef
}

func (f *FuncActual) String() string {
	return f.Ref.String()
}

func NewFuncActual(vr *ValueRef) *FuncActual {
	return &FuncActual{Ref: vr}
}
