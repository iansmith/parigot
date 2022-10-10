package jspatch

import "github.com/iansmith/parigot/abi/client"

//go:wasm-module env
//export syscall/js.valueSetIndex
func ValueSetIndex(int64, int32, int64, int32) {
	client.JSNotImplemented("not implemented: js.valueSetIndex")
}

//go:wasm-module env
//export syscall/js.Value_.Int
func ValueInt(int64, int32, int32) int32 {
	client.JSNotImplemented("not implemented: js.value_.Int\n")
	return 0
}
