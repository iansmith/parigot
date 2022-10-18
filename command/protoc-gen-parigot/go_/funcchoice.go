package go_

import (
	"fmt"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/codegen"
	"strings"
)

func (g *GoText) FuncChoice() *codegen.FuncChooser {
	return &codegen.FuncChooser{
		Bits:                funcChoicesToBits,
		NeedsFill:           funcChoicesNeedsFill,
		NeedsRet:            funcChoicesNeedsRet,
		InputParam:          funcChoicesInputParam,
		NeedsPullApart:      funcChoicesNeedsPullApart,
		Inbound:             funcChoicesInbound,
		Outbound:            funcChoicesOutbound,
		RetError:            funcChoicesRetError,
		RetValue:            funcChoicesRetValue,
		MethodRet:           funcChoicesMethodRet,
		ZeroValueRet:        funcChoicesZeroValueRet,
		MethodParamDecl:     funcChoicesMethodParamDecl,
		OutLocal:            funcChoicesOutLocal,
		MethodCall:          funcChoicesMethodCall,
		DecodeRequired:      funcChoicesDecodeRequired,
		NoDecodeRequired:    funcChoicesNoDecodeRequired,
		MethodParamDeclWasm: funcChoicesMethodParamDeclWasm,
		HasComplexParam:     funcChoicesHasComplexParam,
		MethodCallWasm:      funcChoicesMethodCallWasm,
	}
}

// paramWalker is a utility function for creating a parameter decl list or a
// parameter call list.  It takes a function that returns 1 or 2 argumest. If
// the function returns 1 argument (2nd string is "") paramWalker assumes you
// are creating a call list and with 2 arguments a declaration list. It will
// walk through each of the parameters (CGParameter) and call the function on
// it.  It will return the concatenation of the results of the calls to fn,
// correctly separated.
func paramWalker(b1, b2, b3, b4 bool, method *codegen.WasmMethod,
	fn func(parameter *codegen.CGParameter, lang codegen.LanguageText, protoPkg string) (string, string)) string {
	param := methodToCGParameter(b1, b2, b3, b4, method)
	if len(param) == 0 {
		return ""
	}
	lang := method.Language()
	asString := make([]string, len(param))
	for i, p := range param {
		formal, typ := fn(p, lang, method.ProtoPackage())
		if typ == "" {
			asString[i] = formal
		} else {
			asString[i] = lang.GetFormalTypeCombination(formal, typ)
		}
	}
	return strings.Join(asString, lang.GetFormalArgSeparator())
}

func basicTypeRequiresDecode(s string) bool {
	switch s {
	case "TYPE_STRING", "TYPE_BYTES", "TYPE_BOOL", "TYPE_BYTE":
		return true
	}
	return false
}

// funcChoicesComplexParam is a convenience method for grouping methods in the ABI
func funcChoicesHasComplexParam(b1, b2, b3, b4 bool, m *codegen.WasmMethod) bool {
	if !b1 {
		return false
	}
	decode := false
	_ = paramWalker(b1, b2, b3, b4, m, func(p *codegen.CGParameter, lang codegen.LanguageText, protoPkg string) (string, string) {
		if !p.GetCGType().IsBasic() && p.GetCGType().IsCompositeNoFields() == false {
			panic(fmt.Sprintf("abi types should not be using composite for parameters: %s", m.GetName()))
		}
		if basicTypeRequiresDecode(p.GetCGType().Basic()) {
			decode = true
		}
		return "bogus", "cruft"
	})
	return decode

}

// funcChoicesDecodeRequired is a convenience method for grouping methods in the ABI
func funcChoicesDecodeRequired(b1, b2, b3, b4 bool, m *codegen.WasmMethod) bool {
	if b1 {
		return funcChoicesHasComplexParam(b1, b2, b3, b4, m)
	}
	if b2 {
		detail, scenario := outputTypeInfo(b4, m)
		switch scenario {
		case outTypeCompositeNoFields:
			return false
		case outTypeBasic:
			return basicTypeRequiresDecode(detail.Basic())
		case outTypeComposite:
			panic(fmt.Sprintf("abi types should not be using composites: %s", m.GetName()))
		}
	}
	panic(fmt.Sprintf("method has neither input nor output: %s", m.GetName()))
}

