package go_

import (
	"google.golang.org/protobuf/types/descriptorpb"
	"log"
	"strings"
)

import (
	"text/template"
)

var goFuncMap = template.FuncMap{
	"inputParamName":    inputParamName,
	"outputName":        outputName,
	"outputZeroVal":     outputZeroVal,
	"inputParamNameSet": inputParamNameSet,
	"hasInput":          hasInput,
	"hasOutput":         hasOutput,
}

const emptyType = ".google.protobuf.Empty"

func outputName(currPkg string, m *descriptorpb.MethodDescriptorProto) string {
	return shortenTypeName(currPkg, m.GetOutputType())
}

func shortenTypeName(currPkg string, t string) string {
	if t == emptyType {
		return ""
	}
	if t[0] == '.' {
		t = t[1:]
	}
	if currPkg[len(currPkg)-1] != '.' {
		currPkg += "."
	}
	log.Printf("shorten**: currPkg=%s, t=%s", currPkg, t)
	if strings.HasPrefix(t, currPkg) {
		return strings.TrimPrefix(t, currPkg)
	}
	return t
}

func inputParamName(m *descriptorpb.MethodDescriptorProto) string {
	return m.GetInputType()
}

func outputZeroVal(currPkg string, m *descriptorpb.MethodDescriptorProto) string {
	short := shortenTypeName(currPkg, m.GetOutputType())
	if short == "" {
		return ""
	}

	return short + "{}"
}

func inputParamNameSet(currPkg string, m *descriptorpb.MethodDescriptorProto) string {
	short := shortenTypeName(currPkg, m.GetInputType())
	if short == "" {
		return ""
	}
	return "input " + short
}

func hasInput(m *descriptorpb.MethodDescriptorProto) bool {
	return m.GetInputType() != emptyType
}
func hasOutput(m *descriptorpb.MethodDescriptorProto) bool {
	return m.GetOutputType() != emptyType
}
