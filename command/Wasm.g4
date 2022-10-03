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

typeRef returns [*TypeRef t]:
    Lparen TypeWord Num Rparen
    {
        localctx.SetT(&TypeRef{Num:TokenToInt($Num)})
    }
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
    FuncWord Ident typeRef paramDef? resultDef? localDef? funcBody
    ;

funcBody:
    stmt+
    ;

zeroOp:
    ZeroOpWord
    ;

argOp:
    ArgWord ( StackPointerWord | Num | HexFloatConst) (Lparen (BlockAnnotation|ConstAnnotation) Rparen)?
    ;

callOp:
    CallWord Ident
    ;

callIndirectOp:
    CallIndirectWord typeRef
    ;

loadStore:
    LoadStore (Offset)? (Align)?
    ;

brTable:
    BrTableWord (Num Lparen BlockAnnotation Rparen)+
    ;

intConst:
    ConstIntWord Num
    ;

globalDef:
    GlobalWord TypeAnnotation? Ident typeStmt valueStmt
    ;

typeStmt:
    stmt
    ;

valueStmt:
    stmt
    ;

stmt:
    blockStmt
    | ifStmt
    | loopStmt
    | zeroOp
    | argOp
    | loadStore
    | callOp
    | callIndirectOp
    | brTable
    | globalDef
    ;

blockStmt:
    BlockWord resultDef? stmt+ EndWord
    ;

loopStmt:
    LoopWord stmt+ EndWord
    ;

ifStmt:
    IfWord resultDef? stmt+ (elsePart stmt+)? EndWord
    ;

// slightly hacky: I use this construction to get a call Enter/ExitElsePart on the builder
elsePart:
    ElseWord
    ;

    ////////// older ////////


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
LoopWord: 'loop';
BrTableWord: 'br_table';
GlobalWord: 'global';

TypeName: 'i32' | 'i64' | 'f64' | 'f32';

fragment HexDigit: ('0' .. '9' | 'a'..'f');
HexFloatConst: ('-')? ('0x')? HexDigit+ ('.' HexDigit+)? 'p' ('+' | '-') Digit+ ;

// op with no params (uses stack only)
ZeroOpWord:
    Memory
    | IntegerMath
    | IntegerBitwise
    | IntegerComp
    | IntegerUnary
    | FloatMath
    | FloatUnary
    | TypeRepresentation
    | ControlFlow0
    | Extend
    ;

Memory: 'memory.size'| 'memory.grow'| 'memory.copy' | 'memory.fill';
IntegerMath:
    'i32.add' | 'i32.sub' | 'i32.mul' | 'i32.div' | 'i32.div_s'| 'i32.div_u'|'i32.rem_s'| 'i32.rem_u'|
    'i64.add' | 'i64.sub' | 'i64.mul' | 'i64.div' | 'i64.div_s'| 'i64.div_u'|'i64.rem_s'| 'i64.rem_u';
IntegerBitwise:
    'i32.and' | 'i32.or' | 'i32.xor' | 'i32.shl' | 'i32.shr_u' | 'i32.shr_s' | 'i32.rotl' | 'i32.rotr'|
    'i64.and' | 'i64.or' | 'i64.xor' | 'i64.shl' | 'i64.shr_u' | 'i64.shr_s' | 'i64.rotl' | 'i64.rotr';
IntegerComp:
    'i32.eq' | 'i32.ne' | 'i32.lt_s' | 'i32.lt_u' | 'i32.le_s' | 'i32.le_u' | 'i32.gt_s' | 'i32.gt_u' | 'i32.ge_s' | 'i32.ge_u' |
    'i64.eq' | 'i64.ne' | 'i64.lt_s' | 'i64.lt_u' | 'i64.le_s' | 'i64.le_u' | 'i64.gt_s' | 'i64.gt_u' | 'i64.ge_s' | 'i64.ge_u';
IntegerUnary:
    'i32.clz' | 'i32.ctz' | 'i32.popcn' | 'i32.eqz' |
    'i64.clz' | 'i64.ctz' | 'i64.popcn' | 'i64.eqz';
