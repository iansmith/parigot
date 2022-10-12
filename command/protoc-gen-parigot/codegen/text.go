package codegen

// These interfaces are to specify the code generation functions that need
// to implemented to use the primitives defined in this package. You'll need
// an implementation of this to port to a new language.

type LanguageText interface {
	// ProtoTypeNameToLanguageTypeName takes something like "TYPE_STRING", "TYPE_INT32" or "TYPE_BOOL"
	// and convert it to the appropriate string for the language.
	ProtoTypeNameToLanguageTypeName(string) string
	// AllInputWithFormal returns the input variables of the method with the option to change their
	// formally declared name to "_".  Returns things like (in go):
	// foo typeOfFoo, bar typeOfBar
	AllInputWithFormal(method *WasmMethod, showFormalName bool) string
	// AllInputFormal returns just the names of the formals for this method.  Returns things like (in go)
	// param1, param2, param3
	AllInputFormal(method *WasmMethod) string
	// OutZeroValue is used thwart compilers that can detect reaching the end of a
	// function with no return value.  This is called to generate a dummy value in
	// the current language for the type provided.  This value will never be used.
	// Returns things like (in go):
	// int32(0)
	OutZeroValue(m *WasmMethod) string
	// OutType returns the type name appropriate for the language of the output type.
	// Returns things like (for go):
	// int64
	OutType(m *WasmMethod) string
	// AllInputNumberedParam is used when we want to output the parameters and their
	// values, but we don't know how many parameters there are due to the wasm level's
	// parameters not matching the proto level.  Returns things like (in go):
	// p0 int32, p1 int32, p3 float64
	AllInputNumberedParam(m *WasmMethod) string // xxx dead code?
}

// AbiLanguageText is only for methods that are used by the ABI.  Since the ABI is currently
// implemented in go, only the abi "language" uses this and it generates go code.
// Other languages can safely ignore this.
type AbiLanguageText interface {
	LanguageText
	AllInputParamWithFormalWasmLevel(m *WasmMethod, showFormalName bool) string
	AllInputParamWasmToGoImpl(m *WasmMethod) string
}
