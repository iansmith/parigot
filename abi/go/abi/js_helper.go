package abi

//go:wasm-module env
//export syscall/js.valueGet
func valueGet(int64, int32, int32, int32) int64 {
	JSNotImplemented()
	return 0
}

//go:wasm-module env
//export syscall/js.valuePrepareString
func valuePrepareString(int32, int64, int32) {
	JSNotImplemented()
}

//go:wasm-module env
//export syscall/js.valueLoadString
func valueLoadString(int64, int32, int32, int32, int32) {
	JSNotImplemented()
}

//go:wasm-module env
//export syscall/js.finializeRef
func finalizeRef(int64, int32) {
	JSNotImplemented()
}

//go:wasm-module env
//export syscall/js.stringVal
func stringVal(int32, int32, int32) int64 {
	JSNotImplemented()
	return 0
}

//go:wasm-module env
//export syscall/js.valueSet
func valueSet(int64, int32, int32, int64, int32) {
	JSNotImplemented()
}

//go:wasm-module env
//export syscall/js.valueLength
func valueLength() int32 {
	JSNotImplemented()
	return 0
}

//go:wasm-module env
//export syscall/js.valueIndex
func valueIndex(int64, int32, int32) int64 {
	JSNotImplemented()
	return 0
}

//go:wasm-module env
//export syscall/js.valueCall
func valueCall(int32, int64, int32, int32, int32, int32, int32, int32) {
	JSNotImplemented()
}
