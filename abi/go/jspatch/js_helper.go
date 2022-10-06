package jspatch

import "github.com/iansmith/parigot/abi/go/abi"

//go:wasm-module env
//export syscall/js.valueGet
func ValueGet(int64, int32, int32, int32) int64 {
	abi.OutputString("not implemented: js.ValueGet\n")
	abi.JSNotImplemented()
	return 0
}

//go:wasm-module env
//export syscall/js.valuePrepareString
func ValuePrepareString(int32, int64, int32) {
	abi.OutputString("not implemented: js.valuePrepareString\n")
	abi.JSNotImplemented()
}

//go:wasm-module env
//export syscall/js.valueLoadString
func ValueLoadString(int64, int32, int32, int32, int32) {
	abi.OutputString("not implemented: js.valueLoadString\n")
	abi.JSNotImplemented()
}

//go:wasm-module env
//export syscall/js.finializeRef
func FinalizeRef(int64, int32) {
	abi.OutputString("not implemented: js.finalizeRef\n")
	abi.JSNotImplemented()
}

//go:wasm-module env
//export syscall/js.stringVal
func StringVal(int32, int32, int32) int64 {
	abi.OutputString("not implemented: js.stringVal\n")
	abi.JSNotImplemented()
	return 0
}

//go:wasm-module env
//export Syscall/js.valueSet
func ValueSet(int64, int32, int32, int64, int32) {
	abi.OutputString("not implemented: js.valueSet\n")
	abi.JSNotImplemented()
}

//go:wasm-module env
//export syscall/js.valueLength
func ValueLength(int64, int32) int32 {
	abi.OutputString("not implemented: js.valueLength\n")
	abi.JSNotImplemented()
	return 0
}

//go:wasm-module env
//export syscall/js.valueIndex
func ValueIndex(int64, int32, int32) int64 {
	abi.OutputString("not implemented: js.valueIndex\n")
	abi.JSNotImplemented()
	return 0
}

//go:wasm-module env
//export syscall/js.valueCall
func ValueCall(int32, int64, int32, int32, int32, int32, int32, int32) {
	abi.OutputString("not implemented: js.valueCall\n")
	abi.JSNotImplemented()
}

//go:wasm-module env
//export syscall/js.valueNew
func ValueNew(int32, int64, int32, int32, int32, int32) {
	abi.OutputString("not implemented: js.valueNew\n")
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
