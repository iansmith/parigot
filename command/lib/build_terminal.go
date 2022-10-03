package lib

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"strconv"
)

func numToInt(text string) int {
	i, e := strconv.Atoi(text)
	if e != nil {
		panic("badly formed number") // should never happen
	}
	return i
}

func annoToInt(raw string, p bool) int {
	start := 1
	if p {
		start = 2
	}
	text := raw[start : len(raw)-1]
	i, e := strconv.Atoi(text)
	if e != nil {
		panic("badly formed number in annotation") // should never happen
	}
	return i
}

func annoToString(raw string, p bool) string {
	start := 1
	if p {
		start = 2
	}
	return raw[start : len(raw)-1]
}

// remove leading and trailing quote
func quotedStringToString(node antlr.Token) string {
	return node.GetText()[1 : len(node.GetText())-1]
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
