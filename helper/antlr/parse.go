package antlr

import (
	"github.com/iansmith/parigot/helper"
	"github.com/iansmith/parigot/pbmodel"
	"github.com/iansmith/parigot/ui/parser/tree"
)

func ParseModelSection(sourceFile, pkg string, m *tree.MVCSectionNode) (string, bool) {
	for _, def := range m.ModelDecl {
		_, bad, ok := ParseModelDecl(sourceFile, pkg, def)
		if !ok {
			return bad, false
		}
	}
	return "", true
}

func ParseModelDecl(sourceFile, pkg string, def *tree.ModelDecl) (*tree.ProtobufFileNode, string, bool) {
	builder := pbmodel.NewPb3Builder()
	for _, f := range def.Path {
		rel := helper.RelativePath(f, sourceFile, pkg)
		pf, bad, ok := EvaluateOneFile(rel, pkg, builder)
		if !ok {
			return nil, bad, false
		}
		//log.Printf("ParseModelDef: %s[%s] -> %#v with %d imports", sourceFile, pkg, pf, len(pf.Import))
		def.File = append(def.File, pf)
	}
	return nil, "", true
}
