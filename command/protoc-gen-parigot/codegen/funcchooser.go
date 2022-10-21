package codegen

type QuadOptions func(b1, b2, b3, b4 bool) bool
type QuadWithMethodBool func(b1, b2, b3, b4 bool, m *WasmMethod) bool
type QuadString func(b1, b2, b3, b4 bool) string
type QuadWithMethodString func(b1, b2, b3, b4 bool, m *WasmMethod) string
type FiveWithMethodString func(b1, b2, b3, b4, abi bool, m *WasmMethod) string

type FuncChooser struct {
	Bits                QuadString
	NeedsFill           QuadOptions
	NeedsRet            QuadOptions
	InputParam          QuadWithMethodString
	NeedsPullApart      QuadOptions
	Inbound             QuadString
	Outbound            QuadString
	RetError            QuadWithMethodString
	RetValue            QuadWithMethodString
	MethodRet           FiveWithMethodString
	ZeroValueRet        FiveWithMethodString
	MethodParamDecl     QuadWithMethodString
	OutLocal            QuadWithMethodString
	MethodCall          QuadWithMethodString
	DecodeRequired      QuadWithMethodBool
	NoDecodeRequired    QuadWithMethodBool
	UsesReturnValuePtr  QuadWithMethodBool
	MethodParamDeclWasm QuadWithMethodString
	HasComplexParam     QuadWithMethodBool
	MethodCallWasm      QuadWithMethodString
	InputToSend         QuadWithMethodString
}
