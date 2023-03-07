// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // wcl
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// BasewclListener is a complete listener for a parse tree produced by wcl.
type BasewclListener struct{}

var _ wclListener = &BasewclListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasewclListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasewclListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasewclListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasewclListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterProgram is called when production program is entered.
func (s *BasewclListener) EnterProgram(ctx *ProgramContext) {}

// ExitProgram is called when production program is exited.
func (s *BasewclListener) ExitProgram(ctx *ProgramContext) {}

// EnterGlobal is called when production global is entered.
func (s *BasewclListener) EnterGlobal(ctx *GlobalContext) {}

// ExitGlobal is called when production global is exited.
func (s *BasewclListener) ExitGlobal(ctx *GlobalContext) {}

// EnterExtern is called when production extern is entered.
func (s *BasewclListener) EnterExtern(ctx *ExternContext) {}

// ExitExtern is called when production extern is exited.
func (s *BasewclListener) ExitExtern(ctx *ExternContext) {}

// EnterWcl_section is called when production wcl_section is entered.
func (s *BasewclListener) EnterWcl_section(ctx *Wcl_sectionContext) {}

// ExitWcl_section is called when production wcl_section is exited.
func (s *BasewclListener) ExitWcl_section(ctx *Wcl_sectionContext) {}

// EnterImport_section is called when production import_section is entered.
func (s *BasewclListener) EnterImport_section(ctx *Import_sectionContext) {}

// ExitImport_section is called when production import_section is exited.
func (s *BasewclListener) ExitImport_section(ctx *Import_sectionContext) {}

// EnterCss_section is called when production css_section is entered.
func (s *BasewclListener) EnterCss_section(ctx *Css_sectionContext) {}

// ExitCss_section is called when production css_section is exited.
func (s *BasewclListener) ExitCss_section(ctx *Css_sectionContext) {}

// EnterCss_filespec is called when production css_filespec is entered.
func (s *BasewclListener) EnterCss_filespec(ctx *Css_filespecContext) {}

// ExitCss_filespec is called when production css_filespec is exited.
func (s *BasewclListener) ExitCss_filespec(ctx *Css_filespecContext) {}

// EnterText_section is called when production text_section is entered.
func (s *BasewclListener) EnterText_section(ctx *Text_sectionContext) {}

// ExitText_section is called when production text_section is exited.
func (s *BasewclListener) ExitText_section(ctx *Text_sectionContext) {}

// EnterText_func is called when production text_func is entered.
func (s *BasewclListener) EnterText_func(ctx *Text_funcContext) {}

// ExitText_func is called when production text_func is exited.
func (s *BasewclListener) ExitText_func(ctx *Text_funcContext) {}

// EnterPre_code is called when production pre_code is entered.
func (s *BasewclListener) EnterPre_code(ctx *Pre_codeContext) {}

// ExitPre_code is called when production pre_code is exited.
func (s *BasewclListener) ExitPre_code(ctx *Pre_codeContext) {}

// EnterPost_code is called when production post_code is entered.
func (s *BasewclListener) EnterPost_code(ctx *Post_codeContext) {}

// ExitPost_code is called when production post_code is exited.
func (s *BasewclListener) ExitPost_code(ctx *Post_codeContext) {}

// EnterText_func_local is called when production text_func_local is entered.
func (s *BasewclListener) EnterText_func_local(ctx *Text_func_localContext) {}

// ExitText_func_local is called when production text_func_local is exited.
func (s *BasewclListener) ExitText_func_local(ctx *Text_func_localContext) {}

// EnterText_top is called when production text_top is entered.
func (s *BasewclListener) EnterText_top(ctx *Text_topContext) {}

// ExitText_top is called when production text_top is exited.
func (s *BasewclListener) ExitText_top(ctx *Text_topContext) {}

// EnterText_content is called when production text_content is entered.
func (s *BasewclListener) EnterText_content(ctx *Text_contentContext) {}

// ExitText_content is called when production text_content is exited.
func (s *BasewclListener) ExitText_content(ctx *Text_contentContext) {}

// EnterRawText is called when production RawText is entered.
func (s *BasewclListener) EnterRawText(ctx *RawTextContext) {}

// ExitRawText is called when production RawText is exited.
func (s *BasewclListener) ExitRawText(ctx *RawTextContext) {}

// EnterVarSub is called when production VarSub is entered.
func (s *BasewclListener) EnterVarSub(ctx *VarSubContext) {}

// ExitVarSub is called when production VarSub is exited.
func (s *BasewclListener) ExitVarSub(ctx *VarSubContext) {}

// EnterVar_subs is called when production var_subs is entered.
func (s *BasewclListener) EnterVar_subs(ctx *Var_subsContext) {}

// ExitVar_subs is called when production var_subs is exited.
func (s *BasewclListener) ExitVar_subs(ctx *Var_subsContext) {}

