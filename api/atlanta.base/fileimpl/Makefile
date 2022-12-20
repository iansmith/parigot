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
all: $(GENERATED)/file/file.pb.go main.go build/filesvc.p.wasm

# we only need a single representative file to insure regeneration of all the generated code
$(GENERATED)/file/file.pb.go: ../proto/file/file.proto ../proto/pb/file/file.proto
	@echo
	@echo "\033[92mgenerating interfaces: file =========================================================================\033[0m"
	$(BUF_CMD) lint
	$(BUF_CMD) generate

clean:
	rm -rf $(GENERATED)/*
	rm -rf build/*

build/filesvc.p.wasm: main.go $(PARIGOT_SOURCE)/build/protoc-gen-parigot
	@echo
	@echo "\033[92mbuilding server side: filesvc  ============================================================\033[0m"
	$(GO_CMD) build -a -tags $(BUILD_TAGS) -o build/filesvc.p.wasm main.go

#### Do not remove this line or edit below it.  The rest of this file is computed by jdepp.
