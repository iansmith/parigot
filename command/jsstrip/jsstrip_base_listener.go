// Code generated from command/jsstrip/jsstrip.g4 by ANTLR 4.9. DO NOT EDIT.

package main // jsstrip
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BasejsstripListener is a complete listener for a parse tree produced by jsstripParser.
type BasejsstripListener struct{}

var _ jsstripListener = &BasejsstripListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasejsstripListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasejsstripListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasejsstripListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasejsstripListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterSexpr is called when production sexpr is entered.
func (s *BasejsstripListener) EnterSexpr(ctx *SexprContext) {}

// ExitSexpr is called when production sexpr is exited.
func (s *BasejsstripListener) ExitSexpr(ctx *SexprContext) {}

// EnterList is called when production list is entered.
func (s *BasejsstripListener) EnterList(ctx *ListContext) {}

// ExitList is called when production list is exited.
func (s *BasejsstripListener) ExitList(ctx *ListContext) {}

// EnterMembers is called when production members is entered.
func (s *BasejsstripListener) EnterMembers(ctx *MembersContext) {}

// ExitMembers is called when production members is exited.
func (s *BasejsstripListener) ExitMembers(ctx *MembersContext) {}

// EnterAtom is called when production atom is entered.
func (s *BasejsstripListener) EnterAtom(ctx *AtomContext) {}

// ExitAtom is called when production atom is exited.
func (s *BasejsstripListener) ExitAtom(ctx *AtomContext) {}
