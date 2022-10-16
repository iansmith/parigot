package go_

import (
	"fmt"
	"github.com/iansmith/parigot/command/protoc-gen-parigot/codegen"
	"log"
)

func (g *GoText) FuncChoice() *codegen.FuncChooser {
	return &codegen.FuncChooser{
		Bits:           funcChoicesToBits,
		NeedsFill:      funcChoicesNeedsFill,
		NeedsInputPtr:  funcChoicesNeedsInputPtr,
		NeedsPullApart: funcChoicesNeedsPullApart,
		Inbound:        funcChoicesInbound,
		Outbound:       funcChoicesOutbound,
		RetError:       funcChoicesRetError,
		RetValue:       funcChoicesRetValue,
		MethodRet:      funcChoicesMethodRet,
		ZeroValueRet:   funcChoicesZeroValueRet,
	}
}

func funcChoicesZeroValueRet(_, b2, _, b4 bool, method *codegen.WasmMethod) string {
	if !b2 {
		return "nil"
	}
	t := method.GetCGOutput().GetCGType()
	if b4 {
		v := ""
		if t.IsCompositeNoFields() {
			return "nil"
		}
		f := t.GetCompositeType().GetField()
		inner := codegen.NewCGTypeFromField(f[0], method, method.GetProtoPackage())
		if inner.IsBasic() {
			v = method.GetLanguage().ZeroValuesForProtoTypes(inner.String(""))
		} else {
			v = inner.String(method.GetProtoPackage()) + "Response{}"
		}
		return v + ",nil"
	}
	return t.String(method.GetProtoPackage()) + "Response{}" + ",nil"
}
func funcChoicesNeedsFill(b1, b2, b3, b4 bool) bool {
	return funcChoicesToInt(b1, b2, b3, b4) == 0xa
}
func funcChoicesNeedsInputPtr(b1, b2, b3, b4 bool) bool {
	return funcChoicesToInt(b1, b2, b3, b4) == 0x8
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

func outReturnToString(m *codegen.WasmMethod, abi bool) string {
	field := m.GetCGOutput().GetCGType().GetCompositeType().GetField()
	if len(field) > 1 {
		log.Fatalf("unable to pull up output for method %s because it has more than 1 result", m.GetWasmMethodName())
	}
	inner := codegen.NewCGTypeFromField(field[0], m, m.GetProtoPackage())
	if abi && !inner.IsBasic() {
		panic("all abi functions should be returning basic values")
	}

	basic := inner.String(m.GetProtoPackage())
	return m.GetLanguage().BasicTypeToString(basic, true)
}

func funcChoicesMethodRet(b1, b2, b3, b4 bool, abi bool, m *codegen.WasmMethod) string {
	if abi {
		// sadly, the abi is a special case because it doesn't return error values
		if b2 {
			//b2 implies b4 for the abi
			return outReturnToString(m, true)
		}
		return ""
	}
	if b2 {
		if !b4 {
			return fmt.Sprintf("(*%sResponse,error)", codegen.ToCamelCase(m.GetWasmMethodName()))
		}
		basic := outReturnToString(m, false)
		return fmt.Sprintf("(%s,error)", basic)
	}
	return "error"
}

func funcChoicesRetValue(b1, b2, b3, b4, abi bool, m *codegen.WasmMethod) string {
	if abi {
		// sadly, the abi is a special case because it doesn't return error values
		t := m.GetCGOutput().GetCGType()
		if b2 {
			if b4 == false {
				panic("all the ABI functions must be pulling up their return values")
			}
			if t.IsCompositeNoFields() {
				return ""
			}
			f := t.GetCompositeType().GetField()
			field := f[0]
			inner := codegen.NewCGTypeFromField(field, m, m.GetProtoPackage())
			if !inner.IsBasic() {
				panic("all abi functions should be returning basic values")
			}
			return m.GetLanguage().BasicTypeToString(inner.String(""), true)
		}
	}
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
