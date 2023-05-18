package main

import (
	"context"
	"log"
	"unsafe"

	"github.com/iansmith/parigot/eng"
	"github.com/iansmith/parigot/sys"

	"github.com/tetratelabs/wazero/api"
)

type filePlugin struct{}

var _ = unsafe.Sizeof([]byte{})

var ParigiotInitialize sys.ParigotInit = &filePlugin{}

func (*filePlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "file", "open", open)
	e.AddSupportedFunc(ctx, "file", "load_testt_data", loadTestData)
	return true
}

func open(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("file.open 0x%x", stack)
}

func loadTestData(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("file.dispatch 0x%x", stack)
}
