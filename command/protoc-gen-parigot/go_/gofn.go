package go_

import (
	"google.golang.org/protobuf/types/descriptorpb"
	"strings"
)

import (
	"text/template"
)

var goFuncMap = template.FuncMap{
	"toCamelCase":    toCamelCase,
	"inputParamName": inputParamName,
	"outputName":     outputName,
}

func toCamelCase(snake string) string {
	if len(snake) == 0 {
		return ""
	}
	snake = strings.ToUpper(snake[0:1]) + snake[1:]
	index := strings.Index(snake, "_")
	// allow _ in first & last spot
	for index != -1 && index != len(snake)-1 && index != 0 {
		snake = snake[:index] + strings.ToUpper(snake[index+1:index+2]) + snake[index+2:]
		index = strings.Index(snake, "_")
	}
	return snake
}

func outputName(m *descriptorpb.MethodDescriptorProto) string {
	return m.GetName()
}

func inputParamName(m *descriptorpb.MethodDescriptorProto) string {
	return m.GetName()
}
