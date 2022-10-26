package main

import (
	"flag"
	"github.com/iansmith/parigot/command/transform"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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
var op *string = flag.String("op", "", "name of the operation to perform on the binary")
var dumpStats *bool = flag.Bool("d", false, "dump info about unlink")
var fnName *string = flag.String("f", "", "function name to operate on, for dbgprint")
var replaceFuncNames *string = flag.String("r", "", "function name to operate on, for replacefn or a filename to read a series of replacements from")
var start time.Time

func main() {
	start = time.Now()
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalf("unable to understand arguments, providing an input WASM file is required")
	}

	_, baseArg0 := filepath.Split("surgery")
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
	end := time.Now()
	diff := end.Sub(start)
	//os.RemoveAll(tmp)
	log.Printf("\t surgery: (1) %s -> %s\n", flag.Arg(0), watVersion)
	log.Printf("\t surgery: (2) %s -> %s\n", watVersion, filepath.Join(tmp, parigotFilename))
	log.Printf("\t surgery: (3) %s -> %s\n", filepath.Join(tmp, parigotFilename), *outputFile)

	if *outputFile != "" {
		statOut, err := os.Stat(*outputFile)
		if err != nil {
			log.Printf("unable to stat output file %s: %v", *outputFile, err)
		}
		statIn, err := os.Stat(flag.Arg(0))
		if err != nil {
			log.Printf("unable to stat input file %s: %v", flag.Arg(0), err)
		}
		log.Printf("\t surgery: input file: %s, output file: %s", sizeInBytes(statIn.Size()),
			sizeInBytes(statOut.Size()))
		s := diff.Seconds()
		log.Printf("\t surgery: total time %02.1f seconds", s)
	}

	os.Exit(0)
}

func allPasses() bool {
	return *first == false && *second == false && *third == false
}

func transformation(mod *transform.Module) {
	if *op == "" {
		log.Fatalf("you need to supply the -op parameter to express which op to perform")
	}
	switch {
	case *op == "replacefn":
		if *replaceFuncNames == "" {
			log.Fatalf("operation replacefn needs two functions to operate on, use the -r option and comma separate")
		}
		parts := strings.Split(*replaceFuncNames, ",")
		var fp *os.File
		var repl *replaceFn
		if len(parts) == 1 {
			var err error
			log.Printf("assuming that the -r parameter is a filename: '%s'", *replaceFuncNames)
			fp, err = os.Open(*replaceFuncNames)
			if err != nil {
				log.Fatalf("operation replacefn needs two functions to operate on,"+
					"use the -r option and comma separate, could not understand '%s'", *replaceFuncNames)
			}
			repl = newReplaceFn("", "", fp)
		} else {
			repl = newReplaceFn(parts[0], parts[1], nil)
		}
		repl.validateFunctions(mod)
		changeStatementInModule(mod, repl.replace)
	//case "old":
	//changeToplevelInModule(mod, transform.ImportDefT, )
	//changeStatementInModule(mod, changeCallsToUseTrueABI)
	case *op == "unlink":
		u := newUnlink()
		findToplevelInModule(mod, transform.ImportDefT, u.compileImports)
		findToplevelInModule(mod, transform.TypeDefT, u.compileFuncTypes)
		findToplevelInModule(mod, transform.FuncDefT, u.compileCandidates)
		u.addImportsForCandidates(mod)
		changeToplevelInModule(mod, transform.FuncDefT, u.unlinkStdlib)
		if *dumpStats {
			u.dumpStats()
		}
	default:
		log.Fatalf("unknow op %s", *op)
	}
	if *fnName != "" {
		dbg := newDbgPrint(*fnName)
		log.Printf("instrumenting function '%s'", *fnName)
		dbg.driver(mod)
	}
}

func changeImportsToTrueABI(tl transform.TopLevel) transform.TopLevel {
	idef := tl.(*transform.ImportDef)
	if idef.ModuleName == parigotModule {
		*idef.FuncNameRef.Name = *idef.FuncNameRef.Name + parigotSuffix
	}
	return idef
}

func changeCallsToUseTrueABI(stmt transform.Stmt) transform.Stmt {
	if stmt.StmtType() == transform.OpStmtT &&
		stmt.(transform.Op).OpType() == transform.CallT {
		if strings.HasPrefix(stmt.(*transform.CallOp).Arg, abiGoLinkerPath) {
			stmt.(*transform.CallOp).Arg = stmt.(*transform.CallOp).Arg + parigotSuffix
		}
	}

	return stmt
}
