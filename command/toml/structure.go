package toml

import (
	"io"
	"log"
	"os"
	"strings"

	toml "github.com/pelletier/go-toml"
)

type ProjectDecl struct {
	Dir           string                  // path to the dir containing go.mod for go projects
	MarkerFile    string                  // path to the go.mod indicating project root
	GoModule      string                  // name of the module inside the module file
	ServicesFound map[string]*ServiceDecl // converts a path to toml file into its contents
}

type ServiceDecl struct {
	Name            string // service name
	ProtoFile       string // path to the proto file
	TomlFile        string // path of the toml file we generated
	WasmServiceName string // optional comment on the service name
	GoPackage       string // listed as "go_package" in the proto file
	TargetDir       string // where to place generated output, this is the directory with the p.toml
	Method          map[string]MethodDecl
}
type MethodDecl struct {
	Name   string
	Input  string
	Output string
}

type TomlConfig struct {
	Service map[string]ServiceDecl
}

func ParseTomlFile(tomlPath string, project *ProjectDecl) *TomlConfig {
	fp, err := os.Open(tomlPath)
	if err != nil {
		log.Fatalf("unable to open %s:%v", tomlPath, err)
	}
	defer fp.Close()
	b, err := io.ReadAll(fp)
	if err != nil {
		log.Fatalf("unable to read %s:%v", tomlPath, err)
	}
	var conf TomlConfig
	dec := toml.NewDecoder(strings.NewReader(string(b)))
	dec.Strict(true)
	err = dec.Decode(&conf)
	if err != nil {
		log.Fatalf("unable to understand toml file %s:%v", tomlPath, err)
	}
	return &conf
}
