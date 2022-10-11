package jspatch

//go:wasm-module env
//export syscall/js.valueSetIndex
func ValueSetIndex(int64, int32, int64, int32) {
	abi.JSNotImplemented("not implemented: js.valueSetIndex")
}

//go:wasm-module env
//export syscall/js.Value_.Int
func ValueInt(int64, int32, int32) int32 {
	abi.JSNotImplemented("not implemented: js.value_.Int\n")
	return 0
}
