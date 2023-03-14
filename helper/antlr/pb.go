package antlr

import (
	"log"
	"os"

	v4 "github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/iansmith/parigot/helper"
	"github.com/iansmith/parigot/pbmodel"
	"github.com/iansmith/parigot/ui/parser/tree"
)

func EvaluateOneFile(f string, b *pbmodel.Pb3Builder) (*tree.ProtobufFileNode, string, bool) {
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
		found := helper.FindProtobufFile(out, b.CurrentPkgPrefix)
		if found == "" {
			for _, s := range helper.ParigotImportPath() {
				log.Printf("    %s\n", s)
			}
			log.Printf("(maybe you need to check your PARIGOT_IMPORT_PATH?)")
			os.Exit(1)
		}

		_, failedFile, ok := EvaluateOneFile(found, b)
		if !ok {
			return nil, failedFile, false
		}
	}
	return nil, "", true
}
