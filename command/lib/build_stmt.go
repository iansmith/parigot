package lib

// EnterStmt is called when production stmt is entered.
func (b *Builder) EnterStmt(_ *StmtContext) {}

// ExitStmt is called when production stmt is exited.
func (b *Builder) ExitStmt(_ *StmtContext) {
	// this is null when the most recent statement IS the container
	if b.currentStmt != nil {
		b.currentContainer.AddStmt(b.currentStmt)
		b.currentStmt = nil
	}
	return
}

// EnterBlockStmt is called when production blockStmt is entered.
func (b *Builder) EnterBlockStmt(_ *BlockStmtContext) {
	b.currentNestingLevel++
	block := &BlockStmt{
		PreviousContainer: b.currentContainer,
		PreviousResult:    b.currentResultDef,
	}
	block.PreviousContainer.AddStmt(block)
	b.currentContainer = block
	b.currentResultDef = nil
}

// ExitBlockStmt is called when production block is exited.
func (b *Builder) ExitBlockStmt(_ *BlockStmtContext) {
	// this block statement
	currentBlock := b.currentContainer.(*BlockStmt)
	currentBlock.Result = b.currentResultDef
	currentBlock.nestingLevel = b.currentNestingLevel

	// pop stack of containers
	b.currentContainer = currentBlock.PreviousContainer
	b.currentResultDef = currentBlock.PreviousResult
	b.currentNestingLevel--
}

// EnterIfStmt is called when production ifStmt is entered.
func (b *Builder) EnterIfStmt(_ *IfStmtContext) {
	b.currentNestingLevel++
	ifStmt := &IfStmt{
		PreviousContainer: b.currentContainer,
		PreviousResult:    b.currentResultDef,
	}
	ifStmt.PreviousContainer.AddStmt(ifStmt)
	b.currentContainer = ifStmt
	b.currentResultDef = nil
}

// ExitIfStmt is called when production ifStmt is exited.
func (b *Builder) ExitIfStmt(_ *IfStmtContext) {
	// this block statement
	currentIf := b.currentContainer.(*IfStmt)
	currentIf.Result = b.currentResultDef
	currentIf.nestingLevel = b.currentNestingLevel

	// pop stack of containers
	b.currentContainer = currentIf.PreviousContainer
	b.currentResultDef = currentIf.PreviousResult
	b.currentNestingLevel--
}

// EnterElsePart is called when production elsePart is entered.
func (b *Builder) EnterElsePart(_ *ElsePartContext) {}

// ExitElsePart is called when production elsePart is exited.
func (b *Builder) ExitElsePart(_ *ElsePartContext) {
	currentIf := b.currentContainer.(*IfStmt)
	currentIf.ElsePart = []Stmt{}
}

// EnterLoopStmt is called when production loopStmt is entered.
func (b *Builder) EnterLoopStmt(_ *LoopStmtContext) {
	b.currentNestingLevel++
	loop := &LoopStmt{
		BlockStmt: &BlockStmt{
			PreviousContainer: b.currentContainer,
			PreviousResult:    b.currentResultDef,
		},
	}
	loop.PreviousContainer.AddStmt(loop)
	b.currentContainer = loop
	b.currentResultDef = nil
}

// ExitLoopStmt is called when production loopStmt is exited.
func (b *Builder) ExitLoopStmt(_ *LoopStmtContext) {
	// this loop statement
	currentLoop := b.currentContainer.(*LoopStmt)
	currentLoop.Result = b.currentResultDef
	currentLoop.nestingLevel = b.currentNestingLevel

	// pop stack of containers
	b.currentContainer = currentLoop.PreviousContainer
	b.currentResultDef = currentLoop.PreviousResult
	b.currentNestingLevel--
}
