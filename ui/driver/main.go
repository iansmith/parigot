package driver

import (
	"embed"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	v4 "github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/iansmith/parigot/helper/antlr"
	"github.com/iansmith/parigot/ui/parser"
	"github.com/iansmith/parigot/ui/parser/tree"
)

var langToTempl = map[string]string{
	"go": golang,
}

var language = flag.String("l", "go", "pass the name of a known language to get result in that language")
var outputFile = flag.String("o", "", "output file (default is stdout)")
var gopkg = flag.String("gopkg", "main", "golang package code should be generated for")
var invert = flag.Bool("invert", false, "invert the exit error code (useful only for testing)")

//var buildSuccess = true

//go:embed template/*
var templateFS embed.FS

var exitOk = 0
var exitError = 1

func Main() {
	flag.Parse()
	if flag.NArg() == 0 {
		wclFatalf("provide a filename that contains a web coordination language program")
	}
	// tests use this
	if *invert {
		exitOk = 1
		exitError = 0
	}
	inFile := flag.Arg(0)

	// primary objects
	l := parser.Newwcllex(nil)
	p := parser.Newwcl(nil)
	b := parser.NewWclBuildListener(inFile)
	// antlr setup machinery
	el := antlr.AntlrSetupLexParse(inFile, l.BaseLexer, p.BaseParser)

	// start parsing
	prog := p.Program()
	p.AddParseListener(b)
	v4.ParseTreeWalkerDefault.Walk(b, prog)
	if el.Failed() {
		wclFatalf("failed due to syntax errors")
	}

	// need to clean up some pointers and such
	if err := tree.GProgram.FinalizeSemantics(inFile); err != nil {
		wclFatalf("failed due to semantic checks: %v", err)
	}

	if !tree.GProgram.VarCheck(inFile) {
		wclFatalf("failed check of variables and functions")
	}

	execTemplate(prog, *language)

	// topo, err := graph.TopologicalSort(pbmodel.Pb3Dep)
	// if err != nil {
	// 	log.Fatalf("unable to topologically sort the PB3 dependencies!")
	// }
	// log.Printf("TOPO SORT")
	// for _, t := range topo {
	// 	log.Printf(t)
	// }
	os.Exit(exitOk)
}

func wclFatalf(spec string, rest ...interface{}) {
	wclPrintf(spec, rest...)
	os.Exit(exitError)
}

func wclPrintf(spec string, rest ...interface{}) {
	log.Printf(spec, rest...)
}

func execTemplate(prog parser.IProgramContext, lang string) {
	// create a context for this generate
	t, ok := langToTempl[lang]
	if !ok {
		wclFatalf("unable to find a template for language '%s'", *language)
	}
	ctx := newGenerateContext(t)
	ctx.program = prog.GetP()
	ctx.templateName = t
	ctx.global["import"] = tree.GProgram.ImportSection
	ctx.global["text"] = tree.GProgram.TextSection
	ctx.global["doc"] = tree.GProgram.DocSection
	ctx.global["event"] = tree.GProgram.EventSection
	if tree.GProgram.ModelSection != nil {
		ctx.global["controller"] = tree.GProgram.ModelSection.ControllerDecl
	}
	ctx.global["inputFile"] = flag.Arg(0)
	golang := make(map[string]any)
	ctx.global["golang"] = golang
	golang["package"] = *gopkg
	dir, err := os.MkdirTemp(os.TempDir(), "wcl*")
	if err != nil {
		wclFatalf("unable to create temp dir: %v", err)
	}
	defer func() {
		//log.Printf("cleaning up temp dir %s", dir)
		os.RemoveAll(dir) // clean up
	}()
	//log.Printf("output file is %s\n", filepath.Join(dir, "output_program.go"))
	file := filepath.Join(dir, "output_program.go")
	fp, err := os.Create(file)
	if err != nil {
		wclFatalf("unable to create output file: %v", err)
	}
	err = runTemplate(ctx, fp)
	if err != nil {
		wclFatalf("error trying to execute template %s: %v", ctx.templateName, err)
	}

	cmd := exec.Command("gofmt", file)
	var outFp io.Writer
	if *outputFile != "" {
		outFp, err = os.Create(*outputFile)
		if err != nil {
			wclFatalf("unable to create output file %s: %v", *outputFile, err)
		}
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if outFp != nil {
		cmd.Stdout = outFp
	}
	err = cmd.Run()
	if err != nil {
		wclFatalf("failed to run gofmt: %v, %v", err, outFp)
	}
}
