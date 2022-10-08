package main

import (
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func generateCode(result *SearchResult) {
	tmpl := template.Must(template.ParseFS(templateFS, "template/*.tmpl"))
	for _, project := range result.ProjectsFound {
		for _, svc := range project.ServicesFound {
			err := generateCodeForService(tmpl, svc, project)
			if err != nil {
				log.Fatalf("unable to execute template for service: %v", err)
			}
		}
	}
}

func generateCodeForService(tmpl *template.Template, svc *ServiceDecl, project *ProjectDecl) error {
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
	return tmpl.Execute(fp, data)
}
