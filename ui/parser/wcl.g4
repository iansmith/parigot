parser grammar wcl;
options {
	tokenVocab = wcllex;
}
@parser::header {
	import "github.com/iansmith/parigot/ui/parser/tree"
	import "log"
	var _ = &tree.ProgramNode{}
	var _ = log.Printf
}	

program
	returns[*tree.ProgramNode p]:
	wcl_section
	css_section?

	// the "early sections"	
	(
		import_section
		|extern
		|global
	)*

	mvc_section?
	text_section? 
	css_section?     
	doc_section?
	event_section?
	EOF   
	;

global
	returns [*tree.GlobalSectionNode g]:
	Global param_spec
	;

extern
	returns [*tree.ExternSectionNode e]:
	Extern LParen Id* RParen
	;

wcl_section:
	Wcl Version
	;

import_section
	returns[*tree.ImportSectionNode section]:
	Import uninterp 
	;

css_section:
	CSS (css_filespec)*
	;

css_filespec:
	Plus StringLit 
	;
text_section
	returns[*tree.TextSectionNode section]:
	Text (
		text_func 
	)*;

text_func
	returns[*tree.TextFuncNode f]:
	i = (Id|GrabId) param_spec? text_func_local? pre_code? uninterp post_code?
	;

pre_code 
	returns [[]tree.TextItem item]:
	Pre LCurly uninterp
	;

post_code 
	returns [[]tree.TextItem item]:
	Post LCurly uninterp
	;

text_func_local
	returns [[]*tree.PFormal formal]:
	Local param_spec
	;

//text_top
//	returns[[]tree.TextItem item]:
//	DoubleLess (
//		text_content
//		|
//	) GrabDoubleGreater;

text_content
	returns[[]tree.TextItem item]:
	(
		raw_text_or_sub
	)*;

raw_text_or_sub
	returns [[]tree.TextItem item]:
	RawText 
	| var_subs
	;

var_subs
	returns [[]tree.TextItem item]: 
	(Dollar|GrabDollar) (LCurly|GrabLCurly)
	value_ref
	(RCurly|GrabRCurly)
	;

value_ref
	returns [*tree.ValueRef vr]:
		ident         #value_ref_id 
	|	func_invoc    #value_ref_func
	;

uninterp
	returns[[]tree.TextItem item]:
	DoubleLess {log.Printf("xxx got double less\n")} (
			uninterp_inner 
	)* GrabDoubleGreater {log.Printf("xxx got double greater %+v",localctx.GetItem())}
	;

uninterp_inner 
	returns [[]tree.TextItem Item]:
	RawText
	|var_subs
	;

param_spec
	returns[[]*tree.PFormal formal]: 
	LParen (param_pair (Comma param_pair)*)?  RParen;

param_pair
	returns[*tree.PFormal formal]:
	Id (TypeStarter)? ident;

simple_or_model_param
	returns [*tree.TypeDecl t]:
	id1=Id
	|Colon id2=Id 
	;

doc_section
	returns [*tree.DocSectionNode section]: 
	Doc (doc_func)*;

doc_func
	returns [*tree.DocFuncNode fn]:
	Id doc_func_post
	;

doc_func_post
	returns [*tree.DocFuncNode fn]:
	doc_func_formal doc_func_local? pre_code? doc_elem post_code?
	;

doc_func_local
	returns [[]*tree.PFormal formal]:
	Local param_spec
	;

doc_func_formal
	returns [[]*tree.PFormal formal]:
	param_spec
	|
	;

doc_tag
	returns [*tree.DocTag tag]:
	LessThan value_ref
	doc_id?
	doc_class?
	GreaterThan 
	;

doc_id
	returns [*tree.ValueRef s]:
	Hash value_ref 
	;

doc_class
	returns [[]*tree.ValueRef clazz]:
	(value_ref)+ 
	;

doc_elem
	returns [*tree.DocElement elem]:
	value_ref                    # haveVar
	| doc_tag uninterp?  # haveTag
	| doc_elem_child             # haveList
	;

//doc_elem_content
//	returns [*tree.DocElement element]:
//	doc_elem_text
//	| doc_elem_child 
//	;

//doc_elem_text
//	returns [*tree.FuncInvoc invoc]:
//	text_top      
//	;

doc_elem_child
	returns [*tree.DocElement elem]:
	LParen (doc_elem)* RParen
	;

func_invoc
	returns [*tree.FuncInvoc invoc]:
	(Id|GrabId) {log.Printf("got id or grab id\n")}(LParen|GrabLParen) func_actual_seq (RParen|GrabRParen)
	;


func_actual_seq
	returns [[]*tree.FuncActual actual]:
	( func_actual ( (Comma|GrabComma) func_actual)* )?
	;

func_actual 
	returns [*tree.FuncActual actual]:
	value_ref
	;

event_section
	returns [*tree.EventSectionNode section]:
	Event {log.Printf("GOT @EVENT\n")}(event_spec)*;

event_spec
	returns [*tree.EventSpec spec]:
	selector Id event_call
	;

event_call
	returns [*tree.FuncInvoc invoc]:
	(Arrow)? func_invoc 
	;

selector
	returns [*tree.Selector sel]:
	Hash id=value_ref
	| class=value_ref // must start with a dot
	;

mvc_section
	returns [*tree.MVCSectionNode section]:
	Mvc Model model_decl+
	(View view_decl*)?
	;

model_decl
	returns [*tree.ModelDecl decl]:
	id1=Id filename_seq
	;

view_decl
	returns [*tree.ViewDecl vdecl]:
	vname=Id doc_func_post
	;

filename_seq
	returns [[]string seq]: 
	StringLit (Comma StringLit)*
	;

ident
	returns [*tree.Ident id]:
	(Colon|GrabColon)? (Id|GrabId)
	(
		(dot_qual)*
		| (colon_qual)*
	)
	;

dot_qual 
	returns [*tree.IdentPart part]: 
	(Dot|GrabDot) ident
	;

colon_qual
	returns [*tree.IdentPart part]: 
	(Colon|GrabColon) ident
	;