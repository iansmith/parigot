package helper

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

func AntlrSetupLexParse(path string, lexer *antlr.BaseLexer, parser *antlr.BaseParser) *AntlrErrorListener {
	// load whole input file into memory, convert to input stream
	fp, err := os.Open(path)
	if err != nil {
		wd, _ := os.Getwd()
		AntlrFatalf("%v (wd is %s), %v", flag.Arg(0), wd, err)
	}
	buffer, err := io.ReadAll(fp)
	if err != nil {
		AntlrFatalf("reading %s: %v", flag.Arg(0), err)
	}
	fp.Close()
	input := antlr.NewInputStream(string(buffer))

	// setup the connection between input->lexer->parser
	lexer.SetInputStream(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	parser.SetTokenStream(stream)

	// the diagnostic listener is good for debugging (displays good error msgs)
	el := &AntlrErrorListener{}
	parser.RemoveErrorListeners()
	parser.AddErrorListener(&antlr.DiagnosticErrorListener{
		DefaultErrorListener: &antlr.DefaultErrorListener{},
	})
	parser.AddErrorListener(el)
	return el
}

func AntlrFatalf(spec string, rest ...interface{}) {
	log.Fatalf(spec, rest...)
}

type AntlrErrorListener struct {
	count   int
	failure bool
}

func (a *AntlrErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	if e != nil {
		log.Printf("syntax error %d:%d %s %v (%T)", line, column, msg, offendingSymbol, e)
		a.failure = true
		a.count++
	}
}
func (a *AntlrErrorListener) Failed() bool {
	return a.failure
}

func (el *AntlrErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	//log.Printf("ambiguous alternatives in DFA: %v", ambigAlts)
	el.count++
}

func (a *AntlrErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	//log.Printf("attempting full context in DFA: %v", conflictingAlts)
	a.count++
}
func (a *AntlrErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
	log.Printf("context sensitivity in DFA")
	a.count++
}