// funcChoicesNoDecodeRequired is a convenience method for grouping methods in the ABI
func funcChoicesNoDecodeRequired(b1, b2, b3, b4 bool, m *codegen.WasmMethod) bool {
	return !funcChoicesDecodeRequired(b1, b2, b3, b4, m)
}

// methodToCGParameter returns an array of CGParameter objects corresponding to the
// parameters of the function.  This accounts for parameter pull up.  This function
// can return nil if there are no parameters.
func methodToCGParameter(b1, b2, b3, b4 bool, method *codegen.WasmMethod) []*codegen.CGParameter {
	if !b1 {
		return []*codegen.CGParameter{}
	}
	if b3 {
		//we checked b1 so we know there are fields
		comp := method.CGInput().CGType().CompositeType()
		result := []*codegen.CGParameter{}
		for _, f := range comp.GetField() {
			p := codegen.NewCGParameterFromField(f, method, method.ProtoPackage())
			result = append(result, p)
		}
		return result
	}
	t := method.CGInput().CGType()
	lang := method.Language()
	formal := lang.ToId(t.ShortName(), true, method)
	p := codegen.NewCGParameterFromString(formal, t)
	return []*codegen.CGParameter{p}
}

// funcChoicesMethodParamDeclWasm is used for declaring the parameters of a method declaration
// for the abi.  ABI methods can only use the 4 basic types and have the complex types
// expanded.
func funcChoicesMethodParamDeclWasm(b1, b2, b3, b4 bool, method *codegen.WasmMethod) string {
	n := 0
	count := 0
	return paramWalker(b1, b2, b3, b4, method, func(parameter *codegen.CGParameter, lang codegen.LanguageText, protoPkg string) (string, string) {
		result := ""
		if !parameter.GetCGType().IsBasic() {
			panic("should not be using composite types in the ABI")
		}
		seq := lang.BasicTypeToWasm(parameter.GetCGType().Basic())
		used := len(seq)
		// XXX This special case has something to do with tinygo's mapping of byte slices.
		// XXX note this parameter comes from the _OUTPUT_ type
		// XXX this code is only called in the case of the ABI, so no need to check that
		cgt := method.CGOutput().GetCGType()
		field := cgt.CompositeType().GetField()
		if n == 0 && len(field) == 1 {
			cgt = codegen.NewCGTypeFromField(field[0], method, method.ProtoPackage())
			if cgt.IsBasic() && cgt.Basic() == "TYPE_BYTES" {
				result = "ret int32,"
			}
		}

		for i := 0; i < used; i++ {
			l := letters[i : i+1]
			name := fmt.Sprintf("p%d", n)
			if used > 1 {
				name = fmt.Sprintf("p%d%s", n, l)
			}
			result += fmt.Sprintf("%s %s", name, lang.BasicTypeToString(seq[i], true))
			if i != used-1 {
				result += lang.GetFormalArgSeparator()
			}
		}
		count += used
		n++
		return result, ""
	})
}

// funcChoicesMethodParamDecl is used for declaring the parameters of a method declaration on
// on a service.  the part in caps is the part we are declaring here, using go syntax:
// func foo(A BAR, B BAZ)string
func funcChoicesMethodParamDecl(b1, b2, b3, b4 bool, method *codegen.WasmMethod) string {
	return paramWalker(b1, b2, b3, b4, method, func(parameter *codegen.CGParameter, lang codegen.LanguageText, protoPkg string) (string, string) {
		s := ""
		if parameter.GetCGType().IsBasic() {
			s = lang.BasicTypeToString(parameter.GetCGType().Basic(), true)
		} else {
			cgt := parameter.GetCGType()
			s = lang.ToTypeName(cgt.String(protoPkg), false, method)
		}
		return lang.ToId(parameter.GetFormalName(), true, method), s
	})
}

