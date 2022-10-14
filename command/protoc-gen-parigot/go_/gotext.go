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
	return codegen.OutputType(m,
		func(protoPkg string, method *codegen.WasmMethod, parameter *codegen.CGParameter) string {
			return g.GetReturnValueDecl(protoPkg, method, parameter)
		},
		func(protoPkg string, method *codegen.WasmMethod) string {
			return g.NoReturnValueDecl(protoPkg, method)
		},
	)
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

// OutZeroValue should return a legal value for its type.  This value
// is just for compilers (to keep them quiet) so the value will never be used.
func (g *GoText) OutZeroValue(m *codegen.WasmMethod) string {
	if m.PullParameters() {
		// xxx fix me, this should not be hitting codgen like this
		exp := codegen.ExpandReturnInfoForOutput(m.GetCGOutput(), m, m.GetProtoPackage())
		if exp == nil {
			return ""
		}
		if !exp.GetCGType().IsBasic() {
			return exp.GetCGType().String(m.GetProtoPackage()) + "{}"
		}
		return goZeroValuesForProtoTypes(exp.GetCGType().String(""))
	}
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
		goZeroValuesForProtoTypes(t.String(""))
	}
	return t.String(m.GetProtoPackage()) + "{}"
}

func goZeroValuesForProtoTypes(s string) string {
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

func (g *GoText) BasicTypeToString(s string, panicOnFail bool) string {
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

func (g *GoText) AllInputWithFormalWasmLevel(method *codegen.WasmMethod, showFormalName bool) string {
	result := ""
	file := method.GetParent().GetParent()
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
				result += g.GetFormalTypeSeparator(protoPkg, method, parameter)
			}
			if parameter.GetCGType().IsBasic() {
				result += g.BasicTypeToString(g.GetCGTypeName(protoPkg, method, parameter), true)
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

const letters = "abcdefghijklmnop"

func (g *GoText) AllInputNumberedParam(m *codegen.WasmMethod) string {
	result := ""
	count := 0
	file := m.GetParent().GetParent()
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
				newCG := g.convertGoTypeToWasmType(parameter.GetCGType(), i, m.GetLanguage(), protoPkg)
				param := codegen.NewCGParameterFromString(name, newCG)
				result += param.GetFormalName()
				result += g.GetFormalTypeSeparator(protoPkg, method, param)
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
	file := method.GetParent().GetParent()
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

func NewGoText() *GoText {
	return &GoText{}
}
