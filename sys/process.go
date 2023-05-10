package sys

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	logimpl "github.com/iansmith/parigot/apiimpl/log/go_"
	"github.com/iansmith/parigot/eng"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	syscallmsg "github.com/iansmith/parigot/g/msg/syscall/v1"
	"github.com/iansmith/parigot/sys/dep"
	"github.com/iansmith/parigot/sys/jspatch"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service interface {
	IsServer() bool
	IsLocal() bool
	IsRemote() bool
	GetName() string
	GetArg() []string
	GetEnv() []string
	GetPath() string
	GetModule() eng.Module
}

type ParigotExitCode int

const (
	ExitCodeArgsTooLarge  ParigotExitCode = 252
	ExitCodeNoStartSymbol ParigotExitCode = 251
	ExitCodePanic         ParigotExitCode = 254
	ExitCodeTrapped       ParigotExitCode = 253
	ExitCodeNoMain        ParigotExitCode = 255
)

const useWasmtime = false

// Flip this switch to see debug messages from the process.
var processVerbose = true || envVerbose != ""

var lastProcessId = 7

type Process struct {
	id   int
	path string

	lock sync.Mutex

	module eng.Module
	//linkage      []eng.Extern
	memPtr       uintptr
	instance     eng.Instance
	engine       eng.Engine
	syscall      *syscallReadWrite
	microservice Service
	key          dep.DepKey

	requirementsMet bool
	reachedRunBlock bool
	running         bool
	exited          bool
	exitCode_       int //really only 0-192

	argv       int32 //ptr
	argc       int32
	argvBuffer *bytes.Buffer

	callCh chan *callContext
	runCh  chan bool
}

// NewProcessFromMicroservice does not handle concurrent use. It assumes that each call to this
// method is called from the same thread/goroutine, in sequence.  This is, effectively,
// a loader for the os.  xxxfixme this really should be safe to use in multiple go routines ... then we
// could have a repl??
func NewProcessFromMicroservice(engine eng.Engine, m Service, ctx *DeployContext) (*Process, error) {

	lastProcessId++
	id := lastProcessId
	proc := &Process{
		id:              id,
		engine:          engine,
		module:          m.GetModule(),
		memPtr:          0,
		instance:        nil,
		syscall:         ctx.supportFunc.rt.syscall,
		running:         false,
		reachedRunBlock: false,
		exited:          false,
		microservice:    m,
		path:            m.GetPath(),

		callCh: make(chan *callContext),
		//runCh:  make(chan bool),
	}
	proc.key = NewDepKeyFromProcess(proc)

	// l, memPtr, err := proc.checkLinkage(rt, logViewer, fileSvc, queueSvc)
	// if err != nil {
	// 	return nil, err
	// }
	// if proc.memPtr == 0 && memPtr != 0 {
	// 	proc.memPtr = memPtr
	// }
	// XXX WASMTIME SPECIFIC
	//	ext := make([]eng.Extern, len(l))
	// for i, inst := range l {
	// 	ext[i] = eng.NewWasmtimeLinkage(inst)
	// }
	//proc.linkage = ext
	instance, err := proc.module.NewInstance()
	if err != nil {
		return nil, err
	}
	proc.instance = instance

	memExt, err := instance.GetMemoryExport()
	if err != nil {
		return nil, err
	}
	// memory pointer shenanigans
	memptr := memExt.Memptr()
	ctx.supportFunc.SetMemPtr(memptr)
	ctx.supportFunc.SetProcess(proc)
	proc.SetMemPtr(memptr)

	return proc, nil
}

func (p *Process) RequirementsMet() bool {
	return p.requirementsMet
}
func (p *Process) IsServer() bool {
	// if we have a remote spec, then we are remote
	return p.microservice.IsServer()
}

func (p *Process) Exit() {
	p.lock.Lock()
	defer p.lock.Unlock()
	panic("xxx exit() xxx")
	log.Printf("process %s exiting\n", p)
	p.exited = true
}

