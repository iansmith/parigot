package codegen

// These interfaces are to specify the code generation functions that need
// to implemented to use the primitives defined in this package. You'll need
// an implementation of this to port to a new language.

type LanguageText interface {
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
	// OutZeroValueDecl returns the return value for the "nothing" case for the particular output
	// type. Returns things like (for go):
	// return int64,error
	OutZeroValueDecl(m *WasmMethod) string
	// OutType returns the declaration text for the output type.
	// Returns things like (for go):
	// (int64,error)
	OutType(m *WasmMethod) string
	// OutTypeDecl returns the output type declaration appropriate for the language of the output type.
	// Returns things like (for go):
	// (int64,error)
	OutTypeDecl(m *WasmMethod) string
	// AllInputNumberedParam is used when we want to output the parameters and their
	// values, but we don't know how many parameters there are due to the wasm level's
	// parameters not matching the proto level.  Returns things like (in go):
	// p0 int32, p1 int32, p3 float64
	AllInputNumberedParam(m *WasmMethod) string // xxx dead code?

	// Given a particular language, how many parameters on WASM does it take
	// to encode this type
	GetNumberParametersUsed(*CGType) int

	// GetForalArgSeparator returns the string that separates arguments in
	// declaration or call.
	GetFormalArgSeparator() string

	// BasicTypeToString returns the language specific version of the input
	// or panics because it does not know how.  The calller should insure
	// that the value sent to this function is in fact a basic type like
	// TYPE_STRING or TYPE_INT32.  If the second parameter is true, this
	// function panics on unknown strings, which is usually what you want.
	BasicTypeToString(string, bool) string
}

// AbiLanguageText is only for methods that are used by the ABI.  Since the ABI is currently
// implemented in go, only the abi "language" uses this and it generates go code.
// Other languages can safely ignore this.
type AbiLanguageText interface {
	LanguageText
	AllInputWithFormalWasmLevel(m *WasmMethod, showFormalName bool) string
	AllInputWasmToGoImpl(m *WasmMethod) string
}
