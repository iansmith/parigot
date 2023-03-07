package parser

// IRNode is the parent of all node types, such as TextFuncNode.
type IRNode interface {
	Dump()
}

// Scope is things that can have variable decls and uses.
type Scope interface {
}

const maxScopeStackSize = 128
