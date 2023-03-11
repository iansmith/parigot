//go:build js && browser

package syscall

import _ "unsafe"

//go:noinline
//go:linkname locate go:parigot.locate_
func locate(int32) {
	print("browser parigot.locate_()\n")
}

//go:noinline
//go:linkname dispatch parigot.dispatch_
func dispatch(int32) {
	print("browser parigot.dispatch_()\n")
}

//go:noinline
//go:linkname bindMethod parigot.bind_method_
func bindMethod(int32) {
	print("browser parigot.bind_method_()\n")
}

//go:noinline
//go:linkname exit parigot.exit_
func exit(int32) {
	print("browser parigot.exit_()\n")
}

//go:noinline
//go:linkname blockUntilCall parigot.block_until_call_
func blockUntilCall(int32) {
	print("browser parigot.block_until_call_()\n")
}

//go:noinline
//go:linkname returnValue parigot.return_value_
func returnValue(int32) {
	print("browser parigot.return_value_\n")
}

//go:noinline
//go:linkname export parigot.export_
func export(int32) {
	print("browser parigot.export_\n")
}

//go:noinline
//go:linkname require parigot.require_
func require(int32) {
	print("browser parigot.require_\n")
}

//go:noinline
//go:linkname run parigot.run_
func run(int32) {
	print("browser parigot.run_\n")
}

//go:noinline
//go:linkname backdoorLog parigot.backdoor_log_
func backdoorLog(int32) {
	print("browser parigot.backdoor_log_\n")
}
