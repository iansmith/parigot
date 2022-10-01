package lib

// EnterFuncNameRef is called when production funcNameRef is entered.
func (b *Builder) EnterFuncNameRef(ctx *FuncNameRefContext) {}

// ExitFuncNameRef is called when production funcNameRef is exited.
func (b *Builder) ExitFuncNameRef(ctx *FuncNameRefContext) {
	if b.currentFuncNameRef != nil {
		panic("previous currentFuncNameRef not consumed")
	}
	b.currentFuncNameRef = &FuncNameRef{
		Name: ctx.GetToken(WasmLexerIdent, 0).GetText(),
		Type: b.currentTypeRef,
	}
	b.currentTypeRef = nil
}

// EnterFuncSpec is called when production funcSpec is entered.
func (b *Builder) EnterFuncSpec(ctx *FuncSpecContext) {}

// ExitFuncSpec is called when production funcSpec is exited.
func (b *Builder) ExitFuncSpec(ctx *FuncSpecContext) {
	if b.currentFuncSpec != nil {
		panic("previous currentFuncSpec not consumed")
	}
	b.currentFuncSpec = &FuncSpec{
		Param:  b.currentParamDef,
		Result: b.currentResultDef,
	}
	b.currentParamDef = nil
	b.currentResultDef = nil
}

// EnterParamDef is called when production paramDef is entered.
func (b *Builder) EnterParamDef(_ *ParamDefContext) {}

// ExitParamDef is called when production paramDef is exited.
func (b *Builder) ExitParamDef(ctx *ParamDefContext) {
	seq := getTypeNameSeq(ctx.BaseParserRuleContext)
	b.currentParamDef = &ParamDef{
		Type: &TypeNameSeq{seq},
	}
}

// EnterResultDef is called when production resultDef is entered.
func (b *Builder) EnterResultDef(ctx *ResultDefContext) {}

// ExitResultDef is called when production resultDef is exited.
func (b *Builder) ExitResultDef(ctx *ResultDefContext) {
	seq := getTypeNameSeq(ctx.BaseParserRuleContext)
	b.currentResultDef = &ResultDef{
		Type: &TypeNameSeq{seq},
	}
}

// EnterLocalDef is called when production localDef is entered.
func (b *Builder) EnterLocalDef(ctx *LocalDefContext) {}

// ExitLocalDef is called when production localDef is exited.
func (b *Builder) ExitLocalDef(ctx *LocalDefContext) {
	if b.currentLocalDef != nil {
		panic("previous current local def not consumed")
	}
	seq := getTypeNameSeq(ctx.BaseParserRuleContext)
	b.currentLocalDef = &LocalDef{
		Type: &TypeNameSeq{seq},
	}
}

// EnterTypeRef is called when entering the typeRef production.
func (b *Builder) EnterTypeRef(_ *TypeRefContext) {}

// ExitTypeRef is called when production typeRef is exited.
func (b *Builder) ExitTypeRef(ctx *TypeRefContext) {
	if b.currentTypeRef != nil {
		panic("previous TypeRef was not consumed")
	}
	b.currentTypeRef = &TypeRef{Num: numTerminalToInt(ctx.GetToken(WasmLexerNum, 0).GetSymbol())}
}
