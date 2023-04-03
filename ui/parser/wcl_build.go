package parser

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	antlrhelp "github.com/iansmith/parigot/helper/antlr"

	"github.com/iansmith/parigot/ui/css"
	"github.com/iansmith/parigot/ui/parser/builtin"
	"github.com/iansmith/parigot/ui/parser/tree"
)

type WclBuildListener struct {
	*BasewclListener

	ClassName map[string]struct{}

	SourceCode string
}

var _ wclListener = &WclBuildListener{}

func NewWclBuildListener(sourceCode string) *WclBuildListener {
	return &WclBuildListener{
		BasewclListener: &BasewclListener{},
		ClassName:       make(map[string]struct{}),
		SourceCode:      sourceCode,
	}
}

// EnterEveryRule is called when any rule is entered.
func (s *WclBuildListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	//log.Printf("EnterEveryRule: %v", ctx.GetText())
}

// ExitEveryRule is called when any rule is exited.
func (s *WclBuildListener) ExitEveryRule(ctx antlr.ParserRuleContext) {

}

func (l *WclBuildListener) EnterProgram(c *ProgramContext) {
}

func checkDuplicateEarlySection(c *ProgramContext) bool {
	glob := c.AllGlobal()
	ext := c.AllExtern()
	imp := c.AllImport_section()
	b := bytes.Buffer{}
	result := true

	if len(glob) > 1 {
		b.WriteString("no more than one global section permitted, found global sections")
		for _, g := range glob {
			b.WriteString(fmt.Sprintf("\tline %d, column %d\n", g.GetG().LineNumber, g.GetG().ColumnNumber))
		}
		result = false
	}
	if len(ext) > 1 {
		b := bytes.Buffer{}
		b.WriteString("no more than one extern section permitted, found extern sections")
		for _, e := range ext {
			b.WriteString(fmt.Sprintf("\tline %d, column %d\n", e.GetE().LineNumber, e.GetE().ColumnNumber))
		}
		result = false
	}
	if len(imp) > 1 {
		b := bytes.Buffer{}
		b.WriteString("no more than one preamble section permitted, found preamble sections")
		for _, i := range imp {
			b.WriteString(fmt.Sprintf("\tline %d, column %d\n", i.GetSection().LineNumber, i.GetSection().ColumnNumber))
		}
		result = false
	}
	return result
}

func (l *WclBuildListener) ExitProgram(c *ProgramContext) {
	if !checkDuplicateEarlySection(c) {
		notifyError("failed due to section count",
			c.BaseParserRuleContext, c.GetParser())

	}
	imp := c.AllImport_section()
	if len(imp) == 1 && imp[0].GetSection() != nil {
		tree.GProgram.ImportSection = imp[0].GetSection()
	}
	ext := c.AllExtern()
	glob := c.AllGlobal()
	var extNode *tree.ExternSectionNode
	var globNode *tree.GlobalSectionNode
	if len(ext) == 1 {
		if ext[0].GetE() != nil {
			extNode = ext[0].GetE()
			extNode.Program = tree.GProgram
		}
	}
	if len(glob) == 1 {
		if glob[0].GetG() != nil {
			globNode = glob[0].GetG()
			globNode.Program = tree.GProgram
		}
	}
	tree.GProgram.Global = tree.NewAllGlobal(tree.GProgram, globNode, extNode)
	if c.Text_section() != nil && c.Text_section().GetSection() != nil {
		tree.GProgram.TextSection = c.Text_section().GetSection()
		tree.GProgram.TextSection.Program = tree.GProgram
		tree.GProgram.TextSection.Scope_.Parent_ = tree.GProgram.Global
	}
	if c.Doc_section() != nil && c.Doc_section().GetSection() != nil {
		tree.GProgram.DocSection = c.Doc_section().GetSection()
		tree.GProgram.DocSection.Program = tree.GProgram
		tree.GProgram.DocSection.Scope_.Parent_ = tree.GProgram.Global
	}
	if c.Event_section() != nil && c.Event_section().GetSection() != nil {
		tree.GProgram.EventSection = c.Event_section().GetSection()
		tree.GProgram.EventSection.Program = tree.GProgram
	}

	if c.Mvc_section() != nil && c.Mvc_section().GetSection() != nil {
		tree.GProgram.ModelSection = c.Mvc_section().GetSection()
		tree.GProgram.ModelSection.Program = tree.GProgram
	}
}

