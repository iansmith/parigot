### BEFORE YOU USE THIS MAKEFILE BE SURE YOU HAVE DONE A FULL BUILD OF PARIGOT ABOVE THIS.
### parigot/build MUST BE POPULATED OR THIS WILL DIE A HORRIBLE DEATH.

# this file expects you to have a current directory o•f vvv (the dir containing this file)

PARIGOT_SOURCE=../..

# compiler setup.
#TINYGO_MOD_CACHE="$HOME/tinygo/pkg/mod"
#TINYGO_CMD=GOMODCACHE=$(TINYGO_MOD_CACHE) tinygo #really shouldn't need to change this if you use the tools directory
#TINYGO_CMD=tinygo

# currently no GC and no scheduler
#TINYGO_WASM_OPTS=-target wasm -gc leaking -scheduler none
#TINYGO_WASM_OPTS=-target wasm -gc leaking -scheduler none -opt 0 -x -work
BUILD_TAGS=parigot

GO_CMD=GOOS=js GOARCH=wasm go

## assumes buf is in path.   https://buf.build
BUF_CMD=buf

# anything in a /g directory is generated and should not be checked in. it can deleted at will.
# generated files are place so that correspond exactly to the protobuf "package" statement in
# the .proto file.
GENERATED=proto/g/
all: $(GENERATED)/vvv/businessservicedecl.p.go build/storeclient.p.wasm build/server.p.wasm

build/server.p.wasm: main.go
	@echo
	@echo "\033[92mbuilding server side: vvv (server) ============================================================\033[0m"
	$(GO_CMD) build -a -tags $(BUILD_TAGS) -o build/server.p.wasm main.go gen.go
	wasm2wat --no-check -o build/server.wat build/server.p.wasm

build/storeclient.p.wasm: storeclient/main.go
	@echo
	@echo "\033[92mbuilding client side: vvv (storeclient) ============================================================\033[0m"
	$(GO_CMD) build -a -tags $(BUILD_TAGS) -o build/storeclient.p.wasm storeclient/main.go
	wasm2wat --no-check -o build/storeclient.wat build/storeclient.p.wasm

build/storeclient-repl.wasm: build/storeclient.wasm
	@echo
	@echo "\033[92mfunction surgery (reflection) ======================================================================\033[0m"
	$(PARIGOT_SOURCE)/build/surgery -op replacefn -r replacement.txt \
		-o build/storeclient-repl.wasm -f '$$runtime.initAll' build/storeclient.wasm

build/vvv.wasm: main.go
	@echo
	@echo "\033[92mbuilding app: vvv ==================================================================================\033[0m"
	$(TINYGO_CMD) build  $(TINYGO_WASM_OPTS) -tags $(TINYGO_BUILD_TAGS) -o build/vvv.wasm main.go

# we only need a single representative file to insure regeneration of all the generate services
$(GENERATED)/vvv/businessservicedecl.p.go: proto/business.proto
	@echo
	@echo "\033[92mgenerating interfaces: vvv =========================================================================\033[0m"
	$(BUF_CMD) lint
	$(BUF_CMD) generate

clean:
	rm -rf proto/g/vvv
	rm -rf build/*

#### Do not remove this line or edit below it.  The rest of this file is computed by jdepp.
### jdepp computed dependencies for binary: example/vvv/build/vvv
example/vvv/build/vvv: \
	example/vvv/proto/g/vvv/pb/businesspb.pb.go \
	build/protoc-gen-parigot \
	example/vvv/main.go \
	example/vvv/storeclient/main.go \
	example/vvv/proto/g/vvv/business.pb.go \
	example/vvv/proto/g/vvv/businessservicedecl.p.go

### jdepp computed dependencies for binary: example/vvv/build/storeclient
example/vvv/build/storeclient: \
	example/vvv/storeclient/main.go \
	example/vvv/proto/g/vvv/business.pb.go \
	example/vvv/proto/g/vvv/businessservicedecl.p.go \
	example/vvv/proto/g/vvv/pb/businesspb.pb.go \
	build/protoc-gen-parigot

