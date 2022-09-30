// Code generated from command/jsstrip/jsstrip.g4 by ANTLR 4.9. DO NOT EDIT.

package main

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 9, 67, 8,
	1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9,
	7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 3, 2, 3, 2, 3, 2,
	3, 2, 5, 2, 28, 10, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 4, 3, 4, 3, 5, 3, 5,
	3, 6, 6, 6, 39, 10, 6, 13, 6, 14, 6, 40, 3, 7, 3, 7, 3, 8, 3, 8, 3, 9,
	3, 9, 7, 9, 49, 10, 9, 12, 9, 14, 9, 52, 11, 9, 3, 10, 3, 10, 3, 10, 3,
	10, 3, 11, 3, 11, 3, 11, 7, 11, 61, 10, 11, 12, 11, 14, 11, 64, 11, 11,
	3, 11, 3, 11, 2, 2, 12, 3, 3, 5, 4, 7, 5, 9, 6, 11, 7, 13, 2, 15, 2, 17,
	8, 19, 2, 21, 9, 3, 2, 6, 8, 2, 38, 38, 48, 48, 61, 61, 67, 92, 97, 97,
	99, 124, 9, 2, 38, 38, 48, 48, 50, 59, 61, 61, 67, 92, 97, 97, 99, 124,
	4, 2, 50, 59, 99, 104, 3, 2, 36, 36, 2, 69, 2, 3, 3, 2, 2, 2, 2, 5, 3,
	2, 2, 2, 2, 7, 3, 2, 2, 2, 2, 9, 3, 2, 2, 2, 2, 11, 3, 2, 2, 2, 2, 17,
	3, 2, 2, 2, 2, 21, 3, 2, 2, 2, 3, 27, 3, 2, 2, 2, 5, 31, 3, 2, 2, 2, 7,
	33, 3, 2, 2, 2, 9, 35, 3, 2, 2, 2, 11, 38, 3, 2, 2, 2, 13, 42, 3, 2, 2,
	2, 15, 44, 3, 2, 2, 2, 17, 46, 3, 2, 2, 2, 19, 53, 3, 2, 2, 2, 21, 57,
	3, 2, 2, 2, 23, 28, 7, 34, 2, 2, 24, 25, 7, 15, 2, 2, 25, 28, 7, 12, 2,
	2, 26, 28, 4, 11, 12, 2, 27, 23, 3, 2, 2, 2, 27, 24, 3, 2, 2, 2, 27, 26,
	3, 2, 2, 2, 28, 29, 3, 2, 2, 2, 29, 30, 8, 2, 2, 2, 30, 4, 3, 2, 2, 2,
	31, 32, 7, 42, 2, 2, 32, 6, 3, 2, 2, 2, 33, 34, 7, 43, 2, 2, 34, 8, 3,
	2, 2, 2, 35, 36, 7, 36, 2, 2, 36, 10, 3, 2, 2, 2, 37, 39, 4, 50, 59, 2,
	38, 37, 3, 2, 2, 2, 39, 40, 3, 2, 2, 2, 40, 38, 3, 2, 2, 2, 40, 41, 3,
	2, 2, 2, 41, 12, 3, 2, 2, 2, 42, 43, 9, 2, 2, 2, 43, 14, 3, 2, 2, 2, 44,
	45, 9, 3, 2, 2, 45, 16, 3, 2, 2, 2, 46, 50, 5, 13, 7, 2, 47, 49, 5, 15,
	8, 2, 48, 47, 3, 2, 2, 2, 49, 52, 3, 2, 2, 2, 50, 48, 3, 2, 2, 2, 50, 51,
	3, 2, 2, 2, 51, 18, 3, 2, 2, 2, 52, 50, 3, 2, 2, 2, 53, 54, 7, 94, 2, 2,
	54, 55, 9, 4, 2, 2, 55, 56, 9, 4, 2, 2, 56, 20, 3, 2, 2, 2, 57, 62, 7,
	36, 2, 2, 58, 61, 5, 19, 10, 2, 59, 61, 10, 5, 2, 2, 60, 58, 3, 2, 2, 2,
	60, 59, 3, 2, 2, 2, 61, 64, 3, 2, 2, 2, 62, 60, 3, 2, 2, 2, 62, 63, 3,
	2, 2, 2, 63, 65, 3, 2, 2, 2, 64, 62, 3, 2, 2, 2, 65, 66, 7, 36, 2, 2, 66,
	22, 3, 2, 2, 2, 8, 2, 27, 40, 50, 60, 62, 3, 8, 2, 2,
}

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "", "'('", "')'", "'\"'",
}

var lexerSymbolicNames = []string{
	"", "Whitespace", "Lparen", "Rparen", "Quote", "Num", "Ident", "QuotedString",
}

var lexerRuleNames = []string{
	"Whitespace", "Lparen", "Rparen", "Quote", "Num", "IdentFirst", "IdentAfter",
	"Ident", "HexByteValue", "QuotedString",
}

type jsstripLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

// NewjsstripLexer produces a new lexer instance for the optional input antlr.CharStream.
//
// The *jsstripLexer instance produced may be reused by calling the SetInputStream method.
// The initial lexer configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func NewjsstripLexer(input antlr.CharStream) *jsstripLexer {
	l := new(jsstripLexer)
	lexerDeserializer := antlr.NewATNDeserializer(nil)
	lexerAtn := lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)
	lexerDecisionToDFA := make([]*antlr.DFA, len(lexerAtn.DecisionToState))
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "jsstrip.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// jsstripLexer tokens.
const (
	jsstripLexerWhitespace   = 1
	jsstripLexerLparen       = 2
	jsstripLexerRparen       = 3
	jsstripLexerQuote        = 4
	jsstripLexerNum          = 5
	jsstripLexerIdent        = 6
	jsstripLexerQuotedString = 7
)
