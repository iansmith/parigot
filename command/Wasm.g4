grammar Wasm;

module returns [*Module m]:
    Lparen ModuleWord topLevelSeq Rparen Rparen
    {
        m:=&Module{
            Code: $topLevelSeq.t,
        }
        localctx.SetM(m)
    }
    ;

topLevelSeq returns [[]TopLevel t]:
    | tl+=topLevel (tl+=topLevel)*
    {
        result:=make([]TopLevel,len(localctx.GetTl()))
        for i,n:=range localctx.GetTl() {
            result[i]=n.GetT()
        }
        localctx.SetT(result)
    }
    ;

topLevel returns [TopLevel t]:
    Lparen
    (
        typeDef { $t=$typeDef.t }
        | importDef { $t=$importDef.i }
        | funcDef { $t=$funcDef.f }
        | tableDef { $t=$tableDef.t}
        | memoryDef { $t=$memoryDef.m }
        | globalDef { $t=$globalDef.g }
    )
    Rparen
    ;

typeDef returns [*TypeDef t]:
    TypeWord typeAnno f=funcSpec
    {
        localctx.SetT(
            &TypeDef{
                Annotation: $typeAnno.t,
                Func: localctx.GetF().GetF(),
            },
         );
    }
    ;

importDef returns [*ImportDef i]:
    ImportWord m=QuotedString im=QuotedString funcNameRef
    {
        moduleName:=localctx.GetM().GetText()[1:len(localctx.GetM().GetText())-1]
        importedAs:=localctx.GetIm().GetText()[1:len(localctx.GetIm().GetText())-1]
        localctx.SetI(
            &ImportDef{
                ModuleName:moduleName,
                ImportedAs:importedAs,
            },
         );
    }
    ;

typeRef returns [*TypeRef t]:
    Lparen TypeWord Num Rparen
    {
        localctx.SetT(&TypeRef{Num:numToInt($Num.GetText())})
    }
    ;

typeAnno returns [int t]:
    Lparen TypeAnnotation Rparen
    {
        localctx.SetT(annoToInt($TypeAnnotation.GetText(),false))
    }
    ;

branchAnno returns [int t]:
    Lparen BranchAnnotation Rparen
    {
        localctx.SetT(annoToInt($BranchAnnotation.GetText(),true))
    }
    ;

constAnno returns [string t]:
    Lparen ConstAnnotation Rparen
    {
        localctx.SetT(annoToString($ConstAnnotation.GetText(),true))
    }
    ;

funcSpec returns [*FuncSpec f]:
    Lparen FuncWord paramDef? resultDef? Rparen
    {
        var pd *ParamDef
        var r *ResultDef

        if localctx.Get_paramDef()!=nil {
            pd=$paramDef.p
        }
        if localctx.Get_resultDef()!=nil {
            r=$resultDef.r
        }

        localctx.SetF(&FuncSpec{
            Param: pd,
            Result: r,
        });
    }
    ;

funcNameRef returns [*FuncNameRef f]:
    Lparen FuncWord i=Ident t=typeRef Rparen
    {
        localctx.SetF(
            &FuncNameRef{
                Name: localctx.GetI().GetText(),
                Type: localctx.GetT().GetT(),
            },
        )
    }
    ;

paramDef returns [*ParamDef p]: Lparen ParamWord typeNameSeq Rparen
    {
        result:=TypeNameSeq{Name:$typeNameSeq.t}
        localctx.SetP(&ParamDef{&result})
    }
    ;

typeNameSeq returns [[]string t]:
    tn+=TypeName (tn+=TypeName)*
    {
        result:=make([]string,len(localctx.GetTn()))
        for i,n:=range localctx.GetTn() {
            result[i]=n.GetText()
        }
        $t=result
    }

    ;

resultDef returns [*ResultDef r]: Lparen ResultWord typeNameSeq Rparen
    {
        result:=TypeNameSeq{Name:$typeNameSeq.t}
        localctx.SetR(&ResultDef{&result})
    }
    ;

localDef returns [*LocalDef l]: Lparen LocalWord typeNameSeq Rparen
    {
        result:=TypeNameSeq{Name:$typeNameSeq.t}
        localctx.SetL(&LocalDef{&result})
    }
    ;

