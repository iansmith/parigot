package lib

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"strconv"
)

// remove leading and trailing semicolon, should never fail because it parsed ok
func typeAnnotationTerminalToInt(node antlr.Token) int {
	text := node.GetText()[1 : len(node.GetText())-1]
	i, e := strconv.Atoi(text)
	if e != nil {
		panic("badly formed number in type annotation") // should never happen
	}
	return i
}

// remove leading and trailing semicolon, should never fail because it parsed ok
func numTerminalToInt(node antlr.Token) int {
	text := node.GetText()
	i, e := strconv.Atoi(text)
	if e != nil {
		panic("badly formed number in Num terminal") // should never happen
	}
	return i
}

// remove leading and trailing quote
func quotedStringTerminalToString(node antlr.Token) string {
	return node.GetText()[1 : len(node.GetText())-1]
}

func stringTerminalToString(node antlr.Token) string {
	return node.GetText()
}

func getTypeNameSeq(r *antlr.BaseParserRuleContext) []string {
	tok := r.GetTokens(WasmLexerTypeName)
	result := make([]string, len(tok))
	for i, t := range tok {
		result[i] = t.GetText()
	}
	return result
}

// VisitTerminal is called when a terminal node is visited.
func (b *Builder) VisitTerminal(node antlr.TerminalNode) {
	//	fmt.Printf("xxx terminal %s\n", node.GetText())
}
