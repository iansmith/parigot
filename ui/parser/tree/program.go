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
	if p.DocSection != nil && p.DocSection.Scope_ != nil && p.TextSection != nil && p.TextSection.Scope_ != nil {
		p.DocSection.Scope_.Brother = p.TextSection.Scope_
	}
}

func (p *ProgramNode) VarCheck(filename string) bool {
	if p.DocSection != nil {
		if !p.DocSection.VarCheck(filename) {
			return false
		}
	}
	if p.TextSection != nil {
		if !p.TextSection.VarCheck(filename) {
			return false
		}
	}

	return true
}

// Singleton and global with GProgramNode
func NewProgramNode() *ProgramNode {
	p := &ProgramNode{}
	return p
}
