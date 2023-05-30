package shim

import (
	"context"

	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

	"github.com/tetratelabs/wazero/api"
)

// func fdReadFn(_ context.Context, mod api.Module, params []uint64) syscall.Errno {
//	return fdReadOrPread(mod, params, false)
//}

func FdReadFn() func(ctx context.Context, mod api.Module, params []uint64) {
	return wasi.ExposeFdReadFn
}
func FdWriteFn() func(ctx context.Context, mod api.Module, params []uint64) {
	return wasi.ExposeFdWriteFn
}
func PathOpen() func(ctx context.Context, mod api.Module, params []uint64) {
	return wasi.ExposePathOpen
}
func FdClose() func(ctx context.Context, mod api.Module, params []uint64) {
	return wasi.ExposeFdClose
}
