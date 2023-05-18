package main

import (
	"context"
	"log"

	"reflect"
	"unsafe"

	"github.com/iansmith/parigot/eng"
	//"github.com/iansmith/parigot/g/msg/syscall/v1"

	"github.com/tetratelabs/wazero/api"
)

type syscallPlugin struct{}

var ParigiotInitialize = syscallPlugin{}

func (*syscallPlugin) Init(ctx context.Context, e eng.Engine) bool {
	e.AddSupportedFunc(ctx, "parigot", "locate_", locate)
	e.AddSupportedFunc(ctx, "parigot", "dispatch_", dispatch)
	e.AddSupportedFunc(ctx, "parigot", "block_until_call_", blockUntilCall)
	e.AddSupportedFunc(ctx, "parigot", "bind_method_", bindMethod)
	e.AddSupportedFunc(ctx, "parigot", "run_", run)
	e.AddSupportedFunc(ctx, "parigot", "export", export)
	e.AddSupportedFunc(ctx, "parigot", "return_value_", returnValue)
	e.AddSupportedFunc(ctx, "parigot", "require", require)
	e.AddSupportedFunc(ctx, "parigot", "exit", exit)
	return true
}

func locate(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("locate %s 0x%x", m.Name(), stack)
}

func dispatch(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("dispatch 0x%x", stack)
}

func blockUntilCall(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("blockUntilCall 0x%x", stack)
}
func bindMethod(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("bindMethod 0x%x", stack)
}
func run(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("run 0x%x", stack)
}
func export(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("export 0x%ax", stack)
}
func returnValue(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("returnValue 0x%x", stack)
}

func require(ctx context.Context, m api.Module, stack []uint64) {
	malloc := m.ExportedFunction("malloc")
	//free := m.ExportedFunction("free")
	resultArr, err := malloc.Call(ctx, []uint64{12}...) //size in bytes
	if err != nil {
		log.Fatalf("failed to call malloc: %v", err)
	}
	print("got a value back from malloc ", resultArr[0], "\n")
	storeArea := resultArr[0]

	print("xxx require ", len(stack), " and ", len(resultArr), "\n")
	flatten := m.ExportedFunction("github.com/iansmith/parigot/apiwasm/syscall.RequireFlattenToByte")
	print("xxx require2 ", flatten == nil, "\n")
	resp, err := flatten.Call(ctx, storeArea)
	print("xxx require3 ", len(resp), "\n\n")
	if err != nil {
		log.Printf("xxx error from flatten %v", err)
	}
	print("xxx resp size is ", len(resp), "\n")
	ptr := api.DecodeU32(resp[0])
	b := (*reflect.SliceHeader)(unsafe.Pointer(uintptr(ptr)))
	print("xxx size of request: %d,%d,%p\n", b.Len, b.Cap, b.Data)
	//ptr := uintptr(stack[0])
	//fake := make([]byte, 0)
	// xxx this is terrible, using server side size for wasm side size
	// req := &syscallmsg.RequireRequest{}
	// size := req.ProtoReflect().ProtoMessage().Size()
	// log.Printf("type %T, size %d", size)

	// b, ok := m.Memory().Read(api.DecodeU32(stack[0]), size)
	// if !ok {
	// 	stack[0] = EncodeI32(0)
	// 	return
	// }
	// log.Printf("found b as %d", len(b))
	//req := (*syscallmsg.RequireRequest)(unsafe.Pointer(ptr))
	//	serverSide := copy
	//log.Printf("xxx -- %+v", req)
}
func exit(ctx context.Context, m api.Module, stack []uint64) {
	log.Printf("exit 0x%x", stack)
	panic("exit called ")
}
