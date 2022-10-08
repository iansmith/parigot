TINYGO_MOD_CACHE="/Users/iansmith/tinygo/pkg/mod"

BUILD_PRINT = \e[1;34mBuilding $<\e[0m
GO_CMD=go #really shouldn't need to change this if you use the tools directory
TINYGO_CMD=GOMODCACHE=$(TINYGO_MOD_CACHE) tinygo #really shouldn't need to change this if you use the tools directory
TINYGO_WASM_OPTS=-target wasm -opt 0 -wasm-abi generic
TINYGO_BUILD_TAGS=parigot_abi
TRANSFORM=command/transform
ABI_GO=command/runner/abi.go
PROTOC_PARIGOT_GEN=command/protoc_parigot/proto/gen

all: build/hello-go.p.wasm build/ex1.p.wasm build/runner build/jsstrip build/genabi build/protoc-gen-parigot build/findservices

build/hello-go.wasm: example/hello-go/main.go
	@echo
	@echo "\033[92mexample/hello-go ===================================================================================\033[0m"
	$(TINYGO_CMD) build  $(TINYGO_WASM_OPTS) -tags $(TINYGO_BUILD_TAGS) -o build/hello-go.wasm github.com/iansmith/parigot/example/hello-go

build/ex1.wasm: command/jsstrip/testdata/ex1/ex1.go
	@echo
	@echo "\033[92mcommand/jsstrip/testdata/ex1 =======================================================================\033[0m"
	$(TINYGO_CMD) build  $(TINYGO_WASM_OPTS) -tags $(TINYGO_BUILD_TAGS) -o build/ex1.wasm github.com/iansmith/parigot/command/jsstrip/testdata

build/ex1.p.wasm: build/ex1.wasm
	@echo
	@echo "\033[92mstripping ex1 binary ===============================================================================\033[0m"
	build/jsstrip -o build/ex1.p.wasm build/ex1.wasm

build/hello-go.p.wasm: build/jsstrip build/hello-go.wasm
	@echo
	@echo "\033[92mstripping hello-go binary ==========================================================================\033[0m"
	build/jsstrip -o build/hello-go.p.wasm build/hello-go.wasm

build/runner: sys/cmd/runner/*.go sys/abi_impl/*.go abi/go/abi/*.go sys/cmd/runner/abi.go
	@echo
	@echo "\033[92mrunner =============================================================================================\033[0m"
	$(GO_CMD) build -o build/runner github.com/iansmith/parigot/sys/cmd/runner

sys/cmd/runner/abi.go: build/genabi abi/go/abi/*.go
	@echo
	@echo "\033[92mgenerating ABI wrappers ============================================================================\033[0m"
	build/genabi > sys/cmd/runner/abi.go

build/genabi: sys/cmd/genabi/*.go sys/abi_impl/*.go
	@echo
	@echo "\033[92mABI wrapper generator ==============================================================================\033[0m"
	$(GO_CMD) build -o build/genabi github.com/iansmith/parigot/sys/cmd/genabi

build/findservices: command/findservices/*.go
	@echo
	@echo "\033[92mfind services ======================================================================================\033[0m"
	$(GO_CMD) build -o build/findservices github.com/iansmith/parigot/command/findservices

clean:
	@echo "\033[92mclean ==============================================================================================\033[0m"
	rm -f build/*
	rm -rf $(TINYGO_MOD_CACHE)
	rm -f $(TRANSFORM)/Wasm.* $(TRANSFORM)/WasmLexer.* $(TRANSFORM)/wasm_base_listener.go $(TRANSFORM)/wasm_lexer.go $(TRANSFORM)/wasm_parser.go $(TRANSFORM)/wasm_listener.go
	rm -f $(ABI_GO)
	rm -rf $(PROTOC_PARIGOT_GEN)/*


build/jsstrip: command/Wasm.g4 command/transform/*.go command/jsstrip/*.go command/transform/wasm_parser.go abi/go/abi/*.go
	@echo
	@echo "\033[92mjsstrip ============================================================================================\033[0m"
	go build -o build/jsstrip github.com/iansmith/parigot/command/jsstrip

# only need to run the generator once, not once per file
command/transform/wasm_parser.go: command/Wasm.g4
	@echo
	@echo "\033[92mWASM wat file parser \(via Antlr4 and Wasm.g4\) ======================================================\033[0m"
	pushd command >& /dev/null && java -Xmx500M -cp "../tools/lib/antlr-4.9-complete.jar" org.antlr.v4.Tool -Dlanguage=Go -o transform -package transform Wasm.g4 && popd >& /dev/null

build/protoc-gen-parigot: command/protoc-gen-parigot/*.go command/protoc-gen-parigot/proto/gen/google/protobuf/compiler/plugin.pb.go
	@echo
	@echo "\033[92mprotoc_parigot =====================================================================================\033[0m"
	go build -o build/protoc-gen-parigot github.com/iansmith/parigot/command/protoc-gen-parigot

command/protoc-gen-parigot/proto/gen/google/protobuf/compiler/plugin.pb.go: command/protoc-gen-parigot/proto/google/protobuf/compiler/plugin.proto
	@echo
	@echo "\033[92mgenerating from plugin.proto =======================================================================\033[0m"
	pushd command/protoc-gen-parigot/proto >& /dev/null && buf generate && popd >& /dev/null

#command/protoc-gen-parigot/proto/gen/vvv/vvv.p.go: build/protoc-gen-parigot command/protoc-gen-parigot/proto/vvv/store.proto
#	@echo
#	@echo "\033[92mgenerating test api (vvv) with parigot bindings  ===================================================\033[0m"
#	protoc --go_out=command/protoc-gen-parigot/proto/gen --go_opt=paths=source_relative \
#		--parigot_out=command/protoc-gen-parigot/proto/gen --parigot_opt=paths=source_relative \
#		-I command/protoc-gen-parigot/proto vvv/store.proto
