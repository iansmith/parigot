package go_

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/iansmith/parigot/command/protoc-gen-parigot/codegen"
)

const GenTestEnvVar = "PARIGOT_GEN_TEST"

type FieldSpec struct {
	name, typ string
}

func (f *FieldSpec) Name() string {
	return f.name
}

func (f *FieldSpec) Type() string {
	return f.typ
}

func (g *GoText) FuncChoice() *codegen.FuncChooser {
	return &codegen.FuncChooser{
		Bits:                  funcChoicesToBits,
		NeedsFillIn:           funcChoicesNeedsFillIn,
		NeedsFillOut:          funcChoicesNeedsFillOut,
		NeedsRet:              funcChoicesNeedsRet,
		InputParam:            funcChoicesInputParam,
		OutputParam:           funcChoicesOutputParam,
		NeedsPullApart:        funcChoicesNeedsPullApart,
		Inbound:               funcChoicesInbound,
		Outbound:              funcChoicesOutbound,
		RetError:              funcChoicesRetError,
		RetValue:              funcChoicesRetValue,
		MethodRet:             funcChoicesMethodRet,
		ZeroValueRet:          funcChoicesZeroValueRet,
		MethodParamDecl:       funcChoicesMethodParamDecl,
		OutLocal:              funcChoicesOutLocal,
		MethodCall:            funcChoicesMethodCall,
		DecodeRequired:        funcChoicesDecodeRequired,
		NoDecodeRequired:      funcChoicesNoDecodeRequired,
		MethodParamDeclWasm:   funcChoicesMethodParamDeclWasm,
		HasComplexParam:       funcChoicesHasComplexParam,
		MethodCallWasm:        funcChoicesMethodCallWasm,
		InputToSend:           funcChoicesInputToSend,
		UsesReturnValuePtr:    funcChoicesUsesReturnValuePtr,
		DispatchParam:         funcChoicesDispatchParam,
		DispatchResult:        funcChoicesDispatchResult,
		OutParamDecl:          funcChoicesOutParamDecl,
		BindDirection:         funcChoicesBindDirection,
		GenMethodPossibleTest: funcGenMethodPossibleTest,
	}
}

// func collectImports(m *codegen.WasmMethod) {
// 	in := m.CGInput().CGType()
// 	parts := strings.Split(in.CompositeType().GetFullName(), ".")
// 	log.Printf("IN %s -> %s", in.ShortName(), parts[0:len(parts)-1])
// 	out := m.CGInput().CGType()
// 	parts = strings.Split(out.CompositeType().GetFullName(), ".")
// 	log.Printf("OUT %s -> %s", in.ShortName(), parts[0:len(parts)-1])
// }

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
			asString[i] = lang.FormalTypeCombination(formal, typ)
		}
	}
	return strings.Join(asString, lang.FormalArgSeparator())
}

func basicTypeRequiresDecode(s string) bool {
	switch s {
	case "TYPE_STRING", "TYPE_BYTES", "TYPE_BOOL", "TYPE_BYTE":
		return true
	}
	return false
}

func funcChoicesBindDirection(b1, b2, _, _ bool, m *codegen.WasmMethod) string {
	if !b1 && !b2 {
		panic(fmt.Sprintf("method %s has neither input nor output!", m.WasmMethodName()))
	}
	if b1 && b2 {
		return "Both"
	}
	if b1 {
		return "In"
	}
	return "Out"
}

