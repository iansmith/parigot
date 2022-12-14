### BEFORE YOU USE THIS MAKEFILE BE SURE YOU HAVE DONE A FULL BUILD OF PARIGOT ABOVE THIS.
### parigot/build MUST BE POPULATED OR THIS WILL DIE A HORRIBLE DEATH.

# this file expects you to have a current directory o•f vvv (the dir containing this file)

PARIGOT_SOURCE=../..
BUILD_TAGS=parigot

GO_CMD=GOOS=js GOARCH=wasm go

## assumes buf is in path.   https://buf.build
BUF_CMD=buf

# anything in a /g directory is generated and should not be checked in. it can deleted at will.
# generated files are place so that correspond exactly to the protobuf "package" statement in
# the .proto file.
GENERATED=proto/g/
all: $(GENERATED)/vvv/businessservicedecl.p.go build/storeclient.p.wasm build/server.p.wasm

build/server.p.wasm: main.go $(PARIGOT_SOURCE)/build/protoc-gen-parigot
	@echo
	@echo "\033[92mbuilding server side: vvv (server) ============================================================\033[0m"
	$(GO_CMD) build -a -tags $(BUILD_TAGS) -o build/server.p.wasm main.go
	#wasm2wat --no-check -o build/server.wat build/server.p.wasm

build/storeclient.p.wasm: storeclient/main.go $(PARIGOT_SOURCE)/build/protoc-gen-parigot
	@echo
	@echo "\033[92mbuilding client side: vvv (storeclient) ============================================================\033[0m"
	$(GO_CMD) build -a -tags $(BUILD_TAGS) -o build/storeclient.p.wasm storeclient/main.go
	#wasm2wat --no-check -o build/storeclient.wat build/storeclient.p.wasm

# we only need a single representative file to insure regeneration of all the generated code
$(GENERATED)/vvv/businessservicedecl.p.go: proto/business.proto
	@echo
	@echo "\033[92mgenerating interfaces: vvv =========================================================================\033[0m"
	$(BUF_CMD) lint
	$(BUF_CMD) generate

clean:
	rm -rf proto/g/vvv
	rm -rf build/*

#### Do not remove this line or edit below it.  The rest of this file is computed by jdepp.
### jdepp computed dependencies for binary: build/vv
build/vv: \
	/workspaces/parigot/example/vvv/main.go \
	proto/g/vvv/business.pb.go \
	proto/g/vvv/businessserver.p.go \
	proto/g/vvv/businessservicedecl.p.go \
	proto/g/vvv/pb/businesspb.pb.go \
	/workspaces/parigot/example/vvv/storeclient/main.go

### jdepp computed dependencies for binary: build/storeclient
build/storeclient: \
	proto/g/vvv/pb/businesspb.pb.go \
	storeclient/main.go \
	proto/g/vvv/business.pb.go \
	proto/g/vvv/businessserver.p.go \
	proto/g/vvv/businessservicedecl.p.go

