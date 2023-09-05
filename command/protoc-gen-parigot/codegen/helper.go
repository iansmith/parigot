package codegen

import (
	"bytes"
	"regexp"
	"strings"
	"text/template"
	"unicode"
)

var versionExpr = regexp.MustCompile("v[0-9]+")

// FuncMap is "default" function map for use with the templates. This contains generally
// useful helper functions.
var FuncMap = template.FuncMap{
	"toCamelCase":           ToCamelCase,
	"toSnakeCase":           ToSnakeCase,
	"packageToGoName":       PackageToGoName,
	"toSnakeCaseFirstUpper": ToSnakeCaseFirstUpper,
	"toCamelCaseFirstLower": ToCamelCaseFirstLower,
	"toLowerNoService":      ToLowerNoService,
	"LastSegmentOfPackage":  LastSegmentOfPackage,
	"BasicTypeToString":     BasicTypeToString,
	"RemoveDots":            RemoveDots,
}

func ToCamelCaseFirstLower(snake string) string {
	c := ToCamelCase(snake)
	if len(c) == 1 {
		return strings.ToLower(c)
	}
	return strings.ToLower(c[0:1]) + c[1:]
}
func ToLowerNoService(snake string) string {
	c := ToCamelCase(snake)
	l := strings.ToLower(c)
	if len(c) == 1 {
		return strings.ToLower(c)
	}
	if strings.HasSuffix(l, "service") {
		return strings.TrimSuffix(l, "service")
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

// ToSnakeCaseFirstUpper converts from CamelCase to snake_case but makes sure the first character is upper.
func ToSnakeCaseFirstUpper(camel string) string {
	s := ToSnakeCase(camel)
	if s == "" {
		return s
	}
	first, rest := s[0:1], s[1:]
	return strings.ToUpper(first) + rest
}

// ToSnakeCase converts from CamelCase to snake_case.
func ToSnakeCase(camel string) string {
	if len(camel) < 2 {
		return camel
	}
	inLower := false
	if strings.ToLower(camel[0:1]) == camel[0:1] {
		inLower = true
	}
	result := &bytes.Buffer{}

	for _, r := range camel {
		if !unicode.IsLetter(r) {
			result.WriteRune(r)
			continue
		}
		if inLower {
			if unicode.IsUpper(r) {
				result.WriteString("_")
				inLower = false
			}
			result.WriteRune(unicode.ToLower(r))
			continue
		}
		// we are inLower==false
		if unicode.IsUpper(r) {
			// possible acronym liko blahFooWASM => blah_foo_wasm
			result.WriteRune(unicode.ToLower(r))
			continue
		}
		// lower case when we are in inLower==false
		result.WriteRune(r)
		inLower = true
	}
	return result.String()
}

func PackageToGoName(pkg string) string {
	if !strings.Contains(pkg, ".") {
		return pkg
	}
	part := strings.Split(pkg, ".")
	target := len(part) - 1
	if isVersion(part[len(part)-1]) {
		target--
	}
	return part[target]
}

// LastSegmentOfPackage takes a file to be generated file file/v1/file.proto and converts it to
// v1.
func LastSegmentOfPackage(pkg string) string {
	pkg = strings.TrimSuffix(pkg, ".proto")
	part := strings.Split(pkg, "/")
	if len(part) > 1 {
		possibleVersion := part[len(part)-2]
		if isVersion(possibleVersion) {
			part = part[:len(part)-1]
			pkg = strings.Join(part, ".")
		}
	}
	suffix := ""
	if strings.HasPrefix(pkg, ".msg.") {
		suffix = "msg"
	}
	last := strings.LastIndex(pkg, ".")
	if last == -1 || last == len(pkg)-1 {
		return pkg + suffix
	}
	return pkg[last+1:] + suffix
}
func isVersion(possibleVersion string) bool {
	return versionExpr.MatchString(possibleVersion)
}

func BasicTypeToString(l LanguageText, s string, panicOnFail bool) string {
	return l.BasicTypeToString(s, panicOnFail)
}

// RemoveDots takes a string and returns it without dots.
func RemoveDots(str string) string {
	return strings.Replace(str, ".", "", -1)
}
