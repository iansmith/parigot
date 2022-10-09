package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"log"
	"text/template"
)

var allTemplates = template.Must(template.ParseFS(templateFS, "template/*.tmpl"))

func generateCodeService(outFile *protogen.GeneratedFile, svc *protogen.Service, misc *MiscData) {
	t := allTemplates.Lookup("tomldeclservice.tmpl")
	if t == nil {
		log.Fatalf("could not find template tomldeclservice.tmpl ")
	}
	err := generateCodeForProtoService(t, outFile, svc, misc)
	if err != nil {
		log.Fatalf("unable to execute template for service: %v", err)
	}
}
func generateCodeMessage(outFile *protogen.GeneratedFile, msg *protogen.Message) {
	t := allTemplates.Lookup("tomldeclmessage.tmpl")
	if t == nil {
		log.Fatalf("could not find template tomldeclmessage.tmpl ")
	}
	err := generateCodeForProtoMessage(t, outFile, msg)
	if err != nil {
		log.Fatalf("unable to execute template for service: %v", err)
	}
}

func generateCodeForProtoService(tmpl *template.Template, outFile *protogen.GeneratedFile,
	svc *protogen.Service, misc *MiscData) error {
	data := map[string]interface{}{
		"svc":  svc,
		"misc": misc,
	}

	return tmpl.Execute(outFile, data)
}
func generateCodeForProtoMessage(tmpl *template.Template, outFile *protogen.GeneratedFile,
	msg *protogen.Message) error {
	data := map[string]interface{}{
		"msg": msg,
	}
	log.Printf("---- %s", msg.GoIdent.GoName)
	return tmpl.Execute(outFile, data)
}
