// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // wcl
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type wcl struct {
	*antlr.BaseParser
}

var wclParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func wclParserInit() {
	staticData := &wclParserStaticData
	staticData.literalNames = []string{
		"", "'@text'", "'@css'", "'@preamble'", "'@doc'", "'@local'", "'@global'",
		"'@extern'", "'@pre'", "'@post'", "'@wcl'", "'@event'", "", "", "",
		"", "'('", "')'", "", "','", "'<'", "'>'", "'.'", "'#'", "'-'", "';'",
		"'+'",
	}
	staticData.symbolicNames = []string{
		"", "Text", "CSS", "Import", "Doc", "Local", "Global", "Extern", "Pre",
		"Post", "Wcl", "Event", "Id", "Version", "LCurly", "RCurly", "LParen",
		"RParen", "Dollar", "Comma", "LessThan", "GreaterThan", "Dot", "Hash",
		"Dash", "Semi", "Plus", "BackTick", "StringLit", "DoubleSlashComment",
		"Whitespace", "ContentRawText", "ContentDollar", "ContentBackTick",
		"UninterpRawText", "UninterpDollar", "UninterpLCurly", "UninterpRCurly",
		"VarRCurly", "VarId",
	}
	staticData.ruleNames = []string{
		"program", "global", "extern", "wcl_section", "import_section", "css_section",
		"css_filespec", "text_section", "text_func", "pre_code", "post_code",
		"text_func_local", "text_top", "text_content", "text_content_inner",
		"var_subs", "sub", "uninterp", "uninterp_inner", "uninterp_var", "param_spec",
		"param_pair", "doc_section", "doc_func", "doc_func_local", "doc_func_formal",
		"doc_tag", "id_or_var_ref", "var_ref", "doc_id", "doc_class", "doc_elem",
		"doc_elem_content", "doc_elem_text", "doc_elem_child", "func_invoc",
		"func_actual_seq", "func_actual", "event_section", "event_spec", "event_call",
		"selector",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 39, 346, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 1,
		0, 1, 0, 3, 0, 87, 8, 0, 1, 0, 3, 0, 90, 8, 0, 1, 0, 3, 0, 93, 8, 0, 1,
		0, 3, 0, 96, 8, 0, 1, 0, 3, 0, 99, 8, 0, 1, 0, 3, 0, 102, 8, 0, 1, 0, 3,
		0, 105, 8, 0, 1, 0, 3, 0, 108, 8, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 2,
		1, 2, 1, 2, 5, 2, 118, 8, 2, 10, 2, 12, 2, 121, 9, 2, 1, 2, 1, 2, 1, 3,
		1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 5, 5, 134, 8, 5, 10, 5,
		12, 5, 137, 9, 5, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 5, 7, 144, 8, 7, 10, 7,
		12, 7, 147, 9, 7, 1, 8, 1, 8, 3, 8, 151, 8, 8, 1, 8, 3, 8, 154, 8, 8, 1,
		8, 3, 8, 157, 8, 8, 1, 8, 1, 8, 3, 8, 161, 8, 8, 1, 9, 1, 9, 1, 9, 1, 9,
		1, 10, 1, 10, 1, 10, 1, 10, 1, 11, 1, 11, 1, 11, 1, 12, 1, 12, 1, 12, 3,
		12, 177, 8, 12, 1, 12, 1, 12, 1, 13, 5, 13, 182, 8, 13, 10, 13, 12, 13,
		185, 9, 13, 1, 14, 1, 14, 3, 14, 189, 8, 14, 1, 15, 1, 15, 1, 15, 1, 16,
		1, 16, 1, 16, 1, 17, 4, 17, 198, 8, 17, 11, 17, 12, 17, 199, 1, 17, 1,
		17, 1, 18, 1, 18, 1, 18, 1, 18, 3, 18, 208, 8, 18, 1, 19, 1, 19, 1, 19,
		1, 19, 1, 20, 1, 20, 5, 20, 216, 8, 20, 10, 20, 12, 20, 219, 9, 20, 1,
		20, 1, 20, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 3, 21, 228, 8, 21, 1, 22,
		1, 22, 5, 22, 232, 8, 22, 10, 22, 12, 22, 235, 9, 22, 1, 23, 1, 23, 1,
		23, 3, 23, 240, 8, 23, 1, 23, 3, 23, 243, 8, 23, 1, 23, 1, 23, 3, 23, 247,
		8, 23, 1, 24, 1, 24, 1, 24, 1, 25, 1, 25, 3, 25, 254, 8, 25, 1, 26, 1,
		26, 1, 26, 3, 26, 259, 8, 26, 1, 26, 3, 26, 262, 8, 26, 1, 26, 1, 26, 1,
		27, 1, 27, 3, 27, 268, 8, 27, 1, 28, 1, 28, 1, 28, 1, 28, 1, 29, 1, 29,
		1, 29, 1, 30, 4, 30, 278, 8, 30, 11, 30, 12, 30, 279, 1, 31, 1, 31, 1,
		31, 3, 31, 285, 8, 31, 1, 31, 3, 31, 288, 8, 31, 1, 32, 1, 32, 3, 32, 292,
		8, 32, 1, 33, 1, 33, 3, 33, 296, 8, 33, 1, 34, 1, 34, 5, 34, 300, 8, 34,
		10, 34, 12, 34, 303, 9, 34, 1, 34, 1, 34, 1, 35, 1, 35, 1, 35, 1, 35, 1,
		35, 1, 36, 1, 36, 1, 36, 5, 36, 315, 8, 36, 10, 36, 12, 36, 318, 9, 36,
		3, 36, 320, 8, 36, 1, 37, 1, 37, 1, 38, 1, 38, 5, 38, 326, 8, 38, 10, 38,
		12, 38, 329, 9, 38, 1, 39, 1, 39, 1, 39, 1, 39, 1, 40, 1, 40, 3, 40, 337,
		8, 40, 1, 40, 1, 40, 1, 41, 1, 41, 1, 41, 3, 41, 344, 8, 41, 1, 41, 0,
		0, 42, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34,
		36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70,
		72, 74, 76, 78, 80, 82, 0, 1, 2, 0, 12, 12, 28, 28, 346, 0, 84, 1, 0, 0,
		0, 2, 111, 1, 0, 0, 0, 4, 114, 1, 0, 0, 0, 6, 124, 1, 0, 0, 0, 8, 127,
		1, 0, 0, 0, 10, 131, 1, 0, 0, 0, 12, 138, 1, 0, 0, 0, 14, 141, 1, 0, 0,
		0, 16, 148, 1, 0, 0, 0, 18, 162, 1, 0, 0, 0, 20, 166, 1, 0, 0, 0, 22, 170,
		1, 0, 0, 0, 24, 173, 1, 0, 0, 0, 26, 183, 1, 0, 0, 0, 28, 188, 1, 0, 0,
		0, 30, 190, 1, 0, 0, 0, 32, 193, 1, 0, 0, 0, 34, 197, 1, 0, 0, 0, 36, 207,
		1, 0, 0, 0, 38, 209, 1, 0, 0, 0, 40, 213, 1, 0, 0, 0, 42, 227, 1, 0, 0,
		0, 44, 229, 1, 0, 0, 0, 46, 236, 1, 0, 0, 0, 48, 248, 1, 0, 0, 0, 50, 253,
		1, 0, 0, 0, 52, 255, 1, 0, 0, 0, 54, 267, 1, 0, 0, 0, 56, 269, 1, 0, 0,
		0, 58, 273, 1, 0, 0, 0, 60, 277, 1, 0, 0, 0, 62, 287, 1, 0, 0, 0, 64, 291,
		1, 0, 0, 0, 66, 295, 1, 0, 0, 0, 68, 297, 1, 0, 0, 0, 70, 306, 1, 0, 0,
		0, 72, 319, 1, 0, 0, 0, 74, 321, 1, 0, 0, 0, 76, 323, 1, 0, 0, 0, 78, 330,
		1, 0, 0, 0, 80, 336, 1, 0, 0, 0, 82, 343, 1, 0, 0, 0, 84, 86, 3, 6, 3,
		0, 85, 87, 3, 10, 5, 0, 86, 85, 1, 0, 0, 0, 86, 87, 1, 0, 0, 0, 87, 89,
		1, 0, 0, 0, 88, 90, 3, 8, 4, 0, 89, 88, 1, 0, 0, 0, 89, 90, 1, 0, 0, 0,
		90, 92, 1, 0, 0, 0, 91, 93, 3, 4, 2, 0, 92, 91, 1, 0, 0, 0, 92, 93, 1,
		0, 0, 0, 93, 95, 1, 0, 0, 0, 94, 96, 3, 2, 1, 0, 95, 94, 1, 0, 0, 0, 95,
		96, 1, 0, 0, 0, 96, 98, 1, 0, 0, 0, 97, 99, 3, 14, 7, 0, 98, 97, 1, 0,
		0, 0, 98, 99, 1, 0, 0, 0, 99, 101, 1, 0, 0, 0, 100, 102, 3, 10, 5, 0, 101,
		100, 1, 0, 0, 0, 101, 102, 1, 0, 0, 0, 102, 104, 1, 0, 0, 0, 103, 105,
		3, 44, 22, 0, 104, 103, 1, 0, 0, 0, 104, 105, 1, 0, 0, 0, 105, 107, 1,
		0, 0, 0, 106, 108, 3, 76, 38, 0, 107, 106, 1, 0, 0, 0, 107, 108, 1, 0,
		0, 0, 108, 109, 1, 0, 0, 0, 109, 110, 5, 0, 0, 1, 110, 1, 1, 0, 0, 0, 111,
		112, 5, 6, 0, 0, 112, 113, 3, 40, 20, 0, 113, 3, 1, 0, 0, 0, 114, 115,
		5, 7, 0, 0, 115, 119, 5, 16, 0, 0, 116, 118, 5, 12, 0, 0, 117, 116, 1,
		0, 0, 0, 118, 121, 1, 0, 0, 0, 119, 117, 1, 0, 0, 0, 119, 120, 1, 0, 0,
		0, 120, 122, 1, 0, 0, 0, 121, 119, 1, 0, 0, 0, 122, 123, 5, 17, 0, 0, 123,
		5, 1, 0, 0, 0, 124, 125, 5, 10, 0, 0, 125, 126, 5, 13, 0, 0, 126, 7, 1,
		0, 0, 0, 127, 128, 5, 3, 0, 0, 128, 129, 5, 14, 0, 0, 129, 130, 3, 34,
		17, 0, 130, 9, 1, 0, 0, 0, 131, 135, 5, 2, 0, 0, 132, 134, 3, 12, 6, 0,
		133, 132, 1, 0, 0, 0, 134, 137, 1, 0, 0, 0, 135, 133, 1, 0, 0, 0, 135,
		136, 1, 0, 0, 0, 136, 11, 1, 0, 0, 0, 137, 135, 1, 0, 0, 0, 138, 139, 5,
		26, 0, 0, 139, 140, 5, 28, 0, 0, 140, 13, 1, 0, 0, 0, 141, 145, 5, 1, 0,
		0, 142, 144, 3, 16, 8, 0, 143, 142, 1, 0, 0, 0, 144, 147, 1, 0, 0, 0, 145,
		143, 1, 0, 0, 0, 145, 146, 1, 0, 0, 0, 146, 15, 1, 0, 0, 0, 147, 145, 1,
		0, 0, 0, 148, 150, 5, 12, 0, 0, 149, 151, 3, 40, 20, 0, 150, 149, 1, 0,
		0, 0, 150, 151, 1, 0, 0, 0, 151, 153, 1, 0, 0, 0, 152, 154, 3, 22, 11,
		0, 153, 152, 1, 0, 0, 0, 153, 154, 1, 0, 0, 0, 154, 156, 1, 0, 0, 0, 155,
		157, 3, 18, 9, 0, 156, 155, 1, 0, 0, 0, 156, 157, 1, 0, 0, 0, 157, 158,
		1, 0, 0, 0, 158, 160, 3, 24, 12, 0, 159, 161, 3, 20, 10, 0, 160, 159, 1,
		0, 0, 0, 160, 161, 1, 0, 0, 0, 161, 17, 1, 0, 0, 0, 162, 163, 5, 8, 0,
		0, 163, 164, 5, 14, 0, 0, 164, 165, 3, 34, 17, 0, 165, 19, 1, 0, 0, 0,
		166, 167, 5, 9, 0, 0, 167, 168, 5, 14, 0, 0, 168, 169, 3, 34, 17, 0, 169,
		21, 1, 0, 0, 0, 170, 171, 5, 5, 0, 0, 171, 172, 3, 40, 20, 0, 172, 23,
		1, 0, 0, 0, 173, 176, 5, 27, 0, 0, 174, 177, 3, 26, 13, 0, 175, 177, 1,
		0, 0, 0, 176, 174, 1, 0, 0, 0, 176, 175, 1, 0, 0, 0, 177, 178, 1, 0, 0,
		0, 178, 179, 5, 33, 0, 0, 179, 25, 1, 0, 0, 0, 180, 182, 3, 28, 14, 0,
		181, 180, 1, 0, 0, 0, 182, 185, 1, 0, 0, 0, 183, 181, 1, 0, 0, 0, 183,
		184, 1, 0, 0, 0, 184, 27, 1, 0, 0, 0, 185, 183, 1, 0, 0, 0, 186, 189, 5,
		31, 0, 0, 187, 189, 3, 30, 15, 0, 188, 186, 1, 0, 0, 0, 188, 187, 1, 0,
		0, 0, 189, 29, 1, 0, 0, 0, 190, 191, 5, 32, 0, 0, 191, 192, 3, 32, 16,
		0, 192, 31, 1, 0, 0, 0, 193, 194, 5, 39, 0, 0, 194, 195, 5, 38, 0, 0, 195,
		33, 1, 0, 0, 0, 196, 198, 3, 36, 18, 0, 197, 196, 1, 0, 0, 0, 198, 199,
		1, 0, 0, 0, 199, 197, 1, 0, 0, 0, 199, 200, 1, 0, 0, 0, 200, 201, 1, 0,
		0, 0, 201, 202, 5, 37, 0, 0, 202, 35, 1, 0, 0, 0, 203, 208, 5, 34, 0, 0,
		204, 205, 5, 36, 0, 0, 205, 208, 3, 34, 17, 0, 206, 208, 3, 38, 19, 0,
		207, 203, 1, 0, 0, 0, 207, 204, 1, 0, 0, 0, 207, 206, 1, 0, 0, 0, 208,
		37, 1, 0, 0, 0, 209, 210, 5, 35, 0, 0, 210, 211, 5, 39, 0, 0, 211, 212,
		5, 38, 0, 0, 212, 39, 1, 0, 0, 0, 213, 217, 5, 16, 0, 0, 214, 216, 3, 42,
		21, 0, 215, 214, 1, 0, 0, 0, 216, 219, 1, 0, 0, 0, 217, 215, 1, 0, 0, 0,
		217, 218, 1, 0, 0, 0, 218, 220, 1, 0, 0, 0, 219, 217, 1, 0, 0, 0, 220,
		221, 5, 17, 0, 0, 221, 41, 1, 0, 0, 0, 222, 223, 5, 12, 0, 0, 223, 224,
		5, 12, 0, 0, 224, 228, 5, 19, 0, 0, 225, 226, 5, 12, 0, 0, 226, 228, 5,
		12, 0, 0, 227, 222, 1, 0, 0, 0, 227, 225, 1, 0, 0, 0, 228, 43, 1, 0, 0,
		0, 229, 233, 5, 4, 0, 0, 230, 232, 3, 46, 23, 0, 231, 230, 1, 0, 0, 0,
		232, 235, 1, 0, 0, 0, 233, 231, 1, 0, 0, 0, 233, 234, 1, 0, 0, 0, 234,
		45, 1, 0, 0, 0, 235, 233, 1, 0, 0, 0, 236, 237, 5, 12, 0, 0, 237, 239,
		3, 50, 25, 0, 238, 240, 3, 48, 24, 0, 239, 238, 1, 0, 0, 0, 239, 240, 1,
		0, 0, 0, 240, 242, 1, 0, 0, 0, 241, 243, 3, 18, 9, 0, 242, 241, 1, 0, 0,
		0, 242, 243, 1, 0, 0, 0, 243, 244, 1, 0, 0, 0, 244, 246, 3, 62, 31, 0,
		245, 247, 3, 20, 10, 0, 246, 245, 1, 0, 0, 0, 246, 247, 1, 0, 0, 0, 247,
		47, 1, 0, 0, 0, 248, 249, 5, 5, 0, 0, 249, 250, 3, 40, 20, 0, 250, 49,
		1, 0, 0, 0, 251, 254, 3, 40, 20, 0, 252, 254, 1, 0, 0, 0, 253, 251, 1,
		0, 0, 0, 253, 252, 1, 0, 0, 0, 254, 51, 1, 0, 0, 0, 255, 256, 5, 20, 0,
		0, 256, 258, 3, 54, 27, 0, 257, 259, 3, 58, 29, 0, 258, 257, 1, 0, 0, 0,
		258, 259, 1, 0, 0, 0, 259, 261, 1, 0, 0, 0, 260, 262, 3, 60, 30, 0, 261,
		260, 1, 0, 0, 0, 261, 262, 1, 0, 0, 0, 262, 263, 1, 0, 0, 0, 263, 264,
		5, 21, 0, 0, 264, 53, 1, 0, 0, 0, 265, 268, 5, 12, 0, 0, 266, 268, 3, 56,
		28, 0, 267, 265, 1, 0, 0, 0, 267, 266, 1, 0, 0, 0, 268, 55, 1, 0, 0, 0,
		269, 270, 5, 18, 0, 0, 270, 271, 5, 39, 0, 0, 271, 272, 5, 38, 0, 0, 272,
		57, 1, 0, 0, 0, 273, 274, 5, 23, 0, 0, 274, 275, 5, 12, 0, 0, 275, 59,
		1, 0, 0, 0, 276, 278, 5, 12, 0, 0, 277, 276, 1, 0, 0, 0, 278, 279, 1, 0,
		0, 0, 279, 277, 1, 0, 0, 0, 279, 280, 1, 0, 0, 0, 280, 61, 1, 0, 0, 0,
		281, 288, 3, 56, 28, 0, 282, 284, 3, 52, 26, 0, 283, 285, 3, 64, 32, 0,
		284, 283, 1, 0, 0, 0, 284, 285, 1, 0, 0, 0, 285, 288, 1, 0, 0, 0, 286,
		288, 3, 68, 34, 0, 287, 281, 1, 0, 0, 0, 287, 282, 1, 0, 0, 0, 287, 286,
		1, 0, 0, 0, 288, 63, 1, 0, 0, 0, 289, 292, 3, 66, 33, 0, 290, 292, 3, 68,
		34, 0, 291, 289, 1, 0, 0, 0, 291, 290, 1, 0, 0, 0, 292, 65, 1, 0, 0, 0,
		293, 296, 3, 70, 35, 0, 294, 296, 3, 24, 12, 0, 295, 293, 1, 0, 0, 0, 295,
		294, 1, 0, 0, 0, 296, 67, 1, 0, 0, 0, 297, 301, 5, 16, 0, 0, 298, 300,
		3, 62, 31, 0, 299, 298, 1, 0, 0, 0, 300, 303, 1, 0, 0, 0, 301, 299, 1,
		0, 0, 0, 301, 302, 1, 0, 0, 0, 302, 304, 1, 0, 0, 0, 303, 301, 1, 0, 0,
		0, 304, 305, 5, 17, 0, 0, 305, 69, 1, 0, 0, 0, 306, 307, 5, 12, 0, 0, 307,
		308, 5, 16, 0, 0, 308, 309, 3, 72, 36, 0, 309, 310, 5, 17, 0, 0, 310, 71,
		1, 0, 0, 0, 311, 316, 3, 74, 37, 0, 312, 313, 5, 19, 0, 0, 313, 315, 3,
		74, 37, 0, 314, 312, 1, 0, 0, 0, 315, 318, 1, 0, 0, 0, 316, 314, 1, 0,
		0, 0, 316, 317, 1, 0, 0, 0, 317, 320, 1, 0, 0, 0, 318, 316, 1, 0, 0, 0,
		319, 311, 1, 0, 0, 0, 319, 320, 1, 0, 0, 0, 320, 73, 1, 0, 0, 0, 321, 322,
		7, 0, 0, 0, 322, 75, 1, 0, 0, 0, 323, 327, 5, 11, 0, 0, 324, 326, 3, 78,
		39, 0, 325, 324, 1, 0, 0, 0, 326, 329, 1, 0, 0, 0, 327, 325, 1, 0, 0, 0,
		327, 328, 1, 0, 0, 0, 328, 77, 1, 0, 0, 0, 329, 327, 1, 0, 0, 0, 330, 331,
		3, 82, 41, 0, 331, 332, 5, 12, 0, 0, 332, 333, 3, 80, 40, 0, 333, 79, 1,
		0, 0, 0, 334, 335, 5, 21, 0, 0, 335, 337, 5, 21, 0, 0, 336, 334, 1, 0,
		0, 0, 336, 337, 1, 0, 0, 0, 337, 338, 1, 0, 0, 0, 338, 339, 3, 70, 35,
		0, 339, 81, 1, 0, 0, 0, 340, 341, 5, 23, 0, 0, 341, 344, 5, 12, 0, 0, 342,
		344, 5, 12, 0, 0, 343, 340, 1, 0, 0, 0, 343, 342, 1, 0, 0, 0, 344, 83,
		1, 0, 0, 0, 41, 86, 89, 92, 95, 98, 101, 104, 107, 119, 135, 145, 150,
		153, 156, 160, 176, 183, 188, 199, 207, 217, 227, 233, 239, 242, 246, 253,
		258, 261, 267, 279, 284, 287, 291, 295, 301, 316, 319, 327, 336, 343,
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

// wclInit initializes any static state used to implement wcl. By default the
// static state used to implement the parser is lazily initialized during the first call to
// Newwcl(). You can call this function if you wish to initialize the static state ahead
// of time.
func WclInit() {
	staticData := &wclParserStaticData
	staticData.once.Do(wclParserInit)
}

// Newwcl produces a new parser instance for the optional input antlr.TokenStream.
func Newwcl(input antlr.TokenStream) *wcl {
	WclInit()
	this := new(wcl)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &wclParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "java-escape"

	return this
}

// wcl tokens.
const (
	wclEOF                = antlr.TokenEOF
	wclText               = 1
	wclCSS                = 2
	wclImport             = 3
	wclDoc                = 4
	wclLocal              = 5
	wclGlobal             = 6
	wclExtern             = 7
	wclPre                = 8
	wclPost               = 9
	wclWcl                = 10
	wclEvent              = 11
	wclId                 = 12
	wclVersion            = 13
	wclLCurly             = 14
	wclRCurly             = 15
	wclLParen             = 16
	wclRParen             = 17
	wclDollar             = 18
	wclComma              = 19
	wclLessThan           = 20
	wclGreaterThan        = 21
	wclDot                = 22
	wclHash               = 23
	wclDash               = 24
	wclSemi               = 25
	wclPlus               = 26
	wclBackTick           = 27
	wclStringLit          = 28
	wclDoubleSlashComment = 29
	wclWhitespace         = 30
	wclContentRawText     = 31
	wclContentDollar      = 32
	wclContentBackTick    = 33
	wclUninterpRawText    = 34
	wclUninterpDollar     = 35
	wclUninterpLCurly     = 36
	wclUninterpRCurly     = 37
	wclVarRCurly          = 38
	wclVarId              = 39
)

// wcl rules.
const (
	wclRULE_program            = 0
	wclRULE_global             = 1
	wclRULE_extern             = 2
	wclRULE_wcl_section        = 3
	wclRULE_import_section     = 4
	wclRULE_css_section        = 5
	wclRULE_css_filespec       = 6
	wclRULE_text_section       = 7
	wclRULE_text_func          = 8
	wclRULE_pre_code           = 9
	wclRULE_post_code          = 10
	wclRULE_text_func_local    = 11
	wclRULE_text_top           = 12
	wclRULE_text_content       = 13
	wclRULE_text_content_inner = 14
	wclRULE_var_subs           = 15
	wclRULE_sub                = 16
	wclRULE_uninterp           = 17
	wclRULE_uninterp_inner     = 18
	wclRULE_uninterp_var       = 19
	wclRULE_param_spec         = 20
	wclRULE_param_pair         = 21
	wclRULE_doc_section        = 22
	wclRULE_doc_func           = 23
	wclRULE_doc_func_local     = 24
	wclRULE_doc_func_formal    = 25
	wclRULE_doc_tag            = 26
	wclRULE_id_or_var_ref      = 27
	wclRULE_var_ref            = 28
	wclRULE_doc_id             = 29
	wclRULE_doc_class          = 30
	wclRULE_doc_elem           = 31
	wclRULE_doc_elem_content   = 32
	wclRULE_doc_elem_text      = 33
	wclRULE_doc_elem_child     = 34
	wclRULE_func_invoc         = 35
	wclRULE_func_actual_seq    = 36
	wclRULE_func_actual        = 37
	wclRULE_event_section      = 38
	wclRULE_event_spec         = 39
	wclRULE_event_call         = 40
	wclRULE_selector           = 41
)

// IProgramContext is an interface to support dynamic dispatch.
type IProgramContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetP returns the p attribute.
	GetP() *ProgramNode

	// SetP sets the p attribute.
	SetP(*ProgramNode)

	// IsProgramContext differentiates from other interfaces.
	IsProgramContext()
}

type ProgramContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	p      *ProgramNode
}

