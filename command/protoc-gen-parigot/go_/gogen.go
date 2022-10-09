package go_

import (
	"fmt"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"io"
	"log"
	"text/template"
)

const (
	serviceDecl = "template/go/servicedecl.tmpl"
	messageDecl = "template/go/messagedecl.tmpl"
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
	return []string{serviceDecl, messageDecl}
}
func (g *GoGen) FuncMap() template.FuncMap {
	return goFuncMap
}

func (g *GoGen) Process(proto *descriptorpb.FileDescriptorProto) error {
	g.addType(proto.GetName(), proto)
	return nil
}

type wasmService struct {
	*protogen.Service
	WasmServiceName string
	GoPackage       string
	ProtoFile       string
}

func (g *GoGen) Generate(t *template.Template, proto *descriptorpb.FileDescriptorProto) ([]*util.OutputFile, error) {
	n := util.GenerateOutputFilenameBase(proto) + "svc.p.go"
	if n != "" {
		panic(n)
	}
	log.Printf("services being generated to %s", n)
	//f := util.NewOutputFile(n)
	//f, err := os.Create(filename)
	//if err != nil {
	//	return fmt.Errorf("unable to open %s: %v", filename, err)
	//}
	//defer f.Close()
	//fmt.Fprint(f, "package %s\n", file.GoPackageName)
	////g.P("package " + file.GoPackageName)
	//for _, s := range file.Services {
	//	wasmName := s.GoName
	//	optName := wasmNameFromComment(s.Comments.Leading.String())
	//	if optName != "" {
	//		wasmName = optName
	//	}
	//	w := &wasmService{
	//		ProtoFile:       file.Desc.Path(),
	//		GoPackage:       fmt.Sprint(file.GoPackageName),
	//		WasmServiceName: wasmName,
	//	}
	//	w.Service = s
	//	if err := generateCodeService(f, w, t); err != nil {
	//		return err
	//	}
	//	fmt.Fprint(f, "\n")
	//}
	//
	//filename = file.GeneratedFilenamePrefix + "msg.p.go"
	//log.Printf("services being generated to %s", filename)
	//f, err = os.Create(filename)
	//if err != nil {
	//	return fmt.Errorf("unable to open %s: %v", filename, err)
	//}
	//defer f.Close()
	//
	//for _, msg := range file.Messages {
	//	log.Printf("generating message %s", msg.GoIdent.GoName)
	//	if err := generateCodeMessage(f, msg, t); err != nil {
	//		return err
	//	}
	//	fmt.Fprint(f, "\n")
	//}
	return nil, nil
}

func generateCodeService(w io.WriteCloser, svc *wasmService, tmpl *template.Template) error {
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
