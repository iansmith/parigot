package parser

type ProgramNode struct {
	ImportSection          *ImportSectionNode
	CSSSection             *CSSSectionNode
	TextSection            *TextSectionNode
	DocSection             *DocSectionNode
	Extern                 []string
	Global                 []*PFormal
	NeedBytes, NeedElement bool
}

func (p *ProgramNode) checkGlobalAndExtern(n string) bool {
	if p.Global != nil && len(p.Global) > 0 {
		for _, g := range p.Global {
			if g.Name == n {
				return true
			}
		}
	}
	if p.Extern != nil && len(p.Extern) > 0 {
		for _, e := range p.Extern {
			if e == n {
				return true
			}
		}
	}
	return false
}

func NewProgramNode() *ProgramNode {
	return &ProgramNode{}
}
