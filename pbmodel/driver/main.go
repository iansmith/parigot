package driver

import (
	"flag"
	"log"

	"github.com/iansmith/parigot/helper/antlr"
	"github.com/iansmith/parigot/pbmodel"
)

func Main() {
	flag.Parse()
	if flag.NArg() != 1 {
		antlr.AntlrFatalf("must supply exactly one filename")
	}
	builder := pbmodel.NewPb3Builder()
	pbfile, bad, ok := antlr.EvaluateOneFile(flag.Arg(0), "", builder)
	if !ok {
		antlr.AntlrFatalf("failed due to parse tree walk of proto file '%s'", bad)
	}
	log.Printf("got pbfile: %#v", pbfile)
}
