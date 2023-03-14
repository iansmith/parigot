package antlr

import (
	"log"

	"github.com/iansmith/parigot/helper"
	"github.com/iansmith/parigot/pbmodel"
	"github.com/iansmith/parigot/ui/parser/tree"
)

func ParseModelSection(sourceFile, pkg string, m *tree.ModelSectionNode) (string, bool) {
	for _, def := range m.ModelDef {
		_, bad, ok := ParseModelDef(sourceFile, pkg, def)
		if !ok {
			return bad, false
		}
	}
	return "", true
}

func ParseModelDef(sourceFile, pkg string, def *tree.ModelDef) (*tree.ProtobufFileNode, string, bool) {
	builder := pbmodel.NewPb3Builder()
	for _, f := range def.Path {
		rel := helper.RelativePath(f, sourceFile, pkg)
		log.Printf("xxx out1 is %s, rel is %s, file is %s", f, rel, sourceFile)

		pf, bad, ok := EvaluateOneFile(rel, pkg, builder)
		if !ok {
			return nil, bad, false
		}
		def.File = append(def.File, pf)
	}
	return nil, "", true
}
