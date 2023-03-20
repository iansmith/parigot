package parser

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/iansmith/parigot/ui/css"
	"github.com/iansmith/parigot/ui/parser/tree"
)

const anonPrefix = "_anon"

type WclBuildListener struct {
	*BasewclListener

	anonCount int

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
		notifyError(fmt.Sprintf("failed due to section count'"),
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
		}
	}
	if len(glob) == 1 {
		if glob[0].GetG() != nil {
			globNode = glob[0].GetG()
		}
	}
	tree.GProgram.Global = tree.NewAllGlobal(globNode, extNode)

	if c.Text_section() != nil && c.Text_section().GetSection() != nil {
		tree.GProgram.TextSection = c.Text_section().GetSection()
		tree.GProgram.TextSection.Program = tree.GProgram
	}
	if c.Doc_section() != nil && c.Doc_section().GetSection() != nil {
		tree.GProgram.DocSection = c.Doc_section().GetSection()
		tree.GProgram.TextSection.Program = tree.GProgram
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
	for i, fn := range raw {
		tfn[i] = fn.GetF()
	}
	section.Func = tfn
	c.SetSection(section)
}

// Text_func
func (l *WclBuildListener) EnterText_func(c *Text_funcContext) {
}

func (l *WclBuildListener) ExitText_func(c *Text_funcContext) {
	// fn := tree.NewTextFuncNode()
	// fn.Name = c.Id().GetText()
	// if c.Param_spec() != nil && c.Param_spec().GetFormal() != nil {
	// 	fn.Param = c.Param_spec().GetFormal()
	// }
	// if c.Text_func_local() != nil && c.Text_func_local().GetFormal() != nil {
	// 	fn.Local = c.Text_func_local().GetFormal()
	// }
	// fn.Item_ = c.Text_top().GetItem()

	// if c.Post_code() != nil && c.Post_code().GetItem() != nil {
	// 	fn.PostCode = c.Post_code().GetItem()
	// }
	// c.SetF(fn)
}

// Text_top
// func (l *WclBuildListener) EnterText_top(c *Text_topContext) {
// 	//nothing to do
// }

// func (l *WclBuildListener) ExitText_top(c *Text_topContext) {
// 	if c.Text_content() != nil {
// 		c.SetItem(c.Text_content().GetItem())
// 	}
// }

// Text_content
func (l *WclBuildListener) EnterText_content(c *Text_contentContext) {
}

func (l *WclBuildListener) ExitText_content(c *Text_contentContext) {
	result := []tree.TextItem{}
	for _, t := range c.AllRaw_text_or_sub() {
		result = append(result, t.GetItem()...)
	}
	c.SetItem(result)
}

// raw_text_or_sub
func (l *WclBuildListener) EnterRaw_text_or_sub(c *Raw_text_or_subContext) {
}

func (l *WclBuildListener) ExitRaw_text_or_sub(c *Raw_text_or_subContext) {
	if c.RawText() != nil {
		r := []tree.TextItem{tree.NewTextConstant(c.RawText().GetText(),
			c.RawText().GetSymbol().GetLine(), c.RawText().GetSymbol().GetColumn())}
		c.SetItem(r)
		return
	}
	if c.Var_subs() != nil && c.Var_subs().GetItem() != nil {
		c.SetItem(c.Var_subs().GetItem())
		return
	}
}

// // Text_content_inner.RawText
// func (l *WclBuildListener) EnterRawText(c *RawTextContext) {
// 	//nothing to do
// }
// func (l *WclBuildListener) ExitRawText(c *RawTextContext) {
// 	c.SetItem([]tree.TextItem{tree.NewTextConstant(c.ContentRawText().GetText())})
// }

// Text_content_inner.VarSub
// func (l *WclBuildListener) EnterVarSub(c *VarSubContext) {
// 	//nothing to do
// }
// func (l *WclBuildListener) ExitVarSub(c *VarSubContext) {
// 	c.SetItem(c.Var_subs().GetItem())
// }

