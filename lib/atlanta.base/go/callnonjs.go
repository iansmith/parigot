//go:build !js
// +build !js

// This file is to make the compiler happy. These functions are only actually called from the WASM side but
// the compiler sees us compiling this package with the host system and expects these to have definitions.
// We put a panic each one just to be safe.
package lib

func locate(int32) {
	panic("locate nonjs")
}

func dispatch(int32) {
	panic("dispatch nonjs")
}

func bindMethod(int32) {
	panic("bindMethod nonjs")
}

func exit(int32) {
	panic("exit nonjs")
}

func blockUntilCall(int32) {
	panic("blockUntilCall nonjs")

}
func returnValue(int32) {
	panic("returnValue nonjs")
}

func export(int32) {
	panic("export nonjs")
}

func require(int32) {
	panic("require nonjs")

}

func run(int32) {
	panic("run nonjs")
}