// EnterSub is called when production sub is entered.
func (s *BasewclListener) EnterSub(ctx *SubContext) {}

// ExitSub is called when production sub is exited.
func (s *BasewclListener) ExitSub(ctx *SubContext) {}

// EnterUninterp is called when production uninterp is entered.
func (s *BasewclListener) EnterUninterp(ctx *UninterpContext) {}

// ExitUninterp is called when production uninterp is exited.
func (s *BasewclListener) ExitUninterp(ctx *UninterpContext) {}

// EnterUninterpRawText is called when production UninterpRawText is entered.
func (s *BasewclListener) EnterUninterpRawText(ctx *UninterpRawTextContext) {}

// ExitUninterpRawText is called when production UninterpRawText is exited.
func (s *BasewclListener) ExitUninterpRawText(ctx *UninterpRawTextContext) {}

// EnterUninterpNested is called when production UninterpNested is entered.
func (s *BasewclListener) EnterUninterpNested(ctx *UninterpNestedContext) {}

// ExitUninterpNested is called when production UninterpNested is exited.
func (s *BasewclListener) ExitUninterpNested(ctx *UninterpNestedContext) {}

// EnterUninterpVar is called when production UninterpVar is entered.
func (s *BasewclListener) EnterUninterpVar(ctx *UninterpVarContext) {}

// ExitUninterpVar is called when production UninterpVar is exited.
func (s *BasewclListener) ExitUninterpVar(ctx *UninterpVarContext) {}

// EnterUninterp_var is called when production uninterp_var is entered.
func (s *BasewclListener) EnterUninterp_var(ctx *Uninterp_varContext) {}

// ExitUninterp_var is called when production uninterp_var is exited.
func (s *BasewclListener) ExitUninterp_var(ctx *Uninterp_varContext) {}

// EnterParam_spec is called when production param_spec is entered.
func (s *BasewclListener) EnterParam_spec(ctx *Param_specContext) {}

// ExitParam_spec is called when production param_spec is exited.
func (s *BasewclListener) ExitParam_spec(ctx *Param_specContext) {}

// EnterPair is called when production Pair is entered.
func (s *BasewclListener) EnterPair(ctx *PairContext) {}

// ExitPair is called when production Pair is exited.
func (s *BasewclListener) ExitPair(ctx *PairContext) {}

// EnterLast is called when production Last is entered.
func (s *BasewclListener) EnterLast(ctx *LastContext) {}

// ExitLast is called when production Last is exited.
func (s *BasewclListener) ExitLast(ctx *LastContext) {}

// EnterDoc_section is called when production doc_section is entered.
func (s *BasewclListener) EnterDoc_section(ctx *Doc_sectionContext) {}

// ExitDoc_section is called when production doc_section is exited.
func (s *BasewclListener) ExitDoc_section(ctx *Doc_sectionContext) {}

// EnterDoc_func is called when production doc_func is entered.
func (s *BasewclListener) EnterDoc_func(ctx *Doc_funcContext) {}

// ExitDoc_func is called when production doc_func is exited.
func (s *BasewclListener) ExitDoc_func(ctx *Doc_funcContext) {}

// EnterDoc_func_local is called when production doc_func_local is entered.
func (s *BasewclListener) EnterDoc_func_local(ctx *Doc_func_localContext) {}

// ExitDoc_func_local is called when production doc_func_local is exited.
func (s *BasewclListener) ExitDoc_func_local(ctx *Doc_func_localContext) {}

// EnterDoc_func_formal is called when production doc_func_formal is entered.
func (s *BasewclListener) EnterDoc_func_formal(ctx *Doc_func_formalContext) {}

// ExitDoc_func_formal is called when production doc_func_formal is exited.
func (s *BasewclListener) ExitDoc_func_formal(ctx *Doc_func_formalContext) {}

// EnterDoc_tag is called when production doc_tag is entered.
func (s *BasewclListener) EnterDoc_tag(ctx *Doc_tagContext) {}

// ExitDoc_tag is called when production doc_tag is exited.
func (s *BasewclListener) ExitDoc_tag(ctx *Doc_tagContext) {}

// EnterId_or_var_ref is called when production id_or_var_ref is entered.
func (s *BasewclListener) EnterId_or_var_ref(ctx *Id_or_var_refContext) {}

// ExitId_or_var_ref is called when production id_or_var_ref is exited.
func (s *BasewclListener) ExitId_or_var_ref(ctx *Id_or_var_refContext) {}

// EnterVar_ref is called when production var_ref is entered.
func (s *BasewclListener) EnterVar_ref(ctx *Var_refContext) {}

// ExitVar_ref is called when production var_ref is exited.
func (s *BasewclListener) ExitVar_ref(ctx *Var_refContext) {}

// EnterDoc_id is called when production doc_id is entered.
func (s *BasewclListener) EnterDoc_id(ctx *Doc_idContext) {}

