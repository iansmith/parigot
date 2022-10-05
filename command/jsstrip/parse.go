package main

import (
	"log"

	"github.com/iansmith/parigot/command/transform"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func parse(inputFilename string) *transform.Module {
	// Set up the input
	fs, err := antlr.NewFileStream(inputFilename)
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
	return builder.Module() // only one module right now
}