func NewEmptyProgramContext() *ProgramContext {
	var p = new(ProgramContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_program
	return p
}

func (*ProgramContext) IsProgramContext() {}

func NewProgramContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgramContext {
	var p = new(ProgramContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_program

	return p
}

func (s *ProgramContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgramContext) GetP() *ProgramNode { return s.p }

func (s *ProgramContext) SetP(v *ProgramNode) { s.p = v }

func (s *ProgramContext) Wcl_section() IWcl_sectionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWcl_sectionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWcl_sectionContext)
}

func (s *ProgramContext) EOF() antlr.TerminalNode {
	return s.GetToken(wclEOF, 0)
}

func (s *ProgramContext) AllCss_section() []ICss_sectionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICss_sectionContext); ok {
			len++
		}
	}

	tst := make([]ICss_sectionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICss_sectionContext); ok {
			tst[i] = t.(ICss_sectionContext)
			i++
		}
	}

	return tst
}

func (s *ProgramContext) Css_section(i int) ICss_sectionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICss_sectionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICss_sectionContext)
}

func (s *ProgramContext) Import_section() IImport_sectionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImport_sectionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IImport_sectionContext)
}

func (s *ProgramContext) Extern() IExternContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExternContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExternContext)
}

func (s *ProgramContext) Global() IGlobalContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGlobalContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGlobalContext)
}

func (s *ProgramContext) Text_section() IText_sectionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IText_sectionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IText_sectionContext)
}

func (s *ProgramContext) Doc_section() IDoc_sectionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDoc_sectionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDoc_sectionContext)
}

func (s *ProgramContext) Event_section() IEvent_sectionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEvent_sectionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEvent_sectionContext)
}

func (s *ProgramContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProgramContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProgramContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterProgram(s)
	}
}

func (s *ProgramContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitProgram(s)
	}
}

func (s *ProgramContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitProgram(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Program() (localctx IProgramContext) {
	this := p
	_ = this

	localctx = NewProgramContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, wclRULE_program)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(84)
		p.Wcl_section()
	}
	p.SetState(86)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(85)
			p.Css_section()
		}

	}
	p.SetState(89)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclImport {
		{
			p.SetState(88)
			p.Import_section()
		}

	}
	p.SetState(92)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclExtern {
		{
			p.SetState(91)
			p.Extern()
		}

	}
	p.SetState(95)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclGlobal {
		{
			p.SetState(94)
			p.Global()
		}

	}
	p.SetState(98)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclText {
		{
			p.SetState(97)
			p.Text_section()
		}

	}
	p.SetState(101)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclCSS {
		{
			p.SetState(100)
			p.Css_section()
		}

	}
	p.SetState(104)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclDoc {
		{
			p.SetState(103)
			p.Doc_section()
		}

	}
	p.SetState(107)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclEvent {
		{
			p.SetState(106)
			p.Event_section()
		}

	}
	{
		p.SetState(109)
		p.Match(wclEOF)
	}

	return localctx
}