// Import_section
func (l *WclBuildListener) EnterImport_section(c *Import_sectionContext) {
}
func (l *WclBuildListener) ExitImport_section(c *Import_sectionContext) {
	i := tree.NewImportSectionNode()
	i.TextItem_ = c.Uninterp().GetItem()
	c.SetSection(i)
}

// Text_section
func (l *WclBuildListener) EnterText_section(c *Text_sectionContext) {
}

func (l *WclBuildListener) ExitText_section(c *Text_sectionContext) {
	raw := c.AllText_func()
	tfn := make([]*tree.TextFuncNode, len(raw))
	section := tree.NewTextSectionNode(tree.GProgram)
	count := 0
	for i, fn := range raw {
		tfn[i] = fn.GetF()
		count++
	}
	if len(tfn) > 0 {
		section.Func = tfn
		section.Scope_.TextFn = tfn
	}
	c.SetSection(section)
}

// Text_func
func (l *WclBuildListener) EnterText_func(c *Text_funcContext) {
}

func (l *WclBuildListener) ExitText_func(c *Text_funcContext) {
	fn := tree.NewTextFuncNode()
	fn.Name = c.Id().GetText()

	if c.Param_spec() != nil && c.Param_spec().GetFormal() != nil {
		fn.Param = c.Param_spec().GetFormal()
	}
	if c.Text_func_local() != nil && c.Text_func_local().GetFormal() != nil {
		fn.Local = c.Text_func_local().GetFormal()
	}
	if c.Uninterp() != nil {
		fn.Item_ = c.Uninterp().GetItem()
	}
	if c.Pre_code() != nil && c.Pre_code().GetItem() != nil {
		fn.PreCode = c.Pre_code().GetItem()
	}

	if c.Post_code() != nil && c.Post_code().GetItem() != nil {
		fn.PostCode = c.Post_code().GetItem()
	}
	fn.LineNumber = c.Id().GetSymbol().GetLine()
	fn.ColumnNumber = c.Id().GetSymbol().GetColumn()
	//log.Printf("got text fn %s, pre=%d main=%d", fn.Name, len(fn.PreCode), len(fn.Item_))
	c.SetF(fn)
}

func (l *WclBuildListener) EnterIdent(c *IdentContext) {
}
func (l *WclBuildListener) ExitIdent(c *IdentContext) {
	var line, col int
	var text string
	line, col = c.Id().GetSymbol().GetLine(), c.GetStart().GetColumn()
	text = c.Id().GetSymbol().GetText()
	id := tree.NewIdent(text, c.GetText(), line, col)
	if c.Colon() != nil {
		id.HasStartColon = true
	}

	rawC := c.AllColon_qual()
	rawD := c.AllDot_qual()
	if len(rawC) > 0 && len(rawD) > 0 {
		msg := fmt.Sprintf("identifier '%s' cannot dot and colon qualified parts",
			text)
		notifyError(msg, c.BaseParserRuleContext, c.GetParser())
		return
	}
	prev := &id.Part.Qual
	if len(rawC) > 0 {
		colonPart := make([]*tree.IdentPart, len(rawC))
		for i, p := range rawC {
			colonPart[i] = p.GetPart()
			*prev = p.GetPart()
			prev = &(p.GetPart().Qual)
		}
		c.SetId(id)
		return
	}
	if len(rawD) > 0 {
		dotPart := make([]*tree.IdentPart, len(rawD))
		for i, p := range rawD {
			dotPart[i] = p.GetPart()
			*prev = p.GetPart()
			prev = &(p.GetPart().Qual)
		}
		c.SetId(id)
		return
	}
	c.SetId(id)

}

// dot qual
func (l *WclBuildListener) EnterDot_qual(c *Dot_qualContext) {
}
func (l *WclBuildListener) ExitDot_qual(c *Dot_qualContext) {
	c.SetPart(&tree.IdentPart{
		Id:       c.Id().GetSymbol().GetText(),
		ColonSep: false,
		Qual:     nil,
	})
}

// colon qual
func (l *WclBuildListener) EnterColon_qual(c *Colon_qualContext) {
}
func (l *WclBuildListener) ExitColon_qual(c *Colon_qualContext) {
	c.SetPart(&tree.IdentPart{
		Id:       c.Id().GetSymbol().GetText(),
		ColonSep: true,
		Qual:     nil,
	})
}

