package driver

import (
	"log"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type errorListener struct {
	count int
}

func (el *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	if e != nil {
		log.Printf("syntax error %d:%d %s %v (%T)", line, column, msg, offendingSymbol, e)
		buildSuccess = false
		el.count++
	}
}
func (el *errorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	log.Printf("ambiguous alternatives in DFA: %v", ambigAlts)
	el.count++
}
func (el *errorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	log.Printf("attempting full context in DFA: %v", conflictingAlts)
	el.count++
}
func (el *errorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
	log.Printf("context sensitivity in DFA")
	el.count++
}
