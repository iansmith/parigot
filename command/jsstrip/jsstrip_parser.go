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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 16, 43, 4,
	2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 3, 2, 3, 2, 5, 2, 13, 10,
	2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 21, 10, 3, 3, 4, 6, 4, 24,
	10, 4, 13, 4, 14, 4, 25, 3, 5, 3, 5, 5, 5, 30, 10, 5, 3, 5, 5, 5, 33, 10,
	5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 5, 5, 41, 10, 5, 3, 5, 2, 2, 6,
	2, 4, 6, 8, 2, 2, 2, 49, 2, 12, 3, 2, 2, 2, 4, 20, 3, 2, 2, 2, 6, 23, 3,
	2, 2, 2, 8, 40, 3, 2, 2, 2, 10, 13, 5, 4, 3, 2, 11, 13, 5, 8, 5, 2, 12,
	10, 3, 2, 2, 2, 12, 11, 3, 2, 2, 2, 13, 3, 3, 2, 2, 2, 14, 15, 7, 5, 2,
	2, 15, 21, 7, 6, 2, 2, 16, 17, 7, 5, 2, 2, 17, 18, 5, 6, 4, 2, 18, 19,
	7, 6, 2, 2, 19, 21, 3, 2, 2, 2, 20, 14, 3, 2, 2, 2, 20, 16, 3, 2, 2, 2,
	21, 5, 3, 2, 2, 2, 22, 24, 5, 2, 2, 2, 23, 22, 3, 2, 2, 2, 24, 25, 3, 2,
	2, 2, 25, 23, 3, 2, 2, 2, 25, 26, 3, 2, 2, 2, 26, 7, 3, 2, 2, 2, 27, 29,
	7, 9, 2, 2, 28, 30, 7, 11, 2, 2, 29, 28, 3, 2, 2, 2, 29, 30, 3, 2, 2, 2,
	30, 32, 3, 2, 2, 2, 31, 33, 7, 12, 2, 2, 32, 31, 3, 2, 2, 2, 32, 33, 3,
	2, 2, 2, 33, 41, 3, 2, 2, 2, 34, 41, 7, 8, 2, 2, 35, 41, 7, 16, 2, 2, 36,
	41, 7, 14, 2, 2, 37, 41, 7, 13, 2, 2, 38, 41, 7, 10, 2, 2, 39, 41, 7, 15,
	2, 2, 40, 27, 3, 2, 2, 2, 40, 34, 3, 2, 2, 2, 40, 35, 3, 2, 2, 2, 40, 36,
	3, 2, 2, 2, 40, 37, 3, 2, 2, 2, 40, 38, 3, 2, 2, 2, 40, 39, 3, 2, 2, 2,
	41, 9, 3, 2, 2, 2, 8, 12, 20, 25, 29, 32, 40,
}
var literalNames = []string{
	"", "", "", "'('", "')'", "'\"'",
}
var symbolicNames = []string{
	"", "Whitespace", "Comment", "Lparen", "Rparen", "Quote", "Num", "Ident",
	"HexPointer", "Offset", "Align", "ConstAnnotation", "BlockAnnotation",
	"ConstValue", "QuotedString",
}

var ruleNames = []string{
	"sexpr", "list", "members", "atom",
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
	jsstripParserWhitespace      = 1
	jsstripParserComment         = 2
	jsstripParserLparen          = 3
	jsstripParserRparen          = 4
	jsstripParserQuote           = 5
	jsstripParserNum             = 6
	jsstripParserIdent           = 7
	jsstripParserHexPointer      = 8
	jsstripParserOffset          = 9
	jsstripParserAlign           = 10
	jsstripParserConstAnnotation = 11
	jsstripParserBlockAnnotation = 12
	jsstripParserConstValue      = 13
	jsstripParserQuotedString    = 14
)

// jsstripParser rules.
const (
	jsstripParserRULE_sexpr   = 0
	jsstripParserRULE_list    = 1
	jsstripParserRULE_members = 2
	jsstripParserRULE_atom    = 3
)

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
	p.EnterRule(localctx, 0, jsstripParserRULE_sexpr)

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

	p.SetState(10)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case jsstripParserLparen:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(8)
			p.List()
		}

	case jsstripParserNum, jsstripParserIdent, jsstripParserHexPointer, jsstripParserConstAnnotation, jsstripParserBlockAnnotation, jsstripParserConstValue, jsstripParserQuotedString:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(9)
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
	p.EnterRule(localctx, 2, jsstripParserRULE_list)

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

	p.SetState(18)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(12)
			p.Match(jsstripParserLparen)
		}
		{
			p.SetState(13)
			p.Match(jsstripParserRparen)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(14)
			p.Match(jsstripParserLparen)
		}
		{
			p.SetState(15)
			p.Members()
		}
		{
			p.SetState(16)
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
	p.EnterRule(localctx, 4, jsstripParserRULE_members)
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
	p.SetState(21)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<jsstripParserLparen)|(1<<jsstripParserNum)|(1<<jsstripParserIdent)|(1<<jsstripParserHexPointer)|(1<<jsstripParserConstAnnotation)|(1<<jsstripParserBlockAnnotation)|(1<<jsstripParserConstValue)|(1<<jsstripParserQuotedString))) != 0) {
		{
			p.SetState(20)
			p.Sexpr()
		}

		p.SetState(23)
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

func (s *AtomContext) BlockAnnotation() antlr.TerminalNode {
	return s.GetToken(jsstripParserBlockAnnotation, 0)
}

func (s *AtomContext) ConstAnnotation() antlr.TerminalNode {
	return s.GetToken(jsstripParserConstAnnotation, 0)
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
	p.EnterRule(localctx, 6, jsstripParserRULE_atom)
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

	p.SetState(38)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case jsstripParserIdent:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(25)
			p.Match(jsstripParserIdent)
		}
		p.SetState(27)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == jsstripParserOffset {
			{
				p.SetState(26)
				p.Match(jsstripParserOffset)
			}

		}
		p.SetState(30)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == jsstripParserAlign {
			{
				p.SetState(29)
				p.Match(jsstripParserAlign)
			}

		}

	case jsstripParserNum:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(32)
			p.Match(jsstripParserNum)
		}

	case jsstripParserQuotedString:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(33)
			p.Match(jsstripParserQuotedString)
		}

	case jsstripParserBlockAnnotation:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(34)
			p.Match(jsstripParserBlockAnnotation)
		}

	case jsstripParserConstAnnotation:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(35)
			p.Match(jsstripParserConstAnnotation)
		}

	case jsstripParserHexPointer:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(36)
			p.Match(jsstripParserHexPointer)
		}

	case jsstripParserConstValue:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(37)
			p.Match(jsstripParserConstValue)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}
