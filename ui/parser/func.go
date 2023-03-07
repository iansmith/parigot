package parser

type FuncInvoc struct {
	Name   *DocIdOrVar
	Actual []*FuncActual
}

func NewFuncInvoc(n *DocIdOrVar, actual []*FuncActual) *FuncInvoc {
	return &FuncInvoc{Name: n, Actual: actual}
}

type FuncActual struct {
	Var     string
	Literal string
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
