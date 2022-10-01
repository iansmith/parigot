package lib

// Builder is the "listener" that one uses to build the AST of the wat file
type Builder struct {
	*BaseWasmListener

	// current here means "it is being constructed if the value is not nil"
	currentModule   *Module
	currentTopLevel TopLevel
}

// EnterModule is called when entering the module production.
func (b *Builder) EnterModule(c *ModuleContext) {
	b.currentModule = &Module{}
}

// ExitModule is called when exiting the module production.
func (b *Builder) ExitModule(c *ModuleContext) {
	// post tree build should happen here
}