funcDef returns [*FuncDef f]:
    FuncWord Ident typeRef paramDef? resultDef? localDef? funcBody
    {
        var pd *ParamDef
        var r *ResultDef
        var l *LocalDef

        if localctx.Get_paramDef()!=nil {
            pd=$paramDef.p
        }
        if localctx.Get_resultDef()!=nil {
            r=$resultDef.r
        }
        if localctx.Get_localDef()!=nil {
            l=$localDef.l
        }

        localctx.SetF(
        &FuncDef{
            Name:$Ident.GetText(),
            Type: $typeRef.t,
            Param: pd,
            Result: r,
            Local: l,
            Code: $funcBody.f,
        })
    }
    ;

funcBody returns [[]Stmt f]:
    stmtSeq
    {
        $f=($stmtSeq.s)
    }
    ;

stmtSeq returns [[]Stmt s]:
    st+=stmt (st+=stmt)*
    {
        result:=make([]Stmt,len(localctx.GetSt()))
        for i,n:=range localctx.GetSt() {
            result[i]=n.GetS()
        }
        localctx.SetS(result)
    }
    ;

zeroOp returns [Stmt z]:
    o=ZeroOpWord
    {
        localctx.SetZ(&ZeroOp{localctx.GetO().GetText()})
    }
    ;

argOp returns [Stmt a]:
    o=ArgWord ( s=StackPointerWord | n=Num | h=HexFloatConst) (branchAnno|c=constAnno) ?
    {
        op:=&ArgOp{Op:localctx.GetO().GetText()}
        if localctx.GetS()!=nil {
            op.Special=new(SpecialIdT)
            *op.Special=StackPointer
        }
        if localctx.GetN()!=nil {
            op.IntArg=new(int)
            *op.IntArg=numToInt(localctx.GetN().GetText())
        }
        if localctx.GetH()!=nil {
            op.FloatArg=new(string)
            *op.FloatArg=localctx.GetH().GetText()
        }
        if localctx.Get_branchAnno()!=nil {
            op.BranchAnno = new(int)
            *op.BranchAnno= $branchAnno.t
        }
        if localctx.GetC()!=nil {
            op.ConstAnno = new(string)
            *op.ConstAnno= $constAnno.t
        }

    }
    ;

callOp returns [Stmt c]:
    CallWord i=Ident
    {
        localctx.SetC(
            &CallOp{
                Arg:localctx.GetI().GetText(),
            },
        )
    }
    ;

callIndirectOp returns [Stmt c]:
    CallIndirectWord t=typeRef
    {
        localctx.SetC(
            &IndirectCallOp{
                Type: localctx.GetT().GetT(),
            },
        )
    }
    ;

loadStore returns [Stmt l]:
    lo=LoadStore (o=Offset)? (a=Align)?
    {
        op:=&LoadStoreOp{
            Op:localctx.GetLo().GetText(),
        }
        if localctx.GetO()!=nil {
            op.Offset=new(int)
            *op.Offset=numToInt(localctx.GetO().GetText()[len("offset="):])
        }
        if localctx.GetA()!=nil {
            op.Align=new(int)
            *op.Align=numToInt(localctx.GetA().GetText()[len("align="):])
        }
        localctx.SetL(op)
    }
    ;

brTable returns [Stmt b]:
    BrTableWord br=brTableTargetSeq
    {
        localctx.SetB(
            &BrTableOp{Target:localctx.GetBr().GetB()},
        )
    }
    ;

brTableTargetSeq returns [[]*BranchTarget b]:
    t+=brTableTarget (t+=brTableTarget)*
    ;

brTableTarget returns [*BranchTarget b]:
    n=Num branchAnno
    {
        localctx.SetB(&BranchTarget{
            Num:numToInt(localctx.GetN().GetText()),
            Branch:$branchAnno.t,
            })
    }
    ;

globalDef returns [TopLevel g]:
    GlobalWord typeAnno i=Ident t=stmt v=stmt
    {
        op:=&GlobalOp{
            Name:$i.GetText(),
            Type: $t.s,
            Value: $v.s,
        }
        if localctx.Get_typeAnno!=nil {
            op.Anno=new(int)
            *op.Anno=$typeAnno.t
        }
        localctx.SetG(op)
    }
    ;

