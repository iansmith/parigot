package lib

// Builder is the "listener" that one uses to build the AST of the wat file
type Builder struct {
	*BaseWasmListener

	// current here means "it is being constructed if the value is not nil"
	currentModule *Module

	currentTopLevelDef TopLevelDef

	currentTypeRef     *TypeRef
	currentParamDef    *ParamDef
	currentLocalDef    *LocalDef
	currentFuncNameRef *FuncNameRef
	currentFuncSpec    *FuncSpec
	currentResultDef   *ResultDef
	currentTypeNameSeq *TypeNameSeq

	currentTypeAnnotation *TypeAnnotation

	module []*Module
}

// Module returns the first Module in the list of parsed Modules()
func (b *Builder) Module() *Module {
	return b.module[0]
}
