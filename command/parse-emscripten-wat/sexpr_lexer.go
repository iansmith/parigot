// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package main

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type sexprLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var sexprlexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	channelNames           []string
	modeNames              []string
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func sexprlexerLexerInit() {
	staticData := &sexprlexerLexerStaticData
	staticData.channelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.modeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.literalNames = []string{
		"", "", "", "", "", "", "'('", "')'", "'.'",
	}
	staticData.symbolicNames = []string{
		"", "STRING", "WHITESPACE", "NUMBER", "SYMBOL", "COMMENT_NUM", "LPAREN",
		"RPAREN", "DOT",
	}
	staticData.ruleNames = []string{
		"STRING", "WHITESPACE", "NUMBER", "SYMBOL", "COMMENT_NUM", "LPAREN",
		"RPAREN", "DOT", "SYMBOL_START", "DIGIT",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 8, 82, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 1,
		0, 1, 0, 1, 0, 1, 0, 5, 0, 26, 8, 0, 10, 0, 12, 0, 29, 9, 0, 1, 0, 1, 0,
		1, 1, 4, 1, 34, 8, 1, 11, 1, 12, 1, 35, 1, 1, 1, 1, 1, 2, 3, 2, 41, 8,
		2, 1, 2, 4, 2, 44, 8, 2, 11, 2, 12, 2, 45, 1, 2, 1, 2, 4, 2, 50, 8, 2,
		11, 2, 12, 2, 51, 3, 2, 54, 8, 2, 1, 3, 1, 3, 1, 3, 5, 3, 59, 8, 3, 10,
		3, 12, 3, 62, 9, 3, 1, 4, 1, 4, 4, 4, 66, 8, 4, 11, 4, 12, 4, 67, 1, 4,
		1, 4, 1, 5, 1, 5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 8, 3, 8, 79, 8, 8, 1, 9, 1,
		9, 0, 0, 10, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 17, 0,
		19, 0, 1, 0, 4, 2, 0, 34, 34, 92, 92, 3, 0, 9, 10, 13, 13, 32, 32, 2, 0,
		43, 43, 45, 45, 4, 0, 42, 43, 45, 47, 65, 90, 97, 122, 89, 0, 1, 1, 0,
		0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0,
		0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 1, 21, 1,
		0, 0, 0, 3, 33, 1, 0, 0, 0, 5, 40, 1, 0, 0, 0, 7, 55, 1, 0, 0, 0, 9, 63,
		1, 0, 0, 0, 11, 71, 1, 0, 0, 0, 13, 73, 1, 0, 0, 0, 15, 75, 1, 0, 0, 0,
		17, 78, 1, 0, 0, 0, 19, 80, 1, 0, 0, 0, 21, 27, 5, 34, 0, 0, 22, 23, 5,
		92, 0, 0, 23, 26, 9, 0, 0, 0, 24, 26, 8, 0, 0, 0, 25, 22, 1, 0, 0, 0, 25,
		24, 1, 0, 0, 0, 26, 29, 1, 0, 0, 0, 27, 25, 1, 0, 0, 0, 27, 28, 1, 0, 0,
		0, 28, 30, 1, 0, 0, 0, 29, 27, 1, 0, 0, 0, 30, 31, 5, 34, 0, 0, 31, 2,
		1, 0, 0, 0, 32, 34, 7, 1, 0, 0, 33, 32, 1, 0, 0, 0, 34, 35, 1, 0, 0, 0,
		35, 33, 1, 0, 0, 0, 35, 36, 1, 0, 0, 0, 36, 37, 1, 0, 0, 0, 37, 38, 6,
		1, 0, 0, 38, 4, 1, 0, 0, 0, 39, 41, 7, 2, 0, 0, 40, 39, 1, 0, 0, 0, 40,
		41, 1, 0, 0, 0, 41, 43, 1, 0, 0, 0, 42, 44, 3, 19, 9, 0, 43, 42, 1, 0,
		0, 0, 44, 45, 1, 0, 0, 0, 45, 43, 1, 0, 0, 0, 45, 46, 1, 0, 0, 0, 46, 53,
		1, 0, 0, 0, 47, 49, 5, 46, 0, 0, 48, 50, 3, 19, 9, 0, 49, 48, 1, 0, 0,
		0, 50, 51, 1, 0, 0, 0, 51, 49, 1, 0, 0, 0, 51, 52, 1, 0, 0, 0, 52, 54,
		1, 0, 0, 0, 53, 47, 1, 0, 0, 0, 53, 54, 1, 0, 0, 0, 54, 6, 1, 0, 0, 0,
		55, 60, 3, 17, 8, 0, 56, 59, 3, 17, 8, 0, 57, 59, 3, 19, 9, 0, 58, 56,
		1, 0, 0, 0, 58, 57, 1, 0, 0, 0, 59, 62, 1, 0, 0, 0, 60, 58, 1, 0, 0, 0,
		60, 61, 1, 0, 0, 0, 61, 8, 1, 0, 0, 0, 62, 60, 1, 0, 0, 0, 63, 65, 5, 59,
		0, 0, 64, 66, 3, 19, 9, 0, 65, 64, 1, 0, 0, 0, 66, 67, 1, 0, 0, 0, 67,
		65, 1, 0, 0, 0, 67, 68, 1, 0, 0, 0, 68, 69, 1, 0, 0, 0, 69, 70, 5, 59,
		0, 0, 70, 10, 1, 0, 0, 0, 71, 72, 5, 40, 0, 0, 72, 12, 1, 0, 0, 0, 73,
		74, 5, 41, 0, 0, 74, 14, 1, 0, 0, 0, 75, 76, 5, 46, 0, 0, 76, 16, 1, 0,
		0, 0, 77, 79, 7, 3, 0, 0, 78, 77, 1, 0, 0, 0, 79, 18, 1, 0, 0, 0, 80, 81,
		2, 48, 57, 0, 81, 20, 1, 0, 0, 0, 12, 0, 25, 27, 35, 40, 45, 51, 53, 58,
		60, 67, 78, 1, 6, 0, 0,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// sexprLexerInit initializes any static state used to implement sexprLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewsexprLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func SexprLexerInit() {
	staticData := &sexprlexerLexerStaticData
	staticData.once.Do(sexprlexerLexerInit)
}

// NewsexprLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewsexprLexer(input antlr.CharStream) *sexprLexer {
	SexprLexerInit()
	l := new(sexprLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &sexprlexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	l.channelNames = staticData.channelNames
	l.modeNames = staticData.modeNames
	l.RuleNames = staticData.ruleNames
	l.LiteralNames = staticData.literalNames
	l.SymbolicNames = staticData.symbolicNames
	l.GrammarFileName = "sexpr.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// sexprLexer tokens.
const (
	sexprLexerSTRING      = 1
	sexprLexerWHITESPACE  = 2
	sexprLexerNUMBER      = 3
	sexprLexerSYMBOL      = 4
	sexprLexerCOMMENT_NUM = 5
	sexprLexerLPAREN      = 6
	sexprLexerRPAREN      = 7
	sexprLexerDOT         = 8
)
