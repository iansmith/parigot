package tree

import (
	"fmt"
	"log"
)

const ValueRefTemplate = "TextValueRef"

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
func (t *TextConstant) GetLine() int {
	return t.LineNumber
}
func (t *TextConstant) GetCol() int {
	return t.ColumnNumber
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
	return ValueRefTemplate
}

func (t *TextValueRef) GetLine() int {
	return t.LineNumber
}
func (t *TextValueRef) GetCol() int {
	return t.ColumnNumber
}

// TextItem is an interface that represents the things that we
// know how to place inside a text unit.
type TextItem interface {
	String() string
	VarCtx() *VarCtx
	SubTemplate() string
	GetLine() int
	GetCol() int
}

// TypeName is a potential type starter and an ident
type TypeName struct {
	Type        *Ident
	TypeStarter string
}

func (t *TypeName) String() string {
	return t.TypeStarter + t.Type.Text
}

// PFormal holds a parameter and type pair.
type PFormal struct {
	Message                  *ProtobufMessage
	Name                     string
	TypeName                 *TypeName
	LineNumber, ColumnNumber int
}

func NewPFormal(n string, t *TypeName, ts string, line, col int) *PFormal {
	return &PFormal{Name: n, TypeName: t, LineNumber: line, ColumnNumber: col}
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
	LineNumber, ColumnNumber int
}

func (f *TextFuncNode) CheckDup(filename string) bool {
	if !checkDupParamAndLocal(f.Param, f.Local, filename, f.Name, false) {
		return false
	}
	if !checkParamShadown(f.Param, filename, f.Name, f.Section.Scope_, false) {
		return false
	}
	if !checkLocalShadow(f.Local, f.Param, filename, f.Name, f.Section.Scope_, true) {
		return false
	}

	return true
}

func (t *TextFuncNode) Item() []TextItem {
	return t.Item_
}

func (t *TextFuncNode) SetItem(item []TextItem) {
	t.Item_ = item
}

func (s *TextFuncNode) VarCheck(filename string) bool {
	if !CheckAllItems(s.Name, s.PreCode, s.Local, s.Param, s.Section.Scope_, filename) {
		return false
	}
	if !CheckAllItems(s.Name, s.PostCode, s.Local, s.Param, s.Section.Scope_, filename) {
		return false
	}
	return CheckAllItems(s.Name, s.Item_, s.Local, s.Param, s.Section.Scope_, filename)
}

func (s *TextSectionNode) FinalizeSemantics(path string) error {
	if s == nil {
		return nil //no section no error
	}
	for _, fn := range s.Func {
		fn.Section = s
	}
	s.Scope_.TextFn = s.Func

	// patch up the parameter types
	modelSect := s.Program.ModelSection
	for _, fn := range s.Func {
		for _, param := range fn.Param {
			if param.TypeName.Type.HasStartColon {
				_, msg, err := modelSect.ResolveModelMessageTypeForParam(path, param)
				if err != nil {
					return err
				}
				param.Message = msg
			}
		}
	}
	return nil
}

func NewTextFuncNode() *TextFuncNode {
	return &TextFuncNode{}
}

// TestSection is the collection of text functions.
type TextSectionNode struct {
	Func    []*TextFuncNode
	Program *ProgramNode
	Scope_  *SectionScope
}

func NewTextSectionNode(p *ProgramNode) *TextSectionNode {
	return &TextSectionNode{Program: p, Scope_: NewSectionScope(p.Global)}
}
func IsSelfVar(name string) bool {
	return name == "result"
}

func (t *TextSectionNode) VarCheck(filename string) bool {
	for _, fn := range t.Func {
		if !fn.VarCheck(filename) {
			return false
		}
		if !fn.CheckDup(filename) {
			return false
		}
		seen := make(map[string]*ErrorLoc)
		for _, fn := range t.Func {
			e := &ErrorLoc{filename, fn.LineNumber, fn.ColumnNumber}
			if _, ok := seen[fn.Name]; ok {
				log.Printf("two instances of text func %s found %s and %s", fn.Name, seen[fn.Name].String(), e.String())
				return false
			}
			seen[fn.Name] = e
		}

	}
	return true
}
