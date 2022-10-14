package go_

import (
	"fmt"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/codegen"
	"log"
)

type GoText struct {
}

func (g *GoText) AllInputWithFormal(m *codegen.WasmMethod, showFormalName bool) string {
	result := ""
	result += g.AllInputWithFormalWasmLevel(m, showFormalName)
	return result
}

func (g *GoText) OutType(m *codegen.WasmMethod) string {
	if m.GetCGOutput().IsMultipleReturn() {
		log.Fatalf("unable to process multiple return values (%s) at this time",
			m.GetCGOutput().GetTypeName())
	}
	if m.GetCGOutput().IsEmpty() {
		return ""
	}
	return m.GetInputParam().GetCGType().String(m.GetParent().GetProtoPackage())
	//if m.GetCGOutput().IsMultipleReturn() {
	//	log.Fatalf("unable to process multiple return values (%s) at this time",
	//		m.GetCGOutput().GetTypeName())
	//}
	//if m.GetCGOutput().IsEmpty() {
	//	return ""
	//}
	//// for a typename (which they call a message) we have to figure out how to address the type
	//if m.GetCGOutput().GetParamVar()[0].GetField().GetType().String() == "TYPE_MESSAGE" {
	//	n := ""
	//	if m.HasNoPackageOption() {
	//		n = m.GetFinder().AddressingNameFromMessage("", m.GetCGOutput().GetTyp())
	//	} else {
	//		n = m.GetFinder().AddressingNameFromMessage(m.GetProtoPackage(), m.GetCGOutput().GetTyp())
	//		//return codegen.ComputeMessageName(ourPkg, m.GetCGOutput().GetTypeMessage())
	//	}
	//	return n
	//}
	//result += g.ProtoTypeNameToLanguageTypeName(m.GetCGOutput().GetParamVar()[0].GetField().GetType().String())
	//return result
}

//func (g *GoText) walkInputParams(m *codegen.WasmMethod, separator string, fn func(paramVar *codegen.ParamVar) string) string {
//	result := ""
//	for i, p := range m.GetCGInput().GetParamVar() {
//		result += fn(p)
//		if i != len(m.GetCGInput().GetParamVar())-1 {
//			result += separator
//		}
//	}
//	return result
//}

func (g *GoText) AllInputFormal(method *codegen.WasmMethod) string {
	result := ""
	result += codegen.FuncParamPass(method,
		func(_ *codegen.WasmMethod, _ int, parameter *codegen.CGParameter) string {
			return parameter.GetFormalName()
		},
		func(_ *codegen.WasmMethod, _ bool, _ *codegen.CGParameter) string {
			return ""
		})
	return result
}

// OutZeroValue should return a legal value for its type.  This value
// is just for compilers (to keep them quiet) so the value will never be used.
func (g *GoText) OutZeroValue(m *codegen.WasmMethod) string {
	out := m.GetCGOutput()
	if out == nil || out.IsEmpty() {
		return ""
	}
	if out.IsMultipleReturn() {
		log.Fatalf("unable to process multiple return values (%s) at this time",
			out.GetCGType().ShortName())
	}
	t := out.GetCGType()
	if t.IsBasic() {
		s := t.String("" /* doesn't matter for basic type*/)
		switch s {
		case "TYPE_STRING":
			return "\"\""
		case "TYPE_INT32":
			return "int32(0)"
		case "TYPE_INT64":
			return "int64(0)"
		case "TYPE_FLOAT":
			return "float32(0.0)"
		case "TYPE_DOUBLE":
			return "float64(0.0)"
		case "TYPE_BOOL":
			return "int32(0)"
		case "TYPE_BYTES":
			return "[]byte{0}"
		case "TYPE_BYTE":
			return "byte(0)"
		}
		panic("unable to understand basic type " + s)
	}
	return t.String(m.GetProtoPackage()) + "{}"
}

func (g *GoText) GetFormalName(
	_ string,
	_ *codegen.WasmMethod,
	p *codegen.CGParameter) string {
	return p.GetFormalName()
}

func (g *GoText) GetFormalNameUnused(
	_ string,
	_ *codegen.WasmMethod,
	_ *codegen.CGParameter) string {
	return "_"
}

func (g *GoText) GetFormalTypeSeparator(
	_ string,
	_ *codegen.WasmMethod,
	_ *codegen.CGParameter) string {
	return " "
}

func (g *GoText) GetFormalArgSeparator() string {
	return ","
}

func (g *GoText) GetCGTypeName(
	protoPkg string,
	_ *codegen.WasmMethod,
	p *codegen.CGParameter) string {
	return p.GetCGType().String(protoPkg)
}

func (g *GoText) GetNoInputParams(
	_ string,
	_ *codegen.WasmMethod) string {
	return "()"
}

func (g *GoText) GetNoOutputParams(
	_ string,
	_ *codegen.WasmMethod) string {
	return ""
}

func (g *GoText) CallFuncWithArg(
	funcName string,
	p []*codegen.CGParameter) string {
	paramToFunc := ""
	for i, arg := range p {
		paramToFunc += arg.GetFormalName()
		if i != len(p)-1 {
			paramToFunc += g.GetFormalArgSeparator()
		}
	}
	return fmt.Sprintf("%s(%s)", funcName, paramToFunc)
}

func (g *GoText) GetNumberParametersUsed(
	cgType *codegen.CGType) int {

	if !cgType.IsBasic() {
		return 1
	}
	switch cgType.String("" /*doesn't matter*/) {
	case "TYPE_INT32", "TYPE_INT64", "TYPE_FLOAT",
		"TYPE_DOUBLE", "TYPE_BOOL", "TYPE_BYTE":
		return 1
	case "TYPE_STRING":
		return 2 // ptr + lengeth
	case "TYPE_BYTES":
		return 3 // ptr + length + capacity
	}
	panic("unable to understand number of parameters for " + cgType.String(""))
}