FloatMath:
    'f32.add' | 'f32.sub'| 'f32.mul' | 'f32.div'| 'f32.copysign' | 'f32.eq'| 'f32.ne' | 'f32.lt' | 'f32.le' | 'f32.gt' |'f32.ge' | 'f32.min' | 'f32.max'|
    'f64.add' | 'f64.sub'| 'f64.mul' | 'f64.div'| 'f64.copysign' | 'f64.eq'| 'f64.ne' | 'f64.lt' | 'f64.le' | 'f64.gt' |'f64.ge' | 'f64.min' | 'f64.max';
FloatUnary:
    'f32.abs' | 'f32.neg' | 'f32.ceil' | 'f32.floor' | 'f32.trunc' | 'f32.nearest' | 'f32.sqrt'|
    'f64.abs' | 'f64.neg' | 'f64.ceil' | 'f64.floor' | 'f64.trunc' | 'f64.nearest' | 'f64.sqrt';

TypeRepresentation:
    'i32.wrap_i64' |
    'i32.trunc_f32' | 'i32.trunc_sat_f32' | 'i32.trunc_f64'| 'i32.trunc_sat_f64_u' | 'i32.trunc_sat_f64_s'| 'i32.reinterpret_f32'|
    'i64.trunc_f32' | 'i64.trunc_sat_f32' | 'i64.trunc_f64'| 'i64.trunc_sat_f64_u' | 'i64.trunc_sat_f64_s'| 'i64.reinterpret_f64'|
    'f32.demote_f64' | 'f32.convert_i32_s' | 'f32.convert_i64_s' | 'f32.convert_i32_u' | 'f32.convert_i64_u' | 'f32.reinterpret_i32' |
    'f64.promote_f32' | 'f64.convert_i32_s'| 'f64.convert_i64_s' | 'f64.convert_i32_u' | 'f64.convert_i64_u' | 'f64.reinterpret_i64';

Extend:
    'i64.extend_i32_s' | 'i64.extend_i32_u';

ControlFlow0:
   'nop'|'unreachable'|'drop'|'return'|'select';

BulkMemory:
    'data.drop' | 'elem.drop' | 'table.copy' | 'table.init';

// op with 1 parameter which is an integer
ArgWord:
    Variables
    | Branch1
    | ConstIntWord
    | ConstFloatWord
    ;

ConstIntWord: 'i32.const' | 'i64.const';
ConstFloatWord:  'f32.const' | 'f64.const';

Variables: 'global.get' |'global.set' | 'local.get' | 'local.set' | 'global.tee' | 'local.tee';
Branch1:  'br'|'br_if';

// op 1 parameter which is an ident
CallWord: 'call';

CallIndirectWord: 'call_indirect';

LoadStore:
    'i64.store'| 'i64.store32' | 'i64.store16' | 'i64.store8' |
    'i64.load' | 'i64.load8_s' | 'i64.load8_u' | 'i64.load16_s' | 'i64.load16_u' | 'i64.load32_s' | 'i64.load32_u' |
    'i32.load' | 'i32.load8_s' | 'i32.load8_u' | 'i32.load16_s' | 'i32.load16_u' |
    'i32.store' | 'i32.store8' | 'i32.store16' |
    'f32.store' | 'f64.store' |
    'f32.load' | 'f64.load';

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
fragment Digit: '0'..'9';
ConstValue: ('-')?  '0' | ('-')? Digit '.' Digit+ 'e' ('+'|'-') Digit (Digit)? ;

//HexPointer: ('-')? '0x' ( '0' .. '9')+ 'p+' ('0'..'9')+;
Offset: 'offset=' ( '0' .. '9')+;
Align: 'align=' ( '0' .. '9')+;
//ConstValue: IntConst | FloatConst;

// annotations look like ;blah;
ConstAnnotation: ';' '=' ConstValue ';' ;
BlockAnnotation: ';' '@' Digit+ ';' ;
TypeAnnotation: ';'  Digit+ ';' ;

fragment HexByteValue: '\\' ( '0' .. '9' | 'a' .. 'f') ( '0' .. '9' | 'a' .. 'f');
QuotedString: '"' ( HexByteValue | ~('"') )* '"';

Comment: ';;' ~( '\n' | '\r')* ('\r' '\n'| '\n') -> skip;
