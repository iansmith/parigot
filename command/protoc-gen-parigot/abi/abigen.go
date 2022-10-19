package abi

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/iansmith/parigot/command/protoc-gen-parigot/codegen"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/go_"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"

	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	tinygo = "template/abi/abitinygo.tmpl"
	ide    = "template/abi/abiide.tmpl"
	helper = "template/abi/abiimplhelper.tmpl"
	undef  = "template/abi/wasm-undefined.txt.tmpl"
)

var abiFuncMap = template.FuncMap{}

var resultFile = []string{".p.go", "null.p.go", "helper.p.go", "wasm-undefined.txt"}

type AbiGen struct {
	finder codegen.Finder
	lang   codegen.AbiLanguageText
}

func NewAbiGen(f codegen.Finder) *AbiGen {
	return &AbiGen{finder: f, lang: &go_.GoText{}}
}

func (a *AbiGen) ResultName() []string {
	return resultFile
}

func (a *AbiGen) Process(proto *descriptorpb.FileDescriptorProto) error {
	codegen.AddFileContentToFinder(a.finder, proto, a.lang)
	return nil
}

var text = []string{"fake abi for ide into ", "abi function declarations into",
	"helpers for abi implementation of functions into",
	"names of symbols that can be undefined into"}

func (a *AbiGen) GeneratingMessage() []string {
	return text
}

func (a *AbiGen) Generate(t *template.Template, info *codegen.GenInfo) ([]*util.OutputFile, error) {
	if !strings.HasSuffix(info.GetFile().GetName(), "abi.proto") {
		return nil, fmt.Errorf("unable to understand abi definition with %d services (expecting 1)",
			len(info.GetFile().GetService()))
	}
	return codegen.BasicGenerate(a, t, info)
}

func (a *AbiGen) TemplateName() []string {
	return []string{tinygo, ide, helper, undef}
}
func (a *AbiGen) FuncMap() template.FuncMap {
	return abiFuncMap
}

func (a *AbiGen) LanguageText() codegen.LanguageText {
	return a.lang
}
