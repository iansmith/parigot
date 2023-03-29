// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// at top of file

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type wcllex struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var wcllexLexerStaticData struct {
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

func wcllexLexerInit() {
	staticData := &wcllexLexerStaticData
	staticData.channelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.modeNames = []string{
		"DEFAULT_MODE", "GrabText",
	}
	staticData.literalNames = []string{
		"", "'@text'", "'@css'", "'@preamble'", "'@doc'", "'@local'", "'@global'",
		"'@extern'", "'@pre'", "'@post'", "'@wcl'", "'@event'", "'@model'",
		"'@view'", "'@collection'", "'@controller'", "", "", "", "'<<'", "'->'",
		"'{'", "'}'", "'('", "')'", "", "','", "':'", "'<'", "'>'", "'.'", "'#'",
		"'-'", "'^'", "';'", "'+'", "'`'", "", "", "", "", "'>>'", "", "'\\>'",
	}
	staticData.symbolicNames = []string{
		"", "Text", "CSS", "Import", "Doc", "Local", "Global", "Extern", "Pre",
		"Post", "Wcl", "Event", "Model", "View", "ViewCollection", "Controller",
		"Id", "TypeStarter", "Version", "DoubleLess", "Arrow", "LCurly", "RCurly",
		"LParen", "RParen", "Dollar", "Comma", "Colon", "LessThan", "GreaterThan",
		"Dot", "Hash", "Dash", "Caret", "Semi", "Plus", "BackTick", "StringLit",
		"DoubleSlashComment", "Whitespace", "GrabDollar", "GrabDoubleGreater",
		"RawText", "GrabGreaterThan",
	}
	staticData.ruleNames = []string{
		"Text", "CSS", "Import", "Doc", "Local", "Global", "Extern", "Pre",
		"Post", "Wcl", "Event", "Model", "View", "ViewCollection", "Controller",
		"Id", "TypeStarter", "IdentFirst", "IdentAfter", "Version", "Digit",
		"DoubleLess", "Arrow", "LCurly", "RCurly", "LParen", "RParen", "Dollar",
		"Comma", "Colon", "LessThan", "GreaterThan", "Dot", "Hash", "Dash",
		"Caret", "Semi", "Plus", "BackTick", "StringLit", "Esc", "DoubleSlashComment",
		"Whitespace", "GrabDollar", "GrabGreaterThan", "GrabDoubleGreater",
		"RawText",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 43, 336, 6, -1, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3,
		7, 3, 2, 4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9,
		7, 9, 2, 10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7,
		14, 2, 15, 7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19,
		2, 20, 7, 20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2,
		25, 7, 25, 2, 26, 7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30,
		7, 30, 2, 31, 7, 31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7,
		35, 2, 36, 7, 36, 2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40,
		2, 41, 7, 41, 2, 42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2,
		46, 7, 46, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 3,
		1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5,
		1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6,
		1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 8,
		1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 10, 1,
		10, 1, 10, 1, 10, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 12,
		1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1,
		13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 14, 1, 14, 1, 14, 1, 14,
		1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 15, 1, 15, 5,
		15, 208, 8, 15, 10, 15, 12, 15, 211, 9, 15, 1, 16, 4, 16, 214, 8, 16, 11,
		16, 12, 16, 215, 1, 17, 1, 17, 1, 18, 1, 18, 3, 18, 222, 8, 18, 1, 19,
		4, 19, 225, 8, 19, 11, 19, 12, 19, 226, 1, 19, 1, 19, 4, 19, 231, 8, 19,
		11, 19, 12, 19, 232, 1, 19, 1, 19, 4, 19, 237, 8, 19, 11, 19, 12, 19, 238,
		1, 20, 1, 20, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 22, 1, 22, 1, 22, 1,
		23, 1, 23, 1, 24, 1, 24, 1, 25, 1, 25, 1, 26, 1, 26, 1, 27, 1, 27, 1, 28,
		1, 28, 1, 29, 1, 29, 1, 30, 1, 30, 1, 31, 1, 31, 1, 32, 1, 32, 1, 33, 1,
		33, 1, 34, 1, 34, 1, 35, 1, 35, 1, 36, 1, 36, 1, 37, 1, 37, 1, 38, 1, 38,
		1, 39, 1, 39, 1, 39, 5, 39, 286, 8, 39, 10, 39, 12, 39, 289, 9, 39, 1,
		39, 1, 39, 1, 40, 1, 40, 1, 40, 1, 40, 3, 40, 297, 8, 40, 1, 41, 1, 41,
		1, 41, 1, 41, 4, 41, 303, 8, 41, 11, 41, 12, 41, 304, 1, 41, 1, 41, 1,
		41, 1, 41, 1, 42, 4, 42, 312, 8, 42, 11, 42, 12, 42, 313, 1, 42, 1, 42,
		1, 43, 1, 43, 1, 43, 1, 43, 1, 44, 1, 44, 1, 44, 1, 44, 1, 44, 1, 45, 1,
		45, 1, 45, 1, 45, 1, 45, 1, 46, 4, 46, 333, 8, 46, 11, 46, 12, 46, 334,
		1, 304, 0, 47, 2, 1, 4, 2, 6, 3, 8, 4, 10, 5, 12, 6, 14, 7, 16, 8, 18,
		9, 20, 10, 22, 11, 24, 12, 26, 13, 28, 14, 30, 15, 32, 16, 34, 17, 36,
		0, 38, 0, 40, 18, 42, 0, 44, 19, 46, 20, 48, 21, 50, 22, 52, 23, 54, 24,
		56, 25, 58, 26, 60, 27, 62, 28, 64, 29, 66, 30, 68, 31, 70, 32, 72, 33,
		74, 34, 76, 35, 78, 36, 80, 37, 82, 0, 84, 38, 86, 39, 88, 40, 90, 43,
		92, 41, 94, 42, 2, 0, 1, 6, 3, 0, 42, 42, 91, 91, 93, 93, 4, 0, 45, 45,
		65, 90, 95, 95, 97, 122, 2, 0, 34, 34, 92, 92, 2, 0, 10, 10, 13, 13, 3,
		0, 0, 0, 9, 13, 32, 32, 3, 0, 36, 36, 62, 62, 92, 92, 342, 0, 2, 1, 0,
		0, 0, 0, 4, 1, 0, 0, 0, 0, 6, 1, 0, 0, 0, 0, 8, 1, 0, 0, 0, 0, 10, 1, 0,
		0, 0, 0, 12, 1, 0, 0, 0, 0, 14, 1, 0, 0, 0, 0, 16, 1, 0, 0, 0, 0, 18, 1,
		0, 0, 0, 0, 20, 1, 0, 0, 0, 0, 22, 1, 0, 0, 0, 0, 24, 1, 0, 0, 0, 0, 26,
		1, 0, 0, 0, 0, 28, 1, 0, 0, 0, 0, 30, 1, 0, 0, 0, 0, 32, 1, 0, 0, 0, 0,
		34, 1, 0, 0, 0, 0, 40, 1, 0, 0, 0, 0, 44, 1, 0, 0, 0, 0, 46, 1, 0, 0, 0,
		0, 48, 1, 0, 0, 0, 0, 50, 1, 0, 0, 0, 0, 52, 1, 0, 0, 0, 0, 54, 1, 0, 0,
		0, 0, 56, 1, 0, 0, 0, 0, 58, 1, 0, 0, 0, 0, 60, 1, 0, 0, 0, 0, 62, 1, 0,
		0, 0, 0, 64, 1, 0, 0, 0, 0, 66, 1, 0, 0, 0, 0, 68, 1, 0, 0, 0, 0, 70, 1,
		0, 0, 0, 0, 72, 1, 0, 0, 0, 0, 74, 1, 0, 0, 0, 0, 76, 1, 0, 0, 0, 0, 78,
		1, 0, 0, 0, 0, 80, 1, 0, 0, 0, 0, 84, 1, 0, 0, 0, 0, 86, 1, 0, 0, 0, 1,
		88, 1, 0, 0, 0, 1, 90, 1, 0, 0, 0, 1, 92, 1, 0, 0, 0, 1, 94, 1, 0, 0, 0,
		2, 96, 1, 0, 0, 0, 4, 102, 1, 0, 0, 0, 6, 107, 1, 0, 0, 0, 8, 117, 1, 0,
		0, 0, 10, 122, 1, 0, 0, 0, 12, 129, 1, 0, 0, 0, 14, 137, 1, 0, 0, 0, 16,
		145, 1, 0, 0, 0, 18, 150, 1, 0, 0, 0, 20, 156, 1, 0, 0, 0, 22, 161, 1,
		0, 0, 0, 24, 168, 1, 0, 0, 0, 26, 175, 1, 0, 0, 0, 28, 181, 1, 0, 0, 0,
		30, 193, 1, 0, 0, 0, 32, 205, 1, 0, 0, 0, 34, 213, 1, 0, 0, 0, 36, 217,
		1, 0, 0, 0, 38, 221, 1, 0, 0, 0, 40, 224, 1, 0, 0, 0, 42, 240, 1, 0, 0,
		0, 44, 242, 1, 0, 0, 0, 46, 247, 1, 0, 0, 0, 48, 250, 1, 0, 0, 0, 50, 252,
		1, 0, 0, 0, 52, 254, 1, 0, 0, 0, 54, 256, 1, 0, 0, 0, 56, 258, 1, 0, 0,
		0, 58, 260, 1, 0, 0, 0, 60, 262, 1, 0, 0, 0, 62, 264, 1, 0, 0, 0, 64, 266,
		1, 0, 0, 0, 66, 268, 1, 0, 0, 0, 68, 270, 1, 0, 0, 0, 70, 272, 1, 0, 0,
		0, 72, 274, 1, 0, 0, 0, 74, 276, 1, 0, 0, 0, 76, 278, 1, 0, 0, 0, 78, 280,
		1, 0, 0, 0, 80, 282, 1, 0, 0, 0, 82, 296, 1, 0, 0, 0, 84, 298, 1, 0, 0,
		0, 86, 311, 1, 0, 0, 0, 88, 317, 1, 0, 0, 0, 90, 321, 1, 0, 0, 0, 92, 326,
		1, 0, 0, 0, 94, 332, 1, 0, 0, 0, 96, 97, 5, 64, 0, 0, 97, 98, 5, 116, 0,
		0, 98, 99, 5, 101, 0, 0, 99, 100, 5, 120, 0, 0, 100, 101, 5, 116, 0, 0,
		101, 3, 1, 0, 0, 0, 102, 103, 5, 64, 0, 0, 103, 104, 5, 99, 0, 0, 104,
		105, 5, 115, 0, 0, 105, 106, 5, 115, 0, 0, 106, 5, 1, 0, 0, 0, 107, 108,
		5, 64, 0, 0, 108, 109, 5, 112, 0, 0, 109, 110, 5, 114, 0, 0, 110, 111,
		5, 101, 0, 0, 111, 112, 5, 97, 0, 0, 112, 113, 5, 109, 0, 0, 113, 114,
		5, 98, 0, 0, 114, 115, 5, 108, 0, 0, 115, 116, 5, 101, 0, 0, 116, 7, 1,
		0, 0, 0, 117, 118, 5, 64, 0, 0, 118, 119, 5, 100, 0, 0, 119, 120, 5, 111,
		0, 0, 120, 121, 5, 99, 0, 0, 121, 9, 1, 0, 0, 0, 122, 123, 5, 64, 0, 0,
		123, 124, 5, 108, 0, 0, 124, 125, 5, 111, 0, 0, 125, 126, 5, 99, 0, 0,
		126, 127, 5, 97, 0, 0, 127, 128, 5, 108, 0, 0, 128, 11, 1, 0, 0, 0, 129,
		130, 5, 64, 0, 0, 130, 131, 5, 103, 0, 0, 131, 132, 5, 108, 0, 0, 132,
		133, 5, 111, 0, 0, 133, 134, 5, 98, 0, 0, 134, 135, 5, 97, 0, 0, 135, 136,
		5, 108, 0, 0, 136, 13, 1, 0, 0, 0, 137, 138, 5, 64, 0, 0, 138, 139, 5,
		101, 0, 0, 139, 140, 5, 120, 0, 0, 140, 141, 5, 116, 0, 0, 141, 142, 5,
		101, 0, 0, 142, 143, 5, 114, 0, 0, 143, 144, 5, 110, 0, 0, 144, 15, 1,
		0, 0, 0, 145, 146, 5, 64, 0, 0, 146, 147, 5, 112, 0, 0, 147, 148, 5, 114,
		0, 0, 148, 149, 5, 101, 0, 0, 149, 17, 1, 0, 0, 0, 150, 151, 5, 64, 0,
		0, 151, 152, 5, 112, 0, 0, 152, 153, 5, 111, 0, 0, 153, 154, 5, 115, 0,
		0, 154, 155, 5, 116, 0, 0, 155, 19, 1, 0, 0, 0, 156, 157, 5, 64, 0, 0,
		157, 158, 5, 119, 0, 0, 158, 159, 5, 99, 0, 0, 159, 160, 5, 108, 0, 0,
		160, 21, 1, 0, 0, 0, 161, 162, 5, 64, 0, 0, 162, 163, 5, 101, 0, 0, 163,
		164, 5, 118, 0, 0, 164, 165, 5, 101, 0, 0, 165, 166, 5, 110, 0, 0, 166,
		167, 5, 116, 0, 0, 167, 23, 1, 0, 0, 0, 168, 169, 5, 64, 0, 0, 169, 170,
		5, 109, 0, 0, 170, 171, 5, 111, 0, 0, 171, 172, 5, 100, 0, 0, 172, 173,
		5, 101, 0, 0, 173, 174, 5, 108, 0, 0, 174, 25, 1, 0, 0, 0, 175, 176, 5,
		64, 0, 0, 176, 177, 5, 118, 0, 0, 177, 178, 5, 105, 0, 0, 178, 179, 5,
		101, 0, 0, 179, 180, 5, 119, 0, 0, 180, 27, 1, 0, 0, 0, 181, 182, 5, 64,
		0, 0, 182, 183, 5, 99, 0, 0, 183, 184, 5, 111, 0, 0, 184, 185, 5, 108,
		0, 0, 185, 186, 5, 108, 0, 0, 186, 187, 5, 101, 0, 0, 187, 188, 5, 99,
		0, 0, 188, 189, 5, 116, 0, 0, 189, 190, 5, 105, 0, 0, 190, 191, 5, 111,
		0, 0, 191, 192, 5, 110, 0, 0, 192, 29, 1, 0, 0, 0, 193, 194, 5, 64, 0,
		0, 194, 195, 5, 99, 0, 0, 195, 196, 5, 111, 0, 0, 196, 197, 5, 110, 0,
		0, 197, 198, 5, 116, 0, 0, 198, 199, 5, 114, 0, 0, 199, 200, 5, 111, 0,
		0, 200, 201, 5, 108, 0, 0, 201, 202, 5, 108, 0, 0, 202, 203, 5, 101, 0,
		0, 203, 204, 5, 114, 0, 0, 204, 31, 1, 0, 0, 0, 205, 209, 3, 36, 17, 0,
		206, 208, 3, 38, 18, 0, 207, 206, 1, 0, 0, 0, 208, 211, 1, 0, 0, 0, 209,
		207, 1, 0, 0, 0, 209, 210, 1, 0, 0, 0, 210, 33, 1, 0, 0, 0, 211, 209, 1,
		0, 0, 0, 212, 214, 7, 0, 0, 0, 213, 212, 1, 0, 0, 0, 214, 215, 1, 0, 0,
		0, 215, 213, 1, 0, 0, 0, 215, 216, 1, 0, 0, 0, 216, 35, 1, 0, 0, 0, 217,
		218, 7, 1, 0, 0, 218, 37, 1, 0, 0, 0, 219, 222, 7, 1, 0, 0, 220, 222, 3,
		42, 20, 0, 221, 219, 1, 0, 0, 0, 221, 220, 1, 0, 0, 0, 222, 39, 1, 0, 0,
		0, 223, 225, 3, 42, 20, 0, 224, 223, 1, 0, 0, 0, 225, 226, 1, 0, 0, 0,
		226, 224, 1, 0, 0, 0, 226, 227, 1, 0, 0, 0, 227, 228, 1, 0, 0, 0, 228,
		230, 3, 66, 32, 0, 229, 231, 3, 42, 20, 0, 230, 229, 1, 0, 0, 0, 231, 232,
		1, 0, 0, 0, 232, 230, 1, 0, 0, 0, 232, 233, 1, 0, 0, 0, 233, 234, 1, 0,
		0, 0, 234, 236, 3, 66, 32, 0, 235, 237, 3, 42, 20, 0, 236, 235, 1, 0, 0,
		0, 237, 238, 1, 0, 0, 0, 238, 236, 1, 0, 0, 0, 238, 239, 1, 0, 0, 0, 239,
		41, 1, 0, 0, 0, 240, 241, 2, 48, 57, 0, 241, 43, 1, 0, 0, 0, 242, 243,
		5, 60, 0, 0, 243, 244, 5, 60, 0, 0, 244, 245, 1, 0, 0, 0, 245, 246, 6,
		21, 0, 0, 246, 45, 1, 0, 0, 0, 247, 248, 5, 45, 0, 0, 248, 249, 5, 62,
		0, 0, 249, 47, 1, 0, 0, 0, 250, 251, 5, 123, 0, 0, 251, 49, 1, 0, 0, 0,
		252, 253, 5, 125, 0, 0, 253, 51, 1, 0, 0, 0, 254, 255, 5, 40, 0, 0, 255,
		53, 1, 0, 0, 0, 256, 257, 5, 41, 0, 0, 257, 55, 1, 0, 0, 0, 258, 259, 5,
		36, 0, 0, 259, 57, 1, 0, 0, 0, 260, 261, 5, 44, 0, 0, 261, 59, 1, 0, 0,
		0, 262, 263, 5, 58, 0, 0, 263, 61, 1, 0, 0, 0, 264, 265, 5, 60, 0, 0, 265,
		63, 1, 0, 0, 0, 266, 267, 5, 62, 0, 0, 267, 65, 1, 0, 0, 0, 268, 269, 5,
		46, 0, 0, 269, 67, 1, 0, 0, 0, 270, 271, 5, 35, 0, 0, 271, 69, 1, 0, 0,
		0, 272, 273, 5, 45, 0, 0, 273, 71, 1, 0, 0, 0, 274, 275, 5, 94, 0, 0, 275,
		73, 1, 0, 0, 0, 276, 277, 5, 59, 0, 0, 277, 75, 1, 0, 0, 0, 278, 279, 5,
		43, 0, 0, 279, 77, 1, 0, 0, 0, 280, 281, 5, 96, 0, 0, 281, 79, 1, 0, 0,
		0, 282, 287, 5, 34, 0, 0, 283, 286, 3, 82, 40, 0, 284, 286, 8, 2, 0, 0,
		285, 283, 1, 0, 0, 0, 285, 284, 1, 0, 0, 0, 286, 289, 1, 0, 0, 0, 287,
		285, 1, 0, 0, 0, 287, 288, 1, 0, 0, 0, 288, 290, 1, 0, 0, 0, 289, 287,
		1, 0, 0, 0, 290, 291, 5, 34, 0, 0, 291, 81, 1, 0, 0, 0, 292, 293, 5, 92,
		0, 0, 293, 297, 5, 34, 0, 0, 294, 295, 5, 92, 0, 0, 295, 297, 5, 92, 0,
		0, 296, 292, 1, 0, 0, 0, 296, 294, 1, 0, 0, 0, 297, 83, 1, 0, 0, 0, 298,
		299, 5, 47, 0, 0, 299, 300, 5, 47, 0, 0, 300, 302, 1, 0, 0, 0, 301, 303,
		9, 0, 0, 0, 302, 301, 1, 0, 0, 0, 303, 304, 1, 0, 0, 0, 304, 305, 1, 0,
		0, 0, 304, 302, 1, 0, 0, 0, 305, 306, 1, 0, 0, 0, 306, 307, 7, 3, 0, 0,
		307, 308, 1, 0, 0, 0, 308, 309, 6, 41, 1, 0, 309, 85, 1, 0, 0, 0, 310,
		312, 7, 4, 0, 0, 311, 310, 1, 0, 0, 0, 312, 313, 1, 0, 0, 0, 313, 311,
		1, 0, 0, 0, 313, 314, 1, 0, 0, 0, 314, 315, 1, 0, 0, 0, 315, 316, 6, 42,
		1, 0, 316, 87, 1, 0, 0, 0, 317, 318, 5, 36, 0, 0, 318, 319, 1, 0, 0, 0,
		319, 320, 6, 43, 2, 0, 320, 89, 1, 0, 0, 0, 321, 322, 5, 92, 0, 0, 322,
		323, 5, 62, 0, 0, 323, 324, 1, 0, 0, 0, 324, 325, 6, 44, 3, 0, 325, 91,
		1, 0, 0, 0, 326, 327, 5, 62, 0, 0, 327, 328, 5, 62, 0, 0, 328, 329, 1,
		0, 0, 0, 329, 330, 6, 45, 2, 0, 330, 93, 1, 0, 0, 0, 331, 333, 8, 5, 0,
		0, 332, 331, 1, 0, 0, 0, 333, 334, 1, 0, 0, 0, 334, 332, 1, 0, 0, 0, 334,
		335, 1, 0, 0, 0, 335, 95, 1, 0, 0, 0, 14, 0, 1, 209, 215, 221, 226, 232,
		238, 285, 287, 296, 304, 313, 334, 4, 5, 1, 0, 6, 0, 0, 4, 0, 0, 7, 42,
		0,
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

// wcllexInit initializes any static state used to implement wcllex. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// Newwcllex(). You can call this function if you wish to initialize the static state ahead
// of time.
func WcllexInit() {
	staticData := &wcllexLexerStaticData
	staticData.once.Do(wcllexLexerInit)
}

// Newwcllex produces a new lexer instance for the optional input antlr.CharStream.
func Newwcllex(input antlr.CharStream) *wcllex {
	WcllexInit()
	l := new(wcllex)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &wcllexLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	l.channelNames = staticData.channelNames
	l.modeNames = staticData.modeNames
	l.RuleNames = staticData.ruleNames
	l.LiteralNames = staticData.literalNames
	l.SymbolicNames = staticData.symbolicNames
	l.GrammarFileName = "wcllex.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// wcllex tokens.
const (
	wcllexText               = 1
	wcllexCSS                = 2
	wcllexImport             = 3
	wcllexDoc                = 4
	wcllexLocal              = 5
	wcllexGlobal             = 6
	wcllexExtern             = 7
	wcllexPre                = 8
	wcllexPost               = 9
	wcllexWcl                = 10
	wcllexEvent              = 11
	wcllexModel              = 12
	wcllexView               = 13
	wcllexViewCollection     = 14
	wcllexController         = 15
	wcllexId                 = 16
	wcllexTypeStarter        = 17
	wcllexVersion            = 18
	wcllexDoubleLess         = 19
	wcllexArrow              = 20
	wcllexLCurly             = 21
	wcllexRCurly             = 22
	wcllexLParen             = 23
	wcllexRParen             = 24
	wcllexDollar             = 25
	wcllexComma              = 26
	wcllexColon              = 27
	wcllexLessThan           = 28
	wcllexGreaterThan        = 29
	wcllexDot                = 30
	wcllexHash               = 31
	wcllexDash               = 32
	wcllexCaret              = 33
	wcllexSemi               = 34
	wcllexPlus               = 35
	wcllexBackTick           = 36
	wcllexStringLit          = 37
	wcllexDoubleSlashComment = 38
	wcllexWhitespace         = 39
	wcllexGrabDollar         = 40
	wcllexGrabDoubleGreater  = 41
	wcllexRawText            = 42
	wcllexGrabGreaterThan    = 43
)

// wcllexGrabText is the wcllex mode.
const wcllexGrabText = 1
