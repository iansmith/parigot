package tree

import (
	"fmt"
	"strings"
)

type IdentPart struct {
	Id       string
	ColonSep bool
	Qual     *IdentPart
}

type Ident struct {
	HasStartColon            bool
	LineNumber, ColumnNumber int
	Part                     *IdentPart
	Text                     string
}

func (i *Ident) String() string {
	return i.Text
}

func NewIdent(start string, text string, line, col int) *Ident {
	i := &Ident{Text: text, LineNumber: line, ColumnNumber: col}
	i.Part = &IdentPart{Id: start}
	return i
}

func (i *Ident) IsWellFormed() string {
	if i.Part == nil {
		return fmt.Sprintf("did not find expected id, full text was '%s'", i.Text)
	}
	prev := i.Part
	if prev.ColonSep {
		for prev != nil {
			if prev.Qual != nil {
				return fmt.Sprintf("Badly formed id '%s', cannot mix qualifying identifier with Dot and Colon", i.Text)
			}
			prev = prev.Qual
		}
	} else {
		for prev != nil {
			if prev.ColonSep {
				return fmt.Sprintf("Badly formed id '%s', cannot mix qualifying identifier with Dot and Colon", i.Text)
			}
			prev = prev.Qual
		}
	}
	return ""
}

// Can be a value of a variable or a function call.  It's
// *almost* an expr.
type ValueRef struct {
	Id                       *Ident
	FuncInvoc                *FuncInvoc
	Lit                      string
	LineNumber, ColumnNumber int
}

func NewValueRef(i *Ident, f *FuncInvoc, lit string, line, col int) *ValueRef {
	if lit != "" {
		lit = strings.TrimPrefix(lit, "\"")
		lit = strings.TrimSuffix(lit, "\"")
	}
	return &ValueRef{Id: i, FuncInvoc: f, Lit: lit}
}

func (v *ValueRef) String() string {
	if v.Id != nil {
		return v.Id.String()
	}
	if v.Lit != "" {
		return fmt.Sprintf("\"%s\"", v.Lit)
	}
	return v.FuncInvoc.String()
}

// func (v *ValueRef) checkAllForNameDecl(varName string) {
// 	if v.Id != nil {
// 		v.Id.checkAllForNameDecl(varName)
// 	}
// 	v.FuncInvoc.checkAllForNameDecl(varName)
// }
