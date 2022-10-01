package lib

import "github.com/antlr/antlr4/runtime/Go/antlr"

// EnterZeroOp is called when production zeroOp is entered.
func (b *Builder) EnterZeroOp(_ *ZeroOpContext) {}

// ExitZeroOp is called when production zeroOp is exited.
func (b *Builder) ExitZeroOp(ctx *ZeroOpContext) {
	b.currentStmt = &ZeroOp{ctx.GetToken(WasmLexerZeroOpWord, 0).GetText()}
}

// EnterInt1Op is called when production int1Op is entered.
func (b *Builder) EnterInt1Op(_ *Int1OpContext) {}

// ExitInt1Op is called when production int1Op is exited.
func (b *Builder) ExitInt1Op(ctx *Int1OpContext) {
	op := &Int1Op{
		Op:  ctx.GetToken(WasmLexerInt1OpWord, 0).GetText(),
		Arg: numTerminalToInt(ctx.GetToken(WasmLexerNum, 0).GetSymbol()),
	}
	anno := ctx.GetToken(WasmLexerBlockAnnotation, 0)
	if anno != nil {
		op.BranchTarget = new(int)
		*op.BranchTarget = blockAnnotationTerminalToInt(anno.GetSymbol())
	}
	b.currentStmt = op
}

// EnterI64Store is called when production i64Store is entered.
func (b *Builder) EnterI64Store(_ *I64StoreContext) {}

// ExitI64Store is called when production i64Store is exited.
func (b *Builder) ExitI64Store(ctx *I64StoreContext) {
	op := i64LoadStore(true, ctx.BaseParserRuleContext)
	b.currentStmt = op
}

// EnterI64Load is called when production i64Load is entered.
func (b *Builder) EnterI64Load(ctx *I64LoadContext) {}

// ExitI64Load is called when production i64Load is exited.
func (b *Builder) ExitI64Load(ctx *I64LoadContext) {
	op := i64LoadStore(false, ctx.BaseParserRuleContext)
	b.currentStmt = op
}

// EnterI32Load is called when production i32Load is entered.
func (b *Builder) EnterI32Load(_ *I32LoadContext) {}

// ExitI32Load is called when production i32Load is exited.
func (b *Builder) ExitI32Load(ctx *I32LoadContext) {
	op := i32LoadStore(false, ctx.BaseParserRuleContext)
	b.currentStmt = op
}

// EnterI32Store is called when production i32Store is entered.
func (b *Builder) EnterI32Store(_ *I32StoreContext) {}

// ExitI32Store is called when production i32Store is exited.
func (b *Builder) ExitI32Store(ctx *I32StoreContext) {
	op := i32LoadStore(true, ctx.BaseParserRuleContext)
	b.currentStmt = op
}

// EnterId1Op is called when production id1Op is entered.
func (b *Builder) EnterId1Op(_ *Id1OpContext) {}

// ExitId1Op is called when production id1Op is exited.
func (b *Builder) ExitId1Op(ctx *Id1OpContext) {
	b.currentStmt = &Id1Op{
		Op:  ctx.GetToken(WasmLexerId1OpWord, 0).GetText(),
		Arg: identTerminalToString(ctx.GetToken(WasmLexerIdent, 0).GetSymbol()),
	}
}

// helper for load and store i32
func i32LoadStore(isStore bool, ctx *antlr.BaseParserRuleContext) *I32LoadStore {
	op := &I32LoadStore{
		IsStore: isStore,
	}
	offsetTerminal := ctx.GetToken(WasmLexerOffset, 0)
	if offsetTerminal != nil {
		op.SetOffset(offsetTerminalToInt(offsetTerminal.GetSymbol()))
	}
	return op
}

// helper for load and store i64
func i64LoadStore(isStore bool, ctx *antlr.BaseParserRuleContext) *I64LoadStore {
	op := &I64LoadStore{
		IsStore: isStore,
	}
	offsetTerminal := ctx.GetToken(WasmLexerOffset, 0)
	alignTerminal := ctx.GetToken(WasmLexerAlign, 0)
	if offsetTerminal != nil {
		op.SetOffset(offsetTerminalToInt(offsetTerminal.GetSymbol()))
	}
	if alignTerminal != nil {
		op.SetAlign(alignTerminalToInt(alignTerminal.GetSymbol()))
	}
	return op
}

// EnterBrTable is called when production brTable is entered.
func (b *Builder) EnterBrTable(_ *BrTableContext) {}

// ExitBrTable is called when production brTable is exited.
func (b *Builder) ExitBrTable(ctx *BrTableContext) {
	num := ctx.GetTokens(WasmLexerNum)
	block := ctx.GetTokens(WasmLexerBlockAnnotation)
	if len(num) != len(block) {
		panic("misatched num and block inside BrTable target") // should never happen
	}
	jump := &BrTable{}
	jump.Target = make([]*BranchTarget, len(num))
	for i := range num {
		jump.Target[i].Num = numTerminalToInt(num[i].GetSymbol())
		jump.Target[i].Block = blockAnnotationTerminalToInt(block[i].GetSymbol())
	}
	b.currentStmt = jump
}
