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
	*descriptorpb.ServiceDescriptorProto
	WasmServiceName string
	GoPackage       string
	ProtoFile       string
	ProtoPackage    string
}

func (g *GoGen) Generate(t *template.Template, proto *descriptorpb.FileDescriptorProto) ([]*util.OutputFile, error) {
	svcOutName := util.GenerateOutputFilenameBase(proto) + "svc.p.go"
	log.Printf("services being generated to %s", svcOutName)
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
	}
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
	return []*util.OutputFile{f}, nil
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