// IGlobalContext is an interface to support dynamic dispatch.
type IGlobalContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetG returns the g attribute.
	GetG() []*PFormal

	// SetG sets the g attribute.
	SetG([]*PFormal)

	// IsGlobalContext differentiates from other interfaces.
	IsGlobalContext()
}

type GlobalContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	g      []*PFormal
}

func NewEmptyGlobalContext() *GlobalContext {
	var p = new(GlobalContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_global
	return p
}

func (*GlobalContext) IsGlobalContext() {}

func NewGlobalContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GlobalContext {
	var p = new(GlobalContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_global

	return p
}

func (s *GlobalContext) GetParser() antlr.Parser { return s.parser }

func (s *GlobalContext) GetG() []*PFormal { return s.g }

func (s *GlobalContext) SetG(v []*PFormal) { s.g = v }

func (s *GlobalContext) Global() antlr.TerminalNode {
	return s.GetToken(wclGlobal, 0)
}

func (s *GlobalContext) Param_spec() IParam_specContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParam_specContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParam_specContext)
}

func (s *GlobalContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GlobalContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GlobalContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterGlobal(s)
	}
}

func (s *GlobalContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitGlobal(s)
	}
}

func (s *GlobalContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitGlobal(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Global() (localctx IGlobalContext) {
	this := p
	_ = this

	localctx = NewGlobalContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, wclRULE_global)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(111)
		p.Match(wclGlobal)
	}
	{
		p.SetState(112)
		p.Param_spec()
	}

	return localctx
}

// IExternContext is an interface to support dynamic dispatch.
type IExternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetE returns the e attribute.
	GetE() []string

	// SetE sets the e attribute.
	SetE([]string)

	// IsExternContext differentiates from other interfaces.
	IsExternContext()
}

type ExternContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	e      []string
}

func NewEmptyExternContext() *ExternContext {
	var p = new(ExternContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_extern
	return p
}

func (*ExternContext) IsExternContext() {}

func NewExternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternContext {
	var p = new(ExternContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_extern

	return p
}

func (s *ExternContext) GetParser() antlr.Parser { return s.parser }

func (s *ExternContext) GetE() []string { return s.e }

func (s *ExternContext) SetE(v []string) { s.e = v }

func (s *ExternContext) Extern() antlr.TerminalNode {
	return s.GetToken(wclExtern, 0)
}

func (s *ExternContext) LParen() antlr.TerminalNode {
	return s.GetToken(wclLParen, 0)
}

func (s *ExternContext) RParen() antlr.TerminalNode {
	return s.GetToken(wclRParen, 0)
}

func (s *ExternContext) AllId() []antlr.TerminalNode {
	return s.GetTokens(wclId)
}

func (s *ExternContext) Id(i int) antlr.TerminalNode {
	return s.GetToken(wclId, i)
}

func (s *ExternContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExternContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterExtern(s)
	}
}

func (s *ExternContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitExtern(s)
	}
}

func (s *ExternContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitExtern(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Extern() (localctx IExternContext) {
	this := p
	_ = this

	localctx = NewExternContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, wclRULE_extern)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(114)
		p.Match(wclExtern)
	}
	{
		p.SetState(115)
		p.Match(wclLParen)
	}
	p.SetState(119)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == wclId {
		{
			p.SetState(116)
			p.Match(wclId)
		}

		p.SetState(121)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(122)
		p.Match(wclRParen)
	}

	return localctx
}

// IWcl_sectionContext is an interface to support dynamic dispatch.
type IWcl_sectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWcl_sectionContext differentiates from other interfaces.
	IsWcl_sectionContext()
}

type Wcl_sectionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWcl_sectionContext() *Wcl_sectionContext {
	var p = new(Wcl_sectionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_wcl_section
	return p
}

func (*Wcl_sectionContext) IsWcl_sectionContext() {}

func NewWcl_sectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Wcl_sectionContext {
	var p = new(Wcl_sectionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_wcl_section

	return p
}

func (s *Wcl_sectionContext) GetParser() antlr.Parser { return s.parser }

func (s *Wcl_sectionContext) Wcl() antlr.TerminalNode {
	return s.GetToken(wclWcl, 0)
}

func (s *Wcl_sectionContext) Version() antlr.TerminalNode {
	return s.GetToken(wclVersion, 0)
}

func (s *Wcl_sectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Wcl_sectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Wcl_sectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterWcl_section(s)
	}
}

func (s *Wcl_sectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitWcl_section(s)
	}
}

func (s *Wcl_sectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitWcl_section(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Wcl_section() (localctx IWcl_sectionContext) {
	this := p
	_ = this

	localctx = NewWcl_sectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, wclRULE_wcl_section)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(124)
		p.Match(wclWcl)
	}
	{
		p.SetState(125)
		p.Match(wclVersion)
	}

	return localctx
}

// IImport_sectionContext is an interface to support dynamic dispatch.
type IImport_sectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetSection returns the section attribute.
	GetSection() *ImportSectionNode

	// SetSection sets the section attribute.
	SetSection(*ImportSectionNode)

	// IsImport_sectionContext differentiates from other interfaces.
	IsImport_sectionContext()
}

type Import_sectionContext struct {
	*antlr.BaseParserRuleContext
	parser  antlr.Parser
	section *ImportSectionNode
}

func NewEmptyImport_sectionContext() *Import_sectionContext {
	var p = new(Import_sectionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_import_section
	return p
}

func (*Import_sectionContext) IsImport_sectionContext() {}

func NewImport_sectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Import_sectionContext {
	var p = new(Import_sectionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_import_section

	return p
}

func (s *Import_sectionContext) GetParser() antlr.Parser { return s.parser }

func (s *Import_sectionContext) GetSection() *ImportSectionNode { return s.section }

func (s *Import_sectionContext) SetSection(v *ImportSectionNode) { s.section = v }

func (s *Import_sectionContext) Import() antlr.TerminalNode {
	return s.GetToken(wclImport, 0)
}

func (s *Import_sectionContext) LCurly() antlr.TerminalNode {
	return s.GetToken(wclLCurly, 0)
}

func (s *Import_sectionContext) Uninterp() IUninterpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUninterpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUninterpContext)
}

func (s *Import_sectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Import_sectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Import_sectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterImport_section(s)
	}
}

func (s *Import_sectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitImport_section(s)
	}
}

func (s *Import_sectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitImport_section(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Import_section() (localctx IImport_sectionContext) {
	this := p
	_ = this

	localctx = NewImport_sectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, wclRULE_import_section)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(127)
		p.Match(wclImport)
	}
	{
		p.SetState(128)
		p.Match(wclLCurly)
	}
	{
		p.SetState(129)
		p.Uninterp()
	}

	return localctx
}

// ICss_sectionContext is an interface to support dynamic dispatch.
type ICss_sectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCss_sectionContext differentiates from other interfaces.
	IsCss_sectionContext()
}

type Css_sectionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCss_sectionContext() *Css_sectionContext {
	var p = new(Css_sectionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_css_section
	return p
}

func (*Css_sectionContext) IsCss_sectionContext() {}

func NewCss_sectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Css_sectionContext {
	var p = new(Css_sectionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_css_section

	return p
}

func (s *Css_sectionContext) GetParser() antlr.Parser { return s.parser }

func (s *Css_sectionContext) CSS() antlr.TerminalNode {
	return s.GetToken(wclCSS, 0)
}

func (s *Css_sectionContext) AllCss_filespec() []ICss_filespecContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICss_filespecContext); ok {
			len++
		}
	}

	tst := make([]ICss_filespecContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICss_filespecContext); ok {
			tst[i] = t.(ICss_filespecContext)
			i++
		}
	}

	return tst
}

func (s *Css_sectionContext) Css_filespec(i int) ICss_filespecContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICss_filespecContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICss_filespecContext)
}

func (s *Css_sectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Css_sectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Css_sectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterCss_section(s)
	}
}

func (s *Css_sectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitCss_section(s)
	}
}

func (s *Css_sectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitCss_section(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Css_section() (localctx ICss_sectionContext) {
	this := p
	_ = this

	localctx = NewCss_sectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, wclRULE_css_section)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(131)
		p.Match(wclCSS)
	}
	p.SetState(135)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == wclPlus {
		{
			p.SetState(132)
			p.Css_filespec()
		}

		p.SetState(137)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ICss_filespecContext is an interface to support dynamic dispatch.
type ICss_filespecContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCss_filespecContext differentiates from other interfaces.
	IsCss_filespecContext()
}

type Css_filespecContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCss_filespecContext() *Css_filespecContext {
	var p = new(Css_filespecContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_css_filespec
	return p
}

func (*Css_filespecContext) IsCss_filespecContext() {}

func NewCss_filespecContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Css_filespecContext {
	var p = new(Css_filespecContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_css_filespec

	return p
}

func (s *Css_filespecContext) GetParser() antlr.Parser { return s.parser }

func (s *Css_filespecContext) Plus() antlr.TerminalNode {
	return s.GetToken(wclPlus, 0)
}

func (s *Css_filespecContext) StringLit() antlr.TerminalNode {
	return s.GetToken(wclStringLit, 0)
}

func (s *Css_filespecContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Css_filespecContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Css_filespecContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterCss_filespec(s)
	}
}

func (s *Css_filespecContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitCss_filespec(s)
	}
}

func (s *Css_filespecContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitCss_filespec(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Css_filespec() (localctx ICss_filespecContext) {
	this := p
	_ = this

	localctx = NewCss_filespecContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, wclRULE_css_filespec)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(138)
		p.Match(wclPlus)
	}
	{
		p.SetState(139)
		p.Match(wclStringLit)
	}

	return localctx
}

// IText_sectionContext is an interface to support dynamic dispatch.
type IText_sectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetSection returns the section attribute.
	GetSection() *TextSectionNode

	// SetSection sets the section attribute.
	SetSection(*TextSectionNode)

	// IsText_sectionContext differentiates from other interfaces.
	IsText_sectionContext()
}

type Text_sectionContext struct {
	*antlr.BaseParserRuleContext
	parser  antlr.Parser
	section *TextSectionNode
}

func NewEmptyText_sectionContext() *Text_sectionContext {
	var p = new(Text_sectionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_text_section
	return p
}

func (*Text_sectionContext) IsText_sectionContext() {}

func NewText_sectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Text_sectionContext {
	var p = new(Text_sectionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_text_section

	return p
}

func (s *Text_sectionContext) GetParser() antlr.Parser { return s.parser }

func (s *Text_sectionContext) GetSection() *TextSectionNode { return s.section }

func (s *Text_sectionContext) SetSection(v *TextSectionNode) { s.section = v }

func (s *Text_sectionContext) Text() antlr.TerminalNode {
	return s.GetToken(wclText, 0)
}

func (s *Text_sectionContext) AllText_func() []IText_funcContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IText_funcContext); ok {
			len++
		}
	}

	tst := make([]IText_funcContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IText_funcContext); ok {
			tst[i] = t.(IText_funcContext)
			i++
		}
	}

	return tst
}

func (s *Text_sectionContext) Text_func(i int) IText_funcContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IText_funcContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IText_funcContext)
}

func (s *Text_sectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Text_sectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Text_sectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterText_section(s)
	}
}

func (s *Text_sectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitText_section(s)
	}
}

func (s *Text_sectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitText_section(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Text_section() (localctx IText_sectionContext) {
	this := p
	_ = this

	localctx = NewText_sectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, wclRULE_text_section)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(141)
		p.Match(wclText)
	}
	p.SetState(145)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == wclId {
		{
			p.SetState(142)
			p.Text_func()
		}

		p.SetState(147)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IText_funcContext is an interface to support dynamic dispatch.
type IText_funcContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetI returns the i token.
	GetI() antlr.Token

	// SetI sets the i token.
	SetI(antlr.Token)

	// GetF returns the f attribute.
	GetF() *TextFuncNode

	// SetF sets the f attribute.
	SetF(*TextFuncNode)

	// IsText_funcContext differentiates from other interfaces.
	IsText_funcContext()
}

type Text_funcContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	f      *TextFuncNode
	i      antlr.Token
}

func NewEmptyText_funcContext() *Text_funcContext {
	var p = new(Text_funcContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_text_func
	return p
}

func (*Text_funcContext) IsText_funcContext() {}

func NewText_funcContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Text_funcContext {
	var p = new(Text_funcContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_text_func

	return p
}

func (s *Text_funcContext) GetParser() antlr.Parser { return s.parser }

func (s *Text_funcContext) GetI() antlr.Token { return s.i }

func (s *Text_funcContext) SetI(v antlr.Token) { s.i = v }

func (s *Text_funcContext) GetF() *TextFuncNode { return s.f }

func (s *Text_funcContext) SetF(v *TextFuncNode) { s.f = v }

func (s *Text_funcContext) Text_top() IText_topContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IText_topContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IText_topContext)
}

func (s *Text_funcContext) Id() antlr.TerminalNode {
	return s.GetToken(wclId, 0)
}

func (s *Text_funcContext) Param_spec() IParam_specContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParam_specContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParam_specContext)
}

func (s *Text_funcContext) Text_func_local() IText_func_localContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IText_func_localContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IText_func_localContext)
}

func (s *Text_funcContext) Pre_code() IPre_codeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPre_codeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPre_codeContext)
}

func (s *Text_funcContext) Post_code() IPost_codeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPost_codeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPost_codeContext)
}

func (s *Text_funcContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Text_funcContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Text_funcContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterText_func(s)
	}
}

func (s *Text_funcContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitText_func(s)
	}
}

