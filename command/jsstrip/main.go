package main

import (
	"log"
	"os"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type sexprListener struct {
	*BasejsstripListener
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
	antlr.ParseTreeWalkerDefault.Walk(&sexprListener{}, p.Sexpr())
}
