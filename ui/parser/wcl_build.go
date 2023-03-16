package parser

import (
	"fmt"
	"log"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/iansmith/parigot/ui/css"
	"github.com/iansmith/parigot/ui/parser/builtin"
	"github.com/iansmith/parigot/ui/parser/tree"
)

const anonPrefix = "_anon"

type WclBuildListener struct {
	*BasewclListener

	// we use these when the object does not need a stack (it's a singleton)
	TextSection *tree.TextSectionNode

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
	c.SetP(tree.NewProgramNode())
}

func (l *WclBuildListener) ExitProgram(c *ProgramContext) {
	if c.Import_section() != nil && c.Import_section().GetSection() != nil {
		c.GetP().ImportSection = c.Import_section().GetSection()
	}
	var tsNode *tree.TextSectionNode
	ts := c.Text_section()
	if ts == nil {
		if l.TextSection != nil {
			tsNode = l.TextSection
		}
	}
	if ts != nil && c.Text_section().GetSection() != nil {
		tsNode = c.Text_section().GetSection()
		tsNode.Program = c.GetP()
	}
	c.GetP().TextSection = tsNode
	if c.Doc_section() != nil && c.Doc_section().GetSection() != nil {
		c.GetP().DocSection = c.Doc_section().GetSection()
		c.GetP().DocSection.Program = c.GetP()
	}
	if c.Event_section() != nil && c.Event_section().GetSection() != nil {
		c.GetP().EventSection = c.Event_section().GetSection()
		c.GetP().EventSection.Program = c.GetP()
	}
	if c.Mvc_section() != nil && c.Mvc_section().GetSection() != nil {
		c.GetP().ModelSection = c.Mvc_section().GetSection()
		c.GetP().ModelSection.Program = c.GetP()
	}

	if c.Extern() != nil && c.Extern().GetE() != nil {
		c.GetP().Extern = c.Extern().GetE()
	}
	if c.Global() != nil && c.Global().GetG() != nil {
		c.GetP().Global = c.Global().GetG()
	}

}

// Import_section
func (l *WclBuildListener) EnterImport_section(c *Import_sectionContext) {
	i := tree.NewImportSectionNode()
	c.SetSection(i)
}
func (l *WclBuildListener) ExitImport_section(c *Import_sectionContext) {
	c.GetSection().TextItem_ = c.Uninterp().GetItem()
}

// Text_section
func (l *WclBuildListener) EnterText_section(c *Text_sectionContext) {
	c.SetSection(tree.NewTextSectionNode())
	l.TextSection = c.GetSection() // singletone
}

// Text_func
func (l *WclBuildListener) EnterText_func(c *Text_funcContext) {
	c.SetF(tree.NewTextFuncNode())
}
func (l *WclBuildListener) ExitText_func(c *Text_funcContext) {
	c.GetF().Name = c.Id().GetText()
	if c.GetF() != nil && c.GetF().Item_ != nil {
		c.GetF().Item_ = c.Text_top().GetItem()
	}
	if c.Param_spec() != nil {
		c.GetF().Param = c.Param_spec().GetFormal()
	}
	if c.Text_func_local() != nil {
		c.GetF().Local = c.Text_func_local().GetFormal()
	}
	c.GetF().Item_ = c.Text_top().GetItem()

	// add to list of funcs
	l.TextSection.Func = append(l.TextSection.Func, c.GetF())
	if c.Pre_code() != nil && c.Pre_code().GetItem() != nil {
		c.GetF().PreCode = c.Pre_code().GetItem()
	}
	if c.Post_code() != nil && c.Post_code().GetItem() != nil {
		c.GetF().PostCode = c.Post_code().GetItem()
	}
}

// Text_top
func (l *WclBuildListener) EnterText_top(c *Text_topContext) {
	//nothing to do
}

func (l *WclBuildListener) ExitText_top(c *Text_topContext) {
	if c.Text_content() != nil {
		c.SetItem(c.Text_content().GetItem())
	}
}

// Text_content
func (l *WclBuildListener) EnterText_content(c *Text_contentContext) {
}

func (l *WclBuildListener) ExitText_content(c *Text_contentContext) {
	result := []tree.TextItem{}
	for _, t := range c.AllText_content_inner() {
		result = append(result, t.GetItem()...)
	}
	c.SetItem(result)
}

// Text_content_inner.RawText
func (l *WclBuildListener) EnterRawText(c *RawTextContext) {
	//nothing to do
}
func (l *WclBuildListener) ExitRawText(c *RawTextContext) {
	c.SetItem([]tree.TextItem{tree.NewTextConstant(c.ContentRawText().GetText())})
}

// Text_content_inner.VarSub
func (l *WclBuildListener) EnterVarSub(c *VarSubContext) {
	//nothing to do
}
func (l *WclBuildListener) ExitVarSub(c *VarSubContext) {
	c.SetItem(c.Var_subs().GetItem())
}

