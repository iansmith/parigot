package codegen

import (
	"strings"
	"text/template"
)

// FuncMap is "default" function map for use with the templates. This contains generally
// useful helper functions.
var FuncMap = template.FuncMap{
	"toCamelCase":           ToCamelCase,
	"toCamelCaseFirstLower": ToCamelCaseFirstLower,
	"LastSegmentOfPackage":  LastSegmentOfPackage,
	"BasicTypeToString":     BasicTypeToString,
}

func ToCamelCaseFirstLower(snake string) string {
	c := ToCamelCase(snake)
	if len(c) == 1 {
		return strings.ToLower(c)
	}
	return strings.ToLower(c[0:1]) + c[1:]
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

// LastSegmentOfPackage returns the string after the last dot of a fully spelled out package name.
// If there are no dots or the last dot is the last character, it returns its input.
func LastSegmentOfPackage(pkg string) string {
	last := strings.LastIndex(pkg, ".")
	if last == -1 || last == len(pkg)-1 {
		return pkg
	}
	return pkg[last+1:]
}

func BasicTypeToString(l LanguageText, s string, panicOnFail bool) string {
	return l.BasicTypeToString(s, panicOnFail)
}