// ExitDoc_id is called when production doc_id is exited.
func (s *BasewclListener) ExitDoc_id(ctx *Doc_idContext) {}

// EnterDoc_class is called when production doc_class is entered.
func (s *BasewclListener) EnterDoc_class(ctx *Doc_classContext) {}

// ExitDoc_class is called when production doc_class is exited.
func (s *BasewclListener) ExitDoc_class(ctx *Doc_classContext) {}

// EnterHaveVar is called when production haveVar is entered.
func (s *BasewclListener) EnterHaveVar(ctx *HaveVarContext) {}

// ExitHaveVar is called when production haveVar is exited.
func (s *BasewclListener) ExitHaveVar(ctx *HaveVarContext) {}

// EnterHaveTag is called when production haveTag is entered.
func (s *BasewclListener) EnterHaveTag(ctx *HaveTagContext) {}

// ExitHaveTag is called when production haveTag is exited.
func (s *BasewclListener) ExitHaveTag(ctx *HaveTagContext) {}

// EnterHaveList is called when production haveList is entered.
func (s *BasewclListener) EnterHaveList(ctx *HaveListContext) {}

// ExitHaveList is called when production haveList is exited.
func (s *BasewclListener) ExitHaveList(ctx *HaveListContext) {}

// EnterDoc_elem_content is called when production doc_elem_content is entered.
func (s *BasewclListener) EnterDoc_elem_content(ctx *Doc_elem_contentContext) {}

// ExitDoc_elem_content is called when production doc_elem_content is exited.
func (s *BasewclListener) ExitDoc_elem_content(ctx *Doc_elem_contentContext) {}

// EnterDoc_elem_text_func_call is called when production doc_elem_text_func_call is entered.
func (s *BasewclListener) EnterDoc_elem_text_func_call(ctx *Doc_elem_text_func_callContext) {}

// ExitDoc_elem_text_func_call is called when production doc_elem_text_func_call is exited.
func (s *BasewclListener) ExitDoc_elem_text_func_call(ctx *Doc_elem_text_func_callContext) {}

// EnterDoc_elem_text_anon is called when production doc_elem_text_anon is entered.
func (s *BasewclListener) EnterDoc_elem_text_anon(ctx *Doc_elem_text_anonContext) {}

// ExitDoc_elem_text_anon is called when production doc_elem_text_anon is exited.
func (s *BasewclListener) ExitDoc_elem_text_anon(ctx *Doc_elem_text_anonContext) {}

// EnterDoc_elem_child is called when production doc_elem_child is entered.
func (s *BasewclListener) EnterDoc_elem_child(ctx *Doc_elem_childContext) {}

// ExitDoc_elem_child is called when production doc_elem_child is exited.
func (s *BasewclListener) ExitDoc_elem_child(ctx *Doc_elem_childContext) {}

// EnterFunc_invoc is called when production func_invoc is entered.
func (s *BasewclListener) EnterFunc_invoc(ctx *Func_invocContext) {}

// ExitFunc_invoc is called when production func_invoc is exited.
func (s *BasewclListener) ExitFunc_invoc(ctx *Func_invocContext) {}

// EnterFunc_actual_seq is called when production func_actual_seq is entered.
func (s *BasewclListener) EnterFunc_actual_seq(ctx *Func_actual_seqContext) {}

// ExitFunc_actual_seq is called when production func_actual_seq is exited.
func (s *BasewclListener) ExitFunc_actual_seq(ctx *Func_actual_seqContext) {}

// EnterFunc_actual is called when production func_actual is entered.
func (s *BasewclListener) EnterFunc_actual(ctx *Func_actualContext) {}

// ExitFunc_actual is called when production func_actual is exited.
func (s *BasewclListener) ExitFunc_actual(ctx *Func_actualContext) {}

// EnterEvent_section is called when production event_section is entered.
func (s *BasewclListener) EnterEvent_section(ctx *Event_sectionContext) {}

// ExitEvent_section is called when production event_section is exited.
func (s *BasewclListener) ExitEvent_section(ctx *Event_sectionContext) {}

// EnterEvent_spec is called when production event_spec is entered.
func (s *BasewclListener) EnterEvent_spec(ctx *Event_specContext) {}

// ExitEvent_spec is called when production event_spec is exited.
func (s *BasewclListener) ExitEvent_spec(ctx *Event_specContext) {}

// EnterEvent_call is called when production event_call is entered.
func (s *BasewclListener) EnterEvent_call(ctx *Event_callContext) {}

// ExitEvent_call is called when production event_call is exited.
func (s *BasewclListener) ExitEvent_call(ctx *Event_callContext) {}

// EnterSelector is called when production selector is entered.
func (s *BasewclListener) EnterSelector(ctx *SelectorContext) {}

// ExitSelector is called when production selector is exited.
func (s *BasewclListener) ExitSelector(ctx *SelectorContext) {}
