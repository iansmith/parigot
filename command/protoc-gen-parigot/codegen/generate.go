package codegen

import (
	"fmt"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"
	"io"
	"text/template"
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
// into the given GenInfo object.  It walks the full proto file as specified by
// proto.
func Collect(result *GenInfo, lang LanguageText) *GenInfo {

	for _, s := range result.file.GetService() {
		w := &WasmService{
			ServiceDescriptorProto: s,
			parent:                 result.GetFile(),
			lang:                   lang,
		}
		result.RegisterService(w)
		w.method = make([]*WasmMethod, len(s.GetMethod()))
		for j, m := range s.GetMethod() {
			meth := &WasmMethod{
				MethodDescriptorProto: m,
				parent:                w,
				lang:                  lang,
			}
			w.method[j] = meth
		}
	}
	result.message = make(map[*MessageRecord]*WasmMessage, len(result.file.GetMessageType()))
	for _, m := range result.GetFile().GetMessageType() {
		w := NewWasmMessage( /*parent*/ result.GetFile(), m, lang)
		result.RegisterMessage(w)
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
			in := newInputParam(proto.GetPackage(), m.MethodDescriptorProto.GetInputType(), m)
			out := newOutputParam(proto.GetPackage(), m.MethodDescriptorProto.GetOutputType(), m)
			in.lang = lang
			out.lang = lang
			m.input = in
			m.output = out
		}
	}
	return result
}