// funcChoicesMethodCall is used for calling a method of a service.
// the part in caps is the part we are declaring here, using go syntax:
// foo(A, B, C)
func funcChoicesMethodCall(b1, b2, b3, b4 bool, method *codegen.WasmMethod) string {
	return paramWalker(b1, b2, b3, b4, method, func(parameter *codegen.CGParameter, lang codegen.LanguageText, protoPkg string) (string, string) {
		return lang.ToId(parameter.GetFormalName(), true, method), ""
	})
}

// funcChoicesMethodCallWasm is used by the helper of the ABI implementation.
// It uses some convenience functions to convert (unmarshal) the complex
// ABI arguments and call a Go function with the "nicer" types.
func funcChoicesMethodCallWasm(b1, b2, b3, b4 bool, method *codegen.WasmMethod) string {
	n := 0
	return paramWalker(b1, b2, b3, b4, method, func(parameter *codegen.CGParameter, lang codegen.LanguageText, protoPkg string) (string, string) {
		if !parameter.GetCGType().IsBasic() {
			panic("should not be using composite types in the ABI")
		}
		result := ""
		t := parameter.GetCGType().Basic()
		switch t {
		case "TYPE_STRING":
			result = fmt.Sprintf("strConvert(impl.GetMemPtr(),p%da,p%db)", n, n)
		case "TYPE_BYTES":
			result = fmt.Sprintf("bytesConvert(impl.GetMemPtr(),p%da,p%db,p%dc)", n, n, n)
		case "TYPE_BOOL":
			result = fmt.Sprintf("p%d!=0", n)
		default:
			result = fmt.Sprintf("p%d", n)
		}
		n++
		return result, ""
	})
}

func funcChoicesOutLocal(_, b2, _, b4 bool, method *codegen.WasmMethod) string {
	if !b2 {
		return ""
	}
	detail, scenario := outputTypeInfo(b4, method)
	lang := method.Language()
	result := ""
	switch scenario {
	case outTypeCompositeNoFields:
	case outTypeBasic:
		result = lang.ZeroValuesForProtoTypes(detail.Basic())
	case outTypeComposite:
		result = lang.EmptyComposite(detail.String(method.ProtoPackage()), method)
	}
	return result
}

const outTypeCompositeNoFields = 1
const outTypeBasic = 2
const outTypeComposite = 3

func outputTypeInfo(b4 bool, method *codegen.WasmMethod) (*codegen.CGType, int) {
	t := method.CGOutput().GetCGType()
	if b4 {
		if t.IsCompositeNoFields() {
			return nil, outTypeCompositeNoFields
		}
		f := t.CompositeType().GetField()
		inner := codegen.NewCGTypeFromField(f[0], method, method.ProtoPackage())
		if inner.IsBasic() {
			//v := method.Language().ZeroValuesForProtoTypes(inner.String(""))
			return inner, outTypeBasic
		}
		// we allow this to fall through because we return the same thing
		// as they didn't have b4 set, just on the inner one
		t = inner
	}
	//v := t.String(method.ProtoPackage())
	return t, outTypeComposite

}

func funcChoicesZeroValueRet(_, b2, _, b4 bool, abi bool, method *codegen.WasmMethod) string {
	if !abi {
		panic("unepected call to ZeroValueRet")
	}
	if !b2 {
		return ""
	}
	lang := method.Language()
	outDetail, scenario := outputTypeInfo(b4, method)
	result := ""
	switch scenario {
	case outTypeCompositeNoFields:
	case outTypeBasic:
		result = lang.ZeroValuesForProtoTypes(outDetail.Basic())
	case outTypeComposite:
		s := outDetail.String(method.ProtoPackage())
		result = lang.EmptyComposite(lang.ToTypeName(s, false, method), method)
	}
	return "return " + result
}

