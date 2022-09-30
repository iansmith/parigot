// Code generated from command/jsstrip/jsstrip.g4 by ANTLR 4.9. DO NOT EDIT.

package main // jsstrip
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BasejsstripListener is a complete listener for a parse tree produced by jsstripParser.
type BasejsstripListener struct{}

var _ jsstripListener = &BasejsstripListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasejsstripListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasejsstripListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasejsstripListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasejsstripListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterModule is called when production module is entered.
func (s *BasejsstripListener) EnterModule(ctx *ModuleContext) {}

// ExitModule is called when production module is exited.
func (s *BasejsstripListener) ExitModule(ctx *ModuleContext) {}

// EnterTopLevel is called when production topLevel is entered.
func (s *BasejsstripListener) EnterTopLevel(ctx *TopLevelContext) {}

// ExitTopLevel is called when production topLevel is exited.
func (s *BasejsstripListener) ExitTopLevel(ctx *TopLevelContext) {}

// EnterTypeDef is called when production typeDef is entered.
func (s *BasejsstripListener) EnterTypeDef(ctx *TypeDefContext) {}

// ExitTypeDef is called when production typeDef is exited.
func (s *BasejsstripListener) ExitTypeDef(ctx *TypeDefContext) {}

// EnterImportOp is called when production importOp is entered.
func (s *BasejsstripListener) EnterImportOp(ctx *ImportOpContext) {}

// ExitImportOp is called when production importOp is exited.
func (s *BasejsstripListener) ExitImportOp(ctx *ImportOpContext) {}

// EnterTypeRef is called when production typeRef is entered.
func (s *BasejsstripListener) EnterTypeRef(ctx *TypeRefContext) {}

// ExitTypeRef is called when production typeRef is exited.
func (s *BasejsstripListener) ExitTypeRef(ctx *TypeRefContext) {}

// EnterTypeAnnotation is called when production typeAnnotation is entered.
func (s *BasejsstripListener) EnterTypeAnnotation(ctx *TypeAnnotationContext) {}

// ExitTypeAnnotation is called when production typeAnnotation is exited.
func (s *BasejsstripListener) ExitTypeAnnotation(ctx *TypeAnnotationContext) {}

// EnterFuncSpec is called when production funcSpec is entered.
func (s *BasejsstripListener) EnterFuncSpec(ctx *FuncSpecContext) {}

// ExitFuncSpec is called when production funcSpec is exited.
func (s *BasejsstripListener) ExitFuncSpec(ctx *FuncSpecContext) {}

// EnterFuncRef is called when production funcRef is entered.
func (s *BasejsstripListener) EnterFuncRef(ctx *FuncRefContext) {}

// ExitFuncRef is called when production funcRef is exited.
func (s *BasejsstripListener) ExitFuncRef(ctx *FuncRefContext) {}

// EnterType_ is called when production type_ is entered.
func (s *BasejsstripListener) EnterType_(ctx *Type_Context) {}

// ExitType_ is called when production type_ is exited.
func (s *BasejsstripListener) ExitType_(ctx *Type_Context) {}

// EnterTypeSeq is called when production typeSeq is entered.
func (s *BasejsstripListener) EnterTypeSeq(ctx *TypeSeqContext) {}

// ExitTypeSeq is called when production typeSeq is exited.
func (s *BasejsstripListener) ExitTypeSeq(ctx *TypeSeqContext) {}

// EnterParamDef is called when production paramDef is entered.
func (s *BasejsstripListener) EnterParamDef(ctx *ParamDefContext) {}

// ExitParamDef is called when production paramDef is exited.
func (s *BasejsstripListener) ExitParamDef(ctx *ParamDefContext) {}

// EnterResultDef is called when production resultDef is entered.
func (s *BasejsstripListener) EnterResultDef(ctx *ResultDefContext) {}

// ExitResultDef is called when production resultDef is exited.
func (s *BasejsstripListener) ExitResultDef(ctx *ResultDefContext) {}

// EnterLocalDef is called when production localDef is entered.
func (s *BasejsstripListener) EnterLocalDef(ctx *LocalDefContext) {}

// ExitLocalDef is called when production localDef is exited.
func (s *BasejsstripListener) ExitLocalDef(ctx *LocalDefContext) {}

// EnterFuncDef is called when production funcDef is entered.
func (s *BasejsstripListener) EnterFuncDef(ctx *FuncDefContext) {}

// ExitFuncDef is called when production funcDef is exited.
func (s *BasejsstripListener) ExitFuncDef(ctx *FuncDefContext) {}

// EnterFuncBody is called when production funcBody is entered.
func (s *BasejsstripListener) EnterFuncBody(ctx *FuncBodyContext) {}

// ExitFuncBody is called when production funcBody is exited.
func (s *BasejsstripListener) ExitFuncBody(ctx *FuncBodyContext) {}

// EnterStmt is called when production stmt is entered.
func (s *BasejsstripListener) EnterStmt(ctx *StmtContext) {}

// ExitStmt is called when production stmt is exited.
func (s *BasejsstripListener) ExitStmt(ctx *StmtContext) {}

// EnterBlock is called when production block is entered.
func (s *BasejsstripListener) EnterBlock(ctx *BlockContext) {}

// ExitBlock is called when production block is exited.
func (s *BasejsstripListener) ExitBlock(ctx *BlockContext) {}

// EnterSexpr is called when production sexpr is entered.
func (s *BasejsstripListener) EnterSexpr(ctx *SexprContext) {}

// ExitSexpr is called when production sexpr is exited.
func (s *BasejsstripListener) ExitSexpr(ctx *SexprContext) {}

// EnterList is called when production list is entered.
func (s *BasejsstripListener) EnterList(ctx *ListContext) {}

// ExitList is called when production list is exited.
func (s *BasejsstripListener) ExitList(ctx *ListContext) {}

// EnterMembers is called when production members is entered.
func (s *BasejsstripListener) EnterMembers(ctx *MembersContext) {}

// ExitMembers is called when production members is exited.
func (s *BasejsstripListener) ExitMembers(ctx *MembersContext) {}

// EnterAtom is called when production atom is entered.
func (s *BasejsstripListener) EnterAtom(ctx *AtomContext) {}

// ExitAtom is called when production atom is exited.
func (s *BasejsstripListener) ExitAtom(ctx *AtomContext) {}
