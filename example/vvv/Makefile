### WARNING: This file's dependency knowlege does NOT include dependencies into parigot proper.
### WARNING: This file assumes that parigot is built and unchanging.  If you make a change in
### WARNING: in parigot, you will need to do `make clean` and `make` in this directory.

GO_CMD=GOOS=js GOARCH=wasm go

PARIGOT_SRC=../..

## assumes buf is in path.   https://buf.build
BUF_CMD=buf

# anything in a /g directory is generated and should not be checked in. it can deleted at will.
# generated files are place so that correspond exactly to the protobuf "package" statement in
# the .proto file.
GENERATED=g
all: $(GENERATED)/store/v1/storeservicedecl.p.go build/storeclient.p.wasm build/server.p.wasm

build/server.p.wasm: main.go 
	@echo
	@echo "\033[92mbuilding server side: vvv (server) ============================================================\033[0m"
	$(GO_CMD) build -a -o build/server.p.wasm main.go

build/storeclient.p.wasm: storeclient/main.go 
	@echo
	@echo "\033[92mbuilding client side: vvv (storeclient) ============================================================\033[0m"
	$(GO_CMD) build -a -o build/storeclient.p.wasm storeclient/main.go

# we only need a single representative file to insure regeneration of all the generated code
$(GENERATED)/store/v1/storeservicedecl.p.go: proto/store/v1/store.proto proto/msg/store/v1/store.proto
	@echo
	@echo "\033[92mgenerating interfaces: vvv =========================================================================\033[0m"
	$(BUF_CMD) lint
	$(BUF_CMD) generate

clean:
	rm -rf g/*
	rm -rf build/*

run: build/server.p.wasm build/storeclient.p.wasm $(GENERATED)/store/v1/storeservicedecl.p.go
	$(PARIGOT_SRC)/build/runner vvv.toml
