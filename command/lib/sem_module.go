package lib

import (
	"bytes"
	"fmt"
)

// Module is the unit of linkage.
type Module struct {
	topLevel []TopLevelDef
}

func (m *Module) IndentedString(indented int) string {
	var buf bytes.Buffer
	for i := 0; i < indented; i++ {
		buf.WriteString(" ")
	}
	buf.WriteString(fmt.Sprintf("(module\n"))
	for _, tl := range m.topLevel {
		buf.WriteString(tl.(IndentedStringer).IndentedString(indented + 2))
	}
	buf.WriteString("\n") // for the module closing
	return buf.String()
}

func (m *Module) AddTopLevelDef(def TopLevelDef) {
	m.topLevel = append(m.topLevel, def)
}
