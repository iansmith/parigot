// Code generated from command/jsstrip/jsstrip.g4 by ANTLR 4.9. DO NOT EDIT.

package main // jsstrip
import "github.com/antlr/antlr4/runtime/Go/antlr"

// jsstripListener is a complete listener for a parse tree produced by jsstripParser.
type jsstripListener interface {
	antlr.ParseTreeListener

	// EnterSexpr is called when entering the sexpr production.
	EnterSexpr(c *SexprContext)

	// EnterList is called when entering the list production.
	EnterList(c *ListContext)

	// EnterMembers is called when entering the members production.
	EnterMembers(c *MembersContext)

	// EnterAtom is called when entering the atom production.
	EnterAtom(c *AtomContext)

	// ExitSexpr is called when exiting the sexpr production.
	ExitSexpr(c *SexprContext)

	// ExitList is called when exiting the list production.
	ExitList(c *ListContext)

	// ExitMembers is called when exiting the members production.
	ExitMembers(c *MembersContext)

	// ExitAtom is called when exiting the atom production.
	ExitAtom(c *AtomContext)
}
