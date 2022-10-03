package lib

// Builder is the "listener" that one uses to build the AST of the wat file
type Builder struct {
	*BaseWasmListener
	currentModule       *Module
	currentNestingLevel int
	module              []*Module
}

// Module returns the first Module in the list of parsed Modules()
func (b *Builder) Module() *Module {
	return b.module[0]
}