func (l *WclBuildListener) EnterIdent(c *IdentContext) {
}
func (l *WclBuildListener) ExitIdent(c *IdentContext) {
	//var prev *tree.IdentPart
	// hasDot := false // get a colonoscopy after 45 every 5 years
	// if c.Dot() != nil || c.GrabDot() != nil {
	// 	hasDot = true
	// }
	var line, col int
	if c.Id() != nil {
		line, col = c.Id().GetSymbol().GetLine(), c.GetStart().GetColumn()
	}
	var text string
	if c.Id() != nil {
		text = c.Id().GetText()
	}
	id := tree.NewIdent(text, false, c.GetText(), line, col)
	// if c.Dot_qual() != nil {
	// 	// raw := c.AllDot_qual()
	// 	// for _, r := range raw {
	// 	// 	current := &tree.IdentPart{Qual: r.GetPart()}
	// 	// 	prev.Qual = current
	// 	// 	prev = current
	// 	// }
	// }
	// if c.Colon_qual() != nil {
	// 	// raw := c.AllColon_qual()
	// 	// for _, r := range raw {
	// 	// 	current := &tree.IdentPart{Qual: r.GetPart(), ColonSep: true}
	// 	// 	prev.Qual = current
	// 	// 	prev = current
	// 	// }
	// }

	c.SetId(id)
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
	vr := c.Value_ref().GetVr()
	c.SetItem([]tree.TextItem{
		tree.NewTextValueRef(vr, vr.LineNumber, vr.ColumnNumber),
	})
}

// sub
// func (l *WclBuildListener) EnterSub(c *SubContext) {
// }
// func (l *WclBuildListener) ExitSub(c *SubContext) {
// 	if c.Func_invoc_var() != nil {
// 		log.Printf("sub found invoc %+v", c.Func_invoc_var().GetInvoc())
// 		c.SetItem(tree.NewTextInvoc(c.Func_invoc_var().GetInvoc()))
// 		return
// 	}
// 	//normal case, just a var
// 	c.SetItem(tree.NewTextVar(c.VarId().GetText()))
// }

// UninterpRawText
// func (l *WclBuildListener) EnterUninterpRawText(c *UninterpRawTextContext) {
// }
// func (l *WclBuildListener) ExitUninterpRawText(c *UninterpRawTextContext) {
// 	c.SetItem([]tree.TextItem{tree.NewTextConstant(c.UninterpRawText().GetText())})
// }

// UninterpNested
// func (l *WclBuildListener) EnterUninterpNested(c *UninterpNestedContext) {
// }
// func (l *WclBuildListener) ExitUninterpNested(c *UninterpNestedContext) {
// 	r := append([]tree.TextItem{tree.NewTextConstant("{")}, c.Uninterp().GetItem()...)
// 	r = append(r, tree.NewTextConstant("}"))
// 	c.SetItem(r)

// }

// Var text
// func (l *WclBuildListener) EnterUninterp_var(c *Uninterp_varContext) {
// }
// func (l *WclBuildListener) ExitUninterp_var(c *Uninterp_varContext) {
// 	v := tree.NewTextVar(c.VarId().GetText())
// 	c.SetItem([]tree.TextItem{v})
// }

// Var inside uninterp
// func (l *WclBuildListener) EnterUninterpVar(c *UninterpVarContext) {
// }
// func (l *WclBuildListener) ExitUninterpVar(c *UninterpVarContext) {
// 	c.SetItem(c.Uninterp_var().GetItem())
// }

