package lib

// EnterModule is called when entering the module production.
func (b *Builder) EnterModule(c *ModuleContext) {
	if b.currentModule != nil {
		panic("Module already under construction")
	}
	b.currentModule = &Module{}
}

// ExitModule is called when exiting the module production.
func (b *Builder) ExitModule(c *ModuleContext) {
	b.module = append(b.module, b.currentModule)
	// post tree build should happen here
	b.currentModule = nil
}
