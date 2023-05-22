package sys

import (
	"context"

	"github.com/iansmith/parigot/eng"

	"github.com/tetratelabs/wazero/api"
)

func go_debug(ctx context.Context, m api.Module, stack []uint64) {
	panic("go_debug")
}

func go_runtime_resetMemoryDataView(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.runtime.resetMemoryDataView")
}

func go_runtime_wasm_exit(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.runtime.wasmExit")
}

func go_runtime_wasm_write(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.runtime.wasmWrite")
}

func go_runtime_nanotime1(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.runtime.nanotime1")
}

func go_runtime_walltime(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.runtime.walltime")
}
func go_runtime_scheduleTimeoutEvent(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.runtime.scheduleTimeoutEvent")
}
func go_runtime_clearTimeoutEvent(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.runtime.clearTimeoutEvent")
}
func go_runtime_getRandomData(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.runtime.getRandomData")
}
func go_syscalljs_finalizeRef(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.finalizeRef")
}
func go_syscalljs_stringVal(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.stringVal")
}
func go_syscalljs_valueGet(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.valueGet")
}
func go_syscalljs_valueSet(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.valueSet")
}
func go_syscalljs_valueDelete(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.valueDelete")
}
func go_syscalljs_valueIndex(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.valueIndex")
}
func go_syscalljs_valueCall(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.valueCall")
}
func go_syscalljs_valueSetIndex(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.valueSetIndex")
}
func go_syscalljs_valueInvoke(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.valueInvoke")
}
func go_syscalljs_valueNew(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.valueNew")
}
func go_syscalljs_valueLength(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.valueLength")
}
func go_syscalljs_valuePrepareString(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.valuePrepareString")
}
func go_syscalljs_valueLoadString(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.valueLoadString")
}
func go_syscalljs_valueInstanceOf(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.InstanceOf")
}
func go_syscalljs_copyBytesToGo(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.copyBytesToGo")
}
func go_syscalljs_copyBytesToJS(ctx context.Context, m api.Module, stack []uint64) {
	panic("go.syscall/js.copyBytesToJS")
}

var bg = context.Background()

func InitializePatch(e eng.Engine) {
	e.AddSupportedFunc(bg, "gojs", "debug", go_debug)

	e.AddSupportedFunc_i32_v(bg, "gojs", "runtime.resetMemoryDataView", go_runtime_resetMemoryDataView)
	e.AddSupportedFunc_i32_v(bg, "gojs", "runtime.wasmWrite", go_runtime_wasm_write)
	e.AddSupportedFunc_i32_v(bg, "gojs", "runtime.wasmExit", go_runtime_wasm_exit)
	e.AddSupportedFunc_i32_v(bg, "gojs", "runtime.nanotime1", go_runtime_nanotime1)
	e.AddSupportedFunc_i32_v(bg, "gojs", "runtime.scheduleTimeoutEvent", go_runtime_scheduleTimeoutEvent)
	e.AddSupportedFunc_i32_v(bg, "gojs", "runtime.walltime", go_runtime_walltime)
	e.AddSupportedFunc_i32_v(bg, "gojs", "runtime.clearTimeoutEvent", go_runtime_clearTimeoutEvent)
	e.AddSupportedFunc_i32_v(bg, "gojs", "runtime.getRandomData", go_runtime_getRandomData)

	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.finalizeRef", go_syscalljs_finalizeRef)
	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.stringVal", go_syscalljs_stringVal)
	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.valueGet", go_syscalljs_valueGet)
	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.valueSet", go_syscalljs_valueSet)
	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.valueDelete", go_syscalljs_valueDelete)
	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.valueIndex", go_syscalljs_valueIndex)
	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.valueSetIndex", go_syscalljs_valueSetIndex)
	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.valueCall", go_syscalljs_valueCall)
	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.valueInvoke", go_syscalljs_valueInvoke)
	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.valueNew", go_syscalljs_valueNew)
	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.valueLength", go_syscalljs_valueLength)
	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.valuePrepareString", go_syscalljs_valuePrepareString)
	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.valueInstanceOf", go_syscalljs_valueInstanceOf)
	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.valueLoadString", go_syscalljs_valueLoadString)
	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.copyBytesToGo", go_syscalljs_copyBytesToGo)
	e.AddSupportedFunc_i32_v(bg, "gojs", "syscall/js.copyBytesToJS", go_syscalljs_copyBytesToJS)
}
