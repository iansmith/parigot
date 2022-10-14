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

var generatorMap = map[string]codegen.Generator{}
var abiOnlyMap = map[string]codegen.Generator{}

var save = flag.Bool("s", true, "save a copy of the input to temp dir")
var load = flag.String("l", "", "load a previously saved input (filename)")
var terminal = flag.Bool("t", false, "dump the generated code to stdout instead of using protobuf format")

func main() {
	flag.Parse()

	// info is the root of the all the nodes that we use for code gen
	info := codegen.NewGenInfo()
	generatorMap["go"] = go_.NewGoGen(info)
	abiOnlyMap["abi"] = abi.NewAbiGen(info)

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
	genReq := util.ReadStdinIntoBuffer(source, *save)
	resp := pluginpb.CodeGeneratorResponse{
		Error:             nil,
		SupportedFeatures: nil,
	}
	// generate code, at this point language neutral
	file, err := generateNeutral(info, genReq)

	// set up the response going back out stdout to the protocol buffers compiler
	// response with an error filled in or a set of output files
	if err != nil {
		resp.Error = new(string)
		*resp.Error = err.Error()
	} else {
		resp.File = make([]*pluginpb.CodeGeneratorResponse_File, len(file))
		for i, f := range file {
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
func isGenerate(desc *descriptorpb.FileDescriptorProto, genReq *pluginpb.CodeGeneratorRequest) bool {
	for _, gen := range genReq.GetFileToGenerate() {
		if desc.GetName() == gen {
			return true
		}
	}
	return false
}

// generateNeutral starts the process of code generation and it does not care
// about the languages being used.  those get set at the point we compute genMap
func generateNeutral(info *codegen.GenInfo, genReq *pluginpb.CodeGeneratorRequest) ([]*util.OutputFile, error) {
	fileList := []*util.OutputFile{}
	// walk all the proto files indicated in the request
	for _, desc := range genReq.GetProtoFile() {
		isGenerate := isGenerate(desc, genReq)
		info.SetReqAndFile(genReq, desc) // set to what we are processing now
		genMap := getGeneratorMap(desc)
		// walk all languages, or just the abi if the input turns out to be the abi protos
		for lang, generator := range genMap {
			codegen.Collect(info, generator.LanguageText())
			if isGenerate {
				// skip it? only if no services and no messages xxx will break enums
				if len(info.Service()) == 0 && len(info.Message()) == 0 {
					continue
				}
				// load up templates
				t, err := loadTemplates(generator)
				if err != nil {
					return nil, err
				}
				file, err := generator.Generate(t, info)
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

func getGeneratorMap(desc *descriptorpb.FileDescriptorProto) map[string]codegen.Generator {
	if codegen.IsAbi(desc.GetOptions().String()) {
		return abiOnlyMap // map with JUST the abi generator
	}
	return generatorMap // normal map with one entry per languages

}
