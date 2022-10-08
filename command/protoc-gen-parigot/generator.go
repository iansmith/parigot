package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"log"
	"text/template"
)

func generateCode(outFile *protogen.GeneratedFile, svc *protogen.Service, misc *MiscData) {
	tmpl := template.Must(template.ParseFS(templateFS, "template/*.tmpl"))
	t := tmpl.Lookup("tomldecl.tmpl")
	err := generateCodeForProtoService(t, outFile, svc, misc)
	if err != nil {
		log.Fatalf("unable to execute template for service: %v", err)
	}
}

func generateCodeForProtoService(tmpl *template.Template, outFile *protogen.GeneratedFile, svc *protogen.Service, misc *MiscData) error {
	data := map[string]interface{}{
		"svc":  svc,
		"misc": misc,
	}

	return tmpl.Execute(outFile, data)
}
