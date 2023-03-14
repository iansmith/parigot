package antlr

import (
	"github.com/iansmith/parigot/pbmodel"
	"github.com/iansmith/parigot/ui/parser/tree"
)

func ParseModelSection(m *tree.ModelSectionNode) (string, bool) {
	for _, def := range m.ModelDef {
		_, bad, ok := ParseModelDef(def)
		if !ok {
			return bad, false
		}
	}
	return "", true
}

func ParseModelDef(def *tree.ModelDef) (*tree.ProtobufFileNode, string, bool) {
	builder := pbmodel.NewPb3Builder()
	for _, f := range def.Path {
		pf, bad, ok := EvaluateOneFile(f, builder)
		if !ok {
			return nil, bad, false
		}
		def.File = append(def.File, pf)
	}
	return nil, "", true
}
