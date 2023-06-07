package codegen

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"text/template"

	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"
)

// var msgRegex = regexp.MustCompile(`(.*\/g\/msg\/)([[:alpha][[:alphanum]]*)(\/v[[:digit]]+)`)
var msgRegex = regexp.MustCompile(`(.*g/msg/)([[:alpha:]][[:alnum:]]*)(/v[0-9]+)$`)

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
		log.Printf("xxx -- processing %s", info.GetFileByName(toGen).GetName())
		for i, n := range g.TemplateName() {
			// if len(info.GetFile().GetService()) == 0 && len(info.GetFile().GetMessageType()) == 0 {
			// 	continue
			// }
			//gather imports
			imp := make(map[string]struct{})
			for _, dep := range info.GetFileByName(toGen).GetDependency() {
				importPrefix := ""
				if msgRegex.MatchString(impToPkg[dep]) {
					match := msgRegex.FindStringSubmatch(impToPkg[dep])
					if len(match) != 4 {
						panic(fmt.Sprintf("unable to understand import match result: %+v", match))
					}
					importPrefix = match[2] + "msg " // convention
				}
				imp[importPrefix+"\""+impToPkg[dep]+"\""] = struct{}{}
			}
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
			enumType := []*WasmEnumType{}
			for _, et := range info.GetAllEnumByName(toGen) {
				desc := info.GetFileByName(toGen)
				w := info.FindEnumTypeByName(desc.GetPackage(), et.GetName())
				if w == nil {
					log.Printf("missed on findEnumByName %s.%s", desc.GetPackage(), et.GetName())
				} else {
					log.Printf("found the enum entry %s.%s", desc.GetPackage(), et.GetName())
				}
			}
			if len(wasmService) == 0 {
				path := info.GetFileByName(toGen).GetName()
				et := info.GetAllEnumByName(path)
				log.Printf("warning: no services found in %s (%d)", path,
					len(et))
				// we don't need to do anything, go plugin will do it
				continue
			}
			pkg, err := info.GoPackageOption(wasmService)
			if err != nil {
				return nil, err
			}
			data := map[string]interface{}{
				"file":     toGen,
				"req":      info.GetRequest(),
				"info":     info,
				"package":  pkg,
				"import":   imp,
				"service":  wasmService,
				"enumType": enumType,
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
	et := result.finder.enumType
	for typeName, e := range et {
		for _, name := range e.ReservedName {
			log.Printf("type %s=>name %s (%s)", typeName, name, e.parent.GetName())
			for _, rg := range e.child {
				log.Printf("value %s=>range [%d,%d]", *rg.Name, rg.start, rg.end)

			}
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