// Var_subs
func (l *WclBuildListener) EnterVar_subs(c *Var_subsContext) {
}
func (l *WclBuildListener) ExitVar_subs(c *Var_subsContext) {
	c.SetItem([]tree.TextItem{c.Sub().GetItem()})
}

// sub
func (l *WclBuildListener) EnterSub(c *SubContext) {
}
func (l *WclBuildListener) ExitSub(c *SubContext) {
	if c.Func_invoc_var() != nil {
		log.Printf("sub found invoc %+v", c.Func_invoc_var().GetInvoc())
		c.SetItem(tree.NewTextInvoc(c.Func_invoc_var().GetInvoc()))
		return
	}
	//normal case, just a var
	c.SetItem(tree.NewTextVar(c.VarId().GetText()))
}

// UninterpRawText
func (l *WclBuildListener) EnterUninterpRawText(c *UninterpRawTextContext) {
}
func (l *WclBuildListener) ExitUninterpRawText(c *UninterpRawTextContext) {
	c.SetItem([]tree.TextItem{tree.NewTextConstant(c.UninterpRawText().GetText())})
}

// UninterpNested
func (l *WclBuildListener) EnterUninterpNested(c *UninterpNestedContext) {
}
func (l *WclBuildListener) ExitUninterpNested(c *UninterpNestedContext) {
	r := append([]tree.TextItem{tree.NewTextConstant("{")}, c.Uninterp().GetItem()...)
	r = append(r, tree.NewTextConstant("}"))
	c.SetItem(r)

}

// Var text
func (l *WclBuildListener) EnterUninterp_var(c *Uninterp_varContext) {
}
func (l *WclBuildListener) ExitUninterp_var(c *Uninterp_varContext) {
	v := tree.NewTextVar(c.VarId().GetText())
	c.SetItem([]tree.TextItem{v})
}

// Var inside uninterp
func (l *WclBuildListener) EnterUninterpVar(c *UninterpVarContext) {
}
func (l *WclBuildListener) ExitUninterpVar(c *UninterpVarContext) {
	c.SetItem(c.Uninterp_var().GetItem())
}

// Uninterp
func (l *WclBuildListener) EnterUninterp(c *UninterpContext) {
}
func (l *WclBuildListener) ExitUninterp(c *UninterpContext) {
	result := []tree.TextItem{}
	for _, t := range c.AllUninterp_inner() {
		result = append(result, t.GetItem()...)
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
				c.BaseParserRuleContext, raw, c.GetParser())
			return
		}

		c.SetT(tree.NewTypeDeclModel(part[0], part[1]))
		return
	}
}

func (l *WclBuildListener) ExitParam_pair(c *Param_pairContext) {
	n := c.Id().GetText()
	t := c.Simple_or_model_param().GetT()
	c.SetFormal(tree.NewPFormal(n, t))
}

// Doc_tag is a full tag descriptor
func (s *WclBuildListener) EnterDoc_id(ctx *Doc_idContext) {
}

// Doc_tag is a full tag descriptor
func (s *WclBuildListener) ExitDoc_id(ctx *Doc_idContext) {
	ctx.SetS(ctx.Id().GetText())
}

// Doc_tag is a full tag descriptor
func (s *WclBuildListener) EnterDoc_tag(ctx *Doc_tagContext) {
}

// Doc_tag is a full tag descriptor
func (s *WclBuildListener) ExitDoc_tag(ctx *Doc_tagContext) {
	docId := ""
	if ctx.Doc_id() != nil {
		docId = ctx.Doc_id().GetS()
	}
	cl := []string{}
	if ctx.Doc_class() != nil {
		cl = ctx.Doc_class().GetClazz()
		for _, c := range cl {
			tn := ctx.Id_or_var_ref().GetIdVar().Name
			if !strings.HasPrefix(c, ".") {
				notifyError(fmt.Sprintf("in tag '%s', class name '%s' is invalid, class names must start with a dot ", tn, c),
					ctx.BaseParserRuleContext, ctx.Id_or_var_ref().GetText(), ctx.GetParser())
				return

			}
			if _, ok := s.ClassName[c]; !ok {
				notifyError(fmt.Sprintf("in tag '%s', class name '%s' is not defined the css files declared", tn, c),
					ctx.BaseParserRuleContext, ctx.Id_or_var_ref().GetText(), ctx.GetParser())
				return
			}
		}
	}
	tag, err := tree.NewDocTag(ctx.Id_or_var_ref().GetIdVar(), docId, cl)
	if err != nil {
		notifyError(fmt.Sprintf("unknown tag '%s'",
			ctx.Id_or_var_ref().GetText()),
			ctx.BaseParserRuleContext,
			ctx.Id_or_var_ref().GetText(), ctx.GetParser())
		return
	}
	ctx.SetTag(tag)

}
func notifyError(msg string, ctx *antlr.BaseParserRuleContext, name string, parser antlr.Parser) {
	ex := antlr.NewBaseRecognitionException(msg, parser, parser.GetInputStream(), ctx)
	ctx.SetException(ex)
	parser.NotifyErrorListeners(msg, ctx.GetStart(), ex)
}