// ValueRef ID
func (l *WclBuildListener) EnterValue_ref_id(c *Value_ref_idContext) {
}
func (l *WclBuildListener) ExitValue_ref_id(c *Value_ref_idContext) {
	c.SetVr(&tree.ValueRef{Id: c.Ident().GetId()})
}

// ValueRef lit
func (l *WclBuildListener) EnterValue_ref_lit(c *Value_ref_litContext) {
}
func (l *WclBuildListener) ExitValue_ref_lit(c *Value_ref_litContext) {
	lit := c.StringLit()
	line := lit.GetSymbol().GetLine()
	col := lit.GetSymbol().GetColumn()
	text := lit.GetSymbol().GetText()
	text = strings.TrimPrefix(text, "\"")
	text = strings.TrimSuffix(text, "\"")

	c.SetVr(tree.NewValueRef(nil, nil, text, line, col))
}

// ValueRef Func
func (l *WclBuildListener) EnterValue_ref_func(c *Value_ref_funcContext) {
}
func (l *WclBuildListener) ExitValue_ref_func(c *Value_ref_funcContext) {
	c.SetVr(&tree.ValueRef{FuncInvoc: c.Func_invoc().GetInvoc()})
}

// Var_subs
func (l *WclBuildListener) EnterVar_subs(c *Var_subsContext) {
}
func (l *WclBuildListener) ExitVar_subs(c *Var_subsContext) {
	var valueRef *tree.ValueRef
	if c.GetNormal() == nil {
		if c.GetGrab() == nil {
			panic("neither a normal nor grab found in variable substitution")
		}
		valueRef = c.GetGrab().GetVr()
	} else {
		valueRef = c.GetNormal().GetVr()
	}
	item := tree.NewTextValueRef(valueRef, valueRef.LineNumber, valueRef.ColumnNumber)
	c.SetItem(item)
}

// Uninterp_inner
func (l *WclBuildListener) EnterUninterp_inner(c *Uninterp_innerContext) {
}
func (l *WclBuildListener) ExitUninterp_inner(c *Uninterp_innerContext) {
	subs := c.Var_subs()
	if subs == nil {
		if c.RawText() == nil {
			log.Printf("WARN: uninterpreted text found with no content")
			return
		}
		sym := c.RawText().GetSymbol()
		c.SetItem(tree.NewTextConstant(sym.GetText(), sym.GetLine(), sym.GetColumn()))
		return
	}
	c.SetItem(c.Var_subs().GetItem())
}

// Uninterp
func (l *WclBuildListener) EnterUninterp(c *UninterpContext) {
}
func (l *WclBuildListener) ExitUninterp(c *UninterpContext) {
	result := []tree.TextItem{}
	for _, t := range c.AllUninterp_inner() {
		result = append(result, t.GetItem())
	}

	c.SetItem(result)
}

// ParamSeq
func (l *WclBuildListener) EnterParam_spec(c *Param_specContext) {
}
func (l *WclBuildListener) ExitParam_spec(c *Param_specContext) {
	all := []*tree.PFormal{}
	p := c.AllParam_pair()
	for _, pair := range p {
		all = append(all, pair.GetFormal())
	}
	c.SetFormal(all)
}

func (l *WclBuildListener) ExitParam_pair(c *Param_pairContext) {
	n := c.Id().GetText()
	t := c.Ident().GetId()
	ts := ""
	if c.TypeStarter() != nil {
		ts = c.TypeStarter().GetText()
	}
	c.SetFormal(tree.NewPFormal(n, t, ts, c.Id().GetSymbol().GetLine(), c.Id().GetSymbol().GetColumn()))
}

// Doc_tag is a full tag descriptor
func (s *WclBuildListener) EnterDoc_id(ctx *Doc_idContext) {
}

// Doc_tag is a full tag descriptor
func (s *WclBuildListener) ExitDoc_id(ctx *Doc_idContext) {
	ctx.SetS(ctx.Value_ref().GetVr())
}

// Doc_tag is a full tag descriptor
func (s *WclBuildListener) EnterDoc_tag(ctx *Doc_tagContext) {
}

