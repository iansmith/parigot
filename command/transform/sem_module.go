package transform

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
	buf.WriteString(")\n") // for the module closing
	return buf.String()
}

func (m *Module) AppendTopLevelDef(def TopLevel) {
	if len(m.Code) == 0 || len(m.Code) == 1 {
		m.Code = append(m.Code, def)
		return
	}
	prev := m.Code[0].TopLevelType()
	for i := range m.Code {
		if prev == def.TopLevelType() && m.Code[i].TopLevelType() != def.TopLevelType() {
			//change point
			m.Code = append(m.Code[:i], append([]TopLevel{def}, m.Code[i:]...)...)
			return
		}
	}
	panic("unable to find the sequence of toplevels to append to")
}
