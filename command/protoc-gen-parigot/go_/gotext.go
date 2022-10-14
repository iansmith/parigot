package go_

import (
	"github.com/iansmith/parigot/command/protoc-gen-parigot/codegen"
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
	return ""
	//if !m.PullParameters() {
	//	currentPackage := m.GetParent().GetProtoPackage()
	//	if m.HasNoPackageOption() {
	//		currentPackage = ""
	//	}
	//	a := m.GetFinder().AddressingNameFromMessage(currentPackage, m.GetCGInput().GetTypeMessage())
	//	//n := m.GetCGInput().GetTypeMessage().GetTypeName()
	//	result := m.GetCGInput().GetTypeMessage().GetTypeName()[0:1] + " " + a
	//	return result
	//} else {
	//	result := ""
	//	fields := m.GetCGInput().GetTypeMessage().GetField()
	//	for i, f := range fields {
	//		result += f.GetTypeName() + " " + f.GetTypeName()
	//		if i != len(fields)-1 {
	//			result += ","
	//		}
	//	}
	//	return result
	//}
	//return g.walkInputParams(m, ",", func(p *codegen.ParamVar) string {
	//	result := ""
	//	if showFormalName {
	//		result += p.GetTypeName() + " "
	//	} else {
	//		result += "_" + " "
	//	}
	//	typ := p.GetTypeMessage()
	//	currentPkg := m.GetProtoPackage()
	//	if m.HasNoPackageOption() {
	//		currentPkg = ""
	//	}
	//	//result += codegen.ComputeMessageName(m.GetProtoPackage(), typ)
	//	result += m.GetFinder().AddressingNameFromMessage(currentPkg, typ)
	//	return result
	//})
}

func (g *GoText) OutType(m *codegen.WasmMethod) string {
	result := ""
	return result
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

func (g *GoText) AllInputFormal(m *codegen.WasmMethod) string {
	return ""
	//return g.walkInputParams(m, " ", func(p *codegen.ParamVar) string {
	//	return p.GetTypeName()
	//})
}

func (g *GoText) OutZeroValue(m *codegen.WasmMethod) string {
	result := ""
	//if m.GetCGOutput().IsMultipleReturn() {
	//	log.Fatalf("unable to process multiple return values (%s) at this time",
	//		m.GetCGOutput().GetTypeName())
	//}
	//protoT := m.GetCGOutput().GetParamVar()[0].GetField().GetType().String()
	//goT := g.ProtoTypeNameToLanguageTypeName(protoT)
	//result += goZeroValue(goT)
	return result
}

func (g *GoText) AllInputWithFormalWasmLevel(m *codegen.WasmMethod, showFormalName bool) string {
	result := ""
	//if !m.PullParameters() {
	//	currentPackage := m.GetParent().GetProtoPackage()
	//	if m.HasNoPackageOption() {
	//		currentPackage = ""
	//	}
	//	a := m.GetFinder().AddressingNameFromMessage(currentPackage, m.GetCGInput().GetTypeMessage())
	//	last := codegen.LastSegmentOfPackage(a)
	//	if len(last) == 0 {
	//		panic("unable to understand types! no last part of the package")
	//	}
	//	return strings.ToLower(last[0:1]) + " " + a
	//} else {
	//	result := ""
	//	//currentParam := 0
	//	//currentPackage := m.GetProtoPackage()
	//	// we need to convert the input to the proper types via the CGType
	//	//params := m.InParams()
	//	//newParams := make([]*codegen.ParamVar, len(m.GetCGInput().GetTypeMessage().GetField()))
	//	for _, pv := range m.GetCGInput().GetParamVar() {
	//		var n string
	//		if pv.GetCGType() == nil {
	//			n = pv.GetTypeName()
	//		} else {
	//			n = pv.GetCGType().GetTypeName()
	//		}
	//		p := codegen.NewParamVarWithType(, m.GetProtoPackage(),
	//			n, m.GetFinder())
	//		m.GetCGInput().AddParamVar(p)
	//	}
	//
	//	//for _, p := range newParams {
	//	//	if p.GetTypeMessage() == nil {
	//	//		result += p.GetTypeName() + " " + p.GetField().GetBasicType()
	//	//	} else {
	//	//		result += p.GetTypeName() + " " + p.GetTypeMessage().GetAddressableName(m.GetProtoPackage())
	//	//	}
	//	//}
	//	return result
	//}
	//
	//result := ""
	//currentParam := 0
	//for i, p := range m.GetCGInput().GetParamVar() {
	//	if showFormalName {
	//		result += fmt.Sprintf("p%d", currentParam) + " "
	//	} else {
	//		result += "_" + " "
	//	}
	//	if p.IsStrictWasmType() {
	//		resp := " " + g.ProtoTypeNameToLanguageTypeName(p.TypeFromProto())
	//		parts := strings.Split(resp, " ")
	//		result += resp
	//		if len(parts) != 2 {
	//			log.Printf("%s %d--%s,%+v", m.GetWasmMethodName(), currentParam, result, parts)
	//		}
	//	} else {
	//		if m.HasNoPackageOption() {
	//			log.Printf("overriding input param %s %s", p.Formal(), p.GetTypeMessage().GetTypeName())
	//			n := m.GetFinder().AddressingNameFromMessage("", p.GetTypeMessage())
	//			return n
	//		}
	//		// do our conversion
	//		if !p.IsParigotType() {
	//			log.Fatalf("unable to convert type %s to WASM type", p.TypeFromProto())
	//		}
	//		switch p.TypeFromProto() {
	//		// we convert string to a pair of int32s because that is what the compiler
	//		// outputs at the wasm level.  the two int32s are a pointer and a length
	//		// slices are passed as ptr plus len and cap  all are int32.
	//		case "TYPE_STRING":
	//			result += "int32,"
	//			currentParam++
	//			result += fmt.Sprintf("p%d", currentParam) + " "
	//			result += "int32"
	//		case "TYPE_BOOL":
	//			result += "int32"
	//		case "TYPE_BYTE":
	//			result += "int32"
	//		case "TYPE_BYTES":
	//			result += fmt.Sprintf("int32")
	//			currentParam++
	//			result += fmt.Sprintf(",p%d", currentParam) + " "
	//			result += "int32"
	//			currentParam++
	//			result += fmt.Sprintf(",p%d", currentParam) + " "
	//			result += "int32"
	//		}
	//	}
	//	currentParam++
	//	if i != len(m.GetCGInput().GetParamVar())-1 {
	//		result += ","
	//	}
	//}
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

func (g *GoText) AllInputWasmToGoImpl(m *codegen.WasmMethod) string {
	return ""
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