// Doc_tag is a full tag descriptor
func (s *WclBuildListener) ExitDoc_tag(ctx *Doc_tagContext) {
	var docId *tree.ValueRef
	if ctx.Doc_id() != nil && ctx.Doc_id().GetS() != nil {
		docId = ctx.Doc_id().GetS()
	}
	if ctx.Value_ref() == nil {
		return
	}
	// get the string here because we might need it for err messages
	tn := ctx.Value_ref().GetVr().String()
	numClass := 0
	if ctx.Doc_class() != nil && ctx.Doc_class().GetClazz() != nil {
		numClass = len(ctx.Doc_class().GetClazz())
		vr := ctx.Doc_class().GetClazz()
		for _, ref := range vr {
			if ref.Lit == "" {
				continue
			}
			l := "." + ref.Lit
			_, ok := s.ClassName[l]
			if !ok {
				e := &tree.ErrorLoc{Filename: s.SourceCode, Line: ref.LineNumber, Col: ref.ColumnNumber}
				notifyError(fmt.Sprintf("at %s in tag '%s', class name '%s' is not defined the css files declared", e.String(), tn, ref.Lit),
					ctx.BaseParserRuleContext, ctx.GetParser())
			}
		}
	}
	cl := []*tree.ValueRef{}
	if numClass > 0 {
		cl = make([]*tree.ValueRef, numClass)
		for i, c := range ctx.Doc_class().GetClazz() {
			cl[i] = c
		}
	}
	tag, err := tree.NewDocTag(ctx.Value_ref().GetVr(), docId, cl, s.ClassName)
	if err != nil {
		spec := fmt.Sprintf("%s:%d:%d", s.SourceCode, ctx.Value_ref().GetVr().LineNumber, ctx.Value_ref().GetVr().ColumnNumber)
		notifyError(fmt.Sprintf("%s: %s", spec, err.Error()),
			ctx.BaseParserRuleContext, ctx.GetParser())
		return
	}
	ctx.SetTag(tag)

}
func notifyError(msg string, ctx *antlr.BaseParserRuleContext, parser antlr.Parser) {
	ex := antlr.NewBaseRecognitionException(msg, parser, parser.GetInputStream(), ctx)
	ctx.SetException(ex)
	parser.NotifyErrorListeners(msg, ctx.GetStart(), ex)
}

// Doc_class the part of a doc atom that looks like :foo, describing a css class
func (s *WclBuildListener) EnterDoc_class(ctx *Doc_classContext) {
}

// Doc_class the part of a doc atom that looks like :foo, describing a css class
func (s *WclBuildListener) ExitDoc_class(ctx *Doc_classContext) {
	allVR := ctx.AllValue_ref()
	result := make([]*tree.ValueRef, len(allVR))
	for i, v := range allVR {
		result[i] = v.GetVr()
	}
	ctx.SetClazz(result)
}

// Doc_section
func (s *WclBuildListener) EnterDoc_section(ctx *Doc_sectionContext) {
}

// Doc_sexpr.atom exit
func (s *WclBuildListener) ExitDoc_section(ctx *Doc_sectionContext) {
	raw := ctx.AllDoc_func()
	content := []*tree.DocFuncNode{}
	for _, r := range raw {
		if r == nil || r.GetFn() == nil {
			continue
		}
		content = append(content, r.GetFn())
	}
	if len(content) == 0 {
		return
	}
	section := tree.NewDocSectionNode(tree.GProgram, content)
	ctx.SetSection(section)
}

// func (s *WclBuildListener) EnterDoc_elem_content(ctx *Doc_elem_contentContext) {}

// func (s *WclBuildListener) ExitDoc_elem_content(ctx *Doc_elem_contentContext) {
// 	if ctx.Doc_elem_text() != nil {
// 		ctx.SetElement(&tree.DocElement{TextContent: ctx.Doc_elem_text().GetInvoc()})
// 	} else {
// 		ctx.SetElement(ctx.Doc_elem_child().GetElem())
// 	}
// }
// func (s *WclBuildListener) EnterDoc_elem_text(ctx *Doc_elem_textContext) {}

