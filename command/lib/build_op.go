package lib

// EnterZeroOp is called when production zeroOp is entered.
func (b *Builder) EnterZeroOp(_ *ZeroOpContext) {}

// ExitZeroOp is called when production zeroOp is exited.
func (b *Builder) ExitZeroOp(ctx *ZeroOpContext) {
	b.currentStmt = &ZeroOp{ctx.GetToken(WasmLexerZeroOpWord, 0).GetText()}
}

// EnterArgOp is called when production argOp is entered.
func (b *Builder) EnterInt1Op(_ *ArgOpContext) {
}

// ExitArgOp is called when production argOp is exited.
func (b *Builder) ExitInt1Op(ctx *ArgOpContext) {
	op := &ArgOp{
		Op: ctx.ArgWord().GetText(),
	}
	switch {
	case ctx.Num() != nil:
		op.IntArg = new(int)
		*op.IntArg = numTerminalToInt(ctx.Num().GetSymbol())
	case ctx.HexFloatConst() != nil:
		op.FloatArg = new(string)
		*op.FloatArg = ctx.HexFloatConst().GetText()
	case ctx.StackPointerWord() != nil:
		op.Special = new(SpecialIdT)
		*op.Special = StackPointer
	}

	if ctx.BlockAnnotation() != nil {
		op.BranchAnno = new(int)
		*op.BranchAnno = blockAnnotationTerminalToInt(ctx.BlockAnnotation().GetSymbol())
	}
	if ctx.ConstAnnotation() != nil {
		op.ConstAnno = new(string)
		*op.ConstAnno = constAnnotationTerminalToString(ctx.ConstAnnotation().GetSymbol())
	}
	anno := ctx.GetToken(WasmLexerBlockAnnotation, 0)
	if anno != nil {
		op.BranchAnno = new(int)
		*op.BranchAnno = blockAnnotationTerminalToInt(anno.GetSymbol())
	}
	b.currentStmt = op
}

// EnterLoadStore is called when production loadStore is entered.
func (b *Builder) EnterLoadStore(_ *LoadStoreContext) {}

// ExitLoadStore is called when production i64LoadStore is exited.
func (b *Builder) ExitLoadStore(ctx *LoadStoreContext) {
	t := ctx.GetToken(WasmLexerLoadStore, 0).GetText()
	m := &LoadStoreOp{
		Op: t,
	}
	if ctx.Offset() != nil {
		m.Offset = new(int)
		*m.Offset = offsetTerminalToInt(ctx.Offset().GetSymbol())
	}
	if ctx.Align() != nil {
		m.Align = new(int)
		*m.Align = alignTerminalToInt(ctx.Align().GetSymbol())
	}
	b.currentStmt = m
}

// EnterCallOp is called when production callOp is entered.
func (b *Builder) EnterId1Op(_ *CallOpContext) {}

// ExitCallOp is called when production callOp is exited.
func (b *Builder) ExitId1Op(ctx *CallOpContext) {
	b.currentStmt = &CallOp{
		Op:  ctx.CallWord().GetText(),
		Arg: identTerminalToString(ctx.Ident().GetSymbol()),
	}
}

// EnterBrTable is called when production brTable is entered.
func (b *Builder) EnterBrTable(_ *BrTableContext) {}

// ExitBrTable is called when production brTable is exited.
func (b *Builder) ExitBrTable(ctx *BrTableContext) {
	num := ctx.GetTokens(WasmLexerNum)
	block := ctx.GetTokens(WasmLexerBlockAnnotation)
	if len(num) != len(block) {
		panic("mismatched num and block inside BrTable target") // should never happen
	}
	jump := &BrTable{}
	jump.Target = make([]*BranchTarget, len(num))
	for i := range num {
		jump.Target[i] = &BranchTarget{}
		jump.Target[i].Num = numTerminalToInt(num[i].GetSymbol())
		jump.Target[i].Block = blockAnnotationTerminalToInt(block[i].GetSymbol())
	}
	b.currentStmt = jump
}

// EnterCallIndirectOp is called when production callIndirectOp is entered.
func (b *Builder) EnterCallIndirectOp(_ *CallIndirectOpContext) {}

// ExitCallIndirectOp is called when production callIndirectOp is exited.
func (b *Builder) ExitCallIndirectOp(ctx *CallIndirectOpContext) {
	b.currentStmt = &IndirectCallOp{
		Type: ctx.TypeRef().GetT(),
	}
}

// EnterGlobalDef is called when production globalDef is entered.
func (b *Builder) EnterGlobalDef(_ *GlobalDefContext) {}

// ExitGlobalDef is called when production globalDef is exited.
func (b *Builder) ExitGlobalDef(ctx *GlobalDefContext) {
	op := &GlobalOp{}
	if ctx.TypeAnnotation() != nil {
		op.Anno = new(int)
		*op.Anno = typeAnnotationTerminalToInt(ctx.TypeAnnotation().GetSymbol())
	}
	op.Name = ctx.Ident().GetText()
	op.Type = ctx.TypeStmt().GetText()
	op.Value = ctx.ValueStmt().GetText()
	print("xxx")
}
