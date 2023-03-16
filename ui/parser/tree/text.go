package tree

import (
	"fmt"
)

// TextConstant is a simple string.
type TextConstant struct {
	_VarCtx *VarCtx
	Value   string
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

func NewTextConstant(s string) *TextConstant {
	return &TextConstant{_VarCtx: nil, Value: s}
}
func (t *TextConstant) SubTemplate() string {
	return "TextConstant"
}

// TextVar is a text variable that in source form is ${foo}
type TextVar struct {
	_VarCtx *VarCtx
	Name    string
}

func (t *TextVar) String() string {
	return fmt.Sprintf("${%s}", t.Name)
}

func (t *TextVar) Generate(_ *VarCtx) string {
	return fmt.Sprintf("LookupVar(%s)\n", t.Name)
}

func (t *TextVar) VarCtx() *VarCtx {
	return t._VarCtx
}

func NewTextVar(n string) *TextVar {
	return &TextVar{Name: n}
}
func (t *TextVar) SubTemplate() string {
	return "TextVar"
}

// TextInvoc is a text function call that in source form is ${foo(bar)}
type TextInvoc struct {
	Invoc *FuncInvoc
}

func (t *TextInvoc) String() string {
	return t.Invoc.String()
}

func (t *TextInvoc) Generate(_ *VarCtx) string {
	return fmt.Sprintf("%s\n", t.Invoc.String())
}

func (t *TextInvoc) VarCtx() *VarCtx {
	return nil
}

func NewTextInvoc(i *FuncInvoc) *TextInvoc {
	return &TextInvoc{Invoc: i}
}

func (t *TextInvoc) SubTemplate() string {
	return "TextInvoc"
}

// TextInline is a blob of code to copied into the output.
type TextInline struct {
	Name      string
	_VarCtx   *VarCtx
	TextItem_ []TextItem
}

func (t *TextInline) String() string {
	return "BLEAH NOT AVAILABLE"
}

func (t *TextInline) Generate(_ *VarCtx) string {
	return "BLEAH NOT AVAILABLE Generate"
}

func (t *TextInline) VarCtx() *VarCtx {
	return t._VarCtx
}

func NewTextInline() *TextInline {
	return &TextInline{}
}

func (t *TextInline) SubTemplate() string {
	return "TextInline"
}

// TextItem is an interface that represents the things that we
// know how to place inside a text unit.
type TextItem interface {
	//String() string
	VarCtx() *VarCtx
	SubTemplate() string
}

// TextExpander is something that can have variables uses in it.
type TextExpander interface {
	Item() []TextItem
}

// PFormal holds a parameter and type pair.
type PFormal struct {
	Name string
	Type *TypeDecl
}

func NewPFormal(n string, t *TypeDecl) *PFormal {
	return &PFormal{Name: n, Type: t}
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

func (t *TextFuncNode) CheckForBadVariableUse() string {
	for _, seq := range [][]TextItem{t.PreCode, t.PostCode, t.Item_} {
		for _, item := range seq {
			switch varName := item.(type) {
			case *TextVar:
				msg := t.checkAllForNameDecl(varName.Name)
				if msg != "" {
					return msg
				}
			}
		}
	}
	return ""
}

func (t *TextFuncNode) Item() []TextItem {
	return t.Item_
}

func (t *TextFuncNode) SetItem(item []TextItem) {
	t.Item_ = item
}

func (f *TextFuncNode) checkVar(name string, formal []*PFormal) bool {
	for _, p := range formal {
		if p.Name == name {
			return true
		}
	}
	return false
}

func (f *TextFuncNode) checkLocal(name string) bool {
	return f.checkVar(name, f.Local)
}
func (f *TextFuncNode) checkParam(name string) bool {
	return f.checkVar(name, f.Param)
}

func (f *TextFuncNode) checkGlobalAndExtern(name string) bool {
	return f.Section.Program.checkGlobalAndExtern(name)
}
func (f *TextFuncNode) checkAllForNameDecl(name string) string {
	if IsSelfVar(name) {
		return ""
	}
	found := f.checkLocal(name)
	if found {
		return ""
	}
	found = f.checkParam(name)
	if found {
		return ""
	}
	found = f.checkGlobalAndExtern(name)
	if found {
		return ""
	}
	return fmt.Sprintf("in text function '%s', use of unknown variable '%s'", f.Name, name)
}

func NewTextFuncNode() *TextFuncNode {
	return &TextFuncNode{}
}

// TestSection is the collection of text functions.
type TextSectionNode struct {
	Func    []*TextFuncNode
	Program *ProgramNode
}

func NewTextSectionNode() *TextSectionNode {
	return &TextSectionNode{}
}
func IsSelfVar(name string) bool {
	return name == "result"
}
