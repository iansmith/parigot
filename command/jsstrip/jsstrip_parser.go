// Code generated from command/jsstrip/jsstrip.g4 by ANTLR 4.9. DO NOT EDIT.

package main // jsstrip
import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 28, 177,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 3, 2, 3, 2,
	3, 2, 3, 2, 3, 2, 3, 2, 7, 2, 51, 10, 2, 12, 2, 14, 2, 54, 11, 2, 3, 2,
	3, 2, 3, 3, 3, 3, 3, 3, 5, 3, 61, 10, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3, 5,
	3, 5, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 7, 3, 7, 3, 7,
	3, 7, 3, 8, 3, 8, 3, 8, 5, 8, 84, 10, 8, 3, 8, 5, 8, 87, 10, 8, 3, 8, 3,
	8, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 10, 3, 10,
	3, 11, 6, 11, 103, 10, 11, 13, 11, 14, 11, 104, 3, 12, 3, 12, 3, 12, 3,
	12, 3, 12, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 14, 3, 14, 3, 14, 3, 14,
	3, 14, 3, 15, 3, 15, 3, 15, 3, 15, 3, 15, 5, 15, 127, 10, 15, 3, 15, 3,
	15, 3, 16, 6, 16, 132, 10, 16, 13, 16, 14, 16, 133, 3, 17, 3, 17, 3, 18,
	3, 18, 3, 18, 3, 19, 3, 19, 5, 19, 143, 10, 19, 3, 20, 3, 20, 3, 20, 3,
	20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 5, 20, 157,
	10, 20, 3, 21, 6, 21, 160, 10, 21, 13, 21, 14, 21, 161, 3, 22, 3, 22, 5,
	22, 166, 10, 22, 3, 22, 5, 22, 169, 10, 22, 3, 22, 3, 22, 3, 22, 3, 22,
	5, 22, 175, 10, 22, 3, 22, 2, 2, 23, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20,
	22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 2, 3, 3, 2, 11, 13, 2, 174,
	2, 44, 3, 2, 2, 2, 4, 60, 3, 2, 2, 2, 6, 62, 3, 2, 2, 2, 8, 66, 3, 2, 2,
	2, 10, 71, 3, 2, 2, 2, 12, 76, 3, 2, 2, 2, 14, 80, 3, 2, 2, 2, 16, 90,
	3, 2, 2, 2, 18, 99, 3, 2, 2, 2, 20, 102, 3, 2, 2, 2, 22, 106, 3, 2, 2,
	2, 24, 111, 3, 2, 2, 2, 26, 116, 3, 2, 2, 2, 28, 121, 3, 2, 2, 2, 30, 131,
	3, 2, 2, 2, 32, 135, 3, 2, 2, 2, 34, 137, 3, 2, 2, 2, 36, 142, 3, 2, 2,
	2, 38, 156, 3, 2, 2, 2, 40, 159, 3, 2, 2, 2, 42, 174, 3, 2, 2, 2, 44, 45,
	7, 15, 2, 2, 45, 52, 7, 3, 2, 2, 46, 47, 7, 15, 2, 2, 47, 48, 5, 4, 3,
	2, 48, 49, 7, 16, 2, 2, 49, 51, 3, 2, 2, 2, 50, 46, 3, 2, 2, 2, 51, 54,
	3, 2, 2, 2, 52, 50, 3, 2, 2, 2, 52, 53, 3, 2, 2, 2, 53, 55, 3, 2, 2, 2,
	54, 52, 3, 2, 2, 2, 55, 56, 7, 16, 2, 2, 56, 3, 3, 2, 2, 2, 57, 61, 5,
	6, 4, 2, 58, 61, 5, 8, 5, 2, 59, 61, 5, 28, 15, 2, 60, 57, 3, 2, 2, 2,
	60, 58, 3, 2, 2, 2, 60, 59, 3, 2, 2, 2, 61, 5, 3, 2, 2, 2, 62, 63, 7, 4,
	2, 2, 63, 64, 5, 12, 7, 2, 64, 65, 5, 14, 8, 2, 65, 7, 3, 2, 2, 2, 66,
	67, 7, 8, 2, 2, 67, 68, 7, 27, 2, 2, 68, 69, 7, 27, 2, 2, 69, 70, 5, 16,
	9, 2, 70, 9, 3, 2, 2, 2, 71, 72, 7, 15, 2, 2, 72, 73, 7, 4, 2, 2, 73, 74,
	7, 18, 2, 2, 74, 75, 7, 16, 2, 2, 75, 11, 3, 2, 2, 2, 76, 77, 7, 15, 2,
	2, 77, 78, 7, 26, 2, 2, 78, 79, 7, 16, 2, 2, 79, 13, 3, 2, 2, 2, 80, 81,
	7, 15, 2, 2, 81, 83, 7, 5, 2, 2, 82, 84, 5, 22, 12, 2, 83, 82, 3, 2, 2,
	2, 83, 84, 3, 2, 2, 2, 84, 86, 3, 2, 2, 2, 85, 87, 5, 24, 13, 2, 86, 85,
	3, 2, 2, 2, 86, 87, 3, 2, 2, 2, 87, 88, 3, 2, 2, 2, 88, 89, 7, 16, 2, 2,
	89, 15, 3, 2, 2, 2, 90, 91, 7, 15, 2, 2, 91, 92, 7, 5, 2, 2, 92, 93, 7,
	19, 2, 2, 93, 94, 7, 15, 2, 2, 94, 95, 7, 4, 2, 2, 95, 96, 7, 18, 2, 2,
	96, 97, 7, 16, 2, 2, 97, 98, 7, 16, 2, 2, 98, 17, 3, 2, 2, 2, 99, 100,
	9, 2, 2, 2, 100, 19, 3, 2, 2, 2, 101, 103, 5, 18, 10, 2, 102, 101, 3, 2,
	2, 2, 103, 104, 3, 2, 2, 2, 104, 102, 3, 2, 2, 2, 104, 105, 3, 2, 2, 2,
	105, 21, 3, 2, 2, 2, 106, 107, 7, 15, 2, 2, 107, 108, 7, 6, 2, 2, 108,
	109, 5, 20, 11, 2, 109, 110, 7, 16, 2, 2, 110, 23, 3, 2, 2, 2, 111, 112,
	7, 15, 2, 2, 112, 113, 7, 7, 2, 2, 113, 114, 5, 20, 11, 2, 114, 115, 7,
	16, 2, 2, 115, 25, 3, 2, 2, 2, 116, 117, 7, 15, 2, 2, 117, 118, 7, 9, 2,
	2, 118, 119, 5, 20, 11, 2, 119, 120, 7, 16, 2, 2, 120, 27, 3, 2, 2, 2,
	121, 122, 7, 5, 2, 2, 122, 123, 7, 19, 2, 2, 123, 124, 5, 10, 6, 2, 124,
	126, 5, 22, 12, 2, 125, 127, 5, 26, 14, 2, 126, 125, 3, 2, 2, 2, 126, 127,
	3, 2, 2, 2, 127, 128, 3, 2, 2, 2, 128, 129, 5, 30, 16, 2, 129, 29, 3, 2,
	2, 2, 130, 132, 5, 32, 17, 2, 131, 130, 3, 2, 2, 2, 132, 133, 3, 2, 2,
	2, 133, 131, 3, 2, 2, 2, 133, 134, 3, 2, 2, 2, 134, 31, 3, 2, 2, 2, 135,
	136, 5, 34, 18, 2, 136, 33, 3, 2, 2, 2, 137, 138, 7, 10, 2, 2, 138, 139,
	5, 24, 13, 2, 139, 35, 3, 2, 2, 2, 140, 143, 5, 38, 20, 2, 141, 143, 5,
	42, 22, 2, 142, 140, 3, 2, 2, 2, 142, 141, 3, 2, 2, 2, 143, 37, 3, 2, 2,
	2, 144, 145, 7, 15, 2, 2, 145, 157, 7, 16, 2, 2, 146, 147, 7, 15, 2, 2,
	147, 148, 7, 25, 2, 2, 148, 157, 7, 16, 2, 2, 149, 150, 7, 15, 2, 2, 150,
	151, 7, 24, 2, 2, 151, 157, 7, 16, 2, 2, 152, 153, 7, 15, 2, 2, 153, 154,
	5, 40, 21, 2, 154, 155, 7, 16, 2, 2, 155, 157, 3, 2, 2, 2, 156, 144, 3,
	2, 2, 2, 156, 146, 3, 2, 2, 2, 156, 149, 3, 2, 2, 2, 156, 152, 3, 2, 2,
	2, 157, 39, 3, 2, 2, 2, 158, 160, 5, 36, 19, 2, 159, 158, 3, 2, 2, 2, 160,
	161, 3, 2, 2, 2, 161, 159, 3, 2, 2, 2, 161, 162, 3, 2, 2, 2, 162, 41, 3,
	2, 2, 2, 163, 165, 7, 19, 2, 2, 164, 166, 7, 21, 2, 2, 165, 164, 3, 2,
	2, 2, 165, 166, 3, 2, 2, 2, 166, 168, 3, 2, 2, 2, 167, 169, 7, 22, 2, 2,
	168, 167, 3, 2, 2, 2, 168, 169, 3, 2, 2, 2, 169, 175, 3, 2, 2, 2, 170,
	175, 7, 18, 2, 2, 171, 175, 7, 27, 2, 2, 172, 175, 7, 20, 2, 2, 173, 175,
	7, 23, 2, 2, 174, 163, 3, 2, 2, 2, 174, 170, 3, 2, 2, 2, 174, 171, 3, 2,
	2, 2, 174, 172, 3, 2, 2, 2, 174, 173, 3, 2, 2, 2, 175, 43, 3, 2, 2, 2,
	15, 52, 60, 83, 86, 104, 126, 133, 142, 156, 161, 165, 168, 174,
}
var literalNames = []string{
	"", "'module'", "'type'", "'func'", "'param'", "'result'", "'import'",
	"'local'", "'block'", "'i32'", "'i64'", "'f64'", "", "'('", "')'", "'\"'",
}
var symbolicNames = []string{
	"", "ModuleWord", "TypeWord", "FuncWord", "ParamWord", "ResultWord", "ImportWord",
	"LocalWord", "BlockWord", "I32", "I64", "F64", "Whitespace", "Lparen",
	"Rparen", "Quote", "Num", "Ident", "HexPointer", "Offset", "Align", "ConstValue",
	"ConstAnnotation", "BlockAnnotation", "TypeAnnotation", "QuotedString",
	"Comment",
}

