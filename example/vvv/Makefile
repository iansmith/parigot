### BEFORE YOU USE THIS MAKEFILE BE SURE YOU HAVE DONE A FULL BUILD OF PARIGOT ABOVE THIS.
### parigot/build MUST BE POPULATED OR THIS WILL DIE A HORRIBLE DEATH.

# this file expects you to have a current directory of vvv (the dir containing this file)

PARIGOT_SOURCE=../..

# compiler setup.  note that the -opt 0 is currently mandatory.
TINYGO_CMD=GOMODCACHE=$(TINYGO_MOD_CACHE) tinygo #really shouldn't need to change this if you use the tools directory
TINYGO_WASM_OPTS=-target wasm -opt 0 -wasm-abi generic
TINYGO_BUILD_TAGS=parigot_abi

## assumes buf is in path.   https://buf.build
BUF_CMD=buf

# anything in a /g directory is generated and should not be checked in. it can deleted at will.
# generated files are place so that correspond exactly to the protobuf "package" statement in
# the .proto file.
GENERATED=g/parigot
all: $(GENERATED)/demo/vvv/servicedecl.go

build/vvv.wasm: main.go
	@echo
	@echo "\033[92building app: vvv ===================================================================================\033[0m"
	$(TINYGO_CMD) build  $(TINYGO_WASM_OPTS) -tags $(TINYGO_BUILD_TAGS) -o build/vvv.wasm main.go

# we only need a single representative file to insure regeneration of all the generate services
$(GENERATED)/demo/vvv/servicedecl.go: proto/business.proto
	$(BUF_CMD) generate

clean:
	rm -rf proto/g/parigot/*
