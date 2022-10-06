package main

import (
	"flag"
	"github.com/iansmith/parigot/command/transform"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	parigotModule   = "parigot_abi"
	abiGoLinkerPath = "$github.com/iansmith/parigot/abi/go"
	parigotSuffix   = "_Stub"
)

// command line args
var outputFile *string = flag.String("o", "", "set to the output file you want to produce, otherwise output goes to stdout")

func main() {

	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalf("unable to understand arguments, providing an input WASM file is required")
	}

	_, baseArg0 := filepath.Split("jsstrip")
	tmp, err := os.MkdirTemp("", baseArg0)
	if err != nil {
		log.Fatalf("cannot create temp dir: %v", err)
	}

	watVersion, err := convertWasmToWat(tmp, flag.Arg(0))
	if err != nil {
		os.Exit(1)
	}
	modifiedWat, err := parigotProcessing(watVersion, tmp)
	if err != nil {
		os.Exit(1)
	}
	err = convertWatToWasm(tmp, modifiedWat, outputFile)
	if err != nil {
		os.Exit(1)
	}
	//os.RemoveAll(tmp)
	log.Printf("%s -> %s\n", os.Args[1], watVersion)
	log.Printf("%s -> %s\n", watVersion, filepath.Join(tmp, parigotFilename))
	log.Printf("%s -> %s\n", filepath.Join(tmp, parigotFilename), os.Args[2])
	os.Exit(0)
}

var funcDefCount int

func transformation(mod *transform.Module) {
	changeToplevelInModule(mod, transform.ImportDefT, changeImportsToPointToStub)
	changeStatementInModule(mod, changeCallsToUseStub)
}

func changeImportsToPointToStub(tl transform.TopLevel) transform.TopLevel {
	idef := tl.(*transform.ImportDef)
	if idef.ModuleName == parigotModule {
		idef.FuncNameRef.Name = idef.FuncNameRef.Name + parigotSuffix
	}
	return idef
}

func changeCallsToUseStub(stmt transform.Stmt) transform.Stmt {
	if stmt.StmtType() == transform.OpStmtT &&
		stmt.(transform.Op).OpType() == transform.CallT {
		if strings.HasPrefix(stmt.(*transform.CallOp).Arg, abiGoLinkerPath) {
			stmt.(*transform.CallOp).Arg = stmt.(*transform.CallOp).Arg + parigotSuffix
		}
	}

	return stmt
}