var ruleNames = []string{
	"module", "topLevel", "typeDef", "importOp", "typeRef", "typeAnnotation",
	"funcSpec", "funcRef", "type_", "typeSeq", "paramDef", "resultDef", "localDef",
	"funcDef", "funcBody", "stmt", "block", "sexpr", "list", "members", "atom",
}

type jsstripParser struct {
	*antlr.BaseParser
}

// NewjsstripParser produces a new parser instance for the optional input antlr.TokenStream.
//
// The *jsstripParser instance produced may be reused by calling the SetInputStream method.
// The initial parser configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func NewjsstripParser(input antlr.TokenStream) *jsstripParser {
	this := new(jsstripParser)
	deserializer := antlr.NewATNDeserializer(nil)
	deserializedATN := deserializer.DeserializeFromUInt16(parserATN)
	decisionToDFA := make([]*antlr.DFA, len(deserializedATN.DecisionToState))
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "jsstrip.g4"

	return this
}

// jsstripParser tokens.
const (
	jsstripParserEOF             = antlr.TokenEOF
	jsstripParserModuleWord      = 1
	jsstripParserTypeWord        = 2
	jsstripParserFuncWord        = 3
	jsstripParserParamWord       = 4
	jsstripParserResultWord      = 5
	jsstripParserImportWord      = 6
	jsstripParserLocalWord       = 7
	jsstripParserBlockWord       = 8
	jsstripParserI32             = 9
	jsstripParserI64             = 10
	jsstripParserF64             = 11
	jsstripParserWhitespace      = 12
	jsstripParserLparen          = 13
	jsstripParserRparen          = 14
	jsstripParserQuote           = 15
	jsstripParserNum             = 16
	jsstripParserIdent           = 17
	jsstripParserHexPointer      = 18
	jsstripParserOffset          = 19
	jsstripParserAlign           = 20
	jsstripParserConstValue      = 21
	jsstripParserConstAnnotation = 22
	jsstripParserBlockAnnotation = 23
	jsstripParserTypeAnnotation  = 24
	jsstripParserQuotedString    = 25
	jsstripParserComment         = 26
)

