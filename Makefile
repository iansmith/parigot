TINYGO_MOD_CACHE="/Users/iansmith/tinygo/pkg/mod"

BUILD_PRINT = \e[1;34mBuilding $<\e[0m
GO_CMD=go #really shouldn't need to change this if you use the tools directory

FLAVOR=atlanta.base



all: build/runner \
$(REP_API_NET) \
$(REP_ABI) \
build/protoc-gen-parigot

TRANSFORM_LIB=command/transform/*.go
ABIPATCH_SRC=command/abipatch/*.go
build/abipatch: $(WASM_GRAMMAR) $(TRANSFORM_LIB) $(REP_GEN_WASM) $(ABIPATCH_SRC) $(REP_ABI)
	@echo
	@echo "\033[92mabipatch ===========================================================================================\033[0m"
	go build -o build/abipatch github.com/iansmith/parigot/command/abipatch

# only need to run the generator once, not once per file
REP_GEN_WASM=command/transform/wasm_parser.go
WASM_GRAMMAR=command/Wasm.g4
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

REP_API_NET=g/parigot/net/netsvc.p.go
API_NET_PROTO=api/$(FLAVOR)/proto/net/net.proto
$(REP_API_NET): $(API_NET_PROTO) $(PGP)
	@echo
	@echo "\033[92mgenerating networking (API) =============================================================================\033[0m"
	buf generate
	gofmt -w api/$(FLAVOR)/go/parigot/net/netsvc.p.go

REP_ABI=g/parigot/abi/abi.pb.go
ABI_PROTO=abi/$(FLAVOR)/proto/abi/abi.proto
$(REP_ABI): $(ABI_PROTO)
	@echo
	@echo "\033[92mgenerating parigot ABI =============================================================================\033[0m"
	buf generate

ABI_GO_HELPER=command/runner/g/abihelper.go
RUNNER_SRC=command/runner/*.go
RUNNER=build/runner
$(ABI_GO_HELPER): abi/$(FLAVOR)/proto/abi/abi.proto $(PGP)
	@echo
	@echo "\033[92mgenerating parigot_abi helper for runner ============================================================\033[0m"
	buf generate
	cp g/parigot/abi/abihelper.go $(ABI_GO_HELPER)

$(RUNNER): $(ABI_GO_HELPER) $(RUNNER_SRC)
	@echo
	@echo "\033[92mrunner ==============================================================================================\033[0m"
	go build -o $(RUNNER) github.com/iansmith/parigot/command/runner

clean:
	@echo "\033[92mclean ==============================================================================================\033[0m"
	rm -f build/*
	rm -rf $(TINYGO_MOD_CACHE)
	rm -f $(TRANSFORM)/Wasm.* $(TRANSFORM)/WasmLexer.* $(TRANSFORM)/wasm_base_listener.go $(TRANSFORM)/wasm_lexer.go $(TRANSFORM)/wasm_parser.go $(TRANSFORM)/wasm_listener.go
	rm -f $(ABI_GO_GEN)
	rm -rf $(API_GEN_DIR)


## shorthands
net: $(REP_API_NET)
abi: $(REP_ABI)
protoc: $(PGP)
gen: $(PGP) abi net
	buf generate
runner:build/runner

