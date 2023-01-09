API_VERSION=v1

all: allprotos \
	commands \
	apiimpl \
	methodcalltest

allprotos: g/file/$(API_VERSION)/file.pb.go 
methodcalltest: build/methodcalltest.p.wasm build/methodcallfoo.p.wasm build/methodcallbar.p.wasm
apiimpl: build/file.p.wasm build/log.p.wasm
commands: 	build/protoc-gen-parigot build/runner 


GO_CMD=GOOS=js GOARCH=wasm go

API_PROTO=$(shell find api/proto -type f -regex ".*\.proto")
EXAMPLE_PROTO=$(shell find example -type f -regex ".*\.proto")
TEST_PROTO=$(shell find test -type f -regex ".*\.proto")

SPLIT_UTIL=$(shell find api_impl/splitutil -type f -regex ".*\.go")

## we just use a single representative file for all the generated code
g/file/$(API_VERSION)/file.pb.go: $(API_PROTO) $(EXAMPLE_PROTO) $(TEST_PROTO) build/protoc-gen-parigot 
	rm -rf g/*
	buf lint
	buf generate

# protoc plugin
TEMPLATE=$(shell find command/protoc-gen-parigot -type f -regex ".*\.tmpl")
GENERATOR_SRC=$(shell find command/protoc-gen-parigot -type f -regex ".*\.go")
build/protoc-gen-parigot: $(TEMPLATE) $(GENERATOR_SRC)
	rm -f $@
	go build -o $@ github.com/iansmith/parigot/command/protoc-gen-parigot


# launch a deployment based on a config file with runner
RUNNER_SRC=$(shell find command/runner -type f -regex ".*\.go")
SYS_SRC=$(shell find sys -type f -regex ".*\.go")
build/runner: $(RUNNER_SRC) g/file/$(API_VERSION)/file.pb.go $(SYS_SRC) $(SPLIT_UTIL)
	rm -f $@
	go build -o $@ github.com/iansmith/parigot/command/runner

# implementation of the file service
FILE_SERVICE=$(shell find api_impl/file -type f -regex ".*\.go")
build/file.p.wasm: $(FILE_SERVICE) g/file/$(API_VERSION)/file.pb.go $(SPLIT_UTIL)
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/api_impl/file

# implementation of the log service
LOG_SERVICE=$(shell find api_impl/log -type f -regex ".*\.go")
build/log.p.wasm: $(LOG_SERVICE) g/file/$(API_VERSION)/file.pb.go $(SPLIT_UTIL)
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/api_impl/log

# methodcall test code
METHODCALLTEST=test/func/methodcall/*.go
build/methodcalltest.p.wasm: $(METHODCALLTEST) g/file/$(API_VERSION)/file.pb.go build/methodcallbar.p.wasm build/methodcallfoo.p.wasm build/runner
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/test/func/methodcall

# methodcall service impl: methodcall.FooService
FOO_SERVICE=test/func/methodcall/impl/foo/*.go
build/methodcallfoo.p.wasm: $(FOO_SERVICE) g/file/$(API_VERSION)/file.pb.go test/func/methodcall/methodcall.toml
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/test/func/methodcall/impl/foo

# methodcall service impl: methodcall.BarService
BAR_SERVICE=test/func/methodcall/impl/bar/*.go
build/methodcallbar.p.wasm: $(BAR_SERVICE) g/file/$(API_VERSION)/file.pb.go test/func/methodcall/methodcall.toml
	rm -f $@
	$(GO_CMD) build -a -o $@ github.com/iansmith/parigot/test/func/methodcall/impl/bar

#
# TEST
#
.PHONY: test
test: methodcalltest test/func/methodcall/methodcall.toml all
	build/runner -t test/func/methodcall/methodcall.toml 
#
# CLEAN
#
.PHONY: protoclean
protoclean: 
	rm -rf g/*

.PHONY: clean
clean: protoclean
	rm -f build/*

