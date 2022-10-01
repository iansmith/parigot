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

zeroOp:
    ZeroOpWord
    ;

int1Op:
    Int1OpWord ( StackPointerWord | Num)
    ;

brIfOp:
    BrIfWord Num Lparen BlockAnnotation Rparen
    ;

id1Op:
    Id1OpWord Ident
    ;

i64Store:
    I64StoreWord (Offset)? (Align)?
    ;

i64Load:
    I64LoadWord (Offset)? (Align)?
    ;

i32Store:
    I32StoreWord (Offset)?
    ;

i32Load:
    I32LoadWord (Offset)?
    ;

stmt:
    blockStmt
    | ifStmt
    | zeroOp
    | int1Op
    | i32Store
    | i32Load
    | i64Store
    | i64Load
    | id1Op
    | brIfOp
    ;

blockStmt:
    BlockWord resultDef? stmt+ EndWord
    ;

ifStmt:
    IfWord resultDef? stmt+ (elsePart stmt+)? EndWord
    ;

// slightly hacky: I use this construction to get a call Enter/ExitElsePart on the builder
elsePart:
    ElseWord
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
IfWord: 'if';
ElseWord: 'else';
EndWord: 'end';

fragment I32: 'i32';
fragment I64: 'i64';
fragment F64: 'f64';
TypeName: I32 | I64 | F64;

// opcodish
ZeroOpWord: 'i32.sub' | 'select' | 'i32.eqz' |
    'return' | 'i32.eq' | 'drop' | 'unreachable' | 'i32.add';

Int1OpWord: 'global.get' | 'i32.const' | 'local.tee' | 'local.get' | 'local.set' | 'global.set';
BrIfWord: 'br_if';

Id1OpWord: 'call';

I64StoreWord: 'i64.store';
I64LoadWord: 'i64.load';
I32StoreWord: 'i32.store';
I32LoadWord: 'i32.load';

StackPointerWord: '$__stack_pointer';

// regular Lexer rules after here
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