// jsstripParser rules.
const (
	jsstripParserRULE_module         = 0
	jsstripParserRULE_topLevel       = 1
	jsstripParserRULE_typeDef        = 2
	jsstripParserRULE_importOp       = 3
	jsstripParserRULE_typeRef        = 4
	jsstripParserRULE_typeAnnotation = 5
	jsstripParserRULE_funcSpec       = 6
	jsstripParserRULE_funcRef        = 7
	jsstripParserRULE_type_          = 8
	jsstripParserRULE_typeSeq        = 9
	jsstripParserRULE_paramDef       = 10
	jsstripParserRULE_resultDef      = 11
	jsstripParserRULE_localDef       = 12
	jsstripParserRULE_funcDef        = 13
	jsstripParserRULE_funcBody       = 14
	jsstripParserRULE_stmt           = 15
	jsstripParserRULE_block          = 16
	jsstripParserRULE_sexpr          = 17
	jsstripParserRULE_list           = 18
	jsstripParserRULE_members        = 19
	jsstripParserRULE_atom           = 20
)

// IModuleContext is an interface to support dynamic dispatch.
type IModuleContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsModuleContext differentiates from other interfaces.
	IsModuleContext()
}

type ModuleContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyModuleContext() *ModuleContext {
	var p = new(ModuleContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_module
	return p
}

func (*ModuleContext) IsModuleContext() {}

func NewModuleContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ModuleContext {
	var p = new(ModuleContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_module

	return p
}

func (s *ModuleContext) GetParser() antlr.Parser { return s.parser }

func (s *ModuleContext) AllLparen() []antlr.TerminalNode {
	return s.GetTokens(jsstripParserLparen)
}

func (s *ModuleContext) Lparen(i int) antlr.TerminalNode {
	return s.GetToken(jsstripParserLparen, i)
}

func (s *ModuleContext) ModuleWord() antlr.TerminalNode {
	return s.GetToken(jsstripParserModuleWord, 0)
}

func (s *ModuleContext) AllRparen() []antlr.TerminalNode {
	return s.GetTokens(jsstripParserRparen)
}

func (s *ModuleContext) Rparen(i int) antlr.TerminalNode {
	return s.GetToken(jsstripParserRparen, i)
}

func (s *ModuleContext) AllTopLevel() []ITopLevelContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITopLevelContext)(nil)).Elem())
	var tst = make([]ITopLevelContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITopLevelContext)
		}
	}

	return tst
}

func (s *ModuleContext) TopLevel(i int) ITopLevelContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITopLevelContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITopLevelContext)
}

func (s *ModuleContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ModuleContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ModuleContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterModule(s)
	}
}

func (s *ModuleContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitModule(s)
	}
}

func (p *jsstripParser) Module() (localctx IModuleContext) {
	localctx = NewModuleContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, jsstripParserRULE_module)
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
		p.SetState(42)
		p.Match(jsstripParserLparen)
	}
	{
		p.SetState(43)
		p.Match(jsstripParserModuleWord)
	}
	p.SetState(50)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == jsstripParserLparen {
		{
			p.SetState(44)
			p.Match(jsstripParserLparen)
		}
		{
			p.SetState(45)
			p.TopLevel()
		}
		{
			p.SetState(46)
			p.Match(jsstripParserRparen)
		}

		p.SetState(52)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(53)
		p.Match(jsstripParserRparen)
	}

	return localctx
}

// ITopLevelContext is an interface to support dynamic dispatch.
type ITopLevelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTopLevelContext differentiates from other interfaces.
	IsTopLevelContext()
}

type TopLevelContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTopLevelContext() *TopLevelContext {
	var p = new(TopLevelContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_topLevel
	return p
}

func (*TopLevelContext) IsTopLevelContext() {}

func NewTopLevelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TopLevelContext {
	var p = new(TopLevelContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_topLevel

	return p
}

func (s *TopLevelContext) GetParser() antlr.Parser { return s.parser }

func (s *TopLevelContext) TypeDef() ITypeDefContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeDefContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeDefContext)
}

func (s *TopLevelContext) ImportOp() IImportOpContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImportOpContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImportOpContext)
}

func (s *TopLevelContext) FuncDef() IFuncDefContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFuncDefContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFuncDefContext)
}

func (s *TopLevelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TopLevelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TopLevelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterTopLevel(s)
	}
}

func (s *TopLevelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitTopLevel(s)
	}
}

func (p *jsstripParser) TopLevel() (localctx ITopLevelContext) {
	localctx = NewTopLevelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, jsstripParserRULE_topLevel)

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

	p.SetState(58)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case jsstripParserTypeWord:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(55)
			p.TypeDef()
		}

	case jsstripParserImportWord:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(56)
			p.ImportOp()
		}

	case jsstripParserFuncWord:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(57)
			p.FuncDef()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ITypeDefContext is an interface to support dynamic dispatch.
type ITypeDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypeDefContext differentiates from other interfaces.
	IsTypeDefContext()
}

type TypeDefContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeDefContext() *TypeDefContext {
	var p = new(TypeDefContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_typeDef
	return p
}

func (*TypeDefContext) IsTypeDefContext() {}

func NewTypeDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeDefContext {
	var p = new(TypeDefContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_typeDef

	return p
}

func (s *TypeDefContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeDefContext) TypeWord() antlr.TerminalNode {
	return s.GetToken(jsstripParserTypeWord, 0)
}

func (s *TypeDefContext) TypeAnnotation() ITypeAnnotationContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeAnnotationContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeAnnotationContext)
}

