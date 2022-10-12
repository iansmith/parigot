package abi

import (
	"fmt"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/go_"
	"strings"
	"text/template"

	"github.com/iansmith/parigot/command/protoc-gen-parigot/codegen"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"

	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	tinygo = "template/abi/abitinygo.tmpl"
	ide    = "template/abi/abiide.tmpl"
	helper = "template/abi/abiimplhelper.tmpl"
)

var abiFuncMap = template.FuncMap{}
var resultFile = []string{".p.go", "null.p.go", "helper.go"}

type AbiGen struct {
	types map[string]*descriptorpb.FileDescriptorProto
}

func (a *AbiGen) ResultName() []string {
	return resultFile
}

func (a *AbiGen) addType(name string, fdp *descriptorpb.FileDescriptorProto) {
	if a.types == nil {
		a.types = make(map[string]*descriptorpb.FileDescriptorProto)
	}
	a.types[name] = fdp
}

func (a *AbiGen) Process(proto *descriptorpb.FileDescriptorProto) error {
	a.addType(proto.GetName(), proto)
	return nil
}

var text = []string{"fake abi for ide into ", "abi function declarations into",
	"helpers for abi implementation of functions into"}

func (a *AbiGen) NeedsLocators() bool {
	return false
}

func (a *AbiGen) GeneratingMessage() []string {
	return text
}

func (a *AbiGen) Generate(t *template.Template, info *codegen.GenInfo, loc []string) ([]*util.OutputFile, error) {
	if !strings.HasSuffix(info.GetFile().GetName(), "abi.proto") {
		return nil, fmt.Errorf("unable to understand abi definition with %d services (expecting 1)",
			len(info.GetFile().GetService()))
	}
	return codegen.BasicGenerate(a, t, info, loc)
}

func (a *AbiGen) TemplateName() []string {
	return []string{tinygo, ide, helper}

}
func (a *AbiGen) FuncMap() template.FuncMap {
	return abiFuncMap
}
func (a *AbiGen) LanguageText() codegen.LanguageText {
	return &go_.GoText{}
}
