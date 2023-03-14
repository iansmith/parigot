package tree

import (
	"github.com/iansmith/parigot/pbmodel"
)

type ModelSectionNode struct {
	Program  *ProgramNode
	ModelDef []*ModelDef
}

func NewModelSection() *ModelSectionNode {
	return &ModelSectionNode{}
}

type ModelDef struct {
	Name string
	Path []string
	File []*ProtobufFileNode
}

func NewModelDef() *ModelDef {
	return &ModelDef{}
}

type ProtobufFileNode struct {
	PackageName string
	FileName    string
	GoPkg       string
	Imports     []*ProtobufFileNode
	Message     []*ProtobufMessage
}

func NewProtobufFileNode() *ProtobufFileNode {
	return &ProtobufFileNode{}
}

type ProtobufMessage struct {
	Name string
}

// func (m *ModelSectionNode) Parse() (string, bool) {
// 	for _, def := range m.ModelDef {
// 		_, bad, ok := def.Parse()
// 		if !ok {
// 			return bad, false
// 		}
// 	}
// 	return "", true
// }

// func (m *ModelDef) Parse() (*ProtobufFileNode, string, bool) {
// 	builder := pbmodel.NewPb3Builder()
// 	for _, f := range m.Path {
// 		pf, bad, ok := antlr.EvaluateOneFile(f, builder)
// 		if !ok {
// 			return nil, bad, false
// 		}
// 		m.File = append(m.File, pf)
// 	}
// 	return nil, "", true
// }

func ProtobufNodeFromBuilder(builder *pbmodel.Pb3Builder) *ProtobufFileNode {
	pf := NewProtobufFileNode()
	pf.FileName = builder.CurrentFile
	pf.PackageName = builder.CurrentPackage
	pf.GoPkg = builder.CurrentGoPackage
	return pf
}