// Uninterp
func (l *WclBuildListener) EnterUninterp(c *UninterpContext) {
}
func (l *WclBuildListener) ExitUninterp(c *UninterpContext) {
	result := []tree.TextItem{}
	for _, t := range c.AllUninterp_inner() {
		result = append(result, t.GetItem()...)
	}
	for _, t := range result {
		log.Printf("--%s--", t.String())
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

func (l *WclBuildListener) EnterSimple_or_model_param(c *Simple_or_model_paramContext) {
}
func (l *WclBuildListener) ExitSimple_or_model_param(c *Simple_or_model_paramContext) {
	if c.GetId1() != nil {
		c.SetT(tree.NewTypeDeclSimple(c.GetId1().GetText()))
		return
	}
	if c.GetId2() != nil {
		raw := c.GetId2().GetText()
		part := strings.Split(raw, ".")
		if len(part) != 2 {
			notifyError(fmt.Sprintf("'%s' is not a valid reference to a model message, should be like ':model.message'", raw),
				c.BaseParserRuleContext, c.GetParser())
			return
		}

		c.SetT(tree.NewTypeDeclModel(part[0], part[1]))
		return
	}
}

func (l *WclBuildListener) ExitParam_pair(c *Param_pairContext) {
	n := c.Id().GetText()
	t := c.Ident().GetId()
	ts := ""
	if c.TypeStarter() != nil {
		ts = c.TypeStarter().GetText()
	}
	c.SetFormal(tree.NewPFormal(n, t, ts))
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
				notifyError(fmt.Sprintf("[%s:%d:%d] in tag '%s', class name '%s' is not defined the css files declared", s.SourceCode, ref.LineNumber, ref.ColumnNumber, ref.Lit, tn),
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
	tag, err := tree.NewDocTag(ctx.Value_ref().GetVr(), docId, cl)
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
	content := make([]*tree.DocFuncNode, len(raw))
	for i, r := range raw {
		content[i] = r.GetFn()
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
	ctx.SetElem(&tree.DocElement{Child: result})
}

func (s *WclBuildListener) EnterHaveTag(ctx *HaveTagContext) {}

func (s *WclBuildListener) ExitHaveTag(ctx *HaveTagContext) {
	// elem := &tree.DocElement{Tag: ctx.Doc_tag().GetTag()}
	// elemContent := ctx.Doc_elem_content()
	// if elemContent != nil {
	// 	other := elemContent.GetElement()
	// 	// might be either, but won't be both
	// 	elem.Child = other.Child
	// 	elem.TextContent = other.TextContent
	// }
	// ctx.SetElem(elem)
}

// func (s *WclBuildListener) EnterHaveVar(ctx *HaveVarContext) {}

// func (s *WclBuildListener) ExitHaveVar(ctx *HaveVarContext) {
// 	v := ctx.Value_ref().GetVr()
// 	elem := &tree.DocElement{ValueRef: v}
// 	ctx.SetElem(elem)
// }

func (s *WclBuildListener) EnterHaveList(ctx *HaveListContext) {}

func (s *WclBuildListener) ExitHaveList(ctx *HaveListContext) {
	if ctx.Doc_elem_child() == nil {
		log.Fatalf("no element child")
	}
	elem := ctx.Doc_elem_child().GetElem()
	elem.Tag = nil
	ctx.SetElem(elem)

}

func (s *WclBuildListener) EnterFunc_invoc(ctx *Func_invocContext) {
}

func (s *WclBuildListener) ExitFunc_invoc(ctx *Func_invocContext) {
	actual := ctx.Func_actual_seq().GetActual()
	var name string
	if ctx.Id() != nil {
		name = ctx.Id().GetText()
	}
	var line, col int
	if ctx.Id() != nil {
		line = ctx.Id().GetSymbol().GetLine()
		col = ctx.Id().GetSymbol().GetColumn()
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
	ctx.SetFn(dfunc)
}

func (s *WclBuildListener) EnterDoc_func_post(ctx *Doc_func_postContext) {}

func (s *WclBuildListener) ExitDoc_func_post(ctx *Doc_func_postContext) {
	var f, l []*tree.PFormal
	var pre, post []tree.TextItem

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
	ctx.SetFn(tree.NewDocFuncNode("", f, l, ctx.Doc_elem().GetElem(),
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
	// b := ctx.Dash() != nil && ctx.GreaterThan() != nil
	// f := ctx.Func_invoc().GetInvoc()
	// f.Builtin = b
	// if b {
	// 	s.checkSingleParamIsCSSClass(ctx, f)
	// }
	// ctx.SetInvoc(f)
}

func (s *WclBuildListener) checkSingleParamIsCSSClass(ctx *Event_callContext, f *tree.FuncInvoc) {
	//chk, err := builtin.GetBuiltinChecker(f.Name)
	// if err != nil {
	// 	notifyError(err.Error(),
	// 		ctx.BaseParserRuleContext, ctx.GetParser())
	// 	return
	// }
	detail := fmt.Sprintf("%s:%d:%d", s.SourceCode, f.LineNumber, f.ColumnNumber)
	if len(f.Actual) == 0 {
		notifyError(fmt.Sprintf("%s expected 1 parameter got 0 for %s", detail, f.Name),
			ctx.BaseParserRuleContext, ctx.GetParser())
		return

	}
	if len(f.Actual) > 1 {
		notifyError(fmt.Sprintf("%s number of parameters expected for '%s' is 1", detail, f.Name),
			ctx.BaseParserRuleContext, ctx.GetParser())
		return
	}
	// ok, errText := chk(f.Actual[0].Literal[1 : len(f.Actual[0].Literal)-1])
	// if !ok {
	// 	notifyError(errText,
	// 		ctx.BaseParserRuleContext, name, ctx.GetParser())
	// 	return
	// }
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

func (s *WclBuildListener) ExitView_dedcl(ctx *View_declContext) {
	vdecl := tree.NewViewDecl(ctx.vname.GetText())
	fn := ctx.Doc_func_post().GetFn()
	vdecl.DocFn = fn
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
		//log.Printf("Model '%s'", decl[i].Name)
	}
	section.ModelDecl = decl

	rawView := ctx.AllView_decl()
	vdecl := make([]*tree.ViewDecl, len(rawView))
	for i, v := range rawView {
		vdecl[i] = v.GetVdecl()
		//log.Printf("View '%s'", decl[i].Name)
	}
	section.ViewDecl = vdecl

	ctx.SetSection(section)
}
