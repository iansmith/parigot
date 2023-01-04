GO_CMD=GOOS=js GOARCH=wasm go

## assumes buf is in path.   https://buf.build
BUF_CMD=buf

# anything in a /g directory is generated and should not be checked in. it can deleted at will.
# generated files are place so that correspond exactly to the protobuf "package" statement in
# the .proto file.
GENERATED=g
all: $(GENERATED)/business/v1/businessservicedecl.p.go build/storeclient.p.wasm build/server.p.wasm

build/server.p.wasm: main.go 
	@echo
	@echo "\033[92mbuilding server side: vvv (server) ============================================================\033[0m"
	$(GO_CMD) build -a -o build/server.p.wasm main.go

build/storeclient.p.wasm: storeclient/main.go 
	@echo
	@echo "\033[92mbuilding client side: vvv (storeclient) ============================================================\033[0m"
	$(GO_CMD) build -a -o build/storeclient.p.wasm storeclient/main.go

# we only need a single representative file to insure regeneration of all the generated code
$(GENERATED)/business/v1/businessservicedecl.p.go: proto/business/v1/business.proto proto/msg/business/v1/business.proto
	@echo
	@echo "\033[92mgenerating interfaces: vvv =========================================================================\033[0m"
	$(BUF_CMD) lint
	$(BUF_CMD) generate

clean:
	rm -rf g/*
	rm -rf build/*

#### Do not remove this line or edit below it.  The rest of this file is computed by jdepp.
### jdepp computed dependencies for binary: build/vv
