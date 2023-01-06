API_VERSION=v1

all: g/file/$(API_VERSION)/file.pb.go build/protoc-gen-parigot \
	build/file.p.wasm build/log.p.wasm \
	build/methodcalltest.p.wasm build/methodcallfoo.p.wasm build/methodcallbar.p.wasm

GO_CMD=GOOS=js GOARCH=wasm go

API_PROTO=$(shell find api/proto -type f -regex ".*\.proto")
EXAMPLE_PROTO=$(shell find example -type f -regex ".*\.proto")
TEST_PROTO=$(shell find test -type f -regex ".*\.proto")

## we just use a single representative file for all the generated code
g/file/$(API_VERSION)/file.pb.go: $(API_PROTO) $(EXAMPLE_PROTO) $(TEST_PROTO) build/protoc-gen-parigot 
	rm -rf g/*
	buf lint
	buf generate

.PHONY: protoclean
protoclean: 
	rm -rf g/*

.PHONY: clean
clean: protoclean
	rm -f build/*

TEMPLATE=$(shell find command/protoc-gen-parigot -type f -regex ".*\.tmpl")
GENERATOR_SRC=$(shell find command/protoc-gen-parigot -type f -regex ".*\.go")
build/protoc-gen-parigot: $(TEMPLATE) $(GENERATOR_SRC)
	rm -f $@
	go build -o $@ github.com/iansmith/parigot/command/protoc-gen-parigot

RUNNER_SRC=$(shell find command/runner -type f -regex ".*\.go")
build/runner: $(RUNNER_SRC) g/file/$(API_VERSION)/file.pb.go
	rm -f $@
	go build -o $@ github.com/iansmith/parigot/command/runner

FILE_SERVICE=$(shell find api_impl/file -type f -regex ".*\.go")
build/file.p.wasm: $(FILE_SERVICE) g/file/$(API_VERSION)/file.pb.go
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/api_impl/file

LOG_SERVICE=$(shell find api_impl/log -type f -regex ".*\.go")
build/log.p.wasm: $(LOG_SERVICE) g/file/$(API_VERSION)/file.pb.go
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/api_impl/log

METHODCALLTEST=test/func/methodcall/*.go
build/methodcalltest.p.wasm: $(METHODCALLTEST) g/file/$(API_VERSION)/file.pb.go test/func/methodcall/methodcall.toml
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/test/func/methodcall

FOO_SERVICE=test/func/methodcall/impl/foo/*.go
build/methodcallfoo.p.wasm: $(FOO_SERVICE) g/file/$(API_VERSION)/file.pb.go test/func/methodcall/methodcall.toml
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/test/func/methodcall/impl/foo

BAR_SERVICE=test/func/methodcall/impl/bar/*.go
build/methodcallbar.p.wasm: $(BAR_SERVICE) g/file/$(API_VERSION)/file.pb.go test/func/methodcall/methodcall.toml
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/test/func/methodcall/impl/foo
