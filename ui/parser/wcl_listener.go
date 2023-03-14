// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // wcl
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// wclListener is a complete listener for a parse tree produced by wcl.
type wclListener interface {
	antlr.ParseTreeListener

	// EnterProgram is called when entering the program production.
	EnterProgram(c *ProgramContext)

	// EnterGlobal is called when entering the global production.
	EnterGlobal(c *GlobalContext)

	// EnterExtern is called when entering the extern production.
	EnterExtern(c *ExternContext)

	// EnterWcl_section is called when entering the wcl_section production.
	EnterWcl_section(c *Wcl_sectionContext)

	// EnterImport_section is called when entering the import_section production.
	EnterImport_section(c *Import_sectionContext)

	// EnterCss_section is called when entering the css_section production.
	EnterCss_section(c *Css_sectionContext)

	// EnterCss_filespec is called when entering the css_filespec production.
	EnterCss_filespec(c *Css_filespecContext)

	// EnterText_section is called when entering the text_section production.
	EnterText_section(c *Text_sectionContext)

	// EnterText_func is called when entering the text_func production.
	EnterText_func(c *Text_funcContext)

	// EnterPre_code is called when entering the pre_code production.
	EnterPre_code(c *Pre_codeContext)

	// EnterPost_code is called when entering the post_code production.
	EnterPost_code(c *Post_codeContext)

	// EnterText_func_local is called when entering the text_func_local production.
	EnterText_func_local(c *Text_func_localContext)

	// EnterText_top is called when entering the text_top production.
	EnterText_top(c *Text_topContext)

	// EnterText_content is called when entering the text_content production.
	EnterText_content(c *Text_contentContext)

	// EnterRawText is called when entering the RawText production.
	EnterRawText(c *RawTextContext)

	// EnterVarSub is called when entering the VarSub production.
	EnterVarSub(c *VarSubContext)

	// EnterVar_subs is called when entering the var_subs production.
	EnterVar_subs(c *Var_subsContext)

	// EnterSub is called when entering the sub production.
	EnterSub(c *SubContext)

	// EnterUninterp is called when entering the uninterp production.
	EnterUninterp(c *UninterpContext)

	// EnterUninterpRawText is called when entering the UninterpRawText production.
	EnterUninterpRawText(c *UninterpRawTextContext)

	// EnterUninterpNested is called when entering the UninterpNested production.
	EnterUninterpNested(c *UninterpNestedContext)

	// EnterUninterpVar is called when entering the UninterpVar production.
	EnterUninterpVar(c *UninterpVarContext)

	// EnterUninterp_var is called when entering the uninterp_var production.
	EnterUninterp_var(c *Uninterp_varContext)

	// EnterParam_spec is called when entering the param_spec production.
	EnterParam_spec(c *Param_specContext)

	// EnterPair is called when entering the Pair production.
	EnterPair(c *PairContext)

	// EnterLast is called when entering the Last production.
	EnterLast(c *LastContext)

	// EnterDoc_section is called when entering the doc_section production.
	EnterDoc_section(c *Doc_sectionContext)

	// EnterDoc_func is called when entering the doc_func production.
	EnterDoc_func(c *Doc_funcContext)

	// EnterDoc_func_local is called when entering the doc_func_local production.
	EnterDoc_func_local(c *Doc_func_localContext)

	// EnterDoc_func_formal is called when entering the doc_func_formal production.
	EnterDoc_func_formal(c *Doc_func_formalContext)

	// EnterDoc_tag is called when entering the doc_tag production.
	EnterDoc_tag(c *Doc_tagContext)

	// EnterId_or_var_ref is called when entering the id_or_var_ref production.
	EnterId_or_var_ref(c *Id_or_var_refContext)

	// EnterVar_ref is called when entering the var_ref production.
	EnterVar_ref(c *Var_refContext)

	// EnterDoc_id is called when entering the doc_id production.
	EnterDoc_id(c *Doc_idContext)

	// EnterDoc_class is called when entering the doc_class production.
	EnterDoc_class(c *Doc_classContext)

	// EnterHaveVar is called when entering the haveVar production.
	EnterHaveVar(c *HaveVarContext)

	// EnterHaveTag is called when entering the haveTag production.
	EnterHaveTag(c *HaveTagContext)

	// EnterHaveList is called when entering the haveList production.
	EnterHaveList(c *HaveListContext)

	// EnterDoc_elem_content is called when entering the doc_elem_content production.
	EnterDoc_elem_content(c *Doc_elem_contentContext)

	// EnterDoc_elem_text_func_call is called when entering the doc_elem_text_func_call production.
	EnterDoc_elem_text_func_call(c *Doc_elem_text_func_callContext)

	// EnterDoc_elem_text_anon is called when entering the doc_elem_text_anon production.
	EnterDoc_elem_text_anon(c *Doc_elem_text_anonContext)

	// EnterDoc_elem_child is called when entering the doc_elem_child production.
	EnterDoc_elem_child(c *Doc_elem_childContext)

	// EnterFunc_invoc is called when entering the func_invoc production.
	EnterFunc_invoc(c *Func_invocContext)

	// EnterFunc_actual_seq is called when entering the func_actual_seq production.
	EnterFunc_actual_seq(c *Func_actual_seqContext)

	// EnterFunc_actual is called when entering the func_actual production.
	EnterFunc_actual(c *Func_actualContext)

	// EnterEvent_section is called when entering the event_section production.
	EnterEvent_section(c *Event_sectionContext)

	// EnterEvent_spec is called when entering the event_spec production.
	EnterEvent_spec(c *Event_specContext)

	// EnterEvent_call is called when entering the event_call production.
	EnterEvent_call(c *Event_callContext)

	// EnterSelector is called when entering the selector production.
	EnterSelector(c *SelectorContext)

	// EnterModel_section is called when entering the model_section production.
	EnterModel_section(c *Model_sectionContext)

	// EnterModel_def is called when entering the model_def production.
	EnterModel_def(c *Model_defContext)

	// EnterFilename_seq is called when entering the filename_seq production.
	EnterFilename_seq(c *Filename_seqContext)

	// ExitProgram is called when exiting the program production.
	ExitProgram(c *ProgramContext)

	// ExitGlobal is called when exiting the global production.
	ExitGlobal(c *GlobalContext)

	// ExitExtern is called when exiting the extern production.
	ExitExtern(c *ExternContext)

	// ExitWcl_section is called when exiting the wcl_section production.
	ExitWcl_section(c *Wcl_sectionContext)

	// ExitImport_section is called when exiting the import_section production.
	ExitImport_section(c *Import_sectionContext)

	// ExitCss_section is called when exiting the css_section production.
	ExitCss_section(c *Css_sectionContext)

	// ExitCss_filespec is called when exiting the css_filespec production.
	ExitCss_filespec(c *Css_filespecContext)

	// ExitText_section is called when exiting the text_section production.
	ExitText_section(c *Text_sectionContext)

	// ExitText_func is called when exiting the text_func production.
	ExitText_func(c *Text_funcContext)

	// ExitPre_code is called when exiting the pre_code production.
	ExitPre_code(c *Pre_codeContext)

	// ExitPost_code is called when exiting the post_code production.
	ExitPost_code(c *Post_codeContext)

	// ExitText_func_local is called when exiting the text_func_local production.
	ExitText_func_local(c *Text_func_localContext)

	// ExitText_top is called when exiting the text_top production.
	ExitText_top(c *Text_topContext)

	// ExitText_content is called when exiting the text_content production.
	ExitText_content(c *Text_contentContext)

	// ExitRawText is called when exiting the RawText production.
	ExitRawText(c *RawTextContext)

	// ExitVarSub is called when exiting the VarSub production.
	ExitVarSub(c *VarSubContext)

	// ExitVar_subs is called when exiting the var_subs production.
	ExitVar_subs(c *Var_subsContext)

	// ExitSub is called when exiting the sub production.
	ExitSub(c *SubContext)

	// ExitUninterp is called when exiting the uninterp production.
	ExitUninterp(c *UninterpContext)

	// ExitUninterpRawText is called when exiting the UninterpRawText production.
	ExitUninterpRawText(c *UninterpRawTextContext)

	// ExitUninterpNested is called when exiting the UninterpNested production.
	ExitUninterpNested(c *UninterpNestedContext)

	// ExitUninterpVar is called when exiting the UninterpVar production.
	ExitUninterpVar(c *UninterpVarContext)

	// ExitUninterp_var is called when exiting the uninterp_var production.
	ExitUninterp_var(c *Uninterp_varContext)

	// ExitParam_spec is called when exiting the param_spec production.
	ExitParam_spec(c *Param_specContext)

	// ExitPair is called when exiting the Pair production.
	ExitPair(c *PairContext)

	// ExitLast is called when exiting the Last production.
	ExitLast(c *LastContext)

	// ExitDoc_section is called when exiting the doc_section production.
	ExitDoc_section(c *Doc_sectionContext)

	// ExitDoc_func is called when exiting the doc_func production.
	ExitDoc_func(c *Doc_funcContext)

	// ExitDoc_func_local is called when exiting the doc_func_local production.
	ExitDoc_func_local(c *Doc_func_localContext)

	// ExitDoc_func_formal is called when exiting the doc_func_formal production.
	ExitDoc_func_formal(c *Doc_func_formalContext)

	// ExitDoc_tag is called when exiting the doc_tag production.
	ExitDoc_tag(c *Doc_tagContext)

	// ExitId_or_var_ref is called when exiting the id_or_var_ref production.
	ExitId_or_var_ref(c *Id_or_var_refContext)

	// ExitVar_ref is called when exiting the var_ref production.
	ExitVar_ref(c *Var_refContext)

	// ExitDoc_id is called when exiting the doc_id production.
	ExitDoc_id(c *Doc_idContext)

	// ExitDoc_class is called when exiting the doc_class production.
	ExitDoc_class(c *Doc_classContext)

	// ExitHaveVar is called when exiting the haveVar production.
	ExitHaveVar(c *HaveVarContext)

	// ExitHaveTag is called when exiting the haveTag production.
	ExitHaveTag(c *HaveTagContext)

	// ExitHaveList is called when exiting the haveList production.
	ExitHaveList(c *HaveListContext)

	// ExitDoc_elem_content is called when exiting the doc_elem_content production.
	ExitDoc_elem_content(c *Doc_elem_contentContext)

	// ExitDoc_elem_text_func_call is called when exiting the doc_elem_text_func_call production.
	ExitDoc_elem_text_func_call(c *Doc_elem_text_func_callContext)

	// ExitDoc_elem_text_anon is called when exiting the doc_elem_text_anon production.
	ExitDoc_elem_text_anon(c *Doc_elem_text_anonContext)

	// ExitDoc_elem_child is called when exiting the doc_elem_child production.
	ExitDoc_elem_child(c *Doc_elem_childContext)

	// ExitFunc_invoc is called when exiting the func_invoc production.
	ExitFunc_invoc(c *Func_invocContext)

	// ExitFunc_actual_seq is called when exiting the func_actual_seq production.
	ExitFunc_actual_seq(c *Func_actual_seqContext)

	// ExitFunc_actual is called when exiting the func_actual production.
	ExitFunc_actual(c *Func_actualContext)

	// ExitEvent_section is called when exiting the event_section production.
	ExitEvent_section(c *Event_sectionContext)

	// ExitEvent_spec is called when exiting the event_spec production.
	ExitEvent_spec(c *Event_specContext)

	// ExitEvent_call is called when exiting the event_call production.
	ExitEvent_call(c *Event_callContext)

	// ExitSelector is called when exiting the selector production.
	ExitSelector(c *SelectorContext)

	// ExitModel_section is called when exiting the model_section production.
	ExitModel_section(c *Model_sectionContext)

	// ExitModel_def is called when exiting the model_def production.
	ExitModel_def(c *Model_defContext)

	// ExitFilename_seq is called when exiting the filename_seq production.
	ExitFilename_seq(c *Filename_seqContext)
}