func (s *Text_funcContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitText_func(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Text_func() (localctx IText_funcContext) {
	this := p
	_ = this

	localctx = NewText_funcContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, wclRULE_text_func)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(148)

		var _m = p.Match(wclId)

		localctx.(*Text_funcContext).i = _m
	}
	p.SetState(150)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclLParen {
		{
			p.SetState(149)
			p.Param_spec()
		}

	}
	p.SetState(153)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclLocal {
		{
			p.SetState(152)
			p.Text_func_local()
		}

	}
	p.SetState(156)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclPre {
		{
			p.SetState(155)
			p.Pre_code()
		}

	}
	{
		p.SetState(158)
		p.Text_top()
	}
	p.SetState(160)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclPost {
		{
			p.SetState(159)
			p.Post_code()
		}

	}

	return localctx
}

// IPre_codeContext is an interface to support dynamic dispatch.
type IPre_codeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetItem returns the item attribute.
	GetItem() []TextItem

	// SetItem sets the item attribute.
	SetItem([]TextItem)

	// IsPre_codeContext differentiates from other interfaces.
	IsPre_codeContext()
}

type Pre_codeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	item   []TextItem
}

func NewEmptyPre_codeContext() *Pre_codeContext {
	var p = new(Pre_codeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_pre_code
	return p
}

func (*Pre_codeContext) IsPre_codeContext() {}

func NewPre_codeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Pre_codeContext {
	var p = new(Pre_codeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_pre_code

	return p
}

func (s *Pre_codeContext) GetParser() antlr.Parser { return s.parser }

func (s *Pre_codeContext) GetItem() []TextItem { return s.item }

func (s *Pre_codeContext) SetItem(v []TextItem) { s.item = v }

func (s *Pre_codeContext) Pre() antlr.TerminalNode {
	return s.GetToken(wclPre, 0)
}

func (s *Pre_codeContext) LCurly() antlr.TerminalNode {
	return s.GetToken(wclLCurly, 0)
}

func (s *Pre_codeContext) Uninterp() IUninterpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUninterpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUninterpContext)
}

func (s *Pre_codeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Pre_codeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Pre_codeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterPre_code(s)
	}
}

func (s *Pre_codeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitPre_code(s)
	}
}

func (s *Pre_codeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitPre_code(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Pre_code() (localctx IPre_codeContext) {
	this := p
	_ = this

	localctx = NewPre_codeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, wclRULE_pre_code)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(162)
		p.Match(wclPre)
	}
	{
		p.SetState(163)
		p.Match(wclLCurly)
	}
	{
		p.SetState(164)
		p.Uninterp()
	}

	return localctx
}

// IPost_codeContext is an interface to support dynamic dispatch.
type IPost_codeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetItem returns the item attribute.
	GetItem() []TextItem

	// SetItem sets the item attribute.
	SetItem([]TextItem)

	// IsPost_codeContext differentiates from other interfaces.
	IsPost_codeContext()
}

type Post_codeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	item   []TextItem
}

func NewEmptyPost_codeContext() *Post_codeContext {
	var p = new(Post_codeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_post_code
	return p
}

func (*Post_codeContext) IsPost_codeContext() {}

func NewPost_codeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Post_codeContext {
	var p = new(Post_codeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_post_code

	return p
}

func (s *Post_codeContext) GetParser() antlr.Parser { return s.parser }

func (s *Post_codeContext) GetItem() []TextItem { return s.item }

func (s *Post_codeContext) SetItem(v []TextItem) { s.item = v }

func (s *Post_codeContext) Post() antlr.TerminalNode {
	return s.GetToken(wclPost, 0)
}

func (s *Post_codeContext) LCurly() antlr.TerminalNode {
	return s.GetToken(wclLCurly, 0)
}

func (s *Post_codeContext) Uninterp() IUninterpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUninterpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUninterpContext)
}

func (s *Post_codeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Post_codeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Post_codeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterPost_code(s)
	}
}

func (s *Post_codeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitPost_code(s)
	}
}

func (s *Post_codeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitPost_code(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Post_code() (localctx IPost_codeContext) {
	this := p
	_ = this

	localctx = NewPost_codeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, wclRULE_post_code)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(166)
		p.Match(wclPost)
	}
	{
		p.SetState(167)
		p.Match(wclLCurly)
	}
	{
		p.SetState(168)
		p.Uninterp()
	}

	return localctx
}

// IText_func_localContext is an interface to support dynamic dispatch.
type IText_func_localContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetFormal returns the formal attribute.
	GetFormal() []*PFormal

	// SetFormal sets the formal attribute.
	SetFormal([]*PFormal)

	// IsText_func_localContext differentiates from other interfaces.
	IsText_func_localContext()
}

type Text_func_localContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	formal []*PFormal
}

func NewEmptyText_func_localContext() *Text_func_localContext {
	var p = new(Text_func_localContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_text_func_local
	return p
}

func (*Text_func_localContext) IsText_func_localContext() {}

func NewText_func_localContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Text_func_localContext {
	var p = new(Text_func_localContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_text_func_local

	return p
}

func (s *Text_func_localContext) GetParser() antlr.Parser { return s.parser }

func (s *Text_func_localContext) GetFormal() []*PFormal { return s.formal }

func (s *Text_func_localContext) SetFormal(v []*PFormal) { s.formal = v }

func (s *Text_func_localContext) Local() antlr.TerminalNode {
	return s.GetToken(wclLocal, 0)
}

func (s *Text_func_localContext) Param_spec() IParam_specContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParam_specContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParam_specContext)
}

func (s *Text_func_localContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Text_func_localContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Text_func_localContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterText_func_local(s)
	}
}

func (s *Text_func_localContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitText_func_local(s)
	}
}

func (s *Text_func_localContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitText_func_local(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Text_func_local() (localctx IText_func_localContext) {
	this := p
	_ = this

	localctx = NewText_func_localContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, wclRULE_text_func_local)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(170)
		p.Match(wclLocal)
	}
	{
		p.SetState(171)
		p.Param_spec()
	}

	return localctx
}

// IText_topContext is an interface to support dynamic dispatch.
type IText_topContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetItem returns the item attribute.
	GetItem() []TextItem

	// SetItem sets the item attribute.
	SetItem([]TextItem)

	// IsText_topContext differentiates from other interfaces.
	IsText_topContext()
}

type Text_topContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	item   []TextItem
}

func NewEmptyText_topContext() *Text_topContext {
	var p = new(Text_topContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_text_top
	return p
}

func (*Text_topContext) IsText_topContext() {}

func NewText_topContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Text_topContext {
	var p = new(Text_topContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_text_top

	return p
}

func (s *Text_topContext) GetParser() antlr.Parser { return s.parser }

func (s *Text_topContext) GetItem() []TextItem { return s.item }

func (s *Text_topContext) SetItem(v []TextItem) { s.item = v }

func (s *Text_topContext) BackTick() antlr.TerminalNode {
	return s.GetToken(wclBackTick, 0)
}

func (s *Text_topContext) ContentBackTick() antlr.TerminalNode {
	return s.GetToken(wclContentBackTick, 0)
}

func (s *Text_topContext) Text_content() IText_contentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IText_contentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IText_contentContext)
}

func (s *Text_topContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Text_topContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Text_topContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterText_top(s)
	}
}

func (s *Text_topContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitText_top(s)
	}
}

func (s *Text_topContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitText_top(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Text_top() (localctx IText_topContext) {
	this := p
	_ = this

	localctx = NewText_topContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, wclRULE_text_top)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(173)
		p.Match(wclBackTick)
	}
	p.SetState(176)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 15, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(174)
			p.Text_content()
		}

	case 2:

	}
	{
		p.SetState(178)
		p.Match(wclContentBackTick)
	}

	return localctx
}

// IText_contentContext is an interface to support dynamic dispatch.
type IText_contentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetItem returns the item attribute.
	GetItem() []TextItem

	// SetItem sets the item attribute.
	SetItem([]TextItem)

	// IsText_contentContext differentiates from other interfaces.
	IsText_contentContext()
}

type Text_contentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	item   []TextItem
}

func NewEmptyText_contentContext() *Text_contentContext {
	var p = new(Text_contentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_text_content
	return p
}

func (*Text_contentContext) IsText_contentContext() {}

func NewText_contentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Text_contentContext {
	var p = new(Text_contentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_text_content

	return p
}

func (s *Text_contentContext) GetParser() antlr.Parser { return s.parser }

func (s *Text_contentContext) GetItem() []TextItem { return s.item }

func (s *Text_contentContext) SetItem(v []TextItem) { s.item = v }

func (s *Text_contentContext) AllText_content_inner() []IText_content_innerContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IText_content_innerContext); ok {
			len++
		}
	}

	tst := make([]IText_content_innerContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IText_content_innerContext); ok {
			tst[i] = t.(IText_content_innerContext)
			i++
		}
	}

	return tst
}

func (s *Text_contentContext) Text_content_inner(i int) IText_content_innerContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IText_content_innerContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IText_content_innerContext)
}

func (s *Text_contentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Text_contentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Text_contentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterText_content(s)
	}
}

func (s *Text_contentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitText_content(s)
	}
}

func (s *Text_contentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitText_content(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Text_content() (localctx IText_contentContext) {
	this := p
	_ = this

	localctx = NewText_contentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, wclRULE_text_content)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(183)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == wclContentRawText || _la == wclContentDollar {
		{
			p.SetState(180)
			p.Text_content_inner()
		}

		p.SetState(185)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IText_content_innerContext is an interface to support dynamic dispatch.
type IText_content_innerContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetItem returns the item attribute.
	GetItem() []TextItem

	// SetItem sets the item attribute.
	SetItem([]TextItem)

	// IsText_content_innerContext differentiates from other interfaces.
	IsText_content_innerContext()
}

type Text_content_innerContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	item   []TextItem
}

func NewEmptyText_content_innerContext() *Text_content_innerContext {
	var p = new(Text_content_innerContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_text_content_inner
	return p
}

func (*Text_content_innerContext) IsText_content_innerContext() {}

func NewText_content_innerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Text_content_innerContext {
	var p = new(Text_content_innerContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_text_content_inner

	return p
}

func (s *Text_content_innerContext) GetParser() antlr.Parser { return s.parser }

func (s *Text_content_innerContext) GetItem() []TextItem { return s.item }

func (s *Text_content_innerContext) SetItem(v []TextItem) { s.item = v }

func (s *Text_content_innerContext) CopyFrom(ctx *Text_content_innerContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
	s.item = ctx.item
}

func (s *Text_content_innerContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Text_content_innerContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type VarSubContext struct {
	*Text_content_innerContext
}

func NewVarSubContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *VarSubContext {
	var p = new(VarSubContext)

	p.Text_content_innerContext = NewEmptyText_content_innerContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Text_content_innerContext))

	return p
}

func (s *VarSubContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VarSubContext) Var_subs() IVar_subsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVar_subsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVar_subsContext)
}

func (s *VarSubContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterVarSub(s)
	}
}

func (s *VarSubContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitVarSub(s)
	}
}

func (s *VarSubContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitVarSub(s)

	default:
		return t.VisitChildren(s)
	}
}

type RawTextContext struct {
	*Text_content_innerContext
}

func NewRawTextContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *RawTextContext {
	var p = new(RawTextContext)

	p.Text_content_innerContext = NewEmptyText_content_innerContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Text_content_innerContext))

	return p
}

func (s *RawTextContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RawTextContext) ContentRawText() antlr.TerminalNode {
	return s.GetToken(wclContentRawText, 0)
}

func (s *RawTextContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterRawText(s)
	}
}

func (s *RawTextContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitRawText(s)
	}
}

func (s *RawTextContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitRawText(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Text_content_inner() (localctx IText_content_innerContext) {
	this := p
	_ = this

	localctx = NewText_content_innerContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, wclRULE_text_content_inner)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(188)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case wclContentRawText:
		localctx = NewRawTextContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(186)
			p.Match(wclContentRawText)
		}

	case wclContentDollar:
		localctx = NewVarSubContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(187)
			p.Var_subs()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IVar_subsContext is an interface to support dynamic dispatch.
type IVar_subsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetItem returns the item attribute.
	GetItem() []TextItem

	// SetItem sets the item attribute.
	SetItem([]TextItem)

	// IsVar_subsContext differentiates from other interfaces.
	IsVar_subsContext()
}

type Var_subsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	item   []TextItem
}

func NewEmptyVar_subsContext() *Var_subsContext {
	var p = new(Var_subsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_var_subs
	return p
}

func (*Var_subsContext) IsVar_subsContext() {}

func NewVar_subsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Var_subsContext {
	var p = new(Var_subsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_var_subs

	return p
}

func (s *Var_subsContext) GetParser() antlr.Parser { return s.parser }

func (s *Var_subsContext) GetItem() []TextItem { return s.item }

func (s *Var_subsContext) SetItem(v []TextItem) { s.item = v }

func (s *Var_subsContext) ContentDollar() antlr.TerminalNode {
	return s.GetToken(wclContentDollar, 0)
}

func (s *Var_subsContext) Sub() ISubContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISubContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISubContext)
}

func (s *Var_subsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Var_subsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Var_subsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterVar_subs(s)
	}
}

func (s *Var_subsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitVar_subs(s)
	}
}

func (s *Var_subsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitVar_subs(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Var_subs() (localctx IVar_subsContext) {
	this := p
	_ = this

	localctx = NewVar_subsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, wclRULE_var_subs)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(190)
		p.Match(wclContentDollar)
	}
	{
		p.SetState(191)
		p.Sub()
	}

	return localctx
}

// ISubContext is an interface to support dynamic dispatch.
type ISubContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetItem returns the item attribute.
	GetItem() TextItem

	// SetItem sets the item attribute.
	SetItem(TextItem)

	// IsSubContext differentiates from other interfaces.
	IsSubContext()
}

type SubContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	item   TextItem
}

func NewEmptySubContext() *SubContext {
	var p = new(SubContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_sub
	return p
}

func (*SubContext) IsSubContext() {}

func NewSubContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubContext {
	var p = new(SubContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_sub

	return p
}

func (s *SubContext) GetParser() antlr.Parser { return s.parser }

func (s *SubContext) GetItem() TextItem { return s.item }

func (s *SubContext) SetItem(v TextItem) { s.item = v }

func (s *SubContext) VarId() antlr.TerminalNode {
	return s.GetToken(wclVarId, 0)
}

func (s *SubContext) VarRCurly() antlr.TerminalNode {
	return s.GetToken(wclVarRCurly, 0)
}

func (s *SubContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SubContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterSub(s)
	}
}

func (s *SubContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitSub(s)
	}
}

