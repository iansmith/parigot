package go_

import (
	"github.com/iansmith/parigot/command/protoc-gen-parigot/codegen"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"
	"google.golang.org/protobuf/types/descriptorpb"
	"text/template"
)

const (
	serviceDecl = "template/go/servicedecl.tmpl"
	messageDecl = "template/go/messagedecl.tmpl"
	//simpleLoc   = "template/go/servicesimpleloc.tmpl"
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
	return []string{"service declarations into", "message declarations into"}
	//	"locator declarations into",

}
func (g *GoGen) ResultName() []string {
	return []string{"servicedecl.p.go", "messagedecl.p.go"}
	//	"locdecl.p.go",
}

func (g *GoGen) TemplateName() []string {
	return []string{serviceDecl, messageDecl} // simpleLoc}
}
func (g *GoGen) FuncMap() template.FuncMap {
	return nil
}

func (g *GoGen) Process(pr *descriptorpb.FileDescriptorProto) error {
	codegen.AddFileContentToFinder(g.finder, pr, g.lang)
	return nil
}

func (g *GoGen) Generate(t *template.Template, info *codegen.GenInfo) ([]*util.OutputFile, error) {
	return codegen.BasicGenerate(g, t, info)
}

func (g *GoGen) LanguageText() codegen.LanguageText {
	return g.lang
}
