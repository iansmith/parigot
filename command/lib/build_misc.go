package lib

// Builder is the "listener" that one uses to build the AST of the wat file
type Builder struct {
	*BaseWasmListener

	// current here means "it is being constructed if the value is not nil"
	currentModule *Module

	currentTopLevelDef TopLevelDef

	currentParamDef    *ParamDef
	currentLocalDef    *LocalDef
	currentFuncNameRef *FuncNameRef
	currentFuncSpec    *FuncSpec
	currentResultDef   *ResultDef
	currentTypeNameSeq *TypeNameSeq

	currentTypeAnnotation *TypeAnnotation

	currentContainer    Container
	currentStmt         Stmt // ops are also stmts
	currentNestingLevel int

	module []*Module
}

// Module returns the first Module in the list of parsed Modules()
func (b *Builder) Module() *Module {
	return b.module[0]
}