// func (s *WclBuildListener) ExitDoc_elem_text(ctx *Doc_elem_textContext) {
// 	// if ctx.Text_top() != nil && ctx.Text_top().GetItem() != nil {
// 	// 	invoc := tree.NewFuncInvoc(fmt.Sprintf("%s%04d", anonPrefix, s.anonCount), nil, -1, -1)
// 	// 	s.anonCount++
// 	// 	invoc.AnonBody = ctx.Text_top().GetItem()
// 	// 	ctx.SetInvoc(invoc)
// 	// 	return
// 	// }
// 	// if ctx.Func_invoc() != nil && ctx.Func_invoc().GetInvoc() != nil {
// 	// 	ctx.SetInvoc(ctx.Func_invoc().GetInvoc())
// 	// 	return
// 	// }
// }

func (s *WclBuildListener) EnterDoc_elem_child(ctx *Doc_elem_childContext) {}

func (s *WclBuildListener) ExitDoc_elem_child(ctx *Doc_elem_childContext) {
	raw := ctx.AllDoc_elem()
	result := make([]*tree.DocElement, len(raw))
	for i, elem := range raw {
		result[i] = elem.GetElem()
	}
	ctx.SetElem(result)
}

func (s *WclBuildListener) EnterHaveTag(ctx *HaveTagContext) {}

func (s *WclBuildListener) ExitHaveTag(ctx *HaveTagContext) {
	e := &tree.DocElement{}
	if ctx.Uninterp() != nil {
		e.TextContent = ctx.Uninterp().GetItem()
	}
	if ctx.Doc_elem_child() != nil {
		e.Child = ctx.Doc_elem_child().GetElem()
	}
	e.Tag = ctx.Doc_tag().GetTag()
	ctx.SetElem(e)
}

func (s *WclBuildListener) EnterFunc_invoc(ctx *Func_invocContext) {
}

func (s *WclBuildListener) ExitFunc_invoc(ctx *Func_invocContext) {
	actual := ctx.Func_actual_seq().GetActual()
	var name *tree.Ident
	if ctx.Ident() != nil {
		name = ctx.Ident().GetId()
	}
	var line, col int
	if ctx.Ident() != nil {
		// if ctx.Ident().GetId() == nil {
		// 	log.Printf("xxxx --- ident %+v", ctx.Ident().GetId())
		// }
		line = ctx.Ident().GetId().LineNumber
		col = ctx.Ident().GetId().ColumnNumber
	}
	invoc := tree.NewFuncInvoc(name, actual, line, col)
	ctx.SetInvoc(invoc)
}

func (s *WclBuildListener) EnterFunc_actual(ctx *Func_actualContext) {
}

func (s *WclBuildListener) ExitFunc_actual(ctx *Func_actualContext) {
	if ctx.Value_ref() != nil && ctx.Value_ref().GetVr() != nil {
		ctx.SetActual(tree.NewFuncActual(ctx.Value_ref().GetVr()))
	}
}

func (s *WclBuildListener) EnterFunc_actual_seq(ctx *Func_actual_seqContext) {
}

func (s *WclBuildListener) ExitFunc_actual_seq(ctx *Func_actual_seqContext) {
	raw := ctx.AllFunc_actual()
	result := make([]*tree.FuncActual, len(raw))
	for i, r := range raw {
		result[i] = r.GetActual()
	}
	ctx.SetActual(result)
}

func (s *WclBuildListener) EnterDoc_func_local(ctx *Doc_func_localContext) {}

func (s *WclBuildListener) ExitDoc_func_local(ctx *Doc_func_localContext) {
	if ctx.Param_spec() == nil {
		ctx.SetFormal([]*tree.PFormal{})
	} else {
		ctx.SetFormal(ctx.Param_spec().GetFormal())
	}
}
func (s *WclBuildListener) EnterDoc_func_formal(ctx *Doc_func_formalContext) {}

func (s *WclBuildListener) ExitDoc_func_formal(ctx *Doc_func_formalContext) {
	if ctx.Param_spec() == nil {
		ctx.SetFormal([]*tree.PFormal{})
	} else {
		ctx.SetFormal(ctx.Param_spec().GetFormal())
	}
}

func (s *WclBuildListener) EnterText_func_local(ctx *Text_func_localContext) {}

func (s *WclBuildListener) ExitText_func_local(ctx *Text_func_localContext) {
	if ctx.Param_spec() == nil {
		ctx.SetFormal([]*tree.PFormal{})
	} else {
		ctx.SetFormal(ctx.Param_spec().GetFormal())
	}
}
func (s *WclBuildListener) EnterDoc_func(ctx *Doc_funcContext) {}

