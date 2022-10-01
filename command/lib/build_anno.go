package lib

// EnterTypeAnno is called when production typeAnno is entered.
func (b *Builder) EnterTypeAnno(ctx *TypeAnnoContext) {
}

// ExitTypeAnno is called when production typeAnno is exited.
func (b *Builder) ExitTypeAnno(ctx *TypeAnnoContext) {
	b.currentTypeAnnotation = &TypeAnnotation{
		Number: typeAnnotationTerminalToInt(ctx.GetToken(WasmLexerTypeAnnotation, 0).GetSymbol()),
	}
}
