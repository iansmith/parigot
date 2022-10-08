package main

import (
	"embed"
	_ "embed"
	"flag"
	structure "github.com/iansmith/parigot/command/toml"
)

//go:embed template/*
var templateFS embed.FS

const abiTomlPath = "abi/atlanta1/base/go/parigot/abi/proto/gen/abi/abi.p.toml"
const abiTargetPath = "abi/atlanta1/base/go/parigot/abi"

var abiProject = &structure.ProjectDecl{
	Dir:           "",
	MarkerFile:    "",
	GoModule:      "abi",
	ServicesFound: nil,
}

func main() {
	flag.Parse()

	conf := structure.ParseTomlFile(abiTomlPath, abiProject)
	for _, tomlDecl := range conf.Service {
		tomlDecl.TargetDir = abiTargetPath
		generateCode(&tomlDecl, abiProject)
	}
}
