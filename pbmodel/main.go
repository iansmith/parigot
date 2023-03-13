package pbmodel

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/iansmith/parigot/helper"
)

func pbModelFatalf(spec string, rest ...interface{}) {
	log.Fatalf(spec, rest...)
}
func Main() {
	flag.Parse()
	if flag.NArg() != 1 {
		pbModelFatalf("must supply exactly one filename")
	}
	builder := NewPb3Builder()

	bad, ok := evaluateOneFile(flag.Arg(0), builder)
	if !ok {
		helper.AntlrFatalf("failed due to parse tree walk of proto file '%s'", bad)
	}
}

func evaluateOneFile(f string, b *Pb3Builder) (string, bool) {
	l := Newprotobuf3Lexer(nil)
	p := Newprotobuf3Parser(nil)
	el := helper.AntlrSetupLexParse(f, l.BaseLexer, p.BaseParser)
	b.Reset(f)
	proto := p.Proto()
	antlr.ParseTreeWalkerDefault.Walk(b, proto)
	if el.Failed() {
		return b.currentFile, false
	}
	if b.Failed() {
		return b.currentFile, false
	}
	for _, out := range b.OutgoingImport {
		checked := out
		found := ""
		if len(checked) > 0 && ((checked[0] != '/') || (checked[0] != '\\')) {
			currentPlusImportPath := append([]string{b.currentPkgPrefix}, parigotImportPath()...)
			currentPlusImportPath = append(currentPlusImportPath, "")
			for _, candidate := range currentPlusImportPath {
				path := filepath.Join(candidate, checked)
				_, err := os.Stat(path)
				if err != nil {
					if os.IsNotExist(err) {
						continue
					}
				}
				found = path
				break
			}
		} else {
			// fully qualified path
			_, err := os.Stat(f)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
			}
			found = f
		}
		if found == "" {
			for _, s := range parigotImportPath() {
				log.Printf("    %s\n", s)
			}
			log.Printf("(maybe you need to check your PARIGOT_IMPORT_PATH?)")
			os.Exit(1)
		}

		failedFile, ok := evaluateOneFile(found, b)
		if !ok {
			return failedFile, false
		}
	}
	return "", true
}

const parigotImportVar = "PARIGOT_IMPORT_PATH"

var cachedImportPath []string

func parigotImportPath() []string {
	if cachedImportPath != nil {
		return cachedImportPath
	}
	raw, ok := os.LookupEnv(parigotImportVar)
	if !ok {
		raw, _ = os.Getwd()
	}
	cachedImportPath = filepath.SplitList(raw)
	return cachedImportPath
}
