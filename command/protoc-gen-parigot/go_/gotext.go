package go_

import (
	"fmt"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/codegen"
	"log"
	"strings"
)

type GoText struct {
}

func (g *GoText) AllInputWithFormal(m *codegen.WasmMethod, showFormalName bool) string {
	result := ""
	result += g.AllInputWithFormalWasmLevel(m, showFormalName)
	return result
}

func (g *GoText) OutTypeDecl(m *codegen.WasmMethod) string {
	if m.PullParameters() {
		comp := m.CGOutput().GetCGType().CompositeType()
		if len(comp.GetField()) == 0 {
			return "error"
		}
		param := codegen.NewCGParameterFromField(comp.GetField()[0], m, m.ProtoPackage())
		return fmt.Sprintf("(%s,error)", g.GetReturnValueDecl(m.ProtoPackage(), m, param))
	}
	return g.OutType(m)
}

func (g *GoText) OutType(m *codegen.WasmMethod) string {
	result := ""
	// if we had expansion for the output, would be this call...
	// codegen.ExpandReturnInfoForOutput(m.OutputParam(), m, m.ProtoPackage())
	protoPackage := m.ProtoPackage()
	outType := m.CGOutput().GetCGType()
	if m.CGOutput().IsEmpty() || outType.IsCompositeNoFields() {
		result += ""
	} else {
		result += outType.String(protoPackage)
	}
	return result
}

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

func (g *GoText) ReturnErrorDecl(m *codegen.WasmMethod, msg string) string {
	if m.PullParameters() {
		comp := m.CGOutput().GetCGType().CompositeType()
		err := fmt.Sprintf("parigot.NewFromError(\"%s\",err)", msg)
		if len(comp.GetField()) == 0 {
			return "return " + err
		}
		return fmt.Sprintf("return resp," + err)
	}
	return "return " + g.OutZeroValue(m) + "err"
}

func (g *GoText) ReturnValueDecl(m *codegen.WasmMethod) string {
	if m.PullParameters() {
		comp := m.CGOutput().GetCGType().CompositeType()
		if len(comp.GetField()) == 0 {
			return "return nil"
		}
		return fmt.Sprintf("return %s,nil", g.OutZeroValue(m))
	}
	return "return " + g.OutZeroValue(m) + ",nil"
}

func (g *GoText) OutZeroValueDecl(m *codegen.WasmMethod) string {
	if m.PullOutput() {
		t := m.CGOutput().GetCGType()
		if t.IsCompositeNoFields() {
			return "return nil"
		}
		return fmt.Sprintf("return %s,nil", g.OutZeroValue(m))
	}
	return "return " + g.OutType(m)

}

// OutZeroValue should return a legal value for its type.  This value
// is just for compilers (to keep them quiet) so the value will never be used.
func (g *GoText) OutZeroValue(m *codegen.WasmMethod) string {
	if m.PullOutput() {
		// xxx fix me, this should not be hitting codgen like this
		exp := codegen.ExpandReturnInfoForOutput(m.CGOutput(), m, m.ProtoPackage())
		if exp == nil {
			return ""
		}
		if !exp.GetCGType().IsBasic() {
			return exp.GetCGType().String(m.ProtoPackage()) + "{}"
		}
		return g.ZeroValuesForProtoTypes(exp.GetCGType().String(""))
	}
	out := m.CGOutput()
	if out == nil || out.IsEmpty() {
		return ""
	}
	if out.IsMultipleReturn() {
		log.Fatalf("unable to process multiple return values (%s) at this time",
			out.GetCGType().ShortName())
	}
	t := out.GetCGType()
	if t.IsBasic() {
		g.ZeroValuesForProtoTypes(t.String(""))
	}
	return t.String(m.ProtoPackage()) + "{}"
}

