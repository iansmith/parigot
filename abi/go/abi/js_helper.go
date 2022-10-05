package abi

//go:wasm-module env
//export syscall/js.valueGet
func ValueGet(int64, int32, int32, int32) int64 {
	JSNotImplemented()
	return 0
}

//go:wasm-module env
//export syscall/js.valuePrepareString
func ValuePrepareString(int32, int64, int32) {
	JSNotImplemented()
}

//go:wasm-module env
//export syscall/js.valueLoadString
func ValueLoadString(int64, int32, int32, int32, int32) {
	JSNotImplemented()
}

//go:wasm-module env
//export syscall/js.finializeRef
func FinalizeRef(int64, int32) {
	JSNotImplemented()
}

//go:wasm-module env
//export syscall/js.stringVal
func StringVal(int32, int32, int32) int64 {
	JSNotImplemented()
	return 0
}

//go:wasm-module env
//export Syscall/js.valueSet
func ValueSet(int64, int32, int32, int64, int32) {
	JSNotImplemented()
}

//go:wasm-module env
//export syscall/js.valueLength
func ValueLength() int32 {
	JSNotImplemented()
	return 0
}

//go:wasm-module env
//export syscall/js.valueIndex
func ValueIndex(int64, int32, int32) int64 {
	JSNotImplemented()
	return 0
}

//go:wasm-module env
//export syscall/js.valueCall
func ValueCall(int32, int64, int32, int32, int32, int32, int32, int32) {
	JSNotImplemented()
}
