GO_CMD=go #really shouldn't need to change this if you use the tools directory
FLAVOR=atlanta.base

all: build/runner \
build/protoc-gen-parigot \
build/jdepp \
build/nameserver 

build/jdepp: command/jdepp/main.go
	@echo
	@echo "\033[92mjdepp ============================================================================================\033[0m"
	go build -o build/jdepp github.com/iansmith/parigot/command/jdepp

build/protoc-gen-parigot: \
	command/protoc-gen-parigot/template/go/*.tmpl
	@echo
	@echo "\033[92mprotoc-gen-parigot =================================================================================\033[0m"
	go build -o build/protoc-gen-parigot github.com/iansmith/parigot/command/protoc-gen-parigot

build/runner: sys/atlanta.base/*.go 
	@echo
	@echo "\033[92mrunner =============================================================================================\033[0m"
	go build -o build/runner github.com/iansmith/parigot/command/runner

g/log/logservicedecl.p%go g/net/netservicedecl.p%go: build/protoc-gen-parigot
	@echo
	@echo "\033[92mbuilding parigot interfaces ========================================================================\033[0m"
	buf generate

build/nameserver: command/nameserver/*.go sys/atlanta.base/*.go 
	@echo
	@echo "\033[92mnameserver =============================================================================================\033[0m"
	go build -o build/nameserver github.com/iansmith/parigot/command/nameserver

image:
	@echo
	@echo "\033[92mdocker images =======================================================================================\033[0m"
	# we call this "runner" so we can use the same Dockerfile
	cp build/nameserver container/runner/runner 
	rm -f container/runner/*.p.wasm
	touch data.p.wasm # fake, so we can use same Dockerfile
	# runs in separate shell, so need to cd with && but don't need to "come back"
	cd container/runner && docker build  -t nameserver . 
	rm data.p.wasm
	cp build/runner container/runner/runner
	cp example/vvv/build/server.p.wasm container/runner
	cd container/runner && docker build  -t storeserver . 
	rm container/runner/server.p.wasm
	cp example/vvv/build/storeclient.p.wasm container/runner
	cd container/runner && docker build  -t storeclient . 
	rm container/runner/storeclient.p.wasm

clean:
	@echo "\033[92mclean ==============================================================================================\033[0m"
	rm -f build/*
	rm -rf g/kernel g/log g/net g/*.p.go g/pb
	rm -f $(TRANSFORM)/Wasm.* $(TRANSFORM)/WasmLexer.* $(TRANSFORM)/wasm_base_listener.go $(TRANSFORM)/wasm_lexer.go $(TRANSFORM)/wasm_parser.go $(TRANSFORM)/wasm_listener.go


.PHONY:test
test: $(PGP)
	go run command/protoc-gen-parigot/test/main.go build/protoc-gen-parigot command/protoc-gen-parigot/test/testdata/t0 - abi/abihelper.go

#### Do not remove this line or edit below it.  The rest of this file is computed by jdepp.
### jdepp computed dependencies for binary: build
build: \
	sys/atlanta.base/syscall.go \
	command/runner/fileload.go \
	sys/atlanta.base/memutil.go \
	sys/atlanta.base/process.go \
	sys/atlanta.base/runtime.go \
	sys/atlanta.base/nsproxy.go \
	sys/atlanta.base/func.go \
	sys/atlanta.base/nameserver.go \
	sys/atlanta.base/nameservercore.go \
	sys/atlanta.base/netquic.go \
	command/runner/main.go \
	sys/atlanta.base/netio.go \
	sys/atlanta.base/remote.go \
	sys/atlanta.base/syscallrw.go \
	sys/atlanta.base/local.go

