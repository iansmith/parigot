package antlr

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	v4 "github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

func AntlrSetupLexParse(path string, lexer *v4.BaseLexer, parser *v4.BaseParser) *AntlrErrorListener {
	// load whole input file into memory, convert to input stream
	fp, err := os.Open(path)
	if err != nil {
		wd, _ := os.Getwd()
		AntlrFatalf("%v (working directory is %s), %v", flag.Arg(0), wd, err)
	}
	buffer, err := io.ReadAll(fp)
	if err != nil {
		AntlrFatalf("reading %s: %v", flag.Arg(0), err)
	}
	fp.Close()
	input := v4.NewInputStream(string(buffer))
	// setup the connection between input->lexer->parser
	lexer.SetInputStream(input)
	stream := v4.NewCommonTokenStream(lexer, 0)
	parser.SetTokenStream(stream)

	// the diagnostic listener is good for debugging (displays good error msgs)
	el := &AntlrErrorListener{CurrentFile: path}
	parser.RemoveErrorListeners()
	parser.AddErrorListener(&v4.DiagnosticErrorListener{
		DefaultErrorListener: &v4.DefaultErrorListener{},
	})
	parser.AddErrorListener(el)
	return el
}

func AntlrFatalf(spec string, rest ...interface{}) {
	log.Fatalf(spec, rest...)
}

type AntlrErrorListener struct {
	countSyntaxt int
	countDFA     int
	failure      bool
	CurrentFile  string
}

func (a *AntlrErrorListener) SyntaxError(recognizer v4.Recognizer, offendingSymbol interface{}, line, column int, msg string, e v4.RecognitionException) {
	start := "syntax error "
	if a.CurrentFile != "" {
		start = fmt.Sprintf("syntax error %s:", a.CurrentFile)
	}
	if e != nil {
		log.Printf("%s%d:%d %s %v (%T)", start, line, column, msg, offendingSymbol, e)
		a.failure = true
		a.countSyntaxt++
	}
}
func (a *AntlrErrorListener) Failed() bool {
	return a.failure
}

func (el *AntlrErrorListener) ReportAmbiguity(recognizer v4.Parser, dfa *v4.DFA, startIndex, stopIndex int, exact bool, ambigAlts *v4.BitSet, configs v4.ATNConfigSet) {
	//log.Printf("ambiguous alternatives in DFA: %v", ambigAlts)
	el.countDFA++
}

func (a *AntlrErrorListener) ReportAttemptingFullContext(recognizer v4.Parser, dfa *v4.DFA, startIndex, stopIndex int, conflictingAlts *v4.BitSet, configs v4.ATNConfigSet) {
	//log.Printf("attempting full context in DFA: %v", conflictingAlts)
	a.countDFA++
}
func (a *AntlrErrorListener) ReportContextSensitivity(recognizer v4.Parser, dfa *v4.DFA, startIndex, stopIndex, prediction int, configs v4.ATNConfigSet) {
	//log.Printf("context sensitivity in DFA")
	a.countDFA++
}
