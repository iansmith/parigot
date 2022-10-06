TINYGO_MOD_CACHE="/Users/iansmith/tinygo/pkg/mod"

GO_CMD=go #really shouldn't need to change this if you use the tools directory
TINYGO_CMD=GOMODCACHE=$(TINYGO_MOD_CACHE) tinygo #really shouldn't need to change this if you use the tools directory

TINYGO_WASM_OPTS=-target wasm -wasm-abi generic

TINYGO_BUILD_TAGS=parigot_abi

TRANSFORM=command/transform

ABI_GO=command/runner/abi.go

all: build/hello-go.p.wasm build/runner build/jsstrip build/genabi

example: example-hello-go

build/hello-go.wasm: example/hello-go/main.go
	$(TINYGO_CMD) build  $(TINYGO_WASM_OPTS) -tags $(TINYGO_BUILD_TAGS) -o build/hello-go.wasm github.com/iansmith/parigot/example/hello-go

build/hello-go.p.wasm: build/jsstrip build/hello-go.wasm
	build/jsstrip -o build/hello-go.p.wasm build/hello-go.wasm

build/runner: sys/cmd/runner/*.go sys/abi_impl/*.go abi/go/abi/*.go sys/cmd/runner/abi.go
	$(GO_CMD) build -o build/runner github.com/iansmith/parigot/sys/cmd/runner

sys/cmd/runner/abi.go: build/genabi abi/go/abi/*.go
	build/genabi > sys/cmd/runner/abi.go

build/genabi: sys/cmd/genabi/*.go sys/abi_impl/*.go
	$(GO_CMD) build -o build/genabi github.com/iansmith/parigot/sys/cmd/genabi

clean:
	rm -f build/*
	rm -rf $(TINYGO_MOD_CACHE)
	rm -f $(TRANSFORM)/Wasm.* $(TRANSFORM)/WasmLexer.* $(TRANSFORM)/wasm_base_listener.go $(TRANSFORM)/wasm_lexer.go $(TRANSFORM)/wasm_parser.go $(TRANSFORM)/wasm_listener.go
	rm -f $(ABI_GO)


build/jsstrip: command/Wasm.g4 command/transform/*.go command/jsstrip/*.go command/transform/wasm_parser.go abi/go/abi/*.go
	go build -o build/jsstrip github.com/iansmith/parigot/command/jsstrip

# only need to run the generator once, not once per file
command/transform/wasm_parser.go: command/Wasm.g4
	pushd command >& /dev/null && java -Xmx500M -cp "../tools/lib/antlr-4.9-complete.jar" org.antlr.v4.Tool -Dlanguage=Go -o transform -package transform Wasm.g4 && popd >& /dev/null
