package lib

import (
	"bytes"
	"fmt"
)

// Module is the unit of linkage.
type Module struct {
	Code []TopLevel
}

func (m *Module) IndentedString(indented int) string {
	var buf bytes.Buffer
	for i := 0; i < indented; i++ {
		buf.WriteString(" ")
	}
	buf.WriteString(fmt.Sprintf("(module"))
	for _, tl := range m.Code {
		buf.WriteString("\n" + tl.(IndentedStringer).IndentedString(indented+2))
	}
	buf.WriteString("\n")
	// ugh, leading space... why oh why?
	buf.WriteString(" )\n") // for the module closing
	return buf.String()
}

func (m *Module) AddTopLevelDef(def TopLevel) {
	m.Code = append(m.Code, def)
}
