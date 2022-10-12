package codegen

import (
	"fmt"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"
	"io"
	"log"
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
func Collect(request *pluginpb.CodeGeneratorRequest, proto *descriptorpb.FileDescriptorProto) *GenInfo {
	result := &GenInfo{
		request: request,
		file:    proto,
	}

	result.wasmService = make([]*WasmService, len(result.file.GetService()))
	for i, s := range result.file.GetService() {
		w := &WasmService{
			ServiceDescriptorProto: s,
			parent:                 proto,
		}
		result.wasmService[i] = w
		w.method = make([]*WasmMethod, len(s.GetMethod()))
		for j, m := range s.GetMethod() {
			meth := &WasmMethod{
				MethodDescriptorProto: m,
				parent:                w,
			}
			w.method[j] = meth
			//in := newInParameter(result, meth.MethodDescriptorProto.GetInputType())
			//out := newOutResult(result, meth.MethodDescriptorProto.GetOutputType())
			//meth.input = in
			//meth.output = out
		}
	}
	result.wasmMessage = make([]*WasmMessage, len(result.file.GetMessageType()))
	for i, m := range proto.GetMessageType() {
		w := &WasmMessage{
			DescriptorProto: m,
			parent:          proto,
		}
		result.wasmMessage[i] = w
		w.field = make([]*WasmField, len(w.DescriptorProto.GetField()))
		for j, f := range w.DescriptorProto.GetField() {
			field := &WasmField{
				FieldDescriptorProto: f,
				parent:               w,
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
			in := newInParameter(result, m.MethodDescriptorProto.GetInputType())
			out := newOutResult(result, m.MethodDescriptorProto.GetOutputType())
			m.input = in
			m.output = out
		}
	}

	return result
}

func newOutResult(g *GenInfo, messageName string) *OutputParam {
	msg := g.findMessageByName(messageName)
	if msg == nil {
		log.Fatalf("unable to find output parameter type %s", messageName)
	}
	return &OutputParam{
		name: messageName,
		typ:  msg,
	}
}

func newInParameter(g *GenInfo, messageName string) *InputParam {
	msg := g.findMessageByName(messageName)
	if msg == nil {
		log.Fatalf("unable to find input parameter type %s", messageName)
	}
	return &InputParam{
		name: messageName,
		typ:  msg,
	}
}

func (g *GenInfo) findMessageByName(n string) *WasmMessage {
	for _, m := range g.wasmMessage {
		if m.GetFullName() == n {
			return m
		}
	}
	// why do they do this SOME of the time?
	if len(n) > 0 && n[0] == '.' {
		for _, m := range g.wasmMessage {
			if m.GetFullName() == n[1:] {
				return m
			}
		}
	}
	return nil
}
