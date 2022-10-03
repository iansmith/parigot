package lib

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"strconv"
	"strings"
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

// remove leading and trailing ; and the @
func blockAnnotationTerminalToInt(node antlr.Token) int {
	text := node.GetText()[2 : len(node.GetText())-1]
	i, e := strconv.Atoi(text)
	if e != nil {
		panic("badly formed number in block annotation") // should never happen
	}
	return i
}

// remove leading and trailing ; and the =
func constAnnotationTerminalToString(node antlr.Token) string {
	return node.GetText()[2 : len(node.GetText())-1]
}

// remove annotation cruft
func targetTerminalToBranchTarget(node antlr.Token) *BranchTarget {
	parts := strings.Split(node.GetText(), " ")
	if len(parts) != 2 {
		panic("unable to understand branch target") //should never happen
	}
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		panic("unable to understand branch target field num") //should never happen
	}
	b, err := strconv.Atoi(parts[1][2 : len(parts[1])-1])
	if err != nil {
		panic("unable to understand branch target field block") //should never happen
	}
	return &BranchTarget{Num: n, Block: b}
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

// just return the text
func identTerminalToString(node antlr.Token) string {
	return node.GetText()
}

func offsetTerminalToInt(node antlr.Token) int {
	return alignmentOrOffset(node, "offset=")
}

func alignTerminalToInt(node antlr.Token) int {
	return alignmentOrOffset(node, "align=")
}

func alignmentOrOffset(node antlr.Token, prefix string) int {
	t := node.GetText()[len(prefix):]
	i, err := strconv.Atoi(t)
	if err != nil {
		panic("unable to understand alignment value") // should never happen
	}
	return i
}

func stringTerminalToString(node antlr.Token) string {
	return node.GetText()
}

func getTypeNameSeq(r *antlr.BaseParserRuleContext) []string {
	tok := r.GetTokens(WasmLexerTypeName)
	r.GetStart()
	result := make([]string, len(tok))
	for i, t := range tok {
		result[i] = t.GetText()
	}
	return result
}

func TokenToInt(node antlr.Token) int {
	i, e := strconv.Atoi(node.GetText())
	if e != nil {
		panic("unable to understand token for int conversion") // should never happen
	}
	return i
}

// VisitTerminal is called when a terminal node is visited.
func (b *Builder) VisitTerminal(node antlr.TerminalNode) {
	//	fmt.Printf("xxx terminal %s\n", node.GetText())
}
