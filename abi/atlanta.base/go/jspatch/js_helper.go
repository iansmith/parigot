package jspatch

import "github.com/iansmith/parigot/abi/atlanta1/base/go/parigot/abi"

C	abi.OutputString("not implemented: js.valueNew\n")
	abi.JSNotImplemented()
}

//go:wasm-module env
//export syscall/js.valueSetIndex
func ValueSetIndex(int64, int32, int64, int32) {
	abi.OutputString("not implemented: js.valueSetIndex\n")
	abi.JSNotImplemented()
}

//go:wasm-module env
//export syscall/js.Value_.Int
func ValueInt(int64, int32, int32) int32 {
	abi.OutputString("not implemented: js.value_.Int\n")
	abi.JSNotImplemented()
	return 0
}
