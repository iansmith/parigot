package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/iansmith/parigot/command/protoc-gen-parigot/codegen"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/go_"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var protocVerbose = false

//go:embed template/*
var templateFS embed.FS

var generatorMap = map[string]codegen.Generator{}

var save = flag.Bool("s", true, "save a copy of the input to temp dir")
var load = flag.String("l", "", "load a previously saved input (filename)")
var terminal = flag.Bool("t", false, "dump the generated code to stdout instead of using protobuf format")
var tmpDir = flag.String("d", "tmp", "provide a directory to use as the temp directior, defaults to ./tmp")

func main() {
	flag.Parse()

	// info is the root of the all the nodes that we use for code gen
	info := codegen.NewGenInfo()
	generatorMap["go"] = go_.NewGoGen(info)

	// plugin reads from stdin normally, but we allow input from files and copying
	// the stdin to some output file later testing
	source := io.Reader(os.Stdin)
	if *load != "" {
		*save = false
		fp, err := os.Open(*load)
		if err != nil {
			log.Fatalf("%v", err)
		}
		source = fp
	}
	genReq := util.ReadStdinIntoBuffer(source, *save, *tmpDir, protocVerbose)
	if len(genReq.GetFileToGenerate()) != 1 {
		log.Fatalf("unable to understand the input files, expected 1 file, got %d", len(genReq.GetFileToGenerate()))
	}

	//p := genReq.GetFileToGenerate()[0]
	// if util.IsSystemLibrary(p) {
	// 	log.Printf("skipping system library '%s'", p)
	// 	return
	// }

	importToPackageMap := make(map[string]string)
	for i := range genReq.GetProtoFile() {
		f := genReq.GetProtoFile()[i]
		pkg := f.GetOptions().GetGoPackage()
		if index := strings.LastIndex(pkg, ";"); index != -1 {
			pkg = pkg[:index]
		}
		importToPackageMap[f.GetName()] = pkg
	}
	// compute response
	resp := pluginpb.CodeGeneratorResponse{
		Error:             nil,
		SupportedFeatures: nil,
	}

	// generate code, at this point language neutral
	file, err := generateNeutral(info, genReq, importToPackageMap)

	seen := make(map[string]struct{})
	// clean up the list of files?
	outFile := []*util.OutputFile{}
	for _, f := range file {
		_, foundIt := seen[*f.ToGoogleCGResponseFile().Name]
		if foundIt {
			continue
		}
		seen[*f.ToGoogleCGResponseFile().Name] = struct{}{}
		outFile = append(outFile, f)
	}

	// set up the response going back out stdout to the protocol buffers compiler
	// response with an error filled in or a set of output files
	if err != nil {
		resp.Error = new(string)
		*resp.Error = err.Error()
	} else {
		resp.File = make([]*pluginpb.CodeGeneratorResponse_File, len(outFile))
		for i, f := range outFile {
			resp.File[i] = f.ToGoogleCGResponseFile()
		}
	}

	// send response and exit if not loading from disk
	if *load == "" && !*terminal {
		util.MarshalResponseAndExit(&resp)
	} else {
		util.OutputTerminal(file)
	}

}

// isGenerateTestsIf a specific file given by desc is in the list of files contained
// in the genReq that are required to be generated.  If it is not required to be generated,
// it is being included because it has type information for a dependency of some type
// this is being generated.
func isGenerate(fullName string, genReq *pluginpb.CodeGeneratorRequest) bool {
	for _, gen := range genReq.GetFileToGenerate() {
		if fullName == gen {
			return true
		}
	}
	return false
}