func (s *TypeDefContext) FuncSpec() IFuncSpecContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFuncSpecContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFuncSpecContext)
}

func (s *TypeDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterTypeDef(s)
	}
}

func (s *TypeDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitTypeDef(s)
	}
}

func (p *jsstripParser) TypeDef() (localctx ITypeDefContext) {
	localctx = NewTypeDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, jsstripParserRULE_typeDef)

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
		p.SetState(60)
		p.Match(jsstripParserTypeWord)
	}
	{
		p.SetState(61)
		p.TypeAnnotation()
	}
	{
		p.SetState(62)
		p.FuncSpec()
	}

	return localctx
}

// IImportOpContext is an interface to support dynamic dispatch.
type IImportOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsImportOpContext differentiates from other interfaces.
	IsImportOpContext()
}

type ImportOpContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImportOpContext() *ImportOpContext {
	var p = new(ImportOpContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_importOp
	return p
}

func (*ImportOpContext) IsImportOpContext() {}

func NewImportOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportOpContext {
	var p = new(ImportOpContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_importOp

	return p
}

func (s *ImportOpContext) GetParser() antlr.Parser { return s.parser }

func (s *ImportOpContext) ImportWord() antlr.TerminalNode {
	return s.GetToken(jsstripParserImportWord, 0)
}

func (s *ImportOpContext) AllQuotedString() []antlr.TerminalNode {
	return s.GetTokens(jsstripParserQuotedString)
}

func (s *ImportOpContext) QuotedString(i int) antlr.TerminalNode {
	return s.GetToken(jsstripParserQuotedString, i)
}

func (s *ImportOpContext) FuncRef() IFuncRefContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFuncRefContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFuncRefContext)
}

func (s *ImportOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImportOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImportOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterImportOp(s)
	}
}

func (s *ImportOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitImportOp(s)
	}
}

func (p *jsstripParser) ImportOp() (localctx IImportOpContext) {
	localctx = NewImportOpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, jsstripParserRULE_importOp)

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
		p.SetState(64)
		p.Match(jsstripParserImportWord)
	}
	{
		p.SetState(65)
		p.Match(jsstripParserQuotedString)
	}
	{
		p.SetState(66)
		p.Match(jsstripParserQuotedString)
	}
	{
		p.SetState(67)
		p.FuncRef()
	}

	return localctx
}

// ITypeRefContext is an interface to support dynamic dispatch.
type ITypeRefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypeRefContext differentiates from other interfaces.
	IsTypeRefContext()
}

type TypeRefContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeRefContext() *TypeRefContext {
	var p = new(TypeRefContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_typeRef
	return p
}

func (*TypeRefContext) IsTypeRefContext() {}

func NewTypeRefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeRefContext {
	var p = new(TypeRefContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_typeRef

	return p
}

func (s *TypeRefContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeRefContext) Lparen() antlr.TerminalNode {
	return s.GetToken(jsstripParserLparen, 0)
}

func (s *TypeRefContext) TypeWord() antlr.TerminalNode {
	return s.GetToken(jsstripParserTypeWord, 0)
}

func (s *TypeRefContext) Num() antlr.TerminalNode {
	return s.GetToken(jsstripParserNum, 0)
}

func (s *TypeRefContext) Rparen() antlr.TerminalNode {
	return s.GetToken(jsstripParserRparen, 0)
}

func (s *TypeRefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeRefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeRefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterTypeRef(s)
	}
}

func (s *TypeRefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitTypeRef(s)
	}
}

func (p *jsstripParser) TypeRef() (localctx ITypeRefContext) {
	localctx = NewTypeRefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, jsstripParserRULE_typeRef)

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
		p.SetState(69)
		p.Match(jsstripParserLparen)
	}
	{
		p.SetState(70)
		p.Match(jsstripParserTypeWord)
	}
	{
		p.SetState(71)
		p.Match(jsstripParserNum)
	}
	{
		p.SetState(72)
		p.Match(jsstripParserRparen)
	}

	return localctx
}

// ITypeAnnotationContext is an interface to support dynamic dispatch.
type ITypeAnnotationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypeAnnotationContext differentiates from other interfaces.
	IsTypeAnnotationContext()
}

type TypeAnnotationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeAnnotationContext() *TypeAnnotationContext {
	var p = new(TypeAnnotationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_typeAnnotation
	return p
}

func (*TypeAnnotationContext) IsTypeAnnotationContext() {}

func NewTypeAnnotationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeAnnotationContext {
	var p = new(TypeAnnotationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_typeAnnotation

	return p
}

func (s *TypeAnnotationContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeAnnotationContext) Lparen() antlr.TerminalNode {
	return s.GetToken(jsstripParserLparen, 0)
}

func (s *TypeAnnotationContext) TypeAnnotation() antlr.TerminalNode {
	return s.GetToken(jsstripParserTypeAnnotation, 0)
}

func (s *TypeAnnotationContext) Rparen() antlr.TerminalNode {
	return s.GetToken(jsstripParserRparen, 0)
}

func (s *TypeAnnotationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeAnnotationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeAnnotationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterTypeAnnotation(s)
	}
}

func (s *TypeAnnotationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitTypeAnnotation(s)
	}
}

func (p *jsstripParser) TypeAnnotation() (localctx ITypeAnnotationContext) {
	localctx = NewTypeAnnotationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, jsstripParserRULE_typeAnnotation)

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
		p.SetState(74)
		p.Match(jsstripParserLparen)
	}
	{
		p.SetState(75)
		p.Match(jsstripParserTypeAnnotation)
	}
	{
		p.SetState(76)
		p.Match(jsstripParserRparen)
	}

	return localctx
}

// IFuncSpecContext is an interface to support dynamic dispatch.
type IFuncSpecContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFuncSpecContext differentiates from other interfaces.
	IsFuncSpecContext()
}

type FuncSpecContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFuncSpecContext() *FuncSpecContext {
	var p = new(FuncSpecContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_funcSpec
	return p
}

func (*FuncSpecContext) IsFuncSpecContext() {}

func NewFuncSpecContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncSpecContext {
	var p = new(FuncSpecContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_funcSpec

	return p
}

func (s *FuncSpecContext) GetParser() antlr.Parser { return s.parser }

func (s *FuncSpecContext) Lparen() antlr.TerminalNode {
	return s.GetToken(jsstripParserLparen, 0)
}

func (s *FuncSpecContext) FuncWord() antlr.TerminalNode {
	return s.GetToken(jsstripParserFuncWord, 0)
}

func (s *FuncSpecContext) Rparen() antlr.TerminalNode {
	return s.GetToken(jsstripParserRparen, 0)
}

func (s *FuncSpecContext) ParamDef() IParamDefContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamDefContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamDefContext)
}

func (s *FuncSpecContext) ResultDef() IResultDefContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IResultDefContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IResultDefContext)
}

func (s *FuncSpecContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncSpecContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FuncSpecContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterFuncSpec(s)
	}
}

func (s *FuncSpecContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitFuncSpec(s)
	}
}

func (p *jsstripParser) FuncSpec() (localctx IFuncSpecContext) {
	localctx = NewFuncSpecContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, jsstripParserRULE_funcSpec)
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
		p.SetState(78)
		p.Match(jsstripParserLparen)
	}
	{
		p.SetState(79)
		p.Match(jsstripParserFuncWord)
	}
	p.SetState(81)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(80)
			p.ParamDef()
		}

	}
	p.SetState(84)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == jsstripParserLparen {
		{
			p.SetState(83)
			p.ResultDef()
		}

	}
	{
		p.SetState(86)
		p.Match(jsstripParserRparen)
	}

	return localctx
}

// IFuncRefContext is an interface to support dynamic dispatch.
type IFuncRefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFuncRefContext differentiates from other interfaces.
	IsFuncRefContext()
}

type FuncRefContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFuncRefContext() *FuncRefContext {
	var p = new(FuncRefContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_funcRef
	return p
}

func (*FuncRefContext) IsFuncRefContext() {}

func NewFuncRefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncRefContext {
	var p = new(FuncRefContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_funcRef

	return p
}

func (s *FuncRefContext) GetParser() antlr.Parser { return s.parser }

func (s *FuncRefContext) AllLparen() []antlr.TerminalNode {
	return s.GetTokens(jsstripParserLparen)
}

func (s *FuncRefContext) Lparen(i int) antlr.TerminalNode {
	return s.GetToken(jsstripParserLparen, i)
}

func (s *FuncRefContext) FuncWord() antlr.TerminalNode {
	return s.GetToken(jsstripParserFuncWord, 0)
}

func (s *FuncRefContext) Ident() antlr.TerminalNode {
	return s.GetToken(jsstripParserIdent, 0)
}

func (s *FuncRefContext) TypeWord() antlr.TerminalNode {
	return s.GetToken(jsstripParserTypeWord, 0)
}

func (s *FuncRefContext) Num() antlr.TerminalNode {
	return s.GetToken(jsstripParserNum, 0)
}

func (s *FuncRefContext) AllRparen() []antlr.TerminalNode {
	return s.GetTokens(jsstripParserRparen)
}

func (s *FuncRefContext) Rparen(i int) antlr.TerminalNode {
	return s.GetToken(jsstripParserRparen, i)
}

func (s *FuncRefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncRefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FuncRefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterFuncRef(s)
	}
}

func (s *FuncRefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitFuncRef(s)
	}
}

func (p *jsstripParser) FuncRef() (localctx IFuncRefContext) {
	localctx = NewFuncRefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, jsstripParserRULE_funcRef)

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
		p.SetState(88)
		p.Match(jsstripParserLparen)
	}
	{
		p.SetState(89)
		p.Match(jsstripParserFuncWord)
	}
	{
		p.SetState(90)
		p.Match(jsstripParserIdent)
	}
	{
		p.SetState(91)
		p.Match(jsstripParserLparen)
	}
	{
		p.SetState(92)
		p.Match(jsstripParserTypeWord)
	}
	{
		p.SetState(93)
		p.Match(jsstripParserNum)
	}
	{
		p.SetState(94)
		p.Match(jsstripParserRparen)
	}
	{
		p.SetState(95)
		p.Match(jsstripParserRparen)
	}

	return localctx
}

// IType_Context is an interface to support dynamic dispatch.
type IType_Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsType_Context differentiates from other interfaces.
	IsType_Context()
}

type Type_Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyType_Context() *Type_Context {
	var p = new(Type_Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_type_
	return p
}

func (*Type_Context) IsType_Context() {}

func NewType_Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Type_Context {
	var p = new(Type_Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_type_

	return p
}

func (s *Type_Context) GetParser() antlr.Parser { return s.parser }

func (s *Type_Context) I32() antlr.TerminalNode {
	return s.GetToken(jsstripParserI32, 0)
}

func (s *Type_Context) I64() antlr.TerminalNode {
	return s.GetToken(jsstripParserI64, 0)
}

func (s *Type_Context) F64() antlr.TerminalNode {
	return s.GetToken(jsstripParserF64, 0)
}

func (s *Type_Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Type_Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Type_Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterType_(s)
	}
}

func (s *Type_Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitType_(s)
	}
}

func (p *jsstripParser) Type_() (localctx IType_Context) {
	localctx = NewType_Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, jsstripParserRULE_type_)
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
		p.SetState(97)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<jsstripParserI32)|(1<<jsstripParserI64)|(1<<jsstripParserF64))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// ITypeSeqContext is an interface to support dynamic dispatch.
type ITypeSeqContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypeSeqContext differentiates from other interfaces.
	IsTypeSeqContext()
}

type TypeSeqContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeSeqContext() *TypeSeqContext {
	var p = new(TypeSeqContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_typeSeq
	return p
}

func (*TypeSeqContext) IsTypeSeqContext() {}

func NewTypeSeqContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeSeqContext {
	var p = new(TypeSeqContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_typeSeq

	return p
}

func (s *TypeSeqContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeSeqContext) AllType_() []IType_Context {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IType_Context)(nil)).Elem())
	var tst = make([]IType_Context, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IType_Context)
		}
	}

	return tst
}

