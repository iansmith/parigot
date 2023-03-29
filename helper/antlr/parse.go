package antlr

import (
	"fmt"
	"log"

	"github.com/iansmith/parigot/helper"
	"github.com/iansmith/parigot/pbmodel"
	"github.com/iansmith/parigot/ui/parser/tree"
)

func ParseModelSection(sourceFile, pkg string, m *tree.MVCSectionNode) (*tree.MVCSectionNode, error) {
	for _, def := range m.ModelDecl {
		err := ParseModelDecl(sourceFile, pkg, def)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

// ParseModelDecl walks the source file given and adds the successfully parsed
// ProtobufNode to the list associated with the def given.  If it
// fails, it returns an error otherwise returns nil.
func ParseModelDecl(sourceFile, pkg string, def *tree.ModelDecl) error {
	builder := pbmodel.NewPb3Builder()
	seen := make(map[string]struct{})
	for _, f := range def.Path {
		rel := helper.RelativePath(f, sourceFile, pkg)
		if _, ok := seen[rel]; ok {
			log.Printf("already processed relative path %s in this model", rel)
			continue
		} else {
			seen[rel] = struct{}{}
		}
		pf, bad, ok := EvaluateOneFile(rel, pkg, builder)
		if !ok {
			return fmt.Errorf("%s", bad)
		}
		if pf == nil {
			panic("pf is nil")
		}
		def.File = append(def.File, pf)
	}
	return nil
}
