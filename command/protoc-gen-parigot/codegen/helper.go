package codegen

import (
	"strings"
	"text/template"
)

const serviceIdSuffix = "_sid"

// FuncMap is "default" function map for use with the templates. This contains generally
// useful helper functions.
var FuncMap = template.FuncMap{
	"toCamelCase":          ToCamelCase,
	"lastSegmentOfPackage": lastSegmentOfPackage,
}

// ToCamelCase converts from snake case to camel case. If the first or last character
// is a underscore, it is left alone.
func ToCamelCase(snake string) string {
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

// Last segment of a fully spelled out package name.
func lastSegmentOfPackage(pkg string) string {
	last := strings.LastIndex(pkg, ".")
	if last == -1 || last == len(pkg)-1 {
		return pkg
	}
	return pkg[last+1:]
}
