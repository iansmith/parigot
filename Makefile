
.PHONY: all example example-hello-go runner clean jsstrip

TINYGO_MOD_CACHE="/Users/iansmith/tinygo/pkg/mod"

GO_CMD=go #really shouldn't need to change this if you use the tools directory
TINYGO_CMD=GOMODCACHE=$(TINYGO_MOD_CACHE) tinygo #really shouldn't need to change this if you use the tools directory

TINYGO_WASM_OPTS=-target wasm -wasm-abi generic

TINYGO_BUILD_TAGS=parigot_abi

all: example runner

example: example-hello-go

jsstrip: build/jsstrip

example-hello-go:
	$(TINYGO_CMD) build  $(TINYGO_WASM_OPTS) -tags $(TINYGO_BUILD_TAGS) -o build/hello-go.wasm github.com/iansmith/parigot/example/hello-go

runner:
	$(GO_CMD) build -o build/runner github.com/iansmith/parigot/sys/init

clean:
	rm -f build/*
	rm -rf $(TINYGO_MOD_CACHE)

build/jsstrip: command/Wasm.g4 command/transform/*.go command/jsstrip/main.go command/transform/wasm_parser.go
	go build -o build/jsstrip github.com/iansmith/parigot/command/jsstrip

# only need to run the generator once, not once per file
command/transform/wasm_parser.go: command/Wasm.g4
	pushd command >& /dev/null && java -Xmx500M -cp "../tools/lib/antlr-4.9-complete.jar" org.antlr.v4.Tool -Dlanguage=Go -o transform -package transform Wasm.g4 && popd >& /dev/null
