BUILD_PRINT = \e[1;34mBuilding $<\e[0m
GO_CMD=go #really shouldn't need to change this if you use the tools directory
WASM_GRAMMAR=command/Wasm.g4
TRANSFORM=command/transform
FLAVOR=atlanta.base

all: build/runner \
build/protoc-gen-parigot \
build/surgery \
build/jdepp

build/jdepp: command/jdepp/main.go
	@echo
	@echo "\033[92mjdepp ============================================================================================\033[0m"
	go build -o build/jdepp github.com/iansmith/parigot/command/jdepp

build/surgery: command/transform/wasm_parser.go
	@echo
	@echo "\033[92msurgery ============================================================================================\033[0m"
	go build -o build/surgery github.com/iansmith/parigot/command/surgery

build/protoc-gen-parigot: \
	command/protoc-gen-parigot/template/go/*.tmpl \
	command/protoc-gen-parigot/template/abi/*.tmpl
	@echo
	@echo "\033[92mprotoc-gen-parigot =================================================================================\033[0m"
	go build -o build/protoc-gen-parigot github.com/iansmith/parigot/command/protoc-gen-parigot

command/transform/wasm_parser.go: $(WASM_GRAMMAR)
	@echo
	@echo "\033[92mWASM wat file parser \(via Antlr4 and Wasm.g4\) ====================================================\033[0m"
	cd command; java -Xmx500M -cp "../../tools/lib/antlr-4.9.3-complete.jar" org.antlr.v4.Tool -Dlanguage=Go -o transform -package transform Wasm.g4; cd ..

command/runner/g/abihelper.p.go: g/parigot/abi/abi.p.go
	@echo
	@echo "\033[92mabi helper =========================================================================================\033[0m"
	mv g/parigot/abi/abihelper.p.go command/runner/g/abihelper.p.go
	touch command/runner/g/abihelper.p.go

build/runner: g/parigot/abi/abi.p.go \
	g/parigot/log/logservicedecl.p.go \
	g/parigot/net/netservicedecl.p.go \
	command/runner/g/abihelper.p.go \
	../tools/tinygo0.26/targets/wasm-undefined.txt
	@echo
	@echo "\033[92mrunner =============================================================================================\033[0m"
	go build -o build/runner github.com/iansmith/parigot/command/runner

g/parigot/abi/abi.p%go g/parigot/log/logservicedecl.p%go g/parigot/net/netservicedecl.p%go g/parigot/abi/abiwasm-undefined%txt: build/protoc-gen-parigot
	@echo
	@echo "\033[92mbuilding parigot interfaces ========================================================================\033[0m"
	buf generate

../tools/tinygo0.26/targets/wasm-undefined.txt: g/parigot/abi/abiwasm-undefined.txt
	@echo
	@echo "\033[92mupdating undefined symbols =========================================================================\033[0m"
	cp ../tools/tinygo0.26/targets/wasm-undefined.txt.orig ../tools/tinygo0.26/targets/wasm-undefined.txt
	cat g/parigot/abi/abiwasm-undefined.txt >> ../tools/tinygo0.26/targets/wasm-undefined.txt
	touch ../tools/tinygo0.26/targets/wasm-undefined.txt

clean:
	@echo "\033[92mclean ==============================================================================================\033[0m"
	rm -f build/*
	rm -rf g/parigot/*
	rm -rf command/runner/g/*
	rm -f $(TRANSFORM)/Wasm.* $(TRANSFORM)/WasmLexer.* $(TRANSFORM)/wasm_base_listener.go $(TRANSFORM)/wasm_lexer.go $(TRANSFORM)/wasm_parser.go $(TRANSFORM)/wasm_listener.go


## shorthands
#net: $(REP_API_NET)
#abi: $(REP_ABI)
#protoc: $(PGP)
#gen: $(PGP) abi net
#	buf generate
#runner:build/runner
#surgery:build/surgery

.PHONY:test
test: $(PGP)
	go run command/protoc-gen-parigot/test/main.go build/protoc-gen-parigot command/protoc-gen-parigot/test/testdata/t0 - abi/abihelper.go

#### Do not remove this line or edit below it.  The rest of this file is computed by jdepp.
### jdepp computed dependencies for binary: build/test
build/test: \
	command/protoc-gen-parigot/test/main.go

### jdepp computed dependencies for binary: build/runner
build/runner: \
	command/runner/main.go \
	sys/abiimpl/abiimpl.go \
	abi/atlanta.base/go/jspatch/jshelper.go \
	abi/atlanta.base/go/tinygopatch/tinygohelper.go \
	command/runner/g/abihelper.p.go

### jdepp computed dependencies for binary: build/surgery
build/surgery: \
	command/surgery/dbgprint.go \
	command/transform/build_module.go \
	command/transform/typedescriptor.go \
	command/transform/wasm_listener.go \
	command/surgery/convert.go \
	command/surgery/parse.go \
	command/surgery/unlink.go \
	command/transform/build_stmt.go \
	command/transform/sem_op.go \
	command/transform/sem_toplevel.go \
	command/transform/wasm_base_listener.go \
	command/surgery/testdata/ex1/ex1.go \
	command/surgery/tree.go \
	command/transform/build_misc.go \
	command/transform/build_terminal.go \
	command/transform/sem_func.go \
	command/transform/sem_misc.go \
	command/transform/sem_module.go \
	command/surgery/main.go \
	command/surgery/replacefn.go \
	command/transform/sem_stmt.go \
	command/transform/wasm_lexer.go \
	command/transform/wasm_parser.go \
	lib/base/go/log/log.go \
	command/surgery/testdata/main.go

### jdepp computed dependencies for binary: build/testdata
build/testdata: \
	command/surgery/testdata/ex1/ex1.go \
	lib/base/go/log/log.go \
	command/surgery/testdata/main.go

### jdepp computed dependencies for binary: build/hello-go
build/hello-go: \
	lib/base/go/log/log.go \
	abi/atlanta.base/proto/abi/abi.proto \
	build/protoc-gen-parigot \
	example/hello-go/main.go

### jdepp computed dependencies for binary: build/jdepp
build/jdepp: \
	command/jdepp/main.go

### jdepp computed dependencies for binary: build/protoc-gen-parigot
build/protoc-gen-parigot: \
	command/protoc-gen-parigot/codegen/generate.go \
	command/protoc-gen-parigot/codegen/text.go \
	command/protoc-gen-parigot/codegen/wasmMethod.go \
	command/protoc-gen-parigot/codegen/wasmfield.go \
	command/protoc-gen-parigot/codegen/wasmmessage.go \
	command/protoc-gen-parigot/codegen/finder.go \
	command/protoc-gen-parigot/codegen/geninfo.go \
	command/protoc-gen-parigot/codegen/helper.go \
	command/protoc-gen-parigot/codegen/wasmservice.go \
	command/protoc-gen-parigot/main.go \
	command/protoc-gen-parigot/test/main.go \
	command/protoc-gen-parigot/abi/abigen.go \
	command/protoc-gen-parigot/codegen/cgtype.go \
	command/protoc-gen-parigot/codegen/funcchooser.go \
	command/protoc-gen-parigot/codegen/lang.go \
	command/protoc-gen-parigot/go_/funcchoice.go \
	command/protoc-gen-parigot/go_/gogen.go \
	command/protoc-gen-parigot/go_/gotext.go \
	command/protoc-gen-parigot/codegen/options.go \
	command/protoc-gen-parigot/codegen/pass.go \
	command/protoc-gen-parigot/codegen/wasm.go \
	command/protoc-gen-parigot/util/out.go \
	command/protoc-gen-parigot/util/plugin.go

