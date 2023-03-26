package tree

type Scope interface {
	Parent() Scope
	LookupFunc(*FuncInvoc) bool
	LookupVar(*Ident) *PFormal
}

// these globals are so code that is running DURING the
// construction of these objects can get a reference to them.
var GProgram *ProgramNode = &ProgramNode{}
var GCurrentModel *ModelDecl

type GlobalSectionNode struct {
	Program      *ProgramNode
	LineNumber   int
	ColumnNumber int
	Var          []*PFormal
}

func (g *GlobalSectionNode) LookupVar(id *Ident) *PFormal {
	for _, v := range g.Var {
		if v.Name == id.Part.Id {
			return v
		}
	}
	return nil
}

func NewGlobalSectionNode(p *ProgramNode, ln, col int) *GlobalSectionNode {
	return &GlobalSectionNode{Program: p, LineNumber: ln, ColumnNumber: col}
}

// Parameters of an extern function are not checked or touched.
type ExternSectionNode struct {
	Program                  *ProgramNode
	Name                     []string
	LineNumber, ColumnNumber int
}

func (e *ExternSectionNode) LookupFunc(f *FuncInvoc) bool {
	for _, n := range e.Name {
		if n == f.String() {
			return true //no check of actuals
		}
	}
	return false
}

func NewExternSectionNode(p *ProgramNode, ln, col int) *ExternSectionNode {
	return &ExternSectionNode{Program: p, LineNumber: ln, ColumnNumber: col}
}

// AllGlobal is a scope that is combination of the two "global" sections
type AllGlobal struct {
	G *GlobalSectionNode
	E *ExternSectionNode
}

func NewAllGlobal(g *GlobalSectionNode, e *ExternSectionNode) *AllGlobal {
	return &AllGlobal{G: g, E: e}
}

func (a *AllGlobal) Parent() Scope {
	return nil
}

func (a *AllGlobal) LookupFunc(f *FuncInvoc) bool {
	if a == nil || a.E == nil {
		return false
	}
	return a.E.LookupFunc(f)
}

func (a *AllGlobal) LookupVar(id *Ident) *PFormal {
	if a == nil || a.G == nil {
		return nil
	}
	return a.G.LookupVar(id)
}

type SectionScope struct {
	TextFn  []*TextFuncNode
	DocFn   []*DocFuncNode
	Parent_ Scope
}

func NewSectionScope(a *AllGlobal) *SectionScope {
	return &SectionScope{Parent_: a}
}

func (s *SectionScope) Parent() Scope {
	return s.Parent_
}

func (s *SectionScope) LookupFunc(f *FuncInvoc) bool {
	if s.TextFn != nil {
		for _, fn := range s.TextFn {
			if fn.Name == f.Name.String() {
				return true
			}
		}
	} else if s.DocFn != nil {
		for _, fn := range s.TextFn {
			if fn.Name == f.Name.String() {
				return true
			}
		}
	}
	return s.Parent().LookupFunc(f)
}

func (s *SectionScope) LookupVar(id *Ident) *PFormal {
	return s.Parent().LookupVar(id)
}

type FuncScope struct {
	Parent_ Scope
	Formal  []*PFormal
	Local   []*PFormal
}

func NewFuncScope(p Scope) *FuncScope {
	return &FuncScope{Parent_: p}
}

func (f *FuncScope) Parent() Scope {
	return f.Parent_
}
