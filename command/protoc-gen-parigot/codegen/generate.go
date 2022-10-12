package codegen

import (
	"fmt"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"
	"io"
	"text/template"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// BasicGenerate is the primary code generation driver. Language-specific code
// (in generator.Generate()) is called and that code typically does some setup and
// the calls into this code, passing itself as the g.  This function expects the
// GenInfo to have been created and filled out prior to arriving here.  Because of
// the chaining api, t actually is represents _all_ the templates (not just one) that
// are associated with generator g.   This is called once per .proto file processed.
func BasicGenerate(g Generator, t *template.Template, info *GenInfo, loc []string) ([]*util.OutputFile, error) {
	// run the loop for the templates
	resultName := g.ResultName()
	result := []*util.OutputFile{}
	for i, n := range g.TemplateName() {
		path := util.GenerateOutputFilenameBase(info.GetFile()) + resultName[i]
		f := util.NewOutputFile(path)
		data := map[string]interface{}{
			"file": info.GetFile(),
			"req":  info.GetRequest(),
			"info": info,
		}
		err := executeTemplate(f, t, n, data)
		if err != nil {
			return nil, fmt.Errorf("failed to process template %s: %v", n, err)
		}
		result = append(result, f)
	}
	return result, nil

}

// executeTemplate actually runs a template and makes sure errors are collected.
// The w value in most cases is an instance of *util.Outfile.
func executeTemplate(w io.Writer, t *template.Template, name string, data map[string]interface{}) error {
	if t.Lookup(name) == nil {
		return fmt.Errorf("unable to find teplate: %s", name)
	}
	return t.Lookup(name).Execute(w, data)
}

// Collect is called to gather all the relevant information about the proto files
// into the resultingGenInfo object.  It walks the full proto file spcified by
// proto.
func Collect(request *pluginpb.CodeGeneratorRequest, proto *descriptorpb.FileDescriptorProto,
	lang LanguageText, abilang ABIText) *GenInfo {
	result := &GenInfo{
		request: request,
		file:    proto,
		lang:    lang,
	}

	result.wasmService = make([]*WasmService, len(result.file.GetService()))
	for i, s := range result.file.GetService() {
		w := &WasmService{
			ServiceDescriptorProto: s,
			parent:                 proto,
			lang:                   lang,
		}
		result.wasmService[i] = w
		w.method = make([]*WasmMethod, len(s.GetMethod()))
		for j, m := range s.GetMethod() {
			meth := &WasmMethod{
				MethodDescriptorProto: m,
				parent:                w,
				lang:                  lang,
				abilang:               abilang,
			}
			w.method[j] = meth
		}
	}
	result.wasmMessage = make([]*WasmMessage, len(result.file.GetMessageType()))
	for i, m := range proto.GetMessageType() {
		w := &WasmMessage{
			DescriptorProto: m,
			parent:          proto,
			lang:            lang,
		}
		result.wasmMessage[i] = w
		w.field = make([]*WasmField, len(w.DescriptorProto.GetField()))
		for j, f := range w.DescriptorProto.GetField() {
			field := &WasmField{
				FieldDescriptorProto: f,
				parent:               w,
				lang:                 lang,
			}
			w.field[j] = field
		}
	}
	//
	// Now we have the basic stuff in place we need put things in place that
	// require connections between structures.
	//
	for _, s := range result.GetWasmService() {
		for _, m := range s.GetWasmMethod() {
			in := newInputParam(result, m.MethodDescriptorProto.GetInputType())
			out := newOutputParam(result, m.MethodDescriptorProto.GetOutputType())
			in.lang = lang
			out.lang = lang
			m.input = in
			m.output = out
		}
	}

	return result
}