func (s *SubContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitSub(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Sub() (localctx ISubContext) {
	this := p
	_ = this

	localctx = NewSubContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, wclRULE_sub)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(193)
		p.Match(wclVarId)
	}
	{
		p.SetState(194)
		p.Match(wclVarRCurly)
	}

	return localctx
}

// IUninterpContext is an interface to support dynamic dispatch.
type IUninterpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetItem returns the item attribute.
	GetItem() []TextItem

	// SetItem sets the item attribute.
	SetItem([]TextItem)

	// IsUninterpContext differentiates from other interfaces.
	IsUninterpContext()
}

type UninterpContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	item   []TextItem
}

func NewEmptyUninterpContext() *UninterpContext {
	var p = new(UninterpContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_uninterp
	return p
}

func (*UninterpContext) IsUninterpContext() {}

func NewUninterpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UninterpContext {
	var p = new(UninterpContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_uninterp

	return p
}

func (s *UninterpContext) GetParser() antlr.Parser { return s.parser }

func (s *UninterpContext) GetItem() []TextItem { return s.item }

func (s *UninterpContext) SetItem(v []TextItem) { s.item = v }

func (s *UninterpContext) UninterpRCurly() antlr.TerminalNode {
	return s.GetToken(wclUninterpRCurly, 0)
}

func (s *UninterpContext) AllUninterp_inner() []IUninterp_innerContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IUninterp_innerContext); ok {
			len++
		}
	}

	tst := make([]IUninterp_innerContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IUninterp_innerContext); ok {
			tst[i] = t.(IUninterp_innerContext)
			i++
		}
	}

	return tst
}

func (s *UninterpContext) Uninterp_inner(i int) IUninterp_innerContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUninterp_innerContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUninterp_innerContext)
}

func (s *UninterpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UninterpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UninterpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterUninterp(s)
	}
}

func (s *UninterpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitUninterp(s)
	}
}

func (s *UninterpContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitUninterp(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Uninterp() (localctx IUninterpContext) {
	this := p
	_ = this

	localctx = NewUninterpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, wclRULE_uninterp)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(197)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&120259084288) != 0 {
		{
			p.SetState(196)
			p.Uninterp_inner()
		}

		p.SetState(199)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(201)
		p.Match(wclUninterpRCurly)
	}

	return localctx
}

// IUninterp_innerContext is an interface to support dynamic dispatch.
type IUninterp_innerContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetItem returns the Item attribute.
	GetItem() []TextItem

	// SetItem sets the Item attribute.
	SetItem([]TextItem)

	// IsUninterp_innerContext differentiates from other interfaces.
	IsUninterp_innerContext()
}

type Uninterp_innerContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	Item   []TextItem
}

func NewEmptyUninterp_innerContext() *Uninterp_innerContext {
	var p = new(Uninterp_innerContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_uninterp_inner
	return p
}

func (*Uninterp_innerContext) IsUninterp_innerContext() {}

func NewUninterp_innerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Uninterp_innerContext {
	var p = new(Uninterp_innerContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_uninterp_inner

	return p
}

func (s *Uninterp_innerContext) GetParser() antlr.Parser { return s.parser }

func (s *Uninterp_innerContext) GetItem() []TextItem { return s.Item }

func (s *Uninterp_innerContext) SetItem(v []TextItem) { s.Item = v }

func (s *Uninterp_innerContext) CopyFrom(ctx *Uninterp_innerContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
	s.Item = ctx.Item
}

func (s *Uninterp_innerContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Uninterp_innerContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type UninterpRawTextContext struct {
	*Uninterp_innerContext
}

func NewUninterpRawTextContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *UninterpRawTextContext {
	var p = new(UninterpRawTextContext)

	p.Uninterp_innerContext = NewEmptyUninterp_innerContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Uninterp_innerContext))

	return p
}

func (s *UninterpRawTextContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UninterpRawTextContext) UninterpRawText() antlr.TerminalNode {
	return s.GetToken(wclUninterpRawText, 0)
}

func (s *UninterpRawTextContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterUninterpRawText(s)
	}
}

func (s *UninterpRawTextContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitUninterpRawText(s)
	}
}

func (s *UninterpRawTextContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitUninterpRawText(s)

	default:
		return t.VisitChildren(s)
	}
}

type UninterpNestedContext struct {
	*Uninterp_innerContext
}

func NewUninterpNestedContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *UninterpNestedContext {
	var p = new(UninterpNestedContext)

	p.Uninterp_innerContext = NewEmptyUninterp_innerContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Uninterp_innerContext))

	return p
}

func (s *UninterpNestedContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UninterpNestedContext) UninterpLCurly() antlr.TerminalNode {
	return s.GetToken(wclUninterpLCurly, 0)
}

func (s *UninterpNestedContext) Uninterp() IUninterpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUninterpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUninterpContext)
}

func (s *UninterpNestedContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterUninterpNested(s)
	}
}

func (s *UninterpNestedContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitUninterpNested(s)
	}
}

func (s *UninterpNestedContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitUninterpNested(s)

	default:
		return t.VisitChildren(s)
	}
}

type UninterpVarContext struct {
	*Uninterp_innerContext
}

func NewUninterpVarContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *UninterpVarContext {
	var p = new(UninterpVarContext)

	p.Uninterp_innerContext = NewEmptyUninterp_innerContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Uninterp_innerContext))

	return p
}

func (s *UninterpVarContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UninterpVarContext) Uninterp_var() IUninterp_varContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUninterp_varContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUninterp_varContext)
}

func (s *UninterpVarContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterUninterpVar(s)
	}
}

func (s *UninterpVarContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitUninterpVar(s)
	}
}

func (s *UninterpVarContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitUninterpVar(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Uninterp_inner() (localctx IUninterp_innerContext) {
	this := p
	_ = this

	localctx = NewUninterp_innerContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, wclRULE_uninterp_inner)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(207)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case wclUninterpRawText:
		localctx = NewUninterpRawTextContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(203)
			p.Match(wclUninterpRawText)
		}

	case wclUninterpLCurly:
		localctx = NewUninterpNestedContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(204)
			p.Match(wclUninterpLCurly)
		}
		{
			p.SetState(205)
			p.Uninterp()
		}

	case wclUninterpDollar:
		localctx = NewUninterpVarContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(206)
			p.Uninterp_var()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IUninterp_varContext is an interface to support dynamic dispatch.
type IUninterp_varContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetItem returns the item attribute.
	GetItem() []TextItem

	// SetItem sets the item attribute.
	SetItem([]TextItem)

	// IsUninterp_varContext differentiates from other interfaces.
	IsUninterp_varContext()
}

type Uninterp_varContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	item   []TextItem
}

func NewEmptyUninterp_varContext() *Uninterp_varContext {
	var p = new(Uninterp_varContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_uninterp_var
	return p
}

func (*Uninterp_varContext) IsUninterp_varContext() {}

func NewUninterp_varContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Uninterp_varContext {
	var p = new(Uninterp_varContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_uninterp_var

	return p
}

func (s *Uninterp_varContext) GetParser() antlr.Parser { return s.parser }

func (s *Uninterp_varContext) GetItem() []TextItem { return s.item }

func (s *Uninterp_varContext) SetItem(v []TextItem) { s.item = v }

func (s *Uninterp_varContext) UninterpDollar() antlr.TerminalNode {
	return s.GetToken(wclUninterpDollar, 0)
}

func (s *Uninterp_varContext) VarId() antlr.TerminalNode {
	return s.GetToken(wclVarId, 0)
}

func (s *Uninterp_varContext) VarRCurly() antlr.TerminalNode {
	return s.GetToken(wclVarRCurly, 0)
}

func (s *Uninterp_varContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Uninterp_varContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Uninterp_varContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterUninterp_var(s)
	}
}

func (s *Uninterp_varContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitUninterp_var(s)
	}
}

func (s *Uninterp_varContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitUninterp_var(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Uninterp_var() (localctx IUninterp_varContext) {
	this := p
	_ = this

	localctx = NewUninterp_varContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, wclRULE_uninterp_var)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(209)
		p.Match(wclUninterpDollar)
	}
	{
		p.SetState(210)
		p.Match(wclVarId)
	}
	{
		p.SetState(211)
		p.Match(wclVarRCurly)
	}

	return localctx
}

// IParam_specContext is an interface to support dynamic dispatch.
type IParam_specContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetFormal returns the formal attribute.
	GetFormal() []*PFormal

	// SetFormal sets the formal attribute.
	SetFormal([]*PFormal)

	// IsParam_specContext differentiates from other interfaces.
	IsParam_specContext()
}

type Param_specContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	formal []*PFormal
}

func NewEmptyParam_specContext() *Param_specContext {
	var p = new(Param_specContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_param_spec
	return p
}

func (*Param_specContext) IsParam_specContext() {}

func NewParam_specContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Param_specContext {
	var p = new(Param_specContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_param_spec

	return p
}

func (s *Param_specContext) GetParser() antlr.Parser { return s.parser }

func (s *Param_specContext) GetFormal() []*PFormal { return s.formal }

func (s *Param_specContext) SetFormal(v []*PFormal) { s.formal = v }

func (s *Param_specContext) LParen() antlr.TerminalNode {
	return s.GetToken(wclLParen, 0)
}

func (s *Param_specContext) RParen() antlr.TerminalNode {
	return s.GetToken(wclRParen, 0)
}

func (s *Param_specContext) AllParam_pair() []IParam_pairContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IParam_pairContext); ok {
			len++
		}
	}

	tst := make([]IParam_pairContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IParam_pairContext); ok {
			tst[i] = t.(IParam_pairContext)
			i++
		}
	}

	return tst
}

func (s *Param_specContext) Param_pair(i int) IParam_pairContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParam_pairContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParam_pairContext)
}

func (s *Param_specContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Param_specContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Param_specContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterParam_spec(s)
	}
}

func (s *Param_specContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitParam_spec(s)
	}
}

func (s *Param_specContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitParam_spec(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Param_spec() (localctx IParam_specContext) {
	this := p
	_ = this

	localctx = NewParam_specContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, wclRULE_param_spec)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(213)
		p.Match(wclLParen)
	}
	p.SetState(217)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == wclId {
		{
			p.SetState(214)
			p.Param_pair()
		}

		p.SetState(219)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(220)
		p.Match(wclRParen)
	}

	return localctx
}

// IParam_pairContext is an interface to support dynamic dispatch.
type IParam_pairContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetFormal returns the formal attribute.
	GetFormal() []*PFormal

	// SetFormal sets the formal attribute.
	SetFormal([]*PFormal)

	// IsParam_pairContext differentiates from other interfaces.
	IsParam_pairContext()
}

type Param_pairContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	formal []*PFormal
}

func NewEmptyParam_pairContext() *Param_pairContext {
	var p = new(Param_pairContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_param_pair
	return p
}

func (*Param_pairContext) IsParam_pairContext() {}

func NewParam_pairContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Param_pairContext {
	var p = new(Param_pairContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_param_pair

	return p
}

func (s *Param_pairContext) GetParser() antlr.Parser { return s.parser }

func (s *Param_pairContext) GetFormal() []*PFormal { return s.formal }

func (s *Param_pairContext) SetFormal(v []*PFormal) { s.formal = v }

func (s *Param_pairContext) CopyFrom(ctx *Param_pairContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
	s.formal = ctx.formal
}

func (s *Param_pairContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Param_pairContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type LastContext struct {
	*Param_pairContext
	n antlr.Token
	t antlr.Token
}

func NewLastContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LastContext {
	var p = new(LastContext)

	p.Param_pairContext = NewEmptyParam_pairContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Param_pairContext))

	return p
}

func (s *LastContext) GetN() antlr.Token { return s.n }

func (s *LastContext) GetT() antlr.Token { return s.t }

func (s *LastContext) SetN(v antlr.Token) { s.n = v }

func (s *LastContext) SetT(v antlr.Token) { s.t = v }

func (s *LastContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LastContext) AllId() []antlr.TerminalNode {
	return s.GetTokens(wclId)
}

func (s *LastContext) Id(i int) antlr.TerminalNode {
	return s.GetToken(wclId, i)
}

func (s *LastContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterLast(s)
	}
}

func (s *LastContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitLast(s)
	}
}

func (s *LastContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitLast(s)

	default:
		return t.VisitChildren(s)
	}
}

type PairContext struct {
	*Param_pairContext
	n antlr.Token
	t antlr.Token
}

func NewPairContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PairContext {
	var p = new(PairContext)

	p.Param_pairContext = NewEmptyParam_pairContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Param_pairContext))

	return p
}

func (s *PairContext) GetN() antlr.Token { return s.n }

func (s *PairContext) GetT() antlr.Token { return s.t }

func (s *PairContext) SetN(v antlr.Token) { s.n = v }

func (s *PairContext) SetT(v antlr.Token) { s.t = v }

func (s *PairContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PairContext) Comma() antlr.TerminalNode {
	return s.GetToken(wclComma, 0)
}

func (s *PairContext) AllId() []antlr.TerminalNode {
	return s.GetTokens(wclId)
}

func (s *PairContext) Id(i int) antlr.TerminalNode {
	return s.GetToken(wclId, i)
}

func (s *PairContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterPair(s)
	}
}

func (s *PairContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitPair(s)
	}
}

