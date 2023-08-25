package codegen

import (
	"fmt"
	"io"
	"log"
	"strings"
	"text/template"

	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"
)

var fileSeen = make(map[string]struct{})

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
		fname := info.GetFileByName(toGen).GetName()
		_, ok := fileSeen[fname]
		if ok {
			continue
		} else {
			fileSeen[fname] = struct{}{}
		}
		for i, n := range g.TemplateName() {
			//gather imports
			imp := make(map[string]struct{})
			for _, dep := range info.GetFileByName(toGen).GetDependency() {
				if !IsIgnoredPackage(dep) {
					imp["\""+impToPkg[dep]+"\""] = struct{}{}
				}
			}
			for _, svc := range info.finder.Service() {
				for _, meth := range svc.GetWasmMethod() {
					meth.AddImportsNeeded(imp)
				}
			}
			log.Printf("we have a map %s (%d,%s) %+v", toGen, i, n, imp)
			if !strings.HasSuffix(toGen, ".proto") {
				panic(fmt.Sprintf("unable to understand protocol buffer file with name %s, does not end in .proto", toGen))
			}
			path2 := strings.TrimSuffix(toGen, ".proto") + resultName[i]
			f := util.NewOutputFile(path2)
			wasmService := []*WasmService{}

			for _, pb := range info.GetAllServiceByName(toGen) {
				desc := info.GetFileByName(toGen)
				w := info.FindServiceByName(desc.GetPackage(), pb.GetName())
				if w == nil {
					panic(fmt.Sprintf("can't find service %s", toGen))
				}
				wasmService = append(wasmService, w)
			}
			wasmMessage := []*WasmMessage{}
			for _, pb := range info.GetAllMessageByName(toGen) {
				desc := info.GetFileByName(toGen)
				pbPkg := desc.GetPackage()
				w := info.FindMessageByName(pbPkg, pb.GetName())
				if w == nil {
					panic(fmt.Sprintf("can't find service %s", toGen))
				}
				wasmMessage = append(wasmMessage, w)
			}

			if len(wasmService) == 0 {
				continue
			}
			pkg, err := info.GoPackageOption(wasmService, wasmMessage)
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

func IsIgnoredPackage(s string) bool {
	switch s {
	case "google.golang.org/protobuf/types/known/anypb",
		"google.golang.org/protobuf/types/known/timestamppb",
		"github.com/iansmith/parigot/g/protosupport/v1":
		return false
	}
	return true
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
			w := NewWasmService(result.GetFileByName(f), svc, lang, result.finder)
			result.RegisterService(w)
		}
	}
	for _, f := range file {
		for _, m := range result.GetAllMessageByName(f) {
			w := NewWasmMessage(result.GetFileByName(f), m, lang, result.finder)
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

			hostFuncName := ""
			opt, ok := isStringOptionPresent(m.GetOptions().String(), parigotOptionForHostFuncName)
			if ok {
				hostFuncName = strings.ReplaceAll(opt, "\"", "")
			}
			result.AddMessageType(in.GetTypeName(), m.ProtoPackage(), m.GoPackage(), in.CGType().CompositeType())
			result.AddMessageType(out.GetTypeName(), m.ProtoPackage(), m.GoPackage(), out.GetCGType().CompositeType())
			out.lang = lang
			m.input = in
			m.output = out
			m.HostFuncName = hostFuncName
			m.MarkInputOutputMessages()
		}
	}
	return result
}
