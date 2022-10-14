package codegen

import (
	"fmt"
	"io"
	"log"
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
			in := newInputParam(result.GetFile().GetPackage(), m.GetInputType(), m)
			in.lang = lang
			out := newOutputParam(result.GetFile().GetPackage(), m.GetOutputType(), m, result)
			out.lang = lang
			m.input = in
			m.output = out
		}
	}
	debugDump(result, func(param *CGParameter) string {
		return param.String(result.GetFile().GetPackage()) + "\n"
	})
	return result
}

func debugDump(result *GenInfo, fn func(parameter *CGParameter) string) {
	for _, svc := range result.Service() {
		protoPackage := svc.GetParent().GetPackage()
		for _, method := range svc.GetWasmMethod() {
			in := method.GetInputParam()
			out := method.GetOutputParam()

			// check the types to make sure they are not there
			if in.GetCGType() == nil {
				in.cgType = GetCGTypeForInputParam(in)
			}
			if out.GetCGType() == nil {
				out.cgType = GetCGTypeForOutputParam(out)
			}
			inParam := NewCGParameterFromString(in.GetCGType())
			_ = NewCGParameterFromString(out.GetCGType())

			call := ""
			call += fmt.Sprintf("%s(", method.GetName())
			if !in.IsEmpty() {
				call += inParam.GetCGType().String(protoPackage)
			}
			call += fmt.Sprintf(") ")
			if !out.IsEmpty() {
				call += out.GetCGType().String(protoPackage)
			}
			if !in.IsEmpty() {
				if !in.GetCGType().IsBasic() {
					comp := in.GetCGType().GetCompositeType()
					if len(comp.GetField()) != 0 {
						call += fmt.Sprintf("\n\t---> %s\n", comp.GetName())
						for i, f := range comp.GetField() {
							log.Printf("i is %d", i)
							cgt := NewCGTypeFromBasic(f.GetType().String(), svc.GetLanguage(), svc.GetFinder(), protoPackage)
							cgp := NewCGParameterFromField(f, cgt)
							call += fn(cgp)
						}
					}
				}
			}
			log.Printf("%s", call)
		}
	}
}
