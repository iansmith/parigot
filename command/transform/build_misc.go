package transform

import "github.com/antlr/antlr4/runtime/Go/antlr"

// Builder is the "listener" that one uses to build the AST of the wat file
type Builder struct {
	*BaseWasmListener
	currentNestingLevel int
	module              []*Module
	file                string
	error               antlr.ErrorNode
}

func (b *Builder) Error() antlr.ErrorNode {
	return b.error
}

// Module returns the first Module in the list of parsed Modules()
func (b *Builder) Module() *Module {
	return b.module[0]
}

func (b *Builder) SetFile(f string) {
	b.file = f
}

func (b *Builder) File() string {
	return b.file
}