func (p *Process) String() string {
	p.lock.Lock()
	defer p.lock.Unlock()

	dir, file := filepath.Split(p.path)
	if dir == "" {
		dir = "."
	}

	return fmt.Sprintf("[proc-%d:%s:%s]", p.id, p.microservice.GetName(), file)
}
func (p *Process) SetMemPtr(u uintptr) {
	p.memPtr = u
}
func (p *Process) SetReachedRunBlock(r bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.reachedRunBlock = r
}

func (p *Process) ReachedRunBlock() bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.reachedRunBlock
}
func (p *Process) Running() bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.running
}

func (p *Process) SetExited(e bool) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.exited = e
}
func (p *Process) Exited() bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.exited
}

// func (p *Process) checkLinkage(rt *Runtime, lv *logimpl.LogViewerImpl, fs *fileimpl.FileSvcImpl, queue *queueimpl.QueueSvcImpl) ([]wasmtime.AsExtern, uintptr, error) {

// 	// all available funcs end up in here
// 	//available := make(map[string]*wasmtime.Func)
// 	//availableObj := make(map[string]wasmtime.AsExtern)

// 	//addEmscriptenFuncs(p.parent, available, rt)
// 	addSupportedFunctions(p.engine, rt)
// 	addSplitModeFunctions(p.engine, lv, fs, queue)

// 	result := uintptr(0)
// 	//memPtr := addEmscriptenObjects(p.parent, availableObj)

// 	//	_glob := make(map[string]*wasmtime.Global)
// 	//addEmscriptenGlobals(p.parent, glob)

// 	// result of checking the linkage
// 	linkage := []wasmtime.AsExtern{}

// walk all the module's imports
// 	for _, imp := range p.module.Imports() {
// 		n := "$$ANON$$"
// 		if imp.Name() != nil {
// 			n = *imp.Name()
// 		}
// 		importName := fmt.Sprintf("%s.%s", imp.Module(), n)
// 		ext, ok := available[importName]
// 		if !ok {
// 			// // possibly a global
// 			// g, ok := glob[importName]
// 			// if ok {
// 			// 	external := wasmtime.AsExtern(g)
// 			// 	linkage = append(linkage, external)
// 			// 	continue
// 			// } else {
// 			// 	// might be another type of object
// 			// 	obj, ok := availableObj[importName]
// 			// 	log.Printf("request for import %s, %+v  ok?%v  %v", importName, availableObj, ok, obj)
// 			// 	if ok {
// 			// 		log.Printf("extra link step for %s", importName)
// 			// 		//result = memPtr
// 			// 		linkage = append(linkage, obj)
// 			// 		continue
// 			// 	}
// 			// }
// 			return nil, 0, fmt.Errorf("unable to find linkage for %s in module %s", importName, p.path)
// 		} else {
// 			linkage = append(linkage, ext)
// 		}
// 	}
// 	return linkage, result, nil
//}

// func addEmscriptenObjects(store *wasmtime.Store, obj map[string]wasmtime.AsExtern) uintptr {
// 	mtype := wasmtime.NewMemoryType(320, true, 32767)
// 	mem, err := wasmtime.NewMemory(store, mtype)
// 	if err != nil {
// 		panic("unable to create memory object inside wasmtime")
// 	}
// 	memPtr := uintptr(mem.Data(store))
// 	obj["env.memory"] = wasmtime.AsExtern(mem)
// 	ftype := wasmtime.NewValType(wasmtime.KindFuncref)
// 	ttype := wasmtime.NewTableType(ftype, 6376, true, 8192)
// 	tbl, err := wasmtime.NewTable(store, ttype, wasmtime.ValFuncref(nil))
// 	if err != nil {
// 		panic("unable to create table object inside wasmtime")
// 	}
// 	obj["env.__indirect_function_table"] = tbl
// 	return memPtr
// }

// func addSupportedGlobals(store *wasmtime.Store, glob map[string]*wasmtime.Global) {
// 	valType := wasmtime.NewValType(wasmtime.KindI32)
// 	gType := wasmtime.NewGlobalType(valType, false)
// 	g, err := wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
// 	if err != nil {
// 		log.Fatalf("unable to create global " + err.Error())
// 	}
// 	glob["env.__memory_base"] = g
// }

func (p *Process) SetExitCode(code int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.exitCode_ = code
	p.exited = true
}