func (s *PairContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitPair(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Param_pair() (localctx IParam_pairContext) {
	this := p
	_ = this

	localctx = NewParam_pairContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, wclRULE_param_pair)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(227)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		localctx = NewPairContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(222)

			var _m = p.Match(wclId)

			localctx.(*PairContext).n = _m
		}
		{
			p.SetState(223)

			var _m = p.Match(wclId)

			localctx.(*PairContext).t = _m
		}
		{
			p.SetState(224)
			p.Match(wclComma)
		}

	case 2:
		localctx = NewLastContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(225)

			var _m = p.Match(wclId)

			localctx.(*LastContext).n = _m
		}
		{
			p.SetState(226)

			var _m = p.Match(wclId)

			localctx.(*LastContext).t = _m
		}

	}

	return localctx
}

// IDoc_sectionContext is an interface to support dynamic dispatch.
type IDoc_sectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetSection returns the section attribute.
	GetSection() *DocSectionNode

	// SetSection sets the section attribute.
	SetSection(*DocSectionNode)

	// IsDoc_sectionContext differentiates from other interfaces.
	IsDoc_sectionContext()
}

type Doc_sectionContext struct {
	*antlr.BaseParserRuleContext
	parser  antlr.Parser
	section *DocSectionNode
}

func NewEmptyDoc_sectionContext() *Doc_sectionContext {
	var p = new(Doc_sectionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_doc_section
	return p
}

func (*Doc_sectionContext) IsDoc_sectionContext() {}

func NewDoc_sectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Doc_sectionContext {
	var p = new(Doc_sectionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_doc_section

	return p
}

func (s *Doc_sectionContext) GetParser() antlr.Parser { return s.parser }

func (s *Doc_sectionContext) GetSection() *DocSectionNode { return s.section }

func (s *Doc_sectionContext) SetSection(v *DocSectionNode) { s.section = v }

func (s *Doc_sectionContext) Doc() antlr.TerminalNode {
	return s.GetToken(wclDoc, 0)
}

func (s *Doc_sectionContext) AllDoc_func() []IDoc_funcContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IDoc_funcContext); ok {
			len++
		}
	}

	tst := make([]IDoc_funcContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IDoc_funcContext); ok {
			tst[i] = t.(IDoc_funcContext)
			i++
		}
	}

	return tst
}

func (s *Doc_sectionContext) Doc_func(i int) IDoc_funcContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDoc_funcContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDoc_funcContext)
}

func (s *Doc_sectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Doc_sectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Doc_sectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterDoc_section(s)
	}
}

func (s *Doc_sectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitDoc_section(s)
	}
}

func (s *Doc_sectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitDoc_section(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Doc_section() (localctx IDoc_sectionContext) {
	this := p
	_ = this

	localctx = NewDoc_sectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, wclRULE_doc_section)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(229)
		p.Match(wclDoc)
	}
	p.SetState(233)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == wclId {
		{
			p.SetState(230)
			p.Doc_func()
		}

		p.SetState(235)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IDoc_funcContext is an interface to support dynamic dispatch.
type IDoc_funcContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetFn returns the fn attribute.
	GetFn() *DocFuncNode

	// SetFn sets the fn attribute.
	SetFn(*DocFuncNode)

	// IsDoc_funcContext differentiates from other interfaces.
	IsDoc_funcContext()
}

type Doc_funcContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	fn     *DocFuncNode
}

func NewEmptyDoc_funcContext() *Doc_funcContext {
	var p = new(Doc_funcContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_doc_func
	return p
}

func (*Doc_funcContext) IsDoc_funcContext() {}

func NewDoc_funcContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Doc_funcContext {
	var p = new(Doc_funcContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_doc_func

	return p
}

func (s *Doc_funcContext) GetParser() antlr.Parser { return s.parser }

func (s *Doc_funcContext) GetFn() *DocFuncNode { return s.fn }

func (s *Doc_funcContext) SetFn(v *DocFuncNode) { s.fn = v }

func (s *Doc_funcContext) Id() antlr.TerminalNode {
	return s.GetToken(wclId, 0)
}

func (s *Doc_funcContext) Doc_func_formal() IDoc_func_formalContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDoc_func_formalContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDoc_func_formalContext)
}

func (s *Doc_funcContext) Doc_elem() IDoc_elemContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDoc_elemContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDoc_elemContext)
}

func (s *Doc_funcContext) Doc_func_local() IDoc_func_localContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDoc_func_localContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDoc_func_localContext)
}

func (s *Doc_funcContext) Pre_code() IPre_codeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPre_codeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPre_codeContext)
}

func (s *Doc_funcContext) Post_code() IPost_codeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPost_codeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPost_codeContext)
}

func (s *Doc_funcContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Doc_funcContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Doc_funcContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterDoc_func(s)
	}
}

func (s *Doc_funcContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitDoc_func(s)
	}
}

func (s *Doc_funcContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitDoc_func(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Doc_func() (localctx IDoc_funcContext) {
	this := p
	_ = this

	localctx = NewDoc_funcContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, wclRULE_doc_func)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(236)
		p.Match(wclId)
	}
	{
		p.SetState(237)
		p.Doc_func_formal()
	}
	p.SetState(239)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclLocal {
		{
			p.SetState(238)
			p.Doc_func_local()
		}

	}
	p.SetState(242)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclPre {
		{
			p.SetState(241)
			p.Pre_code()
		}

	}
	{
		p.SetState(244)
		p.Doc_elem()
	}
	p.SetState(246)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclPost {
		{
			p.SetState(245)
			p.Post_code()
		}

	}

	return localctx
}

// IDoc_func_localContext is an interface to support dynamic dispatch.
type IDoc_func_localContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetFormal returns the formal attribute.
	GetFormal() []*PFormal

	// SetFormal sets the formal attribute.
	SetFormal([]*PFormal)

	// IsDoc_func_localContext differentiates from other interfaces.
	IsDoc_func_localContext()
}

type Doc_func_localContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	formal []*PFormal
}

func NewEmptyDoc_func_localContext() *Doc_func_localContext {
	var p = new(Doc_func_localContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_doc_func_local
	return p
}

func (*Doc_func_localContext) IsDoc_func_localContext() {}

func NewDoc_func_localContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Doc_func_localContext {
	var p = new(Doc_func_localContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_doc_func_local

	return p
}

func (s *Doc_func_localContext) GetParser() antlr.Parser { return s.parser }

func (s *Doc_func_localContext) GetFormal() []*PFormal { return s.formal }

func (s *Doc_func_localContext) SetFormal(v []*PFormal) { s.formal = v }

func (s *Doc_func_localContext) Local() antlr.TerminalNode {
	return s.GetToken(wclLocal, 0)
}

func (s *Doc_func_localContext) Param_spec() IParam_specContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParam_specContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParam_specContext)
}

func (s *Doc_func_localContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Doc_func_localContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Doc_func_localContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterDoc_func_local(s)
	}
}

func (s *Doc_func_localContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitDoc_func_local(s)
	}
}

func (s *Doc_func_localContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitDoc_func_local(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Doc_func_local() (localctx IDoc_func_localContext) {
	this := p
	_ = this

	localctx = NewDoc_func_localContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, wclRULE_doc_func_local)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(248)
		p.Match(wclLocal)
	}
	{
		p.SetState(249)
		p.Param_spec()
	}

	return localctx
}

// IDoc_func_formalContext is an interface to support dynamic dispatch.
type IDoc_func_formalContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetFormal returns the formal attribute.
	GetFormal() []*PFormal

	// SetFormal sets the formal attribute.
	SetFormal([]*PFormal)

	// IsDoc_func_formalContext differentiates from other interfaces.
	IsDoc_func_formalContext()
}

type Doc_func_formalContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	formal []*PFormal
}

func NewEmptyDoc_func_formalContext() *Doc_func_formalContext {
	var p = new(Doc_func_formalContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_doc_func_formal
	return p
}

func (*Doc_func_formalContext) IsDoc_func_formalContext() {}

func NewDoc_func_formalContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Doc_func_formalContext {
	var p = new(Doc_func_formalContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_doc_func_formal

	return p
}

func (s *Doc_func_formalContext) GetParser() antlr.Parser { return s.parser }

func (s *Doc_func_formalContext) GetFormal() []*PFormal { return s.formal }

func (s *Doc_func_formalContext) SetFormal(v []*PFormal) { s.formal = v }

func (s *Doc_func_formalContext) Param_spec() IParam_specContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParam_specContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParam_specContext)
}

func (s *Doc_func_formalContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Doc_func_formalContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Doc_func_formalContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterDoc_func_formal(s)
	}
}

func (s *Doc_func_formalContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitDoc_func_formal(s)
	}
}

func (s *Doc_func_formalContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitDoc_func_formal(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Doc_func_formal() (localctx IDoc_func_formalContext) {
	this := p
	_ = this

	localctx = NewDoc_func_formalContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, wclRULE_doc_func_formal)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(253)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(251)
			p.Param_spec()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)

	}

	return localctx
}

// IDoc_tagContext is an interface to support dynamic dispatch.
type IDoc_tagContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetTag returns the tag attribute.
	GetTag() *DocTag

	// SetTag sets the tag attribute.
	SetTag(*DocTag)

	// IsDoc_tagContext differentiates from other interfaces.
	IsDoc_tagContext()
}

type Doc_tagContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	tag    *DocTag
}

func NewEmptyDoc_tagContext() *Doc_tagContext {
	var p = new(Doc_tagContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_doc_tag
	return p
}

func (*Doc_tagContext) IsDoc_tagContext() {}

func NewDoc_tagContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Doc_tagContext {
	var p = new(Doc_tagContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_doc_tag

	return p
}

func (s *Doc_tagContext) GetParser() antlr.Parser { return s.parser }

func (s *Doc_tagContext) GetTag() *DocTag { return s.tag }

func (s *Doc_tagContext) SetTag(v *DocTag) { s.tag = v }

func (s *Doc_tagContext) LessThan() antlr.TerminalNode {
	return s.GetToken(wclLessThan, 0)
}

func (s *Doc_tagContext) Id_or_var_ref() IId_or_var_refContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IId_or_var_refContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IId_or_var_refContext)
}

func (s *Doc_tagContext) GreaterThan() antlr.TerminalNode {
	return s.GetToken(wclGreaterThan, 0)
}

func (s *Doc_tagContext) Doc_id() IDoc_idContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDoc_idContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDoc_idContext)
}

func (s *Doc_tagContext) Doc_class() IDoc_classContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDoc_classContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDoc_classContext)
}

func (s *Doc_tagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Doc_tagContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Doc_tagContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterDoc_tag(s)
	}
}

func (s *Doc_tagContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitDoc_tag(s)
	}
}

func (s *Doc_tagContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitDoc_tag(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Doc_tag() (localctx IDoc_tagContext) {
	this := p
	_ = this

	localctx = NewDoc_tagContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, wclRULE_doc_tag)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(255)
		p.Match(wclLessThan)
	}
	{
		p.SetState(256)
		p.Id_or_var_ref()
	}
	p.SetState(258)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclHash {
		{
			p.SetState(257)
			p.Doc_id()
		}

	}
	p.SetState(261)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclId {
		{
			p.SetState(260)
			p.Doc_class()
		}

	}
	{
		p.SetState(263)
		p.Match(wclGreaterThan)
	}

	return localctx
}

// IId_or_var_refContext is an interface to support dynamic dispatch.
type IId_or_var_refContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetIdVar returns the idVar attribute.
	GetIdVar() *DocIdOrVar

	// SetIdVar sets the idVar attribute.
	SetIdVar(*DocIdOrVar)

	// IsId_or_var_refContext differentiates from other interfaces.
	IsId_or_var_refContext()
}

type Id_or_var_refContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	idVar  *DocIdOrVar
}

func NewEmptyId_or_var_refContext() *Id_or_var_refContext {
	var p = new(Id_or_var_refContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_id_or_var_ref
	return p
}

func (*Id_or_var_refContext) IsId_or_var_refContext() {}

func NewId_or_var_refContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Id_or_var_refContext {
	var p = new(Id_or_var_refContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_id_or_var_ref

	return p
}

func (s *Id_or_var_refContext) GetParser() antlr.Parser { return s.parser }

func (s *Id_or_var_refContext) GetIdVar() *DocIdOrVar { return s.idVar }

func (s *Id_or_var_refContext) SetIdVar(v *DocIdOrVar) { s.idVar = v }

func (s *Id_or_var_refContext) Id() antlr.TerminalNode {
	return s.GetToken(wclId, 0)
}

func (s *Id_or_var_refContext) Var_ref() IVar_refContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVar_refContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVar_refContext)
}

func (s *Id_or_var_refContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Id_or_var_refContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Id_or_var_refContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterId_or_var_ref(s)
	}
}

func (s *Id_or_var_refContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitId_or_var_ref(s)
	}
}

func (s *Id_or_var_refContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitId_or_var_ref(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Id_or_var_ref() (localctx IId_or_var_refContext) {
	this := p
	_ = this

	localctx = NewId_or_var_refContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, wclRULE_id_or_var_ref)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(267)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case wclId:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(265)
			p.Match(wclId)
		}

	case wclDollar:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(266)
			p.Var_ref()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IVar_refContext is an interface to support dynamic dispatch.
type IVar_refContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetV returns the v attribute.
	GetV() *DocIdOrVar

	// SetV sets the v attribute.
	SetV(*DocIdOrVar)

	// IsVar_refContext differentiates from other interfaces.
	IsVar_refContext()
}

type Var_refContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	v      *DocIdOrVar
}

func NewEmptyVar_refContext() *Var_refContext {
	var p = new(Var_refContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_var_ref
	return p
}

func (*Var_refContext) IsVar_refContext() {}

func NewVar_refContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Var_refContext {
	var p = new(Var_refContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_var_ref

	return p
}

func (s *Var_refContext) GetParser() antlr.Parser { return s.parser }

func (s *Var_refContext) GetV() *DocIdOrVar { return s.v }

func (s *Var_refContext) SetV(v *DocIdOrVar) { s.v = v }

func (s *Var_refContext) Dollar() antlr.TerminalNode {
	return s.GetToken(wclDollar, 0)
}

func (s *Var_refContext) VarId() antlr.TerminalNode {
	return s.GetToken(wclVarId, 0)
}

func (s *Var_refContext) VarRCurly() antlr.TerminalNode {
	return s.GetToken(wclVarRCurly, 0)
}

func (s *Var_refContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Var_refContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Var_refContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterVar_ref(s)
	}
}

func (s *Var_refContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitVar_ref(s)
	}
}

