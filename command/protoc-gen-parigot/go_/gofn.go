package go_

import (
	"google.golang.org/protobuf/compiler/protogen"
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

func outputName(m *protogen.Method) string {
	return emtpyToString(m.Input.GoIdent.GoName)
}
func inputParamName(m *protogen.Method) string {
	return emtpyToString(m.Input.GoIdent.GoName)
}

func emtpyToString(s string) string {
	if s == "Empty" {
		return ""
	}
	return s
}
