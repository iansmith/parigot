package main

import (
	"fmt"
	"log"
	"os"

	"github.com/iansmith/parigot/command/transform"

	"github.com/antlr/antlr4/runtime/Go/antlr"
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
	lexer := transform.NewWasmLexer(fs)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := transform.NewWasmParser(stream)

	// Finally parse the expression
	builder := &transform.Builder{}
	antlr.ParseTreeWalkerDefault.Walk(builder, p.Module())
	mod := builder.Module()
	fmt.Printf("%s", mod.IndentedString(0))
}
