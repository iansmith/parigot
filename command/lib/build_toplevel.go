package lib

// EnterTypeDef is called when entering the typeDef production.
func (b *Builder) EnterTypeDef(_ *TypeDefContext) {
}

// ExitTypeDef is called when exiting the typeDef production.
func (b *Builder) ExitTypeDef(_ *TypeDefContext) {
	td := &TypeDef{
		Annotation: b.currentTypeAnnotation,
		Func:       b.currentFuncSpec, //note: funcSpec!
	}
	b.currentTypeAnnotation = nil
	b.currentFuncSpec = nil
	b.currentTopLevelDef = td
}

// EnterFuncDef is called when entering the funcDef production.
func (b *Builder) EnterFuncDef(_ *FuncDefContext) {
	b.currentContainer = &FuncDef{}
}

// ExitFuncDef is called when exiting the funcDef production.
func (b *Builder) ExitFuncDef(ctx *FuncDefContext) {
	fd := b.currentContainer.(*FuncDef)
	fd.Name = stringTerminalToString(ctx.GetToken(WasmLexerIdent, 0).GetSymbol())
	fd.Type = ctx.TypeRef().GetT()
	fd.Param = b.currentParamDef
	fd.Local = b.currentLocalDef
	fd.Result = b.currentResultDef
	b.currentParamDef = nil
	b.currentLocalDef = nil
	b.currentResultDef = nil
	b.currentTopLevelDef = fd
}

// ExitImportDef is called when exiting the importDef production.
func (b *Builder) ExitImportDef(c *ImportDefContext) {
	id := &ImportDef{
		ModuleName:  quotedStringTerminalToString(c.GetToken(WasmLexerQuotedString, 0).GetSymbol()),
		ImportedAs:  quotedStringTerminalToString(c.GetToken(WasmLexerQuotedString, 1).GetSymbol()),
		FuncNameRef: b.currentFuncNameRef,
	}
	b.currentFuncNameRef = nil
	b.currentTopLevelDef = id
}

// EnterImportDef is called when entering the importDef production.
func (b *Builder) EnterImportDef(_ *ImportDefContext) {
}

// EnterTopLevel is called when production topLevel is entered.
func (b *Builder) EnterTopLevel(_ *TopLevelContext) {
}

// ExitTopLevel is called when production topLevel is exited.
func (b *Builder) ExitTopLevel(_ *TopLevelContext) {
	b.currentModule.AddTopLevelDef(b.currentTopLevelDef)
	b.currentTopLevelDef = nil
	return
}
