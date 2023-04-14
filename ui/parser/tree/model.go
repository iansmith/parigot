package tree

import (
	"fmt"
	"log"
	"strings"
)

type MVCSectionNode struct {
	Program        *ProgramNode
	ModelDecl      []*ModelDecl
	ViewDecl       []*ViewDecl
	ControllerDecl []*ControllerDecl
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

func (m *ModelDecl) FinalizeSemantics(filename string) error {

	// this builds a map in each *message* that knows how to resolve each of message's fields to a type (which can be another
	// message)
	for _, file := range m.File {
		for _, msg := range file.Message {
			//log.Printf("\t message %s", msg.Name)
			for _, field := range msg.Field {
				if field.Field != nil {
					if field.Field.TypeBase {
						continue
					}
					pbnode, mesg, err := m.FindMessageTypeByName(field.Field.Type)
					if err == ErrMessageNotFound {
						e := ErrorLoc{Filename: filename, Line: 0, Col: 0} // xxx fixme
						s := fmt.Sprintf("unable to find message named '%s' for field '%s' (in '%s') at %s", field.Field.Type, field.Field.Name, file.Filename, e.String())
						return fmt.Errorf("%s", s)
					} else {
						loc := &TypeLocation{pbnode, mesg}
						msg.Location[field.Name] = loc
					}
				} else {
					var mark, name string
					mark = ""
					name = "[map]" + field.Map.Name
					if !field.Map.ValueTypeBase {
						mark = field.Map.ValueType
					}
					//log.Printf("\t\t%s field is '%s' '%s'-- %v ", msg.Name, name, mark, field.Map.ValueTypeBase)
					_ = fmt.Sprintf("%s,%s", mark, name)
				}
			}
		}
	}
	return nil
}

// ResolveModelMessageByIdent returns either the place where the type is defined (file and message) or a NotFound. It walks all the
// defined models.
func (s *MVCSectionNode) ResolveModelMessageTypeByIdent(filename string, ident *Ident) (*ProtobufFileNode, *ProtobufMessage, error) {
	if !strings.Contains(ident.String(), ":") {
		panic(fmt.Sprintf("unable to understand formal type referenc (no qualifiers) '%s'", ident.String()))
	}
	part := strings.Split(ident.String(), ":")
	if len(part) != 3 {
		panic(fmt.Sprintf("unable to understand model message reference (too many qualifiers) '%s'", ident.String()))
	}
	part = part[1:]
	var found *ModelDecl
	for _, model := range s.ModelDecl {
		//log.Printf("xxxx comparing model.Name %s to %s", model.Name, part[0])
		if model.Name == part[0] {
			found = model
			break
		}
	}
	if found == nil {
		e := ErrorLoc{
			Filename: filename,
			Line:     ident.LineNumber,
			Col:      ident.ColumnNumber,
		}
		return nil, nil, fmt.Errorf("unable find a model message named '%s' at %s", part[0], e.String())
	}
	//continue looking for the message
	for _, f := range found.File {
		for _, msg := range f.Message {
			if msg.Name == part[1] {
				return f, msg, nil
			}
		}
	}
	return nil, nil, ErrMessageNotFound

}

// ResolveModelMessageTypeForParam returns either the place where the type is defined (file and message) or a NotFound. It walks all the
// known models.
func (s *MVCSectionNode) ResolveModelMessageTypeForParam(filename string, param *PFormal) (*ProtobufFileNode, *ProtobufMessage, error) {
	return s.ResolveModelMessageTypeByIdent(filename, param.TypeName.Type)
}

// findMessageTypeByName returns the ProtobufFileNode and the ProtobufMessage associated with the name provided, if the name
// was found. If not, it returns MessageNotFound
func (m *ModelDecl) FindMessageTypeByName(n string) (*ProtobufFileNode, *ProtobufMessage, error) {
	hasDot := strings.Contains(n, ".")
	for _, file := range m.File {
		//log.Printf("--- find message type by name '%s' in '%s  with file imports %+v", n, file.PackageName, file.ImportFile)
		candidate := n
		if hasDot {
			part := strings.Split(candidate, ".")
			candidate = part[len(part)-1]
		}
		for _, msg := range file.Message {
			if hasDot {
				for _, imp := range file.Import {
					for _, message := range imp.Message {
						if message.Name == candidate {
							//log.Printf("got a hit on message %s,%s", imp.Filename, msg.Name)
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
			//return file, msg, nil
		}
	}
	for googleName, msg := range googleTypes {
		if googleName == n {
			return nil, msg, nil
		}
	}
	return nil, nil, ErrMessageNotFound
}

var googleTypes = map[string]*ProtobufMessage{
	"google.protobuf.Any":       {Name: "Any", Package: "google.protobuf", AnyField: true},
	"google.protobuf.Timestamp": {Name: "Timestamp", Package: "google.protobuf", AnyField: true},
}

type ViewDecl struct {
	ModelName *Ident
	Message   *ProtobufMessage
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
		if p.TypeName.Type.HasStartColon {
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
			if modDecl.Name == mt.TypeName.Type.Part.Id {
				if mt.TypeName.Type.Part.Qual == nil {
					log.Printf("view function '%s' has a parameter '%s' which is a model ('%s') but no message selected",
						v.DocFn.Name, mt.Name, modDecl.Name)
					return false
				}
				// we got a suffix
				q := mt.TypeName.Type.Part.Qual
				if q.Qual != nil {
					log.Printf("view function '%s' has a parameter '%s' which is a model ('%s') but message name cannot have qualifier ('%s')",
						v.DocFn.Name, mt.Name, modDecl.Name, mt.TypeName.String())
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

func (v *ViewDecl) FinalizeSemantics(filename string) error {
	_, msg, err := v.Section.ResolveModelMessageTypeByIdent(filename, v.ModelName)
	if err != nil {
		return err
	}
	v.Message = msg
	t := &TypeName{
		TypeStarter: "*",
		Type:        v.ModelName,
	}
	formal := NewPFormal("model", t, t.String(), v.ModelName.LineNumber, v.ModelName.ColumnNumber)
	v.DocFn.Param = append([]*PFormal{formal}, v.DocFn.Param...)
	return nil
}

type ProtobufFileNode struct {
	Package    string
	Filename   string
	GoPkg      string
	LocalGoPkg string
	ImportFile []string
	Import     []*ProtobufFileNode
	Message    []*ProtobufMessage
}

func NewProtobufFileNode() *ProtobufFileNode {
	return &ProtobufFileNode{}
}

type TypeLocation struct {
	File    *ProtobufFileNode
	Message *ProtobufMessage
}
type ProtobufMessage struct {
	Name       string
	Package    string
	LocalGoPkg string
	Field      map[string]*ProtobufMsgElem
	Location   map[string]*TypeLocation
	AnyField   bool // this is true if we are not going to check the fields of this message (usually a google type)
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
	return true
}
func (m *MVCSectionNode) MoveDocFns() {
	for _, view := range m.ViewDecl {
		m.Program.DocSection.AttachViewToSection(view)
	}
	// this is kinda tricky, we cannot just set ViewDecl to nil here
	// although that we be nice.  We need to run some checks that walk
	// the parameters/expressions and make sure they are ok
	//m.ViewDecl = nil
}

func (m *MVCSectionNode) FinalizeSemantics(filename string) error {
	if m == nil {
		return nil
	}

	if m.ModelDecl != nil {
		for _, mod := range m.ModelDecl {
			mod.Section = m
		}

	}
	if m.ViewDecl != nil {
		for _, view := range m.ViewDecl {
			view.Section = m
		}
		// we have to do this after the section gets assigned
		m.MoveDocFns()
	}
	if m.ControllerDecl != nil {
		for _, cont := range m.ControllerDecl {
			cont.Section = m
		}
	}
	// we want all the section pointers set up FIRST, then we can run this code
	if m.ModelDecl != nil {
		for _, mod := range m.ModelDecl {
			if err := mod.FinalizeSemantics(filename); err != nil {
				return err
			}
		}
		for _, v := range m.ViewDecl {
			if err := v.FinalizeSemantics(filename); err != nil {
				return err
			}
		}
		for _, c := range m.ControllerDecl {
			if err := c.FinalizeSemantics(filename); err != nil {
				return err
			}
		}
	}
	return nil
}

type ControllerDecl struct {
	ModelName *Ident
	Message   *ProtobufMessage
	Section   *MVCSectionNode
	Spec      []*EventSpec
}

func NewControllerDecl() *ControllerDecl {
	return &ControllerDecl{}
}

func (c *ControllerDecl) FinalizeSemantics(filename string) error {
	_, msg, err := c.Section.ResolveModelMessageTypeByIdent(filename, c.ModelName)
	if err != nil {
		return err
	}
	c.Message = msg
	for _, spec := range c.Spec {
		act := spec.Function.Actual
		extra := &FuncActual{
			Ref: &ValueRef{
				Lit: "actual",
			},
		}
		spec.Function.Actual = append([]*FuncActual{extra}, act...)
	}
	return nil
}