// funcChoicesComplexParam is a convenience method for grouping methods in the ABI
func funcChoicesHasComplexParam(b1, b2, b3, b4 bool, m *codegen.WasmMethod) bool {
	if !b1 {
		return false
	}
	decode := false
	_ = paramWalker(b1, b2, b3, b4, m, func(p *codegen.CGParameter, lang codegen.LanguageText, protoPkg string) (string, string) {
		if !p.GetCGType().IsBasic() && !p.GetCGType().IsCompositeNoFields() {
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

// funcGenMethodPossibleTest will return true if the generator should emit the method,
// false to omit the method.
//
// A marked method has the suffix Test, and thus by implication its response and request have
// "Test" in the middle of their types:
//
//  rpc LoadTest(msg.file.v1.LoadTestRequest) returns (msg.file.v1.LoadTestResponse)
//
// This function will return true if:
// * the method name does not have the suffix Test OR
// * the user has not passed the environment variable PARIGOT_GEN_TEST with some value (not "") to our generator
// (protoc-gen-parigot).
//
// We return false if the opposite of both conditions are true.

const magicSuffix = "Test"

func funcGenMethodPossibleTest(method *codegen.WasmMethod) bool {
	if !strings.HasSuffix(method.WasmMethodName(), magicSuffix) {
		return true
	}
	return os.Getenv(GenTestEnvVar) == ""
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

// funcChoicesUsesReturnValuePtr returns true if this method has a return value that
// can't be put on the stack and must put in a return value structure.
func funcChoicesUsesReturnValuePtr(b1, b2, b3, b4 bool, method *codegen.WasmMethod) bool {
	firstParamReturn := false
	detail, scenario := outputTypeInfo(b4, method)
	switch scenario {
	case outTypeCompositeNoFields:
	case outTypeBasic:
		switch detail.Basic() {
		case "TYPE_INT32", "TYPE_BYTE", "TYPE_BOOL", "TYPE_FLOAT":
		default:
			firstParamReturn = true
		}
	case outTypeComposite:
		firstParamReturn = true
	}
	return firstParamReturn
}

// funcChoicesMethodParamDeclWasm is used for declaring the parameters of a method declaration
// for the abi.  ABI methods can only use the 4 basic types and have the complex types
// expanded.
func funcChoicesMethodParamDeclWasm(b1, b2, b3, b4 bool, method *codegen.WasmMethod) string {
	n := 0
	count := 0
	firstParamReturn := funcChoicesUsesReturnValuePtr(b1, b2, b3, b4, method)
	result := ""
	if firstParamReturn {
		result += "retVal int32"
		if !method.CGInput().IsEmpty() {
			result += method.Language().FormalArgSeparator()
		}
	}
	rest := paramWalker(b1, b2, b3, b4, method, func(parameter *codegen.CGParameter, lang codegen.LanguageText, protoPkg string) (string, string) {
		res := ""
		seq := lang.BasicTypeToWasm(parameter.GetCGType().Basic())
		used := len(seq)
		for i := 0; i < used; i++ {
			l := letters[i : i+1]
			name := fmt.Sprintf("p%d", n)
			if used > 1 {
				name = fmt.Sprintf("p%d%s", n, l)
			}
			res += fmt.Sprintf("%s %s", name, lang.BasicTypeToString(seq[i], true))
			if i != used-1 {
				res += lang.FormalArgSeparator()
			}
		}
		count += used
		n++
		return res, ""
	})
	return result + rest
}

// funcChoicesMethodParamDecl is used for declaring the parameters of a method declaration on
// on a service.  the part in caps is the part we are declaring here, using go syntax:
// func foo(A BAR, B BAZ)string
func funcChoicesMethodParamDecl(b1, b2, b3, b4 bool, method *codegen.WasmMethod) string {
	result := ""
	if !b1 && !b2 {
		log.Printf("xxx debug! %s: %s[%d],%s[%d]", method.WasmMethodName(),
			method.CGInput().CGType().String(method.ProtoPackage()),
			len(method.CGInput().CGType().CompositeType().GetField()),
			method.CGOutput().GetCGType().String(method.ProtoPackage()),
			len(method.CGOutput().GetCGType().CompositeType().GetField()),
		)
		panic(fmt.Sprintf("method %s has neither input nor output", method.WasmMethodName()))
	}
	if b1 {
		// have input
		t := method.CGInput().CGType()
		if t.IsBasic() {
			panic(fmt.Sprintf("unexpected basic parameter in method %s", method.WasmMethodName()))
		}
		typ := t.String(method.ProtoPackage())
		result += method.Language().FormalTypeCombination("in", "*"+typ)
		// if b2 {
		// 	result += method.Language().FormalArgSeparator()
		// }
	}
	// if b2 {
	// 	//have output
	// 	t := method.CGOutput().GetCGType()
	// 	if t.IsBasic() {
	// 		panic(fmt.Sprintf("unexpected basic return value in method %s", method.WasmMethodName()))
	// 	}
	// 	typ := t.String(method.ProtoPackage())
	// 	result += method.Language().FormalTypeCombination("out", "*"+typ)
	// }
	return result
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
		case "TYPE_INT64":
			result = fmt.Sprintf("int64(p%d)", n)
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
		return inner, outTypeComposite
	}
	if t.IsBasic() {
		panic(fmt.Sprintf("should not have a method returning a basic for %s", method.WasmMethodName()))
	}
	if t.IsCompositeNoFields() {
		return nil, outTypeCompositeNoFields
	}
	return t, outTypeComposite

}

func funcChoicesDispatchParam(b1, b2, b3, b4 bool, method *codegen.WasmMethod) string {
	if b1 {
		return "req"
	}
	return "nil"
}
func funcChoicesDispatchResult(b1, b2, b3, b4 bool, method *codegen.WasmMethod) string {
	if b2 {
		return "resp"
	}
	return "_"
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

func funcChoicesNeedsFillIn(b1, b2, b3, b4 bool) bool {
	return funcChoicesToInt(b1, b2, b3, b4) == 0xa //b1 and b3
}
func funcChoicesNeedsFillOut(b1, b2, b3, b4 bool) bool {
	return funcChoicesToInt(b1, b2, b3, b4) == 0x5 //b2 and b4
}

func funcChoicesNeedsRet(_, b2, _, _ bool) bool {
	return b2
}

//var versionRegexp = regexp.MustCompile(`^(v[0-9]+)(\..*$)`)

func funcChoicesInputParam(b1, _, _, wantDecl bool, m *codegen.WasmMethod) string {
	if !b1 {
		return ""
	}
	//from := m.Parent().ProtoPackage()
	t := m.CGInput().CGType()
	//protoPkg := m.ProtoPackage()
	typeName := t.String(m.ProtoPackage())
	// if strings.HasPrefix(typeName, "v1") {
	// 	log.Printf("typename=%s, protoPkg=%s", typeName, m.ProtoPackage())
	// }
	// log.Printf("funcChoicesInputParam: from=%s,protopkg=%s,typename=%s",
	// 	from, protoPkg, typeName)

	// if versionRegexp.MatchString(typeName) {
	// 	//typeName = t.String(versionRegexp.FindStringSubmatch(protoPkg)[1])
	// 	typeName = t.String(versionRegexp.FindStringSubmatch(typeName)[1])
	// 	log.Printf("--- new TypeName: %s, parts: %+v", typeName,
	// 		versionRegexp.FindStringSubmatch(typeName))
	// }
	return m.Language().ToTypeName(typeName, !wantDecl, m)
}

func funcChoicesOutputParam(_, b2, _, _ bool, m *codegen.WasmMethod) string {
	if !b2 {
		return ""
	}
	t := m.CGOutput().GetCGType()
	return m.Language().ToTypeName(t.String(m.ProtoPackage()), true, m)
}
func funcChoicesOutParamDecl(_, b2, _, _ bool, m *codegen.WasmMethod) string {
	if !b2 {
		return ""
	}
	t := m.CGOutput().GetCGType()
	return m.Language().ToTypeName(t.String(m.ProtoPackage()), false, m) + "{}"
}

func funcChoicesInputToSend(b1, _, b3, _ bool, m *codegen.WasmMethod) string {
	if b1 {
		return "req"
	}
	return "nil"
}

func funcChoicesNeedsPullApart(b1, b2, b3, b4 bool) bool {
	return funcChoicesToInt(b1, b2, b3, b4) == 0x5
}
func funcChoicesRetError(b1, b2, b3, b4 bool, m *codegen.WasmMethod) string {
	if !b2 {
		return ""
	}
	detail, scenario := outputTypeInfo(b4, m)
	switch scenario {
	case outTypeCompositeNoFields:
		return ""
	case outTypeComposite:
		return m.Language().EmptyComposite(detail.String(m.ProtoPackage()), m) + ","
	case outTypeBasic:
		return funcChoicesZeroValueRet(false, false, false, b4, false, m) + ","
	}

	//if b2 {
	//	return fmt.Sprintf("\"%s,err", m.CGInput().CGType().String(m.ProtoPackage()))
	//}
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
	return fmt.Sprintf("(%s,%s)", m.CGOutput().GetCGType().String(m.ProtoPackage()), "error")
}

func funcChoicesMethodRet(b1, b2, b3, b4 bool, abi bool, m *codegen.WasmMethod) string {
	return "error"
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
		return fmt.Sprintf("%sResponse{},nil", m.WasmMethodName())
	}
	return "nil"
}
func funcChoicesInbound(b1, _, _, _ bool, m *codegen.WasmMethod) string {
	if b1 {
		return "in"
	}
	return "nil"
}
func funcChoicesOutbound(b1, b2, b3, b4 bool, m *codegen.WasmMethod) string {
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