func (g *GoText) BasicTypeToReturnExpr(s string, num int, p *codegen.CGParameter) string {
	switch s {
	case "TYPE_STRING":
		fixed := g.replaceFormalName(fmt.Sprintf("memPtr"), p)
		p0 := g.replaceFormalName(fmt.Sprintf("p%da", num), p)
		p1 := g.replaceFormalName(fmt.Sprintf("p%db", num+1), p)
		allArgs := []*codegen.CGParameter{fixed, p0, p1}
		return g.CallFuncWithArg("strConvert", allArgs)
	case "TYPE_INT32":
		return "int32(0)"
	case "TYPE_INT64":
		return "int64(0)"
	case "TYPE_FLOAT":
		return "float32(0.0)"
	case "TYPE_DOUBLE":
		return "float64(0.0)"
	case "TYPE_BOOL":
		return "int32(0)"
	case "TYPE_BYTES":
		return "[]byte{0}"
	case "TYPE_BYTE":
		return "byte(0)"
	}
	panic("unable to convert to return expr:" + s)
}

func (g *GoText) BasicTypeToString(s string) string {
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

func (g *GoText) AllInputWithFormalWasmLevel(method *codegen.WasmMethod, showFormalName bool) string {
	result := ""
	file := method.GetParent().GetParent()
	protoPkg := file.GetPackage()
	result += codegen.FuncParamPass(method,
		func(method *codegen.WasmMethod, _ int, parameter *codegen.CGParameter) string {
			result := ""
			if showFormalName {
				result = g.GetFormalName(protoPkg, method, parameter)
			} else {
				result = g.GetFormalNameUnused(protoPkg, method, parameter)
			}
			result += g.GetFormalTypeSeparator(protoPkg, method, parameter)
			if parameter.GetCGType().IsBasic() {
				result += g.BasicTypeToString(g.GetCGTypeName(protoPkg, method, parameter))
			} else {
				result += g.GetCGTypeName(protoPkg, method, parameter)
			}
			return result
		},
		func(method *codegen.WasmMethod, isInput bool, param *codegen.CGParameter) string {
			if isInput {
				return g.GetNoInputParams(protoPkg, method)
			} else {
				return g.GetNoOutputParams(protoPkg, method)
			}
		},
	)
	return result
}

func (g *GoText) AllInputNumberedParam(m *codegen.WasmMethod) string {
	return ""
	//count := 0
	//result := ""
	//for i, p := range m.GetCGInput().GetParamVar() {
	//	result += fmt.Sprintf("p%d ", count)
	//	result += g.ProtoTypeNameToLanguageTypeName(p.TypeFromProto())
	//	if i != len(m.GetCGInput().GetParamVar())-1 {
	//		result += ","
	//	}
	//	count++
	//}
	//return result
}

func (g *GoText) replaceFormalName(newName string, parameter *codegen.CGParameter) *codegen.CGParameter {
	cgType := parameter.GetCGType()
	return codegen.NewCGParameterFromString(newName, cgType)
}

func (g *GoText) AllInputWasmToGoImpl(method *codegen.WasmMethod) string {
	file := method.GetParent().GetParent()
	protoPkg := file.GetPackage()
	return codegen.FuncParamPass(method,
		func(method *codegen.WasmMethod, num int, parameter *codegen.CGParameter) string {
			//newParam := g.replaceFormalName(fmt.Sprintf("p%d", num), parameter)
			//result = g.GetFormalName(protoPkg, method, newParam)
			//result += g.GetFormalTypeSeparator(protoPkg, method, parameter)
			return g.BasicTypeToReturnExpr(protoPkg, num, parameter)
		},
		func(method *codegen.WasmMethod, isInput bool, parameter *codegen.CGParameter) string {
			if isInput {
				return g.GetNoInputParams(protoPkg, method)
			} else {
				return g.GetNoInputParams(protoPkg, method)
			}
		})

	//if !m.PullParameters() {
	//	//currentPackage := m.GetParent().GetProtoPackage()
	//	//if m.HasNoPackageOption() {
	//	//	currentPackage = ""
	//	//}
	//	n := m.GetCGInput().GetTypeMessage().GetTypeName()
	//	return strings.ToLower(n[0:1])
	//	//return n
	//} else {
	//	result := ""
	//	count := 0
	//	for i, p := range m.GetCGInput().GetParamVar() {
	//		if p.IsStrictWasmType() {
	//			result += fmt.Sprintf("p%d", count)
	//		} else {
	//			if !p.IsParigotType() {
	//				log.Fatalf("unable to generate implementation helper code for type %s", p.TypeFromProto())
	//			}
	//			switch p.TypeFromProto() {
	//			case "TYPE_STRING":
	//				result += "strConvert(impl.GetMemPtr()," + fmt.Sprintf("p%d,p%d)", count, count+1)
	//			case "TYPE_BYTES":
	//				result += "bytesConvert(impl.GetMemPtr()," + fmt.Sprintf("p%d,p%d,p%d)", count, count+1, count+2)
	//			case "TYPE_BOOL":
	//				result += fmt.Sprintf("p%d!=0", count)
	//			}
	//		}
	//		count += 1
	//		if i != len(m.GetCGInput().GetParamVar())-1 {
	//			result += ","
	//		}
	//	}
	//	return result
	//}
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