// Doc_class the part of a doc atom that looks like :foo, describing a css class
func (s *WclBuildListener) EnterDoc_class(ctx *Doc_classContext) {
}

// Doc_class the part of a doc atom that looks like :foo, describing a css class
func (s *WclBuildListener) ExitDoc_class(ctx *Doc_classContext) {
	id := ctx.AllId()
	result := make([]string, len(id))
	for i, t := range id {
		result[i] = t.GetText()
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
	section := tree.NewDocSectionNode(content)
	ctx.SetSection(section)
}

func (s *WclBuildListener) EnterDoc_elem_content(ctx *Doc_elem_contentContext) {}

func (s *WclBuildListener) ExitDoc_elem_content(ctx *Doc_elem_contentContext) {
	if ctx.Doc_elem_text() != nil {
		ctx.SetElement(&tree.DocElement{TextContent: ctx.Doc_elem_text().GetInvoc()})
	} else {
		ctx.SetElement(ctx.Doc_elem_child().GetElem())
	}
}
func (s *WclBuildListener) CreateTextSectionIfNeeded() {
	if s.TextSection == nil {
		s.TextSection = tree.NewTextSectionNode()
	}
}

func (s *WclBuildListener) EnterDoc_elem_text_func_call(ctx *Doc_elem_text_func_callContext) {}

func (s *WclBuildListener) ExitDoc_elem_text_func_call(ctx *Doc_elem_text_func_callContext) {
	ctx.SetInvoc(ctx.Func_invoc().GetInvoc())
}
func (s *WclBuildListener) EnterDoc_elem_text_anon(ctx *Doc_elem_text_anonContext) {}

func (s *WclBuildListener) ExitDoc_elem_text_anon(ctx *Doc_elem_text_anonContext) {
	item := ctx.Text_top().GetItem()
	name := fmt.Sprintf(anonPrefix+"%08d", s.anonCount)
	fn := &tree.TextFuncNode{Name: name, Item_: item}
	s.CreateTextSectionIfNeeded()
	s.TextSection.Func = append(s.TextSection.Func, fn)
	s.anonCount++

	ctx.SetInvoc(tree.NewFuncInvoc(&tree.DocIdOrVar{Name: name, IsVar: false}, nil))
}

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
	elem := &tree.DocElement{Tag: ctx.Doc_tag().GetTag()}
	elemContent := ctx.Doc_elem_content()
	if elemContent != nil {
		other := elemContent.GetElement()
		// might be either, but won't be both
		elem.Child = other.Child
		elem.TextContent = other.TextContent
	}
	ctx.SetElem(elem)
}

func (s *WclBuildListener) EnterHaveVar(ctx *HaveVarContext) {}

func (s *WclBuildListener) ExitHaveVar(ctx *HaveVarContext) {
	v := ctx.Var_ref().GetV()
	if !v.IsVar {
		panic("expected element to be a variable reference: " + v.Name)
	}
	elem := &tree.DocElement{Var: v.Name}
	ctx.SetElem(elem)
}

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
	name := ctx.Id().GetText()
	invoc := tree.NewFuncInvoc(&tree.DocIdOrVar{Name: name, IsVar: false}, actual)
	ctx.SetInvoc(invoc)
}

// This is a copy of func_invoc because we cannot use the same parameters (ctx is per production)
func (s *WclBuildListener) ExitFunc_invoc_var(ctx *Func_invoc_varContext) {
	actual := ctx.Func_actual_seq_var().GetActual()
	name := ctx.VarId().GetText()
	invoc := tree.NewFuncInvoc(&tree.DocIdOrVar{Name: name, IsVar: false}, actual)
	ctx.SetInvoc(invoc)
}

func (s *WclBuildListener) EnterFunc_actual(ctx *Func_actualContext) {
}

func (s *WclBuildListener) ExitFunc_actual(ctx *Func_actualContext) {
	id := ""
	lit := ""
	if ctx.Id() != nil {
		id = ctx.Id().GetText()
	}
	if ctx.StringLit() != nil {
		lit = ctx.StringLit().GetText()
	}
	if id != "" || lit != "" {
		ctx.SetActual(tree.NewFuncActual(id, lit))
	}
}

