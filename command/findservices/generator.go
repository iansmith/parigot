package main

import (
	structure "github.com/iansmith/parigot/command/toml"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func generateCode(result *SearchResult) {
	all, err := templateFS.ReadFile("template/servicedecl.tmpl")
	if err != nil {
		log.Fatalf("cannot parse filesystem of templates:%v", err)
	}
	text := string(all)
	tmpl, err := template.New("servicedecl.tmpl").Funcs(serviceDeclFuncMap).Parse(text)
	if err != nil {
		log.Fatalf("failed to parse servicedecl.tmpl:%v", err)
	}
	for _, project := range result.ProjectsFound {
		for _, svc := range project.ServicesFound {
			err := generateCodeForService(tmpl.Lookup("servicedecl.tmpl").Funcs(serviceDeclFuncMap), svc, project)
			if err != nil {
				log.Fatalf("unable to execute template for service: %v", err)
			}
		}
	}
}

func generateCodeForService(tmpl *template.Template, svc *structure.ServiceDecl, project *structure.ProjectDecl) error {
	target := svc.TargetDir
	path := filepath.Join(target, "servicedecl.go")
	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"svc":     svc,
		"project": project,
	}

	return tmpl.Funcs(serviceDeclFuncMap).Execute(fp, data)
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

var serviceDeclFuncMap = template.FuncMap{
	"toCamelCase": toCamelCase,
}
