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
		// pkgName := splitPackage(toGen)
		// if util.IsSystemLibrary(fname) {
		// 	log.Printf("skipping ")
		// 	continue
		// }
		_, ok := fileSeen[toGen]
		if ok {
			continue
		} else {
			fileSeen[toGen] = struct{}{}
		}
		candPkg := splitPackage(stripProtoToProtoName(toGen))
		if util.IsSystemLibrary(candPkg) {
			continue
		}
		for i, n := range g.TemplateName() {
			//gather imports
			imp := make(map[string]struct{})
			for _, dep := range info.GetFileByName(toGen).GetDependency() {
				candPkg = splitPackage(stripProtoToProtoName(dep))
				if !util.IsSystemLibrary(candPkg) {
					key := "\"" + impToPkg[dep] + "\""
					if key == "\"\"" {
						log.Printf("WARNING: dep produced key that is empty: %s", dep)
					}
				} else {
					continue
				}
			}
			for _, svc := range info.finder.Service() {
				for _, meth := range svc.GetWasmMethod() {
					meth.AddImportsNeeded(imp)
				}
			}
			if !strings.HasSuffix(toGen, ".proto") {
				panic(fmt.Sprintf("unable to understand protocol buffer file with name %s, does not end in .proto", toGen))
			}
			path2 := strings.TrimSuffix(toGen, ".proto") + resultName[i]
			f := util.NewOutputFile(path2)
			wasmService := []*WasmService{}
			nomethod := make(map[string]bool)
			for _, pb := range info.GetAllServiceByName(toGen) {
				desc := info.GetFileByName(toGen)
				w := info.FindServiceByName(desc.GetPackage(), pb.GetName())
				if w == nil {
					panic(fmt.Sprintf("can't find service %s", toGen))
				}
				nomethod[w.GetWasmServiceName()] = w.NoMethod()
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
			noMethodsAtAll := true
			for _, value := range nomethod {
				if !value {
					noMethodsAtAll = false
					break
				}
			}
			data := map[string]interface{}{
				"file":          toGen,
				"req":           info.GetRequest(),
				"info":          info,
				"package":       pkg,
				"import":        imp,
				"service":       wasmService,
				"noMethod":      nomethod,
				"noMethodAtAll": noMethodsAtAll,
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

func addSpecialChars(imp map[string]struct{}, key string) {
	imp["\n"+key+"\n"] = struct{}{}
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
			// if w.ImplementsReverseAPI() != "" {
			// 	log.Printf("xxx implements remote %s: %s", w.GetWasmServiceName(), f)
			// }
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

func stripProtoToProtoName(fullNameAsPath string) string {
	if !strings.HasSuffix(fullNameAsPath, ".proto") {
		panic("unexpected pathname in splitPackage:" + fullNameAsPath)
	}
	p := strings.TrimSuffix(fullNameAsPath, ".proto")
	return strings.ReplaceAll(p, "/", ".")
}

func splitPackage(name string) string {
	part := strings.Split(name, ".")
	if len(part) != 3 {
		panic("unexpected package name:" + name)
	}
	return strings.Join(part[:2], ".")
}
