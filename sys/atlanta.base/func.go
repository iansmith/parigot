package sys

import (
	wasmtime "github.com/bytecodealliance/wasmtime-go"
)

func addSupportedFunctions(store wasmtime.Storelike,
	result map[string]*wasmtime.Func,
	rt *Runtime) {
	result["env.syscall/js.valueSetIndex"] = wasmtime.WrapFunc(store, rt.jsEnv.ValueSetIndex)
	result["go.syscall/js.valueGet"] = wasmtime.WrapFunc(store, rt.jsEnv.ValueGet)
	result["env.syscall/js.valuePrepareString"] = wasmtime.WrapFunc(store, rt.jsEnv.ValuePrepareString)
	result["env.syscall/js.valueLoadString"] = wasmtime.WrapFunc(store, rt.jsEnv.ValueLoadString)
	result["go.syscall/js.finalizeRef"] = wasmtime.WrapFunc(store, rt.jsEnv.FinalizeRef)
	result["go.syscall/js.stringVal"] = wasmtime.WrapFunc(store, rt.jsEnv.StringVal)
	result["go.syscall/js.valueSet"] = wasmtime.WrapFunc(store, rt.jsEnv.ValueSet)
	result["go.syscall/js.valueDelete"] = wasmtime.WrapFunc(store, rt.jsEnv.ValueDelete)
	result["go.syscall/js.valueIndex"] = wasmtime.WrapFunc(store, rt.jsEnv.ValueIndex)
	result["go.syscall/js.valueSetIndex"] = wasmtime.WrapFunc(store, rt.jsEnv.ValueSetIndex)
	result["go.syscall/js.valueLength"] = wasmtime.WrapFunc(store, rt.jsEnv.ValueLength)
	result["go.syscall/js.valuePrepareString"] = wasmtime.WrapFunc(store, rt.jsEnv.ValuePrepareString)
	result["go.syscall/js.valueCall"] = wasmtime.WrapFunc(store, rt.jsEnv.ValueCall)
	result["go.syscall/js.valueInvoke"] = wasmtime.WrapFunc(store, rt.jsEnv.ValueInvoke)
	result["go.syscall/js.valueNew"] = wasmtime.WrapFunc(store, rt.jsEnv.ValueNew)
	result["go.syscall/js.valueLoadString"] = wasmtime.WrapFunc(store, rt.jsEnv.ValueLoadString)
	result["go.syscall/js.valueInstanceOf"] = wasmtime.WrapFunc(store, rt.jsEnv.ValueInstanceOf)
	result["go.syscall/js.copyBytesToGo"] = wasmtime.WrapFunc(store, rt.jsEnv.CopyBytesToGo)
	result["go.syscall/js.copyBytesToJS"] = wasmtime.WrapFunc(store, rt.jsEnv.CopyBytesToJS)
	result["go.runtime.resetMemoryDataView"] = wasmtime.WrapFunc(store, rt.runtimeEnv.ResetMemoryDataView)
	result["go.runtime.wasmExit"] = wasmtime.WrapFunc(store, rt.wasiEnv.WasiExit)
	result["go.runtime.wasmWrite"] = wasmtime.WrapFunc(store, rt.wasiEnv.WasiWrite)
	result["go.runtime.nanotime1"] = wasmtime.WrapFunc(store, rt.runtimeEnv.Nanotime1)
	result["go.runtime.walltime"] = wasmtime.WrapFunc(store, rt.runtimeEnv.WallTime)
	result["go.runtime.scheduleTimeoutEvent"] = wasmtime.WrapFunc(store, rt.runtimeEnv.ScheduleTimeoutEvent)
	result["go.runtime.clearTimeoutEvent"] = wasmtime.WrapFunc(store, rt.runtimeEnv.ClearTimeoutEvent)
	result["go.runtime.getRandomData"] = wasmtime.WrapFunc(store, rt.runtimeEnv.GetRandomData)
	//result["parigot.debugprint"] = wasmtime.WrapFunc(store, rt.syscall.DebugPrint)
	result["go.debug"] = wasmtime.WrapFunc(store, rt.runtimeEnv.GoDebug)

	//system calls
	result["go.parigot.locate_"] = wasmtime.WrapFunc(store, rt.syscall.Locate)
	result["go.parigot.register_"] = wasmtime.WrapFunc(store, rt.syscall.Register)
	// xxx fix me: How are we going to clean up the resources for a particular service when it exits?
	// how do we find the resources associated with the caller of exit().
	result["go.parigot.exit_"] = wasmtime.WrapFunc(store, rt.syscall.Exit)
	result["go.parigot.dispatch_"] = wasmtime.WrapFunc(store, rt.syscall.Dispatch)
	result["go.parigot.bind_method_"] = wasmtime.WrapFunc(store, rt.syscall.BindMethod)
	result["go.parigot.return_value_"] = wasmtime.WrapFunc(store, rt.syscall.ReturnValue)
	result["go.parigot.block_until_call_"] = wasmtime.WrapFunc(store, rt.syscall.BlockUntilCall)
	result["go.parigot.require_"] = wasmtime.WrapFunc(store, rt.syscall.Require)
	result["go.parigot.export_"] = wasmtime.WrapFunc(store, rt.syscall.Export)
	result["go.parigot.run_"] = wasmtime.WrapFunc(store, rt.syscall.Run)
}