func (s *Var_refContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitVar_ref(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Var_ref() (localctx IVar_refContext) {
	this := p
	_ = this

	localctx = NewVar_refContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, wclRULE_var_ref)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(269)
		p.Match(wclDollar)
	}
	{
		p.SetState(270)
		p.Match(wclVarId)
	}
	{
		p.SetState(271)
		p.Match(wclVarRCurly)
	}

	return localctx
}

// IDoc_idContext is an interface to support dynamic dispatch.
type IDoc_idContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetS returns the s attribute.
	GetS() string

	// SetS sets the s attribute.
	SetS(string)

	// IsDoc_idContext differentiates from other interfaces.
	IsDoc_idContext()
}

type Doc_idContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	s      string
}

func NewEmptyDoc_idContext() *Doc_idContext {
	var p = new(Doc_idContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_doc_id
	return p
}

func (*Doc_idContext) IsDoc_idContext() {}

func NewDoc_idContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Doc_idContext {
	var p = new(Doc_idContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_doc_id

	return p
}

func (s *Doc_idContext) GetParser() antlr.Parser { return s.parser }

func (s *Doc_idContext) GetS() string { return s.s }

func (s *Doc_idContext) SetS(v string) { s.s = v }

func (s *Doc_idContext) Hash() antlr.TerminalNode {
	return s.GetToken(wclHash, 0)
}

func (s *Doc_idContext) Id() antlr.TerminalNode {
	return s.GetToken(wclId, 0)
}

func (s *Doc_idContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Doc_idContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Doc_idContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterDoc_id(s)
	}
}

func (s *Doc_idContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitDoc_id(s)
	}
}

func (s *Doc_idContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitDoc_id(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Doc_id() (localctx IDoc_idContext) {
	this := p
	_ = this

	localctx = NewDoc_idContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, wclRULE_doc_id)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(273)
		p.Match(wclHash)
	}
	{
		p.SetState(274)
		p.Match(wclId)
	}

	return localctx
}

// IDoc_classContext is an interface to support dynamic dispatch.
type IDoc_classContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetClazz returns the clazz attribute.
	GetClazz() []string

	// SetClazz sets the clazz attribute.
	SetClazz([]string)

	// IsDoc_classContext differentiates from other interfaces.
	IsDoc_classContext()
}

type Doc_classContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	clazz  []string
}

func NewEmptyDoc_classContext() *Doc_classContext {
	var p = new(Doc_classContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_doc_class
	return p
}

func (*Doc_classContext) IsDoc_classContext() {}

func NewDoc_classContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Doc_classContext {
	var p = new(Doc_classContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_doc_class

	return p
}

func (s *Doc_classContext) GetParser() antlr.Parser { return s.parser }

func (s *Doc_classContext) GetClazz() []string { return s.clazz }

func (s *Doc_classContext) SetClazz(v []string) { s.clazz = v }

func (s *Doc_classContext) AllId() []antlr.TerminalNode {
	return s.GetTokens(wclId)
}

func (s *Doc_classContext) Id(i int) antlr.TerminalNode {
	return s.GetToken(wclId, i)
}

func (s *Doc_classContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Doc_classContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Doc_classContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterDoc_class(s)
	}
}

func (s *Doc_classContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitDoc_class(s)
	}
}

func (s *Doc_classContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitDoc_class(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Doc_class() (localctx IDoc_classContext) {
	this := p
	_ = this

	localctx = NewDoc_classContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, wclRULE_doc_class)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(277)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == wclId {
		{
			p.SetState(276)
			p.Match(wclId)
		}

		p.SetState(279)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IDoc_elemContext is an interface to support dynamic dispatch.
type IDoc_elemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetElem returns the elem attribute.
	GetElem() *DocElement

	// SetElem sets the elem attribute.
	SetElem(*DocElement)

	// IsDoc_elemContext differentiates from other interfaces.
	IsDoc_elemContext()
}

type Doc_elemContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	elem   *DocElement
}

func NewEmptyDoc_elemContext() *Doc_elemContext {
	var p = new(Doc_elemContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_doc_elem
	return p
}

func (*Doc_elemContext) IsDoc_elemContext() {}

func NewDoc_elemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Doc_elemContext {
	var p = new(Doc_elemContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_doc_elem

	return p
}

func (s *Doc_elemContext) GetParser() antlr.Parser { return s.parser }

func (s *Doc_elemContext) GetElem() *DocElement { return s.elem }

func (s *Doc_elemContext) SetElem(v *DocElement) { s.elem = v }

func (s *Doc_elemContext) CopyFrom(ctx *Doc_elemContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
	s.elem = ctx.elem
}

func (s *Doc_elemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Doc_elemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type HaveVarContext struct {
	*Doc_elemContext
}

func NewHaveVarContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *HaveVarContext {
	var p = new(HaveVarContext)

	p.Doc_elemContext = NewEmptyDoc_elemContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Doc_elemContext))

	return p
}

func (s *HaveVarContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HaveVarContext) Var_ref() IVar_refContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVar_refContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVar_refContext)
}

func (s *HaveVarContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterHaveVar(s)
	}
}

func (s *HaveVarContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitHaveVar(s)
	}
}

func (s *HaveVarContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitHaveVar(s)

	default:
		return t.VisitChildren(s)
	}
}

type HaveListContext struct {
	*Doc_elemContext
}

func NewHaveListContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *HaveListContext {
	var p = new(HaveListContext)

	p.Doc_elemContext = NewEmptyDoc_elemContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Doc_elemContext))

	return p
}

func (s *HaveListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HaveListContext) Doc_elem_child() IDoc_elem_childContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDoc_elem_childContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDoc_elem_childContext)
}

func (s *HaveListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterHaveList(s)
	}
}

func (s *HaveListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitHaveList(s)
	}
}

func (s *HaveListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitHaveList(s)

	default:
		return t.VisitChildren(s)
	}
}

type HaveTagContext struct {
	*Doc_elemContext
}

func NewHaveTagContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *HaveTagContext {
	var p = new(HaveTagContext)

	p.Doc_elemContext = NewEmptyDoc_elemContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Doc_elemContext))

	return p
}

func (s *HaveTagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HaveTagContext) Doc_tag() IDoc_tagContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDoc_tagContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDoc_tagContext)
}

func (s *HaveTagContext) Doc_elem_content() IDoc_elem_contentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDoc_elem_contentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDoc_elem_contentContext)
}

func (s *HaveTagContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterHaveTag(s)
	}
}

func (s *HaveTagContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitHaveTag(s)
	}
}

func (s *HaveTagContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitHaveTag(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Doc_elem() (localctx IDoc_elemContext) {
	this := p
	_ = this

	localctx = NewDoc_elemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, wclRULE_doc_elem)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(287)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case wclDollar:
		localctx = NewHaveVarContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(281)
			p.Var_ref()
		}

	case wclLessThan:
		localctx = NewHaveTagContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(282)
			p.Doc_tag()
		}
		p.SetState(284)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 31, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(283)
				p.Doc_elem_content()
			}

		}

	case wclLParen:
		localctx = NewHaveListContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(286)
			p.Doc_elem_child()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IDoc_elem_contentContext is an interface to support dynamic dispatch.
type IDoc_elem_contentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetElement returns the element attribute.
	GetElement() *DocElement

	// SetElement sets the element attribute.
	SetElement(*DocElement)

	// IsDoc_elem_contentContext differentiates from other interfaces.
	IsDoc_elem_contentContext()
}

type Doc_elem_contentContext struct {
	*antlr.BaseParserRuleContext
	parser  antlr.Parser
	element *DocElement
}

func NewEmptyDoc_elem_contentContext() *Doc_elem_contentContext {
	var p = new(Doc_elem_contentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_doc_elem_content
	return p
}

func (*Doc_elem_contentContext) IsDoc_elem_contentContext() {}

func NewDoc_elem_contentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Doc_elem_contentContext {
	var p = new(Doc_elem_contentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_doc_elem_content

	return p
}

func (s *Doc_elem_contentContext) GetParser() antlr.Parser { return s.parser }

func (s *Doc_elem_contentContext) GetElement() *DocElement { return s.element }

func (s *Doc_elem_contentContext) SetElement(v *DocElement) { s.element = v }

func (s *Doc_elem_contentContext) Doc_elem_text() IDoc_elem_textContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDoc_elem_textContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDoc_elem_textContext)
}

func (s *Doc_elem_contentContext) Doc_elem_child() IDoc_elem_childContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDoc_elem_childContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDoc_elem_childContext)
}

func (s *Doc_elem_contentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Doc_elem_contentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Doc_elem_contentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterDoc_elem_content(s)
	}
}

func (s *Doc_elem_contentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitDoc_elem_content(s)
	}
}

func (s *Doc_elem_contentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitDoc_elem_content(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Doc_elem_content() (localctx IDoc_elem_contentContext) {
	this := p
	_ = this

	localctx = NewDoc_elem_contentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, wclRULE_doc_elem_content)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(291)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case wclId, wclBackTick:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(289)
			p.Doc_elem_text()
		}

	case wclLParen:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(290)
			p.Doc_elem_child()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IDoc_elem_textContext is an interface to support dynamic dispatch.
type IDoc_elem_textContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetInvoc returns the invoc attribute.
	GetInvoc() *FuncInvoc

	// SetInvoc sets the invoc attribute.
	SetInvoc(*FuncInvoc)

	// IsDoc_elem_textContext differentiates from other interfaces.
	IsDoc_elem_textContext()
}

type Doc_elem_textContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	invoc  *FuncInvoc
}

func NewEmptyDoc_elem_textContext() *Doc_elem_textContext {
	var p = new(Doc_elem_textContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_doc_elem_text
	return p
}

func (*Doc_elem_textContext) IsDoc_elem_textContext() {}

func NewDoc_elem_textContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Doc_elem_textContext {
	var p = new(Doc_elem_textContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_doc_elem_text

	return p
}

func (s *Doc_elem_textContext) GetParser() antlr.Parser { return s.parser }

func (s *Doc_elem_textContext) GetInvoc() *FuncInvoc { return s.invoc }

func (s *Doc_elem_textContext) SetInvoc(v *FuncInvoc) { s.invoc = v }

func (s *Doc_elem_textContext) CopyFrom(ctx *Doc_elem_textContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
	s.invoc = ctx.invoc
}

func (s *Doc_elem_textContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Doc_elem_textContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type Doc_elem_text_anonContext struct {
	*Doc_elem_textContext
}

func NewDoc_elem_text_anonContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *Doc_elem_text_anonContext {
	var p = new(Doc_elem_text_anonContext)

	p.Doc_elem_textContext = NewEmptyDoc_elem_textContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Doc_elem_textContext))

	return p
}

func (s *Doc_elem_text_anonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Doc_elem_text_anonContext) Text_top() IText_topContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IText_topContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IText_topContext)
}

func (s *Doc_elem_text_anonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterDoc_elem_text_anon(s)
	}
}

func (s *Doc_elem_text_anonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitDoc_elem_text_anon(s)
	}
}

func (s *Doc_elem_text_anonContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitDoc_elem_text_anon(s)

	default:
		return t.VisitChildren(s)
	}
}

type Doc_elem_text_func_callContext struct {
	*Doc_elem_textContext
}

func NewDoc_elem_text_func_callContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *Doc_elem_text_func_callContext {
	var p = new(Doc_elem_text_func_callContext)

	p.Doc_elem_textContext = NewEmptyDoc_elem_textContext()
	p.parser = parser
	p.CopyFrom(ctx.(*Doc_elem_textContext))

	return p
}

func (s *Doc_elem_text_func_callContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Doc_elem_text_func_callContext) Func_invoc() IFunc_invocContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunc_invocContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunc_invocContext)
}

func (s *Doc_elem_text_func_callContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterDoc_elem_text_func_call(s)
	}
}

func (s *Doc_elem_text_func_callContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitDoc_elem_text_func_call(s)
	}
}

func (s *Doc_elem_text_func_callContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitDoc_elem_text_func_call(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Doc_elem_text() (localctx IDoc_elem_textContext) {
	this := p
	_ = this

	localctx = NewDoc_elem_textContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, wclRULE_doc_elem_text)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(295)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case wclId:
		localctx = NewDoc_elem_text_func_callContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(293)
			p.Func_invoc()
		}

	case wclBackTick:
		localctx = NewDoc_elem_text_anonContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(294)
			p.Text_top()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IDoc_elem_childContext is an interface to support dynamic dispatch.
type IDoc_elem_childContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetElem returns the elem attribute.
	GetElem() *DocElement

	// SetElem sets the elem attribute.
	SetElem(*DocElement)

	// IsDoc_elem_childContext differentiates from other interfaces.
	IsDoc_elem_childContext()
}

type Doc_elem_childContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	elem   *DocElement
}

func NewEmptyDoc_elem_childContext() *Doc_elem_childContext {
	var p = new(Doc_elem_childContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_doc_elem_child
	return p
}

func (*Doc_elem_childContext) IsDoc_elem_childContext() {}

func NewDoc_elem_childContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Doc_elem_childContext {
	var p = new(Doc_elem_childContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_doc_elem_child

	return p
}

func (s *Doc_elem_childContext) GetParser() antlr.Parser { return s.parser }

func (s *Doc_elem_childContext) GetElem() *DocElement { return s.elem }

func (s *Doc_elem_childContext) SetElem(v *DocElement) { s.elem = v }

func (s *Doc_elem_childContext) LParen() antlr.TerminalNode {
	return s.GetToken(wclLParen, 0)
}

func (s *Doc_elem_childContext) RParen() antlr.TerminalNode {
	return s.GetToken(wclRParen, 0)
}

func (s *Doc_elem_childContext) AllDoc_elem() []IDoc_elemContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IDoc_elemContext); ok {
			len++
		}
	}

	tst := make([]IDoc_elemContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IDoc_elemContext); ok {
			tst[i] = t.(IDoc_elemContext)
			i++
		}
	}

	return tst
}

