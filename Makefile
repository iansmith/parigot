GO_CMD=go #really shouldn't need to change this if you use the tools directory
WASM_GRAMMAR=command/Wasm.g4
TRANSFORM=command/transform
FLAVOR=atlanta.base

all: build/runner \
build/protoc-gen-parigot \
build/surgery \
build/jdepp \
build/nameserver

build/jdepp: command/jdepp/main.go
	@echo
	@echo "\033[92mjdepp ============================================================================================\033[0m"
	go build -o build/jdepp github.com/iansmith/parigot/command/jdepp

build/surgery: command/transform/wasm_parser.go
	@echo
	@echo "\033[92msurgery ============================================================================================\033[0m"
	go build -o build/surgery github.com/iansmith/parigot/command/surgery

build/protoc-gen-parigot: \
	command/protoc-gen-parigot/template/go/*.tmpl
	@echo
	@echo "\033[92mprotoc-gen-parigot =================================================================================\033[0m"
	go build -o build/protoc-gen-parigot github.com/iansmith/parigot/command/protoc-gen-parigot

command/transform/wasm_parser.go: $(WASM_GRAMMAR)
	@echo
	@echo "\033[92mWASM wat file parser \(via Antlr4 and Wasm.g4\) ====================================================\033[0m"
	cd command; java -Xmx500M -cp "$$PARIGOT_TOOLS/lib/antlr-4.9.3-complete.jar" org.antlr.v4.Tool -Dlanguage=Go -o transform -package transform Wasm.g4; cd ..

build/runner: g/log/logservicedecl.p.go \
	api/atlanta.base/proto/pb/ns/*.proto \
	sys/atlanta.base/*.go 
	@echo
	@echo "\033[92mrunner =============================================================================================\033[0m"
	go build -o build/runner github.com/iansmith/parigot/command/runner

g/log/logservicedecl.p%go g/net/netservicedecl.p%go: build/protoc-gen-parigot
	@echo
	@echo "\033[92mbuilding parigot interfaces ========================================================================\033[0m"
	buf generate

build/nameserver: command/nameserver/*.go sys/atlanta.base/*.go 
	@echo
	@echo "\033[92mnameserver =============================================================================================\033[0m"
	go build -o build/nameserver github.com/iansmith/parigot/command/nameserver

image:
	@echo
	@echo "\033[92mdocker images =======================================================================================\033[0m"
	cp build/nameserver container/runner/runner 
	rm -f container/runner/data.wasm
	touch container/runner/data.wasm
	# runs in separate shell, so need to cd with && but don't need to "come back"
	cd container/runner && docker build  -t nameserver . 
	cp build/runner container/runner/runner
	cp example/vvv/build/server.p.wasm container/runner/data.wasm
	cd container/runner && docker build  -t storeserver . 
	cp example/vvv/build/storeclient.p.wasm container/runner/data.wasm
	cd container/runner && docker build  -t storeclient . 

clean:
	@echo "\033[92mclean ==============================================================================================\033[0m"
	rm -f build/*
	rm -rf g/kernel g/log g/net g/*.p.go g/pb
	rm -f $(TRANSFORM)/Wasm.* $(TRANSFORM)/WasmLexer.* $(TRANSFORM)/wasm_base_listener.go $(TRANSFORM)/wasm_lexer.go $(TRANSFORM)/wasm_parser.go $(TRANSFORM)/wasm_listener.go


.PHONY:test
test: $(PGP)
	go run command/protoc-gen-parigot/test/main.go build/protoc-gen-parigot command/protoc-gen-parigot/test/testdata/t0 - abi/abihelper.go

#### Do not remove this line or edit below it.  The rest of this file is computed by jdepp.
### jdepp computed dependencies for binary: build/testdata
build/testdata: \
	command/surgery/testdata/ex1/ex1.go \
	command/surgery/testdata/main.go

### jdepp computed dependencies for binary: build/terminallog
build/terminallog: \
	sys/atlanta.base/stdlib/terminallog/main.go \
	lib/atlanta.base/go/call.go \
	lib/atlanta.base/go/callimpl.go \
	lib/atlanta.base/go/msgsize.go \
	lib/atlanta.base/go/pctx.go \
	lib/atlanta.base/go/id.go \
	lib/atlanta.base/go/perror.go \
	lib/atlanta.base/go/idconv.go \
	lib/atlanta.base/go/syscallpayload.go \
	sys/atlanta.base/stdlib/terminallog/gen.go \
	lib/atlanta.base/go/calljs.go \
	lib/atlanta.base/go/callnonjs.go \
	lib/atlanta.base/go/client.go \
	lib/atlanta.base/go/helpers.go \
	lib/atlanta.base/go/log.go \
	build/protoc-gen-parigot

### jdepp computed dependencies for binary: build/jdepp
build/jdepp: \
	command/jdepp/main.go

### jdepp computed dependencies for binary: build/nameserver
build/nameserver: \
	lib/atlanta.base/go/log.go \
	sys/atlanta.base/memutil.go \
	lib/atlanta.base/go/call.go \
	lib/atlanta.base/go/helpers.go \
	sys/atlanta.base/remote.go \
	sys/atlanta.base/nameserver_test.go \
	sys/atlanta.base/nameservercore.go \
	lib/atlanta.base/go/id.go \
	lib/atlanta.base/go/msgsize.go \
	lib/atlanta.base/go/idconv.go \
	lib/atlanta.base/go/syscallpayload.go \
	sys/atlanta.base/local.go \
	sys/atlanta.base/process.go \
	command/nameserver/main.go \
	lib/atlanta.base/go/calljs.go \
	lib/atlanta.base/go/perror.go \
	sys/atlanta.base/runtime.go \
	lib/atlanta.base/go/client.go \
	lib/atlanta.base/go/pctx.go \
	sys/atlanta.base/netnameserver.go \
	sys/atlanta.base/nameserver.go \
	sys/atlanta.base/func.go \
	sys/atlanta.base/syscallrw.go \
	build/protoc-gen-parigot \
	lib/atlanta.base/go/callimpl.go \
	lib/atlanta.base/go/callnonjs.go \
	sys/atlanta.base/nsio.go \
	sys/atlanta.base/syscall.go

### jdepp computed dependencies for binary: build/protoc-gen-parigot
build/protoc-gen-parigot: \
	command/protoc-gen-parigot/codegen/wasmMethod.go \
	command/protoc-gen-parigot/go_/funcchoice.go \
	command/protoc-gen-parigot/codegen/cgtype.go \
	command/protoc-gen-parigot/codegen/generate.go \
	command/protoc-gen-parigot/codegen/lang.go \
	command/protoc-gen-parigot/codegen/options.go \
	command/protoc-gen-parigot/codegen/wasmservice.go \
	command/protoc-gen-parigot/codegen/finder.go \
	command/protoc-gen-parigot/codegen/helper.go \
	command/protoc-gen-parigot/codegen/wasm.go \
	command/protoc-gen-parigot/codegen/wasmfield.go \
	command/protoc-gen-parigot/go_/gotext.go \
	command/protoc-gen-parigot/main.go \
	command/protoc-gen-parigot/test/main.go \
	command/protoc-gen-parigot/codegen/funcchooser.go \
	command/protoc-gen-parigot/codegen/geninfo.go \
	command/protoc-gen-parigot/codegen/text.go \
	command/protoc-gen-parigot/util/out.go \
	command/protoc-gen-parigot/codegen/pass.go \
	command/protoc-gen-parigot/codegen/wasmmessage.go \
	command/protoc-gen-parigot/util/plugin.go \
	command/protoc-gen-parigot/go_/gogen.go

### jdepp computed dependencies for binary: build/test
build/test: \
	command/protoc-gen-parigot/test/main.go

### jdepp computed dependencies for binary: build/runner
build/runner: \
	sys/atlanta.base/memutil.go \
	sys/atlanta.base/netnameserver.go \
	sys/atlanta.base/syscallrw.go \
	sys/atlanta.base/func.go \
	sys/atlanta.base/nameserver_test.go \
	sys/atlanta.base/nameservercore.go \
	sys/atlanta.base/nsio.go \
	sys/atlanta.base/remote.go \
	sys/atlanta.base/runtime.go \
	sys/atlanta.base/nameserver.go \
	command/runner/fileload.go \
	command/runner/main.go \
	sys/atlanta.base/local.go \
	sys/atlanta.base/process.go \
	sys/atlanta.base/syscall.go

### jdepp computed dependencies for binary: build/surgery
build/surgery: \
	command/transform/sem_op.go \
	command/transform/typedescriptor.go \
	command/transform/wasm_parser.go \
	command/surgery/testdata/ex1/ex1.go \
	command/surgery/changetype.go \
	command/surgery/convert.go \
	command/transform/build_misc.go \
	command/transform/build_stmt.go \
	command/transform/wasm_lexer.go \
	command/surgery/testdata/main.go \
	command/transform/build_terminal.go \
	command/transform/sem_misc.go \
	command/transform/sem_toplevel.go \
	command/transform/wasm_base_listener.go \
	command/transform/wasm_listener.go \
	command/surgery/dbgprint.go \
	command/surgery/main.go \
	command/transform/build_module.go \
	command/surgery/unlink.go \
	command/transform/sem_func.go \
	command/transform/sem_module.go \
	command/transform/sem_stmt.go \
	command/surgery/parse.go \
	command/surgery/replacefn.go \
	command/surgery/tree.go

