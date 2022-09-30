
.PHONY: all example example-hello-go runner clean jsstrip

TINYGO_MOD_CACHE="/Users/iansmith/tinygo/pkg/mod"

GO_CMD=go #really shouldn't need to change this if you use the tools directory
TINYGO_CMD=GOMODCACHE=$(TINYGO_MOD_CACHE) tinygo #really shouldn't need to change this if you use the tools directory

TINYGO_WASM_OPTS=-target wasm -wasm-abi generic

TINYGO_BUILD_TAGS=parigot_abi

all: jsstrip example runner

example: example-hello-go

example-hello-go:
	$(TINYGO_CMD) build  $(TINYGO_WASM_OPTS) -tags $(TINYGO_BUILD_TAGS) -o build/hello-go.wasm github.com/iansmith/parigot/example/hello-go

runner:
	$(GO_CMD) build -o build/runner github.com/iansmith/parigot/sys/init

clean:
	rm -f build/*
	rm -rf $(TINYGO_MOD_CACHE)

jsstrip:
	java -Xmx500M -cp "tools/lib/antlr-4.9-complete.jar" org.antlr.v4.Tool -Xlog -Dlanguage=Go -package main command/jsstrip/jsstrip.g4