func (s *WclBuildListener) ExitDoc_func(ctx *Doc_funcContext) {
	dfunc := ctx.Doc_func_post().GetFn() // not Gfunk
	dfunc.Name = ctx.Id().GetText()
	dfunc.LineNumber = ctx.Id().GetSymbol().GetLine()
	dfunc.ColumnNumber = ctx.Id().GetSymbol().GetColumn()
	ctx.SetFn(dfunc)
}

func (s *WclBuildListener) EnterDoc_func_post(ctx *Doc_func_postContext) {}

func (s *WclBuildListener) ExitDoc_func_post(ctx *Doc_func_postContext) {
	var f, l []*tree.PFormal
	var pre, post []tree.TextItem
	var elem *tree.DocElement

	if ctx.Doc_func_formal() != nil {
		f = ctx.Doc_func_formal().GetFormal()
	}
	if ctx.Doc_func_local() != nil {
		l = ctx.Doc_func_local().GetFormal()
	}
	if ctx.Pre_code() != nil {
		pre = ctx.Pre_code().GetItem()
	}
	if ctx.Post_code() != nil {
		post = ctx.Post_code().GetItem()
	}
	if ctx.Doc_elem() != nil && ctx.Doc_elem().GetElem() != nil {
		elem = ctx.Doc_elem().GetElem()
	}
	ctx.SetFn(tree.NewDocFuncNode("", f, l, elem,
		pre, post))

}

func (s *WclBuildListener) EnterGlobal(ctx *GlobalContext) {
}

func (s *WclBuildListener) ExitGlobal(ctx *GlobalContext) {
	g := tree.NewGlobalSectionNode(tree.GProgram,
		ctx.Global().GetSymbol().GetLine(), ctx.Global().GetSymbol().GetColumn())
	g.Var = ctx.Param_spec().GetFormal()
	ctx.SetG(g)
}

func (s *WclBuildListener) EnterExtern(ctx *ExternContext) {
}

func (s *WclBuildListener) ExitExtern(ctx *ExternContext) {
	raw := ctx.AllId()
	result := make([]string, len(raw))
	for i, name := range raw {
		result[i] = name.GetText()
	}
	e := tree.NewExternSectionNode(tree.GProgram,
		ctx.Extern().GetSymbol().GetLine(),
		ctx.Extern().GetSymbol().GetColumn())
	e.Name = result
	ctx.SetE(e)
}

func (s *WclBuildListener) EnterPre_code(ctx *Pre_codeContext) {
}

func (s *WclBuildListener) ExitPre_code(ctx *Pre_codeContext) {
	if ctx.Uninterp() != nil {
		ctx.SetItem(ctx.Uninterp().GetItem())
	}
}
func (s *WclBuildListener) EnterPost_code(ctx *Post_codeContext) {
}

func (s *WclBuildListener) ExitPost_code(ctx *Post_codeContext) {
	if ctx.Uninterp() != nil {
		ctx.SetItem(ctx.Uninterp().GetItem())
	}
}

func (s *WclBuildListener) EnterCss_filespec(ctx *Css_filespecContext) {
}

func (s *WclBuildListener) ExitCss_filespec(ctx *Css_filespecContext) {
	raw := ctx.StringLit().GetText()
	path := raw[1 : len(raw)-1]
	className, err := css.ReadCSS(s.SourceCode, path)
	if err != nil {
		log.Fatalf("CSS file %s triggered an error: %v", path, err)
	}
	for k := range className {
		s.ClassName[k] = struct{}{}
	}
}

func (s *WclBuildListener) EnterSelector(ctx *SelectorContext) {
}

func (s *WclBuildListener) ExitSelector(ctx *SelectorContext) {
	if ctx.GetClass() != nil && ctx.GetClass().GetVr() != nil {
		sel := &tree.Selector{Class: ctx.GetClass().GetVr()}
		ctx.SetSel(sel)
		return
	}
	if ctx.GetId() != nil && ctx.GetId().GetVr() != nil {
		sel := &tree.Selector{Id: ctx.GetId().GetVr()}
		ctx.SetSel(sel)
		return
	}

}

func (s *WclBuildListener) EnterEvent_call(ctx *Event_callContext) {
}

