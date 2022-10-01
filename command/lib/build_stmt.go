package lib

// EnterStmt is called when production stmt is entered.
func (b *Builder) EnterStmt(_ *StmtContext) {}

// ExitStmt is called when production stmt is exited.
func (b *Builder) ExitStmt(_ *StmtContext) {
	b.currentContainer.AddStmt(b.currentStmt)
	b.currentStmt = nil
	return
}

// EnterBlockStmt is called when production blockStmt is entered.
func (b *Builder) EnterBlockStmt(_ *BlockStmtContext) {
	prev := b.currentContainer
	b.currentNestingLevel++
	b.currentContainer = &BlockStmt{PreviousContainer: prev}
}

// ExitBlockStmt is called when production block is exited.
func (b *Builder) ExitBlockStmt(ctx *BlockStmtContext) {
	// this block statement
	currentBlock := b.currentContainer.(*BlockStmt)
	currentBlock.Result = b.currentResultDef
	currentBlock.nestingLevel = b.currentNestingLevel

	// setup the statement for ExitStmt
	b.currentStmt = currentBlock

	// pop stack of containers
	b.currentContainer = currentBlock.PreviousContainer
	b.currentNestingLevel--
}

// EnterIfStmt is called when production ifStmt is entered.
func (b *Builder) EnterIfStmt(ctx *IfStmtContext) {
	prev := b.currentContainer
	b.currentNestingLevel++
	b.currentContainer = &IfStmt{PreviousContainer: prev}
}

// ExitIfStmt is called when production ifStmt is exited.
func (b *Builder) ExitIfStmt(ctx *IfStmtContext) {
	// this block statement
	currentIf := b.currentContainer.(*IfStmt)
	currentIf.Result = b.currentResultDef
	currentIf.nestingLevel = b.currentNestingLevel

	// setup the statement for ExitStmt
	b.currentStmt = currentIf

	// pop stack of containers
	b.currentContainer = currentIf.PreviousContainer
	b.currentNestingLevel--
}

// EnterElsePart is called when production elsePart is entered.
func (b *Builder) EnterElsePart(_ *ElsePartContext) {}

// ExitElsePart is called when production elsePart is exited.
func (b *Builder) ExitElsePart(_ *ElsePartContext) {
	currentIf := b.currentContainer.(*IfStmt)
	currentIf.ElsePart = []Stmt{}
}
