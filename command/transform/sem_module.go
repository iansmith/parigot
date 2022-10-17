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
	targetType := def.TopLevelType()
	if len(m.Code) == 0 || len(m.Code) == 1 {
		m.Code = append(m.Code, def)
		return
	}
	found := false
	for i := range m.Code {
		currentType := m.Code[i].TopLevelType()
		if !found && currentType == targetType {
			found = true
		}
		if found && currentType != targetType {
			m.Code = append(m.Code[:i], append([]TopLevel{def}, m.Code[i:]...)...)
			return
		}
	}
	if !found {
		panic("unable to find any top levels of the correct type")
	}
	panic("unable to find the sequence of toplevels to append to")
}