mut returns [Stmt m]:
    MutWord t=TypeName
    {
        localctx.SetM(&MutOp{localctx.GetT().GetText()})
    }
    ;

stmt returns [Stmt s]:
    b=blockStmt {  $s = $blockStmt.b }
    | f=ifStmt  { localctx.SetS(localctx.GetF().GetI()) }
    | l=loopStmt  { localctx.SetS(localctx.GetL().GetL()) }
    | z=zeroOp { localctx.SetS(localctx.GetZ().GetZ()) }
    | a=argOp { localctx.SetS(localctx.GetA().GetA()) }
    | ls=loadStore { localctx.SetS(localctx.GetLs().GetL()) }
    | c=callOp { localctx.SetS(localctx.GetC().GetC()) }
    | i=callIndirectOp { localctx.SetS(localctx.GetI().GetC()) }
    | br=brTable { localctx.SetS(localctx.GetBr().GetB()) }
    ;

blockStmt returns [Stmt b]:
    BlockWord resultDef? stmtSeq EndWord
    {
        block:=&BlockStmt{
            Code:$stmtSeq.s,
        }
        if localctx.Get_resultDef()!=nil {
            block.Result = $resultDef.r
        }
        $b=block
    }
    ;

loopStmt returns [Stmt l]:
    LoopWord stmtSeq EndWord
    {
        localctx.SetL(&LoopStmt{&BlockStmt{Code:$stmtSeq.s}})
    }
    ;

ifStmt returns [Stmt i]:
    IfWord resultDef? s1=stmtSeq (ElseWord s2=stmtSeq)? EndWord
    {
        var r *ResultDef
        if localctx.Get_resultDef() !=nil {
            r = $resultDef.r
        }
        ifStmt:=&IfStmt{
            IfPart:$s1.s,
            Result: r,
        }
        if $ctx.s2!=nil {
            ifStmt.ElsePart=$s2.s
        }
        $i=ifStmt
    }
    ;

tableDef returns [TopLevel t]:
    TableWord typeAnno? min=Num max=Num FuncRefWord
    {
        op:=&TableDef{Min:numToInt($min.GetText()),Max:numToInt($max.GetText())}
        if localctx.Get_typeAnno()!=nil {
            op.Type = new(int)
            *op.Type = $typeAnno.t
        }
        $t=op
    }
    ;

memoryDef returns [TopLevel m]:
    MemoryWord typeAnno? size=Num
    {
        op:=&MemoryDef{Size:numToInt($size.GetText())}
        if localctx.Get_typeAnno()!=nil {
            op.Type = new(int)
            *op.Type = $typeAnno.t
        }
        $m=op
    }
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
LoopWord: 'loop';
BrTableWord: 'br_table';
GlobalWord: 'global';
MutWord: 'mut';
TableWord: 'table';
FuncRefWord: 'funcref';
TypeName: 'i32' | 'i64' | 'f64' | 'f32';
MemoryWord:'memory';

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
Ident:IdentFirst IdentAfter*;

fragment Digit: '0'..'9';
ConstValue: ('-')?  '0' | ('-')? Digit '.' Digit+ 'e' ('+'|'-') Digit (Digit)? ;

//HexPointer: ('-')? '0x' ( '0' .. '9')+ 'p+' ('0'..'9')+;
Offset: 'offset=' ( '0' .. '9')+;
Align: 'align=' ( '0' .. '9')+;
//ConstValue: IntConst | FloatConst;

// annotations look like ;blah;
ConstAnnotation: ';' '=' ConstValue ';' ;
BranchAnnotation: ';' '@' Digit+ ';' ;
TypeAnnotation: ';'  Digit+ ';' ;

fragment HexByteValue: '\\' ( '0' .. '9' | 'a' .. 'f') ( '0' .. '9' | 'a' .. 'f');
QuotedString: '"' ( HexByteValue | ~('"') )* '"';

Comment: ';;' ~( '\n' | '\r')* ('\r' '\n'| '\n') -> skip;