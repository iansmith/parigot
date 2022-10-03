package lib

// EnterBlockStmt is called when production blockStmt is entered.
func (b *Builder) EnterBlockStmt(_ *BlockStmtContext) {
	b.currentNestingLevel++
}

// ExitBlockStmt is called when production block is exited.
func (b *Builder) ExitBlockStmt(bl *BlockStmtContext) {
	bl.GetB().(*BlockStmt).nestingLevel = b.currentNestingLevel
	b.currentNestingLevel--
}

// EnterIfStmt is called when production ifStmt is entered.
func (b *Builder) EnterIfStmt(_ *IfStmtContext) {
	b.currentNestingLevel++
}

// ExitIfStmt is called when production ifStmt is exited.
func (b *Builder) ExitIfStmt(i *IfStmtContext) {
	i.GetI().(*IfStmt).nestingLevel = b.currentNestingLevel
	b.currentNestingLevel--
}

// EnterLoopStmt is called when production loopStmt is entered.
func (b *Builder) EnterLoopStmt(_ *LoopStmtContext) {
	b.currentNestingLevel++
}

// ExitLoopStmt is called when production loopStmt is exited.
func (b *Builder) ExitLoopStmt(l *LoopStmtContext) {
	l.GetL().(*LoopStmt).nestingLevel = b.currentNestingLevel
	b.currentNestingLevel--
}
