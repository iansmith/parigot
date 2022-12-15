### BEFORE YOU USE THIS MAKEFILE BE SURE YOU HAVE DONE A FULL BUILD OF PARIGOT ABOVE THIS.
### parigot/build MUST BE POPULATED OR THIS WILL DIE A HORRIBLE DEATH.

# this file expects you to have a current directory o•f vvv (the dir containing this file)

PARIGOT_SOURCE=../../..
BUILD_TAGS=parigot

GO_CMD=GOOS=js GOARCH=wasm go

## assumes buf is in path.   https://buf.build
BUF_CMD=buf

# anything in a /g directory is generated and should not be checked in. it can deleted at will.
# generated files are place so that correspond exactly to the protobuf "package" statement in
# the .proto file.
GENERATED=../proto/g
all: $(GENERATED)/log/log.pb.go main.go build/logserver.p.wasm

# we only need a single representative file to insure regeneration of all the generated code
$(GENERATED)/log/log.pb.go: ../proto/log/log.proto ../proto/pb/log/log.proto
	@echo
	@echo "\033[92mgenerating interfaces: log =========================================================================\033[0m"
	$(BUF_CMD) lint
	$(BUF_CMD) generate

clean:
	rm -rf $(GENERATED)/log
	rm -rf build/*

build/logserver.p.wasm: main.go $(PARIGOT_SOURCE)/build/protoc-gen-parigot
	@echo
	@echo "\033[92mbuilding server side: logserver  ============================================================\033[0m"
	$(GO_CMD) build -a -tags $(BUILD_TAGS) -o build/logserver.p.wasm main.go

#### Do not remove this line or edit below it.  The rest of this file is computed by jdepp.
### jdepp computed dependencies for binary: api/logimpl/build/logserver.p.wasm
