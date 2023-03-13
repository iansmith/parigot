package driver

import (
	"embed"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/iansmith/parigot/ui/parser"
)

var langToTempl = map[string]string{
	"go": golang,
}

var language = flag.String("l", "go", "pass the name of a known language to get result in that language")
var outputFile = flag.String("o", "", "output file (default is stdout)")
var gopkg = flag.String("gopkg", "main", "golang package code should be generated for")
var invert = flag.Bool("invert", false, "invert the exit error code (useful only for testing)")

var buildSuccess = true

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
	builder, p := readInput(flag.Arg(0))
	prog := p.Program()
	antlr.ParseTreeWalkerDefault.Walk(builder, prog)
	if !buildSuccess {
		wclFatalf("failed due to syntax errors")
	}
	if !parser.NameCheckVisit(prog, builder.ClassName) {
		wclFatalf("failed due to name check")
	}
	execTemplate(prog, *language)
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
	ctx.global["import"] = prog.GetP().ImportSection
	ctx.global["text"] = prog.GetP().TextSection
	ctx.global["doc"] = prog.GetP().DocSection
	ctx.global["event"] = prog.GetP().EventSection
	ctx.global["inputFile"] = flag.Arg(0)
	golang := make(map[string]any)
	ctx.global["golang"] = golang
	golang["package"] = *gopkg
	golang["needBytes"] = prog.GetP().NeedBytes
	golang["needElement"] = prog.GetP().NeedElement
	golang["needEvent"] = prog.GetP().NeedEvent
	// deal with output file
	dir, err := os.MkdirTemp(os.TempDir(), "wcl*")
	if err != nil {
		wclFatalf("unable to create temp dir: %v", err)
	}
	defer func() {
		//log.Printf("cleaning up temp dir %s", dir)
		os.RemoveAll(dir) // clean up
	}()
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

func readInput(path string) (*parser.WclBuildListener, *parser.WCLParser) {
	fp, err := os.Open(path)
	if err != nil {
		wd, _ := os.Getwd()
		wclFatalf("%v (wd is %s), %v", flag.Arg(0), wd, err)
	}
	buffer, err := io.ReadAll(fp)
	if err != nil {
		wclFatalf("reading %s: %v", flag.Arg(0), err)
	}
	fp.Close()
	el := errorListener{0}
	input := antlr.NewInputStream(string(buffer))
	lexer := parser.Newwcllex(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(&el)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.Newwcl(stream)
	p.RemoveErrorListeners()
	// the diagnostic listener is good for debugging (displays good error msgs)
	p.AddErrorListener(&antlr.DiagnosticErrorListener{
		DefaultErrorListener: &antlr.DefaultErrorListener{},
	})
	p.AddErrorListener(&el)
	return parser.NewWclBuildListener(flag.Arg(0)), parser.WCLParserFromWcl(p)

}