func (s *Doc_elem_childContext) Doc_elem(i int) IDoc_elemContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDoc_elemContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDoc_elemContext)
}

func (s *Doc_elem_childContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Doc_elem_childContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Doc_elem_childContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterDoc_elem_child(s)
	}
}

func (s *Doc_elem_childContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitDoc_elem_child(s)
	}
}

func (s *Doc_elem_childContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitDoc_elem_child(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Doc_elem_child() (localctx IDoc_elem_childContext) {
	this := p
	_ = this

	localctx = NewDoc_elem_childContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, wclRULE_doc_elem_child)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(297)
		p.Match(wclLParen)
	}
	p.SetState(301)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1376256) != 0 {
		{
			p.SetState(298)
			p.Doc_elem()
		}

		p.SetState(303)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(304)
		p.Match(wclRParen)
	}

	return localctx
}

// IFunc_invocContext is an interface to support dynamic dispatch.
type IFunc_invocContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetInvoc returns the invoc attribute.
	GetInvoc() *FuncInvoc

	// SetInvoc sets the invoc attribute.
	SetInvoc(*FuncInvoc)

	// IsFunc_invocContext differentiates from other interfaces.
	IsFunc_invocContext()
}

type Func_invocContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	invoc  *FuncInvoc
}

func NewEmptyFunc_invocContext() *Func_invocContext {
	var p = new(Func_invocContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_func_invoc
	return p
}

func (*Func_invocContext) IsFunc_invocContext() {}

func NewFunc_invocContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Func_invocContext {
	var p = new(Func_invocContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_func_invoc

	return p
}

func (s *Func_invocContext) GetParser() antlr.Parser { return s.parser }

func (s *Func_invocContext) GetInvoc() *FuncInvoc { return s.invoc }

func (s *Func_invocContext) SetInvoc(v *FuncInvoc) { s.invoc = v }

func (s *Func_invocContext) Id() antlr.TerminalNode {
	return s.GetToken(wclId, 0)
}

func (s *Func_invocContext) LParen() antlr.TerminalNode {
	return s.GetToken(wclLParen, 0)
}

func (s *Func_invocContext) Func_actual_seq() IFunc_actual_seqContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunc_actual_seqContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunc_actual_seqContext)
}

func (s *Func_invocContext) RParen() antlr.TerminalNode {
	return s.GetToken(wclRParen, 0)
}

func (s *Func_invocContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Func_invocContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Func_invocContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterFunc_invoc(s)
	}
}

func (s *Func_invocContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitFunc_invoc(s)
	}
}

func (s *Func_invocContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitFunc_invoc(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Func_invoc() (localctx IFunc_invocContext) {
	this := p
	_ = this

	localctx = NewFunc_invocContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, wclRULE_func_invoc)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(306)
		p.Match(wclId)
	}
	{
		p.SetState(307)
		p.Match(wclLParen)
	}
	{
		p.SetState(308)
		p.Func_actual_seq()
	}
	{
		p.SetState(309)
		p.Match(wclRParen)
	}

	return localctx
}

// IFunc_actual_seqContext is an interface to support dynamic dispatch.
type IFunc_actual_seqContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetA returns the a rule contexts.
	GetA() IFunc_actualContext

	// GetB returns the b rule contexts.
	GetB() IFunc_actualContext

	// SetA sets the a rule contexts.
	SetA(IFunc_actualContext)

	// SetB sets the b rule contexts.
	SetB(IFunc_actualContext)

	// GetActual returns the actual attribute.
	GetActual() []*FuncActual

	// SetActual sets the actual attribute.
	SetActual([]*FuncActual)

	// IsFunc_actual_seqContext differentiates from other interfaces.
	IsFunc_actual_seqContext()
}

type Func_actual_seqContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	actual []*FuncActual
	a      IFunc_actualContext
	b      IFunc_actualContext
}

func NewEmptyFunc_actual_seqContext() *Func_actual_seqContext {
	var p = new(Func_actual_seqContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_func_actual_seq
	return p
}

func (*Func_actual_seqContext) IsFunc_actual_seqContext() {}

func NewFunc_actual_seqContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Func_actual_seqContext {
	var p = new(Func_actual_seqContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_func_actual_seq

	return p
}

func (s *Func_actual_seqContext) GetParser() antlr.Parser { return s.parser }

func (s *Func_actual_seqContext) GetA() IFunc_actualContext { return s.a }

func (s *Func_actual_seqContext) GetB() IFunc_actualContext { return s.b }

func (s *Func_actual_seqContext) SetA(v IFunc_actualContext) { s.a = v }

func (s *Func_actual_seqContext) SetB(v IFunc_actualContext) { s.b = v }

func (s *Func_actual_seqContext) GetActual() []*FuncActual { return s.actual }

func (s *Func_actual_seqContext) SetActual(v []*FuncActual) { s.actual = v }

func (s *Func_actual_seqContext) AllFunc_actual() []IFunc_actualContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFunc_actualContext); ok {
			len++
		}
	}

	tst := make([]IFunc_actualContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFunc_actualContext); ok {
			tst[i] = t.(IFunc_actualContext)
			i++
		}
	}

	return tst
}

func (s *Func_actual_seqContext) Func_actual(i int) IFunc_actualContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunc_actualContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunc_actualContext)
}

func (s *Func_actual_seqContext) AllComma() []antlr.TerminalNode {
	return s.GetTokens(wclComma)
}

func (s *Func_actual_seqContext) Comma(i int) antlr.TerminalNode {
	return s.GetToken(wclComma, i)
}

func (s *Func_actual_seqContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Func_actual_seqContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Func_actual_seqContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterFunc_actual_seq(s)
	}
}

func (s *Func_actual_seqContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitFunc_actual_seq(s)
	}
}

func (s *Func_actual_seqContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitFunc_actual_seq(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Func_actual_seq() (localctx IFunc_actual_seqContext) {
	this := p
	_ = this

	localctx = NewFunc_actual_seqContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, wclRULE_func_actual_seq)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(319)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclId || _la == wclStringLit {
		{
			p.SetState(311)

			var _x = p.Func_actual()

			localctx.(*Func_actual_seqContext).a = _x
		}
		p.SetState(316)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == wclComma {
			{
				p.SetState(312)
				p.Match(wclComma)
			}
			{
				p.SetState(313)

				var _x = p.Func_actual()

				localctx.(*Func_actual_seqContext).b = _x
			}

			p.SetState(318)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	}

	return localctx
}

// IFunc_actualContext is an interface to support dynamic dispatch.
type IFunc_actualContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetActual returns the actual attribute.
	GetActual() *FuncActual

	// SetActual sets the actual attribute.
	SetActual(*FuncActual)

	// IsFunc_actualContext differentiates from other interfaces.
	IsFunc_actualContext()
}

type Func_actualContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	actual *FuncActual
}

func NewEmptyFunc_actualContext() *Func_actualContext {
	var p = new(Func_actualContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_func_actual
	return p
}

func (*Func_actualContext) IsFunc_actualContext() {}

func NewFunc_actualContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Func_actualContext {
	var p = new(Func_actualContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_func_actual

	return p
}

func (s *Func_actualContext) GetParser() antlr.Parser { return s.parser }

func (s *Func_actualContext) GetActual() *FuncActual { return s.actual }

func (s *Func_actualContext) SetActual(v *FuncActual) { s.actual = v }

func (s *Func_actualContext) Id() antlr.TerminalNode {
	return s.GetToken(wclId, 0)
}

func (s *Func_actualContext) StringLit() antlr.TerminalNode {
	return s.GetToken(wclStringLit, 0)
}

func (s *Func_actualContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Func_actualContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Func_actualContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterFunc_actual(s)
	}
}

func (s *Func_actualContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitFunc_actual(s)
	}
}

func (s *Func_actualContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitFunc_actual(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Func_actual() (localctx IFunc_actualContext) {
	this := p
	_ = this

	localctx = NewFunc_actualContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, wclRULE_func_actual)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(321)
		_la = p.GetTokenStream().LA(1)

		if !(_la == wclId || _la == wclStringLit) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IEvent_sectionContext is an interface to support dynamic dispatch.
type IEvent_sectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEvent_sectionContext differentiates from other interfaces.
	IsEvent_sectionContext()
}

type Event_sectionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEvent_sectionContext() *Event_sectionContext {
	var p = new(Event_sectionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_event_section
	return p
}

func (*Event_sectionContext) IsEvent_sectionContext() {}

func NewEvent_sectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Event_sectionContext {
	var p = new(Event_sectionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_event_section

	return p
}

func (s *Event_sectionContext) GetParser() antlr.Parser { return s.parser }

func (s *Event_sectionContext) Event() antlr.TerminalNode {
	return s.GetToken(wclEvent, 0)
}

func (s *Event_sectionContext) AllEvent_spec() []IEvent_specContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IEvent_specContext); ok {
			len++
		}
	}

	tst := make([]IEvent_specContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IEvent_specContext); ok {
			tst[i] = t.(IEvent_specContext)
			i++
		}
	}

	return tst
}

func (s *Event_sectionContext) Event_spec(i int) IEvent_specContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEvent_specContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEvent_specContext)
}

func (s *Event_sectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Event_sectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Event_sectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterEvent_section(s)
	}
}

func (s *Event_sectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitEvent_section(s)
	}
}

func (s *Event_sectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitEvent_section(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Event_section() (localctx IEvent_sectionContext) {
	this := p
	_ = this

	localctx = NewEvent_sectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, wclRULE_event_section)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(323)
		p.Match(wclEvent)
	}
	p.SetState(327)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == wclId || _la == wclHash {
		{
			p.SetState(324)
			p.Event_spec()
		}

		p.SetState(329)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEvent_specContext is an interface to support dynamic dispatch.
type IEvent_specContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEvent_specContext differentiates from other interfaces.
	IsEvent_specContext()
}

type Event_specContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEvent_specContext() *Event_specContext {
	var p = new(Event_specContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_event_spec
	return p
}

func (*Event_specContext) IsEvent_specContext() {}

func NewEvent_specContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Event_specContext {
	var p = new(Event_specContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_event_spec

	return p
}

func (s *Event_specContext) GetParser() antlr.Parser { return s.parser }

func (s *Event_specContext) Selector() ISelectorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectorContext)
}

func (s *Event_specContext) Id() antlr.TerminalNode {
	return s.GetToken(wclId, 0)
}

func (s *Event_specContext) Event_call() IEvent_callContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEvent_callContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEvent_callContext)
}

func (s *Event_specContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Event_specContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Event_specContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterEvent_spec(s)
	}
}

func (s *Event_specContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitEvent_spec(s)
	}
}

func (s *Event_specContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitEvent_spec(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Event_spec() (localctx IEvent_specContext) {
	this := p
	_ = this

	localctx = NewEvent_specContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, wclRULE_event_spec)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(330)
		p.Selector()
	}
	{
		p.SetState(331)
		p.Match(wclId)
	}
	{
		p.SetState(332)
		p.Event_call()
	}

	return localctx
}

// IEvent_callContext is an interface to support dynamic dispatch.
type IEvent_callContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEvent_callContext differentiates from other interfaces.
	IsEvent_callContext()
}

type Event_callContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEvent_callContext() *Event_callContext {
	var p = new(Event_callContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_event_call
	return p
}

func (*Event_callContext) IsEvent_callContext() {}

func NewEvent_callContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Event_callContext {
	var p = new(Event_callContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_event_call

	return p
}

func (s *Event_callContext) GetParser() antlr.Parser { return s.parser }

func (s *Event_callContext) Func_invoc() IFunc_invocContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunc_invocContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunc_invocContext)
}

func (s *Event_callContext) AllGreaterThan() []antlr.TerminalNode {
	return s.GetTokens(wclGreaterThan)
}

func (s *Event_callContext) GreaterThan(i int) antlr.TerminalNode {
	return s.GetToken(wclGreaterThan, i)
}

func (s *Event_callContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Event_callContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Event_callContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterEvent_call(s)
	}
}

func (s *Event_callContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitEvent_call(s)
	}
}

func (s *Event_callContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitEvent_call(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Event_call() (localctx IEvent_callContext) {
	this := p
	_ = this

	localctx = NewEvent_callContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, wclRULE_event_call)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(336)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == wclGreaterThan {
		{
			p.SetState(334)
			p.Match(wclGreaterThan)
		}
		{
			p.SetState(335)
			p.Match(wclGreaterThan)
		}

	}
	{
		p.SetState(338)
		p.Func_invoc()
	}

	return localctx
}

// ISelectorContext is an interface to support dynamic dispatch.
type ISelectorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSelectorContext differentiates from other interfaces.
	IsSelectorContext()
}

type SelectorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectorContext() *SelectorContext {
	var p = new(SelectorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = wclRULE_selector
	return p
}

func (*SelectorContext) IsSelectorContext() {}

func NewSelectorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectorContext {
	var p = new(SelectorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = wclRULE_selector

	return p
}

func (s *SelectorContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectorContext) Hash() antlr.TerminalNode {
	return s.GetToken(wclHash, 0)
}

func (s *SelectorContext) Id() antlr.TerminalNode {
	return s.GetToken(wclId, 0)
}

func (s *SelectorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SelectorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.EnterSelector(s)
	}
}

func (s *SelectorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(wclListener); ok {
		listenerT.ExitSelector(s)
	}
}

func (s *SelectorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case wclVisitor:
		return t.VisitSelector(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *wcl) Selector() (localctx ISelectorContext) {
	this := p
	_ = this

	localctx = NewSelectorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, wclRULE_selector)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(343)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case wclHash:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(340)
			p.Match(wclHash)
		}
		{
			p.SetState(341)
			p.Match(wclId)
		}

	case wclId:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(342)
			p.Match(wclId)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}