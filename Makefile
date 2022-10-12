TINYGO_MOD_CACHE="/Users/iansmith/tinygo/pkg/mod"

BUILD_PRINT = \e[1;34mBuilding $<\e[0m
GO_CMD=go #really shouldn't need to change this if you use the tools directory
TINYGO_CMD=GOMODCACHE=$(TINYGO_MOD_CACHE) tinygo #really shouldn't need to change this if you use the tools directory
TINYGO_WASM_OPTS=-target wasm -opt 0 -wasm-abi generic
TINYGO_BUILD_TAGS=parigot_abi
TRANSFORM=command/transform
ABI_GO=command/runner/abi.go
PROTOC_PARIGOT_GEN=command/protoc_parigot/proto/gen
PGP=build/protoc-gen-parigot
# tuple is (atlanta.base) (version.variant)
FLAVOR=atlanta.base
WASM_GRAMMAR=command/Wasm.g4
TRANSFORM_LIB=command/transform/*.go
STRUCTURE_LIB=command/toml/*.go
REP_GEN_WASM=command/transform/wasm_parser.go
API_GEN_DIR=api/$(FLAVOR)/go/parigot/*

REP_API_NET=api/$(FLAVOR)/go/parigot/net/netsvc.p.go
API_NET_PROTO=api/$(FLAVOR)/proto/net/net.proto

ABI_GEN=abi/$(FLAVOR)/go/abi/*.go
REP_ABI=abi/g/parigot/abi/abi.pb.go
ABI_PROTO=abi/$(FLAVOR)/proto/abi/abi.proto

all: build/runner \
$(REP_API_NET) \
$(REP_ABI) \
build/protoc-gen-parigot

RUNNER_SRC=$(sys/cmd/runner/*.go)
build/runner: $(RUNNER_SRC ) sys/abi_impl/*.go $(REP_ABI) $(PGP)
	@echo
	@echo "\033[92mrunner =============================================================================================\033[0m"
	$(GO_CMD) build -o build/runner github.com/iansmith/parigot/sys/cmd/runner

ABIPATCH_SRC=command/abipatch/*.go
build/abipatch: $(WASM_GRAMMAR) $(TRANSFORM_LIB) $(REP_GEN_WASM) $(ABIPATCH_SRC) $(REP_ABI)
	@echo
	@echo "\033[92mabipatch ===========================================================================================\033[0m"
	go build -o build/abipatch github.com/iansmith/parigot/command/abipatch

# only need to run the generator once, not once per file
$(REP_GEN_WASM): $(WASM_GRAMMAR)
	@echo
	@echo "\033[92mWASM wat file parser \(via Antlr4 and Wasm.g4\) ======================================================\033[0m"
	pushd command >& /dev/null && java -Xmx500M -cp "../tools/lib/antlr-4.9-complete.jar" org.antlr.v4.Tool -Dlanguage=Go -o transform -package transform Wasm.g4 && popd >& /dev/null

PROTOC_GEN_PARIGOT_SRC=command/protoc-gen-parigot/*.go \
command/protoc-gen-parigot/util/*.go \
command/protoc-gen-parigot/*/*.go \
command/protoc-gen-parigot/template/*/*.tmpl

build/protoc-gen-parigot: $(PROTOC_GEN_PARIGOT_SRC) $(STRUCTURE_LIB)
	@echo
	@echo "\033[92mprotoc-gen-parigot =================================================================================\033[0m"
	go build -o build/protoc-gen-parigot github.com/iansmith/parigot/command/protoc-gen-parigot

$(REP_API_NET): $(API_NET_PROTO) $(PGP)
	buf generate
	gofmt -w api/$(FLAVOR)/go/parigot/net/netsvc.p.go

$(REP_ABI): $(ABI_PROTO)
	@echo
	@echo "\033[92mgenerating parigot_abi =============================================================================\033[0m"
	buf generate

#ABIGEN_SRC=command/abigen/*.go command/abigen/template/*.tmpl
#build/abigen: $(ABIGEN_SRC) $(PGP)
#	@echo
#	@echo "\033[92mabigen ===================================================================================\033[0m"
#	go build -o build/abigen github.com/iansmith/parigot/command/abigen

clean:
	@echo "\033[92mclean ==============================================================================================\033[0m"
	rm -f build/*
	rm -rf $(TINYGO_MOD_CACHE)
	rm -f $(TRANSFORM)/Wasm.* $(TRANSFORM)/WasmLexer.* $(TRANSFORM)/wasm_base_listener.go $(TRANSFORM)/wasm_lexer.go $(TRANSFORM)/wasm_parser.go $(TRANSFORM)/wasm_listener.go
	rm -f $(ABI_GO_GEN)
	rm -rf $(API_GEN_DIR)
	rm -rf $(PROTOC_PARIGOT_GEN)/*


## shorthands
net: $(REP_API_NET)
abi: $(REP_ABI)
protoc: $(PGP)
gen: $(PGP) abi net
	buf generate

