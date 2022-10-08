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
TOOLS=build/jsstrip build/findservices $(PGP) build/runner
# tuple is (atlanta,base,go) (version, variant, lang)
FLAVOR=atlanta1/base/go
WASM_GRAMMAR=command/Wasm.g4
TRANSFORM_LIB=command/transform/*.go
REP_GEN_WASM=command/transform/wasm_parser.go

REP_API_NET=api/$(FLAVOR)/parigot/net/proto/gen/net/servicedecl.go
API_NET_PROTO=api/$(FLAVOR)/parigot/net/proto/net/net.proto

ABI_SOURCE=abi/$(FLAVOR)/abi/*.go
REP_ABI=abi/$(FLAVOR)/parigot/abi/proto/gen/abi/servicedecl.go
ABI_PROTO=abi/$(FLAVOR)/parigot/abi/proto/abi/abi.proto

APP_CODE=build/hello-go.p.wasm \
build/ex1.p.wasm

all: build/runner \
$(TOOLS) \
$(REP_API_NET) \
build/abigen

#build/hello-go.wasm: example/hello-go/main.go
#	@echo
#	@echo "\033[92mexample/hello-go ===================================================================================\033[0m"
#	$(TINYGO_CMD) build  $(TINYGO_WASM_OPTS) -tags $(TINYGO_BUILD_TAGS) -o build/hello-go.wasm github.com/iansmith/parigot/example/hello-go

#build/ex1.wasm: command/jsstrip/testdata/ex1/ex1.go
#	@echo
#	@echo "\033[92mcommand/jsstrip/testdata/ex1 =======================================================================\033[0m"
#	$(TINYGO_CMD) build  $(TINYGO_WASM_OPTS) -tags $(TINYGO_BUILD_TAGS) -o build/ex1.wasm github.com/iansmith/parigot/command/jsstrip/testdata

#build/ex1.p.wasm: build/ex1.wasm
#	@echo
#	@echo "\033[92mstripping ex1 binary ===============================================================================\033[0m"
#	build/jsstrip -o build/ex1.p.wasm build/ex1.wasm

#build/hello-go.p.wasm: build/jsstrip build/hello-go.wasm
#	@echo
#	@echo "\033[92mstripping hello-go binary ==========================================================================\033[0m"
#	build/jsstrip -o build/hello-go.p.wasm build/hello-go.wasm

RUNNER_SRC=$(sys/cmd/runner/*.go sys/cmd/runner/abi.go)
build/runner: $(RUNNER_SRC ) sys/abi_impl/*.go $(REP_ABI)  $(PGP)
	@echo
	@echo "\033[92mrunner =============================================================================================\033[0m"
	$(GO_CMD) build -o build/runner github.com/iansmith/parigot/sys/cmd/runner

#sys/cmd/runner/abi.go: build/genabi abi/atlanta1/base/go/abi/*.go
#	@echo
#	@echo "\033[92mgenerating ABI wrappers ============================================================================\033[0m"
#	build/genabi > sys/cmd/runner/abi.go

#build/genabi: sys/cmd/genabi/*.go sys/abi_impl/*.go
#	@echo
#	@echo "\033[92mABI wrapper generator ==============================================================================\033[0m"
#	$(GO_CMD) build -o build/genabi github.com/iansmith/parigot/sys/cmd/genabi

FIND_SERVICES_SRC=command/findservices/*.go command/findservices/template/*.tmpl
build/findservices: $(FIND_SERVICES_SRC) $(PGP)
	@echo
	@echo "\033[92mfind services ======================================================================================\033[0m"
	$(GO_CMD) build -o build/findservices github.com/iansmith/parigot/command/findservices

JSSTRIP_SRC=command/jsstrip/*.go
build/jsstrip: $(WASM_GRAMMAR) $(TRANSFORM_LIB)  $(REP_GEN_WASM) $(JSSTRIP_SRC) $(ABI_SRC)
	@echo
	@echo "\033[92mjsstrip ============================================================================================\033[0m"
	go build -o build/jsstrip github.com/iansmith/parigot/command/jsstrip

# only need to run the generator once, not once per file
$(REP_GEN_WASM): $(WASM_GRAMMAR)
	@echo
	@echo "\033[92mWASM wat file parser \(via Antlr4 and Wasm.g4\) ======================================================\033[0m"
	pushd command >& /dev/null && java -Xmx500M -cp "../tools/lib/antlr-4.9-complete.jar" org.antlr.v4.Tool -Dlanguage=Go -o transform -package transform Wasm.g4 && popd >& /dev/null

PLUGIN_SRC=command/protoc-gen-parigot/proto/gen/google/protobuf/compiler/plugin.pb.go
PROTOC_GEN_PARIGOT_SRC=command/protoc-gen-parigot/*.go \
$(PLUGIN_SRC)\
command/protoc-gen-parigot/template/*.tmpl

build/protoc-gen-parigot: $(PROTOC_GEN_PARIGOT_SRC)
	@echo
	@echo "\033[92mprotoc_parigot =====================================================================================\033[0m"
	go build -o build/protoc-gen-parigot github.com/iansmith/parigot/command/protoc-gen-parigot

$(PLUGIN_SRC): command/protoc-gen-parigot/proto/google/protobuf/compiler/plugin.proto
	@echo
	@echo "\033[92mgenerating from plugin.proto =======================================================================\033[0m"
	pushd command/protoc-gen-parigot/proto >& /dev/null && buf generate && popd >& /dev/null

#command/protoc-gen-parigot/proto/gen/vvv/vvv.p.go: build/protoc-gen-parigot command/protoc-gen-parigot/proto/vvv/store.proto
#	@echo
#	@echo "\033[92mgenerating test api (vvv) with parigot bindings  ===================================================\033[0m"
#	protoc --go_out=command/protoc-gen-parigot/proto/gen --go_opt=paths=source_relative \
#		--parigot_out=command/protoc-gen-parigot/proto/gen --parigot_opt=paths=source_relative \
#		-I command/protoc-gen-parigot/proto vvv/store.proto

$(REP_API_NET): $(API_NET_PROTO) $(TOOLS)
	@echo
	@echo "\033[92mgenerating parigot_api =============================================================================\033[0m"
	pushd api/$(FLAVOR)/parigot/net/proto >& /dev/null && buf generate && popd >& /dev/null
	build/findservices api

$(REP_ABI): $(ABI_PROTO) $(TOOLS)
	@echo
	@echo "\033[92mgenerating parigot_api =============================================================================\033[0m"
	pushd abi/$(FLAVOR)/parigot/abi/proto >& /dev/null && buf generate && popd >& /dev/null
	build/findservices abi

ABIGEN_SRC=command/abigen/*.go command/abigen/template/*.tmpl
build/abigen: $(ABIGEN_SRC) $(PGP)
	@echo
	@echo "\033[92mabigen ===================================================================================\033[0m"
	go build -o build/abigen github.com/iansmith/parigot/command/abigen

clean:
	@echo "\033[92mclean ==============================================================================================\033[0m"
	rm -f build/*
	rm -rf $(TINYGO_MOD_CACHE)
	rm -f $(TRANSFORM)/Wasm.* $(TRANSFORM)/WasmLexer.* $(TRANSFORM)/wasm_base_listener.go $(TRANSFORM)/wasm_lexer.go $(TRANSFORM)/wasm_parser.go $(TRANSFORM)/wasm_listener.go
	rm -f $(ABI_GO)
	rm -rf $(PROTOC_PARIGOT_GEN)/*

net: $(REP_API_NET)
abi: $(REP_ABI)
