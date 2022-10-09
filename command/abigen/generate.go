package main

import (
	"log"
	"os"
	"path/filepath"
	"text/template"

	structure "github.com/iansmith/parigot/command/toml"
)

// ugh these consts are awful
const AbiIdeIndex = 0
const AbiTinygoIndex = 1
const AbiIde = "template/abiide.tmpl"
const AbiTinygo = "template/abitinygo.tmpl"
const AbiIdeTarget = "abinull.go"
const AbiTinygoTarget = "abi.go"

var templateSource = []string{AbiIde, AbiTinygo}

func generateCode(svc *structure.ServiceDecl, project *structure.ProjectDecl) {
	abiTemplate := make([]*template.Template, len(templateSource))
	for i, path := range templateSource {
		all, err := templateFS.ReadFile(path)
		if err != nil {
			log.Fatalf("cannot parse filesystem of templates:%v", err)
		}
		text := string(all)
		tmpl, err := template.New(path).Parse(text)
		if err != nil {
			log.Fatalf("failed to parse %s:%v", path, err)
		}
		abiTemplate[i] = tmpl
	}
	err := generateIDECode(abiTemplate[AbiIdeIndex].Lookup(AbiIde).Funcs(abiFuncs), svc, project)
	if err != nil {
		log.Fatalf("unable to execute template in abi ide code generation: %v", err)
	}
	err = generateTinygoCode(abiTemplate[AbiTinygoIndex].Lookup(AbiTinygo).Funcs(abiFuncs), svc, project)
	if err != nil {
		log.Fatalf("unable to execute template in abi tinygo code generation: %v", err)
	}
}

func createTarget(path string) *os.File {
	fp, err := os.Create(path)
	if err != nil {
		log.Fatalf("unable to create target of code generation %s:%v", path, err)
	}
	return fp
}

func generateIDECode(t *template.Template, svc *structure.ServiceDecl, proj *structure.ProjectDecl) error {
	path := filepath.Join(svc.TargetDir, AbiIdeTarget)
	return generateABICode(path, t, svc, proj)
}

func generateABICode(path string, t *template.Template, svc *structure.ServiceDecl, proj *structure.ProjectDecl) error {
	fp := createTarget(path)
	defer fp.Close()
	data := map[string]interface{}{
		"svc":  svc,
		"proj": proj,
	}
	return t.Execute(fp, data)
}
func generateTinygoCode(t *template.Template, svc *structure.ServiceDecl, proj *structure.ProjectDecl) error {
	path := filepath.Join(svc.TargetDir, AbiTinygoTarget)
	return generateABICode(path, t, svc, proj)
}

func paramSwap(elem *structure.MethodDecl) string {
	if elem.Input == "Empty" {
		return ""
	}
	return ""
}

var abiFuncs = template.FuncMap{
	"paramSwap": paramSwap,
}
