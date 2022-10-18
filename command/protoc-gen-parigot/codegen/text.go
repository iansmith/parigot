package codegen

// These interfaces are to specify the code generation functions that need
// to implemented to use the primitives defined in this package. You'll need
// an implementation of this to port to a new language.

type LanguageText interface {
	NilValue() string
	// EmptyComposite takes a type name and returns the empty value of that
	// composite. This is the brother of the ZeroValuesForProtoTypes that does
	// the same thing basic types.  This returns Foo{} for the input of Foo in
	// go.  The second values is only for people doing their own namespacing.
	EmptyComposite(typeName string, method *WasmMethod) string
	// GetNumberParametersUsed returns how many parameters on WASM does it take
	// to encode this type.  Note: this is specific to the compiler.
	GetNumberParametersUsed(*CGType) int
	// GetFormalArgSeparator returns the string that separates arguments in
	// declaration or call.  In go this is a comma. e.g foo(a,b,c) because
	// comma.
	GetFormalArgSeparator() string
	// GetFormalTypeCombination returns has the combination of a formal declaration
	// and a type. It can produce its output in any order or drop the type in
	// entirely in languages that don't require.  Note that the text of these
	// two values have already been processed. For ("a","foo") go returns "a foo".
	GetFormalTypeCombination(formal, typ string) string
	// BasicTypeToString returns the language specific version of the input
	// or panics because it does not know how.  The caller should insure
	// that the value sent to this function is in fact a basic type like
	// TYPE_STRING or TYPE_INT32.  If the second parameter is true, this
	// function panics on unknown strings, which is usually what you want.
	BasicTypeToString(string, bool) string
	// BasicTypeToWasm returns the language and compiler specific mapping
	// of basic types to a sequence of the 4 wasm types.  Note that the
	// return value here is still in the form "TYPE_INT32" that is from
	// proto.  To get the final type, use BasicTypeToString.
	BasicTypeToWasm(string) []string
	// ZeroValuesForProtoTypes returns an empty or initial value for a basic type.  This is used where
	// we have to create a "dummy" value for the type. Generally this value is
	// not going to be sued. Returns things like
	// int64(0)
	ZeroValuesForProtoTypes(string) string
	// ToId returns the given value in the appropriate case/style for the language.
	// For go, this returns camel case identifiers. The second parameter is used
	// to indicate this is a parameter of a function that we are creating.
	// The last parameter is only of interest to folks doing their own namespacing;
	// the method you can get the method name and (via the parent) the service name.
	// In go, the first letter of a parameter is not capitalized, as would normall
	// in camel case.
	ToId(string, bool, *WasmMethod) string
	// ToTypeName returns the given value in the appropriate case/style for the language.
	// For go, this is lower case for the components of the string before the
	// last one, then camel case in the last one.  The second parameter indicates
	// that if possible the caller wants a reference type (a pointer).
	// The final argument is only of interest to folks doing their own namespacing;
	// you can get the method name and (via the parent) the service name.
	// For example, the version converts Foo.Bar.Baz=>foo.bar.Baz
	// and foo.bar.baz_fleazil => foo.bar.BazFleazil
	ToTypeName(string, bool, *WasmMethod) string
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
