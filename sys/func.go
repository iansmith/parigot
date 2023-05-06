package sys

import (
	"fmt"
	"runtime/debug"

	filego "github.com/iansmith/parigot/apiimpl/file/go_"
	loggo "github.com/iansmith/parigot/apiimpl/log/go_"
	queuego "github.com/iansmith/parigot/apiimpl/queue/go_"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
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
	result["go.runtime.resetMemoryDataView"] = wrapWithRecover(store, rt.runtimeEnv.ResetMemoryDataView, "ResetMemoryDataView")
	result["go.runtime.wasmExit"] = wrapWithRecover(store, rt.wasiEnv.WasiExit, "WasiExit")
	result["go.runtime.wasmWrite"] = wrapWithRecover(store, rt.wasiEnv.WasiWrite, "WasiWrite")
	result["go.runtime.nanotime1"] = wrapWithRecover(store, rt.runtimeEnv.Nanotime1, "Nanotime1")
	result["go.runtime.walltime"] = wrapWithRecover(store, rt.runtimeEnv.WallTime, "WallTime")
	result["go.runtime.scheduleTimeoutEvent"] = wrapWithRecover(store, rt.runtimeEnv.ScheduleTimeoutEvent, "ScheduleTimeoutEvent")
	result["go.runtime.clearTimeoutEvent"] = wrapWithRecover(store, rt.runtimeEnv.ClearTimeoutEvent, "ClearTimeoutEvent")
	result["go.runtime.getRandomData"] = wasmtime.WrapFunc(store, rt.runtimeEnv.GetRandomData)
	//result["parigot.debugprint"] = wrapWithRecover(store, rt.syscall.DebugPrint)
	result["go.debug"] = wrapWithRecover(store, rt.runtimeEnv.GoDebug, "GoDebug")

	//system calls
	result["go.parigot.locate_"] = wrapWithRecover(store, rt.syscall.Locate, "Locate")
	//result["go.parigot.register_"] = wrapWithRecover(store, rt.syscall.Register,"Register")
	// xxx fix me: How are we going to clean up the resources for a particular service when it exits?
	// how do we find the resources associated with the caller of exit().
	result["go.parigot.exit_"] = wrapWithRecover(store, rt.syscall.Exit, "Exit")
	result["go.parigot.dispatch_"] = wrapWithRecover(store, rt.syscall.Dispatch, "Dispatch")
	result["go.parigot.bind_method_"] = wrapWithRecover(store, rt.syscall.BindMethod, "BindMethod")
	result["go.parigot.return_value_"] = wrapWithRecover(store, rt.syscall.ReturnValue, "ReturnValue")
	result["go.parigot.block_until_call_"] = wrapWithRecover(store, rt.syscall.BlockUntilCall, "BlockUntilCall")
	result["go.parigot.require_"] = wrapWithRecover(store, rt.syscall.Require, "Require")
	result["go.parigot.export_"] = wrapWithRecover(store, rt.syscall.Export, "Export")
	result["go.parigot.run_"] = wrapWithRecover(store, rt.syscall.Run, "Run")

	//backdoor to the logger
	result["go.parigot.backdoor_log_"] = wrapWithRecover(store, rt.syscall.BackdoorLog, "BackdoorLog")
}

func addSplitModeFunctions(store wasmtime.Storelike,
	result map[string]*wasmtime.Func,
	logViewer *loggo.LogViewerImpl,
	fileSvc *filego.FileSvcImpl,
	queueSvc *queuego.QueueSvcImpl) {

	// mixed mode entries: this should be automated (xxxfixmexxx)
	result["go.logviewer.log_request_handler"] = wrapWithRecover(store, logViewer.LogRequestHandler, "LogRequestHandler")
	// mixed mode entries: this should be automated (xxxfixmexxx)
	result["go.filesvc.open"] = wrapWithRecover(store, fileSvc.FileSvcOpen, "FileSvcOpen")
	result["go.filesvc.load"] = wrapWithRecover(store, fileSvc.FileSvcLoad, "FileSvcLoad")

	result["go.queuesvc.create_handler"] = wrapWithRecover(store, queueSvc.QueueSvcCreateQueue, "QueueSvcCreateHandler")
	result["go.queuesvc.delete_handler"] = wrapWithRecover(store, queueSvc.QueueSvcDeleteQueue, "QueueSvcDeleteHandler")
	result["go.queuesvc.length_handler"] = wrapWithRecover(store, queueSvc.QueueSvcLength, "QueueSvcLengthHandler")
	result["go.queuesvc.locate_handler"] = wrapWithRecover(store, queueSvc.QueueSvcLocate, "QueueSvcLocateHandler")
	result["go.queuesvc.mark_done_handler"] = wrapWithRecover(store, queueSvc.QueueSvcMarkDone, "QueueSvcMarkDoneHandler")
	result["go.queuesvc.send_handler"] = wrapWithRecover(store, queueSvc.QueueSvcSend, "QueueSvcSendHandler")
	result["go.queuesvc.receive_handler"] = wrapWithRecover(store, queueSvc.QueueSvcReceive, "QueueSvcReceiveHandler")
}

func wrapWithRecover(store wasmtime.Storelike, fn func(int32), name string) *wasmtime.Func {
	return wasmtime.WrapFunc(store, func(sp int32) {
		defer func() {
			if r := recover(); r != nil {
				print("START RECOVER STACKTRACE " + name + "\n")
				print(fmt.Sprintf("trapped panic: %s %v (%T)\n", name, r, r))
				debug.PrintStack()
				print("END RECOVER+STACKTRACE " + name + "\n")
			}
		}()
		if sp == 0 {
			panic("BAD SP")
		}
		fn(sp)
	})
}