func (s *TypeSeqContext) Type_(i int) IType_Context {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IType_Context)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IType_Context)
}

func (s *TypeSeqContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeSeqContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeSeqContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterTypeSeq(s)
	}
}

func (s *TypeSeqContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitTypeSeq(s)
	}
}

func (p *jsstripParser) TypeSeq() (localctx ITypeSeqContext) {
	localctx = NewTypeSeqContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, jsstripParserRULE_typeSeq)
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
	p.SetState(100)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<jsstripParserI32)|(1<<jsstripParserI64)|(1<<jsstripParserF64))) != 0) {
		{
			p.SetState(99)
			p.Type_()
		}

		p.SetState(102)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IParamDefContext is an interface to support dynamic dispatch.
type IParamDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsParamDefContext differentiates from other interfaces.
	IsParamDefContext()
}

type ParamDefContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParamDefContext() *ParamDefContext {
	var p = new(ParamDefContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_paramDef
	return p
}

func (*ParamDefContext) IsParamDefContext() {}

func NewParamDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamDefContext {
	var p = new(ParamDefContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_paramDef

	return p
}

func (s *ParamDefContext) GetParser() antlr.Parser { return s.parser }

func (s *ParamDefContext) Lparen() antlr.TerminalNode {
	return s.GetToken(jsstripParserLparen, 0)
}

func (s *ParamDefContext) ParamWord() antlr.TerminalNode {
	return s.GetToken(jsstripParserParamWord, 0)
}

func (s *ParamDefContext) TypeSeq() ITypeSeqContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeSeqContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeSeqContext)
}

func (s *ParamDefContext) Rparen() antlr.TerminalNode {
	return s.GetToken(jsstripParserRparen, 0)
}

func (s *ParamDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParamDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParamDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterParamDef(s)
	}
}

func (s *ParamDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitParamDef(s)
	}
}

func (p *jsstripParser) ParamDef() (localctx IParamDefContext) {
	localctx = NewParamDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, jsstripParserRULE_paramDef)

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
		p.SetState(104)
		p.Match(jsstripParserLparen)
	}
	{
		p.SetState(105)
		p.Match(jsstripParserParamWord)
	}
	{
		p.SetState(106)
		p.TypeSeq()
	}
	{
		p.SetState(107)
		p.Match(jsstripParserRparen)
	}

	return localctx
}

// IResultDefContext is an interface to support dynamic dispatch.
type IResultDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsResultDefContext differentiates from other interfaces.
	IsResultDefContext()
}

type ResultDefContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyResultDefContext() *ResultDefContext {
	var p = new(ResultDefContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_resultDef
	return p
}

func (*ResultDefContext) IsResultDefContext() {}

func NewResultDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ResultDefContext {
	var p = new(ResultDefContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_resultDef

	return p
}

func (s *ResultDefContext) GetParser() antlr.Parser { return s.parser }

func (s *ResultDefContext) Lparen() antlr.TerminalNode {
	return s.GetToken(jsstripParserLparen, 0)
}

func (s *ResultDefContext) ResultWord() antlr.TerminalNode {
	return s.GetToken(jsstripParserResultWord, 0)
}

func (s *ResultDefContext) TypeSeq() ITypeSeqContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeSeqContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeSeqContext)
}

func (s *ResultDefContext) Rparen() antlr.TerminalNode {
	return s.GetToken(jsstripParserRparen, 0)
}

func (s *ResultDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ResultDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ResultDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterResultDef(s)
	}
}

func (s *ResultDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitResultDef(s)
	}
}

func (p *jsstripParser) ResultDef() (localctx IResultDefContext) {
	localctx = NewResultDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, jsstripParserRULE_resultDef)

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
		p.SetState(109)
		p.Match(jsstripParserLparen)
	}
	{
		p.SetState(110)
		p.Match(jsstripParserResultWord)
	}
	{
		p.SetState(111)
		p.TypeSeq()
	}
	{
		p.SetState(112)
		p.Match(jsstripParserRparen)
	}

	return localctx
}

// ILocalDefContext is an interface to support dynamic dispatch.
type ILocalDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLocalDefContext differentiates from other interfaces.
	IsLocalDefContext()
}

type LocalDefContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLocalDefContext() *LocalDefContext {
	var p = new(LocalDefContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_localDef
	return p
}

func (*LocalDefContext) IsLocalDefContext() {}

func NewLocalDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LocalDefContext {
	var p = new(LocalDefContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_localDef

	return p
}

func (s *LocalDefContext) GetParser() antlr.Parser { return s.parser }

func (s *LocalDefContext) Lparen() antlr.TerminalNode {
	return s.GetToken(jsstripParserLparen, 0)
}

func (s *LocalDefContext) LocalWord() antlr.TerminalNode {
	return s.GetToken(jsstripParserLocalWord, 0)
}

func (s *LocalDefContext) TypeSeq() ITypeSeqContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeSeqContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeSeqContext)
}

func (s *LocalDefContext) Rparen() antlr.TerminalNode {
	return s.GetToken(jsstripParserRparen, 0)
}

func (s *LocalDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LocalDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LocalDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterLocalDef(s)
	}
}

func (s *LocalDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitLocalDef(s)
	}
}

func (p *jsstripParser) LocalDef() (localctx ILocalDefContext) {
	localctx = NewLocalDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, jsstripParserRULE_localDef)

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
		p.Match(jsstripParserLparen)
	}
	{
		p.SetState(115)
		p.Match(jsstripParserLocalWord)
	}
	{
		p.SetState(116)
		p.TypeSeq()
	}
	{
		p.SetState(117)
		p.Match(jsstripParserRparen)
	}

	return localctx
}

// IFuncDefContext is an interface to support dynamic dispatch.
type IFuncDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFuncDefContext differentiates from other interfaces.
	IsFuncDefContext()
}

type FuncDefContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFuncDefContext() *FuncDefContext {
	var p = new(FuncDefContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_funcDef
	return p
}

func (*FuncDefContext) IsFuncDefContext() {}

func NewFuncDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncDefContext {
	var p = new(FuncDefContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_funcDef

	return p
}

func (s *FuncDefContext) GetParser() antlr.Parser { return s.parser }

func (s *FuncDefContext) FuncWord() antlr.TerminalNode {
	return s.GetToken(jsstripParserFuncWord, 0)
}

func (s *FuncDefContext) Ident() antlr.TerminalNode {
	return s.GetToken(jsstripParserIdent, 0)
}

func (s *FuncDefContext) TypeRef() ITypeRefContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeRefContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeRefContext)
}

func (s *FuncDefContext) ParamDef() IParamDefContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamDefContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamDefContext)
}

func (s *FuncDefContext) FuncBody() IFuncBodyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFuncBodyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFuncBodyContext)
}

func (s *FuncDefContext) LocalDef() ILocalDefContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILocalDefContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILocalDefContext)
}

func (s *FuncDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FuncDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterFuncDef(s)
	}
}

func (s *FuncDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitFuncDef(s)
	}
}

func (p *jsstripParser) FuncDef() (localctx IFuncDefContext) {
	localctx = NewFuncDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, jsstripParserRULE_funcDef)
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
		p.SetState(119)
		p.Match(jsstripParserFuncWord)
	}
	{
		p.SetState(120)
		p.Match(jsstripParserIdent)
	}
	{
		p.SetState(121)
		p.TypeRef()
	}
	{
		p.SetState(122)
		p.ParamDef()
	}
	p.SetState(124)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == jsstripParserLparen {
		{
			p.SetState(123)
			p.LocalDef()
		}

	}
	{
		p.SetState(126)
		p.FuncBody()
	}

	return localctx
}

// IFuncBodyContext is an interface to support dynamic dispatch.
type IFuncBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFuncBodyContext differentiates from other interfaces.
	IsFuncBodyContext()
}

type FuncBodyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFuncBodyContext() *FuncBodyContext {
	var p = new(FuncBodyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_funcBody
	return p
}

func (*FuncBodyContext) IsFuncBodyContext() {}

func NewFuncBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncBodyContext {
	var p = new(FuncBodyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_funcBody

	return p
}

func (s *FuncBodyContext) GetParser() antlr.Parser { return s.parser }

func (s *FuncBodyContext) AllStmt() []IStmtContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStmtContext)(nil)).Elem())
	var tst = make([]IStmtContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStmtContext)
		}
	}

	return tst
}

func (s *FuncBodyContext) Stmt(i int) IStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStmtContext)
}

func (s *FuncBodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FuncBodyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterFuncBody(s)
	}
}

func (s *FuncBodyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitFuncBody(s)
	}
}

func (p *jsstripParser) FuncBody() (localctx IFuncBodyContext) {
	localctx = NewFuncBodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, jsstripParserRULE_funcBody)
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
	p.SetState(129)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == jsstripParserBlockWord {
		{
			p.SetState(128)
			p.Stmt()
		}

		p.SetState(131)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IStmtContext is an interface to support dynamic dispatch.
type IStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStmtContext differentiates from other interfaces.
	IsStmtContext()
}

type StmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStmtContext() *StmtContext {
	var p = new(StmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_stmt
	return p
}

func (*StmtContext) IsStmtContext() {}

func NewStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StmtContext {
	var p = new(StmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_stmt

	return p
}

func (s *StmtContext) GetParser() antlr.Parser { return s.parser }

func (s *StmtContext) Block() IBlockContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBlockContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBlockContext)
}

func (s *StmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterStmt(s)
	}
}

func (s *StmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitStmt(s)
	}
}

func (p *jsstripParser) Stmt() (localctx IStmtContext) {
	localctx = NewStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, jsstripParserRULE_stmt)

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
		p.SetState(133)
		p.Block()
	}

	return localctx
}

// IBlockContext is an interface to support dynamic dispatch.
type IBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBlockContext differentiates from other interfaces.
	IsBlockContext()
}

type BlockContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBlockContext() *BlockContext {
	var p = new(BlockContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_block
	return p
}

func (*BlockContext) IsBlockContext() {}

func NewBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BlockContext {
	var p = new(BlockContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_block

	return p
}

func (s *BlockContext) GetParser() antlr.Parser { return s.parser }

func (s *BlockContext) BlockWord() antlr.TerminalNode {
	return s.GetToken(jsstripParserBlockWord, 0)
}

func (s *BlockContext) ResultDef() IResultDefContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IResultDefContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IResultDefContext)
}

func (s *BlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterBlock(s)
	}
}

func (s *BlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitBlock(s)
	}
}

func (p *jsstripParser) Block() (localctx IBlockContext) {
	localctx = NewBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, jsstripParserRULE_block)

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
		p.SetState(135)
		p.Match(jsstripParserBlockWord)
	}
	{
		p.SetState(136)
		p.ResultDef()
	}

	return localctx
}

// ISexprContext is an interface to support dynamic dispatch.
type ISexprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSexprContext differentiates from other interfaces.
	IsSexprContext()
}

type SexprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySexprContext() *SexprContext {
	var p = new(SexprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_sexpr
	return p
}

func (*SexprContext) IsSexprContext() {}

func NewSexprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SexprContext {
	var p = new(SexprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_sexpr

	return p
}

func (s *SexprContext) GetParser() antlr.Parser { return s.parser }

func (s *SexprContext) List() IListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IListContext)
}

func (s *SexprContext) Atom() IAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *SexprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SexprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SexprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterSexpr(s)
	}
}

func (s *SexprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitSexpr(s)
	}
}

func (p *jsstripParser) Sexpr() (localctx ISexprContext) {
	localctx = NewSexprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, jsstripParserRULE_sexpr)

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

	p.SetState(140)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case jsstripParserLparen:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(138)
			p.List()
		}

	case jsstripParserNum, jsstripParserIdent, jsstripParserHexPointer, jsstripParserConstValue, jsstripParserQuotedString:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(139)
			p.Atom()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IListContext is an interface to support dynamic dispatch.
type IListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsListContext differentiates from other interfaces.
	IsListContext()
}

type ListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyListContext() *ListContext {
	var p = new(ListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_list
	return p
}

func (*ListContext) IsListContext() {}

func NewListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ListContext {
	var p = new(ListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_list

	return p
}

func (s *ListContext) GetParser() antlr.Parser { return s.parser }

func (s *ListContext) Lparen() antlr.TerminalNode {
	return s.GetToken(jsstripParserLparen, 0)
}

func (s *ListContext) Rparen() antlr.TerminalNode {
	return s.GetToken(jsstripParserRparen, 0)
}

func (s *ListContext) BlockAnnotation() antlr.TerminalNode {
	return s.GetToken(jsstripParserBlockAnnotation, 0)
}

func (s *ListContext) ConstAnnotation() antlr.TerminalNode {
	return s.GetToken(jsstripParserConstAnnotation, 0)
}

func (s *ListContext) Members() IMembersContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMembersContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMembersContext)
}

func (s *ListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterList(s)
	}
}

func (s *ListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitList(s)
	}
}

func (p *jsstripParser) List() (localctx IListContext) {
	localctx = NewListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, jsstripParserRULE_list)

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

	p.SetState(154)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(142)
			p.Match(jsstripParserLparen)
		}
		{
			p.SetState(143)
			p.Match(jsstripParserRparen)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(144)
			p.Match(jsstripParserLparen)
		}
		{
			p.SetState(145)
			p.Match(jsstripParserBlockAnnotation)
		}
		{
			p.SetState(146)
			p.Match(jsstripParserRparen)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(147)
			p.Match(jsstripParserLparen)
		}
		{
			p.SetState(148)
			p.Match(jsstripParserConstAnnotation)
		}
		{
			p.SetState(149)
			p.Match(jsstripParserRparen)
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(150)
			p.Match(jsstripParserLparen)
		}
		{
			p.SetState(151)
			p.Members()
		}
		{
			p.SetState(152)
			p.Match(jsstripParserRparen)
		}

	}

	return localctx
}

// IMembersContext is an interface to support dynamic dispatch.
type IMembersContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMembersContext differentiates from other interfaces.
	IsMembersContext()
}

type MembersContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMembersContext() *MembersContext {
	var p = new(MembersContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_members
	return p
}

func (*MembersContext) IsMembersContext() {}

func NewMembersContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MembersContext {
	var p = new(MembersContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_members

	return p
}

func (s *MembersContext) GetParser() antlr.Parser { return s.parser }

func (s *MembersContext) AllSexpr() []ISexprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISexprContext)(nil)).Elem())
	var tst = make([]ISexprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISexprContext)
		}
	}

	return tst
}

func (s *MembersContext) Sexpr(i int) ISexprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISexprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISexprContext)
}

func (s *MembersContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MembersContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MembersContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterMembers(s)
	}
}

func (s *MembersContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitMembers(s)
	}
}

func (p *jsstripParser) Members() (localctx IMembersContext) {
	localctx = NewMembersContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, jsstripParserRULE_members)
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
	p.SetState(157)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<jsstripParserLparen)|(1<<jsstripParserNum)|(1<<jsstripParserIdent)|(1<<jsstripParserHexPointer)|(1<<jsstripParserConstValue)|(1<<jsstripParserQuotedString))) != 0) {
		{
			p.SetState(156)
			p.Sexpr()
		}

		p.SetState(159)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IAtomContext is an interface to support dynamic dispatch.
type IAtomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAtomContext differentiates from other interfaces.
	IsAtomContext()
}

type AtomContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAtomContext() *AtomContext {
	var p = new(AtomContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_atom
	return p
}

func (*AtomContext) IsAtomContext() {}

func NewAtomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AtomContext {
	var p = new(AtomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_atom

	return p
}

func (s *AtomContext) GetParser() antlr.Parser { return s.parser }

func (s *AtomContext) Ident() antlr.TerminalNode {
	return s.GetToken(jsstripParserIdent, 0)
}

func (s *AtomContext) Offset() antlr.TerminalNode {
	return s.GetToken(jsstripParserOffset, 0)
}

func (s *AtomContext) Align() antlr.TerminalNode {
	return s.GetToken(jsstripParserAlign, 0)
}

func (s *AtomContext) Num() antlr.TerminalNode {
	return s.GetToken(jsstripParserNum, 0)
}

func (s *AtomContext) QuotedString() antlr.TerminalNode {
	return s.GetToken(jsstripParserQuotedString, 0)
}

func (s *AtomContext) HexPointer() antlr.TerminalNode {
	return s.GetToken(jsstripParserHexPointer, 0)
}

func (s *AtomContext) ConstValue() antlr.TerminalNode {
	return s.GetToken(jsstripParserConstValue, 0)
}

func (s *AtomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AtomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AtomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterAtom(s)
	}
}

func (s *AtomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitAtom(s)
	}
}

func (p *jsstripParser) Atom() (localctx IAtomContext) {
	localctx = NewAtomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, jsstripParserRULE_atom)
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

	p.SetState(172)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case jsstripParserIdent:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(161)
			p.Match(jsstripParserIdent)
		}
		p.SetState(163)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == jsstripParserOffset {
			{
				p.SetState(162)
				p.Match(jsstripParserOffset)
			}

		}
		p.SetState(166)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == jsstripParserAlign {
			{
				p.SetState(165)
				p.Match(jsstripParserAlign)
			}

		}

	case jsstripParserNum:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(168)
			p.Match(jsstripParserNum)
		}

	case jsstripParserQuotedString:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(169)
			p.Match(jsstripParserQuotedString)
		}

	case jsstripParserHexPointer:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(170)
			p.Match(jsstripParserHexPointer)
		}

	case jsstripParserConstValue:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(171)
			p.Match(jsstripParserConstValue)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}
