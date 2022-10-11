package main

import (
	"github.com/iansmith/parigot/command/protoc-gen-parigot/abi"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/go_"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"
	"google.golang.org/protobuf/types/descriptorpb"
	"text/template"
)

type Generator interface {
	Process(proto *descriptorpb.FileDescriptorProto) error
	Generate(t *template.Template, proto *descriptorpb.FileDescriptorProto, locators []string) ([]*util.OutputFile, error)
	TemplateName() []string
	FuncMap() template.FuncMap
}

var GeneratorMap = map[string]Generator{
	"go": &go_.GoGen{},
}

var AbiOnlyMap = map[string]Generator{
	"abi": &abi.AbiGen{},
}
