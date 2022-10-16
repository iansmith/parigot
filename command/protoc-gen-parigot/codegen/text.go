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
	// int64  xxx?
	OutType(m *WasmMethod) string
	// OutTypeDecl returns the declaration text for the output type.
	// Returns things like (for go):
	// (int64,error)
	OutTypeDecl(m *WasmMethod) string
	// ReturnErrorDecl returns the return statement plus a nil value for any parameters
	// except the error value. Returns things like
	// return int64(0),parigot.NewFromError("my error",err)
	ReturnErrorDecl(m *WasmMethod, msg string) string
	// ReturnValueDecl returns the return statement plus a fixed variable value
	// and a nil to indicate no error.
	ReturnValueDecl(m *WasmMethod) string
	// AllInputNumberedParam is used when we want to output the parameters and their
	// values, but we don't know how many parameters there are due to the wasm level's
	// parameters not matching the proto level.  Returns things like (in go):
	// p0 int32, p1 int32, p3 float64
	AllInputNumberedParam(m *WasmMethod) string // xxx dead code?
	// Given a particular language, how many parameters on WASM does it take
	// to encode this type.  Note: this is specific to the compiler.
	GetNumberParametersUsed(*CGType) int
	// GetForalArgSeparator returns the string that separates arguments in
	// declaration or call.  In go this is a comma.
	GetFormalArgSeparator() string
	// BasicTypeToString returns the language specific version of the input
	// or panics because it does not know how.  The caller should insure
	// that the value sent to this function is in fact a basic type like
	// TYPE_STRING or TYPE_INT32.  If the second parameter is true, this
	// function panics on unknown strings, which is usually what you want.
	BasicTypeToString(string, bool) string
	// Return an empty or initial value for a basic type.  This is used where
	// we have to create a "dummy" value for the type. Generally this value is
	// not going to be sued. Returns things like
	// int64(0)
	ZeroValuesForProtoTypes(string) string
	// Return the given value in the appropriate case/style for the language.
	// For go, this returns camel case identifiers. The second parameter is used
	// to indicate this is a parameter of a function that we are creating.
	// In go, the first letter of a parameter is not capitalized, as would normall
	// in camel case.
	ToId(string, bool) string
	// Table driven maker of choices
	FuncChoice() *FuncChooser
}

// AbiLanguageText is only for methods that are used by the ABI.  Since the ABI is currently
// implemented in go, only the abi "language" uses this and it generates go code.
// Other languages can safely ignore this.
type AbiLanguageText interface {
	LanguageText
	AllInputWithFormalWasmLevel(m *WasmMethod, showFormalName bool) string
	AllInputWasmToGoImpl(m *WasmMethod) string
}

type QuadOptions func(b1, b2, b3, b4 bool) bool
type QuadString func(b1, b2, b3, b4 bool) string
type QuadWithMethodString func(b1, b2, b3, b4 bool, m *WasmMethod) string
type FiveWithMethodString func(b1, b2, b3, b4, abi bool, m *WasmMethod) string
type FuncChooser struct {
	Bits           QuadString
	NeedsFill      QuadOptions
	NeedsInputPtr  QuadOptions
	NeedsPullApart QuadOptions
	Inbound        QuadString
	Outbound       QuadString
	RetError       QuadString
	RetValue       FiveWithMethodString
	MethodRet      FiveWithMethodString
	ZeroValueRet   QuadWithMethodString
}
