package eng

import (
	"context"
	"syscall"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

var bogus = []byte("xxxxxxx")

var fakeArg = [][]byte{
	bogus,
}

var userEntry = []byte("USER=parigot")
var fakeEnv = [][]byte{
	userEntry,
}

var argc = len(fakeArg)
var argLen = uint32(len(fakeArg[0]))

var envc = len(fakeEnv)
var envLen = uint32(len(fakeEnv[0]))

func fakeWasiArgsSizesGetFn(ctx context.Context, mod api.Module, stack []uint64) {
	mem := mod.Memory()
	resultArgc, resultArgvLen := uint32(stack[0]), uint32(stack[1])

	// argc and argv_len offsets are not necessarily sequential, so we have to
	// write them independently.
	if !mem.WriteUint32Le(resultArgc, uint32(argc)) {
		stack[0] = uint64(syscall.EFAULT)
		return
	}
	if !mem.WriteUint32Le(resultArgvLen, argLen+1) {
		stack[0] = uint64(syscall.EFAULT)
		return
	}
	stack[0] = uint64(0)
}

func fakeWasiEnvSizesGetFn(ctx context.Context, mod api.Module, stack []uint64) {
	mem := mod.Memory()
	resultEnvc, resultEnvpLen := uint32(stack[0]), uint32(stack[1])

	// argc and argv_len offsets are not necessarily sequential, so we have to
	// write them independently.
	if !mem.WriteUint32Le(resultEnvc, uint32(envc)) {
		stack[0] = uint64(syscall.EFAULT)
		return
	}
	if !mem.WriteUint32Le(resultEnvpLen, envLen+1) {
		stack[0] = uint64(syscall.EFAULT)
		return
	}
	stack[0] = uint64(0)
}

func fakeWasiArgsGetFn(_ context.Context, mod api.Module, params []uint64) {
	argv, argvBuf := uint32(params[0]), uint32(params[1])
	params[0] = uint64(writeOffsetsAndNullTerminatedValues(mod.Memory(), fakeArg, argv, argvBuf, argLen+1))
}
func fakeWasiEnvGetFn(_ context.Context, mod api.Module, params []uint64) {
	envp, envpBuf := uint32(params[0]), uint32(params[1])
	params[0] = uint64(writeOffsetsAndNullTerminatedValues(mod.Memory(), fakeEnv, envp, envpBuf, argLen+1))
}

// writeOffsetsAndNullTerminatedValues is used to write NUL-terminated values
// for args or environ, given a pre-defined bytesLen (which includes NUL
// terminators).
func writeOffsetsAndNullTerminatedValues(mem api.Memory, values [][]byte, offsets, bytes, bytesLen uint32) syscall.Errno {
	// The caller may not place bytes directly after offsets, so we have to
	// read them independently.
	valuesLen := len(values)
	offsetsLen := uint32(valuesLen * 4) // uint32Le
	offsetsBuf, ok := mem.Read(offsets, offsetsLen)
	if !ok {
		return syscall.EFAULT
	}
	bytesBuf, ok := mem.Read(bytes, bytesLen)
	if !ok {
		return syscall.EFAULT
	}

	// Loop through the values, first writing the location of its data to
	// offsetsBuf[oI], then its NUL-terminated data at bytesBuf[bI]
	var oI, bI uint32
	for _, value := range values {
		// Go can't guarantee inlining as there's not //go:inline directive.
		// This inlines uint32 little-endian encoding instead.
		bytesOffset := bytes + bI
		offsetsBuf[oI] = byte(bytesOffset)
		offsetsBuf[oI+1] = byte(bytesOffset >> 8)
		offsetsBuf[oI+2] = byte(bytesOffset >> 16)
		offsetsBuf[oI+3] = byte(bytesOffset >> 24)
		oI += 4 // size of uint32 we just wrote

		// Write the next value to memory with a NUL terminator
		copy(bytesBuf[bI:], value)
		bI += uint32(len(value))
		bytesBuf[bI] = 0 // NUL terminator
		bI++
	}

	return 0
}

func fakeWasiAddFunc(e *wazeroEng) wazero.HostModuleBuilder {
	//wasi_snapshot_preview1.MustInstantiate(ctx, e.r)
	// Export the default WASI functions.
	wasiBuilder := e.rt.NewHostModuleBuilder("wasi_snapshot_preview1")
	wasi_snapshot_preview1.NewFunctionExporter().ExportFunctions(wasiBuilder)

	// Subsequent calls to NewFunctionBuilder override built-in exports.
	wasiBuilder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(fakeWasiArgsGetFn), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).
		Export("args_get")
	wasiBuilder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(fakeWasiEnvGetFn), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).
		Export("env_get")

	wasiBuilder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(fakeWasiArgsSizesGetFn), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).
		Export("args_sizes_get")

	wasiBuilder.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(fakeWasiEnvSizesGetFn), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).
		Export("env_sizes_get")

	return wasiBuilder
}
