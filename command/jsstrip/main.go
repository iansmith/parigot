package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"

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
	s *sxpr
}

// no e to avoid clash of names
type sxpr struct {
	members []*member
}

type program struct {
	sexpr []*sxpr
}

type sexprListener struct {
	*BasejsstripListener

	// these are the intermediate nodes on the way up the tree construction
	currentAtom       *atom
	currentMember     *member
	currentSexpr      *sxpr
	currentProgram    *program
	currentTypeNumber *int
	currentString     *bytes.Buffer
}

// VisitTerminal is called when a terminal node is visited.
func (s *sexprListener) VisitTerminal(node antlr.TerminalNode) {
	fmt.Printf("ccc terminal=%s\n", node.GetSymbol())
}

// EnterAtom is called when production atom is entered.
func (s *sexprListener) xEnterAtom(_ *AtomContext) {
	s.currentAtom = &atom{}
}

// ExitAtom is called when production atom is exited.
func (s *sexprListener) xExitAtom(ctx *AtomContext) {
	terminal := ctx.GetStart()
	typ := terminal.GetTokenType()
	switch typ {
	case jsstripParserNum:
		n, err := strconv.Atoi(terminal.GetText())
		if err != nil {
			panic("should never happen,bad num")
		}
		s.currentAtom.n = n
	}
	fmt.Printf("xxx %+v, %d\n", ctx.GetStart().GetText(), typ)
}

func main() {

	if len(os.Args) != 2 {
		log.Fatalf("unable to understand arguments, should provide exactly one wasm program as argument")
	}

	// Set up the input
	fs, err := antlr.NewFileStream(os.Args[1])
	if err != nil {
		log.Fatalf("failed trying to open input file, %v", err)
	}
	// make lexer
	lexer := NewjsstripLexer(fs)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := NewjsstripParser(stream)

	// Finally parse the expression
	antlr.ParseTreeWalkerDefault.Walk(&sexprListener{}, p.Module())
}
