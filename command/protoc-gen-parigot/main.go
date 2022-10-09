package main

import (
	"embed"
	"fmt"
	"log"
	"text/template"

	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"

	"google.golang.org/protobuf/types/pluginpb"
)

//go:embed template/*
var templateFS embed.FS

func main() {
	genReq := util.ReadStdinIntoBuffer()
	files, err := generateNeutral(genReq)
	if err != nil {
		log.Fatalf("unable to process plugin request: %v", err)
	}
	log.Printf("created %d files", len(files))
}

func generateNeutral(genReq *pluginpb.CodeGeneratorRequest) ([]*util.OutputFile, error) {
	fileList := []*util.OutputFile{}
	for _, desc := range genReq.GetProtoFile() {
		generate := false
		// figure out which ones to generate and which ones to process
		for _, gen := range genReq.GetFileToGenerate() {
			if desc.GetName() == gen {
				generate = true
				break
			}
		}
		for lang, generator := range GeneratorMap {
			if generate {
				// load up templates
				t, err := loadTemplates(generator)
				if err != nil {
					return nil, err
				}
				file, err := generator.Generate(t, desc)
				if err != nil {
					return nil, err
				}
				if file == nil {
					log.Printf("warning: language %s did not create any output for %s, ignoring",
						lang, desc.GetName())
				} else {
					fileList = append(fileList, file...)
				}
			} else {
				if err := generator.Process(desc); err != nil {
					return nil, err
				}
			}
		}
	}
	return fileList, nil
}

func loadTemplates(generator Generator) (*template.Template, error) {
	// create root template and add functions, if any
	root := template.New("root")
	if generator.FuncMap() != nil {
		root = root.Funcs(generator.FuncMap())
	}
	t := root
	// template loading
	for _, f := range generator.TemplateName() {
		all, err := templateFS.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("unable to read file %s from embedded fs:%v", f, err)
		}
		t, err = root.New(f).Parse(string(all))
		if err != nil {
			return nil, fmt.Errorf("unable to parse template %s:%v", f, err)
		}
	}
	return t, nil
}
