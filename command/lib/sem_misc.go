package lib

import "bytes"

type IndentedStringer interface {
	IndentedString(int) string
}

func NewIndentedBuffer(indented int) *bytes.Buffer {
	var buf bytes.Buffer
	for i := 0; i < indented; i++ {
		buf.WriteString(" ")
	}
	return &buf
}

type WasmTypeName int

func (t WasmTypeName) String() string {
	switch t {
	case i32:
		return "i32"
	case i64:
		return "i64"
	case f64:
		return "f64"
	}
	panic("unknown wasm type name") // should never happen
}

const (
	i32 WasmTypeName = 1
	i64 WasmTypeName = 2
	f64 WasmTypeName = 3
)

type TypeNameSeq struct {
	Name []*WasmTypeName
}

func (t *TypeNameSeq) String() string {
	var buf bytes.Buffer
	for i, tn := range t.Name {
		if i != 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(tn.String())
	}
	return buf.String()
}
