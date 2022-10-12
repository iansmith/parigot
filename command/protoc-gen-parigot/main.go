package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/iansmith/parigot/command/protoc-gen-parigot/abi"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/codegen"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/go_"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/util"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

//go:embed template/*
var templateFS embed.FS

var GeneratorMap = map[string]codegen.Generator{
	"go": &go_.GoGen{},
}

var AbiOnlyMap = map[string]codegen.Generator{
	"abi": &abi.AbiGen{},
}

var save = flag.Bool("s", true, "save a copy of the input to temp dir")
var load = flag.String("l", "", "load a previously saved input (filename)")

func main() {
	flag.Parse()

	source := io.Reader(os.Stdin)
	if *load != "" {
		*save = false
		fp, err := os.Open(*load)
		if err != nil {
			log.Fatalf("%v", err)
		}
		source = fp
	}
	genReq := util.ReadStdinIntoBuffer(source, *save)
	resp := pluginpb.CodeGeneratorResponse{
		Error:             nil,
		SupportedFeatures: nil,
	}
	// generate code, at this point language neutral
	files, err := generateNeutral(genReq, util.LocatorNames(genReq.GetParameter()))

	// set up the response going back out stdout to the protocol buffers compiler
	if err != nil {
		resp.Error = new(string)
		*resp.Error = err.Error()
	} else {
		resp.File = make([]*pluginpb.CodeGeneratorResponse_File, len(files))
		for i, f := range files {
			resp.File[i] = f.ToGoogleCGResponseFile()
		}
	}

	// send response and exit
	util.MarshalResponseAndExit(&resp)
}

// isGenerateTestsIf a specific file given by desc is in the list of files contained
// in the genReq that are required to be generated.  If it is not required to be generated,
// it is being included because it has type information for a dependency of some type
// this is being generated.
func isGenerate(desc *descriptorpb.FileDescriptorProto, genReq *pluginpb.CodeGeneratorRequest) bool {
	for _, gen := range genReq.GetFileToGenerate() {
		if desc.GetName() == gen {
			return true
		}
	}
	return false
}

// generateNeutral starts the process of code generation and it does not care
// about the languages being used.
func generateNeutral(genReq *pluginpb.CodeGeneratorRequest, locators []string) ([]*util.OutputFile, error) {
	fileList := []*util.OutputFile{}
	for _, desc := range genReq.GetProtoFile() {
		info := codegen.Collect(genReq, desc)
		isGenerate := isGenerate(desc, genReq)
		isAbi := info.IsAbi()
		// here is the trick where we pull the switcheroo for code marked as part
		// of the ABI... we call that a special "language"
		langMap := GeneratorMap
		if isGenerate && isAbi {
			langMap = AbiOnlyMap
		} else {
			// xxx fix me DANGER
			continue
		}
		// walk all languages, or just the abi
		for lang, generator := range langMap {
			if isGenerate {
				// load up templates
				t, err := loadTemplates(generator)
				if err != nil {
					return nil, err
				}
				file, err := generator.Generate(t, info, []string{})
				if err != nil {
					return nil, err
				}
				if file == nil {
					log.Printf("warning: language %s did not create any output for %s, ignoring",
						lang, desc.GetName())
				} else {
					fileList = append(fileList, file...)
				}
			} else {
				// process is called when you just might want to look at the types
				// that are being imported
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
