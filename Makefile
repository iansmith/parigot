TINYGO_MOD_CACHE="/Users/iansmith/tinygo/pkg/mod"

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
BINPATCH_SRC=command/binpatch/*.go
build/binpatch: $(WASM_GRAMMAR) $(TRANSFORM_LIB) $(REP_GEN_WASM) $(ABIPATCH_SRC) $(REP_ABI)
	@echo
	@echo "\033[92mabipatch ===========================================================================================\033[0m"
	go build -o build/abipatch github.com/iansmith/parigot/command/abipatch

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
$(REP_ABI): $(ABI_PROTO) $(PGP)
	@echo
	@echo "\033[92mgenerating parigot ABI =============================================================================\033[0m"
	buf generate
	gofmt -w $(ABI_GEN_OUT)/*.p.go

ABI_GO_HELPER=command/runner/g/abihelper.p.go
RUNNER_SRC=command/runner/*.go
RUNNER=build/runner
$(ABI_GO_HELPER): abi/$(FLAVOR)/proto/abi/abi.proto $(PGP)
	@echo
	@echo "\033[92mgenerating parigot_abi helper for runner ============================================================\033[0m"
	buf generate
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

.PHONY:test
test: $(PGP)
	go run command/protoc-gen-parigot/test/main.go build/protoc-gen-parigot command/protoc-gen-parigot/test/testdata/t0 - abi/abihelper.go
