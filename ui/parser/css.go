package parser

import "fmt"

type CSSSectionNode struct {
}

func (i *CSSSectionNode) Dump(indent int) {
	print(fmt.Sprintf("%*s (CSSSectionNode %s)\n", indent, "", "not implemented "))
}
