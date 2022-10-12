package go_

import (
	"fmt"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/codegen"
	"log"
)

type GoText struct{}

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
	}
	panic("unable to convert " + s + " to go type")
}
func (g *GoText) AllInputWithFormal(m *codegen.WasmMethod, showFormalName bool) string {
	result := ""
	for i, p := range m.GetInput().GetParamVar() {
		if showFormalName {
			result += codegen.ToCamelCase(p.GetName()) + " "
		} else {
			result += "_" + " "
		}
		result += g.ProtoTypeNameToLanguageTypeName(p.TypeFromProto())
		if i != len(m.GetInput().GetParamVar())-1 {
			result += ","
		}
	}
	return result
}

func (g *GoText) OutType(m *codegen.WasmMethod) string {
	if m.GetOutput().IsMultipleReturn() {
		log.Fatalf("unable to process multiple return values (%s) at this time",
			m.GetOutput().GetName())
	}
	if m.GetOutput().IsEmpty() {
		return ""
	}
	return g.ProtoTypeNameToLanguageTypeName(m.GetOutput().GetParamVar()[0].GetField().GetType().String())
}

func (g *GoText) AllInputFormal(m *codegen.WasmMethod) string {
	result := ""
	for i, p := range m.GetInput().GetParamVar() {
		result += p.GetName()
		if i != len(m.GetInput().GetParamVar())-1 {
			result += ","
		}
	}
	return result
}

func (g *GoText) OutZeroValue(m *codegen.WasmMethod) string {
	if m.GetOutput().IsMultipleReturn() {
		log.Fatalf("unable to process multiple return values (%s) at this time",
			m.GetOutput().GetName())
	}
	protoT := m.GetOutput().GetParamVar()[0].GetField().GetType().String()
	goT := g.ProtoTypeNameToLanguageTypeName(protoT)
	return goZeroValue(goT)
}

func (g *GoText) AllInputParamWithFormalWasmLevel(m *codegen.WasmMethod, showFormalName bool) string {
	result := ""
	currentParam := 0
	for i, p := range m.GetInput().GetParamVar() {
		if showFormalName {
			result += fmt.Sprintf("p%d", currentParam) + " "
		} else {
			result += "_" + " "
		}
		if p.IsStrictWasmType() {
			result += g.ProtoTypeNameToLanguageTypeName(p.TypeFromProto())
		} else {
			// do our conversion
			// xxx only strings for now
			if !p.IsWasmType() {
				log.Fatalf("unable to convert type %s to WASM type", p.TypeFromProto())
			}
			switch p.TypeFromProto() {
			// we convert string to a pair of int32s because that is what the compiler
			// outputs at the wasm level.  the two int32s are a pointer and a length.
			case "TYPE_STRING":
				result += "int32,"
				currentParam++
				result += fmt.Sprintf("p%d", currentParam) + " "
				result += "int32"
			case "TYPE_BOOL":
				result += "int32"
			}
		}
		currentParam++
		if i != len(m.GetInput().GetParamVar())-1 {
			result += ","
		}
	}
	return result

}

func (g *GoText) AllInputNumberedParam(m *codegen.WasmMethod) string {
	count := 0
	result := ""
	for i, p := range m.GetInput().GetParamVar() {
		result += fmt.Sprintf("p%d ", count)
		result += g.ProtoTypeNameToLanguageTypeName(p.TypeFromProto())
		if i != len(m.GetInput().GetParamVar())-1 {
			result += ","
		}
		count++
	}
	return result
}

func (g *GoText) AllInputParamWasmToGoImpl(m *codegen.WasmMethod) string {
	result := ""
	count := 0
	for i, p := range m.GetInput().GetParamVar() {
		if p.IsStrictWasmType() {
			result += fmt.Sprintf("p%d", count)
		} else {
			if !p.IsWasmType() {
				log.Fatalf("unable to generate implementation helper code for type %s", p.TypeFromProto())
			}
			switch p.TypeFromProto() {
			case "TYPE_STRING":
				result += "strConvert(impl.GetMemPtr()," + fmt.Sprintf("p%d,p%d)", count, count+1)
			case "TYPE_BOOL":
				result += fmt.Sprintf("p%d!=0", count)
				count++
			}
		}
		count += 1
		if i != len(m.GetInput().GetParamVar())-1 {
			result += ","
		}
	}
	return result
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
	}
	panic("unable to get zero value for go type " + s)
}
