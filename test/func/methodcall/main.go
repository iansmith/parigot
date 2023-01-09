package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"
	"unsafe"

	"github.com/iansmith/parigot/api_impl/syscall"
	"github.com/iansmith/parigot/g/log/v1"
	"github.com/iansmith/parigot/g/methodcall/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	methodcallmsg "github.com/iansmith/parigot/g/msg/methodcall/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var logger *log.LogServiceClient
var foo *methodcall.FooServiceClient
var bar *methodcall.BarServiceClient

var callImpl = syscall.NewCallImpl()

var test = []testing.InternalTest{
	{
		Name: "TestFoo",
		F:    TestFoo,
	},
}

func loadCString(rawAddr int32) string {
	var buf bytes.Buffer
	i := int32(0)
	for {
		str := (*byte)(unsafe.Pointer(uintptr(int32(rawAddr + i))))
		buf.WriteByte(*str)
		if *str == 0 {
			break
		}
		i++
	}
	asBytes := buf.Bytes()
	return string(asBytes[:len(asBytes)-1])

}
func loadArgvPointer(argv int32) int32 {
	var buf [4]byte
	for i := int32(0); i < 4; i++ {
		str := (*byte)(unsafe.Pointer(uintptr(int32(argv + i))))
		buf[i] = *str
	}
	return int32(binary.LittleEndian.Uint32(buf[:]))
}

func computeArgvEnvp() ([]string, []string) {
	// xxx hacky: we put argv at a fixed memory location because we have no way to pass param to this function... except via argv!
	argvAddr := loadArgvPointer(0x1000)
	argv := []string{}
	envp := []string{}
	index := int32(0)
	for {
		ptr := loadArgvPointer(argvAddr + (8 * index))
		index++
		//print(fmt.Sprintf("xxx arg ptr is %x\n", ptr))
		if ptr == 0 {
			break
		}
		argv = append(argv, loadCString(ptr))
	}
	for {
		ptr := loadArgvPointer(argvAddr + (8 * index))
		index++
		//print(fmt.Sprintf("xxx env ptr is %x\n", ptr))
		if ptr == 0 {
			break
		}
		envp = append(envp, loadCString(ptr))
		index++
	}
	return argv, envp
}

func main() {
	argv, _ := computeArgvEnvp()
	//print(fmt.Sprintf("after building it ourselves: %v, %v\n", argv, envp))

	os.Args = argv
	flag.Parse()
	for i, arg := range flag.Args() {
		print(fmt.Sprintf("%d:%s\n", i, arg))
	}
	//print(fmt.Sprintf("parsed? %d and %s", flag.NArg(), flag.Arg(0)))
	// //ptr64 := uint64(uintptr(unsafe.Pointer(ptr)))
	// b := make([]byte, 8)
	// binary.LittleEndian.PutUint64(b, uint64(uintptr(unsafe.Pointer(sh.Data))))
	// // t := binary.LittleEndian.Uint32(b[4:])
	// // out := fmt.Sprintf("xxx os arg data (orig 0x%x) raw %x with Len=0x%x and Cap=0x%x\n",
	// // 	ptr, t, sh.Len, sh.Cap)
	// sh.Len = 4
	// sh.Cap = 4
	// print(">>>>", os.Args[0], "--", len(os.Args), ",", os.Args[1], "<<<<", "\n")
	// return
}

func xxxmain() {
	//flag.Parse() <--- can't do this until we get startup args figured out

	//testing.Init()

	if _, err := callImpl.Require1("methodcall", "FooService"); err != nil {
		panic("unable to require foo service: " + err.Error())
	}
	if _, err := callImpl.Require1("methodcall", "BarService"); err != nil {
		panic("unable to require bar service: " + err.Error())
	}
	if _, err := callImpl.Require1("log", "LogService"); err != nil {
		panic("unable to require log service: " + err.Error())
	}
	if _, err := callImpl.Run(&syscallmsg.RunRequest{Wait: true}); err != nil {
		panic("error starting client process:" + err.Error())
	}

	var err error
	logger, err = log.LocateLogService()
	if err != nil {
		panic(fmt.Sprintf("failed to locate LogService:%v", err))
	}
	foo, err = methodcall.LocateFooService(logger)
	if err != nil {
		panic(fmt.Sprintf("failed to locate FooService:%v", err))
	}
	bar, err = methodcall.LocateBarService(logger)
	if err != nil {
		panic(fmt.Sprintf("failed to locate BarService:%v", err))
	}

	// run tests
	ok := testing.RunTests(func(pat, str string) (bool, error) { return true, nil }, test)
	if ok {
		testMessage("result of RunTests %v", ok)
	}

	// cleanup?
	callImpl.Exit(&syscallmsg.ExitRequest{
		Code: 0,
	})
}

func TestFoo(t *testing.T) {
	fn := func(t *testing.T, value0, value1, sum, product int32) {
		req := &methodcallmsg.AddMultiplyRequest{
			Value0: value0,
			Value1: value1,
			IsAdd:  true,
		}
		resp, err := foo.AddMultiply(req)
		if err != nil {
			t.Fatalf("error in AddMultiply [add]:%v", err)
		}
		if resp.Result != sum {
			t.Fatalf("bad result for add, expected %d but got %d", sum, resp.Result)
		}

		req.IsAdd = false
		resp, err = foo.AddMultiply(req)
		if err != nil {
			t.Fatalf("error in AddMultiply [mult]:%v", err)
		}
		if resp.Result != product {
			t.Fatalf("bad result for multiply, expected %d but got %d", product, resp.Result)
		}
		testMessage("OK addMultiply(%d,%d)", value0, value1)
	}
	fn(t, 3, 7, 10, 21)
	fn(t, 0, 1, 1, 0)
	fn(t, -10, 20, 10, -200)
}

func testMessage(spec string, arg ...interface{}) {
	req := &logmsg.LogRequest{
		Stamp:   timestamppb.New(time.Now()), // xxx fix me, should be using the kernel
		Level:   logmsg.LogLevel_LOG_LEVEL_INFO,
		Message: fmt.Sprintf(spec, arg...),
	}
	if logger == nil {
		print("Internal error in methodcall test: logger is nil! "+fmt.Sprintf(spec, arg...), "\n")
		return
	}
	if err := logger.Log(req); err != nil {
		print("methodcall: error in log call:", err.Error(), "\n")
	}
}
