parser grammar wcl;
options {
	tokenVocab = wcllex;
}
@parser::header {
	import "github.com/iansmith/parigot/ui/parser/tree"
	var _ = &tree.ProgramNode{}
}

program
	returns[*tree.ProgramNode p]:
	wcl_section
	css_section?
	import_section?  
	extern?
	global?
	mvc_section?
	text_section? 
	css_section?     
	doc_section?
	event_section?
	EOF   
	;

global
	returns [[]*tree.PFormal g]:
	Global param_spec
	;

extern
	returns [[]string e]:
	Extern LParen Id* RParen
	;

wcl_section:
	Wcl Version
	;

import_section
	returns[*tree.ImportSectionNode section]:
	Import LCurly uninterp 
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
	i = Id param_spec? text_func_local? pre_code? text_top post_code?
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

text_top
	returns[[]tree.TextItem item]:
	BackTick 	(
		text_content
		|
	) ContentBackTick ;

text_content
	returns[[]tree.TextItem item]:
	(
		text_content_inner    
	)*;

text_content_inner
	returns[[]tree.TextItem item]:
		ContentRawText             	#RawText
		| var_subs   				#VarSub
;

var_subs
	returns [[]tree.TextItem item]: 
	ContentDollar sub
	;

sub
	returns [tree.TextItem item]: 
	VarId VarRCurly
	| func_invoc_var VarRCurly
	;

uninterp
	returns[[]tree.TextItem item]:
	(
		uninterp_inner  
	)+ UninterpRCurly;

uninterp_inner 
	returns [[]tree.TextItem Item]:
	UninterpRawText #UninterpRawText
	| UninterpLCurly uninterp  #UninterpNested
	| uninterp_var #UninterpVar
;

uninterp_var
	returns[[]tree.TextItem item]: 
	UninterpDollar VarId VarRCurly;

param_spec
	returns[[]*tree.PFormal formal]: 
	LParen (param_pair (Comma param_pair)*)?  RParen;

param_pair
	returns[*tree.PFormal formal]:
	Id simple_or_model_param;

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
	LessThan id_or_var_ref
	doc_id?
	doc_class?
	GreaterThan 
	;

id_or_var_ref
	returns [*tree.DocIdOrVar idVar]:
	Id
	| var_ref
	;

var_ref
	returns [*tree.DocIdOrVar v]:
	Dollar VarId VarRCurly
	;

doc_id
	returns [string s]:
	Hash Id 
	;

doc_class
	returns [[]string clazz]:
	 // these ids must start with a dot
	(Id)+ 
	;

doc_elem
	returns [*tree.DocElement elem]:
	var_ref                      # haveVar
	| doc_tag doc_elem_content?  # haveTag
	| doc_elem_child             # haveList
	;

doc_elem_content
	returns [*tree.DocElement element]:
	doc_elem_text
	| doc_elem_child 
	;

doc_elem_text
	returns [*tree.FuncInvoc invoc]:
	func_invoc       #doc_elem_text_func_call
	| text_top       #doc_elem_text_anon 
	;

doc_elem_child
	returns [*tree.DocElement elem]:
	LParen (doc_elem)* RParen
	;

func_invoc
	returns [*tree.FuncInvoc invoc]:
	Id LParen func_actual_seq RParen
	;

func_invoc_var
	returns [*tree.FuncInvoc invoc]:
	VarId VarLeftParen func_actual_seq_var VarRightParen
	;

func_actual_seq
	returns [[]*tree.FuncActual actual]:
	( a=func_actual (Comma b=func_actual)* )?
	;

func_actual_seq_var
	returns [[]*tree.FuncActual actual]:
	( a=func_actual_var (Comma b=func_actual_var)* )?
	;

func_actual 
	returns [*tree.FuncActual actual]:
	Id
	| StringLit
	;

func_actual_var
	returns [*tree.FuncActual actual]:
	VarId
	| VarStringLit
	;

event_section
	returns [*tree.EventSectionNode section]:
	Event (event_spec)*;

event_spec
	returns [*tree.EventSpec spec]:
	selector Id event_call
	;

event_call
	returns [*tree.FuncInvoc invoc]:
	(GreaterThan GreaterThan)? func_invoc 
	;

selector
	returns [*tree.Selector sel]:
	Hash IdValue=Id
	| class=Id // must start with a dot
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