// This is a copy of func_actual_var because we cannot use the same parameters (ctx is per production)
func (s *WclBuildListener) ExitFunc_actual_var(ctx *Func_actual_varContext) {
	id := ""
	lit := ""
	if ctx.VarId() != nil {
		id = ctx.VarId().GetText()
	}
	if ctx.VarStringLit() != nil {
		lit = ctx.VarStringLit().GetText()
	}
	if id != "" || lit != "" {
		ctx.SetActual(tree.NewFuncActual(id, lit))
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

// This is a copy of func_actual_seq_var because we cannot use the same parameters (ctx is per production)
func (s *WclBuildListener) ExitFunc_actual_seq_var(ctx *Func_actual_seq_varContext) {
	raw := ctx.AllFunc_actual_var()
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
	ctx.SetG(ctx.Param_spec().GetFormal())
}

func (s *WclBuildListener) EnterExtern(ctx *ExternContext) {
}

func (s *WclBuildListener) ExitExtern(ctx *ExternContext) {
	raw := ctx.AllId()
	result := make([]string, len(raw))
	for i, name := range raw {
		result[i] = name.GetText()
	}
	ctx.SetE(result)
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

func (s *WclBuildListener) EnterId_or_var_ref(ctx *Id_or_var_refContext) {
}

func (s *WclBuildListener) ExitId_or_var_ref(ctx *Id_or_var_refContext) {
	if ctx.Id() != nil {
		ctx.SetIdVar(&tree.DocIdOrVar{Name: ctx.Id().GetText(), IsVar: false})
		return
	}
	ctx.SetIdVar(ctx.Var_ref().GetV())
}

func (s *WclBuildListener) EnterVar_ref(ctx *Var_refContext) {
}

func (s *WclBuildListener) ExitVar_ref(ctx *Var_refContext) {
	ctx.SetV(&tree.DocIdOrVar{Name: ctx.VarId().GetText(), IsVar: true})
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
	id := ""
	clazz := ""
	if ctx.GetIdValue() != nil {
		id = ctx.GetIdValue().GetText()
	}
	if ctx.GetClass() != nil {
		clazz = ctx.GetClass().GetText()
	}
	if clazz != "" && !strings.HasPrefix(clazz, ".") {
		notifyError(fmt.Sprintf("class name '%s' does not start with a dot", clazz),
			ctx.BaseParserRuleContext, ctx.Id().GetText(), ctx.GetParser())
		return
	}
	if id != "" {
		ctx.SetSel(&tree.Selector{Id: id})
	} else {
		ctx.SetSel(&tree.Selector{Class: clazz})
	}

}

func (s *WclBuildListener) EnterEvent_call(ctx *Event_callContext) {
}

func (s *WclBuildListener) ExitEvent_call(ctx *Event_callContext) {
	b := ctx.AllGreaterThan() != nil
	f := ctx.Func_invoc().GetInvoc()
	name := f.Name.Name
	f.Builtin = b
	if b {
		checkSingleParamIsCSSClass(ctx, name, f)
	}
	ctx.SetInvoc(f)
}

func checkSingleParamIsCSSClass(ctx *Event_callContext, name string, f *tree.FuncInvoc) {
	chk, err := builtin.GetBuiltinChecker(name)
	if err != nil {
		notifyError(err.Error(),
			ctx.BaseParserRuleContext, name, ctx.GetParser())
		return
	}
	if len(f.Actual) != 1 {
		notifyError(fmt.Sprintf("number of parameters expected for '%s' is 1", name),
			ctx.BaseParserRuleContext, name, ctx.GetParser())
		return
	}
	ok, errText := chk(f.Actual[0].Literal[1 : len(f.Actual[0].Literal)-1])
	if !ok {
		notifyError(errText,
			ctx.BaseParserRuleContext, name, ctx.GetParser())
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

func (s *WclBuildListener) ExitView_dedcl(ctx *View_declContext) {
	vdecl := tree.NewViewDecl(ctx.vname.GetText())
	fn := ctx.Doc_func_post().GetFn()
	vdecl.DocFn = fn
	ctx.SetVdecl(vdecl)
}

func (s *WclBuildListener) EnterMvc_section(ctx *Mvc_sectionContext) {
}

func (s *WclBuildListener) ExitMvc_section(ctx *Mvc_sectionContext) {
	section := tree.NewMvcSection()
	raw := ctx.AllModel_decl()
	decl := make([]*tree.ModelDecl, len(raw))
	for i, mod := range raw {
		decl[i] = mod.GetDecl()
		log.Printf("Model '%s'", decl[i].Name)
	}
	section.ModelDecl = decl

	rawView := ctx.AllView_decl()
	vdecl := make([]*tree.ViewDecl, len(rawView))
	for i, v := range rawView {
		vdecl[i] = v.GetVdecl()
		log.Printf("View '%s'", decl[i].Name)
	}
	section.ViewDecl = vdecl

	ctx.SetSection(section)
}
