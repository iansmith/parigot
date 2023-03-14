package antlr

import (
	"log"
	"os"

	v4 "github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/iansmith/parigot/helper"
	"github.com/iansmith/parigot/pbmodel"
	"github.com/iansmith/parigot/ui/parser/tree"
)

func EvaluateOneFile(f, pkg string, b *pbmodel.Pb3Builder) (*tree.ProtobufFileNode, string, bool) {
	l := pbmodel.Newprotobuf3Lexer(nil)
	p := pbmodel.Newprotobuf3Parser(nil)
	el := AntlrSetupLexParse(f, l.BaseLexer, p.BaseParser)
	b.Reset(f)
	proto := p.Proto()
	v4.ParseTreeWalkerDefault.Walk(b, proto)
	if el.Failed() {
		return nil, b.CurrentFile, false
	}
	if b.Failed() {
		return nil, b.CurrentFile, false
	}

	for _, out := range b.OutgoingImport {
		rel := helper.RelativePath(out, f, b.CurrentPackage)
		log.Printf("xxx out2 is %s, rel is %s, file is %s", out, rel, f)
		found := helper.FindProtobufFile(rel, b.CurrentPkgPrefix)
		if found == "" {
			log.Printf("unable to find file '%s' as a protobuf file (relative path to '%s')", rel, f)
			for _, s := range helper.ParigotImportPath() {
				log.Printf("    %s\n", s)
			}
			log.Printf("(maybe you need to check your PARIGOT_IMPORT_PATH?)")
			os.Exit(1)
		}

		_, failedFile, ok := EvaluateOneFile(found, pkg, b)
		if !ok {
			return nil, failedFile, false
		}
	}
	return nil, "", true
}