func funcChoicesNeedsFill(b1, b2, b3, b4 bool) bool {
	return funcChoicesToInt(b1, b2, b3, b4) == 0xa
}
func funcChoicesNeedsRet(_, b2, _, _ bool) bool {
	return b2
}

func funcChoicesInputParam(b1, b2, b3, b4 bool, m *codegen.WasmMethod) string {
	choices := funcChoicesToInt(b1, b2, b3, b4)
	if choices&8 == 0 {
		return ""
	}
	if choices == 10 {
		return ""
	}
	t := m.CGInput().CGType()
	// xxx fix me xxx should not have string in here
	return "req:=" + m.Language().ToId(t.String(m.ProtoPackage()), true, m)
}
func funcChoicesNeedsPullApart(b1, b2, b3, b4 bool) bool {
	return funcChoicesToInt(b1, b2, b3, b4) == 0x5
}
func funcChoicesRetError(b1, b2, b3, b4 bool) string {
	if b2 {
		return "nil,err"
	}
	return "err"
}

func outReturnToString(b2, b4 bool, m *codegen.WasmMethod, abi bool) string {
	if !b2 {
		return ""
	}
	field := m.CGOutput().GetCGType().CompositeType().GetField()
	if abi {
		// abi does not return errors
		if len(field) == 0 {
			return ""
		}
		if b4 {
			// we know it has one item and we know it's not a composite because of
			// the check done by outputCodeNeeded
			return m.Language().BasicTypeToString(codegen.NewCGTypeFromField(field[0], m, m.ProtoPackage()).Basic(), true)
		}
		return m.CGOutput().GetCGType().String(m.ProtoPackage())
	}
	// b2 is true, we know from earlier check
	if b4 {
		if len(field) == 0 {
			return "error"
		}
		return m.Language().BasicTypeToString(codegen.NewCGTypeFromField(field[0], m, m.ProtoPackage()).Basic(), true) + ",error"
	}
	return m.CGOutput().GetCGType().String(m.ProtoPackage()) + ",error"
}

func funcChoicesMethodRet(b1, b2, b3, b4 bool, abi bool, m *codegen.WasmMethod) string {

	if abi {
		// sadly, the abi is a special case because it doesn't return error values
		return outReturnToString(b2, b4, m, true)
	}
	return outReturnToString(b2, b4, m, false)
}

func funcChoicesRetValue(b1, b2, b3, b4 bool, m *codegen.WasmMethod) string {
	//if abi {
	//	// sadly, the abi is a special case because it doesn't return error values
	//	//t := m.CGOutput().CGType()
	//	if b2 {
	//		if b4 == false {
	//			panic("all the ABI functions must be pulling up their return values")
	//		}
	//		lang := m.Language()
	//		detail, scenario := outputTypeInfo(b4, m)
	//		result := ""
	//		switch scenario {
	//		case outTypeCompositeNoFields:
	//		case outTypeBasic:
	//			//result = lang.ZeroValuesForProtoTypes(detail)
	//			result = "return " + lang.ZeroValuesForProtoTypes(detail.Basic())
	//		case outTypeComposite:
	//			panic("should not be using composite types in the ABI")
	//		}
	//		return result
	//	}
	//	// nothing to return
	//	return ""
	//}
	if b2 {
		return "resp,nil"
	}
	return "nil"
}
func funcChoicesInbound(b1, b2, b3, b4 bool) string {
	if b1 {
		return "req"
	}
	return "nil"
}
func funcChoicesOutbound(b1, b2, b3, b4 bool) string {
	if b2 {
		return "resp"
	}
	return "nil"
}

func funcChoicesToInt(b1, b2, b3, b4 bool) int {
	result := 0
	if b1 {
		result += 8
	}
	if b2 {
		result += 4
	}
	if b3 {
		result += 2
	}
	if b4 {
		result += 1
	}
	return result
}
func funcChoicesToBits(b1, b2, b3, b4 bool) string {
	return fmt.Sprintf("%b", funcChoicesToInt(b1, b2, b3, b4))
}
