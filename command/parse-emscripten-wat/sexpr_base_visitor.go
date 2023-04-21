// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package main // sexpr
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

type BasesexprVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BasesexprVisitor) VisitSexpr(ctx *SexprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasesexprVisitor) VisitItem(ctx *ItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasesexprVisitor) VisitList_(ctx *List_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasesexprVisitor) VisitAtom(ctx *AtomContext) interface{} {
	return v.VisitChildren(ctx)
}
