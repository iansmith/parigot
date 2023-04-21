// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package main // sexpr
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

type sexprParser struct {
	*antlr.BaseParser
}

var sexprParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func sexprParserInit() {
	staticData := &sexprParserStaticData
	staticData.literalNames = []string{
		"", "", "", "", "", "", "'('", "')'", "'.'",
	}
	staticData.symbolicNames = []string{
		"", "STRING", "WHITESPACE", "NUMBER", "SYMBOL", "COMMENT_NUM", "LPAREN",
		"RPAREN", "DOT",
	}
	staticData.ruleNames = []string{
		"sexpr", "item", "list_", "atom",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 8, 39, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 1, 0, 5, 0,
		10, 8, 0, 10, 0, 12, 0, 13, 9, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 26, 8, 1, 1, 2, 1, 2, 5, 2, 30, 8, 2,
		10, 2, 12, 2, 33, 9, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 0, 0, 4, 0, 2, 4,
		6, 0, 1, 3, 0, 1, 1, 3, 5, 8, 8, 38, 0, 11, 1, 0, 0, 0, 2, 25, 1, 0, 0,
		0, 4, 27, 1, 0, 0, 0, 6, 36, 1, 0, 0, 0, 8, 10, 3, 2, 1, 0, 9, 8, 1, 0,
		0, 0, 10, 13, 1, 0, 0, 0, 11, 9, 1, 0, 0, 0, 11, 12, 1, 0, 0, 0, 12, 14,
		1, 0, 0, 0, 13, 11, 1, 0, 0, 0, 14, 15, 5, 0, 0, 1, 15, 1, 1, 0, 0, 0,
		16, 26, 3, 6, 3, 0, 17, 26, 3, 4, 2, 0, 18, 19, 5, 6, 0, 0, 19, 20, 3,
		2, 1, 0, 20, 21, 5, 8, 0, 0, 21, 22, 3, 2, 1, 0, 22, 23, 5, 7, 0, 0, 23,
		24, 6, 1, -1, 0, 24, 26, 1, 0, 0, 0, 25, 16, 1, 0, 0, 0, 25, 17, 1, 0,
		0, 0, 25, 18, 1, 0, 0, 0, 26, 3, 1, 0, 0, 0, 27, 31, 5, 6, 0, 0, 28, 30,
		3, 2, 1, 0, 29, 28, 1, 0, 0, 0, 30, 33, 1, 0, 0, 0, 31, 29, 1, 0, 0, 0,
		31, 32, 1, 0, 0, 0, 32, 34, 1, 0, 0, 0, 33, 31, 1, 0, 0, 0, 34, 35, 5,
		7, 0, 0, 35, 5, 1, 0, 0, 0, 36, 37, 7, 0, 0, 0, 37, 7, 1, 0, 0, 0, 3, 11,
		25, 31,
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

// sexprParserInit initializes any static state used to implement sexprParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewsexprParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func SexprParserInit() {
	staticData := &sexprParserStaticData
	staticData.once.Do(sexprParserInit)
}

// NewsexprParser produces a new parser instance for the optional input antlr.TokenStream.
func NewsexprParser(input antlr.TokenStream) *sexprParser {
	SexprParserInit()
	this := new(sexprParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &sexprParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "java-escape"

	return this
}

// sexprParser tokens.
const (
	sexprParserEOF         = antlr.TokenEOF
	sexprParserSTRING      = 1
	sexprParserWHITESPACE  = 2
	sexprParserNUMBER      = 3
	sexprParserSYMBOL      = 4
	sexprParserCOMMENT_NUM = 5
	sexprParserLPAREN      = 6
	sexprParserRPAREN      = 7
	sexprParserDOT         = 8
)

// sexprParser rules.
const (
	sexprParserRULE_sexpr = 0
	sexprParserRULE_item  = 1
	sexprParserRULE_list_ = 2
	sexprParserRULE_atom  = 3
)

// ISexprContext is an interface to support dynamic dispatch.
type ISexprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetItem_ returns the item_ attribute.
	GetItem_() []*Item

	// SetItem_ sets the item_ attribute.
	SetItem_([]*Item)

	// IsSexprContext differentiates from other interfaces.
	IsSexprContext()
}

type SexprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	item_  []*Item
}

