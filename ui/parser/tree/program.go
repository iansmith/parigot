package tree

import "log"

type FSemantics interface {
	FinalizeSemantics(path string) error
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

func (p *ProgramNode) FinalizeSemantics(path string) error {
	// ordering of this list matters... models have to be fully put together before dealing with parameters in the text and doc sections
	for _, node := range []any{p.ImportSection, p.CSSSection, p.ModelSection, p.TextSection, p.DocSection,
		p.EventSection} {
		if node == nil {
			log.Printf("nil found in FS list")
			continue
		}
		f, ok := node.(FSemantics)
		if ok && node != nil {
			if err := f.FinalizeSemantics(path); err != nil {
				return err
			}
		}
	}
	if p.DocSection != nil && p.DocSection.Scope_ != nil && p.TextSection != nil && p.TextSection.Scope_ != nil {
		p.DocSection.Scope_.Brother = p.TextSection.Scope_
	}

	if p.DocSection != nil {
		p.DocSection.SetNumber()
	}

	return nil
}

func (p *ProgramNode) VarCheck(filename string) bool {
	// ordering of the section checks matterns because
	// model moves things to doc section and doc section
	// moves things to text section

	if p.ModelSection != nil {
		if !p.ModelSection.VarCheck(filename) {
			return false
		}
	}

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
