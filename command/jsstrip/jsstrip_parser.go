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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 11, 50, 4,
	2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 3, 2, 7, 2, 14,
	10, 2, 12, 2, 14, 2, 17, 11, 2, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 23, 10, 3,
	3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 5, 4, 33, 10, 4, 3, 5,
	6, 5, 36, 10, 5, 13, 5, 14, 5, 37, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6,
	3, 6, 3, 6, 5, 6, 48, 10, 6, 3, 6, 2, 2, 7, 2, 4, 6, 8, 10, 2, 2, 2, 50,
	2, 15, 3, 2, 2, 2, 4, 22, 3, 2, 2, 2, 6, 32, 3, 2, 2, 2, 8, 35, 3, 2, 2,
	2, 10, 47, 3, 2, 2, 2, 12, 14, 5, 4, 3, 2, 13, 12, 3, 2, 2, 2, 14, 17,
	3, 2, 2, 2, 15, 13, 3, 2, 2, 2, 15, 16, 3, 2, 2, 2, 16, 3, 3, 2, 2, 2,
	17, 15, 3, 2, 2, 2, 18, 23, 5, 6, 4, 2, 19, 20, 5, 10, 6, 2, 20, 21, 8,
	3, 1, 2, 21, 23, 3, 2, 2, 2, 22, 18, 3, 2, 2, 2, 22, 19, 3, 2, 2, 2, 23,
	5, 3, 2, 2, 2, 24, 25, 7, 3, 2, 2, 25, 26, 7, 4, 2, 2, 26, 33, 8, 4, 1,
	2, 27, 28, 7, 3, 2, 2, 28, 29, 5, 8, 5, 2, 29, 30, 7, 4, 2, 2, 30, 31,
	8, 4, 1, 2, 31, 33, 3, 2, 2, 2, 32, 24, 3, 2, 2, 2, 32, 27, 3, 2, 2, 2,
	33, 7, 3, 2, 2, 2, 34, 36, 5, 4, 3, 2, 35, 34, 3, 2, 2, 2, 36, 37, 3, 2,
	2, 2, 37, 35, 3, 2, 2, 2, 37, 38, 3, 2, 2, 2, 38, 39, 3, 2, 2, 2, 39, 40,
	8, 5, 1, 2, 40, 9, 3, 2, 2, 2, 41, 42, 7, 7, 2, 2, 42, 48, 8, 6, 1, 2,
	43, 44, 7, 5, 2, 2, 44, 48, 8, 6, 1, 2, 45, 46, 7, 9, 2, 2, 46, 48, 8,
	6, 1, 2, 47, 41, 3, 2, 2, 2, 47, 43, 3, 2, 2, 2, 47, 45, 3, 2, 2, 2, 48,
	11, 3, 2, 2, 2, 7, 15, 22, 32, 37, 47,
}
var literalNames = []string{
	"", "'('", "')'",
}
var symbolicNames = []string{
	"", "", "", "Num", "HexDigit", "Id", "Whitespace", "String_", "HexByteValue",
	"LineComment",
}

var ruleNames = []string{
	"program", "sexpr", "list", "members", "atom",
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
	jsstripParserEOF          = antlr.TokenEOF
	jsstripParserT__0         = 1
	jsstripParserT__1         = 2
	jsstripParserNum          = 3
	jsstripParserHexDigit     = 4
	jsstripParserId           = 5
	jsstripParserWhitespace   = 6
	jsstripParserString_      = 7
	jsstripParserHexByteValue = 8
	jsstripParserLineComment  = 9
)

// jsstripParser rules.
const (
	jsstripParserRULE_program = 0
	jsstripParserRULE_sexpr   = 1
	jsstripParserRULE_list    = 2
	jsstripParserRULE_members = 3
	jsstripParserRULE_atom    = 4
)

// IProgramContext is an interface to support dynamic dispatch.
type IProgramContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsProgramContext differentiates from other interfaces.
	IsProgramContext()
}

type ProgramContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProgramContext() *ProgramContext {
	var p = new(ProgramContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = jsstripParserRULE_program
	return p
}

func (*ProgramContext) IsProgramContext() {}

func NewProgramContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgramContext {
	var p = new(ProgramContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = jsstripParserRULE_program

	return p
}

func (s *ProgramContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgramContext) AllSexpr() []ISexprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISexprContext)(nil)).Elem())
	var tst = make([]ISexprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISexprContext)
		}
	}

	return tst
}

func (s *ProgramContext) Sexpr(i int) ISexprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISexprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISexprContext)
}

func (s *ProgramContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProgramContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProgramContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.EnterProgram(s)
	}
}

func (s *ProgramContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(jsstripListener); ok {
		listenerT.ExitProgram(s)
	}
}

func (p *jsstripParser) Program() (localctx IProgramContext) {
	localctx = NewProgramContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, jsstripParserRULE_program)
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
	p.SetState(13)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<jsstripParserT__0)|(1<<jsstripParserNum)|(1<<jsstripParserId)|(1<<jsstripParserString_))) != 0 {
		{
			p.SetState(10)
			p.Sexpr()
		}

		p.SetState(15)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
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
	p.EnterRule(localctx, 2, jsstripParserRULE_sexpr)

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

	p.SetState(20)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case jsstripParserT__0:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(16)
			p.List()
		}

	case jsstripParserNum, jsstripParserId, jsstripParserString_:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(17)
			p.Atom()
		}
		fmt.Printf("matched sexpr\n")

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
	p.EnterRule(localctx, 4, jsstripParserRULE_list)

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

	p.SetState(30)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(22)
			p.Match(jsstripParserT__0)
		}
		{
			p.SetState(23)
			p.Match(jsstripParserT__1)
		}
		fmt.Printf("matched empty list\n")

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(25)
			p.Match(jsstripParserT__0)
		}
		{
			p.SetState(26)
			p.Members()
		}
		{
			p.SetState(27)
			p.Match(jsstripParserT__1)
		}
		fmt.Printf("matched list\n")

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
	p.EnterRule(localctx, 6, jsstripParserRULE_members)
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
	p.SetState(33)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<jsstripParserT__0)|(1<<jsstripParserNum)|(1<<jsstripParserId)|(1<<jsstripParserString_))) != 0) {
		{
			p.SetState(32)
			p.Sexpr()
		}

		p.SetState(35)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	fmt.Printf("members  1\n")

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

func (s *AtomContext) Id() antlr.TerminalNode {
	return s.GetToken(jsstripParserId, 0)
}

func (s *AtomContext) Num() antlr.TerminalNode {
	return s.GetToken(jsstripParserNum, 0)
}

func (s *AtomContext) String_() antlr.TerminalNode {
	return s.GetToken(jsstripParserString_, 0)
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
	p.EnterRule(localctx, 8, jsstripParserRULE_atom)

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

	p.SetState(45)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case jsstripParserId:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(39)
			p.Match(jsstripParserId)
		}
		fmt.Printf("ID ")

	case jsstripParserNum:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(41)
			p.Match(jsstripParserNum)
		}
		fmt.Printf("NUM ")

	case jsstripParserString_:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(43)
			p.Match(jsstripParserString_)
		}
		fmt.Printf("STRING ")

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}
