package lib

import (
	"bytes"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type atom struct {
	s string
	i string
	n int
	t int
}

type member struct {
	a *atom
	s *sexpr
}

// no e to avoid clash of names
type sexpr struct {
	members []*member
}

type program struct {
	sexpr []*sexpr
}

type SexprListener struct {
	*BaseWasmListener

	// these are the intermediate nodes on the way up the tree construction
	currentAtom       *atom
	currentMember     *member
	currentSexpr      *sexpr
	currentProgram    *program
	currentTypeNumber *int
	currentString     *bytes.Buffer
}

// VisitTerminal is called when a terminal node is visited.
func (s *SexprListener) VisitTerminal(_ antlr.TerminalNode) {
	//fmt.Printf("found terminal=%s\n", node.GetSymbol())
}
