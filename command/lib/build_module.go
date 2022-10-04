package lib

// EnterModule is called when entering the module production.
func (b *Builder) EnterModule(c *ModuleContext) {
}

// ExitModule is called when exiting the module production.
func (b *Builder) ExitModule(c *ModuleContext) {
	b.module = append(b.module, c.GetM())
}
