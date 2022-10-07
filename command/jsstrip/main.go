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
var first *bool = flag.Bool("1", false, "1st pass")
var second *bool = flag.Bool("2", false, "2nd pass")
var third *bool = flag.Bool("3", false, "3rd pass")

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

	watVersion := ""
	if allPasses() || *first {
		watVersion, err = convertWasmToWat(tmp, flag.Arg(0))
		if err != nil {
			os.Exit(1)
		}
		if *first {
			log.Printf("(1) %s -> %s\n", flag.Arg(0), watVersion)
			os.Exit(0)
		}
	}
	modifiedWat := ""
	if allPasses() || *second {
		source := watVersion
		if *second {
			source = flag.Arg(0)
		}
		modifiedWat, err = parigotProcessing(source, tmp)
		if err != nil {
			os.Exit(1)
		}
		if *second {
			log.Printf("(2) %s -> %s\n", source, filepath.Join(tmp, parigotFilename))
			os.Exit(0)
		}
	}
	if allPasses() || *third {
		if *third && outputFile == nil {
			panic("cannot handle third pass without output file")
		}
		source := modifiedWat
		if *third {
			source = flag.Arg(0)
		}
		err = convertWatToWasm(tmp, source, outputFile)
		if err != nil {
			os.Exit(1)
		}
		if *third {
			log.Printf("(3) %s -> %s\n", filepath.Join(tmp, parigotFilename), *outputFile)
			os.Exit(0)
		}
	}
	//os.RemoveAll(tmp)
	log.Printf("(1) %s -> %s\n", flag.Arg(0), watVersion)
	log.Printf("(2) %s -> %s\n", watVersion, filepath.Join(tmp, parigotFilename))
	log.Printf("(3) %s -> %s\n", filepath.Join(tmp, parigotFilename), *outputFile)
	os.Exit(0)
}

func allPasses() bool {
	return *first == false && *second == false && *third == false
}

func transformation(mod *transform.Module) {
	//changeToplevelInModule(mod, transform.ImportDefT, changeImportsToPointToStub)
	//changeStatementInModule(mod, changeCallsToUseStub)
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