// generateNeutral starts the process of code generation and it does not care
// about the languages being used.  those get set at the point we compute genMap
func generateNeutral(info *codegen.GenInfo, genReq *pluginpb.CodeGeneratorRequest, impToPkg map[string]string) ([]*util.OutputFile, error) {
	fileList := []*util.OutputFile{}
	// compute the set of descriptors that will need to be generated... have to do this firest because
	// there can be multiple protos in the same package
	fileToSvc := make(map[string][]*descriptorpb.ServiceDescriptorProto)
	fileToMsg := make(map[string][]*descriptorpb.DescriptorProto)
	nameToFile := make(map[string]*descriptorpb.FileDescriptorProto)
	enumType := make(map[string][]*descriptorpb.EnumDescriptorProto)
	enumTypeToValue := make(map[string][]*descriptorpb.EnumValueDescriptorProto)

	for _, desc := range genReq.GetProtoFile() {
		// if util.IsIgnoredPackage(desc.GetPackage()) {
		// 	log.Printf("NO NO NO skipping package that is defined elsewhere: %s", desc.GetPackage())
		// }
		nameToFile[desc.GetName()] = desc
		isGen := isGenerate(desc.GetName(), genReq)
		var svc []*descriptorpb.ServiceDescriptorProto
		var ok bool
		if isGen {
			for _, s := range desc.GetService() {
				svc, ok = fileToSvc[desc.GetName()]
				if ok {
					svc = append(svc, s)
				} else {
					svc = []*descriptorpb.ServiceDescriptorProto{s}
				}
				fileToSvc[desc.GetName()] = svc
			}
		}
	}
	for _, desc := range genReq.GetProtoFile() {
		msg := desc.GetMessageType()
		for _, mt := range msg {
			var ok bool
			var allMsg []*descriptorpb.DescriptorProto
			allMsg, ok = fileToMsg[desc.GetName()]
			if ok {
				allMsg = append(allMsg, mt)
			} else {
				allMsg = []*descriptorpb.DescriptorProto{mt}
			}
			fileToMsg[desc.GetName()] = allMsg
		}
	}

	info.SetReqAndFileMappings(genReq, nameToFile, fileToSvc, fileToMsg, enumType, enumTypeToValue)
	// walk all the proto files indicated in the request
	for _, desc := range genReq.GetProtoFile() {
		// if util.IsIgnoredPackage(desc.GetPackage()) {
		// 	continue
		// }
		for lang, generator := range generatorMap {
			codegen.Collect(info, generator.LanguageText())
			if info.Contains(desc.GetName()) {
				// inject this desc into the finder
				nSvc := len(info.GetAllServiceByName(desc.GetName()))
				nMsg := len(info.GetAllMessageByName(desc.GetName()))
				if nSvc == 0 && nMsg == 0 {
					continue
				}
				// load up templates
				t, err := loadTemplates(generator)
				if err != nil {
					return nil, err
				}
				file, err := generator.Generate(t, info, impToPkg)
				if err != nil {
					return nil, err
				}
				if file == nil {
					if protocVerbose {
						log.Printf("warning: language %s did not create any output for %s, ignoring",
							lang, desc.GetName())
					}
				} else {
					fileList = append(fileList, file...)
				}
			} else {
				// process is called when you just might want to look at the types
				// that are being imported (this is also used when we will generate)
				if err := generator.Process(desc); err != nil {
					return nil, err
				}
			}
		}
	}
	return fileList, nil
}

// loadTemplates not only loads the templates proper from the embedded FS but
// also sets up the default functions from codegen.FuncMap.  Generators have
// a chance later to add functions if te want.  The list of files to load
// is from generator.TemplateName() and the list of extra functions is
// generator.FuncMap().
func loadTemplates(generator codegen.Generator) (*template.Template, error) {
	// create root template and add functions, if any
	root := template.New("root")
	root = root.Funcs(codegen.FuncMap)
	// these calls are meant to be "chained" so this construction is needed
	// to capture the "new" value of root.
	if generator.FuncMap() != nil {
		root = root.Funcs(generator.FuncMap())
	}
	t := root
	// template loading
	for _, f := range generator.TemplateName() {
		all, err := templateFS.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("unable to read file %s from embedded fs:%v", f, err)
		}
		t, err = root.New(f).Parse(string(all))
		if err != nil {
			return nil, fmt.Errorf("unable to parse template %s:%v", f, err)
		}
	}
	return t, nil
}