func (p *Process) ExitCode() int {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.exitCode_
}
func (p *Process) IsRunning() bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.running
}

func (p *Process) SetRunning(r bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.running = r
}

// Run() is used to let a process proceed with running.  This is
// called when we discover all his requirements have been met.
func (p *Process) Run() {
	p.runCh <- true
}

// Start invokes the wasm interp and returns an error code if this is a "main" process.
func (p *Process) Start() (code int) {
	procPrint("START ", "start process: %s", p)
	var err error
	procPrint("START ", "start of args  %+v", p.microservice.GetArg())
	// startOfArgs := wasmStartAddr + int32(0)
	// if p == nil {
	// 	panic("process is nil!")
	// }
	// p.lock.Lock()
	// defer p.lock.Unlock()

	// procPrint("START ", "get buffer from args and env  %s", p)
	// p.argvBuffer, p.argv, err = GetBufferFromArgsAndEnv(p.microservice, startOfArgs)
	// if err != nil {
	// 	code = int(ExitCodeArgsTooLarge)
	// 	return
	// }
	// p.argc = int32(len(p.microservice.GetArg()))

	var wasmMem *jspatch.WasmMem
	if useWasmtime {
		procPrint("START", "create wasm mem %s (%v)", p, p.memPtr != 0)
		//wasmMem := &jspatch.WasmMem{}
		wasmMem = jspatch.NewWasmMem(p.memPtr)

		// wasmMem.SetInt32(wasmStartAddr-int32(4), p.argv)
		// wasmMem.CopyToMemAddr(startOfArgs, p.argvBuffer.Bytes())
	}

	procPrint("START", "get entry point")
	start, err := p.instance.GetEntryPointExport()
	if err != nil {
		panic(err)
	}
	procPrint("START", "defer %s (%v)", p, start != nil)

	defer func(proc *Process) {
		r := recover()
		if r != nil {
			log.Printf("defer caught it %T, %v", r, r.(string))
			log.Printf("flush")

			procPrint("START ******** ", "INSIDE defer %s, %+v", proc, r)
			e, ok := r.(*syscallmsg.ExitRequest)
			procPrint("Start/Exit ", "INSIDE defer exit req %+v, ok %v", r.(*syscallmsg.ExitRequest), ok)
			if ok {
				code = int(e.GetCode())
				proc.SetExitCode(code)
				procPrint("Start/Exit", "INSIDE DEFER exiting with code %d", e.GetCode())
			} else {
				p.SetExitCode(int(ExitCodePanic))
				code = int(ExitCodePanic)
				print(fmt.Sprintf("golang (not WASM) panic '%v'\n", r))
			}
		}
	}(p)
	procPrint("START ", "calling start func %s", p)
	var info interface{}
	if useWasmtime {
		info = eng.NewWasmtimeMemInfo(wasmMinDataAddr, wasmMem)
	}

	retVal, err := start.Run(p.microservice.GetArg(), info)
	procPrint("END ", "process %s has completed: result=%v, err=%v", p, retVal, err)

	if err != nil {
		p.SetExitCode(int(ExitCodeTrapped))
		procPrint("END ", "process %s trapped: %v, exit code %d", p, err, p.ExitCode())
		return int(ExitCodeTrapped)
	}
	if retVal == nil {
		procPrint("END ", "process %s finished w/no return value (exit code %d)", p, p.ExitCode())
		p.SetExited(true)
		return p.ExitCode()
	}
	procPrint("END ", "process %s finished normally: %+v", p, retVal)
	procPrint("END ", "going to sleep now")
	ch := make(chan struct{})
	<-ch
	return p.ExitCode()
}

func procPrint(method string, spec string, arg ...interface{}) {
	if processVerbose {
		part1 := fmt.Sprintf("PROCESS:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		logimpl.ProcessLogRequest(
			&logmsg.LogRequest{
				Level:   logmsg.LogLevel_LOG_LEVEL_INFO,
				Message: part1 + part2 + "\n",
				Stamp:   timestamppb.Now(), // xxx should use the kernel calls
			}, true, false, false, nil)
	}
}
