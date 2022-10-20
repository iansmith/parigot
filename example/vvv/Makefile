### BEFORE YOU USE THIS MAKEFILE BE SURE YOU HAVE DONE A FULL BUILD OF PARIGOT ABOVE THIS.
### parigot/build MUST BE POPULATED OR THIS WILL DIE A HORRIBLE DEATH.

# this file expects you to have a current directory o•f vvv (the dir containing this file)

PARIGOT_SOURCE=../..

# compiler setup.
#TINYGO_MOD_CACHE="$HOME/tinygo/pkg/mod"
#TINYGO_CMD=GOMODCACHE=$(TINYGO_MOD_CACHE) tinygo #really shouldn't need to change this if you use the tools directory
TINYGO_CMD=tinygo

# currently no GC and no scheduler
TINYGO_WASM_OPTS=-target wasm -gc leaking -scheduler none
#TINYGO_WASM_OPTS=-target wasm -gc leaking -scheduler none -opt 0 -x -work
BUILD_TAGS=parigot_abi

GO_CMD=GOOS=js GOARCH=wasm go

## assumes buf is in path.   https://buf.build
BUF_CMD=buf

# anything in a /g directory is generated and should not be checked in. it can deleted at will.
# generated files are place so that correspond exactly to the protobuf "package" statement in
# the .proto file.
GENERATED=g/parigot
all: $(GENERATED)/demo/vvv/servicedecl.go build/storeclient.wasm

build/storeclient.wasm: storeclient/main.go
	@echo
	@echo "\033[92mbuilding client side: vvv ===================================================================================\033[0m"
	$(TINYGO_CMD) build  $(TINYGO_WASM_OPTS) -tags $(BUILD_TAGS) -o build/vvv_client.wasm storeclient/main.go
	#$(GO_CMD) build -a -tags $(BUILD_TAGS) -o build/vvv_client.wasm storeclient/main.go
	wasm2wat --no-check -o build/vvv_client.wat build/vvv_client.wasm

build/vvv.wasm: main.go
	@echo
	@echo "\033[92mbuilding app: vvv ===================================================================================\033[0m"
	$(TINYGO_CMD) build  $(TINYGO_WASM_OPTS) -tags $(TINYGO_BUILD_TAGS) -o build/vvv.wasm main.go

# we only need a single representative file to insure regeneration of all the generate services
$(GENERATED)/demo/vvv/servicedecl.go: proto/business.proto
	$(BUF_CMD) lint
	$(BUF_CMD) generate

clean:
	rm -rf proto/g/parigot/*
	rm -rf build/*

#### Do not remove this line or edit below it.  The rest of this file is computed by jdepp.
### jdepp computed dependencies for binary: example/vvv/build/vvv
example/vvv/build/vvv: \
	command/protoc-gen-parigot/template/go/servicesimpleloc.tmpl \
	build/protoc-gen-parigot \
	example/vvv/proto/business.proto \
	example/vvv/main.go \
	example/vvv/storeclient/main.go \
	command/protoc-gen-parigot/template/go/messagedecl.tmpl \
	command/protoc-gen-parigot/template/go/servicedecl.tmpl

### jdepp computed dependencies for binary: example/vvv/build/storeclient
example/vvv/build/storeclient: \
	example/vvv/storeclient/main.go \
	command/protoc-gen-parigot/template/go/messagedecl.tmpl \
	command/protoc-gen-parigot/template/go/servicedecl.tmpl \
	command/protoc-gen-parigot/template/go/servicesimpleloc.tmpl \
	build/protoc-gen-parigot \
	example/vvv/proto/business.proto

