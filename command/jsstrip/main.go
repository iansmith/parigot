package main

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/iansmith/parigot/command/lib"
	"log"
	"os"
)

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
	lexer := lib.NewWasmLexer(fs)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := lib.NewWasmParser(stream)

	// Finally parse the expression
	builder := &lib.Builder{}
	antlr.ParseTreeWalkerDefault.Walk(builder, p.Module())
	mod := builder.Module()
	fmt.Printf("%s", mod.IndentedString(0))
}
