// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // wcl
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// A complete Visitor for a parse tree produced by wcl.
type wclVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by wcl#program.
	VisitProgram(ctx *ProgramContext) interface{}

	// Visit a parse tree produced by wcl#global.
	VisitGlobal(ctx *GlobalContext) interface{}

	// Visit a parse tree produced by wcl#extern.
	VisitExtern(ctx *ExternContext) interface{}

	// Visit a parse tree produced by wcl#wcl_section.
	VisitWcl_section(ctx *Wcl_sectionContext) interface{}

	// Visit a parse tree produced by wcl#import_section.
	VisitImport_section(ctx *Import_sectionContext) interface{}

	// Visit a parse tree produced by wcl#css_section.
	VisitCss_section(ctx *Css_sectionContext) interface{}

	// Visit a parse tree produced by wcl#css_filespec.
	VisitCss_filespec(ctx *Css_filespecContext) interface{}

	// Visit a parse tree produced by wcl#text_section.
	VisitText_section(ctx *Text_sectionContext) interface{}

	// Visit a parse tree produced by wcl#text_func.
	VisitText_func(ctx *Text_funcContext) interface{}

	// Visit a parse tree produced by wcl#pre_code.
	VisitPre_code(ctx *Pre_codeContext) interface{}

	// Visit a parse tree produced by wcl#post_code.
	VisitPost_code(ctx *Post_codeContext) interface{}

	// Visit a parse tree produced by wcl#text_func_local.
	VisitText_func_local(ctx *Text_func_localContext) interface{}

	// Visit a parse tree produced by wcl#text_top.
	VisitText_top(ctx *Text_topContext) interface{}

	// Visit a parse tree produced by wcl#text_content.
	VisitText_content(ctx *Text_contentContext) interface{}

	// Visit a parse tree produced by wcl#RawText.
	VisitRawText(ctx *RawTextContext) interface{}

	// Visit a parse tree produced by wcl#VarSub.
	VisitVarSub(ctx *VarSubContext) interface{}

	// Visit a parse tree produced by wcl#var_subs.
	VisitVar_subs(ctx *Var_subsContext) interface{}

	// Visit a parse tree produced by wcl#sub.
	VisitSub(ctx *SubContext) interface{}

	// Visit a parse tree produced by wcl#uninterp.
	VisitUninterp(ctx *UninterpContext) interface{}

	// Visit a parse tree produced by wcl#UninterpRawText.
	VisitUninterpRawText(ctx *UninterpRawTextContext) interface{}

	// Visit a parse tree produced by wcl#UninterpNested.
	VisitUninterpNested(ctx *UninterpNestedContext) interface{}

	// Visit a parse tree produced by wcl#UninterpVar.
	VisitUninterpVar(ctx *UninterpVarContext) interface{}

	// Visit a parse tree produced by wcl#uninterp_var.
	VisitUninterp_var(ctx *Uninterp_varContext) interface{}

	// Visit a parse tree produced by wcl#param_spec.
	VisitParam_spec(ctx *Param_specContext) interface{}

	// Visit a parse tree produced by wcl#Pair.
	VisitPair(ctx *PairContext) interface{}

	// Visit a parse tree produced by wcl#Last.
	VisitLast(ctx *LastContext) interface{}

	// Visit a parse tree produced by wcl#doc_section.
	VisitDoc_section(ctx *Doc_sectionContext) interface{}

	// Visit a parse tree produced by wcl#doc_func.
	VisitDoc_func(ctx *Doc_funcContext) interface{}

	// Visit a parse tree produced by wcl#doc_func_local.
	VisitDoc_func_local(ctx *Doc_func_localContext) interface{}

	// Visit a parse tree produced by wcl#doc_func_formal.
	VisitDoc_func_formal(ctx *Doc_func_formalContext) interface{}

	// Visit a parse tree produced by wcl#doc_tag.
	VisitDoc_tag(ctx *Doc_tagContext) interface{}

	// Visit a parse tree produced by wcl#id_or_var_ref.
	VisitId_or_var_ref(ctx *Id_or_var_refContext) interface{}

	// Visit a parse tree produced by wcl#var_ref.
	VisitVar_ref(ctx *Var_refContext) interface{}

	// Visit a parse tree produced by wcl#doc_id.
	VisitDoc_id(ctx *Doc_idContext) interface{}

	// Visit a parse tree produced by wcl#doc_class.
	VisitDoc_class(ctx *Doc_classContext) interface{}

	// Visit a parse tree produced by wcl#haveVar.
	VisitHaveVar(ctx *HaveVarContext) interface{}

	// Visit a parse tree produced by wcl#haveTag.
	VisitHaveTag(ctx *HaveTagContext) interface{}

	// Visit a parse tree produced by wcl#haveList.
	VisitHaveList(ctx *HaveListContext) interface{}

	// Visit a parse tree produced by wcl#doc_elem_content.
	VisitDoc_elem_content(ctx *Doc_elem_contentContext) interface{}

	// Visit a parse tree produced by wcl#doc_elem_text_func_call.
	VisitDoc_elem_text_func_call(ctx *Doc_elem_text_func_callContext) interface{}

	// Visit a parse tree produced by wcl#doc_elem_text_anon.
	VisitDoc_elem_text_anon(ctx *Doc_elem_text_anonContext) interface{}

	// Visit a parse tree produced by wcl#doc_elem_child.
	VisitDoc_elem_child(ctx *Doc_elem_childContext) interface{}

	// Visit a parse tree produced by wcl#func_invoc.
	VisitFunc_invoc(ctx *Func_invocContext) interface{}

	// Visit a parse tree produced by wcl#func_actual_seq.
	VisitFunc_actual_seq(ctx *Func_actual_seqContext) interface{}

	// Visit a parse tree produced by wcl#func_actual.
	VisitFunc_actual(ctx *Func_actualContext) interface{}

	// Visit a parse tree produced by wcl#event_section.
	VisitEvent_section(ctx *Event_sectionContext) interface{}

	// Visit a parse tree produced by wcl#event_spec.
	VisitEvent_spec(ctx *Event_specContext) interface{}

	// Visit a parse tree produced by wcl#event_call.
	VisitEvent_call(ctx *Event_callContext) interface{}

	// Visit a parse tree produced by wcl#selector.
	VisitSelector(ctx *SelectorContext) interface{}
}
