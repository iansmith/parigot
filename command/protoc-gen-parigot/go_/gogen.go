package go_

import (
	"log"
	"text/template"

	"github.com/iansmith/parigot/command/protoc-gen-parigot/codegen"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	serviceDecl = "template/go/servicedecl.tmpl"
	serverDecl  = "template/go/serverdecl.tmpl"
)

type GoGen struct {
	finder codegen.Finder
	lang   codegen.LanguageText
}

func NewGoGen(finder codegen.Finder) *GoGen {
	gen := &GoGen{
		finder: finder,
		lang:   NewGoText(),
	}
	return gen
}

func (g *GoGen) GeneratingMessage() []string {
	return []string{"service declarations into", "server helpers"}

}
func (g *GoGen) ResultName() []string {
	return []string{"servicedecl.p.go", "server.p.go"}
}

func (g *GoGen) TemplateName() []string {
	return []string{serviceDecl, serverDecl}
}
func (g *GoGen) FuncMap() template.FuncMap {
	return nil
}

func (g *GoGen) Process(pr *descriptorpb.FileDescriptorProto) error {
	log.Printf("xxxx PROCESSing only: %s", pr.GetName())
	codegen.AddFileContentToFinder(g.finder, pr, g.lang)
	return nil
}

func (g *GoGen) Generate(t *template.Template, info *codegen.GenInfo, impToPkg map[string]string) ([]*util.OutputFile, error) {
	return codegen.BasicGenerate(g, t, info, impToPkg)
}

func (g *GoGen) LanguageText() codegen.LanguageText {
	return g.lang
}
