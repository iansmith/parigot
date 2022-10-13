package go_

import (
	"fmt"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/codegen"
	"log"
	"strings"
)

type GoText struct {
}

func (g *GoText) ProtoTypeNameToLanguageTypeName(s string) string {
	switch s {
	case "TYPE_STRING":
		return "string"
	case "TYPE_INT32":
		return "int32"
	case "TYPE_INT64":
		return "int64"
	case "TYPE_FLOAT":
		return "float32"
	case "TYPE_DOUBLE":
		return "float64"
	case "TYPE_BOOL":
		return "bool"
	case "TYPE_BYTES":
		return "[]byte"
	case "TYPE_BYTE":
		return "byte"
	}
	panic("unable to convert " + s + " to go type")
}

func (g *GoText) AllInputWithFormal(m *codegen.WasmMethod, showFormalName bool) string {
	return g.walkInputParams(m, ",", func(p *codegen.ParamVar) string {
		result := ""
		if showFormalName {
			result += codegen.ToCamelCase(p.GetName()) + " "
		} else {
			result += "_" + " "
		}
		typ := p.GetTyp()
		result += codegen.ComputeMessageName(m.GetProtoPackage(), typ)
		return result
	})
}

func (g *GoText) OutType(m *codegen.WasmMethod) string {
	result := ""
	if m.GetCGOutput().IsMultipleReturn() {
		log.Fatalf("unable to process multiple return values (%s) at this time",
			m.GetCGOutput().GetName())
	}
	if m.GetCGOutput().IsEmpty() {
		return ""
	}
	// for a typename (which they call a message) we have to figure out how to address the type
	if m.GetCGOutput().GetParamVar()[0].GetField().GetType().String() == "TYPE_MESSAGE" {
		ourPkg := m.GetParent().GetParent().GetPackage()
		return codegen.ComputeMessageName(ourPkg, m.GetCGOutput().GetTyp())
	}
	result += g.ProtoTypeNameToLanguageTypeName(m.GetCGOutput().GetParamVar()[0].GetField().GetType().String())
	return result
}

func (g *GoText) walkInputParams(m *codegen.WasmMethod, separator string, fn func(paramVar *codegen.ParamVar) string) string {
	result := ""
	for i, p := range m.GetCGInput().GetParamVar() {
		result += fn(p)
		if i != len(m.GetCGInput().GetParamVar())-1 {
			result += separator
		}
	}
	return result
}

func (g *GoText) AllInputFormal(m *codegen.WasmMethod) string {
	return g.walkInputParams(m, " ", func(p *codegen.ParamVar) string {
		return p.GetName()
	})
}

func (g *GoText) OutZeroValue(m *codegen.WasmMethod) string {
	result := ""
	if m.GetCGOutput().IsMultipleReturn() {
		log.Fatalf("unable to process multiple return values (%s) at this time",
			m.GetCGOutput().GetName())
	}
	protoT := m.GetCGOutput().GetParamVar()[0].GetField().GetType().String()
	goT := g.ProtoTypeNameToLanguageTypeName(protoT)
	result += goZeroValue(goT)
	return result
}

func (g *GoText) AllInputWithFormalWasmLevel(m *codegen.WasmMethod, showFormalName bool) string {
	result := ""
	currentParam := 0
	for i, p := range m.GetCGInput().GetParamVar() {
		if showFormalName {
			result += fmt.Sprintf("p%d", currentParam) + " "
		} else {
			result += "_" + " "
		}
		if p.IsStrictWasmType() {
			resp := " " + g.ProtoTypeNameToLanguageTypeName(p.TypeFromProto())
			parts := strings.Split(resp, " ")
			result += resp
			if len(parts) != 2 {
				log.Printf("%s %d--%s,%+v", m.GetWasmMethodName(), currentParam, result, parts)
			}
		} else {
			// do our conversion
			// xxx only strings for now
			if !p.IsParigotType() {
				log.Fatalf("unable to convert type %s to WASM type", p.TypeFromProto())
			}
			switch p.TypeFromProto() {
			// we convert string to a pair of int32s because that is what the compiler
			// outputs at the wasm level.  the two int32s are a pointer and a length
			// slices are passed as ptr plus len and cap  all are int32.
			case "TYPE_STRING":
				result += "int32,"
				currentParam++
				result += fmt.Sprintf("p%d", currentParam) + " "
				result += "int32"
			case "TYPE_BOOL":
				result += "int32"
			case "TYPE_BYTE":
				result += "int32"
			case "TYPE_BYTES":
				result += fmt.Sprintf("int32")
				currentParam++
				result += fmt.Sprintf(",p%d", currentParam) + " "
				result += "int32"
				currentParam++
				result += fmt.Sprintf(",p%d", currentParam) + " "
				result += "int32"
			}
		}
		currentParam++
		if i != len(m.GetCGInput().GetParamVar())-1 {
			result += ","
		}
	}
	return result

}

func (g *GoText) AllInputNumberedParam(m *codegen.WasmMethod) string {
	count := 0
	result := ""
	for i, p := range m.GetCGInput().GetParamVar() {
		result += fmt.Sprintf("p%d ", count)
		result += g.ProtoTypeNameToLanguageTypeName(p.TypeFromProto())
		if i != len(m.GetCGInput().GetParamVar())-1 {
			result += ","
		}
		count++
	}
	return result
}

func (g *GoText) AllInputWasmToGoImpl(m *codegen.WasmMethod) string {
	result := ""
	count := 0
	for i, p := range m.GetCGInput().GetParamVar() {
		if p.IsStrictWasmType() {
			result += fmt.Sprintf("p%d", count)
		} else {
			if !p.IsParigotType() {
				log.Fatalf("unable to generate implementation helper code for type %s", p.TypeFromProto())
			}
			switch p.TypeFromProto() {
			case "TYPE_STRING":
				result += "strConvert(impl.GetMemPtr()," + fmt.Sprintf("p%d,p%d)", count, count+1)
			case "TYPE_BYTES":
				result += "bytesConvert(impl.GetMemPtr()," + fmt.Sprintf("p%d,p%d,p%d)", count, count+1, count+2)
			case "TYPE_BOOL":
				result += fmt.Sprintf("p%d!=0", count)
			}
		}
		count += 1
		if i != len(m.GetCGInput().GetParamVar())-1 {
			result += ","
		}
	}
	return result
}

func NewGoText() *GoText {
	return &GoText{}
}

// goZeroValue returns the simplest, empty value for the given go type.
func goZeroValue(s string) string {
	switch s {
	case "string":
		return ""
	case "int32":
		return "int32(0)"
	case "int64":
		return "int64(0)"
	case "float32":
		return "float32(0.0)"
	case "float64":
		return "float64(0.0)"
	case "bool":
		return "false"
	case "[]byte":
		return "[]byte{}"
	}
	panic("unable to get zero value for go type " + s)
}