func NewEmptySexprContext() *SexprContext {
	var p = new(SexprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = sexprParserRULE_sexpr
	return p
}

func (*SexprContext) IsSexprContext() {}

func NewSexprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SexprContext {
	var p = new(SexprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = sexprParserRULE_sexpr

	return p
}

func (s *SexprContext) GetParser() antlr.Parser { return s.parser }

func (s *SexprContext) GetItem_() []*Item { return s.item_ }

func (s *SexprContext) SetItem_(v []*Item) { s.item_ = v }

func (s *SexprContext) EOF() antlr.TerminalNode {
	return s.GetToken(sexprParserEOF, 0)
}

func (s *SexprContext) AllItem() []IItemContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IItemContext); ok {
			len++
		}
	}

	tst := make([]IItemContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IItemContext); ok {
			tst[i] = t.(IItemContext)
			i++
		}
	}

	return tst
}

func (s *SexprContext) Item(i int) IItemContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IItemContext); ok {
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

	return t.(IItemContext)
}

func (s *SexprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SexprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SexprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sexprListener); ok {
		listenerT.EnterSexpr(s)
	}
}

func (s *SexprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sexprListener); ok {
		listenerT.ExitSexpr(s)
	}
}

func (s *SexprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case sexprVisitor:
		return t.VisitSexpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *sexprParser) Sexpr() (localctx ISexprContext) {
	this := p
	_ = this

	localctx = NewSexprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, sexprParserRULE_sexpr)
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
	p.SetState(11)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&378) != 0 {
		{
			p.SetState(8)
			p.Item()
		}

		p.SetState(13)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(14)
		p.Match(sexprParserEOF)
	}

	return localctx
}

// IItemContext is an interface to support dynamic dispatch.
type IItemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetItem_ returns the item_ attribute.
	GetItem_() *Item

	// SetItem_ sets the item_ attribute.
	SetItem_(*Item)

	// IsItemContext differentiates from other interfaces.
	IsItemContext()
}

type ItemContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	item_  *Item
}

func NewEmptyItemContext() *ItemContext {
	var p = new(ItemContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = sexprParserRULE_item
	return p
}

func (*ItemContext) IsItemContext() {}

func NewItemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ItemContext {
	var p = new(ItemContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = sexprParserRULE_item

	return p
}

func (s *ItemContext) GetParser() antlr.Parser { return s.parser }

func (s *ItemContext) GetItem_() *Item { return s.item_ }

func (s *ItemContext) SetItem_(v *Item) { s.item_ = v }

func (s *ItemContext) Atom() IAtomContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAtomContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *ItemContext) List_() IList_Context {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IList_Context); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IList_Context)
}

func (s *ItemContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(sexprParserLPAREN, 0)
}

func (s *ItemContext) AllItem() []IItemContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IItemContext); ok {
			len++
		}
	}

	tst := make([]IItemContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IItemContext); ok {
			tst[i] = t.(IItemContext)
			i++
		}
	}

	return tst
}

func (s *ItemContext) Item(i int) IItemContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IItemContext); ok {
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

	return t.(IItemContext)
}

func (s *ItemContext) DOT() antlr.TerminalNode {
	return s.GetToken(sexprParserDOT, 0)
}

func (s *ItemContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(sexprParserRPAREN, 0)
}

func (s *ItemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ItemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ItemContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sexprListener); ok {
		listenerT.EnterItem(s)
	}
}

func (s *ItemContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sexprListener); ok {
		listenerT.ExitItem(s)
	}
}

func (s *ItemContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case sexprVisitor:
		return t.VisitItem(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *sexprParser) Item() (localctx IItemContext) {
	this := p
	_ = this

	localctx = NewItemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, sexprParserRULE_item)

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

	p.SetState(25)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(16)
			p.Atom()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(17)
			p.List_()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(18)
			p.Match(sexprParserLPAREN)
		}
		{
			p.SetState(19)
			p.Item()
		}
		{
			p.SetState(20)
			p.Match(sexprParserDOT)
		}
		{
			p.SetState(21)
			p.Item()
		}
		{
			p.SetState(22)
			p.Match(sexprParserRPAREN)
		}
		panic("was not expecting dotted pair")

	}

	return localctx
}

