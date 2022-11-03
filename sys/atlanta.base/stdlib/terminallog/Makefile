### BEFORE YOU USE THIS MAKEFILE BE SURE YOU HAVE DONE A FULL BUILD OF PARIGOT ABOVE THIS.
### parigot/build MUST BE POPULATED OR THIS WILL DIE A HORRIBLE DEATH.

# this file expects you to have a current directory of stdlib/simplelog (the dir containing this file)

PARIGOT_SOURCE=../../../..

BUILD_TAGS=parigot

GO_CMD=GOOS=js GOARCH=wasm go


all: build/terminallog.p.wasm

build/terminallog.p.wasm: main.go gen.go ../../../../build/runner
	@echo
	@echo "\033[92mbuilding server side: simplelog (server) ============================================================\033[0m"
	$(GO_CMD) build -a -tags $(BUILD_TAGS) -o build/terminallog.p.wasm main.go gen.go

clean:
	rm -rf build/*

#### Do not remove this line or edit below it.  The rest of this file is computed by jdepp.