func (s *WclBuildListener) ExitEvent_call(ctx *Event_callContext) {
	b := ctx.Arrow() != nil
	f := ctx.Func_invoc().GetInvoc()
	if b {
		f.Builtin = true
		s.checkBuiltinWithSingleParam(ctx, f)
	}
	ctx.SetInvoc(f)
}

func (s *WclBuildListener) checkBuiltinWithSingleParam(ctx *Event_callContext, f *tree.FuncInvoc) {
	// we cannot check anything related to non-builtins, other than what is checked in the name
	if !f.Builtin {
		return
	}
	checkerFn, err := builtin.GetBuiltinChecker(f.Name.String())
	if err != nil {
		notifyError(err.Error(),
			ctx.BaseParserRuleContext, ctx.GetParser())
		return
	}
	detail := fmt.Sprintf("%s:%d:%d", s.SourceCode, f.LineNumber, f.ColumnNumber)
	ok, msg := checkerFn(f)
	if !ok {
		notifyError(detail+" "+msg, ctx.BaseParserRuleContext, ctx.GetParser())
		return
	}
}

func (s *WclBuildListener) EnterEvent_spec(ctx *Event_specContext) {
}

func (s *WclBuildListener) ExitEvent_spec(ctx *Event_specContext) {
	ctx.SetSpec(&tree.EventSpec{
		Selector:  ctx.Selector().GetSel(),
		EventName: ctx.Id().GetText(),
		Function:  ctx.Event_call().GetInvoc(),
	})
}

func (s *WclBuildListener) EnterEvent_section(ctx *Event_sectionContext) {
}

func (s *WclBuildListener) ExitEvent_section(ctx *Event_sectionContext) {
	raw := ctx.AllEvent_spec()
	e := make([]*tree.EventSpec, len(raw))
	for i, s := range raw {
		e[i] = s.GetSpec()
	}
	ctx.SetSection(&tree.EventSectionNode{Spec: e})
}

func (s *WclBuildListener) EnterFilename_seq(ctx *Filename_seqContext) {
}

func (s *WclBuildListener) ExitFilename_seq(ctx *Filename_seqContext) {
	raw := ctx.AllStringLit()
	rest := make([]string, len(raw))
	for i, s := range raw {
		quoted := s.GetText()
		notQuoted := strings.TrimPrefix(quoted, "\"")
		notQuoted = strings.TrimSuffix(notQuoted, "\"")
		rest[i] = notQuoted
	}
	ctx.SetSeq(rest)
}

func (s *WclBuildListener) EnterModel_decl(ctx *Model_declContext) {
}

func (s *WclBuildListener) ExitModel_decl(ctx *Model_declContext) {
	modelDecl := tree.NewModelDecl()
	modelDecl.Path = ctx.Filename_seq().GetSeq()
	modelDecl.Name = ctx.GetId1().GetText()
	ctx.SetDecl(modelDecl)
}

func (s *WclBuildListener) EnterView_decl(ctx *View_declContext) {
}

func (s *WclBuildListener) ExitView_decl(ctx *View_declContext) {
	vdecl := tree.NewViewDecl()
	fn := ctx.Doc_func_post().GetFn()
	vdecl.DocFn = fn
	vdecl.DocFn.Name = ctx.Id().GetText()
	ctx.SetVdecl(vdecl)
}

func (s *WclBuildListener) EnterMvc_section(ctx *Mvc_sectionContext) {
}

func (s *WclBuildListener) ExitMvc_section(ctx *Mvc_sectionContext) {
	section := tree.NewMvcSectionNode(tree.GProgram)
	raw := ctx.AllModel_decl()
	decl := make([]*tree.ModelDecl, len(raw))
	for i, mod := range raw {
		decl[i] = mod.GetDecl()
	}
	section.ModelDecl = decl

	rawView := ctx.AllView_decl()
	vdecl := make([]*tree.ViewDecl, len(rawView))
	for i, v := range rawView {
		vdecl[i] = v.GetVdecl()
	}
	section.ViewDecl = vdecl

	_, err := antlrhelp.ParseModelSection(s.SourceCode, "", section)
	if err != nil {
		notifyError(fmt.Sprintf("failed due to problem with proto files (originating at '%+v')", section.ModelDecl[0].Path),
			ctx.BaseParserRuleContext, ctx.GetParser())
	}

	ctx.SetSection(section)
}
