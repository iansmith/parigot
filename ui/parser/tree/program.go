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

	//NeedBytes, NeedElement, NeedEvent bool
}

func (p *ProgramNode) FinalizeSemantics() {
	for _, node := range []any{p.ImportSection, p.CSSSection, p.TextSection, p.DocSection,
		p.EventSection, p.ModelSection} {
		f, ok := node.(FSemantics)
		if ok {
			f.FinalizeSemantics()
		}
	}
}

// func (p *ProgramNode) checkGlobalAndExtern(n string) bool {
// 	if p.Global != nil && len(p.Global) > 0 {
// 		for _, g := range p.Global {
// 			if g.Name == n {
// 				return true
// 			}
// 		}
// 	}
// 	if p.Extern != nil && len(p.Extern) > 0 {
// 		for _, e := range p.Extern {
// 			if e == n {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

// Singleton and global with GProgramNode
func NewProgramNode() *ProgramNode {
	p := &ProgramNode{}
	return p
}