// IList_Context is an interface to support dynamic dispatch.
type IList_Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetList returns the list attribute.
	GetList() []*Item

	// SetList sets the list attribute.
	SetList([]*Item)

	// IsList_Context differentiates from other interfaces.
	IsList_Context()
}

type List_Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	list   []*Item
}

func NewEmptyList_Context() *List_Context {
	var p = new(List_Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = sexprParserRULE_list_
	return p
}

func (*List_Context) IsList_Context() {}

func NewList_Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *List_Context {
	var p = new(List_Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = sexprParserRULE_list_

	return p
}

func (s *List_Context) GetParser() antlr.Parser { return s.parser }

func (s *List_Context) GetList() []*Item { return s.list }

func (s *List_Context) SetList(v []*Item) { s.list = v }

func (s *List_Context) LPAREN() antlr.TerminalNode {
	return s.GetToken(sexprParserLPAREN, 0)
}

func (s *List_Context) RPAREN() antlr.TerminalNode {
	return s.GetToken(sexprParserRPAREN, 0)
}

func (s *List_Context) AllItem() []IItemContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IItemContext); ok {
			len++
		}
	}

	tst := make([]IItemContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IItemContext); ok {
			tst[i] = t.(IItemContext)
			i++
		}
	}

	return tst
}

func (s *List_Context) Item(i int) IItemContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IItemContext); ok {
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

	return t.(IItemContext)
}

func (s *List_Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *List_Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *List_Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sexprListener); ok {
		listenerT.EnterList_(s)
	}
}

func (s *List_Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sexprListener); ok {
		listenerT.ExitList_(s)
	}
}

func (s *List_Context) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case sexprVisitor:
		return t.VisitList_(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *sexprParser) List_() (localctx IList_Context) {
	this := p
	_ = this

	localctx = NewList_Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, sexprParserRULE_list_)
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
		p.SetState(27)
		p.Match(sexprParserLPAREN)
	}
	p.SetState(31)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&378) != 0 {
		{
			p.SetState(28)
			p.Item()
		}

		p.SetState(33)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(34)
		p.Match(sexprParserRPAREN)
	}

	return localctx
}

// IAtomContext is an interface to support dynamic dispatch.
type IAtomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetAtom_ returns the atom_ attribute.
	GetAtom_() *Atom

	// SetAtom_ sets the atom_ attribute.
	SetAtom_(*Atom)

	// IsAtomContext differentiates from other interfaces.
	IsAtomContext()
}

type AtomContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	atom_  *Atom
}

func NewEmptyAtomContext() *AtomContext {
	var p = new(AtomContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = sexprParserRULE_atom
	return p
}

func (*AtomContext) IsAtomContext() {}

func NewAtomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AtomContext {
	var p = new(AtomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = sexprParserRULE_atom

	return p
}

func (s *AtomContext) GetParser() antlr.Parser { return s.parser }

func (s *AtomContext) GetAtom_() *Atom { return s.atom_ }

func (s *AtomContext) SetAtom_(v *Atom) { s.atom_ = v }

func (s *AtomContext) STRING() antlr.TerminalNode {
	return s.GetToken(sexprParserSTRING, 0)
}

func (s *AtomContext) SYMBOL() antlr.TerminalNode {
	return s.GetToken(sexprParserSYMBOL, 0)
}

func (s *AtomContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(sexprParserNUMBER, 0)
}

func (s *AtomContext) DOT() antlr.TerminalNode {
	return s.GetToken(sexprParserDOT, 0)
}

func (s *AtomContext) COMMENT_NUM() antlr.TerminalNode {
	return s.GetToken(sexprParserCOMMENT_NUM, 0)
}

func (s *AtomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AtomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AtomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sexprListener); ok {
		listenerT.EnterAtom(s)
	}
}

func (s *AtomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sexprListener); ok {
		listenerT.ExitAtom(s)
	}
}

func (s *AtomContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case sexprVisitor:
		return t.VisitAtom(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *sexprParser) Atom() (localctx IAtomContext) {
	this := p
	_ = this

	localctx = NewAtomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, sexprParserRULE_atom)
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
		p.SetState(36)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&314) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}
