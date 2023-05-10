package sys

import (
	"github.com/iansmith/parigot/eng"
)

func addSupportedFunctions(eng eng.Engine, w *WasmtimeSupportFunc) {
	rt := w.rt
	/*
		eng.AddSupportedFunc("env.syscall/js.valueSetIndex", rt.jsEnv.ValueSetIndex)
		eng.AddSupportedFunc("env.syscall/js.valuePrepareString", rt.jsEnv.ValuePrepareString)
		eng.AddSupportedFunc("env.syscall/js.valueLoadString", rt.jsEnv.ValueLoadString)

		eng.AddSupportedFunc("syscall/js.valueGet", rt.jsEnv.ValueGet)
		eng.AddSupportedFunc("syscall/js.finalizeRef", rt.jsEnv.FinalizeRef)
		eng.AddSupportedFunc("syscall/js.stringVal", rt.jsEnv.StringVal)
		eng.AddSupportedFunc("syscall/js.valueSet", rt.jsEnv.ValueSet)
		eng.AddSupportedFunc("syscall/js.valueDelete", rt.jsEnv.ValueDelete)
		eng.AddSupportedFunc("syscall/js.valueIndex", rt.jsEnv.ValueIndex)
		eng.AddSupportedFunc("syscall/js.valueSetIndex", rt.jsEnv.ValueSetIndex)
		eng.AddSupportedFunc("syscall/js.valueLength", rt.jsEnv.ValueLength)
		eng.AddSupportedFunc("syscall/js.valuePrepareString", rt.jsEnv.ValuePrepareString)
		eng.AddSupportedFunc("syscall/js.valueCall", rt.jsEnv.ValueCall)
		eng.AddSupportedFunc("syscall/js.valueInvoke", rt.jsEnv.ValueInvoke)
		eng.AddSupportedFunc("syscall/js.valueNew", rt.jsEnv.ValueNew)
		eng.AddSupportedFunc("syscall/js.valueLoadString", rt.jsEnv.ValueLoadString)
		eng.AddSupportedFunc("syscall/js.valueInstanceOf", rt.jsEnv.ValueInstanceOf)
		eng.AddSupportedFunc("syscall/js.copyBytesToGo", rt.jsEnv.CopyBytesToGo)
		eng.AddSupportedFunc("syscall/js.copyBytesToJS", rt.jsEnv.CopyBytesToJS)

		eng.AddSupportedFunc("runtime.resetMemoryDataView", rt.runtimeEnv.ResetMemoryDataView)
		eng.AddSupportedFunc("runtime.wasmExit", rt.wasiEnv.WasiExit)
		eng.AddSupportedFunc("runtime.wasmWrite", rt.wasiEnv.WasiWrite)
		eng.AddSupportedFunc("runtime.nanotime1", rt.runtimeEnv.Nanotime1)
		eng.AddSupportedFunc("runtime.walltime", rt.runtimeEnv.WallTime)
		eng.AddSupportedFunc("runtime.scheduleTimeoutEvent", rt.runtimeEnv.ScheduleTimeoutEvent)
		eng.AddSupportedFunc("runtime.clearTimeoutEvent", rt.runtimeEnv.ClearTimeoutEvent)
		eng.AddSupportedFunc("runtime.getRandomData", rt.runtimeEnv.GetRandomData)

		//eng.AddSupportedFunc("parigot.debugprint", wrapWithRecover(store, rt.syscall.DebugPrint)
		eng.AddSupportedFunc("debug", rt.runtimeEnv.GoDebug)
	*/
	//system calls
	eng.AddSupportedFunc("parigot", "locate_", rt.syscall.Locate)

	// xxx fix me: How are we going to clean up the resources for a particular service when it exits?
	// how do we find the resources associated with the caller of exit().

	eng.AddSupportedFunc("parigot", "exit_", rt.syscall.Exit)
	eng.AddSupportedFunc("parigot", "dispatch_", rt.syscall.Dispatch)
	eng.AddSupportedFunc("parigot", "bind_method_", rt.syscall.BindMethod)
	eng.AddSupportedFunc("parigot", "return_value_", rt.syscall.ReturnValue)
	eng.AddSupportedFunc("parigot", "block_until_call_", rt.syscall.BlockUntilCall)
	eng.AddSupportedFunc("parigot", "require_", rt.syscall.Require)
	eng.AddSupportedFunc("parigot", "export_", rt.syscall.Export)
	eng.AddSupportedFunc("parigot", "run_", rt.syscall.Run)

	//backdoor to the logger
	eng.AddSupportedFunc("parigot", "backdoor_log_", rt.syscall.BackdoorLog)
}

func addSplitModeFunctions(eng eng.Engine, w *WasmtimeSupportFunc) {
	// split mode
	// mixed mode entries: this should be automated (xxxfixmexxx)
	eng.AddSupportedFunc("logviewer", "log_request_handler", w.log.LogRequestHandler)
	// mixed mode entries: this should be automated (xxxfixmexxx)
	eng.AddSupportedFunc("filesvc", "open", w.file.FileSvcOpen)
	eng.AddSupportedFunc("filesvc", "load", w.file.FileSvcLoad)

	eng.AddSupportedFunc("queuesvc", "create_handler", w.queue.QueueSvcCreateQueue)
	eng.AddSupportedFunc("queuesvc", "delete_handler", w.queue.QueueSvcDeleteQueue)
	eng.AddSupportedFunc("queuesvc", "length_handler", w.queue.QueueSvcLength)
	eng.AddSupportedFunc("queuesvc", "locate_handler", w.queue.QueueSvcLocate)
	eng.AddSupportedFunc("queuesvc", "mark_done_handler", w.queue.QueueSvcMarkDone)
	eng.AddSupportedFunc("queuesvc", "send_handler", w.queue.QueueSvcSend)
	eng.AddSupportedFunc("queuesvc", "receive_handler", w.queue.QueueSvcReceive)
}

// func wra	pWithRecover(store wasmtime.Storelike, fn func(int32), name string) *wasmtime.Func {
// 	return wasmtime.WrapFunc(store, func(sp int32) {
// 		defer func() {
// 			if r := recover(); r != nil {
// 				print("START RECOVER STACKTRACE " + name + "\n")
// 				print(fmt.Sprintf("trapped panic: %s %v (%T)\n", name, r, r))
// 				debug.PrintStack()
// 				print("END RECOVER+STACKTRACE " + name + "\n")
// 			}
// 		}()
// 		if sp == 0 {
// 			panic("BAD SP")
// 		}
// 		fn(sp)
// 	})
// }
