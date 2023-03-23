package tree

import (
	"fmt"
)

// ////////////////////
// TextConstant is a simple string.
type TextConstant struct {
	_VarCtx                  *VarCtx
	Value                    string
	LineNumber, ColumnNumber int
}

func (t *TextConstant) String() string {
	return t.Value
}

func (t *TextConstant) Generate(_ *VarCtx) string {
	value := t.Value
	return fmt.Sprintf("buf WriteString(%q)\n", value)
}

func (t *TextConstant) VarCtx() *VarCtx {
	return t._VarCtx
}

func NewTextConstant(s string, ln, col int) *TextConstant {
	if s == "\\>" {
		s = ">"
	}
	return &TextConstant{_VarCtx: nil, Value: s, LineNumber: ln, ColumnNumber: col}
}
func (t *TextConstant) SubTemplate() string {
	return "TextConstant"
}

// ////////////////////
// TextValueRef is a  reference to a variable or a function call
type TextValueRef struct {
	_VarCtx                  *VarCtx
	Ref                      *ValueRef
	LineNumber, ColumnNumber int
}

func (t *TextValueRef) String() string {
	return t.Ref.String()
}

func (t *TextValueRef) Generate(_ *VarCtx) string {
	return fmt.Sprintf("LookupVar(%s)\n", "NOT USED")
}

func (t *TextValueRef) VarCtx() *VarCtx {
	return t._VarCtx
}

func NewTextValueRef(vr *ValueRef, ln, col int) *TextValueRef {
	return &TextValueRef{Ref: vr, LineNumber: ln, ColumnNumber: col}
}
func (t *TextValueRef) SubTemplate() string {
	return "TextValueRef"
}

// ////////////////////
// TextInline is a blob of code to copied into the output.
// type TextInline struct {
// 	Name                     string
// 	_VarCtx                  *VarCtx
// 	TextItem_                []TextItem
// 	LineNumber, ColumnNumber int
// }

// func (t *TextInline) String() string {
// 	return "BLEAH NOT AVAILABLE"
// }

// func (t *TextInline) Generate(_ *VarCtx) string {
// 	return "BLEAH NOT AVAILABLE Generate"
// }

// func (t *TextInline) VarCtx() *VarCtx {
// 	return t._VarCtx
// }

// func NewTextInline() *TextInline {
// 	return &TextInline{}
// }

// func (t *TextInline) SubTemplate() string {
// 	return "TextInline"
// }

// TextItem is an interface that represents the things that we
// know how to place inside a text unit.
type TextItem interface {
	String() string
	VarCtx() *VarCtx
	SubTemplate() string
}

// PFormal holds a parameter and type pair.
type PFormal struct {
	Name        string
	Type        *Ident
	TypeStarter string
}

func NewPFormal(n string, t *Ident, ts string) *PFormal {
	return &PFormal{Name: n, Type: t, TypeStarter: ts}
}

// Either Simple is set or both ModelName and ModelMessage are set
type TypeDecl struct {
	Simple       string
	ModelName    string
	ModelMessage string
}

func (t *TypeDecl) String() string {
	if t.Simple != "" {
		return t.Simple
	}
	return fmt.Sprintf("%s_%s", t.ModelName, t.ModelMessage)
}
func NewTypeDeclSimple(s string) *TypeDecl {
	return &TypeDecl{Simple: s}
}

func NewTypeDeclModel(model, message string) *TypeDecl {
	return &TypeDecl{ModelName: model, ModelMessage: message}
}

// TextFuncNode is the that alls the information about a declared
// text function.
type TextFuncNode struct {
	Name                     string
	NumParams                int
	Param                    []*PFormal
	Local                    []*PFormal
	Item_, PreCode, PostCode []TextItem
	Section                  *TextSectionNode
}

// func (t *TextFuncNode) CheckForBadVariableUse() string {
// 	for _, seq := range [][]TextItem{t.PreCode, t.PostCode, t.Item_} {
// 		for _, item := range seq {
// 			switch varName := item.(type) {
// 			case *TextValueRef:
// 				msg := varName.checkAllForNameDecl(varName)
// 				if msg != "" {
// 					return msg
// 				}
// 			}
// 		}
// 	}
// 	return ""
// }

func (t *TextFuncNode) Item() []TextItem {
	return t.Item_
}

func (t *TextFuncNode) SetItem(item []TextItem) {
	t.Item_ = item
}

// func (f *TextFuncNode) checkVar(name string, formal []*PFormal) bool {
// 	for _, p := range formal {
// 		if p.Name == name {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (f *TextFuncNode) checkLocal(name string) bool {
// 	return f.checkVar(name, f.Local)
// }
// func (f *TextFuncNode) checkParam(name string) bool {
// 	return f.checkVar(name, f.Param)
// }

// func (f *TextFuncNode) checkGlobalAndExtern(name string) bool {
// 	return f.Section.Program.checkGlobalAndExtern(name)
// }
// func (f *TextFuncNode) checkAllForNameDecl(name string) string {
// 	if IsSelfVar(name) {
// 		return ""
// 	}
// 	found := f.checkLocal(name)
// 	if found {
// 		return ""
// 	}
// 	found = f.checkParam(name)
// 	if found {
// 		return ""
// 	}
// 	found = f.checkGlobalAndExtern(name)
// 	if found {
// 		return ""
// 	}
// 	return fmt.Sprintf("in text function '%s', use of unknown variable '%s'", f.Name, name)
// }

func NewTextFuncNode() *TextFuncNode {
	return &TextFuncNode{}
}

// TestSection is the collection of text functions.
type TextSectionNode struct {
	Func         []*TextFuncNode
	Program      *ProgramNode
	SectionScope *SectionScope
}

func NewTextSectionNode(p *ProgramNode) *TextSectionNode {
	return &TextSectionNode{Program: p, SectionScope: NewSectionScope(p.Global)}
}
func IsSelfVar(name string) bool {
	return name == "result"
}
