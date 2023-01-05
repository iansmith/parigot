package sys

import (
	"runtime/debug"

	// fileimpl "github.com/iansmith/parigot/api_impl/file"
	// logimpl "github.com/iansmith/parigot/api_impl/log"
	filego "github.com/iansmith/parigot/api_impl/file/go_"
	loggo "github.com/iansmith/parigot/api_impl/log/go_"

	wasmtime "github.com/bytecodealliance/wasmtime-go/v3"
)

func addSupportedFunctions(store wasmtime.Storelike, result map[string]*wasmtime.Func, rt *Runtime) {
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
	result["go.runtime.resetMemoryDataView"] = wrapWithRecover(store, rt.runtimeEnv.ResetMemoryDataView)
	result["go.runtime.wasmExit"] = wrapWithRecover(store, rt.wasiEnv.WasiExit)
	result["go.runtime.wasmWrite"] = wrapWithRecover(store, rt.wasiEnv.WasiWrite)
	result["go.runtime.nanotime1"] = wrapWithRecover(store, rt.runtimeEnv.Nanotime1)
	result["go.runtime.walltime"] = wrapWithRecover(store, rt.runtimeEnv.WallTime)
	result["go.runtime.scheduleTimeoutEvent"] = wrapWithRecover(store, rt.runtimeEnv.ScheduleTimeoutEvent)
	result["go.runtime.clearTimeoutEvent"] = wrapWithRecover(store, rt.runtimeEnv.ClearTimeoutEvent)
	result["go.runtime.getRandomData"] = wrapWithRecover(store, rt.runtimeEnv.GetRandomData)
	//result["parigot.debugprint"] = wrapWithRecover(store, rt.syscall.DebugPrint)
	result["go.debug"] = wrapWithRecover(store, rt.runtimeEnv.GoDebug)

	//system calls
	result["go.parigot.locate_"] = wrapWithRecover(store, rt.syscall.Locate)
	//result["go.parigot.register_"] = wrapWithRecover(store, rt.syscall.Register)
	// xxx fix me: How are we going to clean up the resources for a particular service when it exits?
	// how do we find the resources associated with the caller of exit().
	result["go.parigot.exit_"] = wrapWithRecover(store, rt.syscall.Exit)
	result["go.parigot.dispatch_"] = wrapWithRecover(store, rt.syscall.Dispatch)
	result["go.parigot.bind_method_"] = wrapWithRecover(store, rt.syscall.BindMethod)
	result["go.parigot.return_value_"] = wrapWithRecover(store, rt.syscall.ReturnValue)
	result["go.parigot.block_until_call_"] = wrapWithRecover(store, rt.syscall.BlockUntilCall)
	result["go.parigot.require_"] = wrapWithRecover(store, rt.syscall.Require)
	result["go.parigot.export_"] = wrapWithRecover(store, rt.syscall.Export)
	result["go.parigot.run_"] = wrapWithRecover(store, rt.syscall.Run)

	//backdoor to the logger
	result["go.parigot.backdoor_log_"] = wrapWithRecover(store, rt.syscall.BackdoorLog)
}

func addSplitModeFunctions(store wasmtime.Storelike,
	result map[string]*wasmtime.Func,
	logViewer *loggo.LogViewerImpl,
	fileSvc *filego.FileSvcImpl) {

	// mixed mode entries: this should be automated (xxxfixmexxx)
	result["go.logviewer.log_request_handler"] = wrapWithRecover(store, logViewer.LogRequestHandler)
	// mixed mode entries: this should be automated (xxxfixmexxx)
	result["go.filesvc.open"] = wrapWithRecover(store, fileSvc.FileSvcOpen)
	result["go.filesvc.load"] = wrapWithRecover(store, fileSvc.FileSvcLoad)

}

func wrapWithRecover(store wasmtime.Storelike, fn func(int32)) *wasmtime.Func {
	return wasmtime.WrapFunc(store, func(sp int32) {
		defer func() {
			if r := recover(); r != nil {
				print("RECOVER FROM PANIC\n")
				debug.PrintStack()
				print("END RECOVER+STACKTRACE\n")
			}
			return
		}()
		fn(sp)
	})
}
