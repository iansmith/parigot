// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package main // sexpr
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// A complete Visitor for a parse tree produced by sexprParser.
type sexprVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by sexprParser#sexpr.
	VisitSexpr(ctx *SexprContext) interface{}

	// Visit a parse tree produced by sexprParser#item.
	VisitItem(ctx *ItemContext) interface{}

	// Visit a parse tree produced by sexprParser#list_.
	VisitList_(ctx *List_Context) interface{}

	// Visit a parse tree produced by sexprParser#atom.
	VisitAtom(ctx *AtomContext) interface{}
}
