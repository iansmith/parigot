package tree

type FSemantics interface {
	FinalizeSemantics()
}

type ProgramNode struct {
	ImportSection *ImportSectionNode
	CSSSection    *CSSSectionNode
	TextSection   *TextSectionNode
	DocSection    *DocSectionNode
	EventSection  *EventSectionNode
	ModelSection  *MVCSectionNode
	Global        *AllGlobal
}

func (p *ProgramNode) FinalizeSemantics() {
	for _, node := range []any{p.ImportSection, p.CSSSection, p.TextSection, p.DocSection,
		p.EventSection, p.ModelSection} {
		if node == nil {
			continue
		}
		f, ok := node.(FSemantics)
		if ok && node != nil {
			f.FinalizeSemantics()
		}
	}
}

func (p *ProgramNode) VarCheck(filename string) {
	if p.DocSection != nil {
		p.DocSection.VarCheck(filename)
	}
	if p.TextSection != nil {
		p.TextSection.VarCheck(filename)
	}
}

// Singleton and global with GProgramNode
func NewProgramNode() *ProgramNode {
	p := &ProgramNode{}
	return p
}
