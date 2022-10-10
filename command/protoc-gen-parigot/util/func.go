package util

import (
	"strings"
	"text/template"
)

var FuncMap = template.FuncMap{
	"toCamelCase": toCamelCase,
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

func LastSegmentOfFQProtoPackageName(n string) string {
	last := strings.LastIndex(n, ".")
	if last == -1 {
		return n
	}
	return n[last+1:]
}
