package codegen

import (
	"fmt"
	"io"
	"text/template"

	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"
)

// BasicGenerate is the primary code generation driver. Language-specific code
// (in generator.Generate()) is called and that code typically does some setup and
// the calls into this code, passing itself as the g.  This function expects the
// GenInfo to have been created and filled out prior to arriving here.  Because of
// the chaining api, t actually is represents _all_ the templates (not just one) that
// are associated with generator g.   This is called once per .proto file processed.
func BasicGenerate(g Generator, t *template.Template, info *GenInfo) ([]*util.OutputFile, error) {
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
		w := NewWasmService(result.file, s, lang, result)
		result.RegisterService(w)
		w.method = make([]*WasmMethod, len(s.GetMethod()))
		for j, m := range s.GetMethod() {
			w.method[j] = NewWasmMethod(m, w)
		}
	}
	for _, m := range result.GetFile().GetMessageType() {
		w := NewWasmMessage(result.GetFile(), m, lang, result)
		result.RegisterMessage(w)
		w.field = make([]*WasmField, len(w.DescriptorProto.GetField()))
		for j, f := range w.DescriptorProto.GetField() {
			field := NewWasmField(f, w)
			w.field[j] = field
		}
	}
	//
	// Now we have the basic stuff in place we need put things in place that
	// require connections between structures.  Notably, we have to read in
	// all the message types before we can map parameters to them.
	//
	for _, s := range result.Service() {
		for _, m := range s.GetWasmMethod() {
			in := newInputParam(m)
			out := newOutputParam(m)
			out.lang = lang
			m.input = in
			m.output = out
			m.MarkInputOutputMessages()
		}
	}
	return result
}