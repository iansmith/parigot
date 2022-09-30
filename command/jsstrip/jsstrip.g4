grammar jsstrip;
module:
    Lparen ModuleWord (Lparen topLevel Rparen)* Rparen
    ;

topLevel:
    typeDef
    | importOp
    | funcDef
    ;

typeDef:
    TypeWord typeAnnotation funcSpec
    ;

importOp:
    ImportWord QuotedString QuotedString funcRef
    ;

typeRef:
    Lparen TypeWord Num Rparen
    ;

typeAnnotation:
    Lparen TypeAnnotation Rparen
    ;

funcSpec:
    Lparen FuncWord paramDef? resultDef? Rparen
    ;

funcRef:
    Lparen FuncWord Ident Lparen TypeWord Num Rparen Rparen
    ;

type_:
    I32 | I64 | F64
    ;

typeSeq:
    type_+
    ;

paramDef: Lparen ParamWord typeSeq Rparen
    ;

resultDef: Lparen ResultWord typeSeq Rparen
    ;

localDef: Lparen LocalWord typeSeq Rparen
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

I32: 'i32';
I64: 'i64';
F64: 'f64';

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
