package codegen

import (
	"text/template"

	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"

	"google.golang.org/protobuf/types/descriptorpb"
)

// Generator represents a language-specific collection of files and code that can
// convert a parigot proto file into code for the specific language.
type Generator interface {
	Process(proto *descriptorpb.FileDescriptorProto) error
	Generate(t *template.Template, g *GenInfo, imp2Pkg map[string]string) ([]*util.OutputFile, error)
	TemplateName() []string
	FuncMap() template.FuncMap
	GeneratingMessage() []string
	ResultName() []string
	LanguageText() LanguageText
}
