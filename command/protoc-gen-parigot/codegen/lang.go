package codegen

import (
	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"
	"text/template"

	"google.golang.org/protobuf/types/descriptorpb"
)

// Generator represents a language-specific collection of files and code that can
// convert a parigot proto file into code for the specific language.  This type is
// also "borrowed" for the use of generating ABI interfaces, although it works
// the same way.
type Generator interface {
	Process(proto *descriptorpb.FileDescriptorProto) error
	Generate(t *template.Template, g *GenInfo) ([]*util.OutputFile, error)
	TemplateName() []string
	FuncMap() template.FuncMap
	GeneratingMessage() []string
	ResultName() []string
	LanguageText() LanguageText //golang actually puts a AbiLanguageText in here
}
