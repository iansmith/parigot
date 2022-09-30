// Code generated from command/jsstrip/jsstrip.g4 by ANTLR 4.9. DO NOT EDIT.

package main // jsstrip
import "github.com/antlr/antlr4/runtime/Go/antlr"

// jsstripListener is a complete listener for a parse tree produced by jsstripParser.
type jsstripListener interface {
	antlr.ParseTreeListener

	// EnterModule is called when entering the module production.
	EnterModule(c *ModuleContext)

	// EnterTopLevel is called when entering the topLevel production.
	EnterTopLevel(c *TopLevelContext)

	// EnterTypeDef is called when entering the typeDef production.
	EnterTypeDef(c *TypeDefContext)

	// EnterImportOp is called when entering the importOp production.
	EnterImportOp(c *ImportOpContext)

	// EnterTypeRef is called when entering the typeRef production.
	EnterTypeRef(c *TypeRefContext)

	// EnterTypeAnnotation is called when entering the typeAnnotation production.
	EnterTypeAnnotation(c *TypeAnnotationContext)

	// EnterFuncSpec is called when entering the funcSpec production.
	EnterFuncSpec(c *FuncSpecContext)

	// EnterFuncRef is called when entering the funcRef production.
	EnterFuncRef(c *FuncRefContext)

	// EnterType_ is called when entering the type_ production.
	EnterType_(c *Type_Context)

	// EnterTypeSeq is called when entering the typeSeq production.
	EnterTypeSeq(c *TypeSeqContext)

	// EnterParamDef is called when entering the paramDef production.
	EnterParamDef(c *ParamDefContext)

	// EnterResultDef is called when entering the resultDef production.
	EnterResultDef(c *ResultDefContext)

	// EnterLocalDef is called when entering the localDef production.
	EnterLocalDef(c *LocalDefContext)

	// EnterFuncDef is called when entering the funcDef production.
	EnterFuncDef(c *FuncDefContext)

	// EnterFuncBody is called when entering the funcBody production.
	EnterFuncBody(c *FuncBodyContext)

	// EnterStmt is called when entering the stmt production.
	EnterStmt(c *StmtContext)

	// EnterBlock is called when entering the block production.
	EnterBlock(c *BlockContext)

	// EnterSexpr is called when entering the sexpr production.
	EnterSexpr(c *SexprContext)

	// EnterList is called when entering the list production.
	EnterList(c *ListContext)

	// EnterMembers is called when entering the members production.
	EnterMembers(c *MembersContext)

	// EnterAtom is called when entering the atom production.
	EnterAtom(c *AtomContext)

	// ExitModule is called when exiting the module production.
	ExitModule(c *ModuleContext)

	// ExitTopLevel is called when exiting the topLevel production.
	ExitTopLevel(c *TopLevelContext)

	// ExitTypeDef is called when exiting the typeDef production.
	ExitTypeDef(c *TypeDefContext)

	// ExitImportOp is called when exiting the importOp production.
	ExitImportOp(c *ImportOpContext)

	// ExitTypeRef is called when exiting the typeRef production.
	ExitTypeRef(c *TypeRefContext)

	// ExitTypeAnnotation is called when exiting the typeAnnotation production.
	ExitTypeAnnotation(c *TypeAnnotationContext)

	// ExitFuncSpec is called when exiting the funcSpec production.
	ExitFuncSpec(c *FuncSpecContext)

	// ExitFuncRef is called when exiting the funcRef production.
	ExitFuncRef(c *FuncRefContext)

	// ExitType_ is called when exiting the type_ production.
	ExitType_(c *Type_Context)

	// ExitTypeSeq is called when exiting the typeSeq production.
	ExitTypeSeq(c *TypeSeqContext)

	// ExitParamDef is called when exiting the paramDef production.
	ExitParamDef(c *ParamDefContext)

	// ExitResultDef is called when exiting the resultDef production.
	ExitResultDef(c *ResultDefContext)

	// ExitLocalDef is called when exiting the localDef production.
	ExitLocalDef(c *LocalDefContext)

	// ExitFuncDef is called when exiting the funcDef production.
	ExitFuncDef(c *FuncDefContext)

	// ExitFuncBody is called when exiting the funcBody production.
	ExitFuncBody(c *FuncBodyContext)

	// ExitStmt is called when exiting the stmt production.
	ExitStmt(c *StmtContext)

	// ExitBlock is called when exiting the block production.
	ExitBlock(c *BlockContext)

	// ExitSexpr is called when exiting the sexpr production.
	ExitSexpr(c *SexprContext)

	// ExitList is called when exiting the list production.
	ExitList(c *ListContext)

	// ExitMembers is called when exiting the members production.
	ExitMembers(c *MembersContext)

	// ExitAtom is called when exiting the atom production.
	ExitAtom(c *AtomContext)
}
