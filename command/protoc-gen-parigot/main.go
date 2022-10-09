package main

import (
	"embed"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"log"
	"strings"
)

//go:embed template/*
var templateFS embed.FS

const (
	wasmServiceName = "WasmServiceName"
	wasmFuncName    = "WasmFuncName"
)

var parigotPrefixes = []string{
	"//parigot:",
	"// parigot:",
}

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f)
		}
		return nil
	})
}

type MiscData struct {
	ProtoFile string
	GoPackage string
	WasmName  string
}

// generateFile generates a suitable interface for use with parigot
func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + ".p.toml"

	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	for _, s := range file.Services {
		wasmName := s.GoName
		optName := wasmNameFromComment(s.Comments.Leading.String())
		if optName != "" {
			wasmName = optName
		}
		m := &MiscData{
			ProtoFile: file.Desc.Path(),
			GoPackage: fmt.Sprint(file.GoPackageName),
			WasmName:  wasmName,
		}
		generateCodeService(g, s, m)
		g.P()
	}
	for _, msg := range file.Messages {
		log.Printf("generating message %s", msg.GoIdent.GoName)
		generateCodeMessage(g, msg)
		msg.Fields[0].Desc.I
		g.P()
	}
	g.P()

	return g
}
func wasmNameFromComment(s string) string {
	settings := parigotCommentSettings(parigotCommentLine(s))
	if settings == nil {
		return ""
	}
	wasmName, ok := settings[wasmServiceName]
	if !ok {
		return ""
	}
	return wasmName
}
func parigotCommentLine(allLines string) string {
	lines := strings.Split(allLines, "\n") // xxx break on windows?
	if len(lines) == 0 {
		return ""
	}
	line := ""
	for l := len(lines) - 1; l >= 0; l-- {
		if lines[l] != "" {
			line = lines[l]
			break
		}
	}
	return line
}

func parigotCommentSettings(line string) map[string]string {
	if line == "" {
		return nil
	}
	found := ""
	for _, p := range parigotPrefixes {
		if strings.HasPrefix(line, p) {
			found = strings.TrimPrefix(line, p)
			break
		}
	}
	if found == "" {
		return nil
	}
	parts := strings.Split(found, ",")
	if len(parts) == 0 {
		panic("cant parse comment line: " + line)
	}
	settings := make(map[string]string)
	for _, setting := range parts {
		parts := strings.Split(setting, "=")
		if len(parts) != 2 {
			log.Printf("can't parse parigot comment setting on line:" + line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		settings[key] = val
	}
	return settings
}
