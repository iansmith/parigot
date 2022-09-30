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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 11, 71, 8,
	1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9,
	7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 3, 2, 3, 2, 3, 3, 3, 3, 3, 4,
	6, 4, 27, 10, 4, 13, 4, 14, 4, 28, 3, 5, 6, 5, 32, 10, 5, 13, 5, 14, 5,
	33, 3, 6, 6, 6, 37, 10, 6, 13, 6, 14, 6, 38, 3, 7, 3, 7, 3, 7, 3, 7, 5,
	7, 45, 10, 7, 3, 7, 3, 7, 3, 8, 3, 8, 3, 8, 7, 8, 52, 10, 8, 12, 8, 14,
	8, 55, 11, 8, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 9, 3, 10, 3, 10, 7, 10,
	65, 10, 10, 12, 10, 14, 10, 68, 11, 10, 3, 10, 3, 10, 2, 2, 11, 3, 3, 5,
	4, 7, 5, 9, 6, 11, 7, 13, 8, 15, 9, 17, 10, 19, 11, 3, 2, 6, 4, 2, 50,
	59, 99, 104, 7, 2, 38, 38, 48, 48, 67, 92, 97, 97, 99, 124, 4, 2, 36, 36,
	94, 94, 4, 2, 12, 12, 15, 15, 2, 78, 2, 3, 3, 2, 2, 2, 2, 5, 3, 2, 2, 2,
	2, 7, 3, 2, 2, 2, 2, 9, 3, 2, 2, 2, 2, 11, 3, 2, 2, 2, 2, 13, 3, 2, 2,
	2, 2, 15, 3, 2, 2, 2, 2, 17, 3, 2, 2, 2, 2, 19, 3, 2, 2, 2, 3, 21, 3, 2,
	2, 2, 5, 23, 3, 2, 2, 2, 7, 26, 3, 2, 2, 2, 9, 31, 3, 2, 2, 2, 11, 36,
	3, 2, 2, 2, 13, 44, 3, 2, 2, 2, 15, 48, 3, 2, 2, 2, 17, 58, 3, 2, 2, 2,
	19, 62, 3, 2, 2, 2, 21, 22, 7, 42, 2, 2, 22, 4, 3, 2, 2, 2, 23, 24, 7,
	43, 2, 2, 24, 6, 3, 2, 2, 2, 25, 27, 4, 50, 59, 2, 26, 25, 3, 2, 2, 2,
	27, 28, 3, 2, 2, 2, 28, 26, 3, 2, 2, 2, 28, 29, 3, 2, 2, 2, 29, 8, 3, 2,
	2, 2, 30, 32, 9, 2, 2, 2, 31, 30, 3, 2, 2, 2, 32, 33, 3, 2, 2, 2, 33, 31,
	3, 2, 2, 2, 33, 34, 3, 2, 2, 2, 34, 10, 3, 2, 2, 2, 35, 37, 9, 3, 2, 2,
	36, 35, 3, 2, 2, 2, 37, 38, 3, 2, 2, 2, 38, 36, 3, 2, 2, 2, 38, 39, 3,
	2, 2, 2, 39, 12, 3, 2, 2, 2, 40, 45, 7, 34, 2, 2, 41, 42, 7, 15, 2, 2,
	42, 45, 7, 12, 2, 2, 43, 45, 4, 11, 12, 2, 44, 40, 3, 2, 2, 2, 44, 41,
	3, 2, 2, 2, 44, 43, 3, 2, 2, 2, 45, 46, 3, 2, 2, 2, 46, 47, 8, 7, 2, 2,
	47, 14, 3, 2, 2, 2, 48, 53, 7, 36, 2, 2, 49, 52, 10, 4, 2, 2, 50, 52, 5,
	17, 9, 2, 51, 49, 3, 2, 2, 2, 51, 50, 3, 2, 2, 2, 52, 55, 3, 2, 2, 2, 53,
	51, 3, 2, 2, 2, 53, 54, 3, 2, 2, 2, 54, 56, 3, 2, 2, 2, 55, 53, 3, 2, 2,
	2, 56, 57, 7, 36, 2, 2, 57, 16, 3, 2, 2, 2, 58, 59, 7, 94, 2, 2, 59, 60,
	5, 9, 5, 2, 60, 61, 5, 9, 5, 2, 61, 18, 3, 2, 2, 2, 62, 66, 7, 61, 2, 2,
	63, 65, 10, 5, 2, 2, 64, 63, 3, 2, 2, 2, 65, 68, 3, 2, 2, 2, 66, 64, 3,
	2, 2, 2, 66, 67, 3, 2, 2, 2, 67, 69, 3, 2, 2, 2, 68, 66, 3, 2, 2, 2, 69,
	70, 8, 10, 2, 2, 70, 20, 3, 2, 2, 2, 10, 2, 28, 33, 38, 44, 51, 53, 66,
	3, 8, 2, 2,
}

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "'('", "')'",
}

var lexerSymbolicNames = []string{
	"", "", "", "Num", "HexDigit", "Id", "Whitespace", "String_", "HexByteValue",
	"LineComment",
}

var lexerRuleNames = []string{
	"T__0", "T__1", "Num", "HexDigit", "Id", "Whitespace", "String_", "HexByteValue",
	"LineComment",
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
	jsstripLexerT__0         = 1
	jsstripLexerT__1         = 2
	jsstripLexerNum          = 3
	jsstripLexerHexDigit     = 4
	jsstripLexerId           = 5
	jsstripLexerWhitespace   = 6
	jsstripLexerString_      = 7
	jsstripLexerHexByteValue = 8
	jsstripLexerLineComment  = 9
)
