package codegen

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"
)

// BasicGenerate is the primary code generation driver. Language-specific code
// (in generator.Generate()) is called and that code typically does some setup and
// then calls into this code, passing itself as the g.  This function expects the
// GenInfo to have been created and filled out prior to arriving here.  Because of
// the chaining api, t actually is represents _all_ the templates (not just one) that
// are associated with generator g.   This is called once per .proto file processed.
func BasicGenerate(g Generator, t *template.Template, info *GenInfo, impToPkg map[string]string) ([]*util.OutputFile, error) {
	// run the loop for the templates
	resultName := g.ResultName()
	result := []*util.OutputFile{}
	for _, toGen := range info.request.FileToGenerate {
		for i, n := range g.TemplateName() {
			// if len(info.GetFile().GetService()) == 0 && len(info.GetFile().GetMessageType()) == 0 {
			// 	continue
			// }
			//gather imports
			imp := make(map[string]struct{})
			for _, dep := range info.GetFileByName(toGen).GetDependency() {
				imp["\""+impToPkg[dep]+"\""] = struct{}{}
			}
			//path := util.GenerateOutputFilenameBase(info.GetFile()) + resultName[i]
			if !strings.HasSuffix(toGen, ".proto") {
				panic(fmt.Sprintf("unable to understand protocol buffer file with name %s, does not end in .proto", toGen))
			}
			path2 := strings.TrimSuffix(toGen, ".proto") + resultName[i]
			f := util.NewOutputFile(path2)
			// pbsvc := matchService(toGen, info.request.GetProtoFile())
			// log.Printf("xxx number of pb services for %s: %d", toGen, len(pbsvc))
			// if len(pbsvc) == 1 {
			// 	log.Printf("xxx--> found it in matchService, pbsvc result[0] == %s", pbsvc[0].GetName())
			// }
			wasmService := []*WasmService{}
			for _, pb := range info.GetAllServiceByName(toGen) {
				desc := info.GetFileByName(toGen)
				w := info.FindServiceByName(desc.GetPackage(), pb.GetName())
				if w == nil {
					panic(fmt.Sprintf("can't find service %s", toGen))
				}
				wasmService = append(wasmService, w)
			}
			if len(wasmService) == 0 {
				// we don't need to do anything, go plugin will do it
				continue
			}
			pkg, err := info.GoPackageOption(wasmService)
			if err != nil {
				return nil, err
			}
			data := map[string]interface{}{
				"file":    toGen,
				"req":     info.GetRequest(),
				"info":    info,
				"package": pkg,
				"import":  imp,
				"service": wasmService,
			}
			err = executeTemplate(f, t, n, data)
			if err != nil {
				return nil, fmt.Errorf("failed to process template %s: %v", n, err)
			}
			result = append(result, f)
		}
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
	file := result.GetAllFileName()
	for _, f := range file {
		allSvc := result.GetService(f)
		for _, svc := range allSvc {
			w := NewWasmService(result.GetFileByName(f), svc, lang, result)
			result.RegisterService(w)
		}
	}

	// for _, f:=result.request.FileToGenerate {
	// 	w:=NewWasmService(result.)
	// }
	//we are going to generate the file in result so make sure everything is registered
	// for i, s := range result..GetService() {
	// 	log.Printf("about to create an reg new wasm service [%d]: %s", i, s.GetName())
	// 	w := NewWasmService(result.file, s, lang, result)
	// 	result.RegisterService(w)
	// }
	for _, f := range file {
		for _, m := range result.GetAllMessageByName(f) {
			w := NewWasmMessage(result.GetFileByName(f), m, lang, result)
			result.RegisterMessage(w)
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
			result.AddMessageType(in.GetTypeName(), m.ProtoPackage(), m.GoPackage(), in.CGType().CompositeType())
			result.AddMessageType(out.GetTypeName(), m.ProtoPackage(), m.GoPackage(), out.GetCGType().CompositeType())
			out.lang = lang
			m.input = in
			m.output = out
			m.MarkInputOutputMessages()
		}
	}
	return result
}

// func matchService(toGen string, fd []*descriptorpb.FileDescriptorProto) []*descriptorpb.ServiceDescriptorProto {
// 	for _, protofile := range fd {
// 		fdName := protofile.GetName()
// 		svcs := protofile.GetService()
// 		log.Printf("xxx match service considering match of %s to %s", toGen, fdName)
// 		if toGen == fdName {
// 			return svcs
// 		}
// 	}
// 	return nil
// }
