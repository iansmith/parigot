// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package main // sexpr
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// BasesexprListener is a complete listener for a parse tree produced by sexprParser.
type BasesexprListener struct{}

var _ sexprListener = &BasesexprListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasesexprListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasesexprListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasesexprListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasesexprListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterSexpr is called when production sexpr is entered.
func (s *BasesexprListener) EnterSexpr(ctx *SexprContext) {}

// ExitSexpr is called when production sexpr is exited.
func (s *BasesexprListener) ExitSexpr(ctx *SexprContext) {}

// EnterItem is called when production item is entered.
func (s *BasesexprListener) EnterItem(ctx *ItemContext) {}

// ExitItem is called when production item is exited.
func (s *BasesexprListener) ExitItem(ctx *ItemContext) {}

// EnterList_ is called when production list_ is entered.
func (s *BasesexprListener) EnterList_(ctx *List_Context) {}

// ExitList_ is called when production list_ is exited.
func (s *BasesexprListener) ExitList_(ctx *List_Context) {}

// EnterAtom is called when production atom is entered.
func (s *BasesexprListener) EnterAtom(ctx *AtomContext) {}

// ExitAtom is called when production atom is exited.
func (s *BasesexprListener) ExitAtom(ctx *AtomContext) {}
