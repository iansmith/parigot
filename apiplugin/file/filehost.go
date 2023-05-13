package main

import (
	"context"
	"log"

	"github.com/iansmith/parigot/eng"
	"github.com/iansmith/parigot/sys"
)

type filePlugin struct{}

var ParigiotInitialize sys.ParigotInit = &filePlugin{}

func (*filePlugin) Init(ctx context.Context, e eng.Engine, inst eng.Instance) bool {
	e.AddSupportedFunc(ctx, "file", "open", wrapFunc(inst, open))
	e.AddSupportedFunc(ctx, "file", "load_testt_data", wrapFunc(inst, loadTestData))
	return true
}

func wrapFunc(i eng.Instance, fn func(eng.Instance, int32) int32) func(int32) int32 {
	return func(x int32) int32 {
		return fn(i, x)
	}
}

func open(inst eng.Instance, ptr int32) int32 {
	log.Printf("file.open 0x%x", ptr)
	return 0
}

func loadTestData(inst eng.Instance, ptr int32) int32 {
	log.Printf("file.dispatch 0x%x", ptr)
	return 0
}
