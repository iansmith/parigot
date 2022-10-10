package abi

import (
	"log"
	"text/template"

	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	tinygo = "template/abi/abiTinygo.tmpl"
	ide    = "template/abi/abiIde.tmpl"
)

type AbiGen struct {
	types map[string]*descriptorpb.FileDescriptorProto
}

func (a *AbiGen) addType(name string, fdp *descriptorpb.FileDescriptorProto) {
	if a.types == nil {
		a.types = make(map[string]*descriptorpb.FileDescriptorProto)
	}
	a.types[name] = fdp
}

func (a *AbiGen) Process(proto *descriptorpb.FileDescriptorProto) error {
	a.addType(proto.GetName(), proto)
	return nil
}
func (a *AbiGen) Generate(t *template.Template, proto *descriptorpb.FileDescriptorProto) ([]*util.OutputFile, error) {
	ideOutName := util.GenerateOutputFilenameBase(proto) + "null.p.go"
	//xxxfix me
	globalDescriptor = proto
	log.Printf("fake abi functions being generated to %s", ideOutName)
	f0 := util.NewOutputFile(ideOutName)
	if len(proto.Service) != 1 {
		panic("unexpected service definitions (?) for the abi layer")
	}
	data := map[string]interface{}{
		"abi": proto.Service[0],
	}
	if err := t.Lookup(ide).Execute(f0, data); err != nil {
		log.Fatalf("failed to execute template %s for ide: %v", ide, err)
	}
	tinygoOutName := util.GenerateOutputFilenameBase(proto) + ".p.go"
	log.Printf("abi functions being generated to %s", tinygoOutName)

	f1 := util.NewOutputFile(tinygoOutName)
	if err := t.Lookup(tinygo).Execute(f1, data); err != nil {
		log.Fatalf("failed to execute template %s for ide: %v", tinygoOutName, err)
	}

	return []*util.OutputFile{f0, f1}, nil
}

func (a *AbiGen) TemplateName() []string {
	return []string{tinygo, ide}

}
func (a *AbiGen) FuncMap() template.FuncMap {
	return abiFuncMap
}
