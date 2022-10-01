grammar Wasm;

module:
    Lparen ModuleWord (Lparen topLevel Rparen)* Rparen
    ;

topLevel:
    typeDef
    | importDef
    | funcDef
    ;

typeDef:
    TypeWord typeAnno funcSpec
    ;

importDef:
    ImportWord QuotedString QuotedString funcNameRef
    ;

typeRef:
    Lparen TypeWord Num Rparen
    ;

typeAnno:
    Lparen TypeAnnotation Rparen
    ;

funcSpec:
    Lparen FuncWord paramDef? resultDef? Rparen
    ;

funcNameRef:
    Lparen FuncWord Ident typeRef Rparen
    ;

paramDef: Lparen ParamWord TypeName+ Rparen
    ;

resultDef: Lparen ResultWord TypeName+ Rparen
    ;

localDef: Lparen LocalWord TypeName+ Rparen
    ;

funcDef:
    FuncWord Ident typeRef paramDef localDef? funcBody
    ;

funcBody:
    stmt+
    ;


stmt:
    block
    ;

block:
    BlockWord resultDef
    ;

    ////////// older ////////

sexpr:
   list
   |
   atom
   ;

list:
    Lparen Rparen
    | Lparen BlockAnnotation Rparen
    | Lparen ConstAnnotation Rparen
    | Lparen members Rparen
    ;

members:
    (sexpr)+
    ;

atom:
    Ident Offset? Align?
    | Num
    | QuotedString
    | HexPointer
    | ConstValue
    ;

// keywordish
ModuleWord: 'module' ;
TypeWord: 'type' ;
FuncWord: 'func' ;
ParamWord: 'param' ;
ResultWord: 'result' ;
ImportWord: 'import';
LocalWord: 'local';
BlockWord: 'block';

fragment I32: 'i32';
fragment I64: 'i64';
fragment F64: 'f64';
TypeName: I32 | I64 | F64;

Whitespace: ( ' ' | '\r' '\n' | '\n' | '\t' ) -> skip;

// need to put these simple ones ahead of the complex ones
Lparen: '(';
Rparen: ')';
Quote: '"';

Num: ('-')? ( '0' .. '9')+;
fragment IdentFirst: ('a' .. 'z' | 'A' .. 'Z' | '.' | '$' | '_' | '/' | '*' | '@') ;
fragment IdentAfter: ('a' .. 'z' | 'A' .. 'Z' | '.' | '$' | '_' | '/' | '*' | '@'| '0'..'9');
Ident: IdentFirst IdentAfter* ;

fragment IntConst: ('-')?  ('0' .. '9')+ ( '.' ('0' .. '9')+ ( 'e' ('+' | '-') ('0' .. '9')+)?)? ;
fragment FloatConst: ('-')? ('0x')? ('0' .. '9' | 'a'..'f')+ '.' ('0' .. '9' | 'a' .. 'f')+ 'p' ('+' | '-') ('0' .. '9') ;
HexPointer: ('-')? '0x' ( '0' .. '9')+ 'p+' ('0'..'9')+;
Offset: 'offset=' ( '0' .. '9')+;
Align: 'align=' ( '0' .. '9')+;
ConstValue: IntConst | FloatConst;

// annotations look like ;blah;
ConstAnnotation: ';' '=' IntConst ';' ;
BlockAnnotation: ';' '@' ( '0'..'9')+ ';' ;
TypeAnnotation: ';'  ( '0'..'9')+ ';' ;

fragment HexByteValue: '\\' ( '0' .. '9' | 'a' .. 'f') ( '0' .. '9' | 'a' .. 'f');
QuotedString: '"' ( HexByteValue | ~('"') )* '"';

Comment: ';;' ~( '\n' | '\r')* ('\r' '\n'| '\n') -> skip;