func (g *GoText) ZeroValuesForProtoTypes(s string) string {
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
func (g *GoText) ToId(id string, param bool, _ *codegen.WasmMethod) string {
	if param {
		return codegen.ToCamelCaseFirstLower(id)
	}
	return codegen.ToCamelCase(id)
}

func (g *GoText) EmptyComposite(typeName string, _ *codegen.WasmMethod) string {
	return typeName + "{}"
}

func (g *GoText) NilValue() string {
	return "nil"
}

func (g *GoText) ToTypeName(tn string, ref bool, _ *codegen.WasmMethod) string {
	parts := strings.Split(tn, ".")

	// do simple case
	if len(parts) == 1 {
		name := codegen.ToCamelCase(tn)
		if ref {
			name = "*" + name
		}
		return name
	}

	// complex case
	for i := 0; i < len(parts)-1; i++ {
		parts[i] = strings.ToLower(parts[i])
	}
	parts[len(parts)-1] = strings.ToLower(parts[len(parts)-1])
	name := strings.Join(parts, ".")
	if ref {
		name = "*" + name
	}
	return name
}

func (g *GoText) GetReturnValueDecl(
	protoPkg string,
	_ *codegen.WasmMethod,
	p *codegen.CGParameter) string {
	t := p.GetCGType()
	if t.IsBasic() {
		return g.BasicTypeToString(t.String(""), true)
	}
	return p.GetCGType().String(protoPkg)
}

func (g *GoText) NoReturnValueDecl(
	_ string,
	_ *codegen.WasmMethod) string {
	return ""
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

func (g *GoText) GetFormalTypeCombination(formal, typ string) string {
	return formal + " " + typ
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
func (c *GoText) BasicTypeToString(s string, panicOnFail bool) string {
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
	if !panicOnFail {
		return ""
	}
	panic("unable to convert " + s + " to go type")
}

func (g *GoText) BasicTypeToReturnExpr(s string, num int, p *codegen.CGParameter) string {
	switch s {
	case "TYPE_STRING":
		fixed := g.replaceFormalName(fmt.Sprintf("impl.GetMemPtr()"), p)
		p0 := g.replaceFormalName(fmt.Sprintf("p%da", num), p)
		p1 := g.replaceFormalName(fmt.Sprintf("p%db", num), p)
		allArgs := []*codegen.CGParameter{fixed, p0, p1}
		return g.CallFuncWithArg("strConvert", allArgs)
	case "TYPE_INT32", "TYPE_INT64", "TYPE_FLOAT", "TYPE_DOUBLE":
		return fmt.Sprintf("p%d", num)
	case "TYPE_BOOL":
		return fmt.Sprintf("p%d != 0", num)
	case "TYPE_BYTES":
		fixed := g.replaceFormalName(fmt.Sprintf("impl.GetMemPtr()"), p)
		p0 := g.replaceFormalName(fmt.Sprintf("p%da", num), p)
		p1 := g.replaceFormalName(fmt.Sprintf("p%db", num), p)
		p2 := g.replaceFormalName(fmt.Sprintf("p%db", num), p)
		allArgs := []*codegen.CGParameter{fixed, p0, p1, p2}
		return g.CallFuncWithArg("bytesConvert", allArgs)
	case "TYPE_BYTE":
		return "byte(0)"
	}
	panic("unable to convert to return expr:" + s)
}

func (g *GoText) AllInputWithFormalWasmLevel(method *codegen.WasmMethod, showFormalName bool) string {
	result := ""
	file := method.Parent().GetParent()
	protoPkg := file.GetPackage()
	result += codegen.FuncParamPass(method,
		func(method *codegen.WasmMethod, _ int, parameter *codegen.CGParameter) string {
			result := ""
			if showFormalName {
				// we only show the formal when we have on.. typically outputs don't have a formal
				if parameter.HasFormal() {
					result = g.GetFormalName(protoPkg, method, parameter)
				}
			} else {
				result = g.GetFormalNameUnused(protoPkg, method, parameter)
			}
			if parameter.HasFormal() {
				result += " " // xxx
			}
			//if parameter.CGType().IsBasic() {
			//	result += g.BasicTypeToString(g.GetCGTypeName(protoPkg, method, parameter), true)
			//} else {
			//	result += g.GetCGTypeName(protoPkg, method, parameter)
			//}
			result += parameter.GetCGType().String(protoPkg)
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

const letters = "abcdefghijklmnop"

func (g *GoText) AllInputNumberedParam(m *codegen.WasmMethod) string {
	result := ""
	count := 0
	file := m.Parent().GetParent()
	protoPkg := file.GetPackage()
	result += codegen.FuncParamPass(m,
		func(method *codegen.WasmMethod, n int, parameter *codegen.CGParameter) string {
			result := ""
			used := 0
			if parameter.GetCGType().IsBasic() {
				used = g.GetNumberParametersUsed(parameter.GetCGType())
			} else {
				used = 1
			}
			for i := 0; i < used; i++ {
				l := letters[i : i+1]
				name := fmt.Sprintf("p%d", n)
				if used > 1 {
					name = fmt.Sprintf("p%d%s", n, l)
				}
				newCG := g.convertGoTypeToWasmType(parameter.GetCGType(), i, m.Language(), protoPkg)
				param := codegen.NewCGParameterFromString(name, newCG)
				result += param.GetFormalName()
				result += " " //xxx
				if parameter.GetCGType().IsBasic() {
					result += g.BasicTypeToString(g.GetCGTypeName(protoPkg, method, param), true)
				} else {
					result += g.GetCGTypeName(protoPkg, method, param)
				}
				if i != used-1 {
					result += g.GetFormalArgSeparator()
				}
			}
			count += used
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

func (g *GoText) replaceFormalName(newName string, parameter *codegen.CGParameter) *codegen.CGParameter {
	cgType := parameter.GetCGType()
	return codegen.NewCGParameterFromString(newName, cgType)
}

func (g *GoText) AllInputWasmToGoImpl(method *codegen.WasmMethod) string {
	file := method.Parent().GetParent()
	protoPkg := file.GetPackage()
	return codegen.FuncParamPass(method,
		func(method *codegen.WasmMethod, num int, parameter *codegen.CGParameter) string {
			//newParam := g.replaceFormalName(fmt.Sprintf("p%d", num), parameter)
			//result = g.GetFormalName(protoPkg, method, newParam)
			//result += g.GetFormalTypeSeparator(protoPkg, method, parameter)
			if !parameter.GetCGType().IsBasic() {
				return parameter.GetCGType().String(protoPkg) + "{}"
			}
			return g.BasicTypeToReturnExpr(parameter.GetCGType().String(""), num, parameter)
		},
		func(method *codegen.WasmMethod, isInput bool, parameter *codegen.CGParameter) string {
			if isInput {
				return g.GetNoInputParams(protoPkg, method)
			} else {
				return g.GetNoInputParams(protoPkg, method)
			}
		})
}

// is the magic translation the parameter types used by the compiler
func (g *GoText) convertGoTypeToWasmType(t *codegen.CGType, n int, l codegen.LanguageText, protoPkg string) *codegen.CGType {
	if t.IsBasic() {
		switch t.String("") {
		case "TYPE_STRING":
			// both are INT32
			return codegen.NewCGTypeFromBasic("TYPE_INT32", l, protoPkg)
		case "TYPE_BOOL":
			// bools are 32 bits
			return codegen.NewCGTypeFromBasic("TYPE_INT32", l, protoPkg)
		case "TYPE_BYTES":
			// all are 32 bits
			return codegen.NewCGTypeFromBasic("TYPE_INT32", l, protoPkg)
		case "TYPE_INT32", "TYPE_INT64", "TYPE_FLOAT", "TYPE_DOUBLE":
			return t
		}
		panic("unable to convert simple type %s" + t.String(""))
	}
	panic("we cannot handle composite types yet in the conversion to WASM")
}

func (g *GoText) BasicTypeToWasm(t string) []string {
	switch t {
	case "TYPE_STRING":
		// both are INT32
		return []string{"TYPE_INT32", "TYPE_INT32"}
	case "TYPE_BOOL":
		// bools are 32 bits
		return []string{"TYPE_INT32"}
	case "TYPE_BYTE":
		// bytes are 32 bits, sadly
		return []string{"TYPE_INT32"}
	case "TYPE_BYTES":
		// all are 32 bits
		return []string{"TYPE_INT32", "TYPE_INT32", "TYPE_INT32"}
	case "TYPE_INT32", "TYPE_INT64", "TYPE_FLOAT", "TYPE_DOUBLE":
		return []string{t}
	}
	panic(fmt.Sprintf("unable to convert simple type %s to wasm", t))
}

func NewGoText() *GoText {
	return &GoText{}
}
