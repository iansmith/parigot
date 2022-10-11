package go_

import (
	"fmt"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"io"
	"log"
	"strings"
	"text/template"
)

const (
	serviceDecl = "template/go/servicedecl.tmpl"
	messageDecl = "template/go/messagedecl.tmpl"
	simpleLoc   = "template/go/servicesimpleloc.tmpl"
)

type GoGen struct {
	types map[string]*descriptorpb.FileDescriptorProto
}

func (g *GoGen) addType(name string, fdp *descriptorpb.FileDescriptorProto) {
	if g.types == nil {
		g.types = make(map[string]*descriptorpb.FileDescriptorProto)
	}
	g.types[name] = fdp
}

func (g *GoGen) TemplateName() []string {
	return []string{serviceDecl, messageDecl, simpleLoc}
}
func (g *GoGen) FuncMap() template.FuncMap {
	return goFuncMap
}

func (g *GoGen) Process(proto *descriptorpb.FileDescriptorProto) error {
	g.addType(proto.GetName(), proto)
	return nil
}

type wasmService struct {
	*descriptorpb.ServiceDescriptorProto
	WasmServiceName string
	GoPackage       string
	ProtoFile       string
	ProtoPackage    string
}

func (g *GoGen) Generate(t *template.Template, proto *descriptorpb.FileDescriptorProto,
	locators []string) ([]*util.OutputFile, error) {
	svcOutName := util.GenerateOutputFilenameBase(proto) + "svc.p.go"
	log.Printf("services being generated to %s from proto file %s", svcOutName, proto.GetName())
	f := util.NewOutputFile(svcOutName)
	pkg := proto.GetPackage()
	if strings.LastIndex(pkg, ".") != -1 {
		last := strings.LastIndex(pkg, ".")
		if last != len(pkg)-1 {
			pkg = pkg[last+1:]
		}
	}
	fmt.Fprintf(f, "package %s\n", pkg)
	for _, s := range proto.GetService() {
		wasmName := s.GetName()

		//optName := wasmNameFromComment(s.Comments.Leading.String())
		//if optName != "" {
		//	wasmName = optName
		//}
		w := &wasmService{
			ProtoFile:       proto.GetSourceCodeInfo().String(),
			GoPackage:       proto.GetOptions().GetGoPackage(),
			WasmServiceName: wasmName,
			ProtoPackage:    proto.GetPackage(),
		}
		w.ServiceDescriptorProto = s
		if err := generateCodeService(f, w, t); err != nil {
			return nil, err
		}
		for _, loc := range locators {
			switch loc {
			case "simple":
				if err := generateCodeSimpleLoc(f, w, t.Lookup(simpleLoc)); err != nil {
					return nil, err
				}
			default:
				log.Printf("unknown locator %s, ignoring", loc)
			}
		}
	}
	return []*util.OutputFile{f}, nil
}

func generateCodeSimpleLoc(w io.WriteCloser, svc *wasmService, tmpl *template.Template) error {
	return runTemplateForService(w, svc, tmpl)
}

func generateCodeService(w io.WriteCloser, svc *wasmService, tmpl *template.Template) error {
	return runTemplateForService(w, svc, tmpl)
}

func runTemplateForService(w io.WriteCloser, svc *wasmService, tmpl *template.Template) error {
	data := map[string]interface{}{
		"svc": svc,
	}
	t := tmpl.Lookup(serviceDecl)
	if t == nil {
		return fmt.Errorf("unable to find %s in template set", serviceDecl)
	}
	return t.Execute(w, data)
}

func generateCodeMessage(w io.WriteCloser, msg *protogen.Message, tmpl *template.Template) error {
	data := map[string]interface{}{
		"msg": msg,
	}
	log.Printf("message is %s\n", msg.GoIdent.GoName)
	t := tmpl.Lookup(messageDecl)
	if t == nil {
		return fmt.Errorf("unable to find %s in template set", messageDecl)
	}
	return t.Execute(w, data)
}

func wasmNameFromComment(s string) string {
	settings := util.ParigotCommentSettings(util.ParigotCommentLine(s))
	if settings == nil {
		return ""
	}
	wasmName, ok := settings[util.WasmServiceName]
	if !ok {
		return ""
	}
	return wasmName
}
