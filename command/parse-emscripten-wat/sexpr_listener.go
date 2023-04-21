// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package main // sexpr
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// sexprListener is a complete listener for a parse tree produced by sexprParser.
type sexprListener interface {
	antlr.ParseTreeListener

	// EnterSexpr is called when entering the sexpr production.
	EnterSexpr(c *SexprContext)

	// EnterItem is called when entering the item production.
	EnterItem(c *ItemContext)

	// EnterList_ is called when entering the list_ production.
	EnterList_(c *List_Context)

	// EnterAtom is called when entering the atom production.
	EnterAtom(c *AtomContext)

	// ExitSexpr is called when exiting the sexpr production.
	ExitSexpr(c *SexprContext)

	// ExitItem is called when exiting the item production.
	ExitItem(c *ItemContext)

	// ExitList_ is called when exiting the list_ production.
	ExitList_(c *List_Context)

	// ExitAtom is called when exiting the atom production.
	ExitAtom(c *AtomContext)
}
