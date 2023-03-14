// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // wcl
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

type BasewclVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BasewclVisitor) VisitProgram(ctx *ProgramContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitGlobal(ctx *GlobalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitExtern(ctx *ExternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitWcl_section(ctx *Wcl_sectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitImport_section(ctx *Import_sectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitCss_section(ctx *Css_sectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitCss_filespec(ctx *Css_filespecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitText_section(ctx *Text_sectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitText_func(ctx *Text_funcContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitPre_code(ctx *Pre_codeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitPost_code(ctx *Post_codeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitText_func_local(ctx *Text_func_localContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitText_top(ctx *Text_topContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitText_content(ctx *Text_contentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitRawText(ctx *RawTextContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitVarSub(ctx *VarSubContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitVar_subs(ctx *Var_subsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitSub(ctx *SubContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitUninterp(ctx *UninterpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitUninterpRawText(ctx *UninterpRawTextContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitUninterpNested(ctx *UninterpNestedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitUninterpVar(ctx *UninterpVarContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitUninterp_var(ctx *Uninterp_varContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitParam_spec(ctx *Param_specContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitPair(ctx *PairContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitLast(ctx *LastContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitDoc_section(ctx *Doc_sectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitDoc_func(ctx *Doc_funcContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitDoc_func_local(ctx *Doc_func_localContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitDoc_func_formal(ctx *Doc_func_formalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitDoc_tag(ctx *Doc_tagContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitId_or_var_ref(ctx *Id_or_var_refContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitVar_ref(ctx *Var_refContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitDoc_id(ctx *Doc_idContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitDoc_class(ctx *Doc_classContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitHaveVar(ctx *HaveVarContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitHaveTag(ctx *HaveTagContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitHaveList(ctx *HaveListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitDoc_elem_content(ctx *Doc_elem_contentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitDoc_elem_text_func_call(ctx *Doc_elem_text_func_callContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitDoc_elem_text_anon(ctx *Doc_elem_text_anonContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitDoc_elem_child(ctx *Doc_elem_childContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitFunc_invoc(ctx *Func_invocContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitFunc_actual_seq(ctx *Func_actual_seqContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitFunc_actual(ctx *Func_actualContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitEvent_section(ctx *Event_sectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitEvent_spec(ctx *Event_specContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitEvent_call(ctx *Event_callContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitSelector(ctx *SelectorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitModel_section(ctx *Model_sectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitModel_def(ctx *Model_defContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasewclVisitor) VisitFilename_seq(ctx *Filename_seqContext) interface{} {
	return v.VisitChildren(ctx)
}
