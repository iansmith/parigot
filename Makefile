BUILD_PRINT = \e[1;34mBuilding $<\e[0m
GO_CMD=go #really shouldn't need to change this if you use the tools directory

FLAVOR=atlanta.base

all: build/runner \
$(REP_API_NET) \
$(REP_ABI) \
build/protoc-gen-parigot

# transform library
TRANSFORM_LIB=command/transform/*.go

# only need to run the generator once, not once per file
REP_GEN_WASM=command/transform/wasm_parser.go
WASM_GRAMMAR=command/Wasm.g4
SURGERY_SRC=command/surgery/*.go
build/surgery: $(WASM_GRAMMAR) $(TRANSFORM_LIB) $(REP_GEN_WASM) $(ABIPATCH_SRC) $(REP_ABI) $(SURGERY_SRC)
	@echo
	@echo "\033[92msurgery =============================================================================================\033[0m"
	go build -o build/surgery github.com/iansmith/parigot/command/surgery

$(REP_GEN_WASM): $(WASM_GRAMMAR)
	@echo
	@echo "\033[92mWASM wat file parser \(via Antlr4 and Wasm.g4\) ======================================================\033[0m"
	pushd command >& /dev/null && java -Xmx500M -cp "../tools/lib/antlr-4.9-complete.jar" org.antlr.v4.Tool -Dlanguage=Go -o transform -package transform Wasm.g4 && popd >& /dev/null

PROTOC_GEN_PARIGOT_SRC=command/protoc-gen-parigot/*.go \
command/protoc-gen-parigot/util/*.go \
command/protoc-gen-parigot/codegen/*.go \
command/protoc-gen-parigot/go_/*.go \
command/protoc-gen-parigot/abi/*.go \
command/protoc-gen-parigot/template/abi/*.tmpl \
command/protoc-gen-parigot/template/go/*.tmpl
PGP=build/protoc-gen-parigot
$(PGP): $(PROTOC_GEN_PARIGOT_SRC)
	@echo
	@echo "\033[92mprotoc-gen-parigot =================================================================================\033[0m"
	go build -o build/protoc-gen-parigot github.com/iansmith/parigot/command/protoc-gen-parigot

API_NET_PROTO=api/$(FLAVOR)/proto/net/net.proto
API_NET_GEN_OUT=g/parigot/net
REP_API_NET=$(API_NET_GEN_OUT)/netservicedecl.p.go
$(REP_API_NET): $(API_NET_PROTO) $(PGP)
	@echo
	@echo "\033[92mgenerating networking (API) =============================================================================\033[0m"
	buf generate
	gofmt -w $(REP_API_NET) $(API_NET_GEN_OUT)/netmessagedecl.p.go

ABI_GEN_OUT=g/parigot/abi
REP_ABI=$(ABI_GEN_OUT)/abi.pb.go
ABI_PROTO=abi/$(FLAVOR)/proto/abi/abi.proto
ORIG_UNDEF=$(TINYGOROOT)/targets/wasm-undefined.txt.orig
TINYGO_UNDEF=$(TINYGOROOT)/targets/wasm-undefined.txt
ABI_UNDEF=$(ABI_GEN_OUT)/abiwasm-undefined.txt

$(REP_ABI): $(ABI_PROTO) $(PGP)
	@echo
	@echo "\033[92mgenerating parigot ABI =============================================================================\033[0m"
	buf generate
	gofmt -w $(ABI_GEN_OUT)/*.p.go
	cp $(ORIG_UNDEF) $(TINYGO_UNDEF)
	cat $(ABI_UNDEF) >> $(TINYGO_UNDEF)

ABI_GO_HELPER=command/runner/g/abihelper.p.go
RUNNER_SRC=command/runner/*.go
RUNNER=build/runner
$(ABI_GO_HELPER): abi/$(FLAVOR)/proto/abi/abi.proto $(PGP) \
	abi/atlanta.base/go/jspatch/*.go abi/atlanta.base/go/tinygopatch/*.go \
	$(REP_ABI)
	@echo
	@echo "\033[92mgenerating parigot_abi helper for runner ============================================================\033[0m"
	buf generate
	gofmt -w $(ABI_GEN_OUT)/*.go
	mv $(ABI_GEN_OUT)/abihelper.p.go $(ABI_GO_HELPER)
	gofmt -w $(ABI_GO_HELPER)

$(RUNNER): $(ABI_GO_HELPER) $(RUNNER_SRC) $(PGP)
	@echo
	@echo "\033[92mrunner ==============================================================================================\033[0m"
	go build -o $(RUNNER) github.com/iansmith/parigot/command/runner

clean:
	@echo "\033[92mclean ==============================================================================================\033[0m"
	rm -f build/*
	rm -rf g/parigot/*
	rm -rf command/runner/g/*
	rm -rf $(TINYGO_MOD_CACHE)
	rm -f $(TRANSFORM)/Wasm.* $(TRANSFORM)/WasmLexer.* $(TRANSFORM)/wasm_base_listener.go $(TRANSFORM)/wasm_lexer.go $(TRANSFORM)/wasm_parser.go $(TRANSFORM)/wasm_listener.go


## shorthands
net: $(REP_API_NET)
abi: $(REP_ABI)
protoc: $(PGP)
gen: $(PGP) abi net
	buf generate
runner:build/runner
surgery:build/surgery

.PHONY:test
test: $(PGP)
	go run command/protoc-gen-parigot/test/main.go build/protoc-gen-parigot command/protoc-gen-parigot/test/testdata/t0 - abi/abihelper.go

#### Do not remove this line or edit below it.  The rest of this file is computed by jdepp.
### jdepp computed dependencies for binary: build/testdata
build/testdata: \
	command/surgery/testdata/ex1/ex1.go \
	lib/base/go/log/log.go \
	command/surgery/testdata/main.go

### jdepp computed dependencies for binary: build/hello-go
build/hello-go: \
	example/hello-go/main.go \
	lib/base/go/log/log.go \
	command/protoc-gen-parigot/template/go/messagedecl.tmpl \
	command/protoc-gen-parigot/template/go/servicedecl.tmpl \
	command/protoc-gen-parigot/template/go/servicesimpleloc.tmpl \
	build/protoc-gen-parigot

### jdepp computed dependencies for binary: build/genabi
build/genabi: \
	command/transform/sem_op.go \
	command/transform/typedescriptor.go \
	command/transform/wasm_lexer.go \
	command/transform/wasm_base_listener.go \
	command/transform/wasm_parser.go \
	command/transform/build_misc.go \
	command/transform/build_stmt.go \
	command/transform/sem_module.go \
	command/transform/sem_toplevel.go \
	sys/cmd/genabi/main.go \
	command/transform/sem_misc.go \
	command/transform/sem_stmt.go \
	command/transform/wasm_listener.go \
	command/transform/build_module.go \
	command/transform/build_terminal.go \
	command/transform/sem_func.go

### jdepp computed dependencies for binary: build/jdepp
build/jdepp: \
	command/jdepp/main.go

### jdepp computed dependencies for binary: build/protoc-gen-parigot
build/protoc-gen-parigot: \
	command/protoc-gen-parigot/test/main.go \
	command/protoc-gen-parigot/codegen/geninfo.go \
	command/protoc-gen-parigot/codegen/helper.go \
	command/protoc-gen-parigot/codegen/lang.go \
	command/protoc-gen-parigot/codegen/wasmfield.go \
	command/protoc-gen-parigot/go_/gotext.go \
	command/protoc-gen-parigot/util/plugin.go \
	command/protoc-gen-parigot/codegen/cgtype.go \
	command/protoc-gen-parigot/codegen/funcchooser.go \
	command/protoc-gen-parigot/codegen/generate.go \
	command/protoc-gen-parigot/codegen/wasm.go \
	command/protoc-gen-parigot/codegen/wasmMethod.go \
	command/protoc-gen-parigot/util/out.go \
	command/protoc-gen-parigot/abi/abigen.go \
	command/protoc-gen-parigot/codegen/finder.go \
	command/protoc-gen-parigot/codegen/pass.go \
	command/protoc-gen-parigot/codegen/text.go \
	command/protoc-gen-parigot/codegen/wasmmessage.go \
	command/protoc-gen-parigot/go_/funcchoice.go \
	command/protoc-gen-parigot/codegen/options.go \
	command/protoc-gen-parigot/codegen/wasmservice.go \
	command/protoc-gen-parigot/go_/gogen.go \
	command/protoc-gen-parigot/main.go

### jdepp computed dependencies for binary: build/test
build/test: \
	command/protoc-gen-parigot/test/main.go

### jdepp computed dependencies for binary: build/runner
build/runner: \
	command/runner/main.go \
	sys/abiimpl/abiimpl.go \
	abi/atlanta.base/go/jspatch/jshelper.go \
	abi/atlanta.base/go/tinygopatch/tinygohelper.go

### jdepp computed dependencies for binary: build/surgery
build/surgery: \
	command/surgery/unlink.go \
	command/transform/build_misc.go \
	command/transform/sem_stmt.go \
	command/transform/wasm_listener.go \
	command/surgery/main.go \
	command/transform/build_module.go \
	command/transform/build_stmt.go \
	command/transform/wasm_base_listener.go \
	command/transform/wasm_lexer.go \
	command/transform/sem_func.go \
	command/transform/typedescriptor.go \
	command/surgery/testdata/ex1/ex1.go \
	command/surgery/testdata/main.go \
	command/surgery/convert.go \
	command/surgery/parse.go \
	command/surgery/tree.go \
	command/transform/build_terminal.go \
	command/transform/sem_misc.go \
	command/transform/sem_module.go \
	command/transform/sem_op.go \
	command/transform/sem_toplevel.go \
	command/transform/wasm_parser.go \
	lib/base/go/log/log.go

