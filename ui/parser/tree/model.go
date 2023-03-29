package tree

import (
	"fmt"
	"log"
	"strings"
)

type MVCSectionNode struct {
	Program   *ProgramNode
	ModelDecl []*ModelDecl
	ViewDecl  []*ViewDecl
}

func NewMvcSectionNode(p *ProgramNode) *MVCSectionNode {
	return &MVCSectionNode{Program: p}
}

type ModelDecl struct {
	Name    string
	Path    []string
	File    []*ProtobufFileNode
	Section *MVCSectionNode
}

func NewModelDecl() *ModelDecl {
	m := &ModelDecl{}
	GCurrentModel = m
	return m
}

var ErrMessageNotFound = fmt.Errorf("unable to find message")

func (m *ModelDecl) FinalizeSemantics() {
	for _, file := range m.File {
		for _, msg := range file.Message {
			//log.Printf("\t message %s", msg.Name)
			for _, field := range msg.Field {
				var mark, name string
				if field.Field != nil {
					if field.Field.TypeBase {
						continue
					}
					if !field.Field.TypeBase {
						pbnode, mesg, err := m.findMessageTypeByName(field.Field.Type)
						if err == ErrMessageNotFound {
							log.Printf("unable to find message named '%s' for field '%s' (in '%s')", field.Field.Type, field.Field.Name, file.Filename)
							return
						} else {
							loc := &TypeLocation{pbnode, mesg}
							msg.Location[field.Name] = loc
						}
					}
				} else {
					mark = ""
					name = "[map]" + field.Map.Name
					if !field.Map.ValueTypeBase {
						mark = field.Map.ValueType
					}
					log.Printf("\t\t%s field is '%s' '%s'-- %v ", msg.Name, name, mark, field.Map.ValueTypeBase)
				}
				//x := fmt.Sprintf("%s,%s", mark, name)
			}
		}
	}
}

// findMessageTypeByName returns the ProtobufFileNode and the ProtobufMessage associated with the name provided, if the name
// was found. If not, it returns MessageNotFound

func (m *ModelDecl) findMessageTypeByName(n string) (*ProtobufFileNode, *ProtobufMessage, error) {
	hasDot := strings.Contains(n, ".")
	for _, file := range m.File {
		//log.Printf("--- find message type by name '%s' in '%s',", n, file.PackageName)
		candidate := n
		if hasDot {
			part := strings.Split(candidate, ".")
			candidate = part[len(part)-1]
		}
		for _, msg := range file.Message {
			if hasDot {
				for _, imp := range file.Import {
					//log.Printf("xxx -- reached import %s (%s) from %s", imp.Filename, imp.PackageName, file.Filename)
					for _, message := range imp.Message {
						if message.Name == candidate {
							return imp, msg, nil
						}
					}
				}
			} else {
				// could be in this file
				for _, msg := range file.Message {
					//log.Printf("considering %s in file %s vs %s", msg.Name, file.Filename, n)
					if msg.Name == n {
						return file, msg, nil
					}
				}
			}
			return file, msg, nil
		}
	}
	return nil, nil, ErrMessageNotFound
}

type ViewDecl struct {
	ModelName string
	Section   *MVCSectionNode
	DocFn     *DocFuncNode
}

func NewViewDecl() *ViewDecl {
	return &ViewDecl{}
}
func (v *ViewDecl) CheckModelName() bool {
	foundModel := false
	modelType := []*PFormal{}
	for _, p := range v.DocFn.Param {
		if p.Type.HasStartColon {
			foundModel = true
			modelType = append(modelType, p)
			break
		}
	}
	if !foundModel {
		log.Printf("view function '%s' has no parameters which are a protobuf model (of the form ':model:message')", v.DocFn.Name)
		return false
	}
	for _, mt := range modelType {
		for _, modDecl := range v.Section.ModelDecl {
			if modDecl.Name == mt.Type.Part.Id {
				if mt.Type.Part.Qual == nil {
					log.Printf("view function '%s' has a parameter '%s' which is a model ('%s') but no message selected",
						v.DocFn.Name, mt.Name, modDecl.Name)
					return false
				}
				// we got a suffix
				q := mt.Type.Part.Qual
				if q.Qual != nil {
					log.Printf("view function '%s' has a parameter '%s' which is a model ('%s') but message name cannot have qualifier ('%s')",
						v.DocFn.Name, mt.Name, modDecl.Name, mt.Type.String())
					return false
				}
				for _, protobufNode := range modDecl.File {
					for _, candidate := range protobufNode.Message {
						if candidate.Name == q.Id {
							mt.Message = candidate
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func (v *ViewDecl) FinalizeSemantics() {
}

type ProtobufFileNode struct {
	PackageName string
	Filename    string
	GoPkg       string
	LocalGoPkg  string
	ImportFile  []string
	Import      []*ProtobufFileNode
	Message     []*ProtobufMessage
}

func NewProtobufFileNode() *ProtobufFileNode {
	return &ProtobufFileNode{}
}

type TypeLocation struct {
	File    *ProtobufFileNode
	Message *ProtobufMessage
}
type ProtobufMessage struct {
	Name     string
	Package  string
	Field    map[string]*ProtobufMsgElem
	Location map[string]*TypeLocation
}

func NewProtobufMessage(name string, allField map[string]*ProtobufMsgElem) *ProtobufMessage {
	return &ProtobufMessage{Name: name, Field: allField, Location: make(map[string]*TypeLocation)}
}

type FullIdent struct {
	Part []string
}

func NewFullIdent(c []string) *FullIdent {
	return &FullIdent{Part: c}
}

type OptionTriple struct {
	Name         *FullIdent
	Value, Extra string
}

func NewOptionTriple() *OptionTriple {
	return &OptionTriple{}
}

func (m *MVCSectionNode) VarCheck(filename string) bool {
	for _, v := range m.ViewDecl {
		if !v.CheckModelName() {
			return false
		}
	}
	m.MoveDocFns()
	return true
}
func (m *MVCSectionNode) MoveDocFns() {
	for _, view := range m.ViewDecl {
		m.Program.DocSection.AttachViewToSection(view)
	}
	m.ViewDecl = nil
}
func (m *MVCSectionNode) FinalizeSemantics() {
	if m == nil {
		return
	}
	if m.ModelDecl != nil {
		for _, mod := range m.ModelDecl {
			mod.FinalizeSemantics()
			mod.Section = m
		}

	}
	if m.ViewDecl != nil {
		for _, view := range m.ViewDecl {
			view.FinalizeSemantics()
			view.Section = m
		}
	}
}
